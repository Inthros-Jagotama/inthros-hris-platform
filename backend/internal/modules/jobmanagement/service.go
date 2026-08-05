package jobmanagement

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

type Service struct {
	repo   *Repository
	logger *zap.Logger
	calc   *Calculator
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger, calc: NewCalculator(repo)}
}

// =========================================================================
// Job Titles (9.1)
// =========================================================================

func (s *Service) CreateJobTitle(ctx context.Context, req CreateJobTitleRequest) (*JobTitleResponse, error) {
	t := &JobTitle{
		Name: &req.Name,
	}
	if req.Descriptions != "" {
		t.Descriptions = &req.Descriptions
	}
	if req.Status != 0 {
		t.Status = &req.Status
	}
	t.CreatedBy = authctx.GetUserID(ctx)
	t.UpdatedBy = t.CreatedBy
	if err := s.repo.CreateJobTitle(ctx, t); err != nil {
		return nil, err
	}
	s.logger.Info("Job title created", zap.String("id", t.ID.String()), zap.String("name", req.Name))
	r := toJobTitleResponse(t)
	return &r, nil
}

func (s *Service) GetJobTitleByID(ctx context.Context, id string) (*JobTitleResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	t, err := s.repo.FindJobTitleByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := toJobTitleResponse(t)
	return &r, nil
}

