package competency

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// =========================================================================
// Reporting (plan generik §20)
// =========================================================================

// EmployeeReportDTO adalah laporan individual seorang employee: hasil,
// gap, self vs others, dan komentar/feedback rater.
type EmployeeReportDTO struct {
	EmployeeID    string            `json:"employee_id"`
	TargetID      string            `json:"target_id"`
	EventID       string            `json:"event_id"`
	OverallScore  float64           `json:"overall_score"`
	TotalGap      float64           `json:"total_gap"`
	SelfScore     float64           `json:"self_score"`
	OthersScore   float64           `json:"others_score"`
	PerceptionGap float64           `json:"perception_gap"`
	Competencies  []CompetencyScoreResult `json:"competencies"`
	Strengths     []GapItem         `json:"strengths"`
	DevelopmentAreas []GapItem      `json:"development_areas"`
	Comments      []string          `json:"comments"`
}

// GetEmployeeReport — Employee Report (§20 Employee Report).
func (s *Service) GetEmployeeReport(ctx context.Context, employeeID string) (*EmployeeReportDTO, error) {
	result, err := s.GetEmployeeResult(ctx, employeeID, "")
	if err != nil {
		return nil, err
	}

	report := &EmployeeReportDTO{
		EmployeeID:    result.EmployeeID,
		TargetID:      result.TargetID,
		EventID:       result.EventID,
		OverallScore:  result.OverallScore,
		TotalGap:      result.TotalGap,
		SelfScore:     result.SelfScore,
		OthersScore:   result.OthersScore,
		PerceptionGap: result.PerceptionGap,
		Competencies:  result.Competencies,
		Strengths:     []GapItem{},
		DevelopmentAreas: []GapItem{},
		Comments:      []string{},
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
			report.Strengths = append(report.Strengths, item)
		} else {
			report.DevelopmentAreas = append(report.DevelopmentAreas, item)
		}
	}

	// Komentar rater untuk subject tsb (anonymous aggregate — §19).
	comments, err := s.repo.ListCommentsByTarget(ctx, result.TargetID)
	if err == nil {
		report.Comments = comments
	}
	return report, nil
}

// ManagerReportDTO — ringkasan kompetensi employee di bawah manager.
type ManagerReportDTO struct {
	EventID        string                  `json:"event_id"`
	Employees      []EmployeeSummaryItem   `json:"employees"`
	TotalEmployees int                     `json:"total_employees"`
	AvgScore       float64                 `json:"avg_score"`
}

type EmployeeSummaryItem struct {
	EmployeeID      string  `json:"employee_id"`
	TargetID        string  `json:"target_id"`
	OverallScore    float64 `json:"overall_score"`
	TotalGap        float64 `json:"total_gap"`
	Status          string  `json:"status"`
	RaterCompletion int     `json:"rater_completion"` // persen (0-100)
}

// GetManagerReport — Manager Report (§20 Manager Report): overview seluruh
// employee pada sebuah event (scope manager = event tsb).
func (s *Service) GetManagerReport(ctx context.Context, eventID string) (*ManagerReportDTO, error) {
	eventUID, err := uuid.Parse(eventID)
	if err != nil {
		return nil, fmt.Errorf("invalid event_id: %w", err)
	}
	targets, err := s.repo.FindTargetsByEvent(ctx, eventUID)
	if err != nil {
		return nil, err
	}

	targetIDs := make([]uuid.UUID, 0, len(targets))
	empIDs := make([]uuid.UUID, 0, len(targets))
	for _, t := range targets {
		targetIDs = append(targetIDs, t.ID)
		if t.EmployeeID != nil {
			empIDs = append(empIDs, *t.EmployeeID)
		}
	}
	submitted, err := s.repo.CountSubmittedRatersByTarget(ctx, targetIDs)
	if err != nil {
		return nil, err
	}
	totalRaters, err := s.repo.CountRatersByTarget(ctx, targetIDs)
	if err != nil {
		return nil, err
	}
	empNames, err := s.repo.GetEmployeeNamesByIDs(ctx, empIDs)
	if err != nil {
		return nil, err
	}

	report := &ManagerReportDTO{
		EventID:        eventUID.String(),
		Employees:      make([]EmployeeSummaryItem, 0, len(targets)),
		TotalEmployees: len(targets),
	}
	totalScore := 0.0
	scored := 0
	for _, t := range targets {
		item := EmployeeSummaryItem{
			EmployeeID:   "",
			TargetID:     t.ID.String(),
			Status:       t.Status,
		}
		if t.EmployeeID != nil {
			item.EmployeeID = t.EmployeeID.String()
			item.TargetID = t.ID.String()
			_ = empNames[t.EmployeeID.String()]
		}
		total := totalRaters[t.ID.String()]
		if total > 0 {
			item.RaterCompletion = submitted[t.ID.String()] * 100 / total
		}
		report.Employees = append(report.Employees, item)

		// Skor hanya untuk target finalized (hasil sudah dihitung).
		if t.Status == "finalized" && t.EmployeeID != nil {
			if r, err := s.CalculateTarget(ctx, t.ID.String()); err == nil {
				item.OverallScore = r.OverallScore
				item.TotalGap = r.TotalGap
				report.Employees[len(report.Employees)-1] = item
				totalScore += r.OverallScore
				scored++
			}
		}
	}
	if scored > 0 {
		report.AvgScore = round2(totalScore / float64(scored))
	}
	return report, nil
}

