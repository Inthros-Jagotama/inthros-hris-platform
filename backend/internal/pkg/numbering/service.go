package numbering

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrInvalidDocumentType = errors.New("invalid document type")
	ErrInvalidResetPeriod  = errors.New("invalid reset period")
	ErrSettingNotFound     = errors.New("numbering setting not found")
)

var validDocumentTypes = map[string]bool{
	DocumentTypeEmployeeMovement: true,
	DocumentTypeEmployeeContract: true,
}

var validResetPeriods = map[string]bool{
	ResetPeriodYearly:  true,
	ResetPeriodMonthly: true,
	ResetPeriodNever:   true,
}

type Service struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
	logger     *zap.Logger
	now        func() time.Time
}

func NewService(dbResolver func(ctx context.Context) (*gorm.DB, error), logger *zap.Logger) *Service {
	return &Service{dbResolver: dbResolver, logger: logger, now: time.Now}
}

func (s *Service) getDB(ctx context.Context) (*gorm.DB, error) {
	return s.dbResolver(ctx)
}

func (s *Service) List(ctx context.Context) ([]DocumentNumberingSetting, error) {
	db, err := s.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []DocumentNumberingSetting
	if err := db.Order("document_type").Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to list numbering settings: %w", err)
	}
	return items, nil
}

func (s *Service) Update(ctx context.Context, documentType, formatTemplate, resetPeriod string) (*DocumentNumberingSetting, error) {
	if !validDocumentTypes[documentType] {
		return nil, ErrInvalidDocumentType
	}
	if !validResetPeriods[resetPeriod] {
		return nil, ErrInvalidResetPeriod
	}
	db, err := s.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var setting DocumentNumberingSetting
	if err := db.Where("document_type = ?", documentType).First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSettingNotFound
		}
		return nil, fmt.Errorf("failed to load numbering setting: %w", err)
	}
	setting.FormatTemplate = formatTemplate
	setting.ResetPeriod = resetPeriod
	if err := db.Save(&setting).Error; err != nil {
		return nil, fmt.Errorf("failed to update numbering setting: %w", err)
	}
	return &setting, nil
}

// Preview formats what the next Generate() call would return, without
// mutating last_sequence — used by the settings UI to show a live example.
func (s *Service) Preview(ctx context.Context, documentType string) (string, error) {
	if !validDocumentTypes[documentType] {
		return "", ErrInvalidDocumentType
	}
	db, err := s.getDB(ctx)
	if err != nil {
		return "", err
	}
	var setting DocumentNumberingSetting
	if err := db.Where("document_type = ?", documentType).First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrSettingNotFound
		}
		return "", fmt.Errorf("failed to load numbering setting: %w", err)
	}
	now := s.now()
	nextSeq := setting.LastSequence + 1
	if ResetKeyFor(setting.ResetPeriod, now) != setting.LastResetKey {
		nextSeq = 1
	}
	return FormatTemplate(setting.FormatTemplate, nextSeq, now), nil
}

// Generate atomically increments (and resets, if the period rolled over)
// the sequence for documentType and returns the formatted number. Must be
// safe under concurrent callers, hence the row lock.
//
// Note: Generate commits its own transaction rather than joining the
// caller's — if the enclosing movement/contract create fails after this
// call, the consumed sequence number is not reclaimed. This is a
// documented, accepted tradeoff, not a bug.
func (s *Service) Generate(ctx context.Context, documentType string) (string, error) {
	if !validDocumentTypes[documentType] {
		return "", ErrInvalidDocumentType
	}
	db, err := s.getDB(ctx)
	if err != nil {
		return "", err
	}

	var result string
	txErr := db.Transaction(func(tx *gorm.DB) error {
		var setting DocumentNumberingSetting
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("document_type = ?", documentType).
			First(&setting).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrSettingNotFound
			}
			return fmt.Errorf("failed to load numbering setting: %w", err)
		}

		now := s.now()
		resetKey := ResetKeyFor(setting.ResetPeriod, now)
		if resetKey != setting.LastResetKey {
			setting.LastSequence = 0
			setting.LastResetKey = resetKey
		}
		setting.LastSequence++

		if err := tx.Save(&setting).Error; err != nil {
			return fmt.Errorf("failed to persist numbering sequence: %w", err)
		}

		result = FormatTemplate(setting.FormatTemplate, setting.LastSequence, now)
		return nil
	})
	if txErr != nil {
		if s.logger != nil {
			s.logger.Warn("numbering: failed to generate document number",
				zap.String("document_type", documentType), zap.Error(txErr))
		}
		return "", txErr
	}
	return result, nil
}
