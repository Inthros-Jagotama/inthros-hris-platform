package useraccount

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/auth"
	"github.com/inthros/hris-platform/internal/pkg/authctx"
)

// DefaultRoleName adalah role bawaan tenant yang di-assign ke akun employee.
const DefaultRoleName = "Employee"

// SetupTokenTTL adalah masa berlaku link set-password (48 jam).
const SetupTokenTTL = 48 * time.Hour

// Mailer adalah antarmuka pengiriman email (diimplementasikan oleh pkg/mailer).
type Mailer interface {
	Send(to, subject, bodyHTML string) error
	SetupLink(token, companyID string) string
}

// EmployeeRef adalah subset minimal tabel employees (untuk verifikasi + user_id).
type EmployeeRef struct {
	ID     uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name   string    `gorm:"type:varchar(255);not null" json:"name"`
	Email  *string   `gorm:"type:varchar(255)" json:"email,omitempty"`
	UserID *uuid.UUID `gorm:"column:user_id;type:char(36)" json:"user_id,omitempty"`
}

func (EmployeeRef) TableName() string { return "employees" }

type Service struct {
	repo        *Repository
	authManager *auth.Manager
	mailer      Mailer
	logger      *zap.Logger
}

func NewService(repo *Repository, authManager *auth.Manager, mailer Mailer, logger *zap.Logger) *Service {
	return &Service{repo: repo, authManager: authManager, mailer: mailer, logger: logger}
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func randomToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// hashToken menyimpan token sebagai SHA-256 di database (jangan plaintext).
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// randomPassword menghasilkan password placeholder acak (user akan set sendiri
// melalui link email — password placeholder tidak pernah digunakan untuk login).
func randomPassword() (string, error) {
	tok, err := randomToken()
	if err != nil {
		return "", err
	}
	return tok[:16] + "Aa1!", nil
}

// ── Account creation (authenticated) ─────────────────────────────────────────

// CreateAccount membuat user tenant + assign role Employee + kirim email setup.
// email diambil dari request; fallback ke email personal employee jika request kosong.
func (s *Service) CreateAccount(ctx context.Context, employeeID, email string) (*AccountResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	// Verifikasi employee ada.
	var emp EmployeeRef
	{
		db, err := s.repo.getDB(ctx)
		if err != nil {
			return nil, err
		}
		if err := db.First(&emp, "id = ?", empUID).Error; err != nil {
			return nil, fmt.Errorf("employee not found: %w", err)
		}
	}

	// Akun hanya boleh dibuat sekali per employee.
	if _, err := s.repo.FindAccountByEmployeeID(ctx, empUID); err == nil {
		return nil, errors.New("account already exists for this employee")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) && !isNotFound(err) {
		return nil, err
	}

	if email == "" {
		if emp.Email == nil || *emp.Email == "" {
			return nil, errors.New("email is required (fill employee personal email or provide email)")
		}
		email = *emp.Email
	}

	// Deteksi duplikat email di tabel users.
	if _, err := s.repo.FindUserByEmail(ctx, email); err == nil {
		return nil, fmt.Errorf("email %q is already registered as a user", email)
	} else if !isNotFound(err) {
		return nil, err
	}

	// 0. Ambil company_id dari context (di-set oleh TenantRequired middleware).
	// Dipakai agar route publik set-password bisa resolve tenant DB tanpa JWT.
	companyID, _ := ctx.Value("company_id").(string)
	var companyUUID uuid.UUID
	if companyID != "" {
		companyUUID, _ = uuid.Parse(companyID)
	}

	// 1. Buat user tenant dengan password placeholder.
	placeholder, err := randomPassword()
	if err != nil {
		return nil, err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(placeholder), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	userID := uuid.NewString()
	user := &TenantUser{
		ID:       userID,
		Name:     emp.Name,
		Email:    email,
		Password: string(hashed),
		IsActive: 1,
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 2. Assign role default (Employee).
	roleID, err := s.repo.FindRoleByName(ctx, DefaultRoleName)
	if err != nil {
		return nil, err
	}
	if err := s.repo.AssignRoleToUser(ctx, userID, roleID); err != nil {
		return nil, fmt.Errorf("failed to assign default role: %w", err)
	}

	// 3. Buat record employee_accounts + setup token.
	token, err := randomToken()
	if err != nil {
		return nil, err
	}
	expiry := time.Now().Add(SetupTokenTTL)
	hashedToken := hashToken(token)
	createdBy := authctx.GetUserID(ctx)
	acc := &EmployeeAccount{
		CompanyID:        companyUUID,
		EmployeeID:       empUID,
		UserID:           uuid.MustParse(userID),
		Email:            email,
		SetupToken:       &hashedToken,
		SetupTokenExpiry: &expiry,
		CreatedBy:        createdBy,
		UpdatedBy:        createdBy,
	}
	if err := s.repo.CreateAccount(ctx, acc); err != nil {
		return nil, fmt.Errorf("failed to create employee account: %w", err)
	}

	// 4. Simpan user_id ke employee (jika kolom tersedia).
	s.tryLinkEmployeeUser(ctx, empUID, userID)

	// 5. Kirim email setup (link menyertakan company_id untuk route publik).
	link := s.mailer.SetupLink(token, acc.CompanyID.String())
	if err := s.mailer.Send(email, subjectSetupAccount, bodySetupAccount(emp.Name, link)); err != nil {
		s.logger.Warn("setup email send failed", zap.Error(err), zap.String("email", email))
	}

	s.logger.Info("Employee account created",
		zap.String("employee_id", empUID.String()),
		zap.String("user_id", userID),
		zap.String("email", email),
	)

	return &AccountResponse{
		ID:               acc.ID.String(),
		EmployeeID:       empUID.String(),
		CompanyID:        acc.CompanyID.String(),
		UserID:           userID,
		Email:            email,
		RoleName:         DefaultRoleName,
		SetupToken:       token,
		SetupTokenExpiry: &expiry,
		PasswordSet:      false,
		CreatedAt:        acc.CreatedAt,
		UpdatedAt:        acc.UpdatedAt,
	}, nil
}

// ResendSetupEmail membuat token baru + kirim ulang email setup.
func (s *Service) ResendSetupEmail(ctx context.Context, employeeID string) (*AccountResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	acc, err := s.repo.FindAccountByEmployeeID(ctx, empUID)
	if err != nil {
		return nil, err
	}

	token, err := randomToken()
	if err != nil {
		return nil, err
	}
	expiry := time.Now().Add(SetupTokenTTL)
	hashedToken := hashToken(token)
	acc.SetupToken = &hashedToken
	acc.SetupTokenExpiry = &expiry
	acc.UpdatedBy = authctx.GetUserID(ctx)
	if err := s.repo.UpdateAccount(ctx, acc); err != nil {
		return nil, err
	}

	link := s.mailer.SetupLink(token, acc.CompanyID.String())
	if err := s.mailer.Send(acc.Email, subjectSetupAccount, bodySetupAccount("", link)); err != nil {
		s.logger.Warn("resend setup email failed", zap.Error(err), zap.String("email", acc.Email))
	}

	return &AccountResponse{
		ID:               acc.ID.String(),
		EmployeeID:       acc.EmployeeID.String(),
		UserID:           acc.UserID.String(),
		Email:            acc.Email,
		RoleName:         DefaultRoleName,
		SetupToken:       token,
		SetupTokenExpiry: &expiry,
		PasswordSet:      false,
		CreatedAt:        acc.CreatedAt,
		UpdatedAt:        acc.UpdatedAt,
	}, nil
}

// GetAccountStatus mengembalikan status akun employee tanpa token.
func (s *Service) GetAccountStatus(ctx context.Context, employeeID string) (*AccountResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}
	acc, err := s.repo.FindAccountByEmployeeID(ctx, empUID)
	if err != nil {
		return nil, err
	}
	passwordSet := acc.SetupToken == nil || acc.SetupTokenExpiry == nil
	return &AccountResponse{
		ID:          acc.ID.String(),
		EmployeeID:  acc.EmployeeID.String(),
		UserID:      acc.UserID.String(),
		Email:       acc.Email,
		RoleName:    DefaultRoleName,
		PasswordSet: passwordSet,
		CreatedAt:   acc.CreatedAt,
		UpdatedAt:   acc.UpdatedAt,
	}, nil
}

// GetMyAccount resolves the employee_accounts row for the currently
// authenticated user (via authctx.GetUserID), so FE features can look up
// "my employee_id" without being handed a client-supplied employee_id.
func (s *Service) GetMyAccount(ctx context.Context) (*AccountResponse, error) {
	userID := authctx.GetUserID(ctx)
	if userID == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	acc, err := s.repo.FindAccountByUserID(ctx, *userID)
	if err != nil {
		return nil, err
	}
	passwordSet := acc.SetupToken == nil || acc.SetupTokenExpiry == nil
	return &AccountResponse{
		ID:          acc.ID.String(),
		EmployeeID:  acc.EmployeeID.String(),
		UserID:      acc.UserID.String(),
		Email:       acc.Email,
		RoleName:    DefaultRoleName,
		PasswordSet: passwordSet,
		CreatedAt:   acc.CreatedAt,
		UpdatedAt:   acc.UpdatedAt,
	}, nil
}

// ── Password setup (public, via link email) ──────────────────────────────────

// SetPassword memverifikasi setup token lalu mengeset password baru.
// ctx wajib memuat company_id (handler publik menginject dari employee_accounts).
func (s *Service) SetPassword(ctx context.Context, token, newPassword string) (*AccountResponse, error) {
	if token == "" {
		return nil, errors.New("token is required")
	}
	hashedToken := hashToken(token)
	acc, err := s.repo.FindAccountBySetupToken(ctx, hashedToken)
	if err != nil {
		return nil, err
	}
	if acc.SetupTokenExpiry != nil && time.Now().After(*acc.SetupTokenExpiry) {
		return nil, errors.New("setup link has expired, please request a new one")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.FindUserByID(ctx, acc.UserID.String())
	if err != nil {
		return nil, err
	}
	user.Password = string(hashed)

	// Token sekali pakai — langsung dihapus setelah sukses.
	acc.SetupToken = nil
	acc.SetupTokenExpiry = nil
	acc.UpdatedBy = authctx.GetUserID(ctx)
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateAccount(ctx, acc); err != nil {
		return nil, err
	}

	return &AccountResponse{
		ID:          acc.ID.String(),
		EmployeeID:  acc.EmployeeID.String(),
		UserID:      acc.UserID.String(),
		Email:       acc.Email,
		RoleName:    DefaultRoleName,
		PasswordSet: true,
		CreatedAt:   acc.CreatedAt,
		UpdatedAt:   acc.UpdatedAt,
	}, nil
}

// ── Tenant Auth: login & refresh (public) ────────────────────────────────────

// Login memvalidasi kredensial user tenant (employee) terhadap tenant database
// yang di-resolve dari company slug, lalu mengembalikan JWT dengan company_id,
// role, dan permissions. Bilingual-safe: pesan error generik.
//
// Fallback: bila user tenant tidak ditemukan/password salah, coba platform user
// (company_admin) yang terikat company — agar admin yang sebelumnya login via
// /api/v1/platform/login tetap bisa masuk lewat FE tenant yang sama.
func (s *Service) Login(ctx context.Context, req TenantLoginRequest) (*TenantLoginResponse, error) {
	// 1. Resolve company dari slug ATAU company_id (platform DB).
	//    company_id dipakai development via env (VITE_COMPANY_ID) atau
	//    setelah resolve hostname (endpoint /public/companies/resolve).
	var companyID, status string
	var err error
	if req.CompanyID != "" {
		companyID = req.CompanyID
		status, err = s.repo.FindCompanyByID(ctx, req.CompanyID)
	} else {
		companyID, status, err = s.repo.FindCompanyBySlug(ctx, req.CompanySlug)
	}
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}
	if status != "active" {
		return nil, fmt.Errorf("company is not active")
	}

	// 2. Inject company_id ke context agar tenant DB bisa di-resolve.
	tctx := context.WithValue(ctx, "company_id", companyID)

	// 3. Coba login sebagai user tenant (employee).
	//    Error DB asli (bukan not-found) harus diteruskan, jangan disamarkan
	//    menjadi "invalid email or password" oleh fallback platform.
	user, userErr := s.repo.FindUserByEmail(tctx, req.Email)
	if userErr != nil && !isNotFound(userErr) {
		return nil, userErr
	}
	if userErr == nil && user.IsActive == 1 &&
		bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) == nil {
		return s.loginTenantUser(tctx, user, companyID)
	}

	// 4. Fallback: login sebagai platform user (company_admin) terikat company.
	return s.loginPlatformUser(ctx, companyID, req)
}

// loginTenantUser membuat token untuk user tenant (employee).
func (s *Service) loginTenantUser(tctx context.Context, user *TenantUser, companyID string) (*TenantLoginResponse, error) {
	roles, err := s.repo.FindRolesForUser(tctx, user.ID)
	if err != nil {
		return nil, err
	}
	role := "Employee"
	if len(roles) > 0 {
		role = roles[0]
	}
	permissions, err := s.repo.FindPermissionsForUser(tctx, user.ID)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.authManager.GenerateTokenPair(
		user.ID, companyID, user.Email, user.Name, role, permissions,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}
	_ = s.repo.UpdateLastLogin(tctx, user.ID)

	s.logger.Info("Tenant user logged in",
		zap.String("user_id", user.ID),
		zap.String("email", user.Email),
		zap.String("company_id", companyID),
		zap.String("role", role),
	)

	// Nama company dari platform DB — dipakai FE (sidebar) tanpa fetch tambahan.
	companyName, _ := s.repo.FindCompanyNameByID(tctx, companyID)

	return &TenantLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User: TenantUserResponse{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			Role:        role,
			Permissions: permissions,
			CompanyID:   companyID,
			CompanyName: companyName,
		},
	}, nil
}

