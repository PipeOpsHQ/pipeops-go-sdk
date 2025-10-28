package pipeops

import (
	"context"
	"fmt"
	"net/http"
)

// CloudProviderService handles communication with cloud provider related
// methods of the PipeOps API.
type CloudProviderService struct {
	client *Client
}

// AWS Cloud Provider Methods

// AWSAccount represents an AWS account configuration.
type AWSAccount struct {
	ID            string     `json:"id,omitempty"`
	UUID          string     `json:"uuid,omitempty"`
	AccessKeyID   string     `json:"access_key_id,omitempty"`
	SecretKey     string     `json:"secret_key,omitempty"`
	Region        string     `json:"region,omitempty"`
	WorkspaceUUID string     `json:"workspace_uuid,omitempty"`
	CreatedAt     *Timestamp `json:"created_at,omitempty"`
}

// AWSAccountRequest represents a request to add an AWS account.
type AWSAccountRequest struct {
	AccessKeyID string `json:"access_key_id"`
	SecretKey   string `json:"secret_key"`
	Region      string `json:"region"`
}

// AWSAccountResponse represents AWS account response.
type AWSAccountResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Account AWSAccount `json:"account"`
	} `json:"data"`
}

// AddAWSAccount adds a new AWS account.
func (s *CloudProviderService) AddAWSAccount(ctx context.Context, req *AWSAccountRequest) (*AWSAccountResponse, *http.Response, error) {
	u := "aws/add_account"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	accountResp := new(AWSAccountResponse)
	resp, err := s.client.Do(ctx, httpReq, accountResp)
	if err != nil {
		return nil, resp, err
	}

	return accountResp, resp, nil
}

// DisconnectAWSAccount disconnects an AWS account.
func (s *CloudProviderService) DisconnectAWSAccount(ctx context.Context, accountUUID string) (*http.Response, error) {
	u := fmt.Sprintf("aws/disconnect/%s", accountUUID)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// DeleteAWSAccount deletes an AWS account.
func (s *CloudProviderService) DeleteAWSAccount(ctx context.Context, accountUUID string) (*http.Response, error) {
	u := fmt.Sprintf("aws/%s", accountUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// GCP Cloud Provider Methods

// GCPCredentialRequest represents a request to upload GCP credentials.
type GCPCredentialRequest struct {
	CredentialsJSON string `json:"credentials_json"`
}

// GCPAccountResponse represents GCP account response.
type GCPAccountResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Account map[string]interface{} `json:"account"`
	} `json:"data"`
}

// UploadGCPCredential uploads GCP service account credentials.
func (s *CloudProviderService) UploadGCPCredential(ctx context.Context, workspaceUUID string, req *GCPCredentialRequest) (*GCPAccountResponse, *http.Response, error) {
	u := fmt.Sprintf("gcp/%s/upload-credential", workspaceUUID)

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	accountResp := new(GCPAccountResponse)
	resp, err := s.client.Do(ctx, httpReq, accountResp)
	if err != nil {
		return nil, resp, err
	}

	return accountResp, resp, nil
}

// DeleteGCPAccount deletes a GCP account.
func (s *CloudProviderService) DeleteGCPAccount(ctx context.Context, accountUUID string) (*http.Response, error) {
	u := fmt.Sprintf("gcp/%s", accountUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// Azure Cloud Provider Methods

// AzureCredentialRequest represents a request to add Azure credentials.
type AzureCredentialRequest struct {
	SubscriptionID string `json:"subscription_id"`
	TenantID       string `json:"tenant_id"`
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
}

// AzureAccountResponse represents Azure account response.
type AzureAccountResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Account map[string]interface{} `json:"account"`
	} `json:"data"`
}

// AddAzureAccount adds Azure cloud credentials.
func (s *CloudProviderService) AddAzureAccount(ctx context.Context, req *AzureCredentialRequest) (*AzureAccountResponse, *http.Response, error) {
	u := "azure/add-account"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	accountResp := new(AzureAccountResponse)
	resp, err := s.client.Do(ctx, httpReq, accountResp)
	if err != nil {
		return nil, resp, err
	}

	return accountResp, resp, nil
}

// DeleteAzureAccount deletes an Azure account.
func (s *CloudProviderService) DeleteAzureAccount(ctx context.Context, accountUUID string) (*http.Response, error) {
	u := fmt.Sprintf("azure/%s", accountUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// DigitalOcean Cloud Provider Methods

// DigitalOceanAccountRequest represents a request to add DigitalOcean credentials.
type DigitalOceanAccountRequest struct {
	Token string `json:"token"`
}

// DigitalOceanAccountResponse represents DigitalOcean account response.
type DigitalOceanAccountResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Account map[string]interface{} `json:"account"`
	} `json:"data"`
}

// AddDigitalOceanAccount adds DigitalOcean credentials.
func (s *CloudProviderService) AddDigitalOceanAccount(ctx context.Context, req *DigitalOceanAccountRequest) (*DigitalOceanAccountResponse, *http.Response, error) {
	u := "digitalocean/add-account"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	accountResp := new(DigitalOceanAccountResponse)
	resp, err := s.client.Do(ctx, httpReq, accountResp)
	if err != nil {
		return nil, resp, err
	}

	return accountResp, resp, nil
}

// DeleteDigitalOceanAccount deletes a DigitalOcean account.
func (s *CloudProviderService) DeleteDigitalOceanAccount(ctx context.Context, accountUUID string) (*http.Response, error) {
	u := fmt.Sprintf("digitalocean/%s", accountUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// Huawei Cloud Provider Methods

// HuaweiAccountRequest represents a request to add Huawei credentials.
type HuaweiAccountRequest struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Region    string `json:"region"`
}

// HuaweiAccountResponse represents Huawei account response.
type HuaweiAccountResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Account map[string]interface{} `json:"account"`
	} `json:"data"`
}

