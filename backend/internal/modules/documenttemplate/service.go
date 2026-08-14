package documenttemplate

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

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
	if err := s.checkCodeUnique(ctx, code); err != nil {
		return nil, err
	}

	tpl := &DocumentTemplate{
		ID:           uuid.New().String(),
		Name:         name,
		Code:         code,
		DocumentType: documentType,
		Status:       StatusInactive,
		IsDefault:    false,
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
		ActorID:    &actorID,
	}); err != nil {
		return nil, err
	}
	return tpl, nil
}

func (s *Service) CreateFromDefault(ctx context.Context, documentType, name, code, actorID string) (*DocumentTemplate, error) {
	if !ValidDocumentTypes[documentType] {
		return nil, &InvalidDocumentTypeError{DocumentType: documentType}
	}
	def, err := s.repo.FindDefaultByType(ctx, documentType)
	if err != nil {
		return nil, err
	}
	if err := s.checkCodeUnique(ctx, code); err != nil {
		return nil, err
	}

	tpl := &DocumentTemplate{
		ID:           uuid.New().String(),
		Name:         name,
		Code:         code,
		DocumentType: documentType,
		Content:      def.Content,
		Status:       StatusInactive,
		IsDefault:    false,
		IsActive:     true,
	}
	if err := s.repo.Create(ctx, tpl); err != nil {
		return nil, err
	}
	if err := s.repo.CreateAudit(ctx, nil, &DocumentTemplateAudit{
		ID:         uuid.New().String(),
		TemplateID: tpl.ID,
		Action:     "CREATED_FROM_DEFAULT",
		ActorID:    &actorID,
	}); err != nil {
		return nil, err
	}
	return tpl, nil
}

func (s *Service) Update(ctx context.Context, id, name, description string, actorID string) (*DocumentTemplate, error) {
	tpl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tpl.IsDefault {
		return nil, &ReferenceTemplateImmutableError{Action: "edited"}
	}
	tpl.Name = name
	if description != "" {
		tpl.Description = &description
	}
	if err := s.repo.Update(ctx, tpl); err != nil {
		return nil, err
	}
	if err := s.repo.CreateAudit(ctx, nil, &DocumentTemplateAudit{
		ID:         uuid.New().String(),
		TemplateID: tpl.ID,
		Action:     "UPDATED",
		ActorID:    &actorID,
	}); err != nil {
		return nil, err
	}
	return tpl, nil
}

func (s *Service) UpdateDefaultContent(ctx context.Context, documentType, content, actorID string) (*DocumentTemplate, error) {
	tpl, err := s.repo.FindDefaultByType(ctx, documentType)
	if err != nil {
		return nil, err
	}
	tpl.Content = &content
	if err := s.repo.Update(ctx, tpl); err != nil {
		return nil, err
	}
	if err := s.repo.CreateAudit(ctx, nil, &DocumentTemplateAudit{
		ID:         uuid.New().String(),
		TemplateID: tpl.ID,
		Action:     "DEFAULT_UPDATED",
		ActorID:    &actorID,
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
	if tpl.IsDefault {
		return &ReferenceTemplateImmutableError{Action: "deleted"}
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
		ActorID:    &actorID,
	})
}

func (s *Service) Activate(ctx context.Context, id, actorID string) (*DocumentTemplate, error) {
	target, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if target.IsDefault {
		return nil, &ReferenceTemplateImmutableError{Action: "activated"}
	}

	err = s.repo.WithTx(ctx, func(tx *gorm.DB) error {
		var previous DocumentTemplate
		err := tx.Where("type = ? AND status = ? AND id != ? AND deleted_at IS NULL", target.DocumentType, StatusActive, target.ID).First(&previous).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to find previous active template: %w", err)
		}
		if err == nil {
			previous.Status = StatusInactive
			if err := tx.Save(&previous).Error; err != nil {
				return fmt.Errorf("failed to deactivate previous active template: %w", err)
			}
			if err := s.repo.CreateAudit(ctx, tx, &DocumentTemplateAudit{
				ID:         uuid.New().String(),
				TemplateID: previous.ID,
				Action:     "DEACTIVATED",
				ActorID:    &actorID,
			}); err != nil {
				return err
			}
		}

		target.Status = StatusActive
		if err := tx.Save(target).Error; err != nil {
			return fmt.Errorf("failed to activate template: %w", err)
		}
		return s.repo.CreateAudit(ctx, tx, &DocumentTemplateAudit{
			ID:         uuid.New().String(),
			TemplateID: target.ID,
			Action:     "ACTIVATED",
			ActorID:    &actorID,
		})
	})
	if err != nil {
		return nil, err
	}
	return target, nil
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
		ActorID:    &actorID,
	}); err != nil {
		return nil, err
	}
	return tpl, nil
}

func (s *Service) CreateVersion(ctx context.Context, templateID, content, paperSize, orientation string, margins [4]int, actorID string) (*DocumentTemplateVersion, error) {
	var v DocumentTemplateVersion
	err := s.repo.WithTx(ctx, func(tx *gorm.DB) error {
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
			ActorID:    &actorID,
		})
	})
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *Service) ListVersions(ctx context.Context, templateID string) ([]DocumentTemplateVersion, error) {
	return s.repo.ListVersions(ctx, templateID)
}
