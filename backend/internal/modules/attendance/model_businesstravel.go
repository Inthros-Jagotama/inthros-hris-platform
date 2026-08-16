package attendance

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =========================================================================
// BusinessTravel Status
// =========================================================================

type TravelStatus string

const (
	TravelStatusDraft      TravelStatus = "DRAFT"
	TravelStatusSubmitted  TravelStatus = "SUBMITTED"
	TravelStatusApproved   TravelStatus = "APPROVED"
	TravelStatusRejected   TravelStatus = "REJECTED"
	TravelStatusInProgress TravelStatus = "IN_PROGRESS"
	TravelStatusCompleted  TravelStatus = "COMPLETED"
	TravelStatusClosed     TravelStatus = "CLOSED"
	TravelStatusCancelled  TravelStatus = "CANCELLED"
)

// =========================================================================
// BusinessTravel (core)
// =========================================================================

type BusinessTravel struct {
	ID                  uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	RequestNumber       string         `gorm:"type:varchar(50);not null;uniqueIndex:uq_biztrav_request_number" json:"request_number"`
	RequesterID         uuid.UUID      `gorm:"type:char(36);not null;index:idx_biztrav_requester" json:"requester_id"`
	Title               string         `gorm:"type:varchar(200);not null" json:"title"`
	Purpose             *string        `gorm:"type:varchar(500)" json:"purpose,omitempty"`
	Description         *string        `gorm:"type:text" json:"description,omitempty"`
	StartDate           time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate             time.Time      `gorm:"type:date;not null" json:"end_date"`
	Origin              *string        `gorm:"type:varchar(200)" json:"origin,omitempty"`
	Status              TravelStatus   `gorm:"type:varchar(30);not null;default:DRAFT;index:idx_biztrav_status" json:"status"`
	ApprovalStatus      string         `gorm:"type:varchar(30);not null;default:DRAFT" json:"approval_status"`
	ApprovalInstanceID  *uuid.UUID     `gorm:"type:char(36);index:idx_biztrav_approval_instance" json:"approval_instance_id,omitempty"`
	CreatedBy           *uuid.UUID     `gorm:"type:char(36)" json:"created_by,omitempty"`

	Participants []BusinessTravelParticipant `gorm:"foreignKey:BusinessTravelID" json:"participants,omitempty"`
	Destinations []BusinessTravelDestination `gorm:"foreignKey:BusinessTravelID" json:"destinations,omitempty"`
	Activities   []BusinessTravelActivity    `gorm:"foreignKey:BusinessTravelID" json:"activities,omitempty"`
	Schedules    []BusinessTravelSchedule    `gorm:"foreignKey:BusinessTravelID" json:"schedules,omitempty"`

	DeletedAt gorm.DeletedAt `gorm:"index:idx_biztrav_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (BusinessTravel) TableName() string {
	return "business_travels"
}

func (t *BusinessTravel) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Participants
// =========================================================================

type ParticipantType string

const (
	ParticipantTypeEmployee    ParticipantType = "EMPLOYEE"
	ParticipantTypeNonEmployee ParticipantType = "NON_EMPLOYEE"
)

type ParticipantRole string

const (
	ParticipantRoleLeader     ParticipantRole = "LEADER"
	ParticipantRoleMember     ParticipantRole = "MEMBER"
	ParticipantRoleDriver     ParticipantRole = "DRIVER"
	ParticipantRoleNarasumber ParticipantRole = "NARASUMBER"
	ParticipantRoleClient     ParticipantRole = "CLIENT"
	ParticipantRoleConsultant ParticipantRole = "CONSULTANT"
	ParticipantRoleOther      ParticipantRole = "OTHER"
)

type BusinessTravelParticipant struct {
	ID                uuid.UUID       `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessTravelID  uuid.UUID       `gorm:"type:char(36);not null;index:idx_biztrav_part_travel" json:"business_travel_id"`
	ParticipantType   ParticipantType `gorm:"type:varchar(20);not null;default:EMPLOYEE" json:"participant_type"`
	EmployeeID        *uuid.UUID      `gorm:"type:char(36);index:idx_biztrav_part_employee" json:"employee_id,omitempty"`
	Name              *string         `gorm:"type:varchar(150)" json:"name,omitempty"`
	Organization      *string         `gorm:"type:varchar(150)" json:"organization,omitempty"`
	Position          *string         `gorm:"type:varchar(150)" json:"position,omitempty"`
	IdentityNumber    *string         `gorm:"type:varchar(50)" json:"identity_number,omitempty"`
	Email             *string         `gorm:"type:varchar(150)" json:"email,omitempty"`
	Phone             *string         `gorm:"type:varchar(30)" json:"phone,omitempty"`
	Role              ParticipantRole `gorm:"type:varchar(30);not null;default:MEMBER" json:"role"`
	Notes             *string         `gorm:"type:varchar(500)" json:"notes,omitempty"`
	DeletedAt         gorm.DeletedAt  `gorm:"index:idx_biztrav_part_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

func (BusinessTravelParticipant) TableName() string {
	return "business_travel_participants"
}

func (p *BusinessTravelParticipant) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Destinations
// =========================================================================

type BusinessTravelDestination struct {
	ID               uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessTravelID uuid.UUID      `gorm:"type:char(36);not null;index:idx_biztrav_dest_travel" json:"business_travel_id"`
	Sequence         int            `gorm:"not null;default:1" json:"sequence"`
	Country          *string        `gorm:"type:varchar(100)" json:"country,omitempty"`
	Province         *string        `gorm:"type:varchar(100)" json:"province,omitempty"`
	City             *string        `gorm:"type:varchar(100)" json:"city,omitempty"`
	Location         *string        `gorm:"type:varchar(200)" json:"location,omitempty"`
	ArrivalDate      *time.Time     `gorm:"type:date" json:"arrival_date,omitempty"`
	DepartureDate    *time.Time     `gorm:"type:date" json:"departure_date,omitempty"`
	Purpose          *string        `gorm:"type:varchar(300)" json:"purpose,omitempty"`
	Notes            *string        `gorm:"type:varchar(500)" json:"notes,omitempty"`
	DeletedAt        gorm.DeletedAt `gorm:"index:idx_biztrav_dest_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func (BusinessTravelDestination) TableName() string {
	return "business_travel_destinations"
}

