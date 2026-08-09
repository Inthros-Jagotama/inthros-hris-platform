package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
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
//
// Callers pass notifType + params only, not rendered text — the bilingual
// title/body are rendered per-request from the catalog in i18n.go, in
// whatever language the recipient is actually viewing in at the time (their
// own Accept-Language), which is unknowable at Notify()-time since the
// creator (e.g. an approver submitting an action) and the recipient are
// different people who may be using different languages. The English
// rendering computed here is stored purely as an audit/fallback default.
func (s *Service) Notify(ctx context.Context, recipientUserID uuid.UUID, notifType string, params []string, referenceType string, referenceID uuid.UUID) error {
	defaultTitle, defaultBody := translate(httputil.LangEN, notifType, params, notifType, "")

	var paramsJSON *string
	if len(params) > 0 {
		if b, err := json.Marshal(params); err == nil {
			s := string(b)
			paramsJSON = &s
		}
	}

	n := &Notification{
		RecipientUserID: recipientUserID,
		Type:            notifType,
		Title:           defaultTitle,
		Body:            defaultBody,
		Params:          paramsJSON,
		ReferenceType:   &referenceType,
		ReferenceID:     &referenceID,
	}
	return s.repo.CreateNotification(ctx, n)
}

// ListNotifications returns paginated notifications for the given recipient,
// optionally filtered by read status, with Title/Body re-rendered in `lang`
// (the recipient's own requested language) rather than whatever was stored
// at Notify()-time.
func (s *Service) ListNotifications(ctx context.Context, recipientUserID uuid.UUID, isRead *bool, page, perPage int, lang httputil.Lang) ([]Notification, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	notifications, total, err := s.repo.ListNotificationsByRecipient(ctx, recipientUserID, isRead, page, perPage)
	if err != nil {
		return nil, 0, err
	}
	for i := range notifications {
		applyTranslation(&notifications[i], lang)
	}
	return notifications, total, nil
}

// applyTranslation overrides a notification's Title/Body in place with the
// rendering for `lang`, decoding its stored Params (if any) for body
// placeholder substitution.
func applyTranslation(n *Notification, lang httputil.Lang) {
	var params []string
	if n.Params != nil && *n.Params != "" {
		_ = json.Unmarshal([]byte(*n.Params), &params)
	}
	n.Title, n.Body = translate(lang, n.Type, params, n.Title, n.Body)
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
