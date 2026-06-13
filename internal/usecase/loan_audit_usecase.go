package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"github.com/wiradana/backend/internal/model/converter"
	"github.com/wiradana/backend/internal/repository"
)

var (
	ErrInvalidOrExpiredToken = errors.New("token audit tidak valid atau sudah kedaluwarsa")
)

type LoanAuditUsecase interface {
	CreateToken(ctx context.Context, coopID, loanID, userID string, req *model.CreateAuditTokenRequest) (*model.CreateAuditTokenResponse, error)
	GetAuditDetails(ctx context.Context, tokenRaw string) (*model.LoanAuditDetailResponse, error)
	FlagLog(ctx context.Context, tokenRaw string, req *model.FlagAuditLogRequest) error
}

type loanAuditUsecase struct {
	auditRepo repository.LoanAuditRepository
	loanRepo  repository.LoanRepository
}

func NewLoanAuditUsecase(auditRepo repository.LoanAuditRepository, loanRepo repository.LoanRepository) LoanAuditUsecase {
	return &loanAuditUsecase{
		auditRepo: auditRepo,
		loanRepo:  loanRepo,
	}
}

func hashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func (u *loanAuditUsecase) CreateToken(ctx context.Context, coopID, loanID, userID string, req *model.CreateAuditTokenRequest) (*model.CreateAuditTokenResponse, error) {
	coopUUID, err := uuid.Parse(coopID)
	if err != nil {
		return nil, errors.New("cooperative_id tidak valid")
	}
	loanUUID, err := uuid.Parse(loanID)
	if err != nil {
		return nil, errors.New("loan_id tidak valid")
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("user_id tidak valid")
	}

	// Verify loan exists in cooperative
	_, err = u.loanRepo.FindByID(ctx, coopID, loanID)
	if err != nil {
		return nil, errors.New("loan tidak ditemukan di koperasi ini")
	}

	// Generate secure token
	rawToken := uuid.New().String() + "-" + uuid.New().String()
	tokenHash := hashToken(rawToken)

	expiresAt := time.Now().Add(time.Duration(req.ExpiresInHours) * time.Hour)

	tokenEntity := &entity.LoanAuditToken{
		CooperativeID: coopUUID,
		LoanID:        loanUUID,
		TokenHash:     tokenHash,
		ExpiresAt:     expiresAt,
		CreatedBy:     userUUID,
	}

	if err := u.auditRepo.CreateToken(ctx, tokenEntity); err != nil {
		return nil, err
	}

	auditURL := fmt.Sprintf("/audit/loans?token=%s", rawToken)

	return &model.CreateAuditTokenResponse{
		Token:     rawToken,
		ExpiresAt: expiresAt.Format(time.RFC3339),
		AuditURL:  auditURL,
	}, nil
}

func (u *loanAuditUsecase) GetAuditDetails(ctx context.Context, tokenRaw string) (*model.LoanAuditDetailResponse, error) {
	tokenHash := hashToken(tokenRaw)
	token, err := u.auditRepo.FindTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, ErrInvalidOrExpiredToken
	}

	// Fetch loan
	loanMeta, err := u.loanRepo.FindByID(ctx, "", token.LoanID.String())
	if err != nil {
		return nil, errors.New("data pinjaman tidak ditemukan")
	}

	// Mask PII
	maskedMemberName := maskMemberName(loanMeta.MemberName)
	maskedMemberID := maskString(loanMeta.MemberID.String(), 4, 4)

	// Fetch schedules and convert
	instResponses := make([]model.InstallmentResponse, len(loanMeta.Schedule))
	for i, s := range loanMeta.Schedule {
		instResponses[i] = converter.ToInstallmentResponse(&s.InstallmentSchedule, s.PaidAmount)
	}

	loanResp := converter.ToLoanResponse(&loanMeta.Loan, maskedMemberName, loanMeta.Outstanding, instResponses)
	loanResp.MemberID = maskedMemberID // Mask UUID

	// Fetch logs
	logs, err := u.auditRepo.FindLogsByLoanID(ctx, token.LoanID.String())
	if err != nil {
		return nil, err
	}

	auditLogResponses := make([]model.LoanAuditLogResponse, len(logs))
	for i, l := range logs {
		// Mask raw before/after JSON PII if they contain member details
		beforeData := maskJSONPII(string(l.BeforeData))
		afterData := maskJSONPII(string(l.AfterData))

		respLog := converter.ToLoanAuditLogResponse(&l.LoanAuditLog, l.PerformedByEmail)
		respLog.BeforeData = beforeData
		respLog.AfterData = afterData
		auditLogResponses[i] = respLog
	}

	return &model.LoanAuditDetailResponse{
		Loan:      loanResp,
		AuditLogs: auditLogResponses,
	}, nil
}

