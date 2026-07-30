package organization

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 500
)

// PaginatedResponse DTO untuk response pagination.
type PaginatedResponse struct {
	Success    bool                   `json:"success"`
	Data       []OrganizationResponse `json:"data"`
	Page       int                    `json:"page"`
	PerPage    int                    `json:"per_page"`
	Total      int64                  `json:"total"`
	TotalPages int                    `json:"total_pages"`
}

type Service struct {
	repo   *Repository
	logger *zap.Logger
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// =========================================================================
// Organization CRUD — with automatic history capture
// =========================================================================

func (s *Service) Create(ctx context.Context, req CreateOrganizationRequest) (*OrganizationResponse, error) {
	// Parse required organization_summary_id
	summaryUUID, err := uuid.Parse(req.OrganizationSummaryID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_summary_id: %w", err)
	}

	// ── Field-level length validations (defense-in-depth) ──
	if len(req.Code) > 10 {
		return nil, fmt.Errorf("code must not exceed 10 characters")
	}
	if len(req.Nomenclature) > 255 {
		return nil, fmt.Errorf("nomenclature must not exceed 255 characters")
	}

	// Set current user from context
	currentUserID := authctx.GetUserID(ctx)

	org := &Organization{
		Code:                  req.Code,
		Nomenclature:          req.Nomenclature,
		OrganizationSummaryID: &summaryUUID,
		CreatedBy:             currentUserID,
		UpdatedBy:             currentUserID,
	}

	// Parse optional foreign keys with error handling
	if req.ParentID != nil && *req.ParentID != "" {
		id, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent_id: %w", err)
		}
		org.ParentID = &id
	}
	if req.ZoneID != nil && *req.ZoneID != "" {
		id, err := uuid.Parse(*req.ZoneID)
		if err != nil {
			return nil, fmt.Errorf("invalid zone_id: %w", err)
		}
		org.ZoneID = &id
	}
	if req.JobFamilyID != nil && *req.JobFamilyID != "" {
		id, err := uuid.Parse(*req.JobFamilyID)
		if err != nil {
			return nil, fmt.Errorf("invalid job_family_id: %w", err)
		}
		org.JobFamilyID = &id
	}
	if req.GradingID != nil && *req.GradingID != "" {
		id, err := uuid.Parse(*req.GradingID)
		if err != nil {
			return nil, fmt.Errorf("invalid grading_id: %w", err)
		}
		org.GradingID = &id
	}

	// Generate full_code and level based on parent
	if org.ParentID != nil {
		parent, err := s.repo.FindByID(ctx, *org.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent not found: %w", err)
		}
		org.FullCode = parent.FullCode + org.Code
		org.Level = parent.Level + 1
	} else {
		org.FullCode = org.Code
		org.Level = 0
	}

	// Validate auto-generated full_code length
	if len(org.FullCode) > 50 {
		return nil, fmt.Errorf("generated full_code '%s' exceeds maximum length of 50 characters", org.FullCode)
	}

	// Validate unique full_code within the same summary
	existing, err := s.repo.FindByFullCodeAndSummary(ctx, org.FullCode, *org.OrganizationSummaryID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("full_code '%s' already exists in this organization summary", org.FullCode)
	}

	if err := s.repo.Create(ctx, org); err != nil {
		return nil, err
	}

	// Capture history
	newValues := captureOrgValues(org)
	newJSON, _ := toJSONString(newValues)
	s.captureHistory(ctx, org.ID, ActionCreate, nil, &newJSON)

	s.logger.Info("Organization created",
		zap.String("id", org.ID.String()),
		zap.String("code", org.FullCode),
		zap.String("nomenclature", org.Nomenclature),
	)

	response := org.ToResponse()
	return &response, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*OrganizationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	org, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	response := org.ToResponse()
	return &response, nil
}

