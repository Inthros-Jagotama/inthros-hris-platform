package competency

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// =========================================================================
// Calculation Engine (plan generik §14-§18)
// =========================================================================
//
// Input:
//   - Rater responses per indicator
//   - Rater type (self/superior/peer/subordinate) + rater weight
//   - Template rater type weight (§10)
//   - Template competency required level & weight (§5.2)
//
// Output (per competency): rater type scores, competency score (weighted
// rater-type average), required level, gap, weighted gap, self vs others.
// Overall score = weighted average kompetensi.

type raterTypeScores struct {
	Self        float64 `json:"self"`
	Superior    float64 `json:"superior"`
	Peer        float64 `json:"peer"`
	Subordinate float64 `json:"subordinate"`
	Other       float64 `json:"other"`
}

type CompetencyScoreResult struct {
	CompetencyID    string         `json:"competency_id"`
	CompetencyName  string         `json:"competency_name"`
	RequiredLevel   float64        `json:"required_level"`
	Score           float64        `json:"score"`
	Gap             float64        `json:"gap"`
	WeightedGap     float64        `json:"weighted_gap"`
	Weight          float64        `json:"weight"`
	RaterScores     raterTypeScores `json:"rater_scores"`
}

type AssessmentResult struct {
	TargetID            string                 `json:"target_id"`
	EventID             string                 `json:"event_id"`
	EmployeeID          string                 `json:"employee_id"`
	OrganizationID      string                 `json:"organization_id"`
	OverallScore        float64                `json:"overall_score"`
	TotalGap            float64                `json:"total_gap"`
	SelfScore           float64                `json:"self_score"`
	OthersScore         float64                `json:"others_score"`
	PerceptionGap       float64                `json:"perception_gap"`
	Competencies        []CompetencyScoreResult `json:"competencies"`
}

