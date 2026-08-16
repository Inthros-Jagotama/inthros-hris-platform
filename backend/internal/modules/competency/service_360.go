package competency

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
)

func calcTotalPages(total int64, perPage int) int {
	tp := int(total) / perPage
	if int(total)%perPage > 0 {
		tp++
	}
	return tp
}

// generateIndicatorCode membuat kode indikator dari statement — UPPER_SNAKE_CASE
// (pola generateUniqueCode) dipotong maksimal 40 karakter agar muat kolom
// varchar(50), dengan suffix numerik bila kode dasar sudah terpakai.
func generateIndicatorCode(statement string, existing map[string]bool) string {
	base := generateUniqueCode(statement, existing)
	if len(base) > 40 {
		base = base[:40]
	}
	code := base
	for i := 2; existing[code]; i++ {
		suffix := fmt.Sprintf("_%d", i)
		cut := 40 - len(suffix)
		if cut < 1 {
			cut = 1
		}
		code = base[:cut] + suffix
	}
	return code
}

// toCodeSet mengubah daftar code menjadi set untuk pengecekan duplikat.
func toCodeSet(codes []string) map[string]bool {
	set := make(map[string]bool, len(codes))
	for _, c := range codes {
		set[c] = true
	}
	return set
}

// generateUniqueCode mengubah nama menjadi kode UPPER_SNAKE_CASE (mis.
// "Skala 5 Poin" -> "SKALA_5_POIN") dan menambahkan suffix numerik bila kode
// dasar sudah dipakai — pola sama attendance.generateUniqueCode.
func generateUniqueCode(name string, existing map[string]bool) string {
	var b strings.Builder
	lastUnderscore := false
	for _, r := range strings.ToUpper(strings.TrimSpace(name)) {
		switch {
		case r >= 'A' && r <= 'Z' || r >= '0' && r <= '9':
			b.WriteRune(r)
			lastUnderscore = false
		default:
			if !lastUnderscore && b.Len() > 0 {
				b.WriteRune('_')
				lastUnderscore = true
			}
		}
	}
	base := strings.TrimSuffix(b.String(), "_")
	if base == "" {
		base = "SCALE"
	}
	code := base
	for i := 2; existing[code]; i++ {
		code = fmt.Sprintf("%s_%d", base, i)
	}
	return code
}

// =========================================================================
// Rating Scale
// =========================================================================

func (s *Service) CreateRatingScale(ctx context.Context, req CreateRatingScaleRequest) (*RatingScaleResponse, error) {
	// Code dibuat otomatis dari Name (UPPER_SNAKE_CASE + suffix bila duplikat)
	// — form pembuatan tidak perlu mengisi kode.
	existingCodes, err := s.repo.FindRatingScaleCodes(ctx)
	if err != nil {
		return nil, err
	}
	scale := &CompetencyRatingScale{
		Name: req.Name,
		Code: generateUniqueCode(req.Name, toCodeSet(existingCodes)),
	}
	if req.Description != nil {
		scale.Description = req.Description
	}
	if req.Status != "" {
		scale.Status = req.Status
	}
	scale.CreatedBy = authctx.GetUserID(ctx)
	scale.UpdatedBy = scale.CreatedBy

	items := make([]CompetencyRatingScaleItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, CompetencyRatingScaleItem{
			Value:     it.Value,
			Label:     it.Label,
			Weight:    it.Weight,
			SortOrder: it.SortOrder,
		})
		if it.Description != nil {
			items[len(items)-1].Description = it.Description
		}
	}

	if err := s.repo.CreateRatingScaleWithItems(ctx, scale, items); err != nil {
		return nil, err
	}
	response := scale.ToResponse()
	return &response, nil
}

func (s *Service) GetRatingScaleByID(ctx context.Context, id string) (*RatingScaleResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid rating scale id: %w", err)
	}
	scale, err := s.repo.FindRatingScaleByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := scale.ToResponse()
	return &response, nil
}

