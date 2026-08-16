package competency

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
)

// =========================================================================
// Rater Assignment (§9)
// =========================================================================

// AssignRaters menugaskan satu atau lebih rater ke sebuah assessment target
// (competency_event_target). Validasi: target ada, rater bukan subject
// sendiri kecuali rater_type=self, tidak ada duplikat rater pada target.
func (s *Service) AssignRaters(ctx context.Context, targetID string, req AssignRatersRequest) ([]RaterResponse, error) {
	targetUID, err := uuid.Parse(targetID)
	if err != nil {
		return nil, fmt.Errorf("invalid target id: %w", err)
	}
	target, err := s.repo.FindCompetencyEventTargetByID(ctx, targetUID)
	if err != nil {
		return nil, err
	}
	if target.EmployeeID == nil {
		return nil, fmt.Errorf("assessment target has no employee subject")
	}

	responses := make([]RaterResponse, 0, len(req.Raters))
	for _, r := range req.Raters {
		empUID, err := uuid.Parse(r.RaterEmployeeID)
		if err != nil {
			return nil, fmt.Errorf("invalid rater_employee_id: %w", err)
		}
		// Employee tidak boleh menilai dirinya sendiri kecuali rater_type=self.
		if empUID == *target.EmployeeID && r.RaterType != string(RaterTypeSelf) {
			return nil, fmt.Errorf("employee cannot rate themselves with rater_type %q", r.RaterType)
		}
		// Duplicate rater pada assessment yang sama tidak diperbolehkan.
		existing, err := s.repo.FindRaterByTargetAndEmployee(ctx, targetUID, empUID)
		if err == nil && existing != nil {
			return nil, fmt.Errorf("rater already assigned to this assessment")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		rat := &CompetencyAssessmentRater{
			CompetencyEventTargetID: targetUID,
			RaterEmployeeID:         empUID,
			RaterType:               r.RaterType,
			Status:                  string(RaterStatusAssigned),
		}
		if r.Weight != nil {
			rat.Weight = *r.Weight
		}
		now := time.Now()
		rat.AssignedAt = &now
		if err := s.repo.CreateRater(ctx, rat); err != nil {
			return nil, err
		}
		loaded, err := s.repo.FindRaterByID(ctx, rat.ID)
		if err != nil {
			return nil, err
		}
		responses = append(responses, loaded.ToResponse())
	}
	s.enrichRaterNames(ctx, responses)
	return responses, nil
}

// ListRatersByTarget mengambil seluruh rater sebuah assessment target.
// SuggestedRaters mengembalikan saran rater dari struktur organisasi untuk
// satu target: superior (atasan — parent org) dan subordinates (bawahan —
// subtree org). Rater yang sudah di-assign pada target dikecualikan, sehingga
// hasilnya aman langsung di-assign ulang (AssignRaters menolak duplikat).
func (s *Service) SuggestedRaters(ctx context.Context, targetID string) (*SuggestedRatersDTO, error) {
	uid, err := uuid.Parse(targetID)
	if err != nil {
		return nil, fmt.Errorf("invalid event target id: %w", err)
	}
	target, err := s.repo.FindCompetencyEventTargetByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	result := &SuggestedRatersDTO{Superior: []EmployeeBriefDTO{}, Subordinates: []EmployeeBriefDTO{}}
	if target.EmployeeID == nil {
		return result, nil // target level organisasi — tidak ada subject
	}
	subjectID := *target.EmployeeID

	// Rater yang sudah ada — jangan disarankan ulang.
	assigned := make(map[uuid.UUID]bool)
	existing, err := s.repo.FindRatersByTarget(ctx, uid)
	if err != nil {
		return nil, err
	}
	for _, rat := range existing {
		assigned[rat.RaterEmployeeID] = true
	}

	superiorIDs, err := s.repo.FindSuperiorEmployeeIDsBySubject(ctx, subjectID)
	if err != nil {
		return nil, err
	}
	subordinateIDs, err := s.repo.FindSubordinateEmployeeIDsByManager(ctx, subjectID)
	if err != nil {
		return nil, err
	}

	var suggestSuperior []uuid.UUID
	for _, id := range superiorIDs {
		if !assigned[id] {
			suggestSuperior = append(suggestSuperior, id)
		}
	}
	var suggestSubordinates []uuid.UUID
	for _, id := range subordinateIDs {
		if !assigned[id] {
			suggestSubordinates = append(suggestSubordinates, id)
		}
	}

	// Self: subject menilai dirinya sendiri — disarankan bila template
	// mewajibkan tipe self dan subject belum di-assign (mis. target lama).
	needSelf := false
	if !assigned[subjectID] && s.templateRequiresSelf(ctx, target) {
		needSelf = true
	}

	allIDs := append(suggestSuperior, suggestSubordinates...)
	if needSelf {
		allIDs = append(allIDs, subjectID)
	}
	names, err := s.repo.GetEmployeeNamesByIDs(ctx, allIDs)
	if err != nil {
		return nil, err
	}
	if needSelf {
		result.Self = &EmployeeBriefDTO{ID: subjectID.String(), Name: names[subjectID.String()]}
	}
	for _, id := range suggestSuperior {
		result.Superior = append(result.Superior, EmployeeBriefDTO{ID: id.String(), Name: names[id.String()]})
	}
	for _, id := range suggestSubordinates {
		result.Subordinates = append(result.Subordinates, EmployeeBriefDTO{ID: id.String(), Name: names[id.String()]})
	}
	return result, nil
}

// templateRequiresSelf memeriksa apakah template event mewajibkan tipe rater
// self (required atau min_rater > 0) — dasar saran rater self otomatis.
func (s *Service) templateRequiresSelf(ctx context.Context, target *CompetencyEventTarget) bool {
	event, err := s.repo.FindCompetencyEventByID(ctx, target.CompetencyEventID)
	if err != nil || event.TemplateID == nil {
		return false
	}
	raterTypes, err := s.repo.FindTemplateRaterTypes(ctx, *event.TemplateID)
	if err != nil {
		return false
	}
	for _, rt := range raterTypes {
		if rt.RaterType == string(RaterTypeSelf) && (rt.Required || rt.MinRater > 0) {
			return true
		}
	}
	return false
}

func (s *Service) ListRatersByTarget(ctx context.Context, targetID string) ([]RaterResponse, error) {
	targetUID, err := uuid.Parse(targetID)
	if err != nil {
		return nil, fmt.Errorf("invalid target id: %w", err)
	}
	raters, err := s.repo.FindRatersByTarget(ctx, targetUID)
	if err != nil {
		return nil, err
	}
	responses := make([]RaterResponse, 0, len(raters))
	for _, r := range raters {
		responses = append(responses, r.ToResponse())
	}
	s.enrichRaterNames(ctx, responses)
	return responses, nil
}

// DeleteRater menghapus assignment rater (bila belum submit).
func (s *Service) DeleteRater(ctx context.Context, raterID string) error {
	uid, err := uuid.Parse(raterID)
	if err != nil {
		return fmt.Errorf("invalid rater id: %w", err)
	}
	rat, err := s.repo.FindRaterByID(ctx, uid)
	if err != nil {
		return err
	}
	if rat.Status == string(RaterStatusSubmitted) {
		return fmt.Errorf("cannot delete a rater that already submitted")
	}
	return s.repo.DeleteRater(ctx, uid)
}

// MyAssessments mengambil seluruh assessment yang ditugaskan kepada employee
// dari user yang sedang login (inbox rater — "My Assessment").
func (s *Service) MyAssessments(ctx context.Context) ([]RaterResponse, error) {
	userID := authctx.GetUserID(ctx)
	if userID == nil {
		return nil, fmt.Errorf("authenticated user not found")
	}
	empID, err := s.repo.FindEmployeeIDByUserID(ctx, *userID)
	if err != nil {
		return nil, err
	}
	if empID == nil {
		return nil, fmt.Errorf("no employee account linked to this user")
	}
	raters, err := s.repo.FindRatersByEmployee(ctx, *empID)
	if err != nil {
		return nil, err
	}
	responses := make([]RaterResponse, 0, len(raters))
	for _, r := range raters {
		responses = append(responses, r.ToResponse())
	}
	s.enrichRaterNames(ctx, responses)
	return responses, nil
}

// ManagerAssessments mengembalikan daftar bawahan manager (dari struktur
// organisasi: seluruh employee di subtree organization tempat manager bekerja)
// beserta target assessment mereka — "Manager Assessment". Untuk tiap bawahan
// dicari target pada event (bila eventID diberikan) atau target terbaru
// (finalized dulu, lalu draft), beserta status rater superior manager pada
// target tsb (bila sudah di-assign). Bawahan tanpa target di-skip — belum
// menjadi subject assessment pada scope tsb.
func (s *Service) ManagerAssessments(ctx context.Context, eventID string) ([]ManagerAssessmentItem, error) {
	userID := authctx.GetUserID(ctx)
	if userID == nil {
		return nil, fmt.Errorf("authenticated user not found")
	}
	managerID, err := s.repo.FindEmployeeIDByUserID(ctx, *userID)
	if err != nil {
		return nil, err
	}
	if managerID == nil {
		return nil, fmt.Errorf("no employee account linked to this user")
	}

	subordinateIDs, err := s.repo.FindSubordinateEmployeeIDsByManager(ctx, *managerID)
	if err != nil {
		return nil, err
	}
	names, err := s.repo.GetEmployeeNamesByIDs(ctx, subordinateIDs)
	if err != nil {
		return nil, err
	}

	items := make([]ManagerAssessmentItem, 0, len(subordinateIDs))
	for _, subID := range subordinateIDs {
		// Pilih target: per event bila eventID diberikan, selain itu target
		// terbaru (finalized dulu, lalu draft) — pola sama GetEmployeeResult.
		var target *CompetencyEventTarget
		if eventID != "" {
			eventUID, perr := uuid.Parse(eventID)
			if perr != nil {
				return nil, fmt.Errorf("invalid event_id: %w", perr)
			}
			t, terr := s.repo.FindTargetByEventAndEmployee(ctx, eventUID, subID)
			if terr != nil {
				continue // bukan subject pada event tsb
			}
			target = t
		} else {
			targets, terr := s.repo.FindTargetsByEmployee(ctx, subID)
			if terr != nil || len(targets) == 0 {
				continue
			}
			for i := range targets {
				if targets[i].Status == "finalized" {
					target = &targets[i]
					break
				}
			}
			if target == nil {
				target = &targets[0]
			}
		}

		item := ManagerAssessmentItem{
			EmployeeID:        subID.String(),
			EmployeeName:      names[subID.String()],
			TargetID:          target.ID.String(),
			CompetencyEventID: target.CompetencyEventID.String(),
		}
		// Status rater superior manager pada target (bila sudah di-assign).
		rat, rerr := s.repo.FindRaterByTargetAndEmployee(ctx, target.ID, *managerID)
		if rerr == nil && rat != nil && rat.RaterType == string(RaterTypeSuperior) {
			item.RaterID = rat.ID.String()
			item.RaterStatus = rat.Status
			item.AssignedAt = rat.AssignedAt
			item.SubmittedAt = rat.SubmittedAt
		}
		items = append(items, item)
	}
	return items, nil
}

// GetAssessmentDetail mengambil detail satu assessment milik rater: rater,
// target subject, indicator template, dan response yang sudah tersimpan.
func (s *Service) GetAssessmentDetail(ctx context.Context, raterID string) (*AssessmentDetailDTO, error) {
	uid, err := uuid.Parse(raterID)
	if err != nil {
		return nil, fmt.Errorf("invalid rater id: %w", err)
	}
	rat, err := s.repo.FindRaterByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Target → competency event → template
	target := rat.Target
	if target == nil {
		target, err = s.repo.FindCompetencyEventTargetByID(ctx, rat.CompetencyEventTargetID)
		if err != nil {
			return nil, err
		}
	}
	event, err := s.repo.FindCompetencyEventByID(ctx, target.CompetencyEventID)
	if err != nil {
		return nil, err
	}

	var indicators []TemplateIndicatorResponse
	var tpl *CompetencyAssessmentTemplate
	if event.TemplateID != nil {
		tpl, err = s.repo.FindAssessmentTemplateByID(ctx, *event.TemplateID)
		if err == nil && tpl != nil {
			items, lerr := s.repo.ListTemplateIndicators(ctx, *event.TemplateID)
			if lerr != nil {
				return nil, lerr
			}
			for _, it := range items {
				indicators = append(indicators, it.ToResponse())
			}
		}
	}

	responses, err := s.repo.FindResponsesByRater(ctx, uid)
	if err != nil {
		return nil, err
	}
	respDTOs := make([]AssessmentResponseDTO, 0, len(responses))
	for _, r := range responses {
		respDTOs = append(respDTOs, r.ToDTO())
	}

	raterResp := rat.ToResponse()
	names, _ := s.repo.GetEmployeeNamesByIDs(ctx, []uuid.UUID{rat.RaterEmployeeID})
	if rat.Target != nil && rat.Target.EmployeeID != nil {
		subjectNames, _ := s.repo.GetEmployeeNamesByIDs(ctx, []uuid.UUID{*rat.Target.EmployeeID})
		if n, ok := subjectNames[rat.Target.EmployeeID.String()]; ok {
			raterResp.SubjectEmployeeName = n
		}
	}
	if n, ok := names[rat.RaterEmployeeID.String()]; ok {
		raterResp.RaterEmployeeName = n
	}

	targetResp := target.ToResponse()
	return &AssessmentDetailDTO{
		Rater:      raterResp,
		Target:     &targetResp,
		Indicators: indicators,
		Responses:  respDTOs,
	}, nil
}

// =========================================================================
// Assessment Response (§11)
// =========================================================================

// SaveResponses menyimpan/upsert response rater untuk assessment miliknya.
// Bila rater sudah submit, response tidak bisa diubah (immutable).
func (s *Service) SaveResponses(ctx context.Context, raterID string, req SaveResponsesRequest) ([]AssessmentResponseDTO, error) {
	uid, err := uuid.Parse(raterID)
	if err != nil {
		return nil, fmt.Errorf("invalid rater id: %w", err)
	}
	rat, err := s.repo.FindRaterByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if rat.Status == string(RaterStatusSubmitted) {
		return nil, fmt.Errorf("assessment already submitted and cannot be modified")
	}

	now := time.Now()
	for _, r := range req.Responses {
		indUID, err := uuid.Parse(r.IndicatorID)
		if err != nil {
			return nil, fmt.Errorf("invalid indicator_id: %w", err)
		}
		existing, err := s.repo.FindResponseByRaterAndIndicator(ctx, uid, indUID)
		if err == nil && existing != nil {
			existing.RatingValue = r.RatingValue
			existing.Comment = r.Comment
			existing.SubmittedAt = &now
			if err := s.repo.SaveResponse(ctx, existing); err != nil {
				return nil, err
			}
			continue
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		resp := &CompetencyAssessmentResponse{
			RaterID:     uid,
			IndicatorID: indUID,
			RatingValue: r.RatingValue,
			Comment:     r.Comment,
			SubmittedAt: &now,
		}
		if err := s.repo.SaveResponse(ctx, resp); err != nil {
			return nil, err
		}
	}

	// Mark rater as started (sudah ada aktivitas pengisian).
	if rat.Status == string(RaterStatusAssigned) {
		rat.Status = string(RaterStatusStarted)
		if err := s.repo.UpdateRater(ctx, rat); err != nil {
			return nil, err
		}
	}

	responses, err := s.repo.FindResponsesByRater(ctx, uid)
	if err != nil {
		return nil, err
	}
	out := make([]AssessmentResponseDTO, 0, len(responses))
	for _, r := range responses {
		out = append(out, r.ToDTO())
	}
	return out, nil
}

// SubmitAssessment meng-submit seluruh response rater — status rater menjadi
// submitted dan response terkunci (immutable).
func (s *Service) SubmitAssessment(ctx context.Context, raterID string) (*RaterResponse, error) {
	uid, err := uuid.Parse(raterID)
	if err != nil {
		return nil, fmt.Errorf("invalid rater id: %w", err)
	}
	rat, err := s.repo.FindRaterByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if rat.Status == string(RaterStatusSubmitted) {
		return nil, fmt.Errorf("assessment already submitted")
	}
	responses, err := s.repo.FindResponsesByRater(ctx, uid)
	if err != nil {
		return nil, err
	}
	if len(responses) == 0 {
		return nil, fmt.Errorf("no responses to submit")
	}
	now := time.Now()
	rat.Status = string(RaterStatusSubmitted)
	rat.SubmittedAt = &now
	if err := s.repo.UpdateRater(ctx, rat); err != nil {
		return nil, err
	}
	resp := rat.ToResponse()
	s.enrichRaterNames(ctx, []RaterResponse{resp})
	return &resp, nil
}

// enrichRaterNames mengisi nama rater (dan subject) pada response rater.
func (s *Service) enrichRaterNames(ctx context.Context, responses []RaterResponse) {
	empIDSet := make(map[uuid.UUID]struct{})
	for i := range responses {
		if uid, err := uuid.Parse(responses[i].RaterEmployeeID); err == nil {
			empIDSet[uid] = struct{}{}
		}
		if uid, err := uuid.Parse(responses[i].SubjectEmployeeID); err == nil && uid != uuid.Nil {
			empIDSet[uid] = struct{}{}
		}
	}
	ids := make([]uuid.UUID, 0, len(empIDSet))
	for id := range empIDSet {
		ids = append(ids, id)
	}
	names, err := s.repo.GetEmployeeNamesByIDs(ctx, ids)
	if err != nil {
		return
	}
	for i := range responses {
		if n, ok := names[responses[i].RaterEmployeeID]; ok {
			responses[i].RaterEmployeeName = n
		}
		if n, ok := names[responses[i].SubjectEmployeeID]; ok {
			responses[i].SubjectEmployeeName = n
		}
	}
}