func (s *Service) ListJobTitles(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	titles, total, err := s.repo.FindAllJobTitles(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]JobTitleResponse, 0, len(titles))
	for _, t := range titles {
		responses = append(responses, toJobTitleResponse(&t))
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

func (s *Service) UpdateJobTitle(ctx context.Context, id string, req UpdateJobTitleRequest) (*JobTitleResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	t, err := s.repo.FindJobTitleByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	t.UpdatedBy = authctx.GetUserID(ctx)
	if req.Name != nil {
		t.Name = req.Name
	}
	if req.Descriptions != nil {
		t.Descriptions = req.Descriptions
	}
	if req.Status != nil {
		t.Status = req.Status
	}
	if err := s.repo.UpdateJobTitle(ctx, t); err != nil {
		return nil, err
	}
	r := toJobTitleResponse(t)
	return &r, nil
}

func (s *Service) DeleteJobTitle(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteJobTitle(ctx, uid)
}

// =========================================================================
// Job Title Subs (9.2)
// =========================================================================

func (s *Service) CreateJobTitleSub(ctx context.Context, titleID string, req CreateJobTitleSubRequest) (*JobTitleSubResponse, error) {
	titleUID, err := uuid.Parse(titleID)
	if err != nil {
		return nil, fmt.Errorf("invalid title id: %w", err)
	}
	sub := &JobTitleSub{
		JobManagementTitleID: &titleUID,
		Name:                 &req.Name,
	}
	if req.Descriptions != "" {
		sub.Descriptions = &req.Descriptions
	}
	if req.Status != 0 {
		sub.Status = &req.Status
	}
	sub.CreatedBy = authctx.GetUserID(ctx)
	sub.UpdatedBy = sub.CreatedBy
	// Copy title name if available
	title, err := s.repo.FindJobTitleByID(ctx, titleUID)
	if err == nil && title.Name != nil {
		sub.JobManagementTitleName = title.Name
	}
	if err := s.repo.CreateJobTitleSub(ctx, sub); err != nil {
		return nil, err
	}
	r := toJobTitleSubResponse(sub)
	return &r, nil
}

func (s *Service) GetJobTitleSubByID(ctx context.Context, id string) (*JobTitleSubResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	sub, err := s.repo.FindJobTitleSubByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := toJobTitleSubResponse(sub)
	return &r, nil
}

func (s *Service) ListJobTitleSubs(ctx context.Context, titleID string) ([]JobTitleSubResponse, error) {
	titleUID, err := uuid.Parse(titleID)
	if err != nil {
		return nil, fmt.Errorf("invalid title id: %w", err)
	}
	subs, err := s.repo.FindJobTitleSubsByTitleID(ctx, titleUID)
	if err != nil {
		return nil, err
	}
	responses := make([]JobTitleSubResponse, 0, len(subs))
	for _, sub := range subs {
		responses = append(responses, toJobTitleSubResponse(&sub))
	}
	return responses, nil
}

func (s *Service) UpdateJobTitleSub(ctx context.Context, id string, req UpdateJobTitleSubRequest) (*JobTitleSubResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	sub, err := s.repo.FindJobTitleSubByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	sub.UpdatedBy = authctx.GetUserID(ctx)
	if req.Name != nil {
		sub.Name = req.Name
	}
	if req.Descriptions != nil {
		sub.Descriptions = req.Descriptions
	}
	if req.Status != nil {
		sub.Status = req.Status
	}
	if err := s.repo.UpdateJobTitleSub(ctx, sub); err != nil {
		return nil, err
	}
	r := toJobTitleSubResponse(sub)
	return &r, nil
}

func (s *Service) DeleteJobTitleSub(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteJobTitleSub(ctx, uid)
}

// =========================================================================
// Job Values (9.3)
// =========================================================================

func (s *Service) CreateJobValue(ctx context.Context, req CreateJobValueRequest) (*JobValueResponse, error) {
	v := &JobValue{
		Type: req.Type,
	}
	if req.JobManagementTitleSubID != nil && *req.JobManagementTitleSubID != "" {
		subID, err := uuid.Parse(*req.JobManagementTitleSubID)
		if err != nil {
			return nil, fmt.Errorf("invalid job_management_title_sub_id: %w", err)
		}
		v.JobManagementTitleSubID = &subID
		// Copy sub name
		sub, err := s.repo.FindJobTitleSubByID(ctx, subID)
		if err == nil && sub.Name != nil {
			v.JobManagementTitleSubName = sub.Name
		}
	}
	if req.TypeGroup != nil && *req.TypeGroup != "" {
		v.TypeGroup = req.TypeGroup
	}
	if req.Level != nil {
		v.Level = req.Level
	}
	if req.Descriptions != "" {
		v.Descriptions = &req.Descriptions
	}
	if req.DescriptionGroup != nil && *req.DescriptionGroup != "" {
		v.DescriptionGroup = req.DescriptionGroup
	}
	if req.Note != "" {
		v.Note = &req.Note
	}
	if req.Sort != nil {
		v.Sort = req.Sort
	}
	if req.RefID != nil && *req.RefID != "" {
		refID, err := uuid.Parse(*req.RefID)
		if err != nil {
			return nil, fmt.Errorf("invalid ref_id: %w", err)
		}
		v.RefID = &refID
	}
	if req.RefType != nil && *req.RefType != "" {
		v.RefType = req.RefType
	}
	v.CreatedBy = authctx.GetUserID(ctx)
	v.UpdatedBy = v.CreatedBy
	if err := s.repo.CreateJobValue(ctx, v); err != nil {
		return nil, err
	}
	r := toJobValueResponse(v)
	return &r, nil
}

func (s *Service) GetJobValueByID(ctx context.Context, id string) (*JobValueResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	v, err := s.repo.FindJobValueByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := toJobValueResponse(v)
	return &r, nil
}

// ListJobValuesTree mengembalikan tree job management values:
// type_group → daftar tipe (label = description_group) → options per tipe
// (level + deskripsi). Filter type_group dilakukan di sisi FE.
func (s *Service) ListJobValuesTree(ctx context.Context) (*JobValueTreeResponse, error) {
	values, err := s.repo.FindAllJobValuesForTree(ctx)
	if err != nil {
		return nil, err
	}

	// Urutan tampilan grup tetap (konsisten dengan seed legacy)
	groupOrder := []string{
		"education", "experience",
		"psychological", "technical", "managerial",
		"communication", "problem_solving",
		"financial", "asset",
		"subordinate", "activity", "environment", "risk", "relationship", "frequency",
	}
	groupRank := map[string]int{}
	for i, g := range groupOrder {
		groupRank[g] = i
	}

	groupIdx := map[string]int{}         // type_group -> index di data
	typeIdx := map[string]map[string]int{} // type_group -> type -> index dalam group
	var data []JobValueTreeGroup

	for i := range values {
		v := &values[i]
		if v.TypeGroup == nil || *v.TypeGroup == "" {
			continue // value tanpa type_group tidak masuk tree
		}
		tg := *v.TypeGroup
		gi, ok := groupIdx[tg]
		if !ok {
			gi = len(data)
			groupIdx[tg] = gi
			data = append(data, JobValueTreeGroup{TypeGroup: tg, Types: []JobValueTreeType{}})
			typeIdx[tg] = map[string]int{}
		}
		ti, ok := typeIdx[tg][v.Type]
		if !ok {
			ti = len(data[gi].Types)
			typeIdx[tg][v.Type] = ti
			desc := ""
			if v.DescriptionGroup != nil {
				desc = *v.DescriptionGroup
			}
			data[gi].Types = append(data[gi].Types, JobValueTreeType{
				Type:             v.Type,
				DescriptionGroup: desc,
				Options:          []JobValueTreeOption{},
			})
		}
		opt := JobValueTreeOption{ID: v.ID.String()}
		if v.Level != nil {
			opt.Level = *v.Level
		}
		if v.Descriptions != nil {
			opt.Descriptions = *v.Descriptions
		}
		data[gi].Types[ti].Options = append(data[gi].Types[ti].Options, opt)
	}

	// Sortir grup sesuai urutan tetap (group tidak dikenal di akhir, alfabetis)
	sort.Slice(data, func(a, b int) bool {
		ra, oka := groupRank[data[a].TypeGroup]
		rb, okb := groupRank[data[b].TypeGroup]
		if !oka {
			ra = 999
		}
		if !okb {
			rb = 999
		}
		if ra != rb {
			return ra < rb
		}
		return data[a].TypeGroup < data[b].TypeGroup
	})
	// Sortir tipe dalam group (label description_group lalu type) dan isi label grup
	for gi := range data {
		if len(data[gi].Types) > 0 {
			data[gi].DescriptionGroup = data[gi].Types[0].DescriptionGroup
		}
		sort.Slice(data[gi].Types, func(a, b int) bool {
			if data[gi].Types[a].DescriptionGroup != data[gi].Types[b].DescriptionGroup {
				return data[gi].Types[a].DescriptionGroup < data[gi].Types[b].DescriptionGroup
			}
			return data[gi].Types[a].Type < data[gi].Types[b].Type
		})
	}

	return &JobValueTreeResponse{Success: true, Data: data}, nil
}

func (s *Service) ListJobValues(ctx context.Context, page, perPage int, valueType string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	values, total, err := s.repo.FindAllJobValues(ctx, page, perPage, valueType)
	if err != nil {
		return nil, err
	}
	responses := make([]JobValueResponse, 0, len(values))
	for _, v := range values {
		responses = append(responses, toJobValueResponse(&v))
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

func (s *Service) UpdateJobValue(ctx context.Context, id string, req UpdateJobValueRequest) (*JobValueResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	v, err := s.repo.FindJobValueByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	v.UpdatedBy = authctx.GetUserID(ctx)
	if req.Type != nil {
		v.Type = *req.Type
	}
	if req.TypeGroup != nil {
		v.TypeGroup = req.TypeGroup
	}
	if req.Level != nil {
		v.Level = req.Level
	}
	if req.Descriptions != nil {
		v.Descriptions = req.Descriptions
	}
	if req.DescriptionGroup != nil {
		v.DescriptionGroup = req.DescriptionGroup
	}
	if req.Note != nil {
		v.Note = req.Note
	}
	if req.Sort != nil {
		v.Sort = req.Sort
	}
	if req.RefID != nil {
		if *req.RefID == "" {
			v.RefID = nil
		} else {
			refID, err := uuid.Parse(*req.RefID)
			if err != nil {
				return nil, fmt.Errorf("invalid ref_id: %w", err)
			}
			v.RefID = &refID
		}
	}
	if req.RefType != nil {
		v.RefType = req.RefType
	}
	if err := s.repo.UpdateJobValue(ctx, v); err != nil {
		return nil, err
	}
	r := toJobValueResponse(v)
	return &r, nil
}

func (s *Service) DeleteJobValue(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteJobValue(ctx, uid)
}

// =========================================================================
// Job Value Clusters (9.3b) — mapping type ↔ cluster kompetensi
// =========================================================================

// ListJobValueClusters mengembalikan daftar cluster yang dipetakan ke sebuah tipe
// nilai jabatan (mis. 'technical'), diurutkan alfabetis.
func (s *Service) ListJobValueClusters(ctx context.Context, valueType string) (*JobValueClusterResponse, error) {
	clusters, err := s.repo.FindJobValueClusters(ctx, valueType)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(clusters))
	for _, c := range clusters {
		names = append(names, c.Cluster)
	}
	return &JobValueClusterResponse{Type: valueType, Clusters: names}, nil
}

// UpdateJobValueClusters mengganti seluruh mapping cluster untuk satu tipe dalam
// satu transaksi. Array kosong = hapus semua mapping (aksi valid, bukan error).
func (s *Service) UpdateJobValueClusters(ctx context.Context, valueType string, req UpdateJobValueClustersRequest) (*JobValueClusterResponse, error) {
	// Trim + dedupe; abaikan nilai kosong
	seen := make(map[string]struct{}, len(req.Clusters))
	clean := make([]string, 0, len(req.Clusters))
	for _, c := range req.Clusters {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		if _, dup := seen[c]; dup {
			continue
		}
		seen[c] = struct{}{}
		clean = append(clean, c)
	}
	if err := s.repo.ReplaceJobValueClusters(ctx, valueType, clean); err != nil {
		return nil, err
	}
	return &JobValueClusterResponse{Type: valueType, Clusters: clean}, nil
}

// =========================================================================
// Management entities with shared CRUD pattern (9.4 - 9.15)
// These all follow the same pattern: Create, GetByID, List, Update, Delete
// with nomenclature + full_code + organization_id fields
// =========================================================================

// =========================================================================
// Job Objectives (9.4)
// =========================================================================

func (s *Service) CreateJobObjective(ctx context.Context, req CreateJobObjectiveRequest) (*JobObjectiveResponse, error) {
	o := &JobObjective{
		Nomenclature: req.Nomenclature,
		FullCode:     req.FullCode,
	}
	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}
	o.OrganizationID = &orgID
	if req.Objective != "" {
		o.Objective = &req.Objective
	}
	o.CreatedBy = authctx.GetUserID(ctx)
	o.UpdatedBy = o.CreatedBy
	if err := s.repo.CreateJobObjective(ctx, o); err != nil {
		return nil, err
	}
	r := toJobObjectiveResponse(o)
	return &r, nil
}

func (s *Service) GetJobObjectiveByID(ctx context.Context, id string) (*JobObjectiveResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	o, err := s.repo.FindJobObjectiveByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := toJobObjectiveResponse(o)
	return &r, nil
}

func (s *Service) ListJobObjectives(ctx context.Context, page, perPage int, orgID string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var orgUID *uuid.UUID
	if orgID != "" {
		if id, err := uuid.Parse(orgID); err == nil {
			orgUID = &id
		}
	}
	objectives, total, err := s.repo.FindAllJobObjectives(ctx, page, perPage, orgUID)
	if err != nil {
		return nil, err
	}
	responses := make([]JobObjectiveResponse, 0, len(objectives))
	for _, o := range objectives {
		responses = append(responses, toJobObjectiveResponse(&o))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateJobObjective(ctx context.Context, id string, req UpdateJobObjectiveRequest) (*JobObjectiveResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	o, err := s.repo.FindJobObjectiveByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	o.UpdatedBy = authctx.GetUserID(ctx)
	if req.Nomenclature != nil {
		o.Nomenclature = *req.Nomenclature
	}
	if req.FullCode != nil {
		o.FullCode = *req.FullCode
	}
	if req.Objective != nil {
		o.Objective = req.Objective
	}
	if err := s.repo.UpdateJobObjective(ctx, o); err != nil {
		return nil, err
	}
	r := toJobObjectiveResponse(o)
	return &r, nil
}

func (s *Service) DeleteJobObjective(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteJobObjective(ctx, uid)
}

// =========================================================================
// Job Identifications (9.5)
// =========================================================================

func (s *Service) CreateJobIdentification(ctx context.Context, req CreateJobIdentificationRequest) (*JobIdentificationResponse, error) {
	i := &JobIdentification{
		Nomenclature: req.Nomenclature,
		FullCode:     req.FullCode,
	}
	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}
	i.OrganizationID = &orgID
	gradingID, err := uuid.Parse(req.GradingID)
	if err != nil {
		return nil, fmt.Errorf("invalid grading_id: %w", err)
	}
	i.GradingID = gradingID
	i.CreatedBy = authctx.GetUserID(ctx)
	i.UpdatedBy = i.CreatedBy
	if err := s.repo.CreateJobIdentification(ctx, i); err != nil {
		return nil, err
	}
	r := toJobIdentificationResponse(i)
	return &r, nil
}

func (s *Service) GetJobIdentificationByID(ctx context.Context, id string) (*JobIdentificationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	i, err := s.repo.FindJobIdentificationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := toJobIdentificationResponse(i)
	return &r, nil
}

func (s *Service) ListJobIdentifications(ctx context.Context, page, perPage int, orgID string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var orgUID *uuid.UUID
	if orgID != "" {
		if id, err := uuid.Parse(orgID); err == nil {
			orgUID = &id
		}
	}
	ids, total, err := s.repo.FindAllJobIdentifications(ctx, page, perPage, orgUID)
	if err != nil {
		return nil, err
	}
	responses := make([]JobIdentificationResponse, 0, len(ids))
	for _, i := range ids {
		responses = append(responses, toJobIdentificationResponse(&i))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateJobIdentification(ctx context.Context, id string, req UpdateJobIdentificationRequest) (*JobIdentificationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	i, err := s.repo.FindJobIdentificationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	i.UpdatedBy = authctx.GetUserID(ctx)
	if req.Nomenclature != nil {
		i.Nomenclature = *req.Nomenclature
	}
	if req.FullCode != nil {
		i.FullCode = *req.FullCode
	}
	if req.GradingID != nil && *req.GradingID != "" {
		gid, err := uuid.Parse(*req.GradingID)
		if err != nil {
			return nil, fmt.Errorf("invalid grading_id: %w", err)
		}
		i.GradingID = gid
	}
	if err := s.repo.UpdateJobIdentification(ctx, i); err != nil {
		return nil, err
	}
	r := toJobIdentificationResponse(i)
	return &r, nil
}

func (s *Service) DeleteJobIdentification(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteJobIdentification(ctx, uid)
}

// =========================================================================
// Job Responsibilities (9.6)
// =========================================================================

func (s *Service) CreateJobResponsibility(ctx context.Context, req CreateJobResponsibilityRequest) (*JobResponsibilityResponse, error) {
	r := &JobResponsibility{
		Nomenclature: req.Nomenclature,
		FullCode:     req.FullCode,
	}
	orgID, _ := uuid.Parse(req.OrganizationID)
	r.OrganizationID = &orgID
	if req.MainTask != "" {
		r.MainTask = &req.MainTask
	}
	if req.Activities != "" {
		r.Activities = &req.Activities
	}
	if req.Outputs != "" {
		r.Outputs = &req.Outputs
	}
	if req.SuccessIndicators != "" {
		r.SuccessIndicators = &req.SuccessIndicators
	}
	r.CreatedBy = authctx.GetUserID(ctx)
	r.UpdatedBy = r.CreatedBy
	if err := s.repo.CreateJobResponsibility(ctx, r); err != nil {
		return nil, err
	}
	resp := toJobResponsibilityResponse(r)
	return &resp, nil
}

func (s *Service) GetJobResponsibilityByID(ctx context.Context, id string) (*JobResponsibilityResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	r, err := s.repo.FindJobResponsibilityByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	resp := toJobResponsibilityResponse(r)
	return &resp, nil
}

func (s *Service) ListJobResponsibilities(ctx context.Context, page, perPage int, orgID string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var orgUID *uuid.UUID
	if orgID != "" {
		if id, err := uuid.Parse(orgID); err == nil {
			orgUID = &id
		}
	}
	responsibilities, total, err := s.repo.FindAllJobResponsibilities(ctx, page, perPage, orgUID)
	if err != nil {
		return nil, err
	}
	responses := make([]JobResponsibilityResponse, 0, len(responsibilities))
	for _, r := range responsibilities {
		responses = append(responses, toJobResponsibilityResponse(&r))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateJobResponsibility(ctx context.Context, id string, req UpdateJobResponsibilityRequest) (*JobResponsibilityResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	r, err := s.repo.FindJobResponsibilityByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r.UpdatedBy = authctx.GetUserID(ctx)
	if req.Nomenclature != nil {
		r.Nomenclature = *req.Nomenclature
	}
	if req.FullCode != nil {
		r.FullCode = *req.FullCode
	}
	if req.MainTask != nil {
		r.MainTask = req.MainTask
	}
	if req.Activities != nil {
		r.Activities = req.Activities
	}
	if req.Outputs != nil {
		r.Outputs = req.Outputs
	}
	if req.SuccessIndicators != nil {
		r.SuccessIndicators = req.SuccessIndicators
	}
	if err := s.repo.UpdateJobResponsibility(ctx, r); err != nil {
		return nil, err
	}
	resp := toJobResponsibilityResponse(r)
	return &resp, nil
}

func (s *Service) DeleteJobResponsibility(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteJobResponsibility(ctx, uid)
}

// =========================================================================
// Job Education Experiences (9.7)
// =========================================================================

func (s *Service) CreateJobEducationExperience(ctx context.Context, req CreateJobEducationExperienceRequest) (*JobEducationExperienceResponse, error) {
	e := &JobEducationExperience{
		Nomenclature: req.Nomenclature,
		FullCode:     req.FullCode,
	}
	orgID, _ := uuid.Parse(req.OrganizationID)
	e.OrganizationID = &orgID
	if req.EducationID != nil && *req.EducationID != "" {
		id, _ := uuid.Parse(*req.EducationID)
		e.EducationID = &id
	}
	if req.ExperienceID != nil && *req.ExperienceID != "" {
		id, _ := uuid.Parse(*req.ExperienceID)
		e.ExperienceID = &id
	}
	// Jurusan (multiple) → pivot job_management_majors
	for _, majorID := range req.EducationMajorID {
		if majorID == "" {
			continue
		}
		id, err := uuid.Parse(majorID)
		if err != nil {
			return nil, fmt.Errorf("invalid education_major_id: %w", err)
		}
		e.Majors = append(e.Majors, JobManagementMajor{EducationMajorID: id})
	}
	// Bidang Pekerjaan (multiple) → pivot job_management_job_family
	for _, familyID := range req.JobFamilyID {
		if familyID == "" {
			continue
		}
		id, err := uuid.Parse(familyID)
		if err != nil {
			return nil, fmt.Errorf("invalid job_family_id: %w", err)
		}
		e.JobFamilies = append(e.JobFamilies, JobManagementJobFamily{JobFamilyID: id})
	}
	e.CreatedBy = authctx.GetUserID(ctx)
	e.UpdatedBy = e.CreatedBy
	if err := s.repo.CreateJobEducationExperience(ctx, e); err != nil {
		return nil, err
	}
	s.recalculateScore(ctx, e.OrganizationID)
	// Re-fetch with relations loaded so response includes master names
	// (Education, Experience, Majors, JobFamilies) — create does not Preload them.
	if fetched, err := s.repo.FindJobEducationExperienceByID(ctx, e.ID); err == nil {
		e = fetched
	}
	r := toJobEducationExperienceResponse(e)
	return &r, nil
}

func (s *Service) GetJobEducationExperienceByID(ctx context.Context, id string) (*JobEducationExperienceResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	e, err := s.repo.FindJobEducationExperienceByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := toJobEducationExperienceResponse(e)
	return &r, nil
}

func (s *Service) ListJobEducationExperiences(ctx context.Context, page, perPage int, orgID string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var orgUID *uuid.UUID
	if orgID != "" {
		if id, err := uuid.Parse(orgID); err == nil {
			orgUID = &id
		}
	}
	experiences, total, err := s.repo.FindAllJobEducationExperiences(ctx, page, perPage, orgUID)
	if err != nil {
		return nil, err
	}
	responses := make([]JobEducationExperienceResponse, 0, len(experiences))
	for _, e := range experiences {
		responses = append(responses, toJobEducationExperienceResponse(&e))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateJobEducationExperience(ctx context.Context, id string, req UpdateJobEducationExperienceRequest) (*JobEducationExperienceResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	e, err := s.repo.FindJobEducationExperienceByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	e.UpdatedBy = authctx.GetUserID(ctx)
	if req.Nomenclature != nil {
		e.Nomenclature = *req.Nomenclature
	}
	if req.FullCode != nil {
		e.FullCode = *req.FullCode
	}
	if req.EducationID != nil {
		if *req.EducationID == "" {
			e.EducationID = nil
		} else if id, err := uuid.Parse(*req.EducationID); err == nil {
			e.EducationID = &id
		}
	}
	if req.ExperienceID != nil {
		if *req.ExperienceID == "" {
			e.ExperienceID = nil
		} else if id, err := uuid.Parse(*req.ExperienceID); err == nil {
			e.ExperienceID = &id
		}
	}
	// Jurusan (multiple) — array kosong = hapus semua; nil = biarkan
	if req.EducationMajorID != nil {
		e.Majors = nil
		for _, majorID := range req.EducationMajorID {
			if majorID == "" {
				continue
			}
			id, err := uuid.Parse(majorID)
			if err != nil {
				return nil, fmt.Errorf("invalid education_major_id: %w", err)
			}
			e.Majors = append(e.Majors, JobManagementMajor{EducationMajorID: id})
		}
	}
	// Bidang Pekerjaan (multiple) — array kosong = hapus semua; nil = biarkan
	if req.JobFamilyID != nil {
		e.JobFamilies = nil
		for _, familyID := range req.JobFamilyID {
			if familyID == "" {
				continue
			}
			id, err := uuid.Parse(familyID)
			if err != nil {
				return nil, fmt.Errorf("invalid job_family_id: %w", err)
			}
			e.JobFamilies = append(e.JobFamilies, JobManagementJobFamily{JobFamilyID: id})
		}
	}
	if err := s.repo.UpdateJobEducationExperience(ctx, e); err != nil {
		return nil, err
	}
	s.recalculateScore(ctx, e.OrganizationID)
	if fetched, err := s.repo.FindJobEducationExperienceByID(ctx, e.ID); err == nil {
		e = fetched
	}
	r := toJobEducationExperienceResponse(e)
	return &r, nil
}

func (s *Service) DeleteJobEducationExperience(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if rec, err := s.repo.FindJobEducationExperienceByID(ctx, uid); err == nil {
		defer s.recalculateScore(ctx, rec.OrganizationID)
	}
	return s.repo.DeleteJobEducationExperience(ctx, uid)
}

// =========================================================================
// Job HR Authorities (9.8) — simplified pattern for remaining entities
// =========================================================================

func (s *Service) CreateJobHRAuthority(ctx context.Context, req CreateJobHRAuthorityRequest) (*JobHRAuthorityResponse, error) {
	a := &JobHRAuthority{
		Nomenclature: req.Nomenclature,
		FullCode:     req.FullCode,
	}
	orgID, _ := uuid.Parse(req.OrganizationID)
	a.OrganizationID = &orgID
	if req.Description != "" {
		a.Description = &req.Description
	}
	a.CreatedBy = authctx.GetUserID(ctx)
	a.UpdatedBy = a.CreatedBy
	if err := s.repo.CreateJobHRAuthority(ctx, a); err != nil {
		return nil, err
	}
	r := toJobHRAuthorityResponse(a)
	return &r, nil
}

func (s *Service) GetJobHRAuthorityByID(ctx context.Context, id string) (*JobHRAuthorityResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindJobHRAuthorityByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := toJobHRAuthorityResponse(a)
	return &r, nil
}

func (s *Service) ListJobHRAuthorities(ctx context.Context, page, perPage int, orgID string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var orgUID *uuid.UUID
	if orgID != "" {
		if id, err := uuid.Parse(orgID); err == nil {
			orgUID = &id
		}
	}
	authorities, total, err := s.repo.FindAllJobHRAuthorities(ctx, page, perPage, orgUID)
	if err != nil {
		return nil, err
	}
	responses := make([]JobHRAuthorityResponse, 0, len(authorities))
	for _, a := range authorities {
		responses = append(responses, toJobHRAuthorityResponse(&a))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateJobHRAuthority(ctx context.Context, id string, req UpdateJobHRAuthorityRequest) (*JobHRAuthorityResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindJobHRAuthorityByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	a.UpdatedBy = authctx.GetUserID(ctx)
	if req.Nomenclature != nil {
		a.Nomenclature = *req.Nomenclature
	}
	if req.FullCode != nil {
		a.FullCode = *req.FullCode
	}
	if req.Description != nil {
		a.Description = req.Description
	}
	if err := s.repo.UpdateJobHRAuthority(ctx, a); err != nil {
		return nil, err
	}
	r := toJobHRAuthorityResponse(a)
	return &r, nil
}

func (s *Service) DeleteJobHRAuthority(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteJobHRAuthority(ctx, uid)
}

// =========================================================================
// Job Operational Authorities (9.9)
// =========================================================================

func (s *Service) CreateJobOperationalAuthority(ctx context.Context, req CreateJobOperationalAuthorityRequest) (*JobOperationalAuthorityResponse, error) {
	a := &JobOperationalAuthority{
		Nomenclature: req.Nomenclature,
		FullCode:     req.FullCode,
	}
	orgID, _ := uuid.Parse(req.OrganizationID)
	a.OrganizationID = &orgID
	if req.Description != "" {
		a.Description = &req.Description
	}
	a.CreatedBy = authctx.GetUserID(ctx)
	a.UpdatedBy = a.CreatedBy
	if err := s.repo.CreateJobOperationalAuthority(ctx, a); err != nil {
		return nil, err
	}
	r := toJobOperationalAuthorityResponse(a)
	return &r, nil
}

func (s *Service) GetJobOperationalAuthorityByID(ctx context.Context, id string) (*JobOperationalAuthorityResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindJobOperationalAuthorityByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := toJobOperationalAuthorityResponse(a)
	return &r, nil
}

func (s *Service) ListJobOperationalAuthorities(ctx context.Context, page, perPage int, orgID string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var orgUID *uuid.UUID
	if orgID != "" {
		if id, err := uuid.Parse(orgID); err == nil {
			orgUID = &id
		}
	}
	authorities, total, err := s.repo.FindAllJobOperationalAuthorities(ctx, page, perPage, orgUID)
	if err != nil {
		return nil, err
	}
	responses := make([]JobOperationalAuthorityResponse, 0, len(authorities))
	for _, a := range authorities {
		responses = append(responses, toJobOperationalAuthorityResponse(&a))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateJobOperationalAuthority(ctx context.Context, id string, req UpdateJobOperationalAuthorityRequest) (*JobOperationalAuthorityResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindJobOperationalAuthorityByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	a.UpdatedBy = authctx.GetUserID(ctx)
	if req.Nomenclature != nil {
		a.Nomenclature = *req.Nomenclature
	}
	if req.FullCode != nil {
		a.FullCode = *req.FullCode
	}
	if req.Description != nil {
		a.Description = req.Description
	}
	if err := s.repo.UpdateJobOperationalAuthority(ctx, a); err != nil {
		return nil, err
	}
	r := toJobOperationalAuthorityResponse(a)
	return &r, nil
}

func (s *Service) DeleteJobOperationalAuthority(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteJobOperationalAuthority(ctx, uid)
}

// =========================================================================
// Job Working Activities (9.10)
// =========================================================================

func (s *Service) CreateJobWorkingActivity(ctx context.Context, req CreateJobWorkingActivityRequest) (*JobWorkingActivityResponse, error) {
	a := &JobWorkingActivity{
		Nomenclature: req.Nomenclature,
		FullCode:     req.FullCode,
	}
	orgID, _ := uuid.Parse(req.OrganizationID)
	a.OrganizationID = &orgID
	if req.JobManagementValueID != nil && *req.JobManagementValueID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueID)
		a.JobManagementValueID = &id
	}
	a.CreatedBy = authctx.GetUserID(ctx)
	a.UpdatedBy = a.CreatedBy
	if err := s.repo.CreateJobWorkingActivity(ctx, a); err != nil {
		return nil, err
	}
	s.recalculateScore(ctx, a.OrganizationID)
	r := toJobWorkingActivityResponse(a)
	return &r, nil
}

func (s *Service) GetJobWorkingActivityByID(ctx context.Context, id string) (*JobWorkingActivityResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindJobWorkingActivityByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := toJobWorkingActivityResponse(a)
	return &r, nil
}

func (s *Service) ListJobWorkingActivities(ctx context.Context, page, perPage int, orgID string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var orgUID *uuid.UUID
	if orgID != "" {
		if id, err := uuid.Parse(orgID); err == nil {
			orgUID = &id
		}
	}
	activities, total, err := s.repo.FindAllJobWorkingActivities(ctx, page, perPage, orgUID)
	if err != nil {
		return nil, err
	}
	responses := make([]JobWorkingActivityResponse, 0, len(activities))
	for _, a := range activities {
		responses = append(responses, toJobWorkingActivityResponse(&a))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateJobWorkingActivity(ctx context.Context, id string, req UpdateJobWorkingActivityRequest) (*JobWorkingActivityResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindJobWorkingActivityByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	a.UpdatedBy = authctx.GetUserID(ctx)
	if req.Nomenclature != nil {
		a.Nomenclature = *req.Nomenclature
	}
	if req.FullCode != nil {
		a.FullCode = *req.FullCode
	}
	if req.JobManagementValueID != nil && *req.JobManagementValueID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueID)
		a.JobManagementValueID = &id
	}
	if err := s.repo.UpdateJobWorkingActivity(ctx, a); err != nil {
		return nil, err
	}
	s.recalculateScore(ctx, a.OrganizationID)
	r := toJobWorkingActivityResponse(a)
	return &r, nil
}

func (s *Service) DeleteJobWorkingActivity(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if rec, err := s.repo.FindJobWorkingActivityByID(ctx, uid); err == nil {
		defer s.recalculateScore(ctx, rec.OrganizationID)
	}
	return s.repo.DeleteJobWorkingActivity(ctx, uid)
}

// =========================================================================
// Job Working Risks (9.11)
// =========================================================================

func (s *Service) CreateJobWorkingRisk(ctx context.Context, req CreateJobWorkingRiskRequest) (*JobWorkingRiskResponse, error) {
	r := &JobWorkingRisk{
		Nomenclature: req.Nomenclature,
		FullCode:     req.FullCode,
	}
	orgID, _ := uuid.Parse(req.OrganizationID)
	r.OrganizationID = &orgID
	if req.JobManagementValueEnvironmentID != nil && *req.JobManagementValueEnvironmentID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueEnvironmentID)
		r.JobManagementValueEnvironmentID = &id
	}
	if req.JobManagementValueHazardID != nil && *req.JobManagementValueHazardID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueHazardID)
		r.JobManagementValueHazardID = &id
	}
	r.CreatedBy = authctx.GetUserID(ctx)
	r.UpdatedBy = r.CreatedBy
	if err := s.repo.CreateJobWorkingRisk(ctx, r); err != nil {
		return nil, err
	}
	s.recalculateScore(ctx, r.OrganizationID)
	resp := toJobWorkingRiskResponse(r)
	return &resp, nil
}

func (s *Service) GetJobWorkingRiskByID(ctx context.Context, id string) (*JobWorkingRiskResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	r, err := s.repo.FindJobWorkingRiskByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	resp := toJobWorkingRiskResponse(r)
	return &resp, nil
}

func (s *Service) ListJobWorkingRisks(ctx context.Context, page, perPage int, orgID string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var orgUID *uuid.UUID
	if orgID != "" {
		if id, err := uuid.Parse(orgID); err == nil {
			orgUID = &id
		}
	}
	risks, total, err := s.repo.FindAllJobWorkingRisks(ctx, page, perPage, orgUID)
	if err != nil {
		return nil, err
	}
	responses := make([]JobWorkingRiskResponse, 0, len(risks))
	for _, r := range risks {
		responses = append(responses, toJobWorkingRiskResponse(&r))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateJobWorkingRisk(ctx context.Context, id string, req UpdateJobWorkingRiskRequest) (*JobWorkingRiskResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	r, err := s.repo.FindJobWorkingRiskByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r.UpdatedBy = authctx.GetUserID(ctx)
	if req.Nomenclature != nil {
		r.Nomenclature = *req.Nomenclature
	}
	if req.FullCode != nil {
		r.FullCode = *req.FullCode
	}
	if req.JobManagementValueEnvironmentID != nil && *req.JobManagementValueEnvironmentID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueEnvironmentID)
		r.JobManagementValueEnvironmentID = &id
	}
	if req.JobManagementValueHazardID != nil && *req.JobManagementValueHazardID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueHazardID)
		r.JobManagementValueHazardID = &id
	}
	if err := s.repo.UpdateJobWorkingRisk(ctx, r); err != nil {
		return nil, err
	}
	s.recalculateScore(ctx, r.OrganizationID)
	resp := toJobWorkingRiskResponse(r)
	return &resp, nil
}

func (s *Service) DeleteJobWorkingRisk(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if rec, err := s.repo.FindJobWorkingRiskByID(ctx, uid); err == nil {
		defer s.recalculateScore(ctx, rec.OrganizationID)
	}
	return s.repo.DeleteJobWorkingRisk(ctx, uid)
}

// =========================================================================
// Job Relationships (9.12)
// =========================================================================

func (s *Service) CreateJobRelationship(ctx context.Context, req CreateJobRelationshipRequest) (*JobRelationshipResponse, error) {
	r := &JobRelationship{
		Nomenclature: req.Nomenclature,
		FullCode:     req.FullCode,
	}
	orgID, _ := uuid.Parse(req.OrganizationID)
	r.OrganizationID = &orgID
	if req.JobManagementValueRelationshipID != nil && *req.JobManagementValueRelationshipID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueRelationshipID)
		r.JobManagementValueRelationshipID = &id
	}
	if req.JobManagementValueFrequencyID != nil && *req.JobManagementValueFrequencyID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueFrequencyID)
		r.JobManagementValueFrequencyID = &id
	}
	r.CreatedBy = authctx.GetUserID(ctx)
	r.UpdatedBy = r.CreatedBy
	if err := s.repo.CreateJobRelationship(ctx, r); err != nil {
		return nil, err
	}
	s.recalculateScore(ctx, r.OrganizationID)
	resp := toJobRelationshipResponse(r)
	return &resp, nil
}

