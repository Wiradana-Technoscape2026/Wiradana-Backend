package model

type DashboardStats struct {
	ActiveMembers          int64 `gorm:"column:active_members"`
	TotalMembers           int64 `gorm:"column:total_members"`
	TotalSavings           int64 `gorm:"column:total_savings"`
	ActiveLoans            int64 `gorm:"column:active_loans"`
	ActiveLoansOutstanding int64 `gorm:"column:active_loans_outstanding"`
	OverdueLoans           int64 `gorm:"column:overdue_loans"`
}

type UpcomingInstallment struct {
	InstallmentID string `json:"installment_id"`
	LoanID        string `json:"loan_id"`
	MemberName    string `json:"member_name"`
	PeriodNo      int    `json:"period_no"`
	DueDate       string `json:"due_date"`
	TotalDue      int64  `json:"total_due"`
	Status        string `json:"status"`
}

type PendingApplication struct {
	ID          string `json:"id"`
	MemberName  string `json:"member_name"`
	Amount      int64  `json:"amount"`
	TenorMonths int    `json:"tenor_months"`
	Purpose     string `json:"purpose"`
	Grade       string `json:"grade"`
}

type DashboardResponse struct {
	ActiveMembers             int64                 `json:"active_members"`
	TotalMembers              int64                 `json:"total_members"`
	TotalSavings              int64                 `json:"total_savings"`
	ActiveLoans               int64                 `json:"active_loans"`
	ActiveLoansOutstanding    int64                 `json:"active_loans_outstanding"`
	OverdueLoans              int64                 `json:"overdue_loans"`
	UpcomingInstallmentsCount int64                 `json:"upcoming_installments_count"`
	UpcomingInstallments      []UpcomingInstallment `json:"upcoming_installments"`
	PendingApplicationsCount  int64                 `json:"pending_applications_count"`
	PendingApplications       []PendingApplication  `json:"pending_applications"`
}
