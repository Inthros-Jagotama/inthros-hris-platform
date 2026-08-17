package employee

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
	"github.com/inthros/hris-platform/internal/pkg/crypto"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

type Service struct {
	repo   *Repository
	logger *zap.Logger
	quota  EmployeeQuotaChecker // nil = unlimited (mode saas)
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// SetQuotaChecker menginjeksi batas kuota employee (mode on-premise).
// Dipanggil dari main.go saat lisensi on-premise dimuat.
func (s *Service) SetQuotaChecker(qc EmployeeQuotaChecker) {
	s.quota = qc
}

// checkQuota menolak pembuatan employee jika jumlah saat ini sudah mencapai
// batas maksimal lisensi on-premise. Tanpa checker (mode saas) → selalu lolos.
//
// Catatan: check-then-insert tidak atomik — dua request konkuren bisa lolos
// bersamaan dan melebihi batas sebanyak 1. Diterima untuk MVP; untuk
// enforcement ketat perlu transaksi/kunci di level DB.
func (s *Service) checkQuota(ctx context.Context) error {
	if s.quota == nil {
		return nil
	}
	max := s.quota.MaxEmployees()
	if max <= 0 {
		return nil // unlimited
	}

	count, err := s.repo.CountEmployees(ctx)
	if err != nil {
		return fmt.Errorf("failed to count employees for quota check: %w", err)
	}
	if count >= int64(max) {
		return ErrQuotaExceeded
	}
	return nil
}

// encryptIfEnabled meng-enkripsi *value in-place jika toggle enkripsi
// untuk fieldKey aktif. Nilai nil/kosong tidak diproses.
func (s *Service) encryptIfEnabled(ctx context.Context, fieldKey string, value *string) error {
	if value == nil || *value == "" {
		return nil
	}
	if crypto.LooksEncrypted(*value) {
		return nil
	}
	enabled, err := s.IsFieldEncryptionEnabled(ctx, fieldKey)
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}
	encrypted, err := crypto.EncryptString(*value)
	if err != nil {
		return fmt.Errorf("encrypt %s: %w", fieldKey, err)
	}
	*value = encrypted
	return nil
}

// =========================================================================
// Employee CRUD
// =========================================================================

func (s *Service) Create(ctx context.Context, req CreateEmployeeRequest) (*EmployeeResponse, error) {
	// Enforce kuota max_employees (on-premise) sebelum membuat record baru.
	if err := s.checkQuota(ctx); err != nil {
		return nil, err
	}

	emp := &Employee{
		EmployeeID: req.EmployeeID,
		Name:       req.Name,
		Status:     "active",
	}
	emp.CreatedBy = authctx.GetUserID(ctx)
	emp.UpdatedBy = emp.CreatedBy

	if req.NIK != nil {
		emp.NIK = req.NIK
	}
	if req.FamilyID != nil {
		emp.FamilyID = req.FamilyID
	}
	if req.MotherName != nil {
		emp.MotherName = req.MotherName
	}
	if req.Gender != nil {
		emp.Gender = req.Gender
	}
	if req.NationalityType != nil {
		emp.NationalityType = req.NationalityType
	}
	if req.NationalityID != nil {
		emp.NationalityID = req.NationalityID
	}
	if req.Passport != nil {
		emp.Passport = req.Passport
	}
	if req.POB != nil {
		emp.POB = req.POB
	}
	if req.DOB != nil {
		emp.DOB = req.DOB
	}
	if req.PhoneNumber != nil {
		emp.PhoneNumber = req.PhoneNumber
	}
	if req.Email != nil {
		emp.Email = req.Email
	}
	if req.LinkedIn != nil {
		emp.LinkedIn = req.LinkedIn
	}
	if req.Instagram != nil {
		emp.Instagram = req.Instagram
	}

	if req.ReligionID != nil && *req.ReligionID != "" {
		id, err := uuid.Parse(*req.ReligionID)
		if err != nil {
			return nil, fmt.Errorf("invalid religion_id: %w", err)
		}
		emp.ReligionID = &id
	}
	if req.MaritalStatusID != nil && *req.MaritalStatusID != "" {
		id, err := uuid.Parse(*req.MaritalStatusID)
		if err != nil {
			return nil, fmt.Errorf("invalid marital_status_id: %w", err)
		}
		emp.MaritalStatusID = &id
	}
	// G-4: referensi balik ke aplikasi recruitment asal (offer eksternal diterima).
	if req.RecruitedFromApplicationID != nil && *req.RecruitedFromApplicationID != "" {
		id, err := uuid.Parse(*req.RecruitedFromApplicationID)
		if err != nil {
			return nil, fmt.Errorf("invalid recruited_from_application_id: %w", err)
		}
		emp.RecruitedFromApplicationID = &id
	}

	if err := s.encryptIfEnabled(ctx, "employee.nik", emp.NIK); err != nil {
		return nil, err
	}
	if err := s.encryptIfEnabled(ctx, "employee.passport", emp.Passport); err != nil {
		return nil, err
	}
	if err := s.encryptIfEnabled(ctx, "employee.phone_number", emp.PhoneNumber); err != nil {
		return nil, err
	}
	if err := s.encryptIfEnabled(ctx, "employee.email", emp.Email); err != nil {
		return nil, err
	}

	if err := s.repo.CreateEmployee(ctx, emp); err != nil {
		return nil, err
	}

	s.logger.Info("Employee created",
		zap.String("id", emp.ID.String()),
		zap.String("name", emp.Name),
		zap.String("employee_id", emp.EmployeeID),
	)

	response := emp.ToResponse()
	s.fillRefNames(ctx, emp, &response)
	maskEmployeeResponse(ctx, &response)
	return &response, nil
}