func (s *Service) GetJobRelationshipByID(ctx context.Context, id string) (*JobRelationshipResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	r, err := s.repo.FindJobRelationshipByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	resp := toJobRelationshipResponse(r)
	return &resp, nil
}

func (s *Service) ListJobRelationships(ctx context.Context, page, perPage int, orgID string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var orgUID *uuid.UUID
	if orgID != "" {
		if id, err := uuid.Parse(orgID); err == nil {
			orgUID = &id
		}
	}
	relationships, total, err := s.repo.FindAllJobRelationships(ctx, page, perPage, orgUID)
	if err != nil {
		return nil, err
	}
	responses := make([]JobRelationshipResponse, 0, len(relationships))
	for _, r := range relationships {
		responses = append(responses, toJobRelationshipResponse(&r))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateJobRelationship(ctx context.Context, id string, req UpdateJobRelationshipRequest) (*JobRelationshipResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	r, err := s.repo.FindJobRelationshipByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r.UpdatedBy = authctx.GetUserID(ctx)
	if req.Nomenclature != nil {
		r.Nomenclature = *req.Nomenclature
	}
	if req.FullCode != nil {
		r.FullCode = *req.FullCode
	}
	if req.JobManagementValueRelationshipID != nil && *req.JobManagementValueRelationshipID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueRelationshipID)
		r.JobManagementValueRelationshipID = &id
	}
	if req.JobManagementValueFrequencyID != nil && *req.JobManagementValueFrequencyID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueFrequencyID)
		r.JobManagementValueFrequencyID = &id
	}
	if err := s.repo.UpdateJobRelationship(ctx, r); err != nil {
		return nil, err
	}
	s.recalculateScore(ctx, r.OrganizationID)
	resp := toJobRelationshipResponse(r)
	return &resp, nil
}

