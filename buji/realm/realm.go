// Package realm ports io.buji.pac4j.realm.Pac4jRealm: the Shiro realm that
// authenticates and authorizes subjects from pac4j user profiles carried by a
// Pac4jToken.
package realm

import (
	bujisubject "github.com/buji/pac4j/buji/subject"
	"github.com/buji/pac4j/buji/token"
	"github.com/buji/pac4j/pac4j/profile"
	"github.com/buji/pac4j/shiro/authc"
	"github.com/buji/pac4j/shiro/authz"
	shirorealm "github.com/buji/pac4j/shiro/realm"
	shirosubject "github.com/buji/pac4j/shiro/subject"
)

// SHIRO_PERMISSIONS ports Pac4jRealm.SHIRO_PERMISSIONS: the profile attribute
// name under which string permissions are carried.
const SHIRO_PERMISSIONS = "$shiro_permissions$"

// Pac4jRealm ports io.buji.pac4j.realm.Pac4jRealm. It composes a Shiro
// AuthorizingRealm, supplying itself as the source of authentication and
// authorization info (Go's stand-in for extending AuthorizingRealm).
type Pac4jRealm struct {
	authorizingRealm      *shirorealm.AuthorizingRealm
	principalNameAttribute string
}

// NewPac4jRealm constructs the realm, wiring itself as the authn/authz hook and
// restricting supported tokens to the Pac4jToken type.
func NewPac4jRealm() *Pac4jRealm {
	p := &Pac4jRealm{}
	p.authorizingRealm = shirorealm.NewAuthorizingRealm(p, p)
	p.authorizingRealm.SetTokenSupports(func(tok authc.AuthenticationToken) bool {
		_, ok := tok.(*token.Pac4jToken)
		return ok
	})
	return p
}

// DoGetAuthenticationInfo ports Pac4jRealm.doGetAuthenticationInfo: it builds a
// Pac4jPrincipal from the token's profiles, derives the username via the
// configured principal-name attribute, and returns a principal collection of
// [username, principal] under the realm name with the profiles' hash as
// credentials.
func (p *Pac4jRealm) DoGetAuthenticationInfo(tok authc.AuthenticationToken) (authc.AuthenticationInfo, error) {
	pt := tok.(*token.Pac4jToken)
	profiles := pt.GetProfiles()
	principal := bujisubject.NewPac4jPrincipalWithAttr(profiles, p.principalNameAttribute)
	name := principal.GetName()
	var usernameVal any = nil
	if name != nil {
		usernameVal = *name
	}
	pc := shirosubject.NewSimplePrincipalCollectionFor([]any{usernameVal, principal}, p.authorizingRealm.GetName())
	return authc.NewSimpleAuthenticationInfo(pc, profile.ProfilesHashCode(profiles)), nil
}

// DoGetAuthorizationInfo ports Pac4jRealm.doGetAuthorizationInfo: it aggregates
// roles and string permissions ($shiro_permissions$) across every profile of the
// Pac4jPrincipal found in the collection.
func (p *Pac4jRealm) DoGetAuthorizationInfo(principals shirosubject.PrincipalCollection) authz.AuthorizationInfo {
	sai := authz.NewSimpleAuthorizationInfo()
	pr := principals.OneByType((*bujisubject.Pac4jPrincipal)(nil))
	if pr != nil {
		p2 := pr.(*bujisubject.Pac4jPrincipal)
		for _, prof := range p2.GetProfiles() {
			if prof != nil {
				sai.AddRoles(prof.GetRoles())
				perm := prof.GetAttribute(SHIRO_PERMISSIONS)
				switch v := perm.(type) {
				case []string:
					sai.AddStringPermissions(v)
				case []any:
					for _, x := range v {
						if s, ok := x.(string); ok {
							sai.AddStringPermission(s)
						}
					}
				}
			}
		}
	}
	return sai
}

// GetAuthorizingRealm returns the composed Shiro realm, used to register the
// realm with a security manager.
func (p *Pac4jRealm) GetAuthorizingRealm() *shirorealm.AuthorizingRealm {
	return p.authorizingRealm
}

// SetPrincipalNameAttribute ports Pac4jRealm.setPrincipalNameAttribute.
func (p *Pac4jRealm) SetPrincipalNameAttribute(attr string) {
	p.principalNameAttribute = attr
}

// GetPrincipalNameAttribute ports Pac4jRealm.getPrincipalNameAttribute.
func (p *Pac4jRealm) GetPrincipalNameAttribute() string {
	return p.principalNameAttribute
}

// GetName returns the realm name.
func (p *Pac4jRealm) GetName() string { return p.authorizingRealm.GetName() }

// SetName sets the realm name.
func (p *Pac4jRealm) SetName(name string) { p.authorizingRealm.SetName(name) }
