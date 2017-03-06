package security

import (
	"net/http"
)

// UserDetailsService provides user details for the proper user
type UserDetailsService interface {
	LoadByUsername(string) UserDetails
}

// AuthenticationManager Manages authentications
type AuthenticationManager struct {
	providers []AuthenticationProvider
	userSrv   UserDetailsService
}

// NewAuthenticationManager creates an authentication manager
func NewAuthenticationManager() *AuthenticationManager {
	return &AuthenticationManager{
		providers: make([]AuthenticationProvider, 0, 10),
	}
}

// SetUserDetailsService sets the details service
func (mgr *AuthenticationManager) SetUserDetailsService(srv UserDetailsService) *AuthenticationManager {
	mgr.userSrv = srv
	return mgr
}

// AddProvider adds an authentication provider to the list
func (mgr *AuthenticationManager) AddProvider(auth AuthenticationProvider) *AuthenticationManager {

	capacity := cap(mgr.providers)
	length := len(mgr.providers)

	if length == capacity {
		newProviders := make([]AuthenticationProvider, 0, capacity+10)
		copy(newProviders, mgr.providers)
		mgr.providers = newProviders
	}

	mgr.providers = append(mgr.providers, auth)
	return mgr
}

// Authenticate authenticates the current connection
func (mgr *AuthenticationManager) Authenticate(req *http.Request) (UserDetails, map[string]interface{}, []error) {
	errs := make([]error, 0, len(mgr.providers))

	for _, prov := range mgr.providers {
		userDetails, info, err := prov.Authenticate(req, mgr.userSrv)
		if err != nil {
			errs = append(errs, err)
		} else {
			return userDetails, info, nil
		}
	}
	return nil, nil, errs
}

// Authorize authorizes a user based on roles
func (mgr *AuthenticationManager) Authorize(userDetails UserDetails, roles []Role) bool {
	for _, role := range roles {
		if userDetails.SatisfiesRole(role) {
			return true
		}
	}
	return false
}
