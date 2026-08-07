package approval

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository untuk database operations Approval Engine.
type Repository struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
}

// NewRepository membuat Repository baru.
func NewRepository(dbResolver func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbResolver: dbResolver}
}

func (r *Repository) getDB(ctx context.Context) (*gorm.DB, error) {
	return r.dbResolver(ctx)
}

// =========================================================================
// Approval Flows
// =========================================================================

func (r *Repository) CreateFlow(ctx context.Context, flow *ApprovalFlow) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(flow).Error; err != nil {
		return fmt.Errorf("failed to create approval flow: %w", err)
	}
	return nil
}

func (r *Repository) FindFlowByID(ctx context.Context, id uuid.UUID) (*ApprovalFlow, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var flow ApprovalFlow
	if err := db.Where("id = ? AND deleted_at IS NULL", id).First(&flow).Error; err != nil {
		return nil, fmt.Errorf("approval flow not found: %w", err)
	}
	return &flow, nil
}

func (r *Repository) FindFlowByIDWithSteps(ctx context.Context, id uuid.UUID) (*ApprovalFlow, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var flow ApprovalFlow
	if err := db.Preload("Steps", func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted_at IS NULL").Order("step_order ASC")
	}).Preload("Steps.Organizations").
		Where("id = ? AND deleted_at IS NULL", id).First(&flow).Error; err != nil {
		return nil, fmt.Errorf("approval flow not found: %w", err)
	}
	return &flow, nil
}

// FindActiveFlowByModule resolves the flow a caller should use for a given
// module without the caller having to pick a flow_id manually (mirrors how
// leave/reimbursement/etc. still require an explicit flow_id — this is for
// callers, like the KPI self-assessment two-phase submission, that want
// auto-resolution instead). Picks the highest-version active flow if more
// than one somehow exists for the module.
func (r *Repository) FindActiveFlowByModule(ctx context.Context, module string) (*ApprovalFlow, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var flow ApprovalFlow
	if err := db.Where("module = ? AND is_active = ? AND deleted_at IS NULL", module, true).
		Order("version DESC").First(&flow).Error; err != nil {
		return nil, fmt.Errorf("no active approval flow found for module %q: %w", module, err)
	}
	return &flow, nil
}

func (r *Repository) ListFlows(ctx context.Context, page, perPage int) ([]ApprovalFlow, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var flows []ApprovalFlow
	var total int64

	query := db.Model(&ApprovalFlow{}).Where("deleted_at IS NULL")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count flows: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&flows).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list flows: %w", err)
	}
	return flows, total, nil
}

func (r *Repository) ListFlowsByModule(ctx context.Context, module string) ([]ApprovalFlow, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var flows []ApprovalFlow
	if err := db.Where("module = ? AND is_active = ? AND deleted_at IS NULL", module, true).
		Order("version DESC").Find(&flows).Error; err != nil {
		return nil, fmt.Errorf("failed to list flows by module: %w", err)
	}
	return flows, nil
}

func (r *Repository) UpdateFlow(ctx context.Context, flow *ApprovalFlow) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Save(flow).Error; err != nil {
		return fmt.Errorf("failed to update approval flow: %w", err)
	}
	return nil
}

func (r *Repository) SoftDeleteFlow(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	result := db.Model(&ApprovalFlow{}).Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now)
	if result.Error != nil {
		return fmt.Errorf("failed to delete approval flow: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("approval flow not found")
	}
	// Soft delete all steps under this flow
	db.Model(&ApprovalFlowStep{}).Where("flow_id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now)
	return nil
}

// =========================================================================
// Approval Flow Steps
// =========================================================================

func (r *Repository) CreateStep(ctx context.Context, step *ApprovalFlowStep) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(step).Error; err != nil {
		return fmt.Errorf("failed to create approval step: %w", err)
	}
	return nil
}

func (r *Repository) FindStepByID(ctx context.Context, id uuid.UUID) (*ApprovalFlowStep, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var step ApprovalFlowStep
	if err := db.Where("id = ? AND deleted_at IS NULL", id).First(&step).Error; err != nil {
		return nil, fmt.Errorf("approval step not found: %w", err)
	}
	return &step, nil
}

