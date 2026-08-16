package attendance

import "time"

// =========================================================================
// Business Travel — Request DTOs
// =========================================================================

// CreateBusinessTravelRequest — request perjalanan dinas (§4 plan doc).
// Sengaja TIDAK punya field funding_method/deposit/transfer — funding baru
// ditentukan setelah travel APPROVED.
type CreateBusinessTravelRequest struct {
	Title        string                                   `json:"title" binding:"required"`
	Purpose      string                                   `json:"purpose"`
	Description  string                                   `json:"description"`
	StartDate    string                                   `json:"start_date" binding:"required"`
	EndDate      string                                   `json:"end_date" binding:"required"`
	Origin       string                                   `json:"origin"`
	Participants []CreateBusinessTravelParticipantRequest `json:"participants"`
	Destinations []CreateBusinessTravelDestinationRequest `json:"destinations"`
}

type UpdateBusinessTravelRequest struct {
	Title       *string `json:"title"`
	Purpose     *string `json:"purpose"`
	Description *string `json:"description"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	Origin      *string `json:"origin"`
}

// SubmitBusinessTravelRequest — flow_id opsional, sama seperti overtime/leave;
// jika kosong, di-auto-resolve dari active flow module "business_travel".
type SubmitBusinessTravelRequest struct {
	FlowID *string `json:"flow_id"`
}

type CreateBusinessTravelParticipantRequest struct {
	ParticipantType string `json:"participant_type" binding:"required"`
	EmployeeID      string `json:"employee_id"`
	Name            string `json:"name"`
	Organization    string `json:"organization"`
	Position        string `json:"position"`
	IdentityNumber  string `json:"identity_number"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Role            string `json:"role"`
	Notes           string `json:"notes"`
}

type CreateBusinessTravelDestinationRequest struct {
	Sequence      int    `json:"sequence"`
	Country       string `json:"country"`
	Province      string `json:"province"`
	City          string `json:"city"`
	Location      string `json:"location"`
	ArrivalDate   string `json:"arrival_date"`
	DepartureDate string `json:"departure_date"`
	Purpose       string `json:"purpose"`
	Notes         string `json:"notes"`
}

type CreateBusinessTravelActivityRequest struct {
	ActivityDate string `json:"activity_date" binding:"required"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	Location     string `json:"location"`
	Organizer    string `json:"organizer"`
	Notes        string `json:"notes"`
}

type CreateBusinessTravelScheduleRequest struct {
	ScheduleType       string `json:"schedule_type" binding:"required"`
	DepartureDatetime  string `json:"departure_datetime"`
	ArrivalDatetime    string `json:"arrival_datetime"`
	Origin             string `json:"origin"`
	Destination        string `json:"destination"`
	TransportationType string `json:"transportation_type"`
	Provider           string `json:"provider"`
	BookingReference   string `json:"booking_reference"`
	Notes              string `json:"notes"`
}

// =========================================================================
// Business Travel — Response DTOs
// =========================================================================

type BusinessTravelResponse struct {
	ID                 string                              `json:"id"`
	RequestNumber      string                              `json:"request_number"`
	RequesterID        string                              `json:"requester_id"`
	Title              string                              `json:"title"`
	Purpose            *string                             `json:"purpose,omitempty"`
	Description        *string                             `json:"description,omitempty"`
	StartDate          string                              `json:"start_date"`
	EndDate            string                              `json:"end_date"`
	Origin             *string                             `json:"origin,omitempty"`
	Status             string                              `json:"status"`
	ApprovalStatus     string                              `json:"approval_status"`
	ApprovalInstanceID *string                             `json:"approval_instance_id,omitempty"`
	Participants       []BusinessTravelParticipantResponse `json:"participants,omitempty"`
	Destinations       []BusinessTravelDestinationResponse `json:"destinations,omitempty"`
	CreatedAt          time.Time                           `json:"created_at"`
	UpdatedAt          time.Time                           `json:"updated_at"`
}

type BusinessTravelParticipantResponse struct {
	ID              string  `json:"id"`
	ParticipantType string  `json:"participant_type"`
	EmployeeID      *string `json:"employee_id,omitempty"`
	Name            *string `json:"name,omitempty"`
	Organization    *string `json:"organization,omitempty"`
	Position        *string `json:"position,omitempty"`
	IdentityNumber  *string `json:"identity_number,omitempty"`
	Email           *string `json:"email,omitempty"`
	Phone           *string `json:"phone,omitempty"`
	Role            string  `json:"role"`
	Notes           *string `json:"notes,omitempty"`
}

type BusinessTravelDestinationResponse struct {
	ID            string  `json:"id"`
	Sequence      int     `json:"sequence"`
	Country       *string `json:"country,omitempty"`
	Province      *string `json:"province,omitempty"`
	City          *string `json:"city,omitempty"`
	Location      *string `json:"location,omitempty"`
	ArrivalDate   *string `json:"arrival_date,omitempty"`
	DepartureDate *string `json:"departure_date,omitempty"`
	Purpose       *string `json:"purpose,omitempty"`
	Notes         *string `json:"notes,omitempty"`
}

type BusinessTravelActivityResponse struct {
	ID           string  `json:"id"`
	ActivityDate string  `json:"activity_date"`
	StartTime    *string `json:"start_time,omitempty"`
	EndTime      *string `json:"end_time,omitempty"`
	Title        string  `json:"title"`
	Description  *string `json:"description,omitempty"`
	Location     *string `json:"location,omitempty"`
	Organizer    *string `json:"organizer,omitempty"`
	Notes        *string `json:"notes,omitempty"`
}

type BusinessTravelScheduleResponse struct {
	ID                 string     `json:"id"`
	ScheduleType       string     `json:"schedule_type"`
	DepartureDatetime  *time.Time `json:"departure_datetime,omitempty"`
	ArrivalDatetime    *time.Time `json:"arrival_datetime,omitempty"`
	Origin             *string    `json:"origin,omitempty"`
	Destination        *string    `json:"destination,omitempty"`
	TransportationType string     `json:"transportation_type"`
	Provider           *string    `json:"provider,omitempty"`
	BookingReference   *string    `json:"booking_reference,omitempty"`
	Notes              *string    `json:"notes,omitempty"`
}
