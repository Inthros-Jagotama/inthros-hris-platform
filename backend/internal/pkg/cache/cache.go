// Package cache menyediakan distributed cache dengan Redis Pub/Sub
// untuk invalidasi cache antar instance Go server.
//
// Arsitektur:
//
//	┌─────────────────────┐     Pub/Sub           ┌─────────────────────┐
//	│  Instance A         │◄══════════════════════►│  Instance B         │
//	│  ┌───────────────┐  │  "cache:invalidate"   │  ┌───────────────┐  │
//	│  │ Local LRU     │  │                       │  │ Local LRU     │  │
//	│  │ Cache (fast)  │  │                       │  │ Cache (fast)  │  │
//	│  └───────┬───────┘  │                       │  └───────┬───────┘  │
//	│          │          │                       │          │          │
//	│  ┌───────▼───────┐  │                       │  ┌───────▼───────┐  │
//	│  │ Redis Store   │  │                       │  │ Redis Store   │  │
//	│  │ (shared)      │  │                       │  │ (shared)      │  │
//	│  └───────┬───────┘  │                       │  └───────┬───────┘  │
//	└──────────┼──────────┘                       └──────────┼──────────┘
//	           ▼                                               ▼
//	   ┌────────────────┐                             ┌────────────────┐
//	   │    Redis        │                             │    Redis       │
//	   │ (Cache + PubSub)│                             │ (Cache +PubSub)│
//	   └────────────────┘                             └────────────────┘
//
// Cache levels:
//  1. Local (sync.Map + TTL) — microsecond latency, per-instance
//  2. Redis (go-redis) — millisecond latency, shared across instances
//  3. Database — fallback jika cache miss
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// KeyPrefix untuk namespacing cache keys.
const (
	PrefixPlatform = "hris:platform:" // company, module, license, user
	PrefixTenant   = "hris:tenant:"   // organization, employee, dll
	PrefixCache    = "hris:cache:"    // generic cache
)

// item adalah entry di local cache dengan TTL.
type item struct {
	value     interface{}
	expiresAt time.Time
}

// Cache menyediakan two-tier cache (local + Redis) dengan
// Pub/Sub-based invalidation untuk distributed environment.
type Cache struct {
	client     *redis.Client
	local      sync.Map
	logger     *zap.Logger
	defaultTTL time.Duration
	pubSub     *PubSub
}

// Config untuk cache.
type Config struct {
	RedisAddr     string        // host:port
	RedisPassword string
	RedisDB       int
	DefaultTTL    time.Duration // default: 5 menit
}

// New membuat Cache baru dan memulai Pub/Sub listener.
// Jika Redis tidak tersedia, secara otomatis fallback ke local-only cache
// (tanpa Redis / Pub/Sub) untuk development environment.
func New(cfg Config, logger *zap.Logger) (*Cache, error) {
	if cfg.DefaultTTL == 0 {
		cfg.DefaultTTL = 5 * time.Minute
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Test koneksi — jika gagal, fallback ke local-only cache
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Warn("Redis unavailable, falling back to local-only cache",
			zap.String("redis_addr", cfg.RedisAddr),
			zap.Error(err),
		)
		// Local-only cache (no Redis, no Pub/Sub)
		c := &Cache{
			client:     nil,
			local:      sync.Map{},
			logger:     logger,
			defaultTTL: cfg.DefaultTTL,
		}
		logger.Info("Local-only cache initialized (no Redis)",
			zap.Duration("default_ttl", cfg.DefaultTTL),
		)
		return c, nil
	}

	c := &Cache{
		client:     rdb,
		logger:     logger,
		defaultTTL: cfg.DefaultTTL,
	}

	// Inisialisasi Pub/Sub untuk distributed cache invalidation
	ps, err := NewPubSub(rdb, logger, c.evictLocal)
	if err != nil {
		logger.Warn("Pub/Sub initialization failed, falling back to local-only cache",
			zap.Error(err),
		)
		// Close Redis client since we'll use local-only
		rdb.Close()
		c := &Cache{
			client:     nil,
			local:      sync.Map{},
			logger:     logger,
			defaultTTL: cfg.DefaultTTL,
		}
		logger.Info("Local-only cache initialized (Pub/Sub unavailable)",
			zap.Duration("default_ttl", cfg.DefaultTTL),
		)
		return c, nil
	}
	c.pubSub = ps

	logger.Info("Cache initialized with Redis",
		zap.String("redis_addr", cfg.RedisAddr),
		zap.Duration("default_ttl", cfg.DefaultTTL),
	)

	return c, nil
}

