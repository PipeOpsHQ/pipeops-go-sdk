# Billing Service

The Billing Service manages PipeOps balance, subscriptions, plans, portal access, and workspace payment cards.

## Overview

```go
billingService := client.Billing
```

> Note: the current controller-backed billing reads verified against the PipeOps controller collection are `GetBillingInfo`, `GetBalance`, `GetCurrentSubscription`, `GetPortalURL`, `GetPlans`, `ListWorkspaceCards`, and `GetActiveCard`. The legacy `GetUsage` method still targets `/billing/usage`, which is not exposed by the current controller.

## Key Methods

### Get Billing Info

Retrieve current balance plus the active subscription snapshot in one call:

```go
info, _, err := client.Billing.GetBillingInfo(ctx)
if err != nil {
    log.Fatalf("Failed to get billing info: %v", err)
}

fmt.Printf("Balance: %.2f %s\n", info.Data.Balance.Balance, info.Data.Balance.Currency)
if info.Data.CurrentSubscription != nil {
    fmt.Printf("Plan: %s (%s)\n", info.Data.CurrentSubscription.PlanName, info.Data.CurrentSubscription.PlanTier)
}
```

### Get Balance

```go
balance, _, err := client.Billing.GetBalance(ctx)
if err != nil {
    log.Fatalf("Failed to get balance: %v", err)
}

fmt.Printf("Balance: %.2f %s\n", balance.Data.Balance, balance.Data.Currency)
```

### Get Current Subscription

```go
subscription, _, err := client.Billing.GetCurrentSubscription(ctx)
if err != nil {
    log.Fatalf("Failed to get current subscription: %v", err)
}

fmt.Printf("Plan: %s (%s)\n", subscription.Data.Subscription.PlanName, subscription.Data.Subscription.PlanTier)
fmt.Printf("Status: %s\n", subscription.Data.Subscription.Status)
```

### Get Plans

```go
plans, _, err := client.Billing.GetPlans(ctx)
if err != nil {
    log.Fatalf("Failed to get plans: %v", err)
}

for _, plan := range plans.Data.Plans {
    fmt.Printf("%s: %.2f %s (%s)\n", plan.Name, plan.Price, plan.Currency, plan.Period)
}
```

### Get Billing Portal URL

```go
portal, _, err := client.Billing.GetPortalURL(ctx)
if err != nil {
    log.Fatalf("Failed to get portal URL: %v", err)
}

fmt.Println(portal.Data.PortalURL)
```

### List Workspace Cards

```go
cards, _, err := client.Billing.ListWorkspaceCards(ctx)
if err != nil {
    log.Fatalf("Failed to list workspace cards: %v", err)
}

for _, card := range cards.Data.Cards {
    fmt.Printf("%s ending in %s\n", card.Brand, card.Last4)
}
```

### Get Active Card

```go
card, _, err := client.Billing.GetActiveCard(ctx)
if err != nil {
    log.Fatalf("Failed to get active card: %v", err)
}

fmt.Printf("Active card: %s ending in %s\n", card.Data.Card.Brand, card.Data.Card.Last4)
```

### Subscribe To A Plan

```go
resp, _, err := client.Billing.Subscribe(ctx, &pipeops.SubscribeRequest{PlanID: "startup"})
if err != nil {
    log.Fatalf("Failed to subscribe: %v", err)
}

fmt.Println(resp.Data.CheckoutURL)
```

### Start Trial

```go
_, err = client.Billing.StartTrial(ctx, &pipeops.StartTrialRequest{PlanID: "startup_v1"})
if err != nil {
    log.Fatalf("Failed to start trial: %v", err)
}
```

## See Also

- [Workspaces Service](workspaces.md)
