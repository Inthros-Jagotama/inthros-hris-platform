package notification

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
}

func NewRepository(dbResolver func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbResolver: dbResolver}
}

func (r *Repository) getDB(ctx context.Context) (*gorm.DB, error) {
	return r.dbResolver(ctx)
}

func (r *Repository) CreateNotification(ctx context.Context, n *Notification) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(n).Error
}

func (r *Repository) FindNotificationByID(ctx context.Context, id uuid.UUID) (*Notification, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var n Notification
	if err := db.First(&n, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("notification not found: %w", err)
	}
	return &n, nil
}

func (r *Repository) ListNotificationsByRecipient(ctx context.Context, recipientUserID uuid.UUID, isRead *bool, page, perPage int) ([]Notification, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var notifications []Notification
	var total int64
	query := db.Model(&Notification{}).Where("recipient_user_id = ?", recipientUserID)
	if isRead != nil {
		query = query.Where("is_read = ?", *isRead)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, 0, err
	}
	return notifications, total, nil
}

func (r *Repository) UpdateNotification(ctx context.Context, n *Notification) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(n).Error
}
