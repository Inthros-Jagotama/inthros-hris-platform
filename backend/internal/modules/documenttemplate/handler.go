package documenttemplate

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func actorID(c *gin.Context) string {
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (h *Handler) handleServiceError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *InvalidDocumentTypeError:
		httputil.ErrorSimple(c, http.StatusBadRequest, e.Error())
	case *DuplicateCodeError:
		httputil.ErrorJSON(c, http.StatusConflict, "DUPLICATE_CODE", "documenttemplate.duplicate_code", e.Code)
	case *DuplicateActiveTemplateError:
		httputil.ErrorJSON(c, http.StatusConflict, "DUPLICATE_ACTIVE", "documenttemplate.duplicate_active", e.DocumentType)
	case *ReferenceTemplateImmutableError:
		httputil.ErrorSimple(c, http.StatusForbidden, e.Error())
	default:
		if err == ErrTemplateNotFound || err == ErrVersionNotFound {
			httputil.NotFound(c, err.Error())
			return
		}
		httputil.InternalError(c, err.Error())
	}
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	items, total, err := h.svc.List(c.Request.Context(), page, perPage, c.Query("document_type"), c.Query("status"), c.Query("search"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, TemplateListResponse{Data: items, Total: total, Page: page})
}

func (h *Handler) ListVariables(c *gin.Context) {
	httputil.SuccessJSON(c, VariableRegistry())
}

func (h *Handler) GetByID(c *gin.Context) {
	tpl, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, tpl)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateTemplateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	tpl, err := h.svc.Create(c.Request.Context(), req.Name, req.Code, req.DocumentType, req.Description, actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.CreatedJSON(c, tpl, "documenttemplate.created")
}

func (h *Handler) CreateFromDefault(c *gin.Context) {
	var req CreateFromDefaultRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	tpl, err := h.svc.CreateFromDefault(c.Request.Context(), req.DocumentType, req.Name, req.Code, actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.CreatedJSON(c, tpl, "documenttemplate.created_from_default")
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateTemplateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	name, desc := "", ""
	if req.Name != nil {
		name = *req.Name
	}
	if req.Description != nil {
		desc = *req.Description
	}
	tpl, err := h.svc.Update(c.Request.Context(), c.Param("id"), name, desc, actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.UpdatedJSON(c, tpl, "documenttemplate.updated")
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id"), actorID(c)); err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.DeletedJSON(c, "documenttemplate.deleted")
}

func (h *Handler) Activate(c *gin.Context) {
	tpl, err := h.svc.Activate(c.Request.Context(), c.Param("id"), actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.UpdatedJSON(c, tpl, "documenttemplate.activated")
}

func (h *Handler) Deactivate(c *gin.Context) {
	tpl, err := h.svc.Deactivate(c.Request.Context(), c.Param("id"), actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.UpdatedJSON(c, tpl, "documenttemplate.deactivated")
}

func (h *Handler) ListVersions(c *gin.Context) {
	versions, err := h.svc.ListVersions(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, versions)
}

func (h *Handler) CreateVersion(c *gin.Context) {
	var req CreateVersionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	paperSize, orientation := req.PaperSize, req.Orientation
	if paperSize == "" {
		paperSize = "A4"
	}
	if orientation == "" {
		orientation = "portrait"
	}
	margins := [4]int{req.MarginTop, req.MarginRight, req.MarginBottom, req.MarginLeft}
	for i, m := range margins {
		if m == 0 {
			margins[i] = 20
		}
	}
	v, err := h.svc.CreateVersion(c.Request.Context(), c.Param("id"), req.Content, paperSize, orientation, margins, actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.CreatedJSON(c, v, "documenttemplate.version_created")
}
