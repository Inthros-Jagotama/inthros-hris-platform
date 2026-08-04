package jobmanagement

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =========================================================================
// Calculator tests — port dari contoh perhitungan dokumen
// docs/job-management-score-analysis.md section 10
// =========================================================================

func createValue(t *testing.T, db *gorm.DB, valueType string, level int) *JobValue {
	t.Helper()
	l := level
	v := &JobValue{Type: valueType, Level: &l}
	if err := db.Create(v).Error; err != nil {
		t.Fatalf("create job value (%s level %d): %v", valueType, level, err)
	}
	return v
}

// seedCalcScenario membuat satu organisasi lengkap dengan data section
// sesuai contoh dokumen section 10, kecuali financial (sesuai hasAuthority).
func seedCalcScenario(t *testing.T, db *gorm.DB, hasAuthority bool) uuid.UUID {
	t.Helper()
	orgID := uuid.New()

	// Data referensi (level sesuai contoh dokumen)
	edu := createValue(t, db, "education", 4)                                    // MAP_DEFAULT → 10
	exp := createValue(t, db, "experience", 3)                                   // MAP_DEFAULT → 6
	kec := createValue(t, db, "kecerdasan", 4)                                   // potensi psikologi
	inno := createValue(t, db, "innovation_creativity", 4)                       // potensi psikologi
	tech := createValue(t, db, "competency_based_human_resources_management", 5) // technical
	mgr := createValue(t, db, "integrity", 3)                                    // managerial
	comm := createValue(t, db, "communicating_influencing_skill", 2)             // → 3 poin
	env := createValue(t, db, "thinking_environment", 4)                         // → 10 poin
	chg := createValue(t, db, "thinking_chalenge", 3)                            // → 6 poin
	cash := createValue(t, db, "cash", 4)                                        // EXTENDED → 10
	auth := createValue(t, db, "authority", 5)                                   // EXTENDED → 15
	impact := createValue(t, db, "impact", 4)                                    // EXTENDED → 10
	authUnauthorized := createValue(t, db, "authority_unauthorized", 4)          // EXTENDED → 10
	impactUnauthorized := createValue(t, db, "impact_unauthorized", 3)           // EXTENDED → 6
	assetV := createValue(t, db, "asset", 3)                                     // LINEAR_8 → 3
	assetAuth := createValue(t, db, "asset_authority", 4)                        // DEFAULT → 10
	sub := createValue(t, db, "subordinate", 4)                                  // DEFAULT → 10
	rel := createValue(t, db, "relationship", 4)                                 // DEFAULT → 10
	freq := createValue(t, db, "frequency", 3)                                   // LINEAR_5 → 3
	act := createValue(t, db, "activity", 3)                                     // DEFAULT → 6
	workEnv := createValue(t, db, "environment", 3)                              // LINEAR_5 → 3
	workRisk := createValue(t, db, "risk", 3)                                    // LINEAR_5 → 3

	// 9.7 Pendidikan & Pengalaman
	eduRec := &JobEducationExperience{
		OrganizationID: &orgID, Nomenclature: "N", FullCode: "FC",
		EducationID: &edu.ID, ExperienceID: &exp.ID,
	}
	if err := db.Create(eduRec).Error; err != nil {
		t.Fatalf("create education experience: %v", err)
	}

	// 9.16 Kompetensi Potensi (kecerdasan tanpa competency_id — dokumen 8.3)
	potencyRows := []struct {
		value *JobValue
	}{
		{kec}, {inno}, {tech}, {mgr}, {comm}, {env}, {chg},
	}
	for _, row := range potencyRows {
		rec := &JobPotencyCompetency{
			OrganizationID:       &orgID,
			JobManagementValueID: &row.value.ID,
		}
		if err := db.Create(rec).Error; err != nil {
			t.Fatalf("create potency competency: %v", err)
		}
	}

	// 9.15 Keuangan
	fin := &JobFinancial{
		OrganizationID: &orgID, Nomenclature: "N", FullCode: "FC",
		IsAuthorized: hasAuthority,
	}
	if hasAuthority {
		fin.JobManagementValueCashID = &cash.ID
		fin.JobManagementValueAuthorityID = &auth.ID
		fin.JobManagementValueImpactID = &impact.ID
	} else {
		fin.JobManagementValueAuthorityID = &authUnauthorized.ID
		fin.JobManagementValueImpactID = &impactUnauthorized.ID
	}
	if err := db.Create(fin).Error; err != nil {
		t.Fatalf("create financial: %v", err)
	}

	// 9.14 Aset
	if err := db.Create(&JobAsset{
		OrganizationID:                &orgID,
		JobManagementValueAssetID:     &assetV.ID,
		JobManagementValueAuthorityID: &assetAuth.ID,
	}).Error; err != nil {
		t.Fatalf("create asset: %v", err)
	}

	// 9.13 Bawahan
	if err := db.Create(&JobSubordinateControl{
		OrganizationID:       &orgID,
		JobManagementValueID: &sub.ID,
	}).Error; err != nil {
		t.Fatalf("create subordinate: %v", err)
	}

	// 9.12 Hubungan Kerja
	if err := db.Create(&JobRelationship{
		OrganizationID:                   &orgID,
		JobManagementValueRelationshipID: &rel.ID,
		JobManagementValueFrequencyID:    &freq.ID,
	}).Error; err != nil {
		t.Fatalf("create relationship: %v", err)
	}

	// 9.10 Aktivitas Kerja
	if err := db.Create(&JobWorkingActivity{
		OrganizationID:       &orgID,
		JobManagementValueID: &act.ID,
	}).Error; err != nil {
		t.Fatalf("create activity: %v", err)
	}

	// 9.11 Risiko Kerja
	if err := db.Create(&JobWorkingRisk{
		OrganizationID:                  &orgID,
		JobManagementValueEnvironmentID: &workEnv.ID,
		JobManagementValueHazardID:      &workRisk.ID,
	}).Error; err != nil {
		t.Fatalf("create risk: %v", err)
	}

	return orgID
}

