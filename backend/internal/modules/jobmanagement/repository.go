package jobmanagement

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

// =========================================================================
// Job Titles (9.1)
// =========================================================================

func (r *Repository) CreateJobTitle(ctx context.Context, t *JobTitle) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(t).Error
}

func (r *Repository) FindJobTitleByID(ctx context.Context, id uuid.UUID) (*JobTitle, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var t JobTitle
	if err := db.Preload("Subs").First(&t, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job title not found: %w", err)
	}
	return &t, nil
}

func (r *Repository) FindAllJobTitles(ctx context.Context, page, perPage int) ([]JobTitle, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var titles []JobTitle
	var total int64
	query := db.Model(&JobTitle{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&titles).Error; err != nil {
		return nil, 0, err
	}
	return titles, total, nil
}

func (r *Repository) UpdateJobTitle(ctx context.Context, t *JobTitle) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(t).Error
}

func (r *Repository) DeleteJobTitle(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobTitle{}).Error
}

// =========================================================================
// Job Title Subs (9.2)
// =========================================================================

func (r *Repository) CreateJobTitleSub(ctx context.Context, s *JobTitleSub) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(s).Error
}

func (r *Repository) FindJobTitleSubByID(ctx context.Context, id uuid.UUID) (*JobTitleSub, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var s JobTitleSub
	if err := db.First(&s, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job title sub not found: %w", err)
	}
	return &s, nil
}

func (r *Repository) FindJobTitleSubsByTitleID(ctx context.Context, titleID uuid.UUID) ([]JobTitleSub, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var subs []JobTitleSub
	if err := db.Where("job_management_title_id = ?", titleID).Order("created_at DESC").Find(&subs).Error; err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *Repository) UpdateJobTitleSub(ctx context.Context, s *JobTitleSub) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(s).Error
}

func (r *Repository) DeleteJobTitleSub(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobTitleSub{}).Error
}

// =========================================================================
// Job Values (9.3)
// =========================================================================

func (r *Repository) CreateJobValue(ctx context.Context, v *JobValue) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(v).Error
}

func (r *Repository) FindJobValueByID(ctx context.Context, id uuid.UUID) (*JobValue, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var v JobValue
	if err := db.First(&v, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job value not found: %w", err)
	}
	return &v, nil
}

func (r *Repository) FindAllJobValues(ctx context.Context, page, perPage int, valueType string) ([]JobValue, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var values []JobValue
	var total int64
	query := db.Model(&JobValue{})
	if valueType != "" {
		query = query.Where("type = ?", valueType)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("sort ASC, created_at DESC").Find(&values).Error; err != nil {
		return nil, 0, err
	}
	return values, total, nil
}

// FindAllJobValuesForTree memuat SEMUA job_management_values (tanpa pagination)
// untuk membangun tree type_group → type → options di service.
func (r *Repository) FindAllJobValuesForTree(ctx context.Context) ([]JobValue, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var values []JobValue
	if err := db.Model(&JobValue{}).
		Order("type ASC, sort ASC, created_at DESC").
		Find(&values).Error; err != nil {
		return nil, err
	}
	return values, nil
}

func (r *Repository) FindJobValuesByType(ctx context.Context, valueType string) ([]JobValue, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var values []JobValue
	if err := db.Where("type = ?", valueType).Order("sort ASC").Find(&values).Error; err != nil {
		return nil, err
	}
	return values, nil
}

func (r *Repository) UpdateJobValue(ctx context.Context, v *JobValue) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(v).Error
}

func (r *Repository) DeleteJobValue(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobValue{}).Error
}

// =========================================================================
// Job Value Clusters (9.3b) — mapping type ↔ cluster kompetensi
// =========================================================================

// FindJobValueClusters mengembalikan daftar cluster yang dipetakan ke sebuah tipe,
// diurutkan alfabetis.
func (r *Repository) FindJobValueClusters(ctx context.Context, valueType string) ([]JobManagementValueCluster, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var clusters []JobManagementValueCluster
	if err := db.Where("type = ?", valueType).Order("cluster ASC").Find(&clusters).Error; err != nil {
		return nil, err
	}
	return clusters, nil
}

