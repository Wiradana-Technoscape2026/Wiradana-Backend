package model

type KTPOCRResponse struct {
	NIK        string  `json:"nik"`
	FullName   string  `json:"full_name"`
	Address    string  `json:"address"`
	BirthDate  string  `json:"birth_date"`
	Confidence float64 `json:"confidence"`
	Source     string  `json:"source"`
}
