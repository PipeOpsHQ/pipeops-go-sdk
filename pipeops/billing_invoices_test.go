package pipeops

import (
	"context"
	"net/http"
	"testing"
)

func TestBillingService_ListInvoices_UsesBillingHistoryRoute(t *testing.T) {
	t.Parallel()

	client, err := NewClient("https://api.pipeops.test")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}
	client.SetHTTPClient(&http.Client{
		Transport: billingRoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
			}
			if r.URL.Path != "/billing/history" {
				t.Fatalf("path = %s, want %s", r.URL.Path, "/billing/history")
			}
			return billingJSONResponse(r, http.StatusOK, `{"success":true,"message":"billing history fetched succesfully","data":[{"UID":"inv_1","InvoiceNumber":"PIP-123","Amount":19.99,"Currency":"USD","Status":"paid","Date":"2024-01-02T03:04:05Z","CreatedAt":"2024-01-01T03:04:05Z"}],"meta":{"current_page":1}}`), nil
		}),
	})

	resp, _, err := client.Billing.ListInvoices(context.Background())
	if err != nil {
		t.Fatalf("Billing.ListInvoices error: %v", err)
	}
	if len(resp.Data.Invoices) != 1 {
		t.Fatalf("invoices len = %d, want %d", len(resp.Data.Invoices), 1)
	}
	invoice := resp.Data.Invoices[0]
	if invoice.UUID != "inv_1" {
		t.Fatalf("invoice UUID = %q, want %q", invoice.UUID, "inv_1")
	}
	if invoice.InvoiceNumber != "PIP-123" {
		t.Fatalf("invoice number = %q, want %q", invoice.InvoiceNumber, "PIP-123")
	}
	if invoice.Amount != 19.99 {
		t.Fatalf("amount = %v, want %v", invoice.Amount, 19.99)
	}
	if invoice.Status != "paid" {
		t.Fatalf("status = %q, want %q", invoice.Status, "paid")
	}
}
