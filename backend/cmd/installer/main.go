// CLI Installer untuk HRIS Platform
//
// Menyediakan perintah CLI untuk:
//   - Provision tenant baru (create database, run migrations, seed)
//   - Backup & restore tenant
//   - Health check tenant
//
// Usage: go run ./cmd/installer --help
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/pkg/config"
	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/logger"
	"github.com/inthros/hris-platform/internal/pkg/migrator"
	"github.com/inthros/hris-platform/internal/pkg/tenantseed"
	"go.uber.org/zap"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	// Setiap subcommand memuat konfigurasi sendiri setelah parsing flag-nya,
	// sehingga flag --config bisa diberikan per subcommand:
	//   installer migrate --company=<id> --config=./config/config.yaml
	// Jika --config tidak diberikan, config.Load() mencari file config di
	// direktori saat ini dan ./config (mis. ./config/config.yaml saat
	// dijalankan dari folder backend).
	switch cmd {
	case "provision":
		fs := flag.NewFlagSet("provision", flag.ExitOnError)
		cfgPath := fs.String("config", "", "Path to configuration file")
		companyID := fs.String("company", "", "Company ID to provision")
		dbName := fs.String("db-name", "", "Tenant database name (optional, generated from company slug if empty)")
		fs.Parse(args)

		l, dbManager := initApp(*cfgPath)
		defer l.Sync()
		defer dbManager.CloseAll()
		handleProvision(l, dbManager, *companyID, *dbName)

	case "migrate":
		fs := flag.NewFlagSet("migrate", flag.ExitOnError)
		cfgPath := fs.String("config", "", "Path to configuration file")
		companyID := fs.String("company", "", "Company ID to migrate")
		fs.Parse(args)

		l, dbManager := initApp(*cfgPath)
		defer l.Sync()
		defer dbManager.CloseAll()
		handleTenantMigrate(l, dbManager, *companyID)

	case "encrypt-passwords":
		fs := flag.NewFlagSet("encrypt-passwords", flag.ExitOnError)
		cfgPath := fs.String("config", "", "Path to configuration file")
		fs.Parse(args)

		l, dbManager := initApp(*cfgPath)
		defer l.Sync()
		defer dbManager.CloseAll()
		handleEncryptPasswords(l, dbManager)

	case "seed-modules":
		fs := flag.NewFlagSet("seed-modules", flag.ExitOnError)
		cfgPath := fs.String("config", "", "Path to configuration file")
		fs.Parse(args)

		l, dbManager := initApp(*cfgPath)
		defer l.Sync()
		defer dbManager.CloseAll()
		handleSeedModules(l, dbManager)

	case "seed-data":
		fs := flag.NewFlagSet("seed-data", flag.ExitOnError)
		cfgPath := fs.String("config", "", "Path to configuration file")
		seedCompanyID := fs.String("company", "", "Company ID to seed tenant master data (optional)")
		fs.Parse(args)

		l, dbManager := initApp(*cfgPath)
		defer l.Sync()
		defer dbManager.CloseAll()
		if *seedCompanyID != "" {
			// Seed ke tenant database
			tenantDB, err := dbManager.TenantDB(*seedCompanyID)
			if err != nil {
				l.Fatal("Failed to connect to tenant database", zap.Error(err))
			}
			if err := tenantseed.SeedTenantMasterData(tenantDB, l); err != nil {
				l.Fatal("Failed to seed tenant master data", zap.Error(err))
			}
			if err := tenantseed.SeedTenantRBAC(tenantDB, l); err != nil {
				l.Fatal("Failed to seed tenant RBAC defaults", zap.Error(err))
			}
			l.Info("Tenant master data & RBAC seeding completed", zap.String("company_id", *seedCompanyID))
		} else {
			l.Info("No company specified. Use --company=<id> to seed a specific tenant.")
			l.Info("Example: installer seed-data --company=<uuid> --config=./config/config.yaml")
		}

	default:
		printUsage()
		os.Exit(1)
	}
}

