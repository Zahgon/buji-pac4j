package realm

import (
	"testing"

	bujisubject "github.com/buji/pac4j/buji/subject"
	"github.com/buji/pac4j/pac4j/profile"
	shirosubject "github.com/buji/pac4j/shiro/subject"
)

func contains(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}

func TestPac4jRealmNameAndAttribute(t *testing.T) {
	pr := NewPac4jRealm()
	if pr.GetAuthorizingRealm() == nil {
		t.Fatalf("GetAuthorizingRealm returned nil")
	}
	if pr.GetPrincipalNameAttribute() != "" {
		t.Fatalf("default principalNameAttribute must be empty")
	}
	pr.SetPrincipalNameAttribute("email")
	if pr.GetPrincipalNameAttribute() != "email" {
		t.Fatalf("SetPrincipalNameAttribute failed")
	}
	pr.SetName("customRealm")
	if pr.GetName() != "customRealm" {
		t.Fatalf("SetName/GetName failed, got %q", pr.GetName())
	}
}

func TestPac4jRealmDoGetAuthorizationInfoStringSlice(t *testing.T) {
	pr := NewPac4jRealm()

	p := profile.NewCommonProfile()
	p.SetID("id")
	p.AddRole("admin")
	p.AddAttribute(SHIRO_PERMISSIONS, []string{"read", "write"})
	principal := bujisubject.NewPac4jPrincipal([]profile.UserProfile{p})

	pc := shirosubject.NewSimplePrincipalCollectionFor([]any{"id", principal}, pr.GetName())
	info := pr.DoGetAuthorizationInfo(pc)

	if !contains(info.GetRoles(), "admin") {
		t.Fatalf("expected role admin, got %v", info.GetRoles())
	}
	if !contains(info.GetStringPermissions(), "read") || !contains(info.GetStringPermissions(), "write") {
		t.Fatalf("expected permissions read+write, got %v", info.GetStringPermissions())
	}
}

func TestPac4jRealmDoGetAuthorizationInfoAnySlice(t *testing.T) {
	pr := NewPac4jRealm()

	p := profile.NewCommonProfile()
	p.SetID("id")
	p.AddAttribute(SHIRO_PERMISSIONS, []any{"delete", 42, "purge"})
	principal := bujisubject.NewPac4jPrincipal([]profile.UserProfile{p})

	pc := shirosubject.NewSimplePrincipalCollectionFor([]any{"id", principal}, pr.GetName())
	info := pr.DoGetAuthorizationInfo(pc)

	perms := info.GetStringPermissions()
	if !contains(perms, "delete") || !contains(perms, "purge") {
		t.Fatalf("expected string permissions delete+purge, got %v", perms)
	}
	if contains(perms, "42") {
		t.Fatalf("non-string permission entries must be skipped, got %v", perms)
	}
}

func TestPac4jRealmDoGetAuthorizationInfoNoPrincipal(t *testing.T) {
	pr := NewPac4jRealm()
	pc := shirosubject.NewSimplePrincipalCollectionFor([]any{"just-a-string"}, pr.GetName())
	info := pr.DoGetAuthorizationInfo(pc)
	if len(info.GetRoles()) != 0 || len(info.GetStringPermissions()) != 0 {
		t.Fatalf("expected empty roles/permissions when no Pac4jPrincipal present")
	}
}