func (u *loanAuditUsecase) FlagLog(ctx context.Context, tokenRaw string, req *model.FlagAuditLogRequest) error {
	tokenHash := hashToken(tokenRaw)
	token, err := u.auditRepo.FindTokenByHash(ctx, tokenHash)
	if err != nil {
		return ErrInvalidOrExpiredToken
	}

	// Ensure the log entry belongs to the loan authorized by the token
	// Let's verify this in the repository or fetch the log first
	// We can check it here
	logs, err := u.auditRepo.FindLogsByLoanID(ctx, token.LoanID.String())
	if err != nil {
		return err
	}
	found := false
	for _, l := range logs {
		if l.ID.String() == req.AuditLogID {
			found = true
			break
		}
	}
	if !found {
		return errors.New("audit log tidak valid untuk pinjaman ini")
	}

	return u.auditRepo.FlagLog(ctx, req.AuditLogID, req.FlaggedByName, req.FlaggedReason)
}

// ---- Masking Helpers ----

func maskString(s string, visibleStart, visibleEnd int) string {
	if len(s) <= visibleStart+visibleEnd {
		return s
	}
	maskedLen := len(s) - visibleStart - visibleEnd
	return s[:visibleStart] + strings.Repeat("*", maskedLen) + s[len(s)-visibleEnd:]
}

func maskMemberName(name string) string {
	parts := strings.Split(name, " ")
	for i, part := range parts {
		if len(part) > 2 {
			parts[i] = part[:1] + strings.Repeat("*", len(part)-2) + part[len(part)-1:]
		} else if len(part) > 0 {
			parts[i] = part[:1] + "*"
		}
	}
	return strings.Join(parts, " ")
}

// maskJSONPII attempts to mask fields like "nik", "full_name", "address", "phone_number" inside json dumps if they exist.
func maskJSONPII(jsonStr string) string {
	if jsonStr == "" || jsonStr == "{}" {
		return jsonStr
	}
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return jsonStr
	}

	maskFields := func(m map[string]interface{}) {
		if val, ok := m["full_name"]; ok {
			if s, ok := val.(string); ok {
				m["full_name"] = maskMemberName(s)
			}
		}
		if val, ok := m["nik"]; ok {
			if s, ok := val.(string); ok {
				m["nik"] = maskString(s, 4, 2)
			}
		}
		if val, ok := m["phone_number"]; ok {
			if s, ok := val.(string); ok {
				m["phone_number"] = maskString(s, 4, 3)
			}
		}
		if val, ok := m["address"]; ok {
			if s, ok := val.(string); ok {
				m["address"] = maskAddress(s)
			}
		}
	}

	maskFields(data)

	// Recursively look for nested members or other sub-structures
	for k, v := range data {
		if subMap, ok := v.(map[string]interface{}); ok {
			maskFields(subMap)
			data[k] = subMap
		}
	}

	res, err := json.Marshal(data)
	if err == nil {
		return string(res)
	}
	return jsonStr
}

func maskAddress(addr string) string {
	parts := strings.Split(addr, ",")
	if len(parts) > 1 {
		return "..., " + strings.TrimSpace(parts[len(parts)-1])
	}
	if len(addr) > 5 {
		return "..." + addr[len(addr)-5:]
	}
	return "..."
}
