package main

import (
	"context"
	"fmt"

	"google.golang.org/api/idtoken"
)

type googleVerifier struct{}

// VerifiedEmail validates the token's signature, expiry and audience against
// Google's public keys and returns its email if Google has verified it.
func (googleVerifier) VerifiedEmail(ctx context.Context, token, audience string) (string, error) {
	payload, err := idtoken.Validate(ctx, token, audience)
	if err != nil {
		return "", fmt.Errorf("validating ID token: %w", err)
	}
	if verified, _ := payload.Claims["email_verified"].(bool); !verified {
		return "", fmt.Errorf("email in ID token is not verified")
	}
	email, _ := payload.Claims["email"].(string)
	return email, nil
}
