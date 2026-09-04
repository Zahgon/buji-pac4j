package realm

import (
	"testing"

	"github.com/buji/pac4j/shiro/authc"
	"github.com/buji/pac4j/shiro/authz"
	"github.com/buji/pac4j/shiro/subject"
)

// stubAuthnSource / stubAuthzSource are minimal hook implementations letting us
// exercise AuthorizingRealm's delegation, naming, and token-support behaviour.
type stubAuthnSource struct{}

func (stubAuthnSource) DoGetAuthenticationInfo(token authc.AuthenticationToken) (authc.AuthenticationInfo, error) {
	pc := subject.NewSimplePrincipalCollectionFor([]any{"alice"}, "stubRealm")
	return authc.NewSimpleAuthenticationInfo(pc, 0), nil
}

type stubAuthzSource struct{}

func (stubAuthzSource) DoGetAuthorizationInfo(principals subject.PrincipalCollection) authz.AuthorizationInfo {
	sai := authz.NewSimpleAuthorizationInfo()
	sai.AddRole("stub-role")
	return sai
}

type okToken struct{}

func (okToken) GetPrincipal() any   { return "alice" }
func (okToken) GetCredentials() any { return nil }

func TestAuthorizingRealmNameAndDelegation(t *testing.T) {
	ar := NewAuthorizingRealm(stubAuthnSource{}, stubAuthzSource{})
	if ar.GetName() != "pac4jRealm" {
		t.Fatalf("expected default name pac4jRealm, got %q", ar.GetName())
	}
	ar.SetName("customRealm")
	if ar.GetName() != "customRealm" {
		t.Fatalf("SetName failed, got %q", ar.GetName())
	}

	// default tokenSupports = token != nil
	if !ar.Supports(okToken{}) {
		t.Fatal("expected Supports true for non-nil token")
	}

	info, err := ar.GetAuthenticationInfo(okToken{})
	if err != nil || info == nil {
		t.Fatalf("GetAuthenticationInfo failed: %v", err)
	}

	pc := subject.NewSimplePrincipalCollectionFor([]any{"alice"}, "stubRealm")
	azi := ar.GetAuthorizationInfo(pc)
	if azi == nil || len(azi.GetRoles()) != 1 || azi.GetRoles()[0] != "stub-role" {
		t.Fatal("expected delegated authorization info with stub-role")
	}
}

func TestAuthorizingRealmSetTokenSupports(t *testing.T) {
	ar := NewAuthorizingRealm(stubAuthnSource{}, stubAuthzSource{})
	ar.SetTokenSupports(func(token authc.AuthenticationToken) bool { return false })
	if ar.Supports(okToken{}) {
		t.Fatal("expected Supports false after overriding tokenSupports")
	}
}
