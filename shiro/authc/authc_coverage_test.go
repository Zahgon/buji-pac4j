package authc

import (
	"errors"
	"testing"

	"github.com/buji/pac4j/shiro/subject"
)

// TestSimpleAuthenticationInfo exercises the principals/credentials accessors.
func TestSimpleAuthenticationInfo(t *testing.T) {
	pc := subject.NewSimplePrincipalCollectionFor([]any{"alice"}, "realmA")
	info := NewSimpleAuthenticationInfo(pc, 1234)
	if info.GetPrincipals() != pc {
		t.Fatalf("GetPrincipals mismatch")
	}
	if info.GetCredentials() != 1234 {
		t.Fatalf("GetCredentials: got %v want 1234", info.GetCredentials())
	}
}

// TestAuthenticationException exercises both Error branches and Unwrap.
func TestAuthenticationException(t *testing.T) {
	e1 := &AuthenticationException{Message: "boom"}
	if e1.Error() != "boom" {
		t.Fatalf("Error() with message: got %q", e1.Error())
	}
	if e1.Unwrap() != nil {
		t.Fatalf("Unwrap with no cause should be nil")
	}

	cause := errors.New("root")
	e2 := &AuthenticationException{Cause: cause}
	if e2.Error() != "root" {
		t.Fatalf("Error() should fall back to cause: got %q", e2.Error())
	}
	if !errors.Is(e2, cause) {
		t.Fatalf("errors.Is should find the wrapped cause")
	}
}
