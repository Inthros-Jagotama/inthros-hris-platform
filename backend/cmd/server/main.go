// HRIS Platform - Go Modular Monolith
//
// Entry point untuk aplikasi HRIS platform.
// Melakukan inisialisasi shared infrastructure, mendaftarkan
// semua modul (platform & tenant), dan menjalankan HTTP server.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/auth"
	"github.com/inthros/hris-platform/internal/pkg/authz"
	"github.com/inthros/hris-platform/internal/pkg/cache"
	"github.com/inthros/hris-platform/internal/pkg/config"
	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/logger"
	"github.com/inthros/hris-platform/internal/pkg/middleware"
	"github.com/inthros/hris-platform/internal/pkg/module"
	"github.com/inthros/hris-platform/internal/pkg/router"
	"github.com/inthros/hris-platform/internal/pkg/docs"
	"github.com/inthros/hris-platform/internal/pkg/migrator"
	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"

	"github.com/google/uuid"

	// Platform modules
	"github.com/inthros/hris-platform/internal/platform/company"
	"github.com/inthros/hris-platform/internal/platform/license"
	"github.com/inthros/hris-platform/internal/platform/modulemgmt"
	"github.com/inthros/hris-platform/internal/platform/monitoring"
	pkgmgr "github.com/inthros/hris-platform/internal/platform/package"
	"github.com/inthros/hris-platform/internal/platform/user"

	// Tenant modules
	"github.com/inthros/hris-platform/internal/modules/approval"
	"github.com/inthros/hris-platform/internal/modules/attendance"
	"github.com/inthros/hris-platform/internal/modules/competency"
	"github.com/inthros/hris-platform/internal/modules/employee"
	"github.com/inthros/hris-platform/internal/modules/employeemovement"
	"github.com/inthros/hris-platform/internal/modules/jobmanagement"
	"github.com/inthros/hris-platform/internal/modules/leave"
	"github.com/inthros/hris-platform/internal/modules/organization"
	"github.com/inthros/hris-platform/internal/modules/payroll"
	"github.com/inthros/hris-platform/internal/modules/performance"
	"github.com/inthros/hris-platform/internal/modules/recruitment"
	"github.com/inthros/hris-platform/internal/modules/reimbursement"
	"github.com/inthros/hris-platform/internal/modules/training"
	"github.com/inthros/hris-platform/internal/modules/workforceintelligence"
	"github.com/inthros/hris-platform/internal/modules/careerintelligence"
	"github.com/inthros/hris-platform/internal/modules/setting"
)

// =============================================================================
// Adapters
// =============================================================================

// payrollApprovalAdapter implements payroll.ApprovalEngine using the approval service.
type payrollApprovalAdapter struct {
	approvalSvc *approval.Service
}

