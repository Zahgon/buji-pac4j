// Package subject ports io.buji.pac4j.subject: the Pac4jPrincipal that carries
// pac4j user profiles into the Shiro principal collection, and the subject
// factory that adapts pac4j remember-me state to Shiro.
package subject

import (
	"bytes"
	"encoding/gob"

	"github.com/buji/pac4j/pac4j/profile"
)

// Pac4jPrincipal ports io.buji.pac4j.subject.Pac4jPrincipal: a Shiro principal
// wrapping the authenticated pac4j profiles. It implements value equality over
// its profiles and computes its name either from the primary profile id or from
// a configured principal-name attribute.
type Pac4jPrincipal struct {
	// profiles holds the pac4j user profiles; mutable to allow in-place refresh.
	profiles []profile.UserProfile
	// principalNameAttribute, when non-empty, names the profile attribute whose
	// value is used as the principal name; empty means "use the profile id".
	principalNameAttribute string
}

// NewPac4jPrincipal mirrors Pac4jPrincipal(List<UserProfile>): a principal with
// no configured name attribute.
func NewPac4jPrincipal(profiles []profile.UserProfile) *Pac4jPrincipal {
	return &Pac4jPrincipal{profiles: profiles}
}

// NewPac4jPrincipalWithAttr mirrors
// Pac4jPrincipal(List<UserProfile>, String principalNameAttribute): a blank
// attribute is normalised to empty (CommonHelper.isBlank), otherwise it is
// trimmed.
func NewPac4jPrincipalWithAttr(profiles []profile.UserProfile, principalNameAttribute string) *Pac4jPrincipal {
	attr := ""
	if !profile.IsBlank(principalNameAttribute) {
		attr = trimSpace(principalNameAttribute)
	}
	return &Pac4jPrincipal{profiles: profiles, principalNameAttribute: attr}
}

// trimSpace trims leading and trailing whitespace, matching String.trim().
func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && s[start] <= ' ' {
		start++
	}
	for end > start && s[end-1] <= ' ' {
		end--
	}
	return s[start:end]
}

// GetProfiles returns the wrapped profiles (Pac4jPrincipal.getProfiles()).
func (p *Pac4jPrincipal) GetProfiles() []profile.UserProfile {
	return p.profiles
}

// SetProfiles replaces the wrapped profiles (Pac4jPrincipal.setProfiles()).
// The bridge uses this to refresh a logged-in principal in place.
func (p *Pac4jPrincipal) SetProfiles(profiles []profile.UserProfile) {
	p.profiles = profiles
}

// GetProfile returns the single, flattened profile
// (Pac4jPrincipal.getProfile()). It panics when there is no profile, mirroring
// Optional.get() on an empty Optional; callers that may hold no profiles must
// avoid calling it (as the Java code does).
func (p *Pac4jPrincipal) GetProfile() profile.UserProfile {
	prof, ok := profile.FlatIntoOneProfile(p.profiles)
	if !ok {
		panic("No value present")
	}
	return prof
}

// GetName ports Pac4jPrincipal.getName(). It returns nil when the configured
// principal-name attribute is absent from the profile (the Java method can
// return null), otherwise the id or the attribute's string value.
func (p *Pac4jPrincipal) GetName() *string {
	prof := p.GetProfile()
	if p.principalNameAttribute == "" {
		id := prof.GetID()
		return &id
	}
	attrValue := prof.GetAttribute(p.principalNameAttribute)
	if attrValue == nil {
		return nil
	}
	s := profile.ValueOf(attrValue)
	return &s
}

// String ports Pac4jPrincipal.toString() which returns getName(); a nil name
// renders as the empty string.
func (p *Pac4jPrincipal) String() string {
	name := p.GetName()
	if name == nil {
		return ""
	}
	return *name
}

// Equals ports Pac4jPrincipal.equals(Object): identity plus null-safe profile
// list comparison.
func (p *Pac4jPrincipal) Equals(other *Pac4jPrincipal) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if p.profiles != nil {
		return profile.ProfilesEqual(p.profiles, other.profiles)
	}
	return other.profiles == nil
}

// HashCode ports Pac4jPrincipal.hashCode(): the profile list hash, or 0 when
// the list is nil.
func (p *Pac4jPrincipal) HashCode() int {
	if p.profiles != nil {
		return profile.ProfilesHashCode(p.profiles)
	}
	return 0
}

// gobPrincipal is the serialization image of a Pac4jPrincipal: the concrete
// CommonProfile values plus the name attribute. It exists so the principal can
// round-trip through the Shiro DefaultSerializer while preserving equality.
type gobPrincipal struct {
	Profiles               []profile.GobProfile
	PrincipalNameAttribute string
}

// MarshalBinary encodes the principal for the Shiro serializer, capturing the
// profiles and the name attribute.
func (p *Pac4jPrincipal) MarshalBinary() ([]byte, error) {
	img := gobPrincipal{PrincipalNameAttribute: p.principalNameAttribute}
	img.Profiles = make([]profile.GobProfile, 0, len(p.profiles))
	for _, pr := range p.profiles {
		img.Profiles = append(img.Profiles, profile.ToGobProfile(pr))
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary decodes a principal previously produced by MarshalBinary,
// reconstructing the profiles and name attribute so the result equals the
// original.
func (p *Pac4jPrincipal) UnmarshalBinary(data []byte) error {
	var img gobPrincipal
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&img); err != nil {
		return err
	}
	p.principalNameAttribute = img.PrincipalNameAttribute
	p.profiles = make([]profile.UserProfile, 0, len(img.Profiles))
	for _, gp := range img.Profiles {
		p.profiles = append(p.profiles, profile.FromGobProfile(gp))
	}
	return nil
}
