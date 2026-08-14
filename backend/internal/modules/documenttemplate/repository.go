package documenttemplate

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
)

type Repository struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
}

func NewRepository(dbResolver func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbResolver: dbResolver}
}

func NewTenantDBResolver(dbManager *database.Manager) func(ctx context.Context) (*gorm.DB, error) {
	return func(ctx context.Context) (*gorm.DB, error) {
		companyID, ok := ctx.Value("company_id").(string)
		if !ok || companyID == "" {
			return nil, fmt.Errorf("tenant context not found in request: company_id is required")
		}
		return dbManager.TenantDB(companyID)
	}
}

func (r *Repository) getDB(ctx context.Context) (*gorm.DB, error) {
	return r.dbResolver(ctx)
}

func escapeLike(s string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(s)
}

func (r *Repository) List(ctx context.Context, page, perPage int, documentType, status, search string) ([]DocumentTemplate, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	query := db.Session(&gorm.Session{}).Model(&DocumentTemplate{}).Where("deleted_at IS NULL")
	if documentType != "" {
		query = query.Where("type = ?", documentType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		like := "%" + escapeLike(search) + "%"
		query = query.Where("(LOWER(name) LIKE LOWER(?) OR LOWER(code) LIKE LOWER(?))", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count document templates: %w", err)
	}

	var items []DocumentTemplate
	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list document templates: %w", err)
	}
	return items, total, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*DocumentTemplate, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var tpl DocumentTemplate
	err = db.Where("id = ? AND deleted_at IS NULL", id).First(&tpl).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTemplateNotFound
		}
		return nil, fmt.Errorf("failed to get document template: %w", err)
	}
	return &tpl, nil
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*DocumentTemplate, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var tpl DocumentTemplate
	err = db.Where("code = ? AND deleted_at IS NULL", code).First(&tpl).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTemplateNotFound
		}
		return nil, fmt.Errorf("failed to get document template by code: %w", err)
	}
	return &tpl, nil
}

func (r *Repository) FindActiveByType(ctx context.Context, documentType string) (*DocumentTemplate, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var tpl DocumentTemplate
	err = db.Where("type = ? AND status = ? AND deleted_at IS NULL", documentType, StatusActive).First(&tpl).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTemplateNotFound
		}
		return nil, fmt.Errorf("failed to find active document template: %w", err)
	}
	return &tpl, nil
}

func (r *Repository) Create(ctx context.Context, tpl *DocumentTemplate) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(tpl).Error; err != nil {
		return fmt.Errorf("failed to create document template: %w", err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, tpl *DocumentTemplate) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Save(tpl).Error; err != nil {
		return fmt.Errorf("failed to update document template: %w", err)
	}
	return nil
}

func (r *Repository) SoftDelete(ctx context.Context, id string) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	// Also mutate code so the unconditional unique index on `code` doesn't
	// block reusing the code for a new template after this row is soft-deleted.
	// CONCAT() isn't portable across sqlite (tests), mysql, and postgres, so
	// the suffix is computed in Go and written as a plain value.
	var tpl DocumentTemplate
	if err := db.Select("code").Where("id = ?", id).First(&tpl).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrTemplateNotFound
		}
		return fmt.Errorf("failed to load document template for soft delete: %w", err)
	}
	deletedCode := fmt.Sprintf("%s_deleted_%s", tpl.Code, id)
	if err := db.Model(&DocumentTemplate{}).Where("id = ?", id).Updates(map[string]interface{}{
		"deleted_at": gorm.Expr("CURRENT_TIMESTAMP"),
		"code":       deletedCode,
	}).Error; err != nil {
		return fmt.Errorf("failed to soft delete document template: %w", err)
	}
	return nil
}

func (r *Repository) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(fn)
}

func (r *Repository) CreateVersion(ctx context.Context, tx *gorm.DB, v *DocumentTemplateVersion) error {
	if err := tx.Create(v).Error; err != nil {
		return fmt.Errorf("failed to create document template version: %w", err)
	}
	return nil
}

func (r *Repository) ListVersions(ctx context.Context, templateID string) ([]DocumentTemplateVersion, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var versions []DocumentTemplateVersion
	if err := db.Where("template_id = ?", templateID).Order("version DESC").Find(&versions).Error; err != nil {
		return nil, fmt.Errorf("failed to list document template versions: %w", err)
	}
	return versions, nil
}

func (r *Repository) GetVersion(ctx context.Context, templateID, versionID string) (*DocumentTemplateVersion, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var v DocumentTemplateVersion
	err = db.Where("id = ? AND template_id = ?", versionID, templateID).First(&v).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrVersionNotFound
		}
		return nil, fmt.Errorf("failed to get document template version: %w", err)
	}
	return &v, nil
}

func (r *Repository) NextVersionNumber(ctx context.Context, tx *gorm.DB, templateID string) (int, error) {
	var max int
	err := tx.Model(&DocumentTemplateVersion{}).
		Where("template_id = ?", templateID).
		Select("COALESCE(MAX(version), 0)").
		Scan(&max).Error
	if err != nil {
		return 0, fmt.Errorf("failed to compute next version number: %w", err)
	}
	return max + 1, nil
}

func (r *Repository) CreateAudit(ctx context.Context, tx *gorm.DB, a *DocumentTemplateAudit) error {
	if tx != nil {
		if err := tx.Create(a).Error; err != nil {
			return fmt.Errorf("failed to create document template audit: %w", err)
		}
		return nil
	}
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(a).Error; err != nil {
		return fmt.Errorf("failed to create document template audit: %w", err)
	}
	return nil
}
