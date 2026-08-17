package employee

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/modules/setting"
)

type Repository struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
}

func NewRepository(dbResolver func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbResolver: dbResolver}
}

func (r *Repository) getDB(ctx context.Context) (*gorm.DB, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required for tenant database resolution")
	}
	return r.dbResolver(ctx)
}

// =========================================================================
// Employee
// =========================================================================

func (r *Repository) CreateEmployee(ctx context.Context, emp *Employee) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(emp).Error
}

func (r *Repository) FindEmployeeByID(ctx context.Context, id uuid.UUID) (*Employee, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var emp Employee
	q := db.Preload("Addresses").Preload("EmergencyContacts.RelationshipType").
		Preload("Families.RelationshipType").Preload("Families.Education").
		Preload("Educations.Education").Preload("Educations.EducationMajor").
		Preload("Experiences").Preload("Documents").
		Preload("Insurances.Insurance").Preload("Banks.Bank").
		Preload("Employments.Organization").Preload("Employments.EmploymentStatus")
	if err := q.First(&emp, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}
	return &emp, nil
}

// CountByGender menghitung jumlah karyawan per jenis kelamin (kolom gender:
// "M" / "F" / lainnya/kosong). Dipakai pie chart dashboard Employment.
func (r *Repository) CountByGender(ctx context.Context) (male, female, other int64, err error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return 0, 0, 0, err
	}
	type row struct {
		Gender string
		Count  int64
	}
	var rows []row
	if err := db.Table("employees").
		Where("status = ?", "active").
		Select("COALESCE(gender, '') AS gender, COUNT(*) AS count").
		Group("gender").
		Scan(&rows).Error; err != nil {
		return 0, 0, 0, err
	}
	for _, r := range rows {
		switch r.Gender {
		case "M":
			male = r.Count
		case "F":
			female = r.Count
		default:
			other += r.Count
		}
	}
	return male, female, other, nil
}

// CountByEmploymentStatus menghitung jumlah karyawan per status kepegawaian
// berdasarkan employment berjalan (effective_end_date NULL/''). Karyawan tanpa
// employment berjalan dihitung sebagai unclassified. Nama status diambil dari
// tabel employment_statuses (tenant-configurable).
func (r *Repository) CountByEmploymentStatus(ctx context.Context) ([]EmploymentStatusCount, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	type row struct {
		Name  string
		Count int64
	}
	var rows []row
	// Hanya karyawan aktif. Employment berjalan = effective_end_date NULL
	// (kolom DATE — tidak boleh dibandingkan dengan '' di MySQL).
	if err := db.Table("employments AS emp").
		Joins("JOIN employees AS e ON e.id = emp.employee_id AND e.status = 'active'").
		Joins("LEFT JOIN employment_statuses AS es ON es.id = emp.employment_status_id").
		Where("emp.effective_end_date IS NULL").
		Select("COALESCE(es.name, '') AS name, COUNT(DISTINCT e.id) AS count").
		Group("es.name").
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	groups := make([]EmploymentStatusCount, 0, len(rows))
	for _, r := range rows {
		groups = append(groups, EmploymentStatusCount{Name: r.Name, Count: r.Count})
	}
	var unclassified int64
	if err := db.Table("employees AS e").
		Where("status = ? AND NOT EXISTS (SELECT 1 FROM employments emp WHERE emp.employee_id = e.id AND emp.effective_end_date IS NULL)", "active").
		Count(&unclassified).Error; err != nil {
		return nil, 0, err
	}
	return groups, unclassified, nil
}

