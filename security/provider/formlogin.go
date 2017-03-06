package provider

import (
	"net/http"

	"luxe.technology/rest-utils/security"
)

// FormLoginAuthenticationProvider is an implementation of AuthenticationProvider
// which takes form inputs
type FormLoginAuthenticationProvider struct {
}

// NewFormLoginAuthenticationProvider creates a new FormLoginAuthenticationProvider
func NewFormLoginAuthenticationProvider() *FormLoginAuthenticationProvider {
	return &FormLoginAuthenticationProvider{}
}

// Authenticate authenticates based on the form login.
// provides a UserDetails object if successful
func (p *FormLoginAuthenticationProvider) Authenticate(
	req *http.Request,
	usrSrv security.UserDetailsService,
) (*security.UserDetails, map[string]interface{}, error) {

	// Get details from http.Request

	return nil, nil, nil
}
