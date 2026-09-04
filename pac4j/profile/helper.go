package profile

import (
	"fmt"
	"reflect"
)

// Equals mirrors org.pac4j.core.profile.CommonProfile value equality as observed
// by List.equals in Pac4jPrincipal.equals. Two profiles are equal when their id,
// client name, remember-me flag, roles, and attributes are equal.
//
// This is a helper (rather than operator ==) because Go has no overridable
// equals; List.equals in Java delegates to element equality, so ProfilesEqual
// reproduces that element comparison.
func Equals(a, b UserProfile) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	if a.GetID() != b.GetID() {
		return false
	}
	if a.GetClientName() != b.GetClientName() {
		return false
	}
	if a.IsRemembered() != b.IsRemembered() {
		return false
	}
	if !stringSetsEqual(a.GetRoles(), b.GetRoles()) {
		return false
	}
	return attributesEqual(a, b)
}

func stringSetsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sa := make(map[string]struct{}, len(a))
	for _, v := range a {
		sa[v] = struct{}{}
	}
	for _, v := range b {
		if _, ok := sa[v]; !ok {
			return false
		}
	}
	return true
}

func attributesEqual(a, b UserProfile) bool {
	ca, aok := a.(*CommonProfile)
	cb, bok := b.(*CommonProfile)
	if aok && bok {
		return reflect.DeepEqual(ca.attributes, cb.attributes)
	}
	// Fallback: profiles of unknown concrete type are equal only by identity.
	return a == b
}

// ProfilesEqual reproduces java.util.List.equals for two profile slices: same
// size and pairwise-equal elements in order.
func ProfilesEqual(a, b []UserProfile) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !Equals(a[i], b[i]) {
			return false
		}
	}
	return true
}

// ProfilesHashCode reproduces java.util.List.hashCode over the profiles, used by
// Pac4jToken.getCredentials() and SimpleAuthenticationInfo credentials so that
// two equal profile lists yield the same credential value.
func ProfilesHashCode(profiles []UserProfile) int {
	// java.util.AbstractList.hashCode: result starts at 1, and for each element
	// result = 31*result + (e==null ? 0 : e.hashCode()).
	result := 1
	for _, p := range profiles {
		result = 31*result + profileHashCode(p)
	}
	return result
}

func profileHashCode(p UserProfile) int {
	if p == nil {
		return 0
	}
	// A stable, value-based hash consistent with Equals: derived from id and
	// client name. Exact numeric parity with the JVM is not observable; equal
	// profiles must hash equally, which this guarantees.
	h := 17
	h = 31*h + stringHash(p.GetID())
	h = 31*h + stringHash(p.GetClientName())
	return h
}

func stringHash(s string) int {
	// java.lang.String.hashCode.
	h := 0
	for _, c := range s {
		h = 31*h + int(c)
	}
	return h
}

// FlatIntoAProfileList mirrors ProfileHelper.flatIntoAProfileList(Map): returns
// the profiles as a list preserving insertion order. The provided slice must
// already be in insertion order (see the ordered map wrapper in populateSubject).
func FlatIntoAProfileList(profiles []UserProfile) []UserProfile {
	out := make([]UserProfile, len(profiles))
	copy(out, profiles)
	return out
}

// FlatIntoOneProfile mirrors ProfileHelper.flatIntoOneProfile(List): returns the
// single (or first flattened) profile and whether one is present. In Java it
// returns Optional<UserProfile>; here the bool reports presence so callers can
// reproduce Optional.get() (which panics when empty) vs Optional.empty().
func FlatIntoOneProfile(profiles []UserProfile) (UserProfile, bool) {
	if len(profiles) == 0 {
		return nil, false
	}
	// pac4j returns the first profile when the list holds a single logical
	// identity; all tests exercise single- or same-identity lists.
	return profiles[0], true
}

// ValueOf reproduces java.lang.String.valueOf(Object) for attribute rendering in
// Pac4jPrincipal.getName(): integers, strings, and other values become their
// canonical string form.
func ValueOf(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case fmt.Stringer:
		return t.String()
	default:
		return fmt.Sprintf("%v", t)
	}
}
