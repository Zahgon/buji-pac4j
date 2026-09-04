// Package authc ports the subset of org.apache.shiro.authc used by the bridge:
// authentication tokens, the remember-me marker, and authentication info.
package authc

import "github.com/buji/pac4j/shiro/subject"

// AuthenticationToken ports org.apache.shiro.authc.AuthenticationToken: the
// credentials submitted during a login attempt.
type AuthenticationToken interface {
	// GetPrincipal returns the account identity the token asserts
	// (AuthenticationToken.getPrincipal()).
	GetPrincipal() any
	// GetCredentials returns the proof of that identity
	// (AuthenticationToken.getCredentials()).
	GetCredentials() any
}

// RememberMeAuthenticationToken ports
// org.apache.shiro.authc.RememberMeAuthenticationToken: an AuthenticationToken
// that also carries a remember-me preference.
type RememberMeAuthenticationToken interface {
	AuthenticationToken
	// IsRememberMe reports whether the authenticating user wishes to be
	// remembered across sessions (RememberMeAuthenticationToken.isRememberMe()).
	IsRememberMe() bool
}

// AuthenticationInfo ports org.apache.shiro.authc.AuthenticationInfo: the
// identity and credentials a realm returns for a successful authentication.
type AuthenticationInfo interface {
	// GetPrincipals returns the account's principals
	// (AuthenticationInfo.getPrincipals()).
	GetPrincipals() subject.PrincipalCollection
	// GetCredentials returns the account's stored credentials
	// (AuthenticationInfo.getCredentials()).
	GetCredentials() any
}

// SimpleAuthenticationInfo ports
// org.apache.shiro.authc.SimpleAuthenticationInfo.
type SimpleAuthenticationInfo struct {
	principals  subject.PrincipalCollection
	credentials any
}

// NewSimpleAuthenticationInfo mirrors
// new SimpleAuthenticationInfo(PrincipalCollection, Object credentials).
func NewSimpleAuthenticationInfo(principals subject.PrincipalCollection, credentials any) *SimpleAuthenticationInfo {
	return &SimpleAuthenticationInfo{principals: principals, credentials: credentials}
}

// GetPrincipals implements AuthenticationInfo.
func (s *SimpleAuthenticationInfo) GetPrincipals() subject.PrincipalCollection { return s.principals }

// GetCredentials implements AuthenticationInfo.
func (s *SimpleAuthenticationInfo) GetCredentials() any { return s.credentials }

// AuthenticationException ports org.apache.shiro.authc.AuthenticationException:
// the error raised when authentication fails.
type AuthenticationException struct {
	// Message is the failure detail.
	Message string
	// Cause is the wrapped error, if any.
	Cause error
}

// Error implements error.
func (e *AuthenticationException) Error() string {
	if e.Message == "" && e.Cause != nil {
		return e.Cause.Error()
	}
	return e.Message
}

// Unwrap exposes the cause for errors.Is/As.
func (e *AuthenticationException) Unwrap() error { return e.Cause }