// GetEmploymentStatusStats mengembalikan jumlah karyawan per status kepegawaian
// (employment berjalan) untuk pie chart dashboard Employment.
func (s *Service) GetEmploymentStatusStats(ctx context.Context) (*EmploymentStatusStatsResponse, error) {
	groups, unclassified, err := s.repo.CountByEmploymentStatus(ctx)
	if err != nil {
		return nil, err
	}
	return &EmploymentStatusStatsResponse{Groups: groups, Unclassified: unclassified}, nil
}

// GetGenderStats mengembalikan jumlah karyawan per jenis kelamin untuk pie
// chart dashboard Employment.
func (s *Service) GetGenderStats(ctx context.Context) (*GenderStatsResponse, error) {
	male, female, other, err := s.repo.CountByGender(ctx)
	if err != nil {
		return nil, err
	}
	return &GenderStatsResponse{Male: male, Female: female, Other: other}, nil
}

// fillRefNames mengisi nama referensi (agama, status perkawinan, kewarganegaraan)
// pada response — di-resolve langsung dari tenant DB supaya halaman profile
// menampilkan nama (bukan ID mentah) dan tidak bergantung pada permission
// viewer terhadap endpoint /settings/*.
func (s *Service) fillRefNames(ctx context.Context, emp *Employee, response *EmployeeResponse) {
	religionName, maritalStatusName, nationalityName := s.repo.ResolveEmployeeRefNames(ctx, emp.ReligionID, emp.MaritalStatusID, emp.NationalityID)
	response.ReligionName = religionName
	response.MaritalStatusName = maritalStatusName
	response.NationalityName = nationalityName
}

func (s *Service) GetByID(ctx context.Context, id string) (*EmployeeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	emp, err := s.repo.FindEmployeeByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	response := emp.ToResponse()
	s.fillRefNames(ctx, emp, &response)
	maskEmployeeResponse(ctx, &response)
	return &response, nil
}

