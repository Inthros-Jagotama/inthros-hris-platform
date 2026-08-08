package notification

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func ctx() context.Context {
	return context.Background()
}

func TestRepository_CreateNotification_And_FindByID(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	recipient := uuid.New()
	refID := uuid.New()
	refType := "leave_request"
	n := &Notification{
		RecipientUserID: recipient,
		Type:            "LEAVE_APPROVED",
		Title:           "Leave Approved",
		Body:            "Your leave request has been approved.",
		ReferenceType:   &refType,
		ReferenceID:     &refID,
	}

	if err := repo.CreateNotification(ctx(), n); err != nil {
		t.Fatalf("CreateNotification failed: %v", err)
	}
	if n.ID == uuid.Nil {
		t.Fatal("expected ID to be generated")
	}

	found, err := repo.FindNotificationByID(ctx(), n.ID)
	if err != nil {
		t.Fatalf("FindNotificationByID failed: %v", err)
	}
	if found.Title != "Leave Approved" {
		t.Errorf("expected title 'Leave Approved', got %q", found.Title)
	}
	if found.IsRead {
		t.Error("expected new notification to be unread")
	}
}

func TestRepository_ListNotificationsByRecipient_FiltersAndPaginates(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)

	recipient := uuid.New()
	other := uuid.New()

	for i := 0; i < 3; i++ {
		if err := repo.CreateNotification(ctx(), &Notification{
			RecipientUserID: recipient,
			Type:            "GENERIC",
			Title:           "Notif",
			Body:            "Body",
		}); err != nil {
			t.Fatalf("CreateNotification failed: %v", err)
		}
	}
	if err := repo.CreateNotification(ctx(), &Notification{
		RecipientUserID: other,
		Type:            "GENERIC",
		Title:           "Not mine",
		Body:            "Body",
	}); err != nil {
		t.Fatalf("CreateNotification failed: %v", err)
	}

	list, total, err := repo.ListNotificationsByRecipient(ctx(), recipient, nil, 1, 10)
	if err != nil {
		t.Fatalf("ListNotificationsByRecipient failed: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected 3 notifications for recipient, got %d", total)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 notifications returned, got %d", len(list))
	}

	unread := false
	list[0].IsRead = true
	if err := repo.UpdateNotification(ctx(), &list[0]); err != nil {
		t.Fatalf("UpdateNotification failed: %v", err)
	}

	unreadOnly, unreadTotal, err := repo.ListNotificationsByRecipient(ctx(), recipient, &unread, 1, 10)
	if err != nil {
		t.Fatalf("ListNotificationsByRecipient (unread filter) failed: %v", err)
	}
	if unreadTotal != 2 {
		t.Fatalf("expected 2 unread notifications, got %d", unreadTotal)
	}
	if len(unreadOnly) != 2 {
		t.Fatalf("expected 2 unread notifications returned, got %d", len(unreadOnly))
	}
}
