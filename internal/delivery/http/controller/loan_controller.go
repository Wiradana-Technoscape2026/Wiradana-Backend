package controller

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/wiradana/backend/internal/usecase"
	"github.com/xuri/excelize/v2"
)

type LoanController struct {
	uc usecase.LoanUsecase
}

func NewLoanController(uc usecase.LoanUsecase) *LoanController {
	return &LoanController{uc: uc}
}

func (ctrl *LoanController) List(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	status := c.Query("status")
	loans, err := ctrl.uc.List(c.Context(), coopID, status)
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", err.Error())
	}
	return OKList(c, loans, int64(len(loans)))
}

func (ctrl *LoanController) GetByID(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	loanID := c.Params("id")
	loan, err := ctrl.uc.GetByID(c.Context(), coopID, loanID)
	if err != nil {
		return Fail(c, fiber.StatusNotFound, "NOT_FOUND", err.Error())
	}
	return OK(c, loan)
}

func (ctrl *LoanController) ExportExcel(c *fiber.Ctx) error {
	coopID := c.Locals("cooperative_id").(string)
	loanID := c.Params("id")

	loan, err := ctrl.uc.GetByID(c.Context(), coopID, loanID)
	if err != nil {
		return Fail(c, fiber.StatusNotFound, "NOT_FOUND", err.Error())
	}

	f := excelize.NewFile()
	defer f.Close()

	const sh = "Detail Pinjaman"
	f.SetSheetName("Sheet1", sh)

	// ── Column widths ──────────────────────────────────────────────────────────
	f.SetColWidth(sh, "A", "A", 14)
	f.SetColWidth(sh, "B", "B", 16)
	f.SetColWidth(sh, "C", "C", 20)
	f.SetColWidth(sh, "D", "D", 20)
	f.SetColWidth(sh, "E", "E", 20)
	f.SetColWidth(sh, "F", "F", 14)

	// ── Styles ────────────────────────────────────────────────────────────────
	rpFmt := `"Rp "#,##0`
	rateFmt := `0.0"% / bln"`

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14, Color: "111827"},
	})
	infoStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 10, Color: "6B7280", Italic: true},
	})
	summaryLabelStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 9, Color: "1E40AF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"EFF6FF"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "top", Color: "BFDBFE", Style: 1},
			{Type: "left", Color: "BFDBFE", Style: 1},
			{Type: "right", Color: "BFDBFE", Style: 1},
		},
	})
	summaryRpStyle, _ := f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Bold: true, Size: 12, Color: "1D4ED8"},
		CustomNumFmt: &rpFmt,
		Fill:         excelize.Fill{Type: "pattern", Color: []string{"EFF6FF"}, Pattern: 1},
		Alignment:    &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "BFDBFE", Style: 1},
			{Type: "left", Color: "BFDBFE", Style: 1},
			{Type: "right", Color: "BFDBFE", Style: 1},
		},
	})
	summaryRateStyle, _ := f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Bold: true, Size: 12, Color: "1D4ED8"},
		CustomNumFmt: &rateFmt,
		Fill:         excelize.Fill{Type: "pattern", Color: []string{"EFF6FF"}, Pattern: 1},
		Alignment:    &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "BFDBFE", Style: 1},
			{Type: "left", Color: "BFDBFE", Style: 1},
			{Type: "right", Color: "BFDBFE", Style: 1},
		},
	})
	sectionStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "111827"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	tblHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"1A56DB"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "FFFFFF", Style: 1},
			{Type: "right", Color: "FFFFFF", Style: 1},
		},
	})
	cellCenterStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    []excelize.Border{{Type: "bottom", Color: "E5E7EB", Style: 1}},
	})
	rpCellStyle, _ := f.NewStyle(&excelize.Style{
		CustomNumFmt: &rpFmt,
		Alignment:    &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		Border:       []excelize.Border{{Type: "bottom", Color: "E5E7EB", Style: 1}},
	})
	lunasStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "059669"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    []excelize.Border{{Type: "bottom", Color: "E5E7EB", Style: 1}},
	})
	belumBayarStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Color: "D97706"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    []excelize.Border{{Type: "bottom", Color: "E5E7EB", Style: 1}},
	})
	terlambatStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "DC2626"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    []excelize.Border{{Type: "bottom", Color: "E5E7EB", Style: 1}},
	})

	// ── Row 1: Title ──────────────────────────────────────────────────────────
	f.MergeCell(sh, "A1", "F1")
	f.SetCellValue(sh, "A1", "Detail Pinjaman — "+loan.MemberName)
	f.SetCellStyle(sh, "A1", "F1", titleStyle)
	f.SetRowHeight(sh, 1, 30)

	// ── Row 2: Info ───────────────────────────────────────────────────────────
	paidCount := 0
	for _, s := range loan.Schedule {
		if s.Status == "lunas" {
			paidCount++
		}
	}
	f.MergeCell(sh, "A2", "F2")
	f.SetCellValue(sh, "A2", fmt.Sprintf(
		"Cair %s  ·  %d bln  ·  %d/%d terbayar  ·  Status: %s",
		loan.DisbursedAt, loan.TenorMonths, paidCount, loan.TenorMonths,
		strings.ToUpper(loan.Status[:1])+loan.Status[1:],
	))
	f.SetCellStyle(sh, "A2", "F2", infoStyle)
	f.SetRowHeight(sh, 2, 18)

	// ── Row 3: spacer ─────────────────────────────────────────────────────────
	f.SetRowHeight(sh, 3, 8)

	// ── Rows 4–5: Summary cards ───────────────────────────────────────────────
	// Row 4: labels, Row 5: values — columns A, B, C, D (one per card)
	var installmentPerMonth int64
	if len(loan.Schedule) > 0 {
		installmentPerMonth = loan.Schedule[0].TotalDue
	}

	summaryLabels := []string{"Pokok Pinjaman", "Bunga Flat", "Outstanding", "Angsuran / Bln"}
	summaryCols := []string{"A", "B", "C", "D"}

	for i, col := range summaryCols {
		f.SetCellValue(sh, col+"4", summaryLabels[i])
		f.SetCellStyle(sh, col+"4", col+"4", summaryLabelStyle)
	}
	f.SetRowHeight(sh, 4, 22)

	// Values
	f.SetCellValue(sh, "A5", loan.Principal)
	f.SetCellStyle(sh, "A5", "A5", summaryRpStyle)
	f.SetCellValue(sh, "B5", loan.FlatRateMonthly)
	f.SetCellStyle(sh, "B5", "B5", summaryRateStyle)
	f.SetCellValue(sh, "C5", loan.Outstanding)
	f.SetCellStyle(sh, "C5", "C5", summaryRpStyle)
	f.SetCellValue(sh, "D5", installmentPerMonth)
	f.SetCellStyle(sh, "D5", "D5", summaryRpStyle)
	f.SetRowHeight(sh, 5, 26)

	// ── Row 6: spacer ─────────────────────────────────────────────────────────
	f.SetRowHeight(sh, 6, 8)

	// ── Row 7: section header ─────────────────────────────────────────────────
	f.MergeCell(sh, "A7", "F7")
	f.SetCellValue(sh, "A7", "Jadwal Angsuran")
	f.SetCellStyle(sh, "A7", "F7", sectionStyle)
	f.SetRowHeight(sh, 7, 22)

	// ── Row 8: table headers ──────────────────────────────────────────────────
	tblHeaders := []string{"Periode", "Jatuh Tempo", "Pokok", "Bunga", "Total", "Status"}
	tblCols := []string{"A", "B", "C", "D", "E", "F"}
	for i, col := range tblCols {
		f.SetCellValue(sh, col+"8", tblHeaders[i])
		f.SetCellStyle(sh, col+"8", col+"8", tblHeaderStyle)
	}
	f.SetRowHeight(sh, 8, 22)

	// ── Rows 9+: schedule ─────────────────────────────────────────────────────
	for i, inst := range loan.Schedule {
		row := 9 + i
		rowStr := fmt.Sprintf("%d", row)

		f.SetCellValue(sh, "A"+rowStr, fmt.Sprintf("%02d", inst.PeriodNo))
		f.SetCellStyle(sh, "A"+rowStr, "A"+rowStr, cellCenterStyle)

		f.SetCellValue(sh, "B"+rowStr, inst.DueDate)
		f.SetCellStyle(sh, "B"+rowStr, "B"+rowStr, cellCenterStyle)

		f.SetCellValue(sh, "C"+rowStr, inst.PrincipalDue)
		f.SetCellStyle(sh, "C"+rowStr, "C"+rowStr, rpCellStyle)

		f.SetCellValue(sh, "D"+rowStr, inst.InterestDue)
		f.SetCellStyle(sh, "D"+rowStr, "D"+rowStr, rpCellStyle)

		f.SetCellValue(sh, "E"+rowStr, inst.TotalDue)
		f.SetCellStyle(sh, "E"+rowStr, "E"+rowStr, rpCellStyle)

		statusLabel := statusLabel(inst.Status)
		f.SetCellValue(sh, "F"+rowStr, statusLabel)
		switch inst.Status {
		case "lunas":
			f.SetCellStyle(sh, "F"+rowStr, "F"+rowStr, lunasStyle)
		case "terlambat":
			f.SetCellStyle(sh, "F"+rowStr, "F"+rowStr, terlambatStyle)
		default:
			f.SetCellStyle(sh, "F"+rowStr, "F"+rowStr, belumBayarStyle)
		}

		f.SetRowHeight(sh, row, 20)
	}

	// ── Stream response ───────────────────────────────────────────────────────
	buf, err := f.WriteToBuffer()
	if err != nil {
		return Fail(c, fiber.StatusInternalServerError, "INTERNAL", "gagal generate excel")
	}

	safeName := strings.ReplaceAll(loan.MemberName, " ", "_")
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="pinjaman_%s.xlsx"`, safeName))
	return c.Send(buf.Bytes())
}

func statusLabel(s string) string {
	switch s {
	case "lunas":
		return "Lunas"
	case "terlambat":
		return "Terlambat"
	default:
		return "Belum Bayar"
	}
}