func (s *Service) List(ctx context.Context, page, perPage int, search, status, organizationID string) (*ListResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}

	employees, total, err := s.repo.FindAllEmployees(ctx, page, perPage, search, status, organizationID)
	if err != nil {
		return nil, err
	}

	var responses []EmployeeResponse
	for _, e := range employees {
		r := e.ToResponse()
		maskEmployeeResponse(ctx, &r)
		responses = append(responses, r)
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &ListResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateEmployeeRequest) (*EmployeeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	emp, err := s.repo.FindEmployeeByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	emp.UpdatedBy = authctx.GetUserID(ctx)

	if req.Name != nil {
		emp.Name = *req.Name
	}
	// Field sensitif: lewati penulisan jika request hanya memantulkan nilai
	// ter-mask yang diterima caller saat GET (lihat sensitive_field_update.go).
	applyPtrIfNotMasked(req.NIK, &emp.NIK)
	if req.FamilyID != nil {
		emp.FamilyID = req.FamilyID
	}
	if req.MotherName != nil {
		emp.MotherName = req.MotherName
	}
	if req.Gender != nil {
		emp.Gender = req.Gender
	}
	if req.NationalityType != nil {
		emp.NationalityType = req.NationalityType
	}
	if req.NationalityID != nil {
		emp.NationalityID = req.NationalityID
	}
	applyPtrIfNotMasked(req.Passport, &emp.Passport)
	if req.POB != nil {
		emp.POB = req.POB
	}
	if req.DOB != nil {
		emp.DOB = req.DOB
	}
	applyPtrIfNotMasked(req.PhoneNumber, &emp.PhoneNumber)
	applyPtrIfNotMasked(req.Email, &emp.Email)
	if req.LinkedIn != nil {
		emp.LinkedIn = req.LinkedIn
	}
	if req.Instagram != nil {
		emp.Instagram = req.Instagram
	}
	if req.ReligionID != nil && *req.ReligionID != "" {
		id, err := uuid.Parse(*req.ReligionID)
		if err != nil {
			return nil, fmt.Errorf("invalid religion_id: %w", err)
		}
		emp.ReligionID = &id
	}
	if req.MaritalStatusID != nil && *req.MaritalStatusID != "" {
		id, err := uuid.Parse(*req.MaritalStatusID)
		if err != nil {
			return nil, fmt.Errorf("invalid marital_status_id: %w", err)
		}
		emp.MaritalStatusID = &id
	}
	if req.Status != nil {
		emp.Status = *req.Status
	}

	if err := s.encryptIfEnabled(ctx, "employee.nik", emp.NIK); err != nil {
		return nil, err
	}
	if err := s.encryptIfEnabled(ctx, "employee.passport", emp.Passport); err != nil {
		return nil, err
	}
	if err := s.encryptIfEnabled(ctx, "employee.phone_number", emp.PhoneNumber); err != nil {
		return nil, err
	}
	if err := s.encryptIfEnabled(ctx, "employee.email", emp.Email); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateEmployee(ctx, emp); err != nil {
		return nil, err
	}

	response := emp.ToResponse()
	s.fillRefNames(ctx, emp, &response)
	maskEmployeeResponse(ctx, &response)
	return &response, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid employee id: %w", err)
	}
	return s.repo.DeleteEmployee(ctx, uid)
}

func (s *Service) UpdatePhoto(ctx context.Context, id, photoURL string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid employee id: %w", err)
	}

	emp, err := s.repo.FindEmployeeByID(ctx, uid)
	if err != nil {
		return err
	}
	emp.ProfilePicture = &photoURL
	emp.UpdatedBy = authctx.GetUserID(ctx)

	return s.repo.UpdateEmployee(ctx, emp)
}

func (s *Service) DeletePhoto(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid employee id: %w", err)
	}

	emp, err := s.repo.FindEmployeeByID(ctx, uid)
	if err != nil {
		return err
	}
	emp.ProfilePicture = nil
	emp.UpdatedBy = authctx.GetUserID(ctx)

	return s.repo.UpdateEmployee(ctx, emp)
}

// =========================================================================
// Sub-module CRUD: Addresses
// =========================================================================

func (s *Service) CreateAddress(ctx context.Context, employeeID string, req CreateAddressRequest) (*AddressResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	addr := &EmployeeAddress{
		EmployeeID: &empUID,
		Type:       &req.Type,
		Address:    &req.Address,
	}
	if req.ProvinceID != nil {
		addr.ProvinceID = req.ProvinceID
	}
	if req.RegencyID != nil {
		addr.RegencyID = req.RegencyID
	}
	if req.DistrictID != nil {
		addr.DistrictID = req.DistrictID
	}
	if req.VillageID != nil {
		addr.VillageID = req.VillageID
	}
	if req.PostalCode != nil {
		addr.PostalCode = req.PostalCode
	}
	addr.CreatedBy = authctx.GetUserID(ctx)
	addr.UpdatedBy = addr.CreatedBy

	if err := s.repo.CreateAddress(ctx, addr); err != nil {
		return nil, err
	}

	response := toAddressResponse(addr)
	return &response, nil
}

