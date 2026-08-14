package documenttemplate

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func actorIDPtr(actorID string) *string {
	if actorID == "" {
		return nil
	}
	return &actorID
}

// generateTemplateCode membuat kode template otomatis saat client tidak
// mengirimnya (spec: "kode dibuat otomatis dan tidak perlu ditampilkan di
// form"). Format mengikuti konvensi project: prefix jenis + UUID pendek
// uppercase, mis. TMPL-CONTRACT_AGREEMENT-1A2B3C4D.
func generateTemplateCode(documentType string) string {
	return fmt.Sprintf("TMPL-%s-%s", documentType, strings.ToUpper(uuid.New().String()[:8]))
}

// versionFileURL mengembalikan URL publik untuk konten versi yang berupa path
// file DOCX tersimpan (mis. "document_templates/{uuid}.docx" →
// "/uploads/document_templates/{uuid}.docx", diserve router.Static("/uploads")).
// Konten lama berupa HTML (mis. default content) tidak menghasilkan URL file.
func versionFileURL(content string) string {
	if strings.HasPrefix(content, "document_templates/") {
		return "/uploads/" + content
	}
	return ""
}

// decorateVersionFileURL mengisi FileURL (field response-only, gorm:"-") agar
// frontend bisa menampilkan/download file DOCX tanpa menebak dari content.
func decorateVersionFileURL(v *DocumentTemplateVersion) {
	if v == nil {
		return
	}
	v.FileURL = versionFileURL(v.Content)
}

type Service struct {
	repo   *Repository
	logger *zap.Logger
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) List(ctx context.Context, page, perPage int, documentType, status, search string) ([]DocumentTemplate, int64, error) {
	return s.repo.List(ctx, page, perPage, documentType, status, search)
}

func (s *Service) GetByID(ctx context.Context, id string) (*DocumentTemplate, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) checkCodeUnique(ctx context.Context, code string) error {
	_, err := s.repo.GetByCode(ctx, code)
	if err == nil {
		return &DuplicateCodeError{Code: code}
	}
	if errors.Is(err, ErrTemplateNotFound) {
		return nil
	}
	return err
}

func (s *Service) Create(ctx context.Context, name, code, documentType, description string, actorID string) (*DocumentTemplate, error) {
	if !ValidDocumentTypes[documentType] {
		return nil, &InvalidDocumentTypeError{DocumentType: documentType}
	}
	if code == "" {
		code = generateTemplateCode(documentType)
	}
	if err := s.checkCodeUnique(ctx, code); err != nil {
		return nil, err
	}

	tpl := &DocumentTemplate{
		ID:           uuid.New().String(),
		Name:         name,
		Code:         code,
		DocumentType: documentType,
		Status:       StatusInactive,
		IsActive:     true,
	}
	if description != "" {
		tpl.Description = &description
	}
	if err := s.repo.Create(ctx, tpl); err != nil {
		return nil, err
	}
	if err := s.repo.CreateAudit(ctx, nil, &DocumentTemplateAudit{
		ID:         uuid.New().String(),
		TemplateID: tpl.ID,
		Action:     "CREATED",
		ActorID:    actorIDPtr(actorID),
	}); err != nil {
		return nil, err
	}
	return tpl, nil
}

func (s *Service) Update(ctx context.Context, id string, name, description *string, actorID string) (*DocumentTemplate, error) {
	tpl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if name != nil {
		tpl.Name = *name
	}
	if description != nil {
		tpl.Description = description
	}
	if err := s.repo.Update(ctx, tpl); err != nil {
		return nil, err
	}
	if err := s.repo.CreateAudit(ctx, nil, &DocumentTemplateAudit{
		ID:         uuid.New().String(),
		TemplateID: tpl.ID,
		Action:     "UPDATED",
		ActorID:    actorIDPtr(actorID),
	}); err != nil {
		return nil, err
	}
	return tpl, nil
}

func (s *Service) Delete(ctx context.Context, id, actorID string) error {
	tpl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	// Spec §2.1 restricts deletion to templates "not yet used"; there is no
	// generated_documents writer yet in this phase, so every template is
	// deletable for now. Add a usage check here once Phase 5 introduces
	// document generation.
	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return err
	}
	return s.repo.CreateAudit(ctx, nil, &DocumentTemplateAudit{
		ID:         uuid.New().String(),
		TemplateID: tpl.ID,
		Action:     "DELETED",
		ActorID:    actorIDPtr(actorID),
	})
}