func (d *BusinessTravelDestination) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Activities
// =========================================================================

type BusinessTravelActivity struct {
	ID               uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessTravelID uuid.UUID      `gorm:"type:char(36);not null;index:idx_biztrav_act_travel" json:"business_travel_id"`
	ActivityDate     time.Time      `gorm:"type:date;not null" json:"activity_date"`
	StartTime        *string        `gorm:"type:time" json:"start_time,omitempty"`
	EndTime          *string        `gorm:"type:time" json:"end_time,omitempty"`
	Title            string         `gorm:"type:varchar(200);not null" json:"title"`
	Description      *string        `gorm:"type:text" json:"description,omitempty"`
	Location         *string        `gorm:"type:varchar(200)" json:"location,omitempty"`
	Organizer        *string        `gorm:"type:varchar(150)" json:"organizer,omitempty"`
	Notes            *string        `gorm:"type:varchar(500)" json:"notes,omitempty"`
	DeletedAt        gorm.DeletedAt `gorm:"index:idx_biztrav_act_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func (BusinessTravelActivity) TableName() string {
	return "business_travel_activities"
}

func (a *BusinessTravelActivity) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Schedules / Transportation
// =========================================================================

type ScheduleType string

const (
	ScheduleTypeDeparture ScheduleType = "DEPARTURE"
	ScheduleTypeReturn    ScheduleType = "RETURN"
	ScheduleTypeTransfer  ScheduleType = "TRANSFER"
	ScheduleTypeOther     ScheduleType = "OTHER"
)

type TransportationType string

const (
	TransportationAirplane   TransportationType = "AIRPLANE"
	TransportationTrain      TransportationType = "TRAIN"
	TransportationBus        TransportationType = "BUS"
	TransportationCar        TransportationType = "CAR"
	TransportationCompanyCar TransportationType = "COMPANY_CAR"
	TransportationRentalCar  TransportationType = "RENTAL_CAR"
	TransportationMotorcycle TransportationType = "MOTORCYCLE"
	TransportationOther      TransportationType = "OTHER"
)

type BusinessTravelSchedule struct {
	ID                  uuid.UUID           `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessTravelID    uuid.UUID           `gorm:"type:char(36);not null;index:idx_biztrav_sched_travel" json:"business_travel_id"`
	ScheduleType        ScheduleType        `gorm:"type:varchar(20);not null;default:DEPARTURE" json:"schedule_type"`
	DepartureDatetime   *time.Time          `json:"departure_datetime,omitempty"`
	ArrivalDatetime     *time.Time          `json:"arrival_datetime,omitempty"`
	Origin              *string             `gorm:"type:varchar(200)" json:"origin,omitempty"`
	Destination         *string             `gorm:"type:varchar(200)" json:"destination,omitempty"`
	TransportationType  TransportationType  `gorm:"type:varchar(30);not null;default:OTHER" json:"transportation_type"`
	Provider            *string             `gorm:"type:varchar(150)" json:"provider,omitempty"`
	BookingReference    *string             `gorm:"type:varchar(100)" json:"booking_reference,omitempty"`
	Notes               *string             `gorm:"type:varchar(500)" json:"notes,omitempty"`
	DeletedAt           gorm.DeletedAt      `gorm:"index:idx_biztrav_sched_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
}

func (BusinessTravelSchedule) TableName() string {
	return "business_travel_schedules"
}

func (s *BusinessTravelSchedule) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Expense Category (master)
// =========================================================================

type ExpenseCategory struct {
	ID                uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code              string         `gorm:"type:varchar(50);not null;uniqueIndex:uq_biztrav_expcat_code" json:"code"`
	Name              string         `gorm:"type:varchar(150);not null" json:"name"`
	Description       *string        `gorm:"type:varchar(500)" json:"description,omitempty"`
	RequiresReceipt   bool           `gorm:"not null;default:1" json:"requires_receipt"`
	Reimbursable      bool           `gorm:"not null;default:1" json:"reimbursable"`
	PayrollTreatment  *string        `gorm:"type:varchar(50)" json:"payroll_treatment,omitempty"`
	AccountCode       *string        `gorm:"type:varchar(50)" json:"account_code,omitempty"`
	Active            bool           `gorm:"not null;default:1;index:idx_biztrav_expcat_active" json:"active"`
	DeletedAt         gorm.DeletedAt `gorm:"index:idx_biztrav_expcat_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

func (ExpenseCategory) TableName() string {
	return "business_travel_expense_categories"
}

func (c *ExpenseCategory) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Expense Plan (estimasi biaya saat request)
// =========================================================================

type ExpensePlan struct {
	ID                 uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessTravelID   uuid.UUID      `gorm:"type:char(36);not null;index:idx_biztrav_expplan_travel" json:"business_travel_id"`
	ParticipantID      *uuid.UUID     `gorm:"type:char(36)" json:"participant_id,omitempty"`
	ExpenseCategoryID  uuid.UUID      `gorm:"type:char(36);not null;index:idx_biztrav_expplan_category" json:"expense_category_id"`
	Description        *string        `gorm:"type:varchar(300)" json:"description,omitempty"`
	Quantity           float64        `gorm:"type:decimal(10,2);not null;default:1" json:"quantity"`
	Unit               *string        `gorm:"type:varchar(30)" json:"unit,omitempty"`
	EstimatedAmount    float64        `gorm:"type:decimal(18,2);not null;default:0" json:"estimated_amount"`
	Notes              *string        `gorm:"type:varchar(500)" json:"notes,omitempty"`
	DeletedAt          gorm.DeletedAt `gorm:"index:idx_biztrav_expplan_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

func (ExpensePlan) TableName() string {
	return "business_travel_expense_plans"
}

func (p *ExpensePlan) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Funding Method (master)
// =========================================================================

type FundingMethod struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code        string         `gorm:"type:varchar(50);not null;uniqueIndex:uq_biztrav_fundmethod_code" json:"code"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description *string        `gorm:"type:varchar(300)" json:"description,omitempty"`
	Active      bool           `gorm:"not null;default:1" json:"active"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (FundingMethod) TableName() string {
	return "business_travel_funding_methods"
}

func (m *FundingMethod) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// Default funding method codes (seed) — §16 plan doc.
const (
	FundingMethodDeposit       = "DEPOSIT"
	FundingMethodReimbursement = "REIMBURSEMENT"
	FundingMethodCompanyPaid   = "COMPANY_PAID"
	FundingMethodOther         = "OTHER"
)

// =========================================================================
// Funding
// =========================================================================

type FundingStatus string

const (
	FundingStatusPending    FundingStatus = "PENDING"
	FundingStatusProcessing FundingStatus = "PROCESSING"
	FundingStatusFunded     FundingStatus = "FUNDED"
	FundingStatusCancelled  FundingStatus = "CANCELLED"
	FundingStatusReversed   FundingStatus = "REVERSED"
)

type Funding struct {
	ID                uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessTravelID  uuid.UUID      `gorm:"type:char(36);not null;index:idx_biztrav_funding_travel" json:"business_travel_id"`
	FundingMethodID   uuid.UUID      `gorm:"type:char(36);not null;index:idx_biztrav_funding_method" json:"funding_method_id"`
	ParticipantID     *uuid.UUID     `gorm:"type:char(36)" json:"participant_id,omitempty"`
	Amount            float64        `gorm:"type:decimal(18,2);not null;default:0" json:"amount"`
	FundingDate       *time.Time     `gorm:"type:date" json:"funding_date,omitempty"`
	PaymentMethod     *string        `gorm:"type:varchar(50)" json:"payment_method,omitempty"`
	PaymentReference  *string        `gorm:"type:varchar(100)" json:"payment_reference,omitempty"`
	FundedBy          *uuid.UUID     `gorm:"type:char(36)" json:"funded_by,omitempty"`
	Status            FundingStatus  `gorm:"type:varchar(30);not null;default:PENDING;index:idx_biztrav_funding_status" json:"status"`
	Notes             *string        `gorm:"type:varchar(500)" json:"notes,omitempty"`

	Documents []FundingDocument `gorm:"foreignKey:BusinessTravelFundingID" json:"documents,omitempty"`

	DeletedAt gorm.DeletedAt `gorm:"index:idx_biztrav_funding_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (Funding) TableName() string {
	return "business_travel_fundings"
}

func (f *Funding) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Funding Document (bukti transfer)
// =========================================================================

type FundingDocumentType string

const (
	FundingDocTransferReceipt FundingDocumentType = "TRANSFER_RECEIPT"
	FundingDocPaymentProof    FundingDocumentType = "PAYMENT_PROOF"
	FundingDocOther           FundingDocumentType = "OTHER"
)

type FundingDocument struct {
	ID                         uuid.UUID           `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessTravelFundingID    uuid.UUID           `gorm:"type:char(36);not null;index:idx_biztrav_funddoc_funding" json:"business_travel_funding_id"`
	DocumentType               FundingDocumentType `gorm:"type:varchar(30);not null;default:OTHER" json:"document_type"`
	FileName                   string              `gorm:"type:varchar(255);not null" json:"file_name"`
	FilePath                   string              `gorm:"type:text;not null" json:"file_path"`
	MimeType                   *string             `gorm:"type:varchar(100)" json:"mime_type,omitempty"`
	FileSize                   *int64              `json:"file_size,omitempty"`
	UploadedBy                 *uuid.UUID          `gorm:"type:char(36)" json:"uploaded_by,omitempty"`
	UploadedAt                 *time.Time          `json:"uploaded_at,omitempty"`
	DeletedAt                  gorm.DeletedAt      `json:"deleted_at,omitempty"`
	CreatedAt                  time.Time           `json:"created_at"`
	UpdatedAt                  time.Time           `json:"updated_at"`
}

func (FundingDocument) TableName() string {
	return "business_travel_funding_documents"
}

func (d *FundingDocument) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Expense (actual)
// =========================================================================

type ExpenseStatus string

const (
	ExpenseStatusDraft     ExpenseStatus = "DRAFT"
	ExpenseStatusSubmitted ExpenseStatus = "SUBMITTED"
	ExpenseStatusApproved  ExpenseStatus = "APPROVED"
	ExpenseStatusRejected  ExpenseStatus = "REJECTED"
)

type Expense struct {
	ID                 uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessTravelID   uuid.UUID      `gorm:"type:char(36);not null;index:idx_biztrav_exp_travel" json:"business_travel_id"`
	ParticipantID      *uuid.UUID     `gorm:"type:char(36)" json:"participant_id,omitempty"`
	ExpenseCategoryID  uuid.UUID      `gorm:"type:char(36);not null;index:idx_biztrav_exp_category" json:"expense_category_id"`
	ExpenseDate        time.Time      `gorm:"type:date;not null" json:"expense_date"`
	Description        *string        `gorm:"type:varchar(300)" json:"description,omitempty"`
	Quantity           float64        `gorm:"type:decimal(10,2);not null;default:1" json:"quantity"`
	Unit               *string        `gorm:"type:varchar(30)" json:"unit,omitempty"`
	Amount             float64        `gorm:"type:decimal(18,2);not null;default:0" json:"amount"`
	FundingMethodID    *uuid.UUID     `gorm:"type:char(36);index:idx_biztrav_exp_funding_method" json:"funding_method_id,omitempty"`
	Vendor             *string        `gorm:"type:varchar(150)" json:"vendor,omitempty"`
	ReceiptNumber      *string        `gorm:"type:varchar(100)" json:"receipt_number,omitempty"`
	Status             ExpenseStatus  `gorm:"type:varchar(30);not null;default:DRAFT" json:"status"`
	Notes              *string        `gorm:"type:varchar(500)" json:"notes,omitempty"`

	Documents []ExpenseDocument `gorm:"foreignKey:BusinessTravelExpenseID" json:"documents,omitempty"`

	DeletedAt gorm.DeletedAt `gorm:"index:idx_biztrav_exp_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (Expense) TableName() string {
	return "business_travel_expenses"
}

func (e *Expense) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Expense Document (bukti pengeluaran)
// =========================================================================

type ExpenseDocumentType string

const (
	ExpenseDocReceipt      ExpenseDocumentType = "RECEIPT"
	ExpenseDocInvoice      ExpenseDocumentType = "INVOICE"
	ExpenseDocTicket       ExpenseDocumentType = "TICKET"
	ExpenseDocHotelBill    ExpenseDocumentType = "HOTEL_BILL"
	ExpenseDocBoardingPass ExpenseDocumentType = "BOARDING_PASS"
	ExpenseDocTollReceipt  ExpenseDocumentType = "TOLL_RECEIPT"
	ExpenseDocOther        ExpenseDocumentType = "OTHER"
)

type ExpenseDocument struct {
	ID                        uuid.UUID           `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessTravelExpenseID   uuid.UUID           `gorm:"type:char(36);not null;index:idx_biztrav_expdoc_expense" json:"business_travel_expense_id"`
	DocumentType              ExpenseDocumentType `gorm:"type:varchar(30);not null;default:RECEIPT" json:"document_type"`
	FileName                  string              `gorm:"type:varchar(255);not null" json:"file_name"`
	FilePath                  string              `gorm:"type:text;not null" json:"file_path"`
	MimeType                  *string             `gorm:"type:varchar(100)" json:"mime_type,omitempty"`
	FileSize                  *int64              `json:"file_size,omitempty"`
	UploadedBy                *uuid.UUID          `gorm:"type:char(36)" json:"uploaded_by,omitempty"`
	UploadedAt                *time.Time          `json:"uploaded_at,omitempty"`
	DeletedAt                 gorm.DeletedAt      `json:"deleted_at,omitempty"`
	CreatedAt                 time.Time           `json:"created_at"`
	UpdatedAt                 time.Time           `json:"updated_at"`
}

func (ExpenseDocument) TableName() string {
	return "business_travel_expense_documents"
}

func (d *ExpenseDocument) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Document (bukti perjalanan, level travel)
// =========================================================================

type TravelDocumentType string

const (
	TravelDocOrder            TravelDocumentType = "TRAVEL_ORDER"
	TravelDocInvitation       TravelDocumentType = "INVITATION"
	TravelDocTicket           TravelDocumentType = "TICKET"
	TravelDocBoardingPass     TravelDocumentType = "BOARDING_PASS"
	TravelDocHotel            TravelDocumentType = "HOTEL"
	TravelDocMeetingDocument  TravelDocumentType = "MEETING_DOCUMENT"
	TravelDocAttendanceProof  TravelDocumentType = "ATTENDANCE_PROOF"
	TravelDocPhoto            TravelDocumentType = "PHOTO"
	TravelDocTravelReport     TravelDocumentType = "TRAVEL_REPORT"
	TravelDocOther            TravelDocumentType = "OTHER"
)

type TravelDocument struct {
	ID                uuid.UUID           `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessTravelID  uuid.UUID           `gorm:"type:char(36);not null;index:idx_biztrav_doc_travel" json:"business_travel_id"`
	DocumentType      TravelDocumentType  `gorm:"type:varchar(30);not null;default:OTHER" json:"document_type"`
	FileName          string              `gorm:"type:varchar(255);not null" json:"file_name"`
	FilePath          string              `gorm:"type:text;not null" json:"file_path"`
	MimeType          *string             `gorm:"type:varchar(100)" json:"mime_type,omitempty"`
	FileSize          *int64              `json:"file_size,omitempty"`
	UploadedBy        *uuid.UUID          `gorm:"type:char(36)" json:"uploaded_by,omitempty"`
	UploadedAt        *time.Time          `json:"uploaded_at,omitempty"`
	DeletedAt         gorm.DeletedAt      `gorm:"index:idx_biztrav_doc_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
}

func (TravelDocument) TableName() string {
	return "business_travel_documents"
}

func (d *TravelDocument) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Settlement
// =========================================================================

type SettlementStatus string

const (
	SettlementStatusPending                 SettlementStatus = "PENDING"
	SettlementStatusSubmitted               SettlementStatus = "SUBMITTED"
	SettlementStatusApproved                SettlementStatus = "APPROVED"
	SettlementStatusRefundRequired          SettlementStatus = "REFUND_REQUIRED"
	SettlementStatusReimbursementRequired   SettlementStatus = "REIMBURSEMENT_REQUIRED"
	SettlementStatusBalanced                SettlementStatus = "BALANCED"
	SettlementStatusRefundProcessing        SettlementStatus = "REFUND_PROCESSING"
	SettlementStatusReimbursementProcessing SettlementStatus = "REIMBURSEMENT_PROCESSING"
	SettlementStatusSettled                 SettlementStatus = "SETTLED"
	SettlementStatusRejected                SettlementStatus = "REJECTED"
)

type Settlement struct {
	ID                    uuid.UUID        `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessTravelID      uuid.UUID        `gorm:"type:char(36);not null;index:idx_biztrav_settle_travel" json:"business_travel_id"`
	ParticipantID         *uuid.UUID       `gorm:"type:char(36);index:idx_biztrav_settle_participant" json:"participant_id,omitempty"`
	TotalAdvance          float64          `gorm:"type:decimal(18,2);not null;default:0" json:"total_advance"`
	TotalActualExpense    float64          `gorm:"type:decimal(18,2);not null;default:0" json:"total_actual_expense"`
	TotalCompanyPaid      float64          `gorm:"type:decimal(18,2);not null;default:0" json:"total_company_paid"`
	TotalReimbursement    float64          `gorm:"type:decimal(18,2);not null;default:0" json:"total_reimbursement"`
	TotalRefund           float64          `gorm:"type:decimal(18,2);not null;default:0" json:"total_refund"`
	Balance               float64          `gorm:"type:decimal(18,2);not null;default:0" json:"balance"`
	Status                SettlementStatus `gorm:"type:varchar(30);not null;default:PENDING;index:idx_biztrav_settle_status" json:"status"`
	ApprovalInstanceID    *uuid.UUID       `gorm:"type:char(36);index:idx_biztrav_settle_approval_instance" json:"approval_instance_id,omitempty"`
	SubmittedAt           *time.Time       `json:"submitted_at,omitempty"`
	ApprovedAt            *time.Time       `json:"approved_at,omitempty"`
	SettledAt             *time.Time       `json:"settled_at,omitempty"`
	Notes                 *string          `gorm:"type:varchar(500)" json:"notes,omitempty"`

	Items []SettlementItem `gorm:"foreignKey:BusinessTravelSettlementID" json:"items,omitempty"`

	DeletedAt gorm.DeletedAt `gorm:"index:idx_biztrav_settle_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (Settlement) TableName() string {
	return "business_travel_settlements"
}

func (s *Settlement) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Settlement Item (breakdown per expense/funding method)
// =========================================================================

type SettlementItemType string

const (
	SettlementItemAdvance      SettlementItemType = "ADVANCE"
	SettlementItemActual       SettlementItemType = "ACTUAL"
	SettlementItemCompanyPaid  SettlementItemType = "COMPANY_PAID"
	SettlementItemRefund       SettlementItemType = "REFUND"
	SettlementItemReimbursement SettlementItemType = "REIMBURSEMENT"
)

type SettlementItem struct {
	ID                          uuid.UUID           `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessTravelSettlementID  uuid.UUID           `gorm:"type:char(36);not null;index:idx_biztrav_settleitem_settlement" json:"business_travel_settlement_id"`
	ExpenseID                   *uuid.UUID          `gorm:"type:char(36);index:idx_biztrav_settleitem_expense" json:"expense_id,omitempty"`
	FundingMethodID              *uuid.UUID          `gorm:"type:char(36)" json:"funding_method_id,omitempty"`
	ItemType                     SettlementItemType  `gorm:"type:varchar(30);not null;default:ACTUAL" json:"item_type"`
	Category                     *string             `gorm:"type:varchar(100)" json:"category,omitempty"`
	Amount                       float64             `gorm:"type:decimal(18,2);not null;default:0" json:"amount"`
	Notes                        *string             `gorm:"type:varchar(500)" json:"notes,omitempty"`
	DeletedAt                    gorm.DeletedAt      `json:"deleted_at,omitempty"`
	CreatedAt                    time.Time           `json:"created_at"`
	UpdatedAt                    time.Time           `json:"updated_at"`
}

func (SettlementItem) TableName() string {
	return "business_travel_settlement_items"
}

func (i *SettlementItem) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Refund
// =========================================================================

type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "PENDING"
	RefundStatusConfirmed RefundStatus = "CONFIRMED"
	RefundStatusCancelled RefundStatus = "CANCELLED"
)