func (s *Service) UpdateAddress(ctx context.Context, employeeID, addressID string, req UpdateAddressRequest) (*AddressResponse, error) {
	addrUID, err := uuid.Parse(addressID)
	if err != nil {
		return nil, fmt.Errorf("invalid address id: %w", err)
	}

	addr, err := s.repo.FindAddressByID(ctx, addrUID)
	if err != nil {
		return nil, err
	}
	addr.UpdatedBy = authctx.GetUserID(ctx)

	if req.Type != nil {
		addr.Type = req.Type
	}
	if req.Address != nil {
		addr.Address = req.Address
	}
	if req.ProvinceID != nil {
		addr.ProvinceID = req.ProvinceID
	}
	if req.RegencyID != nil {
		addr.RegencyID = req.RegencyID
	}
	if req.DistrictID != nil {
		addr.DistrictID = req.DistrictID
	}
	if req.VillageID != nil {
		addr.VillageID = req.VillageID
	}
	if req.PostalCode != nil {
		addr.PostalCode = req.PostalCode
	}

	if err := s.repo.UpdateAddress(ctx, addr); err != nil {
		return nil, err
	}

	response := toAddressResponse(addr)
	return &response, nil
}

func (s *Service) DeleteAddress(ctx context.Context, employeeID, addressID string) error {
	addrUID, err := uuid.Parse(addressID)
	if err != nil {
		return fmt.Errorf("invalid address id: %w", err)
	}
	return s.repo.DeleteAddress(ctx, addrUID)
}

// =========================================================================
// Sub-module CRUD: Emergency Contacts
// =========================================================================

func (s *Service) CreateEmergencyContact(ctx context.Context, employeeID string, req CreateEmergencyContactRequest) (*EmergencyContactResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	contact := &EmergencyContact{
		EmployeeID:  &empUID,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
	}
	if req.RelationshipTypeID != nil && *req.RelationshipTypeID != "" {
		id, _ := uuid.Parse(*req.RelationshipTypeID)
		contact.RelationshipTypeID = &id
	}
	if req.Address != nil {
		contact.Address = req.Address
	}
	contact.CreatedBy = authctx.GetUserID(ctx)
	contact.UpdatedBy = contact.CreatedBy

	if err := s.encryptIfEnabled(ctx, "emergency_contact.phone_number", &contact.PhoneNumber); err != nil {
		return nil, err
	}

	if err := s.repo.CreateEmergencyContact(ctx, contact); err != nil {
		return nil, err
	}

	response := toEmergencyContactResponse(contact)
	maskEmergencyContactResponse(ctx, &response)
	return &response, nil
}

func (s *Service) UpdateEmergencyContact(ctx context.Context, employeeID, contactID string, req UpdateEmergencyContactRequest) (*EmergencyContactResponse, error) {
	contUID, err := uuid.Parse(contactID)
	if err != nil {
		return nil, fmt.Errorf("invalid contact id: %w", err)
	}

	contact, err := s.repo.FindEmergencyContactByID(ctx, contUID)
	if err != nil {
		return nil, err
	}
	contact.UpdatedBy = authctx.GetUserID(ctx)

	if req.Name != nil {
		contact.Name = *req.Name
	}
	applyIfNotMasked(req.PhoneNumber, &contact.PhoneNumber)
	if req.RelationshipTypeID != nil && *req.RelationshipTypeID != "" {
		id, _ := uuid.Parse(*req.RelationshipTypeID)
		contact.RelationshipTypeID = &id
	}
	if req.Address != nil {
		contact.Address = req.Address
	}

	if err := s.encryptIfEnabled(ctx, "emergency_contact.phone_number", &contact.PhoneNumber); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateEmergencyContact(ctx, contact); err != nil {
		return nil, err
	}

	response := toEmergencyContactResponse(contact)
	maskEmergencyContactResponse(ctx, &response)
	return &response, nil
}

func (s *Service) DeleteEmergencyContact(ctx context.Context, employeeID, contactID string) error {
	contUID, err := uuid.Parse(contactID)
	if err != nil {
		return fmt.Errorf("invalid contact id: %w", err)
	}
	return s.repo.DeleteEmergencyContact(ctx, contUID)
}

// =========================================================================
// Sub-module CRUD: Families
// =========================================================================

func (s *Service) CreateFamily(ctx context.Context, employeeID string, req CreateFamilyRequest) (*FamilyResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	fam := &EmployeeFamily{
		EmployeeID: &empUID,
		Name:       req.Name,
	}
	if req.NIK != nil {
		fam.NIK = req.NIK
	}
	if req.DOB != nil {
		fam.DOB = req.DOB
	}
	if req.RelationshipTypeID != nil && *req.RelationshipTypeID != "" {
		id, _ := uuid.Parse(*req.RelationshipTypeID)
		fam.RelationshipTypeID = &id
	}
	if req.EducationID != nil && *req.EducationID != "" {
		id, _ := uuid.Parse(*req.EducationID)
		fam.EducationID = &id
	}
	fam.CreatedBy = authctx.GetUserID(ctx)
	fam.UpdatedBy = fam.CreatedBy

	if err := s.encryptIfEnabled(ctx, "employee_family.nik", fam.NIK); err != nil {
		return nil, err
	}

	if err := s.repo.CreateFamily(ctx, fam); err != nil {
		return nil, err
	}

	response := toFamilyResponse(fam)
	maskFamilyResponse(ctx, &response)
	return &response, nil
}

