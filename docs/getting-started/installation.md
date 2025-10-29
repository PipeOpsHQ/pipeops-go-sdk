# Installation

## Requirements

- Go 1.19 or higher
- Git (for go get)

## Install via go get

The simplest way to install the PipeOps Go SDK is using `go get`:

```bash
go get github.com/PipeOpsHQ/pipeops-go-sdk
```

This will download the SDK and its dependencies into your Go workspace.

## Install a Specific Version

To install a specific version of the SDK:

```bash
go get github.com/PipeOpsHQ/pipeops-go-sdk@v1.0.0
```

## Using Go Modules

If you're using Go modules (recommended), simply import the package in your code:

```go
import "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
```

Then run:

```bash
go mod tidy
```

This will automatically download the SDK and add it to your `go.mod` file.

## Manual Installation

You can also manually clone the repository:

```bash
git clone https://github.com/PipeOpsHQ/pipeops-go-sdk.git
cd pipeops-go-sdk
go install
```

## Verify Installation

Create a simple test file to verify the installation:

```go title="test.go"
package main

import (
    "fmt"
    "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func main() {
    client, err := pipeops.NewClient("")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Client created successfully: %T\n", client)
}
```

Run the test:

```bash
go run test.go
```

You should see output like:

```
Client created successfully: *pipeops.Client
```

## Dependencies

The SDK has minimal dependencies:

- `github.com/google/go-querystring` - URL query string encoding

All dependencies are automatically managed by Go modules.

## Updating the SDK

To update to the latest version:

```bash
go get -u github.com/PipeOpsHQ/pipeops-go-sdk
```

Or to update to a specific version:

```bash
go get -u github.com/PipeOpsHQ/pipeops-go-sdk@v1.1.0
```

## Next Steps

Now that you have the SDK installed, proceed to the [Quick Start Guide](quickstart.md) to learn how to use it.
