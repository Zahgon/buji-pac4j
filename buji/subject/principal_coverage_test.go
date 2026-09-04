package subject

import (
	"testing"

	"github.com/buji/pac4j/pac4j/profile"
)

func TestPac4jPrincipalStringAndHashCode(t *testing.T) {
	p := profile.NewCommonProfile()
	p.SetID("x")
	principal := NewPac4jPrincipal([]profile.UserProfile{p})

	if principal.String() != "x" {
		t.Fatalf("expected String()==x, got %q", principal.String())
	}

	// HashCode of two equal principals must match.
	p2 := profile.NewCommonProfile()
	p2.SetID("x")
	principal2 := NewPac4jPrincipal([]profile.UserProfile{p2})
	if principal.HashCode() != principal2.HashCode() {
		t.Fatalf("HashCode of equal principals must match")
	}

	// nil profiles => HashCode 0, String "" (name nil path via missing attr).
	empty := &Pac4jPrincipal{}
	if empty.HashCode() != 0 {
		t.Fatalf("expected HashCode 0 for nil profiles, got %d", empty.HashCode())
	}
}

func TestPac4jPrincipalStringNilName(t *testing.T) {
	p := profile.NewCommonProfile()
	p.SetID("id")
	// attribute "display_name" absent => GetName() nil => String() "".
	principal := NewPac4jPrincipalWithAttr([]profile.UserProfile{p}, "display_name")
	if principal.String() != "" {
		t.Fatalf("expected empty String() when name attribute absent, got %q", principal.String())
	}
}