func (s *Service) Activate(ctx context.Context, id, actorID string) (*DocumentTemplate, error) {
	var target DocumentTemplate
	err := s.repo.WithTx(ctx, func(tx *gorm.DB) error {
		// Locking read: serializes concurrent Activate calls for the same
		// document type against each other so the "one ACTIVE per type"
		// invariant holds on MySQL, which has no partial unique index backstop.
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND deleted_at IS NULL", id).First(&target).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTemplateNotFound
			}
			return fmt.Errorf("failed to load template for activation: %w", err)
		}

		// Lock all other rows of this document type so a concurrent Activate
		// on a sibling template of the same type serializes behind this tx.
		var others []DocumentTemplate
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("type = ? AND id <> ? AND deleted_at IS NULL", target.DocumentType, target.ID).
			Find(&others).Error; err != nil {
			return fmt.Errorf("failed to lock sibling templates: %w", err)
		}

		// Set-based deactivation of ALL other active rows of this type (self
		// heals if the invariant was ever violated), not just one.
		if err := tx.Model(&DocumentTemplate{}).
			Where("type = ? AND status = ? AND id <> ? AND deleted_at IS NULL", target.DocumentType, StatusActive, target.ID).
			Update("status", StatusInactive).Error; err != nil {
			return fmt.Errorf("failed to deactivate previous active templates: %w", err)
		}
		for _, previous := range others {
			if previous.Status != StatusActive {
				continue
			}
			if err := s.repo.CreateAudit(ctx, tx, &DocumentTemplateAudit{
				ID:         uuid.New().String(),
				TemplateID: previous.ID,
				Action:     "DEACTIVATED",
				ActorID:    actorIDPtr(actorID),
			}); err != nil {
				return err
			}
		}

		if err := tx.Model(&DocumentTemplate{}).Where("id = ?", target.ID).Update("status", StatusActive).Error; err != nil {
			return fmt.Errorf("failed to activate template: %w", err)
		}
		target.Status = StatusActive
		return s.repo.CreateAudit(ctx, tx, &DocumentTemplateAudit{
			ID:         uuid.New().String(),
			TemplateID: target.ID,
			Action:     "ACTIVATED",
			ActorID:    actorIDPtr(actorID),
		})
	})
	if err != nil {
		return nil, err
	}
	return &target, nil
}

func (s *Service) Deactivate(ctx context.Context, id, actorID string) (*DocumentTemplate, error) {
	tpl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	tpl.Status = StatusInactive
	if err := s.repo.Update(ctx, tpl); err != nil {
		return nil, err
	}
	if err := s.repo.CreateAudit(ctx, nil, &DocumentTemplateAudit{
		ID:         uuid.New().String(),
		TemplateID: tpl.ID,
		Action:     "DEACTIVATED",
		ActorID:    actorIDPtr(actorID),
	}); err != nil {
		return nil, err
	}
	return tpl, nil
}

func (s *Service) CreateVersion(ctx context.Context, templateID, content, paperSize, orientation string, margins [4]int, fileName, actorID string) (*DocumentTemplateVersion, error) {
	var v DocumentTemplateVersion
	err := s.repo.WithTx(ctx, func(tx *gorm.DB) error {
		var tpl DocumentTemplate
		if err := tx.Where("id = ? AND deleted_at IS NULL", templateID).First(&tpl).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTemplateNotFound
			}
			return fmt.Errorf("failed to load template for version creation: %w", err)
		}

		next, err := s.repo.NextVersionNumber(ctx, tx, templateID)
		if err != nil {
			return err
		}
		v = DocumentTemplateVersion{
			ID:           uuid.New().String(),
			TemplateID:   templateID,
			Version:      next,
			Content:      content,
			PaperSize:    paperSize,
			Orientation:  orientation,
			MarginTop:    margins[0],
			MarginRight:  margins[1],
			MarginBottom: margins[2],
			MarginLeft:   margins[3],
			CreatedBy:    &actorID,
		}
		if fileName != "" {
			v.FileName = &fileName
		}
		if err := s.repo.CreateVersion(ctx, tx, &v); err != nil {
			return err
		}
		// Update the template's active_version_id directly inside this
		// transaction. Repository.Update resolves its own non-transactional
		// DB handle, so calling it here would not be part of the transaction;
		// using tx.Model(...).Update(...) directly is the less invasive
		// option compared to adding a tx-aware overload to Repository.Update.
		if err := tx.Model(&DocumentTemplate{}).Where("id = ?", templateID).Update("active_version_id", v.ID).Error; err != nil {
			return fmt.Errorf("failed to update template active_version_id: %w", err)
		}
		return s.repo.CreateAudit(ctx, tx, &DocumentTemplateAudit{
			ID:         uuid.New().String(),
			TemplateID: templateID,
			VersionID:  &v.ID,
			Action:     "VERSION_CREATED",
			ActorID:    actorIDPtr(actorID),
		})
	})
	if err != nil {
		return nil, err
	}
	decorateVersionFileURL(&v)
	return &v, nil
}

func (s *Service) ListVersions(ctx context.Context, templateID string) ([]DocumentTemplateVersion, error) {
	versions, err := s.repo.ListVersions(ctx, templateID)
	if err != nil {
		return nil, err
	}
	for i := range versions {
		decorateVersionFileURL(&versions[i])
	}
	return versions, nil
}

func (s *Service) GetVersion(ctx context.Context, templateID, versionID string) (*DocumentTemplateVersion, error) {
	// Pastikan template-nya ada (404 terpisah dari version-not-found).
	if _, err := s.repo.GetByID(ctx, templateID); err != nil {
		return nil, err
	}
	v, err := s.repo.GetVersion(ctx, templateID, versionID)
	if err != nil {
		return nil, err
	}
	decorateVersionFileURL(v)
	return v, nil
}