// ReplaceJobValueClusters mengganti seluruh mapping untuk satu tipe dalam satu
// transaksi (delete semua + create ulang). Array kosong = hapus semua mapping.
func (r *Repository) ReplaceJobValueClusters(ctx context.Context, valueType string, clusters []string) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("type = ?", valueType).Delete(&JobManagementValueCluster{}).Error; err != nil {
			return err
		}
		if len(clusters) == 0 {
			return nil
		}
		records := make([]JobManagementValueCluster, 0, len(clusters))
		for _, c := range clusters {
			records = append(records, JobManagementValueCluster{Type: valueType, Cluster: c})
		}
		return tx.Create(&records).Error
	})
}

// =========================================================================
// Job Objectives (9.4)
// =========================================================================

func (r *Repository) CreateJobObjective(ctx context.Context, o *JobObjective) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(o).Error
}

func (r *Repository) FindJobObjectiveByID(ctx context.Context, id uuid.UUID) (*JobObjective, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var o JobObjective
	if err := db.First(&o, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job objective not found: %w", err)
	}
	return &o, nil
}

func (r *Repository) FindAllJobObjectives(ctx context.Context, page, perPage int, orgID *uuid.UUID) ([]JobObjective, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var objectives []JobObjective
	var total int64
	query := db.Model(&JobObjective{})
	if orgID != nil {
		query = query.Where("organization_id = ?", orgID.String())
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("full_code ASC").Find(&objectives).Error; err != nil {
		return nil, 0, err
	}
	return objectives, total, nil
}

func (r *Repository) UpdateJobObjective(ctx context.Context, o *JobObjective) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(o).Error
}

func (r *Repository) DeleteJobObjective(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobObjective{}).Error
}

// =========================================================================
// Job Identifications (9.5)
// =========================================================================

func (r *Repository) CreateJobIdentification(ctx context.Context, i *JobIdentification) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(i).Error
}

func (r *Repository) FindJobIdentificationByID(ctx context.Context, id uuid.UUID) (*JobIdentification, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var i JobIdentification
	if err := db.First(&i, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job identification not found: %w", err)
	}
	return &i, nil
}

func (r *Repository) FindAllJobIdentifications(ctx context.Context, page, perPage int, orgID *uuid.UUID) ([]JobIdentification, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var ids []JobIdentification
	var total int64
	query := db.Model(&JobIdentification{})
	if orgID != nil {
		query = query.Where("organization_id = ?", orgID.String())
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("full_code ASC").Find(&ids).Error; err != nil {
		return nil, 0, err
	}
	return ids, total, nil
}

func (r *Repository) UpdateJobIdentification(ctx context.Context, i *JobIdentification) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(i).Error
}

func (r *Repository) DeleteJobIdentification(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobIdentification{}).Error
}

// =========================================================================
// Job Responsibilities (9.6)
// =========================================================================

func (r *Repository) CreateJobResponsibility(ctx context.Context, resp *JobResponsibility) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(resp).Error
}

func (r *Repository) FindJobResponsibilityByID(ctx context.Context, id uuid.UUID) (*JobResponsibility, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var resp JobResponsibility
	if err := db.First(&resp, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job responsibility not found: %w", err)
	}
	return &resp, nil
}

func (r *Repository) FindAllJobResponsibilities(ctx context.Context, page, perPage int, orgID *uuid.UUID) ([]JobResponsibility, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var responsibilities []JobResponsibility
	var total int64
	query := db.Model(&JobResponsibility{})
	if orgID != nil {
		query = query.Where("organization_id = ?", orgID.String())
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("full_code ASC").Find(&responsibilities).Error; err != nil {
		return nil, 0, err
	}
	return responsibilities, total, nil
}

func (r *Repository) UpdateJobResponsibility(ctx context.Context, resp *JobResponsibility) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(resp).Error
}

func (r *Repository) DeleteJobResponsibility(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobResponsibility{}).Error
}

// =========================================================================
// Job Education Experiences (9.7)
// =========================================================================

