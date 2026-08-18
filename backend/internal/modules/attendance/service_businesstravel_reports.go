package attendance

import (
	"context"
)

// =========================================================================
// Business Travel Reports
// =========================================================================

func (s *Service) GetBusinessTravelReport(ctx context.Context, fromDate, toDate, status string) ([]BusinessTravelReportRow, error) {
	return s.repo.FindTravelReport(ctx, fromDate, toDate, status)
}

func (s *Service) GetBusinessTravelFundingReport(ctx context.Context, fromDate, toDate, status string) ([]BusinessTravelFundingReportRow, error) {
	return s.repo.FindFundingReport(ctx, fromDate, toDate, status)
}

func (s *Service) GetBusinessTravelAdvanceReport(ctx context.Context, fromDate, toDate, status string) ([]BusinessTravelAdvanceReportRow, error) {
	return s.repo.FindAdvanceReport(ctx, fromDate, toDate, status)
}

func (s *Service) GetBusinessTravelReimbursementReport(ctx context.Context, fromDate, toDate, status string) ([]BusinessTravelReimbursementReportRow, error) {
	return s.repo.FindReimbursementReport(ctx, fromDate, toDate, status)
}

func (s *Service) GetBusinessTravelRefundReport(ctx context.Context, fromDate, toDate, status string) ([]BusinessTravelRefundReportRow, error) {
	return s.repo.FindRefundReport(ctx, fromDate, toDate, status)
}

func (s *Service) GetBusinessTravelCostReport(ctx context.Context, fromDate, toDate, status string) ([]BusinessTravelCostReportRow, error) {
	return s.repo.FindTravelCostReport(ctx, fromDate, toDate, status)
}
