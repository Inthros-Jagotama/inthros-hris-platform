package training

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

// ApprovalEngine — narrow interface ke Central Approval Engine (pola leave/
// reimbursement/attendance/employeemovement — plan §15/§45).
type ApprovalEngine interface {
	CreateApprovalInstance(ctx context.Context, module, documentID, flowID string) (string, error)
	GetApprovalInstanceStatus(ctx context.Context, instanceID string) (string, error)
	// GetActiveFlowIDForModule lets a training request auto-resolve which flow
	// to route through when the client doesn't supply one explicitly.
	GetActiveFlowIDForModule(ctx context.Context, module string) (string, error)
}

// Notifier — narrow interface ke module notification (pola leave).
type Notifier interface {
	Notify(ctx context.Context, recipientUserID uuid.UUID, notifType string, params []string, referenceType string, referenceID uuid.UUID) error
}

type Service struct {
	repo           *Repository
	logger         *zap.Logger
	approvalEngine ApprovalEngine
	notifier       Notifier
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// SetApprovalEngine wires the central approval module into this service.
func (s *Service) SetApprovalEngine(ae ApprovalEngine) {
	s.approvalEngine = ae
}

// SetNotifier wires the notification module into this service.
func (s *Service) SetNotifier(n Notifier) {
	s.notifier = n
}

// =========================================================================
// Training Categories
// =========================================================================

func (s *Service) CreateCategory(ctx context.Context, req CreateTrainingCategoryRequest) (*TrainingCategoryResponse, error) {
	// Kode kategori: bila tidak dikirim, di-generate otomatis CAT-{sekuens}.
	code := strings.TrimSpace(req.Code)
	if code == "" {
		var err error
		code, err = s.repo.NextCategoryCode(ctx)
		if err != nil {
			return nil, err
		}
	}

	cat := &TrainingCategory{
		Code:     code,
		Name:     req.Name,
		IsActive: true,
	}
	if req.Description != nil {
		cat.Description = req.Description
	}
	if req.IsActive != nil {
		cat.IsActive = *req.IsActive
	}
	if err := s.repo.CreateCategory(ctx, cat); err != nil {
		return nil, err
	}
	s.logger.Info("Training category created", zap.String("id", cat.ID.String()), zap.String("code", cat.Code))
	return categoryToResponse(cat), nil
}

func (s *Service) GetCategoryByID(ctx context.Context, id string) (*TrainingCategoryResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	cat, err := s.repo.FindCategoryByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return categoryToResponse(cat), nil
}

func (s *Service) ListCategories(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	cats, total, err := s.repo.ListCategories(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingCategoryResponse, 0, len(cats))
	for _, c := range cats {
		responses = append(responses, *categoryToResponse(&c))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateCategory(ctx context.Context, id string, req UpdateTrainingCategoryRequest) (*TrainingCategoryResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	cat, err := s.repo.FindCategoryByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Code != nil {
		cat.Code = *req.Code
	}
	if req.Name != nil {
		cat.Name = *req.Name
	}
	if req.Description != nil {
		cat.Description = req.Description
	}
	if req.IsActive != nil {
		cat.IsActive = *req.IsActive
	}
	if err := s.repo.UpdateCategory(ctx, cat); err != nil {
		return nil, err
	}
	return categoryToResponse(cat), nil
}

func (s *Service) DeleteCategory(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCategory(ctx, uid)
}

// =========================================================================
// Training Courses
// =========================================================================

func (s *Service) CreateCourse(ctx context.Context, req CreateTrainingCourseRequest) (*TrainingCourseResponse, error) {
	catID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category_id: %w", err)
	}
	// Verify category exists (kode kategorinya dipakai sebagai prefix kode kursus)
	cat, err := s.repo.FindCategoryByID(ctx, catID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	// Kode kursus: bila tidak dikirim, di-generate otomatis {KODE_KATEGORI}-{sekuens}.
	code := strings.TrimSpace(req.Code)
	if code == "" {
		code, err = s.repo.NextCourseCode(ctx, cat.Code)
		if err != nil {
			return nil, err
		}
	}

	course := &TrainingCourse{
		CategoryID:  catID,
		Code:        code,
		Name:        req.Name,
		IsCertified: false,
		IsActive:    true,
	}
	if req.Description != nil {
		course.Description = req.Description
	}
	if req.DurationHour != nil {
		course.DurationHour = req.DurationHour
	}
	if req.MinScore != nil {
		course.MinScore = req.MinScore
	}
	if req.Cost != nil {
		course.Cost = req.Cost
	}
	if req.IsCertified != nil {
		course.IsCertified = *req.IsCertified
	}
	if req.ExternalVendor != nil {
		course.ExternalVendor = req.ExternalVendor
	}
	// P0-BE (plan §7)
	if req.CourseType != nil {
		ct := CourseType(*req.CourseType)
		course.CourseType = &ct
	}
	if req.DeliveryType != nil {
		dt := DeliveryType(*req.DeliveryType)
		course.DeliveryType = &dt
	}
	if req.IsMandatory != nil {
		course.IsMandatory = *req.IsMandatory
	}

	if err := s.repo.CreateCourse(ctx, course); err != nil {
		return nil, err
	}
	s.logger.Info("Training course created", zap.String("id", course.ID.String()), zap.String("code", course.Code))
	return courseToResponse(course), nil
}

func (s *Service) GetCourseByID(ctx context.Context, id string) (*TrainingCourseResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	course, err := s.repo.FindCourseByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return courseToResponse(course), nil
}

func (s *Service) ListCourses(ctx context.Context, categoryID *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var catUUID *uuid.UUID
	if categoryID != nil && *categoryID != "" {
		uid, err := uuid.Parse(*categoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category_id: %w", err)
		}
		catUUID = &uid
	}
	courses, total, err := s.repo.ListCourses(ctx, catUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingCourseResponse, 0, len(courses))
	for _, c := range courses {
		responses = append(responses, *courseToResponse(&c))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateCourse(ctx context.Context, id string, req UpdateTrainingCourseRequest) (*TrainingCourseResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	course, err := s.repo.FindCourseByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.CategoryID != nil {
		catID, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category_id: %w", err)
		}
		course.CategoryID = catID
	}
	if req.Code != nil {
		course.Code = *req.Code
	}
	if req.Name != nil {
		course.Name = *req.Name
	}
	if req.Description != nil {
		course.Description = req.Description
	}
	if req.DurationHour != nil {
		course.DurationHour = req.DurationHour
	}
	if req.MinScore != nil {
		course.MinScore = req.MinScore
	}
	if req.Cost != nil {
		course.Cost = req.Cost
	}
	if req.IsCertified != nil {
		course.IsCertified = *req.IsCertified
	}
	if req.ExternalVendor != nil {
		course.ExternalVendor = req.ExternalVendor
	}
	if req.IsActive != nil {
		course.IsActive = *req.IsActive
	}
	// P0-BE (plan §7)
	if req.CourseType != nil {
		ct := CourseType(*req.CourseType)
		course.CourseType = &ct
	}
	if req.DeliveryType != nil {
		dt := DeliveryType(*req.DeliveryType)
		course.DeliveryType = &dt
	}
	if req.IsMandatory != nil {
		course.IsMandatory = *req.IsMandatory
	}
	if err := s.repo.UpdateCourse(ctx, course); err != nil {
		return nil, err
	}
	return courseToResponse(course), nil
}

func (s *Service) DeleteCourse(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCourse(ctx, uid)
}

// =========================================================================
// Training Sessions
// =========================================================================

func (s *Service) CreateSession(ctx context.Context, req CreateTrainingSessionRequest) (*TrainingSessionResponse, error) {
	courseID, err := uuid.Parse(req.CourseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course_id: %w", err)
	}
	// Verify course exists
	course, err := s.repo.FindCourseByID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	// Kode sesi: bila tidak dikirim, di-generate otomatis {KODE_KURSUS}-{sekuens}.
	sessionCode := strings.TrimSpace(req.SessionCode)
	if sessionCode == "" {
		sessionCode, err = s.repo.NextSessionCode(ctx, course.Code)
		if err != nil {
			return nil, err
		}
	}

	maxQuota := 30
	if req.MaxQuota > 0 {
		maxQuota = req.MaxQuota
	}

	sess := &TrainingSession{
		CourseID:    courseID,
		SessionCode: sessionCode,
		TrainerName: req.TrainerName,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		MaxQuota:    maxQuota,
		Status:      SessStatusScheduled,
	}
	if req.Location != nil {
		sess.Location = req.Location
	}
	if err := applySessionEnhancement(sess, req.ProviderType, req.DeliveryMode, req.ProviderID,
		req.StartDatetime, req.EndDatetime, req.MeetingURL, req.RegistrationDeadline); err != nil {
		return nil, err
	}

	if err := s.repo.CreateSession(ctx, sess); err != nil {
		return nil, err
	}
	s.logger.Info("Training session created", zap.String("id", sess.ID.String()), zap.String("code", sess.SessionCode))
	return sessionToResponse(sess), nil
}

func (s *Service) GetSessionByID(ctx context.Context, id string) (*TrainingSessionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	sess, err := s.repo.FindSessionByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return sessionToResponse(sess), nil
}

func (s *Service) ListSessions(ctx context.Context, courseID *string, status *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var courseUUID *uuid.UUID
	if courseID != nil && *courseID != "" {
		uid, err := uuid.Parse(*courseID)
		if err != nil {
			return nil, fmt.Errorf("invalid course_id: %w", err)
		}
		courseUUID = &uid
	}
	sessions, total, err := s.repo.ListSessions(ctx, courseUUID, status, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingSessionResponse, 0, len(sessions))
	for _, s := range sessions {
		responses = append(responses, *sessionToResponse(&s))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateSession(ctx context.Context, id string, req UpdateTrainingSessionRequest) (*TrainingSessionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	sess, err := s.repo.FindSessionByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.SessionCode != nil {
		sess.SessionCode = *req.SessionCode
	}
	if req.TrainerName != nil {
		sess.TrainerName = *req.TrainerName
	}
	if req.Location != nil {
		sess.Location = req.Location
	}
	if req.StartDate != nil {
		sess.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		sess.EndDate = *req.EndDate
	}
	if req.MaxQuota != nil {
		sess.MaxQuota = *req.MaxQuota
	}
	// P0-BE (plan §14) — hanya field yang dikirim yang di-update.
	if err := applySessionEnhancement(sess, req.ProviderType, req.DeliveryMode, req.ProviderID,
		req.StartDatetime, req.EndDatetime, req.MeetingURL, req.RegistrationDeadline); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateSession(ctx, sess); err != nil {
		return nil, err
	}
	return sessionToResponse(sess), nil
}

func (s *Service) UpdateSessionStatus(ctx context.Context, id, status string) (*TrainingSessionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	sess, err := s.repo.FindSessionByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	sess.Status = SessionStatus(status)
	if err := s.repo.UpdateSession(ctx, sess); err != nil {
		return nil, err
	}
	s.logger.Info("Training session status updated", zap.String("id", id), zap.String("status", status))
	return sessionToResponse(sess), nil
}

func (s *Service) DeleteSession(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteSession(ctx, uid)
}

// =========================================================================
// Training Participants
// =========================================================================

func (s *Service) CreateParticipant(ctx context.Context, req CreateTrainingParticipantRequest) (*TrainingParticipantResponse, error) {
	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}

	// Verify session exists and check quota
	sess, err := s.repo.FindSessionByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}
	if sess.Status == SessStatusCancelled {
		return nil, fmt.Errorf("cannot register to cancelled session")
	}
	// P0-BE (plan §18): cegah employee terdaftar 2x di session yang sama.
	existing, err := s.repo.FindParticipantBySessionAndEmployee(ctx, sessionID, empID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("employee already registered in this session")
	}
	count, err := s.repo.CountParticipantsBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if count >= int64(sess.MaxQuota) {
		return nil, fmt.Errorf("session quota full (%d/%d)", count, sess.MaxQuota)
	}

	p := &TrainingParticipant{
		SessionID:          sessionID,
		EmployeeID:         empID,
		AttendanceStatus:   AttendStatusPresent,
		RegistrationStatus: RegStatusRegistered,
	}
	now := time.Now()
	if req.RegistrationStatus != nil {
		p.RegistrationStatus = RegistrationStatus(*req.RegistrationStatus)
	}
	// Status yang menandakan peserta terdaftar → set registered_at.
	if p.RegistrationStatus == RegStatusRegistered || p.RegistrationStatus == RegStatusApproved {
		p.RegisteredAt = &now
	}
	if p.RegistrationStatus == RegStatusApproved {
		p.ApprovedAt = &now
	}
	if err := s.repo.CreateParticipant(ctx, p); err != nil {
		return nil, err
	}
	s.logger.Info("Participant registered for training", zap.String("session_id", req.SessionID), zap.String("employee_id", req.EmployeeID))
	return participantToResponse(p), nil
}

func (s *Service) GetParticipantByID(ctx context.Context, id string) (*TrainingParticipantResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindParticipantByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return participantToResponse(p), nil
}

func (s *Service) ListParticipants(ctx context.Context, sessionID *string, employeeID *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var sessUUID, empUUID *uuid.UUID
	if sessionID != nil && *sessionID != "" {
		uid, err := uuid.Parse(*sessionID)
		if err != nil {
			return nil, fmt.Errorf("invalid session_id: %w", err)
		}
		sessUUID = &uid
	}
	if employeeID != nil && *employeeID != "" {
		uid, err := uuid.Parse(*employeeID)
		if err != nil {
			return nil, fmt.Errorf("invalid employee_id: %w", err)
		}
		empUUID = &uid
	}
	participants, total, err := s.repo.ListParticipants(ctx, sessUUID, empUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingParticipantResponse, 0, len(participants))
	for _, p := range participants {
		responses = append(responses, *participantToResponse(&p))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateParticipant(ctx context.Context, id string, req UpdateTrainingParticipantRequest) (*TrainingParticipantResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindParticipantByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.AttendanceStatus != nil {
		p.AttendanceStatus = AttendanceStatus(*req.AttendanceStatus)
	}
	if req.Score != nil {
		p.Score = *req.Score
		// Auto-mark completion with today's date when score is entered.
		completed := time.Now().Format("2006-01-02")
		p.CompletedAt = &completed
		p.CompletionDate = &completed
		score := *req.Score
		p.FinalScore = &score
		// Lulus jika >= min_score course (default lulus saat course tak punya min_score).
		passed := true
		if sess, err := s.repo.FindSessionByID(ctx, p.SessionID); err == nil && sess != nil {
			if course, err := s.repo.FindCourseByID(ctx, sess.CourseID); err == nil && course != nil && course.MinScore != nil {
				passed = score >= *course.MinScore
			}
		}
		p.Passed = &passed
		// Status completion sinkron dengan hasil: FAILED bila tidak lulus.
		if passed {
			p.CompletionStatus = CompletionCompleted
		} else {
			p.CompletionStatus = CompletionFailed
		}
	}
	// P0-BE (plan §18): field completion opsional.
	if req.CompletionStatus != nil {
		p.CompletionStatus = CompletionStatus(*req.CompletionStatus)
	}
	if req.FinalScore != nil {
		p.FinalScore = req.FinalScore
	}
	if req.Passed != nil {
		p.Passed = req.Passed
	}
	if req.Remarks != nil {
		p.Remarks = req.Remarks
	}
	if err := s.repo.UpdateParticipant(ctx, p); err != nil {
		return nil, err
	}
	return participantToResponse(p), nil
}

func (s *Service) DeleteParticipant(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteParticipant(ctx, uid)
}

// =========================================================================
// Training Materials
// =========================================================================

func (s *Service) CreateMaterial(ctx context.Context, req CreateTrainingMaterialRequest) (*TrainingMaterialResponse, error) {
	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}
	// Verify session exists
	if _, err := s.repo.FindSessionByID(ctx, sessionID); err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	m := &TrainingMaterial{
		SessionID: sessionID,
		Title:     req.Title,
		FileURL:   req.FileURL,
		FileType:  req.FileType,
	}
	if req.SortOrder != nil {
		m.SortOrder = *req.SortOrder
	}
	// P0-BE (plan §20)
	m.Description = req.Description
	if req.IsRequired != nil {
		m.IsRequired = *req.IsRequired
	}
	if req.AvailableFrom != nil {
		avail, err := parseTimePtr(req.AvailableFrom)
		if err != nil {
			return nil, err
		}
		m.AvailableFrom = avail
	}
	if err := s.repo.CreateMaterial(ctx, m); err != nil {
		return nil, err
	}
	return materialToResponse(m), nil
}

func (s *Service) ListMaterials(ctx context.Context, sessionID string) ([]TrainingMaterialResponse, error) {
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}
	mats, err := s.repo.ListMaterials(ctx, sid)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingMaterialResponse, 0, len(mats))
	for _, m := range mats {
		responses = append(responses, *materialToResponse(&m))
	}
	return responses, nil
}

func (s *Service) UpdateMaterial(ctx context.Context, id string, req UpdateTrainingMaterialRequest) (*TrainingMaterialResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	m, err := s.repo.FindMaterialByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		m.Title = *req.Title
	}
	if req.FileURL != nil {
		m.FileURL = req.FileURL
	}
	if req.FileType != nil {
		m.FileType = req.FileType
	}
	if req.SortOrder != nil {
		m.SortOrder = *req.SortOrder
	}
	// P0-BE (plan §20)
	if req.Description != nil {
		m.Description = req.Description
	}
	if req.IsRequired != nil {
		m.IsRequired = *req.IsRequired
	}
	if req.AvailableFrom != nil {
		avail, err := parseTimePtr(req.AvailableFrom)
		if err != nil {
			return nil, err
		}
		m.AvailableFrom = avail
	}
	if err := s.repo.UpdateMaterial(ctx, m); err != nil {
		return nil, err
	}
	return materialToResponse(m), nil
}

func (s *Service) DeleteMaterial(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteMaterial(ctx, uid)
}

// =========================================================================
// Training Evaluations
// =========================================================================

func (s *Service) CreateEvaluation(ctx context.Context, req CreateTrainingEvaluationRequest) (*TrainingEvaluationResponse, error) {
	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}

	e := &TrainingEvaluation{
		SessionID:  sessionID,
		EmployeeID: empID,
		Rating:     req.Rating,
		Feedback:   req.Feedback,
	}
	if err := s.repo.CreateEvaluation(ctx, e); err != nil {
		return nil, err
	}
	s.logger.Info("Training evaluation submitted", zap.String("session_id", req.SessionID), zap.Int("rating", req.Rating))
	return evaluationToResponse(e), nil
}

func (s *Service) GetEvaluationByID(ctx context.Context, id string) (*TrainingEvaluationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	e, err := s.repo.FindEvaluationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return evaluationToResponse(e), nil
}

func (s *Service) ListEvaluations(ctx context.Context, sessionID *string, employeeID *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var sessUUID, empUUID *uuid.UUID
	if sessionID != nil && *sessionID != "" {
		uid, err := uuid.Parse(*sessionID)
		if err != nil {
			return nil, fmt.Errorf("invalid session_id: %w", err)
		}
		sessUUID = &uid
	}
	if employeeID != nil && *employeeID != "" {
		uid, err := uuid.Parse(*employeeID)
		if err != nil {
			return nil, fmt.Errorf("invalid employee_id: %w", err)
		}
		empUUID = &uid
	}
	evals, total, err := s.repo.ListEvaluations(ctx, sessUUID, empUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingEvaluationResponse, 0, len(evals))
	for _, e := range evals {
		responses = append(responses, *evaluationToResponse(&e))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateEvaluation(ctx context.Context, id string, req UpdateTrainingEvaluationRequest) (*TrainingEvaluationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	e, err := s.repo.FindEvaluationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Rating != nil {
		e.Rating = *req.Rating
	}
	if req.Feedback != nil {
		e.Feedback = req.Feedback
	}
	if err := s.repo.UpdateEvaluation(ctx, e); err != nil {
		return nil, err
	}
	return evaluationToResponse(e), nil
}

func (s *Service) DeleteEvaluation(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteEvaluation(ctx, uid)
}

// =========================================================================
// Training Certificates
// =========================================================================

func (s *Service) CreateCertificate(ctx context.Context, req CreateTrainingCertificateRequest) (*TrainingCertificateResponse, error) {
	partID, err := uuid.Parse(req.ParticipantID)
	if err != nil {
		return nil, fmt.Errorf("invalid participant_id: %w", err)
	}
	// Verify participant exists
	if _, err := s.repo.FindParticipantByID(ctx, partID); err != nil {
		return nil, fmt.Errorf("participant not found: %w", err)
	}

	cert := &TrainingCertificate{
		ParticipantID: partID,
		CertificateNo: req.CertificateNo,
		IssuedDate:    req.IssuedDate,
		ExpiryDate:    req.ExpiryDate,
	}
	if err := s.repo.CreateCertificate(ctx, cert); err != nil {
		return nil, err
	}
	s.logger.Info("Training certificate issued", zap.String("cert_no", req.CertificateNo))
	return certificateToResponse(cert), nil
}

func (s *Service) GetCertificateByID(ctx context.Context, id string) (*TrainingCertificateResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	cert, err := s.repo.FindCertificateByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return certificateToResponse(cert), nil
}

func (s *Service) ListCertificates(ctx context.Context, participantID *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var partUUID *uuid.UUID
	if participantID != nil && *participantID != "" {
		uid, err := uuid.Parse(*participantID)
		if err != nil {
			return nil, fmt.Errorf("invalid participant_id: %w", err)
		}
		partUUID = &uid
	}
	certs, total, err := s.repo.ListCertificates(ctx, partUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingCertificateResponse, 0, len(certs))
	for _, c := range certs {
		responses = append(responses, *certificateToResponse(&c))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateCertificate(ctx context.Context, id string, req UpdateTrainingCertificateRequest) (*TrainingCertificateResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	cert, err := s.repo.FindCertificateByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.CertificateNo != nil {
		cert.CertificateNo = *req.CertificateNo
	}
	if req.IssuedDate != nil {
		cert.IssuedDate = *req.IssuedDate
	}
	if req.ExpiryDate != nil {
		cert.ExpiryDate = req.ExpiryDate
	}
	if err := s.repo.UpdateCertificate(ctx, cert); err != nil {
		return nil, err
	}
	return certificateToResponse(cert), nil
}

func (s *Service) DeleteCertificate(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCertificate(ctx, uid)
}

// =========================================================================
// Helpers
// =========================================================================

func calcTotalPages(total int64, perPage int) int {
	pages := int(math.Ceil(float64(total) / float64(perPage)))
	if pages < 1 {
		return 1
	}
	return pages
}

// =========================================================================
// Response converters
// =========================================================================

func categoryToResponse(c *TrainingCategory) *TrainingCategoryResponse {
	desc := ""
	if c.Description != nil {
		desc = *c.Description
	}
	return &TrainingCategoryResponse{
		ID:          c.ID.String(),
		Code:        c.Code,
		Name:        c.Name,
		Description: desc,
		IsActive:    c.IsActive,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func courseToResponse(c *TrainingCourse) *TrainingCourseResponse {
	desc := ""
	if c.Description != nil {
		desc = *c.Description
	}
	dur := float64(0)
	if c.DurationHour != nil {
		dur = *c.DurationHour
	}
	ms := float64(0)
	if c.MinScore != nil {
		ms = *c.MinScore
	}
	cost := float64(0)
	if c.Cost != nil {
		cost = *c.Cost
	}
	vendor := ""
	if c.ExternalVendor != nil {
		vendor = *c.ExternalVendor
	}
	return &TrainingCourseResponse{
		ID:             c.ID.String(),
		CategoryID:     c.CategoryID.String(),
		Code:           c.Code,
		Name:           c.Name,
		Description:    desc,
		DurationHour:   dur,
		MinScore:       ms,
		Cost:           cost,
		IsCertified:    c.IsCertified,
		ExternalVendor: vendor,
		CourseType:     strPtr(c.CourseType),
		DeliveryType:   strPtr(c.DeliveryType),
		IsMandatory:    c.IsMandatory,
		IsActive:       c.IsActive,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}

func sessionToResponse(s *TrainingSession) *TrainingSessionResponse {
	loc := ""
	if s.Location != nil {
		loc = *s.Location
	}
	return &TrainingSessionResponse{
		ID:                   s.ID.String(),
		CourseID:             s.CourseID.String(),
		SessionCode:          s.SessionCode,
		TrainerName:          s.TrainerName,
		ProviderType:         strPtr(s.ProviderType),
		DeliveryMode:         strPtr(s.DeliveryMode),
		ProviderID:           uuidPtr(s.ProviderID),
		StartDatetime:        formatTimePtr(s.StartDatetime),
		EndDatetime:          formatTimePtr(s.EndDatetime),
		MeetingURL:           strPtr(s.MeetingURL),
		RegistrationDeadline: formatTimePtr(s.RegistrationDeadline),
		Location:             loc,
		StartDate:            s.StartDate,
		EndDate:              s.EndDate,
		MaxQuota:             s.MaxQuota,
		Status:               string(s.Status),
		CreatedAt:            s.CreatedAt,
		UpdatedAt:            s.UpdatedAt,
	}
}

func participantToResponse(p *TrainingParticipant) *TrainingParticipantResponse {
	completed := ""
	if p.CompletedAt != nil {
		completed = *p.CompletedAt
	}
	return &TrainingParticipantResponse{
		ID:                 p.ID.String(),
		SessionID:          p.SessionID.String(),
		EmployeeID:         p.EmployeeID.String(),
		RegistrationStatus: string(p.RegistrationStatus),
		RegisteredAt:       formatTimePtr(p.RegisteredAt),
		ApprovedAt:         formatTimePtr(p.ApprovedAt),
		AttendanceStatus:   string(p.AttendanceStatus),
		Score:              p.Score,
		CompletionStatus:   string(p.CompletionStatus),
		CompletionDate:     strPtr(p.CompletionDate),
		FinalScore:         floatPtr(p.FinalScore),
		Passed:             boolPtr(p.Passed),
		Remarks:            strPtr(p.Remarks),
		CompletedAt:        completed,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}
}

func materialToResponse(m *TrainingMaterial) *TrainingMaterialResponse {
	url := ""
	if m.FileURL != nil {
		url = *m.FileURL
	}
	ft := ""
	if m.FileType != nil {
		ft = *m.FileType
	}
	return &TrainingMaterialResponse{
		ID:            m.ID.String(),
		SessionID:     m.SessionID.String(),
		Title:         m.Title,
		Description:   strPtr(m.Description),
		IsRequired:    m.IsRequired,
		AvailableFrom: formatTimePtr(m.AvailableFrom),
		FileURL:       url,
		FileType:      ft,
		SortOrder:     m.SortOrder,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func evaluationToResponse(e *TrainingEvaluation) *TrainingEvaluationResponse {
	fb := ""
	if e.Feedback != nil {
		fb = *e.Feedback
	}
	return &TrainingEvaluationResponse{
		ID:         e.ID.String(),
		SessionID:  e.SessionID.String(),
		EmployeeID: e.EmployeeID.String(),
		Rating:     e.Rating,
		Feedback:   fb,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

func certificateToResponse(c *TrainingCertificate) *TrainingCertificateResponse {
	exp := ""
	if c.ExpiryDate != nil {
		exp = *c.ExpiryDate
	}
	resp := &TrainingCertificateResponse{
		ID:            c.ID.String(),
		ParticipantID: c.ParticipantID.String(),
		CertificateNo: c.CertificateNo,
		IssuedDate:    c.IssuedDate,
		ExpiryDate:    exp,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
	if c.CertificationID != nil {
		resp.CertificationID = c.CertificationID.String()
	}
	if c.CertificateFileURL != nil {
		resp.CertificateFileURL = *c.CertificateFileURL
	}
	return resp
}

func providerToResponse(p *TrainingProvider) *TrainingProviderResponse {
	return &TrainingProviderResponse{
		ID:          p.ID.String(),
		Code:        p.Code,
		Name:        p.Name,
		Type:        string(p.Type),
		ContactName: strPtr(p.ContactName),
		Email:       strPtr(p.Email),
		Phone:       strPtr(p.Phone),
		Address:     strPtr(p.Address),
		Website:     strPtr(p.Website),
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func trainerToResponse(t *TrainingTrainer) *TrainingTrainerResponse {
	return &TrainingTrainerResponse{
		ID:         t.ID.String(),
		Type:       string(t.Type),
		EmployeeID: uuidPtr(t.EmployeeID),
		ProviderID: uuidPtr(t.ProviderID),
		Name:       t.Name,
		Email:      strPtr(t.Email),
		Phone:      strPtr(t.Phone),
		Bio:        strPtr(t.Bio),
		IsActive:   t.IsActive,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

func sessionTrainerToResponse(st *TrainingSessionTrainer) *TrainingSessionTrainerResponse {
	return &TrainingSessionTrainerResponse{
		ID:        st.ID.String(),
		SessionID: st.SessionID.String(),
		TrainerID: st.TrainerID.String(),
		Role:      string(st.Role),
		CreatedAt: st.CreatedAt,
		UpdatedAt: st.UpdatedAt,
	}
}

func attendanceToResponse(a *TrainingAttendance) *TrainingAttendanceResponse {
	return &TrainingAttendanceResponse{
		ID:             a.ID.String(),
		ParticipantID:  a.ParticipantID.String(),
		AttendanceDate: a.AttendanceDate,
		CheckIn:        formatTimePtr(a.CheckIn),
		CheckOut:       formatTimePtr(a.CheckOut),
		Status:         string(a.Status),
		Remarks:        strPtr(a.Remarks),
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}
}

func sessionAttendanceRowToResponse(r attendanceScanRow) *SessionAttendanceRow {
	return &SessionAttendanceRow{
		AttendanceID:   r.AttendanceID,
		ParticipantID:  r.ParticipantID,
		EmployeeID:     r.EmployeeID,
		AttendanceDate: r.AttendanceDate,
		Status:         r.Status,
		CheckIn:        formatTimePtr(r.CheckIn),
		CheckOut:       formatTimePtr(r.CheckOut),
		Remarks:        strPtr(r.Remarks),
	}
}

func assessmentToResponse(a *TrainingAssessment) *TrainingAssessmentResponse {
	return &TrainingAssessmentResponse{
		ID:           a.ID.String(),
		SessionID:    a.SessionID.String(),
		Name:         a.Name,
		Type:         string(a.Type),
		MaxScore:     a.MaxScore,
		PassingScore: a.PassingScore,
		AttemptLimit: a.AttemptLimit,
		IsRequired:   a.IsRequired,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

func assessmentResultToResponse(r *TrainingAssessmentResult) *TrainingAssessmentResultResponse {
	return &TrainingAssessmentResultResponse{
		ID:            r.ID.String(),
		AssessmentID:  r.AssessmentID.String(),
		ParticipantID: r.ParticipantID.String(),
		Score:         r.Score,
		Passed:        r.Passed,
		Attempt:       r.Attempt,
		CompletedAt:   formatTimePtr(r.CompletedAt),
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

// =========================================================================
// Helpers: pointer converters & datetime parse/format
// =========================================================================

func strPtr[T ~string](v *T) string {
	if v == nil {
		return ""
	}
	return string(*v)
}

func uuidPtr(v *uuid.UUID) string {
	if v == nil {
		return ""
	}
	return v.String()
}

func floatPtr(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func boolPtr(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

// parseTimePtr mem-parsing string datetime ke *time.Time.
// Format diterima: RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05", atau date saja.
func parseTimePtr(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	layouts := []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02"}
	for _, l := range layouts {
		if t, err := time.ParseInLocation(l, *s, time.Local); err == nil {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("invalid datetime: %s", *s)
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02T15:04:05")
}

// applySessionEnhancement mengisi field P0-BE session + validasi (plan §14/§32).
// Hanya field non-nil yang diubah (agar aman untuk update parsial).
func applySessionEnhancement(sess *TrainingSession, providerType, deliveryMode, providerID,
	startDatetime, endDatetime, meetingURL, registrationDeadline *string) error {
	if providerType != nil {
		pt := ProviderType(*providerType)
		sess.ProviderType = &pt
	}
	if deliveryMode != nil {
		dm := DeliveryMode(*deliveryMode)
		sess.DeliveryMode = &dm
	}
	if providerID != nil {
		if *providerID == "" {
			sess.ProviderID = nil
		} else {
			pid, err := uuid.Parse(*providerID)
			if err != nil {
				return fmt.Errorf("invalid provider_id: %w", err)
			}
			sess.ProviderID = &pid
		}
	}
	if startDatetime != nil {
		t, err := parseTimePtr(startDatetime)
		if err != nil {
			return err
		}
		sess.StartDatetime = t
	}
	if endDatetime != nil {
		t, err := parseTimePtr(endDatetime)
		if err != nil {
			return err
		}
		sess.EndDatetime = t
	}
	if meetingURL != nil {
		sess.MeetingURL = meetingURL
	}
	if registrationDeadline != nil {
		t, err := parseTimePtr(registrationDeadline)
		if err != nil {
			return err
		}
		sess.RegistrationDeadline = t
	}
	// Validasi (plan §32): provider wajib untuk EXTERNAL.
	if sess.ProviderType != nil && *sess.ProviderType == ProviderTypeExternal && sess.ProviderID == nil {
		return fmt.Errorf("provider_id is required for EXTERNAL training")
	}
	// Validasi: start <= end.
	if sess.StartDatetime != nil && sess.EndDatetime != nil && sess.StartDatetime.After(*sess.EndDatetime) {
		return fmt.Errorf("start_datetime must be before end_datetime")
	}
	// Validasi: registration deadline < start.
	if sess.RegistrationDeadline != nil && sess.StartDatetime != nil &&
		sess.RegistrationDeadline.After(*sess.StartDatetime) {
		return fmt.Errorf("registration_deadline must be before start_datetime")
	}
	return nil
}

// =========================================================================
// Training Providers (P0-BE — plan §11)
// =========================================================================

func (s *Service) CreateProvider(ctx context.Context, req CreateTrainingProviderRequest) (*TrainingProviderResponse, error) {
	// Kode penyelenggara: bila tidak dikirim, di-generate otomatis PRV-{sekuens}.
	code := strings.TrimSpace(req.Code)
	if code == "" {
		var err error
		code, err = s.repo.NextProviderCode(ctx)
		if err != nil {
			return nil, err
		}
	}

	p := &TrainingProvider{
		Code:     code,
		Name:     req.Name,
		Type:     ProviderTypeExternal,
		IsActive: true,
	}
	if req.Type != nil {
		p.Type = ProviderType(*req.Type)
	}
	if req.ContactName != nil {
		p.ContactName = req.ContactName
	}
	if req.Email != nil {
		p.Email = req.Email
	}
	if req.Phone != nil {
		p.Phone = req.Phone
	}
	if req.Address != nil {
		p.Address = req.Address
	}
	if req.Website != nil {
		p.Website = req.Website
	}
	if req.IsActive != nil {
		p.IsActive = *req.IsActive
	}
	if err := s.repo.CreateProvider(ctx, p); err != nil {
		return nil, err
	}
	s.logger.Info("Training provider created", zap.String("id", p.ID.String()), zap.String("code", p.Code))
	return providerToResponse(p), nil
}

func (s *Service) GetProviderByID(ctx context.Context, id string) (*TrainingProviderResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindProviderByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return providerToResponse(p), nil
}

func (s *Service) ListProviders(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	providers, total, err := s.repo.ListProviders(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingProviderResponse, 0, len(providers))
	for _, p := range providers {
		responses = append(responses, *providerToResponse(&p))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateProvider(ctx context.Context, id string, req UpdateTrainingProviderRequest) (*TrainingProviderResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindProviderByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Code != nil {
		p.Code = *req.Code
	}
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Type != nil {
		p.Type = ProviderType(*req.Type)
	}
	if req.ContactName != nil {
		p.ContactName = req.ContactName
	}
	if req.Email != nil {
		p.Email = req.Email
	}
	if req.Phone != nil {
		p.Phone = req.Phone
	}
	if req.Address != nil {
		p.Address = req.Address
	}
	if req.Website != nil {
		p.Website = req.Website
	}
	if req.IsActive != nil {
		p.IsActive = *req.IsActive
	}
	if err := s.repo.UpdateProvider(ctx, p); err != nil {
		return nil, err
	}
	return providerToResponse(p), nil
}

func (s *Service) DeleteProvider(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteProvider(ctx, uid)
}

// =========================================================================
// Training Trainers (P0-BE — plan §12)
// =========================================================================

func (s *Service) CreateTrainer(ctx context.Context, req CreateTrainingTrainerRequest) (*TrainingTrainerResponse, error) {
	t := &TrainingTrainer{
		Type:     TrainerType(req.Type),
		Name:     req.Name,
		IsActive: true,
	}
	if req.EmployeeID != nil && *req.EmployeeID != "" {
		uid, err := uuid.Parse(*req.EmployeeID)
		if err != nil {
			return nil, fmt.Errorf("invalid employee_id: %w", err)
		}
		t.EmployeeID = &uid
	}
	if req.ProviderID != nil && *req.ProviderID != "" {
		uid, err := uuid.Parse(*req.ProviderID)
		if err != nil {
			return nil, fmt.Errorf("invalid provider_id: %w", err)
		}
		t.ProviderID = &uid
	}
	// Validasi (plan §32): trainer harus sesuai type.
	if t.Type == TrainerTypeInternal && t.EmployeeID == nil {
		return nil, fmt.Errorf("employee_id is required for INTERNAL trainer")
	}
	if t.Type == TrainerTypeExternal && t.ProviderID == nil {
		return nil, fmt.Errorf("provider_id is required for EXTERNAL trainer")
	}
	if req.Email != nil {
		t.Email = req.Email
	}
	if req.Phone != nil {
		t.Phone = req.Phone
	}
	if req.Bio != nil {
		t.Bio = req.Bio
	}
	if req.IsActive != nil {
		t.IsActive = *req.IsActive
	}
	if err := s.repo.CreateTrainer(ctx, t); err != nil {
		return nil, err
	}
	s.logger.Info("Training trainer created", zap.String("id", t.ID.String()), zap.String("name", t.Name))
	return trainerToResponse(t), nil
}

func (s *Service) GetTrainerByID(ctx context.Context, id string) (*TrainingTrainerResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	t, err := s.repo.FindTrainerByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return trainerToResponse(t), nil
}

func (s *Service) ListTrainers(ctx context.Context, trainerType *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	trainers, total, err := s.repo.ListTrainers(ctx, trainerType, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingTrainerResponse, 0, len(trainers))
	for _, t := range trainers {
		responses = append(responses, *trainerToResponse(&t))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateTrainer(ctx context.Context, id string, req UpdateTrainingTrainerRequest) (*TrainingTrainerResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	t, err := s.repo.FindTrainerByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Type != nil {
		t.Type = TrainerType(*req.Type)
	}
	if req.EmployeeID != nil {
		if *req.EmployeeID == "" {
			t.EmployeeID = nil
		} else {
			eid, err := uuid.Parse(*req.EmployeeID)
			if err != nil {
				return nil, fmt.Errorf("invalid employee_id: %w", err)
			}
			t.EmployeeID = &eid
		}
	}
	if req.ProviderID != nil {
		if *req.ProviderID == "" {
			t.ProviderID = nil
		} else {
			pid, err := uuid.Parse(*req.ProviderID)
			if err != nil {
				return nil, fmt.Errorf("invalid provider_id: %w", err)
			}
			t.ProviderID = &pid
		}
	}
	if t.Type == TrainerTypeInternal && t.EmployeeID == nil {
		return nil, fmt.Errorf("employee_id is required for INTERNAL trainer")
	}
	if t.Type == TrainerTypeExternal && t.ProviderID == nil {
		return nil, fmt.Errorf("provider_id is required for EXTERNAL trainer")
	}
	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.Email != nil {
		t.Email = req.Email
	}
	if req.Phone != nil {
		t.Phone = req.Phone
	}
	if req.Bio != nil {
		t.Bio = req.Bio
	}
	if req.IsActive != nil {
		t.IsActive = *req.IsActive
	}
	if err := s.repo.UpdateTrainer(ctx, t); err != nil {
		return nil, err
	}
	return trainerToResponse(t), nil
}

func (s *Service) DeleteTrainer(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteTrainer(ctx, uid)
}

// =========================================================================
// Training Session Trainers (P0-BE — plan §13)
// =========================================================================

func (s *Service) AddSessionTrainer(ctx context.Context, sessionID string, req AddSessionTrainerRequest) (*TrainingSessionTrainerResponse, error) {
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}
	tid, err := uuid.Parse(req.TrainerID)
	if err != nil {
		return nil, fmt.Errorf("invalid trainer_id: %w", err)
	}
	if _, err := s.repo.FindSessionByID(ctx, sid); err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}
	if _, err := s.repo.FindTrainerByID(ctx, tid); err != nil {
		return nil, fmt.Errorf("trainer not found: %w", err)
	}
	existing, err := s.repo.FindSessionTrainerBySessionAndTrainer(ctx, sid, tid)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("trainer already assigned to this session")
	}
	st := &TrainingSessionTrainer{
		SessionID: sid,
		TrainerID: tid,
		Role:      SessionTrainerRoleMain,
	}
	if req.Role != nil {
		st.Role = SessionTrainerRole(*req.Role)
	}
	if err := s.repo.CreateSessionTrainer(ctx, st); err != nil {
		return nil, err
	}
	return sessionTrainerToResponse(st), nil
}

func (s *Service) ListSessionTrainers(ctx context.Context, sessionID string) ([]TrainingSessionTrainerResponse, error) {
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}
	items, err := s.repo.ListSessionTrainers(ctx, sid)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingSessionTrainerResponse, 0, len(items))
	for _, st := range items {
		responses = append(responses, *sessionTrainerToResponse(&st))
	}
	return responses, nil
}

func (s *Service) RemoveSessionTrainer(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteSessionTrainer(ctx, uid)
}

// =========================================================================
// Training Attendances (P0-BE — plan §19)
// =========================================================================

// MarkAttendance — upsert attendance per participant per hari; optional array.
func (s *Service) MarkAttendance(ctx context.Context, sessionID string, reqs []MarkTrainingAttendanceRequest) ([]TrainingAttendanceResponse, error) {
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}
	if _, err := s.repo.FindSessionByID(ctx, sid); err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}
	var responses []TrainingAttendanceResponse
	for _, req := range reqs {
		att, err := s.markOneAttendance(ctx, sid, req)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *attendanceToResponse(att))
	}
	return responses, nil
}

func (s *Service) markOneAttendance(ctx context.Context, sessionID uuid.UUID, req MarkTrainingAttendanceRequest) (*TrainingAttendance, error) {
	pid, err := uuid.Parse(req.ParticipantID)
	if err != nil {
		return nil, fmt.Errorf("invalid participant_id: %w", err)
	}
	// Peserta harus terdaftar pada session ini.
	part, err := s.repo.FindParticipantByID(ctx, pid)
	if err != nil {
		return nil, fmt.Errorf("participant not found: %w", err)
	}
	if part.SessionID != sessionID {
		return nil, fmt.Errorf("participant does not belong to this session")
	}
	// Validasi status (plan §33): PRESENT | ABSENT | LATE | EXCUSED.
	status := AttendStatusPresent
	if req.Status != nil {
		switch AttendanceStatus(*req.Status) {
		case AttendStatusPresent, AttendStatusAbsent, AttendStatusLate, AttendStatusExcused:
			status = AttendanceStatus(*req.Status)
		default:
			return nil, fmt.Errorf("invalid attendance status: %s", *req.Status)
		}
	}
	// Validasi tanggal wajib format YYYY-MM-DD.
	if _, err := time.ParseInLocation("2006-01-02", req.AttendanceDate, time.Local); err != nil {
		return nil, fmt.Errorf("invalid attendance_date: %s", req.AttendanceDate)
	}
	checkIn, err := parseTimePtr(req.CheckIn)
	if err != nil {
		return nil, err
	}
	checkOut, err := parseTimePtr(req.CheckOut)
	if err != nil {
		return nil, err
	}

	// Upsert: satu baris per (participant, attendance_date).
	existing, err := s.repo.FindAttendanceByParticipantAndDate(ctx, pid, req.AttendanceDate)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		existing = &TrainingAttendance{
			ParticipantID:  pid,
			AttendanceDate: req.AttendanceDate,
		}
	}
	existing.CheckIn = checkIn
	existing.CheckOut = checkOut
	existing.Status = status
	existing.Remarks = req.Remarks
	if existing.ID == uuid.Nil {
		if err := s.repo.CreateAttendance(ctx, existing); err != nil {
			return nil, err
		}
	} else {
		if err := s.repo.UpdateAttendance(ctx, existing); err != nil {
			return nil, err
		}
	}
	// Sinkronisasi status agregat ke kolom legacy participant (compatibility, plan §19).
	part.AttendanceStatus = status
	if err := s.repo.UpdateParticipant(ctx, part); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) UpdateAttendance(ctx context.Context, id string, req UpdateTrainingAttendanceRequest) (*TrainingAttendanceResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindAttendanceByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.CheckIn != nil {
		ci, err := parseTimePtr(req.CheckIn)
		if err != nil {
			return nil, err
		}
		a.CheckIn = ci
	}
	if req.CheckOut != nil {
		co, err := parseTimePtr(req.CheckOut)
		if err != nil {
			return nil, err
		}
		a.CheckOut = co
	}
	if req.Status != nil {
		a.Status = AttendanceStatus(*req.Status)
	}
	if req.Remarks != nil {
		a.Remarks = req.Remarks
	}
	if err := s.repo.UpdateAttendance(ctx, a); err != nil {
		return nil, err
	}
	return attendanceToResponse(a), nil
}

func (s *Service) ListAttendanceBySession(ctx context.Context, sessionID string) ([]SessionAttendanceRow, error) {
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}
	rows, err := s.repo.ListAttendanceBySession(ctx, sid)
	if err != nil {
		return nil, err
	}
	responses := make([]SessionAttendanceRow, 0, len(rows))
	for _, r := range rows {
		responses = append(responses, *sessionAttendanceRowToResponse(r))
	}
	return responses, nil
}

// =========================================================================
// Training Assessments (P0-BE — plan §21)
// =========================================================================

func (s *Service) CreateAssessment(ctx context.Context, req CreateTrainingAssessmentRequest) (*TrainingAssessmentResponse, error) {
	sid, err := uuid.Parse(req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}
	if _, err := s.repo.FindSessionByID(ctx, sid); err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}
	a := &TrainingAssessment{
		SessionID:    sid,
		Name:         req.Name,
		Type:         AssessTypeOther,
		MaxScore:     100,
		PassingScore: 60,
		AttemptLimit: 1,
		IsRequired:   true,
	}
	if req.Type != nil {
		a.Type = AssessmentType(*req.Type)
	}
	if req.MaxScore != nil {
		a.MaxScore = *req.MaxScore
	}
	if req.PassingScore != nil {
		a.PassingScore = *req.PassingScore
	}
	if req.AttemptLimit != nil {
		a.AttemptLimit = *req.AttemptLimit
	}
	if req.IsRequired != nil {
		a.IsRequired = *req.IsRequired
	}
	// Validasi (plan §32): passing score <= max score.
	if a.PassingScore > a.MaxScore {
		return nil, fmt.Errorf("passing_score must be less than or equal to max_score")
	}
	if a.AttemptLimit < 1 {
		return nil, fmt.Errorf("attempt_limit must be at least 1")
	}
	if err := s.repo.CreateAssessment(ctx, a); err != nil {
		return nil, err
	}
	s.logger.Info("Training assessment created", zap.String("id", a.ID.String()), zap.String("session_id", req.SessionID))
	return assessmentToResponse(a), nil
}

func (s *Service) ListAssessmentsBySession(ctx context.Context, sessionID string) ([]TrainingAssessmentResponse, error) {
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}
	items, err := s.repo.ListAssessmentsBySession(ctx, sid)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingAssessmentResponse, 0, len(items))
	for _, a := range items {
		responses = append(responses, *assessmentToResponse(&a))
	}
	return responses, nil
}

func (s *Service) SubmitAssessmentResult(ctx context.Context, assessmentID string, req SubmitAssessmentResultRequest) (*TrainingAssessmentResultResponse, error) {
	aid, err := uuid.Parse(assessmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid assessment_id: %w", err)
	}
	a, err := s.repo.FindAssessmentByID(ctx, aid)
	if err != nil {
		return nil, err
	}
	pid, err := uuid.Parse(req.ParticipantID)
	if err != nil {
		return nil, fmt.Errorf("invalid participant_id: %w", err)
	}
	part, err := s.repo.FindParticipantByID(ctx, pid)
	if err != nil {
		return nil, fmt.Errorf("participant not found: %w", err)
	}
	if part.SessionID != a.SessionID {
		return nil, fmt.Errorf("participant does not belong to the assessment session")
	}
	// Validasi (plan §32): score tidak boleh melebihi max score.
	if req.Score > a.MaxScore {
		return nil, fmt.Errorf("score exceeds max_score (%.2f)", a.MaxScore)
	}
	attempts, err := s.repo.CountAssessmentAttempts(ctx, aid, pid)
	if err != nil {
		return nil, err
	}
	if attempts >= int64(a.AttemptLimit) {
		return nil, fmt.Errorf("attempt limit reached (%d)", a.AttemptLimit)
	}
	res := &TrainingAssessmentResult{
		AssessmentID:  aid,
		ParticipantID: pid,
		Score:         req.Score,
		Passed:        req.Score >= a.PassingScore,
		Attempt:       int(attempts) + 1,
	}
	if req.CompletedAt != nil {
		ct, err := parseTimePtr(req.CompletedAt)
		if err != nil {
			return nil, err
		}
		res.CompletedAt = ct
	}
	if err := s.repo.CreateAssessmentResult(ctx, res); err != nil {
		return nil, err
	}
	s.logger.Info("Training assessment result submitted", zap.String("assessment_id", assessmentID), zap.String("participant_id", req.ParticipantID))
	return assessmentResultToResponse(res), nil
}

// =========================================================================
// Training Plans (P1-BE — plan §16)
// =========================================================================

func (s *Service) CreatePlan(ctx context.Context, req CreateTrainingPlanRequest) (*TrainingPlanResponse, error) {
	// Kode plan: bila tidak dikirim, di-generate otomatis TP-{tahun}-{sekuens}.
	code := strings.TrimSpace(req.Code)
	if code == "" {
		var err error
		code, err = s.repo.NextPlanCode(ctx, req.Year)
		if err != nil {
			return nil, err
		}
	}

	p := &TrainingPlan{
		Code:   code,
		Name:   req.Name,
		Year:   req.Year,
		Status: PlanStatusDraft,
	}
	if req.Description != nil {
		p.Description = req.Description
	}
	if req.Status != nil {
		p.Status = PlanStatus(*req.Status)
	}
	if err := s.repo.CreatePlan(ctx, p); err != nil {
		return nil, err
	}
	s.logger.Info("Training plan created", zap.String("id", p.ID.String()), zap.String("code", p.Code))
	return planToResponse(p), nil
}

func (s *Service) GetPlanByID(ctx context.Context, id string) (*TrainingPlanResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindPlanByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return planToResponse(p), nil
}

func (s *Service) ListPlans(ctx context.Context, year *int, status *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	plans, total, err := s.repo.ListPlans(ctx, year, status, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingPlanResponse, 0, len(plans))
	for _, p := range plans {
		responses = append(responses, *planToResponse(&p))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdatePlan(ctx context.Context, id string, req UpdateTrainingPlanRequest) (*TrainingPlanResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindPlanByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Code != nil {
		p.Code = *req.Code
	}
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Year != nil {
		p.Year = *req.Year
	}
	if req.Description != nil {
		p.Description = req.Description
	}
	if req.Status != nil {
		p.Status = PlanStatus(*req.Status)
	}
	if err := s.repo.UpdatePlan(ctx, p); err != nil {
		return nil, err
	}
	return planToResponse(p), nil
}

func (s *Service) DeletePlan(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeletePlan(ctx, uid)
}

// =========================================================================
// Training Plan Items (P1-BE — plan §16)
// =========================================================================

func (s *Service) CreatePlanItem(ctx context.Context, planID string, req CreateTrainingPlanItemRequest) (*TrainingPlanItemResponse, error) {
	pid, err := uuid.Parse(planID)
	if err != nil {
		return nil, fmt.Errorf("invalid training_plan_id: %w", err)
	}
	cid, err := uuid.Parse(req.CourseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course_id: %w", err)
	}
	if _, err := s.repo.FindPlanByID(ctx, pid); err != nil {
		return nil, fmt.Errorf("training plan not found: %w", err)
	}
	if _, err := s.repo.FindCourseByID(ctx, cid); err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}
	item := &TrainingPlanItem{
		TrainingPlanID: pid,
		CourseID:       cid,
		Priority:       PriorityMedium,
	}
	if req.TargetDate != nil {
		item.TargetDate = req.TargetDate
	}
	if req.TargetParticipants != nil {
		item.TargetParticipants = req.TargetParticipants
	}
	if req.EstimatedCost != nil {
		item.EstimatedCost = req.EstimatedCost
	}
	if req.Priority != nil {
		item.Priority = PriorityLevel(*req.Priority)
	}
	if err := s.repo.CreatePlanItem(ctx, item); err != nil {
		return nil, err
	}
	return planItemToResponse(item), nil
}

func (s *Service) ListPlanItems(ctx context.Context, planID string) ([]TrainingPlanItemResponse, error) {
	pid, err := uuid.Parse(planID)
	if err != nil {
		return nil, fmt.Errorf("invalid training_plan_id: %w", err)
	}
	items, err := s.repo.ListPlanItems(ctx, pid)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingPlanItemResponse, 0, len(items))
	for _, i := range items {
		responses = append(responses, *planItemToResponse(&i))
	}
	return responses, nil
}

func (s *Service) UpdatePlanItem(ctx context.Context, id string, req UpdateTrainingPlanItemRequest) (*TrainingPlanItemResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	item, err := s.repo.FindPlanItemByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.CourseID != nil {
		cid, err := uuid.Parse(*req.CourseID)
		if err != nil {
			return nil, fmt.Errorf("invalid course_id: %w", err)
		}
		item.CourseID = cid
	}
	if req.TargetDate != nil {
		item.TargetDate = req.TargetDate
	}
	if req.TargetParticipants != nil {
		item.TargetParticipants = req.TargetParticipants
	}
	if req.EstimatedCost != nil {
		item.EstimatedCost = req.EstimatedCost
	}
	if req.Priority != nil {
		item.Priority = PriorityLevel(*req.Priority)
	}
	if err := s.repo.UpdatePlanItem(ctx, item); err != nil {
		return nil, err
	}
	return planItemToResponse(item), nil
}

func (s *Service) DeletePlanItem(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeletePlanItem(ctx, uid)
}

// =========================================================================
// Training Needs (P1-BE — plan §17)
// =========================================================================

func (s *Service) CreateNeed(ctx context.Context, req CreateTrainingNeedRequest) (*TrainingNeedResponse, error) {
	n := &TrainingNeed{
		Priority:   PriorityMedium,
		SourceType: NeedSourceManual,
		Status:     NeedStatusOpen,
	}
	if err := applyNeedFields(n, req.EmployeeID, req.OrganizationID, req.PositionID, req.CourseID,
		req.Reason, req.Priority, req.SourceType, req.SourceID, req.Status); err != nil {
		return nil, err
	}
	if err := s.repo.CreateNeed(ctx, n); err != nil {
		return nil, err
	}
	return needToResponse(n), nil
}

// CreateOnboardingNeed (S-7 — onboarding → training handoff) membuat kebutuhan
// training untuk employee yang baru menyelesaikan onboarding. Dipanggil oleh
// Recruitment via interface narrow (adapter di cmd/server/main.go) — Training
// tetap source of truth training; Recruitment hanya menghasilkan kebutuhan
// (handoff), tidak mengeksekusi training. source_id = employee_onboarding.id
// sehingga asal kebutuhan terlacak ke onboarding spesifik.
func (s *Service) CreateOnboardingNeed(ctx context.Context, employeeID, onboardingID uuid.UUID, reason string) (*TrainingNeedResponse, error) {
	r := reason
	if r == "" {
		r = "Onboarding completed — training plan handoff"
	}
	srcID := onboardingID
	n := &TrainingNeed{
		EmployeeID: &employeeID,
		Priority:   PriorityMedium,
		SourceType: NeedSourceOnboarding,
		SourceID:   &srcID,
		Status:     NeedStatusOpen,
		Reason:     &r,
	}
	if err := s.repo.CreateNeed(ctx, n); err != nil {
		return nil, err
	}
	return needToResponse(n), nil
}

func (s *Service) GetNeedByID(ctx context.Context, id string) (*TrainingNeedResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	n, err := s.repo.FindNeedByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return needToResponse(n), nil
}

func (s *Service) ListNeeds(ctx context.Context, employeeID, courseID *string, status *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var empUUID, cUUID *uuid.UUID
	if employeeID != nil && *employeeID != "" {
		uid, err := uuid.Parse(*employeeID)
		if err != nil {
			return nil, fmt.Errorf("invalid employee_id: %w", err)
		}
		empUUID = &uid
	}
	if courseID != nil && *courseID != "" {
		uid, err := uuid.Parse(*courseID)
		if err != nil {
			return nil, fmt.Errorf("invalid course_id: %w", err)
		}
		cUUID = &uid
	}
	needs, total, err := s.repo.ListNeeds(ctx, empUUID, cUUID, status, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingNeedResponse, 0, len(needs))
	for _, n := range needs {
		responses = append(responses, *needToResponse(&n))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

// DeleteNeed — soft delete training need (plan §17).
func (s *Service) DeleteNeed(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if _, err := s.repo.FindNeedByID(ctx, uid); err != nil {
		return fmt.Errorf("need not found: %w", err)
	}
	return s.repo.DeleteNeed(ctx, uid)
}

func (s *Service) UpdateNeed(ctx context.Context, id string, req UpdateTrainingNeedRequest) (*TrainingNeedResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	n, err := s.repo.FindNeedByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if err := applyNeedFields(n, req.EmployeeID, req.OrganizationID, req.PositionID, req.CourseID,
		req.Reason, req.Priority, req.SourceType, req.SourceID, req.Status); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateNeed(ctx, n); err != nil {
		return nil, err
	}
	return needToResponse(n), nil
}

// applyNeedFields — hanya field non-nil yang diubah (aman untuk update parsial).
func applyNeedFields(n *TrainingNeed, employeeID, organizationID, positionID, courseID, reason,
	priority, sourceType, sourceID, status *string) error {
	if employeeID != nil {
		if *employeeID == "" {
			n.EmployeeID = nil
		} else {
			uid, err := uuid.Parse(*employeeID)
			if err != nil {
				return fmt.Errorf("invalid employee_id: %w", err)
			}
			n.EmployeeID = &uid
		}
	}
	if organizationID != nil {
		if *organizationID == "" {
			n.OrganizationID = nil
		} else {
			uid, err := uuid.Parse(*organizationID)
			if err != nil {
				return fmt.Errorf("invalid organization_id: %w", err)
			}
			n.OrganizationID = &uid
		}
	}
	if positionID != nil {
		if *positionID == "" {
			n.PositionID = nil
		} else {
			uid, err := uuid.Parse(*positionID)
			if err != nil {
				return fmt.Errorf("invalid position_id: %w", err)
			}
			n.PositionID = &uid
		}
	}
	if courseID != nil {
		if *courseID == "" {
			n.CourseID = nil
		} else {
			uid, err := uuid.Parse(*courseID)
			if err != nil {
				return fmt.Errorf("invalid course_id: %w", err)
			}
			n.CourseID = &uid
		}
	}
	if reason != nil {
		n.Reason = reason
	}
	if priority != nil {
		n.Priority = PriorityLevel(*priority)
	}
	if sourceType != nil {
		n.SourceType = NeedSourceType(*sourceType)
	}
	if sourceID != nil {
		if *sourceID == "" {
			n.SourceID = nil
		} else {
			uid, err := uuid.Parse(*sourceID)
			if err != nil {
				return fmt.Errorf("invalid source_id: %w", err)
			}
			n.SourceID = &uid
		}
	}
	if status != nil {
		n.Status = NeedStatus(*status)
	}
	return nil
}

// =========================================================================
// Training Requests (P1-BE — plan §15, Central Approval)
// =========================================================================

func (s *Service) CreateRequest(ctx context.Context, req CreateTrainingRequestRequest) (*TrainingRequestResponse, error) {
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	courseID, err := uuid.Parse(req.CourseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course_id: %w", err)
	}
	if _, err := s.repo.FindCourseByID(ctx, courseID); err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}
	tr := &TrainingRequest{
		EmployeeID:    empID,
		CourseID:      courseID,
		RequestedDate: req.RequestedDate,
		Priority:      PriorityMedium,
		Status:        ReqStatusDraft,
	}
	if req.SessionID != nil {
		sid, err := uuid.Parse(*req.SessionID)
		if err != nil {
			return nil, fmt.Errorf("invalid session_id: %w", err)
		}
		tr.SessionID = &sid
	}
	if req.Reason != nil {
		tr.Reason = req.Reason
	}
	if req.Priority != nil {
		tr.Priority = PriorityLevel(*req.Priority)
	}
	if err := s.repo.CreateRequest(ctx, tr); err != nil {
		return nil, err
	}
	return requestToResponse(tr), nil
}

func (s *Service) GetRequestByID(ctx context.Context, id string) (*TrainingRequestResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	tr, err := s.repo.FindRequestByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return requestToResponse(tr), nil
}

func (s *Service) ListRequests(ctx context.Context, employeeID *string, status *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var empUUID *uuid.UUID
	if employeeID != nil && *employeeID != "" {
		uid, err := uuid.Parse(*employeeID)
		if err != nil {
			return nil, fmt.Errorf("invalid employee_id: %w", err)
		}
		empUUID = &uid
	}
	reqs, total, err := s.repo.ListRequests(ctx, empUUID, status, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingRequestResponse, 0, len(reqs))
	for _, r := range reqs {
		responses = append(responses, *requestToResponse(&r))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

// SubmitRequest — buat approval instance via Central Approval (plan §15).
// Status request berubah menjadi PENDING_APPROVAL setelah instance dibuat.
func (s *Service) SubmitRequest(ctx context.Context, id string, flowID *string) (*TrainingRequestResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	tr, err := s.repo.FindRequestByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if tr.Status != ReqStatusDraft && tr.Status != ReqStatusRejected {
		return nil, fmt.Errorf("only DRAFT or REJECTED requests can be submitted")
	}
	if s.approvalEngine == nil {
		return nil, fmt.Errorf("approval engine is not configured")
	}
	// Auto-resolve flow bila client tidak mengirim flow_id (pola leave).
	resolvedFlow := ""
	if flowID != nil && *flowID != "" {
		resolvedFlow = *flowID
	} else if f, err := s.approvalEngine.GetActiveFlowIDForModule(ctx, "training_request"); err == nil {
		resolvedFlow = f
	}
	if resolvedFlow == "" {
		tr.Status = ReqStatusSubmitted
		if err := s.repo.UpdateRequest(ctx, tr); err != nil {
			return nil, err
		}
		return requestToResponse(tr), nil
	}
	instanceID, err := s.approvalEngine.CreateApprovalInstance(ctx, "training_request", tr.ID.String(), resolvedFlow)
	if err != nil {
		return nil, err
	}
	parsed, err := uuid.Parse(instanceID)
	if err != nil {
		return nil, fmt.Errorf("invalid approval instance id: %w", err)
	}
	tr.ApprovalInstanceID = &parsed
	tr.Status = ReqStatusPendingApproval
	if err := s.repo.UpdateRequest(ctx, tr); err != nil {
		return nil, err
	}
	s.logger.Info("Training request submitted for approval", zap.String("id", tr.ID.String()), zap.String("instance_id", instanceID))
	return requestToResponse(tr), nil
}

// CancelRequest — batalkan request (tanpa approval instance baru; pola leave cancel).
func (s *Service) CancelRequest(ctx context.Context, id string) (*TrainingRequestResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	tr, err := s.repo.FindRequestByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if tr.Status == ReqStatusApproved || tr.Status == ReqStatusCancelled {
		return nil, fmt.Errorf("request cannot be cancelled in status %s", tr.Status)
	}
	tr.Status = ReqStatusCancelled
	if err := s.repo.UpdateRequest(ctx, tr); err != nil {
		return nil, err
	}
	return requestToResponse(tr), nil
}

// HandleApprovalStatusChange — callback push dari Central Approval (plan §15/§45).
// APPROVED → status request APPROVED + auto-enroll participant ke session (bila ada).
// REJECTED → status REJECTED; CANCELLED → status CANCELLED.
func (s *Service) HandleApprovalStatusChange(ctx context.Context, documentID uuid.UUID, status string, note string) error {
	tr, err := s.repo.FindRequestByID(ctx, documentID)
	if err != nil {
		return err
	}
	if tr.Status != ReqStatusPendingApproval {
		return nil
	}
	if status == "PENDING" {
		if note != "" {
			tr.SupervisorNote = &note
			return s.repo.UpdateRequest(ctx, tr)
		}
		return nil
	}
	now := time.Now()
	switch status {
	case "APPROVED":
		tr.Status = ReqStatusApproved
		tr.ApprovedAt = &now
		if note != "" {
			tr.SupervisorNote = &note
		}
	case "REJECTED":
		tr.Status = ReqStatusRejected
		tr.RejectedAt = &now
		if note != "" {
			tr.SupervisorNote = &note
		}
	case "CANCELLED":
		tr.Status = ReqStatusCancelled
	default:
		return nil
	}
	if err := s.repo.UpdateRequest(ctx, tr); err != nil {
		return err
	}
	s.logger.Info("Training request status updated via approval status handler",
		zap.String("training_request_id", tr.ID.String()),
		zap.String("approval_status", status))

	// Auto-enroll participant ke session saat APPROVED (plan §15/§30).
	if status == "APPROVED" && tr.SessionID != nil && tr.ApprovalInstanceID != nil {
		if err := s.autoEnrollApprovedRequest(ctx, tr); err != nil {
			s.logger.Warn("auto-enroll failed after approval", zap.String("training_request_id", tr.ID.String()), zap.Error(err))
		}
	}
	return nil
}

// autoEnrollApprovedRequest — register employee ke session saat request disetujui.
func (s *Service) autoEnrollApprovedRequest(ctx context.Context, tr *TrainingRequest) error {
	existing, err := s.repo.FindParticipantBySessionAndEmployee(ctx, *tr.SessionID, tr.EmployeeID)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil // sudah terdaftar — idempotent
	}
	sess, err := s.repo.FindSessionByID(ctx, *tr.SessionID)
	if err != nil {
		return err
	}
	count, err := s.repo.CountParticipantsBySession(ctx, *tr.SessionID)
	if err != nil {
		return err
	}
	if count >= int64(sess.MaxQuota) {
		return fmt.Errorf("session quota full (%d/%d)", count, sess.MaxQuota)
	}
	now := time.Now()
	p := &TrainingParticipant{
		SessionID:          *tr.SessionID,
		EmployeeID:         tr.EmployeeID,
		RegistrationStatus: RegStatusApproved,
		AttendanceStatus:   AttendStatusPresent,
		CompletionStatus:   CompletionNotStarted,
		RegisteredAt:       &now,
		ApprovedAt:         &now,
	}
	return s.repo.CreateParticipant(ctx, p)
}

// =========================================================================
// Course Objectives (P1-BE — plan §8)
// =========================================================================

func (s *Service) CreateCourseObjective(ctx context.Context, courseID string, req CreateCourseObjectiveRequest) (*CourseObjectiveResponse, error) {
	cid, err := uuid.Parse(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course_id: %w", err)
	}
	if _, err := s.repo.FindCourseByID(ctx, cid); err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}
	o := &TrainingCourseObjective{CourseID: cid, Objective: req.Objective}
	if req.SortOrder != nil {
		o.SortOrder = *req.SortOrder
	}
	if err := s.repo.CreateCourseObjective(ctx, o); err != nil {
		return nil, err
	}
	return courseObjectiveToResponse(o), nil
}

func (s *Service) ListCourseObjectives(ctx context.Context, courseID string) ([]CourseObjectiveResponse, error) {
	cid, err := uuid.Parse(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course_id: %w", err)
	}
	items, err := s.repo.ListCourseObjectives(ctx, cid)
	if err != nil {
		return nil, err
	}
	responses := make([]CourseObjectiveResponse, 0, len(items))
	for _, o := range items {
		responses = append(responses, *courseObjectiveToResponse(&o))
	}
	return responses, nil
}

func (s *Service) UpdateCourseObjective(ctx context.Context, id string, req UpdateCourseObjectiveRequest) (*CourseObjectiveResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	o, err := s.repo.FindCourseObjectiveByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Objective != nil {
		o.Objective = *req.Objective
	}
	if req.SortOrder != nil {
		o.SortOrder = *req.SortOrder
	}
	if err := s.repo.UpdateCourseObjective(ctx, o); err != nil {
		return nil, err
	}
	return courseObjectiveToResponse(o), nil
}

func (s *Service) DeleteCourseObjective(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCourseObjective(ctx, uid)
}

// =========================================================================
// Course Competencies (P1-BE — plan §9)
// =========================================================================

func (s *Service) CreateCourseCompetency(ctx context.Context, courseID string, req CreateCourseCompetencyRequest) (*CourseCompetencyResponse, error) {
	cid, err := uuid.Parse(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course_id: %w", err)
	}
	compID, err := uuid.Parse(req.CompetencyID)
	if err != nil {
		return nil, fmt.Errorf("invalid competency_id: %w", err)
	}
	if _, err := s.repo.FindCourseByID(ctx, cid); err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}
	c := &TrainingCourseCompetency{CourseID: cid, CompetencyID: compID}
	if req.TargetLevel != nil {
		c.TargetLevel = req.TargetLevel
	}
	if err := s.repo.CreateCourseCompetency(ctx, c); err != nil {
		return nil, err
	}
	return courseCompetencyToResponse(c), nil
}

func (s *Service) ListCourseCompetencies(ctx context.Context, courseID string) ([]CourseCompetencyResponse, error) {
	cid, err := uuid.Parse(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course_id: %w", err)
	}
	items, err := s.repo.ListCourseCompetencies(ctx, cid)
	if err != nil {
		return nil, err
	}
	responses := make([]CourseCompetencyResponse, 0, len(items))
	for _, c := range items {
		responses = append(responses, *courseCompetencyToResponse(&c))
	}
	return responses, nil
}

func (s *Service) DeleteCourseCompetency(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCourseCompetency(ctx, uid)
}

// =========================================================================
// Course Prerequisites (P1-BE — plan §10)
// =========================================================================

func (s *Service) CreateCoursePrerequisite(ctx context.Context, courseID string, req CreateCoursePrerequisiteRequest) (*CoursePrerequisiteResponse, error) {
	cid, err := uuid.Parse(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course_id: %w", err)
	}
	if _, err := s.repo.FindCourseByID(ctx, cid); err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}
	p := &TrainingCoursePrerequisite{
		CourseID:         cid,
		PrerequisiteType: PrerequisiteType(req.PrerequisiteType),
		IsRequired:       true,
	}
	if req.PrerequisiteID != nil && *req.PrerequisiteID != "" {
		pid, err := uuid.Parse(*req.PrerequisiteID)
		if err != nil {
			return nil, fmt.Errorf("invalid prerequisite_id: %w", err)
		}
		p.PrerequisiteID = &pid
	}
	if req.IsRequired != nil {
		p.IsRequired = *req.IsRequired
	}
	if err := s.repo.CreateCoursePrerequisite(ctx, p); err != nil {
		return nil, err
	}
	return coursePrerequisiteToResponse(p), nil
}

func (s *Service) ListCoursePrerequisites(ctx context.Context, courseID string) ([]CoursePrerequisiteResponse, error) {
	cid, err := uuid.Parse(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course_id: %w", err)
	}
	items, err := s.repo.ListCoursePrerequisites(ctx, cid)
	if err != nil {
		return nil, err
	}
	responses := make([]CoursePrerequisiteResponse, 0, len(items))
	for _, p := range items {
		responses = append(responses, *coursePrerequisiteToResponse(&p))
	}
	return responses, nil
}

func (s *Service) DeleteCoursePrerequisite(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCoursePrerequisite(ctx, uid)
}

// =========================================================================
// Training Mandatories (P1-BE — plan §25)
// =========================================================================

func (s *Service) CreateMandatory(ctx context.Context, req CreateTrainingMandatoryRequest) (*TrainingMandatoryResponse, error) {
	cid, err := uuid.Parse(req.CourseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course_id: %w", err)
	}
	if _, err := s.repo.FindCourseByID(ctx, cid); err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}
	m := &TrainingMandatory{CourseID: cid, IsActive: true}
	if err := applyMandatoryFields(m, req.OrganizationID, req.PositionID, req.EmploymentStatusID,
		req.DueDays, req.ValidityPeriodMonth, req.IsActive); err != nil {
		return nil, err
	}
	if err := s.repo.CreateMandatory(ctx, m); err != nil {
		return nil, err
	}
	return mandatoryToResponse(m), nil
}

func (s *Service) GetMandatoryByID(ctx context.Context, id string) (*TrainingMandatoryResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	m, err := s.repo.FindMandatoryByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return mandatoryToResponse(m), nil
}

func (s *Service) ListMandatories(ctx context.Context, courseID *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var cUUID *uuid.UUID
	if courseID != nil && *courseID != "" {
		uid, err := uuid.Parse(*courseID)
		if err != nil {
			return nil, fmt.Errorf("invalid course_id: %w", err)
		}
		cUUID = &uid
	}
	items, total, err := s.repo.ListMandatories(ctx, cUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingMandatoryResponse, 0, len(items))
	for _, m := range items {
		responses = append(responses, *mandatoryToResponse(&m))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateMandatory(ctx context.Context, id string, req UpdateTrainingMandatoryRequest) (*TrainingMandatoryResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	m, err := s.repo.FindMandatoryByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.CourseID != nil {
		cid, err := uuid.Parse(*req.CourseID)
		if err != nil {
			return nil, fmt.Errorf("invalid course_id: %w", err)
		}
		m.CourseID = cid
	}
	if err := applyMandatoryFields(m, req.OrganizationID, req.PositionID, req.EmploymentStatusID,
		req.DueDays, req.ValidityPeriodMonth, req.IsActive); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateMandatory(ctx, m); err != nil {
		return nil, err
	}
	return mandatoryToResponse(m), nil
}

func (s *Service) DeleteMandatory(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteMandatory(ctx, uid)
}

// applyMandatoryFields — hanya field non-nil yang diubah.
func applyMandatoryFields(m *TrainingMandatory, organizationID, positionID, employmentStatusID *string,
	dueDays, validityPeriodMonth *int, isActive *bool) error {
	if organizationID != nil {
		if *organizationID == "" {
			m.OrganizationID = nil
		} else {
			uid, err := uuid.Parse(*organizationID)
			if err != nil {
				return fmt.Errorf("invalid organization_id: %w", err)
			}
			m.OrganizationID = &uid
		}
	}
	if positionID != nil {
		if *positionID == "" {
			m.PositionID = nil
		} else {
			uid, err := uuid.Parse(*positionID)
			if err != nil {
				return fmt.Errorf("invalid position_id: %w", err)
			}
			m.PositionID = &uid
		}
	}
	if employmentStatusID != nil {
		if *employmentStatusID == "" {
			m.EmploymentStatusID = nil
		} else {
			uid, err := uuid.Parse(*employmentStatusID)
			if err != nil {
				return fmt.Errorf("invalid employment_status_id: %w", err)
			}
			m.EmploymentStatusID = &uid
		}
	}
	if dueDays != nil {
		m.DueDays = dueDays
	}
	if validityPeriodMonth != nil {
		m.ValidityPeriodMonth = validityPeriodMonth
	}
	if isActive != nil {
		m.IsActive = *isActive
	}
	return nil
}

// =========================================================================
// Training Session Costs (P1-BE — plan §26)
// =========================================================================

func (s *Service) CreateSessionCost(ctx context.Context, sessionID string, req CreateTrainingSessionCostRequest) (*TrainingSessionCostResponse, error) {
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}
	if _, err := s.repo.FindSessionByID(ctx, sid); err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}
	c := &TrainingSessionCost{SessionID: sid, CostType: CostType(req.CostType)}
	if req.Description != nil {
		c.Description = req.Description
	}
	if req.Amount != nil {
		c.Amount = *req.Amount
	}
	if err := s.repo.CreateSessionCost(ctx, c); err != nil {
		return nil, err
	}
	return sessionCostToResponse(c), nil
}

func (s *Service) ListSessionCosts(ctx context.Context, sessionID string) ([]TrainingSessionCostResponse, error) {
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}
	items, err := s.repo.ListSessionCosts(ctx, sid)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingSessionCostResponse, 0, len(items))
	for _, c := range items {
		responses = append(responses, *sessionCostToResponse(&c))
	}
	return responses, nil
}

func (s *Service) UpdateSessionCost(ctx context.Context, id string, req UpdateTrainingSessionCostRequest) (*TrainingSessionCostResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	c, err := s.repo.FindSessionCostByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.CostType != nil {
		c.CostType = CostType(*req.CostType)
	}
	if req.Description != nil {
		c.Description = req.Description
	}
	if req.Amount != nil {
		c.Amount = *req.Amount
	}
	if err := s.repo.UpdateSessionCost(ctx, c); err != nil {
		return nil, err
	}
	return sessionCostToResponse(c), nil
}

func (s *Service) DeleteSessionCost(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteSessionCost(ctx, uid)
}

// =========================================================================
// Training Documents (P1-BE — plan §27)
// =========================================================================

func (s *Service) CreateDocument(ctx context.Context, sessionID string, req CreateTrainingDocumentRequest) (*TrainingDocumentResponse, error) {
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}
	if _, err := s.repo.FindSessionByID(ctx, sid); err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}
	d := &TrainingDocument{
		SessionID:    sid,
		DocumentType: DocumentType(req.DocumentType),
		FileName:     &req.FileName,
		FileURL:      &req.FileURL,
	}
	if req.FileName == "" {
		d.FileName = nil
	}
	if err := s.repo.CreateDocument(ctx, d); err != nil {
		return nil, err
	}
	return documentToResponse(d), nil
}

func (s *Service) ListDocuments(ctx context.Context, sessionID string) ([]TrainingDocumentResponse, error) {
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}
	items, err := s.repo.ListDocuments(ctx, sid)
	if err != nil {
		return nil, err
	}
	responses := make([]TrainingDocumentResponse, 0, len(items))
	for _, d := range items {
		responses = append(responses, *documentToResponse(&d))
	}
	return responses, nil
}

func (s *Service) DeleteDocument(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteDocument(ctx, uid)
}

// =========================================================================
// P1-BE response converters
// =========================================================================

func planToResponse(p *TrainingPlan) *TrainingPlanResponse {
	desc := ""
	if p.Description != nil {
		desc = *p.Description
	}
	return &TrainingPlanResponse{
		ID:          p.ID.String(),
		Code:        p.Code,
		Name:        p.Name,
		Year:        p.Year,
		Description: desc,
		Status:      string(p.Status),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func planItemToResponse(i *TrainingPlanItem) *TrainingPlanItemResponse {
	targetDate := ""
	if i.TargetDate != nil {
		targetDate = *i.TargetDate
	}
	cost := float64(0)
	if i.EstimatedCost != nil {
		cost = *i.EstimatedCost
	}
	return &TrainingPlanItemResponse{
		ID:                 i.ID.String(),
		TrainingPlanID:     i.TrainingPlanID.String(),
		CourseID:           i.CourseID.String(),
		TargetDate:         targetDate,
		TargetParticipants: i.TargetParticipants,
		EstimatedCost:      cost,
		Priority:           string(i.Priority),
		CreatedAt:          i.CreatedAt,
		UpdatedAt:          i.UpdatedAt,
	}
}

func needToResponse(n *TrainingNeed) *TrainingNeedResponse {
	return &TrainingNeedResponse{
		ID:             n.ID.String(),
		EmployeeID:     uuidPtr(n.EmployeeID),
		OrganizationID: uuidPtr(n.OrganizationID),
		PositionID:     uuidPtr(n.PositionID),
		CourseID:       uuidPtr(n.CourseID),
		Reason:         strPtr(n.Reason),
		Priority:       string(n.Priority),
		SourceType:     string(n.SourceType),
		SourceID:       uuidPtr(n.SourceID),
		Status:         string(n.Status),
		CreatedAt:      n.CreatedAt,
		UpdatedAt:      n.UpdatedAt,
	}
}

func requestToResponse(r *TrainingRequest) *TrainingRequestResponse {
	return &TrainingRequestResponse{
		ID:                 r.ID.String(),
		EmployeeID:         r.EmployeeID.String(),
		CourseID:           r.CourseID.String(),
		SessionID:          uuidPtr(r.SessionID),
		RequestedDate:      r.RequestedDate,
		Reason:             strPtr(r.Reason),
		Priority:           string(r.Priority),
		Status:             string(r.Status),
		ApprovalInstanceID: uuidPtr(r.ApprovalInstanceID),
		ApprovedAt:         formatTimePtr(r.ApprovedAt),
		RejectedAt:         formatTimePtr(r.RejectedAt),
		SupervisorNote:     strPtr(r.SupervisorNote),
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
	}
}

func courseObjectiveToResponse(o *TrainingCourseObjective) *CourseObjectiveResponse {
	return &CourseObjectiveResponse{
		ID:        o.ID.String(),
		CourseID:  o.CourseID.String(),
		Objective: o.Objective,
		SortOrder: o.SortOrder,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func courseCompetencyToResponse(c *TrainingCourseCompetency) *CourseCompetencyResponse {
	return &CourseCompetencyResponse{
		ID:           c.ID.String(),
		CourseID:     c.CourseID.String(),
		CompetencyID: c.CompetencyID.String(),
		TargetLevel:  c.TargetLevel,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func coursePrerequisiteToResponse(p *TrainingCoursePrerequisite) *CoursePrerequisiteResponse {
	return &CoursePrerequisiteResponse{
		ID:               p.ID.String(),
		CourseID:         p.CourseID.String(),
		PrerequisiteType: string(p.PrerequisiteType),
		PrerequisiteID:   uuidPtr(p.PrerequisiteID),
		IsRequired:       p.IsRequired,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
	}
}

func mandatoryToResponse(m *TrainingMandatory) *TrainingMandatoryResponse {
	return &TrainingMandatoryResponse{
		ID:                  m.ID.String(),
		CourseID:            m.CourseID.String(),
		OrganizationID:      uuidPtr(m.OrganizationID),
		PositionID:          uuidPtr(m.PositionID),
		EmploymentStatusID:  uuidPtr(m.EmploymentStatusID),
		DueDays:             m.DueDays,
		ValidityPeriodMonth: m.ValidityPeriodMonth,
		IsActive:            m.IsActive,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}
}

func sessionCostToResponse(c *TrainingSessionCost) *TrainingSessionCostResponse {
	return &TrainingSessionCostResponse{
		ID:          c.ID.String(),
		SessionID:   c.SessionID.String(),
		CostType:    string(c.CostType),
		Description: strPtr(c.Description),
		Amount:      c.Amount,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func documentToResponse(d *TrainingDocument) *TrainingDocumentResponse {
	return &TrainingDocumentResponse{
		ID:           d.ID.String(),
		SessionID:    d.SessionID.String(),
		DocumentType: string(d.DocumentType),
		FileName:     strPtr(d.FileName),
		FileURL:      strPtr(d.FileURL),
		UploadedBy:   uuidPtr(d.UploadedBy),
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

// =========================================================================
// Evaluation Forms (P2-BE — plan §22)
// =========================================================================

func (s *Service) CreateEvaluationForm(ctx context.Context, req CreateEvaluationFormRequest) (*EvaluationFormResponse, error) {
	sid, err := uuid.Parse(req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}
	if _, err := s.repo.FindSessionByID(ctx, sid); err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}
	f := &TrainingEvaluationForm{SessionID: sid, Name: req.Name, IsActive: true}
	if req.IsActive != nil {
		f.IsActive = *req.IsActive
	}
	if err := s.repo.CreateEvaluationForm(ctx, f); err != nil {
		return nil, err
	}
	return evaluationFormToResponse(f), nil
}

func (s *Service) GetEvaluationFormByID(ctx context.Context, id string) (*EvaluationFormResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	f, err := s.repo.FindEvaluationFormByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return evaluationFormToResponse(f), nil
}

func (s *Service) ListEvaluationForms(ctx context.Context, sessionID *string, page, perPage int) (*PaginatedResponse, error) {
	var sessUUID *uuid.UUID
	if sessionID != nil && *sessionID != "" {
		uid, err := uuid.Parse(*sessionID)
		if err != nil {
			return nil, fmt.Errorf("invalid session_id: %w", err)
		}
		sessUUID = &uid
	}
	forms, total, err := s.repo.ListEvaluationForms(ctx, sessUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]EvaluationFormResponse, 0, len(forms))
	for _, f := range forms {
		responses = append(responses, *evaluationFormToResponse(&f))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateEvaluationForm(ctx context.Context, id string, req UpdateEvaluationFormRequest) (*EvaluationFormResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	f, err := s.repo.FindEvaluationFormByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		f.Name = *req.Name
	}
	if req.IsActive != nil {
		f.IsActive = *req.IsActive
	}
	if err := s.repo.UpdateEvaluationForm(ctx, f); err != nil {
		return nil, err
	}
	return evaluationFormToResponse(f), nil
}

func (s *Service) DeleteEvaluationForm(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if _, err := s.repo.FindEvaluationFormByID(ctx, uid); err != nil {
		return fmt.Errorf("evaluation form not found: %w", err)
	}
	return s.repo.DeleteEvaluationForm(ctx, uid)
}

// GetEvaluationFormBySession — form + questions untuk session (dipakai FE
// saat menampilkan form evaluasi di detail session).
func (s *Service) GetEvaluationFormBySession(ctx context.Context, sessionID string) (*EvaluationFormWithQuestionsResponse, error) {
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}
	f, err := s.repo.FindEvaluationFormBySession(ctx, sid)
	if err != nil {
		return nil, err
	}
	qs, err := s.repo.ListEvaluationQuestions(ctx, f.ID)
	if err != nil {
		return nil, err
	}
	resp := &EvaluationFormWithQuestionsResponse{
		Form:      *evaluationFormToResponse(f),
		Questions: make([]EvaluationQuestionResponse, 0, len(qs)),
	}
	for _, q := range qs {
		resp.Questions = append(resp.Questions, *evaluationQuestionToResponse(&q))
	}
	return resp, nil
}

// =========================================================================
// Evaluation Questions (P2-BE — plan §22)
// =========================================================================

func (s *Service) CreateEvaluationQuestion(ctx context.Context, formID string, req CreateEvaluationQuestionRequest) (*EvaluationQuestionResponse, error) {
	fid, err := uuid.Parse(formID)
	if err != nil {
		return nil, fmt.Errorf("invalid form_id: %w", err)
	}
	if _, err := s.repo.FindEvaluationFormByID(ctx, fid); err != nil {
		return nil, fmt.Errorf("evaluation form not found: %w", err)
	}
	q := &TrainingEvaluationQuestion{
		FormID:       fid,
		Question:     req.Question,
		QuestionType: EvaluationQuestionType(req.QuestionType),
		IsRequired:   true,
	}
	if req.SortOrder != nil {
		q.SortOrder = *req.SortOrder
	}
	if req.IsRequired != nil {
		q.IsRequired = *req.IsRequired
	}
	if err := s.repo.CreateEvaluationQuestion(ctx, q); err != nil {
		return nil, err
	}
	return evaluationQuestionToResponse(q), nil
}

func (s *Service) ListEvaluationQuestions(ctx context.Context, formID string) ([]EvaluationQuestionResponse, error) {
	fid, err := uuid.Parse(formID)
	if err != nil {
		return nil, fmt.Errorf("invalid form_id: %w", err)
	}
	qs, err := s.repo.ListEvaluationQuestions(ctx, fid)
	if err != nil {
		return nil, err
	}
	responses := make([]EvaluationQuestionResponse, 0, len(qs))
	for _, q := range qs {
		responses = append(responses, *evaluationQuestionToResponse(&q))
	}
	return responses, nil
}

func (s *Service) UpdateEvaluationQuestion(ctx context.Context, id string, req UpdateEvaluationQuestionRequest) (*EvaluationQuestionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	q, err := s.repo.FindEvaluationQuestionByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Question != nil {
		q.Question = *req.Question
	}
	if req.QuestionType != nil {
		q.QuestionType = EvaluationQuestionType(*req.QuestionType)
	}
	if req.SortOrder != nil {
		q.SortOrder = *req.SortOrder
	}
	if req.IsRequired != nil {
		q.IsRequired = *req.IsRequired
	}
	if err := s.repo.UpdateEvaluationQuestion(ctx, q); err != nil {
		return nil, err
	}
	return evaluationQuestionToResponse(q), nil
}

func (s *Service) DeleteEvaluationQuestion(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if _, err := s.repo.FindEvaluationQuestionByID(ctx, uid); err != nil {
		return fmt.Errorf("evaluation question not found: %w", err)
	}
	return s.repo.DeleteEvaluationQuestion(ctx, uid)
}

// =========================================================================
// Evaluation Answers (P2-BE — plan §22)
// =========================================================================

// SubmitEvaluationAnswers — simpan jawaban peserta untuk form tertentu.
// Setiap jawaban di-upsert (satu jawaban per pertanyaan per peserta).
func (s *Service) SubmitEvaluationAnswers(ctx context.Context, formID string, participantID string, req SubmitEvaluationAnswersRequest) ([]EvaluationAnswerResponse, error) {
	fid, err := uuid.Parse(formID)
	if err != nil {
		return nil, fmt.Errorf("invalid form_id: %w", err)
	}
	pid, err := uuid.Parse(participantID)
	if err != nil {
		return nil, fmt.Errorf("invalid participant_id: %w", err)
	}
	if _, err := s.repo.FindEvaluationFormByID(ctx, fid); err != nil {
		return nil, fmt.Errorf("evaluation form not found: %w", err)
	}
	if _, err := s.repo.FindParticipantByID(ctx, pid); err != nil {
		return nil, fmt.Errorf("participant not found: %w", err)
	}
	responses := make([]EvaluationAnswerResponse, 0, len(req.Answers))
	for _, in := range req.Answers {
		qid, err := uuid.Parse(in.QuestionID)
		if err != nil {
			return nil, fmt.Errorf("invalid question_id: %w", err)
		}
		// Pastikan pertanyaan milik form.
		q, err := s.repo.FindEvaluationQuestionByID(ctx, qid)
		if err != nil {
			return nil, fmt.Errorf("question not found: %w", err)
		}
		if q.FormID != fid {
			return nil, fmt.Errorf("question %s does not belong to form %s", qid, fid)
		}
		a := &TrainingEvaluationAnswer{
			QuestionID:    qid,
			ParticipantID: pid,
			Answer:        in.Answer,
		}
		if err := s.repo.UpsertEvaluationAnswer(ctx, a); err != nil {
			return nil, err
		}
		responses = append(responses, *evaluationAnswerToResponse(a))
	}
	return responses, nil
}

func (s *Service) ListEvaluationAnswers(ctx context.Context, questionID *string, participantID *string) ([]EvaluationAnswerResponse, error) {
	var qUUID, pUUID *uuid.UUID
	if questionID != nil && *questionID != "" {
		uid, err := uuid.Parse(*questionID)
		if err != nil {
			return nil, fmt.Errorf("invalid question_id: %w", err)
		}
		qUUID = &uid
	}
	if participantID != nil && *participantID != "" {
		uid, err := uuid.Parse(*participantID)
		if err != nil {
			return nil, fmt.Errorf("invalid participant_id: %w", err)
		}
		pUUID = &uid
	}
	answers, err := s.repo.ListEvaluationAnswers(ctx, qUUID, pUUID)
	if err != nil {
		return nil, err
	}
	responses := make([]EvaluationAnswerResponse, 0, len(answers))
	for _, a := range answers {
		responses = append(responses, *evaluationAnswerToResponse(&a))
	}
	return responses, nil
}

// =========================================================================
// Effectiveness Assessments (P2-BE — plan §23)
// =========================================================================

func (s *Service) CreateEffectivenessAssessment(ctx context.Context, req CreateEffectivenessAssessmentRequest) (*EffectivenessAssessmentResponse, error) {
	pid, err := uuid.Parse(req.ParticipantID)
	if err != nil {
		return nil, fmt.Errorf("invalid participant_id: %w", err)
	}
	if _, err := s.repo.FindParticipantByID(ctx, pid); err != nil {
		return nil, fmt.Errorf("participant not found: %w", err)
	}
	a := &TrainingEffectivenessAssessment{
		ParticipantID:  pid,
		AssessmentDate: req.AssessmentDate,
	}
	if req.AssessorEmployeeID != nil {
		aid, err := uuid.Parse(*req.AssessorEmployeeID)
		if err != nil {
			return nil, fmt.Errorf("invalid assessor_employee_id: %w", err)
		}
		a.AssessorEmployeeID = &aid
	}
	a.BeforeScore = req.BeforeScore
	a.AfterScore = req.AfterScore
	a.EffectivenessScore = req.EffectivenessScore
	a.Remarks = req.Remarks
	if err := s.repo.CreateEffectivenessAssessment(ctx, a); err != nil {
		return nil, err
	}
	return effectivenessToResponse(a), nil
}

func (s *Service) GetEffectivenessAssessmentByID(ctx context.Context, id string) (*EffectivenessAssessmentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindEffectivenessAssessmentByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return effectivenessToResponse(a), nil
}

func (s *Service) ListEffectivenessAssessments(ctx context.Context, participantID *string, page, perPage int) (*PaginatedResponse, error) {
	var pUUID *uuid.UUID
	if participantID != nil && *participantID != "" {
		uid, err := uuid.Parse(*participantID)
		if err != nil {
			return nil, fmt.Errorf("invalid participant_id: %w", err)
		}
		pUUID = &uid
	}
	items, total, err := s.repo.ListEffectivenessAssessments(ctx, pUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]EffectivenessAssessmentResponse, 0, len(items))
	for _, a := range items {
		responses = append(responses, *effectivenessToResponse(&a))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateEffectivenessAssessment(ctx context.Context, id string, req UpdateEffectivenessAssessmentRequest) (*EffectivenessAssessmentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindEffectivenessAssessmentByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.AssessmentDate != nil {
		a.AssessmentDate = *req.AssessmentDate
	}
	if req.AssessorEmployeeID != nil {
		if *req.AssessorEmployeeID == "" {
			a.AssessorEmployeeID = nil
		} else {
			aid, err := uuid.Parse(*req.AssessorEmployeeID)
			if err != nil {
				return nil, fmt.Errorf("invalid assessor_employee_id: %w", err)
			}
			a.AssessorEmployeeID = &aid
		}
	}
	if req.BeforeScore != nil {
		a.BeforeScore = req.BeforeScore
	}
	if req.AfterScore != nil {
		a.AfterScore = req.AfterScore
	}
	if req.EffectivenessScore != nil {
		a.EffectivenessScore = req.EffectivenessScore
	}
	if req.Remarks != nil {
		a.Remarks = req.Remarks
	}
	if err := s.repo.UpdateEffectivenessAssessment(ctx, a); err != nil {
		return nil, err
	}
	return effectivenessToResponse(a), nil
}

func (s *Service) DeleteEffectivenessAssessment(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if _, err := s.repo.FindEffectivenessAssessmentByID(ctx, uid); err != nil {
		return fmt.Errorf("effectiveness assessment not found: %w", err)
	}
	return s.repo.DeleteEffectivenessAssessment(ctx, uid)
}

// =========================================================================
// Certifications (P2-BE — plan §24)
// =========================================================================

func (s *Service) CreateCertification(ctx context.Context, req CreateCertificationRequest) (*CertificationResponse, error) {
	// Kode sertifikasi: bila tidak dikirim, di-generate otomatis CERT-{sekuens}.
	code := strings.TrimSpace(req.Code)
	if code == "" {
		var err error
		code, err = s.repo.NextCertificationCode(ctx)
		if err != nil {
			return nil, err
		}
	}

	c := &TrainingCertification{Code: code, Name: req.Name, IsActive: true}
	c.IssuingBody = req.IssuingBody
	c.ValidityPeriodMonth = req.ValidityPeriodMonth
	// Satuan masa berlaku: default 'month'; 'year'/'month' bila dikirim.
	c.ValidityPeriodUnit = "month"
	if req.ValidityPeriodUnit != nil && *req.ValidityPeriodUnit != "" {
		c.ValidityPeriodUnit = *req.ValidityPeriodUnit
	}
	if req.RenewalRequired != nil {
		c.RenewalRequired = *req.RenewalRequired
	}
	if req.IsActive != nil {
		c.IsActive = *req.IsActive
	}
	if err := s.repo.CreateCertification(ctx, c); err != nil {
		return nil, err
	}
	return certificationToResponse(c), nil
}

func (s *Service) GetCertificationByID(ctx context.Context, id string) (*CertificationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	c, err := s.repo.FindCertificationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return certificationToResponse(c), nil
}

func (s *Service) ListCertifications(ctx context.Context, isActive *bool, page, perPage int) (*PaginatedResponse, error) {
	items, total, err := s.repo.ListCertifications(ctx, isActive, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]CertificationResponse, 0, len(items))
	for _, c := range items {
		responses = append(responses, *certificationToResponse(&c))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateCertification(ctx context.Context, id string, req UpdateCertificationRequest) (*CertificationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	c, err := s.repo.FindCertificationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Code != nil {
		c.Code = *req.Code
	}
	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.IssuingBody != nil {
		c.IssuingBody = req.IssuingBody
	}
	if req.ValidityPeriodMonth != nil {
		c.ValidityPeriodMonth = req.ValidityPeriodMonth
	}
	if req.ValidityPeriodUnit != nil {
		c.ValidityPeriodUnit = *req.ValidityPeriodUnit
	}
	if req.RenewalRequired != nil {
		c.RenewalRequired = *req.RenewalRequired
	}
	if req.IsActive != nil {
		c.IsActive = *req.IsActive
	}
	if err := s.repo.UpdateCertification(ctx, c); err != nil {
		return nil, err
	}
	return certificationToResponse(c), nil
}

func (s *Service) DeleteCertification(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if _, err := s.repo.FindCertificationByID(ctx, uid); err != nil {
		return fmt.Errorf("certification not found: %w", err)
	}
	return s.repo.DeleteCertification(ctx, uid)
}

// =========================================================================
// Certificates — generate dari completion (P2-BE — plan §24)
// =========================================================================

// GenerateCertificate — buat sertifikat dari participant yang COMPLETED.
// CertificateNo dibuat otomatis (TRN-{participantID:8}).
func (s *Service) GenerateCertificate(ctx context.Context, participantID string, req GenerateCertificateRequest) (*TrainingCertificateResponse, error) {
	pid, err := uuid.Parse(participantID)
	if err != nil {
		return nil, fmt.Errorf("invalid participant_id: %w", err)
	}
	p, err := s.repo.FindParticipantByID(ctx, pid)
	if err != nil {
		return nil, fmt.Errorf("participant not found: %w", err)
	}
	if p.CompletionStatus != CompletionCompleted {
		return nil, fmt.Errorf("certificate can only be generated for completed participants")
	}
	// Idempotent: bila sudah ada, update saja.
	cert, err := s.repo.FindCertificateByParticipant(ctx, pid)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if cert == nil {
		cert = &TrainingCertificate{ParticipantID: pid}
	}
	if cert.CertificateNo == "" {
		cert.CertificateNo = fmt.Sprintf("TRN-%s", strings.ToUpper(pid.String()[:8]))
	}
	if req.CertificationID != nil && *req.CertificationID != "" {
		cid, err := uuid.Parse(*req.CertificationID)
		if err != nil {
			return nil, fmt.Errorf("invalid certification_id: %w", err)
		}
		if _, err := s.repo.FindCertificationByID(ctx, cid); err != nil {
			return nil, fmt.Errorf("certification not found: %w", err)
		}
		cert.CertificationID = &cid
	}
	cert.CertificateFileURL = req.CertificateFileURL
	if req.ExpiryDate != nil {
		cert.ExpiryDate = req.ExpiryDate
	}
	now := time.Now()
	cert.IssuedDate = now.Format("2006-01-02")
	if cert.ID == uuid.Nil {
		if err := s.repo.CreateCertificate(ctx, cert); err != nil {
			return nil, err
		}
	} else {
		if err := s.repo.UpdateCertificate(ctx, cert); err != nil {
			return nil, err
		}
	}
	return certificateToResponse(cert), nil
}

func (s *Service) UpdateCertificateFile(ctx context.Context, id string, req UpdateCertificateRequest) (*TrainingCertificateResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	cert, err := s.repo.FindCertificateByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.CertificateFileURL != nil {
		cert.CertificateFileURL = req.CertificateFileURL
	}
	if req.ExpiryDate != nil {
		cert.ExpiryDate = req.ExpiryDate
	}
	if err := s.repo.UpdateCertificate(ctx, cert); err != nil {
		return nil, err
	}
	return certificateToResponse(cert), nil
}

// =========================================================================
// Reports & History (P2-BE — plan §38)
// =========================================================================

func (s *Service) GetTrainingHistory(ctx context.Context, employeeID string) ([]TrainingHistoryResponse, error) {
	uid, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	return s.repo.HistoryByEmployee(ctx, uid)
}

func (s *Service) GetParticipationReport(ctx context.Context, sessionStatus *string) ([]ParticipationReportRow, error) {
	return s.repo.ParticipationReport(ctx, sessionStatus)
}

func (s *Service) GetCostReport(ctx context.Context) ([]CostReportRow, error) {
	rows, err := s.repo.CostReport(ctx)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		if rows[i].ParticipantCount > 0 {
			rows[i].CostPerParticipant = round2(rows[i].TotalCost / float64(rows[i].ParticipantCount))
		}
	}
	return rows, nil
}

func (s *Service) GetComplianceReport(ctx context.Context) ([]ComplianceReportRow, error) {
	return s.repo.ComplianceReport(ctx)
}

func (s *Service) GetDashboardReport(ctx context.Context) (*DashboardReport, error) {
	return s.repo.DashboardReport(ctx)
}

// =========================================================================
// Converters P2
// =========================================================================

func evaluationFormToResponse(f *TrainingEvaluationForm) *EvaluationFormResponse {
	return &EvaluationFormResponse{
		ID:        f.ID.String(),
		SessionID: f.SessionID.String(),
		Name:      f.Name,
		IsActive:  f.IsActive,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func evaluationQuestionToResponse(q *TrainingEvaluationQuestion) *EvaluationQuestionResponse {
	return &EvaluationQuestionResponse{
		ID:           q.ID.String(),
		FormID:       q.FormID.String(),
		Question:     q.Question,
		QuestionType: string(q.QuestionType),
		SortOrder:    q.SortOrder,
		IsRequired:   q.IsRequired,
		CreatedAt:    q.CreatedAt,
		UpdatedAt:    q.UpdatedAt,
	}
}

func evaluationAnswerToResponse(a *TrainingEvaluationAnswer) *EvaluationAnswerResponse {
	return &EvaluationAnswerResponse{
		ID:            a.ID.String(),
		QuestionID:    a.QuestionID.String(),
		ParticipantID: a.ParticipantID.String(),
		Answer:        a.Answer,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

func effectivenessToResponse(a *TrainingEffectivenessAssessment) *EffectivenessAssessmentResponse {
	resp := &EffectivenessAssessmentResponse{
		ID:                 a.ID.String(),
		ParticipantID:      a.ParticipantID.String(),
		AssessmentDate:     a.AssessmentDate,
		BeforeScore:        a.BeforeScore,
		AfterScore:         a.AfterScore,
		EffectivenessScore: a.EffectivenessScore,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}
	if a.AssessorEmployeeID != nil {
		resp.AssessorEmployeeID = a.AssessorEmployeeID.String()
	}
	if a.Remarks != nil {
		resp.Remarks = *a.Remarks
	}
	return resp
}

func certificationToResponse(c *TrainingCertification) *CertificationResponse {
	resp := &CertificationResponse{
		ID:                  c.ID.String(),
		Code:                c.Code,
		Name:                c.Name,
		ValidityPeriodMonth: c.ValidityPeriodMonth,
		ValidityPeriodUnit:  c.ValidityPeriodUnit,
		RenewalRequired:     c.RenewalRequired,
		IsActive:            c.IsActive,
		CreatedAt:           c.CreatedAt,
		UpdatedAt:           c.UpdatedAt,
	}
	if c.IssuingBody != nil {
		resp.IssuingBody = *c.IssuingBody
	}
	return resp
}
