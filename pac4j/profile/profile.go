// Package profile ports the subset of org.pac4j.core.profile used by buji-pac4j.
//
// It provides the UserProfile contract, the CommonProfile implementation, the
// Pac4jConstants values, and the ProfileHelper flattening utilities that the
// bridge relies upon.
package profile

import (
	"sort"
)

// Pac4jConstants mirrors org.pac4j.core.util.Pac4jConstants values referenced by
// the migrated code and tests.
const (
	// USERNAME is org.pac4j.core.util.Pac4jConstants.USERNAME.
	USERNAME = "username"
)

// UserProfile mirrors org.pac4j.core.profile.UserProfile: the contract the
// bridge consumes when transferring pac4j identity into the Shiro subject.
type UserProfile interface {
	// GetID returns the profile identifier (UserProfile.getId()).
	GetID() string
	// SetID sets the profile identifier (UserProfile.setId(String)).
	SetID(id string)
	// GetClientName returns the authenticating client name
	// (UserProfile.getClientName()).
	GetClientName() string
	// SetClientName sets the authenticating client name
	// (UserProfile.setClientName(String)).
	SetClientName(clientName string)
	// GetAttribute returns the attribute value or nil if absent
	// (UserProfile.getAttribute(String)).
	GetAttribute(name string) any
	// AddAttribute stores an attribute (UserProfile.addAttribute(String,Object)).
	AddAttribute(name string, value any)
	// GetRoles returns the set of role names granted to the profile
	// (UserProfile.getRoles()); iteration order is deterministic (sorted) to
	// mirror the observable behavior of a Shiro role set.
	GetRoles() []string
	// AddRole grants a role to the profile (UserProfile.addRole(String)).
	AddRole(role string)
	// IsRemembered reports whether the profile was restored from a
	// remember-me token (UserProfile.isRemembered()).
	IsRemembered() bool
	// SetRemembered sets the remember-me flag (UserProfile.setRemembered(boolean)).
	SetRemembered(remembered bool)
}

// CommonProfile ports org.pac4j.core.profile.CommonProfile: the default profile
// implementation used throughout the tests.
type CommonProfile struct {
	id         string
	clientName string
	attributes map[string]any
	roles      map[string]struct{}
	remembered bool
}

// NewCommonProfile constructs an empty CommonProfile, matching the Java no-arg
// constructor semantics (empty attributes/roles, not-remembered).
func NewCommonProfile() *CommonProfile {
	return &CommonProfile{
		attributes: make(map[string]any),
		roles:      make(map[string]struct{}),
	}
}

// GetID implements UserProfile.
func (p *CommonProfile) GetID() string { return p.id }

// SetID implements UserProfile.
func (p *CommonProfile) SetID(id string) { p.id = id }

// GetClientName implements UserProfile.
func (p *CommonProfile) GetClientName() string { return p.clientName }

// SetClientName implements UserProfile.
func (p *CommonProfile) SetClientName(clientName string) { p.clientName = clientName }

// GetAttribute implements UserProfile. It returns nil when the attribute is
// absent, mirroring Java's null return.
func (p *CommonProfile) GetAttribute(name string) any {
	if p.attributes == nil {
		return nil
	}
	v, ok := p.attributes[name]
	if !ok {
		return nil
	}
	return v
}

// AddAttribute implements UserProfile.
func (p *CommonProfile) AddAttribute(name string, value any) {
	if p.attributes == nil {
		p.attributes = make(map[string]any)
	}
	p.attributes[name] = value
}

// GetRoles implements UserProfile, returning a sorted snapshot for a
// deterministic, set-like iteration order.
func (p *CommonProfile) GetRoles() []string {
	out := make([]string, 0, len(p.roles))
	for r := range p.roles {
		out = append(out, r)
	}
	sort.Strings(out)
	return out
}

// AddRole implements UserProfile.
func (p *CommonProfile) AddRole(role string) {
	if p.roles == nil {
		p.roles = make(map[string]struct{})
	}
	p.roles[role] = struct{}{}
}

// IsRemembered implements UserProfile.
func (p *CommonProfile) IsRemembered() bool { return p.remembered }

// SetRemembered implements UserProfile.
func (p *CommonProfile) SetRemembered(remembered bool) { p.remembered = remembered }
