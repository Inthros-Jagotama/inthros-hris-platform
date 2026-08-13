package recruitment

import "time"

// =========================================================================
// Job Requisition DTOs
// =========================================================================

type CreateRequisitionRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	Title          string `json:"title" binding:"required,max=255"`
	// G-2: requisition_number opsional — bila kosong auto-generated
	// REQ-YYYYMM-XXXX oleh service.
	RequisitionNumber    string  `json:"requisition_number" binding:"omitempty,max=50"`
	Priority             *string `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH URGENT"`
	PositionID           *string `json:"position_id"`
	Department           string  `json:"department" binding:"omitempty,max=150"`
	EmploymentType       string  `json:"employment_type" binding:"omitempty,max=50"`
	Location             string  `json:"location" binding:"omitempty,max=255"`
	MinSalary            float64 `json:"min_salary"`
	MaxSalary            float64 `json:"max_salary"`
	Description          string  `json:"description"`
	Requirements         string  `json:"requirements"`
	Responsibilities     string  `json:"responsibilities"`
	SlotsAvailable       *int    `json:"slots_available"`
	RequestedBy          *string `json:"requested_by"`
	TargetStartDate      *string `json:"target_start_date"`
	ReasonType           *string `json:"reason_type" binding:"omitempty,oneof=NEW_POSITION REPLACEMENT EXPANSION WORKFORCE_GAP SUCCESSION_GAP"`
	WorkforceGapID       *string `json:"workforce_gap_id"`
	WorkforcePlanID      *string `json:"workforce_plan_id"`
	SuccessionPositionID *string `json:"succession_position_id"`
}

type UpdateRequisitionRequest struct {
	Title             *string  `json:"title" binding:"omitempty,max=255"`
	RequisitionNumber *string  `json:"requisition_number" binding:"omitempty,max=50"`
	Priority          *string  `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH URGENT"`
	PositionID        *string  `json:"position_id"`
	Department        *string  `json:"department" binding:"omitempty,max=150"`
	EmploymentType    *string  `json:"employment_type" binding:"omitempty,max=50"`
	Location          *string  `json:"location" binding:"omitempty,max=255"`
	MinSalary         *float64 `json:"min_salary"`
	MaxSalary         *float64 `json:"max_salary"`
	Description       *string  `json:"description"`
	Requirements      *string  `json:"requirements"`
	Responsibilities  *string  `json:"responsibilities"`
	SlotsAvailable    *int     `json:"slots_available"`
	Status            *string  `json:"status" binding:"omitempty,oneof=DRAFT OPEN IN_PROGRESS FILLED CANCELLED"`
	// Catatan G-1: SUBMITTED/REJECTED tidak boleh di-set manual via update —
	// keduanya hanya dihasilkan alur approval (submit + push-callback).
	TargetStartDate      *string `json:"target_start_date"`
	ReasonType           *string `json:"reason_type" binding:"omitempty,oneof=NEW_POSITION REPLACEMENT EXPANSION WORKFORCE_GAP SUCCESSION_GAP"`
	WorkforceGapID       *string `json:"workforce_gap_id"`
	WorkforcePlanID      *string `json:"workforce_plan_id"`
	SuccessionPositionID *string `json:"succession_position_id"`
}

// SubmitRequisitionRequest mengirim requisition draft ke Central Approval
// (plan G-1). FlowID opsional — bila kosong, flow aktif untuk modul
// "recruitment" di-auto-resolve (pola employeemovement G-3).
type SubmitRequisitionRequest struct {
	FlowID *string `json:"flow_id"`
}