func (s *Service) DeleteJobRelationship(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if rec, err := s.repo.FindJobRelationshipByID(ctx, uid); err == nil {
		defer s.recalculateScore(ctx, rec.OrganizationID)
	}
	return s.repo.DeleteJobRelationship(ctx, uid)
}

// =========================================================================
// Job Relationship Details (9.12b)
// =========================================================================

func (s *Service) CreateJobRelationshipDetail(ctx context.Context, relationshipID string, req CreateJobRelationshipDetailRequest) (*JobRelationshipDetailResponse, error) {
	relUID, err := uuid.Parse(relationshipID)
	if err != nil {
		return nil, fmt.Errorf("invalid relationship id: %w", err)
	}
	d := &JobManagementRelationshipDetail{
		JobManagementRelationshipID: relUID,
	}
	if req.OrganizationID != nil && *req.OrganizationID != "" {
		id, err := uuid.Parse(*req.OrganizationID)
		if err != nil {
			return nil, fmt.Errorf("invalid organization_id: %w", err)
		}
		d.OrganizationID = &id
	}
	if req.Activity != nil && *req.Activity != "" {
		d.Activity = req.Activity
	}
	d.CreatedBy = authctx.GetUserID(ctx)
	d.UpdatedBy = d.CreatedBy
	if err := s.repo.CreateJobRelationshipDetail(ctx, d); err != nil {
		return nil, err
	}
	if fetched, err := s.repo.FindJobRelationshipDetailByID(ctx, d.ID); err == nil {
		d = fetched
	}
	r := toJobRelationshipDetailResponse(d)
	return &r, nil
}

