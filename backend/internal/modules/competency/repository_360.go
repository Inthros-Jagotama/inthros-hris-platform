package competency

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =========================================================================
// Rating Scale
// =========================================================================

func (r *Repository) CreateRatingScaleWithItems(ctx context.Context, scale *CompetencyRatingScale, items []CompetencyRatingScaleItem) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(scale).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].ScaleID = scale.ID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) FindRatingScaleByID(ctx context.Context, id uuid.UUID) (*CompetencyRatingScale, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var s CompetencyRatingScale
	if err := db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC, value ASC")
	}).First(&s, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("rating scale not found: %w", err)
	}
	return &s, nil
}

func (r *Repository) FindAllRatingScales(ctx context.Context, page, perPage int, status string) ([]CompetencyRatingScale, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []CompetencyRatingScale
	var total int64

	query := db.Model(&CompetencyRatingScale{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("name ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateRatingScale(ctx context.Context, scale *CompetencyRatingScale) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(scale).Error
}

func (r *Repository) ReplaceScaleItems(ctx context.Context, scaleID uuid.UUID, items []CompetencyRatingScaleItem) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("scale_id = ?", scaleID).Delete(&CompetencyRatingScaleItem{}).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].ScaleID = scaleID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) DeleteRatingScale(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&CompetencyRatingScale{}).Error
}

// =========================================================================
// Assessment Template
// =========================================================================

func (r *Repository) CreateAssessmentTemplate(ctx context.Context, tpl *CompetencyAssessmentTemplate) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(tpl).Error
}

func (r *Repository) FindAssessmentTemplateByID(ctx context.Context, id uuid.UUID) (*CompetencyAssessmentTemplate, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var t CompetencyAssessmentTemplate
	if err := db.
		Preload("Competencies.Competency").
		Preload("RaterTypes").
		First(&t, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("assessment template not found: %w", err)
	}
	return &t, nil
}

func (r *Repository) FindAllAssessmentTemplates(ctx context.Context, page, perPage int, status string) ([]CompetencyAssessmentTemplate, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []CompetencyAssessmentTemplate
	var total int64

	query := db.Model(&CompetencyAssessmentTemplate{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("name ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateAssessmentTemplate(ctx context.Context, tpl *CompetencyAssessmentTemplate) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(tpl).Error
}

func (r *Repository) ReplaceTemplateCompetencies(ctx context.Context, templateID uuid.UUID, items []CompetencyAssessmentTemplateCompetency) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("template_id = ?", templateID).Delete(&CompetencyAssessmentTemplateCompetency{}).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].TemplateID = templateID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) ReplaceTemplateRaterTypes(ctx context.Context, templateID uuid.UUID, items []CompetencyAssessmentTemplateRaterType) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("template_id = ?", templateID).Delete(&CompetencyAssessmentTemplateRaterType{}).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].TemplateID = templateID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) DeleteAssessmentTemplate(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&CompetencyAssessmentTemplate{}).Error
}

// =========================================================================
// Indicator
// =========================================================================

func (r *Repository) CreateIndicator(ctx context.Context, ind *CompetencyIndicator) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(ind).Error
}

func (r *Repository) FindIndicatorByID(ctx context.Context, id uuid.UUID) (*CompetencyIndicator, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var ind CompetencyIndicator
	if err := db.Preload("Competency").First(&ind, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("competency indicator not found: %w", err)
	}
	return &ind, nil
}

func (r *Repository) FindAllIndicators(ctx context.Context, page, perPage int, competencyID, status string) ([]CompetencyIndicator, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []CompetencyIndicator
	var total int64

	query := db.Model(&CompetencyIndicator{})
	if competencyID != "" {
		if uid, perr := uuid.Parse(competencyID); perr == nil {
			query = query.Where("competency_id = ?", uid)
		}
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Preload("Competency").Offset(offset).Limit(perPage).Order("sort_order ASC, created_at ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateIndicator(ctx context.Context, ind *CompetencyIndicator) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(ind).Error
}

func (r *Repository) DeleteIndicator(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&CompetencyIndicator{}).Error
}

// ReplaceTemplateIndicators mengganti seluruh indicator milik template.
func (r *Repository) ReplaceTemplateIndicators(ctx context.Context, templateID uuid.UUID, items []CompetencyAssessmentTemplateIndicator) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("template_id = ?", templateID).Delete(&CompetencyAssessmentTemplateIndicator{}).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].TemplateID = templateID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ListTemplateIndicators mengambil indicator milik template beserta statement-nya.
func (r *Repository) ListTemplateIndicators(ctx context.Context, templateID uuid.UUID) ([]CompetencyAssessmentTemplateIndicator, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var list []CompetencyAssessmentTemplateIndicator
	if err := db.Preload("Indicator").
		Where("template_id = ?", templateID).
		Order("sort_order ASC, created_at ASC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
