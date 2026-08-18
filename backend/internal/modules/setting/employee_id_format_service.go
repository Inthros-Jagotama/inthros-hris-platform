package setting

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/inthros/hris-platform/internal/pkg/numbering"
)

var (
	ErrInvalidEmployeeIDGenerationMode = errors.New("invalid employee id generation mode")
	ErrInvalidEmployeeIDResetPeriod    = errors.New("invalid employee id reset period")
)

var validEmployeeIDGenerationModes = map[string]bool{
	EmployeeIDGenerationModeManual: true,
	EmployeeIDGenerationModeHybrid: true,
	EmployeeIDGenerationModeAuto:   true,
}

var validEmployeeIDResetPeriods = map[string]bool{
	numbering.ResetPeriodYearly:  true,
	numbering.ResetPeriodMonthly: true,
	numbering.ResetPeriodNever:   true,
}

// EmployeeIDFormatService owns the employee_id generation-mode setting and
// the atomic sequence generator behind AUTO/HYBRID mode. It is a parallel,
// smaller counterpart to numbering.Service — it reuses numbering's pure
// FormatTemplate/ResetKeyFor helpers but keeps its own table, since
// employee_id's MANUAL/HYBRID/AUTO mode concept doesn't fit numbering's
// closed document-type enum.
type EmployeeIDFormatService struct {
	repo       *EmployeeIDFormatRepository
	dbResolver TenantDBFunc
	logger     *zap.Logger
	now        func() time.Time
}

func NewEmployeeIDFormatService(repo *EmployeeIDFormatRepository, dbResolver TenantDBFunc, logger *zap.Logger) *EmployeeIDFormatService {
	return &EmployeeIDFormatService{repo: repo, dbResolver: dbResolver, logger: logger, now: time.Now}
}

// GetSetting returns the current setting (auto-creating a default row if
// none exists yet).
func (s *EmployeeIDFormatService) GetSetting(ctx context.Context) (*EmployeeIDFormatSetting, error) {
	return s.repo.GetOrCreate(ctx)
}

// UpdateSetting validates and persists the generation mode, format template,
// and reset period. It never touches last_sequence/last_reset_key.
func (s *EmployeeIDFormatService) UpdateSetting(ctx context.Context, mode, formatTemplate, resetPeriod string) (*EmployeeIDFormatSetting, error) {
	if !validEmployeeIDGenerationModes[mode] {
		return nil, ErrInvalidEmployeeIDGenerationMode
	}
	if !validEmployeeIDResetPeriods[resetPeriod] {
		return nil, ErrInvalidEmployeeIDResetPeriod
	}

	setting, err := s.repo.GetOrCreate(ctx)
	if err != nil {
		return nil, err
	}
	setting.GenerationMode = mode
	setting.FormatTemplate = formatTemplate
	setting.ResetPeriod = resetPeriod
	if err := s.repo.Update(ctx, setting); err != nil {
		return nil, err
	}
	return setting, nil
}

// Preview formats what the next Generate() call would return, without
// mutating last_sequence/last_reset_key — read-only, no row lock, no save.
func (s *EmployeeIDFormatService) Preview(ctx context.Context) (string, error) {
	setting, err := s.repo.GetOrCreate(ctx)
	if err != nil {
		return "", err
	}
	now := s.now()
	nextSeq := setting.LastSequence + 1
	if numbering.ResetKeyFor(setting.ResetPeriod, now) != setting.LastResetKey {
		nextSeq = 1
	}
	return numbering.FormatTemplate(setting.FormatTemplate, nextSeq, now), nil
}

// Generate atomically increments (and resets, if the period rolled over) the
// employee_id sequence and returns the formatted ID. Safe under concurrent
// callers via row lock.
//
// Note: like numbering.Service.Generate, this commits its own transaction
// rather than joining the caller's (employee create) transaction — if
// CreateEmployee fails after this call, the consumed sequence number is not
// reclaimed. This is a documented, accepted tradeoff, not a bug.
func (s *EmployeeIDFormatService) Generate(ctx context.Context) (string, error) {
	db, err := s.dbResolver(ctx)
	if err != nil {
		return "", err
	}

	var result string
	txErr := db.Transaction(func(tx *gorm.DB) error {
		var setting EmployeeIDFormatSetting
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&setting).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// No row yet — create default inline within the same tx so
				// Generate() works even before any GET has touched the setting.
				setting = EmployeeIDFormatSetting{
					ID:             uuid.NewString(),
					GenerationMode: EmployeeIDGenerationModeManual,
					FormatTemplate: defaultEmployeeIDFormatTemplate,
					ResetPeriod:    defaultEmployeeIDResetPeriod,
				}
				if genErr := tx.Create(&setting).Error; genErr != nil {
					return fmt.Errorf("failed to create default employee id format setting: %w", genErr)
				}
			} else {
				return fmt.Errorf("failed to load employee id format setting: %w", err)
			}
		}

		now := s.now()
		resetKey := numbering.ResetKeyFor(setting.ResetPeriod, now)
		if resetKey != setting.LastResetKey {
			setting.LastSequence = 0
			setting.LastResetKey = resetKey
		}
		setting.LastSequence++

		if err := tx.Save(&setting).Error; err != nil {
			return fmt.Errorf("failed to persist employee id sequence: %w", err)
		}

		result = numbering.FormatTemplate(setting.FormatTemplate, setting.LastSequence, now)
		return nil
	})
	if txErr != nil {
		if s.logger != nil {
			s.logger.Warn("employee_id_format: failed to generate employee id", zap.Error(txErr))
		}
		return "", txErr
	}
	return result, nil
}