type Refund struct {
	ID                uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessTravelID  uuid.UUID      `gorm:"type:char(36);not null;index:idx_biztrav_refund_travel" json:"business_travel_id"`
	SettlementID      *uuid.UUID     `gorm:"type:char(36);index:idx_biztrav_refund_settlement" json:"settlement_id,omitempty"`
	ParticipantID     *uuid.UUID     `gorm:"type:char(36)" json:"participant_id,omitempty"`
	RefundAmount      float64        `gorm:"type:decimal(18,2);not null;default:0" json:"refund_amount"`
	RefundDate        *time.Time     `gorm:"type:date" json:"refund_date,omitempty"`
	RefundReference   *string        `gorm:"type:varchar(100)" json:"refund_reference,omitempty"`
	RefundedBy        *uuid.UUID     `gorm:"type:char(36)" json:"refunded_by,omitempty"`
	RefundDocument    *string        `gorm:"type:text" json:"refund_document,omitempty"`
	Status            RefundStatus   `gorm:"type:varchar(30);not null;default:PENDING;index:idx_biztrav_refund_status" json:"status"`
	Notes             *string        `gorm:"type:varchar(500)" json:"notes,omitempty"`
	DeletedAt         gorm.DeletedAt `gorm:"index:idx_biztrav_refund_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

func (Refund) TableName() string {
	return "business_travel_refunds"
}

func (r *Refund) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Reimbursement (internal — additional claim hasil settlement)
// Integrasi ke module Reimbursement standalone hanya jika tenant subscribe
// module tersebut (lihat docs/module-attendance-business-travel-development-plan.md §54.7).
// =========================================================================

type TravelReimbursementStatus string

const (
	TravelReimbStatusRequested TravelReimbursementStatus = "REQUESTED"
	TravelReimbStatusApproved  TravelReimbursementStatus = "APPROVED"
	TravelReimbStatusProcessing TravelReimbursementStatus = "PROCESSING"
	TravelReimbStatusPaid      TravelReimbursementStatus = "PAID"
	TravelReimbStatusRejected  TravelReimbursementStatus = "REJECTED"
	TravelReimbStatusCancelled TravelReimbursementStatus = "CANCELLED"
)

type TravelReimbursement struct {
	ID                uuid.UUID                  `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessTravelID  uuid.UUID                  `gorm:"type:char(36);not null;index:idx_biztrav_reimb_travel" json:"business_travel_id"`
	ParticipantID     *uuid.UUID                 `gorm:"type:char(36)" json:"participant_id,omitempty"`
	SettlementID      *uuid.UUID                 `gorm:"type:char(36);index:idx_biztrav_reimb_settlement" json:"settlement_id,omitempty"`
	Amount            float64                    `gorm:"type:decimal(18,2);not null;default:0" json:"amount"`
	Status            TravelReimbursementStatus  `gorm:"type:varchar(30);not null;default:REQUESTED;index:idx_biztrav_reimb_status" json:"status"`
	RequestedAt       *time.Time                 `json:"requested_at,omitempty"`
	ApprovedAt        *time.Time                 `json:"approved_at,omitempty"`
	PaidAt            *time.Time                 `json:"paid_at,omitempty"`
	PaymentReference  *string                    `gorm:"type:varchar(100)" json:"payment_reference,omitempty"`
	PaidBy            *uuid.UUID                 `gorm:"type:char(36)" json:"paid_by,omitempty"`
	Notes             *string                    `gorm:"type:varchar(500)" json:"notes,omitempty"`
	DeletedAt         gorm.DeletedAt             `gorm:"index:idx_biztrav_reimb_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"updated_at"`
}

