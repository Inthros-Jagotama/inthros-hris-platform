package documenttemplate

type CreateTemplateRequest struct {
	Name         string `json:"name" binding:"required,max=255"`
	Code         string `json:"code" binding:"required,max=100"`
	DocumentType string `json:"document_type" binding:"required"`
	Description  string `json:"description,omitempty" binding:"max=1000"`
}

type CreateFromDefaultRequest struct {
	DocumentType string `json:"document_type" binding:"required"`
	Name         string `json:"name" binding:"required,max=255"`
	Code         string `json:"code" binding:"required,max=100"`
}

type UpdateTemplateRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,max=255"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=1000"`
}

type UpdateDefaultContentRequest struct {
	Content string `json:"content" binding:"required"`
}

type CreateVersionRequest struct {
	Content      string `json:"content" binding:"required"`
	PaperSize    string `json:"paper_size,omitempty" binding:"omitempty,oneof=A4 A5 Letter Legal"`
	Orientation  string `json:"orientation,omitempty" binding:"omitempty,oneof=portrait landscape"`
	MarginTop    int    `json:"margin_top,omitempty"`
	MarginRight  int    `json:"margin_right,omitempty"`
	MarginBottom int    `json:"margin_bottom,omitempty"`
	MarginLeft   int    `json:"margin_left,omitempty"`
}

type TemplateListResponse struct {
	Data  []DocumentTemplate `json:"data"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
}