func (a *payrollApprovalAdapter) CreateApprovalInstance(ctx context.Context, module, documentID, flowID string) (string, error) {
	resp, err := a.approvalSvc.CreateInstance(ctx, approval.CreateInstanceRequest{
		Module:     module,
		DocumentID: documentID,
		FlowID:     flowID,
	})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (a *payrollApprovalAdapter) GetApprovalInstanceStatus(ctx context.Context, instanceID string) (string, error) {
	resp, err := a.approvalSvc.GetInstanceByID(ctx, instanceID)
	if err != nil {
		return "", err
	}
	return resp.Status, nil
}

// licenseCreatorAdapter implements company.LicenseCreator using the license service.
// Digunakan untuk auto-create license saat signup company dengan package.
type licenseCreatorAdapter struct {
	licenseSvc *license.Service
}

func (a *licenseCreatorAdapter) CreateFromPackage(companyID, packageID string) (licenseID, licenseKey, planType string, err error) {
	resp, err := a.licenseSvc.CreateLicense(license.CreateLicenseRequest{
		CompanyID: companyID,
		PlanType:  "pro",
		StartDate: time.Now().Format("2006-01-02"),
		EndDate:   time.Now().AddDate(1, 0, 0).Format("2006-01-02"),
		PackageID: packageID,
	})
	if err != nil {
		return "", "", "", err
	}
	return resp.ID, resp.LicenseKey, resp.PlanType, nil
}

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	migrateDown := flag.Bool("migrate-down", false, "Rollback all applied migrations and exit")
	migrateTo := flag.String("migrate-to", "", "Rollback migrations to specific version (exclusive) and exit")
	flag.Parse()

	// Validate flags: --migrate-down and --migrate-to are mutually exclusive
	if *migrateDown && *migrateTo != "" {
		fmt.Fprintln(os.Stderr, "ERROR: --migrate-down and --migrate-to are mutually exclusive")
		os.Exit(1)
	}

	// 1. Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize logger
	l := logger.New(cfg.Logger.Level, cfg.Logger.Format, "hris-platform")
	defer l.Sync()

	// 3. Initialize distributed cache (Redis + local in-memory + Pub/Sub invalidation)
	cacheManager, err := cache.New(cache.Config{
		RedisAddr:     cfg.Redis.RedisAddr(),
		RedisPassword: cfg.Redis.Password,
		RedisDB:       cfg.Redis.DB,
		DefaultTTL:    time.Duration(cfg.Cache.DefaultTTL) * time.Second,
	}, l)
	if err != nil {
		l.Fatal("Failed to initialize cache", zap.Error(err))
	}
	defer cacheManager.Close()

	// 4. Initialize database manager (multi-tenant, multi-driver)
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
		LogLevel:                4, // Warn
	}, l)
	if err != nil {
		l.Fatal("Failed to initialize database manager", zap.Error(err))
	}
	defer dbManager.CloseAll()

	// 5. Handle migration CLI commands (run and exit without starting server)
	if *migrateDown || *migrateTo != "" {
		cacheManager.Close()
		runMigrationCommand(l, dbManager, *migrateDown, *migrateTo)
		return
	}

	// 5. Initialize JWT auth manager
	authManager := auth.NewManager(auth.Config{
		Secret:          cfg.JWT.Secret,
		AccessTokenTTL:  time.Duration(cfg.JWT.AccessTokenTTL) * time.Minute,
		RefreshTokenTTL: time.Duration(cfg.JWT.RefreshTokenTTL) * time.Hour,
		Issuer:          cfg.JWT.Issuer,
	})

	// 6. Initialize module registry
	var platformModules []module.ModuleRegistration
	var tenantModules []module.ModuleRegistration

	// Create auth middleware once (reused across all platform modules)
	authMW := middleware.AuthJWT(authManager, l)

	// 6a. Register platform modules (ordered by priority)
	// Note: rbacMW diinisialisasi setelah SQL migrations karena membutuhkan
	// tabel RBAC di database. Tapi kita pass nil dulu dan update setelah init.
	//
	// User service dibuat terlebih dahulu karena dibutuhkan oleh
	// Company module untuk membuat company_admin user saat create company.
	userRepo := user.NewRepository(dbManager.PlatformDB())
	userSvc := user.NewService(userRepo, authManager, l)

	// License service dibuat secara eksternal agar bisa di-share
	// dengan LicenseCreator adapter untuk company signup flow.
	licenseRepo := license.NewRepository(dbManager.PlatformDB())
	licenseSvc := license.NewService(licenseRepo, dbManager, l)

	// LicenseCreator adapter wraps license.Service untuk auto-create
	// license saat signup company dengan package.
	licenseCreator := &licenseCreatorAdapter{licenseSvc: licenseSvc}

	// Package service & handler dibuat secara eksternal agar bisa di-share
	// dengan tenant route untuk menampilkan public packages.
	pkgRepo := pkgmgr.NewRepository(dbManager.PlatformDB())
	pkgSvc := pkgmgr.NewService(pkgRepo, l)
	pkgHandler := pkgmgr.NewHandler(pkgSvc)

	// Module Management service dibuat secara eksternal agar bisa di-share
	// dengan subscribe flow untuk auto-activate modules dari package.
	modulemgmtRepo := modulemgmt.NewRepository(dbManager.PlatformDB())
	modulemgmtSvc := modulemgmt.NewService(modulemgmtRepo, dbManager, l)

	// Company module menerima UserCreator (user service) untuk
	// auto-create admin user saat create company, dan LicenseCreator
	// untuk auto-create license dari package.
	platformModules = append(platformModules,
		module.ModuleRegistration{
			Module:   company.NewModule(dbManager, userSvc, licenseCreator, l, authMW, nil),
			TargetDB: module.TargetPlatform,
			Priority: 1,
		},
		module.ModuleRegistration{
			Module:   user.NewModuleWithService(userSvc, l, authMW, nil),
			TargetDB: module.TargetPlatform,
			Priority: 2,
		},
		module.ModuleRegistration{
			Module:   modulemgmt.NewModuleWithService(modulemgmtSvc, l, authMW, nil),
			TargetDB: module.TargetPlatform,
			Priority: 3,
		},
		module.ModuleRegistration{
			Module:   license.NewModuleWithService(licenseSvc, l, authMW, nil),
			TargetDB: module.TargetPlatform,
			Priority: 4,
		},
		module.ModuleRegistration{
			Module:   pkgmgr.NewModuleWithService(pkgSvc, l, authMW, nil),
			TargetDB: module.TargetPlatform,
			Priority: 5,
		},
	)

	// 6b. Create approval engine adapter (for payroll integration)
	// Approval service needs a tenant DB resolver like other modules
	approvalResolver := func(ctx context.Context) (*gorm.DB, error) {
		companyID, ok := ctx.Value("company_id").(string)
		if !ok || companyID == "" {
			return nil, fmt.Errorf("tenant context not found in request: company_id is required")
		}
		return dbManager.TenantDB(companyID)
	}
	approvalLogger := l.Named("approval")
	approvalRepo := approval.NewRepository(approvalResolver)
	approvalSvc := approval.NewService(approvalRepo, approvalLogger)

	// Create the ApprovalEngine adapter that the payroll module expects
	payrollApprovalEngine := &payrollApprovalAdapter{
		approvalSvc: approvalSvc,
	}

	// 6c. Register tenant modules (ordered by priority & dependency)
	tenantModules = append(tenantModules,
		module.ModuleRegistration{
			Module:   organization.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 1,
		},
		module.ModuleRegistration{
			Module:   employee.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 2,
		},
		module.ModuleRegistration{
			Module:   jobmanagement.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 3,
		},
		module.ModuleRegistration{
			Module:   competency.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 4,
		},
		module.ModuleRegistration{
			Module:   employeemovement.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 5,
		},
		module.ModuleRegistration{
			Module:   attendance.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 6,
		},
		module.ModuleRegistration{
			Module:   approval.NewModuleWithService(dbManager, l, approvalSvc),
			TargetDB: module.TargetTenant,
			Priority: 7,
		},
		module.ModuleRegistration{
			Module:   payroll.NewModule(dbManager, l, payrollApprovalEngine),
			TargetDB: module.TargetTenant,
			Priority: 8,
		},
		module.ModuleRegistration{
			Module:   leave.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 9,
		},
		module.ModuleRegistration{
			Module:   performance.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 10,
		},
		module.ModuleRegistration{
			Module:   recruitment.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 11,
		},
		module.ModuleRegistration{
			Module:   reimbursement.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 12,
		},
		module.ModuleRegistration{
			Module:   training.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 13,
		},
		module.ModuleRegistration{
			Module:   workforceintelligence.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 14,
		},
		module.ModuleRegistration{
			Module:   careerintelligence.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 15,
		},
		module.ModuleRegistration{
			Module:   setting.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 16,
		},
	)

	// 7. Run SQL file migrations for platform modules
	l.Info("Running platform SQL migrations...")
	platformMigrator := migrator.New(dbManager.PlatformDB(), l, migrator.MigrationsFS, migrator.RootPlatform)
	if err := platformMigrator.Up(); err != nil {
		l.Fatal("Platform SQL migration failed", zap.Error(err))
	}

	// 8. Run AutoMigrate for platform modules (sync GORM models to schema)
	l.Info("Running platform AutoMigrate...")
	for _, reg := range platformModules {
		if err := reg.Module.Migrate(dbManager.PlatformDB()); err != nil {
			l.Fatal("Platform AutoMigrate failed",
				zap.String("module", reg.Module.Info().Name),
				zap.Error(err),
			)
		}
		l.Info("Platform AutoMigrate completed",
			zap.String("module", reg.Module.Info().Name),
		)
	}

	// Note: Tenant migrations run during tenant provisioning,
	// not at startup. Each tenant gets its own database.
	l.Info("Running SQL seeders...")
	seederMigrator := migrator.New(dbManager.PlatformDB(), l, migrator.MigrationsFS, migrator.RootSeeders)
	if err := seederMigrator.Up(); err != nil {
		l.Warn("SQL seeder warning", zap.Error(err))
	}

	// 10. Initialize database-backed RBAC enforcer (seeds defaults if empty)
	l.Info("Initializing database-backed RBAC enforcer...")
	rbacEnforcer, err := authz.NewEnforcerFromDB(dbManager.PlatformDB())
	if err != nil {
		l.Fatal("Failed to initialize RBAC enforcer", zap.Error(err))
	}
	rbacMW := authz.NewMiddleware(authz.MiddlewareConfig{Enforcer: rbacEnforcer})

	// Update platform modules dengan RBAC middleware yang sebenarnya
	// (sebelumnya modules dibuat dengan nil karena RBAC butuh DB siap)
	l.Info("Wiring RBAC middleware to platform modules...")
	for _, reg := range platformModules {
		if setter, ok := reg.Module.(interface{ SetRBACMiddleware(gin.HandlerFunc) }); ok {
			setter.SetRBACMiddleware(rbacMW)
			l.Debug("RBAC middleware wired",
				zap.String("module", reg.Module.Info().Name),
			)
		}
	}

	// 10a. Create RBAC service & handler for management API
	rbacLogger := l.Named("rbac")
	rbacRepo := authz.NewRepository(dbManager.PlatformDB())
	rbacService := authz.NewService(rbacRepo, rbacEnforcer, rbacLogger)
	rbacHandler := authz.NewHandler(rbacService)

	// 10b. Run module seeders for platform modules
	l.Info("Running platform module seeders...")
	for _, reg := range platformModules {
		if err := reg.Module.Seed(dbManager.PlatformDB()); err != nil {
			l.Warn("Platform seeder warning",
				zap.String("module", reg.Module.Info().Name),
				zap.Error(err),
			)
		}
	}

	// 11. Setup router and middleware
	r := router.Setup(
		router.Config{Mode: cfg.Server.Mode},
		middleware.AuthJWT(authManager, l),
		middleware.TenantRequired(),
		rbacMW,
		middleware.CORS(middleware.CORSConfig{
			AllowedOrigins:   cfg.CORS.AllowedOrigins,
			AllowedMethods:   cfg.CORS.AllowedMethods,
			AllowedHeaders:   cfg.CORS.AllowedHeaders,
			AllowCredentials: cfg.CORS.AllowCredentials,
			MaxAge:           cfg.CORS.MaxAge,
		}),
		middleware.RequestLogger(l),
		middleware.Recovery(l),
		platformModules,
		tenantModules,
	)

	// Register platform RBAC management routes (standalone)
	platformGroup := r.Group("/api/v1/platform")
	platformGroup.Use(authMW)
	platformGroup.Use(rbacMW)
	authz.RegisterRoutes(platformGroup, rbacHandler)

	// Register platform monitoring routes (standalone, no module interface needed)
	monitoringHandler := monitoring.NewHandler(dbManager, cacheManager, l)
	monitoring.RegisterRoutes(platformGroup, monitoringHandler, authMW, rbacMW)

	// Register tenant packages routes — authenticated tenant users can browse
	// published packages and subscribe (upgrade/downgrade) their company license.
	tenantPkgGroup := r.Group("/api/v1/tenant")
	tenantPkgGroup.Use(authMW)
	tenantPkgGroup.Use(middleware.TenantRequired())
	{
		tenantPkgGroup.GET("/packages", pkgHandler.ListPublishedPackages)
		tenantPkgGroup.POST("/packages/:id/subscribe", func(c *gin.Context) {
			packageID := c.Param("id")
			companyID := middleware.GetCompanyID(c)

			// Validate package exists and is published
			pkg, err := pkgSvc.GetPackage(packageID)
			if err != nil {
				httputil.NotFound(c, "Package not found")
				return
			}
			if pkg.Status != "published" {
				httputil.ErrorRaw(c, 400, "VALIDATION_ERROR", "Package is not available for subscription")
				return
			}

			// Create license via adapter (creates a new license for this company + package)
			licenseID, licenseKey, planType, err := licenseCreator.CreateFromPackage(companyID, packageID)
			if err != nil {
				httputil.InternalError(c, err.Error())
				return
			}

			// Auto-activate all modules included in the package
			var activatedModules []string
			if len(pkg.Modules) > 0 {
				for _, m := range pkg.Modules {
					_, err := modulemgmtSvc.ActivateModule(m.ModuleID, companyID)
					if err != nil {
						l.Warn("Failed to auto-activate module during subscribe",
							zap.String("module_id", m.ModuleID),
							zap.String("module_name", m.ModuleName),
							zap.String("company_id", companyID),
							zap.Error(err),
						)
						continue
					}
					activatedModules = append(activatedModules, m.ModuleName)
				}
			}

			httputil.CreatedJSON(c, gin.H{
				"license_id":        licenseID,
				"license_key":       licenseKey,
				"plan_type":         planType,
				"package_id":        packageID,
				"package_name":      pkg.Name,
				"activated_modules": activatedModules,
			}, "package.subscribed")
		})

		tenantPkgGroup.POST("/packages/:id/unsubscribe", func(c *gin.Context) {
			packageID := c.Param("id")
			companyID := middleware.GetCompanyID(c)

			// Validate package exists
			pkg, err := pkgSvc.GetPackage(packageID)
			if err != nil {
				httputil.NotFound(c, "Package not found")
				return
			}

			// Deactivate all modules from the package
			var deactivatedModules []string
			if len(pkg.Modules) > 0 {
				for _, m := range pkg.Modules {
					_, err := modulemgmtSvc.DeactivateModule(m.ModuleID, companyID)
					if err != nil {
						l.Warn("Failed to deactivate module during unsubscribe",
							zap.String("module_id", m.ModuleID),
							zap.String("module_name", m.ModuleName),
							zap.String("company_id", companyID),
							zap.Error(err),
						)
						continue
					}
					deactivatedModules = append(deactivatedModules, m.ModuleName)
				}
			}

			// Suspend the active license associated with this package
			// Find license by company ID and package ID via license service
			cid, _ := uuid.Parse(companyID)
			pid, _ := uuid.Parse(packageID)
			license, err := licenseRepo.FindByCompanyIDAndPackageID(cid, pid)
			if err == nil && license != nil {
				license.Status = "suspended"
				if err := licenseRepo.Update(license); err != nil {
					l.Warn("Failed to suspend license during unsubscribe",
						zap.String("license_id", license.ID.String()),
						zap.Error(err),
					)
				}
			} else {
				l.Warn("No active license found for company and package during unsubscribe",
					zap.String("company_id", companyID),
					zap.String("package_id", packageID),
				)
			}

			httputil.MessageJSON(c, "package.unsubscribed")
		})
	}

	// Register public packages endpoint (no auth)
	r.GET("/api/v1/public/packages", pkgHandler.ListPublishedPackages)

	// Register Scalar API Documentation
	r.GET("/docs", docs.ScalarUIHandler())
	r.GET("/openapi.json", docs.OpenAPIHandler())

	// 12. Start server
	l.Info("Starting HRIS Platform server",
		zap.String("port", cfg.Server.Port),
		zap.String("mode", cfg.Server.Mode),
	)
	l.Info("API Documentation available at",
		zap.String("url", "/docs"),
	)

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		l.Fatal("Failed to start server", zap.Error(err))
	}
}

