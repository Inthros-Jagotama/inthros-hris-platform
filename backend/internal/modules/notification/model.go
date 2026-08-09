package notification

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification represents an in-app notification delivered to a platform user.
// Recipients are keyed by user_id (not employee_id), matching the recipient
// identity convention already used by the Approval module.
type Notification struct {
	ID              uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	RecipientUserID uuid.UUID `gorm:"type:char(36);not null;index:idx_notification_recipient" json:"recipient_user_id"`
	Type            string    `gorm:"type:varchar(50);not null" json:"type"`
	// Title/Body store the English rendering computed once at Notify()-time
	// (fallback + audit/debug legibility, e.g. reading the raw table).
	// ListNotifications overrides these per-response using Type+Params
	// translated into the requesting recipient's own language (see i18n.go)
	// — Title/Body here are NOT the source of truth for what a user sees.
	Title string `gorm:"type:varchar(255);not null" json:"title"`
	Body  string `gorm:"type:varchar(1000);not null" json:"body"`
	// Params is a JSON-encoded []string used to fill the body template's
	// %s placeholders when re-rendering in a different language (e.g. a
	// step name or module name) — nil/empty for notification types whose
	// body has no dynamic content.
	Params        *string    `gorm:"type:text" json:"-"`
	ReferenceType *string    `gorm:"type:varchar(50)" json:"reference_type,omitempty"`
	ReferenceID   *uuid.UUID `gorm:"type:char(36)" json:"reference_id,omitempty"`
	IsRead        bool       `gorm:"not null;default:false" json:"is_read"`
	ReadAt        *time.Time `json:"read_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (Notification) TableName() string {
	return "notifications"
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
