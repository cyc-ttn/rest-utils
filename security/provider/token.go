package provider

import (
	"net/http"
	"strings"

	"luxe.technology/rest-utils/security"
)

// TokenAuthenticationProvider is an implementation of AuthenticationProvider
// which determines whether a token is part of the request
type TokenAuthenticationProvider struct {
	Repository *security.UserTokenRepository
}

// NewTokenAuthenticationProvider creates a new FormLoginAuthenticationProvider
func NewTokenAuthenticationProvider(repo *security.UserTokenRepository) *TokenAuthenticationProvider {
	return &TokenAuthenticationProvider{
		Repository: repo,
	}
}

// Authenticate authenticates based on the token
// provides a UserDetails object if successful
func (p *TokenAuthenticationProvider) Authenticate(
	req *http.Request,
	usrSrv security.UserDetailsService,
) (security.UserDetails, map[string]interface{}, error) {

	// Get the token from the Authentication Header
	tokenStr := GetBearerToken(req)
	if tokenStr == "" {
		return nil, nil, ErrAuthenticationInvalid
	}

	// Get details from http.Request
	token, err := p.Repository.FindAndVerifyToken(tokenStr)
	if err != nil {
		return nil, nil, err
	}

	return usrSrv.RefreshUserDetails(token.User), nil, nil
}

// GetBearerToken Returns the Bearer Token from request
func GetBearerToken(req *http.Request) string {
	// Get the token from the Authentication Header
	authHeader, ok := req.Header["Authorization"]
	if !ok {
		return ""
	}

	for _, value := range authHeader {
		parts := strings.Split(value, " ")
		if strings.ToLower(parts[0]) == "bearer" && len(parts) >= 2 {
			return strings.TrimSpace(parts[1])
		}
	}

	return ""
}
