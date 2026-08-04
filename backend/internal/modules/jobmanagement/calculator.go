package jobmanagement

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// =========================================================================
// Job Value Calculator — port dari JobValueCalculator.php (legacy Laravel)
//
// Referensi: docs/job-management-score-analysis.md
// Mesin ini menghitung skor jabatan dari 10 komponen section job management
// (level dari job_management_values), lalu mengagregasi menjadi:
//   - base_score (tanpa komponen keuangan)
//   - job_value_with_financial   (jika has_financial_authority)
//   - job_value_without_financial (jika !has_financial_authority)
//   - sub_component_points (poin per sub-komponen, dipakai untuk is_complete)
//   - components (rincian lengkap tiap komponen, untuk tampilan)
//
// Catatan mapping ke Go (dokumen 8.2/8.6): di Go, slug `code` legacy diangkat
// menjadi kolom `type` pada job_management_values. Nama grup legacy
// (Psychological, Technical, Managerial, ...) tidak ada — jadi kalkulator
// mengelompokkan record potency competencies lewat daftar tetap slug tipe.
// =========================================================================

// ---------------------------------------------------------------------------
// Tabel pemetaan level → poin (dokumen 3.1–3.5)
// ---------------------------------------------------------------------------

// mapDefault — mayoritas komponen (5 level: 1,3,6,10,15)
var mapDefault = map[int]uint64{1: 1, 2: 3, 3: 6, 4: 10, 5: 15}

// mapExtended — komponen 8 level (1,3,6,10,15,21,28,36)
var mapExtended = map[int]uint64{1: 1, 2: 3, 3: 6, 4: 10, 5: 15, 6: 21, 7: 28, 8: 36}

// mapLinear5 — 5 level linier (1..5)
var mapLinear5 = map[int]uint64{1: 1, 2: 2, 3: 3, 4: 4, 5: 5}

// mapLinear8 — 8 level linier (1..8)
var mapLinear8 = map[int]uint64{1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8}

// mapCommunication — komunikasi & influencing (3 level: 1,3,6)
var mapCommunication = map[int]uint64{1: 1, 2: 3, 3: 6}

// ---------------------------------------------------------------------------
// Slug tipe job_management_values per kelompok (dokumen 8.2 & 8.6)
// ---------------------------------------------------------------------------

// psychologicalTypes — Potensi Psikologi (migration 041).
// Kecerdasan (kecerdasan) termasuk potensi psikologi (dokumen 8.3).
var psychologicalTypes = map[string]struct{}{
	"kecerdasan":            {},
	"innovation_creativity": {},
	"self_confidence":       {},
	"flexibility":           {},
	"tenacity":              {},
	"continuous_learning":   {},
}

// technicalTypes — Kompetensi Technical (16 tipe, migration 042)
var technicalTypes = map[string]struct{}{
	"competency_based_human_resources_management": {},
	"competency_development":                      {},
	"people_development":                          {},
	"career_management":                           {},
	"hr_assessment":                               {},
	"recruitement_selection":                      {},
	"job_analysis_evaluation":                     {},
	"organizational_development":                  {},
	"human_resources_information_system":          {},
	"workload_analysis":                           {},
	"performance_apraisal":                        {},
	"remuneration_manajemen":                      {},
	"reward_punisment_management":                 {},
	"health_safety_environment":                   {},
	"hubungan_industrial":                         {},
	"budgeting":                                   {},
}

// managerialTypes — Kompetensi Managerial (6 tipe, migration 042)
var managerialTypes = map[string]struct{}{
	"integrity":               {},
	"achievement_orientation": {},
	"building_partnership":    {},
	"planning_organizing":     {},
	"leadership":              {},
	"developing_others":       {},
}

// Slug tipe khusus (komponen tunggal)
const (
	typeCommunicationSkill  = "communicating_influencing_skill"
	typeThinkingEnvironment = "thinking_environment"
	typeThinkingChallenge   = "thinking_chalenge"
)

// ---------------------------------------------------------------------------
// Struktur hasil komponen (JSON: key persis seperti legacy)
// ---------------------------------------------------------------------------

type educationExperienceComponent struct {
	EducationLevel   *int   `json:"education_level"`
	ExperienceLevel  *int   `json:"experience_level"`
	EducationPoints  uint64 `json:"education_points"`
	ExperiencePoints uint64 `json:"experience_points"`
	Score            uint64 `json:"score"`
}

