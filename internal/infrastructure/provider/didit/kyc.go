package didit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
)

type DiditSessionReq struct {
	WorkflowID string `json:"workflow_id"`
	VendorData string `json:"vendor_data"`
	Callback   string `json:"callback,omitempty"`
	Metadata   string `json:"metadata,omitempty"`
}

type diditSessionRes struct {
	SessionID string `json:"session_id"`
	URL       string `json:"url"`
}

type kycProviderImpl struct {
	cfg config.KYCConfig
}

func NewKYCProvider(cfg config.KYCConfig) port.KYCProvider {
	return &kycProviderImpl{cfg}
}

func (k *kycProviderImpl) CreateVerificationSession(ctx context.Context, vendorData string, callbackURL string) (*port.VerificationSession, error) {
	if k.cfg.AccessToken == "" {
		return nil, fmt.Errorf("didit access token is required")
	}
	if k.cfg.WorkFlowID == "" {
		return nil, fmt.Errorf("didit workflow id is required")
	}

	metadataBytes, err := json.Marshal(map[string]any{
		"source":    "web-app",
		"timestamp": time.Now().UnixMilli(),
	})
	if err != nil {
		return nil, err
	}

	reqBody := DiditSessionReq{
		WorkflowID: k.cfg.WorkFlowID,
		VendorData: vendorData,
		Callback:   callbackURL,
		Metadata:   string(metadataBytes),
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://verification.didit.me/v3/session/", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", k.cfg.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("didit create verification session failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var sessionRes diditSessionRes
	if err = json.Unmarshal(respBody, &sessionRes); err != nil {
		return nil, err
	}
	if sessionRes.SessionID == "" || sessionRes.URL == "" {
		return nil, fmt.Errorf("didit create verification session invalid response")
	}

	return &port.VerificationSession{
		SessionID:       sessionRes.SessionID,
		VerificationURL: sessionRes.URL,
	}, nil
}
