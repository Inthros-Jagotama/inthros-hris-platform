package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// maxMarkAllPageSize bounds how many unread notifications MarkAllAsRead
// processes in a single call.
const maxMarkAllPageSize = 10000

type Service struct {
	repo   *Repository
	logger *zap.Logger
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// Notify creates a new notification for a recipient. It is called by other
// modules through their local Notifier interface — there is no public HTTP
// endpoint for creating notifications.
func (s *Service) Notify(ctx context.Context, recipientUserID uuid.UUID, notifType, title, body, referenceType string, referenceID uuid.UUID) error {
	n := &Notification{
		RecipientUserID: recipientUserID,
		Type:            notifType,
		Title:           title,
		Body:            body,
		ReferenceType:   &referenceType,
		ReferenceID:     &referenceID,
	}
	return s.repo.CreateNotification(ctx, n)
}

// ListNotifications returns paginated notifications for the given recipient,
// optionally filtered by read status.
func (s *Service) ListNotifications(ctx context.Context, recipientUserID uuid.UUID, isRead *bool, page, perPage int) ([]Notification, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	return s.repo.ListNotificationsByRecipient(ctx, recipientUserID, isRead, page, perPage)
}

// MarkAsRead marks a single notification as read. Only the recipient of the
// notification is allowed to mark it as read.
func (s *Service) MarkAsRead(ctx context.Context, notificationID, recipientUserID uuid.UUID) error {
	n, err := s.repo.FindNotificationByID(ctx, notificationID)
	if err != nil {
		return err
	}
	if n.RecipientUserID != recipientUserID {
		return fmt.Errorf("notification does not belong to this user")
	}
	if n.IsRead {
		return nil
	}
	now := time.Now()
	n.IsRead = true
	n.ReadAt = &now
	return s.repo.UpdateNotification(ctx, n)
}

// MarkAllAsRead marks every unread notification belonging to recipientUserID as read.
func (s *Service) MarkAllAsRead(ctx context.Context, recipientUserID uuid.UUID) error {
	unread := false
	notifications, _, err := s.repo.ListNotificationsByRecipient(ctx, recipientUserID, &unread, 1, maxMarkAllPageSize)
	if err != nil {
		return err
	}
	now := time.Now()
	for i := range notifications {
		notifications[i].IsRead = true
		notifications[i].ReadAt = &now
		if err := s.repo.UpdateNotification(ctx, &notifications[i]); err != nil {
			return err
		}
	}
	return nil
}

// GetUnreadCount returns the number of unread notifications for recipientUserID.
func (s *Service) GetUnreadCount(ctx context.Context, recipientUserID uuid.UUID) (int64, error) {
	unread := false
	_, total, err := s.repo.ListNotificationsByRecipient(ctx, recipientUserID, &unread, 1, 1)
	if err != nil {
		return 0, err
	}
	return total, nil
}
