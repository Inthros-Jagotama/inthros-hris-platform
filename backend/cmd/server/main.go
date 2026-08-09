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

	"github.com/gin-gonic/gin"
	"github.com/inthros/hris-platform/internal/pkg/auth"
	"github.com/inthros/hris-platform/internal/pkg/authz"
	"github.com/inthros/hris-platform/internal/pkg/cache"
	"github.com/inthros/hris-platform/internal/pkg/config"
	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/docs"
	"github.com/inthros/hris-platform/internal/pkg/logger"
	"github.com/inthros/hris-platform/internal/pkg/mailer"
	"github.com/inthros/hris-platform/internal/pkg/middleware"
	"github.com/inthros/hris-platform/internal/pkg/migrator"
	"github.com/inthros/hris-platform/internal/pkg/module"
	"github.com/inthros/hris-platform/internal/pkg/onpremise"
	"github.com/inthros/hris-platform/internal/pkg/router"
	"github.com/inthros/hris-platform/internal/pkg/upload"

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
	"github.com/inthros/hris-platform/internal/modules/careerintelligence"
	"github.com/inthros/hris-platform/internal/modules/competency"
	"github.com/inthros/hris-platform/internal/modules/employee"
	"github.com/inthros/hris-platform/internal/modules/employeemovement"
	"github.com/inthros/hris-platform/internal/modules/jobmanagement"
	"github.com/inthros/hris-platform/internal/modules/leave"
	"github.com/inthros/hris-platform/internal/modules/notification"
	"github.com/inthros/hris-platform/internal/modules/organization"
	"github.com/inthros/hris-platform/internal/modules/payroll"
	"github.com/inthros/hris-platform/internal/modules/performance"
	"github.com/inthros/hris-platform/internal/modules/rbac"
	"github.com/inthros/hris-platform/internal/modules/recruitment"
	"github.com/inthros/hris-platform/internal/modules/reimbursement"
	"github.com/inthros/hris-platform/internal/modules/setting"
	"github.com/inthros/hris-platform/internal/modules/training"
	"github.com/inthros/hris-platform/internal/modules/useraccount"
	"github.com/inthros/hris-platform/internal/modules/workforceintelligence"
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