func (s *Service) UpdateFamily(ctx context.Context, employeeID, familyID string, req UpdateFamilyRequest) (*FamilyResponse, error) {
	famUID, err := uuid.Parse(familyID)
	if err != nil {
		return nil, fmt.Errorf("invalid family id: %w", err)
	}

	fam, err := s.repo.FindFamilyByID(ctx, famUID)
	if err != nil {
		return nil, err
	}
	fam.UpdatedBy = authctx.GetUserID(ctx)

	if req.Name != nil {
		fam.Name = *req.Name
	}
	applyPtrIfNotMasked(req.NIK, &fam.NIK)
	if req.DOB != nil {
		fam.DOB = req.DOB
	}
	if req.RelationshipTypeID != nil && *req.RelationshipTypeID != "" {
		id, _ := uuid.Parse(*req.RelationshipTypeID)
		fam.RelationshipTypeID = &id
	}
	if req.EducationID != nil && *req.EducationID != "" {
		id, _ := uuid.Parse(*req.EducationID)
		fam.EducationID = &id
	}

	if err := s.encryptIfEnabled(ctx, "employee_family.nik", fam.NIK); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateFamily(ctx, fam); err != nil {
		return nil, err
	}

	response := toFamilyResponse(fam)
	maskFamilyResponse(ctx, &response)
	return &response, nil
}

func (s *Service) DeleteFamily(ctx context.Context, employeeID, familyID string) error {
	famUID, err := uuid.Parse(familyID)
	if err != nil {
		return fmt.Errorf("invalid family id: %w", err)
	}
	return s.repo.DeleteFamily(ctx, famUID)
}

// =========================================================================
// Sub-module CRUD: Educations
// =========================================================================

func (s *Service) CreateEducation(ctx context.Context, employeeID string, req CreateEducationRequest) (*EducationResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	edu := &EmployeeEducation{
		EmployeeID: &empUID,
		Name:       req.Name,
	}
	if req.EducationID != nil && *req.EducationID != "" {
		id, _ := uuid.Parse(*req.EducationID)
		edu.EducationID = &id
	}
	if req.EducationMajorID != nil && *req.EducationMajorID != "" {
		id, _ := uuid.Parse(*req.EducationMajorID)
		edu.EducationMajorID = &id
	}
	if req.Major != nil {
		edu.Major = req.Major
	}
	if req.GradYear != nil {
		edu.GradYear = req.GradYear
	}
	edu.CreatedBy = authctx.GetUserID(ctx)
	edu.UpdatedBy = edu.CreatedBy

	if err := s.repo.CreateEducation(ctx, edu); err != nil {
		return nil, err
	}

	response := toEducationResponse(edu)
	return &response, nil
}

func (s *Service) UpdateEducation(ctx context.Context, employeeID, educationID string, req UpdateEducationRequest) (*EducationResponse, error) {
	eduUID, err := uuid.Parse(educationID)
	if err != nil {
		return nil, fmt.Errorf("invalid education id: %w", err)
	}

	edu, err := s.repo.FindEducationByID(ctx, eduUID)
	if err != nil {
		return nil, err
	}
	edu.UpdatedBy = authctx.GetUserID(ctx)

	if req.Name != nil {
		edu.Name = *req.Name
	}
	if req.EducationID != nil && *req.EducationID != "" {
		id, _ := uuid.Parse(*req.EducationID)
		edu.EducationID = &id
	}
	if req.EducationMajorID != nil {
		if *req.EducationMajorID == "" {
			edu.EducationMajorID = nil
		} else if id, err := uuid.Parse(*req.EducationMajorID); err == nil {
			edu.EducationMajorID = &id
		}
	}
	if req.Major != nil {
		edu.Major = req.Major
	}
	if req.GradYear != nil {
		edu.GradYear = req.GradYear
	}

	if err := s.repo.UpdateEducation(ctx, edu); err != nil {
		return nil, err
	}

	response := toEducationResponse(edu)
	return &response, nil
}

