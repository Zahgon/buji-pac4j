// Package subject tests port io.buji.pac4j.subject.Pac4jPrincipalTests: they
// verify principal-name resolution across attribute configurations and the
// serialization round-trip.
package subject

import (
	"testing"

	shiroio "github.com/buji/pac4j/shiro/io"

	"github.com/buji/pac4j/pac4j/profile"
)

const (
	profileID    = "123"
	testUsername = "superman"
	testEmail    = "clark.kent@dailyplanet.org"
)

// createProfiles ports Pac4jPrincipalTests.createProfiles(): a single
// CommonProfile with an id and a fixed set of attributes.
func createProfiles() []profile.UserProfile {
	p := profile.NewCommonProfile()
	p.SetID(profileID)
	p.AddAttribute(profile.USERNAME, testUsername)
	p.AddAttribute("family_name", "Kent")
	p.AddAttribute("first_name", "Clark")
	p.AddAttribute("email", testEmail)
	p.AddAttribute("age", 21)
	return []profile.UserProfile{p}
}

// requireName dereferences a non-nil principal name or fails the test.
func requireName(t *testing.T, principal *Pac4jPrincipal) string {
	t.Helper()
	name := principal.GetName()
	if name == nil {
		t.Fatalf("expected a non-nil principal name, got nil")
	}
	return *name
}

// TestSerialize ports testSerialize: an (empty-profile) principal survives a
// serializer round-trip equal to the original.
func TestSerialize(t *testing.T) {
	profiles := []profile.UserProfile{}
	principal := NewPac4jPrincipal(profiles)

	serializer := shiroio.NewDefaultSerializer()
	data, err := serializer.Serialize(principal)
	if err != nil {
		t.Fatalf("serialize failed: %v", err)
	}

	principal2 := &Pac4jPrincipal{}
	if err := serializer.Deserialize(data, principal2); err != nil {
		t.Fatalf("deserialize failed: %v", err)
	}

	if !principal2.Equals(principal) {
		t.Fatalf("deserialized principal not equal to original")
	}
}

// TestNoAttribute ports testNoAttribute: with no name attribute the name is the
// profile id.
func TestNoAttribute(t *testing.T) {
	principal := NewPac4jPrincipal(createProfiles())
	if got := requireName(t, principal); got != profileID {
		t.Fatalf("getName() = %q, want %q", got, profileID)
	}
}

// TestBlankAttribute ports testBlankAttribute: a blank attribute falls back to
// the profile id.
func TestBlankAttribute(t *testing.T) {
	principal := NewPac4jPrincipalWithAttr(createProfiles(), " ")
	if got := requireName(t, principal); got != profileID {
		t.Fatalf("getName() = %q, want %q", got, profileID)
	}
}

// TestNullAttribute ports testNullAttribute: a null (empty) attribute falls back
// to the profile id.
func TestNullAttribute(t *testing.T) {
	principal := NewPac4jPrincipalWithAttr(createProfiles(), "")
	if got := requireName(t, principal); got != profileID {
		t.Fatalf("getName() = %q, want %q", got, profileID)
	}
}

// TestLeftPaddedAttribute ports testLeftPaddedAttribute: a left-padded attribute
// name is trimmed before lookup.
func TestLeftPaddedAttribute(t *testing.T) {
	principal := NewPac4jPrincipalWithAttr(createProfiles(), "  username")
	if got := requireName(t, principal); got != testUsername {
		t.Fatalf("getName() = %q, want %q", got, testUsername)
	}
}

// TestRightPaddedAttribute ports testRightPaddedAttribute: a right-padded
// attribute name is trimmed before lookup.
func TestRightPaddedAttribute(t *testing.T) {
	principal := NewPac4jPrincipalWithAttr(createProfiles(), "username ")
	if got := requireName(t, principal); got != testUsername {
		t.Fatalf("getName() = %q, want %q", got, testUsername)
	}
}

// TestUsernameAttribute ports testUsernameAttribute: the username attribute value
// is used as the name.
func TestUsernameAttribute(t *testing.T) {
	principal := NewPac4jPrincipalWithAttr(createProfiles(), "username")
	if got := requireName(t, principal); got != testUsername {
		t.Fatalf("getName() = %q, want %q", got, testUsername)
	}
}

// TestEmailAttribute ports testEmailAttribute: the email attribute value is used
// as the name.
func TestEmailAttribute(t *testing.T) {
	principal := NewPac4jPrincipalWithAttr(createProfiles(), "email")
	if got := requireName(t, principal); got != testEmail {
		t.Fatalf("getName() = %q, want %q", got, testEmail)
	}
}

// TestNonExistantAttribute ports testNonExistantAttribute: an absent attribute
// yields a nil name.
func TestNonExistantAttribute(t *testing.T) {
	principal := NewPac4jPrincipalWithAttr(createProfiles(), "display_name")
	if name := principal.GetName(); name != nil {
		t.Fatalf("getName() = %q, want nil", *name)
	}
}

// TestIntegerAttribute ports testIntegerAttribute: a non-string attribute value
// is rendered via String.valueOf semantics.
func TestIntegerAttribute(t *testing.T) {
	principal := NewPac4jPrincipalWithAttr(createProfiles(), "age")
	if got := requireName(t, principal); got != "21" {
		t.Fatalf("getName() = %q, want %q", got, "21")
	}
}