// ResolveEmployeeRefNames mengambil nama referensi (agama, status perkawinan,
// kewarganegaraan) untuk ditampilkan pada response detail employee.
// Nama di-resolve langsung dari tabel tenant (religions / marital_statuses /
// nationalities) supaya halaman profile tidak bergantung pada permission
// viewer terhadap endpoint /settings/* — dan tidak menampilkan ID mentah
// ketika referensi tidak ditemukan (mengembalikan string kosong).
func (r *Repository) ResolveEmployeeRefNames(ctx context.Context, religionID, maritalStatusID *uuid.UUID, nationalityID *string) (religionName, maritalStatusName, nationalityName string) {
	db, err := r.getDB(ctx)
	if err != nil {
		return "", "", ""
	}
	if religionID != nil {
		var row setting.Religion
		if err := db.First(&row, "id = ?", *religionID).Error; err == nil {
			religionName = row.Name
		}
	}
	if maritalStatusID != nil {
		var row setting.MaritalStatus
		if err := db.First(&row, "id = ?", *maritalStatusID).Error; err == nil {
			maritalStatusName = row.Name
		}
	}
	if nationalityID != nil {
		// NationalityID employee menyimpan kode (char(2)) — join via kolom code.
		var row setting.Nationality
		if err := db.First(&row, "code = ?", *nationalityID).Error; err == nil {
			nationalityName = row.Name
		}
	}
	return religionName, maritalStatusName, nationalityName
}

func (r *Repository) FindEmployeeByEmployeeID(ctx context.Context, employeeID string) (*Employee, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var emp Employee
	if err := db.First(&emp, "employee_id = ?", employeeID).Error; err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}
	return &emp, nil
}

// CountEmployees menghitung jumlah total employee (untuk kuota on-premise).
func (r *Repository) CountEmployees(ctx context.Context) (int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return 0, err
	}
	var total int64
	if err := db.Model(&Employee{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *Repository) FindAllEmployees(ctx context.Context, page, perPage int, search, status, organizationID string) ([]Employee, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var employees []Employee
	var total int64

	query := db.Model(&Employee{})

	// Apply search filter
	if search != "" {
		like := "%" + search + "%"
		query = query.Where(
		db.Where("name LIKE ?", like).
				Or("employee_id LIKE ?", like).
				Or("nik LIKE ?", like),
		)
	}

	// Apply status filter
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Apply organization filter — employee dengan employment saat ini di org tsb.
	if organizationID != "" {
		query = query.Where("id IN (SELECT employee_id FROM employments WHERE organization_id = ? AND effective_end_date IS NULL)", organizationID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("name ASC").Find(&employees).Error; err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}

func (r *Repository) UpdateEmployee(ctx context.Context, emp *Employee) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(emp).Error
}

// SetEmployeeStatus updates only the employee's status column (used by
// movement execution to mark offboarded/retired employees inactive — avoids
// the heavy Preload chain in FindEmployeeByID). Returns an error when the
// employee does not exist.
func (r *Repository) SetEmployeeStatus(ctx context.Context, id uuid.UUID, status string) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return setEmployeeStatus(db, id, status)
}

// SetEmployeeStatusTx is the transactional variant of SetEmployeeStatus — runs
// on the caller's transaction (used by employeemovement.ExecuteMovementTx).
func (r *Repository) SetEmployeeStatusTx(_ context.Context, db *gorm.DB, id uuid.UUID, status string) error {
	return setEmployeeStatus(db, id, status)
}

func setEmployeeStatus(db *gorm.DB, id uuid.UUID, status string) error {
	result := db.Model(&Employee{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to set employee status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}
	return nil
}

func (r *Repository) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&Employee{}).Error
}

// =========================================================================
// Addresses (nested under employee)
// =========================================================================

func (r *Repository) CreateAddress(ctx context.Context, addr *EmployeeAddress) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(addr).Error
}

func (r *Repository) FindAddressByID(ctx context.Context, id uuid.UUID) (*EmployeeAddress, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var addr EmployeeAddress
	if err := db.First(&addr, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("address not found: %w", err)
	}
	return &addr, nil
}

func (r *Repository) UpdateAddress(ctx context.Context, addr *EmployeeAddress) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(addr).Error
}

func (r *Repository) DeleteAddress(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeeAddress{}).Error
}

