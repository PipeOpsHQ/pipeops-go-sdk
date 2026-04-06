package pipeops

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type billingRoundTripperFunc func(*http.Request) (*http.Response, error)

func (f billingRoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func billingJSONResponse(req *http.Request, statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}
}

func TestBillingServiceControllerResponseCompatibility(t *testing.T) {
	t.Parallel()

	client, err := NewClient("https://api.pipeops.test")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}
	client.SetHTTPClient(&http.Client{
		Transport: billingRoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			switch r.URL.Path {
			case "/workspace":
				return billingJSONResponse(r, http.StatusOK, `{"data":[{"UUID":"w1"}],"message":"ok","success":true}`), nil
			case "/billing/balance":
				return billingJSONResponse(r, http.StatusOK, `{"data":{"Balance":"0.01","Currency":"USD"},"message":"ok","success":true}`), nil
			case "/billing/subscriptions/current":
				return billingJSONResponse(r, http.StatusOK, `{"data":{"UID":"sub_123","PlanTier":"startup","PlanName":"Start-up","Amount":"34.99","BillingType":"trial","Status":"active"},"message":"ok","success":true}`), nil
			case "/billing/portal":
				return billingJSONResponse(r, http.StatusOK, `{"data":"https://billing.example/session","message":"portal created successfully","success":true}`), nil
			case "/billing/plans":
				return billingJSONResponse(r, http.StatusOK, `{"data":[{"Name":"Custom","Price":550000,"Period":"monthly","Currency":"NGN","ConcurrentBuild":1,"FreeTrialDays":0}],"message":"ok","success":true}`), nil
			case "/billing/workspace/cards":
				if got := r.URL.Query().Get("workspace_uuid"); got != "w1" {
					t.Fatalf("workspace_uuid = %q, want %q", got, "w1")
				}
				return billingJSONResponse(r, http.StatusOK, `{"data":[{"UID":"card_123","Provider":"stripe","Last4Digit":"4242","ExpiryMonth":"12","ExpiryYear":"2034","CardType":"visa","Country":"US","IsActiveBillingCard":false}],"message":"provider profiles fetched","success":true}`), nil
			case "/billing/workspace/cards/active":
				if got := r.URL.Query().Get("workspace_uuid"); got != "w1" {
					t.Fatalf("workspace_uuid = %q, want %q", got, "w1")
				}
				return billingJSONResponse(r, http.StatusOK, `{"data":{"UID":"card_active","Provider":"stripe","Last4Digit":"4242","ExpiryMonth":"12","ExpiryYear":"2034","CardType":"visa","Country":"US","IsActiveBillingCard":true},"message":"billing card fetched succesfully.","success":true}`), nil
			case "/billing/cards":
				if r.Method != http.MethodPost {
					t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
				}
				return billingJSONResponse(r, http.StatusOK, `{"data":{"checkout_url":"https://checkout.example/card"},"message":"checkout url generated successfully","success":true}`), nil
			default:
				t.Fatalf("unexpected request: %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
				return nil, nil
			}
		}),
	})

	balance, _, err := client.Billing.GetBalance(context.Background())
	if err != nil {
		t.Fatalf("GetBalance error: %v", err)
	}
	if balance.Data.Balance != 0.01 {
		t.Fatalf("Balance = %v, want %v", balance.Data.Balance, 0.01)
	}
	if balance.Data.Currency != "USD" {
		t.Fatalf("Currency = %q, want %q", balance.Data.Currency, "USD")
	}

	subscription, _, err := client.Billing.GetCurrentSubscription(context.Background())
	if err != nil {
		t.Fatalf("GetCurrentSubscription error: %v", err)
	}
	if subscription.Data.Subscription.PlanTier != "startup" {
		t.Fatalf("PlanTier = %q, want %q", subscription.Data.Subscription.PlanTier, "startup")
	}
	if subscription.Data.Subscription.Amount != 34.99 {
		t.Fatalf("Amount = %v, want %v", subscription.Data.Subscription.Amount, 34.99)
	}
	if subscription.Data.Subscription.BillingType != "trial" {
		t.Fatalf("BillingType = %q, want %q", subscription.Data.Subscription.BillingType, "trial")
	}

	portal, _, err := client.Billing.GetPortalURL(context.Background())
	if err != nil {
		t.Fatalf("GetPortalURL error: %v", err)
	}
	if portal.Data.PortalURL != "https://billing.example/session" {
		t.Fatalf("PortalURL = %q, want %q", portal.Data.PortalURL, "https://billing.example/session")
	}

	plans, _, err := client.Billing.GetPlans(context.Background())
	if err != nil {
		t.Fatalf("GetPlans error: %v", err)
	}
	if len(plans.Data.Plans) != 1 {
		t.Fatalf("plans len = %d, want 1", len(plans.Data.Plans))
	}
	if plans.Data.Plans[0].Name != "Custom" {
		t.Fatalf("plan name = %q, want %q", plans.Data.Plans[0].Name, "Custom")
	}
	if plans.Data.Plans[0].Price != 550000 {
		t.Fatalf("plan price = %v, want %v", plans.Data.Plans[0].Price, 550000.0)
	}
	if plans.Data.Plans[0].Period != "monthly" {
		t.Fatalf("plan period = %q, want %q", plans.Data.Plans[0].Period, "monthly")
	}

	cards, _, err := client.Billing.ListWorkspaceCards(context.Background())
	if err != nil {
		t.Fatalf("ListWorkspaceCards error: %v", err)
	}
	if len(cards.Data.Cards) != 1 {
		t.Fatalf("cards len = %d, want 1", len(cards.Data.Cards))
	}
	if cards.Data.Cards[0].UUID != "card_123" {
		t.Fatalf("card UUID = %q, want %q", cards.Data.Cards[0].UUID, "card_123")
	}
	if cards.Data.Cards[0].Last4 != "4242" {
		t.Fatalf("card Last4 = %q, want %q", cards.Data.Cards[0].Last4, "4242")
	}

	activeCard, _, err := client.Billing.GetActiveCard(context.Background())
	if err != nil {
		t.Fatalf("GetActiveCard error: %v", err)
	}
	if activeCard.Data.Card.UUID != "card_active" {
		t.Fatalf("active card UUID = %q, want %q", activeCard.Data.Card.UUID, "card_active")
	}
	if !activeCard.Data.Card.IsActive {
		t.Fatal("expected active card to be marked active")
	}

	addCard, _, err := client.Billing.AddCard(context.Background(), &AddCardRequest{Token: "tok"})
	if err != nil {
		t.Fatalf("AddCard error: %v", err)
	}
	if addCard.Data.CheckoutURL != "https://checkout.example/card" {
		t.Fatalf("checkout URL = %q, want %q", addCard.Data.CheckoutURL, "https://checkout.example/card")
	}
}

