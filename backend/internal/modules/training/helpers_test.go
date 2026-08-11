package training

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"
)

// testDB menyiapkan database SQLite in-memory unik per-panggilan.
func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	// Gunakan URI unik per test untuk mencegah sharing database antar test
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid.New().String())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(
		&TrainingCategory{},
		&TrainingCourse{},
		&TrainingSession{},
		&TrainingParticipant{},
		&TrainingMaterial{},
		&TrainingEvaluation{},
		&TrainingCertificate{},
		// P0-BE models
		&TrainingProvider{},
		&TrainingTrainer{},
		&TrainingSessionTrainer{},
		&TrainingAttendance{},
		&TrainingAssessment{},
		&TrainingAssessmentResult{},
		// P1-BE models
		&TrainingPlan{},
		&TrainingPlanItem{},
		&TrainingNeed{},
		&TrainingRequest{},
		&TrainingCourseObjective{},
		&TrainingCourseCompetency{},
		&TrainingCoursePrerequisite{},
		&TrainingMandatory{},
		&TrainingSessionCost{},
		&TrainingDocument{},
		// P2-BE models
		&TrainingEvaluationForm{},
		&TrainingEvaluationQuestion{},
		&TrainingEvaluationAnswer{},
		&TrainingEffectivenessAssessment{},
		&TrainingCertification{},
	); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}
	return db
}

// testLogger menyiapkan logger untuk testing (discard output).
func testLogger() *zap.Logger {
	return zap.NewNop()
}

// testRepo menyiapkan repository dengan SQLite in-memory.
func testRepo(t *testing.T) *Repository {
	t.Helper()
	db := testDB(t)
	return &Repository{
		dbFunc: func(ctx context.Context) (*gorm.DB, error) {
			return db.WithContext(ctx), nil
		},
	}
}

// testSvc menyiapkan service dengan repository dan logger test.
func testSvc(t *testing.T) *Service {
	t.Helper()
	repo := testRepo(t)
	logger := testLogger()
	return NewService(repo, logger)
}

// testCtx menyiapkan context untuk testing.
func testCtx() context.Context {
	return context.Background()
}

// TestMain untuk setup global testing.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

// seedCategory menyiapkan data kategori untuk testing.
func seedCategory(t *testing.T, svc *Service) string {
	t.Helper()
	desc := "Test category"
	resp, err := svc.CreateCategory(testCtx(), CreateTrainingCategoryRequest{
		Code:        "TECH",
		Name:        "Technical",
		Description: &desc,
	})
	if err != nil {
		t.Fatalf("failed to seed category: %v", err)
	}
	return resp.ID
}

// seedCourse menyiapkan data course untuk testing.
func seedCourse(t *testing.T, svc *Service, categoryID string) string {
	t.Helper()
	dur := 8.0
	resp, err := svc.CreateCourse(testCtx(), CreateTrainingCourseRequest{
		CategoryID:   categoryID,
		Code:         "GOLANG-101",
		Name:         "Golang Fundamentals",
		DurationHour: &dur,
	})
	if err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}
	return resp.ID
}

// seedSession menyiapkan data session untuk testing.
func seedSession(t *testing.T, svc *Service, courseID string) string {
	t.Helper()
	resp, err := svc.CreateSession(testCtx(), CreateTrainingSessionRequest{
		CourseID:    courseID,
		SessionCode: "CLS-001",
		TrainerName: "John Doe",
		StartDate:   "2026-08-01",
		EndDate:     "2026-08-05",
		MaxQuota:    30,
	})
	if err != nil {
		t.Fatalf("failed to seed session: %v", err)
	}
	return resp.ID
}

// validatePagination adalah helper untuk validasi pagination response.
func validatePagination(resp *PaginatedResponse, expectedTotal int64) error {
	if !resp.Success {
		return fmt.Errorf("expected success=true, got false")
	}
	if resp.Total != expectedTotal {
		return fmt.Errorf("expected total=%d, got %d", expectedTotal, resp.Total)
	}
	return nil
}