// =========================================================================
// Emergency Contacts
// =========================================================================

func (r *Repository) CreateEmergencyContact(ctx context.Context, contact *EmergencyContact) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(contact).Error
}

func (r *Repository) FindEmergencyContactByID(ctx context.Context, id uuid.UUID) (*EmergencyContact, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var contact EmergencyContact
	if err := db.First(&contact, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("emergency contact not found: %w", err)
	}
	return &contact, nil
}

func (r *Repository) UpdateEmergencyContact(ctx context.Context, contact *EmergencyContact) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(contact).Error
}

func (r *Repository) DeleteEmergencyContact(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmergencyContact{}).Error
}

// =========================================================================
// Families
// =========================================================================

func (r *Repository) CreateFamily(ctx context.Context, fam *EmployeeFamily) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(fam).Error
}

func (r *Repository) FindFamilyByID(ctx context.Context, id uuid.UUID) (*EmployeeFamily, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var fam EmployeeFamily
	if err := db.First(&fam, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("family not found: %w", err)
	}
	return &fam, nil
}

func (r *Repository) UpdateFamily(ctx context.Context, fam *EmployeeFamily) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(fam).Error
}

func (r *Repository) DeleteFamily(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeeFamily{}).Error
}

// =========================================================================
// Educations
// =========================================================================

func (r *Repository) CreateEducation(ctx context.Context, edu *EmployeeEducation) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(edu).Error
}

func (r *Repository) FindEducationByID(ctx context.Context, id uuid.UUID) (*EmployeeEducation, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var edu EmployeeEducation
	if err := db.First(&edu, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("education not found: %w", err)
	}
	return &edu, nil
}

func (r *Repository) UpdateEducation(ctx context.Context, edu *EmployeeEducation) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(edu).Error
}

func (r *Repository) DeleteEducation(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeeEducation{}).Error
}

// =========================================================================
// Experiences
// =========================================================================

func (r *Repository) CreateExperience(ctx context.Context, exp *EmployeeExperience) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(exp).Error
}

func (r *Repository) FindExperienceByID(ctx context.Context, id uuid.UUID) (*EmployeeExperience, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var exp EmployeeExperience
	if err := db.First(&exp, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("experience not found: %w", err)
	}
	return &exp, nil
}

func (r *Repository) UpdateExperience(ctx context.Context, exp *EmployeeExperience) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(exp).Error
}

func (r *Repository) DeleteExperience(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeeExperience{}).Error
}

// =========================================================================
// Documents
// =========================================================================

func (r *Repository) CreateDocument(ctx context.Context, doc *EmployeeDocument) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(doc).Error
}

func (r *Repository) FindDocumentByID(ctx context.Context, id uuid.UUID) (*EmployeeDocument, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var doc EmployeeDocument
	if err := db.First(&doc, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}
	return &doc, nil
}

func (r *Repository) UpdateDocument(ctx context.Context, doc *EmployeeDocument) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(doc).Error
}

func (r *Repository) DeleteDocument(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeeDocument{}).Error
}

// =========================================================================
// Insurances
// =========================================================================

func (r *Repository) CreateInsurance(ctx context.Context, ins *EmployeeInsurance) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(ins).Error
}

func (r *Repository) FindInsuranceByID(ctx context.Context, id uuid.UUID) (*EmployeeInsurance, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var ins EmployeeInsurance
	if err := db.First(&ins, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("insurance not found: %w", err)
	}
	return &ins, nil
}

func (r *Repository) UpdateInsurance(ctx context.Context, ins *EmployeeInsurance) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(ins).Error
}

func (r *Repository) DeleteInsurance(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeeInsurance{}).Error
}

// =========================================================================
// Banks
// =========================================================================

func (r *Repository) CreateBank(ctx context.Context, bank *EmployeeBankAccount) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(bank).Error
}

