package jobmanagement

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/modules/setting"
)

// =========================================================================
// 9.1 Job Management Titles (Jenis Jabatan)
// =========================================================================
type JobTitle struct {
	ID           uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	Name         *string    `gorm:"type:varchar(100)" json:"name,omitempty"`
	Descriptions *string    `gorm:"type:text" json:"descriptions,omitempty"`
	Status       *int8      `gorm:"type:tinyint" json:"status,omitempty"`
	CreatedBy    *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy    *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Relations
	Subs []JobTitleSub `gorm:"foreignKey:JobManagementTitleID" json:"subs,omitempty"`
}

func (JobTitle) TableName() string { return "job_management_titles" }

func (t *JobTitle) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.2 Job Management Title Subs (Sub Jenis Jabatan)
// =========================================================================
type JobTitleSub struct {
	ID                    uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	JobManagementTitleID  *uuid.UUID `gorm:"type:char(36);index" json:"job_management_title_id,omitempty"`
	JobManagementTitleName *string   `gorm:"type:varchar(100)" json:"job_management_title_name,omitempty"`
	Name                  *string    `gorm:"type:varchar(100)" json:"name,omitempty"`
	Descriptions          *string    `gorm:"type:text" json:"descriptions,omitempty"`
	Status                *int8      `gorm:"type:tinyint" json:"status,omitempty"`
	CreatedBy             *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy             *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

func (JobTitleSub) TableName() string { return "job_management_title_subs" }

func (s *JobTitleSub) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.3 Job Management Values (Nilai Jabatan)
// =========================================================================
type JobValue struct {
	ID                      uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	JobManagementTitleSubID *uuid.UUID `gorm:"type:char(36);index" json:"job_management_title_sub_id,omitempty"`
	JobManagementTitleSubName *string  `gorm:"type:varchar(100)" json:"job_management_title_sub_name,omitempty"`
	Type                    string     `gorm:"type:varchar(255);not null" json:"type"`
	TypeGroup               *string    `gorm:"type:varchar(255)" json:"type_group,omitempty"`
	Level                   *int       `gorm:"type:int" json:"level,omitempty"`
	Descriptions            *string    `gorm:"type:text" json:"descriptions,omitempty"`
	DescriptionGroup        *string    `gorm:"type:text" json:"description_group,omitempty"`
	Note                    *string    `gorm:"type:text" json:"note,omitempty"`
	Sort                    *int       `gorm:"type:int" json:"sort,omitempty"`
	RefID                   *uuid.UUID `gorm:"type:char(36);index" json:"ref_id,omitempty"`
	RefType                 *string    `gorm:"type:varchar(100)" json:"ref_type,omitempty"`
	CreatedBy               *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy               *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

func (JobValue) TableName() string { return "job_management_values" }

func (v *JobValue) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.3b Job Management Value Clusters (Mapping type ↔ cluster kompetensi)
// =========================================================================
// Menyimpan daftar cluster kompetensi (dari tabel competencies) yang dipetakan
// ke sebuah tipe nilai jabatan (mis. 'technical'). Dipakai card Kompetensi
// Teknis sebagai filter cluster yang valid — diatur di halaman Mapping Job Value.
type JobManagementValueCluster struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Type      string    `gorm:"type:varchar(255);not null;uniqueIndex:uk_jmvc_type_cluster" json:"type"`
	Cluster   string    `gorm:"type:varchar(255);not null;uniqueIndex:uk_jmvc_type_cluster" json:"cluster"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (JobManagementValueCluster) TableName() string { return "job_management_value_clusters" }

func (c *JobManagementValueCluster) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.4 Job Management Objectives (Tujuan Jabatan)
// =========================================================================
type JobObjective struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature   string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode       string     `gorm:"type:varchar(20);not null" json:"full_code"`
	Objective      *string    `gorm:"type:text" json:"objective,omitempty"`
	CreatedBy      *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (JobObjective) TableName() string { return "job_management_objectives" }

func (o *JobObjective) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.5 Job Management Identifications (Identitas Jabatan)
// =========================================================================
type JobIdentification struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature   string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode       string     `gorm:"type:varchar(20);not null" json:"full_code"`
	GradingID      uuid.UUID  `gorm:"type:char(36);not null" json:"grading_id"`
	CreatedBy      *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (JobIdentification) TableName() string { return "job_management_identifications" }

func (i *JobIdentification) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.6 Job Management Responsibilities (Tanggung Jawab)
// =========================================================================
type JobResponsibility struct {
	ID                uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID    *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature      string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode          string     `gorm:"type:varchar(20);not null" json:"full_code"`
	MainTask          *string    `gorm:"type:text" json:"main_task,omitempty"`
	Activities        *string    `gorm:"type:text" json:"activities,omitempty"`
	Outputs           *string    `gorm:"type:text" json:"outputs,omitempty"`
	SuccessIndicators *string    `gorm:"type:text" json:"success_indicators,omitempty"`
	CreatedBy         *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy         *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func (JobResponsibility) TableName() string { return "job_management_responsibilities" }

func (r *JobResponsibility) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.7 Job Management Education Experiences (Pendidikan & Pengalaman)
// =========================================================================
// Terkait ke master module setting:
//   EducationID   → job_management_values.id (type=education) — Pendidikan
//   ExperienceID  → job_management_values.id (type=experience) — Pengalaman Kerja
//   Jurusan        → many-to-many via job_management_majors (education_majors.id)
//   Bidang Pekerjaan → many-to-many via job_management_job_family (job_families.id)
type JobEducationExperience struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature   string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode       string     `gorm:"type:varchar(20);not null" json:"full_code"`
	EducationID    *uuid.UUID `gorm:"type:char(36);index" json:"education_id,omitempty"`
	ExperienceID   *uuid.UUID `gorm:"type:char(36);index" json:"experience_id,omitempty"`
	CreatedBy      *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relasi untuk menampilkan nama:
	//   Education / Experience → job_management_values (di-preload dengan scope type)
	//   Jurusan & Bidang Pekerjaan → tabel pivot (multiple select)
	Education  *JobValue                 `gorm:"foreignKey:EducationID" json:"-"`
	Experience *JobValue                 `gorm:"foreignKey:ExperienceID" json:"-"`
	Majors     []JobManagementMajor      `gorm:"foreignKey:JobEducationExperienceID" json:"majors,omitempty"`
	JobFamilies []JobManagementJobFamily `gorm:"foreignKey:JobEducationExperienceID" json:"job_families,omitempty"`
}

func (JobEducationExperience) TableName() string { return "job_management_education_experiences" }

func (e *JobEducationExperience) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// JobManagementMajor — pivot table relasi job_management_education_experiences
// dengan master jurusan (education_majors). Satu record boleh punya banyak jurusan.
type JobManagementMajor struct {
	ID                           uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	JobEducationExperienceID     uuid.UUID  `gorm:"column:job_management_education_experience_id;type:char(36);not null;index" json:"job_management_education_experience_id,omitempty"`
	EducationMajorID             uuid.UUID  `gorm:"type:char(36);not null;index" json:"education_major_id,omitempty"`
	CreatedAt                    time.Time  `json:"created_at"`
	UpdatedAt                    time.Time  `json:"updated_at"`

	EducationMajor *setting.EducationMajor `gorm:"foreignKey:EducationMajorID" json:"-"`
}

func (JobManagementMajor) TableName() string { return "job_management_majors" }

func (m *JobManagementMajor) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// JobManagementJobFamily — pivot table relasi job_management_education_experiences
// dengan master bidang pekerjaan (job_families). Satu record boleh punya banyak bidang pekerjaan.
type JobManagementJobFamily struct {
	ID                           uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	JobEducationExperienceID     uuid.UUID  `gorm:"column:job_management_education_experience_id;type:char(36);not null;index" json:"job_management_education_experience_id,omitempty"`
	JobFamilyID                  uuid.UUID  `gorm:"type:char(36);not null;index" json:"job_family_id,omitempty"`
	CreatedAt                    time.Time  `json:"created_at"`
	UpdatedAt                    time.Time  `json:"updated_at"`

	JobFamily *setting.JobFamily `gorm:"foreignKey:JobFamilyID" json:"-"`
}

func (JobManagementJobFamily) TableName() string { return "job_management_job_family" }

func (jf *JobManagementJobFamily) BeforeCreate(tx *gorm.DB) error {
	if jf.ID == uuid.Nil {
		jf.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.8 Job Management HR Authorities (Kewenangan SDM)
// =========================================================================
type JobHRAuthority struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature   string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode       string     `gorm:"type:varchar(20);not null" json:"full_code"`
	Description    *string    `gorm:"type:text" json:"description,omitempty"`
	CreatedBy      *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (JobHRAuthority) TableName() string { return "job_management_hr_authorities" }

func (a *JobHRAuthority) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.9 Job Management Operational Authorities (Kewenangan Operasional)
// =========================================================================
type JobOperationalAuthority struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature   string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode       string     `gorm:"type:varchar(20);not null" json:"full_code"`
	Description    *string    `gorm:"type:text" json:"description,omitempty"`
	CreatedBy      *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (JobOperationalAuthority) TableName() string { return "job_management_operational_authorities" }

func (a *JobOperationalAuthority) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.10 Job Management Working Activities (Aktivitas Kerja)
// =========================================================================
type JobWorkingActivity struct {
	ID                  uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID      *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature        string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode            string     `gorm:"type:varchar(20);not null" json:"full_code"`
	JobManagementValueID *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_id,omitempty"`
	CreatedBy           *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy           *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func (JobWorkingActivity) TableName() string { return "job_management_working_activities" }

func (a *JobWorkingActivity) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.11 Job Management Working Risks (Risiko Kerja)
// =========================================================================
type JobWorkingRisk struct {
	ID                              uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID                  *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature                    string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode                        string     `gorm:"type:varchar(20);not null" json:"full_code"`
	JobManagementValueEnvironmentID *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_environment_id,omitempty"`
	JobManagementValueHazardID      *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_hazard_id,omitempty"`
	CreatedBy                       *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy                       *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedAt                       time.Time  `json:"updated_at"`
}

func (JobWorkingRisk) TableName() string { return "job_management_working_risks" }

func (r *JobWorkingRisk) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.12 Job Management Relationships (Hubungan Kerja)
// =========================================================================
type JobRelationship struct {
	ID                              uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID                  *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature                    string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode                        string     `gorm:"type:varchar(20);not null" json:"full_code"`
	JobManagementValueRelationshipID *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_relationship_id,omitempty"`
	JobManagementValueFrequencyID   *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_frequency_id,omitempty"`
	CreatedBy                       *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy                       *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedAt                       time.Time  `json:"updated_at"`
}

func (JobRelationship) TableName() string { return "job_management_relationships" }

func (r *JobRelationship) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.12b Job Management Relationship Details (Detail Hubungan Kerja)
// =========================================================================
// Rincian aktivitas per job_management_relationships.
//   organization_id → organizations.id (organisasi dengan summary yang sama —
//                     bagian dari "Hubungan Kerja")
//   activity        → deskripsi "Aktivitas Hubungan Dalam Rangka" (Ruang Lingkup)
type JobManagementRelationshipDetail struct {
	ID                        uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	JobManagementRelationshipID uuid.UUID `gorm:"column:job_management_relationship_id;type:char(36);not null;index" json:"job_management_relationship_id,omitempty"`
	OrganizationID            *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Activity                  *string    `gorm:"type:text" json:"activity,omitempty"`
	CreatedBy                 *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy                 *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`

	// Relasi untuk menampilkan nama organisasi pada detail
	Organization *OrganizationRef `gorm:"foreignKey:OrganizationID" json:"-"`
}

func (JobManagementRelationshipDetail) TableName() string { return "job_management_relationship_details" }

func (d *JobManagementRelationshipDetail) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// OrganizationRef — representasi minimal organisasi untuk keperluan relasi detail
// (menggunakan tabel organizations tanpa menimbulkan circular import).
type OrganizationRef struct {
	ID           string         `gorm:"column:id;type:char(36);primaryKey" json:"id"`
	Nomenclature string         `gorm:"column:nomenclature;type:varchar(255)" json:"nomenclature,omitempty"`
	FullCode     string         `gorm:"column:full_code;type:varchar(50)" json:"full_code,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (OrganizationRef) TableName() string { return "organizations" }

// =========================================================================
// 9.13 Job Management Subordinate Controls (Bawahan yang Dikendalikan)
// =========================================================================
type JobSubordinateControl struct {
	ID                  uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID      *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature        string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode            string     `gorm:"type:varchar(20);not null" json:"full_code"`
	JobManagementValueID *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_id,omitempty"`
	CreatedBy           *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy           *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func (JobSubordinateControl) TableName() string { return "job_management_subordinate_controls" }

func (c *JobSubordinateControl) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.14 Job Management Assets (Aset Jabatan)
// =========================================================================
type JobAsset struct {
	ID                          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID              *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature                string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode                    string     `gorm:"type:varchar(20);not null" json:"full_code"`
	JobManagementValueAssetID   *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_asset_id,omitempty"`
	JobManagementValueAuthorityID *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_authority_id,omitempty"`
	CreatedBy                   *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy                   *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at"`
}

func (JobAsset) TableName() string { return "job_management_assets" }

func (a *JobAsset) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.15 Job Management Financials (Keuangan Jabatan)
// =========================================================================
type JobFinancial struct {
	ID                          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID              *uuid.UUID `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	Nomenclature                string     `gorm:"type:varchar(50);not null" json:"nomenclature"`
	FullCode                    string     `gorm:"type:varchar(20);not null" json:"full_code"`
	IsAuthorized                bool       `gorm:"type:tinyint(1);not null;default:0" json:"is_authorized"`
	JobManagementValueCashID    *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_cash_id,omitempty"`
	JobManagementValueAuthorityID *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_authority_id,omitempty"`
	JobManagementValueImpactID  *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_impact_id,omitempty"`
	CreatedBy                   *uuid.UUID `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy                   *uuid.UUID `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at"`
}

func (JobFinancial) TableName() string { return "job_management_financials" }

func (f *JobFinancial) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.16 Job Management Potency Competencies (Kompetensi Potensi)
// =========================================================================
type JobPotencyCompetency struct {
	ID                  uuid.UUID   `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID      *uuid.UUID  `gorm:"type:char(36);index" json:"organization_id,omitempty"`
	JobManagementValueID *uuid.UUID `gorm:"type:char(36);index" json:"job_management_value_id,omitempty"`
	CompetencyID        *uuid.UUID  `gorm:"type:char(36);index" json:"competency_id,omitempty"`
	Weight              *float64    `gorm:"type:decimal(8,2)" json:"weight,omitempty"`
	CreatedBy           *uuid.UUID  `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy           *uuid.UUID  `gorm:"type:char(36)" json:"updated_by,omitempty"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
}

func (JobPotencyCompetency) TableName() string { return "job_management_potency_competencies" }

func (c *JobPotencyCompetency) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.17 Job Management Scores (Skor Jabatan)
// =========================================================================
type JobScore struct {
	ID                    uuid.UUID       `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID        *uuid.UUID      `gorm:"type:char(36);not null;uniqueIndex" json:"organization_id"`
	JobValueWithFinancial   uint64        `gorm:"type:bigint unsigned;not null;default:0" json:"job_value_with_financial"`
	JobValueWithoutFinancial uint64       `gorm:"type:bigint unsigned;not null;default:0" json:"job_value_without_financial"`
	HasFinancialAuthority  bool            `gorm:"type:tinyint(1);not null;default:0" json:"has_financial_authority"`
	Components             *string         `gorm:"type:json" json:"components,omitempty"`
	SubComponentPoints     *string         `gorm:"type:json" json:"sub_component_points,omitempty"`
	CalculatedAt           *time.Time      `json:"calculated_at,omitempty"`
	IsComplete             bool            `gorm:"type:tinyint(1);not null;default:0" json:"is_complete"`
	CompletedAt            *time.Time      `json:"completed_at,omitempty"`
	CreatedAt              time.Time       `json:"created_at"`
	UpdatedAt              time.Time       `json:"updated_at"`
}

func (JobScore) TableName() string { return "job_management_scores" }

func (s *JobScore) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// 9.18 Job Management Competency Groups (Bobot Kompetensi per Organisasi)
// =========================================================================
type JobCompetencyGroup struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID *uuid.UUID `gorm:"type:char(36);not null" json:"organization_id"`
	Category       string     `gorm:"type:varchar(20);not null" json:"category"`
	Weight         float64    `gorm:"type:decimal(8,2);not null;default:0" json:"weight"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (JobCompetencyGroup) TableName() string { return "job_management_competency_groups" }

func (g *JobCompetencyGroup) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}
