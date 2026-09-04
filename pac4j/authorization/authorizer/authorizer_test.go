package authorizer

import (
	"testing"

	"github.com/buji/pac4j/pac4j/profile"
)

func rememberedProfile(remembered bool) profile.UserProfile {
	p := profile.NewCommonProfile()
	p.SetID("id")
	p.SetRemembered(remembered)
	return p
}

// TestIsRememberedAuthorizer covers empty, all-remembered and mixed cases.
func TestIsRememberedAuthorizer(t *testing.T) {
	a := NewIsRememberedAuthorizer()
	if a.IsAuthorized(nil) {
		t.Fatalf("empty list must not be authorized")
	}
	if !a.IsAuthorized([]profile.UserProfile{rememberedProfile(true)}) {
		t.Fatalf("all-remembered must be authorized")
	}
	if a.IsAuthorized([]profile.UserProfile{rememberedProfile(false)}) {
		t.Fatalf("non-remembered must not be authorized")
	}
	if a.IsAuthorized([]profile.UserProfile{nil}) {
		t.Fatalf("nil profile must not be authorized")
	}
}

// TestIsFullyAuthenticatedAuthorizer covers empty, fully-authenticated and
// all-remembered cases.
func TestIsFullyAuthenticatedAuthorizer(t *testing.T) {
	a := NewIsFullyAuthenticatedAuthorizer()
	if a.IsAuthorized(nil) {
		t.Fatalf("empty list must not be authorized")
	}
	if !a.IsAuthorized([]profile.UserProfile{rememberedProfile(false)}) {
		t.Fatalf("non-remembered profile must be fully authenticated")
	}
	if a.IsAuthorized([]profile.UserProfile{rememberedProfile(true)}) {
		t.Fatalf("all-remembered must not be fully authenticated")
	}
}