type potentialsComponent struct {
	AverageLevel *float64 `json:"average_level"`
	Items        []int    `json:"items"`
	Score        uint64   `json:"score"`
}

type competenciesComponent struct {
	TechnicalAverageLevel  *float64 `json:"technical_average_level"`
	ManagerialAverageLevel *float64 `json:"managerial_average_level"`
	CommunicationLevel     *int     `json:"communication_level"`
	TechnicalPoints        uint64   `json:"technical_points"`
	ManagerialPoints       uint64   `json:"managerial_points"`
	CommunicationPoints    uint64   `json:"communication_points"`
	Score                  uint64   `json:"score"`
	// Nilai agregat ditambahkan setelah semua komponen dihitung (legacy):
	BaseScore           uint64 `json:"base_score"`
	PotentialScore      uint64 `json:"potential_score"`
	ProblemSolvingScore uint64 `json:"problem_solving_score"`
}

type problemSolvingComponent struct {
	EnvironmentLevel  *int   `json:"environment_level"`
	ChallengeLevel    *int   `json:"challenge_level"`
	EnvironmentPoints uint64 `json:"environment_points"`
	ChallengePoints   uint64 `json:"challenge_points"`
	Score             uint64 `json:"score"`
}

type financialAuthorityComponent struct {
	HasAuthority    bool   `json:"has_authority"`
	MoneyLevel      *int   `json:"money_level"`
	AuthorityLevel  *int   `json:"authority_level"`
	ImpactLevel     *int   `json:"impact_level"`
	MoneyPoints     uint64 `json:"money_points"`
	AuthorityPoints uint64 `json:"authority_points"`
	ImpactPoints    uint64 `json:"impact_points"`
	Score           uint64 `json:"score"`
	AlternateScore  uint64 `json:"alternate_score"`
}

type assetAuthorityComponent struct {
	AssetValueLevel      *int   `json:"asset_value_level"`
	AssetAuthorityLevel  *int   `json:"asset_authority_level"`
	AssetValuePoints     uint64 `json:"asset_value_points"`
	AssetAuthorityPoints uint64 `json:"asset_authority_points"`
	Score                uint64 `json:"score"`
}

// singleLevelComponent — dipakai subordinate control & working activity
type singleLevelComponent struct {
	Level  *int   `json:"level"`
	Points uint64 `json:"points"`
	Score  uint64 `json:"score"`
}

type workScopeComponent struct {
	ScopeLevel      *int   `json:"scope_level"`
	FrequencyLevel  *int   `json:"frequency_level"`
	ScopePoints     uint64 `json:"scope_points"`
	FrequencyPoints uint64 `json:"frequency_points"`
	Score           uint64 `json:"score"`
}

type workRiskComponent struct {
	EnvironmentLevel  *int   `json:"environment_level"`
	HazardLevel       *int   `json:"hazard_level"`
	EnvironmentPoints uint64 `json:"environment_points"`
	HazardPoints      uint64 `json:"hazard_points"`
	Score             uint64 `json:"score"`
}

type scoreComponents struct {
	EducationExperience educationExperienceComponent `json:"education_experience"`
	Potentials          potentialsComponent          `json:"potentials"`
	Competencies        competenciesComponent        `json:"competencies"`
	ProblemSolving      problemSolvingComponent      `json:"problem_solving"`
	FinancialAuthority  financialAuthorityComponent  `json:"financial_authority"`
	AssetAuthority      assetAuthorityComponent      `json:"asset_authority"`
	SubordinateControl  singleLevelComponent         `json:"subordinate_control"`
	WorkScope           workScopeComponent           `json:"work_scope"`
	WorkActivity        singleLevelComponent         `json:"work_activity"`
	WorkRisk            workRiskComponent            `json:"work_risk"`
}