func TestCalculator_WithFinancial(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	calc := NewCalculator(repo)
	ctx := context.Background()

	orgID := seedCalcScenario(t, db, true)
	result, err := calc.CalculateForOrganization(ctx, orgID)
	if err != nil {
		t.Fatalf("CalculateForOrganization failed: %v", err)
	}

	// Contoh dokumen section 10: base = 485, with_financial = 1985
	if result.Totals.WithFinancial != 1985 {
		t.Errorf("expected with_financial=1985, got %d", result.Totals.WithFinancial)
	}
	if result.Totals.WithoutFinancial != 0 {
		t.Errorf("expected without_financial=0 (has authority), got %d", result.Totals.WithoutFinancial)
	}
	if !result.HasFinancialAuthority {
		t.Error("expected has_financial_authority=true")
	}
	if !result.IsComplete {
		t.Error("expected is_complete=true (semua komponen terisi)")
	}

	// Spot-check sub-componen
	sub := result.SubComponents
	if sub.Education != 10 || sub.Experience != 6 {
		t.Errorf("unexpected education/experience points: %d/%d", sub.Education, sub.Experience)
	}
	if sub.Potential != 10 {
		t.Errorf("expected potential=10, got %d", sub.Potential)
	}
	if sub.CompetencyTechnical != 15 || sub.CompetencyManagerial != 6 || sub.CompetencyCommunication != 3 {
		t.Errorf("unexpected competency points: %d/%d/%d",
			sub.CompetencyTechnical, sub.CompetencyManagerial, sub.CompetencyCommunication)
	}
	if sub.CompetencyTotal != 340 {
		t.Errorf("expected competency_total=340, got %d", sub.CompetencyTotal)
	}
	if sub.ProblemSolving != 60 {
		t.Errorf("expected problem_solving=60, got %d", sub.ProblemSolving)
	}
	if sub.AssetManagement != 30 {
		t.Errorf("expected asset_management=30, got %d", sub.AssetManagement)
	}
	if sub.SubordinateControl != 10 {
		t.Errorf("expected subordinate_control=10, got %d", sub.SubordinateControl)
	}
	if sub.WorkScope != 30 {
		t.Errorf("expected work_scope=30, got %d", sub.WorkScope)
	}
	if sub.WorkActivity != 6 {
		t.Errorf("expected work_activity=6, got %d", sub.WorkActivity)
	}
	if sub.WorkRisk != 9 {
		t.Errorf("expected work_risk=9, got %d", sub.WorkRisk)
	}
	if sub.FinancialWithAuthority != 1500 || sub.FinancialWithoutAuthority != 0 {
		t.Errorf("unexpected financial sub points: %d/%d",
			sub.FinancialWithAuthority, sub.FinancialWithoutAuthority)
	}
}

