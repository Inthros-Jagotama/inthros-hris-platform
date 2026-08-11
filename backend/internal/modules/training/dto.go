package training

import "time"

// =========================================================================
// Training Category DTOs
// =========================================================================

type CreateTrainingCategoryRequest struct {
	Code        string  `json:"code" binding:"required,max=20"`
	Name        string  `json:"name" binding:"required,max=150"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateTrainingCategoryRequest struct {
	Code        *string `json:"code" binding:"omitempty,max=20"`
	Name        *string `json:"name" binding:"omitempty,max=150"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type TrainingCategoryResponse struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// =========================================================================
// Training Course DTOs
// =========================================================================

type CreateTrainingCourseRequest struct {
	CategoryID     string   `json:"category_id" binding:"required"`
	Code           string   `json:"code" binding:"required,max=20"`
	Name           string   `json:"name" binding:"required,max=200"`
	Description    *string  `json:"description"`
	DurationHour   *float64 `json:"duration_hour"`
	MinScore       *float64 `json:"min_score"`
	Cost           *float64 `json:"cost"`
	IsCertified    *bool    `json:"is_certified"`
	ExternalVendor *string  `json:"external_vendor" binding:"omitempty,max=200"`
	// P0-BE (plan §7)
	CourseType   *string `json:"course_type" binding:"omitempty,oneof=TECHNICAL SOFT_SKILL COMPLIANCE MANAGEMENT CERTIFICATION OTHER"`
	DeliveryType *string `json:"delivery_type" binding:"omitempty,oneof=IN_HOUSE EXTERNAL BOTH"`
	IsMandatory  *bool   `json:"is_mandatory"`
}

type UpdateTrainingCourseRequest struct {
	CategoryID     *string  `json:"category_id"`
	Code           *string  `json:"code" binding:"omitempty,max=20"`
	Name           *string  `json:"name" binding:"omitempty,max=200"`
	Description    *string  `json:"description"`
	DurationHour   *float64 `json:"duration_hour"`
	MinScore       *float64 `json:"min_score"`
	Cost           *float64 `json:"cost"`
	IsCertified    *bool    `json:"is_certified"`
	ExternalVendor *string  `json:"external_vendor" binding:"omitempty,max=200"`
	IsActive       *bool    `json:"is_active"`
	// P0-BE (plan §7)
	CourseType   *string `json:"course_type" binding:"omitempty,oneof=TECHNICAL SOFT_SKILL COMPLIANCE MANAGEMENT CERTIFICATION OTHER"`
	DeliveryType *string `json:"delivery_type" binding:"omitempty,oneof=IN_HOUSE EXTERNAL BOTH"`
	IsMandatory  *bool   `json:"is_mandatory"`
}

type TrainingCourseResponse struct {
	ID             string    `json:"id"`
	CategoryID     string    `json:"category_id"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	Description    string    `json:"description,omitempty"`
	DurationHour   float64   `json:"duration_hour,omitempty"`
	MinScore       float64   `json:"min_score,omitempty"`
	Cost           float64   `json:"cost,omitempty"`
	IsCertified    bool      `json:"is_certified"`
	ExternalVendor string    `json:"external_vendor,omitempty"`
	// P0-BE (plan §7)
	CourseType   string `json:"course_type,omitempty"`
	DeliveryType string `json:"delivery_type,omitempty"`
	IsMandatory  bool   `json:"is_mandatory"`
	IsActive     bool   `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// =========================================================================
// Training Session DTOs
// =========================================================================

type CreateTrainingSessionRequest struct {
	CourseID    string  `json:"course_id" binding:"required"`
	SessionCode string  `json:"session_code" binding:"required,max=20"`
	TrainerName string  `json:"trainer_name" binding:"required,max=200"`
	Location    *string `json:"location" binding:"omitempty,max=255"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     string  `json:"end_date" binding:"required"`
	MaxQuota    int     `json:"max_quota" binding:"omitempty,min=1"`
	// P0-BE (plan §14)
	ProviderType         *string `json:"provider_type" binding:"omitempty,oneof=IN_HOUSE EXTERNAL"`
	DeliveryMode         *string `json:"delivery_mode" binding:"omitempty,oneof=ONSITE ONLINE HYBRID SELF_PACED"`
	ProviderID           *string `json:"provider_id"`
	StartDatetime        *string `json:"start_datetime"`
	EndDatetime          *string `json:"end_datetime"`
	MeetingURL           *string `json:"meeting_url"`
	RegistrationDeadline *string `json:"registration_deadline"`
}

type UpdateTrainingSessionRequest struct {
	SessionCode *string `json:"session_code" binding:"omitempty,max=20"`
	TrainerName *string `json:"trainer_name" binding:"omitempty,max=200"`
	Location    *string `json:"location" binding:"omitempty,max=255"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	MaxQuota    *int    `json:"max_quota" binding:"omitempty,min=1"`
	// P0-BE (plan §14)
	ProviderType         *string `json:"provider_type" binding:"omitempty,oneof=IN_HOUSE EXTERNAL"`
	DeliveryMode         *string `json:"delivery_mode" binding:"omitempty,oneof=ONSITE ONLINE HYBRID SELF_PACED"`
	ProviderID           *string `json:"provider_id"`
	StartDatetime        *string `json:"start_datetime"`
	EndDatetime          *string `json:"end_datetime"`
	MeetingURL           *string `json:"meeting_url"`
	RegistrationDeadline *string `json:"registration_deadline"`
}

type UpdateSessionStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=DRAFT SCHEDULED REGISTRATION_OPEN FULL IN_PROGRESS COMPLETED CANCELLED"`
}

type TrainingSessionResponse struct {
	ID          string    `json:"id"`
	CourseID    string    `json:"course_id"`
	SessionCode string    `json:"session_code"`
	TrainerName string    `json:"trainer_name"`
	// P0-BE (plan §14)
	ProviderType         string `json:"provider_type,omitempty"`
	DeliveryMode         string `json:"delivery_mode,omitempty"`
	ProviderID           string `json:"provider_id,omitempty"`
	StartDatetime        string `json:"start_datetime,omitempty"`
	EndDatetime          string `json:"end_datetime,omitempty"`
	MeetingURL           string `json:"meeting_url,omitempty"`
	RegistrationDeadline string `json:"registration_deadline,omitempty"`
	Location    string    `json:"location,omitempty"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
	MaxQuota    int       `json:"max_quota"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// =========================================================================
// Training Participant DTOs
// =========================================================================

type CreateTrainingParticipantRequest struct {
	SessionID  string `json:"session_id" binding:"required"`
	EmployeeID string `json:"employee_id" binding:"required"`
	// P0-BE (plan §18) — default REGISTERED bila kosong.
	RegistrationStatus *string `json:"registration_status" binding:"omitempty,oneof=NOMINATED REQUESTED APPROVED REGISTERED WAITLISTED CANCELLED"`
}

type UpdateTrainingParticipantRequest struct {
	AttendanceStatus *string `json:"attendance_status" binding:"omitempty,oneof=PRESENT ABSENT LATE EXCUSED"`
	Score            *float64 `json:"score"`
	// P0-BE (plan §18) — completion fields.
	CompletionStatus *string  `json:"completion_status" binding:"omitempty,oneof=NOT_STARTED IN_PROGRESS COMPLETED FAILED"`
	FinalScore       *float64 `json:"final_score"`
	Passed           *bool    `json:"passed"`
	Remarks          *string  `json:"remarks"`
}

type TrainingParticipantResponse struct {
	ID                 string    `json:"id"`
	SessionID          string    `json:"session_id"`
	EmployeeID         string    `json:"employee_id"`
	RegistrationStatus string    `json:"registration_status"`
	RegisteredAt       string    `json:"registered_at,omitempty"`
	ApprovedAt         string    `json:"approved_at,omitempty"`
	AttendanceStatus   string    `json:"attendance_status"`
	Score              float64   `json:"score"`
	CompletionStatus   string    `json:"completion_status"`
	CompletionDate     string    `json:"completion_date,omitempty"`
	FinalScore         float64   `json:"final_score,omitempty"`
	Passed             bool      `json:"passed,omitempty"`
	Remarks            string    `json:"remarks,omitempty"`
	CompletedAt        string    `json:"completed_at,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// =========================================================================
// Training Material DTOs
// =========================================================================

type CreateTrainingMaterialRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Title     string `json:"title" binding:"required,max=200"`
	// P0-BE (plan §20)
	Description   *string `json:"description"`
	IsRequired    *bool   `json:"is_required"`
	AvailableFrom *string `json:"available_from"`
	FileURL   *string `json:"file_url"`
	FileType  *string `json:"file_type" binding:"omitempty,max=50"`
	SortOrder *int    `json:"sort_order"`
}

type UpdateTrainingMaterialRequest struct {
	Title     *string `json:"title" binding:"omitempty,max=200"`
	// P0-BE (plan §20)
	Description   *string `json:"description"`
	IsRequired    *bool   `json:"is_required"`
	AvailableFrom *string `json:"available_from"`
	FileURL   *string `json:"file_url"`
	FileType  *string `json:"file_type" binding:"omitempty,max=50"`
	SortOrder *int    `json:"sort_order"`
}

type TrainingMaterialResponse struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Title     string    `json:"title"`
	// P0-BE (plan §20)
	Description   string `json:"description,omitempty"`
	IsRequired    bool   `json:"is_required"`
	AvailableFrom string `json:"available_from,omitempty"`
	FileURL   string    `json:"file_url,omitempty"`
	FileType  string    `json:"file_type,omitempty"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// =========================================================================
// Training Evaluation DTOs
// =========================================================================

type CreateTrainingEvaluationRequest struct {
	SessionID  string `json:"session_id" binding:"required"`
	EmployeeID string `json:"employee_id" binding:"required"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	Feedback   *string `json:"feedback"`
}

type UpdateTrainingEvaluationRequest struct {
	Rating   *int    `json:"rating" binding:"omitempty,min=1,max=5"`
	Feedback *string `json:"feedback"`
}

type TrainingEvaluationResponse struct {
	ID         string    `json:"id"`
	SessionID  string    `json:"session_id"`
	EmployeeID string    `json:"employee_id"`
	Rating     int       `json:"rating"`
	Feedback   string    `json:"feedback,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// =========================================================================
// Training Certificate DTOs
// =========================================================================

type CreateTrainingCertificateRequest struct {
	ParticipantID string `json:"participant_id" binding:"required"`
	CertificateNo string `json:"certificate_no" binding:"required,max=50"`
	IssuedDate    string `json:"issued_date" binding:"required"`
	ExpiryDate    *string `json:"expiry_date"`
}

type UpdateTrainingCertificateRequest struct {
	CertificateNo *string `json:"certificate_no" binding:"omitempty,max=50"`
	IssuedDate    *string `json:"issued_date"`
	ExpiryDate    *string `json:"expiry_date"`
}

type TrainingCertificateResponse struct {
	ID            string    `json:"id"`
	ParticipantID string    `json:"participant_id"`
	CertificateNo string    `json:"certificate_no"`
	IssuedDate    string    `json:"issued_date"`
	ExpiryDate    string    `json:"expiry_date,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// =========================================================================
// Training Provider DTOs (P0-BE — plan §11)
// =========================================================================

type CreateTrainingProviderRequest struct {
	Code        string  `json:"code" binding:"required,max=20"`
	Name        string  `json:"name" binding:"required,max=200"`
	Type        *string `json:"type" binding:"omitempty,oneof=INTERNAL EXTERNAL"`
	ContactName *string `json:"contact_name"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	Address     *string `json:"address"`
	Website     *string `json:"website"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateTrainingProviderRequest struct {
	Code        *string `json:"code" binding:"omitempty,max=20"`
	Name        *string `json:"name" binding:"omitempty,max=200"`
	Type        *string `json:"type" binding:"omitempty,oneof=INTERNAL EXTERNAL"`
	ContactName *string `json:"contact_name"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	Address     *string `json:"address"`
	Website     *string `json:"website"`
	IsActive    *bool   `json:"is_active"`
}

type TrainingProviderResponse struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	ContactName string    `json:"contact_name,omitempty"`
	Email       string    `json:"email,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	Address     string    `json:"address,omitempty"`
	Website     string    `json:"website,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// =========================================================================
// Training Trainer DTOs (P0-BE — plan §12)
// =========================================================================

type CreateTrainingTrainerRequest struct {
	Type       string  `json:"type" binding:"required,oneof=INTERNAL EXTERNAL"`
	EmployeeID *string `json:"employee_id"`
	ProviderID *string `json:"provider_id"`
	Name       string  `json:"name" binding:"required,max=200"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	Bio        *string `json:"bio"`
	IsActive   *bool   `json:"is_active"`
}

type UpdateTrainingTrainerRequest struct {
	Type       *string `json:"type" binding:"omitempty,oneof=INTERNAL EXTERNAL"`
	EmployeeID *string `json:"employee_id"`
	ProviderID *string `json:"provider_id"`
	Name       *string `json:"name" binding:"omitempty,max=200"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	Bio        *string `json:"bio"`
	IsActive   *bool   `json:"is_active"`
}

type TrainingTrainerResponse struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	EmployeeID string    `json:"employee_id,omitempty"`
	ProviderID string    `json:"provider_id,omitempty"`
	Name       string    `json:"name"`
	Email      string    `json:"email,omitempty"`
	Phone      string    `json:"phone,omitempty"`
	Bio        string    `json:"bio,omitempty"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// =========================================================================
// Training Session Trainer DTOs (P0-BE — plan §13)
// =========================================================================

type AddSessionTrainerRequest struct {
	TrainerID string  `json:"trainer_id" binding:"required"`
	Role      *string `json:"role" binding:"omitempty,oneof=MAIN ASSISTANT"`
}

type TrainingSessionTrainerResponse struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	TrainerID string    `json:"trainer_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// =========================================================================
// Training Attendance DTOs (P0-BE — plan §19)
// =========================================================================

type MarkTrainingAttendanceRequest struct {
	ParticipantID  string  `json:"participant_id" binding:"required"`
	AttendanceDate string  `json:"attendance_date" binding:"required"`
	CheckIn        *string `json:"check_in"`
	CheckOut       *string `json:"check_out"`
	Status         *string `json:"status" binding:"omitempty,oneof=PRESENT ABSENT LATE EXCUSED"`
	Remarks        *string `json:"remarks"`
}

type UpdateTrainingAttendanceRequest struct {
	CheckIn  *string `json:"check_in"`
	CheckOut *string `json:"check_out"`
	Status   *string `json:"status" binding:"omitempty,oneof=PRESENT ABSENT LATE EXCUSED"`
	Remarks  *string `json:"remarks"`
}

type TrainingAttendanceResponse struct {
	ID             string `json:"id"`
	ParticipantID  string `json:"participant_id"`
	AttendanceDate string `json:"attendance_date"`
	CheckIn        string `json:"check_in,omitempty"`
	CheckOut       string `json:"check_out,omitempty"`
	Status         string `json:"status"`
	Remarks        string `json:"remarks,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// SessionAttendanceRow — baris attendance per peserta dalam satu session
// (hasil join training_attendances ↔ training_participants).
type SessionAttendanceRow struct {
	AttendanceID   string `json:"attendance_id,omitempty"`
	ParticipantID  string `json:"participant_id"`
	EmployeeID     string `json:"employee_id"`
	AttendanceDate string `json:"attendance_date"`
	Status         string `json:"status"`
	CheckIn        string `json:"check_in,omitempty"`
	CheckOut       string `json:"check_out,omitempty"`
	Remarks        string `json:"remarks,omitempty"`
}

// =========================================================================
// Training Assessment DTOs (P0-BE — plan §21)
// =========================================================================

type CreateTrainingAssessmentRequest struct {
	SessionID    string   `json:"session_id" binding:"required"`
	Name         string   `json:"name" binding:"required,max=200"`
	Type         *string  `json:"type" binding:"omitempty,oneof=PRE_TEST POST_TEST FINAL PRACTICAL OTHER"`
	MaxScore     *float64 `json:"max_score"`
	PassingScore *float64 `json:"passing_score"`
	AttemptLimit *int     `json:"attempt_limit"`
	IsRequired   *bool    `json:"is_required"`
}

type TrainingAssessmentResponse struct {
	ID           string    `json:"id"`
	SessionID    string    `json:"session_id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	MaxScore     float64   `json:"max_score"`
	PassingScore float64   `json:"passing_score"`
	AttemptLimit int       `json:"attempt_limit"`
	IsRequired   bool      `json:"is_required"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SubmitAssessmentResultRequest struct {
	ParticipantID string  `json:"participant_id" binding:"required"`
	Score         float64 `json:"score" binding:"required"`
	CompletedAt   *string `json:"completed_at"`
}

type TrainingAssessmentResultResponse struct {
	ID            string    `json:"id"`
	AssessmentID  string    `json:"assessment_id"`
	ParticipantID string    `json:"participant_id"`
	Score         float64   `json:"score"`
	Passed        bool      `json:"passed"`
	Attempt       int       `json:"attempt"`
	CompletedAt   string    `json:"completed_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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
// =========================================================================
// Training Plan DTOs (P1-BE — plan §16)
// =========================================================================

type CreateTrainingPlanRequest struct {
	Code        string  `json:"code" binding:"required,max=30"`
	Name        string  `json:"name" binding:"required,max=200"`
	Year        int     `json:"year" binding:"required,min=2000,max=2100"`
	Description *string `json:"description"`
	Status      *string `json:"status" binding:"omitempty,oneof=DRAFT ACTIVE ARCHIVED"`
}

type UpdateTrainingPlanRequest struct {
	Code        *string `json:"code" binding:"omitempty,max=30"`
	Name        *string `json:"name" binding:"omitempty,max=200"`
	Year        *int    `json:"year" binding:"omitempty,min=2000,max=2100"`
	Description *string `json:"description"`
	Status      *string `json:"status" binding:"omitempty,oneof=DRAFT ACTIVE ARCHIVED"`
}

type TrainingPlanResponse struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Year        int       `json:"year"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTrainingPlanItemRequest struct {
	CourseID           string   `json:"course_id" binding:"required"`
	TargetDate         *string  `json:"target_date"`
	TargetParticipants *int     `json:"target_participants"`
	EstimatedCost      *float64 `json:"estimated_cost"`
	Priority           *string  `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH URGENT"`
}

type UpdateTrainingPlanItemRequest struct {
	CourseID           *string  `json:"course_id"`
	TargetDate         *string  `json:"target_date"`
	TargetParticipants *int     `json:"target_participants"`
	EstimatedCost      *float64 `json:"estimated_cost"`
	Priority           *string  `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH URGENT"`
}

type TrainingPlanItemResponse struct {
	ID                 string    `json:"id"`
	TrainingPlanID     string    `json:"training_plan_id"`
	CourseID           string    `json:"course_id"`
	TargetDate         string    `json:"target_date,omitempty"`
	TargetParticipants *int      `json:"target_participants,omitempty"`
	EstimatedCost      float64   `json:"estimated_cost,omitempty"`
	Priority           string    `json:"priority"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// =========================================================================
// Training Need DTOs (P1-BE — plan §17)
// =========================================================================

type CreateTrainingNeedRequest struct {
	EmployeeID     *string `json:"employee_id"`
	OrganizationID *string `json:"organization_id"`
	PositionID     *string `json:"position_id"`
	CourseID       *string `json:"course_id"`
	Reason         *string `json:"reason"`
	Priority       *string `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH URGENT"`
	SourceType     *string `json:"source_type" binding:"omitempty,oneof=MANUAL PERFORMANCE COMPETENCY CAREER SUCCESSION COMPLIANCE WORKFORCE"`
	SourceID       *string `json:"source_id"`
	Status         *string `json:"status" binding:"omitempty,oneof=OPEN PLANNED FULFILLED CANCELLED"`
}

type UpdateTrainingNeedRequest struct {
	EmployeeID     *string `json:"employee_id"`
	OrganizationID *string `json:"organization_id"`
	PositionID     *string `json:"position_id"`
	CourseID       *string `json:"course_id"`
	Reason         *string `json:"reason"`
	Priority       *string `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH URGENT"`
	SourceType     *string `json:"source_type" binding:"omitempty,oneof=MANUAL PERFORMANCE COMPETENCY CAREER SUCCESSION COMPLIANCE WORKFORCE"`
	SourceID       *string `json:"source_id"`
	Status         *string `json:"status" binding:"omitempty,oneof=OPEN PLANNED FULFILLED CANCELLED"`
}

type TrainingNeedResponse struct {
	ID             string    `json:"id"`
	EmployeeID     string    `json:"employee_id,omitempty"`
	OrganizationID string    `json:"organization_id,omitempty"`
	PositionID     string    `json:"position_id,omitempty"`
	CourseID       string    `json:"course_id,omitempty"`
	Reason         string    `json:"reason,omitempty"`
	Priority       string    `json:"priority"`
	SourceType     string    `json:"source_type"`
	SourceID       string    `json:"source_id,omitempty"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// =========================================================================
// Training Request DTOs (P1-BE — plan §15, Central Approval)
// =========================================================================

type CreateTrainingRequestRequest struct {
	EmployeeID    string  `json:"employee_id" binding:"required"`
	CourseID      string  `json:"course_id" binding:"required"`
	SessionID     *string `json:"session_id"`
	RequestedDate string  `json:"requested_date" binding:"required"`
	Reason        *string `json:"reason"`
	Priority      *string `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH URGENT"`
}

type SubmitTrainingRequestRequest struct {
	FlowID *string `json:"flow_id"`
}

type TrainingRequestResponse struct {
	ID                 string    `json:"id"`
	EmployeeID         string    `json:"employee_id"`
	CourseID           string    `json:"course_id"`
	SessionID          string    `json:"session_id,omitempty"`
	RequestedDate      string    `json:"requested_date"`
	Reason             string    `json:"reason,omitempty"`
	Priority           string    `json:"priority"`
	Status             string    `json:"status"`
	ApprovalInstanceID string    `json:"approval_instance_id,omitempty"`
	ApprovedAt         string    `json:"approved_at,omitempty"`
	RejectedAt         string    `json:"rejected_at,omitempty"`
	SupervisorNote     string    `json:"supervisor_note,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// =========================================================================
// Course Sub-resource DTOs (P1-BE — plan §8/§9/§10)
// =========================================================================

type CreateCourseObjectiveRequest struct {
	Objective string `json:"objective" binding:"required"`
	SortOrder *int   `json:"sort_order"`
}

type UpdateCourseObjectiveRequest struct {
	Objective *string `json:"objective"`
	SortOrder *int    `json:"sort_order"`
}

type CourseObjectiveResponse struct {
	ID        string    `json:"id"`
	CourseID  string    `json:"course_id"`
	Objective string    `json:"objective"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCourseCompetencyRequest struct {
	CompetencyID string `json:"competency_id" binding:"required"`
	TargetLevel  *int   `json:"target_level"`
}

type CourseCompetencyResponse struct {
	ID           string    `json:"id"`
	CourseID     string    `json:"course_id"`
	CompetencyID string    `json:"competency_id"`
	TargetLevel  *int      `json:"target_level,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateCoursePrerequisiteRequest struct {
	PrerequisiteType string  `json:"prerequisite_type" binding:"required,oneof=COURSE COMPETENCY CERTIFICATION EXPERIENCE"`
	PrerequisiteID   *string `json:"prerequisite_id"`
	IsRequired       *bool   `json:"is_required"`
}

type CoursePrerequisiteResponse struct {
	ID               string    `json:"id"`
	CourseID         string    `json:"course_id"`
	PrerequisiteType string    `json:"prerequisite_type"`
	PrerequisiteID   string    `json:"prerequisite_id,omitempty"`
	IsRequired       bool      `json:"is_required"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// =========================================================================
// Training Mandatory DTOs (P1-BE — plan §25)
// =========================================================================

type CreateTrainingMandatoryRequest struct {
	CourseID            string  `json:"course_id" binding:"required"`
	OrganizationID      *string `json:"organization_id"`
	PositionID          *string `json:"position_id"`
	EmploymentStatusID  *string `json:"employment_status_id"`
	DueDays             *int    `json:"due_days"`
	ValidityPeriodMonth *int    `json:"validity_period_month"`
	IsActive            *bool   `json:"is_active"`
}

type UpdateTrainingMandatoryRequest struct {
	CourseID            *string `json:"course_id"`
	OrganizationID      *string `json:"organization_id"`
	PositionID          *string `json:"position_id"`
	EmploymentStatusID  *string `json:"employment_status_id"`
	DueDays             *int    `json:"due_days"`
	ValidityPeriodMonth *int    `json:"validity_period_month"`
	IsActive            *bool   `json:"is_active"`
}

type TrainingMandatoryResponse struct {
	ID                  string    `json:"id"`
	CourseID            string    `json:"course_id"`
	OrganizationID      string    `json:"organization_id,omitempty"`
	PositionID          string    `json:"position_id,omitempty"`
	EmploymentStatusID  string    `json:"employment_status_id,omitempty"`
	DueDays             *int      `json:"due_days,omitempty"`
	ValidityPeriodMonth *int      `json:"validity_period_month,omitempty"`
	IsActive            bool      `json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// =========================================================================
// Training Session Cost DTOs (P1-BE — plan §26)
// =========================================================================

type CreateTrainingSessionCostRequest struct {
	CostType    string   `json:"cost_type" binding:"required,oneof=TRAINER PROVIDER VENUE MATERIAL CERTIFICATION TRAVEL ACCOMMODATION OTHER"`
	Description *string  `json:"description"`
	Amount      *float64 `json:"amount" binding:"omitempty,min=0"`
}

type UpdateTrainingSessionCostRequest struct {
	CostType    *string  `json:"cost_type" binding:"omitempty,oneof=TRAINER PROVIDER VENUE MATERIAL CERTIFICATION TRAVEL ACCOMMODATION OTHER"`
	Description *string  `json:"description"`
	Amount      *float64 `json:"amount" binding:"omitempty,min=0"`
}

type TrainingSessionCostResponse struct {
	ID          string    `json:"id"`
	SessionID   string    `json:"session_id"`
	CostType    string    `json:"cost_type"`
	Description string    `json:"description,omitempty"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// =========================================================================
// Training Document DTOs (P1-BE — plan §27)
// =========================================================================

type CreateTrainingDocumentRequest struct {
	DocumentType string `json:"document_type" binding:"required,oneof=PROPOSAL QUOTATION ATTENDANCE_SHEET INVOICE CONTRACT TRAINING_REPORT OTHER"`
	FileName     string `json:"file_name" binding:"omitempty,max=255"`
	FileURL      string `json:"file_url" binding:"required"`
}

type TrainingDocumentResponse struct {
	ID           string    `json:"id"`
	SessionID    string    `json:"session_id"`
	DocumentType string    `json:"document_type"`
	FileName     string    `json:"file_name,omitempty"`
	FileURL      string    `json:"file_url"`
	UploadedBy   string    `json:"uploaded_by,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