// HRReportDTO — HR Report (§20 HR Report): distribusi, completion, gap.
type HRReportDTO struct {
	EventID           string               `json:"event_id"`
	TotalTargets      int64                `json:"total_targets"`
	FinalizedTargets  int64                `json:"finalized_targets"`
	RaterCompletion   int                  `json:"rater_completion"` // persen
	AvgScore          float64              `json:"avg_score"`
	TopStrengths      []GapItem            `json:"top_strengths"`
	TopDevelopmentGaps []GapItem           `json:"top_development_gaps"`
}

// GetHRReport — HR Report untuk sebuah event (distribusi, completion,
// strength & development gaps agregat).
func (s *Service) GetHRReport(ctx context.Context, eventID string) (*HRReportDTO, error) {
	eventUID, err := uuid.Parse(eventID)
	if err != nil {
		return nil, fmt.Errorf("invalid event_id: %w", err)
	}
	targets, err := s.repo.FindTargetsByEvent(ctx, eventUID)
	if err != nil {
		return nil, err
	}
	finalizedCount, err := s.repo.CountScoresByEvent(ctx, eventUID)
	if err != nil {
		return nil, err
	}

	targetIDs := make([]uuid.UUID, 0, len(targets))
	for _, t := range targets {
		targetIDs = append(targetIDs, t.ID)
	}
	submitted, err := s.repo.CountSubmittedRatersByTarget(ctx, targetIDs)
	if err != nil {
		return nil, err
	}
	totalRaters, err := s.repo.CountRatersByTarget(ctx, targetIDs)
	if err != nil {
		return nil, err
	}

	report := &HRReportDTO{
		EventID:          eventUID.String(),
		TotalTargets:     int64(len(targets)),
		FinalizedTargets: finalizedCount,
		TopStrengths:     []GapItem{},
		TopDevelopmentGaps: []GapItem{},
	}

	totalSub, totalRat := 0, 0
	for id := range submitted {
		totalSub += submitted[id]
	}
	for id := range totalRaters {
		totalRat += totalRaters[id]
	}
	if totalRat > 0 {
		report.RaterCompletion = totalSub * 100 / totalRat
	}

	// Agregat strength/gap dari seluruh target finalized.
	scoreSum := 0.0
	scored := 0
	gapByCompetency := make(map[string]GapItem)
	for _, t := range targets {
		if t.Status != "finalized" {
			continue
		}
		r, err := s.CalculateTarget(ctx, t.ID.String())
		if err != nil {
			continue
		}
		scoreSum += r.OverallScore
		scored++
		for _, c := range r.Competencies {
			existing, ok := gapByCompetency[c.CompetencyID]
			if !ok {
				existing = GapItem{
					CompetencyID:   c.CompetencyID,
					CompetencyName: c.CompetencyName,
				}
			}
			existing.Gap += c.Gap
			existing.Score += c.Score
			gapByCompetency[c.CompetencyID] = existing
		}
	}
	if scored > 0 {
		report.AvgScore = round2(scoreSum / float64(scored))
	}
	for _, item := range gapByCompetency {
		item.Gap = round2(item.Gap)
		item.Score = round2(item.Score)
		if item.Gap >= 0 {
			report.TopStrengths = append(report.TopStrengths, item)
		} else {
			report.TopDevelopmentGaps = append(report.TopDevelopmentGaps, item)
		}
	}
	return report, nil
}