// initApp memuat konfigurasi, logger, dan database manager untuk satu subcommand.
// configPath kosong → config.Load() mencari file config di direktori saat ini
// dan ./config (mis. ./config/config.yaml saat dijalankan dari folder backend).
func initApp(configPath string) (*zap.Logger, *database.Manager) {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Init logger
	l := logger.New(cfg.Logger.Level, cfg.Logger.Format, "hris-installer")

	// Init database manager
	dbManager, err := database.NewManager(&database.Config{
		Driver:                  cfg.Database.Driver,
		PlatformDSN:             cfg.Database.PlatformDSN(),
		PlatformHost:            cfg.Database.PlatformHost,
		PlatformPort:            cfg.Database.PlatformPort,
		PlatformUser:            cfg.Database.PlatformUser,
		PlatformPassword:        cfg.Database.PlatformPassword,
		PlatformSSLMode:         cfg.Database.PlatformSSLMode,
		TenantHost:              cfg.Database.TenantHost,
		TenantPort:              cfg.Database.TenantPort,
		TenantSuperUser:         cfg.Database.TenantSuperUser,
		TenantSuperPass:         cfg.Database.TenantSuperPass,
		TenantSSLMode:           cfg.Database.TenantSSLMode,
		MaxOpenConns:            cfg.Database.MaxOpenConns,
		MaxIdleConns:            cfg.Database.MaxIdleConns,
		ConnMaxLifetimeMs:       cfg.Database.ConnMaxLifetimeMs,
		TenantMaxOpenConns:      cfg.Database.TenantMaxOpenConns,
		TenantMaxIdleConns:      cfg.Database.TenantMaxIdleConns,
		TenantConnMaxLifetimeMs: cfg.Database.TenantConnMaxLifetimeMs,
		TenantConnMaxIdleTimeMs: cfg.Database.TenantConnMaxIdleTimeMs,
		LogLevel:                4,
	}, l)
	if err != nil {
		l.Fatal("Failed to initialize database manager", zap.Error(err))
	}

	return l, dbManager
}

func printUsage() {
	fmt.Println("HRIS Platform CLI Installer")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  installer provision --company=<id> [--db-name=<name>] [--config=<path>]")
	fmt.Println("  installer migrate --company=<id> [--config=<path>]")
	fmt.Println("  installer encrypt-passwords [--config=<path>]")
	fmt.Println("  installer seed-modules [--config=<path>]")
	fmt.Println("  installer seed-data [--company=<id>] [--config=<path>]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  provision          Provision a new tenant (create database + run migrations + seed data)")
	fmt.Println("  migrate            Run pending tenant migrations for an existing company")
	fmt.Println("  encrypt-passwords  Encrypt legacy plaintext passwords in tenant_connections")
	fmt.Println("  seed-modules       Register all platform & tenant modules into database")
	fmt.Println("  seed-data          Seed master reference data + RBAC defaults into a tenant DB")
	fmt.Println("")
	fmt.Println("Environment:")
	fmt.Println("  HRIS_ENCRYPTION_KEY    Required for encrypt-passwords (32-byte hex, 64 chars)")
}

func handleProvision(l *zap.Logger, dbManager *database.Manager, companyID, dbName string) {
	if companyID == "" {
		log.Fatal("company is required")
	}

	l.Info("Starting tenant provisioning from CLI",
		zap.String("company_id", companyID),
	)

	if dbName == "" {
		dbName = fmt.Sprintf("hris_tenant_%s", companyID[:8])
	}

	// 1. Create database
	conn, err := dbManager.ProvisionTenant(companyID, dbName, "root", "", dbManager.Driver())
	if err != nil {
		l.Fatal("Failed to create tenant database", zap.Error(err))
	}

	// 2. Simpan TenantConnection
	if err := dbManager.SaveTenantConnection(conn); err != nil {
		l.Fatal("Failed to save tenant connection", zap.Error(err))
	}
	l.Info("Tenant connection saved")

	// 3. Dapatkan koneksi ke tenant DB
	tenantDB, err := dbManager.TenantDB(companyID)
	if err != nil {
		l.Fatal("Failed to connect to tenant database", zap.Error(err))
	}

	// 4. Jalankan tenant migrations (pilih dialect sesuai driver)
	l.Info("Running tenant SQL migrations...")
	tenantRoot := migrator.TenantRootPath(dbManager.Driver())
	tenantMigrator := migrator.New(tenantDB, l, migrator.MigrationsFS, tenantRoot)
	if err := tenantMigrator.Up(); err != nil {
		l.Fatal("Tenant migration failed", zap.Error(err))
	}

	// 5. Seed master reference data (religions, educations, provinces, districts, villages, etc.)
	l.Info("Seeding tenant master data...")
	if err := tenantseed.SeedTenantMasterData(tenantDB, l); err != nil {
		l.Fatal("Failed to seed tenant master data", zap.Error(err))
	}

	// 6. Seed default tenant RBAC (roles Admin/Employee + permissions)
	if err := tenantseed.SeedTenantRBAC(tenantDB, l); err != nil {
		l.Fatal("Failed to seed tenant RBAC defaults", zap.Error(err))
	}

	l.Info("Tenant provisioning completed successfully",
		zap.String("company_id", companyID),
		zap.String("db_name", dbName),
	)
}

