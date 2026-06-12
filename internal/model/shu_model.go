package model

type ShuPeriodResponse struct {
	ID           string  `json:"id"`
	Year         int     `json:"year"`
	TotalShu     int64   `json:"total_shu"`
	PctJasaModal float64 `json:"pct_jasa_modal"`
	PctJasaUsaha float64 `json:"pct_jasa_usaha"`
	Status       string  `json:"status"`
}

type CreateShuPeriodRequest struct {
	Year         int     `json:"year" validate:"required,gte=2000"`
	TotalShu     int64   `json:"total_shu" validate:"required,gt=0"`
	PctJasaModal float64 `json:"pct_jasa_modal" validate:"required,gte=0,lte=100"`
	PctJasaUsaha float64 `json:"pct_jasa_usaha" validate:"required,gte=0,lte=100"`
}

type ShuDistributionResponse struct {
	ID          string `json:"id"`
	ShuPeriodID string `json:"shu_period_id"`
	MemberID    string `json:"member_id"`
	MemberName  string `json:"member_name"`
	JasaModal   int64  `json:"jasa_modal"`
	JasaUsaha   int64  `json:"jasa_usaha"`
	TotalShu    int64  `json:"total_shu"`
}

type ShuDistributionDetail struct {
	ID          string `json:"id"`
	ShuPeriodID string `json:"shu_period_id"`
	MemberID    string `json:"member_id"`
	MemberName  string `json:"member_name"`
	JasaModal   int64  `json:"jasa_modal"`
	JasaUsaha   int64  `json:"jasa_usaha"`
	TotalShu    int64  `json:"total_shu"`
}

type CalculateShuResponse struct {
	Period        ShuPeriodResponse        `json:"period"`
	Distributions []ShuDistributionDetail  `json:"distributions"`
}
