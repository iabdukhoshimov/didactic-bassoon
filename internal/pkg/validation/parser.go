package validation

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"gitlab.com/tour/internal/pkg/security"
)

func ExtractUserIDFromToken(ctx context.Context, r *http.Request, issuer security.IssuerInterface) (string, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, "Bearer"))

	claims, err := issuer.ParseAccessToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("error parsing token: %w", err)
	}

	userID := fmt.Sprintf("%v", claims.Subject)
	if userID == "" {
		return "", fmt.Errorf("user ID not found in token claims")
	}

	return userID, nil
}
