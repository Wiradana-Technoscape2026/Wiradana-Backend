package model

type DashboardNotification struct {
	Type       string `json:"type"`
	MemberName string `json:"member_name"`
	DueDate    string `json:"due_date"`
}

type DashboardResponse struct {
	TotalMembers  int64                   `json:"total_members"`
	TotalSavings  int64                   `json:"total_savings"`
	ActiveLoans   int64                   `json:"active_loans"`
	OverdueLoans  int64                   `json:"overdue_loans"`
	Notifications []DashboardNotification `json:"notifications"`
}