// List mengembalikan daftar organisasi dengan pagination, opsional filter summary_id dan active_only.
func (s *Service) List(ctx context.Context, page, perPage int, summaryID string, activeOnly bool) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}

	orgs, total, err := s.repo.FindAll(ctx, page, perPage, summaryID, activeOnly)
	if err != nil {
		return nil, err
	}

	var responses []OrganizationResponse
	for _, o := range orgs {
		responses = append(responses, o.ToResponse())
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) GetTree(ctx context.Context, summaryID string) ([]OrganizationResponse, error) {
	tree, err := s.repo.FindTree(ctx, summaryID)
	if err != nil {
		return nil, err
	}

	responses := make([]OrganizationResponse, 0, len(tree))
	for _, org := range tree {
		responses = append(responses, org.ToResponse())
	}
	return responses, nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateOrganizationRequest) (*OrganizationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	org, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	// ── Field-level length validations (defense-in-depth) ──
	if req.Code != nil && len(*req.Code) > 10 {
		return nil, fmt.Errorf("code must not exceed 10 characters")
	}
	if req.Nomenclature != nil && len(*req.Nomenclature) > 255 {
		return nil, fmt.Errorf("nomenclature must not exceed 255 characters")
	}

	// Set current user from context
	currentUserID := authctx.GetUserID(ctx)
	org.UpdatedBy = currentUserID

	// Capture old values before update
	oldValues := captureOrgValues(org)
	oldJSON, _ := toJSONString(oldValues)

	if req.Code != nil {
		org.Code = *req.Code
		// Recalculate full_code if code changed
		if org.ParentID != nil {
			parent, err := s.repo.FindByID(ctx, *org.ParentID)
			if err == nil {
				org.FullCode = parent.FullCode + org.Code
			}
		} else {
			org.FullCode = org.Code
		}
	}
	if req.Nomenclature != nil {
		org.Nomenclature = *req.Nomenclature
	}
	if req.ZoneID != nil {
		if *req.ZoneID == "" {
			org.ZoneID = nil
		} else {
			id, err := uuid.Parse(*req.ZoneID)
			if err != nil {
				return nil, fmt.Errorf("invalid zone_id: %w", err)
			}
			org.ZoneID = &id
		}
	}
	if req.JobFamilyID != nil {
		if *req.JobFamilyID == "" {
			org.JobFamilyID = nil
		} else {
			id, err := uuid.Parse(*req.JobFamilyID)
			if err != nil {
				return nil, fmt.Errorf("invalid job_family_id: %w", err)
			}
			org.JobFamilyID = &id
		}
	}
	if req.GradingID != nil {
		if *req.GradingID == "" {
			org.GradingID = nil
		} else {
			id, err := uuid.Parse(*req.GradingID)
			if err != nil {
				return nil, fmt.Errorf("invalid grading_id: %w", err)
			}
			org.GradingID = &id
		}
	}

	// Validate auto-generated full_code after update
	if len(org.FullCode) > 50 {
		return nil, fmt.Errorf("generated full_code '%s' exceeds maximum length of 50 characters", org.FullCode)
	}

	// Validate unique full_code within the same summary (exclude self)
	if org.OrganizationSummaryID != nil {
		existing, err := s.repo.FindByFullCodeAndSummaryExcludeSelf(ctx, org.FullCode, *org.OrganizationSummaryID, uid)
		if err == nil && existing != nil {
			return nil, fmt.Errorf("full_code '%s' already exists in this organization summary", org.FullCode)
		}
	}

	if err := s.repo.Update(ctx, org); err != nil {
		return nil, err
	}

	// Capture new values after update
	newValues := captureOrgValues(org)
	newJSON, _ := toJSONString(newValues)
	s.captureHistory(ctx, org.ID, ActionUpdate, &oldJSON, &newJSON)

	response := org.ToResponse()
	return &response, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	org, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return err
	}

	// Capture old values before delete
	oldValues := captureOrgValues(org)
	oldJSON, _ := toJSONString(oldValues)

	if err := s.repo.SoftDelete(ctx, uid); err != nil {
		return err
	}

	// Capture history for delete
	s.captureHistory(ctx, org.ID, ActionDelete, &oldJSON, nil)

	return nil
}

// =========================================================================
// History Service Methods
// =========================================================================

