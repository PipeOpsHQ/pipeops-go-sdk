package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
)

func main() {
	// Create a new client
	client, err := pipeops.NewClient("")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// OAuth 2.0 Authorization Code Flow Example

	// Step 1: Generate authorization URL
	authURL, err := client.OAuth.Authorize(&pipeops.AuthorizeOptions{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		RedirectURI:  "http://localhost:3000/callback",
		ResponseType: "code",
		Scope:        "user:read user:write",
		State:        "random_state_string",
	})
	if err != nil {
		log.Fatalf("Failed to generate auth URL: %v", err)
	}

	fmt.Printf("Please visit this URL to authorize:\n%s\n\n", authURL)
	fmt.Print("Enter the authorization code from the callback: ")

	var authCode string
	if _, err := fmt.Scanln(&authCode); err != nil {
		log.Fatalf("Failed to read authorization code: %v", err)
	}

	// Step 2: Exchange authorization code for access token
	tokenResp, _, err := client.OAuth.ExchangeCodeForToken(ctx, &pipeops.TokenRequest{
		GrantType:    "authorization_code",
		Code:         authCode,
		RedirectURI:  "http://localhost:3000/callback",
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
	})
	if err != nil {
		log.Fatalf("Failed to exchange code for token: %v", err)
	}

	fmt.Printf("Access Token: %s\n", tokenResp.AccessToken)
	fmt.Printf("Token Type: %s\n", tokenResp.TokenType)
	fmt.Printf("Expires In: %d seconds\n", tokenResp.ExpiresIn)

	// Set the OAuth access token for API requests
	client.SetToken(tokenResp.AccessToken)

	// Step 3: Get user information
	userInfo, _, err := client.OAuth.GetUserInfo(ctx)
	if err != nil {
		log.Fatalf("Failed to get user info: %v", err)
	}

	fmt.Printf("\nUser Information:\n")
	fmt.Printf("  Email: %s\n", userInfo.Data.Email)
	fmt.Printf("  Name: %s\n", userInfo.Data.Name)
	fmt.Printf("  Email Verified: %t\n", userInfo.Data.EmailVerified)

	// Optional: Get consent page (if implemented)
	consent, _, err := client.OAuth.GetConsent(ctx)
	if err != nil {
		log.Printf("Consent endpoint: %v", err)
	} else {
		fmt.Printf("\nConsent: %s\n", consent.Message)
	}
}