func handleEncryptPasswords(l *zap.Logger, dbManager *database.Manager) {
	l.Info("Starting legacy password encryption...")
	l.Warn("Ensure HRIS_ENCRYPTION_KEY environment variable is set before running this command")
	l.Warn("Passwords encrypted with a different key will NOT be recoverable!")

	encrypted, errCount, err := dbManager.EncryptLegacyPasswords()
	if err != nil {
		l.Fatal("Failed to encrypt legacy passwords", zap.Error(err))
	}

	l.Info("Legacy password encryption completed",
		zap.Int("passwords_encrypted", encrypted),
		zap.Int("errors", errCount),
	)

	if encrypted == 0 && errCount == 0 {
		l.Info("No legacy plaintext passwords found — all credentials are already encrypted")
	}
}

func handleTenantMigrate(l *zap.Logger, dbManager *database.Manager, companyID string) {
	if companyID == "" {
		log.Fatal("company is required")
	}

	l.Info("Running tenant migration upgrade",
		zap.String("company_id", companyID),
	)

	tenantDB, err := dbManager.TenantDB(companyID)
	if err != nil {
		l.Fatal("Failed to connect to tenant database", zap.Error(err))
	}

	tenantRoot := migrator.TenantRootPath(dbManager.Driver())
	tenantMigrator := migrator.New(tenantDB, l, migrator.MigrationsFS, tenantRoot)
	if err := tenantMigrator.Up(); err != nil {
		l.Fatal("Tenant migration failed", zap.Error(err))
	}

	l.Info("Tenant migration completed successfully")
}

// moduleDef untuk data modul yang akan di-seed.
type moduleDef struct {
	Name        string
	Slug        string
	Version     string
	Description string
	ModuleType  string // "platform" atau "tenant"
	IsCore      bool
	DependsOn   string
}

