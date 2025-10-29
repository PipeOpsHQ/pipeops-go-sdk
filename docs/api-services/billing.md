# Billing Service

The Billing Service manages subscriptions, payments, and invoices.

## Overview

```go
// Access the billing service
billingService := client.Billing
```

## Methods

### Get Balance

Get current account balance:

```go
balance, _, err := client.Billing.GetBalance(ctx)
if err != nil {
    log.Fatalf("Failed to get balance: %v", err)
}

fmt.Printf("Balance: $%.2f\n", balance.Data.Balance)
fmt.Printf("Credit: $%.2f\n", balance.Data.Credit)
```

### List Invoices

List all invoices:

```go
invoices, _, err := client.Billing.ListInvoices(ctx, &pipeops.InvoiceListOptions{
    Page:  1,
    Limit: 20,
})
if err != nil {
    log.Fatalf("Failed to list invoices: %v", err)
}

for _, invoice := range invoices.Data.Invoices {
    fmt.Printf("Invoice: %s - $%.2f - %s\n", 
        invoice.ID, invoice.Amount, invoice.Status)
}
```

### Get Invoice

Get specific invoice:

```go
invoice, _, err := client.Billing.GetInvoice(ctx, "invoice-uuid")
if err != nil {
    log.Fatalf("Failed to get invoice: %v", err)
}

fmt.Printf("Amount: $%.2f\n", invoice.Data.Invoice.Amount)
```

### Add Payment Card

Add a payment method:

```go
card, _, err := client.Billing.AddCard(ctx, &pipeops.AddCardRequest{
    Token: "stripe-card-token",
})
if err != nil {
    log.Fatalf("Failed to add card: %v", err)
}

fmt.Printf("Card added: ****%s\n", card.Data.Card.Last4)
```

### Get Active Card

Get the active payment card:

```go
card, _, err := client.Billing.GetActiveCard(ctx)
if err != nil {
    log.Fatalf("Failed to get card: %v", err)
}

fmt.Printf("Active card: ****%s\n", card.Data.Card.Last4)
```

### Delete Card

Remove a payment card:

```go
_, err := client.Billing.DeleteCard(ctx, "card-uuid")
```

### Add Credit

Add credit to account:

```go
credit, _, err := client.Billing.AddCredit(ctx, &pipeops.CreditRequest{
    Amount: 100.00,
})
```

### Apply Discount

Apply a discount code:

```go
_, err := client.Billing.ApplyDiscount(ctx, &pipeops.ApplyDiscountRequest{
    Code: "DISCOUNT20",
})
```

### Get Subscription

Get current subscription details:

```go
subscription, _, err := client.Billing.GetSubscription(ctx)
if err != nil {
    log.Fatalf("Failed to get subscription: %v", err)
}

fmt.Printf("Plan: %s\n", subscription.Data.Subscription.Plan)
fmt.Printf("Status: %s\n", subscription.Data.Subscription.Status)
```

### Cancel Subscription

Cancel an active subscription:

```go
_, err := client.Billing.CancelSubscription(ctx, "subscription-uuid")
```

### Export Invoices

Export invoices in various formats:

```go
export, _, err := client.Billing.ExportInvoices(ctx, &pipeops.ExportInvoicesRequest{
    Format:    "pdf",
    StartDate: "2024-01-01",
    EndDate:   "2024-12-31",
})
```

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func main() {
    client, _ := pipeops.NewClient("")
    client.SetToken("your-token")
    
    ctx := context.Background()
    
    // Get current balance
    balance, _, err := client.Billing.GetBalance(ctx)
    if err != nil {
        log.Fatalf("Failed to get balance: %v", err)
    }
    
    fmt.Printf("Account Balance: $%.2f\n", balance.Data.Balance)
    
    // List recent invoices
    invoices, _, err := client.Billing.ListInvoices(ctx, &pipeops.InvoiceListOptions{
        Limit: 5,
    })
    if err != nil {
        log.Fatalf("Failed to list invoices: %v", err)
    }
    
    fmt.Printf("\nRecent Invoices:\n")
    for _, invoice := range invoices.Data.Invoices {
        fmt.Printf("- %s: $%.2f (%s)\n", 
            invoice.Date, invoice.Amount, invoice.Status)
    }
}
```

## See Also

- [Workspaces Service](workspaces.md)
- [Admin Service](admin.md)
