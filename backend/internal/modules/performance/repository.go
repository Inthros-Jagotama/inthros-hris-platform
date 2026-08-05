package performance

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	dbFunc func(ctx context.Context) (*gorm.DB, error)
}

func NewRepository(dbFunc func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbFunc: dbFunc}
}

func (r *Repository) db(ctx context.Context) (*gorm.DB, error) {
	return r.dbFunc(ctx)
}

// =========================================================================
// Performance Periods
// =========================================================================

func (r *Repository) CreatePerformancePeriod(ctx context.Context, p *PerformancePeriod) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(p).Error
}

func (r *Repository) FindPerformancePeriodByID(ctx context.Context, id uuid.UUID) (*PerformancePeriod, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var p PerformancePeriod
	if err := db.WithContext(ctx).First(&p, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("performance period not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *Repository) ListPerformancePeriods(ctx context.Context, page, perPage int) ([]PerformancePeriod, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []PerformancePeriod
	var total int64
	query := db.WithContext(ctx).Model(&PerformancePeriod{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("year DESC, created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdatePerformancePeriod(ctx context.Context, p *PerformancePeriod) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(p).Error
}

func (r *Repository) DeletePerformancePeriod(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&PerformancePeriod{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("performance period not found")
	}
	return result.Error
}

// =========================================================================
// Performance Perspectives
// =========================================================================

func (r *Repository) CreatePerformancePerspective(ctx context.Context, p *PerformancePerspective) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(p).Error
}

func (r *Repository) FindPerformancePerspectiveByID(ctx context.Context, id uuid.UUID) (*PerformancePerspective, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var p PerformancePerspective
	if err := db.WithContext(ctx).First(&p, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("performance perspective not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *Repository) ListPerformancePerspectives(ctx context.Context, page, perPage int) ([]PerformancePerspective, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []PerformancePerspective
	var total int64
	query := db.WithContext(ctx).Model(&PerformancePerspective{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("sort_order ASC, name ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdatePerformancePerspective(ctx context.Context, p *PerformancePerspective) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(p).Error
}

func (r *Repository) DeletePerformancePerspective(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&PerformancePerspective{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("performance perspective not found")
	}
	return result.Error
}

// =========================================================================
// Performance Templates
// =========================================================================

func (r *Repository) CreatePerformanceTemplate(ctx context.Context, t *PerformanceTemplate) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(t).Error
}

func (r *Repository) FindPerformanceTemplateByID(ctx context.Context, id uuid.UUID) (*PerformanceTemplate, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var t PerformanceTemplate
	if err := db.WithContext(ctx).First(&t, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("performance template not found")
		}
		return nil, err
	}
	return &t, nil
}

func (r *Repository) ListPerformanceTemplates(ctx context.Context, orgID *uuid.UUID, page, perPage int) ([]PerformanceTemplate, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []PerformanceTemplate
	var total int64
	query := db.WithContext(ctx).Model(&PerformanceTemplate{})
	if orgID != nil {
		query = query.Where("organization_id = ?", *orgID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdatePerformanceTemplate(ctx context.Context, t *PerformanceTemplate) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(t).Error
}

func (r *Repository) DeletePerformanceTemplate(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&PerformanceTemplate{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("performance template not found")
	}
	return result.Error
}

// =========================================================================
// Performance Indicators
// =========================================================================

func (r *Repository) CreatePerformanceIndicator(ctx context.Context, i *PerformanceIndicator) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(i).Error
}

func (r *Repository) FindPerformanceIndicatorByID(ctx context.Context, id uuid.UUID) (*PerformanceIndicator, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var i PerformanceIndicator
	if err := db.WithContext(ctx).First(&i, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("performance indicator not found")
		}
		return nil, err
	}
	return &i, nil
}

func (r *Repository) ListPerformanceIndicators(ctx context.Context, templateID uuid.UUID, page, perPage int) ([]PerformanceIndicator, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []PerformanceIndicator
	var total int64
	query := db.WithContext(ctx).Model(&PerformanceIndicator{}).
		Where("performance_template_id = ?", templateID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("sort_order ASC, created_at ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdatePerformanceIndicator(ctx context.Context, i *PerformanceIndicator) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(i).Error
}

func (r *Repository) DeletePerformanceIndicator(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&PerformanceIndicator{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("performance indicator not found")
	}
	return result.Error
}

// =========================================================================
// Performance Evaluations
// =========================================================================

func (r *Repository) CreatePerformanceEvaluation(ctx context.Context, e *PerformanceEvaluation) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(e).Error
}

func (r *Repository) FindPerformanceEvaluationByID(ctx context.Context, id uuid.UUID) (*PerformanceEvaluation, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var e PerformanceEvaluation
	if err := db.WithContext(ctx).Preload("Details").First(&e, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("performance evaluation not found")
		}
		return nil, err
	}
	return &e, nil
}

func (r *Repository) ListPerformanceEvaluations(ctx context.Context, employeeID, periodID *uuid.UUID, status *string, page, perPage int) ([]PerformanceEvaluation, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []PerformanceEvaluation
	var total int64
	query := db.WithContext(ctx).Model(&PerformanceEvaluation{})
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	if periodID != nil {
		query = query.Where("period_id = ?", *periodID)
	}
	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdatePerformanceEvaluation(ctx context.Context, e *PerformanceEvaluation) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(e).Error
}

func (r *Repository) DeletePerformanceEvaluation(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&PerformanceEvaluation{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("performance evaluation not found")
	}
	return result.Error
}

// =========================================================================
// Performance Evaluation Details
// =========================================================================

func (r *Repository) CreateEvaluationDetail(ctx context.Context, d *PerformanceEvaluationDetail) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(d).Error
}

func (r *Repository) FindEvaluationDetailByID(ctx context.Context, id uuid.UUID) (*PerformanceEvaluationDetail, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var d PerformanceEvaluationDetail
	if err := db.WithContext(ctx).First(&d, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("evaluation detail not found")
		}
		return nil, err
	}
	return &d, nil
}

func (r *Repository) ListEvaluationDetails(ctx context.Context, evalID uuid.UUID) ([]PerformanceEvaluationDetail, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []PerformanceEvaluationDetail
	if err := db.WithContext(ctx).Where("performance_evaluation_id = ?", evalID).
		Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateEvaluationDetail(ctx context.Context, d *PerformanceEvaluationDetail) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(d).Error
}

func (r *Repository) DeleteEvaluationDetail(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&PerformanceEvaluationDetail{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("evaluation detail not found")
	}
	return result.Error
}

// =========================================================================
// Performance Targets
// =========================================================================

func (r *Repository) CreatePerformanceTarget(ctx context.Context, t *PerformanceTarget) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(t).Error
}

func (r *Repository) FindPerformanceTargetByID(ctx context.Context, id uuid.UUID) (*PerformanceTarget, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var t PerformanceTarget
	if err := db.WithContext(ctx).First(&t, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("performance target not found")
		}
		return nil, err
	}
	return &t, nil
}

func (r *Repository) ListPerformanceTargets(ctx context.Context, evalID uuid.UUID) ([]PerformanceTarget, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []PerformanceTarget
	var total int64
	query := db.WithContext(ctx).Model(&PerformanceTarget{}).
		Where("performance_evaluation_id = ?", evalID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdatePerformanceTarget(ctx context.Context, t *PerformanceTarget) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(t).Error
}

func (r *Repository) DeletePerformanceTarget(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&PerformanceTarget{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("performance target not found")
	}
	return result.Error
}

// =========================================================================
// Aggregation
// =========================================================================

func (r *Repository) UpdateEvaluationFinalScore(ctx context.Context, evalID uuid.UUID) (float64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return 0, err
	}
	var totalScore float64
	if err := db.WithContext(ctx).Model(&PerformanceEvaluationDetail{}).
		Where("performance_evaluation_id = ?", evalID).
		Select("COALESCE(SUM(score), 0)").Scan(&totalScore).Error; err != nil {
		return 0, err
	}
	return totalScore, nil
}

// =========================================================================
// Performance Progress
// =========================================================================

func (r *Repository) CreatePerformanceProgress(ctx context.Context, p *PerformanceProgress) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(p).Error
}

func (r *Repository) FindPerformanceProgressByID(ctx context.Context, id uuid.UUID) (*PerformanceProgress, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var p PerformanceProgress
	if err := db.WithContext(ctx).First(&p, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("performance progress not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *Repository) ListPerformanceProgressByDetailID(ctx context.Context, detailID uuid.UUID) ([]PerformanceProgress, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []PerformanceProgress
	if err := db.WithContext(ctx).Where("evaluation_detail_id = ?", detailID).
		Order("progress_date DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdatePerformanceProgress(ctx context.Context, p *PerformanceProgress) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(p).Error
}

func (r *Repository) DeletePerformanceProgress(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&PerformanceProgress{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("performance progress not found")
	}
	return result.Error
}

// =========================================================================
// Performance Comments
// =========================================================================

func (r *Repository) CreatePerformanceComment(ctx context.Context, c *PerformanceComment) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(c).Error
}

func (r *Repository) FindPerformanceCommentByID(ctx context.Context, id uuid.UUID) (*PerformanceComment, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var c PerformanceComment
	if err := db.WithContext(ctx).First(&c, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("performance comment not found")
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) ListPerformanceCommentsByEvaluationID(ctx context.Context, evalID uuid.UUID) ([]PerformanceComment, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []PerformanceComment
	if err := db.WithContext(ctx).Where("evaluation_id = ?", evalID).
		Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdatePerformanceComment(ctx context.Context, c *PerformanceComment) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(c).Error
}

func (r *Repository) DeletePerformanceComment(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&PerformanceComment{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("performance comment not found")
	}
	return result.Error
}

// =========================================================================
// Performance Attachments
// =========================================================================

func (r *Repository) CreatePerformanceAttachment(ctx context.Context, a *PerformanceAttachment) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(a).Error
}

func (r *Repository) FindPerformanceAttachmentByID(ctx context.Context, id uuid.UUID) (*PerformanceAttachment, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var a PerformanceAttachment
	if err := db.WithContext(ctx).First(&a, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("performance attachment not found")
		}
		return nil, err
	}
	return &a, nil
}

func (r *Repository) ListPerformanceAttachmentsByDetailID(ctx context.Context, detailID uuid.UUID) ([]PerformanceAttachment, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []PerformanceAttachment
	if err := db.WithContext(ctx).Where("evaluation_detail_id = ?", detailID).
		Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdatePerformanceAttachment(ctx context.Context, a *PerformanceAttachment) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(a).Error
}

func (r *Repository) DeletePerformanceAttachment(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&PerformanceAttachment{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("performance attachment not found")
	}
	return result.Error
}

// =========================================================================
// Performance Ratings
// =========================================================================

func (r *Repository) CreatePerformanceRating(ctx context.Context, rt *PerformanceRating) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(rt).Error
}

func (r *Repository) FindPerformanceRatingByID(ctx context.Context, id uuid.UUID) (*PerformanceRating, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var rt PerformanceRating
	if err := db.WithContext(ctx).First(&rt, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("performance rating not found")
		}
		return nil, err
	}
	return &rt, nil
}

func (r *Repository) ListPerformanceRatings(ctx context.Context, page, perPage int) ([]PerformanceRating, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []PerformanceRating
	var total int64
	query := db.WithContext(ctx).Model(&PerformanceRating{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("sort_order ASC, name ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdatePerformanceRating(ctx context.Context, rt *PerformanceRating) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(rt).Error
}

func (r *Repository) DeletePerformanceRating(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&PerformanceRating{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("performance rating not found")
	}
	return result.Error
}

// =========================================================================
// Performance Indicator Formulas
// =========================================================================

func (r *Repository) CreatePerformanceIndicatorFormula(ctx context.Context, f *PerformanceIndicatorFormula) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(f).Error
}

func (r *Repository) FindPerformanceIndicatorFormulaByID(ctx context.Context, id uuid.UUID) (*PerformanceIndicatorFormula, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var f PerformanceIndicatorFormula
	if err := db.WithContext(ctx).First(&f, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("performance indicator formula not found")
		}
		return nil, err
	}
	return &f, nil
}

func (r *Repository) ListPerformanceIndicatorFormulas(ctx context.Context, page, perPage int) ([]PerformanceIndicatorFormula, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []PerformanceIndicatorFormula
	var total int64
	query := db.WithContext(ctx).Model(&PerformanceIndicatorFormula{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("sort_order ASC, name ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdatePerformanceIndicatorFormula(ctx context.Context, f *PerformanceIndicatorFormula) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(f).Error
}

func (r *Repository) DeletePerformanceIndicatorFormula(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&PerformanceIndicatorFormula{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("performance indicator formula not found")
	}
	return result.Error
}

// =========================================================================
// Performance Logs
// =========================================================================

func (r *Repository) CreatePerformanceLog(ctx context.Context, l *PerformanceLog) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(l).Error
}

func (r *Repository) FindPerformanceLogByID(ctx context.Context, id uuid.UUID) (*PerformanceLog, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var l PerformanceLog
	if err := db.WithContext(ctx).First(&l, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("performance log not found")
		}
		return nil, err
	}
	return &l, nil
}

func (r *Repository) ListPerformanceLogs(ctx context.Context, page, perPage int) ([]PerformanceLog, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []PerformanceLog
	var total int64
	query := db.WithContext(ctx).Model(&PerformanceLog{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) ListPerformanceLogsByEvaluationID(ctx context.Context, evalID uuid.UUID, page, perPage int) ([]PerformanceLog, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []PerformanceLog
	var total int64
	query := db.WithContext(ctx).Model(&PerformanceLog{}).Where("evaluation_id = ?", evalID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// =========================================================================
// Phase 3 - Business Process Repository Methods
// =========================================================================

// ListIndicatorsByTemplateID returns all indicators for a template (for snapshot)
func (r *Repository) ListIndicatorsByTemplateID(ctx context.Context, templateID uuid.UUID) ([]PerformanceIndicator, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []PerformanceIndicator
	if err := db.WithContext(ctx).Where("performance_template_id = ?", templateID).
		Order("sort_order ASC, created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// CreateEvaluationWithDetails creates evaluation and its details in a transaction
func (r *Repository) CreateEvaluationWithDetails(ctx context.Context, eval *PerformanceEvaluation, details []PerformanceEvaluationDetail) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(eval).Error; err != nil {
			return err
		}
		for i := range details {
			details[i].PerformanceEvaluationID = eval.ID
			if err := tx.Create(&details[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// FindRatingByScore finds the rating that matches the given score
func (r *Repository) FindRatingByScore(ctx context.Context, score float64) (*PerformanceRating, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var rating PerformanceRating
	if err := db.WithContext(ctx).
		Where("min_score <= ? AND max_score >= ?", score, score).
		Order("sort_order ASC").
		First(&rating).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No matching rating found
		}
		return nil, err
	}
	return &rating, nil
}

// UpdateEvaluationWithRating updates evaluation final score and rating
func (r *Repository) UpdateEvaluationWithRating(ctx context.Context, evalID uuid.UUID, finalScore float64, ratingID *uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	updates := map[string]interface{}{
		"final_score": finalScore,
		"rating_id":   ratingID,
	}
	return db.WithContext(ctx).Model(&PerformanceEvaluation{}).Where("id = ?", evalID).Updates(updates).Error
}

// BulkUpdateEvaluationDetails updates multiple evaluation details
func (r *Repository) BulkUpdateEvaluationDetails(ctx context.Context, details []PerformanceEvaluationDetail) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, d := range details {
			if err := tx.Save(&d).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetEvaluationProgressSummary returns progress statistics for an evaluation
func (r *Repository) GetEvaluationProgressSummary(ctx context.Context, evalID uuid.UUID) (total, completed, inProgress, notStarted int, avgAchievement float64, err error) {
	db, err := r.db(ctx)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	var details []PerformanceEvaluationDetail
	if err := db.WithContext(ctx).Where("performance_evaluation_id = ?", evalID).Find(&details).Error; err != nil {
		return 0, 0, 0, 0, 0, err
	}

	total = len(details)
	var totalAchievement float64
	for _, d := range details {
		if d.Actual > 0 && d.Achievement > 0 {
			completed++
			totalAchievement += d.Achievement
		} else if d.Actual > 0 {
			inProgress++
			totalAchievement += d.Achievement
		} else {
			notStarted++
		}
	}

	if completed+inProgress > 0 {
		avgAchievement = totalAchievement / float64(completed+inProgress)
	}

	return total, completed, inProgress, notStarted, avgAchievement, nil
}

// GetLatestProgressByDetailID returns the latest progress entry for a detail
func (r *Repository) GetLatestProgressByDetailID(ctx context.Context, detailID uuid.UUID) (*PerformanceProgress, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var p PerformanceProgress
	if err := db.WithContext(ctx).Where("evaluation_detail_id = ?", detailID).
		Order("progress_date DESC, created_at DESC").
		First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// =========================================================================
// Phase 4 - Dashboard Repository Methods
// =========================================================================

// GetActivePeriod returns the currently active performance period
func (r *Repository) GetActivePeriod(ctx context.Context) (*PerformancePeriod, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var p PerformancePeriod
	if err := db.WithContext(ctx).Where("status = ?", "active").
		Order("year DESC, created_at DESC").
		First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// GetEmployeeEvaluation returns evaluation for employee in a period
func (r *Repository) GetEmployeeEvaluation(ctx context.Context, employeeID, periodID uuid.UUID) (*PerformanceEvaluation, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var e PerformanceEvaluation
	if err := db.WithContext(ctx).
		Where("employee_id = ? AND period_id = ?", employeeID, periodID).
		First(&e).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

// GetRecentLogsByEvaluationID returns recent activity logs for an evaluation
func (r *Repository) GetRecentLogsByEvaluationID(ctx context.Context, evalID uuid.UUID, limit int) ([]PerformanceLog, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var logs []PerformanceLog
	if err := db.WithContext(ctx).Where("evaluation_id = ?", evalID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// ListEvaluationsByOrganization returns all evaluations for an organization in a period
func (r *Repository) ListEvaluationsByOrganization(ctx context.Context, orgID, periodID uuid.UUID) ([]PerformanceEvaluation, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []PerformanceEvaluation
	query := db.WithContext(ctx).Where("organization_id = ?", orgID)
	if periodID != uuid.Nil {
		query = query.Where("period_id = ?", periodID)
	}
	if err := query.Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// ListEvaluationsBySupervisor returns all evaluations supervised by a manager
func (r *Repository) ListEvaluationsBySupervisor(ctx context.Context, supervisorID, periodID uuid.UUID) ([]PerformanceEvaluation, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []PerformanceEvaluation
	query := db.WithContext(ctx).Where("supervisor_id = ?", supervisorID)
	if periodID != uuid.Nil {
		query = query.Where("period_id = ?", periodID)
	}
	if err := query.Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetPendingReviews returns evaluations pending manager review (SUBMITTED status)
func (r *Repository) GetPendingReviews(ctx context.Context, supervisorID, periodID uuid.UUID) ([]PerformanceEvaluation, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []PerformanceEvaluation
	query := db.WithContext(ctx).Where("supervisor_id = ? AND status = ?", supervisorID, "SUBMITTED")
	if periodID != uuid.Nil {
		query = query.Where("period_id = ?", periodID)
	}
	if err := query.Order("submitted_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// CountEvaluationsByStatus counts evaluations by status for a period
func (r *Repository) CountEvaluationsByStatus(ctx context.Context, periodID uuid.UUID) (draft, submitted, approved, completed int64, err error) {
	db, err := r.db(ctx)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	query := db.WithContext(ctx).Model(&PerformanceEvaluation{})
	if periodID != uuid.Nil {
		query = query.Where("period_id = ?", periodID)
	}

	query.Where("status = ?", "DRAFT").Count(&draft)
	query.Where("status = ?", "SUBMITTED").Count(&submitted)
	query.Where("status = ?", "APPROVED").Count(&approved)
	query.Where("status = ?", "COMPLETED").Count(&completed)

	return draft, submitted, approved, completed, nil
}

// GetAverageScoreByPeriod returns average final score for a period
func (r *Repository) GetAverageScoreByPeriod(ctx context.Context, periodID uuid.UUID) (avgScore float64, avgAchievement float64, err error) {
	db, err := r.db(ctx)
	if err != nil {
		return 0, 0, err
	}

	type Result struct {
		AvgScore       float64
		AvgAchievement float64
	}
	var result Result

	query := db.WithContext(ctx).Model(&PerformanceEvaluation{}).
		Select("COALESCE(AVG(final_score), 0) as avg_score")
	if periodID != uuid.Nil {
		query = query.Where("period_id = ? AND status IN ?", periodID, []string{"APPROVED", "COMPLETED"})
	} else {
		query = query.Where("status IN ?", []string{"APPROVED", "COMPLETED"})
	}
	if err := query.Scan(&result).Error; err != nil {
		return 0, 0, err
	}

	// Calculate average achievement from details
	detailQuery := db.WithContext(ctx).Model(&PerformanceEvaluationDetail{}).
		Select("COALESCE(AVG(achievement), 0) as avg_achievement").
		Joins("JOIN performance_evaluations e ON e.id = performance_evaluation_details.performance_evaluation_id")
	if periodID != uuid.Nil {
		detailQuery = detailQuery.Where("e.period_id = ? AND e.status IN ?", periodID, []string{"APPROVED", "COMPLETED"})
	} else {
		detailQuery = detailQuery.Where("e.status IN ?", []string{"APPROVED", "COMPLETED"})
	}
	if err := detailQuery.Scan(&result.AvgAchievement).Error; err != nil {
		return result.AvgScore, 0, err
	}

	return result.AvgScore, result.AvgAchievement, nil
}

// GetRatingDistribution returns count of evaluations per rating for a period
func (r *Repository) GetRatingDistribution(ctx context.Context, periodID uuid.UUID) ([]struct {
	RatingID uuid.UUID
	Count    int64
}, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}

	type RatingCount struct {
		RatingID uuid.UUID
		Count    int64
	}
	var results []RatingCount

	query := db.WithContext(ctx).Model(&PerformanceEvaluation{}).
		Select("rating_id, COUNT(*) as count").
		Where("rating_id IS NOT NULL AND status IN ?", []string{"APPROVED", "COMPLETED"}).
		Group("rating_id")
	if periodID != uuid.Nil {
		query = query.Where("period_id = ?", periodID)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Convert to expected return type
	var output []struct {
		RatingID uuid.UUID
		Count    int64
	}
	for _, r := range results {
		output = append(output, struct {
			RatingID uuid.UUID
			Count    int64
		}{r.RatingID, r.Count})
	}
	return output, nil
}

// GetOrganizationStats returns performance stats grouped by organization
func (r *Repository) GetOrganizationStats(ctx context.Context, periodID uuid.UUID) ([]struct {
	OrganizationID uuid.UUID
	TotalCount     int64
	CompletedCount int64
	AvgScore       float64
}, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}

	type OrgStats struct {
		OrganizationID uuid.UUID
		TotalCount     int64
		CompletedCount int64
		AvgScore       float64
	}
	var results []OrgStats

	query := db.WithContext(ctx).Model(&PerformanceEvaluation{}).
		Select(`organization_id,
			COUNT(*) as total_count,
			SUM(CASE WHEN status IN ('APPROVED', 'COMPLETED') THEN 1 ELSE 0 END) as completed_count,
			COALESCE(AVG(CASE WHEN status IN ('APPROVED', 'COMPLETED') THEN final_score END), 0) as avg_score`).
		Group("organization_id")
	if periodID != uuid.Nil {
		query = query.Where("period_id = ?", periodID)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	var output []struct {
		OrganizationID uuid.UUID
		TotalCount     int64
		CompletedCount int64
		AvgScore       float64
	}
	for _, r := range results {
		output = append(output, struct {
			OrganizationID uuid.UUID
			TotalCount     int64
			CompletedCount int64
			AvgScore       float64
		}{r.OrganizationID, r.TotalCount, r.CompletedCount, r.AvgScore})
	}
	return output, nil
}

// GetTopPerformers returns top N performers for a period
func (r *Repository) GetTopPerformers(ctx context.Context, periodID uuid.UUID, limit int) ([]PerformanceEvaluation, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}

	var list []PerformanceEvaluation
	query := db.WithContext(ctx).
		Where("status IN ? AND final_score > 0", []string{"APPROVED", "COMPLETED"})
	if periodID != uuid.Nil {
		query = query.Where("period_id = ?", periodID)
	}
	if err := query.Order("final_score DESC").Limit(limit).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetBottomPerformers returns bottom N performers for a period
func (r *Repository) GetBottomPerformers(ctx context.Context, periodID uuid.UUID, limit int) ([]PerformanceEvaluation, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}

	var list []PerformanceEvaluation
	query := db.WithContext(ctx).
		Where("status IN ? AND final_score > 0", []string{"APPROVED", "COMPLETED"})
	if periodID != uuid.Nil {
		query = query.Where("period_id = ?", periodID)
	}
	if err := query.Order("final_score ASC").Limit(limit).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetTrendData returns performance trends across periods
func (r *Repository) GetTrendData(ctx context.Context, limit int) ([]struct {
	PeriodID     uuid.UUID
	AvgScore     float64
	TotalCount   int64
	CompletedCount int64
}, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}

	type TrendStats struct {
		PeriodID       uuid.UUID
		AvgScore       float64
		TotalCount     int64
		CompletedCount int64
	}
	var results []TrendStats

	if err := db.WithContext(ctx).Model(&PerformanceEvaluation{}).
		Select(`period_id,
			COALESCE(AVG(CASE WHEN status IN ('APPROVED', 'COMPLETED') THEN final_score END), 0) as avg_score,
			COUNT(*) as total_count,
			SUM(CASE WHEN status IN ('APPROVED', 'COMPLETED') THEN 1 ELSE 0 END) as completed_count`).
		Group("period_id").
		Order("period_id DESC").
		Limit(limit).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	var output []struct {
		PeriodID       uuid.UUID
		AvgScore       float64
		TotalCount     int64
		CompletedCount int64
	}
	for _, r := range results {
		output = append(output, struct {
			PeriodID       uuid.UUID
			AvgScore       float64
			TotalCount     int64
			CompletedCount int64
		}{r.PeriodID, r.AvgScore, r.TotalCount, r.CompletedCount})
	}
	return output, nil
}