func (s *Service) GetJobRelationshipDetailByID(ctx context.Context, id string) (*JobRelationshipDetailResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	d, err := s.repo.FindJobRelationshipDetailByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := toJobRelationshipDetailResponse(d)
	return &r, nil
}

func (s *Service) ListJobRelationshipDetails(ctx context.Context, relationshipID string) ([]JobRelationshipDetailResponse, error) {
	relUID, err := uuid.Parse(relationshipID)
	if err != nil {
		return nil, fmt.Errorf("invalid relationship id: %w", err)
	}
	details, err := s.repo.FindAllJobRelationshipDetails(ctx, relUID)
	if err != nil {
		return nil, err
	}
	responses := make([]JobRelationshipDetailResponse, 0, len(details))
	for _, d := range details {
		responses = append(responses, toJobRelationshipDetailResponse(&d))
	}
	return responses, nil
}

func (s *Service) UpdateJobRelationshipDetail(ctx context.Context, id string, req UpdateJobRelationshipDetailRequest) (*JobRelationshipDetailResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	d, err := s.repo.FindJobRelationshipDetailByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	d.UpdatedBy = authctx.GetUserID(ctx)
	if req.OrganizationID != nil {
		if *req.OrganizationID == "" {
			d.OrganizationID = nil
		} else {
			orgID, err := uuid.Parse(*req.OrganizationID)
			if err != nil {
				return nil, fmt.Errorf("invalid organization_id: %w", err)
			}
			d.OrganizationID = &orgID
		}
	}
	if req.Activity != nil {
		if *req.Activity == "" {
			d.Activity = nil
		} else {
			d.Activity = req.Activity
		}
	}
	if err := s.repo.UpdateJobRelationshipDetail(ctx, d); err != nil {
		return nil, err
	}
	if fetched, err := s.repo.FindJobRelationshipDetailByID(ctx, d.ID); err == nil {
		d = fetched
	}
	r := toJobRelationshipDetailResponse(d)
	return &r, nil
}

