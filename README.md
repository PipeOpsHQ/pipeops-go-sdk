# PipeOps Go SDK

A Go SDK for interacting with the PipeOps Control Plane API.

## Installation

```bash
go get github.com/PipeOpsHQ/pipeops-go-sdk
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func main() {
    // Create a new client
    client := pipeops.NewClient("https://api.pipeops.io")
    
    // Login to get an authentication token
    ctx := context.Background()
    loginReq := &pipeops.LoginRequest{
        Email:    "your-email@example.com",
        Password: "your-password",
    }
    
    resp, err := client.Auth.Login(ctx, loginReq)
    if err != nil {
        log.Fatal(err)
    }
    
    // Set the token for authenticated requests
    client.SetToken(resp.Data.Token)
    
    // Now you can make authenticated API calls
    // ...
}
```

## Features

This SDK provides access to the following PipeOps API categories:

- Authentication
- Projects
- Deployments
- Servers/Clusters
- Environments
- Teams
- Workspaces
- Billing
- Cloud Providers (AWS, Azure, GCP, DigitalOcean, Huawei)
- Add-Ons
- Webhooks
- User Settings
- Admin Operations
- And more...

## Documentation

For detailed API documentation, please refer to the [PipeOps API Documentation](https://api.pipeops.io/docs).

## License

This SDK is distributed under the terms of the license specified in the LICENSE file.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