func (r *Repository) FindBankByID(ctx context.Context, id uuid.UUID) (*EmployeeBankAccount, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var bank EmployeeBankAccount
	if err := db.First(&bank, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("bank not found: %w", err)
	}
	return &bank, nil
}

func (r *Repository) UpdateBank(ctx context.Context, bank *EmployeeBankAccount) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(bank).Error
}

func (r *Repository) DeleteBank(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&EmployeeBankAccount{}).Error
}

// =========================================================================
// Employments
// =========================================================================

func (r *Repository) CreateEmployment(ctx context.Context, emp *Employment) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return createEmployment(db, emp)
}

// CreateEmploymentTx is the transactional variant of CreateEmployment — runs
// on the caller's transaction (used by employeemovement.ExecuteMovementTx so
// the new employment commits/rolls back with the rest of the execution).
func (r *Repository) CreateEmploymentTx(_ context.Context, db *gorm.DB, emp *Employment) error {
	return createEmployment(db, emp)
}

func createEmployment(db *gorm.DB, emp *Employment) error {
	return db.Create(emp).Error
}

func (r *Repository) FindEmploymentByID(ctx context.Context, id uuid.UUID) (*Employment, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var emp Employment
	if err := db.First(&emp, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("employment not found: %w", err)
	}
	return &emp, nil
}

// FindActiveEmploymentByEmployeeID returns the employee's currently active
// employment — the most recent record with no effective_end_date. Returns
// nil (without error) when the employee has no active employment.
func (r *Repository) FindActiveEmploymentByEmployeeID(ctx context.Context, employeeID uuid.UUID) (*Employment, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	return findActiveEmploymentByEmployeeID(db, employeeID)
}

// FindActiveEmploymentByEmployeeIDTx is the transactional variant of
// FindActiveEmploymentByEmployeeID — it runs the same lookup on the given
// *gorm.DB (a transaction opened by a consumer module, e.g. employeemovement's
// ExecuteMovementTx) so the read is part of the caller's atomic unit of work.
func (r *Repository) FindActiveEmploymentByEmployeeIDTx(_ context.Context, db *gorm.DB, employeeID uuid.UUID) (*Employment, error) {
	return findActiveEmploymentByEmployeeID(db, employeeID)
}

func findActiveEmploymentByEmployeeID(db *gorm.DB, employeeID uuid.UUID) (*Employment, error) {
	var emp Employment
	err := db.
		Where("employee_id = ? AND effective_end_date IS NULL", employeeID).
		Order("effective_date DESC, created_at DESC").
		First(&emp).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find active employment: %w", err)
	}
	return &emp, nil
}

func (r *Repository) UpdateEmployment(ctx context.Context, emp *Employment) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(emp).Error
}

// CloseEmployment sets the employment's effective_end_date (used by movement
// execution to close the previous employment the day before the new one takes
// effect). Returns an error if the employment does not exist or is already
// closed (defensive guard: never overwrite an existing end date).
func (r *Repository) CloseEmployment(ctx context.Context, id uuid.UUID, effectiveEndDate string) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return closeEmployment(db, id, effectiveEndDate)
}

// CloseEmploymentTx is the transactional variant of CloseEmployment — runs on
// the caller's transaction so the close is part of the same atomic unit of
// work (used by employeemovement.ExecuteMovementTx).
func (r *Repository) CloseEmploymentTx(_ context.Context, db *gorm.DB, id uuid.UUID, effectiveEndDate string) error {
	return closeEmployment(db, id, effectiveEndDate)
}

func closeEmployment(db *gorm.DB, id uuid.UUID, effectiveEndDate string) error {
	result := db.Model(&Employment{}).
		Where("id = ? AND effective_end_date IS NULL", id).
		Update("effective_end_date", effectiveEndDate)
	if result.Error != nil {
		return fmt.Errorf("failed to close employment: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employment not found or already closed")
	}
	return nil
}

func (r *Repository) DeleteEmployment(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&Employment{}).Error
}
