package recruitment

import "time"

// =========================================================================
// Job Requisition DTOs
// =========================================================================

type CreateRequisitionRequest struct {
	OrganizationID   string  `json:"organization_id" binding:"required"`
	Title            string  `json:"title" binding:"required,max=255"`
	// G-2: requisition_number opsional — bila kosong auto-generated
	// REQ-YYYYMM-XXXX oleh service.
	RequisitionNumber string `json:"requisition_number" binding:"omitempty,max=50"`
	Priority         *string `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH URGENT"`
	PositionID       *string `json:"position_id"`
	Department       string  `json:"department" binding:"omitempty,max=150"`
	EmploymentType   string  `json:"employment_type" binding:"omitempty,max=50"`
	Location         string  `json:"location" binding:"omitempty,max=255"`
	MinSalary        float64 `json:"min_salary"`
	MaxSalary        float64 `json:"max_salary"`
	Description      string  `json:"description"`
	Requirements     string  `json:"requirements"`
	Responsibilities string  `json:"responsibilities"`
	SlotsAvailable   *int    `json:"slots_available"`
	RequestedBy      *string `json:"requested_by"`
	TargetStartDate  *string `json:"target_start_date"`
	ReasonType       *string `json:"reason_type" binding:"omitempty,oneof=NEW_POSITION REPLACEMENT EXPANSION WORKFORCE_GAP SUCCESSION_GAP"`
	WorkforceGapID   *string `json:"workforce_gap_id"`
	WorkforcePlanID  *string `json:"workforce_plan_id"`
	SuccessionPositionID *string `json:"succession_position_id"`
}

type UpdateRequisitionRequest struct {
	Title            *string  `json:"title" binding:"omitempty,max=255"`
	RequisitionNumber *string `json:"requisition_number" binding:"omitempty,max=50"`
	Priority         *string  `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH URGENT"`
	PositionID       *string  `json:"position_id"`
	Department       *string  `json:"department" binding:"omitempty,max=150"`
	EmploymentType   *string  `json:"employment_type" binding:"omitempty,max=50"`
	Location         *string  `json:"location" binding:"omitempty,max=255"`
	MinSalary        *float64 `json:"min_salary"`
	MaxSalary        *float64 `json:"max_salary"`
	Description      *string  `json:"description"`
	Requirements     *string  `json:"requirements"`
	Responsibilities *string  `json:"responsibilities"`
	SlotsAvailable   *int     `json:"slots_available"`
	Status           *string  `json:"status" binding:"omitempty,oneof=DRAFT OPEN IN_PROGRESS FILLED CANCELLED"`
	// Catatan G-1: SUBMITTED/REJECTED tidak boleh di-set manual via update —
	// keduanya hanya dihasilkan alur approval (submit + push-callback).
	TargetStartDate  *string  `json:"target_start_date"`
	ReasonType       *string  `json:"reason_type" binding:"omitempty,oneof=NEW_POSITION REPLACEMENT EXPANSION WORKFORCE_GAP SUCCESSION_GAP"`
	WorkforceGapID   *string  `json:"workforce_gap_id"`
	WorkforcePlanID  *string  `json:"workforce_plan_id"`
	SuccessionPositionID *string `json:"succession_position_id"`
}

// SubmitRequisitionRequest mengirim requisition draft ke Central Approval
// (plan G-1). FlowID opsional — bila kosong, flow aktif untuk modul
// "recruitment" di-auto-resolve (pola employeemovement G-3).
type SubmitRequisitionRequest struct {
	FlowID *string `json:"flow_id"`
}

