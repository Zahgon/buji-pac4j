// Package subject ports the subset of org.apache.shiro.subject used by the
// bridge: the PrincipalCollection contract, its simple implementation, and the
// Subject/SubjectContext contracts.
package subject

import "reflect"

// PrincipalCollection ports org.apache.shiro.subject.PrincipalCollection: an
// immutable-ish view over the principals of one or more realms.
type PrincipalCollection interface {
	// GetPrimaryPrincipal returns the "primary" principal — the first principal
	// added (PrincipalCollection.getPrimaryPrincipal()).
	GetPrimaryPrincipal() any
	// OneByType returns the single principal assignable to the given type, or
	// nil when absent (PrincipalCollection.oneByType(Class)). The target is a
	// pointer to a typed nil used purely to convey the desired type, e.g.
	// (*Pac4jPrincipal)(nil).
	OneByType(target any) any
	// FromRealm returns the principals contributed by the named realm
	// (PrincipalCollection.fromRealm(String)).
	FromRealm(realmName string) []any
	// AsList returns all principals across realms in insertion order.
	AsList() []any
	// IsEmpty reports whether the collection holds no principals.
	IsEmpty() bool
}

// SimplePrincipalCollection ports
// org.apache.shiro.subject.SimplePrincipalCollection.
type SimplePrincipalCollection struct {
	// realmPrincipals maps realm name to that realm's principals, preserving
	// insertion order within each realm.
	realmPrincipals map[string][]any
	// realmOrder records realm insertion order for primary-principal selection.
	realmOrder []string
}

// NewSimplePrincipalCollection constructs an empty collection.
func NewSimplePrincipalCollection() *SimplePrincipalCollection {
	return &SimplePrincipalCollection{realmPrincipals: make(map[string][]any)}
}

// NewSimplePrincipalCollectionFor mirrors
// new SimplePrincipalCollection(Collection principals, String realmName): adds
// each principal under the given realm, preserving order.
func NewSimplePrincipalCollectionFor(principals []any, realmName string) *SimplePrincipalCollection {
	c := NewSimplePrincipalCollection()
	for _, p := range principals {
		c.Add(p, realmName)
	}
	return c
}

// Add appends a principal under a realm (SimplePrincipalCollection.add).
func (c *SimplePrincipalCollection) Add(principal any, realmName string) {
	if principal == nil {
		return
	}
	if _, ok := c.realmPrincipals[realmName]; !ok {
		c.realmOrder = append(c.realmOrder, realmName)
	}
	c.realmPrincipals[realmName] = append(c.realmPrincipals[realmName], principal)
}

// GetPrimaryPrincipal implements PrincipalCollection: the first principal of the
// first realm.
func (c *SimplePrincipalCollection) GetPrimaryPrincipal() any {
	for _, realm := range c.realmOrder {
		ps := c.realmPrincipals[realm]
		if len(ps) > 0 {
			return ps[0]
		}
	}
	return nil
}

// OneByType implements PrincipalCollection: returns the first principal whose
// dynamic type is assignable to the type of target.
func (c *SimplePrincipalCollection) OneByType(target any) any {
	want := reflect.TypeOf(target)
	for _, realm := range c.realmOrder {
		for _, p := range c.realmPrincipals[realm] {
			if p == nil {
				continue
			}
			if reflect.TypeOf(p) == want {
				return p
			}
			// Also allow assignability (interface targets).
			if want != nil && want.Kind() == reflect.Interface && reflect.TypeOf(p).AssignableTo(want) {
				return p
			}
		}
	}
	return nil
}

// FromRealm implements PrincipalCollection.
func (c *SimplePrincipalCollection) FromRealm(realmName string) []any {
	return c.realmPrincipals[realmName]
}

// AsList implements PrincipalCollection.
func (c *SimplePrincipalCollection) AsList() []any {
	out := make([]any, 0)
	for _, realm := range c.realmOrder {
		out = append(out, c.realmPrincipals[realm]...)
	}
	return out
}

// IsEmpty implements PrincipalCollection.
func (c *SimplePrincipalCollection) IsEmpty() bool {
	return len(c.AsList()) == 0
}