// CalculateTarget menghitung hasil 360 untuk sebuah assessment target dari
// response rater yang sudah submit. Murni logic — tidak menyentuh DB selain
// input lookup; persisten dilakukan FinalizeTarget.
func (s *Service) CalculateTarget(ctx context.Context, targetID string) (*AssessmentResult, error) {
	targetUID, err := uuid.Parse(targetID)
	if err != nil {
		return nil, fmt.Errorf("invalid target id: %w", err)
	}
	target, err := s.repo.FindCompetencyEventTargetByID(ctx, targetUID)
	if err != nil {
		return nil, err
	}
	event, err := s.repo.FindCompetencyEventByID(ctx, target.CompetencyEventID)
	if err != nil {
		return nil, err
	}
	if event.TemplateID == nil {
		return nil, fmt.Errorf("assessment event has no template configured")
	}

	raters, err := s.repo.FindAllRatersByTarget(ctx, targetUID)
	if err != nil {
		return nil, err
	}
	tplCompetencies, err := s.repo.FindTemplateCompetencies(ctx, *event.TemplateID)
	if err != nil {
		return nil, err
	}
	tplRaterTypes, err := s.repo.FindTemplateRaterTypes(ctx, *event.TemplateID)
	if err != nil {
		return nil, err
	}
	tplIndicators, err := s.repo.ListTemplateIndicatorsByCompetency(ctx, *event.TemplateID)
	if err != nil {
		return nil, err
	}

	// Mapping indicator -> competency (dari template indicators).
	indicatorCompetency := make(map[string]string)
	for _, ti := range tplIndicators {
		if ti.Indicator != nil {
			indicatorCompetency[ti.IndicatorID.String()] = ti.Indicator.CompetencyID.String()
		}
	}

	// Mapping rater type -> weight template (§10).
	typeWeight := make(map[string]float64)
	for _, rt := range tplRaterTypes {
		typeWeight[rt.RaterType] = rt.Weight
	}

	result := &AssessmentResult{
		TargetID:       target.ID.String(),
		EventID:        event.ID.String(),
		EmployeeID:     target.EmployeeID.String(),
		OrganizationID: target.OrganizationID.String(),
		Competencies:   make([]CompetencyScoreResult, 0, len(tplCompetencies)),
	}

	totalWeightedScore := 0.0
	totalWeightSum := 0.0
	totalGap := 0.0
	var selfScores, otherScores []float64

	for _, tc := range tplCompetencies {
		compID := tc.CompetencyID.String()

		// Kumpulkan response per rater type untuk indicator competency ini.
		typeRatings := make(map[string][]float64) // rater_type -> ratings per rater average
		raterCountByType := make(map[string]int)
		for _, rat := range raters {
			if rat.Status != string(RaterStatusSubmitted) {
				continue
			}
			var ratings []float64
			for _, resp := range rat.Responses {
				if indicatorCompetency[resp.IndicatorID.String()] != compID {
					continue
				}
				ratings = append(ratings, float64(resp.RatingValue))
			}
			if len(ratings) == 0 {
				continue
			}
			avg := avgOf(ratings)
			// Weight rater-level bila di-set (default 1), lalu rata-rata.
			w := rat.Weight
			if w <= 0 {
				w = 1
			}
			typeRatings[rat.RaterType] = append(typeRatings[rat.RaterType], avg*w)
			raterCountByType[rat.RaterType]++
		}

		// Skor per rater type: rata-rata skor rater (sudah di-weight rater).
		scoresByType := make(map[string]float64)
		for rtype, vals := range typeRatings {
			scoresByType[rtype] = avgOf(vals)
		}

		// Competency score = weighted average across rater types.
		weightedSum := 0.0
		weightSum := 0.0
		for rtype, score := range scoresByType {
			w := typeWeight[rtype]
			weightedSum += score * w
			weightSum += w
		}
		compScore := 0.0
		if weightSum > 0 {
			compScore = weightedSum / weightSum
		}

		requiredLevel := 0.0
		if tc.RequiredLevel != nil {
			requiredLevel = float64(*tc.RequiredLevel)
		}
		gap := compScore - requiredLevel

		compWeight := tc.Weight
		if compWeight <= 0 {
			compWeight = 1
		}

		cr := CompetencyScoreResult{
			CompetencyID:  compID,
			RequiredLevel: requiredLevel,
			Score:         round2(compScore),
			Gap:           round2(gap),
			WeightedGap:   round2(gap * compWeight),
			Weight:        compWeight,
		}
		if tc.Competency != nil {
			cr.CompetencyName = tc.Competency.Name
		}
		if v, ok := scoresByType[string(RaterTypeSelf)]; ok {
			cr.RaterScores.Self = round2(v)
			selfScores = append(selfScores, v)
		}
		if v, ok := scoresByType[string(RaterTypeSuperior)]; ok {
			cr.RaterScores.Superior = round2(v)
			otherScores = append(otherScores, v)
		}
		if v, ok := scoresByType[string(RaterTypePeer)]; ok {
			cr.RaterScores.Peer = round2(v)
			otherScores = append(otherScores, v)
		}
		if v, ok := scoresByType[string(RaterTypeSubordinate)]; ok {
			cr.RaterScores.Subordinate = round2(v)
			otherScores = append(otherScores, v)
		}
		if v, ok := scoresByType[string(RaterTypeOther)]; ok {
			cr.RaterScores.Other = round2(v)
			otherScores = append(otherScores, v)
		}

		result.Competencies = append(result.Competencies, cr)
		totalWeightedScore += compScore * compWeight
		totalWeightSum += compWeight
		totalGap += gap * compWeight
	}

	if totalWeightSum > 0 {
		result.OverallScore = round2(totalWeightedScore / totalWeightSum)
		result.TotalGap = round2(totalGap / totalWeightSum)
	}
	if len(selfScores) > 0 {
		result.SelfScore = round2(avgOf(selfScores))
	}
	if len(otherScores) > 0 {
		result.OthersScore = round2(avgOf(otherScores))
	}
	result.PerceptionGap = round2(result.SelfScore - result.OthersScore)

	return result, nil
}