type RequisitionResponse struct {
	ID                string    `json:"id"`
	OrganizationID    string    `json:"organization_id"`
	Title             string    `json:"title"`
	RequisitionNumber string    `json:"requisition_number,omitempty"`
	Priority          string    `json:"priority"`
	PositionID        string    `json:"position_id,omitempty"`
	OpenedAt          *int64    `json:"opened_at,omitempty"`
	Department        string    `json:"department,omitempty"`
	EmploymentType    string    `json:"employment_type,omitempty"`
	Location          string    `json:"location,omitempty"`
	MinSalary         float64   `json:"min_salary"`
	MaxSalary         float64   `json:"max_salary"`
	Description       string    `json:"description"`
	Requirements      string    `json:"requirements"`
	Responsibilities  string    `json:"responsibilities"`
	SlotsAvailable    int       `json:"slots_available"`
	SlotsFilled       int       `json:"slots_filled"`
	Status            string    `json:"status"`
	RequestedBy       string    `json:"requested_by,omitempty"`
	ApprovedBy        string    `json:"approved_by,omitempty"`
	ApprovalInstanceID string   `json:"approval_instance_id,omitempty"`
	ReasonType        string    `json:"reason_type,omitempty"`
	WorkforceGapID    string    `json:"workforce_gap_id,omitempty"`
	WorkforcePlanID   string    `json:"workforce_plan_id,omitempty"`
	SuccessionPositionID string `json:"succession_position_id,omitempty"`
	TargetStartDate   string    `json:"target_start_date,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
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
}

type CandidateResponse struct {
	ID             string    `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone,omitempty"`
	Address        string    `json:"address,omitempty"`
	CurrentCompany string    `json:"current_company,omitempty"`
	CurrentTitle   string    `json:"current_title,omitempty"`
	ResumeURL      string    `json:"resume_url,omitempty"`
	PortfolioURL   string    `json:"portfolio_url,omitempty"`
	LinkedInURL    string    `json:"linkedin_url,omitempty"`
	Source         string    `json:"source"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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
	InterviewerID   *string `json:"interviewer_id"`
	Stage           *string `json:"stage" binding:"omitempty,max=50"`
	ScheduledAt     *int64  `json:"scheduled_at"`
	DurationMinutes *int    `json:"duration_minutes"`
	Location        *string `json:"location"`
	MeetingLink     *string `json:"meeting_link"`
	Status          *string `json:"status" binding:"omitempty,oneof=SCHEDULED COMPLETED CANCELLED RESCHEDULED"`
	Score           *float64 `json:"score"`
	Feedback        *string `json:"feedback"`
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
// Onboarding Task Template DTOs
// =========================================================================

type CreateOnboardingTaskTemplateRequest struct {
	Name         string `json:"name" binding:"required,max=255"`
	Description  string `json:"description"`
	Category     string `json:"category" binding:"omitempty,max=50"`
	DayOffset    *int   `json:"day_offset"`
	AssignedRole string `json:"assigned_role" binding:"omitempty,max=50"`
	IsMandatory  *bool  `json:"is_mandatory"`
}

type UpdateOnboardingTaskTemplateRequest struct {
	Name         *string `json:"name" binding:"omitempty,max=255"`
	Description  *string `json:"description"`
	Category     *string `json:"category" binding:"omitempty,max=50"`
	DayOffset    *int    `json:"day_offset"`
	AssignedRole *string `json:"assigned_role" binding:"omitempty,max=50"`
	IsMandatory  *bool   `json:"is_mandatory"`
}

type OnboardingTaskTemplateResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	DayOffset    int       `json:"day_offset"`
	AssignedRole string    `json:"assigned_role"`
	IsMandatory  bool      `json:"is_mandatory"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// =========================================================================
// Employee Onboarding DTOs
// =========================================================================

type CreateEmployeeOnboardingRequest struct {
	EmployeeID    string `json:"employee_id" binding:"required"`
	ApplicationID string `json:"application_id" binding:"required"`
	StartDate     string `json:"start_date" binding:"required"`
	BuddyID       *string `json:"buddy_id"`
	Notes         string `json:"notes"`
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
	ID                    string    `json:"id"`
	EmployeeOnboardingID string    `json:"employee_onboarding_id"`
	TemplateID            string    `json:"template_id,omitempty"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	AssignedTo            string    `json:"assigned_to,omitempty"`
	IsCompleted           bool      `json:"is_completed"`
	Notes                 string    `json:"notes"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
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