func (s *Service) DeleteJobRelationshipDetail(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteJobRelationshipDetail(ctx, uid)
}

// =========================================================================
// Job Subordinate Controls (9.13)
// =========================================================================

func (s *Service) CreateJobSubordinateControl(ctx context.Context, req CreateJobSubordinateControlRequest) (*JobSubordinateControlResponse, error) {
	c := &JobSubordinateControl{
		Nomenclature: req.Nomenclature,
		FullCode:     req.FullCode,
	}
	orgID, _ := uuid.Parse(req.OrganizationID)
	c.OrganizationID = &orgID
	if req.JobManagementValueID != nil && *req.JobManagementValueID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueID)
		c.JobManagementValueID = &id
	}
	c.CreatedBy = authctx.GetUserID(ctx)
	c.UpdatedBy = c.CreatedBy
	if err := s.repo.CreateJobSubordinateControl(ctx, c); err != nil {
		return nil, err
	}
	s.recalculateScore(ctx, c.OrganizationID)
	r := toJobSubordinateControlResponse(c)
	return &r, nil
}

func (s *Service) GetJobSubordinateControlByID(ctx context.Context, id string) (*JobSubordinateControlResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	c, err := s.repo.FindJobSubordinateControlByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := toJobSubordinateControlResponse(c)
	return &r, nil
}

func (s *Service) ListJobSubordinateControls(ctx context.Context, page, perPage int, orgID string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var orgUID *uuid.UUID
	if orgID != "" {
		if id, err := uuid.Parse(orgID); err == nil {
			orgUID = &id
		}
	}
	controls, total, err := s.repo.FindAllJobSubordinateControls(ctx, page, perPage, orgUID)
	if err != nil {
		return nil, err
	}
	responses := make([]JobSubordinateControlResponse, 0, len(controls))
	for _, c := range controls {
		responses = append(responses, toJobSubordinateControlResponse(&c))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateJobSubordinateControl(ctx context.Context, id string, req UpdateJobSubordinateControlRequest) (*JobSubordinateControlResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	c, err := s.repo.FindJobSubordinateControlByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	c.UpdatedBy = authctx.GetUserID(ctx)
	if req.Nomenclature != nil {
		c.Nomenclature = *req.Nomenclature
	}
	if req.FullCode != nil {
		c.FullCode = *req.FullCode
	}
	if req.JobManagementValueID != nil && *req.JobManagementValueID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueID)
		c.JobManagementValueID = &id
	}
	if err := s.repo.UpdateJobSubordinateControl(ctx, c); err != nil {
		return nil, err
	}
	s.recalculateScore(ctx, c.OrganizationID)
	r := toJobSubordinateControlResponse(c)
	return &r, nil
}

func (s *Service) DeleteJobSubordinateControl(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if rec, err := s.repo.FindJobSubordinateControlByID(ctx, uid); err == nil {
		defer s.recalculateScore(ctx, rec.OrganizationID)
	}
	return s.repo.DeleteJobSubordinateControl(ctx, uid)
}

// =========================================================================
// Job Assets (9.14)
// =========================================================================

func (s *Service) CreateJobAsset(ctx context.Context, req CreateJobAssetRequest) (*JobAssetResponse, error) {
	a := &JobAsset{
		Nomenclature: req.Nomenclature,
		FullCode:     req.FullCode,
	}
	orgID, _ := uuid.Parse(req.OrganizationID)
	a.OrganizationID = &orgID
	if req.JobManagementValueAssetID != nil && *req.JobManagementValueAssetID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueAssetID)
		a.JobManagementValueAssetID = &id
	}
	if req.JobManagementValueAuthorityID != nil && *req.JobManagementValueAuthorityID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueAuthorityID)
		a.JobManagementValueAuthorityID = &id
	}
	a.CreatedBy = authctx.GetUserID(ctx)
	a.UpdatedBy = a.CreatedBy
	if err := s.repo.CreateJobAsset(ctx, a); err != nil {
		return nil, err
	}
	s.recalculateScore(ctx, a.OrganizationID)
	r := toJobAssetResponse(a)
	return &r, nil
}

func (s *Service) GetJobAssetByID(ctx context.Context, id string) (*JobAssetResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindJobAssetByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := toJobAssetResponse(a)
	return &r, nil
}

func (s *Service) ListJobAssets(ctx context.Context, page, perPage int, orgID string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var orgUID *uuid.UUID
	if orgID != "" {
		if id, err := uuid.Parse(orgID); err == nil {
			orgUID = &id
		}
	}
	assets, total, err := s.repo.FindAllJobAssets(ctx, page, perPage, orgUID)
	if err != nil {
		return nil, err
	}
	responses := make([]JobAssetResponse, 0, len(assets))
	for _, a := range assets {
		responses = append(responses, toJobAssetResponse(&a))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateJobAsset(ctx context.Context, id string, req UpdateJobAssetRequest) (*JobAssetResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindJobAssetByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	a.UpdatedBy = authctx.GetUserID(ctx)
	if req.Nomenclature != nil {
		a.Nomenclature = *req.Nomenclature
	}
	if req.FullCode != nil {
		a.FullCode = *req.FullCode
	}
	if req.JobManagementValueAssetID != nil && *req.JobManagementValueAssetID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueAssetID)
		a.JobManagementValueAssetID = &id
	}
	if req.JobManagementValueAuthorityID != nil && *req.JobManagementValueAuthorityID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueAuthorityID)
		a.JobManagementValueAuthorityID = &id
	}
	if err := s.repo.UpdateJobAsset(ctx, a); err != nil {
		return nil, err
	}
	s.recalculateScore(ctx, a.OrganizationID)
	r := toJobAssetResponse(a)
	return &r, nil
}

