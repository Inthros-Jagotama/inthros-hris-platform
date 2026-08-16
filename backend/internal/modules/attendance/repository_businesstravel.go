package attendance

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// =========================================================================
// Business Travel (core)
// =========================================================================

func (r *Repository) CreateBusinessTravel(ctx context.Context, t *BusinessTravel) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(t).Error
}

func (r *Repository) FindBusinessTravelByID(ctx context.Context, id uuid.UUID) (*BusinessTravel, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var t BusinessTravel
	if err := db.Preload("Participants").Preload("Destinations").First(&t, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("business travel not found: %w", err)
	}
	return &t, nil
}

func (r *Repository) UpdateBusinessTravel(ctx context.Context, t *BusinessTravel) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(t).Error
}

func (r *Repository) ListBusinessTravels(ctx context.Context, requesterID *uuid.UUID, status string, page, perPage int) ([]BusinessTravel, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var travels []BusinessTravel
	var total int64
	query := db.Model(&BusinessTravel{})
	if requesterID != nil {
		query = query.Where("requester_id = ?", *requesterID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&travels).Error; err != nil {
		return nil, 0, err
	}
	return travels, total, nil
}

// =========================================================================
// Participants
// =========================================================================

func (r *Repository) CreateParticipant(ctx context.Context, p *BusinessTravelParticipant) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(p).Error
}

func (r *Repository) ListParticipantsByTravel(ctx context.Context, travelID uuid.UUID) ([]BusinessTravelParticipant, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var participants []BusinessTravelParticipant
	if err := db.Where("business_travel_id = ?", travelID).Find(&participants).Error; err != nil {
		return nil, err
	}
	return participants, nil
}

func (r *Repository) FindParticipantByID(ctx context.Context, id uuid.UUID) (*BusinessTravelParticipant, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var p BusinessTravelParticipant
	if err := db.First(&p, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("participant not found: %w", err)
	}
	return &p, nil
}

// =========================================================================
// Destinations
// =========================================================================

func (r *Repository) CreateDestination(ctx context.Context, d *BusinessTravelDestination) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(d).Error
}

func (r *Repository) ListDestinationsByTravel(ctx context.Context, travelID uuid.UUID) ([]BusinessTravelDestination, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var destinations []BusinessTravelDestination
	if err := db.Where("business_travel_id = ?", travelID).Order("sequence ASC").Find(&destinations).Error; err != nil {
		return nil, err
	}
	return destinations, nil
}

// =========================================================================
// Activities
// =========================================================================

func (r *Repository) CreateActivity(ctx context.Context, a *BusinessTravelActivity) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(a).Error
}

func (r *Repository) ListActivitiesByTravel(ctx context.Context, travelID uuid.UUID) ([]BusinessTravelActivity, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var activities []BusinessTravelActivity
	if err := db.Where("business_travel_id = ?", travelID).Order("activity_date ASC").Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

// =========================================================================
// Schedules
// =========================================================================

func (r *Repository) CreateSchedule(ctx context.Context, s *BusinessTravelSchedule) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(s).Error
}

func (r *Repository) ListSchedulesByTravel(ctx context.Context, travelID uuid.UUID) ([]BusinessTravelSchedule, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var schedules []BusinessTravelSchedule
	if err := db.Where("business_travel_id = ?", travelID).Order("created_at ASC").Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// =========================================================================
// Funding Methods (master)
// =========================================================================

func (r *Repository) CreateFundingMethod(ctx context.Context, m *FundingMethod) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(m).Error
}

func (r *Repository) ListFundingMethods(ctx context.Context, activeOnly bool) ([]FundingMethod, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	query := db.Model(&FundingMethod{})
	if activeOnly {
		query = query.Where("active = ?", true)
	}
	var methods []FundingMethod
	if err := query.Order("name ASC").Find(&methods).Error; err != nil {
		return nil, err
	}
	return methods, nil
}

func (r *Repository) FindFundingMethodByID(ctx context.Context, id uuid.UUID) (*FundingMethod, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var m FundingMethod
	if err := db.First(&m, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("funding method not found: %w", err)
	}
	return &m, nil
}

// =========================================================================
// Fundings
// =========================================================================

func (r *Repository) CreateFunding(ctx context.Context, f *Funding) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(f).Error
}

func (r *Repository) FindFundingByID(ctx context.Context, id uuid.UUID) (*Funding, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var f Funding
	if err := db.Preload("Documents").First(&f, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("funding not found: %w", err)
	}
	return &f, nil
}

func (r *Repository) UpdateFunding(ctx context.Context, f *Funding) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(f).Error
}

func (r *Repository) ListFundingsByTravel(ctx context.Context, travelID uuid.UUID) ([]Funding, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var fundings []Funding
	if err := db.Preload("Documents").Where("business_travel_id = ?", travelID).Order("created_at ASC").Find(&fundings).Error; err != nil {
		return nil, err
	}
	return fundings, nil
}

// =========================================================================
// Funding Documents
// =========================================================================

func (r *Repository) CreateFundingDocument(ctx context.Context, d *FundingDocument) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(d).Error
}

// =========================================================================
// Expense Categories (master)
// =========================================================================

func (r *Repository) CreateExpenseCategory(ctx context.Context, c *ExpenseCategory) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(c).Error
}