func (s *Service) GetHistory(ctx context.Context, orgID string, page, perPage int) (*PaginatedHistoryResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}

	var histories []OrganizationHistory
	var total int64

	if orgID != "" {
		uid, err := uuid.Parse(orgID)
		if err != nil {
			return nil, fmt.Errorf("invalid organization id: %w", err)
		}
		histories, total, err = s.repo.FindHistoryByOrgID(ctx, uid, page, perPage)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		histories, total, err = s.repo.FindAllHistory(ctx, page, perPage)
		if err != nil {
			return nil, err
		}
	}

	responses := make([]HistoryResponse, len(histories))
	for i, h := range histories {
		responses[i] = HistoryResponse{
			ID:             h.ID.String(),
			OrganizationID: h.OrganizationID.String(),
			Action:         string(h.Action),
			OldValues:      h.OldValues,
			NewValues:      h.NewValues,
			CreatedAt:      h.CreatedAt,
		}
		if h.ChangedBy != nil {
			cb := h.ChangedBy.String()
			responses[i].ChangedBy = &cb
		}
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedHistoryResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// =========================================================================
// Version Service Methods
// =========================================================================

func (s *Service) CreateVersion(ctx context.Context, req CreateVersionRequest) (*VersionResponse, error) {
	// Get current full tree
	tree, err := s.GetTree(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get organization tree: %w", err)
	}

	// Also get flat list for complete data
	allOrgs, err := s.repo.FindAllOrganizationsFlat(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}

	nodes := make([]organizationNode, len(allOrgs))
	for i, o := range allOrgs {
		nodes[i] = captureOrgValues(&o)
	}

	snapshot := versionData{
		Tree:  tree,
		Nodes: nodes,
	}

	snapshotJSON, err := json.Marshal(snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	version := &OrganizationVersion{
		VersionName: req.VersionName,
		Description: req.Description,
		Snapshot:    string(snapshotJSON),
		Status:      VersionStatusActive,
		NodeCount:   len(allOrgs),
	}

	if err := s.repo.CreateVersion(ctx, version); err != nil {
		return nil, err
	}

	s.logger.Info("Organization version created",
		zap.String("version_id", version.ID.String()),
		zap.String("name", version.VersionName),
		zap.Int("node_count", version.NodeCount),
	)

	resp := versionToResponse(version, false)
	return &resp, nil
}

func (s *Service) ListVersions(ctx context.Context, page, perPage int) (*PaginatedVersionResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}

	versions, total, err := s.repo.FindAllVersions(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	responses := make([]VersionResponse, len(versions))
	for i, v := range versions {
		responses[i] = versionToResponse(&v, false)
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedVersionResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) GetVersion(ctx context.Context, id string, includeSnapshot bool) (*VersionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid version id: %w", err)
	}

	version, err := s.repo.FindVersionByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	resp := versionToResponse(version, includeSnapshot)
	return &resp, nil
}

func (s *Service) DiffVersions(ctx context.Context, sourceID, targetID string) (*DiffResponse, error) {
	srcUID, err := uuid.Parse(sourceID)
	if err != nil {
		return nil, fmt.Errorf("invalid source version id: %w", err)
	}
	tgtUID, err := uuid.Parse(targetID)
	if err != nil {
		return nil, fmt.Errorf("invalid target version id: %w", err)
	}

	source, err := s.repo.FindVersionByID(ctx, srcUID)
	if err != nil {
		return nil, err
	}
	target, err := s.repo.FindVersionByID(ctx, tgtUID)
	if err != nil {
		return nil, err
	}

	// Parse snapshots
	var srcData, tgtData versionData
	if err := json.Unmarshal([]byte(source.Snapshot), &srcData); err != nil {
		return nil, fmt.Errorf("failed to parse source snapshot: %w", err)
	}
	if err := json.Unmarshal([]byte(target.Snapshot), &tgtData); err != nil {
		return nil, fmt.Errorf("failed to parse target snapshot: %w", err)
	}

	// Build maps for comparison
	srcMap := make(map[string]organizationNode)
	for _, n := range srcData.Nodes {
		srcMap[n.ID] = n
	}
	tgtMap := make(map[string]organizationNode)
	for _, n := range tgtData.Nodes {
		tgtMap[n.ID] = n
	}

	var changes []DiffEntry
	addedCount := 0
	removedCount := 0
	modifiedCount := 0

	// Find added and modified nodes
	for id, tgtNode := range tgtMap {
		if srcNode, exists := srcMap[id]; exists {
			// Compare fields
			diffs := compareNodes(srcNode, tgtNode)
			for _, d := range diffs {
				changes = append(changes, DiffEntry{
					Type:     d.Type,
					OrgID:    id,
					Code:     tgtNode.Code,
					Name:     tgtNode.Nomenclature,
					Field:    d.Field,
					OldValue: d.OldValue,
					NewValue: d.NewValue,
				})
				modifiedCount++
			}
		} else {
			changes = append(changes, DiffEntry{
				Type:  "ADDED",
				OrgID: id,
				Code:  tgtNode.Code,
				Name:  tgtNode.Nomenclature,
			})
			addedCount++
		}
	}

	// Find removed nodes
	for id, srcNode := range srcMap {
		if _, exists := tgtMap[id]; !exists {
			changes = append(changes, DiffEntry{
				Type:  "REMOVED",
				OrgID: id,
				Code:  srcNode.Code,
				Name:  srcNode.Nomenclature,
			})
			removedCount++
		}
	}

	// Sort changes by type then code
	sort.Slice(changes, func(i, j int) bool {
		if changes[i].Type != changes[j].Type {
			return changes[i].Type < changes[j].Type
		}
		return changes[i].Code < changes[j].Code
	})

	return &DiffResponse{
		SourceVersionID: sourceID,
		TargetVersionID: targetID,
		SourceName:      source.VersionName,
		TargetName:      target.VersionName,
		Changes:         changes,
		AddedCount:      addedCount,
		RemovedCount:    removedCount,
		ModifiedCount:   modifiedCount,
	}, nil
}

func (s *Service) RestoreVersion(ctx context.Context, versionID string) (*VersionResponse, error) {
	uid, err := uuid.Parse(versionID)
	if err != nil {
		return nil, fmt.Errorf("invalid version id: %w", err)
	}

	version, err := s.repo.FindVersionByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Parse snapshot
	var data versionData
	if err := json.Unmarshal([]byte(version.Snapshot), &data); err != nil {
		return nil, fmt.Errorf("failed to parse snapshot: %w", err)
	}

	// Recreate organizations from snapshot nodes, preserving original UUIDs
	// so that parent-child relationships (ParentID references) remain valid.
	// Uses atomic RestoreAllFromSnapshot (delete + create in one transaction).
	newOrgs := make([]Organization, 0, len(data.Nodes))
	for _, n := range data.Nodes {
		// ── Field-level length validations (defense-in-depth) ──
		if len(n.Code) > 10 {
			return nil, fmt.Errorf("snapshot node '%s': code must not exceed 10 characters", n.Code)
		}
		if len(n.FullCode) > 50 {
			return nil, fmt.Errorf("snapshot node '%s': full_code must not exceed 50 characters", n.Code)
		}
		if len(n.Nomenclature) > 255 {
			return nil, fmt.Errorf("snapshot node '%s': nomenclature must not exceed 255 characters", n.Nomenclature)
		}

		orgID, err := uuid.Parse(n.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID in snapshot for node %s: %w", n.Code, err)
		}
		newOrg := Organization{
			ID:           orgID, // Preserve original UUID for FK integrity
			Code:         n.Code,
			FullCode:     n.FullCode,
			Nomenclature: n.Nomenclature,
			Level:        n.Level,
			SortOrder:    n.SortOrder,
		}
		// Parse UUID references (these reference other restored orgs, so must use snapshot UUIDs)
		if n.ParentID != "" {
			pid, err := uuid.Parse(n.ParentID)
			if err == nil {
				newOrg.ParentID = &pid
			}
		}
		if n.ZoneID != "" {
			zid, err := uuid.Parse(n.ZoneID)
			if err == nil {
				newOrg.ZoneID = &zid
			}
		}
		if n.JobFamilyID != "" {
			jid, err := uuid.Parse(n.JobFamilyID)
			if err == nil {
				newOrg.JobFamilyID = &jid
			}
		}
		if n.GradingID != "" {
			gid, err := uuid.Parse(n.GradingID)
			if err == nil {
				newOrg.GradingID = &gid
			}
		}
		newOrgs = append(newOrgs, newOrg)
	}

	// Atomic restore: delete all + create all in single transaction
	if err := s.repo.RestoreAllFromSnapshot(ctx, newOrgs); err != nil {
		return nil, fmt.Errorf("failed to restore organizations: %w", err)
	}

	// Mark version as used for restore (archive it)
	version.Status = VersionStatusArchived
	_ = s.repo.UpdateVersion(ctx, version)

	s.logger.Info("Organization structure restored from version",
		zap.String("version_id", version.ID.String()),
		zap.String("version_name", version.VersionName),
		zap.Int("nodes_restored", len(newOrgs)),
	)

	resp := versionToResponse(version, false)
	return &resp, nil
}

func (s *Service) CloneVersion(ctx context.Context, req CloneRequest) (*CloneResponse, error) {
	// Get current tree
	tree, err := s.GetTree(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get organization tree: %w", err)
	}

	// Get flat list
	allOrgs, err := s.repo.FindAllOrganizationsFlat(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}

	nodes := make([]organizationNode, len(allOrgs))
	for i, o := range allOrgs {
		nodes[i] = captureOrgValues(&o)
	}

	snapshot := versionData{
		Tree:  tree,
		Nodes: nodes,
	}

	snapshotJSON, err := json.Marshal(snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	version := &OrganizationVersion{
		VersionName: req.VersionName + " (Clone)",
		Description: req.Description,
		Snapshot:    string(snapshotJSON),
		Status:      VersionStatusDraft,
		NodeCount:   len(allOrgs),
	}

	if err := s.repo.CreateVersion(ctx, version); err != nil {
		return nil, err
	}

	s.logger.Info("Organization tree cloned to draft version",
		zap.String("version_id", version.ID.String()),
		zap.String("name", version.VersionName),
		zap.Int("nodes_cloned", version.NodeCount),
	)

	resp := versionToResponse(version, false)
	return &CloneResponse{
		Version:     resp,
		NodesCloned: version.NodeCount,
	}, nil
}

// =========================================================================
// Internal Helpers
// =========================================================================

type organizationNode struct {
	ID                    string `json:"id"`
	Code                  string `json:"code"`
	FullCode              string `json:"full_code"`
	Nomenclature          string `json:"nomenclature"`
	ParentID              string `json:"parent_id"`
	ZoneID                string `json:"zone_id"`
	JobFamilyID           string `json:"job_family_id"`
	GradingID             string `json:"grading_id"`
	Level                 int    `json:"level"`
	SortOrder             int    `json:"sort_order"`
	OrganizationSummaryID string `json:"organization_summary_id"`
}

type versionData struct {
	Tree  []OrganizationResponse `json:"tree"`
	Nodes []organizationNode     `json:"nodes"`
}

type nodeDiff struct {
	Type     string `json:"type"`
	Field    string `json:"field"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}

func captureOrgValues(o *Organization) organizationNode {
	node := organizationNode{
		ID:           o.ID.String(),
		Code:         o.Code,
		FullCode:     o.FullCode,
		Nomenclature: o.Nomenclature,
		Level:        o.Level,
		SortOrder:    o.SortOrder,
	}
	if o.ParentID != nil {
		node.ParentID = o.ParentID.String()
	}
	if o.ZoneID != nil {
		node.ZoneID = o.ZoneID.String()
	}
	if o.JobFamilyID != nil {
		node.JobFamilyID = o.JobFamilyID.String()
	}
	if o.GradingID != nil {
		node.GradingID = o.GradingID.String()
	}
	if o.OrganizationSummaryID != nil {
		node.OrganizationSummaryID = o.OrganizationSummaryID.String()
	}
	return node
}

func compareNodes(old, new organizationNode) []nodeDiff {
	var diffs []nodeDiff
	if old.Code != new.Code {
		diffs = append(diffs, nodeDiff{Type: "MODIFIED", Field: "code", OldValue: old.Code, NewValue: new.Code})
	}
	if old.Nomenclature != new.Nomenclature {
		diffs = append(diffs, nodeDiff{Type: "MODIFIED", Field: "nomenclature", OldValue: old.Nomenclature, NewValue: new.Nomenclature})
	}
	if old.ParentID != new.ParentID {
		diffs = append(diffs, nodeDiff{Type: "MODIFIED", Field: "parent_id", OldValue: old.ParentID, NewValue: new.ParentID})
	}
	if old.ZoneID != new.ZoneID {
		diffs = append(diffs, nodeDiff{Type: "MODIFIED", Field: "zone_id", OldValue: old.ZoneID, NewValue: new.ZoneID})
	}
	if old.JobFamilyID != new.JobFamilyID {
		diffs = append(diffs, nodeDiff{Type: "MODIFIED", Field: "job_family_id", OldValue: old.JobFamilyID, NewValue: new.JobFamilyID})
	}
	if old.SortOrder != new.SortOrder {
		diffs = append(diffs, nodeDiff{Type: "MODIFIED", Field: "sort_order", OldValue: fmt.Sprintf("%d", old.SortOrder), NewValue: fmt.Sprintf("%d", new.SortOrder)})
	}
	return diffs
}

func toJSONString(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *Service) captureHistory(ctx context.Context, orgID uuid.UUID, action ActionType, oldValues, newValues *string) {
	history := &OrganizationHistory{
		OrganizationID: orgID,
		Action:         action,
		OldValues:      oldValues,
		NewValues:      newValues,
	}
	if err := s.repo.CreateHistory(ctx, history); err != nil {
		s.logger.Error("Failed to capture organization history",
			zap.String("org_id", orgID.String()),
			zap.String("action", string(action)),
			zap.Error(err),
		)
	}
}

func versionToResponse(v *OrganizationVersion, includeSnapshot bool) VersionResponse {
	resp := VersionResponse{
		ID:          v.ID.String(),
		VersionName: v.VersionName,
		Description: v.Description,
		Status:      string(v.Status),
		NodeCount:   v.NodeCount,
		CreatedAt:   v.CreatedAt,
		UpdatedAt:   v.UpdatedAt,
	}
	if v.CreatedBy != nil {
		cb := v.CreatedBy.String()
		resp.CreatedBy = &cb
	}
	if v.ParentVersionID != nil {
		pvid := v.ParentVersionID.String()
		resp.ParentVersionID = &pvid
	}
	if includeSnapshot {
		resp.Snapshot = &v.Snapshot
	}
	return resp
}
