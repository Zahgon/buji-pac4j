// Package token ports io.buji.pac4j.token: the authentication token carrying the
// pac4j profiles into the Shiro authentication pipeline.
package token

import "github.com/buji/pac4j/pac4j/profile"

// Pac4jToken ports io.buji.pac4j.token.Pac4jToken, an
// org.apache.shiro.authc.RememberMeAuthenticationToken whose payload is the list
// of pac4j user profiles resolved for the current request.
type Pac4jToken struct {
	profiles     []profile.UserProfile
	isRemembered bool
}

// NewPac4jToken constructs a token from the given profiles and remember-me flag.
func NewPac4jToken(profiles []profile.UserProfile, isRemembered bool) *Pac4jToken {
	return &Pac4jToken{profiles: profiles, isRemembered: isRemembered}
}

// GetProfiles returns the carried profiles (Pac4jToken.getProfiles()).
func (t *Pac4jToken) GetProfiles() []profile.UserProfile { return t.profiles }

// GetPrincipal ports Pac4jToken.getPrincipal(): it returns
// ProfileHelper.flatIntoOneProfile(profiles) — the flattened single profile when
// present, or nil when the list is empty (Java returns the Optional; the observable
// value is the contained profile or "absent").
func (t *Pac4jToken) GetPrincipal() any {
	if prof, ok := profile.FlatIntoOneProfile(t.profiles); ok {
		return prof
	}
	return nil
}

// GetCredentials ports Pac4jToken.getCredentials(): the profiles' hash code, so
// Shiro's credential matching keys on the concrete profile set.
func (t *Pac4jToken) GetCredentials() any { return profile.ProfilesHashCode(t.profiles) }

// IsRememberMe ports Pac4jToken.isRememberMe().
func (t *Pac4jToken) IsRememberMe() bool { return t.isRemembered }
