package services

import (
	"context"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"net/url"
)

// AuthService provides functions for checking if a user is authenticated
// or extracting user info from it.
type AuthService struct {
	serviceUrl url.URL
	audience   string
}

// NewAuthService creates a new AuthenticationService that asks the auth0 url
func NewAuthService(auth0 url.URL, audience string) *AuthService {
	return &AuthService{
		serviceUrl: auth0,
		audience:   audience,
	}
}

// GetAuthMiddleware returns a middleware that checks if a JWT Token is present and valid according to the given secret.
// Can be used to enforce Authentication on routes that need it.
func (s *AuthService) GetAuthMiddleware(secret []byte) (*jwtmiddleware.JWTMiddleware, error) {
	keyFunc := func(ctx context.Context) (interface{}, error) {
		return secret, nil
	}

	jwtValidator, err := validator.New(
		keyFunc,
		validator.HS256,
		s.serviceUrl.String(),
		[]string{s.audience},
	)
	if err != nil {
		return nil, err
	}

	middleware := jwtmiddleware.New(jwtValidator.ValidateToken)
	return middleware, nil
}
