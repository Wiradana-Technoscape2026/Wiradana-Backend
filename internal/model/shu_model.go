package model

type ShuPeriodResponse struct {
	ID           string  `json:"id"`
	Year         int     `json:"year"`
	TotalShu     int64   `json:"total_shu"`
	PctJasaModal float64 `json:"pct_jasa_modal"`
	PctJasaUsaha float64 `json:"pct_jasa_usaha"`
	Status       string  `json:"status"`
}

type ShuDistributionResponse struct {
	ID          string `json:"id"`
	ShuPeriodID string `json:"shu_period_id"`
	MemberID    string `json:"member_id"`
	JasaModal   int64  `json:"jasa_modal"`
	JasaUsaha   int64  `json:"jasa_usaha"`
	TotalShu    int64  `json:"total_shu"`
}
