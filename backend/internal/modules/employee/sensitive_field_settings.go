package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
)

// SensitiveFieldSetting adalah baris toggle enkripsi-at-rest untuk satu
// field sensitif. Dibuat oleh migration SQL (151_sensitive_field_settings),
// bukan AutoMigrate — lihat catatan Migrate() di module.go.
type SensitiveFieldSetting struct {
	ID                  string    `gorm:"type:char(36);primaryKey" json:"id"`
	FieldKey            string    `gorm:"type:varchar(100);uniqueIndex" json:"field_key"`
	IsEncryptionEnabled bool      `gorm:"column:is_encryption_enabled" json:"is_encryption_enabled"`
	UpdatedBy           *string   `gorm:"type:char(36)" json:"updated_by,omitempty"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (SensitiveFieldSetting) TableName() string { return "sensitive_field_settings" }

// ListSensitiveFieldSettings mengembalikan seluruh baris setting field sensitif.
func (r *Repository) ListSensitiveFieldSettings(ctx context.Context) ([]SensitiveFieldSetting, error) {
	db, err := r.dbResolver(ctx)
	if err != nil {
		return nil, err
	}
	var settings []SensitiveFieldSetting
	if err := db.WithContext(ctx).Order("field_key").Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

// SetSensitiveFieldEnabled meng-update toggle enkripsi satu field.
func (r *Repository) SetSensitiveFieldEnabled(ctx context.Context, fieldKey string, enabled bool, updatedBy *uuid.UUID) error {
	db, err := r.dbResolver(ctx)
	if err != nil {
		return err
	}
	updates := map[string]interface{}{
		"is_encryption_enabled": enabled,
		"updated_at":            time.Now(),
	}
	if updatedBy != nil {
		id := updatedBy.String()
		updates["updated_by"] = id
	}
	result := db.WithContext(ctx).Model(&SensitiveFieldSetting{}).
		Where("field_key = ?", fieldKey).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("sensitive field setting not found: %s", fieldKey)
	}
	return nil
}

// ListSensitiveFieldSettings (Service) — passthrough untuk handler.
func (s *Service) ListSensitiveFieldSettings(ctx context.Context) ([]SensitiveFieldSetting, error) {
	return s.repo.ListSensitiveFieldSettings(ctx)
}

// SetSensitiveFieldEnabled mengubah toggle enkripsi untuk satu field,
// setelah memvalidasi field_key ada di SensitiveFieldRegistry.
func (s *Service) SetSensitiveFieldEnabled(ctx context.Context, fieldKey string, enabled bool) error {
	if _, ok := FieldDef(fieldKey); !ok {
		return fmt.Errorf("unknown sensitive field key: %s", fieldKey)
	}
	updatedBy := authctx.GetUserID(ctx)
	return s.repo.SetSensitiveFieldEnabled(ctx, fieldKey, enabled, updatedBy)
}

// IsFieldEncryptionEnabled mengecek apakah field tertentu sedang di-toggle
// aktif enkripsinya. Dipanggil sebelum menulis nilai field sensitif.
func (s *Service) IsFieldEncryptionEnabled(ctx context.Context, fieldKey string) (bool, error) {
	settings, err := s.repo.ListSensitiveFieldSettings(ctx)
	if err != nil {
		return false, err
	}
	for _, st := range settings {
		if st.FieldKey == fieldKey {
			return st.IsEncryptionEnabled, nil
		}
	}
	return false, fmt.Errorf("unknown sensitive field key: %s", fieldKey)
}