func TestBillingServiceControllerRouteCompatibility(t *testing.T) {
	t.Parallel()

	client, err := NewClient("https://api.pipeops.test")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}
	client.SetHTTPClient(&http.Client{
		Transport: billingRoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			switch r.URL.Path {
			case "/workspace":
				return billingJSONResponse(r, http.StatusOK, `{"data":[{"UUID":"w1"}],"message":"ok","success":true}`), nil
			case "/billing/subscribe":
				if r.Method != http.MethodPost {
					t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
				}
				if got := r.URL.Query().Get("plan"); got != "startup" {
					t.Fatalf("plan = %q, want %q", got, "startup")
				}
				body := []byte{}
				if r.Body != nil {
					var err error
					body, err = io.ReadAll(r.Body)
					if err != nil {
						t.Fatalf("read body error: %v", err)
					}
				}
				if strings.TrimSpace(string(body)) != "" {
					t.Fatalf("expected empty body, got %q", string(body))
				}
				return billingJSONResponse(r, http.StatusOK, `{"data":{"checkout_url":"https://checkout.example/subscribe","message":"subscription processing"},"message":"ok","success":true}`), nil
			case "/billing/start-trial":
				if r.Method != http.MethodPost {
					t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
				}
				if got := r.URL.Query().Get("plan"); got != "startup_v1" {
					t.Fatalf("plan = %q, want %q", got, "startup_v1")
				}
				body := []byte{}
				if r.Body != nil {
					var err error
					body, err = io.ReadAll(r.Body)
					if err != nil {
						t.Fatalf("read body error: %v", err)
					}
				}
				if strings.TrimSpace(string(body)) != "" {
					t.Fatalf("expected empty body, got %q", string(body))
				}
				return billingJSONResponse(r, http.StatusOK, `{"data":{},"message":"trial is being processed","success":true}`), nil
			case "/billing/subscriptions/workspace/team-seat":
				if got := r.URL.Query().Get("workspace"); got != "w1" {
					t.Fatalf("workspace = %q, want %q", got, "w1")
				}
				return billingJSONResponse(r, http.StatusOK, `{"data":[{"UID":"sub_team","PlanName":"Team Seat Topup","Quantity":10,"Status":"active"}],"message":"ok","success":true}`), nil
			default:
				t.Fatalf("unexpected request: %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
				return nil, nil
			}
		}),
	})

	subscribeResp, _, err := client.Billing.Subscribe(context.Background(), &SubscribeRequest{PlanID: "startup"})
	if err != nil {
		t.Fatalf("Subscribe error: %v", err)
	}
	if subscribeResp.Data.CheckoutURL != "https://checkout.example/subscribe" {
		t.Fatalf("checkout URL = %q, want %q", subscribeResp.Data.CheckoutURL, "https://checkout.example/subscribe")
	}

	if _, err := client.Billing.StartTrial(context.Background(), &StartTrialRequest{PlanID: "startup_v1"}); err != nil {
		t.Fatalf("StartTrial error: %v", err)
	}

	teamSeatResp, _, err := client.Billing.GetTeamSeatSubscription(context.Background())
	if err != nil {
		t.Fatalf("GetTeamSeatSubscription error: %v", err)
	}
	if len(teamSeatResp.Data.Subscriptions) != 1 {
		t.Fatalf("team seat subscriptions len = %d, want 1", len(teamSeatResp.Data.Subscriptions))
	}
	if teamSeatResp.Data.Subscriptions[0].PlanName != "Team Seat Topup" {
		t.Fatalf("plan name = %q, want %q", teamSeatResp.Data.Subscriptions[0].PlanName, "Team Seat Topup")
	}
}