func (s *Service) DeleteJobAsset(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if rec, err := s.repo.FindJobAssetByID(ctx, uid); err == nil {
		defer s.recalculateScore(ctx, rec.OrganizationID)
	}
	return s.repo.DeleteJobAsset(ctx, uid)
}

// =========================================================================
// Job Financials (9.15)
// =========================================================================

func (s *Service) CreateJobFinancial(ctx context.Context, req CreateJobFinancialRequest) (*JobFinancialResponse, error) {
	f := &JobFinancial{
		Nomenclature: req.Nomenclature,
		FullCode:     req.FullCode,
		IsAuthorized: req.IsAuthorized,
	}
	orgID, _ := uuid.Parse(req.OrganizationID)
	f.OrganizationID = &orgID
	if req.JobManagementValueCashID != nil && *req.JobManagementValueCashID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueCashID)
		f.JobManagementValueCashID = &id
	}
	if req.JobManagementValueAuthorityID != nil && *req.JobManagementValueAuthorityID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueAuthorityID)
		f.JobManagementValueAuthorityID = &id
	}
	if req.JobManagementValueImpactID != nil && *req.JobManagementValueImpactID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueImpactID)
		f.JobManagementValueImpactID = &id
	}
	f.CreatedBy = authctx.GetUserID(ctx)
	f.UpdatedBy = f.CreatedBy
	if err := s.repo.CreateJobFinancial(ctx, f); err != nil {
		return nil, err
	}
	s.recalculateScore(ctx, f.OrganizationID)
	r := toJobFinancialResponse(f)
	return &r, nil
}

func (s *Service) GetJobFinancialByID(ctx context.Context, id string) (*JobFinancialResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	f, err := s.repo.FindJobFinancialByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := toJobFinancialResponse(f)
	return &r, nil
}

func (s *Service) ListJobFinancials(ctx context.Context, page, perPage int, orgID string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var orgUID *uuid.UUID
	if orgID != "" {
		if id, err := uuid.Parse(orgID); err == nil {
			orgUID = &id
		}
	}
	financials, total, err := s.repo.FindAllJobFinancials(ctx, page, perPage, orgUID)
	if err != nil {
		return nil, err
	}
	responses := make([]JobFinancialResponse, 0, len(financials))
	for _, f := range financials {
		responses = append(responses, toJobFinancialResponse(&f))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateJobFinancial(ctx context.Context, id string, req UpdateJobFinancialRequest) (*JobFinancialResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	f, err := s.repo.FindJobFinancialByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	f.UpdatedBy = authctx.GetUserID(ctx)
	if req.Nomenclature != nil {
		f.Nomenclature = *req.Nomenclature
	}
	if req.FullCode != nil {
		f.FullCode = *req.FullCode
	}
	if req.IsAuthorized != nil {
		f.IsAuthorized = *req.IsAuthorized
	}
	// String kosong ("") = clear field; nil = biarkan nilai lama.
	if req.JobManagementValueCashID != nil {
		if *req.JobManagementValueCashID == "" {
			f.JobManagementValueCashID = nil
		} else if id, err := uuid.Parse(*req.JobManagementValueCashID); err == nil {
			f.JobManagementValueCashID = &id
		}
	}
	if req.JobManagementValueAuthorityID != nil {
		if *req.JobManagementValueAuthorityID == "" {
			f.JobManagementValueAuthorityID = nil
		} else if id, err := uuid.Parse(*req.JobManagementValueAuthorityID); err == nil {
			f.JobManagementValueAuthorityID = &id
		}
	}
	if req.JobManagementValueImpactID != nil {
		if *req.JobManagementValueImpactID == "" {
			f.JobManagementValueImpactID = nil
		} else if id, err := uuid.Parse(*req.JobManagementValueImpactID); err == nil {
			f.JobManagementValueImpactID = &id
		}
	}
	if err := s.repo.UpdateJobFinancial(ctx, f); err != nil {
		return nil, err
	}
	s.recalculateScore(ctx, f.OrganizationID)
	r := toJobFinancialResponse(f)
	return &r, nil
}

func (s *Service) DeleteJobFinancial(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if rec, err := s.repo.FindJobFinancialByID(ctx, uid); err == nil {
		defer s.recalculateScore(ctx, rec.OrganizationID)
	}
	return s.repo.DeleteJobFinancial(ctx, uid)
}

// =========================================================================
// Job Potency Competencies (9.16)
// =========================================================================

func (s *Service) CreateJobPotencyCompetency(ctx context.Context, req CreateJobPotencyCompetencyRequest) (*JobPotencyCompetencyResponse, error) {
	c := &JobPotencyCompetency{}
	orgID, _ := uuid.Parse(req.OrganizationID)
	c.OrganizationID = &orgID
	if req.JobManagementValueID != nil && *req.JobManagementValueID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueID)
		c.JobManagementValueID = &id
	}
	if req.CompetencyID != nil && *req.CompetencyID != "" {
		id, _ := uuid.Parse(*req.CompetencyID)
		c.CompetencyID = &id
	}
	if req.Weight != nil {
		c.Weight = req.Weight
	}
	c.CreatedBy = authctx.GetUserID(ctx)
	c.UpdatedBy = c.CreatedBy
	if err := s.repo.CreateJobPotencyCompetency(ctx, c); err != nil {
		return nil, err
	}
	s.recalculateScore(ctx, c.OrganizationID)
	r := toJobPotencyCompetencyResponse(c)
	return &r, nil
}

func (s *Service) GetJobPotencyCompetencyByID(ctx context.Context, id string) (*JobPotencyCompetencyResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	c, err := s.repo.FindJobPotencyCompetencyByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := toJobPotencyCompetencyResponse(c)
	return &r, nil
}

