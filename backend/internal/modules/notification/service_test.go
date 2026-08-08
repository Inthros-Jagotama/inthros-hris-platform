package notification

import (
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func newTestService() (*Service, *Repository, func()) {
	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	return svc, repo, func() {
		cleanup()
		_ = logger.Sync()
	}
}

func TestService_Notify_CreatesUnreadNotification(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	recipient := uuid.New()
	refID := uuid.New()
	if err := svc.Notify(ctx(), recipient, "LEAVE_APPROVED", "Leave Approved", "Your leave was approved.", "leave", refID); err != nil {
		t.Fatalf("Notify failed: %v", err)
	}

	list, total, err := svc.ListNotifications(ctx(), recipient, nil, 1, 10)
	if err != nil {
		t.Fatalf("ListNotifications failed: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected 1 notification, got %d", total)
	}
	if list[0].IsRead {
		t.Error("expected notification to be unread")
	}
	if *list[0].ReferenceType != "leave" || *list[0].ReferenceID != refID {
		t.Error("expected reference_type/reference_id to be persisted")
	}
}

func TestService_MarkAsRead_OnlyAffectsOwnNotification(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	owner := uuid.New()
	other := uuid.New()
	if err := svc.Notify(ctx(), owner, "GENERIC", "Title", "Body", "generic", uuid.New()); err != nil {
		t.Fatalf("Notify failed: %v", err)
	}
	list, _, _ := repo.ListNotificationsByRecipient(ctx(), owner, nil, 1, 10)
	notifID := list[0].ID

	if err := svc.MarkAsRead(ctx(), notifID, other); err == nil {
		t.Fatal("expected error marking another user's notification as read")
	}

	if err := svc.MarkAsRead(ctx(), notifID, owner); err != nil {
		t.Fatalf("MarkAsRead failed: %v", err)
	}
	updated, err := repo.FindNotificationByID(ctx(), notifID)
	if err != nil {
		t.Fatalf("FindNotificationByID failed: %v", err)
	}
	if !updated.IsRead || updated.ReadAt == nil {
		t.Error("expected notification to be marked read with read_at set")
	}
}

func TestService_MarkAllAsRead_And_GetUnreadCount(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	recipient := uuid.New()
	for i := 0; i < 3; i++ {
		if err := svc.Notify(ctx(), recipient, "GENERIC", "Title", "Body", "generic", uuid.New()); err != nil {
			t.Fatalf("Notify failed: %v", err)
		}
	}

	count, err := svc.GetUnreadCount(ctx(), recipient)
	if err != nil {
		t.Fatalf("GetUnreadCount failed: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected 3 unread, got %d", count)
	}

	if err := svc.MarkAllAsRead(ctx(), recipient); err != nil {
		t.Fatalf("MarkAllAsRead failed: %v", err)
	}

	count, err = svc.GetUnreadCount(ctx(), recipient)
	if err != nil {
		t.Fatalf("GetUnreadCount failed: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 unread after MarkAllAsRead, got %d", count)
	}
}
