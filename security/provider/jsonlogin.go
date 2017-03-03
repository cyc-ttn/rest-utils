package provider

import (
  "net/http"

  "luxe.technology/rest-utils/controller"
  "luxe.technology/rest-utils/security"
)

// JSONLoginAuthenticationProvider is an implementation of AuthenticationProvider
// which takes form inputs
type JSONLoginAuthenticationProvider struct {
  userIDField       string
  passwordField     string
}

// NewJSONLoginAuthenticationProvider creates a new FormLoginAuthenticationProvider
func NewJSONLoginAuthenticationProvider() *JSONLoginAuthenticationProvider{
  return &JSONLoginAuthenticationProvider{
    userIDField: "Username",
    passwordField: "Password",
  }
}

// Authenticate authenticates based on the form login.
// provides a UserDetails object if successful
func (p * JSONLoginAuthenticationProvider) Authenticate (
  req * http.Request,
  usrSrv security.UserDetailsService,
) security.UserDetails {

  // Get details from http.Request
  body, err := controller.JSONFromRequest(req)
  if err != nil { return nil }

  idAsIf, ok := body[ p.userIDField ]
  if !ok { return nil }

  id, ok := idAsIf.(string)
  if !ok{ return nil }

  password, ok := body[ p.passwordField ]
  if !ok { return nil }

  details := usrSrv.LoadByUsername(id)
  if details == nil || !details.Verify(password) {
    return nil
  }

  return details
}
