package careerintelligence

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

type Service struct {
	repo   *Repository
	logger *zap.Logger
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
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
// Talent Map
// =========================================================================

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

func (s *Service) CreateCareerPath(ctx context.Context, req CreateCareerPathRequest) (*CareerPathResponse, error) {
	srcID, err := uuid.Parse(req.SourceTitleID)
	if err != nil {
		return nil, fmt.Errorf("invalid source_title_id: %w", err)
	}
	tgtID, err := uuid.Parse(req.TargetTitleID)
	if err != nil {
		return nil, fmt.Errorf("invalid target_title_id: %w", err)
	}

	// Edge CI → path 2-langkah pada skema terpadu: step 1 = source,
	// step 2 = target dengan seluruh atribut edge CI (path_type, typical_tenure,
	// requirements, competencies) disimpan pada step target.
	tenure := req.TypicalTenure
	cp := &CareerPath{
		Name:     s.buildCareerPathName(ctx, req),
		IsActive: true,
	}
	steps := []CareerPathStep{
		{PositionID: srcID, Sequence: 1},
		{
			PositionID:    tgtID,
			Sequence:      2,
			PathType:      req.PathType,
			TypicalTenure: &tenure,
			Requirements:  req.Requirements,
			Competencies:  req.Competencies,
		},
	}
	if err := s.repo.CreateCareerPathTx(ctx, cp, steps); err != nil {
		return nil, err
	}

	s.logger.Info("Career path created",
		zap.String("path_id", cp.ID.String()),
		zap.String("path_type", req.PathType),
	)

	return s.careerPathToResponse(cp, steps), nil
}

// buildCareerPathName membuat nama path unik untuk edge CI. Bila klien tidak
// mengirim name, dipakai "<PATH_TYPE>: <source> → <target>" (fallback ke
// potongan UUID jika nama posisi tidak tersedia). Akhiran unik ditambahkan
// bila nama sudah dipakai (menghormati uk_career_paths_name).
func (s *Service) buildCareerPathName(ctx context.Context, req CreateCareerPathRequest) string {
	if req.Name != "" {
		return req.Name
	}
	srcName, _ := s.repo.GetPositionTitle(ctx, mustUUID(req.SourceTitleID))
	tgtName, _ := s.repo.GetPositionTitle(ctx, mustUUID(req.TargetTitleID))
	base := fmt.Sprintf("%s: %s → %s", req.PathType, srcName, tgtName)
	if base == "" {
		base = fmt.Sprintf("%s-%s-%s", req.PathType, req.SourceTitleID[:8], req.TargetTitleID[:8])
	}
	// Pastikan unik: coba base, lalu base-2, base-3, ...
	candidate := base
	for i := 2; ; i++ {
		if _, err := s.repo.FindCareerPathByName(ctx, candidate); err != nil {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
}

func mustUUID(s string) uuid.UUID {
	uid, _ := uuid.Parse(s)
	return uid
}

func (s *Service) ListCareerPaths(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListCareerPaths(ctx, page, perPage)
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
		responses = append(responses, *s.careerPathToResponse(&cp, steps))
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
func (s *Service) careerPathToResponse(cp *CareerPath, steps []CareerPathStep) *CareerPathResponse {
	resp := &CareerPathResponse{
		ID:         cp.ID.String(),
		Name:       cp.Name,
		IsActive:   cp.IsActive,
		CreatedAt:  cp.CreatedAt,
		UpdatedAt:  cp.UpdatedAt,
		Steps:      make([]CareerPathStepResponse, 0, len(steps)),
	}
	for _, st := range steps {
		stepResp := CareerPathStepResponse{
			ID:             st.ID.String(),
			PositionID:     st.PositionID.String(),
			Sequence:       st.Sequence,
			PathType:       st.PathType,
			TypicalTenure:  st.TypicalTenure,
			Competencies:   st.Competencies,
			Certifications: st.Certifications,
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
