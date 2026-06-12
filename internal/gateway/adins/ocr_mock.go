package adins

import (
	"context"
	"strings"
)

// MockKTPOCRGateway returns deterministic mock KTP data.
// Used when OCR_API_KEY is not configured (demo/dev fallback).
// Source field is "MOCK_LITEDMS" to match the framing in WIRADANA.md §10.1.
type MockKTPOCRGateway struct{}

// mockKTPData maps lowercase filename keywords to demo KTP results.
var mockKTPData = []struct {
	keyword string
	result  KTPOCRResult
}{
	{"budi", KTPOCRResult{
		NIK: "3201010101010001", FullName: "Budi Santosa",
		Address: "Jl. Melati No. 10, Bandung", BirthDate: "1990-01-15",
		Confidence: 0.97, Source: "MOCK_LITEDMS",
	}},
	{"siti", KTPOCRResult{
		NIK: "3201010101010002", FullName: "Siti Aminah",
		Address: "Jl. Kenanga No. 21, Bandung", BirthDate: "1992-06-20",
		Confidence: 0.95, Source: "MOCK_LITEDMS",
	}},
	{"asep", KTPOCRResult{
		NIK: "3201010101010004", FullName: "Asep Kurniawan",
		Address: "Jl. Cempaka No. 3, Bandung", BirthDate: "1987-03-22",
		Confidence: 0.96, Source: "MOCK_LITEDMS",
	}},
}

// defaultMockKTP is returned when no filename keyword matches.
var defaultMockKTP = KTPOCRResult{
	NIK: "3201019999999999", FullName: "Anggota Demo",
	Address: "Jl. Demo No. 1, Bandung", BirthDate: "1990-01-01",
	Confidence: 0.92, Source: "MOCK_LITEDMS",
}

func (g *MockKTPOCRGateway) ExtractKTP(_ context.Context, _ []byte, filename string) (*KTPOCRResult, error) {
	lower := strings.ToLower(filename)
	for _, entry := range mockKTPData {
		if strings.Contains(lower, entry.keyword) {
			result := entry.result
			return &result, nil
		}
	}
	result := defaultMockKTP
	return &result, nil
}
