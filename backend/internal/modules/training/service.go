package training

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

type Service struct {
	repo   *Repository
	logger *zap.Logger
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// =========================================================================
// Training Categories
// =========================================================================

func (s *Service) CreateCategory(ctx context.Context, req CreateTrainingCategoryRequest) (*TrainingCategoryResponse, error) {
	cat := &TrainingCategory{
		Code:     req.Code,
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
	// Verify category exists
	if _, err := s.repo.FindCategoryByID(ctx, catID); err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	course := &TrainingCourse{
		CategoryID:    catID,
		Code:          req.Code,
		Name:          req.Name,
		IsCertified:   false,
		IsActive:      true,
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
	if _, err := s.repo.FindCourseByID(ctx, courseID); err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	maxQuota := 30
	if req.MaxQuota > 0 {
		maxQuota = req.MaxQuota
	}

	sess := &TrainingSession{
		CourseID:    courseID,
		SessionCode: req.SessionCode,
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
		ID:                  s.ID.String(),
		CourseID:            s.CourseID.String(),
		SessionCode:         s.SessionCode,
		TrainerName:         s.TrainerName,
		ProviderType:        strPtr(s.ProviderType),
		DeliveryMode:        strPtr(s.DeliveryMode),
		ProviderID:          uuidPtr(s.ProviderID),
		StartDatetime:       formatTimePtr(s.StartDatetime),
		EndDatetime:         formatTimePtr(s.EndDatetime),
		MeetingURL:          strPtr(s.MeetingURL),
		RegistrationDeadline: formatTimePtr(s.RegistrationDeadline),
		Location:            loc,
		StartDate:           s.StartDate,
		EndDate:             s.EndDate,
		MaxQuota:            s.MaxQuota,
		Status:              string(s.Status),
		CreatedAt:           s.CreatedAt,
		UpdatedAt:           s.UpdatedAt,
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
	return &TrainingCertificateResponse{
		ID:            c.ID.String(),
		ParticipantID: c.ParticipantID.String(),
		CertificateNo: c.CertificateNo,
		IssuedDate:    c.IssuedDate,
		ExpiryDate:    exp,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
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
	p := &TrainingProvider{
		Code:     req.Code,
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