func (r *Repository) ListExpenseCategories(ctx context.Context, activeOnly bool) ([]ExpenseCategory, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	query := db.Model(&ExpenseCategory{})
	if activeOnly {
		query = query.Where("active = ?", true)
	}
	var categories []ExpenseCategory
	if err := query.Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *Repository) FindExpenseCategoryByID(ctx context.Context, id uuid.UUID) (*ExpenseCategory, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var c ExpenseCategory
	if err := db.First(&c, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("expense category not found: %w", err)
	}
	return &c, nil
}

// =========================================================================
// Expenses (actual)
// =========================================================================

func (r *Repository) CreateExpense(ctx context.Context, e *Expense) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(e).Error
}

func (r *Repository) FindExpenseByID(ctx context.Context, id uuid.UUID) (*Expense, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var e Expense
	if err := db.Preload("Documents").First(&e, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("expense not found: %w", err)
	}
	return &e, nil
}

func (r *Repository) UpdateExpense(ctx context.Context, e *Expense) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(e).Error
}

func (r *Repository) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	result := db.Where("id = ?", id).Delete(&Expense{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("expense not found")
	}
	return nil
}

func (r *Repository) ListExpensesByTravel(ctx context.Context, travelID uuid.UUID) ([]Expense, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var expenses []Expense
	if err := db.Preload("Documents").Where("business_travel_id = ?", travelID).Order("expense_date ASC").Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

// =========================================================================
// Expense Documents
// =========================================================================

func (r *Repository) CreateExpenseDocument(ctx context.Context, d *ExpenseDocument) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(d).Error
}

// =========================================================================
// Settlement
// =========================================================================

func (r *Repository) CreateSettlement(ctx context.Context, s *Settlement) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(s).Error
}

func (r *Repository) FindSettlementByID(ctx context.Context, id uuid.UUID) (*Settlement, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var s Settlement
	if err := db.Preload("Items").First(&s, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("settlement not found: %w", err)
	}
	return &s, nil
}

func (r *Repository) UpdateSettlement(ctx context.Context, s *Settlement) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(s).Error
}

func (r *Repository) ListSettlementsByTravel(ctx context.Context, travelID uuid.UUID) ([]Settlement, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var settlements []Settlement
	if err := db.Preload("Items").Where("business_travel_id = ?", travelID).Order("created_at DESC").Find(&settlements).Error; err != nil {
		return nil, err
	}
	return settlements, nil
}

func (r *Repository) CreateSettlementItem(ctx context.Context, i *SettlementItem) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(i).Error
}

// =========================================================================
// Refund & Reimbursement (hasil settlement, §35-36)
// =========================================================================

func (r *Repository) CreateRefund(ctx context.Context, rf *Refund) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(rf).Error
}

func (r *Repository) CreateTravelReimbursement(ctx context.Context, tr *TravelReimbursement) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(tr).Error
}

func (r *Repository) FindRefundByID(ctx context.Context, id uuid.UUID) (*Refund, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rf Refund
	if err := db.First(&rf, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("refund not found: %w", err)
	}
	return &rf, nil
}

func (r *Repository) UpdateRefund(ctx context.Context, rf *Refund) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(rf).Error
}

func (r *Repository) ListRefundsByTravel(ctx context.Context, travelID uuid.UUID) ([]Refund, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var refunds []Refund
	if err := db.Where("business_travel_id = ?", travelID).Order("created_at ASC").Find(&refunds).Error; err != nil {
		return nil, err
	}
	return refunds, nil
}

func (r *Repository) FindTravelReimbursementByID(ctx context.Context, id uuid.UUID) (*TravelReimbursement, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var tr TravelReimbursement
	if err := db.First(&tr, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("reimbursement not found: %w", err)
	}
	return &tr, nil
}

func (r *Repository) UpdateTravelReimbursement(ctx context.Context, tr *TravelReimbursement) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(tr).Error
}

func (r *Repository) ListTravelReimbursementsByTravel(ctx context.Context, travelID uuid.UUID) ([]TravelReimbursement, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var reimbursements []TravelReimbursement
	if err := db.Where("business_travel_id = ?", travelID).Order("created_at ASC").Find(&reimbursements).Error; err != nil {
		return nil, err
	}
	return reimbursements, nil
}

// =========================================================================
// Travel Documents
// =========================================================================

func (r *Repository) CreateTravelDocument(ctx context.Context, d *TravelDocument) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(d).Error
}

func (r *Repository) ListTravelDocumentsByTravel(ctx context.Context, travelID uuid.UUID) ([]TravelDocument, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var docs []TravelDocument
	if err := db.Where("business_travel_id = ?", travelID).Order("created_at ASC").Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *Repository) DeleteTravelDocument(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	result := db.Where("id = ?", id).Delete(&TravelDocument{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("travel document not found")
	}
	return nil
}

// FindBusinessTravelByIDForOwnership loads a bare BusinessTravel (no
// preloads) — used by sub-resource creators (activity/schedule) that only
// need to check the parent travel exists and read its status, without the
// cost of preloading participants/destinations.
func (r *Repository) FindBusinessTravelByIDForOwnership(ctx context.Context, id uuid.UUID) (*BusinessTravel, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var t BusinessTravel
	if err := db.First(&t, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("business travel not found: %w", err)
	}
	return &t, nil
}
