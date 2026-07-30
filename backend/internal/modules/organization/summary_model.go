package organization

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrganizationSummary merepresentasikan ringkasan organisasi (SK pendirian).
// Satu summary dapat memiliki banyak Organization nodes.
type OrganizationSummary struct {
	ID         uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Code       string         `gorm:"type:varchar(7);not null;uniqueIndex:idx_orgsum_code" json:"code"`
	DecreeNo   string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_orgsum_decree" json:"decree_no"`
	DecreeDate string         `gorm:"type:date;not null" json:"decree_date"`
	Status     string         `gorm:"type:varchar(20);default:inactive" json:"status"`
	CreatedBy  *uuid.UUID     `gorm:"type:char(36)" json:"created_by,omitempty"`
	UpdatedBy  *uuid.UUID     `gorm:"type:char(36)" json:"updated_by,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`

	// Relasi hasMany
	Organizations []Organization `gorm:"foreignKey:OrganizationSummaryID" json:"organizations,omitempty"`
}

func (OrganizationSummary) TableName() string {
	return "organization_summaries"
}

func (s *OrganizationSummary) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// ── DTOs ──

type CreateOrganizationSummaryRequest struct {
	Code        string `json:"code" binding:"required,max=7"`
	DecreeNo    string `json:"decree_no" binding:"required,max=20"`
	DecreeDate  string `json:"decree_date" binding:"required"` // format: YYYY-MM-DD
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
	CloneFromID string `json:"clone_from_id" binding:"omitempty,uuid"`
}

type UpdateOrganizationSummaryRequest struct {
	Code       *string `json:"code" binding:"omitempty,max=7"`
	DecreeNo   *string `json:"decree_no" binding:"omitempty,max=20"`
	DecreeDate *string `json:"decree_date" binding:"omitempty"`
	Status     *string `json:"status" binding:"omitempty"`
}

type OrganizationSummaryResponse struct {
	ID           string `json:"id"`
	Code         string `json:"code"`
	DecreeNo     string `json:"decree_no"`
	DecreeDate   string `json:"decree_date"`
	Status       string `json:"status"`
	OrgCount     int    `json:"org_count"`
	ClonedFromID string `json:"cloned_from_id,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type PaginatedSummaryResponse struct {
	Success    bool                         `json:"success"`
	Data       []OrganizationSummaryResponse `json:"data"`
	Page       int                          `json:"page"`
	PerPage    int                          `json:"per_page"`
	Total      int64                        `json:"total"`
	TotalPages int                          `json:"total_pages"`
}

func (s *OrganizationSummary) ToResponse() OrganizationSummaryResponse {
	resp := OrganizationSummaryResponse{
		ID:         s.ID.String(),
		Code:       s.Code,
		DecreeNo:   s.DecreeNo,
		DecreeDate: s.DecreeDate,
		Status:     s.Status,
		CreatedAt:  s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  s.UpdatedAt.Format(time.RFC3339),
	}
	if len(s.Organizations) > 0 {
		resp.OrgCount = len(s.Organizations)
	}
	return resp
}