func (r *Repository) CreateJobEducationExperience(ctx context.Context, e *JobEducationExperience) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		// Simpan entity utama dulu (tanpa relasi pivot)
		if err := tx.Omit("Majors", "JobFamilies").Create(e).Error; err != nil {
			return err
		}
		// Simpan relasi pivot jurusan & bidang pekerjaan
		for i := range e.Majors {
			e.Majors[i].JobEducationExperienceID = e.ID
		}
		for i := range e.JobFamilies {
			e.JobFamilies[i].JobEducationExperienceID = e.ID
		}
		if len(e.Majors) > 0 {
			if err := tx.Create(&e.Majors).Error; err != nil {
				return err
			}
		}
		if len(e.JobFamilies) > 0 {
			if err := tx.Create(&e.JobFamilies).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) FindJobEducationExperienceByID(ctx context.Context, id uuid.UUID) (*JobEducationExperience, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var e JobEducationExperience
	// Education / Experience → job_management_values dengan scope type
	if err := db.Preload("Education", "type = ?", "education").Preload("Experience", "type = ?", "experience").Preload("Majors.EducationMajor").Preload("JobFamilies.JobFamily").First(&e, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job education experience not found: %w", err)
	}
	return &e, nil
}

func (r *Repository) FindAllJobEducationExperiences(ctx context.Context, page, perPage int, orgID *uuid.UUID) ([]JobEducationExperience, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var experiences []JobEducationExperience
	var total int64
	query := db.Model(&JobEducationExperience{})
	if orgID != nil {
		query = query.Where("organization_id = ?", orgID.String())
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Preload("Education", "type = ?", "education").Preload("Experience", "type = ?", "experience").Preload("Majors.EducationMajor").Preload("JobFamilies.JobFamily").Offset(offset).Limit(perPage).Order("full_code ASC").Find(&experiences).Error; err != nil {
		return nil, 0, err
	}
	return experiences, total, nil
}

func (r *Repository) UpdateJobEducationExperience(ctx context.Context, e *JobEducationExperience) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		// Hapus relasi pivot lama, lalu simpan ulang sesuai request
		if err := tx.Where("job_management_education_experience_id = ?", e.ID.String()).Delete(&JobManagementMajor{}).Error; err != nil {
			return err
		}
		if err := tx.Where("job_management_education_experience_id = ?", e.ID.String()).Delete(&JobManagementJobFamily{}).Error; err != nil {
			return err
		}
		for i := range e.Majors {
			e.Majors[i].JobEducationExperienceID = e.ID
		}
		for i := range e.JobFamilies {
			e.JobFamilies[i].JobEducationExperienceID = e.ID
		}
		if len(e.Majors) > 0 {
			if err := tx.Create(&e.Majors).Error; err != nil {
				return err
			}
		}
		if len(e.JobFamilies) > 0 {
			if err := tx.Create(&e.JobFamilies).Error; err != nil {
				return err
			}
		}
		// Update entity utama (tanpa menyentuh relasi pivot)
		// Education/Experience wajib di-Omit: keduanya belongs-to yang di-Preload
		// (FindJobEducationExperienceByID). Tanpa Omit, GORM akan mengembalikan FK
		// (education_id/experience_id) ke nilai association lama saat Save, sehingga
		// perubahan FK tidak pernah tersimpan.
		return tx.Omit("Majors", "JobFamilies", "Education", "Experience").Save(e).Error
	})
}

func (r *Repository) DeleteJobEducationExperience(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		// Hapus relasi pivot dulu (atau biarkan ON DELETE CASCADE di DB)
		if err := tx.Where("job_management_education_experience_id = ?", id.String()).Delete(&JobManagementMajor{}).Error; err != nil {
			return err
		}
		if err := tx.Where("job_management_education_experience_id = ?", id.String()).Delete(&JobManagementJobFamily{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&JobEducationExperience{}).Error
	})
}

// =========================================================================
// Job HR Authorities (9.8)
// =========================================================================

func (r *Repository) CreateJobHRAuthority(ctx context.Context, a *JobHRAuthority) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(a).Error
}

func (r *Repository) FindJobHRAuthorityByID(ctx context.Context, id uuid.UUID) (*JobHRAuthority, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var a JobHRAuthority
	if err := db.First(&a, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job HR authority not found: %w", err)
	}
	return &a, nil
}

func (r *Repository) FindAllJobHRAuthorities(ctx context.Context, page, perPage int, orgID *uuid.UUID) ([]JobHRAuthority, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var authorities []JobHRAuthority
	var total int64
	query := db.Model(&JobHRAuthority{})
	if orgID != nil {
		query = query.Where("organization_id = ?", orgID.String())
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("full_code ASC").Find(&authorities).Error; err != nil {
		return nil, 0, err
	}
	return authorities, total, nil
}

func (r *Repository) UpdateJobHRAuthority(ctx context.Context, a *JobHRAuthority) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(a).Error
}