// handleSeedModules mendaftarkan semua modul platform & tenant ke database.
// Me-skip modul yang sudah terdaftar (berdasarkan slug).
func handleSeedModules(l *zap.Logger, dbManager *database.Manager) {
	l.Info("Seeding modules into platform database...")

	modules := []moduleDef{
		// Platform modules (core)
		{Name: "Company Management", Slug: "company", Version: "1.0.0", Description: "Manage companies/tenants and their lifecycle", ModuleType: "platform", IsCore: true, DependsOn: ""},
		{Name: "Platform Users & Auth", Slug: "platform-users", Version: "1.0.0", Description: "Platform user authentication, authorization, and user management", ModuleType: "platform", IsCore: true, DependsOn: ""},
		{Name: "Module Management", Slug: "module-management", Version: "1.0.0", Description: "Manage platform modules and their activation for companies", ModuleType: "platform", IsCore: true, DependsOn: ""},
		{Name: "License Management", Slug: "license-management", Version: "1.0.0", Description: "Manage company licenses, plans, and subscription billing", ModuleType: "platform", IsCore: true, DependsOn: "company"},
		{Name: "Package Management", Slug: "package-management", Version: "1.0.0", Description: "Bundle tenant modules with pricing, dependency validation, and publishing", ModuleType: "platform", IsCore: true, DependsOn: "module-management"},

		// Tenant modules
		{Name: "Organization Management", Slug: "organization", Version: "1.0.0", Description: "Manage organizational structure, departments, and positions", ModuleType: "tenant", DependsOn: ""},
		{Name: "Employee Management", Slug: "employee", Version: "1.0.0", Description: "Employee master data, documents, and lifecycle", ModuleType: "tenant", DependsOn: "organization,setting"},
		{Name: "Job Management", Slug: "jobmanagement", Version: "1.0.0", Description: "Job titles, job grades, and position management", ModuleType: "tenant", DependsOn: "organization"},
		{Name: "Competency Management", Slug: "competency", Version: "1.0.0", Description: "Skills, competency dictionaries, and assessments", ModuleType: "tenant", DependsOn: "organization,employee"},
		{Name: "Employee Movement", Slug: "employeemovement", Version: "1.0.0", Description: "Promotions, transfers, resignations, and mutations", ModuleType: "tenant", DependsOn: "employee,organization,setting"},
		{Name: "Attendance Management", Slug: "attendance", Version: "1.0.0", Description: "Time tracking, attendance logs, and overtime", ModuleType: "tenant", DependsOn: "employee,organization"},
		{Name: "Approval Engine", Slug: "approval", Version: "1.0.0", Description: "Approval workflows, multi-level approvals, and delegation", ModuleType: "tenant", DependsOn: "employee,organization"},
		{Name: "Payroll Management", Slug: "payroll", Version: "1.0.0", Description: "Payroll processing, tax calculations, and salary components", ModuleType: "tenant", DependsOn: "employee,organization,setting"},
		{Name: "Leave Management", Slug: "leave", Version: "1.0.0", Description: "Leave requests, balances, and calendars", ModuleType: "tenant", DependsOn: "employee,organization"},
		{Name: "Performance Management", Slug: "performance", Version: "1.0.0", Description: "Performance reviews, KPI tracking, and feedback", ModuleType: "tenant", DependsOn: "organization,employee,jobmanagement,competency,setting"},
		{Name: "Recruitment & Onboarding", Slug: "recruitment", Version: "1.0.0", Description: "Job postings, candidate management, and onboarding", ModuleType: "tenant", DependsOn: "organization,employee,setting"},
		{Name: "Reimbursement", Slug: "reimbursement", Version: "1.0.0", Description: "Expense reimbursement requests and approvals", ModuleType: "tenant", DependsOn: "employee"},
		{Name: "Training Management", Slug: "training", Version: "1.0.0", Description: "Training programs, enrollments, and certifications", ModuleType: "tenant", DependsOn: "employee,organization,setting"},
		{Name: "Workforce Intelligence", Slug: "workforce-intelligence", Version: "1.0.0", Description: "Workforce analytics, planning, and strategic dashboards", ModuleType: "tenant", DependsOn: "organization,employee,attendance,leave,payroll,performance,competency,training,recruitment,employeemovement"},
		{Name: "Career Intelligence", Slug: "career-intelligence", Version: "1.0.0", Description: "Career development, succession planning, and talent management", ModuleType: "tenant", DependsOn: "organization,employee,jobmanagement,competency"},
		{Name: "Settings (Master Data)", Slug: "setting", Version: "1.0.0", Description: "Manage tenant reference data — zones, provinces, regencies, districts, villages, educations, religions, marital statuses, relationship types, banks, employment statuses, nationalities, job families, and salary grades", ModuleType: "tenant", DependsOn: ""},
	}

	db := dbManager.PlatformDB()
	registered := 0
	skipped := 0
	errors := 0

	for _, m := range modules {
		// Cek apakah slug sudah ada
		var count int64
		db.Table("modules").Where("slug = ?", m.Slug).Count(&count)
		if count > 0 {
			l.Debug("Module already registered, skipping", zap.String("slug", m.Slug))
			skipped++
			continue
		}

		// Insert module
		moduleType := m.ModuleType
		if moduleType == "" {
			moduleType = "tenant" // default
		}
		result := db.Table("modules").Create(map[string]interface{}{
			"id":          uuid.New().String(),
			"name":        m.Name,
			"slug":        m.Slug,
			"version":     m.Version,
			"description": m.Description,
			"module_type":  moduleType,
			"is_core":      0, // All tenant modules default to non-core
			"depends_on":   m.DependsOn,
		})

		if m.IsCore {
			db.Table("modules").Where("slug = ?", m.Slug).Update("is_core", 1)
		}

		if result.Error != nil {
			l.Error("Failed to register module", zap.String("slug", m.Slug), zap.Error(result.Error))
			errors++
			continue
		}

		l.Info("Module registered", zap.String("slug", m.Slug), zap.String("name", m.Name))
		registered++
	}

	l.Info("Module seeding completed",
		zap.Int("registered", registered),
		zap.Int("skipped", skipped),
		zap.Int("errors", errors),
	)
}
