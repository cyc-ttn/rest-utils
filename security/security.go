package security

import (
	"context"
	"errors"
	"net/http"
)

type key string

// PrincipalKey for use in storing in context
const PrincipalKey key = "Principal"

// AuthenticationProvider interface
// provides different methods of authentication
type AuthenticationProvider interface {
	Authenticate(*http.Request, UserDetailsService) (UserDetails, map[string]interface{}, error)
}

// AuthenticatorMiddleware authenticates a request
// and passes to the next middleware if the authentication succeeds
// otherwise, allows the user to define an action
func AuthenticatorMiddleware(
	mgr *AuthenticationManager,
	failure func(http.ResponseWriter, *http.Request, []error),
) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		userDetails, info, errs := mgr.Authenticate(r)
		if errs != nil {
			failure(rw, r, errs)
		} else {
			newCtx := context.WithValue(r.Context(), PrincipalKey, map[string]interface{}{"user": userDetails, "info": info})
			next(rw, r.WithContext(newCtx))
		}
	}
}

// AuthorizationMiddleware authorizes a request
func AuthorizationMiddleware(
	mgr *AuthenticationManager,
	roles []Role,
	failure func(http.ResponseWriter, *http.Request, []error),
) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		defErr := []error{ errors.New("You are unauthorized to view this resource") }

		// Retrieve authentication from context
		value := r.Context().Value(PrincipalKey).(map[string]interface{})

		if value == nil {
			failure(rw, r, defErr)
			return
		}

		usr := value["user"].(UserDetails)
		if mgr.Authorize(usr, roles) {
			next(rw, r)
		} else {
			failure(rw, r, defErr)
		}
	}
}

// GetUserFromRequest - Gets the user from the request!
func GetUserFromRequest(req * http.Request) UserDetails {
	principalAsIf := req.Context().Value(PrincipalKey)
	if principalAsIf == nil { return nil }

	principal := principalAsIf.(map[string]interface{})
	return principal["user"].(UserDetails)
}