func (r *Repository) DeleteJobHRAuthority(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobHRAuthority{}).Error
}

// =========================================================================
// Job Operational Authorities (9.9)
// =========================================================================

func (r *Repository) CreateJobOperationalAuthority(ctx context.Context, a *JobOperationalAuthority) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(a).Error
}

func (r *Repository) FindJobOperationalAuthorityByID(ctx context.Context, id uuid.UUID) (*JobOperationalAuthority, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var a JobOperationalAuthority
	if err := db.First(&a, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job operational authority not found: %w", err)
	}
	return &a, nil
}

func (r *Repository) FindAllJobOperationalAuthorities(ctx context.Context, page, perPage int, orgID *uuid.UUID) ([]JobOperationalAuthority, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var authorities []JobOperationalAuthority
	var total int64
	query := db.Model(&JobOperationalAuthority{})
	if orgID != nil {
		query = query.Where("organization_id = ?", orgID.String())
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("full_code ASC").Find(&authorities).Error; err != nil {
		return nil, 0, err
	}
	return authorities, total, nil
}

func (r *Repository) UpdateJobOperationalAuthority(ctx context.Context, a *JobOperationalAuthority) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(a).Error
}

func (r *Repository) DeleteJobOperationalAuthority(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobOperationalAuthority{}).Error
}

// =========================================================================
// Job Working Activities (9.10)
// =========================================================================

func (r *Repository) CreateJobWorkingActivity(ctx context.Context, a *JobWorkingActivity) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(a).Error
}

func (r *Repository) FindJobWorkingActivityByID(ctx context.Context, id uuid.UUID) (*JobWorkingActivity, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var a JobWorkingActivity
	if err := db.First(&a, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job working activity not found: %w", err)
	}
	return &a, nil
}

func (r *Repository) FindAllJobWorkingActivities(ctx context.Context, page, perPage int, orgID *uuid.UUID) ([]JobWorkingActivity, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var activities []JobWorkingActivity
	var total int64
	query := db.Model(&JobWorkingActivity{})
	if orgID != nil {
		query = query.Where("organization_id = ?", orgID.String())
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("full_code ASC").Find(&activities).Error; err != nil {
		return nil, 0, err
	}
	return activities, total, nil
}

func (r *Repository) UpdateJobWorkingActivity(ctx context.Context, a *JobWorkingActivity) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(a).Error
}

func (r *Repository) DeleteJobWorkingActivity(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobWorkingActivity{}).Error
}

// =========================================================================
// Job Working Risks (9.11)
// =========================================================================

func (r *Repository) CreateJobWorkingRisk(ctx context.Context, risk *JobWorkingRisk) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(risk).Error
}

func (r *Repository) FindJobWorkingRiskByID(ctx context.Context, id uuid.UUID) (*JobWorkingRisk, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var risk JobWorkingRisk
	if err := db.First(&risk, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job working risk not found: %w", err)
	}
	return &risk, nil
}

func (r *Repository) FindAllJobWorkingRisks(ctx context.Context, page, perPage int, orgID *uuid.UUID) ([]JobWorkingRisk, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var risks []JobWorkingRisk
	var total int64
	query := db.Model(&JobWorkingRisk{})
	if orgID != nil {
		query = query.Where("organization_id = ?", orgID.String())
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("full_code ASC").Find(&risks).Error; err != nil {
		return nil, 0, err
	}
	return risks, total, nil
}

func (r *Repository) UpdateJobWorkingRisk(ctx context.Context, risk *JobWorkingRisk) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(risk).Error
}

func (r *Repository) DeleteJobWorkingRisk(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobWorkingRisk{}).Error
}

// =========================================================================
// Job Relationships (9.12)
// =========================================================================

func (r *Repository) CreateJobRelationship(ctx context.Context, rel *JobRelationship) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(rel).Error
}

func (r *Repository) FindJobRelationshipByID(ctx context.Context, id uuid.UUID) (*JobRelationship, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rel JobRelationship
	if err := db.First(&rel, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job relationship not found: %w", err)
	}
	return &rel, nil
}

func (r *Repository) FindAllJobRelationships(ctx context.Context, page, perPage int, orgID *uuid.UUID) ([]JobRelationship, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var relationships []JobRelationship
	var total int64
	query := db.Model(&JobRelationship{})
	if orgID != nil {
		query = query.Where("organization_id = ?", orgID.String())
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("full_code ASC").Find(&relationships).Error; err != nil {
		return nil, 0, err
	}
	return relationships, total, nil
}

func (r *Repository) UpdateJobRelationship(ctx context.Context, rel *JobRelationship) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(rel).Error
}

