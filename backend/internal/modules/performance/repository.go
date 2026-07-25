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
