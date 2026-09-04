// Package realm provides the Realm and AuthorizingRealm security contracts used
// by the bridge.
//
// The template-method pattern is modelled explicitly through pluggable hooks:
// AuthorizingRealm holds an authenticator hook (AuthenticationInfoSource) and an
// authorizer hook (AuthorizationInfoSource). A concrete realm supplies these
// hooks; a decorating realm can wrap another realm's hook, for example to count
// authentications before delegating.
package realm

import (
	"github.com/buji/pac4j/shiro/authc"
	"github.com/buji/pac4j/shiro/authz"
	"github.com/buji/pac4j/shiro/subject"
)

// AuthenticationInfoSource ports the protected
// AuthorizingRealm.doGetAuthenticationInfo(AuthenticationToken) template method.
type AuthenticationInfoSource interface {
	// DoGetAuthenticationInfo resolves authentication info for a token, or returns
	// an error mirroring a thrown AuthenticationException.
	DoGetAuthenticationInfo(token authc.AuthenticationToken) (authc.AuthenticationInfo, error)
}

// AuthorizationInfoSource ports the protected
// AuthorizingRealm.doGetAuthorizationInfo(PrincipalCollection) template method.
type AuthorizationInfoSource interface {
	// DoGetAuthorizationInfo resolves authorization info for the principals.
	DoGetAuthorizationInfo(principals subject.PrincipalCollection) authz.AuthorizationInfo
}

// Realm ports org.apache.shiro.realm.Realm: a security data source consulted by
// the security manager during authentication.
type Realm interface {
	// GetName returns the realm name (Realm.getName()).
	GetName() string
	// Supports reports whether this realm can process the token
	// (Realm.supports(AuthenticationToken)).
	Supports(token authc.AuthenticationToken) bool
	// GetAuthenticationInfo authenticates the token, returning info or an error
	// (Realm.getAuthenticationInfo(AuthenticationToken)).
	GetAuthenticationInfo(token authc.AuthenticationToken) (authc.AuthenticationInfo, error)
}

// AuthorizingRealm ports org.apache.shiro.realm.AuthorizingRealm: a realm that
// supplies both authentication and authorization info. The token-class filter and
// the two template-method hooks are configurable to reproduce Java subclassing.
type AuthorizingRealm struct {
	// name is the realm name; the security manager labels principals with it.
	name string
	// tokenSupports decides Realm.supports; it mirrors
	// setAuthenticationTokenClass by matching a token's dynamic type.
	tokenSupports func(token authc.AuthenticationToken) bool
	// authnSource provides doGetAuthenticationInfo.
	authnSource AuthenticationInfoSource
	// authzSource provides doGetAuthorizationInfo.
	authzSource AuthorizationInfoSource
}

// NewAuthorizingRealm constructs a base realm with the given hooks. The name
// defaults to a stable value and can be overridden with SetName.
func NewAuthorizingRealm(authn AuthenticationInfoSource, authzSrc AuthorizationInfoSource) *AuthorizingRealm {
	return &AuthorizingRealm{
		name:        "pac4jRealm",
		tokenSupports: func(token authc.AuthenticationToken) bool { return token != nil },
		authnSource: authn,
		authzSource: authzSrc,
	}
}

// SetName sets the realm name (Realm name is used as the realm key when building
// the principal collection).
func (r *AuthorizingRealm) SetName(name string) { r.name = name }

// GetName implements Realm.
func (r *AuthorizingRealm) GetName() string { return r.name }

// SetTokenSupports configures Realm.supports, reproducing
// setAuthenticationTokenClass(Class): only tokens accepted by fn are supported.
func (r *AuthorizingRealm) SetTokenSupports(fn func(token authc.AuthenticationToken) bool) {
	r.tokenSupports = fn
}

// Supports implements Realm.
func (r *AuthorizingRealm) Supports(token authc.AuthenticationToken) bool {
	if r.tokenSupports == nil {
		return token != nil
	}
	return r.tokenSupports(token)
}

// GetAuthenticationInfo implements Realm by delegating to the configured
// doGetAuthenticationInfo hook (AuthorizingRealm.getAuthenticationInfo ultimately
// calls the protected template method).
func (r *AuthorizingRealm) GetAuthenticationInfo(token authc.AuthenticationToken) (authc.AuthenticationInfo, error) {
	return r.authnSource.DoGetAuthenticationInfo(token)
}

// GetAuthorizationInfo delegates to the configured doGetAuthorizationInfo hook.
func (r *AuthorizingRealm) GetAuthorizationInfo(principals subject.PrincipalCollection) authz.AuthorizationInfo {
	return r.authzSource.DoGetAuthorizationInfo(principals)
}
