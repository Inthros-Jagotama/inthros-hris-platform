package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
)

func (s *Service) CreateSummary(ctx context.Context, req CreateOrganizationSummaryRequest) (*OrganizationSummaryResponse, error) {
	currentUserID := authctx.GetUserID(ctx)

	summary := &OrganizationSummary{
		Code:       req.Code,
		DecreeNo:   req.DecreeNo,
		DecreeDate: req.DecreeDate,
		Status:     "active",
		CreatedBy:  currentUserID,
		UpdatedBy:  currentUserID,
	}

	if err := s.repo.CreateSummary(ctx, summary); err != nil {
		return nil, err
	}

	s.logger.Info("Organization summary created",
		zap.String("id", summary.ID.String()),
		zap.String("code", summary.Code),
		zap.String("decree_no", summary.DecreeNo),
	)

	resp := summary.ToResponse()
	return &resp, nil
}

func (s *Service) GetSummaryByID(ctx context.Context, id string) (*OrganizationSummaryResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	summary, err := s.repo.FindSummaryByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	resp := summary.ToResponse()
	return &resp, nil
}

func (s *Service) ListSummaries(ctx context.Context, page, perPage int) (*PaginatedSummaryResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}

	summaries, total, err := s.repo.FindAllSummaries(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	responses := make([]OrganizationSummaryResponse, len(summaries))
	for i, sum := range summaries {
		resp := sum.ToResponse()
		// Count organizations associated with this summary
		if len(sum.Organizations) > 0 {
			resp.OrgCount = len(sum.Organizations)
		}
		responses[i] = resp
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedSummaryResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) UpdateSummary(ctx context.Context, id string, req UpdateOrganizationSummaryRequest) (*OrganizationSummaryResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	summary, err := s.repo.FindSummaryByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Set current user from context
	currentUserID := authctx.GetUserID(ctx)
	summary.UpdatedBy = currentUserID

	if req.Code != nil {
		summary.Code = *req.Code
	}
	if req.DecreeNo != nil {
		summary.DecreeNo = *req.DecreeNo
	}
	if req.DecreeDate != nil {
		summary.DecreeDate = *req.DecreeDate
	}
	if req.Status != nil {
		summary.Status = *req.Status
	}

	if err := s.repo.UpdateSummary(ctx, summary); err != nil {
		return nil, err
	}

	s.logger.Info("Organization summary updated",
		zap.String("id", summary.ID.String()),
		zap.String("code", summary.Code),
	)

	resp := summary.ToResponse()
	return &resp, nil
}

func (s *Service) DeleteSummary(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	// Check if summary has organizations
	summary, err := s.repo.FindSummaryByID(ctx, uid)
	if err != nil {
		return err
	}

	if len(summary.Organizations) > 0 {
		return fmt.Errorf("cannot delete summary with %d organization(s) attached", len(summary.Organizations))
	}

	if err := s.repo.SoftDeleteSummary(ctx, uid); err != nil {
		return err
	}

	s.logger.Info("Organization summary deleted",
		zap.String("id", uid.String()),
		zap.String("code", summary.Code),
	)

	return nil
}

func (s *Service) GetSummaryStats(ctx context.Context) (map[string]interface{}, error) {
	_, _, err := s.repo.FindAllSummaries(ctx, 1, 1)
	if err != nil {
		return nil, err
	}

	// Get total count
	db, err := s.repo.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var totalSummaries int64
	db.Model(&OrganizationSummary{}).Count(&totalSummaries)

	// Get org count (all orgs)
	var totalOrgs int64
	db.Model(&Organization{}).Count(&totalOrgs)

	// Get max depth
	var maxDepth int
	type depthResult struct {
		Level int
	}
	var deepest depthResult
	db.Model(&Organization{}).Select("MAX(level) as level").Scan(&deepest)
	maxDepth = deepest.Level

	stats := map[string]interface{}{
		"total_summaries": totalSummaries,
		"total_orgs":      totalOrgs,
		"max_depth":       maxDepth,
		"updated_at":      time.Now().Format(time.RFC3339),
	}

	return stats, nil
}
