package recruitment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"

	"github.com/inthros/hris-platform/internal/modules/setting"
)

func setupTestDB() (*gorm.DB, func(ctx context.Context) (*gorm.DB, error), func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}

	if err := db.AutoMigrate(
		&JobRequisition{},
		&Candidate{},
		&JobApplication{},
		&Interview{},
		&JobOffer{},
		&OnboardingTaskTemplate{},
		&EmployeeOnboarding{},
		&OnboardingTaskItem{},
		&RecruitmentStage{},
		&ApplicationStageHistory{},
		&CandidateEducation{},
		&CandidateWorkExperience{},
		&setting.EducationMajor{},
	); err != nil {
		panic(fmt.Sprintf("failed to migrate test db: %v", err))
	}

	dbResolver := func(ctx context.Context) (*gorm.DB, error) {
		return db, nil
	}

	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	return db, dbResolver, cleanup
}

// seedDefaultRecruitmentStages seeds the standard recruitment stage catalog
// (G-5 state machine). Required for tests that exercise
// transitionApplicationStatus / CreateApplication initial history /
// GetApplicationHistory, since setupTestDB only migrates the schema and
// does not seed rows. Not called by setupTestDB itself because some
// repository tests create their own stage rows and would collide with the
// unique index on recruitment_stages.code.
func seedDefaultRecruitmentStages(db *gorm.DB) {
	defaultStages := []RecruitmentStage{
		{Code: "NEW", Name: "New", SortOrder: 1},
		{Code: "SCREENED", Name: "Screened", SortOrder: 2},
		{Code: "SHORTLISTED", Name: "Shortlisted", SortOrder: 3},
		{Code: "INTERVIEWED", Name: "Interviewed", SortOrder: 4},
		{Code: "OFFERED", Name: "Offered", SortOrder: 5},
		{Code: "ACCEPTED", Name: "Accepted", SortOrder: 6},
		{Code: "REJECTED", Name: "Rejected", SortOrder: 7},
		{Code: "WITHDRAWN", Name: "Withdrawn", SortOrder: 8},
	}
	for i := range defaultStages {
		if err := db.Create(&defaultStages[i]).Error; err != nil {
			panic(fmt.Sprintf("failed to seed recruitment stage %s: %v", defaultStages[i].Code, err))
		}
	}
}

func createTestUUID() string {
	return uuid.New().String()
}

func createTestOrgID() string {
	return uuid.New().String()
}

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
