package model

// CreateShuPeriodRequest — POST /shu-periods (api_planning §3.7)
type CreateShuPeriodRequest struct {
	Year         int   `json:"year"          validate:"required,gt=2000"`
	TotalShu     int64 `json:"total_shu"     validate:"required,gt=0"`
	PctJasaModal float64 `json:"pct_jasa_modal" validate:"required,gte=0,lte=100"`
	PctJasaUsaha float64 `json:"pct_jasa_usaha" validate:"required,gte=0,lte=100"`
}

// ShuPeriodResponse — api_planning §2.10
type ShuPeriodResponse struct {
	ID            string    `json:"id"`
	Year          int       `json:"year"`
	TotalShu      int64     `json:"total_shu"`
	PctJasaModal  float64   `json:"pct_jasa_modal"`
	PctJasaUsaha  float64   `json:"pct_jasa_usaha"`
	Status        string    `json:"status"`
}

// ShuDistributionResponse — api_planning §2.10
type ShuDistributionResponse struct {
	ID          string `json:"id"`
	MemberID    string `json:"member_id"`
	MemberName  string `json:"member_name"`
	JasaModal   int64  `json:"jasa_modal"`
	JasaUsaha   int64  `json:"jasa_usaha"`
	TotalShu    int64  `json:"total_shu"`
}

// CalculateShuResponse — POST /shu-periods/:id/calculate
type CalculateShuResponse struct {
	Period        ShuPeriodResponse        `json:"period"`
	Distributions []ShuDistributionResponse `json:"distributions"`
}

// PortalSHUResponse — GET /portal/shu
type PortalSHUResponse struct {
	EstimatedShu int64                     `json:"estimated_shu"`
	History      []ShuDistributionResponse `json:"history"`
}
