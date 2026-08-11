package training

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	dbFunc TenantDBFunc
}

func NewRepository(dbFunc TenantDBFunc) *Repository {
	return &Repository{dbFunc: dbFunc}
}

func (r *Repository) db(ctx context.Context) (*gorm.DB, error) {
	return r.dbFunc(ctx)
}

// =========================================================================
// Training Categories
// =========================================================================

func (r *Repository) CreateCategory(ctx context.Context, c *TrainingCategory) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(c).Error
}

func (r *Repository) FindCategoryByID(ctx context.Context, id uuid.UUID) (*TrainingCategory, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var cat TrainingCategory
	if err := db.WithContext(ctx).First(&cat, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("training category not found")
		}
		return nil, err
	}
	return &cat, nil
}

func (r *Repository) ListCategories(ctx context.Context, page, perPage int) ([]TrainingCategory, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var cats []TrainingCategory
	var total int64
	query := db.WithContext(ctx).Model(&TrainingCategory{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("code ASC").Find(&cats).Error; err != nil {
		return nil, 0, err
	}
	return cats, total, nil
}

func (r *Repository) UpdateCategory(ctx context.Context, c *TrainingCategory) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(c).Error
}

func (r *Repository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&TrainingCategory{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("training category not found")
	}
	return result.Error
}

// =========================================================================
// Training Courses
// =========================================================================

func (r *Repository) CreateCourse(ctx context.Context, c *TrainingCourse) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(c).Error
}

func (r *Repository) FindCourseByID(ctx context.Context, id uuid.UUID) (*TrainingCourse, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var course TrainingCourse
	if err := db.WithContext(ctx).First(&course, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("training course not found")
		}
		return nil, err
	}
	return &course, nil
}

func (r *Repository) ListCourses(ctx context.Context, categoryID *uuid.UUID, page, perPage int) ([]TrainingCourse, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var courses []TrainingCourse
	var total int64
	query := db.WithContext(ctx).Model(&TrainingCourse{})
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("code ASC").Find(&courses).Error; err != nil {
		return nil, 0, err
	}
	return courses, total, nil
}

func (r *Repository) UpdateCourse(ctx context.Context, c *TrainingCourse) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(c).Error
}

func (r *Repository) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&TrainingCourse{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("training course not found")
	}
	return result.Error
}

// =========================================================================
// Training Sessions
// =========================================================================

func (r *Repository) CreateSession(ctx context.Context, s *TrainingSession) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(s).Error
}

func (r *Repository) FindSessionByID(ctx context.Context, id uuid.UUID) (*TrainingSession, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var sess TrainingSession
	if err := db.WithContext(ctx).First(&sess, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("training session not found")
		}
		return nil, err
	}
	return &sess, nil
}

func (r *Repository) ListSessions(ctx context.Context, courseID *uuid.UUID, status *string, page, perPage int) ([]TrainingSession, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var sessions []TrainingSession
	var total int64
	query := db.WithContext(ctx).Model(&TrainingSession{})
	if courseID != nil {
		query = query.Where("course_id = ?", *courseID)
	}
	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("start_date DESC").Find(&sessions).Error; err != nil {
		return nil, 0, err
	}
	return sessions, total, nil
}

func (r *Repository) UpdateSession(ctx context.Context, s *TrainingSession) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(s).Error
}

func (r *Repository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&TrainingSession{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("training session not found")
	}
	return result.Error
}

// =========================================================================
// Training Participants
// =========================================================================

func (r *Repository) CreateParticipant(ctx context.Context, p *TrainingParticipant) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(p).Error
}

func (r *Repository) FindParticipantByID(ctx context.Context, id uuid.UUID) (*TrainingParticipant, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var p TrainingParticipant
	if err := db.WithContext(ctx).First(&p, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("training participant not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *Repository) ListParticipants(ctx context.Context, sessionID *uuid.UUID, employeeID *uuid.UUID, page, perPage int) ([]TrainingParticipant, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var participants []TrainingParticipant
	var total int64
	query := db.WithContext(ctx).Model(&TrainingParticipant{})
	if sessionID != nil {
		query = query.Where("session_id = ?", *sessionID)
	}
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&participants).Error; err != nil {
		return nil, 0, err
	}
	return participants, total, nil
}

func (r *Repository) UpdateParticipant(ctx context.Context, p *TrainingParticipant) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(p).Error
}

