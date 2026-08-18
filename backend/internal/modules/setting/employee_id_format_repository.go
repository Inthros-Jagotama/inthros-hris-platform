package setting

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	defaultEmployeeIDFormatTemplate = "EMP-{year}-{sequence:4}"
	defaultEmployeeIDResetPeriod    = "yearly"
)

// EmployeeIDFormatRepository persists the singleton employee_id_format_settings
// row for the current tenant.
type EmployeeIDFormatRepository struct {
	dbResolver TenantDBFunc
}

func NewEmployeeIDFormatRepository(dbResolver TenantDBFunc) *EmployeeIDFormatRepository {
	return &EmployeeIDFormatRepository{dbResolver: dbResolver}
}

func (r *EmployeeIDFormatRepository) getDB(ctx context.Context) (*gorm.DB, error) {
	return r.dbResolver(ctx)
}

// GetOrCreate returns the tenant's employee_id_format_settings row, creating
// a default row (MANUAL / EMP-{year}-{sequence:4} / yearly) on first read if
// none exists yet — same auto-create-default-on-read convention as other
// tenant singleton settings in this codebase.
func (r *EmployeeIDFormatRepository) GetOrCreate(ctx context.Context) (*EmployeeIDFormatSetting, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}

	var setting EmployeeIDFormatSetting
	err = db.First(&setting).Error
	if err == nil {
		return &setting, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to load employee id format setting: %w", err)
	}

	setting = EmployeeIDFormatSetting{
		ID:             uuid.NewString(),
		GenerationMode: EmployeeIDGenerationModeManual,
		FormatTemplate: defaultEmployeeIDFormatTemplate,
		ResetPeriod:    defaultEmployeeIDResetPeriod,
		LastSequence:   0,
	}
	if err := db.Create(&setting).Error; err != nil {
		return nil, fmt.Errorf("failed to create default employee id format setting: %w", err)
	}
	return &setting, nil
}

// Update persists changes to the singleton row (caller must have populated
// ID via a prior GetOrCreate).
func (r *EmployeeIDFormatRepository) Update(ctx context.Context, setting *EmployeeIDFormatSetting) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Save(setting).Error; err != nil {
		return fmt.Errorf("failed to update employee id format setting: %w", err)
	}
	return nil
}