func (TravelReimbursement) TableName() string {
	return "business_travel_reimbursements"
}

func (r *TravelReimbursement) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Audit Log
// =========================================================================

type AuditLog struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	EntityType string    `gorm:"type:varchar(50);not null" json:"entity_type"`
	EntityID   uuid.UUID `gorm:"type:char(36);not null" json:"entity_id"`
	Action     string    `gorm:"type:varchar(50);not null" json:"action"`
	OldValue   *string   `gorm:"type:text" json:"old_value,omitempty"`
	NewValue   *string   `gorm:"type:text" json:"new_value,omitempty"`
	UserID     *uuid.UUID `gorm:"type:char(36)" json:"user_id,omitempty"`
	IPAddress  *string   `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	CreatedAt  time.Time `gorm:"index:idx_biztrav_audit_created_at" json:"created_at"`
}

func (AuditLog) TableName() string {
	return "business_travel_audit_logs"
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// Attendance Rule (§37 plan doc)
// =========================================================================

type AttendanceRuleType string

const (
	AttendanceRuleFullDay      AttendanceRuleType = "FULL_DAY"
	AttendanceRuleHalfDay      AttendanceRuleType = "HALF_DAY"
	AttendanceRuleTravelDay    AttendanceRuleType = "TRAVEL_DAY"
	AttendanceRuleNonWorkingDay AttendanceRuleType = "NON_WORKING_DAY"
)

type AttendanceRule struct {
	ID                uuid.UUID          `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessTravelID  *uuid.UUID         `gorm:"type:char(36);index:idx_biztrav_attrule_travel" json:"business_travel_id,omitempty"`
	RuleType          AttendanceRuleType `gorm:"type:varchar(30);not null;default:FULL_DAY" json:"rule_type"`
	Description       *string            `gorm:"type:varchar(300)" json:"description,omitempty"`
	IsDefault         bool               `gorm:"not null;default:0" json:"is_default"`
	Active            bool               `gorm:"not null;default:1;index:idx_biztrav_attrule_active" json:"active"`
	DeletedAt         gorm.DeletedAt     `gorm:"index:idx_biztrav_attrule_deleted_at" json:"deleted_at,omitempty"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

func (AttendanceRule) TableName() string {
	return "business_travel_attendance_rules"
}

func (r *AttendanceRule) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
