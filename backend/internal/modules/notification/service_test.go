package notification

import (
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
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
	if err := svc.Notify(ctx(), recipient, "LEAVE_APPROVED", nil, "leave", refID); err != nil {
		t.Fatalf("Notify failed: %v", err)
	}

	list, total, err := svc.ListNotifications(ctx(), recipient, nil, 1, 10, httputil.LangEN)
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

// =========================================================================
// Bilingual title/body rendering
// =========================================================================

// TestService_ListNotifications_RendersTitleBodyPerRequestedLanguage guards
// the core bilingual requirement: the SAME stored notification renders
// different title/body text depending on the language passed to
// ListNotifications (the recipient's own Accept-Language at request time),
// not whatever language happened to be active when Notify() was called.
func TestService_ListNotifications_RendersTitleBodyPerRequestedLanguage(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	recipient := uuid.New()
	if err := svc.Notify(ctx(), recipient, "LEAVE_APPROVED", nil, "leave", uuid.New()); err != nil {
		t.Fatalf("Notify failed: %v", err)
	}

	en, _, err := svc.ListNotifications(ctx(), recipient, nil, 1, 10, httputil.LangEN)
	if err != nil {
		t.Fatalf("ListNotifications (en) failed: %v", err)
	}
	if len(en) != 1 || en[0].Title != "Leave Request Approved" {
		t.Fatalf("expected English title 'Leave Request Approved', got %+v", en)
	}

	id, _, err := svc.ListNotifications(ctx(), recipient, nil, 1, 10, httputil.LangID)
	if err != nil {
		t.Fatalf("ListNotifications (id) failed: %v", err)
	}
	if len(id) != 1 || id[0].Title != "Permohonan Cuti Disetujui" {
		t.Fatalf("expected Indonesian title 'Permohonan Cuti Disetujui', got %+v", id)
	}
	if id[0].Body != "Permohonan cuti Anda telah disetujui." {
		t.Errorf("expected Indonesian body, got %q", id[0].Body)
	}
}

// TestService_ListNotifications_RendersParamsIntoBodyTemplate guards
// parameterized notification types (e.g. approval task assignment, which
// embeds the module and step name) — params must be persisted and
// substituted correctly regardless of which language ends up being rendered.
func TestService_ListNotifications_RendersParamsIntoBodyTemplate(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	recipient := uuid.New()
	if err := svc.Notify(ctx(), recipient, "APPROVAL_TASK_ASSIGNED", []string{"leave", "Manager Approval"}, "leave", uuid.New()); err != nil {
		t.Fatalf("Notify failed: %v", err)
	}

	en, _, err := svc.ListNotifications(ctx(), recipient, nil, 1, 10, httputil.LangEN)
	if err != nil {
		t.Fatalf("ListNotifications (en) failed: %v", err)
	}
	wantEN := "A leave request needs your approval (Manager Approval)."
	if len(en) != 1 || en[0].Body != wantEN {
		t.Fatalf("expected body %q, got %+v", wantEN, en)
	}

	id, _, err := svc.ListNotifications(ctx(), recipient, nil, 1, 10, httputil.LangID)
	if err != nil {
		t.Fatalf("ListNotifications (id) failed: %v", err)
	}
	wantID := "Permohonan leave memerlukan persetujuan Anda (Manager Approval)."
	if len(id) != 1 || id[0].Body != wantID {
		t.Fatalf("expected body %q, got %+v", wantID, id)
	}
}

// TestService_ListNotifications_UnknownType_FallsBackToStoredDefault
// ensures a notifType with no catalog entry (e.g. a module that hasn't been
// added to i18n.go yet) degrades gracefully to the English text computed
// once at Notify()-time, instead of erroring or showing blank content.
func TestService_ListNotifications_UnknownType_FallsBackToStoredDefault(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	recipient := uuid.New()
	if err := svc.Notify(ctx(), recipient, "SOME_FUTURE_TYPE", nil, "generic", uuid.New()); err != nil {
		t.Fatalf("Notify failed: %v", err)
	}

	id, _, err := svc.ListNotifications(ctx(), recipient, nil, 1, 10, httputil.LangID)
	if err != nil {
		t.Fatalf("ListNotifications (id) failed: %v", err)
	}
	if len(id) != 1 || id[0].Title != "SOME_FUTURE_TYPE" {
		t.Fatalf("expected fallback title 'SOME_FUTURE_TYPE' (unknown type, no catalog entry), got %+v", id)
	}
}

func TestService_MarkAsRead_OnlyAffectsOwnNotification(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	owner := uuid.New()
	other := uuid.New()
	if err := svc.Notify(ctx(), owner, "GENERIC", nil, "generic", uuid.New()); err != nil {
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
		if err := svc.Notify(ctx(), recipient, "GENERIC", nil, "generic", uuid.New()); err != nil {
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