func (s *Service) DeleteEducation(ctx context.Context, employeeID, educationID string) error {
	eduUID, err := uuid.Parse(educationID)
	if err != nil {
		return fmt.Errorf("invalid education id: %w", err)
	}
	return s.repo.DeleteEducation(ctx, eduUID)
}

// =========================================================================
// Sub-module CRUD: Experiences
// =========================================================================

func (s *Service) CreateExperience(ctx context.Context, employeeID string, req CreateExperienceRequest) (*ExperienceResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	exp := &EmployeeExperience{
		EmployeeID: &empUID,
		Company:    req.Company,
	}
	if req.Position != nil {
		exp.Position = req.Position
	}
	if req.StartYear != nil {
		exp.StartYear = req.StartYear
	}
	if req.EndYear != nil {
		exp.EndYear = req.EndYear
	}
	exp.CreatedBy = authctx.GetUserID(ctx)
	exp.UpdatedBy = exp.CreatedBy

	if err := s.repo.CreateExperience(ctx, exp); err != nil {
		return nil, err
	}

	response := toExperienceResponse(exp)
	return &response, nil
}

func (s *Service) UpdateExperience(ctx context.Context, employeeID, experienceID string, req UpdateExperienceRequest) (*ExperienceResponse, error) {
	expUID, err := uuid.Parse(experienceID)
	if err != nil {
		return nil, fmt.Errorf("invalid experience id: %w", err)
	}

	exp, err := s.repo.FindExperienceByID(ctx, expUID)
	if err != nil {
		return nil, err
	}
	exp.UpdatedBy = authctx.GetUserID(ctx)

	if req.Company != nil {
		exp.Company = *req.Company
	}
	if req.Position != nil {
		exp.Position = req.Position
	}
	if req.StartYear != nil {
		exp.StartYear = req.StartYear
	}
	if req.EndYear != nil {
		exp.EndYear = req.EndYear
	}

	if err := s.repo.UpdateExperience(ctx, exp); err != nil {
		return nil, err
	}

	response := toExperienceResponse(exp)
	return &response, nil
}

func (s *Service) DeleteExperience(ctx context.Context, employeeID, experienceID string) error {
	expUID, err := uuid.Parse(experienceID)
	if err != nil {
		return fmt.Errorf("invalid experience id: %w", err)
	}
	return s.repo.DeleteExperience(ctx, expUID)
}

// =========================================================================
// Sub-module CRUD: Documents
// =========================================================================

func (s *Service) CreateDocument(ctx context.Context, employeeID string, req CreateDocumentRequest) (*DocumentResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	doc := &EmployeeDocument{
		EmployeeID: &empUID,
		Name:       req.Name,
		File:       req.File,
	}
	if req.Note != nil {
		doc.Note = req.Note
	}
	doc.CreatedBy = authctx.GetUserID(ctx)
	doc.UpdatedBy = doc.CreatedBy

	if err := s.repo.CreateDocument(ctx, doc); err != nil {
		return nil, err
	}

	response := toDocumentResponse(doc)
	return &response, nil
}

func (s *Service) UpdateDocument(ctx context.Context, employeeID, documentID string, req UpdateDocumentRequest) (*DocumentResponse, error) {
	docUID, err := uuid.Parse(documentID)
	if err != nil {
		return nil, fmt.Errorf("invalid document id: %w", err)
	}

	doc, err := s.repo.FindDocumentByID(ctx, docUID)
	if err != nil {
		return nil, err
	}
	doc.UpdatedBy = authctx.GetUserID(ctx)

	if req.Name != nil {
		doc.Name = *req.Name
	}
	if req.File != nil {
		doc.File = *req.File
	}
	if req.Note != nil {
		doc.Note = req.Note
	}

	if err := s.repo.UpdateDocument(ctx, doc); err != nil {
		return nil, err
	}

	response := toDocumentResponse(doc)
	return &response, nil
}

func (s *Service) DeleteDocument(ctx context.Context, employeeID, documentID string) error {
	docUID, err := uuid.Parse(documentID)
	if err != nil {
		return fmt.Errorf("invalid document id: %w", err)
	}
	return s.repo.DeleteDocument(ctx, docUID)
}

// CreateDocumentRecord creates a document record directly (for upload flow).
func (s *Service) CreateDocumentRecord(ctx context.Context, doc *EmployeeDocument) error {
	return s.repo.CreateDocument(ctx, doc)
}