// loginPlatformUser membuat token untuk platform user (company_admin) yang
// terikat company tertentu. Dipakai sebagai fallback tenant login.
func (s *Service) loginPlatformUser(ctx context.Context, companyID string, req TenantLoginRequest) (*TenantLoginResponse, error) {
	puser, err := s.repo.FindPlatformUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}
	if puser.CompanyID == nil || *puser.CompanyID != companyID {
		return nil, fmt.Errorf("invalid email or password")
	}
	if !puser.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}
	if bcrypt.CompareHashAndPassword([]byte(puser.PasswordHash), []byte(req.Password)) != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	permissions, err := s.repo.FindPlatformPermissionsForRole(ctx, puser.Role)
	if err != nil {
		return nil, err
	}
	accessToken, refreshToken, err := s.authManager.GenerateTokenPair(
		puser.ID, companyID, puser.Email, puser.Name, puser.Role, permissions,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	s.logger.Info("Platform user logged in via tenant endpoint",
		zap.String("user_id", puser.ID),
		zap.String("email", puser.Email),
		zap.String("company_id", companyID),
		zap.String("role", puser.Role),
	)

	// Nama company dari platform DB — dipakai FE (sidebar) tanpa fetch tambahan.
	companyName, _ := s.repo.FindCompanyNameByID(ctx, companyID)

	return &TenantLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User: TenantUserResponse{
			ID:          puser.ID,
			Name:        puser.Name,
			Email:       puser.Email,
			Role:        puser.Role,
			Permissions: permissions,
			CompanyID:   companyID,
			CompanyName: companyName,
		},
	}, nil
}

// Refresh memvalidasi refresh token dan mengembalikan access token baru.
func (s *Service) Refresh(refreshToken string) (*TenantRefreshResponse, error) {
	claims, err := s.authManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}
	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("token is not a refresh token")
	}

	newAccess, err := s.authManager.RefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return &TenantRefreshResponse{
		AccessToken: newAccess,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}, nil
}

// tryLinkEmployeeUser mengisi employees.user_id (jika kolom tersedia).
func (s *Service) tryLinkEmployeeUser(ctx context.Context, employeeID uuid.UUID, userID string) {
	db, err := s.repo.getDB(ctx)
	if err != nil {
		return
	}
	uid := uuid.MustParse(userID)
	// Update hanya jika kolom user_id ada — gagal diam-diam jika tidak.
	_ = db.Model(&EmployeeRef{}).Where("id = ?", employeeID).Update("user_id", uid).Error
}

// isNotFound mengecek apakah error adalah GORM record not found (termasuk wrapper).
func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, gorm.ErrRecordNotFound)
}
