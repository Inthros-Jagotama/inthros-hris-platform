package organization

import "time"

type CreateOrganizationRequest struct {
	OrganizationSummaryID string  `json:"organization_summary_id" binding:"required,uuid"`
	Code                  string  `json:"code" binding:"required,max=10"`
	Nomenclature          string  `json:"nomenclature" binding:"required,max=255"`
	ParentID              *string `json:"parent_id"`
	ZoneID                *string `json:"zone_id"`
	JobFamilyID           *string `json:"job_family_id"`
	GradingID             *string `json:"grading_id"`
}

type UpdateOrganizationRequest struct {
	Code         *string `json:"code" binding:"omitempty,max=10"`
	Nomenclature *string `json:"nomenclature" binding:"omitempty,max=255"`
	ParentID     *string `json:"parent_id"`
	ZoneID       *string `json:"zone_id"`
	JobFamilyID  *string `json:"job_family_id"`
	GradingID    *string `json:"grading_id"`
}

type OrganizationResponse struct {
	ID                   string     `json:"id"`
	Code                 string     `json:"code"`
	FullCode             string     `json:"full_code"`
	Nomenclature         string     `json:"nomenclature"`
	ParentID             *string    `json:"parent_id,omitempty"`
	ZoneID               *string    `json:"zone_id,omitempty"`
	JobFamilyID          *string    `json:"job_family_id,omitempty"`
	GradingID            *string    `json:"grading_id,omitempty"`
	Level                int        `json:"level"`
	SortOrder            int        `json:"sort_order"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	Children             []OrganizationResponse `json:"children,omitempty"`
}

func (o *Organization) ToResponse() OrganizationResponse {
	resp := OrganizationResponse{
		ID:          o.ID.String(),
		Code:        o.Code,
		FullCode:    o.FullCode,
		Nomenclature: o.Nomenclature,
		Level:       o.Level,
		SortOrder:   o.SortOrder,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}

	if o.ParentID != nil {
		pid := o.ParentID.String()
		resp.ParentID = &pid
	}
	if o.ZoneID != nil {
		zid := o.ZoneID.String()
		resp.ZoneID = &zid
	}
	if o.JobFamilyID != nil {
		jid := o.JobFamilyID.String()
		resp.JobFamilyID = &jid
	}
	if o.GradingID != nil {
		gid := o.GradingID.String()
		resp.GradingID = &gid
	}

	// Convert children recursively
	if len(o.Children) > 0 {
		resp.Children = make([]OrganizationResponse, 0, len(o.Children))
		for _, child := range o.Children {
			resp.Children = append(resp.Children, child.ToResponse())
		}
	}

	return resp
}

// =========================================================================
// Organization History DTOs
// =========================================================================

type HistoryResponse struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Action         string    `json:"action"`
	OldValues      *string   `json:"old_values,omitempty"`
	NewValues      *string   `json:"new_values,omitempty"`
	ChangedBy      *string   `json:"changed_by,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type PaginatedHistoryResponse struct {
	Success    bool              `json:"success"`
	Data       []HistoryResponse `json:"data"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	Total      int64             `json:"total"`
	TotalPages int               `json:"total_pages"`
}

// =========================================================================
// Organization Version DTOs
// =========================================================================

type CreateVersionRequest struct {
	VersionName string  `json:"version_name" binding:"required,max=100"`
	Description *string `json:"description"`
}

type VersionResponse struct {
	ID              string    `json:"id"`
	VersionName     string    `json:"version_name"`
	Description     *string   `json:"description,omitempty"`
	Snapshot        *string   `json:"snapshot,omitempty"`
	Status          string    `json:"status"`
	CreatedBy       *string   `json:"created_by,omitempty"`
	ParentVersionID *string   `json:"parent_version_id,omitempty"`
	NodeCount       int       `json:"node_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type PaginatedVersionResponse struct {
	Success    bool              `json:"success"`
	Data       []VersionResponse `json:"data"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	Total      int64             `json:"total"`
	TotalPages int               `json:"total_pages"`
}

// =========================================================================
// Clone & Diff DTOs
// =========================================================================

type CloneRequest struct {
	VersionName string  `json:"version_name" binding:"required,max=100"`
	Description *string `json:"description"`
}

type CloneResponse struct {
	Version     VersionResponse   `json:"version"`
	NodesCloned int               `json:"nodes_cloned"`
}

type DiffEntry struct {
	Type     string `json:"type"`     // ADDED, REMOVED, MODIFIED
	OrgID    string `json:"org_id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Field    string `json:"field,omitempty"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value,omitempty"`
}

type DiffResponse struct {
	SourceVersionID string      `json:"source_version_id"`
	TargetVersionID string      `json:"target_version_id"`
	SourceName      string      `json:"source_name"`
	TargetName      string      `json:"target_name"`
	Changes         []DiffEntry `json:"changes"`
	AddedCount      int         `json:"added_count"`
	RemovedCount    int         `json:"removed_count"`
	ModifiedCount   int         `json:"modified_count"`
}

type RestoreRequest struct {
	Confirm bool `json:"confirm"`
}

type ApplyVersionRequest struct {
	Confirm bool `json:"confirm"`
}
