package pipeops

// TeamService handles communication with the team related
// methods of the PipeOps API.
type TeamService struct {
	client *Client
}

// WorkspaceService handles communication with the workspace related
// methods of the PipeOps API.
type WorkspaceService struct {
	client *Client
}

// BillingService handles communication with the billing related
// methods of the PipeOps API.
type BillingService struct {
	client *Client
}

// AddOnService handles communication with the add-on related
// methods of the PipeOps API.
type AddOnService struct {
	client *Client
}

// WebhookService handles communication with the webhook related
// methods of the PipeOps API.
type WebhookService struct {
	client *Client
}

// UserService handles communication with the user settings related
// methods of the PipeOps API.
type UserService struct {
	client *Client
}

// AdminService handles communication with the admin related
// methods of the PipeOps API.
type AdminService struct {
	client *Client
}