// FinalizeTarget menghitung hasil target lalu menulis snapshot ke
// competency_scores + competency_score_details (per employee per event —
// §15/§16). Dipanggil setelah approval APPROVED (Phase 4 status handler).
func (s *Service) FinalizeTarget(ctx context.Context, targetID string) (*AssessmentResult, error) {
	result, err := s.CalculateTarget(ctx, targetID)
	if err != nil {
		return nil, err
	}

	eventUID, _ := uuid.Parse(result.EventID)
	empUID, _ := uuid.Parse(result.EmployeeID)
	orgUID, _ := uuid.Parse(result.OrganizationID)

	now := time.Now()
	score, err := s.repo.FindScoreByEventAndEmployee(ctx, eventUID, empUID)
	if err != nil {
		// Belum ada — buat baru.
		score = &CompetencyScore{
			OrganizationID:      orgUID,
			EmployeeID:          &empUID,
			CompetencyEventID:   &eventUID,
			TotalGapPercentage:  result.TotalGap,
			TotalGradePercentage: result.OverallScore,
			AssessedAt:          &now,
		}
		if err := s.repo.CreateCompetencyScore(ctx, score); err != nil {
			return nil, err
		}
	} else {
		score.TotalGapPercentage = result.TotalGap
		score.TotalGradePercentage = result.OverallScore
		score.AssessedAt = &now
		if err := s.repo.UpdateCompetencyScore(ctx, score); err != nil {
			return nil, err
		}
	}

	details := make([]CompetencyScoreDetail, 0, len(result.Competencies))
	for _, c := range result.Competencies {
		compUID, _ := uuid.Parse(c.CompetencyID)
		level := int(c.Score)
		stdLevel := int(c.RequiredLevel)
		details = append(details, CompetencyScoreDetail{
			CompetencyID:          compUID,
			Type:                  "technical",
			StandardLevel:         &stdLevel,
			EmployeeLevel:         &level,
			GapPercentage:         c.Gap,
			WeightedGapPercentage: c.WeightedGap,
		})
	}
	if err := s.repo.ReplaceScoreDetails(ctx, score.ID, details); err != nil {
		return nil, err
	}

	s.logger.Info("Competency assessment finalized",
		zap.String("target_id", result.TargetID),
		zap.String("employee_id", result.EmployeeID),
		zap.Float64("overall_score", result.OverallScore),
	)
	return result, nil
}

// =========================================================================
// Result & Gap Queries (§17/§18, §22 Result)
// =========================================================================

// GetEmployeeResult mengembalikan hasil 360 terbaru untuk seorang employee
// (opsional difilter per event). Bila hasil belum pernah di-finalize, dihitung
// on-the-fly dari response — hasil tersimpan dipakai bila sudah ada.
func (s *Service) GetEmployeeResult(ctx context.Context, employeeID string, eventID string) (*AssessmentResult, error) {
	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	if eventID != "" {
		// Hasil per event: cari target employee di event tsb.
		eventUID, err := uuid.Parse(eventID)
		if err != nil {
			return nil, fmt.Errorf("invalid event_id: %w", err)
		}
		target, err := s.repo.FindTargetByEventAndEmployee(ctx, eventUID, empUID)
		if err != nil {
			return nil, err
		}
		return s.CalculateTarget(ctx, target.ID.String())
	}
	// Tanpa event: target terbaru employee (finalized dulu, lalu draft).
	targets, err := s.repo.FindTargetsByEmployee(ctx, empUID)
	if err != nil {
		return nil, err
	}
	if len(targets) == 0 {
		return nil, fmt.Errorf("no competency assessment found for this employee")
	}
	var best *CompetencyEventTarget
	for i := range targets {
		if targets[i].Status == "finalized" {
			best = &targets[i]
			break
		}
	}
	if best == nil {
		best = &targets[0]
	}
	return s.CalculateTarget(ctx, best.ID.String())
}

// GetEmployeeGap mengembalikan analisis gap per competency (strength vs
// development area — §17) untuk seorang employee.
func (s *Service) GetEmployeeGap(ctx context.Context, employeeID string) (*GapAnalysisResponse, error) {
	result, err := s.GetEmployeeResult(ctx, employeeID, "")
	if err != nil {
		return nil, err
	}
	resp := &GapAnalysisResponse{
		TargetID:       result.TargetID,
		EmployeeID:     result.EmployeeID,
		OverallScore:   result.OverallScore,
		TotalGap:       result.TotalGap,
		SelfScore:      result.SelfScore,
		OthersScore:    result.OthersScore,
		PerceptionGap:  result.PerceptionGap,
		Strengths:      make([]GapItem, 0),
		DevelopmentAreas: make([]GapItem, 0),
	}
	for _, c := range result.Competencies {
		item := GapItem{
			CompetencyID:   c.CompetencyID,
			CompetencyName: c.CompetencyName,
			RequiredLevel:  c.RequiredLevel,
			Score:          c.Score,
			Gap:            c.Gap,
			WeightedGap:    c.WeightedGap,
		}
		if c.Score >= c.RequiredLevel {
			resp.Strengths = append(resp.Strengths, item)
		} else {
			resp.DevelopmentAreas = append(resp.DevelopmentAreas, item)
		}
	}
	return resp, nil
}

// =========================================================================
// Helpers
// =========================================================================

func avgOf(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}
