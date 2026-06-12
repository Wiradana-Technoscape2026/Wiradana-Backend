package converter

import (
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
)

func ToShuPeriodResponse(sp *entity.ShuPeriod) model.ShuPeriodResponse {
	return model.ShuPeriodResponse{
<<<<<<< HEAD
		ID:            sp.ID.String(),
		Year:          sp.Year,
		TotalShu:      sp.TotalShu,
		PctJasaModal:  sp.PctJasaModal,
		PctJasaUsaha:  sp.PctJasaUsaha,
		Status:        sp.Status,
	}
}

func ToShuDistributionResponse(sd *entity.ShuDistribution, memberName string) model.ShuDistributionResponse {
	return model.ShuDistributionResponse{
		ID:          sd.ID.String(),
		MemberID:    sd.MemberID.String(),
		MemberName:  memberName,
=======
		ID:           sp.ID.String(),
		Year:         sp.Year,
		TotalShu:     sp.TotalShu,
		PctJasaModal: sp.PctJasaModal,
		PctJasaUsaha: sp.PctJasaUsaha,
		Status:       sp.Status,
	}
}

func ToShuDistributionResponse(sd *entity.ShuDistribution) model.ShuDistributionResponse {
	return model.ShuDistributionResponse{
		ID:          sd.ID.String(),
		ShuPeriodID: sd.ShuPeriodID.String(),
		MemberID:    sd.MemberID.String(),
>>>>>>> e6c7f422c936b4876b95b9366e0dc7eebfff82ed
		JasaModal:   sd.JasaModal,
		JasaUsaha:   sd.JasaUsaha,
		TotalShu:    sd.TotalShu,
	}
}