func TestBillingServiceGetBillingInfo(t *testing.T) {
	t.Parallel()

	t.Run("with current subscription", func(t *testing.T) {
		client, err := NewClient("https://api.pipeops.test")
		if err != nil {
			t.Fatalf("NewClient error: %v", err)
		}
		client.SetHTTPClient(&http.Client{
			Transport: billingRoundTripperFunc(func(r *http.Request) (*http.Response, error) {
				switch r.URL.Path {
				case "/billing/balance":
					return billingJSONResponse(r, http.StatusOK, `{"data":{"Balance":"10.00","Currency":"USD"},"message":"ok","success":true}`), nil
				case "/billing/subscriptions/current":
					return billingJSONResponse(r, http.StatusOK, `{"data":{"UID":"sub_1","PlanTier":"startup","Status":"active"},"message":"ok","success":true}`), nil
				default:
					t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
					return nil, nil
				}
			}),
		})

		billingInfo, _, err := client.Billing.GetBillingInfo(context.Background())
		if err != nil {
			t.Fatalf("GetBillingInfo error: %v", err)
		}
		if billingInfo.Data.Balance.Balance != 10 {
			t.Fatalf("balance = %v, want %v", billingInfo.Data.Balance.Balance, 10.0)
		}
		if billingInfo.Data.CurrentSubscription == nil {
			t.Fatal("expected current subscription")
		}
		if billingInfo.Data.CurrentSubscription.PlanTier != "startup" {
			t.Fatalf("plan tier = %q, want %q", billingInfo.Data.CurrentSubscription.PlanTier, "startup")
		}
	})

	t.Run("without current subscription", func(t *testing.T) {
		client, err := NewClient("https://api.pipeops.test")
		if err != nil {
			t.Fatalf("NewClient error: %v", err)
		}
		client.SetHTTPClient(&http.Client{
			Transport: billingRoundTripperFunc(func(r *http.Request) (*http.Response, error) {
				switch r.URL.Path {
				case "/billing/balance":
					return billingJSONResponse(r, http.StatusOK, `{"data":{"Balance":"3.50","Currency":"USD"},"message":"ok","success":true}`), nil
				case "/billing/subscriptions/current":
					return billingJSONResponse(r, http.StatusNotFound, `{"message":"not found"}`), nil
				default:
					t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
					return nil, nil
				}
			}),
		})

		billingInfo, _, err := client.Billing.GetBillingInfo(context.Background())
		if err != nil {
			t.Fatalf("GetBillingInfo error: %v", err)
		}
		if billingInfo.Data.CurrentSubscription != nil {
			t.Fatalf("expected nil current subscription, got %#v", billingInfo.Data.CurrentSubscription)
		}
	})
}

func TestBillingServiceGetUsageController404(t *testing.T) {
	t.Parallel()

	client, err := NewClient("https://api.pipeops.test")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}
	client.SetHTTPClient(&http.Client{
		Transport: billingRoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path != "/billing/usage" {
				t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
			}
			return billingJSONResponse(r, http.StatusNotFound, `{"message":"not found"}`), nil
		}),
	})

	_, _, err = client.Billing.GetUsage(context.Background())
	if err == nil {
		t.Fatal("expected GetUsage error")
	}
	if !strings.Contains(err.Error(), "GetBillingInfo") {
		t.Fatalf("expected actionable error, got %v", err)
	}
}
