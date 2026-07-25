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
	count, err := s.repo.CountParticipantsBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if count >= int64(sess.MaxQuota) {
		return nil, fmt.Errorf("session quota full (%d/%d)", count, sess.MaxQuota)
	}

	p := &TrainingParticipant{
		SessionID:        sessionID,
		EmployeeID:       empID,
		AttendanceStatus: AttendStatusPresent,
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
		// Auto-mark completed with today's date when score is entered
		completed := time.Now().Format("2006-01-02")
		p.CompletedAt = &completed
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
		ID:          s.ID.String(),
		CourseID:    s.CourseID.String(),
		SessionCode: s.SessionCode,
		TrainerName: s.TrainerName,
		Location:    loc,
		StartDate:   s.StartDate,
		EndDate:     s.EndDate,
		MaxQuota:    s.MaxQuota,
		Status:      string(s.Status),
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func participantToResponse(p *TrainingParticipant) *TrainingParticipantResponse {
	completed := ""
	if p.CompletedAt != nil {
		completed = *p.CompletedAt
	}
	return &TrainingParticipantResponse{
		ID:               p.ID.String(),
		SessionID:        p.SessionID.String(),
		EmployeeID:       p.EmployeeID.String(),
		AttendanceStatus: string(p.AttendanceStatus),
		Score:            p.Score,
		CompletedAt:      completed,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
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
		ID:        m.ID.String(),
		SessionID: m.SessionID.String(),
		Title:     m.Title,
		FileURL:   url,
		FileType:  ft,
		SortOrder: m.SortOrder,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
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