// subComponentPoints — poin sub-komponen (dokumen 5.3), dipakai is_complete
type subComponentPoints struct {
	Education                 uint64 `json:"education"`
	Experience                uint64 `json:"experience"`
	Potential                 uint64 `json:"potential"`
	CompetencyTechnical       uint64 `json:"competency_technical"`
	CompetencyManagerial      uint64 `json:"competency_managerial"`
	CompetencyCommunication   uint64 `json:"competency_communication"`
	CompetencyTotal           uint64 `json:"competency_total"`
	ProblemSolving            uint64 `json:"problem_solving"`
	FinancialWithAuthority    uint64 `json:"financial_with_authority"`
	FinancialWithoutAuthority uint64 `json:"financial_without_authority"`
	AssetManagement           uint64 `json:"asset_management"`
	SubordinateControl        uint64 `json:"subordinate_control"`
	WorkScope                 uint64 `json:"work_scope"`
	WorkActivity              uint64 `json:"work_activity"`
	WorkRisk                  uint64 `json:"work_risk"`
}

type scoreTotals struct {
	WithFinancial    uint64 `json:"with_financial"`
	WithoutFinancial uint64 `json:"without_financial"`
}

// JobScoreResult — hasil kalkulasi satu organisasi (dokumen 5)
type JobScoreResult struct {
	Components            scoreComponents    `json:"components"`
	Totals                scoreTotals        `json:"totals"`
	HasFinancialAuthority bool               `json:"has_financial_authority"`
	SubComponents         subComponentPoints `json:"sub_components"`
	// IsComplete — status kelengkapan skor (dokumen 5.4). Dipersist ke kolom
	// is_complete & completed_at (migration 051, legacy persistResults).
	IsComplete bool `json:"is_complete"`
}

// ---------------------------------------------------------------------------
// Calculator
// ---------------------------------------------------------------------------

// Calculator membaca data section per organisasi dan menghitung skor.
// Data level diambil dari job_management_values (join manual — tidak
// mengandalkan relasi GORM, sesuai rekomendasi dokumen 8.7.1 opsi b).
type Calculator struct {
	repo *Repository
}

func NewCalculator(repo *Repository) *Calculator {
	return &Calculator{repo: repo}
}

// CalculateForOrganization menghitung skor untuk satu organisasi.
func (c *Calculator) CalculateForOrganization(ctx context.Context, orgID uuid.UUID) (*JobScoreResult, error) {
	results, err := c.CalculateForOrganizationIDs(ctx, []uuid.UUID{orgID})
	if err != nil {
		return nil, err
	}
	if r, ok := results[orgID]; ok {
		return r, nil
	}
	return nil, fmt.Errorf("no job score result for organization %s", orgID)
}

