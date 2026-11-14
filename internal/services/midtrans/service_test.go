package midtrans

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSnapTransaction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}

		username, password, ok := r.BasicAuth()
		if !ok || username != "server-key" || password != "" {
			t.Errorf("basic auth = %q/%q/%v, want server-key with empty password", username, password, ok)
		}

		var body struct {
			TransactionDetails struct {
				OrderID     string  `json:"order_id"`
				GrossAmount float64 `json:"gross_amount"`
			} `json:"transaction_details"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if body.TransactionDetails.OrderID != "BOOK-1" {
			t.Errorf("order_id = %q, want BOOK-1", body.TransactionDetails.OrderID)
		}
		if body.TransactionDetails.GrossAmount != 150000 {
			t.Errorf("gross_amount = %v, want 150000", body.TransactionDetails.GrossAmount)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"token":"snap-token","redirect_url":"https://app.sandbox.midtrans.com/snap/v4/redirection/snap-token"}`))
	}))
	defer server.Close()

	service := NewService(Config{ServerKey: "server-key", SnapURL: server.URL})
	result, err := service.CreateSnapTransaction(context.Background(), SnapRequest{OrderID: "BOOK-1", GrossAmount: 150000})
	if err != nil {
		t.Fatalf("CreateSnapTransaction() error = %v", err)
	}
	if result.Token != "snap-token" {
		t.Fatalf("token = %q, want snap-token", result.Token)
	}
}

func TestCreateSnapTransactionRequiresConfiguration(t *testing.T) {
	_, err := NewService(Config{}).CreateSnapTransaction(context.Background(), SnapRequest{OrderID: "BOOK-1", GrossAmount: 1000})
	if !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("error = %v, want ErrNotConfigured", err)
	}
}

func TestCreateSnapTransactionGatewayError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"error_messages":["transaction_details.order_id is required"]}`))
	}))
	defer server.Close()

	service := NewService(Config{ServerKey: "server-key", SnapURL: server.URL})
	_, err := service.CreateSnapTransaction(context.Background(), SnapRequest{OrderID: "BOOK-1", GrossAmount: 1000})
	if !errors.Is(err, ErrGateway) {
		t.Fatalf("error = %v, want ErrGateway", err)
	}
}

func TestRefund(t *testing.T) {
	tests := []struct {
		name             string
		transactionState string
		wantApproved     bool
		wantPending      bool
	}{
		{name: "refunded", transactionState: "refunded", wantApproved: true},
		{name: "partial refunded", transactionState: "partial_refunded", wantApproved: true},
		{name: "refund pending", transactionState: "refund_pending", wantPending: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/v2/BOOK-1/refund" {
					t.Errorf("path = %q, want /v2/BOOK-1/refund", r.URL.Path)
				}

				username, _, ok := r.BasicAuth()
				if !ok || username != "server-key" {
					t.Errorf("basic auth username = %q/%v, want server-key", username, ok)
				}

				var body struct {
					RefundKey string  `json:"refund_key"`
					Amount    float64 `json:"amount"`
					Reason    string  `json:"reason"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("decode body: %v", err)
				}
				if body.RefundKey != "refund-1" || body.Amount != 50000 || body.Reason != "event cancelled" {
					t.Errorf("refund body = %+v, want refund-1/50000/event cancelled", body)
				}

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"status_code":"200","status_message":"Success","transaction_status":"` + tt.transactionState + `"}`))
			}))
			defer server.Close()

			service := NewService(Config{ServerKey: "server-key", CoreURL: server.URL})
			result, err := service.Refund(context.Background(), RefundRequest{
				OrderID:   "BOOK-1",
				RefundKey: "refund-1",
				Amount:    50000,
				Reason:    "event cancelled",
			})
			if err != nil {
				t.Fatalf("Refund() error = %v", err)
			}
			if result.Approved != tt.wantApproved {
				t.Fatalf("Approved = %v, want %v", result.Approved, tt.wantApproved)
			}
			if result.Pending != tt.wantPending {
				t.Fatalf("Pending = %v, want %v", result.Pending, tt.wantPending)
			}
		})
	}
}

func TestRefundGatewayError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"status_code":"400","status_message":"Refund amount exceeds remaining balance"}`))
	}))
	defer server.Close()

	service := NewService(Config{ServerKey: "server-key", CoreURL: server.URL})
	_, err := service.Refund(context.Background(), RefundRequest{OrderID: "BOOK-1", RefundKey: "refund-1", Amount: 50000})
	if !errors.Is(err, ErrGateway) {
		t.Fatalf("error = %v, want ErrGateway", err)
	}
}

func TestVerifySignature(t *testing.T) {
	service := NewService(Config{ServerKey: "server-key"})

	sum := sha512.Sum512([]byte("BOOK-1200" + "150000.00" + "server-key"))
	valid := hex.EncodeToString(sum[:])

	if !service.VerifySignature("BOOK-1", "200", "150000.00", valid) {
		t.Fatal("VerifySignature() = false, want true")
	}
	if service.VerifySignature("BOOK-1", "200", "150000.00", "invalid") {
		t.Fatal("VerifySignature() = true for invalid signature")
	}
	if NewService(Config{}).VerifySignature("BOOK-1", "200", "150000.00", valid) {
		t.Fatal("VerifySignature() = true without a server key")
	}
}

func TestMapTransactionStatus(t *testing.T) {
	tests := []struct {
		transaction, fraud, want string
	}{
		{transaction: "capture", fraud: "accept", want: "success"},
		{transaction: "capture", fraud: "challenge", want: "pending"},
		{transaction: "capture", fraud: "deny", want: "failed"},
		{transaction: "settlement", want: "success"},
		{transaction: "pending", want: "pending"},
		{transaction: "deny", want: "failed"},
		{transaction: "cancel", want: "failed"},
		{transaction: "expire", want: "expired"},
		{transaction: "refund", want: "refunded"},
		{transaction: "partial_refund", want: "partially_refunded"},
		{transaction: "unknown", want: "pending"},
	}

	for _, tt := range tests {
		if got := MapTransactionStatus(tt.transaction, tt.fraud); got != tt.want {
			t.Errorf("MapTransactionStatus(%q, %q) = %q, want %q", tt.transaction, tt.fraud, got, tt.want)
		}
	}
}
