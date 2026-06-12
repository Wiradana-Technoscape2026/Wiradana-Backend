package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WhatsAppGateway struct {
	token   string
	phoneID string
	baseURL string
}

func NewWhatsAppGateway(token, phoneID string) Gateway {
	return &WhatsAppGateway{
		token:   token,
		phoneID: phoneID,
		baseURL: "https://graph.facebook.com/v19.0",
	}
}

func (g *WhatsAppGateway) Source() string {
	return "WA_CLOUD_SANDBOX"
}

func (g *WhatsAppGateway) Send(ctx context.Context, input Input) (Result, error) {
	payload := map[string]any{
		"messaging_product": "whatsapp",
		"to":                input.ToPhone,
		"type":              "text",
		"text":              map[string]string{"body": input.Message},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return Result{}, err
	}

	url := fmt.Sprintf("%s/%s/messages", g.baseURL, g.phoneID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return Result{}, fmt.Errorf("whatsapp api error %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Messages []struct {
			ID string `json:"id"`
		} `json:"messages"`
	}
	_ = json.Unmarshal(respBody, &result)

	msgID := ""
	if len(result.Messages) > 0 {
		msgID = result.Messages[0].ID
	}

	return Result{MessageID: msgID, Source: "WA_CLOUD_SANDBOX"}, nil
}
