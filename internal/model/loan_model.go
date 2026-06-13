package model

// ---- Loan Config ----

type LoanConfigResponse struct {
	ID              string  `json:"id"`
	FlatRateMonthly float64 `json:"flat_rate_monthly"`
	MaxPlafond      int64   `json:"max_plafond"`
	PenaltyDaily    int64   `json:"penalty_daily"`
}

type UpdateLoanConfigRequest struct {
	FlatRateMonthly *float64 `json:"flat_rate_monthly" validate:"omitempty,gt=0"`
	MaxPlafond      *int64   `json:"max_plafond" validate:"omitempty,gt=0"`
	PenaltyDaily    *int64   `json:"penalty_daily" validate:"omitempty,gte=0"`
}

// ---- Credit Assessment ----

type CreditAssessmentResponse struct {
	ID             string             `json:"id"`
	ApplicationID  string             `json:"application_id"`
	Score          int                `json:"score"`
	Grade          string             `json:"grade"`
	Recommendation string             `json:"recommendation"`
	LimitSuggested int64              `json:"limit_suggested"`
	Features       map[string]float64 `json:"features"`
	Reasons        []string           `json:"reasons"`
	Source         string             `json:"source"`
}

// ---- Loan Application ----

type CreateLoanApplicationRequest struct {
	MemberID    string `json:"member_id" validate:"required,uuid"`
	Amount      int64  `json:"amount" validate:"required,gt=0"`
	TenorMonths int    `json:"tenor_months" validate:"required,gt=0"`
	Purpose     string `json:"purpose"`
}

type CreatePortalLoanApplicationRequest struct {
	Amount      int64  `json:"amount" validate:"required,gt=0"`
	TenorMonths int    `json:"tenor_months" validate:"required,gt=0"`
	Purpose     string `json:"purpose"`
}

type RejectLoanApplicationRequest struct {
	Reason string `json:"reason"`
}

type LoanApplicationResponse struct {
	ID             string                    `json:"id"`
	MemberID       string                    `json:"member_id"`
	MemberName     string                    `json:"member_name"`
	Amount         int64                     `json:"amount"`
	TenorMonths    int                       `json:"tenor_months"`
	Purpose        string                    `json:"purpose"`
	Status         string                    `json:"status"`
	ApprovedBy     *string                   `json:"approved_by"`
	ApprovedByName *string                   `json:"approved_by_name"`
	ApprovedAt     *string                   `json:"approved_at"`
	CreatedAt      string                    `json:"created_at"`
	Assessment     *CreditAssessmentResponse `json:"assessment"`
}

// ---- Installment ----

type InstallmentResponse struct {
	ID           string `json:"id"`
	LoanID       string `json:"loan_id"`
	PeriodNo     int    `json:"period_no"`
	DueDate      string `json:"due_date"`
	PrincipalDue int64  `json:"principal_due"`
	InterestDue  int64  `json:"interest_due"`
	TotalDue     int64  `json:"total_due"`
	PaidAmount   int64  `json:"paid_amount"`
	Status       string `json:"status"`
}

// ---- Loan ----

type LoanResponse struct {
	ID              string                `json:"id"`
	ApplicationID   string                `json:"application_id"`
	MemberID        string                `json:"member_id"`
	MemberName      string                `json:"member_name"`
	Principal       int64                 `json:"principal"`
	FlatRateMonthly float64               `json:"flat_rate_monthly"`
	TenorMonths     int                   `json:"tenor_months"`
	Status          string                `json:"status"`
	DisbursedAt     string                `json:"disbursed_at"`
	Outstanding     int64                 `json:"outstanding"`
	Schedule        []InstallmentResponse `json:"schedule,omitempty"`
}

// ---- Payment ----

type PaymentResponse struct {
	ID         string `json:"id"`
	ScheduleID string `json:"schedule_id"`
	Amount     int64  `json:"amount"`
	Penalty    int64  `json:"penalty"`
	PaidAt     string `json:"paid_at"`
}

type PayInstallmentRequest struct {
	Amount  int64 `json:"amount" validate:"required,gt=0"`
	Penalty int64 `json:"penalty" validate:"gte=0"`
}

type PayInstallmentResponse struct {
	Payment     PaymentResponse     `json:"payment"`
	Installment InstallmentResponse `json:"installment"`
	Loan        LoanResponse        `json:"loan"`
}

type ApproveApplicationResponse struct {
	Application LoanApplicationResponse `json:"application"`
	Loan        LoanResponse            `json:"loan"`
}

// ---- Credit Scoring Integration Endpoint ----

type CreditScoringRequest struct {
	MemberID       string             `json:"member_id"`
	Features       map[string]float64 `json:"features"`
	JumlahDiajukan int64              `json:"jumlah_diajukan"`
	TenorBulan     int                `json:"tenor_bulan"`
	TotalSimpanan  int64              `json:"total_simpanan"`
	MaxPlafond     int64              `json:"max_plafond"`
}

type CreditScoringResponse struct {
	Score            int      `json:"score"`
	Grade            string   `json:"grade"`
	Recommendation   string   `json:"recommendation"`
	LimitRekomendasi int64    `json:"limit_rekomendasi"`
	Reasons          []string `json:"reasons"`
	Source           string   `json:"source"`
}

// ---- Loan Audit ----

type LoanAuditLogResponse struct {
	ID               string `json:"id"`
	LoanID           string `json:"loan_id"`
	Action           string `json:"action"`
	PerformedByEmail string `json:"performed_by_email"`
	PerformedAt      string `json:"performed_at"`
	BeforeData       string `json:"before_data"` // JSON string
	AfterData        string `json:"after_data"`  // JSON string
	Note             string `json:"note"`
	IsFlagged        bool   `json:"is_flagged"`
	FlaggedByName    string `json:"flagged_by_name,omitempty"`
	FlaggedAt        string `json:"flagged_at,omitempty"`
	FlaggedReason    string `json:"flagged_reason,omitempty"`
}

type CreateAuditTokenRequest struct {
	ExpiresInHours int `json:"expires_in_hours" validate:"required,gt=0"`
}

type CreateAuditTokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
	AuditURL  string `json:"audit_url"`
}

type FlagAuditLogRequest struct {
	AuditLogID    string `json:"audit_log_id" validate:"required,uuid"`
	FlaggedByName string `json:"flagged_by_name" validate:"required"`
	FlaggedReason string `json:"flagged_reason" validate:"required"`
}

type LoanAuditDetailResponse struct {
	Loan      LoanResponse           `json:"loan"`
	AuditLogs []LoanAuditLogResponse `json:"audit_logs"`
}