// AddHuaweiAccount adds Huawei cloud credentials.
func (s *CloudProviderService) AddHuaweiAccount(ctx context.Context, req *HuaweiAccountRequest) (*HuaweiAccountResponse, *http.Response, error) {
	u := "huawei/add-account"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	accountResp := new(HuaweiAccountResponse)
	resp, err := s.client.Do(ctx, httpReq, accountResp)
	if err != nil {
		return nil, resp, err
	}

	return accountResp, resp, nil
}

// DeleteHuaweiAccount deletes a Huawei account.
func (s *CloudProviderService) DeleteHuaweiAccount(ctx context.Context, accountUUID string) (*http.Response, error) {
	u := fmt.Sprintf("huawei/%s", accountUUID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// AWS Calculator Methods

// EC2CalculatorRequest represents an EC2 cost calculator request.
type EC2CalculatorRequest struct {
	InstanceType string `json:"instance_type"`
	Region       string `json:"region"`
	Hours        int    `json:"hours,omitempty"`
}

// CalculatorResponse represents a calculator response.
type CalculatorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Cost float64 `json:"cost"`
	} `json:"data"`
}

// CalculateEC2Cost calculates EC2 costs.
func (s *CloudProviderService) CalculateEC2Cost(ctx context.Context, req *EC2CalculatorRequest) (*CalculatorResponse, *http.Response, error) {
	u := "aws/ec2-calculator"

	httpReq, err := s.client.NewRequest(http.MethodPost, u, req)
	if err != nil {
		return nil, nil, err
	}

	calcResp := new(CalculatorResponse)
	resp, err := s.client.Do(ctx, httpReq, calcResp)
	if err != nil {
		return nil, resp, err
	}

	return calcResp, resp, nil
}

// GetAWSReference retrieves AWS reference data.
func (s *CloudProviderService) GetAWSReference(ctx context.Context) (*http.Response, error) {
	u := "aws/reference"

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

// ELBCalculatorRequest represents ELB cost calculator request.
type ELBCalculatorRequest struct {
	LoadBalancerType string `json:"load_balancer_type"`
	Region           string `json:"region"`
	Hours            int    `json:"hours,omitempty"`
}

// CalculateELBCost calculates ELB costs.
func (s *CloudProviderService) CalculateELBCost(ctx context.Context, req *ELBCalculatorRequest) (*CalculatorResponse, *http.Response, error) {
	u := "aws/elb-calculator"

	httpReq, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	calcResp := new(CalculatorResponse)
	resp, err := s.client.Do(ctx, httpReq, calcResp)
	if err != nil {
		return nil, resp, err
	}

	return calcResp, resp, nil
}

// EBSCalculatorRequest represents EBS cost calculator request.
type EBSCalculatorRequest struct {
	VolumeType string `json:"volume_type"`
	SizeGB     int    `json:"size_gb"`
	Region     string `json:"region"`
}

// CalculateEBSCost calculates EBS costs.
func (s *CloudProviderService) CalculateEBSCost(ctx context.Context, req *EBSCalculatorRequest) (*CalculatorResponse, *http.Response, error) {
	u := "aws/ebs-calculator"

	httpReq, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	calcResp := new(CalculatorResponse)
	resp, err := s.client.Do(ctx, httpReq, calcResp)
	if err != nil {
		return nil, resp, err
	}

	return calcResp, resp, nil
}
