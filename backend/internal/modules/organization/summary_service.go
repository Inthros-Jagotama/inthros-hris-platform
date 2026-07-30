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

	// ── Validate only one active ──
	status := req.Status
	if status == "" {
		status = "inactive"
	}
	if status == "active" {
		activeCount, err := s.repo.CountActiveSummaries(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to check active summaries: %w", err)
		}
		if activeCount > 0 {
			return nil, fmt.Errorf("only one organization summary can be active at a time. Deactivate the current active summary first")
		}
	}

	summary := &OrganizationSummary{
		Code:       req.Code,
		DecreeNo:   req.DecreeNo,
		DecreeDate: req.DecreeDate,
		Status:     status,
		CreatedBy:  currentUserID,
		UpdatedBy:  currentUserID,
	}

	if err := s.repo.CreateSummary(ctx, summary); err != nil {
		return nil, err
	}

	var clonedFromID string

	// ── Clone organizations from source summary ──
	if req.CloneFromID != "" {
		sourceUUID, err := uuid.Parse(req.CloneFromID)
		if err != nil {
			return nil, fmt.Errorf("invalid clone_from_id: %w", err)
		}

		// Get all organizations from source summary
		sourceOrgs, err := s.repo.FindAllBySummaryID(ctx, sourceUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch source organizations: %w", err)
		}

		if len(sourceOrgs) > 0 {
			// Build UUID mapping: old UUID → new UUID
			uuidMap := make(map[string]string) // old uuid string → new uuid string
			for _, org := range sourceOrgs {
				uuidMap[org.ID.String()] = uuid.New().String()
			}

			// Create new organizations with mapped UUIDs and parent references
			newOrgs := make([]Organization, 0, len(sourceOrgs))
			for _, src := range sourceOrgs {
				newID, _ := uuid.Parse(uuidMap[src.ID.String()])
				newSummaryID := summary.ID

				newOrg := Organization{
					ID:                   newID,
					OrganizationSummaryID: &newSummaryID,
					Code:                 src.Code,
					FullCode:             src.FullCode,
					Nomenclature:         src.Nomenclature,
					Level:                src.Level,
					SortOrder:            src.SortOrder,
				}

				// Map ZoneID (referensi ke setting — tetap sama)
				if src.ZoneID != nil {
					newOrg.ZoneID = src.ZoneID
				}
				// Map JobFamilyID (referensi ke setting — tetap sama)
				if src.JobFamilyID != nil {
					newOrg.JobFamilyID = src.JobFamilyID
				}
				// Map GradingID (referensi ke setting — tetap sama)
				if src.GradingID != nil {
					newOrg.GradingID = src.GradingID
				}
				// Map ParentID (referensi ke org yang sudah di-clone)
				if src.ParentID != nil {
					if newParentIDStr, ok := uuidMap[src.ParentID.String()]; ok {
						newParentUUID, _ := uuid.Parse(newParentIDStr)
						newOrg.ParentID = &newParentUUID
					}
				}
				// Map CreatedBy/UpdatedBy
				newOrg.CreatedBy = currentUserID
				newOrg.UpdatedBy = currentUserID

				newOrgs = append(newOrgs, newOrg)
			}

			if err := s.repo.BulkCreateOrganizations(ctx, newOrgs); err != nil {
				return nil, fmt.Errorf("failed to clone organizations: %w", err)
			}

			clonedFromID = req.CloneFromID

			s.logger.Info("Organizations cloned from summary",
				zap.String("source_summary_id", req.CloneFromID),
				zap.String("target_summary_id", summary.ID.String()),
				zap.Int("orgs_cloned", len(newOrgs)),
			)
		}
	}

	s.logger.Info("Organization summary created",
		zap.String("id", summary.ID.String()),
		zap.String("code", summary.Code),
		zap.String("decree_no", summary.DecreeNo),
	)

	resp := summary.ToResponse()
	resp.ClonedFromID = clonedFromID
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
		newStatus := *req.Status
		// ── Validate only one active ──
		if newStatus == "active" && summary.Status != "active" {
			activeCount, err := s.repo.CountActiveSummaries(ctx, &uid)
			if err != nil {
				return nil, fmt.Errorf("failed to check active summaries: %w", err)
			}
			if activeCount > 0 {
				return nil, fmt.Errorf("only one organization summary can be active at a time. Deactivate the current active summary first")
			}
		}
		summary.Status = newStatus
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
