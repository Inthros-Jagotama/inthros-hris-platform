package setting

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 500
)

type Service struct {
	repo   *Repository
	logger *zap.Logger
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// ── Zone CRUD ──
func (s *Service) CreateZone(ctx context.Context, req CreateZoneRequest) (*ZoneResponse, error) {
	isActive := true
	if req.IsActive != nil { isActive = *req.IsActive }
	zone := &Zone{Code: req.Code, Name: req.Name, Zone: req.Name, Region: req.Region, IsActive: isActive, SortOrder: req.SortOrder}
	if err := s.repo.CreateZone(ctx, zone); err != nil { return nil, err }
	s.logger.Info("Zone created", zap.String("id", zone.ID.String()), zap.String("code", zone.Code))
	resp := zone.ToResponse()
	return &resp, nil
}
func (s *Service) GetZoneByID(ctx context.Context, id string) (*ZoneResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	zone, err := s.repo.FindZoneByID(ctx, uid)
	if err != nil { return nil, err }
	resp := zone.ToResponse()
	return &resp, nil
}
func (s *Service) ListZones(ctx context.Context, page, perPage int) (*ZonePaginatedResponse, error) {
	if page < 1 { page = defaultPage }
	if perPage < 1 { perPage = defaultPerPage }
	if perPage > maxPerPage { perPage = maxPerPage }
	zones, total, err := s.repo.FindAllZones(ctx, page, perPage)
	if err != nil { return nil, err }
	responses := make([]ZoneResponse, len(zones))
	for i, z := range zones { responses[i] = z.ToResponse() }
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	return &ZonePaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}
func (s *Service) ListZonesActive(ctx context.Context) ([]ZoneResponse, error) {
	zones, err := s.repo.FindAllZonesActive(ctx)
	if err != nil { return nil, err }
	responses := make([]ZoneResponse, len(zones))
	for i, z := range zones { responses[i] = z.ToResponse() }
	return responses, nil
}
func (s *Service) UpdateZone(ctx context.Context, id string, req UpdateZoneRequest) (*ZoneResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	zone, err := s.repo.FindZoneByID(ctx, uid)
	if err != nil { return nil, err }
	if req.Code != nil { zone.Code = *req.Code }
	if req.Name != nil { zone.Name = *req.Name; zone.Zone = *req.Name }
	if req.Region != nil { zone.Region = *req.Region }
	if req.IsActive != nil { zone.IsActive = *req.IsActive }
	if req.SortOrder != nil { zone.SortOrder = *req.SortOrder }
	if err := s.repo.UpdateZone(ctx, zone); err != nil { return nil, err }
	resp := zone.ToResponse()
	return &resp, nil
}
func (s *Service) DeleteZone(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil { return fmt.Errorf("invalid id: %w", err) }
	return s.repo.DeleteZone(ctx, uid)
}

// ── Province CRUD ──
func (s *Service) CreateProvince(ctx context.Context, req CreateProvinceRequest) (*ProvinceResponse, error) {
	p := &Province{ID: req.Code, Code: req.Code, Name: req.Name}
	if err := s.repo.CreateProvince(ctx, p); err != nil { return nil, err }
	s.logger.Info("Province created", zap.String("id", p.ID), zap.String("code", p.Code))
	resp := p.ToResponse()
	return &resp, nil
}
func (s *Service) GetProvinceByID(ctx context.Context, id string) (*ProvinceResponse, error) {
	p, err := s.repo.FindProvinceByID(ctx, id)
	if err != nil { return nil, err }
	resp := p.ToResponse()
	return &resp, nil
}
func (s *Service) ListProvinces(ctx context.Context, page, perPage int) (*ProvincePaginatedResponse, error) {
	if page < 1 { page = defaultPage }
	if perPage < 1 { perPage = defaultPerPage }
	if perPage > maxPerPage { perPage = maxPerPage }
	provinces, total, err := s.repo.FindAllProvinces(ctx, page, perPage)
	if err != nil { return nil, err }
	responses := make([]ProvinceResponse, len(provinces))
	for i, p := range provinces { responses[i] = p.ToResponse() }
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	return &ProvincePaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}
func (s *Service) ListAllProvinces(ctx context.Context) ([]ProvinceResponse, error) {
	provinces, err := s.repo.FindAllProvincesList(ctx)
	if err != nil { return nil, err }
	responses := make([]ProvinceResponse, len(provinces))
	for i, p := range provinces { responses[i] = p.ToResponse() }
	return responses, nil
}
func (s *Service) UpdateProvince(ctx context.Context, id string, req UpdateProvinceRequest) (*ProvinceResponse, error) {
	p, err := s.repo.FindProvinceByID(ctx, id)
	if err != nil { return nil, err }
	if req.Code != nil { p.Code = *req.Code; p.ID = *req.Code }
	if req.Name != nil { p.Name = *req.Name }
	if err := s.repo.UpdateProvince(ctx, p); err != nil { return nil, err }
	resp := p.ToResponse()
	return &resp, nil
}
func (s *Service) DeleteProvince(ctx context.Context, id string) error {
	return s.repo.DeleteProvince(ctx, id)
}

// ── Regency CRUD ──
func (s *Service) CreateRegency(ctx context.Context, req CreateRegencyRequest) (*RegencyResponse, error) {
	reg := &Regency{ID: req.Code, Code: req.Code, Name: req.Name, ProvinceID: req.ProvinceID}
	if err := s.repo.CreateRegency(ctx, reg); err != nil { return nil, err }
	s.logger.Info("Regency created", zap.String("id", reg.ID), zap.String("code", reg.Code))
	resp := reg.ToResponse()
	return &resp, nil
}
func (s *Service) GetRegencyByID(ctx context.Context, id string) (*RegencyResponse, error) {
	reg, err := s.repo.FindRegencyByID(ctx, id)
	if err != nil { return nil, err }
	resp := reg.ToResponse()
	return &resp, nil
}
func (s *Service) ListRegencies(ctx context.Context, page, perPage int) (*RegencyPaginatedResponse, error) {
	if page < 1 { page = defaultPage }
	if perPage < 1 { perPage = defaultPerPage }
	if perPage > maxPerPage { perPage = maxPerPage }
	regencies, total, err := s.repo.FindAllRegencies(ctx, page, perPage)
	if err != nil { return nil, err }
	responses := make([]RegencyResponse, len(regencies))
	for i, r := range regencies { responses[i] = r.ToResponse() }
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	return &RegencyPaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}
func (s *Service) ListRegenciesByProvince(ctx context.Context, provinceID string) ([]RegencyResponse, error) {
	regencies, err := s.repo.FindRegenciesByProvince(ctx, provinceID)
	if err != nil { return nil, err }
	responses := make([]RegencyResponse, len(regencies))
	for i, r := range regencies { responses[i] = r.ToResponse() }
	return responses, nil
}
func (s *Service) UpdateRegency(ctx context.Context, id string, req UpdateRegencyRequest) (*RegencyResponse, error) {
	reg, err := s.repo.FindRegencyByID(ctx, id)
	if err != nil { return nil, err }
	if req.Code != nil { reg.Code = *req.Code; reg.ID = *req.Code }
	if req.Name != nil { reg.Name = *req.Name }
	if req.ProvinceID != nil { reg.ProvinceID = *req.ProvinceID }
	if err := s.repo.UpdateRegency(ctx, reg); err != nil { return nil, err }
	resp := reg.ToResponse()
	return &resp, nil
}
func (s *Service) DeleteRegency(ctx context.Context, id string) error {
	return s.repo.DeleteRegency(ctx, id)
}

// ── District CRUD ──
func (s *Service) CreateDistrict(ctx context.Context, req CreateDistrictRequest) (*DistrictResponse, error) {
	d := &District{ID: req.Code, Code: req.Code, Name: req.Name, RegencyID: req.RegencyID}
	if err := s.repo.CreateDistrict(ctx, d); err != nil { return nil, err }
	s.logger.Info("District created", zap.String("id", d.ID), zap.String("code", d.Code))
	resp := d.ToResponse()
	return &resp, nil
}
func (s *Service) GetDistrictByID(ctx context.Context, id string) (*DistrictResponse, error) {
	d, err := s.repo.FindDistrictByID(ctx, id)
	if err != nil { return nil, err }
	resp := d.ToResponse()
	return &resp, nil
}
func (s *Service) ListDistricts(ctx context.Context, page, perPage int) (*DistrictPaginatedResponse, error) {
	if page < 1 { page = defaultPage }
	if perPage < 1 { perPage = defaultPerPage }
	if perPage > maxPerPage { perPage = maxPerPage }
	districts, total, err := s.repo.FindAllDistricts(ctx, page, perPage)
	if err != nil { return nil, err }
	responses := make([]DistrictResponse, len(districts))
	for i, d := range districts { responses[i] = d.ToResponse() }
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	return &DistrictPaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}
func (s *Service) ListDistrictsByRegency(ctx context.Context, regencyID string) ([]DistrictResponse, error) {
	districts, err := s.repo.FindDistrictsByRegency(ctx, regencyID)
	if err != nil { return nil, err }
	responses := make([]DistrictResponse, len(districts))
	for i, d := range districts { responses[i] = d.ToResponse() }
	return responses, nil
}
func (s *Service) UpdateDistrict(ctx context.Context, id string, req UpdateDistrictRequest) (*DistrictResponse, error) {
	d, err := s.repo.FindDistrictByID(ctx, id)
	if err != nil { return nil, err }
	if req.Code != nil { d.Code = *req.Code; d.ID = *req.Code }
	if req.Name != nil { d.Name = *req.Name }
	if req.RegencyID != nil { d.RegencyID = *req.RegencyID }
	if err := s.repo.UpdateDistrict(ctx, d); err != nil { return nil, err }
	resp := d.ToResponse()
	return &resp, nil
}
func (s *Service) DeleteDistrict(ctx context.Context, id string) error {
	return s.repo.DeleteDistrict(ctx, id)
}

// ── Village CRUD ──
func (s *Service) CreateVillage(ctx context.Context, req CreateVillageRequest) (*VillageResponse, error) {
	v := &Village{ID: req.Code, Code: req.Code, Name: req.Name, DistrictID: req.DistrictID}
	if err := s.repo.CreateVillage(ctx, v); err != nil { return nil, err }
	s.logger.Info("Village created", zap.String("id", v.ID), zap.String("code", v.Code))
	resp := v.ToResponse()
	return &resp, nil
}
func (s *Service) GetVillageByID(ctx context.Context, id string) (*VillageResponse, error) {
	v, err := s.repo.FindVillageByID(ctx, id)
	if err != nil { return nil, err }
	resp := v.ToResponse()
	return &resp, nil
}
func (s *Service) ListVillages(ctx context.Context, page, perPage int) (*VillagePaginatedResponse, error) {
	if page < 1 { page = defaultPage }
	if perPage < 1 { perPage = defaultPerPage }
	if perPage > maxPerPage { perPage = maxPerPage }
	villages, total, err := s.repo.FindAllVillages(ctx, page, perPage)
	if err != nil { return nil, err }
	responses := make([]VillageResponse, len(villages))
	for i, v := range villages { responses[i] = v.ToResponse() }
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	return &VillagePaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}
func (s *Service) ListVillagesByDistrict(ctx context.Context, districtID string) ([]VillageResponse, error) {
	villages, err := s.repo.FindVillagesByDistrict(ctx, districtID)
	if err != nil { return nil, err }
	responses := make([]VillageResponse, len(villages))
	for i, v := range villages { responses[i] = v.ToResponse() }
	return responses, nil
}
func (s *Service) UpdateVillage(ctx context.Context, id string, req UpdateVillageRequest) (*VillageResponse, error) {
	v, err := s.repo.FindVillageByID(ctx, id)
	if err != nil { return nil, err }
	if req.Code != nil { v.Code = *req.Code; v.ID = *req.Code }
	if req.Name != nil { v.Name = *req.Name }
	if req.DistrictID != nil { v.DistrictID = *req.DistrictID }
	if err := s.repo.UpdateVillage(ctx, v); err != nil { return nil, err }
	resp := v.ToResponse()
	return &resp, nil
}
func (s *Service) DeleteVillage(ctx context.Context, id string) error {
	return s.repo.DeleteVillage(ctx, id)
}

// ── Education CRUD ──
func (s *Service) CreateEducation(ctx context.Context, req CreateEducationRequest) (*EducationResponse, error) {
	e := &Education{Code: req.Code, Name: req.Name, SortOrder: req.SortOrder}
	if err := s.repo.CreateEducation(ctx, e); err != nil { return nil, err }
	s.logger.Info("Education created", zap.String("id", e.ID.String()), zap.String("code", e.Code))
	resp := e.ToResponse()
	return &resp, nil
}
func (s *Service) GetEducationByID(ctx context.Context, id string) (*EducationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	e, err := s.repo.FindEducationByID(ctx, uid)
	if err != nil { return nil, err }
	resp := e.ToResponse()
	return &resp, nil
}
func (s *Service) ListEducations(ctx context.Context, page, perPage int) (*EducationPaginatedResponse, error) {
	if page < 1 { page = defaultPage }
	if perPage < 1 { perPage = defaultPerPage }
	if perPage > maxPerPage { perPage = maxPerPage }
	list, total, err := s.repo.FindAllEducations(ctx, page, perPage)
	if err != nil { return nil, err }
	responses := make([]EducationResponse, len(list))
	for i, e := range list { responses[i] = e.ToResponse() }
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	return &EducationPaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}
func (s *Service) UpdateEducation(ctx context.Context, id string, req UpdateEducationRequest) (*EducationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	e, err := s.repo.FindEducationByID(ctx, uid)
	if err != nil { return nil, err }
	if req.Code != nil { e.Code = *req.Code }
	if req.Name != nil { e.Name = *req.Name }
	if req.SortOrder != nil { e.SortOrder = *req.SortOrder }
	if err := s.repo.UpdateEducation(ctx, e); err != nil { return nil, err }
	resp := e.ToResponse()
	return &resp, nil
}
func (s *Service) DeleteEducation(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil { return fmt.Errorf("invalid id: %w", err) }
	return s.repo.DeleteEducation(ctx, uid)
}

// ── Religion CRUD ──
func (s *Service) CreateReligion(ctx context.Context, req CreateReligionRequest) (*ReligionResponse, error) {
	rel := &Religion{Code: req.Code, Name: req.Name, SortOrder: req.SortOrder}
	if err := s.repo.CreateReligion(ctx, rel); err != nil { return nil, err }
	s.logger.Info("Religion created", zap.String("id", rel.ID.String()), zap.String("code", rel.Code))
	resp := rel.ToResponse()
	return &resp, nil
}
func (s *Service) GetReligionByID(ctx context.Context, id string) (*ReligionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	rel, err := s.repo.FindReligionByID(ctx, uid)
	if err != nil { return nil, err }
	resp := rel.ToResponse()
	return &resp, nil
}
func (s *Service) ListReligions(ctx context.Context, page, perPage int) (*ReligionPaginatedResponse, error) {
	if page < 1 { page = defaultPage }
	if perPage < 1 { perPage = defaultPerPage }
	if perPage > maxPerPage { perPage = maxPerPage }
	list, total, err := s.repo.FindAllReligions(ctx, page, perPage)
	if err != nil { return nil, err }
	responses := make([]ReligionResponse, len(list))
	for i, r := range list { responses[i] = r.ToResponse() }
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	return &ReligionPaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}
func (s *Service) UpdateReligion(ctx context.Context, id string, req UpdateReligionRequest) (*ReligionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	rel, err := s.repo.FindReligionByID(ctx, uid)
	if err != nil { return nil, err }
	if req.Code != nil { rel.Code = *req.Code }
	if req.Name != nil { rel.Name = *req.Name }
	if req.SortOrder != nil { rel.SortOrder = *req.SortOrder }
	if err := s.repo.UpdateReligion(ctx, rel); err != nil { return nil, err }
	resp := rel.ToResponse()
	return &resp, nil
}
func (s *Service) DeleteReligion(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil { return fmt.Errorf("invalid id: %w", err) }
	return s.repo.DeleteReligion(ctx, uid)
}

// ── MaritalStatus CRUD ──
func (s *Service) CreateMaritalStatus(ctx context.Context, req CreateMaritalStatusRequest) (*MaritalStatusResponse, error) {
	m := &MaritalStatus{Code: req.Code, Name: req.Name, SortOrder: req.SortOrder}
	if err := s.repo.CreateMaritalStatus(ctx, m); err != nil { return nil, err }
	s.logger.Info("MaritalStatus created", zap.String("id", m.ID.String()), zap.String("code", m.Code))
	resp := m.ToResponse()
	return &resp, nil
}
func (s *Service) GetMaritalStatusByID(ctx context.Context, id string) (*MaritalStatusResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	m, err := s.repo.FindMaritalStatusByID(ctx, uid)
	if err != nil { return nil, err }
	resp := m.ToResponse()
	return &resp, nil
}
func (s *Service) ListMaritalStatuses(ctx context.Context, page, perPage int) (*MaritalStatusPaginatedResponse, error) {
	if page < 1 { page = defaultPage }
	if perPage < 1 { perPage = defaultPerPage }
	if perPage > maxPerPage { perPage = maxPerPage }
	list, total, err := s.repo.FindAllMaritalStatuses(ctx, page, perPage)
	if err != nil { return nil, err }
	responses := make([]MaritalStatusResponse, len(list))
	for i, m := range list { responses[i] = m.ToResponse() }
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	return &MaritalStatusPaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}
func (s *Service) UpdateMaritalStatus(ctx context.Context, id string, req UpdateMaritalStatusRequest) (*MaritalStatusResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	m, err := s.repo.FindMaritalStatusByID(ctx, uid)
	if err != nil { return nil, err }
	if req.Code != nil { m.Code = *req.Code }
	if req.Name != nil { m.Name = *req.Name }
	if req.SortOrder != nil { m.SortOrder = *req.SortOrder }
	if err := s.repo.UpdateMaritalStatus(ctx, m); err != nil { return nil, err }
	resp := m.ToResponse()
	return &resp, nil
}
func (s *Service) DeleteMaritalStatus(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil { return fmt.Errorf("invalid id: %w", err) }
	return s.repo.DeleteMaritalStatus(ctx, uid)
}

// ── Bank CRUD ──
func (s *Service) CreateBank(ctx context.Context, req CreateBankRequest) (*BankResponse, error) {
	b := &Bank{Code: req.Code, Name: req.Name, SortOrder: req.SortOrder}
	if err := s.repo.CreateBank(ctx, b); err != nil { return nil, err }
	s.logger.Info("Bank created", zap.String("id", b.ID.String()), zap.String("code", b.Code))
	resp := b.ToResponse()
	return &resp, nil
}
func (s *Service) GetBankByID(ctx context.Context, id string) (*BankResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	b, err := s.repo.FindBankByID(ctx, uid)
	if err != nil { return nil, err }
	resp := b.ToResponse()
	return &resp, nil
}
func (s *Service) ListBanks(ctx context.Context, page, perPage int) (*BankPaginatedResponse, error) {
	if page < 1 { page = defaultPage }
	if perPage < 1 { perPage = defaultPerPage }
	if perPage > maxPerPage { perPage = maxPerPage }
	list, total, err := s.repo.FindAllBanks(ctx, page, perPage)
	if err != nil { return nil, err }
	responses := make([]BankResponse, len(list))
	for i, b := range list { responses[i] = b.ToResponse() }
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	return &BankPaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}
func (s *Service) UpdateBank(ctx context.Context, id string, req UpdateBankRequest) (*BankResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	b, err := s.repo.FindBankByID(ctx, uid)
	if err != nil { return nil, err }
	if req.Code != nil { b.Code = *req.Code }
	if req.Name != nil { b.Name = *req.Name }
	if req.SortOrder != nil { b.SortOrder = *req.SortOrder }
	if err := s.repo.UpdateBank(ctx, b); err != nil { return nil, err }
	resp := b.ToResponse()
	return &resp, nil
}
func (s *Service) DeleteBank(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil { return fmt.Errorf("invalid id: %w", err) }
	return s.repo.DeleteBank(ctx, uid)
}

// ── Nationality CRUD ──
func (s *Service) CreateNationality(ctx context.Context, req CreateNationalityRequest) (*NationalityResponse, error) {
	n := &Nationality{Code: req.Code, Name: req.Name, SortOrder: req.SortOrder}
	if err := s.repo.CreateNationality(ctx, n); err != nil { return nil, err }
	s.logger.Info("Nationality created", zap.String("id", n.ID.String()), zap.String("code", n.Code))
	resp := n.ToResponse()
	return &resp, nil
}
func (s *Service) GetNationalityByID(ctx context.Context, id string) (*NationalityResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	n, err := s.repo.FindNationalityByID(ctx, uid)
	if err != nil { return nil, err }
	resp := n.ToResponse()
	return &resp, nil
}
func (s *Service) ListNationalities(ctx context.Context, page, perPage int) (*NationalityPaginatedResponse, error) {
	if page < 1 { page = defaultPage }
	if perPage < 1 { perPage = defaultPerPage }
	if perPage > maxPerPage { perPage = maxPerPage }
	list, total, err := s.repo.FindAllNationalities(ctx, page, perPage)
	if err != nil { return nil, err }
	responses := make([]NationalityResponse, len(list))
	for i, n := range list { responses[i] = n.ToResponse() }
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	return &NationalityPaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}
func (s *Service) UpdateNationality(ctx context.Context, id string, req UpdateNationalityRequest) (*NationalityResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	n, err := s.repo.FindNationalityByID(ctx, uid)
	if err != nil { return nil, err }
	if req.Code != nil { n.Code = *req.Code }
	if req.Name != nil { n.Name = *req.Name }
	if req.SortOrder != nil { n.SortOrder = *req.SortOrder }
	if err := s.repo.UpdateNationality(ctx, n); err != nil { return nil, err }
	resp := n.ToResponse()
	return &resp, nil
}
func (s *Service) DeleteNationality(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil { return fmt.Errorf("invalid id: %w", err) }
	return s.repo.DeleteNationality(ctx, uid)
}

// ── RelationshipType CRUD ──
func (s *Service) CreateRelationshipType(ctx context.Context, req CreateRelationshipTypeRequest) (*RelationshipTypeResponse, error) {
	rt := &RelationshipType{Code: req.Code, Name: req.Name, SortOrder: req.SortOrder}
	if err := s.repo.CreateRelationshipType(ctx, rt); err != nil { return nil, err }
	s.logger.Info("RelationshipType created", zap.String("id", rt.ID.String()), zap.String("code", rt.Code))
	resp := rt.ToResponse()
	return &resp, nil
}
func (s *Service) GetRelationshipTypeByID(ctx context.Context, id string) (*RelationshipTypeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	rt, err := s.repo.FindRelationshipTypeByID(ctx, uid)
	if err != nil { return nil, err }
	resp := rt.ToResponse()
	return &resp, nil
}
func (s *Service) ListRelationshipTypes(ctx context.Context, page, perPage int) (*RelationshipTypePaginatedResponse, error) {
	if page < 1 { page = defaultPage }
	if perPage < 1 { perPage = defaultPerPage }
	if perPage > maxPerPage { perPage = maxPerPage }
	list, total, err := s.repo.FindAllRelationshipTypes(ctx, page, perPage)
	if err != nil { return nil, err }
	responses := make([]RelationshipTypeResponse, len(list))
	for i, rt := range list { responses[i] = rt.ToResponse() }
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	return &RelationshipTypePaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}
func (s *Service) UpdateRelationshipType(ctx context.Context, id string, req UpdateRelationshipTypeRequest) (*RelationshipTypeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	rt, err := s.repo.FindRelationshipTypeByID(ctx, uid)
	if err != nil { return nil, err }
	if req.Code != nil { rt.Code = *req.Code }
	if req.Name != nil { rt.Name = *req.Name }
	if req.SortOrder != nil { rt.SortOrder = *req.SortOrder }
	if err := s.repo.UpdateRelationshipType(ctx, rt); err != nil { return nil, err }
	resp := rt.ToResponse()
	return &resp, nil
}
func (s *Service) DeleteRelationshipType(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil { return fmt.Errorf("invalid id: %w", err) }
	return s.repo.DeleteRelationshipType(ctx, uid)
}

// ── EmploymentStatus CRUD ──
func (s *Service) CreateEmploymentStatus(ctx context.Context, req CreateEmploymentStatusRequest) (*EmploymentStatusResponse, error) {
	es := &EmploymentStatus{Code: req.Code, Name: req.Name, SortOrder: req.SortOrder}
	if err := s.repo.CreateEmploymentStatus(ctx, es); err != nil { return nil, err }
	s.logger.Info("EmploymentStatus created", zap.String("id", es.ID.String()), zap.String("code", es.Code))
	resp := es.ToResponse()
	return &resp, nil
}
func (s *Service) GetEmploymentStatusByID(ctx context.Context, id string) (*EmploymentStatusResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	es, err := s.repo.FindEmploymentStatusByID(ctx, uid)
	if err != nil { return nil, err }
	resp := es.ToResponse()
	return &resp, nil
}
func (s *Service) ListEmploymentStatuses(ctx context.Context, page, perPage int) (*EmploymentStatusPaginatedResponse, error) {
	if page < 1 { page = defaultPage }
	if perPage < 1 { perPage = defaultPerPage }
	if perPage > maxPerPage { perPage = maxPerPage }
	list, total, err := s.repo.FindAllEmploymentStatuses(ctx, page, perPage)
	if err != nil { return nil, err }
	responses := make([]EmploymentStatusResponse, len(list))
	for i, es := range list { responses[i] = es.ToResponse() }
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	return &EmploymentStatusPaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}
func (s *Service) UpdateEmploymentStatus(ctx context.Context, id string, req UpdateEmploymentStatusRequest) (*EmploymentStatusResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	es, err := s.repo.FindEmploymentStatusByID(ctx, uid)
	if err != nil { return nil, err }
	if req.Code != nil { es.Code = *req.Code }
	if req.Name != nil { es.Name = *req.Name }
	if req.SortOrder != nil { es.SortOrder = *req.SortOrder }
	if err := s.repo.UpdateEmploymentStatus(ctx, es); err != nil { return nil, err }
	resp := es.ToResponse()
	return &resp, nil
}
func (s *Service) DeleteEmploymentStatus(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil { return fmt.Errorf("invalid id: %w", err) }
	return s.repo.DeleteEmploymentStatus(ctx, uid)
}

// ── JobFamily CRUD ──
func (s *Service) CreateJobFamily(ctx context.Context, req CreateJobFamilyRequest) (*JobFamilyResponse, error) {
	jf := &JobFamily{Code: req.Code, Name: req.Name, Description: req.Description, SortOrder: req.SortOrder}
	if err := s.repo.CreateJobFamily(ctx, jf); err != nil { return nil, err }
	s.logger.Info("JobFamily created", zap.String("id", jf.ID.String()), zap.String("code", jf.Code))
	resp := jf.ToResponse()
	return &resp, nil
}
func (s *Service) GetJobFamilyByID(ctx context.Context, id string) (*JobFamilyResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	jf, err := s.repo.FindJobFamilyByID(ctx, uid)
	if err != nil { return nil, err }
	resp := jf.ToResponse()
	return &resp, nil
}
func (s *Service) ListJobFamilies(ctx context.Context, page, perPage int) (*JobFamilyPaginatedResponse, error) {
	if page < 1 { page = defaultPage }
	if perPage < 1 { perPage = defaultPerPage }
	if perPage > maxPerPage { perPage = maxPerPage }
	list, total, err := s.repo.FindAllJobFamilies(ctx, page, perPage)
	if err != nil { return nil, err }
	responses := make([]JobFamilyResponse, len(list))
	for i, jf := range list { responses[i] = jf.ToResponse() }
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	return &JobFamilyPaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}
func (s *Service) UpdateJobFamily(ctx context.Context, id string, req UpdateJobFamilyRequest) (*JobFamilyResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	jf, err := s.repo.FindJobFamilyByID(ctx, uid)
	if err != nil { return nil, err }
	if req.Code != nil { jf.Code = *req.Code }
	if req.Name != nil { jf.Name = *req.Name }
	if req.Description != nil { jf.Description = *req.Description }
	if req.SortOrder != nil { jf.SortOrder = *req.SortOrder }
	if err := s.repo.UpdateJobFamily(ctx, jf); err != nil { return nil, err }
	resp := jf.ToResponse()
	return &resp, nil
}
func (s *Service) DeleteJobFamily(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil { return fmt.Errorf("invalid id: %w", err) }
	return s.repo.DeleteJobFamily(ctx, uid)
}

// ── SalaryGrade CRUD ──
func (s *Service) CreateSalaryGrade(ctx context.Context, req CreateSalaryGradeRequest) (*SalaryGradeResponse, error) {
	sg := &SalaryGrade{Code: req.Code, Name: req.Name, Description: req.Description, MinAmount: req.MinAmount, MaxAmount: req.MaxAmount, SortOrder: req.SortOrder}
	if err := s.repo.CreateSalaryGrade(ctx, sg); err != nil { return nil, err }
	s.logger.Info("SalaryGrade created", zap.String("id", sg.ID.String()), zap.String("code", sg.Code))
	resp := sg.ToResponse()
	return &resp, nil
}
func (s *Service) GetSalaryGradeByID(ctx context.Context, id string) (*SalaryGradeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	sg, err := s.repo.FindSalaryGradeByID(ctx, uid)
	if err != nil { return nil, err }
	resp := sg.ToResponse()
	return &resp, nil
}
func (s *Service) ListSalaryGrades(ctx context.Context, page, perPage int) (*SalaryGradePaginatedResponse, error) {
	if page < 1 { page = defaultPage }
	if perPage < 1 { perPage = defaultPerPage }
	if perPage > maxPerPage { perPage = maxPerPage }
	list, total, err := s.repo.FindAllSalaryGrades(ctx, page, perPage)
	if err != nil { return nil, err }
	responses := make([]SalaryGradeResponse, len(list))
	for i, sg := range list { responses[i] = sg.ToResponse() }
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	return &SalaryGradePaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}
func (s *Service) UpdateSalaryGrade(ctx context.Context, id string, req UpdateSalaryGradeRequest) (*SalaryGradeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	sg, err := s.repo.FindSalaryGradeByID(ctx, uid)
	if err != nil { return nil, err }
	if req.Code != nil { sg.Code = *req.Code }
	if req.Name != nil { sg.Name = *req.Name }
	if req.Description != nil { sg.Description = *req.Description }
	if req.MinAmount != nil { sg.MinAmount = *req.MinAmount }
	if req.MaxAmount != nil { sg.MaxAmount = *req.MaxAmount }
	if req.SortOrder != nil { sg.SortOrder = *req.SortOrder }
	if err := s.repo.UpdateSalaryGrade(ctx, sg); err != nil { return nil, err }
	resp := sg.ToResponse()
	return &resp, nil
}
func (s *Service) DeleteSalaryGrade(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil { return fmt.Errorf("invalid id: %w", err) }
	return s.repo.DeleteSalaryGrade(ctx, uid)
}

// ── TER CRUD ──
func (s *Service) CreateTER(ctx context.Context, req CreateTERRequest) (*TERResponse, error) {
	t := &TER{Group: req.Group, BrutoMin: req.BrutoMin, BrutoMax: req.BrutoMax, Rate: req.Rate}
	if err := s.repo.CreateTER(ctx, t); err != nil { return nil, err }
	s.logger.Info("TER created", zap.String("id", t.ID.String()), zap.String("group", t.Group))
	resp := t.ToResponse()
	return &resp, nil
}
func (s *Service) GetTERByID(ctx context.Context, id string) (*TERResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	t, err := s.repo.FindTERByID(ctx, uid)
	if err != nil { return nil, err }
	resp := t.ToResponse()
	return &resp, nil
}
func (s *Service) ListTERs(ctx context.Context, page, perPage int) (*TERPaginatedResponse, error) {
	if page < 1 { page = defaultPage }
	if perPage < 1 { perPage = defaultPerPage }
	if perPage > maxPerPage { perPage = maxPerPage }
	list, total, err := s.repo.FindAllTERs(ctx, page, perPage)
	if err != nil { return nil, err }
	responses := make([]TERResponse, len(list))
	for i, t := range list { responses[i] = t.ToResponse() }
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	return &TERPaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}
func (s *Service) UpdateTER(ctx context.Context, id string, req UpdateTERRequest) (*TERResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	t, err := s.repo.FindTERByID(ctx, uid)
	if err != nil { return nil, err }
	if req.Group != nil { t.Group = *req.Group }
	if req.BrutoMin != nil { t.BrutoMin = req.BrutoMin }
	if req.BrutoMax != nil { t.BrutoMax = req.BrutoMax }
	if req.Rate != nil { t.Rate = *req.Rate }
	if err := s.repo.UpdateTER(ctx, t); err != nil { return nil, err }
	resp := t.ToResponse()
	return &resp, nil
}
func (s *Service) DeleteTER(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil { return fmt.Errorf("invalid id: %w", err) }
	return s.repo.DeleteTER(ctx, uid)
}

// ── PTKP CRUD ──
func (s *Service) CreatePTKP(ctx context.Context, req CreatePTKPRequest) (*PTKPResponse, error) {
	p := &PTKP{Name: req.Name, PTKP: req.PTKP, Group: req.Group}
	if err := s.repo.CreatePTKP(ctx, p); err != nil { return nil, err }
	s.logger.Info("PTKP created", zap.String("id", p.ID.String()), zap.String("name", p.Name))
	resp := p.ToResponse()
	return &resp, nil
}
func (s *Service) GetPTKPByID(ctx context.Context, id string) (*PTKPResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	p, err := s.repo.FindPTKPByID(ctx, uid)
	if err != nil { return nil, err }
	resp := p.ToResponse()
	return &resp, nil
}
func (s *Service) ListPTKPs(ctx context.Context, page, perPage int) (*PTKPPaginatedResponse, error) {
	if page < 1 { page = defaultPage }
	if perPage < 1 { perPage = defaultPerPage }
	if perPage > maxPerPage { perPage = maxPerPage }
	list, total, err := s.repo.FindAllPTKPs(ctx, page, perPage)
	if err != nil { return nil, err }
	responses := make([]PTKPResponse, len(list))
	for i, p := range list { responses[i] = p.ToResponse() }
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 { totalPages++ }
	return &PTKPPaginatedResponse{Success: true, Data: responses, Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}, nil
}
func (s *Service) UpdatePTKP(ctx context.Context, id string, req UpdatePTKPRequest) (*PTKPResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, fmt.Errorf("invalid id: %w", err) }
	p, err := s.repo.FindPTKPByID(ctx, uid)
	if err != nil { return nil, err }
	if req.Name != nil { p.Name = *req.Name }
	if req.PTKP != nil { p.PTKP = *req.PTKP }
	if req.Group != nil { p.Group = *req.Group }
	if err := s.repo.UpdatePTKP(ctx, p); err != nil { return nil, err }
	resp := p.ToResponse()
	return &resp, nil
}
func (s *Service) DeletePTKP(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil { return fmt.Errorf("invalid id: %w", err) }
	return s.repo.DeletePTKP(ctx, uid)
}