func (r *Repository) ListStepsByFlowID(ctx context.Context, flowID uuid.UUID) ([]ApprovalFlowStep, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var steps []ApprovalFlowStep
	if err := db.Preload("Organizations").
		Where("flow_id = ? AND deleted_at IS NULL", flowID).
		Order("step_order ASC").Find(&steps).Error; err != nil {
		return nil, fmt.Errorf("failed to list steps: %w", err)
	}
	return steps, nil
}

func (r *Repository) GetMaxStepOrder(ctx context.Context, flowID uuid.UUID) (int, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return 0, err
	}
	var maxOrder int
	if err := db.Model(&ApprovalFlowStep{}).
		Select("COALESCE(MAX(step_order), 0)").
		Where("flow_id = ? AND deleted_at IS NULL", flowID).
		Scan(&maxOrder).Error; err != nil {
		return 0, fmt.Errorf("failed to get max step order: %w", err)
	}
	return maxOrder, nil
}

func (r *Repository) UpdateStep(ctx context.Context, step *ApprovalFlowStep) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Save(step).Error; err != nil {
		return fmt.Errorf("failed to update approval step: %w", err)
	}
	return nil
}

func (r *Repository) SoftDeleteStep(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	result := db.Model(&ApprovalFlowStep{}).Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now)
	if result.Error != nil {
		return fmt.Errorf("failed to delete approval step: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("approval step not found")
	}
	return nil
}

// =========================================================================
// Approval Instances
// =========================================================================

func (r *Repository) CreateInstance(ctx context.Context, instance *ApprovalInstance) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(instance).Error; err != nil {
		return fmt.Errorf("failed to create approval instance: %w", err)
	}
	return nil
}

func (r *Repository) FindInstanceByID(ctx context.Context, id uuid.UUID) (*ApprovalInstance, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var instance ApprovalInstance
	if err := db.Where("id = ? AND deleted_at IS NULL", id).First(&instance).Error; err != nil {
		return nil, fmt.Errorf("approval instance not found: %w", err)
	}
	return &instance, nil
}

func (r *Repository) FindInstanceByIDWithRelations(ctx context.Context, id uuid.UUID) (*ApprovalInstance, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var instance ApprovalInstance
	if err := db.Preload("Flow").
		Preload("Actions", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL").Order("step_order ASC")
		}).
		Where("id = ? AND deleted_at IS NULL", id).First(&instance).Error; err != nil {
		return nil, fmt.Errorf("approval instance not found: %w", err)
	}

	// Load steps from flow
	if instance.Flow != nil {
		steps, err := r.ListStepsByFlowID(ctx, instance.FlowID)
		if err != nil {
			return nil, err
		}
		instance.Steps = steps
	}

	return &instance, nil
}

func (r *Repository) FindInstanceByDocument(ctx context.Context, module string, documentID uuid.UUID) (*ApprovalInstance, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var instance ApprovalInstance
	if err := db.Where("module = ? AND document_id = ? AND deleted_at IS NULL", module, documentID).
		First(&instance).Error; err != nil {
		return nil, fmt.Errorf("approval instance not found: %w", err)
	}
	return &instance, nil
}

func (r *Repository) ListInstances(ctx context.Context, page, perPage int, module, status string) ([]ApprovalInstance, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var instances []ApprovalInstance
	var total int64

	query := db.Model(&ApprovalInstance{}).Where("deleted_at IS NULL")
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count instances: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Preload("Flow").
		Offset(offset).Limit(perPage).Order("created_at DESC").Find(&instances).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list instances: %w", err)
	}
	return instances, total, nil
}

func (r *Repository) UpdateInstance(ctx context.Context, instance *ApprovalInstance) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Save(instance).Error; err != nil {
		return fmt.Errorf("failed to update approval instance: %w", err)
	}
	return nil
}

