package model

type KTPOCRResponse struct {
	NIK        string  `json:"nik"`
	FullName   string  `json:"full_name"`
	Address    string  `json:"address"`
	Pekerjaan  string  `json:"pekerjaan"`
	BirthDate  string  `json:"birth_date"`
	Confidence float64 `json:"confidence"`
	Source     string  `json:"source"`
}
