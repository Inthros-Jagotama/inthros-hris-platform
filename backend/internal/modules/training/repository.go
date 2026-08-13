package training

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
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

// NextCategoryCode menghasilkan kode kategori otomatis dengan pola CAT-{NNN}
// (mis. CAT-001). Sekuens dihitung dari kode kategori yang sudah ada sehingga
// tetap unik pada unique index.
func (r *Repository) NextCategoryCode(ctx context.Context) (string, error) {
	db, err := r.db(ctx)
	if err != nil {
		return "", err
	}

	var codes []string
	if err := db.WithContext(ctx).Model(&TrainingCategory{}).
		Where("code LIKE ?", "CAT-%").
		Pluck("code", &codes).Error; err != nil {
		return "", err
	}

	maxSeq := 0
	for _, c := range codes {
		n, err := strconv.Atoi(strings.TrimPrefix(c, "CAT-"))
		if err == nil && n > maxSeq {
			maxSeq = n
		}
	}
	return fmt.Sprintf("CAT-%03d", maxSeq+1), nil
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

// NextCourseCode menghasilkan kode kursus otomatis dengan pola {KODE_KATEGORI}-{NNN}
// (mis. TECH-001). Sekuens dihitung dari kode kursus yang sudah ada dengan prefix
// yang sama (global, bukan per-kategori) sehingga kode tetap unik pada unique index.
func (r *Repository) NextCourseCode(ctx context.Context, categoryCode string) (string, error) {
	db, err := r.db(ctx)
	if err != nil {
		return "", err
	}

	prefix := strings.ToUpper(strings.TrimSpace(categoryCode))
	if prefix == "" {
		prefix = "TRN"
	}
	// Batasi prefix agar panjang kode (varchar(20)) tidak terlampaui: 16 + "-" + 3 digit.
	if len(prefix) > 16 {
		prefix = prefix[:16]
	}

	var codes []string
	if err := db.WithContext(ctx).Model(&TrainingCourse{}).
		Where("code LIKE ?", prefix+"-%").
		Pluck("code", &codes).Error; err != nil {
		return "", err
	}

	maxSeq := 0
	for _, c := range codes {
		n, err := strconv.Atoi(strings.TrimPrefix(c, prefix+"-"))
		if err == nil && n > maxSeq {
			maxSeq = n
		}
	}
	return fmt.Sprintf("%s-%03d", prefix, maxSeq+1), nil
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

// NextSessionCode menghasilkan kode sesi otomatis dengan pola {KODE_KURSUS}-{NNN}
// (mis. TECH-001-001). Sekuens dihitung dari kode sesi yang sudah ada dengan prefix
// yang sama sehingga kode tetap unik.
func (r *Repository) NextSessionCode(ctx context.Context, courseCode string) (string, error) {
	db, err := r.db(ctx)
	if err != nil {
		return "", err
	}

	prefix := strings.ToUpper(strings.TrimSpace(courseCode))
	if prefix == "" {
		prefix = "SES"
	}
	// Batasi prefix agar panjang kode (varchar(20)) tidak terlampaui: 16 + "-" + 3 digit.
	if len(prefix) > 16 {
		prefix = prefix[:16]
	}

	var codes []string
	if err := db.WithContext(ctx).Model(&TrainingSession{}).
		Where("session_code LIKE ?", prefix+"-%").
		Pluck("session_code", &codes).Error; err != nil {
		return "", err
	}

	maxSeq := 0
	for _, c := range codes {
		n, err := strconv.Atoi(strings.TrimPrefix(c, prefix+"-"))
		if err == nil && n > maxSeq {
			maxSeq = n
		}
	}
	return fmt.Sprintf("%s-%03d", prefix, maxSeq+1), nil
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

// NextProviderCode menghasilkan kode penyelenggara otomatis dengan pola PRV-{NNN}
// (mis. PRV-001). Sekuens dihitung dari kode yang sudah ada sehingga tetap unik
// pada unique index.
func (r *Repository) NextProviderCode(ctx context.Context) (string, error) {
	db, err := r.db(ctx)
	if err != nil {
		return "", err
	}

	var codes []string
	if err := db.WithContext(ctx).Model(&TrainingProvider{}).
		Where("code LIKE ?", "PRV-%").
		Pluck("code", &codes).Error; err != nil {
		return "", err
	}

	maxSeq := 0
	for _, c := range codes {
		n, err := strconv.Atoi(strings.TrimPrefix(c, "PRV-"))
		if err == nil && n > maxSeq {
			maxSeq = n
		}
	}
	return fmt.Sprintf("PRV-%03d", maxSeq+1), nil
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
		Select("ta.id AS attendance_id, ta.participant_id, tp.employee_id, "+
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

// =========================================================================
// Training Plans (P1-BE — plan §16)
// =========================================================================

func (r *Repository) CreatePlan(ctx context.Context, p *TrainingPlan) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(p).Error
}

// NextPlanCode menghasilkan kode plan otomatis dengan pola TP-{tahun}-{NNN}
// (mis. TP-2026-001). Sekuens dihitung dari kode plan yang sudah ada dengan
// prefix tahun yang sama sehingga kode tetap unik.
func (r *Repository) NextPlanCode(ctx context.Context, year int) (string, error) {
	db, err := r.db(ctx)
	if err != nil {
		return "", err
	}

	prefix := fmt.Sprintf("TP-%d-", year)

	var codes []string
	if err := db.WithContext(ctx).Model(&TrainingPlan{}).
		Where("code LIKE ?", prefix+"%").
		Pluck("code", &codes).Error; err != nil {
		return "", err
	}

	maxSeq := 0
	for _, c := range codes {
		n, err := strconv.Atoi(strings.TrimPrefix(c, prefix))
		if err == nil && n > maxSeq {
			maxSeq = n
		}
	}
	return fmt.Sprintf("%s%03d", prefix, maxSeq+1), nil
}

func (r *Repository) FindPlanByID(ctx context.Context, id uuid.UUID) (*TrainingPlan, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var p TrainingPlan
	if err := db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) ListPlans(ctx context.Context, year *int, status *string, page, perPage int) ([]TrainingPlan, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	q := db.WithContext(ctx).Model(&TrainingPlan{})
	if year != nil {
		q = q.Where("year = ?", *year)
	}
	if status != nil && *status != "" {
		q = q.Where("status = ?", *status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []TrainingPlan
	if err := q.Order("year DESC, created_at DESC").
		Offset((page - 1) * perPage).Limit(perPage).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) UpdatePlan(ctx context.Context, p *TrainingPlan) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(p).Error
}

func (r *Repository) DeletePlan(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Delete(&TrainingPlan{}, "id = ?", id).Error
}

// =========================================================================
// Training Plan Items (P1-BE — plan §16)
// =========================================================================

func (r *Repository) CreatePlanItem(ctx context.Context, item *TrainingPlanItem) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(item).Error
}

func (r *Repository) FindPlanItemByID(ctx context.Context, id uuid.UUID) (*TrainingPlanItem, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var item TrainingPlanItem
	if err := db.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository) ListPlanItems(ctx context.Context, planID uuid.UUID) ([]TrainingPlanItem, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var items []TrainingPlanItem
	if err := db.WithContext(ctx).
		Where("training_plan_id = ?", planID).
		Order("created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) UpdatePlanItem(ctx context.Context, item *TrainingPlanItem) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(item).Error
}

func (r *Repository) DeletePlanItem(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Delete(&TrainingPlanItem{}, "id = ?", id).Error
}

// =========================================================================
// Training Needs (P1-BE — plan §17)
// =========================================================================

func (r *Repository) CreateNeed(ctx context.Context, n *TrainingNeed) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(n).Error
}

func (r *Repository) FindNeedByID(ctx context.Context, id uuid.UUID) (*TrainingNeed, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var n TrainingNeed
	if err := db.WithContext(ctx).First(&n, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *Repository) ListNeeds(ctx context.Context, employeeID, courseID *uuid.UUID, status *string, page, perPage int) ([]TrainingNeed, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	q := db.WithContext(ctx).Model(&TrainingNeed{})
	if employeeID != nil {
		q = q.Where("employee_id = ?", *employeeID)
	}
	if courseID != nil {
		q = q.Where("course_id = ?", *courseID)
	}
	if status != nil && *status != "" {
		q = q.Where("status = ?", *status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []TrainingNeed
	if err := q.Order("created_at DESC").
		Offset((page - 1) * perPage).Limit(perPage).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) UpdateNeed(ctx context.Context, n *TrainingNeed) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(n).Error
}

// =========================================================================
// Training Requests (P1-BE — plan §15, Central Approval)
// =========================================================================

func (r *Repository) CreateRequest(ctx context.Context, req *TrainingRequest) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(req).Error
}

func (r *Repository) FindRequestByID(ctx context.Context, id uuid.UUID) (*TrainingRequest, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var req TrainingRequest
	if err := db.WithContext(ctx).First(&req, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *Repository) ListRequests(ctx context.Context, employeeID *uuid.UUID, status *string, page, perPage int) ([]TrainingRequest, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	q := db.WithContext(ctx).Model(&TrainingRequest{})
	if employeeID != nil {
		q = q.Where("employee_id = ?", *employeeID)
	}
	if status != nil && *status != "" {
		q = q.Where("status = ?", *status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []TrainingRequest
	if err := q.Order("created_at DESC").
		Offset((page - 1) * perPage).Limit(perPage).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) UpdateRequest(ctx context.Context, req *TrainingRequest) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(req).Error
}

// =========================================================================
// Course Sub-resources (P1-BE — plan §8/§9/§10)
// =========================================================================

func (r *Repository) CreateCourseObjective(ctx context.Context, o *TrainingCourseObjective) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(o).Error
}

func (r *Repository) FindCourseObjectiveByID(ctx context.Context, id uuid.UUID) (*TrainingCourseObjective, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var o TrainingCourseObjective
	if err := db.WithContext(ctx).First(&o, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *Repository) ListCourseObjectives(ctx context.Context, courseID uuid.UUID) ([]TrainingCourseObjective, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var items []TrainingCourseObjective
	if err := db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Order("sort_order ASC, created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) UpdateCourseObjective(ctx context.Context, o *TrainingCourseObjective) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(o).Error
}

func (r *Repository) DeleteCourseObjective(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Delete(&TrainingCourseObjective{}, "id = ?", id).Error
}

func (r *Repository) CreateCourseCompetency(ctx context.Context, c *TrainingCourseCompetency) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(c).Error
}

func (r *Repository) FindCourseCompetencyByID(ctx context.Context, id uuid.UUID) (*TrainingCourseCompetency, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var c TrainingCourseCompetency
	if err := db.WithContext(ctx).First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) ListCourseCompetencies(ctx context.Context, courseID uuid.UUID) ([]TrainingCourseCompetency, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var items []TrainingCourseCompetency
	if err := db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Order("created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) DeleteCourseCompetency(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Delete(&TrainingCourseCompetency{}, "id = ?", id).Error
}

func (r *Repository) CreateCoursePrerequisite(ctx context.Context, p *TrainingCoursePrerequisite) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(p).Error
}

func (r *Repository) FindCoursePrerequisiteByID(ctx context.Context, id uuid.UUID) (*TrainingCoursePrerequisite, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var p TrainingCoursePrerequisite
	if err := db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) ListCoursePrerequisites(ctx context.Context, courseID uuid.UUID) ([]TrainingCoursePrerequisite, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var items []TrainingCoursePrerequisite
	if err := db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Order("created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) DeleteCoursePrerequisite(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Delete(&TrainingCoursePrerequisite{}, "id = ?", id).Error
}

// =========================================================================
// Training Mandatories (P1-BE — plan §25)
// =========================================================================

func (r *Repository) CreateMandatory(ctx context.Context, m *TrainingMandatory) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(m).Error
}

func (r *Repository) FindMandatoryByID(ctx context.Context, id uuid.UUID) (*TrainingMandatory, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var m TrainingMandatory
	if err := db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Repository) ListMandatories(ctx context.Context, courseID *uuid.UUID, page, perPage int) ([]TrainingMandatory, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	q := db.WithContext(ctx).Model(&TrainingMandatory{})
	if courseID != nil {
		q = q.Where("course_id = ?", *courseID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []TrainingMandatory
	if err := q.Order("created_at DESC").
		Offset((page - 1) * perPage).Limit(perPage).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) UpdateMandatory(ctx context.Context, m *TrainingMandatory) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(m).Error
}

func (r *Repository) DeleteMandatory(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Delete(&TrainingMandatory{}, "id = ?", id).Error
}

// =========================================================================
// Training Need — repository P1 (plan §17)
// =========================================================================

// DeleteNeed — soft delete training need (plan §17: tabel punya deleted_at).
func (r *Repository) DeleteNeed(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Delete(&TrainingNeed{}, "id = ?", id).Error
}

// =========================================================================
// Training Session Costs (P1-BE — plan §26)
// =========================================================================

func (r *Repository) CreateSessionCost(ctx context.Context, c *TrainingSessionCost) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(c).Error
}

func (r *Repository) FindSessionCostByID(ctx context.Context, id uuid.UUID) (*TrainingSessionCost, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var c TrainingSessionCost
	if err := db.WithContext(ctx).First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) ListSessionCosts(ctx context.Context, sessionID uuid.UUID) ([]TrainingSessionCost, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var items []TrainingSessionCost
	if err := db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) UpdateSessionCost(ctx context.Context, c *TrainingSessionCost) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(c).Error
}

func (r *Repository) DeleteSessionCost(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Delete(&TrainingSessionCost{}, "id = ?", id).Error
}

// =========================================================================
// Training Documents (P1-BE — plan §27)
// =========================================================================

func (r *Repository) CreateDocument(ctx context.Context, d *TrainingDocument) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(d).Error
}

func (r *Repository) FindDocumentByID(ctx context.Context, id uuid.UUID) (*TrainingDocument, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var d TrainingDocument
	if err := db.WithContext(ctx).First(&d, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *Repository) ListDocuments(ctx context.Context, sessionID uuid.UUID) ([]TrainingDocument, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var items []TrainingDocument
	if err := db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) DeleteDocument(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Delete(&TrainingDocument{}, "id = ?", id).Error
}

// =========================================================================
// Evaluation Forms — repository P2 (plan §22)
// =========================================================================

func (r *Repository) CreateEvaluationForm(ctx context.Context, f *TrainingEvaluationForm) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(f).Error
}

func (r *Repository) FindEvaluationFormByID(ctx context.Context, id uuid.UUID) (*TrainingEvaluationForm, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var f TrainingEvaluationForm
	if err := db.WithContext(ctx).Where("id = ?", id).First(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *Repository) FindEvaluationFormBySession(ctx context.Context, sessionID uuid.UUID) (*TrainingEvaluationForm, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var f TrainingEvaluationForm
	if err := db.WithContext(ctx).Where("session_id = ?", sessionID).Order("created_at ASC").First(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *Repository) ListEvaluationForms(ctx context.Context, sessionID *uuid.UUID, page, perPage int) ([]TrainingEvaluationForm, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var forms []TrainingEvaluationForm
	var total int64
	query := db.WithContext(ctx).Model(&TrainingEvaluationForm{})
	if sessionID != nil {
		query = query.Where("session_id = ?", *sessionID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&forms).Error; err != nil {
		return nil, 0, err
	}
	return forms, total, nil
}

func (r *Repository) UpdateEvaluationForm(ctx context.Context, f *TrainingEvaluationForm) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(f).Error
}

func (r *Repository) DeleteEvaluationForm(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Delete(&TrainingEvaluationForm{}, "id = ?", id).Error
}

// =========================================================================
// Evaluation Questions — repository P2 (plan §22)
// =========================================================================

func (r *Repository) CreateEvaluationQuestion(ctx context.Context, q *TrainingEvaluationQuestion) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(q).Error
}

func (r *Repository) FindEvaluationQuestionByID(ctx context.Context, id uuid.UUID) (*TrainingEvaluationQuestion, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var q TrainingEvaluationQuestion
	if err := db.WithContext(ctx).Where("id = ?", id).First(&q).Error; err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *Repository) ListEvaluationQuestions(ctx context.Context, formID uuid.UUID) ([]TrainingEvaluationQuestion, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var qs []TrainingEvaluationQuestion
	if err := db.WithContext(ctx).Where("form_id = ?", formID).Order("sort_order ASC, created_at ASC").Find(&qs).Error; err != nil {
		return nil, err
	}
	return qs, nil
}

func (r *Repository) UpdateEvaluationQuestion(ctx context.Context, q *TrainingEvaluationQuestion) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(q).Error
}

func (r *Repository) DeleteEvaluationQuestion(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Delete(&TrainingEvaluationQuestion{}, "id = ?", id).Error
}

// =========================================================================
// Evaluation Answers — repository P2 (plan §22)
// =========================================================================

// UpsertEvaluationAnswer — satu jawaban per (question_id, participant_id);
// insert bila belum ada, update bila sudah.
func (r *Repository) UpsertEvaluationAnswer(ctx context.Context, a *TrainingEvaluationAnswer) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	var existing TrainingEvaluationAnswer
	err = db.WithContext(ctx).Where("question_id = ? AND participant_id = ?", a.QuestionID, a.ParticipantID).First(&existing).Error
	if err == nil {
		existing.Answer = a.Answer
		return db.WithContext(ctx).Save(&existing).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	return db.WithContext(ctx).Create(a).Error
}

func (r *Repository) ListEvaluationAnswers(ctx context.Context, questionID *uuid.UUID, participantID *uuid.UUID) ([]TrainingEvaluationAnswer, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var answers []TrainingEvaluationAnswer
	query := db.WithContext(ctx).Model(&TrainingEvaluationAnswer{})
	if questionID != nil {
		query = query.Where("question_id = ?", *questionID)
	}
	if participantID != nil {
		query = query.Where("participant_id = ?", *participantID)
	}
	if err := query.Order("created_at ASC").Find(&answers).Error; err != nil {
		return nil, err
	}
	return answers, nil
}

// =========================================================================
// Effectiveness Assessments — repository P2 (plan §23)
// =========================================================================

func (r *Repository) CreateEffectivenessAssessment(ctx context.Context, a *TrainingEffectivenessAssessment) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(a).Error
}

func (r *Repository) FindEffectivenessAssessmentByID(ctx context.Context, id uuid.UUID) (*TrainingEffectivenessAssessment, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var a TrainingEffectivenessAssessment
	if err := db.WithContext(ctx).Where("id = ?", id).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) ListEffectivenessAssessments(ctx context.Context, participantID *uuid.UUID, page, perPage int) ([]TrainingEffectivenessAssessment, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var items []TrainingEffectivenessAssessment
	var total int64
	query := db.WithContext(ctx).Model(&TrainingEffectivenessAssessment{})
	if participantID != nil {
		query = query.Where("participant_id = ?", *participantID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("assessment_date DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) UpdateEffectivenessAssessment(ctx context.Context, a *TrainingEffectivenessAssessment) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(a).Error
}

func (r *Repository) DeleteEffectivenessAssessment(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Delete(&TrainingEffectivenessAssessment{}, "id = ?", id).Error
}

// =========================================================================
// Certifications — repository P2 (plan §24)
// =========================================================================

func (r *Repository) CreateCertification(ctx context.Context, c *TrainingCertification) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(c).Error
}

// NextCertificationCode menghasilkan kode sertifikasi otomatis dengan pola CERT-{NNN}
// (mis. CERT-001). Sekuens dihitung dari kode sertifikasi yang sudah ada sehingga
// tetap unik pada unique index.
func (r *Repository) NextCertificationCode(ctx context.Context) (string, error) {
	db, err := r.db(ctx)
	if err != nil {
		return "", err
	}

	var codes []string
	if err := db.WithContext(ctx).Model(&TrainingCertification{}).
		Where("code LIKE ?", "CERT-%").
		Pluck("code", &codes).Error; err != nil {
		return "", err
	}

	maxSeq := 0
	for _, c := range codes {
		n, err := strconv.Atoi(strings.TrimPrefix(c, "CERT-"))
		if err == nil && n > maxSeq {
			maxSeq = n
		}
	}
	return fmt.Sprintf("CERT-%03d", maxSeq+1), nil
}

func (r *Repository) FindCertificationByID(ctx context.Context, id uuid.UUID) (*TrainingCertification, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var c TrainingCertification
	if err := db.WithContext(ctx).Where("id = ?", id).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) ListCertifications(ctx context.Context, isActive *bool, page, perPage int) ([]TrainingCertification, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var items []TrainingCertification
	var total int64
	query := db.WithContext(ctx).Model(&TrainingCertification{})
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) UpdateCertification(ctx context.Context, c *TrainingCertification) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(c).Error
}

func (r *Repository) DeleteCertification(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Delete(&TrainingCertification{}, "id = ?", id).Error
}

// =========================================================================
// Certificates — repository P2 (plan §24)
// =========================================================================

func (r *Repository) FindCertificateByParticipant(ctx context.Context, participantID uuid.UUID) (*TrainingCertificate, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var c TrainingCertificate
	if err := db.WithContext(ctx).Where("participant_id = ?", participantID).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// round2 membulatkan float ke 2 desimal (helper report).
func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

// =========================================================================
// Reports & History — repository P2 (plan §38)
// =========================================================================

// HistoryByEmployee — riwayat training per employee (dengan nama course & session).
func (r *Repository) HistoryByEmployee(ctx context.Context, employeeID uuid.UUID) ([]TrainingHistoryResponse, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []TrainingHistoryResponse
	query := `
		SELECT tp.id AS participant_id, tp.employee_id, tp.session_id,
		       ts.course_id, c.name AS course_name,
		       ts.session_code, ts.start_date, ts.end_date,
		       tp.attendance_status, tp.score, tp.completion_status,
		       COALESCE(tp.completion_date, '') AS completion_date,
		       COALESCE(tc.certificate_no, '') AS certificate_no,
		       COALESCE(tc.id, '') AS certificate_id
		FROM training_participants tp
		JOIN training_sessions ts ON ts.id = tp.session_id
		JOIN training_courses c ON c.id = ts.course_id
		LEFT JOIN training_certificates tc ON tc.participant_id = tp.id
		WHERE tp.employee_id = ?
		ORDER BY ts.start_date DESC`
	if err := db.WithContext(ctx).Raw(query, employeeID).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// ParticipationReport — semua partisipasi + status.
func (r *Repository) ParticipationReport(ctx context.Context, sessionStatus *string) ([]ParticipationReportRow, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []ParticipationReportRow
	query := `
		SELECT tp.employee_id, e.name AS employee_name,
		       COALESCE(o.nomenclature, '') AS organization_name,
		       ts.course_id, c.name AS course_name,
		       ts.session_code, ts.status AS session_status,
		       tp.attendance_status, tp.score, tp.completion_status
		FROM training_participants tp
		JOIN training_sessions ts ON ts.id = tp.session_id
		JOIN training_courses c ON c.id = ts.course_id
		LEFT JOIN employees e ON e.id = tp.employee_id
		LEFT JOIN employments em ON em.employee_id = tp.employee_id AND em.effective_end_date IS NULL
		LEFT JOIN organizations o ON o.id = em.organization_id`
	args := []interface{}{}
	if sessionStatus != nil && *sessionStatus != "" {
		query += ` WHERE ts.status = ?`
		args = append(args, *sessionStatus)
	}
	query += ` ORDER BY ts.start_date DESC`
	if err := db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// CostReport — biaya per session + cost per participant.
func (r *Repository) CostReport(ctx context.Context) ([]CostReportRow, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []CostReportRow
	query := `
		SELECT ts.id AS session_id, ts.session_code, ts.course_id, c.name AS course_name,
		       COALESCE(tp.name, '') AS provider_name,
		       COALESCE((SELECT SUM(amount) FROM training_session_costs sc WHERE sc.session_id = ts.id AND sc.deleted_at IS NULL), 0) AS total_cost,
		       (SELECT COUNT(*) FROM training_participants tp2 WHERE tp2.session_id = ts.id AND tp2.deleted_at IS NULL) AS participant_count
		FROM training_sessions ts
		JOIN training_courses c ON c.id = ts.course_id
		LEFT JOIN training_providers tp ON tp.id = ts.provider_id
		WHERE ts.deleted_at IS NULL
		ORDER BY ts.start_date DESC`
	if err := db.WithContext(ctx).Raw(query).Scan(&rows).Error; err != nil {
		return nil, err
	}
	// cost per participant dihitung di service (hindari pembagian 0).
	return rows, nil
}

// ComplianceReport — mandatory training compliance per employee.
// Catatan (deliberate simplification, plan §38): scoping target mandatory
// (organization_id/position_id/employment_status_id) belum di-enforce di query —
// setiap employee aktif dipasangkan ke semua mandatory aktif. Disempurnakan di
// iterasi berikutnya bila kebutuhan compliance by-target muncul.
func (r *Repository) ComplianceReport(ctx context.Context) ([]ComplianceReportRow, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []ComplianceReportRow
	query := `
		SELECT e.id AS employee_id, e.name AS employee_name,
		       COALESCE(o.nomenclature, '') AS organization_name,
		       c.id AS course_id, c.name AS course_name,
		       '' AS due_date,
		       COALESCE(tp.completion_status, 'NOT_STARTED') AS completion_status,
		       CASE WHEN tp.id IS NULL THEN 'NOT_COMPLETED' ELSE 'COMPLETED' END AS status
		FROM training_mandatories tm
		JOIN training_courses c ON c.id = tm.course_id
		CROSS JOIN employees e
		LEFT JOIN employments em ON em.employee_id = e.id AND em.effective_end_date IS NULL
		LEFT JOIN organizations o ON o.id = em.organization_id
		LEFT JOIN training_participants tp ON tp.employee_id = e.id AND tp.session_id IN (
			SELECT id FROM training_sessions WHERE course_id = c.id
		) AND tp.completion_status = 'COMPLETED'
		WHERE tm.is_active = TRUE AND tm.deleted_at IS NULL
		ORDER BY o.nomenclature, e.name`
	if err := db.WithContext(ctx).Raw(query).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// DashboardReport — ringkasan analitik training.
func (r *Repository) DashboardReport(ctx context.Context) (*DashboardReport, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	out := &DashboardReport{}
	if err := db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM training_courses WHERE deleted_at IS NULL`).Scan(&out.TotalCourses).Error; err != nil {
		return nil, err
	}
	if err := db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM training_sessions WHERE deleted_at IS NULL`).Scan(&out.TotalSessions).Error; err != nil {
		return nil, err
	}
	if err := db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM training_participants WHERE deleted_at IS NULL`).Scan(&out.TotalParticipants).Error; err != nil {
		return nil, err
	}
	if err := db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM training_providers WHERE deleted_at IS NULL`).Scan(&out.TotalProviders).Error; err != nil {
		return nil, err
	}
	if err := db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM training_requests WHERE deleted_at IS NULL`).Scan(&out.TotalRequests).Error; err != nil {
		return nil, err
	}
	if err := db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM training_requests WHERE status = 'APPROVED' AND deleted_at IS NULL`).Scan(&out.ApprovedRequests).Error; err != nil {
		return nil, err
	}
	if err := db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM training_requests WHERE status = 'PENDING_APPROVAL' AND deleted_at IS NULL`).Scan(&out.PendingRequests).Error; err != nil {
		return nil, err
	}
	// Completion rate: peserta COMPLETED / total peserta.
	if out.TotalParticipants > 0 {
		var completed int64
		if err := db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM training_participants WHERE completion_status = 'COMPLETED' AND deleted_at IS NULL`).Scan(&completed).Error; err != nil {
			return nil, err
		}
		out.CompletionRate = round2(float64(completed) / float64(out.TotalParticipants) * 100)
	}
	// Pass rate: peserta passed / total peserta.
	var passed int64
	if err := db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM training_participants WHERE passed = TRUE AND deleted_at IS NULL`).Scan(&passed).Error; err != nil {
		return nil, err
	}
	if out.TotalParticipants > 0 {
		out.PassRate = round2(float64(passed) / float64(out.TotalParticipants) * 100)
	}
	if err := db.WithContext(ctx).Raw(`SELECT COALESCE(SUM(amount),0) FROM training_session_costs WHERE deleted_at IS NULL`).Scan(&out.TotalTrainingCost).Error; err != nil {
		return nil, err
	}
	if err := db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM training_certificates WHERE deleted_at IS NULL`).Scan(&out.CertificatesIssued).Error; err != nil {
		return nil, err
	}
	return out, nil
}
