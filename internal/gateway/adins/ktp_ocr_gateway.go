package adins

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"
	"time"
)

type apiCoIDGateway struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewAPICoIDGateway(apiKey, baseURL string) KTPOCRGateway {
	return &apiCoIDGateway{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

type apiCoIDResponse struct {
	IsSuccess bool           `json:"is_success"`
	Message   string         `json:"message"`
	Data      apiCoIDKTPData `json:"data"`
}

type apiCoIDKTPData struct {
	NIK              string `json:"nik"`
	Nama             string `json:"nama"`
	TempatLahir      string `json:"tempat_lahir"`
	TanggalLahir     string `json:"tanggal_lahir"`
	JenisKelamin     string `json:"jenis_kelamin"`
	Alamat           string `json:"alamat"`
	RTRW             string `json:"rt_rw"`
	Kelurahan        string `json:"kelurahan"`
	Kecamatan        string `json:"kecamatan"`
	Kota             string `json:"kota"`
	Provinsi         string `json:"provinsi"`
	Agama            string `json:"agama"`
	StatusPerkawinan string `json:"status_perkawinan"`
	Pekerjaan        string `json:"pekerjaan"`
	Kewarganegaraan  string `json:"kewarganegaraan"`
	BerlakuHingga    string `json:"berlaku_hingga"`
}

func (g *apiCoIDGateway) ExtractKTP(ctx context.Context, image []byte, filename string) (*KTPOCRResult, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := createImageFormFile(writer, "file", filename, image)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat multipart form: %w", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(image)); err != nil {
		return nil, fmt.Errorf("gagal menulis file ke form: %w", err)
	}
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/ocr/ktp-extract", body)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat request: %w", err)
	}
	req.Header.Set("x-api-co-id", g.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal menghubungi layanan OCR: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca response OCR: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusBadRequest, http.StatusUnprocessableEntity:
		if message := apiCoIDMessage(respBody); message != "" {
			return nil, fmt.Errorf("%w: %s", ErrOCRNotReadable, message)
		}
		return nil, ErrOCRNotReadable
	case http.StatusUnauthorized:
		return nil, ErrOCRUnauthorized
	case http.StatusPaymentRequired: // 402
		return nil, ErrOCRInsufficientCredits
	case http.StatusOK:
		// handled below
	default:
		return nil, fmt.Errorf("layanan OCR mengembalikan status %d: %s", resp.StatusCode, string(respBody))
	}

	var apiResp apiCoIDResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("gagal parse response OCR: %w", err)
	}
	if !apiResp.IsSuccess {
		return nil, fmt.Errorf("OCR gagal: %s", apiResp.Message)
	}

	result := mapToKTPResult(&apiResp.Data)
	return result, nil
}

func createImageFormFile(writer *multipart.Writer, fieldname, filename string, image []byte) (io.Writer, error) {
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", multipart.FileContentDisposition(fieldname, filename))
	header.Set("Content-Type", imageContentType(image, filename))
	return writer.CreatePart(header)
}

func imageContentType(image []byte, filename string) string {
	if len(image) > 0 {
		contentType := http.DetectContentType(image)
		if contentType == "image/jpeg" || contentType == "image/png" {
			return contentType
		}
	}

	switch strings.ToLower(filepath.Ext(filename)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
}

func apiCoIDMessage(body []byte) string {
	var resp struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return ""
	}
	return strings.TrimSpace(resp.Message)
}

func mapToKTPResult(d *apiCoIDKTPData) *KTPOCRResult {
	birthDate := parseBirthDate(d.TanggalLahir)
	address := buildAddress(d)

	keyFields := []string{d.NIK, d.Nama, d.Alamat, d.TanggalLahir}
	filled := 0
	for _, f := range keyFields {
		if f != "" {
			filled++
		}
	}
	confidence := 0.7 + float64(filled)/float64(len(keyFields))*0.27

	return &KTPOCRResult{
		NIK:        d.NIK,
		FullName:   d.Nama,
		Address:    address,
		Pekerjaan:  d.Pekerjaan,
		BirthDate:  birthDate,
		Confidence: confidence,
		Source:     "API_CO_ID",
	}
}

// parseBirthDate converts "DD-MM-YYYY" to "YYYY-MM-DD".
func parseBirthDate(raw string) string {
	t, err := time.Parse("02-01-2006", raw)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func buildAddress(d *apiCoIDKTPData) string {
	parts := []string{}
	if d.Alamat != "" {
		parts = append(parts, d.Alamat)
	}
	if d.RTRW != "" {
		parts = append(parts, "RT/RW "+d.RTRW)
	}
	if d.Kelurahan != "" {
		parts = append(parts, "Kel. "+d.Kelurahan)
	}
	if d.Kecamatan != "" {
		parts = append(parts, "Kec. "+d.Kecamatan)
	}
	if d.Kota != "" {
		parts = append(parts, d.Kota)
	}
	if d.Provinsi != "" {
		parts = append(parts, d.Provinsi)
	}
	return strings.Join(parts, ", ")
}