func (r *Repository) CancelInstance(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	return db.Transaction(func(tx *gorm.DB) error {
		// Update instance status
		result := tx.Model(&ApprovalInstance{}).
			Where("id = ? AND status = ? AND deleted_at IS NULL", id, InstanceStatusPending).
			Updates(map[string]interface{}{
				"status":     InstanceStatusCancelled,
				"updated_at": now,
			})
		if result.Error != nil {
			return fmt.Errorf("failed to cancel instance: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("instance not found or not in PENDING status")
		}

		// Cancel all pending tasks
		tx.Model(&ApprovalTask{}).
			Where("instance_id = ? AND status = ?", id, TaskStatusPending).
			Updates(map[string]interface{}{
				"status":     TaskStatusCancelled,
				"updated_at": now,
			})

		return nil
	})
}

// =========================================================================
// Approval Actions
// =========================================================================

func (r *Repository) CreateAction(ctx context.Context, action *ApprovalAction) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(action).Error; err != nil {
		return fmt.Errorf("failed to create approval action: %w", err)
	}
	return nil
}

func (r *Repository) ListActionsByInstanceID(ctx context.Context, instanceID uuid.UUID) ([]ApprovalAction, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var actions []ApprovalAction
	if err := db.Where("instance_id = ?", instanceID).
		Order("created_at ASC").Find(&actions).Error; err != nil {
		return nil, fmt.Errorf("failed to list actions: %w", err)
	}
	return actions, nil
}

// =========================================================================
// Approval Tasks
// =========================================================================

func (r *Repository) CreateTasks(ctx context.Context, tasks []ApprovalTask) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		return nil
	}
	if err := db.Create(&tasks).Error; err != nil {
		return fmt.Errorf("failed to create approval tasks: %w", err)
	}
	return nil
}

func (r *Repository) FindPendingTasksByAssignee(ctx context.Context, assigneeType string, assigneeID uuid.UUID, page, perPage int) ([]ApprovalTask, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var tasks []ApprovalTask
	var total int64

	query := db.Model(&ApprovalTask{}).
		Where("assignee_type = ? AND assignee_id = ? AND status = ? AND deleted_at IS NULL",
			assigneeType, assigneeID, TaskStatusPending)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count tasks: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list tasks: %w", err)
	}
	return tasks, total, nil
}

func (r *Repository) FindTasksByInstanceID(ctx context.Context, instanceID uuid.UUID) ([]ApprovalTask, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var tasks []ApprovalTask
	if err := db.Where("instance_id = ? AND deleted_at IS NULL", instanceID).
		Order("step_order ASC").Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	return tasks, nil
}

func (r *Repository) UpdateTaskStatus(ctx context.Context, id uuid.UUID, status TaskStatus) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	result := db.Model(&ApprovalTask{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": now,
		})
	if result.Error != nil {
		return fmt.Errorf("failed to update task status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

// =========================================================================
// Transactional operations
// =========================================================================

// CreateInstanceWithTasks membuat instance dan tasks dalam satu transaction.
func (r *Repository) CreateInstanceWithTasks(ctx context.Context, instance *ApprovalInstance, tasks []ApprovalTask) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// Create instance
		if err := tx.Create(instance).Error; err != nil {
			return fmt.Errorf("failed to create approval instance: %w", err)
		}

		// Create tasks for first step
		for i := range tasks {
			tasks[i].InstanceID = instance.ID
		}
		if len(tasks) > 0 {
			if err := tx.Create(&tasks).Error; err != nil {
				return fmt.Errorf("failed to create approval tasks: %w", err)
			}
		}

		return nil
	})
}

// ApproveStep melakukan approve untuk step tertentu dan membuat tasks untuk step berikutnya.
func (r *Repository) ApproveStep(ctx context.Context, instanceID uuid.UUID, currentStep int, nextStep int, nextTasks []ApprovalTask) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// Update instance to next step or approve if no more steps
		// Note: Service layer already validates step conditions before calling this
		// So we don't need to check pending tasks here
		if nextStep > currentStep {
			// Move to next step
			if err := tx.Model(&ApprovalInstance{}).
				Where("id = ?", instanceID).
				Updates(map[string]interface{}{
					"current_step": nextStep,
					"updated_at":   time.Now(),
				}).Error; err != nil {
				return fmt.Errorf("failed to update instance step: %w", err)
			}

			// Create tasks for next step
			for i := range nextTasks {
				nextTasks[i].InstanceID = instanceID
			}
			if len(nextTasks) > 0 {
				if err := tx.Create(&nextTasks).Error; err != nil {
					return fmt.Errorf("failed to create next step tasks: %w", err)
				}
			}
		} else {
			// Approve instance completely
			if err := tx.Model(&ApprovalInstance{}).
				Where("id = ?", instanceID).
				Updates(map[string]interface{}{
					"status":     InstanceStatusApproved,
					"updated_at": time.Now(),
				}).Error; err != nil {
				return fmt.Errorf("failed to approve instance: %w", err)
			}

			// Persist any trailing WATCHER-step tasks (already DONE, informational
			// only) that were resolved while walking past the final gate.
			for i := range nextTasks {
				nextTasks[i].InstanceID = instanceID
			}
			if len(nextTasks) > 0 {
				if err := tx.Create(&nextTasks).Error; err != nil {
					return fmt.Errorf("failed to create trailing watcher tasks: %w", err)
				}
			}
		}

		return nil
	})
}