// GetActiveFlowIDForModule lets a consumer (e.g. KPI self-assessment) auto-
// resolve which flow to use for a module instead of picking a flow_id
// manually.
func (a *payrollApprovalAdapter) GetActiveFlowIDForModule(ctx context.Context, module string) (string, error) {
	resp, err := a.approvalSvc.GetActiveFlowByModule(ctx, module)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

// CancelApprovalInstance membatalkan instance approval yang masih aktif
// (dipakai attendance saat request lembur dibatalkan sebelum isian aktual,
// §32b docs/module-attendance-plan.md) supaya task approver tidak menggantung.
func (a *payrollApprovalAdapter) CancelApprovalInstance(ctx context.Context, instanceID string) error {
	return a.approvalSvc.CancelInstance(ctx, instanceID)
}

// licenseCreatorAdapter implements company.LicenseCreator using the license service.
// Digunakan untuk auto-create license saat signup company dengan package.
type licenseCreatorAdapter struct {
	licenseSvc *license.Service
}

func (a *licenseCreatorAdapter) CreateFromPackage(companyID, packageID string) (licenseID, licenseKey, planType string, err error) {
	resp, err := a.licenseSvc.CreateLicense(license.CreateLicenseRequest{
		CompanyID: companyID,
		PlanType:  "subscription",
		StartDate: time.Now().Format("2006-01-02"),
		EndDate:   time.Now().AddDate(1, 0, 0).Format("2006-01-02"),
		PackageID: packageID,
	})
	if err != nil {
		return "", "", "", err
	}
	return resp.ID, resp.LicenseKey, resp.PlanType, nil
}

// packageModuleAdapter implements license.PackageModuleManager using the package service
// and module management service. Digunakan untuk auto-activate/deactivate modules
// saat license dibuat/diupdate dengan package.
type packageModuleAdapter struct {
	pkgSvc        *pkgmgr.Service
	modulemgmtSvc *modulemgmt.Service
	logger        *zap.Logger
}

func (a *packageModuleAdapter) ActivatePackageModules(packageID, companyID string) ([]string, error) {
	pkg, err := a.pkgSvc.GetPackage(packageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get package: %w", err)
	}

	var activated []string
	for _, m := range pkg.Modules {
		_, err := a.modulemgmtSvc.ActivateModule(m.ModuleID, companyID)
		if err != nil {
			a.logger.Warn("Failed to activate module during package activation",
				zap.String("module_id", m.ModuleID),
				zap.String("module_name", m.ModuleName),
				zap.String("company_id", companyID),
				zap.Error(err),
			)
			continue
		}
		activated = append(activated, m.ModuleName)
	}
	return activated, nil
}

func (a *packageModuleAdapter) DeactivatePackageModules(packageID, companyID string) ([]string, error) {
	pkg, err := a.pkgSvc.GetPackage(packageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get package: %w", err)
	}

	var deactivated []string
	for _, m := range pkg.Modules {
		_, err := a.modulemgmtSvc.DeactivateModule(m.ModuleID, companyID)
		if err != nil {
			a.logger.Warn("Failed to deactivate module during package deactivation",
				zap.String("module_id", m.ModuleID),
				zap.String("module_name", m.ModuleName),
				zap.String("company_id", companyID),
				zap.Error(err),
			)
			continue
		}
		deactivated = append(deactivated, m.ModuleName)
	}
	return deactivated, nil
}

// companyModuleListerAdapter membungkus modulemgmt.Service agar memenuhi
// interface middleware.CompanyModuleLister — hanya slug modul yang Enabled
// yang dikembalikan (sumber data untuk PlatformLicenseMiddleware mode SaaS).
type companyModuleListerAdapter struct {
	svc *modulemgmt.Service
}

// EnabledModuleSlugs mengembalikan daftar slug modul aktif untuk company.
func (a companyModuleListerAdapter) EnabledModuleSlugs(companyID string) ([]string, error) {
	mods, err := a.svc.ListCompanyModules(companyID)
	if err != nil {
		return nil, err
	}

	var slugs []string
	for _, m := range mods {
		if m.Enabled {
			slugs = append(slugs, m.ModuleSlug)
		}
	}
	return slugs, nil
}

// approvalModuleCheckerAdapter membungkus modulemgmt.Service agar memenuhi
// approval.ModuleSubscriptionChecker — dipakai approval flow CRUD untuk
// memvalidasi bahwa module target sudah disubscribe tenant sebelum flow
// dibuat/diaktifkan.
type approvalModuleCheckerAdapter struct {
	svc *modulemgmt.Service
}

func (a approvalModuleCheckerAdapter) IsModuleActive(companyID, moduleSlug string) (bool, error) {
	mods, err := a.svc.ListCompanyModules(companyID)
	if err != nil {
		return false, err
	}
	for _, m := range mods {
		if m.ModuleSlug == moduleSlug && m.Enabled {
			return true, nil
		}
	}
	return false, nil
}

func (a approvalModuleCheckerAdapter) ListActiveModules(companyID string) ([]string, error) {
	mods, err := a.svc.ListCompanyModules(companyID)
	if err != nil {
		return nil, err
	}
	var slugs []string
	for _, m := range mods {
		if m.Enabled {
			slugs = append(slugs, m.ModuleSlug)
		}
	}
	return slugs, nil
}

// onPremiseLister mengimplementasikan middleware.CompanyModuleLister untuk mode
// On-Premise: daftar modul yang diizinkan diambil dari file .lic RSA
// (berlaku untuk semua company pada deployment tunggal tersebut).
type onPremiseLister struct {
	lic *onpremise.License
}

// EnabledModuleSlugs mengembalikan allowed_modules dari lisensi on-premise.
func (o onPremiseLister) EnabledModuleSlugs(companyID string) ([]string, error) {
	return o.lic.AllowedModules, nil
}

// onPremiseQuotaChecker mengimplementasikan employee.EmployeeQuotaChecker untuk
// mode on-premise — batas max_employees diambil dari file .lic RSA.
type onPremiseQuotaChecker struct {
	lic *onpremise.License
}

// MaxEmployees mengembalikan batas jumlah employee dari lisensi.
// Nilai <= 0 berarti tanpa batas (unlimited).
func (o onPremiseQuotaChecker) MaxEmployees() int {
	if o.lic == nil {
		return 0
	}
	return o.lic.MaxEmployees
}

// newEmployeeModule membuat employee module. Bila mode on-premise aktif
// (lic != nil), quota checker max_employees di-injeksi ke employee service
// agar pembuatan employee ditolak saat batas tercapai.
func newEmployeeModule(dbManager *database.Manager, l *zap.Logger, lic *onpremise.License) module.Module {
	m := employee.NewModule(dbManager, l)
	if lic != nil {
		if setter, ok := m.(interface {
			SetQuotaChecker(employee.EmployeeQuotaChecker)
		}); ok {
			setter.SetQuotaChecker(onPremiseQuotaChecker{lic: lic})
		} else {
			// Fail-loud: jika module tidak mengimplementasikan setter,
			// kuota tidak ter-enforce secara senyap (hindari misconfig).
			l.Warn("Employee module does not implement SetQuotaChecker — quota enforcement disabled")
		}
	}
	return m
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

	// 4b. Initialize SMTP mailer (untuk link set-password employee via email)
	mailerSvc := mailer.New(mailer.Config{
		Host:            cfg.SMTP.Host,
		Port:            cfg.SMTP.Port,
		Username:        cfg.SMTP.Username,
		Password:        cfg.SMTP.Password,
		From:            cfg.SMTP.From,
		FrontendBaseURL: cfg.SMTP.FrontendBaseURL,
	})

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

	// Package service & handler dibuat secara eksternal agar bisa di-share
	// dengan tenant route untuk menampilkan public packages.
	pkgRepo := pkgmgr.NewRepository(dbManager.PlatformDB())
	pkgSvc := pkgmgr.NewService(pkgRepo, l)
	pkgHandler := pkgmgr.NewHandler(pkgSvc)

	// Module Management service dibuat secara eksternal agar bisa di-share
	// dengan subscribe flow untuk auto-activate modules dari package.
	modulemgmtRepo := modulemgmt.NewRepository(dbManager.PlatformDB())
	modulemgmtSvc := modulemgmt.NewService(modulemgmtRepo, dbManager, l)
	// Injeksi cache manager untuk invalidasi license cache saat modul
	// diaktifkan/dinonaktifkan (PlatformLicenseMiddleware membaca dari cache ini).
	modulemgmtSvc.SetCacheManager(cacheManager)

	// PackageModuleAdapter — license service akan auto-activate modules
	// saat license dibuat/diupdate dengan package.
	pkgModuleMgr := &packageModuleAdapter{
		pkgSvc:        pkgSvc,
		modulemgmtSvc: modulemgmtSvc,
		logger:        l.Named("pkg-module-mgr"),
	}

	// License service dibuat secara eksternal agar bisa di-share
	// dengan LicenseCreator adapter untuk company signup flow.
	// Dilengkapi dengan PackageModuleManager untuk auto-activate modules.
	licenseRepo := license.NewRepository(dbManager.PlatformDB())
	licenseSvc := license.NewService(licenseRepo, dbManager, pkgModuleMgr, l)

	// LicenseCreator adapter wraps license.Service untuk auto-create
	// license saat signup company dengan package.
	licenseCreator := &licenseCreatorAdapter{licenseSvc: licenseSvc}

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
	approvalSvc.SetModuleChecker(approvalModuleCheckerAdapter{svc: modulemgmtSvc})

	// Create the ApprovalEngine adapter shared by every consumer module
	// (payroll, leave, ...) — structurally satisfies each module's own
	// narrow ApprovalEngine interface since they all share the same method set.
	sharedApprovalEngine := &payrollApprovalAdapter{
		approvalSvc: approvalSvc,
	}

	// Construct the payroll service up front (instead of inside payroll.NewModule)
	// so its push-based approval status handler can be registered with
	// approvalSvc before the module is mounted.
	payrollResolver := payroll.NewTenantDBResolver(dbManager)
	payrollRepo := payroll.NewRepository(payrollResolver)
	payrollSvc := payroll.NewService(payrollRepo, l.Named("payroll"))
	payrollSvc.SetApprovalEngine(sharedApprovalEngine)
	approvalSvc.RegisterStatusHandler("payroll", func(ctx context.Context, documentID uuid.UUID, status approval.InstanceStatus, note string) error {
		return payrollSvc.HandleApprovalStatusChange(ctx, documentID, string(status), note)
	})

	// Construct the notification service up front (instead of inside
	// notification.NewModule) so it can be wired as a Notifier into consumer
	// modules (Leave, etc.) before those modules are mounted.
	notificationResolver := notification.NewTenantDBResolver(dbManager)
	notificationRepo := notification.NewRepository(notificationResolver)
	notificationSvc := notification.NewService(notificationRepo, l.Named("notification"))

	// Wire the central Approval module itself as a notification producer —
	// every module routed through it (leave, payroll, attendance, ...)
	// automatically gets "you have a pending approval" pushes for new
	// USER/ROLE/ORGANIZATION-resolved tasks, without each consumer module
	// wiring this separately.
	approvalSvc.SetNotifier(notificationSvc)

	// Construct the leave service up front (instead of inside leave.NewModule)
	// so its push-based approval status handler can be registered with
	// approvalSvc before the module is mounted. Reuses the same generic
	// approval adapter as payroll.
	leaveResolver := leave.NewTenantDBResolver(dbManager)
	leaveRepo := leave.NewRepository(leaveResolver)
	leaveSvc := leave.NewService(leaveRepo, l.Named("leave"))
	leaveSvc.SetApprovalEngine(sharedApprovalEngine)
	leaveSvc.SetNotifier(notificationSvc)
	approvalSvc.RegisterStatusHandler("leave", func(ctx context.Context, documentID uuid.UUID, status approval.InstanceStatus, note string) error {
		return leaveSvc.HandleApprovalStatusChange(ctx, documentID, string(status), note)
	})
	// Separate setting.Service instance used only for its ListHolidayDatesInRange
	// read (company holiday calendar), so leave's working-day calculation can
	// exclude holidays without importing setting's module wiring directly.
	settingRepoForLeave := setting.NewRepository(setting.NewTenantDBResolver(dbManager))
	leaveSvc.SetHolidayProvider(setting.NewService(settingRepoForLeave, l.Named("setting")))

	// Construct the reimbursement service up front (instead of inside
	// reimbursement.NewModule) so its push-based approval status handler can
	// be registered with approvalSvc before the module is mounted.
	reimbursementResolver := reimbursement.NewTenantDBResolver(dbManager)
	reimbursementRepo := reimbursement.NewRepository(reimbursementResolver)
	reimbursementSvc := reimbursement.NewService(reimbursementRepo, l.Named("reimbursement"))
	reimbursementSvc.SetApprovalEngine(sharedApprovalEngine)
	approvalSvc.RegisterStatusHandler("reimbursement", func(ctx context.Context, documentID uuid.UUID, status approval.InstanceStatus, note string) error {
		return reimbursementSvc.HandleApprovalStatusChange(ctx, documentID, string(status), note)
	})

	// Construct the employeemovement service up front (instead of inside
	// employeemovement.NewModule) so its push-based approval status handler
	// can be registered with approvalSvc before the module is mounted.
	employeeMovementResolver := employeemovement.NewTenantDBResolver(dbManager)
	employeeMovementRepo := employeemovement.NewRepository(employeeMovementResolver)
	employeeMovementSvc := employeemovement.NewService(employeeMovementRepo, l.Named("employeemovement"))
	employeeMovementSvc.SetApprovalEngine(sharedApprovalEngine)
	approvalSvc.RegisterStatusHandler("employeemovement", func(ctx context.Context, documentID uuid.UUID, status approval.InstanceStatus, note string) error {
		return employeeMovementSvc.HandleApprovalStatusChange(ctx, documentID, string(status), note)
	})

	// Construct the attendance service up front (instead of inside
	// attendance.NewModule) so its push-based approval status handler can be
	// registered with approvalSvc before the module is mounted.
	attendanceResolver := attendance.NewTenantDBResolver(dbManager)
	attendanceRepo := attendance.NewRepository(attendanceResolver)
	attendanceSvc := attendance.NewService(attendanceRepo, l.Named("attendance"))
	attendanceSvc.SetApprovalEngine(sharedApprovalEngine)
	attendanceSvc.SetNotifier(notificationSvc)
	approvalSvc.RegisterStatusHandler("attendance", func(ctx context.Context, documentID uuid.UUID, status approval.InstanceStatus, note string) error {
		return attendanceSvc.HandleApprovalStatusChange(ctx, documentID, string(status), note)
	})
	// Push approved leave onto Attendance's daily session (§26/§50 of
	// docs/module-attendance-plan.md) so a day fully covered by leave is
	// reflected without needing a check-in event to trigger it.
	leaveSvc.SetAttendanceSessionUpdater(attendanceSvc)

	// Construct the performance service up front (instead of inside
	// performance.NewModule) so its two push-based approval status handlers
	// (KPI target approval, KPI realization approval — two independent
	// checkpoints on the same evaluation) can be registered with approvalSvc
	// before the module is mounted.
	performanceResolver := performance.NewTenantDBResolver(dbManager)
	performanceRepo := performance.NewRepository(performanceResolver)
	performanceSvc := performance.NewService(performanceRepo, l.Named("performance"))
	performanceSvc.SetApprovalEngine(sharedApprovalEngine)
	performanceSvc.SetNotifier(notificationSvc)
	approvalSvc.RegisterStatusHandler(performance.ApprovalModuleKPITarget, func(ctx context.Context, documentID uuid.UUID, status approval.InstanceStatus, note string) error {
		return performanceSvc.HandleTargetApprovalStatusChange(ctx, documentID, string(status), note)
	})
	approvalSvc.RegisterStatusHandler(performance.ApprovalModuleKPIRealization, func(ctx context.Context, documentID uuid.UUID, status approval.InstanceStatus, note string) error {
		return performanceSvc.HandleRealizationApprovalStatusChange(ctx, documentID, string(status), note)
	})

	// Construct the OKR service up front (instead of inside performance.NewModule)
	// so its two push-based approval status handlers (Key Result approval,
	// assessment approval — two independent checkpoints on the same
	// evaluation) can be registered with approvalSvc before the module is
	// mounted. Reuses the same shared approval adapter as performance (KPI).
	okrRepo := performance.NewOKRRepository()
	okrSvc := performance.NewOKRService(okrRepo, performanceResolver)
	okrSvc.SetApprovalEngine(sharedApprovalEngine)
	okrSvc.SetNotifier(notificationSvc)
	approvalSvc.RegisterStatusHandler(performance.ApprovalModuleOKRKeyResult, func(ctx context.Context, documentID uuid.UUID, status approval.InstanceStatus, note string) error {
		return okrSvc.HandleKeyResultApprovalStatusChange(ctx, documentID, string(status), note)
	})
	approvalSvc.RegisterStatusHandler(performance.ApprovalModuleOKRAssessment, func(ctx context.Context, documentID uuid.UUID, status approval.InstanceStatus, note string) error {
		return okrSvc.HandleAssessmentApprovalStatusChange(ctx, documentID, string(status), note)
	})

	// 6b-2. Load deployment license (mode on-premise) SEBELUM registrasi tenant
	// modules, agar employee module dapat menerima quota checker max_employees
	// dari file .lic. Pada mode saas, licenseLister memakai company_modules DB.
	var onPremiseLic *onpremise.License
	var licenseLister middleware.CompanyModuleLister
	if cfg.License.DeploymentMode == onpremise.ModeOnPremise {
		pubKey, err := os.ReadFile(cfg.License.PublicKeyFile)
		if err != nil {
			l.Fatal("On-premise public key not found", zap.String("path", cfg.License.PublicKeyFile), zap.Error(err))
		}
		lic, err := onpremise.ReadLicenseFile(cfg.License.LicenseFile, string(pubKey))
		if err != nil {
			l.Fatal("On-premise license invalid", zap.String("path", cfg.License.LicenseFile), zap.Error(err))
		}
		onPremiseLic = lic
		licenseLister = onPremiseLister{lic: lic}
		l.Info("On-premise license loaded",
			zap.String("company", lic.CompanyName),
			zap.Time("expires_at", lic.ExpiresAt),
			zap.Int("max_employees", lic.MaxEmployees),
		)
	} else {
		licenseLister = companyModuleListerAdapter{svc: modulemgmtSvc}
	}

	// 6c. Register tenant modules (ordered by priority & dependency)
	tenantModules = append(tenantModules,
		module.ModuleRegistration{
			Module:   organization.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 1,
		},
		module.ModuleRegistration{
			Module:   newEmployeeModule(dbManager, l, onPremiseLic),
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
			Module:   employeemovement.NewModuleWithService(l, employeeMovementSvc),
			TargetDB: module.TargetTenant,
			Priority: 5,
		},
		module.ModuleRegistration{
			Module:   attendance.NewModuleWithService(l, attendanceSvc),
			TargetDB: module.TargetTenant,
			Priority: 6,
		},
		module.ModuleRegistration{
			Module:   approval.NewModuleWithService(dbManager, l, approvalSvc),
			TargetDB: module.TargetTenant,
			Priority: 7,
		},
		module.ModuleRegistration{
			Module:   payroll.NewModuleWithService(l, payrollSvc),
			TargetDB: module.TargetTenant,
			Priority: 8,
		},
		module.ModuleRegistration{
			Module:   leave.NewModuleWithService(l, leaveSvc),
			TargetDB: module.TargetTenant,
			Priority: 9,
		},
		module.ModuleRegistration{
			Module:   performance.NewModuleWithServices(dbManager, l, performanceSvc, okrSvc),
			TargetDB: module.TargetTenant,
			Priority: 10,
		},
		module.ModuleRegistration{
			Module:   recruitment.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 11,
		},
		module.ModuleRegistration{
			Module:   reimbursement.NewModuleWithService(l, reimbursementSvc),
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
		module.ModuleRegistration{
			Module:   rbac.NewModule(dbManager, l),
			TargetDB: module.TargetTenant,
			Priority: 17,
		},
		module.ModuleRegistration{
			Module:   useraccount.NewModule(dbManager, authManager, mailerSvc, l),
			TargetDB: module.TargetTenant,
			Priority: 18,
		},
		module.ModuleRegistration{
			Module:   notification.NewModuleWithService(l, notificationSvc),
			TargetDB: module.TargetTenant,
			Priority: 19,
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
	// PlatformLicenseMiddleware: guard lisensi modul tenant (validasi via Redis cache).
	// Diletakkan setelah TenantRequired agar company_id tersedia, dan sebelum RBAC agar
	// error lisensi tampil lebih dulu.
	//
	// licenseLister sudah di-inisialisasi sebelum registrasi tenant modules
	// (bagian 6b-2) — sumber lisensi bergantung pada DeploymentMode:
	//   - saas (default): company_modules dari platform DB (modulemgmt service)
	//   - on_premise: file .lic RSA (expires_at, allowed_modules, max_employees)
	licenseMW := middleware.LicenseMiddleware(cacheManager, licenseLister, l)

	// TenantResolver middleware: auto-detect tenant (company) dari Host header /
	// X-Tenant-ID saat mengakses /api/v1/tenant/** (mode SaaS — tiap company punya
	// URL/subdomain sendiri), sehingga FE tidak perlu mengirim company_slug/company_id
	// manual. Company module mengekspos ResolveByHost untuk memenuhi interface
	// middleware.HostCompanyResolver. JWT claims (company_id) tetap menang.
	var tenantResolver gin.HandlerFunc
	for _, reg := range platformModules {
		if mod, ok := reg.Module.(interface{ ResolveByHost(string) (string, error) }); ok {
			tenantResolver = middleware.TenantResolver(mod, l.Named("tenant-resolver"))
			l.Info("TenantResolver middleware wired (auto-detect company from Host)")
			break
		}
	}

	r := router.Setup(
		router.Config{Mode: cfg.Server.Mode},
		middleware.AuthJWT(authManager, l),
		tenantResolver,
		middleware.TenantRequired(),
		licenseMW,
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
	if tenantResolver != nil {
		tenantPkgGroup.Use(tenantResolver)
	}
	tenantPkgGroup.Use(middleware.TenantRequired())
	{
		// Generic file upload (lampiran modul tenant — konsumen pertama: isian
		// aktual lembur, §32b docs/module-attendance-plan.md). Auth tenant saja;
		// file diserve publik via router r.Static("/uploads", upload_dir).
		uploadHandler := upload.NewHandler(cfg.Storage.UploadDir)
		tenantPkgGroup.POST("/uploads", uploadHandler.Upload)

		// Self-service company detail endpoint (dipakai halaman Detail Perusahaan FE tenant).
		// Terdaftar di tenantPkgGroup (tanpa RBAC/license guard) karena data company
		// milik user sendiri — pola sama dengan /company-modules & /packages.
		for _, reg := range platformModules {
			if mod, ok := reg.Module.(interface{ PublicHandler() *company.Handler }); ok {
				tenantPkgGroup.GET("/companies/me", mod.PublicHandler().GetCurrent)
				tenantPkgGroup.PUT("/companies/me", mod.PublicHandler().UpdateCurrent)
				l.Info("Registered tenant self-service company endpoints (GET/PUT /companies/me)")
				break
			}
		}

		// Endpoint untuk mendapatkan daftar module aktif company (digunakan frontend tenant untuk menu filtering)
		tenantPkgGroup.GET("/company-modules", func(c *gin.Context) {
			companyID := middleware.GetCompanyID(c)
			if companyID == "" {
				httputil.Unauthorized(c, "Company ID is required")
				return
			}
			modules, err := modulemgmtSvc.ListCompanyModules(companyID)
			if err != nil {
				httputil.InternalError(c, "Failed to fetch company modules")
				return
			}
			// Hanya kembalikan slug untuk module yang enabled
			var activeModuleSlugs []string
			for _, m := range modules {
				if m.Enabled {
					activeModuleSlugs = append(activeModuleSlugs, m.ModuleSlug)
				}
			}
			httputil.SuccessJSON(c, gin.H{
				"company_id": companyID,
				"modules":    activeModuleSlugs,
			})
		})
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
			// License service akan auto-activate modules dari package via PackageModuleManager
			licenseID, licenseKey, planType, err := licenseCreator.CreateFromPackage(companyID, packageID)
			if err != nil {
				httputil.InternalError(c, err.Error())
				return
			}

			// Invalidate license cache agar PlatformLicenseMiddleware langsung
			// mengenali modul baru setelah subscribe.
			if err := cacheManager.Invalidate(c.Request.Context(), middleware.LicenseCacheKey(companyID)); err != nil {
				l.Warn("Failed to invalidate license cache after subscribe",
					zap.String("company_id", companyID),
					zap.Error(err),
				)
			}

			// Modules sudah auto-activated oleh license service via PackageModuleManager.
			// Ambil daftar nama module dari package untuk response.
			var activatedModules []string
			for _, m := range pkg.Modules {
				activatedModules = append(activatedModules, m.ModuleName)
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
			_, err := pkgSvc.GetPackage(packageID)
			if err != nil {
				httputil.NotFound(c, "Package not found")
				return
			}

			// Deactivate all modules from the package via PackageModuleManager
			if _, err := pkgModuleMgr.DeactivatePackageModules(packageID, companyID); err != nil {
				l.Warn("Failed to deactivate package modules during unsubscribe",
					zap.String("package_id", packageID),
					zap.String("company_id", companyID),
					zap.Error(err),
				)
			}

			// Invalidate license cache agar modul nonaktif langsung diblokir middleware.
			if err := cacheManager.Invalidate(c.Request.Context(), middleware.LicenseCacheKey(companyID)); err != nil {
				l.Warn("Failed to invalidate license cache after unsubscribe",
					zap.String("company_id", companyID),
					zap.Error(err),
				)
			}

			// Suspend the active license associated with this package
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

	// Register public company resolve endpoint (no auth) —
	// menentukan company dari hostname/subdomain URL aplikasi (mode SaaS).
	// Dipakai tenant FE sebelum login untuk prefill company.
	for _, reg := range platformModules {
		if mod, ok := reg.Module.(interface{ PublicHandler() *company.Handler }); ok {
			r.GET("/api/v1/public/companies/resolve", mod.PublicHandler().ResolveByHost)
			l.Info("Registered public company resolve endpoint")
			break
		}
	}

	// Register public employee account setup endpoint (no auth) —
	// tujuan link email untuk set password akun login.
	// Path: POST /api/v1/public/account/setup-password
	//
	// Register juga tenant auth publik (login & refresh) di /api/v1/tenant/auth —
	// user tenant (employee) login via company_slug, di luar middleware auth.
	for _, reg := range tenantModules {
		if mod, ok := reg.Module.(interface{ PublicHandler() *useraccount.Handler }); ok {
			publicGroup := r.Group("/api/v1/public/account")
			publicGroup.POST("/setup-password", mod.PublicHandler().SetPassword)
			l.Info("Registered public account setup endpoint")

			tenantAuth := r.Group("/api/v1/tenant/auth")
			// TenantResolver juga dipasang di login agar FE bisa login tanpa
			// mengirim company_slug/company_id — tenant di-resolve dari Host.
			if tenantResolver != nil {
				tenantAuth.Use(tenantResolver)
			}
			tenantAuth.POST("/login", mod.PublicHandler().Login)
			tenantAuth.POST("/refresh", mod.PublicHandler().Refresh)
			l.Info("Registered public tenant auth endpoints (login/refresh)")
			break
		}
	}

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