type RequisitionResponse struct {
	ID                   string    `json:"id"`
	OrganizationID       string    `json:"organization_id"`
	Title                string    `json:"title"`
	RequisitionNumber    string    `json:"requisition_number,omitempty"`
	Priority             string    `json:"priority"`
	PositionID           string    `json:"position_id,omitempty"`
	OpenedAt             *int64    `json:"opened_at,omitempty"`
	Department           string    `json:"department,omitempty"`
	EmploymentType       string    `json:"employment_type,omitempty"`
	Location             string    `json:"location,omitempty"`
	MinSalary            float64   `json:"min_salary"`
	MaxSalary            float64   `json:"max_salary"`
	Description          string    `json:"description"`
	Requirements         string    `json:"requirements"`
	Responsibilities     string    `json:"responsibilities"`
	SlotsAvailable       int       `json:"slots_available"`
	SlotsFilled          int       `json:"slots_filled"`
	Status               string    `json:"status"`
	RequestedBy          string    `json:"requested_by,omitempty"`
	ApprovedBy           string    `json:"approved_by,omitempty"`
	ApprovalInstanceID   string    `json:"approval_instance_id,omitempty"`
	ReasonType           string    `json:"reason_type,omitempty"`
	WorkforceGapID       string    `json:"workforce_gap_id,omitempty"`
	WorkforcePlanID      string    `json:"workforce_plan_id,omitempty"`
	SuccessionPositionID string    `json:"succession_position_id,omitempty"`
	TargetStartDate      string    `json:"target_start_date,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// =========================================================================
// Job Offer DTOs (G-3)
// =========================================================================

type CreateOfferRequest struct {
	ApplicationID  string  `json:"application_id" binding:"required"`
	EmploymentType string  `json:"employment_type" binding:"omitempty,max=50"`
	Salary         float64 `json:"salary"`
	Allowances     float64 `json:"allowances"`
	Benefits       string  `json:"benefits"`
	StartDate      string  `json:"start_date" binding:"omitempty,max=10"`
	ExpiryDate     string  `json:"expiry_date" binding:"omitempty,max=10"`
}

type UpdateOfferRequest struct {
	EmploymentType *string  `json:"employment_type" binding:"omitempty,max=50"`
	Salary         *float64 `json:"salary"`
	Allowances     *float64 `json:"allowances"`
	Benefits       *string  `json:"benefits"`
	StartDate      *string  `json:"start_date" binding:"omitempty,max=10"`
	ExpiryDate     *string  `json:"expiry_date" binding:"omitempty,max=10"`
}

// SubmitOfferRequest mengirim offer draft ke Central Approval (G-3). FlowID
// opsional — bila kosong, flow aktif modul "recruitment_offer" di-auto-resolve.
type SubmitOfferRequest struct {
	FlowID *string `json:"flow_id"`
}

type OfferResponse struct {
	ID                 string    `json:"id"`
	ApplicationID      string    `json:"application_id"`
	OfferNumber        string    `json:"offer_number,omitempty"`
	EmploymentType     string    `json:"employment_type"`
	Salary             float64   `json:"salary"`
	Allowances         float64   `json:"allowances"`
	Benefits           string    `json:"benefits"`
	StartDate          string    `json:"start_date"`
	ExpiryDate         string    `json:"expiry_date"`
	Status             string    `json:"status"`
	ApprovalInstanceID string    `json:"approval_instance_id,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// =========================================================================
// Candidate DTOs
// =========================================================================

type CreateCandidateRequest struct {
	FirstName      string  `json:"first_name" binding:"required,max=100"`
	LastName       string  `json:"last_name" binding:"required,max=100"`
	Email          string  `json:"email" binding:"required,email,max=255"`
	Phone          string  `json:"phone" binding:"omitempty,max=50"`
	Address        string  `json:"address"`
	CurrentCompany *string `json:"current_company" binding:"omitempty,max=255"`
	CurrentTitle   *string `json:"current_title" binding:"omitempty,max=255"`
	ResumeURL      *string `json:"resume_url"`
	PortfolioURL   *string `json:"portfolio_url"`
	LinkedInURL    *string `json:"linkedin_url"`
	Source         *string `json:"source" binding:"omitempty,max=50"`
	Notes          string  `json:"notes"`
	// G-4: jenis kandidat (EXTERNAL default | INTERNAL) + referensi employee
	// untuk kandidat internal (hasil seleksi → Employee Movement).
	CandidateType   *string `json:"candidate_type" binding:"omitempty,oneof=EXTERNAL INTERNAL"`
	EmployeeID      *string `json:"employee_id" binding:"omitempty"`
	CandidateNumber *string `json:"candidate_number" binding:"omitempty,max=50"`
}

type UpdateCandidateRequest struct {
	FirstName      *string `json:"first_name" binding:"omitempty,max=100"`
	LastName       *string `json:"last_name" binding:"omitempty,max=100"`
	Email          *string `json:"email" binding:"omitempty,email,max=255"`
	Phone          *string `json:"phone" binding:"omitempty,max=50"`
	Address        *string `json:"address"`
	CurrentCompany *string `json:"current_company" binding:"omitempty,max=255"`
	CurrentTitle   *string `json:"current_title" binding:"omitempty,max=255"`
	ResumeURL      *string `json:"resume_url"`
	PortfolioURL   *string `json:"portfolio_url"`
	LinkedInURL    *string `json:"linkedin_url"`
	Source         *string `json:"source" binding:"omitempty,max=50"`
	Notes          *string `json:"notes"`
	// G-4: jenis kandidat + referensi employee (internal hire).
	CandidateType   *string `json:"candidate_type" binding:"omitempty,oneof=EXTERNAL INTERNAL"`
	EmployeeID      *string `json:"employee_id" binding:"omitempty"`
	CandidateNumber *string `json:"candidate_number" binding:"omitempty,max=50"`
}

type CandidateResponse struct {
	ID             string `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone,omitempty"`
	Address        string `json:"address,omitempty"`
	CurrentCompany string `json:"current_company,omitempty"`
	CurrentTitle   string `json:"current_title,omitempty"`
	ResumeURL      string `json:"resume_url,omitempty"`
	PortfolioURL   string `json:"portfolio_url,omitempty"`
	LinkedInURL    string `json:"linkedin_url,omitempty"`
	Source         string `json:"source"`
	Notes          string `json:"notes"`
	// G-4: jenis kandidat + referensi employee (internal hire).
	CandidateType   string    `json:"candidate_type"`
	EmployeeID      string    `json:"employee_id,omitempty"`
	CandidateNumber string    `json:"candidate_number,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// =========================================================================
// Candidate Education DTOs (G-6)
// =========================================================================

type CreateCandidateEducationRequest struct {
	EducationID      *string  `json:"education_id" binding:"omitempty"`
	InstitutionName  string   `json:"institution_name" binding:"required,max=255"`
	EducationMajorID *string  `json:"education_major_id" binding:"omitempty"`
	Major            *string  `json:"major" binding:"omitempty,max=255"`
	GPA              *float64 `json:"gpa"`
	StartYear        *int     `json:"start_year"`
	EndYear          *int     `json:"end_year"`
	IsHighest        bool     `json:"is_highest"`
	Notes            string   `json:"notes"`
}

type UpdateCandidateEducationRequest struct {
	EducationID      *string  `json:"education_id"`
	InstitutionName  *string  `json:"institution_name" binding:"omitempty,max=255"`
	EducationMajorID *string  `json:"education_major_id"`
	Major            *string  `json:"major" binding:"omitempty,max=255"`
	GPA              *float64 `json:"gpa"`
	StartYear        *int     `json:"start_year"`
	EndYear          *int     `json:"end_year"`
	IsHighest        *bool    `json:"is_highest"`
	Notes            *string  `json:"notes"`
}

type CandidateEducationResponse struct {
	ID               string  `json:"id"`
	CandidateID      string  `json:"candidate_id"`
	EducationID      string  `json:"education_id,omitempty"`
	EducationMajorID string  `json:"education_major_id,omitempty"`
	MajorName        string  `json:"major_name,omitempty"`
	InstitutionName  string  `json:"institution_name"`
	Major            string  `json:"major,omitempty"`
	GPA              float64 `json:"gpa,omitempty"`
	StartYear        int     `json:"start_year,omitempty"`
	EndYear          int     `json:"end_year,omitempty"`
	IsHighest        bool    `json:"is_highest"`
	Notes            string  `json:"notes,omitempty"`
}

// =========================================================================
// Candidate Work Experience DTOs (G-6)
// =========================================================================

type CreateCandidateWorkExperienceRequest struct {
	CompanyName    string  `json:"company_name" binding:"required,max=255"`
	JobTitle       string  `json:"job_title" binding:"required,max=255"`
	EmploymentType *string `json:"employment_type" binding:"omitempty,max=50"`
	StartDate      string  `json:"start_date" binding:"required,max=10"`
	EndDate        *string `json:"end_date" binding:"omitempty,max=10"`
	IsCurrent      bool    `json:"is_current"`
	Description    string  `json:"description"`
}

type UpdateCandidateWorkExperienceRequest struct {
	CompanyName    *string `json:"company_name" binding:"omitempty,max=255"`
	JobTitle       *string `json:"job_title" binding:"omitempty,max=255"`
	EmploymentType *string `json:"employment_type" binding:"omitempty,max=50"`
	StartDate      *string `json:"start_date" binding:"omitempty,max=10"`
	EndDate        *string `json:"end_date" binding:"omitempty,max=10"`
	IsCurrent      *bool   `json:"is_current"`
	Description    *string `json:"description"`
}

type CandidateWorkExperienceResponse struct {
	ID             string `json:"id"`
	CandidateID    string `json:"candidate_id"`
	CompanyName    string `json:"company_name"`
	JobTitle       string `json:"job_title"`
	EmploymentType string `json:"employment_type,omitempty"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date,omitempty"`
	IsCurrent      bool   `json:"is_current"`
	Description    string `json:"description,omitempty"`
}

// =========================================================================
// Candidate Skill DTOs (G-6)
// =========================================================================

type CreateCandidateSkillRequest struct {
	CompetencyID string `json:"competency_id" binding:"required"`
	Level        *int   `json:"level"`
	Notes        string `json:"notes"`
}

type UpdateCandidateSkillRequest struct {
	Level *int    `json:"level"`
	Notes *string `json:"notes"`
}

type CandidateSkillResponse struct {
	ID             string `json:"id"`
	CandidateID    string `json:"candidate_id"`
	CompetencyID   string `json:"competency_id"`
	CompetencyName string `json:"competency_name,omitempty"`
	Level          int    `json:"level,omitempty"`
	Notes          string `json:"notes,omitempty"`
}

// =========================================================================
// Candidate Certification DTOs (G-6)
// =========================================================================

type CreateCandidateCertificationRequest struct {
	Name                string  `json:"name" binding:"required,max=255"`
	IssuingOrganization *string `json:"issuing_organization" binding:"omitempty,max=255"`
	IssueDate           *string `json:"issue_date" binding:"omitempty,max=10"`
	ExpiryDate          *string `json:"expiry_date" binding:"omitempty,max=10"`
	CredentialURL       *string `json:"credential_url"`
	Notes               string  `json:"notes"`
}

type UpdateCandidateCertificationRequest struct {
	Name                *string `json:"name" binding:"omitempty,max=255"`
	IssuingOrganization *string `json:"issuing_organization" binding:"omitempty,max=255"`
	IssueDate           *string `json:"issue_date" binding:"omitempty,max=10"`
	ExpiryDate          *string `json:"expiry_date" binding:"omitempty,max=10"`
	CredentialURL       *string `json:"credential_url"`
	Notes               *string `json:"notes"`
}

type CandidateCertificationResponse struct {
	ID                  string `json:"id"`
	CandidateID         string `json:"candidate_id"`
	Name                string `json:"name"`
	IssuingOrganization string `json:"issuing_organization,omitempty"`
	IssueDate           string `json:"issue_date,omitempty"`
	ExpiryDate          string `json:"expiry_date,omitempty"`
	CredentialURL       string `json:"credential_url,omitempty"`
	Notes               string `json:"notes,omitempty"`
}

// =========================================================================
// Candidate Document DTOs (G-6)
// =========================================================================

type CreateCandidateDocumentRequest struct {
	DocumentType string `json:"document_type" binding:"omitempty,oneof=RESUME COVER_LETTER CERTIFICATE PORTFOLIO IDENTITY OTHER"`
	Name         string `json:"name" binding:"required,max=255"`
	FileURL      string `json:"file_url" binding:"required"`
	Notes        string `json:"notes"`
}

type UpdateCandidateDocumentRequest struct {
	DocumentType *string `json:"document_type" binding:"omitempty,oneof=RESUME COVER_LETTER CERTIFICATE PORTFOLIO IDENTITY OTHER"`
	Name         *string `json:"name" binding:"omitempty,max=255"`
	FileURL      *string `json:"file_url" binding:"omitempty"`
	Notes        *string `json:"notes"`
}

type CandidateDocumentResponse struct {
	ID           string `json:"id"`
	CandidateID  string `json:"candidate_id"`
	DocumentType string `json:"document_type"`
	Name         string `json:"name"`
	FileURL      string `json:"file_url"`
	Notes        string `json:"notes,omitempty"`
}

// =========================================================================
// Candidate Consent DTOs (G-6) — append-only, no Update
// =========================================================================

type CreateCandidateConsentRequest struct {
	Action string `json:"action" binding:"required,oneof=GRANTED REVOKED"`
	Notes  string `json:"notes"`
}

type CandidateConsentResponse struct {
	ID          string `json:"id"`
	CandidateID string `json:"candidate_id"`
	Action      string `json:"action"`
	Notes       string `json:"notes,omitempty"`
	ChangedBy   string `json:"changed_by,omitempty"`
	ChangedAt   int64  `json:"changed_at"`
}

// =========================================================================
// Job Application DTOs
// =========================================================================

type CreateApplicationRequest struct {
	RequisitionID string `json:"requisition_id" binding:"required"`
	CandidateID   string `json:"candidate_id" binding:"required"`
	Notes         string `json:"notes"`
}

type UpdateApplicationStatusRequest struct {
	Status          string `json:"status" binding:"required,oneof=NEW SCREENED SHORTLISTED INTERVIEWED OFFERED ACCEPTED REJECTED WITHDRAWN"`
	RejectionReason string `json:"rejection_reason"`
	Notes           string `json:"notes"`
}

type ApplicationResponse struct {
	ID              string    `json:"id"`
	RequisitionID   string    `json:"requisition_id"`
	CandidateID     string    `json:"candidate_id"`
	Status          string    `json:"status"`
	RejectionReason string    `json:"rejection_reason,omitempty"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type StageRef struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type StageHistoryResponse struct {
	ID        string    `json:"id"`
	FromStage *StageRef `json:"from_stage"`
	ToStage   StageRef  `json:"to_stage"`
	ChangedBy string    `json:"changed_by,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	ChangedAt int64     `json:"changed_at"`
}

// =========================================================================
// Application Screening DTOs (G-7 sub-project 1)
// =========================================================================

type CreateApplicationScreeningRequest struct {
	ScreenedBy string   `json:"screened_by"`
	ScreenedAt int64    `json:"screened_at"`
	Score      *float64 `json:"score"`
	Result     string   `json:"result" binding:"omitempty,oneof=PASS FAIL HOLD"`
	Notes      string   `json:"notes"`
}

type UpdateApplicationScreeningRequest struct {
	ScreenedBy *string  `json:"screened_by"`
	ScreenedAt *int64   `json:"screened_at"`
	Score      *float64 `json:"score"`
	Result     *string  `json:"result" binding:"omitempty,oneof=PASS FAIL HOLD"`
	Notes      *string  `json:"notes"`
}

type ApplicationScreeningResponse struct {
	ID            string    `json:"id"`
	ApplicationID string    `json:"application_id"`
	ScreenedBy    string    `json:"screened_by,omitempty"`
	ScreenedAt    int64     `json:"screened_at"`
	Score         float64   `json:"score,omitempty"`
	Result        string    `json:"result"`
	Notes         string    `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// =========================================================================
// Recruitment Assessment DTOs (G-7 sub-project 2)
// =========================================================================

type CreateAssessmentRequest struct {
	RequisitionID string `json:"requisition_id"`
	Name          string `json:"name" binding:"required,max=255"`
	Type          string `json:"type" binding:"omitempty,oneof=TECHNICAL PSYCHOLOGICAL COGNITIVE PERSONALITY CASE_STUDY CODING LANGUAGE OTHER"`
	ScheduledAt   int64  `json:"scheduled_at"`
	Location      string `json:"location"`
	MeetingLink   string `json:"meeting_link"`
	Notes         string `json:"notes"`
}

type UpdateAssessmentRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=255"`
	Type        *string `json:"type" binding:"omitempty,oneof=TECHNICAL PSYCHOLOGICAL COGNITIVE PERSONALITY CASE_STUDY CODING LANGUAGE OTHER"`
	ScheduledAt *int64  `json:"scheduled_at"`
	Location    *string `json:"location"`
	MeetingLink *string `json:"meeting_link"`
	Notes       *string `json:"notes"`
}

type AssessmentResponse struct {
	ID            string    `json:"id"`
	RequisitionID string    `json:"requisition_id,omitempty"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	ScheduledAt   int64     `json:"scheduled_at"`
	Location      string    `json:"location,omitempty"`
	MeetingLink   string    `json:"meeting_link,omitempty"`
	Notes         string    `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AddAssessmentParticipantRequest struct {
	ApplicationID string `json:"application_id" binding:"required"`
}

type UpdateAssessmentParticipantRequest struct {
	Status         *string  `json:"status" binding:"omitempty,oneof=INVITED COMPLETED NO_SHOW"`
	Score          *float64 `json:"score"`
	Result         *string  `json:"result" binding:"omitempty,oneof=PASS FAIL HOLD"`
	Recommendation *string  `json:"recommendation"`
}

type AssessmentParticipantResponse struct {
	ID             string    `json:"id"`
	AssessmentID   string    `json:"assessment_id"`
	ApplicationID  string    `json:"application_id"`
	Status         string    `json:"status"`
	Score          float64   `json:"score,omitempty"`
	Result         string    `json:"result,omitempty"`
	Recommendation string    `json:"recommendation,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// =========================================================================
// Job Requisition Requirement + Competency DTOs (G-9 sub-project 1)
// =========================================================================

type CreateRequisitionRequirementRequest struct {
	RequirementType string   `json:"requirement_type" binding:"required,max=50"`
	Name            string   `json:"name" binding:"required,max=255"`
	Description     string   `json:"description"`
	MinimumValue    *float64 `json:"minimum_value"`
	MaximumValue    *float64 `json:"maximum_value"`
	IsRequired      *bool    `json:"is_required"`
	SortOrder       *int     `json:"sort_order"`
}

type UpdateRequisitionRequirementRequest struct {
	RequirementType *string  `json:"requirement_type" binding:"omitempty,max=50"`
	Name            *string  `json:"name" binding:"omitempty,max=255"`
	Description     *string  `json:"description"`
	MinimumValue    *float64 `json:"minimum_value"`
	MaximumValue    *float64 `json:"maximum_value"`
	IsRequired      *bool    `json:"is_required"`
	SortOrder       *int     `json:"sort_order"`
}

type RequisitionRequirementResponse struct {
	ID              string    `json:"id"`
	RequisitionID   string    `json:"requisition_id"`
	RequirementType string    `json:"requirement_type"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	MinimumValue    float64   `json:"minimum_value,omitempty"`
	MaximumValue    float64   `json:"maximum_value,omitempty"`
	IsRequired      bool      `json:"is_required"`
	SortOrder       int       `json:"sort_order"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateRequisitionCompetencyRequest struct {
	CompetencyID  string   `json:"competency_id" binding:"required"`
	RequiredLevel *int     `json:"required_level"`
	IsRequired    *bool    `json:"is_required"`
	Weight        *float64 `json:"weight"`
}

type UpdateRequisitionCompetencyRequest struct {
	RequiredLevel *int     `json:"required_level"`
	IsRequired    *bool    `json:"is_required"`
	Weight        *float64 `json:"weight"`
}

type RequisitionCompetencyResponse struct {
	ID            string    `json:"id"`
	RequisitionID string    `json:"requisition_id"`
	CompetencyID  string    `json:"competency_id"`
	RequiredLevel int       `json:"required_level,omitempty"`
	IsRequired    bool      `json:"is_required"`
	Weight        float64   `json:"weight,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// =========================================================================
// Candidate Match Score DTOs (G-9 sub-project 2)
// =========================================================================

type MatchScoreCompetencyBreakdown struct {
	CompetencyID   string  `json:"competency_id"`
	CompetencyName string  `json:"competency_name,omitempty"`
	RequiredLevel  int     `json:"required_level"`
	CandidateLevel int     `json:"candidate_level"`
	Weight         float64 `json:"weight"`
	Contribution   float64 `json:"contribution"`
}

type MatchScoreResponse struct {
	ApplicationID string                          `json:"application_id"`
	CandidateID   string                          `json:"candidate_id"`
	RequisitionID string                          `json:"requisition_id"`
	Score         *float64                        `json:"score"`
	Breakdown     []MatchScoreCompetencyBreakdown `json:"breakdown"`
	Note          string                          `json:"note,omitempty"`
}

// =========================================================================
// Interview DTOs
// =========================================================================

type CreateInterviewRequest struct {
	ApplicationID   string `json:"application_id" binding:"required"`
	InterviewerID   string `json:"interviewer_id" binding:"required"`
	Stage           string `json:"stage" binding:"omitempty,max=50"`
	ScheduledAt     int64  `json:"scheduled_at" binding:"required"`
	DurationMinutes *int   `json:"duration_minutes"`
	Location        string `json:"location"`
	MeetingLink     string `json:"meeting_link"`
}

type UpdateInterviewRequest struct {
	InterviewerID   *string  `json:"interviewer_id"`
	Stage           *string  `json:"stage" binding:"omitempty,max=50"`
	ScheduledAt     *int64   `json:"scheduled_at"`
	DurationMinutes *int     `json:"duration_minutes"`
	Location        *string  `json:"location"`
	MeetingLink     *string  `json:"meeting_link"`
	Status          *string  `json:"status" binding:"omitempty,oneof=SCHEDULED COMPLETED CANCELLED RESCHEDULED"`
	Score           *float64 `json:"score"`
	Feedback        *string  `json:"feedback"`
}

type InterviewResponse struct {
	ID              string    `json:"id"`
	ApplicationID   string    `json:"application_id"`
	InterviewerID   string    `json:"interviewer_id"`
	Stage           string    `json:"stage"`
	DurationMinutes int       `json:"duration_minutes"`
	Location        string    `json:"location,omitempty"`
	MeetingLink     string    `json:"meeting_link,omitempty"`
	Status          string    `json:"status"`
	Score           float64   `json:"score,omitempty"`
	Feedback        string    `json:"feedback"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// =========================================================================
// Interviewer DTOs (G-8 — multi-interviewer)
// =========================================================================

type AddInterviewerRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Role       string `json:"role"`
}

type InterviewerResponse struct {
	ID          string    `json:"id"`
	InterviewID string    `json:"interview_id"`
	EmployeeID  string    `json:"employee_id"`
	Role        string    `json:"role,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// =========================================================================
// Interview Scorecard Item DTOs (G-8)
// =========================================================================

type AddScorecardItemRequest struct {
	Criterion string   `json:"criterion" binding:"required,max=255"`
	Weight    float64  `json:"weight"`
	Score     *float64 `json:"score"`
	Notes     string   `json:"notes"`
}

type UpdateScorecardItemRequest struct {
	Criterion *string  `json:"criterion" binding:"omitempty,max=255"`
	Weight    *float64 `json:"weight"`
	Score     *float64 `json:"score"`
	Notes     *string  `json:"notes"`
}

type ScorecardItemResponse struct {
	ID          string    `json:"id"`
	InterviewID string    `json:"interview_id"`
	Criterion   string    `json:"criterion"`
	Weight      float64   `json:"weight"`
	Score       float64   `json:"score,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// =========================================================================
// Onboarding Task Template DTOs
// =========================================================================

type CreateOnboardingTaskTemplateRequest struct {
	Name           string  `json:"name" binding:"required,max=255"`
	Description    string  `json:"description"`
	Category       string  `json:"category" binding:"omitempty,max=50"`
	DayOffset      *int    `json:"day_offset"`
	AssignedRole   string  `json:"assigned_role" binding:"omitempty,max=50"`
	IsMandatory    *bool   `json:"is_mandatory"`
	OrganizationID *string `json:"organization_id"`
	PositionID     *string `json:"position_id"`
	EmploymentType *string `json:"employment_type"`
}

type UpdateOnboardingTaskTemplateRequest struct {
	Name           *string `json:"name" binding:"omitempty,max=255"`
	Description    *string `json:"description"`
	Category       *string `json:"category" binding:"omitempty,max=50"`
	DayOffset      *int    `json:"day_offset"`
	AssignedRole   *string `json:"assigned_role" binding:"omitempty,max=50"`
	IsMandatory    *bool   `json:"is_mandatory"`
	OrganizationID *string `json:"organization_id"`
	PositionID     *string `json:"position_id"`
	EmploymentType *string `json:"employment_type"`
}

type OnboardingTaskTemplateResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Category       string    `json:"category"`
	DayOffset      int       `json:"day_offset"`
	AssignedRole   string    `json:"assigned_role"`
	IsMandatory    bool      `json:"is_mandatory"`
	OrganizationID string    `json:"organization_id,omitempty"`
	PositionID     string    `json:"position_id,omitempty"`
	EmploymentType string    `json:"employment_type,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// =========================================================================
// Employee Onboarding DTOs
// =========================================================================

type CreateEmployeeOnboardingRequest struct {
	EmployeeID    string  `json:"employee_id" binding:"required"`
	ApplicationID string  `json:"application_id" binding:"required"`
	StartDate     string  `json:"start_date" binding:"required"`
	BuddyID       *string `json:"buddy_id"`
	Notes         string  `json:"notes"`
}

type UpdateEmployeeOnboardingRequest struct {
	StartDate *string `json:"start_date"`
	BuddyID   *string `json:"buddy_id"`
	Status    *string `json:"status" binding:"omitempty,oneof=PENDING IN_PROGRESS COMPLETED"`
	Notes     *string `json:"notes"`
}

type EmployeeOnboardingResponse struct {
	ID            string    `json:"id"`
	EmployeeID    string    `json:"employee_id"`
	ApplicationID string    `json:"application_id"`
	StartDate     string    `json:"start_date"`
	Status        string    `json:"status"`
	BuddyID       string    `json:"buddy_id,omitempty"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// =========================================================================
// Onboarding Task Item DTOs
// =========================================================================

type CreateOnboardingTaskItemRequest struct {
	EmployeeOnboardingID string  `json:"employee_onboarding_id" binding:"required"`
	TemplateID           *string `json:"template_id"`
	Name                 string  `json:"name" binding:"required,max=255"`
	Description          string  `json:"description"`
	AssignedTo           *string `json:"assigned_to"`
	DueDate              *int64  `json:"due_date"`
}

type UpdateOnboardingTaskItemRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=255"`
	Description *string `json:"description"`
	AssignedTo  *string `json:"assigned_to"`
	DueDate     *int64  `json:"due_date"`
	IsCompleted *bool   `json:"is_completed"`
	Notes       *string `json:"notes"`
}

type OnboardingTaskItemResponse struct {
	ID                   string    `json:"id"`
	EmployeeOnboardingID string    `json:"employee_onboarding_id"`
	TemplateID           string    `json:"template_id,omitempty"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	AssignedTo           string    `json:"assigned_to,omitempty"`
	IsCompleted          bool      `json:"is_completed"`
	Notes                string    `json:"notes"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// =========================================================================
// Pagination
// =========================================================================

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}
