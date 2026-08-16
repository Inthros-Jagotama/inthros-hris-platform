package attendance

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// =========================================================================
// Business Travel (core)
// =========================================================================

func (r *Repository) CreateBusinessTravel(ctx context.Context, t *BusinessTravel) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(t).Error
}

func (r *Repository) FindBusinessTravelByID(ctx context.Context, id uuid.UUID) (*BusinessTravel, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var t BusinessTravel
	if err := db.Preload("Participants").Preload("Destinations").First(&t, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("business travel not found: %w", err)
	}
	return &t, nil
}

func (r *Repository) UpdateBusinessTravel(ctx context.Context, t *BusinessTravel) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(t).Error
}

func (r *Repository) ListBusinessTravels(ctx context.Context, requesterID *uuid.UUID, status string, page, perPage int) ([]BusinessTravel, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var travels []BusinessTravel
	var total int64
	query := db.Model(&BusinessTravel{})
	if requesterID != nil {
		query = query.Where("requester_id = ?", *requesterID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&travels).Error; err != nil {
		return nil, 0, err
	}
	return travels, total, nil
}

// =========================================================================
// Participants
// =========================================================================

func (r *Repository) CreateParticipant(ctx context.Context, p *BusinessTravelParticipant) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(p).Error
}

func (r *Repository) ListParticipantsByTravel(ctx context.Context, travelID uuid.UUID) ([]BusinessTravelParticipant, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var participants []BusinessTravelParticipant
	if err := db.Where("business_travel_id = ?", travelID).Find(&participants).Error; err != nil {
		return nil, err
	}
	return participants, nil
}

// =========================================================================
// Destinations
// =========================================================================

func (r *Repository) CreateDestination(ctx context.Context, d *BusinessTravelDestination) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(d).Error
}

func (r *Repository) ListDestinationsByTravel(ctx context.Context, travelID uuid.UUID) ([]BusinessTravelDestination, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var destinations []BusinessTravelDestination
	if err := db.Where("business_travel_id = ?", travelID).Order("sequence ASC").Find(&destinations).Error; err != nil {
		return nil, err
	}
	return destinations, nil
}
