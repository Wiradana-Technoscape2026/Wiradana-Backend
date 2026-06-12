package converter

import (
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
	"gorm.io/datatypes"
)

func ToMemberResponse(m *entity.Member, summary *model.SavingsSummary) model.MemberResponse {
	birthDate := ""
	if m.BirthDate != nil {
		birthDate = m.BirthDate.Format("2006-01-02")
	}

	address := ""
	if m.Address != nil {
		address = *m.Address
	}

	attrs := m.CustomAttributes
	if attrs == nil {
		attrs = datatypes.JSON("{}")
	}

	s := model.SavingsSummary{}
	if summary != nil {
		s = *summary
	}
	s.Total = s.Pokok + s.Wajib + s.Sukarela

	return model.MemberResponse{
		ID:               m.ID.String(),
		NIK:              m.NIK,
		FullName:         m.FullName,
		Address:          address,
		BirthDate:        birthDate,
		Status:           m.Status,
		CustomAttributes: attrs,
		JoinedAt:         m.JoinedAt,
		SavingsSummary:   s,
	}
}
