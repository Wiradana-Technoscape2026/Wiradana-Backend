package model

type ReportSummary struct {
	TotalSavings           int64 `json:"total_savings"`
	TotalInterestCollected int64 `json:"total_interest_collected"`
	TotalOutstanding       int64 `json:"total_outstanding"`
	TotalMembers           int64 `json:"total_members"`
	ActiveMembers          int64 `json:"active_members"`
	ActiveLoansCount       int64 `json:"active_loans_count"`
	OverdueLoansCount      int64 `json:"overdue_loans_count"`
	TotalDisbursed         int64 `json:"total_disbursed"`
}