// =========================================================================
// Approval Flow Step Organizations
// =========================================================================

// ReplaceStepOrganizations mengganti seluruh target Organization untuk sebuah
// step (dipakai saat create/update step dengan approver_type ORGANIZATION/BOTH).
func (r *Repository) ReplaceStepOrganizations(ctx context.Context, stepID uuid.UUID, orgIDs []uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("step_id = ?", stepID).Delete(&ApprovalFlowStepOrganization{}).Error; err != nil {
			return fmt.Errorf("failed to clear step organizations: %w", err)
		}
		if len(orgIDs) == 0 {
			return nil
		}
		rows := make([]ApprovalFlowStepOrganization, 0, len(orgIDs))
		for _, orgID := range orgIDs {
			rows = append(rows, ApprovalFlowStepOrganization{StepID: stepID, OrganizationID: orgID})
		}
		if err := tx.Create(&rows).Error; err != nil {
			return fmt.Errorf("failed to create step organizations: %w", err)
		}
		return nil
	})
}

func (r *Repository) ListStepOrganizations(ctx context.Context, stepID uuid.UUID) ([]ApprovalFlowStepOrganization, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rows []ApprovalFlowStepOrganization
	if err := db.Where("step_id = ?", stepID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("failed to list step organizations: %w", err)
	}
	return rows, nil
}

// =========================================================================
// Organization-hierarchy resolution (raw queries — approval must not import
// the employee/organization packages directly, same reasoning the performance
// module already applies for cross-module lookups)
// =========================================================================

// GetOrganizationParentID mengembalikan parent_id dari sebuah Organization (nil jika root/tidak ada).
func (r *Repository) GetOrganizationParentID(ctx context.Context, orgID uuid.UUID) (*uuid.UUID, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var parentIDStr *string
	if err := db.Table("organizations").
		Where("id = ? AND deleted_at IS NULL", orgID).
		Select("parent_id").
		Scan(&parentIDStr).Error; err != nil {
		return nil, fmt.Errorf("failed to resolve parent organization: %w", err)
	}
	if parentIDStr == nil || *parentIDStr == "" {
		return nil, nil
	}
	parentID, err := uuid.Parse(*parentIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid parent organization id: %w", err)
	}
	return &parentID, nil
}

// GetSubmitterOrganizationID resolve dari platform user (pembuat instance) ke
// Organization tempat dia bekerja saat ini: user -> employee_accounts -> employee
// -> employments (current) -> organization_id.
func (r *Repository) GetSubmitterOrganizationID(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var orgIDStr *string
	err = db.Table("employee_accounts AS ea").
		Joins("JOIN employments AS emp ON emp.employee_id = ea.employee_id").
		Where("ea.user_id = ? AND (emp.effective_end_date IS NULL)", userID).
		Order("emp.effective_date DESC").
		Limit(1).
		Select("emp.organization_id").
		Scan(&orgIDStr).Error
	if err != nil {
		return nil, fmt.Errorf("failed to resolve submitter organization: %w", err)
	}
	if orgIDStr == nil || *orgIDStr == "" {
		return nil, nil
	}
	orgID, err := uuid.Parse(*orgIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid organization id: %w", err)
	}
	return &orgID, nil
}

// submitterInfo carries the display fields the "My Tasks" inbox needs about
// whoever submitted an instance (name, employee code, current org name).
type submitterInfo struct {
	Name             string
	EmployeeCode     string
	OrganizationName string
}