func TestCalculator_WithoutFinancial(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	calc := NewCalculator(repo)
	ctx := context.Background()

	orgID := seedCalcScenario(t, db, false)
	result, err := calc.CalculateForOrganization(ctx, orgID)
	if err != nil {
		t.Fatalf("CalculateForOrganization failed: %v", err)
	}

	// financial = EXTENDED(4)=10 × EXTENDED(3)=6 = 60; base = 485
	if result.Totals.WithFinancial != 0 {
		t.Errorf("expected with_financial=0 (no authority), got %d", result.Totals.WithFinancial)
	}
	if result.Totals.WithoutFinancial != 545 {
		t.Errorf("expected without_financial=545, got %d", result.Totals.WithoutFinancial)
	}
	if result.HasFinancialAuthority {
		t.Error("expected has_financial_authority=false")
	}
	sub := result.SubComponents
	if sub.FinancialWithoutAuthority != 60 || sub.FinancialWithAuthority != 0 {
		t.Errorf("unexpected financial sub points: %d/%d",
			sub.FinancialWithAuthority, sub.FinancialWithoutAuthority)
	}
}

func TestCalculator_EmptyOrg(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	calc := NewCalculator(repo)
	ctx := context.Background()

	orgID := uuid.New()
	result, err := calc.CalculateForOrganization(ctx, orgID)
	if err != nil {
		t.Fatalf("CalculateForOrganization failed: %v", err)
	}
	if result.Totals.WithFinancial != 0 || result.Totals.WithoutFinancial != 0 {
		t.Errorf("expected zero totals for empty org, got %d/%d",
			result.Totals.WithFinancial, result.Totals.WithoutFinancial)
	}
	if result.IsComplete {
		t.Error("expected is_complete=false for empty org")
	}
}

func TestCalculator_CommunicationDefault(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	calc := NewCalculator(repo)
	ctx := context.Background()

	orgID := seedCalcScenario(t, db, true)
	// Hapus record komunikasi → level default 1 (1 poin)
	if err := db.Where("job_management_value_id IN (?)",
		db.Model(&JobValue{}).Select("id").Where("type = ?", "communicating_influencing_skill")).
		Delete(&JobPotencyCompetency{}).Error; err != nil {
		t.Fatalf("delete communication record: %v", err)
	}

	result, err := calc.CalculateForOrganization(ctx, orgID)
	if err != nil {
		t.Fatalf("CalculateForOrganization failed: %v", err)
	}
	if result.Components.Competencies.CommunicationPoints != 1 {
		t.Errorf("expected communication_points=1 (default), got %d",
			result.Components.Competencies.CommunicationPoints)
	}
	// base baru: 60 + (15*6*1 + 10 + 60) + 30 + 10 + 30 + 6 + 9 = 305
	if result.Totals.WithFinancial != 1805 {
		t.Errorf("expected with_financial=1805, got %d", result.Totals.WithFinancial)
	}
}

// =========================================================================
// Integrasi service — kontrak recalculate (dokumen 9 langkah 3)
// =========================================================================

