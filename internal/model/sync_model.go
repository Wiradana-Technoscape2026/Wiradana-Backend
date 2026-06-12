package model

import "time"

type SyncPullResponse struct {
	Cursor               time.Time                 `json:"cursor"`
	Members              []MemberResponse          `json:"members"`
	SavingsTransactions  []SavingsTransactionResponse `json:"savings_transactions"`
	LoanApplications     []LoanApplicationResponse `json:"loan_applications"`
	Loans                []LoanResponse            `json:"loans"`
	InstallmentSchedules []InstallmentResponse     `json:"installment_schedules"`
	Payments             []PaymentResponse         `json:"payments"`
	ShuPeriods           []ShuPeriodResponse       `json:"shu_periods"`
	ShuDistributions     []ShuDistributionResponse `json:"shu_distributions"`
	LoanConfig           *LoanConfigResponse       `json:"loan_config"`
}

type MutationRequest struct {
	ID              string                 `json:"id" validate:"required,uuid"`
	Type            string                 `json:"type" validate:"required"`
	Payload         map[string]interface{} `json:"payload" validate:"required"`
	ClientTimestamp *time.Time             `json:"client_timestamp"`
}

type SyncPushRequest struct {
	Mutations []MutationRequest `json:"mutations" validate:"required,min=1,dive"`
}

type MutationResult struct {
	ID       string  `json:"id"`
	Status   string  `json:"status"`
	ResultID *string `json:"result_id,omitempty"`
	Error    *string `json:"error,omitempty"`
}

type SyncPushResponse struct {
	Results []MutationResult `json:"results"`
}
