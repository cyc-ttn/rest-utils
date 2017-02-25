package provider

import (
  "net/http"

  "luxe.technology/rest-utils/security"
)

// TokenAuthenticationProvider is an implementation of AuthenticationProvider
// which determines whether a token is part of the request
type TokenAuthenticationProvider struct {

}

// NewTokenAuthenticationProvider creates a new FormLoginAuthenticationProvider
func NewTokenAuthenticationProvider() *TokenAuthenticationProvider{
  return &TokenAuthenticationProvider{}
}

// Authenticate authenticates based on the token
// provides a UserDetails object if successful
func (p * TokenAuthenticationProvider) Authenticate (
  req * http.Request,
  usrSrv security.UserDetailsService,
) security.UserDetails {

  // Get details from http.Request

  return nil
}
