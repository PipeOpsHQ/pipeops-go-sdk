// Example: list and create a sandbox via the PipeOps BFF.
//
//	export PIPEOPS_TOKEN=sat_...   # or user JWT
//	export PIPEOPS_WORKSPACE_UUID=...
//	go run ./examples/sandboxes
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func main() {
	token := os.Getenv("PIPEOPS_TOKEN")
	ws := os.Getenv("PIPEOPS_WORKSPACE_UUID")
	if token == "" || ws == "" {
		log.Fatal("set PIPEOPS_TOKEN and PIPEOPS_WORKSPACE_UUID")
	}

	client, err := pipeops.NewClient(os.Getenv("PIPEOPS_BASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	client.SetToken(token)

	ctx := context.Background()
	opts := &pipeops.SandboxWorkspaceOptions{WorkspaceUUID: ws}

	list, _, err := client.Sandboxes.List(ctx, opts)
	if err != nil {
		log.Fatalf("list: %v", err)
	}
	fmt.Printf("sandboxes: %d\n", len(list.Data))
	for _, s := range list.Data {
		fmt.Printf("  %s  %s  %s\n", s.ID, s.Name, s.Status)
	}

	if os.Getenv("PIPEOPS_SANDBOX_CREATE") != "1" {
		fmt.Println("set PIPEOPS_SANDBOX_CREATE=1 to create a sandbox")
		return
	}

	created, _, err := client.Sandboxes.Create(ctx, opts, &pipeops.CreateSandboxRequest{
		Name:  "sdk-example",
		Image: "ubuntu",
		Role:  "standard",
	})
	if err != nil {
		log.Fatalf("create: %v", err)
	}
	fmt.Printf("created: %s status=%s\n", created.Data.ID, created.Data.Status)
}
