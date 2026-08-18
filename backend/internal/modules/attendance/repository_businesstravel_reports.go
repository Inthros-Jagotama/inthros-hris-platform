package attendance

import (
	"context"
)

// =========================================================================
// Business Travel Reports (§ report endpoints — permission
// "business_travel.report", digate via requireBusinessTravelReport() di
// routes.go). Semua query di bawah ini agregat org-wide (lintas requester),
// difilter rentang tanggal + status opsional.
// =========================================================================

// BusinessTravelReportRow — satu baris laporan daftar perjalanan dinas.
type BusinessTravelReportRow struct {
	ID               string  `json:"id"`
	RequestNumber    string  `json:"request_number"`
	Title            string  `json:"title"`
	RequesterID      string  `json:"requester_id"`
	StartDate        string  `json:"start_date"`
	EndDate          string  `json:"end_date"`
	Status           string  `json:"status"`
	ApprovalStatus   string  `json:"approval_status"`
	DestinationCity  *string `json:"destination_city"`
	DestinationCount int     `json:"destination_count"`
}

// FindTravelReport mengambil daftar perjalanan dinas dalam rentang start_date,
// beserta ringkasan destinasi (kota pertama + jumlah destinasi).
func (r *Repository) FindTravelReport(ctx context.Context, fromDate, toDate, status string) ([]BusinessTravelReportRow, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rows []BusinessTravelReportRow
	query := db.Table("business_travels AS t").
		Select(`t.id, t.request_number, t.title, t.requester_id, t.start_date, t.end_date, t.status, t.approval_status,
			(SELECT d.city FROM business_travel_destinations d WHERE d.business_travel_id = t.id AND d.deleted_at IS NULL ORDER BY d.sequence ASC LIMIT 1) AS destination_city,
			(SELECT COUNT(*) FROM business_travel_destinations d2 WHERE d2.business_travel_id = t.id AND d2.deleted_at IS NULL) AS destination_count`).
		Where("t.deleted_at IS NULL").
		Where("t.start_date >= ? AND t.start_date <= ?", fromDate, toDate)
	if status != "" {
		query = query.Where("t.status = ?", status)
	}
	if err := query.Order("t.start_date DESC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// BusinessTravelFundingReportRow — satu baris laporan funding perjalanan dinas.
type BusinessTravelFundingReportRow struct {
	TravelID        string  `json:"travel_id"`
	RequestNumber   string  `json:"request_number"`
	ParticipantID   *string `json:"participant_id"`
	FundingMethodID string  `json:"funding_method_id"`
	FundingMethod   string  `json:"funding_method"`
	Amount          float64 `json:"amount"`
	FundingDate     *string `json:"funding_date"`
	FundedBy        *string `json:"funded_by"`
	Status          string  `json:"status"`
}

func (r *Repository) FindFundingReport(ctx context.Context, fromDate, toDate, status string) ([]BusinessTravelFundingReportRow, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rows []BusinessTravelFundingReportRow
	query := db.Table("business_travel_fundings AS f").
		Select(`f.business_travel_id AS travel_id, t.request_number, f.participant_id, f.funding_method_id,
			m.name AS funding_method, f.amount, f.funding_date, f.funded_by, f.status`).
		Joins("JOIN business_travels t ON t.id = f.business_travel_id").
		Joins("JOIN business_travel_funding_methods m ON m.id = f.funding_method_id").
		Where("f.deleted_at IS NULL").
		Where("f.funding_date >= ? AND f.funding_date <= ?", fromDate, toDate)
	if status != "" {
		query = query.Where("f.status = ?", status)
	}
	if err := query.Order("f.funding_date DESC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// BusinessTravelAdvanceReportRow — ringkasan advance vs actual expense per
// settlement (participant/travel).
type BusinessTravelAdvanceReportRow struct {
	TravelID           string  `json:"travel_id"`
	RequestNumber      string  `json:"request_number"`
	ParticipantID      *string `json:"participant_id"`
	TotalAdvance       float64 `json:"total_advance"`
	TotalActualExpense float64 `json:"total_actual_expense"`
	Remaining          float64 `json:"remaining"`
	SettlementStatus   string  `json:"settlement_status"`
}

func (r *Repository) FindAdvanceReport(ctx context.Context, fromDate, toDate, status string) ([]BusinessTravelAdvanceReportRow, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rows []BusinessTravelAdvanceReportRow
	query := db.Table("business_travel_settlements AS s").
		Select(`s.business_travel_id AS travel_id, t.request_number, s.participant_id,
			s.total_advance, s.total_actual_expense, (s.total_advance - s.total_actual_expense) AS remaining,
			s.status AS settlement_status`).
		Joins("JOIN business_travels t ON t.id = s.business_travel_id").
		Where("s.deleted_at IS NULL").
		Where("s.created_at >= ? AND s.created_at <= ?", fromDate, toDate)
	if status != "" {
		query = query.Where("s.status = ?", status)
	}
	if err := query.Order("s.created_at DESC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// BusinessTravelReimbursementReportRow — laporan klaim reimbursement travel.
type BusinessTravelReimbursementReportRow struct {
	TravelID         string  `json:"travel_id"`
	RequestNumber    string  `json:"request_number"`
	ParticipantID    *string `json:"participant_id"`
	Amount           float64 `json:"amount"`
	ApprovedAt       *string `json:"approved_at"`
	PaidAt           *string `json:"paid_at"`
	PaymentReference *string `json:"payment_reference"`
	Status           string  `json:"status"`
}

func (r *Repository) FindReimbursementReport(ctx context.Context, fromDate, toDate, status string) ([]BusinessTravelReimbursementReportRow, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rows []BusinessTravelReimbursementReportRow
	query := db.Table("business_travel_reimbursements AS rb").
		Select(`rb.business_travel_id AS travel_id, t.request_number, rb.participant_id, rb.amount,
			rb.approved_at, rb.paid_at, rb.payment_reference, rb.status`).
		Joins("JOIN business_travels t ON t.id = rb.business_travel_id").
		Where("rb.deleted_at IS NULL").
		Where("rb.created_at >= ? AND rb.created_at <= ?", fromDate, toDate)
	if status != "" {
		query = query.Where("rb.status = ?", status)
	}
	if err := query.Order("rb.created_at DESC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// BusinessTravelRefundReportRow — laporan refund kelebihan advance.
type BusinessTravelRefundReportRow struct {
	TravelID      string  `json:"travel_id"`
	RequestNumber string  `json:"request_number"`
	ParticipantID *string `json:"participant_id"`
	Advance       float64 `json:"advance"`
	Actual        float64 `json:"actual"`
	RefundAmount  float64 `json:"refund_amount"`
	RefundDate    *string `json:"refund_date"`
	Status        string  `json:"status"`
}

func (r *Repository) FindRefundReport(ctx context.Context, fromDate, toDate, status string) ([]BusinessTravelRefundReportRow, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rows []BusinessTravelRefundReportRow
	query := db.Table("business_travel_refunds AS rf").
		Select(`rf.business_travel_id AS travel_id, t.request_number, rf.participant_id,
			COALESCE(s.total_advance, 0) AS advance, COALESCE(s.total_actual_expense, 0) AS actual,
			rf.refund_amount, rf.refund_date, rf.status`).
		Joins("JOIN business_travels t ON t.id = rf.business_travel_id").
		Joins("LEFT JOIN business_travel_settlements s ON s.id = rf.settlement_id").
		Where("rf.deleted_at IS NULL").
		Where("rf.created_at >= ? AND rf.created_at <= ?", fromDate, toDate)
	if status != "" {
		query = query.Where("rf.status = ?", status)
	}
	if err := query.Order("rf.created_at DESC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// BusinessTravelCostReportRow — total biaya per perjalanan dinas (company
// paid, advance, reimbursement, actual cost keseluruhan).
type BusinessTravelCostReportRow struct {
	TravelID           string  `json:"travel_id"`
	RequestNumber      string  `json:"request_number"`
	Title              string  `json:"title"`
	TotalCompanyPaid   float64 `json:"total_company_paid"`
	TotalAdvance       float64 `json:"total_advance"`
	TotalReimbursement float64 `json:"total_reimbursement"`
	TotalActualCost    float64 `json:"total_actual_cost"`
}

func (r *Repository) FindTravelCostReport(ctx context.Context, fromDate, toDate, status string) ([]BusinessTravelCostReportRow, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rows []BusinessTravelCostReportRow
	query := db.Table("business_travels AS t").
		Select(`t.id AS travel_id, t.request_number, t.title,
			COALESCE(SUM(s.total_company_paid), 0) AS total_company_paid,
			COALESCE(SUM(s.total_advance), 0) AS total_advance,
			COALESCE(SUM(s.total_reimbursement), 0) AS total_reimbursement,
			COALESCE(SUM(s.total_actual_expense), 0) AS total_actual_cost`).
		Joins("LEFT JOIN business_travel_settlements s ON s.business_travel_id = t.id AND s.deleted_at IS NULL").
		Where("t.deleted_at IS NULL").
		Where("t.start_date >= ? AND t.start_date <= ?", fromDate, toDate).
		Group("t.id, t.request_number, t.title")
	if status != "" {
		query = query.Where("t.status = ?", status)
	}
	if err := query.Order("t.start_date DESC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
