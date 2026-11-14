package midtrans

import (
	"bytes"
	"context"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	sandboxSnapURL    = "https://app.sandbox.midtrans.com/snap/v1/transactions"
	productionSnapURL = "https://app.midtrans.com/snap/v1/transactions"
	sandboxCoreURL    = "https://api.sandbox.midtrans.com"
	productionCoreURL = "https://api.midtrans.com"
)

var (
	ErrNotConfigured    = errors.New("midtrans is not configured")
	ErrInvalidSignature = errors.New("invalid midtrans signature")
	ErrGateway          = errors.New("midtrans request failed")
)

type Config struct {
	KeyName      string
	ServerKey    string
	ClientKey    string
	IsProduction bool
	SnapURL      string
	CoreURL      string
}

type Service struct {
	config  Config
	client  *http.Client
	snapURL string
	coreURL string
}

func NewService(config Config) *Service {
	snapURL := strings.TrimSpace(config.SnapURL)
	if snapURL == "" {
		if config.IsProduction {
			snapURL = productionSnapURL
		} else {
			snapURL = sandboxSnapURL
		}
	}

	coreURL := strings.TrimRight(strings.TrimSpace(config.CoreURL), "/")
	if coreURL == "" {
		if config.IsProduction {
			coreURL = productionCoreURL
		} else {
			coreURL = sandboxCoreURL
		}
	}

	return &Service{
		config:  config,
		client:  &http.Client{Timeout: 30 * time.Second},
		snapURL: snapURL,
		coreURL: coreURL,
	}
}

func (s *Service) IsConfigured() bool {
	return s != nil && strings.TrimSpace(s.config.ServerKey) != ""
}

func (s *Service) ClientKey() string {
	if s == nil {
		return ""
	}
	return s.config.ClientKey
}

type SnapRequest struct {
	OrderID     string
	GrossAmount float64
}

type SnapResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

func (s *Service) CreateSnapTransaction(ctx context.Context, request SnapRequest) (*SnapResponse, error) {
	if !s.IsConfigured() {
		return nil, ErrNotConfigured
	}
	if strings.TrimSpace(request.OrderID) == "" {
		return nil, fmt.Errorf("midtrans order_id is required")
	}
	if request.GrossAmount <= 0 {
		return nil, fmt.Errorf("midtrans gross_amount must be greater than zero")
	}

	body := map[string]any{
		"transaction_details": map[string]any{
			"order_id":     request.OrderID,
			"gross_amount": math.Round(request.GrossAmount*100) / 100,
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("encode midtrans request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, s.snapURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build midtrans request: %w", err)
	}
	httpRequest.SetBasicAuth(s.config.ServerKey, "")
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := s.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("send midtrans request: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read midtrans response: %w", err)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("%w: status %d: %s", ErrGateway, response.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	result := &SnapResponse{}
	if err := json.Unmarshal(responseBody, result); err != nil {
		return nil, fmt.Errorf("decode midtrans response: %w", err)
	}
	if result.Token == "" || result.RedirectURL == "" {
		return nil, fmt.Errorf("midtrans response did not contain a token and redirect_url")
	}

	return result, nil
}

type RefundRequest struct {
	OrderID   string
	RefundKey string
	Amount    float64
	Reason    string
}

type RefundResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	RefundKey         string `json:"refund_key"`
	RefundAmount      string `json:"refund_amount"`
}

type RefundResult struct {
	Approved bool
	Pending  bool
	Status   string
	Message  string
}

func (s *Service) Refund(ctx context.Context, request RefundRequest) (*RefundResult, error) {
	if !s.IsConfigured() {
		return nil, ErrNotConfigured
	}
	if strings.TrimSpace(request.OrderID) == "" {
		return nil, fmt.Errorf("midtrans order_id is required")
	}
	if request.Amount <= 0 {
		return nil, fmt.Errorf("midtrans refund amount must be greater than zero")
	}

	body := map[string]any{
		"refund_key": request.RefundKey,
		"amount":     math.Round(request.Amount*100) / 100,
	}
	if strings.TrimSpace(request.Reason) != "" {
		body["reason"] = request.Reason
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("encode midtrans refund request: %w", err)
	}

	endpoint := s.coreURL + "/v2/" + url.PathEscape(request.OrderID) + "/refund"
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build midtrans refund request: %w", err)
	}
	httpRequest.SetBasicAuth(s.config.ServerKey, "")
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := s.client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("send midtrans refund request: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read midtrans refund response: %w", err)
	}

	result := &RefundResponse{}
	if err := json.Unmarshal(responseBody, result); err != nil {
		return nil, fmt.Errorf("decode midtrans refund response: %w", err)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		message := strings.TrimSpace(result.StatusMessage)
		if message == "" {
			message = strings.TrimSpace(string(responseBody))
		}
		return nil, fmt.Errorf("%w: status %d: %s", ErrGateway, response.StatusCode, message)
	}

	status := strings.ToLower(strings.TrimSpace(result.TransactionStatus))
	refund := &RefundResult{Status: status, Message: result.StatusMessage}
	switch status {
	case "refunded", "partial_refunded":
		refund.Approved = true
	case "refund_pending", "pending":
		refund.Pending = true
	}

	return refund, nil
}

func (s *Service) VerifySignature(orderID, statusCode, grossAmount, signature string) bool {
	if !s.IsConfigured() || strings.TrimSpace(signature) == "" {
		return false
	}

	sum := sha512.Sum512([]byte(orderID + statusCode + grossAmount + s.config.ServerKey))
	expected := hex.EncodeToString(sum[:])
	return subtle.ConstantTimeCompare([]byte(expected), []byte(strings.ToLower(signature))) == 1
}

func MapTransactionStatus(transactionStatus, fraudStatus string) string {
	switch strings.ToLower(strings.TrimSpace(transactionStatus)) {
	case "capture":
		switch strings.ToLower(strings.TrimSpace(fraudStatus)) {
		case "challenge":
			return "pending"
		case "deny":
			return "failed"
		default:
			return "success"
		}
	case "settlement":
		return "success"
	case "pending":
		return "pending"
	case "deny", "cancel", "failure":
		return "failed"
	case "expire":
		return "expired"
	case "refund":
		return "refunded"
	case "partial_refund":
		return "partially_refunded"
	case "chargeback", "partial_chargeback":
		return "refunded"
	default:
		return "pending"
	}
}

func ParseTransactionTime(value string) *time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	location := time.FixedZone("WIB", 7*60*60)
	for _, layout := range []string{"2006-01-02 15:04:05", time.RFC3339} {
		parsed, err := time.ParseInLocation(layout, value, location)
		if err == nil {
			return &parsed
		}
	}
	return nil
}