func (r *Repository) DeleteJobRelationship(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobRelationship{}).Error
}

// =========================================================================
// Job Relationship Details (9.12b)
// =========================================================================

func (r *Repository) CreateJobRelationshipDetail(ctx context.Context, d *JobManagementRelationshipDetail) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	// Omit Organization (hanya representasi minimal untuk preload) agar GORM
	// tidak ikut menulis/membuat ulang baris di tabel organizations.
	return db.Omit("Organization").Create(d).Error
}

func (r *Repository) FindJobRelationshipDetailByID(ctx context.Context, id uuid.UUID) (*JobManagementRelationshipDetail, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var d JobManagementRelationshipDetail
	if err := db.Preload("Organization").First(&d, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job relationship detail not found: %w", err)
	}
	return &d, nil
}

func (r *Repository) FindAllJobRelationshipDetails(ctx context.Context, relationshipID uuid.UUID) ([]JobManagementRelationshipDetail, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var details []JobManagementRelationshipDetail
	if err := db.Preload("Organization").Where("job_management_relationship_id = ?", relationshipID.String()).Order("created_at ASC").Find(&details).Error; err != nil {
		return nil, err
	}
	return details, nil
}

func (r *Repository) UpdateJobRelationshipDetail(ctx context.Context, d *JobManagementRelationshipDetail) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	// Omit Organization — field ini hanya representasi minimal (OrganizationRef) hasil
	// preload untuk menampilkan nama org. db.Save() tanpa Omit akan ikut meng-update
	// tabel organizations dan memicu error MySQL 1364 (kolom code tanpa default).
	return db.Omit("Organization").Save(d).Error
}

func (r *Repository) DeleteJobRelationshipDetail(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobManagementRelationshipDetail{}).Error
}

// =========================================================================
// Job Subordinate Controls (9.13)
// =========================================================================

func (r *Repository) CreateJobSubordinateControl(ctx context.Context, c *JobSubordinateControl) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(c).Error
}

func (r *Repository) FindJobSubordinateControlByID(ctx context.Context, id uuid.UUID) (*JobSubordinateControl, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var c JobSubordinateControl
	if err := db.First(&c, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job subordinate control not found: %w", err)
	}
	return &c, nil
}

func (r *Repository) FindAllJobSubordinateControls(ctx context.Context, page, perPage int, orgID *uuid.UUID) ([]JobSubordinateControl, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var controls []JobSubordinateControl
	var total int64
	query := db.Model(&JobSubordinateControl{})
	if orgID != nil {
		query = query.Where("organization_id = ?", orgID.String())
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("full_code ASC").Find(&controls).Error; err != nil {
		return nil, 0, err
	}
	return controls, total, nil
}

func (r *Repository) UpdateJobSubordinateControl(ctx context.Context, c *JobSubordinateControl) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(c).Error
}

func (r *Repository) DeleteJobSubordinateControl(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobSubordinateControl{}).Error
}

// =========================================================================
// Job Assets (9.14)
// =========================================================================

func (r *Repository) CreateJobAsset(ctx context.Context, a *JobAsset) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(a).Error
}

func (r *Repository) FindJobAssetByID(ctx context.Context, id uuid.UUID) (*JobAsset, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var a JobAsset
	if err := db.First(&a, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job asset not found: %w", err)
	}
	return &a, nil
}

func (r *Repository) FindAllJobAssets(ctx context.Context, page, perPage int, orgID *uuid.UUID) ([]JobAsset, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var assets []JobAsset
	var total int64
	query := db.Model(&JobAsset{})
	if orgID != nil {
		query = query.Where("organization_id = ?", orgID.String())
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("full_code ASC").Find(&assets).Error; err != nil {
		return nil, 0, err
	}
	return assets, total, nil
}

func (r *Repository) UpdateJobAsset(ctx context.Context, a *JobAsset) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(a).Error
}

func (r *Repository) DeleteJobAsset(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobAsset{}).Error
}

// =========================================================================
// Job Financials (9.15)
// =========================================================================

func (r *Repository) CreateJobFinancial(ctx context.Context, f *JobFinancial) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(f).Error
}