// TestService_UpsertJobScore_EmptyBody_Recalculates — PUT body kosong
// ({ components: null } dari tombol Recalculate FE) harus memicu kalkulasi
// server-side dan menyimpan hasilnya.
func TestService_UpsertJobScore_EmptyBody_Recalculates(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	svc := NewService(NewRepository(dbResolver), nil)
	ctx := context.Background()

	orgID := seedCalcScenario(t, db, true)

	// Body kosong → recalculate
	resp, err := svc.UpsertJobScore(ctx, orgID.String(), UpdateJobScoreRequest{})
	if err != nil {
		t.Fatalf("UpsertJobScore (empty body) failed: %v", err)
	}
	if resp.JobValueWithFinancial != 1985 {
		t.Errorf("expected recalculated job_value_with_financial=1985, got %d", resp.JobValueWithFinancial)
	}
	if resp.CalculatedAt == nil {
		t.Error("expected calculated_at to be set after recalculation")
	}
	if resp.SubComponentPoints == "" || resp.Components == "" {
		t.Error("expected components & sub_component_points to be persisted")
	}
	// Skenario lengkap → is_complete=true & completed_at terisi
	if !resp.IsComplete || resp.CompletedAt == nil {
		t.Error("expected is_complete=true and completed_at set for complete scenario")
	}

	// Verifikasi tersimpan di DB
	found, err := svc.repo.FindJobScoreByOrganizationID(ctx, orgID)
	if err != nil {
		t.Fatalf("find job score failed: %v", err)
	}
	if found.JobValueWithFinancial != 1985 {
		t.Errorf("expected stored job_value_with_financial=1985, got %d", found.JobValueWithFinancial)
	}
	if !found.IsComplete || found.CompletedAt == nil {
		t.Error("expected stored is_complete=true and completed_at set")
	}
}

// TestService_UpsertJobScore_EmptyBody_WithManualValuesKeepsLegacy — PUT dengan
// nilai eksplisit tetap memakai jalur lama (tidak di-recalculate).
func TestService_UpsertJobScore_WithValuesNoRecalc(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	svc := NewService(NewRepository(dbResolver), nil)
	ctx := context.Background()

	orgID := createTestOrgID()
	resp, err := svc.UpsertJobScore(ctx, orgID, UpdateJobScoreRequest{
		JobValueWithFinancial: uint64Ptr(777),
	})
	if err != nil {
		t.Fatalf("UpsertJobScore failed: %v", err)
	}
	if resp.JobValueWithFinancial != 777 {
		t.Errorf("expected 777 (manual path preserved), got %d", resp.JobValueWithFinancial)
	}
}

// TestService_RecalculateJobScores — re-kalkulasi massal (batch).
func TestService_RecalculateJobScores(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	svc := NewService(NewRepository(dbResolver), nil)
	ctx := context.Background()

	orgA := seedCalcScenario(t, db, true)
	orgB := seedCalcScenario(t, db, false)

	responses, err := svc.RecalculateJobScores(ctx, []string{orgA.String(), orgB.String()})
	if err != nil {
		t.Fatalf("RecalculateJobScores failed: %v", err)
	}
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	if responses[orgA.String()].JobValueWithFinancial != 1985 {
		t.Errorf("org A: expected with_financial=1985, got %d", responses[orgA.String()].JobValueWithFinancial)
	}
	if responses[orgB.String()].JobValueWithoutFinancial != 545 {
		t.Errorf("org B: expected without_financial=545, got %d", responses[orgB.String()].JobValueWithoutFinancial)
	}
}

// =========================================================================
// Hook otomatis — simpan/ubah/hapus section → kalkulator dijalankan
// =========================================================================

func findValueIDByType(t *testing.T, db *gorm.DB, valueType string) string {
	t.Helper()
	var v JobValue
	if err := db.Where("type = ?", valueType).First(&v).Error; err != nil {
		t.Fatalf("find job value type %s: %v", valueType, err)
	}
	return v.ID.String()
}

