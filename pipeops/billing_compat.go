package pipeops

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type billingEnvelope struct {
	Status  string          `json:"status,omitempty"`
	Success *bool           `json:"success,omitempty"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type billingPlanOptions struct {
	Plan string `url:"plan"`
}

type billingWorkspaceOption struct {
	Workspace string `url:"workspace"`
}

// BillingInfoResponse represents controller-backed billing information.
type BillingInfoResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Balance             BillingBalance `json:"balance"`
		CurrentSubscription *Subscription  `json:"current_subscription,omitempty"`
	} `json:"data"`
}

// BillingBalance represents the current wallet balance snapshot.
type BillingBalance struct {
	Balance  float64 `json:"balance,omitempty"`
	Currency string  `json:"currency,omitempty"`
}

// GetBillingInfo retrieves controller-backed billing balance and current subscription information.
func (s *BillingService) GetBillingInfo(ctx context.Context) (*BillingInfoResponse, *http.Response, error) {
	balanceResp, resp, err := s.GetBalance(ctx)
	if err != nil {
		return nil, resp, err
	}

	info := &BillingInfoResponse{
		Status:  coalesceNonEmpty(balanceResp.Status, "success"),
		Message: "Billing information retrieved successfully",
	}
	info.Data.Balance.Balance = balanceResp.Data.Balance
	info.Data.Balance.Currency = balanceResp.Data.Currency

	subscriptionResp, subscriptionHTTPResp, err := s.GetCurrentSubscription(ctx)
	if err != nil {
		if isNotFound(err) {
			return info, resp, nil
		}
		return nil, subscriptionHTTPResp, err
	}

	subscription := subscriptionResp.Data.Subscription
	if !isEmptySubscription(subscription) {
		info.Data.CurrentSubscription = &subscription
	}
	if subscriptionHTTPResp != nil {
		resp = subscriptionHTTPResp
	}

	return info, resp, nil
}

func (c *Card) UnmarshalJSON(data []byte) error {
	type cardWire struct {
		ID           jsonID          `json:"id,omitempty"`
		IDAlt        jsonID          `json:"ID,omitempty"`
		UUID         string          `json:"uuid,omitempty"`
		UUIDAlt      string          `json:"UID,omitempty"`
		Provider     string          `json:"provider,omitempty"`
		ProviderAlt  string          `json:"Provider,omitempty"`
		Last4        string          `json:"last4,omitempty"`
		Last4Alt     string          `json:"Last4Digit,omitempty"`
		Brand        string          `json:"brand,omitempty"`
		BrandAlt     string          `json:"Brand,omitempty"`
		CardType     string          `json:"card_type,omitempty"`
		CardTypeAlt  string          `json:"CardType,omitempty"`
		ExpMonth     json.RawMessage `json:"exp_month,omitempty"`
		ExpMonthAlt  json.RawMessage `json:"ExpiryMonth,omitempty"`
		ExpYear      json.RawMessage `json:"exp_year,omitempty"`
		ExpYearAlt   json.RawMessage `json:"ExpiryYear,omitempty"`
		IsDefault    *bool           `json:"is_default,omitempty"`
		IsDefaultAlt *bool           `json:"IsDefault,omitempty"`
		IsActive     *bool           `json:"is_active,omitempty"`
		IsActiveAlt  *bool           `json:"IsActiveBillingCard,omitempty"`
		Channel      string          `json:"channel,omitempty"`
		ChannelAlt   string          `json:"Channel,omitempty"`
		Bank         string          `json:"bank,omitempty"`
		BankAlt      string          `json:"Bank,omitempty"`
		Country      string          `json:"country,omitempty"`
		CountryAlt   string          `json:"Country,omitempty"`
		CreatedAt    *Timestamp      `json:"created_at,omitempty"`
		CreatedAtAlt *Timestamp      `json:"CreatedAt,omitempty"`
	}

	var wire cardWire
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}

	c.ID = coalesceNonEmpty(wire.ID.String(), wire.IDAlt.String())
	c.UUID = coalesceNonEmpty(wire.UUID, wire.UUIDAlt)
	c.Provider = coalesceNonEmpty(wire.Provider, wire.ProviderAlt)
	c.Last4 = coalesceNonEmpty(wire.Last4, wire.Last4Alt)
	c.Brand = coalesceNonEmpty(wire.Brand, wire.BrandAlt, wire.CardType, wire.CardTypeAlt)
	c.CardType = coalesceNonEmpty(wire.CardType, wire.CardTypeAlt, c.Brand)
	if value, err := decodeFlexibleInt(firstNonEmptyRaw(wire.ExpMonth, wire.ExpMonthAlt)); err == nil {
		c.ExpMonth = value
	}
	if value, err := decodeFlexibleInt(firstNonEmptyRaw(wire.ExpYear, wire.ExpYearAlt)); err == nil {
		c.ExpYear = value
	}
	c.IsDefault = firstBool(wire.IsDefault, wire.IsDefaultAlt, wire.IsActive, wire.IsActiveAlt)
	c.IsActive = firstBool(wire.IsActive, wire.IsActiveAlt, wire.IsDefault, wire.IsDefaultAlt)
	c.Channel = coalesceNonEmpty(wire.Channel, wire.ChannelAlt)
	c.Bank = coalesceNonEmpty(wire.Bank, wire.BankAlt)
	c.Country = coalesceNonEmpty(wire.Country, wire.CountryAlt)
	c.CreatedAt = firstTimestamp(wire.CreatedAt, wire.CreatedAtAlt)

	return nil
}

func (s *Subscription) UnmarshalJSON(data []byte) error {
	type subscriptionWire struct {
		ID               jsonID          `json:"id,omitempty"`
		IDAlt            jsonID          `json:"ID,omitempty"`
		UUID             string          `json:"uuid,omitempty"`
		UUIDAlt          string          `json:"UID,omitempty"`
		PlanID           jsonID          `json:"plan_id,omitempty"`
		PlanIDAlt        jsonID          `json:"PlanID,omitempty"`
		PlanName         string          `json:"plan_name,omitempty"`
		PlanNameAlt      string          `json:"PlanName,omitempty"`
		PlanTier         string          `json:"plan_tier,omitempty"`
		PlanTierAlt      string          `json:"PlanTier,omitempty"`
		Status           string          `json:"status,omitempty"`
		StatusAlt        string          `json:"Status,omitempty"`
		BillingType      string          `json:"billing_type,omitempty"`
		BillingTypeAlt   string          `json:"BillingType,omitempty"`
		BillingStatus    string          `json:"billing_status,omitempty"`
		BillingStatusAlt string          `json:"BillingStatus,omitempty"`
		PaymentMethod    string          `json:"payment_method,omitempty"`
		PaymentMethodAlt string          `json:"PaymentMethod,omitempty"`
		PlanPeriod       string          `json:"plan_period,omitempty"`
		PlanPeriodAlt    string          `json:"PlanPeriod,omitempty"`
		Provider         string          `json:"provider,omitempty"`
		ProviderAlt      string          `json:"Provider,omitempty"`
		Description      string          `json:"description,omitempty"`
		DescriptionAlt   string          `json:"Description,omitempty"`
		Quantity         json.RawMessage `json:"quantity,omitempty"`
		QuantityAlt      json.RawMessage `json:"Quantity,omitempty"`
		StartDate        *Timestamp      `json:"start_date,omitempty"`
		StartDateAlt     *Timestamp      `json:"SubStartDate,omitempty"`
		EndDate          *Timestamp      `json:"end_date,omitempty"`
		EndDateAlt       *Timestamp      `json:"SubEndDate,omitempty"`
		Date             *Timestamp      `json:"date,omitempty"`
		DateAlt          *Timestamp      `json:"Date,omitempty"`
		Amount           json.RawMessage `json:"amount,omitempty"`
		AmountAlt        json.RawMessage `json:"Amount,omitempty"`
		Currency         string          `json:"currency,omitempty"`
		CurrencyAlt      string          `json:"Currency,omitempty"`
		CreatedAt        *Timestamp      `json:"created_at,omitempty"`
		CreatedAtAlt     *Timestamp      `json:"CreatedAt,omitempty"`
		UpdatedAt        *Timestamp      `json:"updated_at,omitempty"`
		UpdatedAtAlt     *Timestamp      `json:"UpdatedAt,omitempty"`
	}

	var wire subscriptionWire
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}

	s.ID = coalesceNonEmpty(wire.ID.String(), wire.IDAlt.String())
	s.UUID = coalesceNonEmpty(wire.UUID, wire.UUIDAlt)
	s.PlanID = coalesceNonEmpty(wire.PlanID.String(), wire.PlanIDAlt.String())
	s.PlanName = coalesceNonEmpty(wire.PlanName, wire.PlanNameAlt)
	s.PlanTier = coalesceNonEmpty(wire.PlanTier, wire.PlanTierAlt)
	s.Status = coalesceNonEmpty(wire.Status, wire.StatusAlt, wire.BillingStatus, wire.BillingStatusAlt)
	s.BillingType = coalesceNonEmpty(wire.BillingType, wire.BillingTypeAlt)
	s.BillingStatus = coalesceNonEmpty(wire.BillingStatus, wire.BillingStatusAlt)
	s.PaymentMethod = coalesceNonEmpty(wire.PaymentMethod, wire.PaymentMethodAlt)
	s.PlanPeriod = coalesceNonEmpty(wire.PlanPeriod, wire.PlanPeriodAlt)
	s.Provider = coalesceNonEmpty(wire.Provider, wire.ProviderAlt)
	s.Description = coalesceNonEmpty(wire.Description, wire.DescriptionAlt)
	if value, err := decodeFlexibleInt(firstNonEmptyRaw(wire.Quantity, wire.QuantityAlt)); err == nil {
		s.Quantity = value
	}
	s.StartDate = firstTimestamp(wire.StartDate, wire.StartDateAlt)
	s.EndDate = firstTimestamp(wire.EndDate, wire.EndDateAlt)
	s.Date = firstTimestamp(wire.Date, wire.DateAlt)
	if value, err := decodeFlexibleFloat(firstNonEmptyRaw(wire.Amount, wire.AmountAlt)); err == nil {
		s.Amount = value
	}
	s.Currency = coalesceNonEmpty(wire.Currency, wire.CurrencyAlt)
	s.CreatedAt = firstTimestamp(wire.CreatedAt, wire.CreatedAtAlt)
	s.UpdatedAt = firstTimestamp(wire.UpdatedAt, wire.UpdatedAtAlt)

	return nil
}

func (p *Plan) UnmarshalJSON(data []byte) error {
	type planWire struct {
		ID                 jsonID          `json:"id,omitempty"`
		IDAlt              jsonID          `json:"ID,omitempty"`
		UUID               string          `json:"uuid,omitempty"`
		UUIDAlt            string          `json:"UID,omitempty"`
		Name               string          `json:"name,omitempty"`
		NameAlt            string          `json:"Name,omitempty"`
		Description        string          `json:"description,omitempty"`
		DescriptionAlt     string          `json:"Description,omitempty"`
		Price              json.RawMessage `json:"price,omitempty"`
		PriceAlt           json.RawMessage `json:"Price,omitempty"`
		Currency           string          `json:"currency,omitempty"`
		CurrencyAlt        string          `json:"Currency,omitempty"`
		Interval           string          `json:"interval,omitempty"`
		Period             string          `json:"period,omitempty"`
		PeriodAlt          string          `json:"Period,omitempty"`
		Active             *bool           `json:"active,omitempty"`
		ActiveAlt          *bool           `json:"Active,omitempty"`
		FreeTrialDays      json.RawMessage `json:"free_trial_days,omitempty"`
		FreeTrialDaysAlt   json.RawMessage `json:"FreeTrialDays,omitempty"`
		ConcurrentBuild    json.RawMessage `json:"concurrent_build,omitempty"`
		ConcurrentBuildAlt json.RawMessage `json:"ConcurrentBuild,omitempty"`
	}

	var wire planWire
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}

	p.ID = coalesceNonEmpty(wire.ID.String(), wire.IDAlt.String())
	p.UUID = coalesceNonEmpty(wire.UUID, wire.UUIDAlt)
	p.Name = coalesceNonEmpty(wire.Name, wire.NameAlt)
	p.Description = coalesceNonEmpty(wire.Description, wire.DescriptionAlt)
	if value, err := decodeFlexibleFloat(firstNonEmptyRaw(wire.Price, wire.PriceAlt)); err == nil {
		p.Price = value
	}
	p.Currency = coalesceNonEmpty(wire.Currency, wire.CurrencyAlt)
	p.Interval = coalesceNonEmpty(wire.Interval, wire.Period, wire.PeriodAlt)
	p.Period = coalesceNonEmpty(wire.Period, wire.PeriodAlt, wire.Interval)
	p.Active = firstBool(wire.Active, wire.ActiveAlt)
	if value, err := decodeFlexibleInt(firstNonEmptyRaw(wire.FreeTrialDays, wire.FreeTrialDaysAlt)); err == nil {
		p.FreeTrialDays = value
	}
	if value, err := decodeFlexibleInt(firstNonEmptyRaw(wire.ConcurrentBuild, wire.ConcurrentBuildAlt)); err == nil {
		p.ConcurrentBuild = value
	}

	return nil
}

func (r *CardsResponse) UnmarshalJSON(data []byte) error {
	envelope, err := parseBillingEnvelope(data)
	if err != nil {
		return err
	}

	r.Status = envelopeStatus(envelope)
	r.Message = envelope.Message
	if isNullData(envelope.Data) {
		return nil
	}

	var cards []Card
	if err := json.Unmarshal(envelope.Data, &cards); err == nil {
		r.Data.Cards = cards
		return nil
	}

	var wrapped struct {
		Cards    []Card `json:"cards,omitempty"`
		CardsAlt []Card `json:"Cards,omitempty"`
	}
	if err := json.Unmarshal(envelope.Data, &wrapped); err != nil {
		return err
	}
	if len(wrapped.Cards) > 0 {
		r.Data.Cards = wrapped.Cards
		return nil
	}
	r.Data.Cards = wrapped.CardsAlt
	return nil
}

func (r *CardResponse) UnmarshalJSON(data []byte) error {
	envelope, err := parseBillingEnvelope(data)
	if err != nil {
		return err
	}

	r.Status = envelopeStatus(envelope)
	r.Message = envelope.Message
	if isNullData(envelope.Data) {
		return nil
	}

	var wrapped struct {
		Card           *Card  `json:"card,omitempty"`
		CardAlt        *Card  `json:"Card,omitempty"`
		CheckoutURL    string `json:"checkout_url,omitempty"`
		CheckoutURLAlt string `json:"CheckoutURL,omitempty"`
		Message        string `json:"message,omitempty"`
		MessageAlt     string `json:"Message,omitempty"`
	}
	if err := json.Unmarshal(envelope.Data, &wrapped); err == nil {
		if wrapped.Card != nil {
			r.Data.Card = *wrapped.Card
		} else if wrapped.CardAlt != nil {
			r.Data.Card = *wrapped.CardAlt
		}
		r.Data.CheckoutURL = coalesceNonEmpty(wrapped.CheckoutURL, wrapped.CheckoutURLAlt)
		r.Data.Message = coalesceNonEmpty(wrapped.Message, wrapped.MessageAlt)
		if r.Data.Card.UUID != "" || r.Data.Card.ID != "" || r.Data.CheckoutURL != "" {
			return nil
		}
	}

	var card Card
	if err := json.Unmarshal(envelope.Data, &card); err == nil {
		r.Data.Card = card
	}
	return nil
}

func (r *SubscriptionsResponse) UnmarshalJSON(data []byte) error {
	envelope, err := parseBillingEnvelope(data)
	if err != nil {
		return err
	}

	r.Status = envelopeStatus(envelope)
	r.Message = envelope.Message
	if isNullData(envelope.Data) {
		return nil
	}

	var subscriptions []Subscription
	if err := json.Unmarshal(envelope.Data, &subscriptions); err == nil {
		r.Data.Subscriptions = subscriptions
		return nil
	}

	var wrapped struct {
		Subscriptions    []Subscription `json:"subscriptions,omitempty"`
		SubscriptionsAlt []Subscription `json:"Subscriptions,omitempty"`
	}
	if err := json.Unmarshal(envelope.Data, &wrapped); err != nil {
		return err
	}
	if len(wrapped.Subscriptions) > 0 {
		r.Data.Subscriptions = wrapped.Subscriptions
		return nil
	}
	r.Data.Subscriptions = wrapped.SubscriptionsAlt
	return nil
}

func (r *SubscriptionResponse) UnmarshalJSON(data []byte) error {
	envelope, err := parseBillingEnvelope(data)
	if err != nil {
		return err
	}

	r.Status = envelopeStatus(envelope)
	r.Message = envelope.Message
	if isNullData(envelope.Data) {
		return nil
	}

	var wrapped struct {
		Subscription     *Subscription  `json:"subscription,omitempty"`
		SubscriptionAlt  *Subscription  `json:"Subscription,omitempty"`
		Subscriptions    []Subscription `json:"subscriptions,omitempty"`
		SubscriptionsAlt []Subscription `json:"Subscriptions,omitempty"`
		CheckoutURL      string         `json:"checkout_url,omitempty"`
		CheckoutURLAlt   string         `json:"CheckoutURL,omitempty"`
		Message          string         `json:"message,omitempty"`
		MessageAlt       string         `json:"Message,omitempty"`
	}
	if err := json.Unmarshal(envelope.Data, &wrapped); err == nil {
		if wrapped.Subscription != nil {
			r.Data.Subscription = *wrapped.Subscription
		} else if wrapped.SubscriptionAlt != nil {
			r.Data.Subscription = *wrapped.SubscriptionAlt
		}
		if len(wrapped.Subscriptions) > 0 {
			r.Data.Subscriptions = wrapped.Subscriptions
		} else {
			r.Data.Subscriptions = wrapped.SubscriptionsAlt
		}
		r.Data.CheckoutURL = coalesceNonEmpty(wrapped.CheckoutURL, wrapped.CheckoutURLAlt)
		r.Data.Message = coalesceNonEmpty(wrapped.Message, wrapped.MessageAlt)
		if !isEmptySubscription(r.Data.Subscription) || len(r.Data.Subscriptions) > 0 || r.Data.CheckoutURL != "" {
			return nil
		}
	}

	var subscription Subscription
	if err := json.Unmarshal(envelope.Data, &subscription); err == nil && !isEmptySubscription(subscription) {
		r.Data.Subscription = subscription
		return nil
	}

	var subscriptions []Subscription
	if err := json.Unmarshal(envelope.Data, &subscriptions); err == nil {
		r.Data.Subscriptions = subscriptions
	}
	return nil
}

func (r *PlansResponse) UnmarshalJSON(data []byte) error {
	envelope, err := parseBillingEnvelope(data)
	if err != nil {
		return err
	}

	r.Status = envelopeStatus(envelope)
	r.Message = envelope.Message
	if isNullData(envelope.Data) {
		return nil
	}

	var plans []Plan
	if err := json.Unmarshal(envelope.Data, &plans); err == nil {
		r.Data.Plans = plans
		return nil
	}

	var wrapped struct {
		Plans    []Plan `json:"plans,omitempty"`
		PlansAlt []Plan `json:"Plans,omitempty"`
	}
	if err := json.Unmarshal(envelope.Data, &wrapped); err != nil {
		return err
	}
	if len(wrapped.Plans) > 0 {
		r.Data.Plans = wrapped.Plans
		return nil
	}
	r.Data.Plans = wrapped.PlansAlt
	return nil
}

func (r *BalanceResponse) UnmarshalJSON(data []byte) error {
	envelope, err := parseBillingEnvelope(data)
	if err != nil {
		return err
	}

	r.Status = envelopeStatus(envelope)
	r.Message = envelope.Message
	if isNullData(envelope.Data) {
		return nil
	}

	var payload struct {
		Balance     json.RawMessage `json:"balance,omitempty"`
		BalanceAlt  json.RawMessage `json:"Balance,omitempty"`
		Currency    string          `json:"currency,omitempty"`
		CurrencyAlt string          `json:"Currency,omitempty"`
	}
	if err := json.Unmarshal(envelope.Data, &payload); err != nil {
		return err
	}
	if value, err := decodeFlexibleFloat(firstNonEmptyRaw(payload.Balance, payload.BalanceAlt)); err == nil {
		r.Data.Balance = value
	}
	r.Data.Currency = coalesceNonEmpty(payload.Currency, payload.CurrencyAlt)
	return nil
}

func (r *PortalResponse) UnmarshalJSON(data []byte) error {
	envelope, err := parseBillingEnvelope(data)
	if err != nil {
		return err
	}

	r.Status = envelopeStatus(envelope)
	r.Message = envelope.Message
	if isNullData(envelope.Data) {
		return nil
	}

	var url string
	if err := json.Unmarshal(envelope.Data, &url); err == nil {
		r.Data.PortalURL = url
		return nil
	}

	var wrapped struct {
		PortalURL    string `json:"portal_url,omitempty"`
		PortalURLAlt string `json:"PortalURL,omitempty"`
	}
	if err := json.Unmarshal(envelope.Data, &wrapped); err != nil {
		return err
	}
	r.Data.PortalURL = coalesceNonEmpty(wrapped.PortalURL, wrapped.PortalURLAlt)
	return nil
}

func parseBillingEnvelope(data []byte) (billingEnvelope, error) {
	var envelope billingEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return billingEnvelope{}, err
	}
	return envelope, nil
}

func envelopeStatus(envelope billingEnvelope) string {
	if envelope.Status != "" {
		return envelope.Status
	}
	if envelope.Success != nil {
		return statusFromSuccess(*envelope.Success)
	}
	return ""
}

func isNullData(data json.RawMessage) bool {
	trimmed := strings.TrimSpace(string(data))
	return trimmed == "" || trimmed == "null"
}

func firstNonEmptyRaw(values ...json.RawMessage) json.RawMessage {
	for _, value := range values {
		if !isNullData(value) {
			return value
		}
	}
	return nil
}

func firstTimestamp(values ...*Timestamp) *Timestamp {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

func firstBool(values ...*bool) bool {
	for _, value := range values {
		if value != nil {
			return *value
		}
	}
	return false
}

func isEmptySubscription(subscription Subscription) bool {
	return subscription.UUID == "" && subscription.PlanName == "" && subscription.PlanTier == "" && subscription.Status == ""
}

func decodeFlexibleFloat(data json.RawMessage) (float64, error) {
	if isNullData(data) {
		return 0, nil
	}

	var floatValue float64
	if err := json.Unmarshal(data, &floatValue); err == nil {
		return floatValue, nil
	}

	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err == nil {
		if stringValue == "" {
			return 0, nil
		}
		return strconv.ParseFloat(stringValue, 64)
	}

	var numberValue json.Number
	if err := json.Unmarshal(data, &numberValue); err == nil {
		return numberValue.Float64()
	}

	return 0, nil
}

func decodeFlexibleInt(data json.RawMessage) (int, error) {
	if isNullData(data) {
		return 0, nil
	}

	var intValue int
	if err := json.Unmarshal(data, &intValue); err == nil {
		return intValue, nil
	}

	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err == nil {
		if stringValue == "" {
			return 0, nil
		}
		parsed, err := strconv.Atoi(stringValue)
		if err != nil {
			return 0, err
		}
		return parsed, nil
	}

	var numberValue json.Number
	if err := json.Unmarshal(data, &numberValue); err == nil {
		parsed, err := strconv.Atoi(numberValue.String())
		if err == nil {
			return parsed, nil
		}
		floatValue, err := numberValue.Float64()
		if err != nil {
			return 0, err
		}
		return int(floatValue), nil
	}

	return 0, nil
}