func (r *Repository) FindJobFinancialByID(ctx context.Context, id uuid.UUID) (*JobFinancial, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var f JobFinancial
	if err := db.First(&f, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job financial not found: %w", err)
	}
	return &f, nil
}

func (r *Repository) FindAllJobFinancials(ctx context.Context, page, perPage int, orgID *uuid.UUID) ([]JobFinancial, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var financials []JobFinancial
	var total int64
	query := db.Model(&JobFinancial{})
	if orgID != nil {
		query = query.Where("organization_id = ?", orgID.String())
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("full_code ASC").Find(&financials).Error; err != nil {
		return nil, 0, err
	}
	return financials, total, nil
}

func (r *Repository) UpdateJobFinancial(ctx context.Context, f *JobFinancial) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(f).Error
}

func (r *Repository) DeleteJobFinancial(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobFinancial{}).Error
}

// =========================================================================
// Job Potency Competencies (9.16)
// =========================================================================

func (r *Repository) CreateJobPotencyCompetency(ctx context.Context, c *JobPotencyCompetency) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(c).Error
}

func (r *Repository) FindJobPotencyCompetencyByID(ctx context.Context, id uuid.UUID) (*JobPotencyCompetency, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var c JobPotencyCompetency
	if err := db.Preload("Competency").Preload("JobManagementValue").First(&c, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job potency competency not found: %w", err)
	}
	return &c, nil
}

func (r *Repository) FindAllJobPotencyCompetencies(ctx context.Context, page, perPage int, orgID *uuid.UUID) ([]JobPotencyCompetency, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var competencies []JobPotencyCompetency
	var total int64
	query := db.Model(&JobPotencyCompetency{})
	if orgID != nil {
		query = query.Where("organization_id = ?", orgID.String())
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Preload("Competency").Preload("JobManagementValue").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&competencies).Error; err != nil {
		return nil, 0, err
	}
	return competencies, total, nil
}

func (r *Repository) UpdateJobPotencyCompetency(ctx context.Context, c *JobPotencyCompetency) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(c).Error
}

func (r *Repository) DeleteJobPotencyCompetency(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobPotencyCompetency{}).Error
}

// =========================================================================
// Job Scores (9.17)
// =========================================================================

func (r *Repository) UpsertJobScore(ctx context.Context, s *JobScore) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if s.OrganizationID == nil {
		return fmt.Errorf("organization_id is required for job score")
	}
	// Find existing by organization_id
	var existing JobScore
	result := db.Where("organization_id = ?", s.OrganizationID.String()).First(&existing)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return db.Create(s).Error
		}
		return result.Error
	}
	// Update existing
	s.ID = existing.ID
	s.CreatedAt = existing.CreatedAt
	return db.Save(s).Error
}

func (r *Repository) FindJobScoreByOrganizationID(ctx context.Context, orgID uuid.UUID) (*JobScore, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var s JobScore
	if err := db.Where("organization_id = ?", orgID.String()).First(&s).Error; err != nil {
		return nil, fmt.Errorf("job score not found: %w", err)
	}
	return &s, nil
}

func (r *Repository) FindAllJobScores(ctx context.Context, page, perPage int) ([]JobScore, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var scores []JobScore
	var total int64
	query := db.Model(&JobScore{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&scores).Error; err != nil {
		return nil, 0, err
	}
	return scores, total, nil
}

// =========================================================================
// Job Competency Groups (9.18)
// =========================================================================

func (r *Repository) CreateJobCompetencyGroup(ctx context.Context, g *JobCompetencyGroup) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(g).Error
}

func (r *Repository) FindJobCompetencyGroupByID(ctx context.Context, id uuid.UUID) (*JobCompetencyGroup, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var g JobCompetencyGroup
	if err := db.First(&g, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("job competency group not found: %w", err)
	}
	return &g, nil
}

func (r *Repository) FindJobCompetencyGroupsByOrganization(ctx context.Context, orgID uuid.UUID) ([]JobCompetencyGroup, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var groups []JobCompetencyGroup
	if err := db.Where("organization_id = ?", orgID.String()).Order("category ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *Repository) UpdateJobCompetencyGroup(ctx context.Context, g *JobCompetencyGroup) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(g).Error
}

func (r *Repository) DeleteJobCompetencyGroup(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&JobCompetencyGroup{}).Error
}