func TestService_SectionCreate_TriggersRecalc(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	svc := NewService(NewRepository(dbResolver), nil)
	ctx := context.Background()

	orgID := seedCalcScenario(t, db, true)
	// Hapus financial langsung (simulasi org belum punya keuangan, tanpa hook)
	if err := db.Where("organization_id = ?", orgID.String()).Delete(&JobFinancial{}).Error; err != nil {
		t.Fatalf("delete financial: %v", err)
	}
	if _, err := svc.repo.FindJobScoreByOrganizationID(ctx, orgID); err == nil {
		t.Fatal("expected no score yet before creating financial")
	}

	cashID := findValueIDByType(t, db, "cash")
	authID := findValueIDByType(t, db, "authority")
	impactID := findValueIDByType(t, db, "impact")

	if _, err := svc.CreateJobFinancial(ctx, CreateJobFinancialRequest{
		OrganizationID:                orgID.String(),
		IsAuthorized:                  true,
		JobManagementValueCashID:      &cashID,
		JobManagementValueAuthorityID: &authID,
		JobManagementValueImpactID:    &impactID,
	}); err != nil {
		t.Fatalf("CreateJobFinancial failed: %v", err)
	}

	found, err := svc.repo.FindJobScoreByOrganizationID(ctx, orgID)
	if err != nil {
		t.Fatalf("expected score after section create, got: %v", err)
	}
	if found.JobValueWithFinancial != 1985 {
		t.Errorf("expected 1985 after creating financial, got %d", found.JobValueWithFinancial)
	}
}

func TestService_SectionUpdate_TriggersRecalc(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	svc := NewService(NewRepository(dbResolver), nil)
	ctx := context.Background()

	orgID := seedCalcScenario(t, db, true)

	var fin JobFinancial
	if err := db.Where("organization_id = ?", orgID.String()).First(&fin).Error; err != nil {
		t.Fatalf("find financial: %v", err)
	}
	// Mirip perilaku FE: saat unauthorized, field diisi value tipe
	// authority_unauthorized / impact_unauthorized dan cash dikosongkan.
	isAuthorized := false
	authU := findValueIDByType(t, db, "authority_unauthorized") // level 4 → 10 poin
	impactU := findValueIDByType(t, db, "impact_unauthorized")  // level 3 → 6 poin
	empty := ""
	if _, err := svc.UpdateJobFinancial(ctx, fin.ID.String(), UpdateJobFinancialRequest{
		IsAuthorized:                  &isAuthorized,
		JobManagementValueCashID:      &empty,
		JobManagementValueAuthorityID: &authU,
		JobManagementValueImpactID:    &impactU,
	}); err != nil {
		t.Fatalf("UpdateJobFinancial failed: %v", err)
	}

	found, err := svc.repo.FindJobScoreByOrganizationID(ctx, orgID)
	if err != nil {
		t.Fatalf("find score: %v", err)
	}
	// financial = 10 × 6 = 60 → without_financial = 485 + 60 = 545
	if found.JobValueWithFinancial != 0 || found.JobValueWithoutFinancial != 545 {
		t.Errorf("expected 0/545 after update (no authority), got %d/%d",
			found.JobValueWithFinancial, found.JobValueWithoutFinancial)
	}
}