func (s *Service) ListRatingScales(ctx context.Context, page, perPage int, status string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.FindAllRatingScales(ctx, page, perPage, status)
	if err != nil {
		return nil, err
	}
	responses := make([]RatingScaleResponse, 0, len(list))
	for _, scale := range list {
		responses = append(responses, scale.ToResponse())
	}
	totalPages := calcTotalPages(total, perPage)
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateRatingScale(ctx context.Context, id string, req UpdateRatingScaleRequest) (*RatingScaleResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid rating scale id: %w", err)
	}
	scale, err := s.repo.FindRatingScaleByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	scale.UpdatedBy = authctx.GetUserID(ctx)
	if req.Name != nil {
		scale.Name = *req.Name
	}
	if req.Code != nil {
		scale.Code = *req.Code
	}
	if req.Description != nil {
		scale.Description = req.Description
	}
	if req.Status != nil {
		scale.Status = *req.Status
	}
	if err := s.repo.UpdateRatingScale(ctx, scale); err != nil {
		return nil, err
	}
	if req.Items != nil {
		items := make([]CompetencyRatingScaleItem, 0, len(req.Items))
		for _, it := range req.Items {
			item := CompetencyRatingScaleItem{
				Value:     it.Value,
				Label:     it.Label,
				Weight:    it.Weight,
				SortOrder: it.SortOrder,
			}
			if it.Description != nil {
				item.Description = it.Description
			}
			items = append(items, item)
		}
		if err := s.repo.ReplaceScaleItems(ctx, uid, items); err != nil {
			return nil, err
		}
	}
	updated, err := s.repo.FindRatingScaleByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := updated.ToResponse()
	return &response, nil
}

func (s *Service) DeleteRatingScale(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid rating scale id: %w", err)
	}
	return s.repo.DeleteRatingScale(ctx, uid)
}

// =========================================================================
// Assessment Template
// =========================================================================

func (s *Service) CreateAssessmentTemplate(ctx context.Context, req CreateAssessmentTemplateRequest) (*AssessmentTemplateResponse, error) {
	// Code dibuat otomatis dari Name (UPPER_SNAKE_CASE + suffix bila duplikat)
	// — form pembuatan tidak perlu mengisi kode.
	existingCodes, err := s.repo.FindAssessmentTemplateCodes(ctx)
	if err != nil {
		return nil, err
	}
	tpl := &CompetencyAssessmentTemplate{
		Name: req.Name,
		Code: generateUniqueCode(req.Name, toCodeSet(existingCodes)),
	}
	if req.Description != nil {
		tpl.Description = req.Description
	}
	if req.Status != "" {
		tpl.Status = req.Status
	}
	if req.ScaleID != nil && *req.ScaleID != "" {
		if uid, perr := uuid.Parse(*req.ScaleID); perr == nil {
			tpl.ScaleID = &uid
		}
	}
	tpl.CreatedBy = authctx.GetUserID(ctx)
	tpl.UpdatedBy = tpl.CreatedBy

	if err := s.repo.CreateAssessmentTemplate(ctx, tpl); err != nil {
		return nil, err
	}
	if err := s.replaceTemplateChildren(ctx, tpl.ID, req.Competencies, req.RaterTypes, nil); err != nil {
		return nil, err
	}
	loaded, err := s.repo.FindAssessmentTemplateByID(ctx, tpl.ID)
	if err != nil {
		return nil, err
	}
	response := loaded.ToResponse()
	return &response, nil
}

func (s *Service) GetAssessmentTemplateByID(ctx context.Context, id string) (*AssessmentTemplateResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid assessment template id: %w", err)
	}
	tpl, err := s.repo.FindAssessmentTemplateByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := tpl.ToResponse()
	return &response, nil
}

