// Package authorizer ports the pac4j authorizers used by ShiroHelper to decide
// between a full authentication and a remember-me authentication.
package authorizer

import "github.com/buji/pac4j/pac4j/profile"

// Authorizer ports org.pac4j.core.authorization.authorizer.Authorizer: a
// predicate over the current profiles.
type Authorizer interface {
	// IsAuthorized reports whether the profiles satisfy this authorizer.
	// The WebContext and SessionStore are unused by the authorizers the bridge
	// relies on (they inspect only the profiles), matching the Java calls
	// isAuthorized(null, null, profiles).
	IsAuthorized(profiles []profile.UserProfile) bool
}

// IsRememberedAuthorizer ports
// org.pac4j.core.authorization.authorizer.IsRememberedAuthorizer: authorized when
// there is at least one profile and every profile is a remember-me profile.
type IsRememberedAuthorizer struct{}

// NewIsRememberedAuthorizer constructs the authorizer.
func NewIsRememberedAuthorizer() *IsRememberedAuthorizer { return &IsRememberedAuthorizer{} }

// IsAuthorized implements Authorizer. Mirrors AbstractCheckAuthenticationAuthorizer
// semantics: requires a non-empty profile list, and every profile remembered.
func (a *IsRememberedAuthorizer) IsAuthorized(profiles []profile.UserProfile) bool {
	if len(profiles) == 0 {
		return false
	}
	for _, p := range profiles {
		if p == nil || !p.IsRemembered() {
			return false
		}
	}
	return true
}

// IsFullyAuthenticatedAuthorizer ports
// org.pac4j.core.authorization.authorizer.IsFullyAuthenticatedAuthorizer:
// authorized when there is at least one profile and not all profiles are
// remembered (i.e. at least one fully-authenticated profile).
type IsFullyAuthenticatedAuthorizer struct{}

// NewIsFullyAuthenticatedAuthorizer constructs the authorizer.
func NewIsFullyAuthenticatedAuthorizer() *IsFullyAuthenticatedAuthorizer {
	return &IsFullyAuthenticatedAuthorizer{}
}

// IsAuthorized implements Authorizer. Requires a non-empty profile list where at
// least one profile is not remembered.
func (a *IsFullyAuthenticatedAuthorizer) IsAuthorized(profiles []profile.UserProfile) bool {
	if len(profiles) == 0 {
		return false
	}
	for _, p := range profiles {
		if p != nil && !p.IsRemembered() {
			return true
		}
	}
	return false
}