func (r *Repository) DeleteParticipant(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&TrainingParticipant{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("training participant not found")
	}
	return result.Error
}

// =========================================================================
// Training Materials
// =========================================================================

func (r *Repository) CreateMaterial(ctx context.Context, m *TrainingMaterial) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(m).Error
}

func (r *Repository) FindMaterialByID(ctx context.Context, id uuid.UUID) (*TrainingMaterial, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var mat TrainingMaterial
	if err := db.WithContext(ctx).First(&mat, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("training material not found")
		}
		return nil, err
	}
	return &mat, nil
}

func (r *Repository) ListMaterials(ctx context.Context, sessionID uuid.UUID) ([]TrainingMaterial, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var mats []TrainingMaterial
	if err := db.WithContext(ctx).Where("session_id = ?", sessionID).
		Order("sort_order ASC, created_at ASC").Find(&mats).Error; err != nil {
		return nil, err
	}
	return mats, nil
}

func (r *Repository) UpdateMaterial(ctx context.Context, m *TrainingMaterial) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(m).Error
}

func (r *Repository) DeleteMaterial(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&TrainingMaterial{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("training material not found")
	}
	return result.Error
}

// =========================================================================
// Training Evaluations
// =========================================================================

func (r *Repository) CreateEvaluation(ctx context.Context, e *TrainingEvaluation) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(e).Error
}

func (r *Repository) FindEvaluationByID(ctx context.Context, id uuid.UUID) (*TrainingEvaluation, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var eval TrainingEvaluation
	if err := db.WithContext(ctx).First(&eval, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("training evaluation not found")
		}
		return nil, err
	}
	return &eval, nil
}

func (r *Repository) ListEvaluations(ctx context.Context, sessionID *uuid.UUID, employeeID *uuid.UUID, page, perPage int) ([]TrainingEvaluation, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var evals []TrainingEvaluation
	var total int64
	query := db.WithContext(ctx).Model(&TrainingEvaluation{})
	if sessionID != nil {
		query = query.Where("session_id = ?", *sessionID)
	}
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&evals).Error; err != nil {
		return nil, 0, err
	}
	return evals, total, nil
}

func (r *Repository) UpdateEvaluation(ctx context.Context, e *TrainingEvaluation) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(e).Error
}

func (r *Repository) DeleteEvaluation(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&TrainingEvaluation{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("training evaluation not found")
	}
	return result.Error
}

// =========================================================================
// Training Certificates
// =========================================================================

func (r *Repository) CreateCertificate(ctx context.Context, c *TrainingCertificate) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(c).Error
}

func (r *Repository) FindCertificateByID(ctx context.Context, id uuid.UUID) (*TrainingCertificate, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var cert TrainingCertificate
	if err := db.WithContext(ctx).First(&cert, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("training certificate not found")
		}
		return nil, err
	}
	return &cert, nil
}

