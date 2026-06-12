package model

// DashboardResponse — api_planning §3.8
type DashboardResponse struct {
	TotalMembers  int64            `json:"total_members"`
	TotalSavings  int64            `json:"total_savings"`
	ActiveLoans   int64            `json:"active_loans"`
	OverdueLoans  int64            `json:"overdue_loans"`
	Notifications []Notification   `json:"notifications"`
}

// Notification — installment due within 3 days
type Notification struct {
	Type       string `json:"type"`
	MemberName string `json:"member_name"`
	DueDate    string `json:"due_date"`
	PeriodNo   int    `json:"period_no"`
	TotalDue   int64  `json:"total_due"`
	LoanID     string `json:"loan_id"`
}