// =========================================================================
// Public API
// =========================================================================

// Get mengambil data dari cache. Cari di local cache dulu, baru Redis.
func (c *Cache) Get(ctx context.Context, key string) ([]byte, bool) {
	// 1. Cek local cache
	if val, ok := c.getLocal(key); ok {
		return val, true
	}

	// 2. Cek Redis (jika tersedia)
	if c.client != nil {
		val, err := c.client.Get(ctx, key).Bytes()
		if err == nil {
			// Simpan ke local cache untuk akses berikutnya
			c.setLocal(key, val, c.defaultTTL)
			return val, true
		}
	}

	return nil, false
}

// Set menyimpan data ke Redis dan local cache.
func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if ttl == 0 {
		ttl = c.defaultTTL
	}

	if c.client != nil {
		if err := c.client.Set(ctx, key, value, ttl).Err(); err != nil {
			return fmt.Errorf("cache: set failed for %s: %w", key, err)
		}
	}

	c.setLocal(key, value, ttl)
	return nil
}

// SetJSON menyimpan data dalam format JSON.
func (c *Cache) SetJSON(ctx context.Context, key string, data interface{}, ttl time.Duration) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cache: json marshal failed for %s: %w", key, err)
	}
	return c.Set(ctx, key, bytes, ttl)
}

// Invalidate menghapus key dari semua cache (lokal + Redis + instance lain via Pub/Sub).
func (c *Cache) Invalidate(ctx context.Context, key string) error {
	// 1. Hapus dari Redis (jika tersedia)
	if c.client != nil {
		if err := c.client.Del(ctx, key).Err(); err != nil {
			return fmt.Errorf("cache: invalidate failed for %s: %w", key, err)
		}
	}

	// 2. Hapus dari local cache
	c.local.Delete(key)

	// 3. Broadcast ke instance lain via Pub/Sub
	if c.pubSub != nil {
		c.pubSub.PublishInvalidate(ctx, key)
	}

	return nil
}

// InvalidatePrefix menghapus semua key dengan prefix tertentu.
// Contoh: InvalidatePrefix("hris:platform:company:") akan menghapus
// semua cache company.
//
// Jika Redis tidak tersedia (local-only mode), hanya menghapus dari local cache.
func (c *Cache) InvalidatePrefix(ctx context.Context, prefix string) error {
	// Local-only mode: hanya hapus dari local cache (kita tidak bisa scan Redis)
	if c.client == nil {
		c.logger.Debug("Cache invalidate prefix (local-only mode)",
			zap.String("prefix", prefix),
		)
		return nil
	}

	iter := c.client.Scan(ctx, 0, prefix+"*", 100).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("cache: scan failed for prefix %s: %w", prefix, err)
	}

	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("cache: invalidate prefix %s failed: %w", prefix, err)
		}

		// Hapus dari local cache
		for _, k := range keys {
			c.local.Delete(k)
		}

		// Broadcast ke instance lain
		if c.pubSub != nil {
			c.pubSub.PublishInvalidate(ctx, keys...)
		}
	}

	return nil
}

// Close menutup koneksi Redis dan Pub/Sub listener.
func (c *Cache) Close() error {
	if c.pubSub != nil {
		c.pubSub.Close()
	}
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// =========================================================================
// Local cache helpers
// =========================================================================

func (c *Cache) getLocal(key string) ([]byte, bool) {
	val, ok := c.local.Load(key)
	if !ok {
		return nil, false
	}

	it := val.(item)
	if time.Now().After(it.expiresAt) {
		c.local.Delete(key)
		return nil, false
	}

	return it.value.([]byte), true
}

func (c *Cache) setLocal(key string, value []byte, ttl time.Duration) {
	c.local.Store(key, item{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	})
}

// evictLocal dipanggil oleh Pub/Sub saat menerima pesan invalidasi
// dari instance lain.
func (c *Cache) evictLocal(keys ...string) {
	for _, key := range keys {
		c.local.Delete(key)
	}
	c.logger.Debug("Cache evicted via Pub/Sub",
		zap.Int("keys", len(keys)),
	)
}

// =========================================================================
// Health check
// =========================================================================

// Ping memeriksa koneksi ke Redis. Di local-only mode, selalu return nil.
func (c *Cache) Ping(ctx context.Context) error {
	if c.client == nil {
		return nil
	}
	return c.client.Ping(ctx).Err()
}
