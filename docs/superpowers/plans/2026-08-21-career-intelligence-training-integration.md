# Career Intelligence Training Integration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Give Career Intelligence read-only access to Training & Development data so it can show an employee's training profile, flag which completed trainings are relevant to their career path, compute the training-driven part of a competency gap, and recommend courses that close that gap — without Career Intelligence writing to or duplicating any Training table.

**Architecture:** Training stays the sole owner and writer of its own tables. Two new read-only Training endpoints expose an employee's aggregate training profile and a competency→course reverse lookup (both additive, no migrations — every column needed already exists). Career Intelligence consumes Training exclusively through a narrow `TrainingProvider` Go interface + adapter (the exact pattern already used for `PerformanceProvider`/`CompetencyProvider` in this module), wired in `cmd/server/main.go`. Two new Career Intelligence endpoints combine that Training data with Career Intelligence's own competency-requirement data (`GetOrgCompetencyRequirements`/`GetEmployeeCompetencyLevels`, already built) to produce the training profile view and the gap-driven course recommendation list.

**Tech Stack:** Go (Gin, GORM, MySQL/PostgreSQL dual-dialect), Vue 3 + PrimeVue frontend.

**Spec:** `docs/career-intelligence-training-enhancement-plan.md` (business requirements — this plan implements P0 items 1–4 of that spec's §27: Employee Training Profile, Training History for Career, Training & Competency Gap, Career Training Recommendation. Items 5/6 of P0 — Career Development Plan integration and full multi-factor Career Readiness scoring — are explicitly OUT of scope for this plan; see "Deferred Scope" below.)

## Global Constraints

- No new database tables or migrations — every field this plan needs already exists in `training_participants`, `training_courses`, `training_certificates`, `training_course_competencies`, or Career Intelligence's existing `competency_score_details`-based repository methods.
- Cross-module reads only through a narrow Go interface + adapter wired in `main.go` — never import `training` package types directly inside `careerintelligence`, and never have `careerintelligence` issue raw SQL against `training_*` tables (Training owns that schema).
- All new handlers return errors via the existing per-module helpers (`httputil.SuccessJSON`/`httputil.ErrorJSON`/`httputil.InternalError` in `training`; `respondSuccess`/`respondError` in `careerintelligence`) — do not invent new response envelopes.
- Every new route with a dynamic `:id`/`:employeeId` segment must be registered after any static-path routes on the same prefix (Gin routing constraint already followed throughout both modules' `routes.go`).
- Commit after each task once its tests pass — do not batch multiple tasks into one commit.

## Deferred Scope (explicitly not built by this plan)

- **Career Development Plan integration** (spec §10 / P0 item 6): there is no `development_plan_item`-equivalent entity anywhere in this codebase (Succession Plans' `development_plan` is a free-text field, not a structured, linkable plan). Building this needs a new table + a product decision on where a "career development plan" entity should live (Career Intelligence, most likely) — out of scope here; flag to the user as a follow-up plan if wanted.
- **Multi-factor Career Readiness score** (spec §11/§25, P1 item 7): combining performance + competency + potential + experience + training + certification into one weighted "readiness %" needs a configurable weighting scheme, which needs its own settings entity (same shape as `TalentMapSettings` from the earlier Talent Map work) and a product decision on default weights. Out of scope here.
- Talent Mapping / Succession Planning training-evidence columns (spec §15/§16, P2): deferred — P0 first.

---

### Task 1: Training — Employee Training Summary (aggregate profile numbers)

**Files:**
- Modify: `backend/internal/modules/training/dto.go` (add `EmployeeTrainingSummaryResponse` near `TrainingHistoryResponse` at line ~939)
- Modify: `backend/internal/modules/training/repository.go` (add `GetEmployeeTrainingSummary` near `HistoryByEmployee` at line ~1938)
- Modify: `backend/internal/modules/training/service.go` (add `GetEmployeeTrainingSummary` near `GetTrainingHistory` at line ~3871)
- Modify: `backend/internal/modules/training/handler.go` (add `GetEmployeeTrainingSummary` near `GetTrainingHistory` at line ~1673)
- Modify: `backend/internal/modules/training/routes.go` (register route near line 219, `trn.GET("/history", ...)`)
- Test: `backend/internal/modules/training/repository_test.go` (new test function, follow existing test file's `setupTestDB`/fixture conventions in that file)

**Interfaces:**
- Produces: `func (r *Repository) GetEmployeeTrainingSummary(ctx context.Context, employeeID uuid.UUID) (*EmployeeTrainingSummaryResponse, error)` — used directly by Task 4's `TrainingProvider` adapter.
- Produces: `func (s *Service) GetEmployeeTrainingSummary(ctx context.Context, employeeID string) (*EmployeeTrainingSummaryResponse, error)` — used directly by Task 4's `TrainingProvider` adapter.
- Produces: `GET /api/v1/tenant/training/employees/:employeeId/summary` → `200 {"success":true,"data":EmployeeTrainingSummaryResponse}`.

- [ ] **Step 1: Add the response DTO**

In `backend/internal/modules/training/dto.go`, immediately before `type TrainingHistoryResponse struct {` (line 939), add:

```go
// EmployeeTrainingSummaryResponse — aggregate training profile numbers for
// one employee (Career Intelligence Training Enhancement plan §5 — Employee
// Training Profile). Read-only projection over training_participants /
// training_courses / training_certificates / training_course_competencies;
// no new table.
type EmployeeTrainingSummaryResponse struct {
	EmployeeID               string  `json:"employee_id"`
	TotalTraining            int     `json:"total_training"`
	Completed                int     `json:"completed"`
	Failed                   int     `json:"failed"`
	TrainingHours            float64 `json:"training_hours"`
	AverageScore             float64 `json:"average_score"`
	CertificationCount       int     `json:"certification_count"`
	CompetencyTrainingCount  int     `json:"competency_training_count"`
}
```

- [ ] **Step 2: Write the failing repository test**

In `backend/internal/modules/training/repository_test.go`, add (matching this file's existing fixture-seeding style — check the top of the file for `setupTestDB`/course+session+participant seed helpers already used by neighboring tests and reuse them rather than re-deriving seed SQL):

```go
func TestRepository_GetEmployeeTrainingSummary(t *testing.T) {
	repo, cleanup := newTestRepository(t)
	defer cleanup()

	employeeID := uuid.New()
	courseID := seedTrainingCourse(t, repo, "Leadership Development", 16) // name, duration_hour
	competencyID := uuid.New()
	seedCourseCompetency(t, repo, courseID, competencyID, 4)
	sessionID := seedTrainingSession(t, repo, courseID)
	participantID := seedTrainingParticipant(t, repo, sessionID, employeeID, "COMPLETED", 92)
	seedTrainingCertificate(t, repo, participantID)

	summary, err := repo.GetEmployeeTrainingSummary(context.Background(), employeeID)
	if err != nil {
		t.Fatalf("GetEmployeeTrainingSummary failed: %v", err)
	}
	if summary.TotalTraining != 1 {
		t.Errorf("expected total_training 1, got %d", summary.TotalTraining)
	}
	if summary.Completed != 1 {
		t.Errorf("expected completed 1, got %d", summary.Completed)
	}
	if summary.TrainingHours != 16 {
		t.Errorf("expected training_hours 16, got %.2f", summary.TrainingHours)
	}
	if summary.AverageScore != 92 {
		t.Errorf("expected average_score 92, got %.2f", summary.AverageScore)
	}
	if summary.CertificationCount != 1 {
		t.Errorf("expected certification_count 1, got %d", summary.CertificationCount)
	}
	if summary.CompetencyTrainingCount != 1 {
		t.Errorf("expected competency_training_count 1, got %d", summary.CompetencyTrainingCount)
	}
}
```

If this test file has no `seedTrainingCourse`/`seedCourseCompetency`/`seedTrainingSession`/`seedTrainingParticipant`/`seedTrainingCertificate` helpers yet, write them as small `t.Helper()` functions right above this test that `db.Create(&TrainingCourse{...})` etc. directly against the models already defined in `model.go` (`TrainingCourse`, `TrainingCourseCompetency`, `TrainingSession`, `TrainingParticipant`, `TrainingCertificate`) — do not invent new fields, use exactly the fields listed in this plan's research (`DurationHour`, `CourseID`, `CompetencyID`, `TargetLevel`, `EmployeeID`, `SessionID`, `CompletionStatus`, `Score`, `ParticipantID`).

- [ ] **Step 3: Run the test to verify it fails**

Run: `cd backend && go test ./internal/modules/training/... -run TestRepository_GetEmployeeTrainingSummary -v`
Expected: FAIL — `GetEmployeeTrainingSummary` undefined.

- [ ] **Step 4: Implement the repository method**

In `backend/internal/modules/training/repository.go`, immediately before `func (r *Repository) HistoryByEmployee(...)` (line 1938), add:

```go
// GetEmployeeTrainingSummary aggregates one employee's training profile
// numbers (Career Intelligence Training Enhancement plan §5). Read-only
// projection, no new table.
func (r *Repository) GetEmployeeTrainingSummary(ctx context.Context, employeeID uuid.UUID) (*EmployeeTrainingSummaryResponse, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	summary := &EmployeeTrainingSummaryResponse{EmployeeID: employeeID.String()}
	query := `
		SELECT
			COUNT(DISTINCT tp.id) AS total_training,
			COALESCE(SUM(CASE WHEN tp.completion_status = 'COMPLETED' THEN 1 ELSE 0 END), 0) AS completed,
			COALESCE(SUM(CASE WHEN tp.completion_status = 'FAILED' THEN 1 ELSE 0 END), 0) AS failed,
			COALESCE(SUM(c.duration_hour), 0) AS training_hours,
			COALESCE(AVG(NULLIF(tp.score, 0)), 0) AS average_score,
			COUNT(DISTINCT tc.id) AS certification_count,
			COUNT(DISTINCT tcc.competency_id) AS competency_training_count
		FROM training_participants tp
		JOIN training_sessions ts ON ts.id = tp.session_id
		JOIN training_courses c ON c.id = ts.course_id
		LEFT JOIN training_certificates tc ON tc.participant_id = tp.id
		LEFT JOIN training_course_competencies tcc ON tcc.course_id = c.id AND tcc.deleted_at IS NULL
		WHERE tp.employee_id = ? AND tp.deleted_at IS NULL`
	if err := db.WithContext(ctx).Raw(query, employeeID).Scan(summary).Error; err != nil {
		return nil, fmt.Errorf("failed to aggregate training summary: %w", err)
	}
	return summary, nil
}
```

- [ ] **Step 5: Run the test to verify it passes**

Run: `cd backend && go test ./internal/modules/training/... -run TestRepository_GetEmployeeTrainingSummary -v`
Expected: PASS

- [ ] **Step 6: Add the service wrapper**

In `backend/internal/modules/training/service.go`, immediately before `func (s *Service) GetTrainingHistory(...)` (line 3871), add:

```go
func (s *Service) GetEmployeeTrainingSummary(ctx context.Context, employeeID string) (*EmployeeTrainingSummaryResponse, error) {
	uid, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	return s.repo.GetEmployeeTrainingSummary(ctx, uid)
}
```

- [ ] **Step 7: Add the handler**

In `backend/internal/modules/training/handler.go`, immediately before `func (h *Handler) GetTrainingHistory(...)` (line 1673), add:

```go
func (h *Handler) GetEmployeeTrainingSummary(c *gin.Context) {
	employeeID := c.Param("employeeId")
	resp, err := h.svc.GetEmployeeTrainingSummary(c.Request.Context(), employeeID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
```

- [ ] **Step 8: Register the route**

In `backend/internal/modules/training/routes.go`, immediately before `trn.GET("/history", handler.GetTrainingHistory)` (line 219), add:

```go
trn.GET("/employees/:employeeId/summary", handler.GetEmployeeTrainingSummary)
```

- [ ] **Step 9: Build and run the full training test suite**

Run: `cd backend && go build ./... && go vet ./... && go test ./internal/modules/training/...`
Expected: build clean, vet clean, all tests PASS (including the new one).

- [ ] **Step 10: Commit**

```bash
cd backend
git add internal/modules/training/dto.go internal/modules/training/repository.go internal/modules/training/repository_test.go internal/modules/training/service.go internal/modules/training/handler.go internal/modules/training/routes.go
git commit -m "$(cat <<'EOF'
feat(training): tambah endpoint employee training summary

GET /training/employees/:employeeId/summary -- agregasi total/completed/
failed training, jam training, skor rata-rata, jumlah sertifikasi, dan
jumlah training yang terhubung ke competency. Read-only projection atas
tabel training yang sudah ada, tidak ada migrasi baru. Dipakai Career
Intelligence Training Enhancement plan (docs/career-intelligence-training-enhancement-plan.md
§5) via narrow TrainingProvider interface (Task 3).

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>
EOF
)"
```

---

### Task 2: Training — reverse competency→course lookup

**Files:**
- Modify: `backend/internal/modules/training/dto.go` (add `CourseCompetencyMatchResponse` near `TrainingCourseCompetency`-related DTOs)
- Modify: `backend/internal/modules/training/repository.go` (add `ListCoursesByCompetencyIDs` near `ListCourseCompetencies` at line ~1359)
- Modify: `backend/internal/modules/training/service.go` (add `ListCoursesByCompetencyIDs`)
- Modify: `backend/internal/modules/training/handler.go` (add `ListCoursesByCompetencyIDs`)
- Modify: `backend/internal/modules/training/routes.go` (register near the existing `/courses/:id/competencies` routes)
- Test: `backend/internal/modules/training/repository_test.go`

**Interfaces:**
- Consumes: `TrainingCourseCompetency{ID, CourseID, CompetencyID, TargetLevel *int}` (existing model, `model.go:673-683`).
- Produces: `func (r *Repository) ListCoursesByCompetencyIDs(ctx context.Context, competencyIDs []uuid.UUID) ([]CourseCompetencyMatchResponse, error)` — used directly by Task 5's recommendation logic via the `TrainingProvider` adapter (Task 3).
- Produces: `GET /api/v1/tenant/training/courses/by-competency?competency_ids=<uuid>,<uuid>` (query param, comma-separated) → `200 {"success":true,"data":[]CourseCompetencyMatchResponse}`.

- [ ] **Step 1: Add the response DTO**

In `backend/internal/modules/training/dto.go`, add (place near other course-competency DTOs; if none exist yet, place next to `TrainingHistoryResponse`):

```go
// CourseCompetencyMatchResponse is one course that develops a requested
// competency, used by Career Intelligence to recommend training that closes
// a competency gap (Career Intelligence Training Enhancement plan §7/§8).
type CourseCompetencyMatchResponse struct {
	CourseID     string `json:"course_id"`
	CourseName   string `json:"course_name"`
	CompetencyID string `json:"competency_id"`
	TargetLevel  *int   `json:"target_level,omitempty"`
	IsMandatory  bool   `json:"is_mandatory"`
	IsCertified  bool   `json:"is_certified"`
}
```

- [ ] **Step 2: Write the failing repository test**

In `backend/internal/modules/training/repository_test.go`, add:

```go
func TestRepository_ListCoursesByCompetencyIDs(t *testing.T) {
	repo, cleanup := newTestRepository(t)
	defer cleanup()

	competencyID := uuid.New()
	otherCompetencyID := uuid.New()
	courseID := seedTrainingCourse(t, repo, "Leadership Development", 16)
	seedCourseCompetency(t, repo, courseID, competencyID, 4)

	matches, err := repo.ListCoursesByCompetencyIDs(context.Background(), []uuid.UUID{competencyID})
	if err != nil {
		t.Fatalf("ListCoursesByCompetencyIDs failed: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].CourseID != courseID.String() {
		t.Errorf("expected course id %s, got %s", courseID, matches[0].CourseID)
	}

	noMatches, err := repo.ListCoursesByCompetencyIDs(context.Background(), []uuid.UUID{otherCompetencyID})
	if err != nil {
		t.Fatalf("ListCoursesByCompetencyIDs (no match) failed: %v", err)
	}
	if len(noMatches) != 0 {
		t.Errorf("expected 0 matches for unrelated competency, got %d", len(noMatches))
	}
}
```

- [ ] **Step 3: Run the test to verify it fails**

Run: `cd backend && go test ./internal/modules/training/... -run TestRepository_ListCoursesByCompetencyIDs -v`
Expected: FAIL — `ListCoursesByCompetencyIDs` undefined.

- [ ] **Step 4: Implement the repository method**

In `backend/internal/modules/training/repository.go`, immediately after `ListCourseCompetencies` (ends around line 1371), add:

```go
// ListCoursesByCompetencyIDs finds every active course that develops any of
// the given competencies (reverse lookup of training_course_competencies),
// used to recommend training that closes a Career Intelligence competency
// gap. Empty input returns an empty slice, not an error.
func (r *Repository) ListCoursesByCompetencyIDs(ctx context.Context, competencyIDs []uuid.UUID) ([]CourseCompetencyMatchResponse, error) {
	if len(competencyIDs) == 0 {
		return []CourseCompetencyMatchResponse{}, nil
	}
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []CourseCompetencyMatchResponse
	err = db.WithContext(ctx).Table("training_course_competencies tcc").
		Select("c.id AS course_id, c.name AS course_name, tcc.competency_id AS competency_id, tcc.target_level AS target_level, c.is_mandatory AS is_mandatory, c.is_certified AS is_certified").
		Joins("JOIN training_courses c ON c.id = tcc.course_id AND c.is_active = ?", true).
		Where("tcc.competency_id IN ? AND tcc.deleted_at IS NULL", competencyIDs).
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list courses by competency: %w", err)
	}
	return rows, nil
}
```

- [ ] **Step 5: Run the test to verify it passes**

Run: `cd backend && go test ./internal/modules/training/... -run TestRepository_ListCoursesByCompetencyIDs -v`
Expected: PASS

- [ ] **Step 6: Add the service wrapper**

In `backend/internal/modules/training/service.go`, add near `GetEmployeeTrainingSummary` (from Task 1):

```go
func (s *Service) ListCoursesByCompetencyIDs(ctx context.Context, competencyIDs []uuid.UUID) ([]CourseCompetencyMatchResponse, error) {
	return s.repo.ListCoursesByCompetencyIDs(ctx, competencyIDs)
}
```

- [ ] **Step 7: Add the handler**

In `backend/internal/modules/training/handler.go`, add near `GetEmployeeTrainingSummary` (from Task 1):

```go
func (h *Handler) ListCoursesByCompetencyIDs(c *gin.Context) {
	raw := c.Query("competency_ids")
	if raw == "" {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "training.competency_ids_required")
		return
	}
	parts := strings.Split(raw, ",")
	ids := make([]uuid.UUID, 0, len(parts))
	for _, p := range parts {
		id, err := uuid.Parse(strings.TrimSpace(p))
		if err != nil {
			httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "training.invalid_competency_id")
			return
		}
		ids = append(ids, id)
	}
	resp, err := h.svc.ListCoursesByCompetencyIDs(c.Request.Context(), ids)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
```

Check the top of `handler.go` already imports `"strings"` — if not, add it to the import block.

- [ ] **Step 8: Register the route**

In `backend/internal/modules/training/routes.go`, add near the existing course-competency routes (search for `/courses/:id/competencies` in that file and add immediately after that group, before any `/:id`-only catch-all on `/courses`):

```go
trn.GET("/courses/by-competency", handler.ListCoursesByCompetencyIDs)
```

Confirm this line is registered before any `trn.GET("/courses/:id", ...)` route in the same file — Gin matches static segments first only if declared first in some versions; follow this file's existing convention of static-before-`:id` (already true for `/history` before `/:id` at line 219).

- [ ] **Step 9: Build and run the full training test suite**

Run: `cd backend && go build ./... && go vet ./... && go test ./internal/modules/training/...`
Expected: build clean, vet clean, all tests PASS.

- [ ] **Step 10: Commit**

```bash
cd backend
git add internal/modules/training/dto.go internal/modules/training/repository.go internal/modules/training/repository_test.go internal/modules/training/service.go internal/modules/training/handler.go internal/modules/training/routes.go
git commit -m "$(cat <<'EOF'
feat(training): tambah reverse lookup competency -> course

GET /training/courses/by-competency?competency_ids=a,b -- course aktif
mana saja yang mengembangkan salah satu competency yang diminta (reverse
lookup training_course_competencies). Dipakai Career Intelligence untuk
merekomendasikan training yang menutup competency gap (docs/career-intelligence-training-enhancement-plan.md
§7/§8) via narrow TrainingProvider interface (Task 3).

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>
EOF
)"
```

---

### Task 3: Career Intelligence — `TrainingProvider` interface + wiring

**Files:**
- Modify: `backend/internal/modules/careerintelligence/service.go` (add interface + field + setter, near `PerformanceProvider`/`CompetencyProvider` at the top of the file)
- Modify: `backend/cmd/server/main.go` (add adapter type near `performanceEligibilityAdapter`/`competencyEligibilityAdapter`, wire after line 1251)
- Test: `backend/internal/modules/careerintelligence/service_test.go` (fake provider + one assertion that a service method calls it — full behavioral tests happen in Tasks 4/5, this task just proves the wiring compiles and a fake can be injected)

**Interfaces:**
- Consumes: `training.Service.GetEmployeeTrainingSummary(ctx, employeeID string) (*training.EmployeeTrainingSummaryResponse, error)` and `training.Service.GetTrainingHistory(ctx, employeeID string) ([]training.TrainingHistoryResponse, error)` and `training.Service.ListCoursesByCompetencyIDs(ctx, competencyIDs []uuid.UUID) ([]training.CourseCompetencyMatchResponse, error)` (Tasks 1 & 2, plus `GetTrainingHistory` which already exists at `training/service.go:3871`).
- Produces: `careerintelligence.TrainingProvider` interface + `Service.SetTrainingProvider(p TrainingProvider)` + `Service.trainingProvider` field — consumed directly by Tasks 4 & 5.

- [ ] **Step 1: Define the narrow interface in careerintelligence**

In `backend/internal/modules/careerintelligence/service.go`, immediately after the existing `CompetencyProvider` interface definition (right before `type Service struct {`), add:

```go
// TrainingProvider membaca data training read-only dari modul Training —
// careerintelligence tidak pernah menulis ke tabel training_* (Career
// Intelligence Training Enhancement plan, docs/career-intelligence-training-enhancement-plan.md).
// Shape-nya sengaja pakai tipe lokal (bukan import training package) supaya
// careerintelligence tidak coupled ke training's internal types.
type TrainingProvider interface {
	GetTrainingSummary(ctx context.Context, employeeID string) (*TrainingSummary, error)
	GetTrainingHistory(ctx context.Context, employeeID string) ([]TrainingHistoryItem, error)
	ListCoursesByCompetencyIDs(ctx context.Context, competencyIDs []uuid.UUID) ([]RecommendedCourse, error)
}

// TrainingSummary mirrors training.EmployeeTrainingSummaryResponse's fields
// (careerintelligence's own copy, decoupled from the training package).
type TrainingSummary struct {
	TotalTraining           int
	Completed                int
	Failed                   int
	TrainingHours            float64
	AverageScore             float64
	CertificationCount       int
	CompetencyTrainingCount  int
}

// TrainingHistoryItem mirrors training.TrainingHistoryResponse's fields
// needed by Career Intelligence.
type TrainingHistoryItem struct {
	CourseID         string
	CourseName       string
	StartDate        string
	CompletionStatus string
	Score            float64
	CertificateNo    string
}

// RecommendedCourse mirrors training.CourseCompetencyMatchResponse's fields.
type RecommendedCourse struct {
	CourseID     string
	CourseName   string
	CompetencyID string
	TargetLevel  *int
	IsMandatory  bool
	IsCertified  bool
}
```

- [ ] **Step 2: Add the field and setter on Service**

In `backend/internal/modules/careerintelligence/service.go`, in the `Service` struct (where `performanceProvider`/`competencyProvider`/`notifier` fields already live), add:

```go
	trainingProvider TrainingProvider
```

Immediately after the existing `SetNotifier` method, add:

```go
// SetTrainingProvider wires the Training data source used by
// GetEmployeeTrainingProfile (Task 4) and GetTrainingRecommendations
// (Task 5). Must be called before those methods are used.
func (s *Service) SetTrainingProvider(p TrainingProvider) {
	s.trainingProvider = p
}
```

- [ ] **Step 3: Write the adapter in main.go**

In `backend/cmd/server/main.go`, immediately after the `competencyEligibilityAdapter` type definition (near `performanceEligibilityAdapter`/`competencyEligibilityAdapter`, before the `okrEligibilityAdapter` block), add:

```go
// trainingProviderAdapter implements careerintelligence.TrainingProvider so
// Career Intelligence can read training profile/history/course-competency
// data without importing the training package's own types (docs/career-intelligence-training-enhancement-plan.md).
type trainingProviderAdapter struct {
	svc *training.Service
}

func (a trainingProviderAdapter) GetTrainingSummary(ctx context.Context, employeeID string) (*careerintelligence.TrainingSummary, error) {
	s, err := a.svc.GetEmployeeTrainingSummary(ctx, employeeID)
	if err != nil {
		return nil, err
	}
	return &careerintelligence.TrainingSummary{
		TotalTraining:           s.TotalTraining,
		Completed:               s.Completed,
		Failed:                  s.Failed,
		TrainingHours:           s.TrainingHours,
		AverageScore:            s.AverageScore,
		CertificationCount:      s.CertificationCount,
		CompetencyTrainingCount: s.CompetencyTrainingCount,
	}, nil
}

func (a trainingProviderAdapter) GetTrainingHistory(ctx context.Context, employeeID string) ([]careerintelligence.TrainingHistoryItem, error) {
	rows, err := a.svc.GetTrainingHistory(ctx, employeeID)
	if err != nil {
		return nil, err
	}
	items := make([]careerintelligence.TrainingHistoryItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, careerintelligence.TrainingHistoryItem{
			CourseID:         r.CourseID,
			CourseName:       r.CourseName,
			StartDate:        r.StartDate,
			CompletionStatus: r.CompletionStatus,
			Score:            r.Score,
			CertificateNo:    r.CertificateNo,
		})
	}
	return items, nil
}

func (a trainingProviderAdapter) ListCoursesByCompetencyIDs(ctx context.Context, competencyIDs []uuid.UUID) ([]careerintelligence.RecommendedCourse, error) {
	rows, err := a.svc.ListCoursesByCompetencyIDs(ctx, competencyIDs)
	if err != nil {
		return nil, err
	}
	courses := make([]careerintelligence.RecommendedCourse, 0, len(rows))
	for _, r := range rows {
		courses = append(courses, careerintelligence.RecommendedCourse{
			CourseID:     r.CourseID,
			CourseName:   r.CourseName,
			CompetencyID: r.CompetencyID,
			TargetLevel:  r.TargetLevel,
			IsMandatory:  r.IsMandatory,
			IsCertified:  r.IsCertified,
		})
	}
	return courses, nil
}
```

- [ ] **Step 4: Wire the adapter**

In `backend/cmd/server/main.go`, immediately after `ciSvc.SetNotifier(notificationSvc)` (line 1251), add:

```go
	// Career Intelligence Training Enhancement (docs/career-intelligence-training-enhancement-plan.md):
	// training data flows into Career Intelligence read-only via this adapter.
	ciSvc.SetTrainingProvider(trainingProviderAdapter{svc: trainingSvc})
```

- [ ] **Step 5: Verify it builds**

Run: `cd backend && go build ./...`
Expected: builds clean. (No behavior to test yet — Tasks 4/5 exercise this wiring. This step only proves the interface/adapter compile against each other.)

- [ ] **Step 6: Write a fake provider for future tests**

In `backend/internal/modules/careerintelligence/service_test.go`, add (this fake is consumed by Tasks 4 & 5's tests):

```go
// fakeTrainingProvider implements TrainingProvider for tests.
type fakeTrainingProvider struct {
	summary   *TrainingSummary
	history   []TrainingHistoryItem
	courses   []RecommendedCourse
	summaryErr, historyErr, coursesErr error
}

func (f *fakeTrainingProvider) GetTrainingSummary(_ context.Context, _ string) (*TrainingSummary, error) {
	return f.summary, f.summaryErr
}

func (f *fakeTrainingProvider) GetTrainingHistory(_ context.Context, _ string) ([]TrainingHistoryItem, error) {
	return f.history, f.historyErr
}

func (f *fakeTrainingProvider) ListCoursesByCompetencyIDs(_ context.Context, _ []uuid.UUID) ([]RecommendedCourse, error) {
	return f.courses, f.coursesErr
}
```

- [ ] **Step 7: Run the full careerintelligence test suite**

Run: `cd backend && go build ./... && go vet ./... && go test ./internal/modules/careerintelligence/...`
Expected: build clean, vet clean, all existing tests still PASS (the fake provider isn't used by any test yet — this just confirms it compiles).

- [ ] **Step 8: Commit**

```bash
cd backend
git add internal/modules/careerintelligence/service.go internal/modules/careerintelligence/service_test.go cmd/server/main.go
git commit -m "$(cat <<'EOF'
feat(career-intelligence): TrainingProvider interface + wiring

Narrow interface + adapter (pola sama PerformanceProvider/CompetencyProvider)
supaya careerintelligence bisa baca training profile/history/course-competency
dari modul Training tanpa import training package langsung. Belum dipakai
di endpoint manapun -- itu Task 4 (training profile) & Task 5 (recommendation).

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>
EOF
)"
```

---

### Task 4: Career Intelligence — Employee Training Profile endpoint

**Files:**
- Modify: `backend/internal/modules/careerintelligence/dto.go` (add `TrainingProfileResponse`)
- Modify: `backend/internal/modules/careerintelligence/service.go` (add `GetEmployeeTrainingProfile`)
- Modify: `backend/internal/modules/careerintelligence/handler.go` (add `GetEmployeeTrainingProfile`)
- Modify: `backend/internal/modules/careerintelligence/routes.go` (register route)
- Test: `backend/internal/modules/careerintelligence/service_test.go`

**Interfaces:**
- Consumes: `Service.trainingProvider TrainingProvider` (Task 3), `TrainingSummary`/`TrainingHistoryItem` types (Task 3).
- Produces: `func (s *Service) GetEmployeeTrainingProfile(ctx context.Context, employeeID string) (*TrainingProfileResponse, error)`.
- Produces: `GET /api/v1/tenant/career-intelligence/employees/:employeeId/training-profile` → `200 {"success":true,"data":TrainingProfileResponse}`.

- [ ] **Step 1: Add the response DTO**

In `backend/internal/modules/careerintelligence/dto.go`, add near `EmployeeTalentProfileResponse`:

```go
// TrainingProfileResponse — Employee Training Profile (Career Intelligence
// Training Enhancement plan §5). Read-only view over Training data via
// TrainingProvider; Career Intelligence stores none of this itself.
type TrainingProfileResponse struct {
	EmployeeID              string                     `json:"employee_id"`
	TotalTraining           int                        `json:"total_training"`
	Completed               int                        `json:"completed"`
	Failed                  int                        `json:"failed"`
	TrainingHours           float64                    `json:"training_hours"`
	AverageScore            float64                    `json:"average_score"`
	CertificationCount      int                        `json:"certification_count"`
	CompetencyTrainingCount int                        `json:"competency_training_count"`
	History                 []TrainingHistoryItemResponse `json:"history"`
}

type TrainingHistoryItemResponse struct {
	CourseID         string  `json:"course_id"`
	CourseName       string  `json:"course_name"`
	StartDate        string  `json:"start_date"`
	CompletionStatus string  `json:"completion_status"`
	Score            float64 `json:"score"`
	CertificateNo    string  `json:"certificate_no,omitempty"`
}
```

- [ ] **Step 2: Write the failing service test**

In `backend/internal/modules/careerintelligence/service_test.go`, add:

```go
func TestService_GetEmployeeTrainingProfile_Success(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	fake := &fakeTrainingProvider{
		summary: &TrainingSummary{
			TotalTraining: 3, Completed: 2, Failed: 1,
			TrainingHours: 40, AverageScore: 85,
			CertificationCount: 1, CompetencyTrainingCount: 2,
		},
		history: []TrainingHistoryItem{
			{CourseID: uuidStr(), CourseName: "Leadership Development", StartDate: "2026-01-10", CompletionStatus: "COMPLETED", Score: 92},
		},
	}
	svc.SetTrainingProvider(fake)

	profile, err := svc.GetEmployeeTrainingProfile(ctx(), uuidStr())
	if err != nil {
		t.Fatalf("GetEmployeeTrainingProfile failed: %v", err)
	}
	if profile.TotalTraining != 3 {
		t.Errorf("expected total_training 3, got %d", profile.TotalTraining)
	}
	if len(profile.History) != 1 {
		t.Fatalf("expected 1 history item, got %d", len(profile.History))
	}
	if profile.History[0].CourseName != "Leadership Development" {
		t.Errorf("expected course name 'Leadership Development', got '%s'", profile.History[0].CourseName)
	}
}

func TestService_GetEmployeeTrainingProfile_NoProvider(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()

	// TrainingProvider never wired -- must fail clearly, not panic.
	if _, err := svc.GetEmployeeTrainingProfile(ctx(), uuidStr()); err == nil {
		t.Error("expected error when training provider is not configured")
	}
}
```

- [ ] **Step 3: Run the tests to verify they fail**

Run: `cd backend && go test ./internal/modules/careerintelligence/... -run TestService_GetEmployeeTrainingProfile -v`
Expected: FAIL — `GetEmployeeTrainingProfile` undefined.

- [ ] **Step 4: Implement the service method**

In `backend/internal/modules/careerintelligence/service.go`, add (near `GetEmployeeTalentProfile`):

```go
// GetEmployeeTrainingProfile — Employee Training Profile (Career
// Intelligence Training Enhancement plan §5): summary numbers + relevant
// history, sourced entirely from Training via trainingProvider.
func (s *Service) GetEmployeeTrainingProfile(ctx context.Context, employeeID string) (*TrainingProfileResponse, error) {
	if s.trainingProvider == nil {
		return nil, fmt.Errorf("training profile belum dikonfigurasi (training provider)")
	}
	summary, err := s.trainingProvider.GetTrainingSummary(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to read training summary: %w", err)
	}
	history, err := s.trainingProvider.GetTrainingHistory(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to read training history: %w", err)
	}
	historyResp := make([]TrainingHistoryItemResponse, 0, len(history))
	for _, h := range history {
		historyResp = append(historyResp, TrainingHistoryItemResponse{
			CourseID:         h.CourseID,
			CourseName:       h.CourseName,
			StartDate:        h.StartDate,
			CompletionStatus: h.CompletionStatus,
			Score:            h.Score,
			CertificateNo:    h.CertificateNo,
		})
	}
	return &TrainingProfileResponse{
		EmployeeID:              employeeID,
		TotalTraining:           summary.TotalTraining,
		Completed:               summary.Completed,
		Failed:                  summary.Failed,
		TrainingHours:           summary.TrainingHours,
		AverageScore:            summary.AverageScore,
		CertificationCount:      summary.CertificationCount,
		CompetencyTrainingCount: summary.CompetencyTrainingCount,
		History:                 historyResp,
	}, nil
}
```

- [ ] **Step 5: Run the tests to verify they pass**

Run: `cd backend && go test ./internal/modules/careerintelligence/... -run TestService_GetEmployeeTrainingProfile -v`
Expected: PASS (both tests).

- [ ] **Step 6: Add the handler**

In `backend/internal/modules/careerintelligence/handler.go`, add (near `GetEmployeeTalentProfile`):

```go
func (h *Handler) GetEmployeeTrainingProfile(c *gin.Context) {
	result, err := h.svc.GetEmployeeTrainingProfile(c.Request.Context(), c.Param("employeeId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	respondSuccess(c, result)
}
```

- [ ] **Step 7: Register the route**

In `backend/internal/modules/careerintelligence/routes.go`, add (place among the other `employees/:employeeId/...` style routes if any exist, otherwise near `/talent-maps/employee/:employeeId`):

```go
ci.GET("/employees/:employeeId/training-profile", handler.GetEmployeeTrainingProfile)
```

- [ ] **Step 8: Build and run the full careerintelligence test suite**

Run: `cd backend && go build ./... && go vet ./... && go test ./internal/modules/careerintelligence/...`
Expected: build clean, vet clean, all tests PASS.

- [ ] **Step 9: Commit**

```bash
cd backend
git add internal/modules/careerintelligence/dto.go internal/modules/careerintelligence/service.go internal/modules/careerintelligence/service_test.go internal/modules/careerintelligence/handler.go internal/modules/careerintelligence/routes.go
git commit -m "$(cat <<'EOF'
feat(career-intelligence): endpoint Employee Training Profile

GET /career-intelligence/employees/:employeeId/training-profile --
ringkasan (total/completed/failed/jam/skor rata-rata/sertifikasi/jumlah
training terhubung competency) + riwayat training, dibaca read-only dari
Training via TrainingProvider (Task 3). Career Intelligence Training
Enhancement plan §5 (docs/career-intelligence-training-enhancement-plan.md).

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>
EOF
)"
```

---

### Task 5: Career Intelligence — Training & Competency Gap + Recommendation endpoint

**Files:**
- Modify: `backend/internal/modules/careerintelligence/dto.go` (add `TrainingRecommendationResponse`, `TrainingRecommendationItem`)
- Modify: `backend/internal/modules/careerintelligence/service.go` (add `GetTrainingRecommendations`)
- Modify: `backend/internal/modules/careerintelligence/handler.go` (add `GetTrainingRecommendations`)
- Modify: `backend/internal/modules/careerintelligence/routes.go` (register route)
- Test: `backend/internal/modules/careerintelligence/service_test.go`

**Interfaces:**
- Consumes: `Repository.GetOrgCompetencyRequirements(ctx, orgID uuid.UUID) ([]CompetencyRequirement, error)` (existing, `repository.go` — built during the Gap Analysis real-data work), `Repository.GetEmployeeCompetencyLevels(ctx, employeeID uuid.UUID) (map[uuid.UUID]int, error)` (existing, same file), `Repository.GetPositionTitle(ctx, titleID uuid.UUID) (string, error)` (existing), `Service.trainingProvider.ListCoursesByCompetencyIDs` (Task 3).
- Produces: `func (s *Service) GetTrainingRecommendations(ctx context.Context, employeeID, targetOrgID string) (*TrainingRecommendationResponse, error)`.
- Produces: `GET /api/v1/tenant/career-intelligence/employees/:employeeId/training-recommendations?target_title_id=<uuid>` → `200 {"success":true,"data":TrainingRecommendationResponse}`.

- [ ] **Step 1: Add the response DTOs**

In `backend/internal/modules/careerintelligence/dto.go`, add near `GapAnalysisResponse`:

```go
// TrainingRecommendationResponse — Training & Competency Gap + Career
// Training Recommendation (Career Intelligence Training Enhancement plan
// §7/§8). Gap computation reuses the same GetOrgCompetencyRequirements /
// GetEmployeeCompetencyLevels data as GetGapAnalysis; recommendations are
// courses (via TrainingProvider) that develop a gapped competency.
type TrainingRecommendationResponse struct {
	EmployeeID    string                       `json:"employee_id"`
	TargetTitle   string                       `json:"target_title"`
	Recommendations []TrainingRecommendationItem `json:"recommendations"`
}

type TrainingRecommendationItem struct {
	CompetencyID   string `json:"competency_id"`
	CompetencyName string `json:"competency_name"`
	CurrentLevel   int    `json:"current_level"`
	RequiredLevel  int    `json:"required_level"`
	Gap            int    `json:"gap"`
	Priority       string `json:"priority"` // HIGH / MEDIUM / LOW -- same bands as GetGapAnalysis
	CourseID       string `json:"course_id,omitempty"`
	CourseName     string `json:"course_name,omitempty"`
	IsMandatory    bool   `json:"is_mandatory,omitempty"`
	IsCertified    bool   `json:"is_certified,omitempty"`
}
```

- [ ] **Step 2: Write the failing service test**

In `backend/internal/modules/careerintelligence/service_test.go`, add:

```go
func TestService_GetTrainingRecommendations_Success(t *testing.T) {
	svc, repo, _, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	orgID := uuid.New()
	seedGapAnalysisCompetencyData(t, repo, employeeID, orgID) // reuse Gap Analysis fixture (Leadership req 4/actual 4, Communication req 3/actual 2)

	fake := &fakeTrainingProvider{
		courses: []RecommendedCourse{
			{CourseID: uuid.New().String(), CourseName: "Communication Skills", CompetencyID: "", IsMandatory: false, IsCertified: true},
		},
	}
	svc.SetTrainingProvider(fake)

	resp, err := svc.GetTrainingRecommendations(ctx(), employeeID.String(), orgID.String())
	if err != nil {
		t.Fatalf("GetTrainingRecommendations failed: %v", err)
	}
	// Leadership (4 required, 4 actual) is met -> no recommendation.
	// Communication (3 required, 2 actual) has a gap -> one recommendation.
	if len(resp.Recommendations) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(resp.Recommendations))
	}
	item := resp.Recommendations[0]
	if item.Gap != 1 {
		t.Errorf("expected gap 1, got %d", item.Gap)
	}
	if item.CourseName != "Communication Skills" {
		t.Errorf("expected course 'Communication Skills', got '%s'", item.CourseName)
	}
}

func TestService_GetTrainingRecommendations_NoRequirements(t *testing.T) {
	svc, _, _, cleanup := newTestService()
	defer cleanup()
	svc.SetTrainingProvider(&fakeTrainingProvider{})

	// competency_score_details table doesn't exist in this test's DB (never
	// seeded) -- must surface a clear error, same contract as GetGapAnalysis.
	if _, err := svc.GetTrainingRecommendations(ctx(), uuidStr(), uuidStr()); err == nil {
		t.Error("expected error when target org has no competency assessment data")
	}
}
```

This reuses the `seedGapAnalysisCompetencyData` fixture already written in `helpers_test.go` for the Gap Analysis real-data work (2 competencies: Leadership required 4/employee 4 — met; Communication required 3/employee 2 — gap 1).

- [ ] **Step 3: Run the tests to verify they fail**

Run: `cd backend && go test ./internal/modules/careerintelligence/... -run TestService_GetTrainingRecommendations -v`
Expected: FAIL — `GetTrainingRecommendations` undefined.

- [ ] **Step 4: Implement the service method**

In `backend/internal/modules/careerintelligence/service.go`, add (near `GetGapAnalysis`):

```go
// GetTrainingRecommendations — Training & Competency Gap (plan §7) +
// Career Training Recommendation (plan §8). Reuses the exact same
// requirement/actual-level data as GetGapAnalysis (target org's latest
// finalized competency_scores as requirement, employee's latest as actual)
// so the two features never disagree about what the gap is; the only new
// step is looking up which course (if any) develops each gapped
// competency via trainingProvider.
func (s *Service) GetTrainingRecommendations(ctx context.Context, employeeID, targetOrgID string) (*TrainingRecommendationResponse, error) {
	empID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	targetID, err := uuid.Parse(targetOrgID)
	if err != nil {
		return nil, fmt.Errorf("invalid target_title_id: %w", err)
	}
	if s.trainingProvider == nil {
		return nil, fmt.Errorf("training recommendation belum dikonfigurasi (training provider)")
	}

	targetTitle, _ := s.repo.GetPositionTitle(ctx, targetID)

	requirements, err := s.repo.GetOrgCompetencyRequirements(ctx, targetID)
	if err != nil {
		return nil, err
	}
	if len(requirements) == 0 {
		return nil, fmt.Errorf("target position belum memiliki data assessment competency (belum ada standar kompetensi tercatat untuk organisasi ini)")
	}

	employeeLevels, err := s.repo.GetEmployeeCompetencyLevels(ctx, empID)
	if err != nil {
		return nil, err
	}

	gapCompetencyIDs := make([]uuid.UUID, 0)
	items := make([]TrainingRecommendationItem, 0)
	for _, r := range requirements {
		level, has := employeeLevels[r.CompetencyID]
		if has && level >= r.StandardLevel {
			continue // met -- no recommendation needed
		}
		gap := r.StandardLevel - level // level is 0 when !has
		priority := "MEDIUM"
		if gap >= 2 {
			priority = "HIGH"
		} else if gap <= 0 {
			priority = "LOW"
		}
		items = append(items, TrainingRecommendationItem{
			CompetencyID:   r.CompetencyID.String(),
			CompetencyName: r.CompetencyName,
			CurrentLevel:   level,
			RequiredLevel:  r.StandardLevel,
			Gap:            gap,
			Priority:       priority,
		})
		gapCompetencyIDs = append(gapCompetencyIDs, r.CompetencyID)
	}

	if len(gapCompetencyIDs) > 0 {
		courses, err := s.trainingProvider.ListCoursesByCompetencyIDs(ctx, gapCompetencyIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to read course recommendations: %w", err)
		}
		courseByCompetency := make(map[string]RecommendedCourse, len(courses))
		for _, c := range courses {
			// First course found for a competency wins -- good enough for
			// P0; ranking multiple candidate courses is out of scope here.
			if _, exists := courseByCompetency[c.CompetencyID]; !exists {
				courseByCompetency[c.CompetencyID] = c
			}
		}
		for i := range items {
			if c, ok := courseByCompetency[items[i].CompetencyID]; ok {
				items[i].CourseID = c.CourseID
				items[i].CourseName = c.CourseName
				items[i].IsMandatory = c.IsMandatory
				items[i].IsCertified = c.IsCertified
			}
		}
	}

	return &TrainingRecommendationResponse{
		EmployeeID:      employeeID,
		TargetTitle:     targetTitle,
		Recommendations: items,
	}, nil
}
```

- [ ] **Step 5: Run the tests to verify they pass**

Run: `cd backend && go test ./internal/modules/careerintelligence/... -run TestService_GetTrainingRecommendations -v`
Expected: PASS (both tests).

Note: the success test's fake course has `CompetencyID: ""` which will never match the real gapped competency's UUID from the fixture, so `CourseName` would NOT actually populate with the fake setup as written above. Fix the test fixture before running: the fake's `CompetencyID` must equal the Communication competency's real UUID from `seedGapAnalysisCompetencyData`. Since that fixture generates the competency UUID internally and doesn't return it, either (a) extend `seedGapAnalysisCompetencyData` in `helpers_test.go` to return the two competency UUIDs it creates (breaking change to its signature, update Task's own reference in `GetGapAnalysis`'s existing tests too), or (b) don't assert `CourseName` in this test and instead assert course-matching behavior separately with a fully self-contained fixture (not reusing `seedGapAnalysisCompetencyData`) that constructs its own known competency UUID end-to-end. Prefer (b) — write a dedicated fixture inline in the test rather than modifying the shared helper's signature (avoids touching Gap Analysis's existing passing tests). Rewrite Step 2's success test to seed its own data directly:

```go
func TestService_GetTrainingRecommendations_Success(t *testing.T) {
	svc, repo, db, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	orgID := uuid.New()
	leadershipID := uuid.New()
	communicationID := uuid.New()

	if err := db.Exec("CREATE TABLE IF NOT EXISTS competencies (id CHAR(36) PRIMARY KEY, name VARCHAR(255))").Error; err != nil {
		t.Fatalf("failed to create competencies table: %v", err)
	}
	if err := db.Exec("CREATE TABLE IF NOT EXISTS competency_scores (id CHAR(36) PRIMARY KEY, organization_id CHAR(36), employee_id CHAR(36), assessed_at TIMESTAMP)").Error; err != nil {
		t.Fatalf("failed to create competency_scores table: %v", err)
	}
	if err := db.Exec("CREATE TABLE IF NOT EXISTS competency_score_details (id CHAR(36) PRIMARY KEY, competency_score_id CHAR(36), competency_id CHAR(36), standard_level INT, employee_level INT)").Error; err != nil {
		t.Fatalf("failed to create competency_score_details table: %v", err)
	}
	db.Table("competencies").Create([]map[string]interface{}{
		{"id": leadershipID.String(), "name": "Leadership"},
		{"id": communicationID.String(), "name": "Communication"},
	})
	orgScoreID := uuid.New()
	db.Table("competency_scores").Create(map[string]interface{}{"id": orgScoreID.String(), "organization_id": orgID.String(), "assessed_at": "2026-01-01"})
	db.Table("competency_score_details").Create([]map[string]interface{}{
		{"id": uuid.New().String(), "competency_score_id": orgScoreID.String(), "competency_id": leadershipID.String(), "standard_level": 4},
		{"id": uuid.New().String(), "competency_score_id": orgScoreID.String(), "competency_id": communicationID.String(), "standard_level": 3},
	})
	empScoreID := uuid.New()
	db.Table("competency_scores").Create(map[string]interface{}{"id": empScoreID.String(), "organization_id": uuid.New().String(), "employee_id": employeeID.String(), "assessed_at": "2026-01-05"})
	db.Table("competency_score_details").Create([]map[string]interface{}{
		{"id": uuid.New().String(), "competency_score_id": empScoreID.String(), "competency_id": leadershipID.String(), "employee_level": 4},
		{"id": uuid.New().String(), "competency_score_id": empScoreID.String(), "competency_id": communicationID.String(), "employee_level": 2},
	})

	fake := &fakeTrainingProvider{
		courses: []RecommendedCourse{
			{CourseID: uuid.New().String(), CourseName: "Communication Skills", CompetencyID: communicationID.String(), IsMandatory: false, IsCertified: true},
		},
	}
	svc.SetTrainingProvider(fake)

	resp, err := svc.GetTrainingRecommendations(ctx(), employeeID.String(), orgID.String())
	if err != nil {
		t.Fatalf("GetTrainingRecommendations failed: %v", err)
	}
	if len(resp.Recommendations) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(resp.Recommendations))
	}
	item := resp.Recommendations[0]
	if item.Gap != 1 {
		t.Errorf("expected gap 1, got %d", item.Gap)
	}
	if item.CourseName != "Communication Skills" {
		t.Errorf("expected course 'Communication Skills', got '%s'", item.CourseName)
	}
	_ = repo // repo unused directly in this test; kept for newTestService's signature
}
```

- [ ] **Step 6: Add the handler**

In `backend/internal/modules/careerintelligence/handler.go`, add (near `GetGapAnalysis`):

```go
func (h *Handler) GetTrainingRecommendations(c *gin.Context) {
	employeeID := c.Param("employeeId")
	targetTitleID := c.Query("target_title_id")
	if targetTitleID == "" {
		respondError(c, http.StatusBadRequest, "target_title_id is required")
		return
	}
	result, err := h.svc.GetTrainingRecommendations(c.Request.Context(), employeeID, targetTitleID)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	respondSuccess(c, result)
}
```

- [ ] **Step 7: Register the route**

In `backend/internal/modules/careerintelligence/routes.go`, add next to the `training-profile` route from Task 4:

```go
ci.GET("/employees/:employeeId/training-recommendations", handler.GetTrainingRecommendations)
```

- [ ] **Step 8: Build and run the full careerintelligence test suite**

Run: `cd backend && go build ./... && go vet ./... && go test ./internal/modules/careerintelligence/...`
Expected: build clean, vet clean, all tests PASS.

- [ ] **Step 9: Commit**

```bash
cd backend
git add internal/modules/careerintelligence/dto.go internal/modules/careerintelligence/service.go internal/modules/careerintelligence/service_test.go internal/modules/careerintelligence/handler.go internal/modules/careerintelligence/routes.go
git commit -m "$(cat <<'EOF'
feat(career-intelligence): endpoint Training & Competency Gap + Recommendation

GET /career-intelligence/employees/:employeeId/training-recommendations?target_title_id=...
-- gap kompetensi (reuse data yang sama persis dengan GetGapAnalysis) plus
course yang menutup tiap gap (via TrainingProvider.ListCoursesByCompetencyIDs,
Task 2/3). Priority HIGH/MEDIUM/LOW pakai band yang sama dengan GetGapAnalysis.
Career Intelligence Training Enhancement plan §7/§8
(docs/career-intelligence-training-enhancement-plan.md).

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>
EOF
)"
```

---

### Task 6: Frontend — surface training profile + recommendations in Gap Analysis page

**Files:**
- Modify: `frontend/tenant/src/views/modules/career-intelligence/GapAnalysis.vue`
- Modify: `frontend/tenant/src/locales/en.json`
- Modify: `frontend/tenant/src/locales/id.json`

**Interfaces:**
- Consumes: `GET /api/v1/tenant/career-intelligence/employees/:employeeId/training-profile` (Task 4), `GET /api/v1/tenant/career-intelligence/employees/:employeeId/training-recommendations?target_title_id=...` (Task 5).

- [ ] **Step 1: Read the current file to find the exact analyze-result section**

Read `frontend/tenant/src/views/modules/career-intelligence/GapAnalysis.vue` in full before editing — this task's exact insertion point depends on that file's current `analyze()` function and result-rendering template block, which must be read fresh (do not guess line numbers here; the file was last touched in an earlier session and may have shifted).

- [ ] **Step 2: Add training profile + recommendations to the analyze flow**

In the `<script setup>` block, alongside the existing `analyze()` function that calls `/paths/gap-analysis`, add two more refs and fetch calls fired in parallel with the existing gap analysis call (use `Promise.allSettled` so a training-side failure doesn't block the existing competency gap result from rendering):

```js
const trainingProfile = ref(null)
const trainingRecommendations = ref(null)

async function loadTrainingData(employeeId, targetTitleId) {
  const [profileRes, recRes] = await Promise.allSettled([
    api.get(`/api/v1/tenant/career-intelligence/employees/${employeeId}/training-profile`),
    api.get(`/api/v1/tenant/career-intelligence/employees/${employeeId}/training-recommendations`, { params: { target_title_id: targetTitleId } })
  ])
  trainingProfile.value = profileRes.status === 'fulfilled' ? (profileRes.value.data?.data || null) : null
  trainingRecommendations.value = recRes.status === 'fulfilled' ? (recRes.value.data?.data || null) : null
}
```

Call `loadTrainingData(form.value.employee_id, form.value.target_title_id)` at the same point the existing `analyze()` function calls the gap-analysis endpoint (read the existing function first — call it right after that existing API call succeeds, inside the same handler, not as a separate button).

- [ ] **Step 3: Render the training profile summary**

In the template, in the results section (after the existing gap-analysis result block), add:

```html
<div v-if="trainingProfile" class="rounded-lg border border-gray-200 dark:border-gray-700 p-4 mt-4">
  <p class="text-xs uppercase tracking-wide text-gray-400 font-medium mb-3">{{ t('gap_analysis.training_profile') }}</p>
  <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 text-sm">
    <div><span class="text-gray-400">{{ t('gap_analysis.total_training') }}</span><p class="font-semibold text-navy-800 dark:text-gray-100">{{ trainingProfile.total_training }}</p></div>
    <div><span class="text-gray-400">{{ t('gap_analysis.completed') }}</span><p class="font-semibold text-navy-800 dark:text-gray-100">{{ trainingProfile.completed }}</p></div>
    <div><span class="text-gray-400">{{ t('gap_analysis.training_hours') }}</span><p class="font-semibold text-navy-800 dark:text-gray-100">{{ trainingProfile.training_hours }}</p></div>
    <div><span class="text-gray-400">{{ t('gap_analysis.average_score') }}</span><p class="font-semibold text-navy-800 dark:text-gray-100">{{ trainingProfile.average_score?.toFixed?.(0) ?? trainingProfile.average_score }}</p></div>
  </div>
</div>
```

- [ ] **Step 4: Render training recommendations per gapped competency**

In the same template section, add:

```html
<div v-if="trainingRecommendations?.recommendations?.length" class="rounded-lg border border-gray-200 dark:border-gray-700 p-4 mt-4">
  <p class="text-xs uppercase tracking-wide text-gray-400 font-medium mb-3">{{ t('gap_analysis.recommended_training') }}</p>
  <div class="space-y-2">
    <div v-for="item in trainingRecommendations.recommendations" :key="item.competency_id" class="flex items-center justify-between text-sm border-b border-gray-100 dark:border-gray-700 pb-2 last:border-0 last:pb-0">
      <div>
        <p class="font-medium text-navy-800 dark:text-gray-100">{{ item.competency_name }}</p>
        <p class="text-xs text-gray-400">{{ t('gap_analysis.gap_level', { current: item.current_level, required: item.required_level }) }}</p>
      </div>
      <div class="text-right">
        <Tag :value="item.priority" :severity="item.priority === 'HIGH' ? 'danger' : item.priority === 'MEDIUM' ? 'warning' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" />
        <p v-if="item.course_name" class="text-xs text-teal-600 dark:text-teal-400 mt-1">{{ item.course_name }}</p>
        <p v-else class="text-xs text-gray-400 mt-1">{{ t('gap_analysis.no_course_available') }}</p>
      </div>
    </div>
  </div>
</div>
```

- [ ] **Step 5: Add the i18n keys**

In `frontend/tenant/src/locales/en.json`, inside the existing `"gap_analysis": {` block, add:

```json
"training_profile": "Training Profile",
"total_training": "Total Training",
"completed": "Completed",
"training_hours": "Training Hours",
"average_score": "Average Score",
"recommended_training": "Recommended Training",
"gap_level": "Current level {current}, target {required}",
"no_course_available": "No matching course found"
```

In `frontend/tenant/src/locales/id.json`, inside the existing `"gap_analysis": {` block, add:

```json
"training_profile": "Profil Training",
"total_training": "Total Training",
"completed": "Selesai",
"training_hours": "Jam Training",
"average_score": "Skor Rata-rata",
"recommended_training": "Training Direkomendasikan",
"gap_level": "Level saat ini {current}, target {required}",
"no_course_available": "Belum ada course yang cocok"
```

- [ ] **Step 6: Build the frontend**

Run: `cd frontend/tenant && npm run build`
Expected: build succeeds with no errors.

- [ ] **Step 7: Manual smoke test**

Run: `cd frontend/tenant && npm run dev` and open `/career-intelligence/gap-analysis` in a browser. Select an employee and a target position that has existing competency assessment data (per this session's earlier finding, the dev tenant DB currently has none — seed at least one `competency_scores` row for a test organization first, or note in the PR that this can only be smoke-tested once real assessment data exists). Confirm: analyzing shows the existing gap-analysis result unchanged, plus the new training profile numbers and recommended-training list beneath it, and that a missing/failed training-side call doesn't break the page (per the `Promise.allSettled` in Step 2).

- [ ] **Step 8: Commit**

```bash
cd frontend/tenant
git add src/views/modules/career-intelligence/GapAnalysis.vue src/locales/en.json src/locales/id.json
git commit -m "$(cat <<'EOF'
feat(career-intelligence): tampilkan training profile + rekomendasi di Gap Analysis

GapAnalysis.vue sekarang juga menampilkan ringkasan training employee dan
daftar course yang direkomendasikan untuk menutup tiap competency gap
(GET .../training-profile dan .../training-recommendations, Task 4/5).
Dimuat paralel dengan gap analysis via Promise.allSettled supaya
kegagalan sisi training tidak menghalangi hasil gap analysis existing.

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>
EOF
)"
```

---

## Self-Review Notes

**Spec coverage** (against `docs/career-intelligence-training-enhancement-plan.md` §27 P0 list):
1. Employee Training Profile → Task 1 (backend) + Task 4 (Career Intelligence endpoint) + Task 6 (FE display).
2. Training History for Career → Task 4 (`history` field on the profile response, sourced from the already-existing `GetTrainingHistory`).
3. Training ↔ Competency Analysis / Career Training Gap → Task 5.
4. Career Training Recommendation → Task 5 (same endpoint, `recommendations[].course_*` fields).
5. Career Development Plan integration → explicitly deferred (see "Deferred Scope").
6. Career Readiness (P1, not P0 per the doc's own numbering ambiguity — treated as deferred here since it needs a new weighting/settings entity) → explicitly deferred.

**Out of scope reminder for whoever picks this plan up next:** Talent Mapping evidence (§15), Succession Planning training readiness columns (§16), certification eligibility (§12), and training-effectiveness-as-evidence (§13) are P1/P2 in the spec's own priority table and are not covered by this plan — write a follow-up plan for those once P0 is live and validated with real data.