// runMigrationCommand mengeksekusi perintah migration CLI dan exit.
// Digunakan untuk --migrate-down dan --migrate-to flags.
func runMigrationCommand(l *zap.Logger, dbManager *database.Manager, down bool, to string) {
	l.Info("Migration command detected, running in CLI mode")

	platformMigrator := migrator.New(dbManager.PlatformDB(), l, migrator.MigrationsFS, migrator.RootPlatform)
	seederMigrator := migrator.New(dbManager.PlatformDB(), l, migrator.MigrationsFS, migrator.RootSeeders)

	if down {
		// Rollback all: seeders first (reverse), then platform
		l.Info("Rolling back all seeders...")
		if err := seederMigrator.Down(); err != nil {
			l.Fatal("Seeder rollback failed", zap.Error(err))
		}

		l.Info("Rolling back all platform migrations...")
		if err := platformMigrator.Down(); err != nil {
			l.Fatal("Platform migration rollback failed", zap.Error(err))
		}
	} else if to != "" {
		// Rollback to specific version
		l.Info("Rolling back platform migrations to version",
			zap.String("target", to))
		if err := platformMigrator.DownTo(to); err != nil {
			l.Fatal("Platform migration rollback failed", zap.Error(err))
		}

		l.Info("Rolling back seeders to version",
			zap.String("target", to))
		if err := seederMigrator.DownTo(to); err != nil {
			l.Warn("Seeder rollback warning", zap.Error(err))
		}
	}

	l.Info("Migration command completed successfully")
}
