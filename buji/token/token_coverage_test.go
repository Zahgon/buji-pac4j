package token

import (
	"testing"

	"github.com/buji/pac4j/pac4j/profile"
)

func TestPac4jTokenEmptyProfiles(t *testing.T) {
	tok := NewPac4jToken(nil, false)
	if tok.GetPrincipal() != nil {
		t.Fatalf("expected nil principal for empty profiles, got %v", tok.GetPrincipal())
	}
	if tok.GetCredentials() != profile.ProfilesHashCode(nil) {
		t.Fatalf("credentials mismatch for empty profiles")
	}
	if tok.IsRememberMe() {
		t.Fatalf("expected isRemembered=false")
	}
}

func TestPac4jTokenWithProfiles(t *testing.T) {
	p := profile.NewCommonProfile()
	p.SetID("id-1")
	p.SetClientName("clientName")
	profiles := []profile.UserProfile{p}

	tok := NewPac4jToken(profiles, true)

	if got := tok.GetProfiles(); len(got) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(got))
	}
	prin := tok.GetPrincipal()
	if prin == nil {
		t.Fatalf("expected non-nil principal for non-empty profiles")
	}
	up, ok := prin.(profile.UserProfile)
	if !ok || up.GetID() != "id-1" {
		t.Fatalf("expected flattened profile id-1, got %v", prin)
	}
	if tok.GetCredentials() != profile.ProfilesHashCode(profiles) {
		t.Fatalf("credentials must equal ProfilesHashCode(profiles)")
	}
	if !tok.IsRememberMe() {
		t.Fatalf("expected isRemembered=true")
	}
}
