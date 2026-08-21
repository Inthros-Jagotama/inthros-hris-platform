package careerintelligence

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

// PerformanceProvider membaca skor final performance evaluation terakhir
// milik employee (skala 0-100) dari modul Performance Management. Narrow-
// interface-plus-adapter pattern (sama seperti employeemovement's
// PerformanceProvider) — careerintelligence tidak menghitung KPI/OKR sendiri.
type PerformanceProvider interface {
	LatestFinalScore(ctx context.Context, employeeID uuid.UUID) (float64, bool, error)
}

// CompetencyProvider membaca skor kompetensi terakhir milik employee (skala
// 0-100) dari modul Competency.
type CompetencyProvider interface {
	LatestScore(ctx context.Context, employeeID uuid.UUID) (float64, bool, error)
}

type Service struct {
	repo                *Repository
	logger              *zap.Logger
	performanceProvider PerformanceProvider
	competencyProvider  CompetencyProvider
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// SetPerformanceProvider wires the Performance Management score source used
// by GenerateTalentMap. Must be called before GenerateTalentMap is used.
func (s *Service) SetPerformanceProvider(p PerformanceProvider) {
	s.performanceProvider = p
}

// SetCompetencyProvider wires the Competency score source used by
// GenerateTalentMap. Must be called before GenerateTalentMap is used.
func (s *Service) SetCompetencyProvider(p CompetencyProvider) {
	s.competencyProvider = p
}

// =========================================================================
// Grid position calculator
// =========================================================================

func computeGridPosition(performance, potential string) string {
	// 9-box grid: performance (rows) x potential (columns)
	// Performance: LOW=3, MEDIUM=2, HIGH=1
	// Potential:   LOW=3, MEDIUM=2, HIGH=1
	pMap := map[string]int{"HIGH": 1, "MEDIUM": 2, "LOW": 3}
	row := pMap[performance]
	col := pMap[potential]
	return fmt.Sprintf("9-BOX-%d", (row-1)*3+col)
}

func gridQuadrantLabel(position string) (string, string) {
	labels := map[string][2]string{
		"9-BOX-1": {"High Performer - High Potential", "Star performers ready for accelerated growth"},
		"9-BOX-2": {"High Performer - Medium Potential", "Solid performers with growth potential"},
		"9-BOX-3": {"High Performer - Low Potential", "Expert contributors at career maturity"},
		"9-BOX-4": {"Medium Performer - High Potential", "Rising stars needing development"},
		"9-BOX-5": {"Medium Performer - Medium Potential", "Core contributors, steady performers"},
		"9-BOX-6": {"Medium Performer - Low Potential", "Consistent performers in current role"},
		"9-BOX-7": {"Low Performer - High Potential", "Diamond in the rough, needs coaching"},
		"9-BOX-8": {"Low Performer - Medium Potential", "Needs improvement plan and monitoring"},
		"9-BOX-9": {"Low Performer - Low Potential", "Performance improvement plan required"},
	}
	if label, ok := labels[position]; ok {
		return label[0], label[1]
	}
	return "Unknown", ""
}

// =========================================================================
// Talent Map Settings
// =========================================================================

func talentMapSettingsToResponse(s *TalentMapSettings) *TalentMapSettingsResponse {
	return &TalentMapSettingsResponse{
		PerformanceLowMax:  s.PerformanceLowMax,
		PerformanceHighMin: s.PerformanceHighMin,
		PotentialLowMax:    s.PotentialLowMax,
		PotentialHighMin:   s.PotentialHighMin,
	}
}

func (s *Service) GetTalentMapSettings(ctx context.Context) (*TalentMapSettingsResponse, error) {
	settings, err := s.repo.GetOrCreateTalentMapSettings(ctx)
	if err != nil {
		return nil, err
	}
	return talentMapSettingsToResponse(settings), nil
}

func (s *Service) UpdateTalentMapSettings(ctx context.Context, req UpdateTalentMapSettingsRequest) (*TalentMapSettingsResponse, error) {
	if req.PerformanceLowMax >= req.PerformanceHighMin {
		return nil, fmt.Errorf("performance_low_max harus lebih kecil dari performance_high_min")
	}
	if req.PotentialLowMax >= req.PotentialHighMin {
		return nil, fmt.Errorf("potential_low_max harus lebih kecil dari potential_high_min")
	}
	settings, err := s.repo.GetOrCreateTalentMapSettings(ctx)
	if err != nil {
		return nil, err
	}
	settings.PerformanceLowMax = req.PerformanceLowMax
	settings.PerformanceHighMin = req.PerformanceHighMin
	settings.PotentialLowMax = req.PotentialLowMax
	settings.PotentialHighMin = req.PotentialHighMin
	if err := s.repo.UpdateTalentMapSettings(ctx, settings); err != nil {
		return nil, err
	}
	return talentMapSettingsToResponse(settings), nil
}

// bandScore converts a raw 0-100 score into LOW/MEDIUM/HIGH using the
// tenant's configured thresholds: score < lowMax => LOW;
// lowMax <= score < highMin => MEDIUM; score >= highMin => HIGH.
func bandScore(score, lowMax, highMin float64) string {
	if score < lowMax {
		return "LOW"
	}
	if score < highMin {
		return "MEDIUM"
	}
	return "HIGH"
}

// =========================================================================
// Talent Map
// =========================================================================

// GenerateTalentMap membuat talent map otomatis dari skor Performance
// Management (performance) & Competency (potential) milik employee, tanpa
// input manual. Skor terbaru masing-masing diambil via provider, lalu
// dibanding LOW/MEDIUM/HIGH memakai TalentMapSettings tenant.
func (s *Service) GenerateTalentMap(ctx context.Context, req GenerateTalentMapRequest) (*TalentMapResponse, error) {
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	if s.performanceProvider == nil || s.competencyProvider == nil {
		return nil, fmt.Errorf("talent map generation belum dikonfigurasi (performance/competency provider)")
	}

	perfScore, found, err := s.performanceProvider.LatestFinalScore(ctx, empID)
	if err != nil {
		return nil, fmt.Errorf("failed to read performance score: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("karyawan belum memiliki hasil evaluasi performance yang selesai (COMPLETED/ACTUAL_APPROVED)")
	}

	potScore, found, err := s.competencyProvider.LatestScore(ctx, empID)
	if err != nil {
		return nil, fmt.Errorf("failed to read competency score: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("karyawan belum memiliki hasil assessment competency")
	}

	settings, err := s.repo.GetOrCreateTalentMapSettings(ctx)
	if err != nil {
		return nil, err
	}
	performance := bandScore(perfScore, settings.PerformanceLowMax, settings.PerformanceHighMin)
	potential := bandScore(potScore, settings.PotentialLowMax, settings.PotentialHighMin)

	gridPos := computeGridPosition(performance, potential)
	tm := &CareerTalentMap{
		EmployeeID:   empID,
		Period:       req.Period,
		Performance:  performance,
		Potential:    potential,
		GridPosition: gridPos,
		Notes:        req.Notes,
		AssessedAt:   time.Now(),
	}
	if err := s.repo.CreateTalentMap(ctx, tm); err != nil {
		return nil, err
	}
	s.logger.Info("Talent map generated",
		zap.String("employee", req.EmployeeID),
		zap.Float64("performance_score", perfScore),
		zap.Float64("potential_score", potScore),
		zap.String("position", gridPos))
	return talentMapToResponse(tm), nil
}

func (s *Service) CreateTalentMap(ctx context.Context, req CreateTalentMapRequest) (*TalentMapResponse, error) {
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	gridPos := computeGridPosition(req.Performance, req.Potential)
	tm := &CareerTalentMap{
		EmployeeID:   empID,
		Period:       req.Period,
		Performance:  req.Performance,
		Potential:    req.Potential,
		GridPosition: gridPos,
		Notes:        req.Notes,
		AssessedAt:   time.Now(),
	}
	if err := s.repo.CreateTalentMap(ctx, tm); err != nil {
		return nil, err
	}
	s.logger.Info("Talent map created", zap.String("employee", req.EmployeeID), zap.String("position", gridPos))
	return talentMapToResponse(tm), nil
}

func (s *Service) GetTalentMapByID(ctx context.Context, id string) (*TalentMapResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	tm, err := s.repo.FindTalentMapByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return talentMapToResponse(tm), nil
}

func (s *Service) ListTalentMaps(ctx context.Context, period, employeeID string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListTalentMaps(ctx, period, employeeID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]TalentMapResponse, 0, len(list))
	for _, tm := range list {
		responses = append(responses, *talentMapToResponse(&tm))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateTalentMap(ctx context.Context, id string, req UpdateTalentMapRequest) (*TalentMapResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	tm, err := s.repo.FindTalentMapByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Performance != nil {
		tm.Performance = *req.Performance
	}
	if req.Potential != nil {
		tm.Potential = *req.Potential
	}
	tm.GridPosition = computeGridPosition(tm.Performance, tm.Potential)
	if req.Notes != nil {
		tm.Notes = *req.Notes
	}
	if err := s.repo.UpdateTalentMap(ctx, tm); err != nil {
		return nil, err
	}
	return talentMapToResponse(tm), nil
}

func (s *Service) DeleteTalentMap(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteTalentMap(ctx, uid)
}

func (s *Service) GetTalentGrid(ctx context.Context, period string) (*TalentGridResponse, error) {
	list, err := s.repo.GetTalentGrid(ctx, period)
	if err != nil {
		return nil, err
	}
	quadrantMap := make(map[string]int)
	for _, tm := range list {
		quadrantMap[tm.GridPosition]++
	}
	var quadrants []TalentQuadrant
	for pos := 1; pos <= 9; pos++ {
		key := fmt.Sprintf("9-BOX-%d", pos)
		label, desc := gridQuadrantLabel(key)
		quadrants = append(quadrants, TalentQuadrant{
			Label:       label,
			Position:    key,
			Count:       quadrantMap[key],
			Description: desc,
		})
	}
	return &TalentGridResponse{
		Period:    period,
		Quadrants: quadrants,
		Total:     len(list),
	}, nil
}

func (s *Service) GetEmployeeTalentProfile(ctx context.Context, employeeID string) (*EmployeeTalentProfileResponse, error) {
	empID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	// Get current talent map
	maps, _, err := s.repo.ListTalentMaps(ctx, "", employeeID, 1, 1)
	var currentMap *TalentMapResponse
	if err == nil && len(maps) > 0 {
		currentMap = talentMapToResponse(&maps[0])
	}
	// Get history
	history, _ := s.repo.GetTalentHistoryByEmployee(ctx, empID)
	historyResp := make([]TalentMapResponse, 0, len(history))
	for _, tm := range history {
		historyResp = append(historyResp, *talentMapToResponse(&tm))
	}
	// Get interests
	interests, _ := s.repo.GetInterestsByEmployee(ctx, empID)
	interestResp := make([]CareerInterestResponse, 0, len(interests))
	for _, ci := range interests {
		interestResp = append(interestResp, *careerInterestToResponse(&ci))
	}
	// Determine readiness (simplified)
	readyFor := []string{}
	if currentMap != nil && currentMap.GridPosition == "9-BOX-1" {
		readyFor = append(readyFor, "Leadership Track", "Senior Position")
	} else if currentMap != nil && currentMap.GridPosition == "9-BOX-4" {
		readyFor = append(readyFor, "Team Lead", "Specialist Track")
	}
	return &EmployeeTalentProfileResponse{
		EmployeeID: employeeID,
		CurrentMap: currentMap,
		History:    historyResp,
		Interests:  interestResp,
		ReadyFor:   readyFor,
	}, nil
}

// =========================================================================
// Career Interest
// =========================================================================

func (s *Service) CreateCareerInterest(ctx context.Context, req CreateCareerInterestRequest) (*CareerInterestResponse, error) {
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	ci := &CareerInterest{
		EmployeeID:       empID,
		InterestType:     req.InterestType,
		TargetPosition:   req.TargetPosition,
		TargetDepartment: req.TargetDepartment,
		Motivation:       req.Motivation,
		ReadinessLevel:   req.ReadinessLevel,
		IsActive:         true,
		RecordedAt:       time.Now(),
	}
	if err := s.repo.CreateCareerInterest(ctx, ci); err != nil {
		return nil, err
	}
	return careerInterestToResponse(ci), nil
}

func (s *Service) ListCareerInterests(ctx context.Context, employeeID string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListCareerInterests(ctx, employeeID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]CareerInterestResponse, 0, len(list))
	for _, ci := range list {
		responses = append(responses, *careerInterestToResponse(&ci))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) GetEmployeeCareerInterests(ctx context.Context, employeeID string) ([]CareerInterestResponse, error) {
	empID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	list, err := s.repo.GetInterestsByEmployee(ctx, empID)
	if err != nil {
		return nil, err
	}
	responses := make([]CareerInterestResponse, 0, len(list))
	for _, ci := range list {
		responses = append(responses, *careerInterestToResponse(&ci))
	}
	return responses, nil
}

// =========================================================================
// Career Path
// =========================================================================

// CreateCareerPathLadder membuat jenjang penuh (nama + langkah berurutan) pada
// skema terpadu. Inilah endpoint utama FE Career Paths (strategical).
// Validasi: minimal 1 step, sequence unik per path, posisi unik per path, dan
// semua position_id merujuk posisi yang ada (pola EM §12.9).
func (s *Service) CreateCareerPathLadder(ctx context.Context, req CreateCareerPathLadderRequest) (*CareerPathResponse, error) {
	cp := &CareerPath{
		Name:        req.Name,
		IsActive:    true,
		Description: "",
	}
	if req.Description != nil {
		cp.Description = *req.Description
	}
	if req.IsActive != nil {
		cp.IsActive = *req.IsActive
	}
	if cp.ID == uuid.Nil {
		cp.ID = uuid.New()
	}

	steps, err := s.buildCareerPathSteps(ctx, cp.ID, req.Steps)
	if err != nil {
		return nil, err
	}
	if err := s.repo.CreateCareerPathTx(ctx, cp, steps); err != nil {
		return nil, err
	}
	s.logger.Info("Career path (ladder) created",
		zap.String("path_id", cp.ID.String()),
		zap.String("name", cp.Name),
		zap.Int("steps", len(steps)),
	)
	return s.careerPathToResponse(ctx, cp, steps), nil
}

// GetCareerPathByID mengambil satu path lengkap dengan steps terurut sequence.
func (s *Service) GetCareerPathByID(ctx context.Context, id string) (*CareerPathResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	cp, err := s.repo.FindCareerPathByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	steps, err := s.repo.ListCareerPathStepsByPathID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return s.careerPathToResponse(ctx, cp, steps), nil
}

// UpdateCareerPath memperbarui header path dan mengganti seluruh steps-nya
// (semantik full-replace).
func (s *Service) UpdateCareerPath(ctx context.Context, id string, req UpdateCareerPathRequest) (*CareerPathResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	cp, err := s.repo.FindCareerPathByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		cp.Name = *req.Name
	}
	if req.Description != nil {
		cp.Description = *req.Description
	}
	if req.IsActive != nil {
		cp.IsActive = *req.IsActive
	}
	steps, err := s.buildCareerPathSteps(ctx, cp.ID, req.Steps)
	if err != nil {
		return nil, err
	}
	if err := s.repo.UpdateCareerPathTx(ctx, cp, steps); err != nil {
		return nil, err
	}
	s.logger.Info("Career path updated",
		zap.String("path_id", cp.ID.String()),
		zap.String("name", cp.Name),
	)
	return s.careerPathToResponse(ctx, cp, steps), nil
}

// buildCareerPathSteps membangun model steps dari request ladder dan
// memvalidasinya: sequence unik, posisi unik, dan semua position_id ada
// (pola EM buildAndValidateCareerPathSteps).
func (s *Service) buildCareerPathSteps(ctx context.Context, careerPathID uuid.UUID, reqSteps []CreateCareerPathStepRequest) ([]CareerPathStep, error) {
	steps := make([]CareerPathStep, 0, len(reqSteps))
	positionIDs := make([]uuid.UUID, 0, len(reqSteps))
	seenSeq := map[int]bool{}
	seenPos := map[string]bool{}
	for _, rs := range reqSteps {
		posID, err := uuid.Parse(rs.PositionID)
		if err != nil {
			return nil, fmt.Errorf("invalid position_id: %w", err)
		}
		if seenSeq[rs.Sequence] {
			return nil, fmt.Errorf("duplicate sequence %d in career path steps", rs.Sequence)
		}
		if seenPos[posID.String()] {
			return nil, fmt.Errorf("position %s appears more than once in career path", posID.String())
		}
		seenSeq[rs.Sequence] = true
		seenPos[posID.String()] = true
		positionIDs = append(positionIDs, posID)
		steps = append(steps, CareerPathStep{
			ID:                   uuid.New(),
			CareerPathID:         careerPathID,
			PositionID:           posID,
			Sequence:             rs.Sequence,
			MinimumServiceMonths: rs.MinimumServiceMonths,
			Requirements:         rs.Requirements,
		})
	}
	// Semua posisi harus eksis di tabel positions (bukan hanya format UUID).
	names, err := s.repo.GetPositionNamesByIDs(ctx, positionIDs)
	if err != nil {
		return nil, err
	}
	if len(names) != len(positionIDs) {
		return nil, fmt.Errorf("one or more position ids do not exist")
	}
	return steps, nil
}

func (s *Service) ListCareerPaths(ctx context.Context, page, perPage int, keyword string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListCareerPaths(ctx, page, perPage, keyword)
	if err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, 0, len(list))
	for _, cp := range list {
		ids = append(ids, cp.ID)
	}
	stepsByPath, err := s.repo.ListCareerPathStepsByPathIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	responses := make([]CareerPathResponse, 0, len(list))
	for _, cp := range list {
		steps := stepsByPath[cp.ID]
		responses = append(responses, *s.careerPathToResponse(ctx, &cp, steps))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) DeleteCareerPath(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCareerPath(ctx, uid)
}

func (s *Service) GetGapAnalysis(ctx context.Context, req GapAnalysisRequest) (*GapAnalysisResponse, error) {
	_, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	targetID, err := uuid.Parse(req.TargetTitleID)
	if err != nil {
		return nil, fmt.Errorf("invalid target_title_id: %w", err)
	}

	// Get target position name
	targetTitle, _ := s.repo.GetPositionTitle(ctx, targetID)

	// Simplified gap analysis — would query competency module in production
	matchedSkills := 3
	totalRequired := 8
	gapPct := float64(totalRequired-matchedSkills) / float64(totalRequired) * 100

	return &GapAnalysisResponse{
		EmployeeID:    req.EmployeeID,
		TargetTitle:   targetTitle,
		MatchedSkills: matchedSkills,
		TotalRequired: totalRequired,
		GapPercentage: gapPct,
		Recommendations: []GapRecommendation{
			{Category: "TRAINING", Description: "Leadership Development Program", Priority: "HIGH"},
			{Category: "EXPERIENCE", Description: "Cross-functional project assignment", Priority: "MEDIUM"},
			{Category: "CERTIFICATION", Description: "Professional certification in target domain", Priority: "MEDIUM"},
		},
		EstimatedTimeline: "12-18 months",
	}, nil
}

// =========================================================================
// Internal Candidate Eligibility (S-4 — CI → Recruitment)
// =========================================================================

// GetEligibleEmployeesForPath mengembalikan employee internal yang eligible
// untuk target position dari sebuah career path: employee dengan employment
// aktif yang saat ini memegang posisi pada source step (semua step sebelum
// step terakhir) dari path tersebut. Target = step terakhir. CI menentukan
// eligibility — Recruitment tidak menghitung sendiri (plan S-4).
func (s *Service) GetEligibleEmployeesForPath(ctx context.Context, pathID string) ([]EligibleEmployeeResponse, error) {
	uid, err := uuid.Parse(pathID)
	if err != nil {
		return nil, fmt.Errorf("invalid path_id: %w", err)
	}
	cp, err := s.repo.FindCareerPathByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	steps, err := s.repo.ListCareerPathStepsByPathID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if len(steps) < 2 {
		return nil, fmt.Errorf("career path must have at least 2 steps (source + target) to compute eligibility")
	}
	// Step terakhir = target; semua step sebelumnya = source steps.
	target := steps[len(steps)-1]
	sourceSteps := steps[:len(steps)-1]

	sourcePositionIDs := make([]uuid.UUID, 0, len(sourceSteps))
	for _, st := range sourceSteps {
		sourcePositionIDs = append(sourcePositionIDs, st.PositionID)
	}
	rows, err := s.repo.ListEligibleEmployeesByPositionIDs(ctx, sourcePositionIDs)
	if err != nil {
		return nil, err
	}

	// Resolve nama target + source positions.
	allIDs := append(sourcePositionIDs, target.PositionID)
	names, err := s.repo.GetPositionNamesByIDs(ctx, allIDs)
	if err != nil {
		return nil, err
	}
	// Peta posisi → sequence (untuk source_step_sequence per employee).
	seqByPosition := make(map[string]int, len(sourceSteps))
	for _, st := range sourceSteps {
		seqByPosition[st.PositionID.String()] = st.Sequence
	}

	responses := make([]EligibleEmployeeResponse, 0, len(rows))
	for _, rw := range rows {
		responses = append(responses, EligibleEmployeeResponse{
			EmployeeID:          rw.EmployeeID,
			Name:                rw.Name,
			CurrentPositionID:   rw.PositionID,
			CurrentPositionName: rw.PositionName,
			SourceStepSequence:  seqByPosition[rw.PositionID],
			TargetPositionID:    target.PositionID.String(),
			TargetPositionName:  names[target.PositionID.String()],
			PathID:              cp.ID.String(),
			PathName:            cp.Name,
		})
	}
	return responses, nil
}

// GetEligibleEmployeesByPosition mengembalikan employee eligible untuk sebuah
// target position — menjembatani recruitment.InternalCandidateProvider (plan
// S-4): mencari semua career path aktif yang menuju posisi tsb, lalu
// mengumpulkan eligible employees dari seluruh path tersebut.
func (s *Service) GetEligibleEmployeesByPosition(ctx context.Context, targetPositionID string) ([]EligibleEmployeeResponse, error) {
	targetID, err := uuid.Parse(targetPositionID)
	if err != nil {
		return nil, fmt.Errorf("invalid position_id: %w", err)
	}
	paths, err := s.repo.FindCareerPathsByTarget(ctx, targetID)
	if err != nil {
		return nil, err
	}
	var result []EligibleEmployeeResponse
	for _, p := range paths {
		emps, err := s.GetEligibleEmployeesForPath(ctx, p.ID.String())
		if err != nil {
			// Path tanpa ≥2 steps dilewati (bukan kegagalan total); error lain
			// (mis. DB) tetap dicatat agar tidak tersembunyi.
			if !strings.Contains(err.Error(), "at least 2 steps") {
				s.logger.Warn("skip career path in eligible-employees aggregation",
					zap.String("path_id", p.ID.String()),
					zap.Error(err))
			}
			continue
		}
		result = append(result, emps...)
	}
	if result == nil {
		result = []EligibleEmployeeResponse{}
	}
	return result, nil
}

// =========================================================================
// Succession Plan
// =========================================================================

func (s *Service) CreateSuccessionPlan(ctx context.Context, req CreateSuccessionPlanRequest) (*SuccessionPlanResponse, error) {
	posID, err := uuid.Parse(req.PositionID)
	if err != nil {
		return nil, fmt.Errorf("invalid position_id: %w", err)
	}
	succID, err := uuid.Parse(req.SuccessorID)
	if err != nil {
		return nil, fmt.Errorf("invalid successor_id: %w", err)
	}
	sp := &CareerSuccessionPlan{
		PositionID:      posID,
		SuccessorID:     succID,
		ReadinessLevel:  req.ReadinessLevel,
		PriorityOrder:   1,
		DevelopmentPlan: req.DevelopmentPlan,
		Notes:           req.Notes,
		Status:          "ACTIVE",
	}
	if req.PriorityOrder != nil {
		sp.PriorityOrder = *req.PriorityOrder
	}
	if req.TargetDate != "" {
		t, err := time.Parse("2006-01-02", req.TargetDate)
		if err == nil {
			sp.TargetDate = &t
		}
	}
	if err := s.repo.CreateSuccessionPlan(ctx, sp); err != nil {
		return nil, err
	}
	s.logger.Info("Succession plan created",
		zap.String("position", req.PositionID),
		zap.String("successor", req.SuccessorID),
		zap.String("readiness", req.ReadinessLevel))
	return successionPlanToResponse(sp), nil
}

func (s *Service) ListSuccessionPlans(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListSuccessionPlans(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]SuccessionPlanResponse, 0, len(list))
	for _, sp := range list {
		responses = append(responses, *successionPlanToResponse(&sp))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) GetSuccessionPlanByID(ctx context.Context, id string) (*SuccessionPlanResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	sp, err := s.repo.FindSuccessionPlanByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return successionPlanToResponse(sp), nil
}

func (s *Service) UpdateSuccessionPlan(ctx context.Context, id string, req UpdateSuccessionPlanRequest) (*SuccessionPlanResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	sp, err := s.repo.FindSuccessionPlanByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.ReadinessLevel != nil {
		sp.ReadinessLevel = *req.ReadinessLevel
	}
	if req.PriorityOrder != nil {
		sp.PriorityOrder = *req.PriorityOrder
	}
	if req.TargetDate != nil {
		t, err := time.Parse("2006-01-02", *req.TargetDate)
		if err == nil {
			sp.TargetDate = &t
		}
	}
	if req.DevelopmentPlan != nil {
		sp.DevelopmentPlan = *req.DevelopmentPlan
	}
	if req.Notes != nil {
		sp.Notes = *req.Notes
	}
	if err := s.repo.UpdateSuccessionPlan(ctx, sp); err != nil {
		return nil, err
	}
	return successionPlanToResponse(sp), nil
}

func (s *Service) DeleteSuccessionPlan(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteSuccessionPlan(ctx, uid)
}

// GetSuccessionGaps (S-5) menandai posisi kunci tanpa successor siap (succession
// gap). Posisi kunci = posisi yang memiliki ≥1 succession plan ACTIVE. Gap
// terjadi ketika TIDAK ada satupun successor dengan readiness READY_NOW —
// artinya suksesi internal tidak tersedia dan external recruitment menjadi
// fallback. Mengembalikan seluruh posisi kunci + status gap (bukan hanya gap),
// sehingga konsumen bisa melihat coverage suksesi sekaligus kebutuhan fallback.
func (s *Service) GetSuccessionGaps(ctx context.Context) ([]SuccessionGapResponse, error) {
	rows, err := s.repo.ListSuccessionGapPositions(ctx)
	if err != nil {
		return nil, err
	}
	responses := make([]SuccessionGapResponse, 0, len(rows))
	for _, row := range rows {
		hasReady := row.ReadyNowCount > 0
		responses = append(responses, SuccessionGapResponse{
			PositionID:                  row.PositionID,
			PositionTitle:               row.PositionTitle,
			OrganizationID:              row.OrganizationID,
			SuccessorCount:              row.SuccessorCount,
			ReadyNowCount:               row.ReadyNowCount,
			HasReadySuccessor:           hasReady,
			RequiresExternalRecruitment: !hasReady,
		})
	}
	return responses, nil
}

// CheckSuccessionGapByPosition (S-5) mengembalikan true bila posisi adalah
// posisi kunci tanpa successor siap (succession gap) — dipakai Recruitment
// sebagai fallback external recruitment melalui narrow interface
// SuccessionGapProvider (wiring di cmd/server/main.go).
func (s *Service) CheckSuccessionGapByPosition(ctx context.Context, positionID string) (bool, error) {
	uid, err := uuid.Parse(positionID)
	if err != nil {
		return false, fmt.Errorf("invalid position_id: %w", err)
	}
	return s.repo.CheckSuccessionGapByPosition(ctx, uid)
}

// =========================================================================
// Helpers
// =========================================================================

// =========================================================================
// Response helpers
// =========================================================================

func talentMapToResponse(tm *CareerTalentMap) *TalentMapResponse {
	return &TalentMapResponse{
		ID:           tm.ID.String(),
		EmployeeID:   tm.EmployeeID.String(),
		Period:       tm.Period,
		Performance:  tm.Performance,
		Potential:    tm.Potential,
		GridPosition: tm.GridPosition,
		Notes:        tm.Notes,
		AssessorID:   tm.AssessorID.String(),
		AssessedAt:   tm.AssessedAt.Format("2006-01-02"),
		CreatedAt:    tm.CreatedAt,
		UpdatedAt:    tm.UpdatedAt,
	}
}

func careerInterestToResponse(ci *CareerInterest) *CareerInterestResponse {
	return &CareerInterestResponse{
		ID:               ci.ID.String(),
		EmployeeID:       ci.EmployeeID.String(),
		InterestType:     ci.InterestType,
		TargetPosition:   ci.TargetPosition,
		TargetDepartment: ci.TargetDepartment,
		Motivation:       ci.Motivation,
		ReadinessLevel:   ci.ReadinessLevel,
		IsActive:         ci.IsActive,
		RecordedAt:       ci.RecordedAt.Format("2006-01-02"),
		CreatedAt:        ci.CreatedAt,
	}
}

// careerPathToResponse membangun respons dari header + steps terpadu.
// Untuk edge CI (path 2-langkah): source = step 1, target = step terakhir,
// atribut edge (path_type/typical_tenure/requirements/competencies/
// certifications) diambil dari step target. Steps lain (EM jenjang) tetap
// ditampilkan apa adanya.
func (s *Service) careerPathToResponse(ctx context.Context, cp *CareerPath, steps []CareerPathStep) *CareerPathResponse {
	resp := &CareerPathResponse{
		ID:          cp.ID.String(),
		Name:        cp.Name,
		Description: cp.Description,
		IsActive:    cp.IsActive,
		CreatedAt:   cp.CreatedAt,
		UpdatedAt:   cp.UpdatedAt,
		Steps:       make([]CareerPathStepResponse, 0, len(steps)),
	}
	// Batch-resolve position names (pola sama EM enrichCareerPathSteps).
	posIDs := make([]uuid.UUID, 0, len(steps))
	for _, st := range steps {
		posIDs = append(posIDs, st.PositionID)
	}
	names, _ := s.repo.GetPositionNamesByIDs(ctx, posIDs)
	for _, st := range steps {
		stepResp := CareerPathStepResponse{
			ID:                   st.ID.String(),
			PositionID:           st.PositionID.String(),
			PositionName:         names[st.PositionID.String()],
			Sequence:             st.Sequence,
			MinimumServiceMonths: st.MinimumServiceMonths,
			Requirements:         st.Requirements,
			PathType:             st.PathType,
			TypicalTenure:        st.TypicalTenure,
			Competencies:         st.Competencies,
			Certifications:       st.Certifications,
		}
		resp.Steps = append(resp.Steps, stepResp)
	}
	if len(steps) > 0 {
		resp.SourceTitleID = steps[0].PositionID.String()
		target := steps[len(steps)-1]
		resp.TargetTitleID = target.PositionID.String()
		resp.PathType = target.PathType
		if target.TypicalTenure != nil {
			resp.TypicalTenure = *target.TypicalTenure
		}
		resp.Requirements = target.Requirements
		resp.Competencies = target.Competencies
		resp.Certifications = target.Certifications
	}
	return resp
}

func successionPlanToResponse(sp *CareerSuccessionPlan) *SuccessionPlanResponse {
	resp := &SuccessionPlanResponse{
		ID:              sp.ID.String(),
		PositionID:      sp.PositionID.String(),
		SuccessorID:     sp.SuccessorID.String(),
		ReadinessLevel:  sp.ReadinessLevel,
		PriorityOrder:   sp.PriorityOrder,
		DevelopmentPlan: sp.DevelopmentPlan,
		Notes:           sp.Notes,
		Status:          sp.Status,
		CreatedAt:       sp.CreatedAt,
		UpdatedAt:       sp.UpdatedAt,
	}
	if sp.TargetDate != nil {
		t := sp.TargetDate.Format("2006-01-02")
		resp.TargetDate = &t
	}
	return resp
}

func calcTotalPages(total int64, perPage int) int {
	if total == 0 {
		return 0
	}
	return int((total + int64(perPage) - 1) / int64(perPage))
}