func (s *Service) ListAssessmentTemplates(ctx context.Context, page, perPage int, status string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.FindAllAssessmentTemplates(ctx, page, perPage, status)
	if err != nil {
		return nil, err
	}
	responses := make([]AssessmentTemplateResponse, 0, len(list))
	for _, tpl := range list {
		responses = append(responses, tpl.ToResponse())
	}
	totalPages := calcTotalPages(total, perPage)
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateAssessmentTemplate(ctx context.Context, id string, req UpdateAssessmentTemplateRequest) (*AssessmentTemplateResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid assessment template id: %w", err)
	}
	tpl, err := s.repo.FindAssessmentTemplateByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	tpl.UpdatedBy = authctx.GetUserID(ctx)
	if req.Name != nil {
		tpl.Name = *req.Name
	}
	if req.Code != nil {
		tpl.Code = *req.Code
	}
	if req.Description != nil {
		tpl.Description = req.Description
	}
	if req.Status != nil {
		tpl.Status = *req.Status
	}
	if req.ScaleID != nil {
		if *req.ScaleID == "" {
			tpl.ScaleID = nil
		} else if uid2, perr := uuid.Parse(*req.ScaleID); perr == nil {
			tpl.ScaleID = &uid2
		}
	}
	if err := s.repo.UpdateAssessmentTemplate(ctx, tpl); err != nil {
		return nil, err
	}
	if req.Competencies != nil || req.RaterTypes != nil {
		if err := s.replaceTemplateChildren(ctx, uid, req.Competencies, req.RaterTypes, nil); err != nil {
			return nil, err
		}
	}
	loaded, err := s.repo.FindAssessmentTemplateByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := loaded.ToResponse()
	return &response, nil
}

func (s *Service) DeleteAssessmentTemplate(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid assessment template id: %w", err)
	}
	return s.repo.DeleteAssessmentTemplate(ctx, uid)
}

// replaceTemplateChildren mengganti children template secara idempotent:
// competencies, rater types, dan (opsional) indicators. Setiap slice yang
// nil dibiarkan tidak tersentuh.
func (s *Service) replaceTemplateChildren(ctx context.Context, templateID uuid.UUID, competencies []TemplateCompetencyRequest, raterTypes []TemplateRaterTypeRequest, indicators []TemplateIndicatorRequest) error {
	if competencies != nil {
		items := make([]CompetencyAssessmentTemplateCompetency, 0, len(competencies))
		for _, c := range competencies {
			compUID, err := uuid.Parse(c.CompetencyID)
			if err != nil {
				return fmt.Errorf("invalid competency_id: %w", err)
			}
			item := CompetencyAssessmentTemplateCompetency{
				CompetencyID: compUID,
				Weight:       c.Weight,
				SortOrder:    c.SortOrder,
			}
			if c.RequiredLevel != nil {
				item.RequiredLevel = c.RequiredLevel
			}
			items = append(items, item)
		}
		if err := s.repo.ReplaceTemplateCompetencies(ctx, templateID, items); err != nil {
			return err
		}
	}
	if raterTypes != nil {
		items := make([]CompetencyAssessmentTemplateRaterType, 0, len(raterTypes))
		for _, rt := range raterTypes {
			item := CompetencyAssessmentTemplateRaterType{
				RaterType: rt.RaterType,
				Weight:    rt.Weight,
				MinRater:  rt.MinRater,
				MaxRater:  rt.MaxRater,
				Required:  rt.Required,
				Anonymous: rt.Anonymous,
			}
			items = append(items, item)
		}
		if err := s.repo.ReplaceTemplateRaterTypes(ctx, templateID, items); err != nil {
			return err
		}
	}
	if indicators != nil {
		items := make([]CompetencyAssessmentTemplateIndicator, 0, len(indicators))
		for _, ind := range indicators {
			indUID, err := uuid.Parse(ind.IndicatorID)
			if err != nil {
				return fmt.Errorf("invalid indicator_id: %w", err)
			}
			items = append(items, CompetencyAssessmentTemplateIndicator{
				IndicatorID: indUID,
				Weight:      ind.Weight,
				SortOrder:   ind.SortOrder,
			})
		}
		if err := s.repo.ReplaceTemplateIndicators(ctx, templateID, items); err != nil {
			return err
		}
	}
	return nil
}

// =========================================================================
// Indicator
// =========================================================================