func (r *Repository) ListCertificates(ctx context.Context, participantID *uuid.UUID, page, perPage int) ([]TrainingCertificate, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var certs []TrainingCertificate
	var total int64
	query := db.WithContext(ctx).Model(&TrainingCertificate{})
	if participantID != nil {
		query = query.Where("participant_id = ?", *participantID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("issued_date DESC").Find(&certs).Error; err != nil {
		return nil, 0, err
	}
	return certs, total, nil
}

func (r *Repository) UpdateCertificate(ctx context.Context, c *TrainingCertificate) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(c).Error
}

func (r *Repository) DeleteCertificate(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&TrainingCertificate{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("training certificate not found")
	}
	return result.Error
}

// =========================================================================
// Aggregation
// =========================================================================

func (r *Repository) CountParticipantsBySession(ctx context.Context, sessionID uuid.UUID) (int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return 0, err
	}
	var count int64
	if err := db.WithContext(ctx).Model(&TrainingParticipant{}).
		Where("session_id = ?", sessionID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindParticipantBySessionAndEmployee — cek duplikasi enrollment (record aktif).
func (r *Repository) FindParticipantBySessionAndEmployee(ctx context.Context, sessionID, employeeID uuid.UUID) (*TrainingParticipant, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var p TrainingParticipant
	if err := db.WithContext(ctx).
		Where("session_id = ? AND employee_id = ?", sessionID, employeeID).
		First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// =========================================================================
// Training Providers (P0-BE — plan §11)
// =========================================================================

func (r *Repository) CreateProvider(ctx context.Context, p *TrainingProvider) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(p).Error
}

func (r *Repository) FindProviderByID(ctx context.Context, id uuid.UUID) (*TrainingProvider, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var p TrainingProvider
	if err := db.WithContext(ctx).First(&p, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("training provider not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *Repository) ListProviders(ctx context.Context, page, perPage int) ([]TrainingProvider, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var providers []TrainingProvider
	var total int64
	query := db.WithContext(ctx).Model(&TrainingProvider{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("code ASC").Find(&providers).Error; err != nil {
		return nil, 0, err
	}
	return providers, total, nil
}

func (r *Repository) UpdateProvider(ctx context.Context, p *TrainingProvider) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(p).Error
}

func (r *Repository) DeleteProvider(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&TrainingProvider{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("training provider not found")
	}
	return result.Error
}

// =========================================================================
// Training Trainers (P0-BE — plan §12)
// =========================================================================

func (r *Repository) CreateTrainer(ctx context.Context, t *TrainingTrainer) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(t).Error
}

func (r *Repository) FindTrainerByID(ctx context.Context, id uuid.UUID) (*TrainingTrainer, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var t TrainingTrainer
	if err := db.WithContext(ctx).First(&t, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("training trainer not found")
		}
		return nil, err
	}
	return &t, nil
}

func (r *Repository) ListTrainers(ctx context.Context, trainerType *string, page, perPage int) ([]TrainingTrainer, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var trainers []TrainingTrainer
	var total int64
	query := db.WithContext(ctx).Model(&TrainingTrainer{})
	if trainerType != nil && *trainerType != "" {
		query = query.Where("type = ?", *trainerType)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("name ASC").Find(&trainers).Error; err != nil {
		return nil, 0, err
	}
	return trainers, total, nil
}

func (r *Repository) UpdateTrainer(ctx context.Context, t *TrainingTrainer) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(t).Error
}

func (r *Repository) DeleteTrainer(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&TrainingTrainer{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("training trainer not found")
	}
	return result.Error
}

// =========================================================================
// Training Session Trainers (P0-BE — plan §13)
// =========================================================================

func (r *Repository) CreateSessionTrainer(ctx context.Context, st *TrainingSessionTrainer) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(st).Error
}

func (r *Repository) FindSessionTrainerByID(ctx context.Context, id uuid.UUID) (*TrainingSessionTrainer, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var st TrainingSessionTrainer
	if err := db.WithContext(ctx).First(&st, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("training session trainer not found")
		}
		return nil, err
	}
	return &st, nil
}

// FindSessionTrainerBySessionAndTrainer — cegah trainer ditambahkan 2x di session.
func (r *Repository) FindSessionTrainerBySessionAndTrainer(ctx context.Context, sessionID, trainerID uuid.UUID) (*TrainingSessionTrainer, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var st TrainingSessionTrainer
	if err := db.WithContext(ctx).
		Where("session_id = ? AND trainer_id = ?", sessionID, trainerID).
		First(&st).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &st, nil
}

func (r *Repository) ListSessionTrainers(ctx context.Context, sessionID uuid.UUID) ([]TrainingSessionTrainer, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var items []TrainingSessionTrainer
	if err := db.WithContext(ctx).Where("session_id = ?", sessionID).
		Order("role ASC, created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) DeleteSessionTrainer(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&TrainingSessionTrainer{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("training session trainer not found")
	}
	return result.Error
}

// =========================================================================
// Training Attendances (P0-BE — plan §19)
// =========================================================================

// FindAttendanceByParticipantAndDate — untuk upsert satu baris per participant per hari.
func (r *Repository) FindAttendanceByParticipantAndDate(ctx context.Context, participantID uuid.UUID, date string) (*TrainingAttendance, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var a TrainingAttendance
	if err := db.WithContext(ctx).
		Where("participant_id = ? AND attendance_date = ?", participantID, date).
		First(&a).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *Repository) CreateAttendance(ctx context.Context, a *TrainingAttendance) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(a).Error
}

