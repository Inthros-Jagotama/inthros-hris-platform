package setting

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

// ── Generic CRUD helpers ──

func (r *Repository) create(ctx context.Context, model interface{}) error {
	db, err := r.getDB(ctx)
	if err != nil { return err }
	return db.Create(model).Error
}

func (r *Repository) findByID(ctx context.Context, model interface{}, id interface{}, table string) error {
	db, err := r.getDB(ctx)
	if err != nil { return err }
	return db.First(model, "id = ?", id).Error
}

func (r *Repository) findAll(ctx context.Context, model interface{}, table string, page, perPage int, order string) (int64, error) {
	db, err := r.getDB(ctx)
	if err != nil { return 0, err }
	var total int64
	query := db.Model(model)
	if err := query.Count(&total).Error; err != nil { return 0, err }
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order(order).Find(model).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *Repository) findAllByParent(ctx context.Context, model interface{}, parentField string, parentID interface{}, order string) error {
	db, err := r.getDB(ctx)
	if err != nil { return err }
	return db.Where(fmt.Sprintf("%s = ?", parentField), parentID).Order(order).Find(model).Error
}

func (r *Repository) update(ctx context.Context, model interface{}) error {
	db, err := r.getDB(ctx)
	if err != nil { return err }
	return db.Save(model).Error
}

func (r *Repository) softDelete(ctx context.Context, model interface{}, id interface{}) error {
	db, err := r.getDB(ctx)
	if err != nil { return err }
	return db.Where("id = ?", id).Delete(model).Error
}