func (s *Service) ListJobPotencyCompetencies(ctx context.Context, page, perPage int, orgID string) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var orgUID *uuid.UUID
	if orgID != "" {
		if id, err := uuid.Parse(orgID); err == nil {
			orgUID = &id
		}
	}
	competencies, total, err := s.repo.FindAllJobPotencyCompetencies(ctx, page, perPage, orgUID)
	if err != nil {
		return nil, err
	}
	responses := make([]JobPotencyCompetencyResponse, 0, len(competencies))
	for _, c := range competencies {
		responses = append(responses, toJobPotencyCompetencyResponse(&c))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

func (s *Service) UpdateJobPotencyCompetency(ctx context.Context, id string, req UpdateJobPotencyCompetencyRequest) (*JobPotencyCompetencyResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	c, err := s.repo.FindJobPotencyCompetencyByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	c.UpdatedBy = authctx.GetUserID(ctx)
	if req.JobManagementValueID != nil && *req.JobManagementValueID != "" {
		id, _ := uuid.Parse(*req.JobManagementValueID)
		c.JobManagementValueID = &id
	}
	if req.CompetencyID != nil && *req.CompetencyID != "" {
		id, _ := uuid.Parse(*req.CompetencyID)
		c.CompetencyID = &id
	}
	if req.Weight != nil {
		c.Weight = req.Weight
	}
	if err := s.repo.UpdateJobPotencyCompetency(ctx, c); err != nil {
		return nil, err
	}
	s.recalculateScore(ctx, c.OrganizationID)
	r := toJobPotencyCompetencyResponse(c)
	return &r, nil
}

func (s *Service) DeleteJobPotencyCompetency(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if rec, err := s.repo.FindJobPotencyCompetencyByID(ctx, uid); err == nil {
		defer s.recalculateScore(ctx, rec.OrganizationID)
	}
	return s.repo.DeleteJobPotencyCompetency(ctx, uid)
}

// recalculateScore menghitung ulang skor jabatan organisasi setelah sebuah
// section (yang memengaruhi skor) disimpan/diubah/dihapus. Kegagalan menghitung
// TIDAK menggagalkan operasi utama section — cukup dicatat di log.
//
// Catatan perf: frontend menyimpan beberapa baris potency dengan request
// per-baris (saveCardRows), sehingga satu klik save bisa memicu beberapa kali
// recalc — tradeoff yang diterima untuk CRUD admin skala kecil.
func (s *Service) recalculateScore(ctx context.Context, orgID *uuid.UUID) {
	if orgID == nil || *orgID == uuid.Nil {
		return
	}
	if _, err := s.RecalculateJobScore(ctx, orgID.String()); err != nil && s.logger != nil {
		s.logger.Warn("recalculate job score after section save failed",
			zap.String("organization_id", orgID.String()), zap.Error(err))
	}
}

// =========================================================================
// Job Scores (9.17)
// =========================================================================

func (s *Service) UpsertJobScore(ctx context.Context, orgID string, req UpdateJobScoreRequest) (*JobScoreResponse, error) {
	uid, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}
	// Kontrak recalculate: body kosong (semua field nil — mis. { components: null }
	// dari tombol "Recalculate" di JobScoreSection) → hitung ulang server-side.
	if req.JobValueWithFinancial == nil && req.JobValueWithoutFinancial == nil &&
		req.HasFinancialAuthority == nil && req.Components == nil && req.SubComponentPoints == nil {
		return s.RecalculateJobScore(ctx, orgID)
	}
	score := &JobScore{
		OrganizationID: &uid,
	}
	if req.JobValueWithFinancial != nil {
		score.JobValueWithFinancial = *req.JobValueWithFinancial
	}
	if req.JobValueWithoutFinancial != nil {
		score.JobValueWithoutFinancial = *req.JobValueWithoutFinancial
	}
	if req.HasFinancialAuthority != nil {
		score.HasFinancialAuthority = *req.HasFinancialAuthority
	}
	if req.Components != nil {
		score.Components = req.Components
	}
	if req.SubComponentPoints != nil {
		score.SubComponentPoints = req.SubComponentPoints
		// Jalur manual: turunkan is_complete/completed_at dari sub_component_points
		// yang dikirim (bila valid) — konsisten dengan legacy yang selalu menghitungnya.
		var sub subComponentPoints
		if err := json.Unmarshal([]byte(*req.SubComponentPoints), &sub); err == nil {
			score.IsComplete = isResultComplete(sub)
			if score.IsComplete {
				now := time.Now()
				score.CompletedAt = &now
			}
		}
	}
	if err := s.repo.UpsertJobScore(ctx, score); err != nil {
		return nil, err
	}
	r := toJobScoreResponse(score)
	return &r, nil
}

func (s *Service) GetJobScoreByOrganization(ctx context.Context, orgID string) (*JobScoreResponse, error) {
	uid, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}
	score, err := s.repo.FindJobScoreByOrganizationID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := toJobScoreResponse(score)
	return &r, nil
}

// RecalculateJobScore menghitung ulang skor jabatan dari data section
// (menggunakan Calculator) lalu menyimpannya ke job_management_scores.
func (s *Service) RecalculateJobScore(ctx context.Context, orgID string) (*JobScoreResponse, error) {
	uid, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}
	result, err := s.calc.CalculateForOrganization(ctx, uid)
	if err != nil {
		return nil, err
	}
	score, err := s.scoreFromResult(uid, result)
	if err != nil {
		return nil, err
	}
	if err := s.repo.UpsertJobScore(ctx, score); err != nil {
		return nil, err
	}
	r := toJobScoreResponse(score)
	return &r, nil
}

// RecalculateJobScores menghitung ulang skor untuk banyak organisasi sekaligus
// (re-kalkulasi massal, mis. saat seed/migrasi data — seperti legacy
// calculateForOrganizationIds).
func (s *Service) RecalculateJobScores(ctx context.Context, orgIDs []string) (map[string]*JobScoreResponse, error) {
	ids := make([]uuid.UUID, 0, len(orgIDs))
	for _, o := range orgIDs {
		id, err := uuid.Parse(o)
		if err != nil {
			return nil, fmt.Errorf("invalid organization_id: %w", err)
		}
		ids = append(ids, id)
	}
	results, err := s.calc.CalculateForOrganizationIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	responses := make(map[string]*JobScoreResponse, len(results))
	for orgID, result := range results {
		score, err := s.scoreFromResult(orgID, result)
		if err != nil {
			return nil, err
		}
		if err := s.repo.UpsertJobScore(ctx, score); err != nil {
			return nil, err
		}
		r := toJobScoreResponse(score)
		responses[orgID.String()] = &r
	}
	return responses, nil
}

// scoreFromResult membungkus hasil kalkulasi menjadi model JobScore siap simpan.
func (s *Service) scoreFromResult(orgID uuid.UUID, result *JobScoreResult) (*JobScore, error) {
	componentsJSON, err := json.Marshal(result.Components)
	if err != nil {
		return nil, fmt.Errorf("marshal components: %w", err)
	}
	subJSON, err := json.Marshal(result.SubComponents)
	if err != nil {
		return nil, fmt.Errorf("marshal sub_component_points: %w", err)
	}
	componentsStr := string(componentsJSON)
	subStr := string(subJSON)
	now := time.Now()
	// completed_at diisi hanya jika skor lengkap (sama seperti legacy persistResults)
	var completedAt *time.Time
	if result.IsComplete {
		completedAt = &now
	}
	return &JobScore{
		OrganizationID:           &orgID,
		JobValueWithFinancial:    result.Totals.WithFinancial,
		JobValueWithoutFinancial: result.Totals.WithoutFinancial,
		HasFinancialAuthority:    result.HasFinancialAuthority,
		Components:               &componentsStr,
		SubComponentPoints:       &subStr,
		CalculatedAt:             &now,
		IsComplete:               result.IsComplete,
		CompletedAt:              completedAt,
	}, nil
}

func (s *Service) ListJobScores(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	scores, total, err := s.repo.FindAllJobScores(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]JobScoreResponse, 0, len(scores))
	for _, s := range scores {
		responses = append(responses, toJobScoreResponse(&s))
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	return &PaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}

// =========================================================================
// Job Competency Groups (9.18)
// =========================================================================

func (s *Service) CreateJobCompetencyGroup(ctx context.Context, req CreateJobCompetencyGroupRequest) (*JobCompetencyGroupResponse, error) {
	uid, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}
	g := &JobCompetencyGroup{
		OrganizationID: &uid,
		Category:       req.Category,
		Weight:         req.Weight,
	}
	if err := s.repo.CreateJobCompetencyGroup(ctx, g); err != nil {
		return nil, err
	}
	r := toJobCompetencyGroupResponse(g)
	return &r, nil
}

func (s *Service) GetJobCompetencyGroupByID(ctx context.Context, id string) (*JobCompetencyGroupResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	g, err := s.repo.FindJobCompetencyGroupByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := toJobCompetencyGroupResponse(g)
	return &r, nil
}

func (s *Service) ListJobCompetencyGroups(ctx context.Context, orgID string) ([]JobCompetencyGroupResponse, error) {
	uid, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}
	groups, err := s.repo.FindJobCompetencyGroupsByOrganization(ctx, uid)
	if err != nil {
		return nil, err
	}
	responses := make([]JobCompetencyGroupResponse, 0, len(groups))
	for _, g := range groups {
		responses = append(responses, toJobCompetencyGroupResponse(&g))
	}
	return responses, nil
}

func (s *Service) UpdateJobCompetencyGroup(ctx context.Context, id string, req UpdateJobCompetencyGroupRequest) (*JobCompetencyGroupResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	g, err := s.repo.FindJobCompetencyGroupByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Category != nil {
		g.Category = *req.Category
	}
	if req.Weight != nil {
		g.Weight = *req.Weight
	}
	if err := s.repo.UpdateJobCompetencyGroup(ctx, g); err != nil {
		return nil, err
	}
	r := toJobCompetencyGroupResponse(g)
	return &r, nil
}

func (s *Service) DeleteJobCompetencyGroup(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteJobCompetencyGroup(ctx, uid)
}