func (r *Repository) FindAttendanceByID(ctx context.Context, id uuid.UUID) (*TrainingAttendance, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var a TrainingAttendance
	if err := db.WithContext(ctx).First(&a, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("training attendance not found")
		}
		return nil, err
	}
	return &a, nil
}

func (r *Repository) UpdateAttendance(ctx context.Context, a *TrainingAttendance) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(a).Error
}

// attendanceScanRow — baris hasil scan join attendance ↔ participant.
// (CheckIn/CheckOut bertipe *time.Time agar scan kolom timestamp aman lintas driver.)
type attendanceScanRow struct {
	AttendanceID   string
	ParticipantID  string
	EmployeeID     string
	AttendanceDate string
	Status         string
	CheckIn        *time.Time
	CheckOut       *time.Time
	Remarks        *string
}

// ListAttendanceBySession — attendance per session (join participant untuk employee_id).
func (r *Repository) ListAttendanceBySession(ctx context.Context, sessionID uuid.UUID) ([]attendanceScanRow, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []attendanceScanRow
	query := db.WithContext(ctx).
		Table("training_attendances ta").
		Select("ta.id AS attendance_id, ta.participant_id, tp.employee_id, " +
			"ta.attendance_date, ta.status, ta.check_in, ta.check_out, ta.remarks").
		Joins("JOIN training_participants tp ON tp.id = ta.participant_id").
		Where("tp.session_id = ?", sessionID).
		Order("ta.attendance_date ASC, tp.employee_id ASC")
	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// =========================================================================
// Training Assessments (P0-BE — plan §21)
// =========================================================================

func (r *Repository) CreateAssessment(ctx context.Context, a *TrainingAssessment) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(a).Error
}

func (r *Repository) FindAssessmentByID(ctx context.Context, id uuid.UUID) (*TrainingAssessment, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var a TrainingAssessment
	if err := db.WithContext(ctx).First(&a, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("training assessment not found")
		}
		return nil, err
	}
	return &a, nil
}

func (r *Repository) ListAssessmentsBySession(ctx context.Context, sessionID uuid.UUID) ([]TrainingAssessment, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var items []TrainingAssessment
	if err := db.WithContext(ctx).Where("session_id = ?", sessionID).
		Order("created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) CreateAssessmentResult(ctx context.Context, res *TrainingAssessmentResult) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(res).Error
}

// CountAssessmentAttempts — jumlah attempt yang sudah dicatat untuk (assessment, participant).
func (r *Repository) CountAssessmentAttempts(ctx context.Context, assessmentID, participantID uuid.UUID) (int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return 0, err
	}
	var count int64
	if err := db.WithContext(ctx).Model(&TrainingAssessmentResult{}).
		Where("assessment_id = ? AND participant_id = ?", assessmentID, participantID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository) AvgEvaluationRatingBySession(ctx context.Context, sessionID uuid.UUID) (float64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return 0, err
	}
	var avg float64
	if err := db.WithContext(ctx).Model(&TrainingEvaluation{}).
		Where("session_id = ?", sessionID).
		Select("COALESCE(AVG(rating), 0)").Scan(&avg).Error; err != nil {
		return 0, err
	}
	return avg, nil
}
