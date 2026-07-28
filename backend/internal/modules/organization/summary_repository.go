package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) CreateSummary(ctx context.Context, s *OrganizationSummary) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(s).Error
}

func (r *Repository) FindSummaryByID(ctx context.Context, id uuid.UUID) (*OrganizationSummary, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var s OrganizationSummary
	if err := db.Preload("Organizations").First(&s, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("organization summary not found: %w", err)
	}
	return &s, nil
}

func (r *Repository) FindAllSummaries(ctx context.Context, page, perPage int) ([]OrganizationSummary, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var summaries []OrganizationSummary
	var total int64

	query := db.Model(&OrganizationSummary{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&summaries).Error; err != nil {
		return nil, 0, err
	}

	return summaries, total, nil
}

func (r *Repository) UpdateSummary(ctx context.Context, s *OrganizationSummary) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(s).Error
}

func (r *Repository) SoftDeleteSummary(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Where("id = ?", id).Delete(&OrganizationSummary{}).Error
}