// UpdateDocumentFile updates document metadata and file path.
func (s *Service) UpdateDocumentFile(ctx context.Context, documentID, name, filePath string, note *string) error {
	docUID, err := uuid.Parse(documentID)
	if err != nil {
		return fmt.Errorf("invalid document id: %w", err)
	}

	doc, err := s.repo.FindDocumentByID(ctx, docUID)
	if err != nil {
		return err
	}
	doc.UpdatedBy = authctx.GetUserID(ctx)

	if name != "" {
		doc.Name = name
	}
	if filePath != "" {
		doc.File = filePath
	}
	if note != nil {
		doc.Note = note
	}

	return s.repo.UpdateDocument(ctx, doc)
}

// =========================================================================
// Sub-module CRUD: Insurances
// =========================================================================

func (s *Service) CreateInsurance(ctx context.Context, employeeID string, req CreateInsuranceRequest) (*InsuranceResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	ins := &EmployeeInsurance{
		EmployeeID: &empUID,
		Number:     req.Number,
	}
	if req.InsuranceID != nil && *req.InsuranceID != "" {
		if id, err := uuid.Parse(*req.InsuranceID); err == nil {
			ins.InsuranceID = &id
		}
	}
	if req.Type != nil {
		ins.Type = req.Type
	}
	ins.CreatedBy = authctx.GetUserID(ctx)
	ins.UpdatedBy = ins.CreatedBy

	if err := s.repo.CreateInsurance(ctx, ins); err != nil {
		return nil, err
	}

	response := toInsuranceResponse(ins)
	return &response, nil
}

func (s *Service) UpdateInsurance(ctx context.Context, employeeID, insuranceID string, req UpdateInsuranceRequest) (*InsuranceResponse, error) {
	insUID, err := uuid.Parse(insuranceID)
	if err != nil {
		return nil, fmt.Errorf("invalid insurance id: %w", err)
	}

	ins, err := s.repo.FindInsuranceByID(ctx, insUID)
	if err != nil {
		return nil, err
	}
	ins.UpdatedBy = authctx.GetUserID(ctx)

	if req.Number != nil {
		ins.Number = *req.Number
	}
	if req.InsuranceID != nil {
		if *req.InsuranceID == "" {
			ins.InsuranceID = nil
		} else if id, err := uuid.Parse(*req.InsuranceID); err == nil {
			ins.InsuranceID = &id
		}
	}
	if req.Type != nil {
		ins.Type = req.Type
	}

	if err := s.repo.UpdateInsurance(ctx, ins); err != nil {
		return nil, err
	}

	response := toInsuranceResponse(ins)
	return &response, nil
}

func (s *Service) DeleteInsurance(ctx context.Context, employeeID, insuranceID string) error {
	insUID, err := uuid.Parse(insuranceID)
	if err != nil {
		return fmt.Errorf("invalid insurance id: %w", err)
	}
	return s.repo.DeleteInsurance(ctx, insUID)
}

// =========================================================================
// Sub-module CRUD: Banks
// =========================================================================

func (s *Service) CreateBank(ctx context.Context, employeeID string, req CreateBankRequest) (*BankResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	bank := &EmployeeBankAccount{
		EmployeeID:    &empUID,
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
	}
	if req.BankID != nil && *req.BankID != "" {
		id, _ := uuid.Parse(*req.BankID)
		bank.BankID = &id
	}
	bank.CreatedBy = authctx.GetUserID(ctx)
	bank.UpdatedBy = bank.CreatedBy

	if err := s.encryptIfEnabled(ctx, "employee_bank_account.account_number", &bank.AccountNumber); err != nil {
		return nil, err
	}
	if err := s.encryptIfEnabled(ctx, "employee_bank_account.account_name", &bank.AccountName); err != nil {
		return nil, err
	}

	if err := s.repo.CreateBank(ctx, bank); err != nil {
		return nil, err
	}

	response := toBankResponse(bank)
	maskBankResponse(ctx, &response)
	return &response, nil
}