// GetSubmitterInfoByUserIDs resolves display info (name, employee code,
// current organization name) for a batch of submitters at once — used to
// enrich the "My Tasks" list without an N+1 query per row. Same
// employee_accounts -> employments -> organizations raw-query pattern as
// GetSubmitterOrganizationID (approval must not import the employee/
// organization packages directly).
func (r *Repository) GetSubmitterInfoByUserIDs(ctx context.Context, userIDs []uuid.UUID) (map[string]submitterInfo, error) {
	result := make(map[string]submitterInfo, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
	}
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	type row struct {
		UserID           string
		Name             string
		EmployeeCode     string
		OrganizationName string
	}
	var rows []row
	err = db.Table("employee_accounts AS ea").
		Joins("JOIN employees AS emp ON emp.id = ea.employee_id").
		Joins("JOIN employments AS em ON em.employee_id = ea.employee_id AND em.effective_end_date IS NULL").
		Joins("JOIN organizations AS o ON o.id = em.organization_id").
		Where("ea.user_id IN ?", userIDs).
		Select("ea.user_id AS user_id, emp.name AS name, emp.employee_id AS employee_code, o.nomenclature AS organization_name").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to resolve submitter info: %w", err)
	}
	for _, rrow := range rows {
		result[rrow.UserID] = submitterInfo{Name: rrow.Name, EmployeeCode: rrow.EmployeeCode, OrganizationName: rrow.OrganizationName}
	}
	return result, nil
}

// GetInstancesByIDs batch-fetches instances (for FlowID/CreatedBy) without
// their steps/actions — used to enrich task lists.
func (r *Repository) GetInstancesByIDs(ctx context.Context, ids []uuid.UUID) (map[string]ApprovalInstance, error) {
	result := make(map[string]ApprovalInstance, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var instances []ApprovalInstance
	if err := db.Where("id IN ? AND deleted_at IS NULL", ids).Find(&instances).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch instances: %w", err)
	}
	for _, inst := range instances {
		result[inst.ID.String()] = inst
	}
	return result, nil
}

// GetFlowNamesByIDs batch-fetches flow names for a set of flow IDs.
func (r *Repository) GetFlowNamesByIDs(ctx context.Context, ids []uuid.UUID) (map[string]string, error) {
	result := make(map[string]string, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var flows []ApprovalFlow
	if err := db.Where("id IN ? AND deleted_at IS NULL", ids).Find(&flows).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch flows: %w", err)
	}
	for _, f := range flows {
		result[f.ID.String()] = f.Name
	}
	return result, nil
}

// GetUserIDsByOrganization mengembalikan platform user ID (bukan employee ID) dari
// employee yang saat ini menempati sebuah Organization: organization_id -> employments
// (current) -> employee_id -> employee_accounts -> user_id.
func (r *Repository) GetUserIDsByOrganization(ctx context.Context, orgID uuid.UUID) ([]uuid.UUID, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var userIDStrs []string
	err = db.Table("employments AS emp").
		Joins("JOIN employee_accounts AS ea ON ea.employee_id = emp.employee_id").
		Where("emp.organization_id = ? AND (emp.effective_end_date IS NULL)", orgID).
		Distinct().
		Pluck("ea.user_id", &userIDStrs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to resolve organization assignees: %w", err)
	}
	userIDs := make([]uuid.UUID, 0, len(userIDStrs))
	for _, s := range userIDStrs {
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("invalid user id: %w", err)
		}
		userIDs = append(userIDs, id)
	}
	return userIDs, nil
}

// RejectInstance melakukan reject pada instance dan cancel semua pending tasks.
func (r *Repository) RejectInstance(ctx context.Context, instanceID uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// Update instance status to rejected
		result := tx.Model(&ApprovalInstance{}).
			Where("id = ? AND status = ? AND deleted_at IS NULL", instanceID, InstanceStatusPending).
			Updates(map[string]interface{}{
				"status":     InstanceStatusRejected,
				"updated_at": time.Now(),
			})
		if result.Error != nil {
			return fmt.Errorf("failed to reject instance: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("instance not found or not in PENDING status")
		}

		// Cancel all pending tasks
		tx.Model(&ApprovalTask{}).
			Where("instance_id = ? AND status = ?", instanceID, TaskStatusPending).
			Updates(map[string]interface{}{
				"status":     TaskStatusCancelled,
				"updated_at": time.Now(),
			})

		return nil
	})
}
