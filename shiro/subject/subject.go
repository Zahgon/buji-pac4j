// This file ports the org.apache.shiro.subject.Subject and SubjectContext
// contracts used by the bridge.
package subject

import (
	"github.com/buji/pac4j/shiro/session"
)

// PRINCIPALS_SESSION_KEY ports
// org.apache.shiro.subject.support.DefaultSubjectContext.PRINCIPALS_SESSION_KEY:
// the session attribute key under which the principal collection is stored.
const PRINCIPALS_SESSION_KEY = "org.apache.shiro.subject.support.DefaultSubjectContext_PRINCIPALS_SESSION_KEY"

// AUTHENTICATED_SESSION_KEY ports
// DefaultSubjectContext.AUTHENTICATED_SESSION_KEY.
const AUTHENTICATED_SESSION_KEY = "org.apache.shiro.subject.support.DefaultSubjectContext_AUTHENTICATED_SESSION_KEY"

// AuthenticationToken is duplicated here as a narrow interface to avoid an import
// cycle with package authc; any authc.AuthenticationToken satisfies it.
type AuthenticationToken interface {
	GetPrincipal() any
	GetCredentials() any
}

// Subject ports org.apache.shiro.subject.Subject: the security view of the
// currently acting entity.
type Subject interface {
	// GetPrincipal returns the primary principal (Subject.getPrincipal()).
	GetPrincipal() any
	// GetPrincipals returns all principals, or nil when none
	// (Subject.getPrincipals()).
	GetPrincipals() PrincipalCollection
	// IsAuthenticated reports full authentication (Subject.isAuthenticated()).
	IsAuthenticated() bool
	// IsRemembered reports remembered-but-not-authenticated state
	// (Subject.isRemembered()).
	IsRemembered() bool
	// Login authenticates the subject with the token, throwing on failure
	// (Subject.login(AuthenticationToken)).
	Login(token AuthenticationToken) error
	// Logout clears the authenticated state (Subject.logout()).
	Logout()
	// GetSession returns the session, creating it when create is true
	// (Subject.getSession(boolean)).
	GetSession(create bool) session.Session
	// HasRole reports whether the subject has the named role
	// (Subject.hasRole(String)).
	HasRole(role string) bool
	// IsPermitted reports whether the subject holds the permission
	// (Subject.isPermitted(String)).
	IsPermitted(permission string) bool
}

// SubjectContext ports org.apache.shiro.subject.SubjectContext: the mutable bag
// of state used to build a Subject.
type SubjectContext interface {
	// IsAuthenticated reports the authenticated flag being carried
	// (SubjectContext.isAuthenticated()).
	IsAuthenticated() bool
	// SetAuthenticated sets the authenticated flag
	// (SubjectContext.setAuthenticated(boolean)).
	SetAuthenticated(authenticated bool)
	// GetAuthenticationToken returns the token driving construction
	// (SubjectContext.getAuthenticationToken()).
	GetAuthenticationToken() AuthenticationToken
	// SetAuthenticationToken sets the token.
	SetAuthenticationToken(token AuthenticationToken)
	// SetSessionCreationEnabled toggles whether the built subject may create a
	// session (SubjectContext.setSessionCreationEnabled(boolean)).
	SetSessionCreationEnabled(enabled bool)
	// IsSessionCreationEnabled reports the flag above.
	IsSessionCreationEnabled() bool
	// GetPrincipals returns any principals already resolved.
	GetPrincipals() PrincipalCollection
	// SetPrincipals sets resolved principals.
	SetPrincipals(principals PrincipalCollection)
	// GetSession returns any existing session carried by the context.
	GetSession() session.Session
	// SetSession sets the session.
	SetSession(s session.Session)
}
