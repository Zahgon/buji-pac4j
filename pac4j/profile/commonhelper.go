package profile

import "strings"

// IsBlank ports org.pac4j.core.util.CommonHelper.isBlank(String): true when the
// string is empty or contains only whitespace. Used by Pac4jPrincipal to
// normalise the principalNameAttribute.
func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

// AreEquals ports org.pac4j.core.util.CommonHelper.areEquals(Object, Object):
// null-safe equality. Here specialised for the string values (client name, id,
// computed principal name) that the bridge compares.
func AreEquals(a, b *string) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return false
	}
	return *a == *b
}

// AreEqualsString is a convenience for comparing two plain strings with
// CommonHelper.areEquals semantics (non-nullable form).
func AreEqualsString(a, b string) bool { return a == b }

// AssertNotNil ports org.pac4j.core.util.CommonHelper.assertNotNull(name, obj):
// panics with a message when obj is nil, mirroring the Java TechnicalException
// thrown for a missing required argument.
func AssertNotNil(name string, obj any) {
	if obj == nil {
		panic("assertNotNull: " + name + " cannot be null")
	}
}
