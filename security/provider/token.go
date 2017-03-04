package provider

import (
  "net/http"
  "strings"

  "luxe.technology/rest-utils/security"
)

// TokenAuthenticationProvider is an implementation of AuthenticationProvider
// which determines whether a token is part of the request
type TokenAuthenticationProvider struct {
  repository    * security.UserTokenRepository
}

// NewTokenAuthenticationProvider creates a new FormLoginAuthenticationProvider
func NewTokenAuthenticationProvider(repo * security.UserTokenRepository) *TokenAuthenticationProvider{
  return &TokenAuthenticationProvider{
    repository: repo,
  }
}

// Authenticate authenticates based on the token
// provides a UserDetails object if successful
func (p * TokenAuthenticationProvider) Authenticate (
  req * http.Request,
  usrSrv security.UserDetailsService,
) security.UserDetails {

  // Get the token from the Authentication Header
  tokenStr := p.getBearerToken(req)
  if tokenStr == "" {
    return nil
  }

  // Get details from http.Request
  token, err := p.repository.FindAndVerifyToken(tokenStr)
  if err != nil {
    return nil
  }

  return token.User
}

func (p * TokenAuthenticationProvider) getBearerToken(req * http.Request) string {
  // Get the token from the Authentication Header
  authHeader, ok := req.Header["Authorization"]
  if !ok {
    return ""
  }

  for _, value := range authHeader {
    parts := strings.Split(value, " ")
    if strings.ToLower(parts[0]) == "bearer" {
      return strings.TrimSpace(parts[1])
    }
  }

  return ""
}