// CalculateForOrganizationIDs menghitung skor untuk banyak organisasi sekaligus
// (mirip calculateForOrganizationIds legacy) — untuk re-kalkulasi massal.
func (c *Calculator) CalculateForOrganizationIDs(ctx context.Context, orgIDs []uuid.UUID) (map[uuid.UUID]*JobScoreResult, error) {
	// Filter & unique
	seen := make(map[uuid.UUID]struct{}, len(orgIDs))
	ids := make([]uuid.UUID, 0, len(orgIDs))
	for _, id := range orgIDs {
		if id == uuid.Nil {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}

	results := make(map[uuid.UUID]*JobScoreResult, len(ids))
	if len(ids) == 0 {
		return results, nil
	}

	data, err := c.loadCalcData(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, orgID := range ids {
		results[orgID] = c.calculateForSingleOrganization(orgID, data)
	}
	return results, nil
}

// ---------------------------------------------------------------------------
// Data loading (batch per organisasi)
// ---------------------------------------------------------------------------

type calcData struct {
	eduExp       map[uuid.UUID]*JobEducationExperience
	potency      map[uuid.UUID][]JobPotencyCompetency
	financials   map[uuid.UUID]*JobFinancial
	assets       map[uuid.UUID]*JobAsset
	subordinates map[uuid.UUID]*JobSubordinateControl
	relations    map[uuid.UUID]*JobRelationship
	activities   map[uuid.UUID]*JobWorkingActivity
	risks        map[uuid.UUID]*JobWorkingRisk
	values       map[uuid.UUID]*JobValue
}

func (c *Calculator) loadCalcData(ctx context.Context, orgIDs []uuid.UUID) (*calcData, error) {
	db, err := c.repo.getDB(ctx)
	if err != nil {
		return nil, err
	}
	idStrings := make([]string, 0, len(orgIDs))
	for _, id := range orgIDs {
		idStrings = append(idStrings, id.String())
	}

	d := &calcData{
		eduExp:       make(map[uuid.UUID]*JobEducationExperience, len(orgIDs)),
		potency:      make(map[uuid.UUID][]JobPotencyCompetency, len(orgIDs)),
		financials:   make(map[uuid.UUID]*JobFinancial, len(orgIDs)),
		assets:       make(map[uuid.UUID]*JobAsset, len(orgIDs)),
		subordinates: make(map[uuid.UUID]*JobSubordinateControl, len(orgIDs)),
		relations:    make(map[uuid.UUID]*JobRelationship, len(orgIDs)),
		activities:   make(map[uuid.UUID]*JobWorkingActivity, len(orgIDs)),
		risks:        make(map[uuid.UUID]*JobWorkingRisk, len(orgIDs)),
		values:       make(map[uuid.UUID]*JobValue),
	}

	var eduExp []JobEducationExperience
	if err := db.Where("organization_id IN ?", idStrings).Find(&eduExp).Error; err != nil {
		return nil, fmt.Errorf("load job education experiences: %w", err)
	}
	for i := range eduExp {
		if eduExp[i].OrganizationID != nil {
			d.eduExp[*eduExp[i].OrganizationID] = &eduExp[i]
		}
	}

	var potency []JobPotencyCompetency
	if err := db.Where("organization_id IN ?", idStrings).Find(&potency).Error; err != nil {
		return nil, fmt.Errorf("load job potency competencies: %w", err)
	}
	for i := range potency {
		if potency[i].OrganizationID != nil {
			orgID := *potency[i].OrganizationID
			d.potency[orgID] = append(d.potency[orgID], potency[i])
		}
	}

	var financials []JobFinancial
	if err := db.Where("organization_id IN ?", idStrings).Find(&financials).Error; err != nil {
		return nil, fmt.Errorf("load job financials: %w", err)
	}
	for i := range financials {
		if financials[i].OrganizationID != nil {
			d.financials[*financials[i].OrganizationID] = &financials[i]
		}
	}

	var assets []JobAsset
	if err := db.Where("organization_id IN ?", idStrings).Find(&assets).Error; err != nil {
		return nil, fmt.Errorf("load job assets: %w", err)
	}
	for i := range assets {
		if assets[i].OrganizationID != nil {
			d.assets[*assets[i].OrganizationID] = &assets[i]
		}
	}

	var subordinates []JobSubordinateControl
	if err := db.Where("organization_id IN ?", idStrings).Find(&subordinates).Error; err != nil {
		return nil, fmt.Errorf("load job subordinate controls: %w", err)
	}
	for i := range subordinates {
		if subordinates[i].OrganizationID != nil {
			d.subordinates[*subordinates[i].OrganizationID] = &subordinates[i]
		}
	}

	var relations []JobRelationship
	if err := db.Where("organization_id IN ?", idStrings).Find(&relations).Error; err != nil {
		return nil, fmt.Errorf("load job relationships: %w", err)
	}
	for i := range relations {
		if relations[i].OrganizationID != nil {
			d.relations[*relations[i].OrganizationID] = &relations[i]
		}
	}

	var activities []JobWorkingActivity
	if err := db.Where("organization_id IN ?", idStrings).Find(&activities).Error; err != nil {
		return nil, fmt.Errorf("load job working activities: %w", err)
	}
	for i := range activities {
		if activities[i].OrganizationID != nil {
			d.activities[*activities[i].OrganizationID] = &activities[i]
		}
	}

	var risks []JobWorkingRisk
	if err := db.Where("organization_id IN ?", idStrings).Find(&risks).Error; err != nil {
		return nil, fmt.Errorf("load job working risks: %w", err)
	}
	for i := range risks {
		if risks[i].OrganizationID != nil {
			d.risks[*risks[i].OrganizationID] = &risks[i]
		}
	}

	// Kumpulkan semua job_management_values yang dirujuk → load sekali (anti N+1)
	valueIDs := make(map[uuid.UUID]struct{})
	addValueID := func(id *uuid.UUID) {
		if id != nil {
			valueIDs[*id] = struct{}{}
		}
	}
	for _, r := range eduExp {
		addValueID(r.EducationID)
		addValueID(r.ExperienceID)
	}
	for _, r := range potency {
		addValueID(r.JobManagementValueID)
	}
	for _, r := range financials {
		addValueID(r.JobManagementValueCashID)
		addValueID(r.JobManagementValueAuthorityID)
		addValueID(r.JobManagementValueImpactID)
	}
	for _, r := range assets {
		addValueID(r.JobManagementValueAssetID)
		addValueID(r.JobManagementValueAuthorityID)
	}
	for _, r := range subordinates {
		addValueID(r.JobManagementValueID)
	}
	for _, r := range relations {
		addValueID(r.JobManagementValueRelationshipID)
		addValueID(r.JobManagementValueFrequencyID)
	}
	for _, r := range activities {
		addValueID(r.JobManagementValueID)
	}
	for _, r := range risks {
		addValueID(r.JobManagementValueEnvironmentID)
		addValueID(r.JobManagementValueHazardID)
	}
	if len(valueIDs) > 0 {
		idList := make([]string, 0, len(valueIDs))
		for id := range valueIDs {
			idList = append(idList, id.String())
		}
		var values []JobValue
		if err := db.Where("id IN ?", idList).Find(&values).Error; err != nil {
			return nil, fmt.Errorf("load job management values: %w", err)
		}
		for i := range values {
			d.values[values[i].ID] = &values[i]
		}
	}

	return d, nil
}

// ---------------------------------------------------------------------------
// Perhitungan per organisasi
// ---------------------------------------------------------------------------

func (c *Calculator) calculateForSingleOrganization(orgID uuid.UUID, d *calcData) *JobScoreResult {
	educationComp := calcEducationExperience(d.eduExp[orgID], d.values)
	potentialComp := calcPotentials(d.potency[orgID], d.values)
	competencyComp := calcCompetencies(d.potency[orgID], d.values)
	problemComp := calcProblemSolving(d.potency[orgID], d.values)
	financialComp := calcFinancialAuthority(d.financials[orgID], d.values)
	assetComp := calcAssetAuthority(d.assets[orgID], d.values)
	subordinateComp := calcSubordinateControl(d.subordinates[orgID], d.values)
	relationshipComp := calcWorkScope(d.relations[orgID], d.values)
	activityComp := calcWorkActivity(d.activities[orgID], d.values)
	riskComp := calcWorkRisk(d.risks[orgID], d.values)

	// kompetensi agregat = base (teknikal×managerial×komunikasi) + potensi + problem solving
	competencyBaseScore := competencyComp.Score
	potentialScore := potentialComp.Score
	problemSolvingScore := problemComp.Score
	competencyAggregate := competencyBaseScore + potentialScore + problemSolvingScore

	competencyComp.BaseScore = competencyBaseScore
	competencyComp.PotentialScore = potentialScore
	competencyComp.ProblemSolvingScore = problemSolvingScore
	competencyComp.Score = competencyAggregate

	sub := subComponentPoints{
		Education:               educationComp.EducationPoints,
		Experience:              educationComp.ExperiencePoints,
		Potential:               potentialScore,
		CompetencyTechnical:     competencyComp.TechnicalPoints,
		CompetencyManagerial:    competencyComp.ManagerialPoints,
		CompetencyCommunication: competencyComp.CommunicationPoints,
		CompetencyTotal:         competencyAggregate,
		ProblemSolving:          problemSolvingScore,
		AssetManagement:         assetComp.Score,
		SubordinateControl:      subordinateComp.Score,
		WorkScope:               relationshipComp.Score,
		WorkActivity:            activityComp.Score,
		WorkRisk:                riskComp.Score,
	}
	if financialComp.HasAuthority {
		sub.FinancialWithAuthority = financialComp.Score
	} else {
		sub.FinancialWithoutAuthority = financialComp.Score
	}

	baseScore := educationComp.Score +
		competencyAggregate +
		assetComp.Score +
		subordinateComp.Score +
		relationshipComp.Score +
		activityComp.Score +
		riskComp.Score

	totals := scoreTotals{}
	if financialComp.HasAuthority {
		totals = scoreTotals{WithFinancial: baseScore + financialComp.Score, WithoutFinancial: 0}
	} else {
		totals = scoreTotals{WithFinancial: 0, WithoutFinancial: baseScore + financialComp.Score}
	}

	return &JobScoreResult{
		Components: scoreComponents{
			EducationExperience: educationComp,
			Potentials:          potentialComp,
			Competencies:        competencyComp,
			ProblemSolving:      problemComp,
			FinancialAuthority:  financialComp,
			AssetAuthority:      assetComp,
			SubordinateControl:  subordinateComp,
			WorkScope:           relationshipComp,
			WorkActivity:        activityComp,
			WorkRisk:            riskComp,
		},
		Totals:                totals,
		HasFinancialAuthority: financialComp.HasAuthority,
		SubComponents:         sub,
		IsComplete:            isResultComplete(sub),
	}
}

// ---------------------------------------------------------------------------
// Komponen 4.1–4.10
// ---------------------------------------------------------------------------

// levelOf mengambil level job_management_values berdasarkan ID yang dirujuk.
func levelOf(values map[uuid.UUID]*JobValue, id *uuid.UUID) *int {
	if id == nil {
		return nil
	}
	if v, ok := values[*id]; ok && v != nil {
		return v.Level
	}
	return nil
}

// calcEducationExperience — 4.1 Pendidikan × Pengalaman (MAP_DEFAULT)
func calcEducationExperience(rec *JobEducationExperience, values map[uuid.UUID]*JobValue) educationExperienceComponent {
	var comp educationExperienceComponent
	if rec == nil {
		return comp
	}
	comp.EducationLevel = levelOf(values, rec.EducationID)
	comp.ExperienceLevel = levelOf(values, rec.ExperienceID)
	comp.EducationPoints = mapPoints(mapDefault, comp.EducationLevel)
	comp.ExperiencePoints = mapPoints(mapDefault, comp.ExperienceLevel)
	if comp.EducationPoints > 0 && comp.ExperiencePoints > 0 {
		comp.Score = comp.EducationPoints * comp.ExperiencePoints
	}
	return comp
}

// calcPotentials — 4.2 Potensi Psikologi (rata-rata ceil, MAP_DEFAULT)
func calcPotentials(records []JobPotencyCompetency, values map[uuid.UUID]*JobValue) potentialsComponent {
	var comp potentialsComponent
	var levels []int
	for _, r := range records {
		if r.JobManagementValueID == nil {
			continue
		}
		v, ok := values[*r.JobManagementValueID]
		if !ok || v == nil {
			continue
		}
		if _, isPsych := psychologicalTypes[v.Type]; !isPsych {
			continue
		}
		if v.Level != nil && *v.Level > 0 {
			levels = append(levels, *v.Level)
			comp.Items = append(comp.Items, *v.Level)
		}
	}
	if avg := ceilAverage(levels); avg != nil {
		comp.AverageLevel = avg
		l := int(*avg)
		comp.Score = mapPoints(mapDefault, &l)
	}
	return comp
}

// calcCompetencies — 4.3 Technical × Managerial × Communication
//
// Catatan deviasi vs legacy (dokumen 8.7.3): legacy punya fallback prioritas #2
// via competency_id → competencies.field (Technical Competency / Manajerial) +
// competency_values.level. Di Go, JobPotencyCompetency TIDAK punya FK ke
// competency_values, jadi fallback itu tidak dapat diimplementasikan — level
// hanya dari job_management_value_id (prioritas #1). Aman secara praktis karena
// frontend selalu mengirim job_management_value_id (termasuk Kecerdasan yang
// tanpa competency_id, dokumen 8.3).
func calcCompetencies(records []JobPotencyCompetency, values map[uuid.UUID]*JobValue) competenciesComponent {
	var comp competenciesComponent
	var technicalLevels, managerialLevels []int
	var communicationLevel *int

	for _, r := range records {
		if r.JobManagementValueID == nil {
			continue
		}
		v, ok := values[*r.JobManagementValueID]
		if !ok || v == nil || v.Level == nil || *v.Level <= 0 {
			continue
		}
		switch {
		case containsType(technicalTypes, v.Type):
			technicalLevels = append(technicalLevels, *v.Level)
		case containsType(managerialTypes, v.Type):
			managerialLevels = append(managerialLevels, *v.Level)
		case v.Type == typeCommunicationSkill && communicationLevel == nil:
			communicationLevel = v.Level
		}
	}

	if avg := ceilAverage(technicalLevels); avg != nil {
		comp.TechnicalAverageLevel = avg
		l := int(*avg)
		comp.TechnicalPoints = mapPoints(mapExtended, &l)
	}
	if avg := ceilAverage(managerialLevels); avg != nil {
		comp.ManagerialAverageLevel = avg
		l := int(*avg)
		comp.ManagerialPoints = mapPoints(mapDefault, &l)
	}
	// Communication: default level 1 bila tidak ditemukan (legacy)
	commLevel := communicationLevel
	if commLevel == nil {
		one := 1
		commLevel = &one
	}
	comp.CommunicationLevel = communicationLevel
	comp.CommunicationPoints = mapPoints(mapCommunication, commLevel)

	if comp.TechnicalPoints > 0 && comp.ManagerialPoints > 0 && comp.CommunicationPoints > 0 {
		comp.Score = comp.TechnicalPoints * comp.ManagerialPoints * comp.CommunicationPoints
	}
	return comp
}

// calcProblemSolving — 4.4 Environment × Challenge
func calcProblemSolving(records []JobPotencyCompetency, values map[uuid.UUID]*JobValue) problemSolvingComponent {
	var comp problemSolvingComponent
	for _, r := range records {
		if r.JobManagementValueID == nil {
			continue
		}
		v, ok := values[*r.JobManagementValueID]
		if !ok || v == nil {
			continue
		}
		switch v.Type {
		case typeThinkingEnvironment:
			if comp.EnvironmentLevel == nil {
				comp.EnvironmentLevel = v.Level
			}
		case typeThinkingChallenge:
			if comp.ChallengeLevel == nil {
				comp.ChallengeLevel = v.Level
			}
		}
	}
	comp.EnvironmentPoints = mapPoints(mapExtended, comp.EnvironmentLevel)
	comp.ChallengePoints = mapPoints(mapDefault, comp.ChallengeLevel)
	if comp.EnvironmentPoints > 0 && comp.ChallengePoints > 0 {
		comp.Score = comp.EnvironmentPoints * comp.ChallengePoints
	}
	return comp
}

// calcFinancialAuthority — 4.5 Kewenangan Keuangan
func calcFinancialAuthority(rec *JobFinancial, values map[uuid.UUID]*JobValue) financialAuthorityComponent {
	var comp financialAuthorityComponent
	if rec == nil {
		return comp
	}
	comp.HasAuthority = rec.IsAuthorized
	comp.MoneyLevel = levelOf(values, rec.JobManagementValueCashID)
	comp.AuthorityLevel = levelOf(values, rec.JobManagementValueAuthorityID)
	comp.ImpactLevel = levelOf(values, rec.JobManagementValueImpactID)

	if comp.HasAuthority {
		comp.MoneyPoints = mapPoints(mapExtended, comp.MoneyLevel)
	}
	comp.AuthorityPoints = mapPoints(mapExtended, comp.AuthorityLevel)
	comp.ImpactPoints = mapPoints(mapExtended, comp.ImpactLevel)

	if comp.HasAuthority {
		if comp.MoneyPoints > 0 && comp.AuthorityPoints > 0 && comp.ImpactPoints > 0 {
			comp.Score = comp.MoneyPoints * comp.AuthorityPoints * comp.ImpactPoints
		}
	} else {
		if comp.AuthorityPoints > 0 && comp.ImpactPoints > 0 {
			comp.Score = comp.AuthorityPoints * comp.ImpactPoints
		}
	}
	comp.AlternateScore = comp.Score
	return comp
}

// calcAssetAuthority — 4.6 Kewenangan Aset (LINEAR_8 × DEFAULT)
func calcAssetAuthority(rec *JobAsset, values map[uuid.UUID]*JobValue) assetAuthorityComponent {
	var comp assetAuthorityComponent
	if rec == nil {
		return comp
	}
	comp.AssetValueLevel = levelOf(values, rec.JobManagementValueAssetID)
	comp.AssetAuthorityLevel = levelOf(values, rec.JobManagementValueAuthorityID)
	comp.AssetValuePoints = mapPoints(mapLinear8, comp.AssetValueLevel)
	comp.AssetAuthorityPoints = mapPoints(mapDefault, comp.AssetAuthorityLevel)
	if comp.AssetValuePoints > 0 && comp.AssetAuthorityPoints > 0 {
		comp.Score = comp.AssetValuePoints * comp.AssetAuthorityPoints
	}
	return comp
}

// calcSubordinateControl — 4.7 Kendali Bawahan (MAP_DEFAULT)
func calcSubordinateControl(rec *JobSubordinateControl, values map[uuid.UUID]*JobValue) singleLevelComponent {
	var comp singleLevelComponent
	if rec == nil {
		return comp
	}
	comp.Level = levelOf(values, rec.JobManagementValueID)
	comp.Points = mapPoints(mapDefault, comp.Level)
	comp.Score = comp.Points
	return comp
}

// calcWorkScope — 4.8 Ruang Lingkup Kerja (DEFAULT × LINEAR_5)
func calcWorkScope(rec *JobRelationship, values map[uuid.UUID]*JobValue) workScopeComponent {
	var comp workScopeComponent
	if rec == nil {
		return comp
	}
	comp.ScopeLevel = levelOf(values, rec.JobManagementValueRelationshipID)
	comp.FrequencyLevel = levelOf(values, rec.JobManagementValueFrequencyID)
	comp.ScopePoints = mapPoints(mapDefault, comp.ScopeLevel)
	comp.FrequencyPoints = mapPoints(mapLinear5, comp.FrequencyLevel)
	if comp.ScopePoints > 0 && comp.FrequencyPoints > 0 {
		comp.Score = comp.ScopePoints * comp.FrequencyPoints
	}
	return comp
}

// calcWorkActivity — 4.9 Aktivitas Kerja (MAP_DEFAULT)
func calcWorkActivity(rec *JobWorkingActivity, values map[uuid.UUID]*JobValue) singleLevelComponent {
	var comp singleLevelComponent
	if rec == nil {
		return comp
	}
	comp.Level = levelOf(values, rec.JobManagementValueID)
	comp.Points = mapPoints(mapDefault, comp.Level)
	comp.Score = comp.Points
	return comp
}

// calcWorkRisk — 4.10 Risiko Kerja (LINEAR_5 × LINEAR_5)
func calcWorkRisk(rec *JobWorkingRisk, values map[uuid.UUID]*JobValue) workRiskComponent {
	var comp workRiskComponent
	if rec == nil {
		return comp
	}
	comp.EnvironmentLevel = levelOf(values, rec.JobManagementValueEnvironmentID)
	comp.HazardLevel = levelOf(values, rec.JobManagementValueHazardID)
	comp.EnvironmentPoints = mapPoints(mapLinear5, comp.EnvironmentLevel)
	comp.HazardPoints = mapPoints(mapLinear5, comp.HazardLevel)
	if comp.EnvironmentPoints > 0 && comp.HazardPoints > 0 {
		comp.Score = comp.EnvironmentPoints * comp.HazardPoints
	}
	return comp
}

// ---------------------------------------------------------------------------
// Helper
// ---------------------------------------------------------------------------

// mapPoints memetakan level → poin. Level di luar tabel dipatok ke level
// tertinggi yang tersedia (sama seperti legacy mapPoints).
func mapPoints(m map[int]uint64, level *int) uint64 {
	if level == nil || *level <= 0 {
		return 0
	}
	l := *level
	if _, ok := m[l]; !ok {
		maxKey := 0
		for k := range m {
			if k > maxKey {
				maxKey = k
			}
		}
		l = maxKey
	}
	return m[l]
}

// ceilAverage menghitung ceil(rata-rata) level; nil bila tidak ada level.
func ceilAverage(levels []int) *float64 {
	if len(levels) == 0 {
		return nil
	}
	sum := 0
	for _, l := range levels {
		sum += l
	}
	avg := (sum + len(levels) - 1) / len(levels)
	f := float64(avg)
	return &f
}

func containsType(set map[string]struct{}, t string) bool {
	_, ok := set[t]
	return ok
}

// isResultComplete — status kelengkapan (dokumen 5.4): semua sub-komponen > 0
// kecuali pasangan finansial (cukup salah satu yang > 0).
func isResultComplete(sub subComponentPoints) bool {
	required := []uint64{
		sub.Education,
		sub.Experience,
		sub.Potential,
		sub.CompetencyTechnical,
		sub.CompetencyManagerial,
		sub.ProblemSolving,
		sub.AssetManagement,
		sub.SubordinateControl,
		sub.WorkScope,
		sub.WorkActivity,
		sub.WorkRisk,
	}
	for _, v := range required {
		if v <= 0 {
			return false
		}
	}
	return sub.FinancialWithAuthority > 0 || sub.FinancialWithoutAuthority > 0
}