func (r *Repository) findByCode(ctx context.Context, model interface{}, code string, table string) (bool, error) {
	db, err := r.getDB(ctx)
	if err != nil { return false, err }
	var count int64
	if err := db.Model(model).Where("code = ?", code).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) findByCodeExcludeSelf(ctx context.Context, model interface{}, code string, id interface{}, table string) (bool, error) {
	db, err := r.getDB(ctx)
	if err != nil { return false, err }
	var count int64
	if err := db.Model(model).Where("code = ? AND id != ?", code, id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ── Zone ──
func (r *Repository) CreateZone(ctx context.Context, zone *Zone) error { return r.create(ctx, zone) }
func (r *Repository) FindZoneByID(ctx context.Context, id uuid.UUID) (*Zone, error) {
	var zone Zone
	if err := r.findByID(ctx, &zone, id, "zones"); err != nil {
		return nil, fmt.Errorf("zone not found: %w", err)
	}
	return &zone, nil
}
func (r *Repository) FindAllZones(ctx context.Context, page, perPage int) ([]Zone, int64, error) {
	var zones []Zone
	total, err := r.findAll(ctx, &zones, "zones", page, perPage, "sort_order ASC, code ASC")
	return zones, total, err
}
func (r *Repository) FindAllZonesActive(ctx context.Context) ([]Zone, error) {
	db, err := r.getDB(ctx)
	if err != nil { return nil, err }
	var zones []Zone
	if err := db.Where("is_active = ?", true).Order("sort_order ASC, code ASC").Find(&zones).Error; err != nil {
		return nil, err
	}
	return zones, nil
}
func (r *Repository) UpdateZone(ctx context.Context, zone *Zone) error { return r.update(ctx, zone) }
func (r *Repository) DeleteZone(ctx context.Context, id uuid.UUID) error { return r.softDelete(ctx, &Zone{}, id) }

// ── Province ──
func (r *Repository) CreateProvince(ctx context.Context, p *Province) error { return r.create(ctx, p) }
func (r *Repository) FindProvinceByID(ctx context.Context, id string) (*Province, error) {
	var p Province
	if err := r.findByID(ctx, &p, id, "provinces"); err != nil {
		return nil, fmt.Errorf("province not found: %w", err)
	}
	return &p, nil
}
func (r *Repository) FindAllProvinces(ctx context.Context, page, perPage int) ([]Province, int64, error) {
	var provinces []Province
	total, err := r.findAll(ctx, &provinces, "provinces", page, perPage, "code ASC")
	return provinces, total, err
}
func (r *Repository) FindAllProvincesList(ctx context.Context) ([]Province, error) {
	db, err := r.getDB(ctx)
	if err != nil { return nil, err }
	var provinces []Province
	if err := db.Order("code ASC").Find(&provinces).Error; err != nil { return nil, err }
	return provinces, nil
}
func (r *Repository) UpdateProvince(ctx context.Context, p *Province) error { return r.update(ctx, p) }
func (r *Repository) DeleteProvince(ctx context.Context, id string) error { return r.softDelete(ctx, &Province{}, id) }

// ── Regency ──
func (r *Repository) CreateRegency(ctx context.Context, reg *Regency) error { return r.create(ctx, reg) }
func (r *Repository) FindRegencyByID(ctx context.Context, id string) (*Regency, error) {
	var reg Regency
	if err := r.findByID(ctx, &reg, id, "regencies"); err != nil {
		return nil, fmt.Errorf("regency not found: %w", err)
	}
	return &reg, nil
}
func (r *Repository) FindAllRegencies(ctx context.Context, page, perPage int) ([]Regency, int64, error) {
	var regencies []Regency
	total, err := r.findAll(ctx, &regencies, "regencies", page, perPage, "code ASC")
	return regencies, total, err
}
func (r *Repository) FindRegenciesByProvince(ctx context.Context, provinceID string) ([]Regency, error) {
	var regencies []Regency
	if err := r.findAllByParent(ctx, &regencies, "province_id", provinceID, "code ASC"); err != nil {
		return nil, err
	}
	return regencies, nil
}
func (r *Repository) UpdateRegency(ctx context.Context, reg *Regency) error { return r.update(ctx, reg) }
func (r *Repository) DeleteRegency(ctx context.Context, id string) error { return r.softDelete(ctx, &Regency{}, id) }

// ── District ──
func (r *Repository) CreateDistrict(ctx context.Context, d *District) error { return r.create(ctx, d) }
func (r *Repository) FindDistrictByID(ctx context.Context, id string) (*District, error) {
	var d District
	if err := r.findByID(ctx, &d, id, "districts"); err != nil {
		return nil, fmt.Errorf("district not found: %w", err)
	}
	return &d, nil
}
func (r *Repository) FindAllDistricts(ctx context.Context, page, perPage int) ([]District, int64, error) {
	var districts []District
	total, err := r.findAll(ctx, &districts, "districts", page, perPage, "code ASC")
	return districts, total, err
}
func (r *Repository) FindDistrictsByRegency(ctx context.Context, regencyID string) ([]District, error) {
	var districts []District
	if err := r.findAllByParent(ctx, &districts, "regency_id", regencyID, "code ASC"); err != nil {
		return nil, err
	}
	return districts, nil
}
func (r *Repository) UpdateDistrict(ctx context.Context, d *District) error { return r.update(ctx, d) }
func (r *Repository) DeleteDistrict(ctx context.Context, id string) error { return r.softDelete(ctx, &District{}, id) }

// ── Village ──
func (r *Repository) CreateVillage(ctx context.Context, v *Village) error { return r.create(ctx, v) }
func (r *Repository) FindVillageByID(ctx context.Context, id string) (*Village, error) {
	var v Village
	if err := r.findByID(ctx, &v, id, "villages"); err != nil {
		return nil, fmt.Errorf("village not found: %w", err)
	}
	return &v, nil
}
func (r *Repository) FindAllVillages(ctx context.Context, page, perPage int) ([]Village, int64, error) {
	var villages []Village
	total, err := r.findAll(ctx, &villages, "villages", page, perPage, "code ASC")
	return villages, total, err
}
func (r *Repository) FindVillagesByDistrict(ctx context.Context, districtID string) ([]Village, error) {
	var villages []Village
	if err := r.findAllByParent(ctx, &villages, "district_id", districtID, "code ASC"); err != nil {
		return nil, err
	}
	return villages, nil
}
func (r *Repository) UpdateVillage(ctx context.Context, v *Village) error { return r.update(ctx, v) }
func (r *Repository) DeleteVillage(ctx context.Context, id string) error { return r.softDelete(ctx, &Village{}, id) }
func (r *Repository) SearchVillages(ctx context.Context, query string, limit int) ([]Village, error) {
	db, err := r.getDB(ctx)
	if err != nil { return nil, err }
	var villages []Village
	if err := db.Where("name LIKE ?", "%"+query+"%").Preload("District.Regency.Province").Limit(limit).Find(&villages).Error; err != nil {
		return nil, err
	}
	return villages, nil
}

// ── Education ──
func (r *Repository) CreateEducation(ctx context.Context, e *Education) error { return r.create(ctx, e) }
func (r *Repository) FindEducationByID(ctx context.Context, id uuid.UUID) (*Education, error) {
	var e Education
	if err := r.findByID(ctx, &e, id, "educations"); err != nil {
		return nil, fmt.Errorf("education not found: %w", err)
	}
	return &e, nil
}
func (r *Repository) FindAllEducations(ctx context.Context, page, perPage int) ([]Education, int64, error) {
	var list []Education
	total, err := r.findAll(ctx, &list, "educations", page, perPage, "sort_order ASC, code ASC")
	return list, total, err
}
func (r *Repository) UpdateEducation(ctx context.Context, e *Education) error { return r.update(ctx, e) }
func (r *Repository) DeleteEducation(ctx context.Context, id uuid.UUID) error { return r.softDelete(ctx, &Education{}, id) }

// ── Religion ──
func (r *Repository) CreateReligion(ctx context.Context, rel *Religion) error { return r.create(ctx, rel) }
func (r *Repository) FindReligionByID(ctx context.Context, id uuid.UUID) (*Religion, error) {
	var rel Religion
	if err := r.findByID(ctx, &rel, id, "religions"); err != nil {
		return nil, fmt.Errorf("religion not found: %w", err)
	}
	return &rel, nil
}
func (r *Repository) FindAllReligions(ctx context.Context, page, perPage int) ([]Religion, int64, error) {
	var list []Religion
	total, err := r.findAll(ctx, &list, "religions", page, perPage, "sort_order ASC, code ASC")
	return list, total, err
}
func (r *Repository) UpdateReligion(ctx context.Context, rel *Religion) error { return r.update(ctx, rel) }
func (r *Repository) DeleteReligion(ctx context.Context, id uuid.UUID) error { return r.softDelete(ctx, &Religion{}, id) }

// ── MaritalStatus ──
func (r *Repository) CreateMaritalStatus(ctx context.Context, m *MaritalStatus) error { return r.create(ctx, m) }
func (r *Repository) FindMaritalStatusByID(ctx context.Context, id uuid.UUID) (*MaritalStatus, error) {
	var m MaritalStatus
	if err := r.findByID(ctx, &m, id, "marital_statuses"); err != nil {
		return nil, fmt.Errorf("marital status not found: %w", err)
	}
	return &m, nil
}
func (r *Repository) FindAllMaritalStatuses(ctx context.Context, page, perPage int) ([]MaritalStatus, int64, error) {
	var list []MaritalStatus
	total, err := r.findAll(ctx, &list, "marital_statuses", page, perPage, "sort_order ASC, code ASC")
	return list, total, err
}
func (r *Repository) UpdateMaritalStatus(ctx context.Context, m *MaritalStatus) error { return r.update(ctx, m) }
func (r *Repository) DeleteMaritalStatus(ctx context.Context, id uuid.UUID) error { return r.softDelete(ctx, &MaritalStatus{}, id) }

// ── RelationshipType ──
func (r *Repository) CreateRelationshipType(ctx context.Context, rt *RelationshipType) error { return r.create(ctx, rt) }
func (r *Repository) FindRelationshipTypeByID(ctx context.Context, id uuid.UUID) (*RelationshipType, error) {
	var rt RelationshipType
	if err := r.findByID(ctx, &rt, id, "relationship_types"); err != nil {
		return nil, fmt.Errorf("relationship type not found: %w", err)
	}
	return &rt, nil
}
func (r *Repository) FindAllRelationshipTypes(ctx context.Context, page, perPage int) ([]RelationshipType, int64, error) {
	var list []RelationshipType
	total, err := r.findAll(ctx, &list, "relationship_types", page, perPage, "sort_order ASC, code ASC")
	return list, total, err
}
func (r *Repository) UpdateRelationshipType(ctx context.Context, rt *RelationshipType) error { return r.update(ctx, rt) }
func (r *Repository) DeleteRelationshipType(ctx context.Context, id uuid.UUID) error { return r.softDelete(ctx, &RelationshipType{}, id) }

// ── EmploymentStatus ──
func (r *Repository) CreateEmploymentStatus(ctx context.Context, es *EmploymentStatus) error { return r.create(ctx, es) }
func (r *Repository) FindEmploymentStatusByID(ctx context.Context, id uuid.UUID) (*EmploymentStatus, error) {
	var es EmploymentStatus
	if err := r.findByID(ctx, &es, id, "employment_statuses"); err != nil {
		return nil, fmt.Errorf("employment status not found: %w", err)
	}
	return &es, nil
}
func (r *Repository) FindAllEmploymentStatuses(ctx context.Context, page, perPage int) ([]EmploymentStatus, int64, error) {
	var list []EmploymentStatus
	total, err := r.findAll(ctx, &list, "employment_statuses", page, perPage, "sort_order ASC, code ASC")
	return list, total, err
}
func (r *Repository) UpdateEmploymentStatus(ctx context.Context, es *EmploymentStatus) error { return r.update(ctx, es) }
func (r *Repository) DeleteEmploymentStatus(ctx context.Context, id uuid.UUID) error { return r.softDelete(ctx, &EmploymentStatus{}, id) }

// ── Bank ──
func (r *Repository) CreateBank(ctx context.Context, b *Bank) error { return r.create(ctx, b) }
func (r *Repository) FindBankByID(ctx context.Context, id uuid.UUID) (*Bank, error) {
	var b Bank
	if err := r.findByID(ctx, &b, id, "banks"); err != nil {
		return nil, fmt.Errorf("bank not found: %w", err)
	}
	return &b, nil
}
func (r *Repository) FindAllBanks(ctx context.Context, page, perPage int) ([]Bank, int64, error) {
	var list []Bank
	total, err := r.findAll(ctx, &list, "banks", page, perPage, "sort_order ASC, code ASC")
	return list, total, err
}
func (r *Repository) UpdateBank(ctx context.Context, b *Bank) error { return r.update(ctx, b) }
func (r *Repository) DeleteBank(ctx context.Context, id uuid.UUID) error { return r.softDelete(ctx, &Bank{}, id) }

// ── Nationality ──
func (r *Repository) CreateNationality(ctx context.Context, n *Nationality) error { return r.create(ctx, n) }
func (r *Repository) FindNationalityByID(ctx context.Context, id uuid.UUID) (*Nationality, error) {
	var n Nationality
	if err := r.findByID(ctx, &n, id, "nationalities"); err != nil {
		return nil, fmt.Errorf("nationality not found: %w", err)
	}
	return &n, nil
}
func (r *Repository) FindAllNationalities(ctx context.Context, page, perPage int) ([]Nationality, int64, error) {
	var list []Nationality
	total, err := r.findAll(ctx, &list, "nationalities", page, perPage, "sort_order ASC, code ASC")
	return list, total, err
}
func (r *Repository) UpdateNationality(ctx context.Context, n *Nationality) error { return r.update(ctx, n) }
func (r *Repository) DeleteNationality(ctx context.Context, id uuid.UUID) error { return r.softDelete(ctx, &Nationality{}, id) }

// ── JobFamily ──
func (r *Repository) CreateJobFamily(ctx context.Context, jf *JobFamily) error { return r.create(ctx, jf) }
func (r *Repository) FindJobFamilyByID(ctx context.Context, id uuid.UUID) (*JobFamily, error) {
	var jf JobFamily
	if err := r.findByID(ctx, &jf, id, "job_families"); err != nil {
		return nil, fmt.Errorf("job family not found: %w", err)
	}
	return &jf, nil
}
func (r *Repository) FindAllJobFamilies(ctx context.Context, page, perPage int) ([]JobFamily, int64, error) {
	var list []JobFamily
	total, err := r.findAll(ctx, &list, "job_families", page, perPage, "sort_order ASC, code ASC")
	return list, total, err
}
func (r *Repository) UpdateJobFamily(ctx context.Context, jf *JobFamily) error { return r.update(ctx, jf) }
func (r *Repository) DeleteJobFamily(ctx context.Context, id uuid.UUID) error { return r.softDelete(ctx, &JobFamily{}, id) }

// ── Grading ──
func (r *Repository) CreateGrading(ctx context.Context, g *Grading) error { return r.create(ctx, g) }
func (r *Repository) FindGradingByID(ctx context.Context, id uuid.UUID) (*Grading, error) {
	var g Grading
	if err := r.findByID(ctx, &g, id, "gradings"); err != nil {
		return nil, fmt.Errorf("grading not found: %w", err)
	}
	return &g, nil
}
func (r *Repository) FindAllGradings(ctx context.Context, page, perPage int) ([]Grading, int64, error) {
	var list []Grading
	total, err := r.findAll(ctx, &list, "gradings", page, perPage, "sort_order ASC, code ASC")
	return list, total, err
}
func (r *Repository) UpdateGrading(ctx context.Context, g *Grading) error { return r.update(ctx, g) }
func (r *Repository) DeleteGrading(ctx context.Context, id uuid.UUID) error { return r.softDelete(ctx, &Grading{}, id) }

// ── SalaryGrade ──
func (r *Repository) CreateSalaryGrade(ctx context.Context, sg *SalaryGrade) error { return r.create(ctx, sg) }
func (r *Repository) FindSalaryGradeByID(ctx context.Context, id uuid.UUID) (*SalaryGrade, error) {
	var sg SalaryGrade
	if err := r.findByID(ctx, &sg, id, "salary_grades"); err != nil {
		return nil, fmt.Errorf("salary grade not found: %w", err)
	}
	return &sg, nil
}
func (r *Repository) FindAllSalaryGrades(ctx context.Context, page, perPage int) ([]SalaryGrade, int64, error) {
	var list []SalaryGrade
	total, err := r.findAll(ctx, &list, "salary_grades", page, perPage, "sort_order ASC, code ASC")
	return list, total, err
}
func (r *Repository) UpdateSalaryGrade(ctx context.Context, sg *SalaryGrade) error { return r.update(ctx, sg) }
func (r *Repository) DeleteSalaryGrade(ctx context.Context, id uuid.UUID) error { return r.softDelete(ctx, &SalaryGrade{}, id) }

// ── TER ──
func (r *Repository) CreateTER(ctx context.Context, t *TER) error { return r.create(ctx, t) }
func (r *Repository) FindTERByID(ctx context.Context, id uuid.UUID) (*TER, error) {
	var t TER
	if err := r.findByID(ctx, &t, id, "ters"); err != nil {
		return nil, fmt.Errorf("TER not found: %w", err)
	}
	return &t, nil
}
func (r *Repository) FindAllTERs(ctx context.Context, page, perPage int) ([]TER, int64, error) {
	var list []TER
	total, err := r.findAll(ctx, &list, "ters", page, perPage, "`group` ASC, bruto_min ASC")
	return list, total, err
}
func (r *Repository) UpdateTER(ctx context.Context, t *TER) error { return r.update(ctx, t) }
func (r *Repository) DeleteTER(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil { return err }
	return db.Delete(&TER{}, "id = ?", id).Error
}

// ── PTKP ──
func (r *Repository) CreatePTKP(ctx context.Context, p *PTKP) error { return r.create(ctx, p) }
func (r *Repository) FindPTKPByID(ctx context.Context, id uuid.UUID) (*PTKP, error) {
	var p PTKP
	if err := r.findByID(ctx, &p, id, "ptkps"); err != nil {
		return nil, fmt.Errorf("PTKP not found: %w", err)
	}
	return &p, nil
}
func (r *Repository) FindAllPTKPs(ctx context.Context, page, perPage int) ([]PTKP, int64, error) {
	var list []PTKP
	total, err := r.findAll(ctx, &list, "ptkps", page, perPage, "`group` ASC, ptkp ASC")
	return list, total, err
}
func (r *Repository) UpdatePTKP(ctx context.Context, p *PTKP) error { return r.update(ctx, p) }
func (r *Repository) DeletePTKP(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil { return err }
	return db.Delete(&PTKP{}, "id = ?", id).Error
}