func (s *Service) CreateIndicator(ctx context.Context, req CreateIndicatorRequest) (*IndicatorResponse, error) {
	compUID, err := uuid.Parse(req.CompetencyID)
	if err != nil {
		return nil, fmt.Errorf("invalid competency_id: %w", err)
	}
	ind := &CompetencyIndicator{
		CompetencyID: compUID,
		Statement:    req.Statement,
	}
	// Code dibuat otomatis dari Statement (UPPER_SNAKE_CASE + suffix bila
	// duplikat) — form pembuatan tidak perlu mengisi kode. Kode manual dari
	// request tetap dihormati (backward compatible).
	if req.Code != nil {
		ind.Code = req.Code
	} else {
		existingCodes, err := s.repo.FindAllIndicatorCodes(ctx)
		if err != nil {
			return nil, err
		}
		code := generateIndicatorCode(req.Statement, toCodeSet(existingCodes))
		ind.Code = &code
	}
	if req.Description != nil {
		ind.Description = req.Description
	}
	if req.Status != "" {
		ind.Status = req.Status
	}
	ind.SortOrder = req.SortOrder

	if err := s.repo.CreateIndicator(ctx, ind); err != nil {
		return nil, err
	}
	loaded, err := s.repo.FindIndicatorByID(ctx, ind.ID)
	if err != nil {
		return nil, err
	}
	response := loaded.ToResponse()
	return &response, nil
}

func (s *Service) GetIndicatorByID(ctx context.Context, id string) (*IndicatorResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid indicator id: %w", err)
	}
	ind, err := s.repo.FindIndicatorByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := ind.ToResponse()
	return &response, nil
}

func (s *Service) ListIndicators(ctx context.Context, page, perPage int, competencyID, status string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.FindAllIndicators(ctx, page, perPage, competencyID, status)
	if err != nil {
		return nil, err
	}
	responses := make([]IndicatorResponse, 0, len(list))
	for _, ind := range list {
		responses = append(responses, ind.ToResponse())
	}
	totalPages := calcTotalPages(total, perPage)
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateIndicator(ctx context.Context, id string, req UpdateIndicatorRequest) (*IndicatorResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid indicator id: %w", err)
	}
	ind, err := s.repo.FindIndicatorByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.CompetencyID != nil {
		if compUID, perr := uuid.Parse(*req.CompetencyID); perr == nil {
			ind.CompetencyID = compUID
		}
	}
	if req.Code != nil {
		ind.Code = req.Code
	}
	if req.Statement != nil {
		ind.Statement = *req.Statement
	}
	if req.Description != nil {
		ind.Description = req.Description
	}
	if req.Status != nil {
		ind.Status = *req.Status
	}
	if req.SortOrder != nil {
		ind.SortOrder = *req.SortOrder
	}
	if err := s.repo.UpdateIndicator(ctx, ind); err != nil {
		return nil, err
	}
	loaded, err := s.repo.FindIndicatorByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	response := loaded.ToResponse()
	return &response, nil
}

func (s *Service) DeleteIndicator(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid indicator id: %w", err)
	}
	return s.repo.DeleteIndicator(ctx, uid)
}

// =========================================================================
// Template Indicators
// =========================================================================

func (s *Service) SetTemplateIndicators(ctx context.Context, templateID string, req []TemplateIndicatorRequest) ([]TemplateIndicatorResponse, error) {
	uid, err := uuid.Parse(templateID)
	if err != nil {
		return nil, fmt.Errorf("invalid assessment template id: %w", err)
	}
	if _, err := s.repo.FindAssessmentTemplateByID(ctx, uid); err != nil {
		return nil, err
	}
	if err := s.replaceTemplateChildren(ctx, uid, nil, nil, req); err != nil {
		return nil, err
	}
	items, err := s.repo.ListTemplateIndicators(ctx, uid)
	if err != nil {
		return nil, err
	}
	responses := make([]TemplateIndicatorResponse, 0, len(items))
	for _, it := range items {
		responses = append(responses, it.ToResponse())
	}
	return responses, nil
}

func (s *Service) ListTemplateIndicators(ctx context.Context, templateID string) ([]TemplateIndicatorResponse, error) {
	uid, err := uuid.Parse(templateID)
	if err != nil {
		return nil, fmt.Errorf("invalid assessment template id: %w", err)
	}
	items, err := s.repo.ListTemplateIndicators(ctx, uid)
	if err != nil {
		return nil, err
	}
	responses := make([]TemplateIndicatorResponse, 0, len(items))
	for _, it := range items {
		responses = append(responses, it.ToResponse())
	}
	return responses, nil
}
