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

// =========================================================================
// Rater Assignment
// =========================================================================

// FindRatersByTarget mengambil seluruh rater untuk satu assessment target
// beserta relasi target-nya.
func (r *Repository) FindRatersByTarget(ctx context.Context, targetID uuid.UUID) ([]CompetencyAssessmentRater, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var list []CompetencyAssessmentRater
	if err := db.Preload("Target").
		Where("competency_event_target_id = ?", targetID).
		Order("rater_type ASC, created_at ASC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// FindRatersByEmployee mengambil seluruh assessment yang ditugaskan kepada
// seorang employee sebagai rater (untuk "My Assessment").
func (r *Repository) FindRatersByEmployee(ctx context.Context, employeeID uuid.UUID) ([]CompetencyAssessmentRater, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var list []CompetencyAssessmentRater
	if err := db.Preload("Target").
		Where("rater_employee_id = ?", employeeID).
		Order("created_at DESC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// FindRaterByID mengambil satu rater beserta target-nya.
func (r *Repository) FindRaterByID(ctx context.Context, id uuid.UUID) (*CompetencyAssessmentRater, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rat CompetencyAssessmentRater
	if err := db.Preload("Target").First(&rat, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("assessment rater not found: %w", err)
	}
	return &rat, nil
}

// FindRaterByTargetAndEmployee memastikan rater belum duplikat pada target.
func (r *Repository) FindRaterByTargetAndEmployee(ctx context.Context, targetID, employeeID uuid.UUID) (*CompetencyAssessmentRater, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rat CompetencyAssessmentRater
	err = db.Where("competency_event_target_id = ? AND rater_employee_id = ?", targetID, employeeID).First(&rat).Error
	if err != nil {
		return nil, err
	}
	return &rat, nil
}

func (r *Repository) CreateRater(ctx context.Context, rat *CompetencyAssessmentRater) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(rat).Error
}

func (r *Repository) UpdateRater(ctx context.Context, rat *CompetencyAssessmentRater) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(rat).Error
}

func (r *Repository) DeleteRater(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&CompetencyAssessmentRater{}).Error
}

// DeleteRatersByTarget menghapus seluruh rater pada target (untuk re-assign).
func (r *Repository) DeleteRatersByTarget(ctx context.Context, targetID uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("competency_event_target_id = ?", targetID).Delete(&CompetencyAssessmentRater{}).Error
}

// =========================================================================
// Assessment Response
// =========================================================================

func (r *Repository) FindResponsesByRater(ctx context.Context, raterID uuid.UUID) ([]CompetencyAssessmentResponse, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var list []CompetencyAssessmentResponse
	if err := db.Preload("Indicator").
		Where("rater_id = ?", raterID).
		Order("created_at ASC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// FindResponseByRaterAndIndicator memeriksa response yang sudah ada untuk
// pasangan (rater, indicator) — dipakai upsert.
func (r *Repository) FindResponseByRaterAndIndicator(ctx context.Context, raterID, indicatorID uuid.UUID) (*CompetencyAssessmentResponse, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var resp CompetencyAssessmentResponse
	err = db.Where("rater_id = ? AND indicator_id = ?", raterID, indicatorID).First(&resp).Error
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *Repository) SaveResponse(ctx context.Context, resp *CompetencyAssessmentResponse) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(resp).Error
}

// DeleteResponsesByRater menghapus seluruh response milik rater (reopen).
func (r *Repository) DeleteResponsesByRater(ctx context.Context, raterID uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("rater_id = ?", raterID).Delete(&CompetencyAssessmentResponse{}).Error
}

// =========================================================================
// Employee resolution (tanpa import package employee — hindari circular dep)
// =========================================================================

// GetEmployeeNamesByIDs mengambil nama karyawan untuk sekumpulan employee ID
// via raw query ke tabel employees.
func (r *Repository) GetEmployeeNamesByIDs(ctx context.Context, ids []uuid.UUID) (map[string]string, error) {
	result := make(map[string]string, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	type row struct {
		ID   string
		Name string
	}
	var rows []row
	if err := db.WithContext(ctx).Table("employees").
		Select("id, name").
		Where("id IN ?", ids).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, rrow := range rows {
		result[rrow.ID] = rrow.Name
	}
	return result, nil
}

// FindEmployeeIDByUserID resolve platform user (karyawan yang login) ke
// employee_id via employee_accounts (user_id -> employee_id). Mengembalikan
// nil bila user tidak punya akun employee terkait — pola sama reimbursement.
func (r *Repository) FindEmployeeIDByUserID(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var empIDStrs []string
	err = db.WithContext(ctx).Table("employee_accounts").
		Where("user_id = ?", userID).
		Limit(1).
		Pluck("employee_id", &empIDStrs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to resolve employee id: %w", err)
	}
	if len(empIDStrs) == 0 || empIDStrs[0] == "" {
		return nil, nil
	}
	empID, err := uuid.Parse(empIDStrs[0])
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}
	return &empID, nil
}

// =========================================================================
// Calculation inputs (§14)
// =========================================================================

// FindAllRatersByTarget mengambil seluruh rater + response-nya untuk target
// — input lengkap calculation engine (rater type, weight, response per
// indicator).
func (r *Repository) FindAllRatersByTarget(ctx context.Context, targetID uuid.UUID) ([]CompetencyAssessmentRater, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var list []CompetencyAssessmentRater
	if err := db.Preload("Responses.Indicator").
		Where("competency_event_target_id = ?", targetID).
		Order("rater_type ASC, created_at ASC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// FindTemplateCompetencies mengambil template competencies beserta nama
// competency — input required level & weight per competency (§14).
func (r *Repository) FindTemplateCompetencies(ctx context.Context, templateID uuid.UUID) ([]CompetencyAssessmentTemplateCompetency, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var list []CompetencyAssessmentTemplateCompetency
	if err := db.Preload("Competency").
		Where("template_id = ?", templateID).
		Order("sort_order ASC, created_at ASC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// FindTemplateRaterTypes mengambil konfigurasi weight rater type per template
// (§10) — input weighting calculation engine.
func (r *Repository) FindTemplateRaterTypes(ctx context.Context, templateID uuid.UUID) ([]CompetencyAssessmentTemplateRaterType, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var list []CompetencyAssessmentTemplateRaterType
	if err := db.Where("template_id = ?", templateID).
		Order("created_at ASC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// ListTemplateIndicatorsByCompetency mengambil template indicators + competency
// id dari indicator — memetakan indicator ke competency untuk agregasi.
func (r *Repository) ListTemplateIndicatorsByCompetency(ctx context.Context, templateID uuid.UUID) ([]CompetencyAssessmentTemplateIndicator, error) {
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

// FindScoreByEventAndEmployee mengambil competency score untuk (event,
// employee) — dipakai upsert hasil finalisasi.
func (r *Repository) FindScoreByEventAndEmployee(ctx context.Context, eventID, employeeID uuid.UUID) (*CompetencyScore, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var s CompetencyScore
	err = db.Preload("Details").
		Where("competency_event_id = ? AND employee_id = ?", eventID, employeeID).
		First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// FindTargetByEventAndEmployee mengambil assessment target seorang employee
// pada sebuah event.
func (r *Repository) FindTargetByEventAndEmployee(ctx context.Context, eventID, employeeID uuid.UUID) (*CompetencyEventTarget, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var t CompetencyEventTarget
	if err := db.Where("competency_event_id = ? AND employee_id = ?", eventID, employeeID).First(&t).Error; err != nil {
		return nil, fmt.Errorf("assessment target not found: %w", err)
	}
	return &t, nil
}

// FindTargetsByEmployee mengambil seluruh assessment target seorang employee
// (dipakai pencarian hasil terbaru — finalized dulu).
func (r *Repository) FindTargetsByEmployee(ctx context.Context, employeeID uuid.UUID) ([]CompetencyEventTarget, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var list []CompetencyEventTarget
	if err := db.Where("employee_id = ?", employeeID).
		Order("created_at DESC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// ReplaceScoreDetails menghapus detail lama dan menulis ulang detail skor
// untuk satu competency score (hasil calculation bersifat snapshot).
func (r *Repository) ReplaceScoreDetails(ctx context.Context, scoreID uuid.UUID, details []CompetencyScoreDetail) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("competency_score_id = ?", scoreID).Delete(&CompetencyScoreDetail{}).Error; err != nil {
			return err
		}
		for i := range details {
			details[i].CompetencyScoreID = scoreID
			if err := tx.Create(&details[i]).Error; err != nil {
				return err
		}
		}
		return nil
	})
}
