package recruitment

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/modules/approval"
	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// =========================================================================
// Job Requisitions
// =========================================================================

func (h *Handler) CreateRequisition(c *gin.Context) {
	var req CreateRequisitionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateRequisition(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListRequisitions(c *gin.Context) {
	page, perPage := parsePagination(c)
	orgID := c.Query("organization_id")
	status := c.Query("status")
	var orgPtr *string
	if orgID != "" {
		orgPtr = &orgID
	}
	resp, err := h.svc.ListRequisitions(c.Request.Context(), orgPtr, &status, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetRequisitionByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetRequisitionByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateRequisition(c *gin.Context) {
	id := c.Param("id")
	var req UpdateRequisitionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateRequisition(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteRequisition(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteRequisition(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// SubmitRequisition menangani POST /api/v1/tenant/recruitment/requisitions/:id/submit
// Routes the requisition through the central approval module (plan G-1).
func (h *Handler) SubmitRequisition(c *gin.Context) {
	id := c.Param("id")
	var req SubmitRequisitionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	flowID := ""
	if req.FlowID != nil {
		flowID = *req.FlowID
	}
	resp, err := h.svc.SubmitRequisition(c.Request.Context(), id, flowID)
	if err != nil {
		// Approval routing/assignee failures (approval.RoutingError) get a
		// bilingual 400 so the user sees why their submission didn't reach
		// an approver instead of a raw internal error.
		if approval.EmitRoutingError(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =========================================================================
// Job Offers (G-3)
// =========================================================================

func (h *Handler) CreateOffer(c *gin.Context) {
	var req CreateOfferRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateOffer(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListOffers(c *gin.Context) {
	page, perPage := parsePagination(c)
	applicationID := c.Query("application_id")
	status := c.Query("status")
	var appPtr *string
	if applicationID != "" {
		appPtr = &applicationID
	}
	resp, err := h.svc.ListOffers(c.Request.Context(), appPtr, &status, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetOfferByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetOfferByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateOffer(c *gin.Context) {
	id := c.Param("id")
	var req UpdateOfferRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateOffer(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteOffer(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteOffer(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// SubmitOffer menangani POST /api/v1/tenant/recruitment/offers/:id/submit
// Routes the offer through the central approval module (plan G-3).
func (h *Handler) SubmitOffer(c *gin.Context) {
	id := c.Param("id")
	var req SubmitOfferRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	flowID := ""
	if req.FlowID != nil {
		flowID = *req.FlowID
	}
	resp, err := h.svc.SubmitOffer(c.Request.Context(), id, flowID)
	if err != nil {
		// Approval routing/assignee failures (approval.RoutingError) get a
		// bilingual 400 so the user sees why their submission didn't reach
		// an approver instead of a raw internal error.
		if approval.EmitRoutingError(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) SendOffer(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.SendOffer(c.Request.Context(), id)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) AcceptOffer(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.AcceptOffer(c.Request.Context(), id)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) RejectOffer(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.RejectOffer(c.Request.Context(), id)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) WithdrawOffer(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.WithdrawOffer(c.Request.Context(), id)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =========================================================================
// Candidates
// =========================================================================

func (h *Handler) CreateCandidate(c *gin.Context) {
	var req CreateCandidateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCandidate(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCandidates(c *gin.Context) {
	page, perPage := parsePagination(c)
	search := c.Query("search")
	var searchPtr *string
	if search != "" {
		searchPtr = &search
	}
	resp, err := h.svc.ListCandidates(c.Request.Context(), searchPtr, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetCandidateByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetCandidateByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateCandidate(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCandidateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCandidate(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteCandidate(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCandidate(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) CreateCandidateEducation(c *gin.Context) {
	candidateID := c.Param("id")
	var req CreateCandidateEducationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCandidateEducation(c.Request.Context(), candidateID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCandidateEducations(c *gin.Context) {
	candidateID := c.Param("id")
	resp, err := h.svc.ListCandidateEducations(c.Request.Context(), candidateID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateCandidateEducation(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCandidateEducationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCandidateEducation(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteCandidateEducation(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCandidateEducation(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) CreateCandidateWorkExperience(c *gin.Context) {
	candidateID := c.Param("id")
	var req CreateCandidateWorkExperienceRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCandidateWorkExperience(c.Request.Context(), candidateID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCandidateWorkExperiences(c *gin.Context) {
	candidateID := c.Param("id")
	resp, err := h.svc.ListCandidateWorkExperiences(c.Request.Context(), candidateID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateCandidateWorkExperience(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCandidateWorkExperienceRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCandidateWorkExperience(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteCandidateWorkExperience(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCandidateWorkExperience(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) CreateCandidateSkill(c *gin.Context) {
	candidateID := c.Param("id")
	var req CreateCandidateSkillRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCandidateSkill(c.Request.Context(), candidateID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCandidateSkills(c *gin.Context) {
	candidateID := c.Param("id")
	resp, err := h.svc.ListCandidateSkills(c.Request.Context(), candidateID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateCandidateSkill(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCandidateSkillRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCandidateSkill(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteCandidateSkill(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCandidateSkill(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) CreateCandidateCertification(c *gin.Context) {
	candidateID := c.Param("id")
	var req CreateCandidateCertificationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCandidateCertification(c.Request.Context(), candidateID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCandidateCertifications(c *gin.Context) {
	candidateID := c.Param("id")
	resp, err := h.svc.ListCandidateCertifications(c.Request.Context(), candidateID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateCandidateCertification(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCandidateCertificationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCandidateCertification(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteCandidateCertification(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCandidateCertification(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) CreateCandidateDocument(c *gin.Context) {
	candidateID := c.Param("id")
	var req CreateCandidateDocumentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCandidateDocument(c.Request.Context(), candidateID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCandidateDocuments(c *gin.Context) {
	candidateID := c.Param("id")
	resp, err := h.svc.ListCandidateDocuments(c.Request.Context(), candidateID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateCandidateDocument(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCandidateDocumentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCandidateDocument(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteCandidateDocument(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCandidateDocument(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) CreateApplicationScreening(c *gin.Context) {
	applicationID := c.Param("id")
	var req CreateApplicationScreeningRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateApplicationScreening(c.Request.Context(), applicationID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListApplicationScreenings(c *gin.Context) {
	applicationID := c.Param("id")
	resp, err := h.svc.ListApplicationScreenings(c.Request.Context(), applicationID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateApplicationScreening(c *gin.Context) {
	id := c.Param("id")
	var req UpdateApplicationScreeningRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateApplicationScreening(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteApplicationScreening(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteApplicationScreening(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) CreateAssessment(c *gin.Context) {
	var req CreateAssessmentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateAssessment(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListAssessments(c *gin.Context) {
	resp, err := h.svc.ListAssessments(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) GetAssessmentByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetAssessmentByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateAssessment(c *gin.Context) {
	id := c.Param("id")
	var req UpdateAssessmentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateAssessment(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteAssessment(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteAssessment(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) AddAssessmentParticipant(c *gin.Context) {
	assessmentID := c.Param("id")
	var req AddAssessmentParticipantRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.AddAssessmentParticipant(c.Request.Context(), assessmentID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListAssessmentParticipants(c *gin.Context) {
	assessmentID := c.Param("id")
	resp, err := h.svc.ListAssessmentParticipants(c.Request.Context(), assessmentID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateAssessmentParticipant(c *gin.Context) {
	id := c.Param("id")
	var req UpdateAssessmentParticipantRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateAssessmentParticipant(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteAssessmentParticipant(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteAssessmentParticipant(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) AddInterviewer(c *gin.Context) {
	interviewID := c.Param("id")
	var req AddInterviewerRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.AddInterviewer(c.Request.Context(), interviewID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListInterviewers(c *gin.Context) {
	interviewID := c.Param("id")
	resp, err := h.svc.ListInterviewers(c.Request.Context(), interviewID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) RemoveInterviewer(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.RemoveInterviewer(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) AddScorecardItem(c *gin.Context) {
	interviewID := c.Param("id")
	var req AddScorecardItemRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.AddScorecardItem(c.Request.Context(), interviewID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListScorecardItems(c *gin.Context) {
	interviewID := c.Param("id")
	resp, err := h.svc.ListScorecardItems(c.Request.Context(), interviewID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateScorecardItem(c *gin.Context) {
	id := c.Param("id")
	var req UpdateScorecardItemRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateScorecardItem(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteScorecardItem(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteScorecardItem(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) CompleteInterview(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.CompleteInterview(c.Request.Context(), id)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) CreateRequisitionRequirement(c *gin.Context) {
	requisitionID := c.Param("id")
	var req CreateRequisitionRequirementRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateRequisitionRequirement(c.Request.Context(), requisitionID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListRequisitionRequirements(c *gin.Context) {
	requisitionID := c.Param("id")
	resp, err := h.svc.ListRequisitionRequirements(c.Request.Context(), requisitionID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateRequisitionRequirement(c *gin.Context) {
	id := c.Param("id")
	var req UpdateRequisitionRequirementRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateRequisitionRequirement(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteRequisitionRequirement(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteRequisitionRequirement(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) CreateRequisitionCompetency(c *gin.Context) {
	requisitionID := c.Param("id")
	var req CreateRequisitionCompetencyRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateRequisitionCompetency(c.Request.Context(), requisitionID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListRequisitionCompetencies(c *gin.Context) {
	requisitionID := c.Param("id")
	resp, err := h.svc.ListRequisitionCompetencies(c.Request.Context(), requisitionID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateRequisitionCompetency(c *gin.Context) {
	id := c.Param("id")
	var req UpdateRequisitionCompetencyRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateRequisitionCompetency(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteRequisitionCompetency(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteRequisitionCompetency(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) GetCandidateMatchScore(c *gin.Context) {
	applicationID := c.Param("id")
	resp, err := h.svc.GetCandidateMatchScore(c.Request.Context(), applicationID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetEligibleInternalCandidates menangani GET /recruitment/eligible-internal-candidates
// (plan S-4) — membaca employee internal yang eligible untuk target position
// dari Career Intelligence (CI menentukan eligibility; Recruitment hanya
// mengeksekusi aplikasi internal).
func (h *Handler) GetEligibleInternalCandidates(c *gin.Context) {
	positionID := c.Query("position_id")
	if positionID == "" {
		httputil.BadRequest(c, "position_id is required")
		return
	}
	result, err := h.svc.GetEligibleInternalCandidates(c.Request.Context(), positionID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, result)
}

// =========================================================================
// Job Applications
// =========================================================================

func (h *Handler) CreateApplication(c *gin.Context) {
	var req CreateApplicationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateApplication(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListApplications(c *gin.Context) {
	page, perPage := parsePagination(c)
	reqID := c.Query("requisition_id")
	candID := c.Query("candidate_id")
	status := c.Query("status")
	var reqPtr, candPtr *string
	if reqID != "" {
		reqPtr = &reqID
	}
	if candID != "" {
		candPtr = &candID
	}
	resp, err := h.svc.ListApplications(c.Request.Context(), reqPtr, candPtr, &status, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetApplicationByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetApplicationByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateApplicationStatus(c *gin.Context) {
	id := c.Param("id")
	var req UpdateApplicationStatusRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	var changedBy *uuid.UUID
	if userIDStr := c.GetString("user_id"); userIDStr != "" {
		if uid, err := uuid.Parse(userIDStr); err == nil {
			changedBy = &uid
		}
	}
	resp, err := h.svc.UpdateApplicationStatus(c.Request.Context(), id, req.Status, req.RejectionReason, req.Notes, changedBy)
	if err != nil {
		if errors.Is(err, ErrInvalidStatusTransition) {
			httputil.BadRequest(c, err.Error())
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) CreateCandidateConsent(c *gin.Context) {
	candidateID := c.Param("id")
	var req CreateCandidateConsentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	var changedBy *uuid.UUID
	if userIDStr := c.GetString("user_id"); userIDStr != "" {
		if uid, err := uuid.Parse(userIDStr); err == nil {
			changedBy = &uid
		}
	}
	resp, err := h.svc.CreateCandidateConsent(c.Request.Context(), candidateID, req, changedBy)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCandidateConsents(c *gin.Context) {
	candidateID := c.Param("id")
	resp, err := h.svc.ListCandidateConsents(c.Request.Context(), candidateID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) GetApplicationHistory(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetApplicationHistory(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteApplication(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteApplication(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Interviews
// =========================================================================

func (h *Handler) CreateInterview(c *gin.Context) {
	var req CreateInterviewRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateInterview(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListInterviews(c *gin.Context) {
	page, perPage := parsePagination(c)
	appID := c.Query("application_id")
	intvID := c.Query("interviewer_id")
	var appPtr, intvPtr *string
	if appID != "" {
		appPtr = &appID
	}
	if intvID != "" {
		intvPtr = &intvID
	}
	resp, err := h.svc.ListInterviews(c.Request.Context(), appPtr, intvPtr, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetInterviewByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetInterviewByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateInterview(c *gin.Context) {
	id := c.Param("id")
	var req UpdateInterviewRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateInterview(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteInterview(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteInterview(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Onboarding Task Templates
// =========================================================================

func (h *Handler) CreateOnboardingTaskTemplate(c *gin.Context) {
	var req CreateOnboardingTaskTemplateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateOnboardingTaskTemplate(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListOnboardingTaskTemplates(c *gin.Context) {
	page, perPage := parsePagination(c)
	category := c.Query("category")
	var catPtr *string
	if category != "" {
		catPtr = &category
	}
	resp, err := h.svc.ListOnboardingTaskTemplates(c.Request.Context(), catPtr, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateOnboardingTaskTemplate(c *gin.Context) {
	id := c.Param("id")
	var req UpdateOnboardingTaskTemplateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateOnboardingTaskTemplate(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteOnboardingTaskTemplate(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteOnboardingTaskTemplate(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Employee Onboarding
// =========================================================================

func (h *Handler) CreateEmployeeOnboarding(c *gin.Context) {
	var req CreateEmployeeOnboardingRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateEmployeeOnboarding(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListEmployeeOnboardings(c *gin.Context) {
	page, perPage := parsePagination(c)
	status := c.Query("status")
	resp, err := h.svc.ListEmployeeOnboardings(c.Request.Context(), &status, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetEmployeeOnboardingByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetEmployeeOnboardingByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateEmployeeOnboarding(c *gin.Context) {
	id := c.Param("id")
	var req UpdateEmployeeOnboardingRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateEmployeeOnboarding(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteEmployeeOnboarding(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteEmployeeOnboarding(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Onboarding Task Items
// =========================================================================

func (h *Handler) CreateOnboardingTaskItem(c *gin.Context) {
	var req CreateOnboardingTaskItemRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateOnboardingTaskItem(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListOnboardingTaskItems(c *gin.Context) {
	onboardingID := c.Param("id")
	items, err := h.svc.ListOnboardingTaskItems(c.Request.Context(), onboardingID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, items)
}

func (h *Handler) UpdateOnboardingTaskItem(c *gin.Context) {
	id := c.Param("id")
	var req UpdateOnboardingTaskItemRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateOnboardingTaskItem(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteOnboardingTaskItem(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteOnboardingTaskItem(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Helpers
// =========================================================================

func parsePagination(c *gin.Context) (int, int) {
	page := 1
	perPage := 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if pp := c.Query("per_page"); pp != "" {
		if v, err := strconv.Atoi(pp); err == nil && v > 0 && v <= 100 {
			perPage = v
		}
	}
	return page, perPage
}
