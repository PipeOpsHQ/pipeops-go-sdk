# Webhooks Service

The Webhooks Service manages webhook configuration and delivery.

## Overview

```go
// Access the webhooks service
webhooksService := client.Webhooks
```

## Methods

### List Webhooks

List all webhooks:

```go
webhooks, _, err := client.Webhooks.List(ctx)
if err != nil {
    log.Fatalf("Failed to list webhooks: %v", err)
}

for _, webhook := range webhooks.Data.Webhooks {
    fmt.Printf("- %s: %s\n", webhook.Name, webhook.URL)
}
```

### Create Webhook

Create a new webhook:

```go
webhook, _, err := client.Webhooks.Create(ctx, &pipeops.CreateWebhookRequest{
    URL:    "https://myapp.com/webhook",
    Events: []string{"project.deployed", "project.failed", "project.started"},
    Secret: "webhook-secret",
})
if err != nil {
    log.Fatalf("Failed to create webhook: %v", err)
}

fmt.Printf("Created webhook: %s\n", webhook.Data.Webhook.UUID)
```

### Get Webhook

Get webhook details:

```go
webhook, _, err := client.Webhooks.Get(ctx, "webhook-uuid")
if err != nil {
    log.Fatalf("Failed to get webhook: %v", err)
}

fmt.Printf("Webhook: %s\n", webhook.Data.Webhook.URL)
```

### Update Webhook

Update webhook configuration:

```go
updated, _, err := client.Webhooks.Update(ctx, webhookUUID, &pipeops.UpdateWebhookRequest{
    URL:    "https://myapp.com/new-webhook",
    Events: []string{"project.deployed"},
})
```

### Delete Webhook

Delete a webhook:

```go
_, err := client.Webhooks.Delete(ctx, "webhook-uuid")
```

### Test Webhook

Test webhook delivery:

```go
_, err := client.Webhooks.TestWebhook(ctx, "webhook-uuid")
if err != nil {
    log.Fatalf("Test failed: %v", err)
}

fmt.Println("Test webhook sent")
```

### Get Webhook Deliveries

Get webhook delivery history:

```go
deliveries, _, err := client.Webhooks.GetWebhookDeliveries(ctx, "webhook-uuid")
if err != nil {
    log.Fatalf("Failed to get deliveries: %v", err)
}

for _, delivery := range deliveries.Data.Deliveries {
    fmt.Printf("Delivery: %s - Status: %d\n", delivery.ID, delivery.StatusCode)
}
```

### Retry Webhook Delivery

Retry a failed webhook delivery:

```go
_, err := client.Webhooks.RetryWebhookDelivery(ctx, webhookUUID, deliveryID)
```

## Data Types

```go
type Webhook struct {
    ID     string   `json:"id,omitempty"`
    UUID   string   `json:"uuid,omitempty"`
    URL    string   `json:"url,omitempty"`
    Events []string `json:"events,omitempty"`
    Secret string   `json:"secret,omitempty"`
}
```

## Available Events

- `project.deployed` - Project deployment completed
- `project.failed` - Project deployment failed
- `project.started` - Project deployment started
- `project.stopped` - Project stopped
- `project.restarted` - Project restarted
- `project.deleted` - Project deleted

## See Also

- [Projects Service](projects.md)