func (s *Service) UpdateBank(ctx context.Context, employeeID, bankID string, req UpdateBankRequest) (*BankResponse, error) {
	bankUID, err := uuid.Parse(bankID)
	if err != nil {
		return nil, fmt.Errorf("invalid bank id: %w", err)
	}

	bank, err := s.repo.FindBankByID(ctx, bankUID)
	if err != nil {
		return nil, err
	}
	bank.UpdatedBy = authctx.GetUserID(ctx)

	if req.BankID != nil && *req.BankID != "" {
		id, _ := uuid.Parse(*req.BankID)
		bank.BankID = &id
	}
	applyIfNotMasked(req.AccountNumber, &bank.AccountNumber)
	applyIfNotMasked(req.AccountName, &bank.AccountName)

	if err := s.encryptIfEnabled(ctx, "employee_bank_account.account_number", &bank.AccountNumber); err != nil {
		return nil, err
	}
	if err := s.encryptIfEnabled(ctx, "employee_bank_account.account_name", &bank.AccountName); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateBank(ctx, bank); err != nil {
		return nil, err
	}

	response := toBankResponse(bank)
	maskBankResponse(ctx, &response)
	return &response, nil
}

func (s *Service) DeleteBank(ctx context.Context, employeeID, bankID string) error {
	bankUID, err := uuid.Parse(bankID)
	if err != nil {
		return fmt.Errorf("invalid bank id: %w", err)
	}
	return s.repo.DeleteBank(ctx, bankUID)
}

// =========================================================================
// Sub-module CRUD: Employments
// =========================================================================

func (s *Service) CreateEmployment(ctx context.Context, employeeID string, req CreateEmploymentRequest) (*EmploymentResponse, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	empl := &Employment{
		EmployeeID:          &empUID,
		DecisionLetterNumber: req.DecisionLetterNumber,
		DecisionLetterDate:   req.DecisionLetterDate,
		EffectiveDate:        req.EffectiveDate,
	}
	if req.OrganizationID != nil && *req.OrganizationID != "" {
		id, _ := uuid.Parse(*req.OrganizationID)
		empl.OrganizationID = &id
	}
	if req.PositionID != nil && *req.PositionID != "" {
		id, _ := uuid.Parse(*req.PositionID)
		empl.PositionID = &id
	}
	if req.EmploymentStatusID != nil && *req.EmploymentStatusID != "" {
		id, _ := uuid.Parse(*req.EmploymentStatusID)
		empl.EmploymentStatusID = &id
	}
	if req.EffectiveEndDate != nil {
		empl.EffectiveEndDate = req.EffectiveEndDate
	}
	empl.CreatedBy = authctx.GetUserID(ctx)
	empl.UpdatedBy = empl.CreatedBy

	if err := s.repo.CreateEmployment(ctx, empl); err != nil {
		return nil, err
	}

	response := toEmploymentResponse(empl)
	return &response, nil
}

func (s *Service) UpdateEmployment(ctx context.Context, employeeID, employmentID string, req UpdateEmploymentRequest) (*EmploymentResponse, error) {
	emplUID, err := uuid.Parse(employmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid employment id: %w", err)
	}

	empl, err := s.repo.FindEmploymentByID(ctx, emplUID)
	if err != nil {
		return nil, err
	}
	empl.UpdatedBy = authctx.GetUserID(ctx)

	if req.DecisionLetterNumber != nil {
		empl.DecisionLetterNumber = *req.DecisionLetterNumber
	}
	if req.DecisionLetterDate != nil {
		empl.DecisionLetterDate = *req.DecisionLetterDate
	}
	if req.EffectiveDate != nil {
		empl.EffectiveDate = *req.EffectiveDate
	}
	if req.OrganizationID != nil && *req.OrganizationID != "" {
		id, _ := uuid.Parse(*req.OrganizationID)
		empl.OrganizationID = &id
	}
	if req.PositionID != nil && *req.PositionID != "" {
		id, _ := uuid.Parse(*req.PositionID)
		empl.PositionID = &id
	}
	if req.EmploymentStatusID != nil && *req.EmploymentStatusID != "" {
		id, _ := uuid.Parse(*req.EmploymentStatusID)
		empl.EmploymentStatusID = &id
	}
	if req.EffectiveEndDate != nil {
		empl.EffectiveEndDate = req.EffectiveEndDate
	}

	if err := s.repo.UpdateEmployment(ctx, empl); err != nil {
		return nil, err
	}

	response := toEmploymentResponse(empl)
	return &response, nil
}

func (s *Service) DeleteEmployment(ctx context.Context, employeeID, employmentID string) error {
	emplUID, err := uuid.Parse(employmentID)
	if err != nil {
		return fmt.Errorf("invalid employment id: %w", err)
	}
	return s.repo.DeleteEmployment(ctx, emplUID)
}