func TestService_SectionDelete_TriggersRecalc(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	svc := NewService(NewRepository(dbResolver), nil)
	ctx := context.Background()

	orgID := seedCalcScenario(t, db, true)

	// Hapus kedua record potensi psikologi (kecerdasan & innovation_creativity)
	var pots []JobPotencyCompetency
	if err := db.Where("organization_id = ?", orgID.String()).Find(&pots).Error; err != nil {
		t.Fatalf("find potency: %v", err)
	}
	deleted := 0
	for _, p := range pots {
		if p.JobManagementValueID == nil {
			continue
		}
		var v JobValue
		if err := db.Where("id = ?", p.JobManagementValueID.String()).First(&v).Error; err == nil {
			if _, ok := psychologicalTypes[v.Type]; ok {
				if err := svc.DeleteJobPotencyCompetency(ctx, p.ID.String()); err != nil {
					t.Fatalf("delete potency: %v", err)
				}
				deleted++
			}
		}
	}
	if deleted < 2 {
		t.Fatalf("expected to delete 2 psychological records, deleted %d", deleted)
	}

	found, err := svc.repo.FindJobScoreByOrganizationID(ctx, orgID)
	if err != nil {
		t.Fatalf("find score: %v", err)
	}
	// potensi=0 → aggregate = 270+0+60 = 330; base = 60+330+30+10+30+6+9 = 475;
	// with_financial = 475 + 1500 = 1975
	if found.JobValueWithFinancial != 1975 {
		t.Errorf("expected 1975 after deleting potentials, got %d", found.JobValueWithFinancial)
	}
	// Potensi kosong → skor tidak lengkap → is_complete=false, completed_at nil
	if found.IsComplete {
		t.Error("expected is_complete=false after deleting all potentials")
	}
	if found.CompletedAt != nil {
		t.Error("expected completed_at nil when score is incomplete")
	}
}

// =========================================================================
// Unit test helper murni
// =========================================================================

func TestMapPoints(t *testing.T) {
	if got := mapPoints(mapDefault, intPtr(3)); got != 6 {
		t.Errorf("mapDefault(3) = %d, want 6", got)
	}
	if got := mapPoints(mapDefault, nil); got != 0 {
		t.Errorf("mapDefault(nil) = %d, want 0", got)
	}
	if got := mapPoints(mapDefault, intPtr(0)); got != 0 {
		t.Errorf("mapDefault(0) = %d, want 0", got)
	}
	// Level di luar tabel dipatok ke level tertinggi (legacy)
	if got := mapPoints(mapDefault, intPtr(99)); got != 15 {
		t.Errorf("mapDefault(99) = %d, want 15 (clamped)", got)
	}
	if got := mapPoints(mapExtended, intPtr(8)); got != 36 {
		t.Errorf("mapExtended(8) = %d, want 36", got)
	}
}

func TestCeilAverage(t *testing.T) {
	if avg := ceilAverage([]int{4, 5}); avg == nil || *avg != 5 {
		t.Errorf("ceilAverage([4,5]) = %v, want 5", avg)
	}
	if avg := ceilAverage([]int{3, 4, 4}); avg == nil || *avg != 4 {
		t.Errorf("ceilAverage([3,4,4]) = %v, want 4", avg)
	}
	if avg := ceilAverage(nil); avg != nil {
		t.Errorf("ceilAverage(nil) = %v, want nil", avg)
	}
}

func TestIsResultComplete(t *testing.T) {
	full := subComponentPoints{
		Education: 10, Experience: 6, Potential: 10,
		CompetencyTechnical: 15, CompetencyManagerial: 6,
		ProblemSolving: 60, AssetManagement: 30,
		SubordinateControl: 10, WorkScope: 30, WorkActivity: 6, WorkRisk: 9,
		FinancialWithAuthority: 1500,
	}
	if !isResultComplete(full) {
		t.Error("expected complete for full sub components")
	}

	// Salah satu required kosong → tidak lengkap
	missing := full
	missing.WorkRisk = 0
	if isResultComplete(missing) {
		t.Error("expected incomplete when work_risk=0")
	}

	// Finansial dua-duanya 0 → tidak lengkap
	noFin := full
	noFin.FinancialWithAuthority = 0
	noFin.FinancialWithoutAuthority = 0
	if isResultComplete(noFin) {
		t.Error("expected incomplete when both financial sub points are 0")
	}

	// Pasangan finansial OR: without cukup
	withoutOnly := full
	withoutOnly.FinancialWithAuthority = 0
	withoutOnly.FinancialWithoutAuthority = 60
	if !isResultComplete(withoutOnly) {
		t.Error("expected complete when financial_without_authority > 0")
	}
}
