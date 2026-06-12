package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wiradana/backend/internal/entity"
	gw "github.com/wiradana/backend/internal/gateway/notification"
	"github.com/wiradana/backend/internal/repository"
)

type NotificationLogResponse struct {
	ID          string    `json:"id"`
	MemberName  string    `json:"member_name"`
	PhoneNumber string    `json:"phone_number"`
	Channel     string    `json:"channel"`
	EventType   string    `json:"event_type"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
	Source      string    `json:"source"`
	CreatedAt   time.Time `json:"created_at"`
}

type NotificationUsecase interface {
	SendDueReminders(ctx context.Context, daysAhead int) error
	SendDueTodayAlerts(ctx context.Context) error
	SendOverdueAlerts(ctx context.Context) error
	SendSavingsConfirmation(ctx context.Context, coopID, memberID, txID string, amount int64, savingsType string) error
	SendLoanDisbursed(ctx context.Context, coopID, memberID, loanID string, principal int64, firstDueDate time.Time) error
	ListLogs(ctx context.Context, coopID, eventType string, page, limit int) ([]NotificationLogResponse, int64, error)
}

type notificationUsecase struct {
	repo       repository.NotificationRepository
	memberRepo repository.MemberRepository
	gateway    gw.Gateway
}

func NewNotificationUsecase(repo repository.NotificationRepository, memberRepo repository.MemberRepository, gateway gw.Gateway) NotificationUsecase {
	return &notificationUsecase{repo: repo, memberRepo: memberRepo, gateway: gateway}
}

func formatRp(amount int64) string {
	s := fmt.Sprintf("%d", amount)
	n := len(s)
	result := make([]byte, n+((n-1)/3))
	j := len(result) - 1
	for i := n - 1; i >= 0; i-- {
		result[j] = s[i]
		j--
		pos := n - i
		if pos%3 == 0 && i > 0 {
			result[j] = '.'
			j--
		}
	}
	return "Rp" + string(result)
}

func (u *notificationUsecase) sendAndLog(ctx context.Context, log *entity.NotificationLog) {
	result, err := u.gateway.Send(ctx, gw.Input{
		ToPhone: log.PhoneNumber,
		Message: log.Message,
	})
	if err != nil {
		log.Status = "failed"
		errMsg := err.Error()
		log.ErrorMsg = &errMsg
	} else {
		log.Status = "sent"
		log.Source = result.Source
	}
	_ = u.repo.CreateLog(ctx, log)
}

func (u *notificationUsecase) SendDueReminders(ctx context.Context, daysAhead int) error {
	wib := time.FixedZone("WIB", 7*3600)
	targetDate := time.Now().In(wib).AddDate(0, 0, daysAhead)

	rows, err := u.repo.FindDueInstallments(ctx, targetDate)
	if err != nil {
		return err
	}

	var eventType string
	var msgTmpl string
	switch daysAhead {
	case 3:
		eventType = "reminder_3d"
		msgTmpl = "Yth. %s, angsuran %s jatuh tempo dalam 3 hari (%s). Segera bersiap untuk melakukan pembayaran. - Koperasi Wiradana"
	case 1:
		eventType = "reminder_1d"
		msgTmpl = "Yth. %s, angsuran %s jatuh tempo BESOK (%s). Segera bayar untuk menghindari keterlambatan. - Koperasi Wiradana"
	default:
		eventType = fmt.Sprintf("reminder_%dd", daysAhead)
		msgTmpl = "Yth. %s, angsuran %s jatuh tempo dalam %d hari (%s). - Koperasi Wiradana"
	}

	for _, row := range rows {
		instID := row.InstallmentID
		msg := fmt.Sprintf(msgTmpl, row.MemberName, formatRp(row.TotalDue), row.DueDate.Format("02 Jan 2006"))
		log := &entity.NotificationLog{
			CooperativeID: row.CooperativeID,
			MemberID:      row.MemberID,
			PhoneNumber:   row.PhoneNumber,
			EventType:     eventType,
			RefID:         &instID,
			Message:       msg,
			Source:        u.gateway.Source(),
		}
		u.sendAndLog(ctx, log)
	}
	return nil
}

func (u *notificationUsecase) SendDueTodayAlerts(ctx context.Context) error {
	wib := time.FixedZone("WIB", 7*3600)
	today := time.Now().In(wib)

	rows, err := u.repo.FindDueInstallments(ctx, today)
	if err != nil {
		return err
	}

	for _, row := range rows {
		instID := row.InstallmentID
		msg := fmt.Sprintf(
			"Yth. %s, angsuran %s JATUH TEMPO HARI INI (%s). Bayar sekarang untuk menghindari denda. - Koperasi Wiradana",
			row.MemberName, formatRp(row.TotalDue), row.DueDate.Format("02 Jan 2006"),
		)
		log := &entity.NotificationLog{
			CooperativeID: row.CooperativeID,
			MemberID:      row.MemberID,
			PhoneNumber:   row.PhoneNumber,
			EventType:     "due_today",
			RefID:         &instID,
			Message:       msg,
			Source:        u.gateway.Source(),
		}
		u.sendAndLog(ctx, log)
	}
	return nil
}

func (u *notificationUsecase) SendOverdueAlerts(ctx context.Context) error {
	rows, err := u.repo.FindOverdueInstallments(ctx)
	if err != nil {
		return err
	}

	for _, row := range rows {
		instID := row.InstallmentID
		msg := fmt.Sprintf(
			"Yth. %s, angsuran %s telah MELEWATI jatuh tempo (%s). Segera hubungi koperasi untuk menghindari penalti lebih lanjut. - Koperasi Wiradana",
			row.MemberName, formatRp(row.TotalDue), row.DueDate.Format("02 Jan 2006"),
		)
		log := &entity.NotificationLog{
			CooperativeID: row.CooperativeID,
			MemberID:      row.MemberID,
			PhoneNumber:   row.PhoneNumber,
			EventType:     "overdue",
			RefID:         &instID,
			Message:       msg,
			Source:        u.gateway.Source(),
		}
		u.sendAndLog(ctx, log)
	}
	return nil
}

func (u *notificationUsecase) SendSavingsConfirmation(ctx context.Context, coopID, memberID, txID string, amount int64, savingsType string) error {
	member, err := u.memberRepo.FindByID(ctx, coopID, memberID)
	if err != nil || member.PhoneNumber == nil || *member.PhoneNumber == "" {
		return nil
	}

	coopUUID, _ := uuid.Parse(coopID)
	txUUID, _ := uuid.Parse(txID)
	wib := time.FixedZone("WIB", 7*3600)
	dateStr := time.Now().In(wib).Format("02 Jan 2006")

	msg := fmt.Sprintf(
		"Yth. %s, setoran simpanan %s %s berhasil dicatat pada %s. Terima kasih atas kepercayaan Anda. - Koperasi Wiradana",
		member.FullName, savingsType, formatRp(amount), dateStr,
	)
	notifLog := &entity.NotificationLog{
		CooperativeID: coopUUID,
		MemberID:      member.ID,
		PhoneNumber:   *member.PhoneNumber,
		EventType:     "savings_confirmed",
		RefID:         &txUUID,
		Message:       msg,
		Source:        u.gateway.Source(),
	}
	u.sendAndLog(ctx, notifLog)
	return nil
}

func (u *notificationUsecase) SendLoanDisbursed(ctx context.Context, coopID, memberID, loanID string, principal int64, firstDueDate time.Time) error {
	member, err := u.memberRepo.FindByID(ctx, coopID, memberID)
	if err != nil || member.PhoneNumber == nil || *member.PhoneNumber == "" {
		return nil
	}

	coopUUID, _ := uuid.Parse(coopID)
	loanUUID, _ := uuid.Parse(loanID)

	msg := fmt.Sprintf(
		"Yth. %s, pinjaman %s telah disetujui dan dicairkan. Angsuran perdana jatuh tempo pada %s. Gunakan dengan bijak. - Koperasi Wiradana",
		member.FullName, formatRp(principal), firstDueDate.Format("02 Jan 2006"),
	)
	notifLog := &entity.NotificationLog{
		CooperativeID: coopUUID,
		MemberID:      member.ID,
		PhoneNumber:   *member.PhoneNumber,
		EventType:     "loan_disbursed",
		RefID:         &loanUUID,
		Message:       msg,
		Source:        u.gateway.Source(),
	}
	u.sendAndLog(ctx, notifLog)
	return nil
}

func (u *notificationUsecase) ListLogs(ctx context.Context, coopID, eventType string, page, limit int) ([]NotificationLogResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	rows, total, err := u.repo.ListLogs(ctx, coopID, eventType, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]NotificationLogResponse, len(rows))
	for i, r := range rows {
		result[i] = NotificationLogResponse{
			ID:          r.ID.String(),
			MemberName:  r.MemberName,
			PhoneNumber: r.PhoneNumber,
			Channel:     r.Channel,
			EventType:   r.EventType,
			Message:     r.Message,
			Status:      r.Status,
			Source:      r.Source,
			CreatedAt:   r.CreatedAt,
		}
	}
	return result, total, nil
}
