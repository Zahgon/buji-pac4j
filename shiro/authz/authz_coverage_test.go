package authz

import "testing"

func contains(list []string, want string) bool {
	for _, v := range list {
		if v == want {
			return true
		}
	}
	return false
}

// TestSimpleAuthorizationInfo exercises role and permission accumulation and lookup.
func TestSimpleAuthorizationInfo(t *testing.T) {
	sai := NewSimpleAuthorizationInfo()
	if len(sai.GetRoles()) != 0 || len(sai.GetStringPermissions()) != 0 {
		t.Fatalf("new authorization info should be empty")
	}

	sai.AddRoles([]string{"admin", "user"})
	sai.AddRole("auditor")
	sai.AddStringPermissions([]string{"read", "write"})
	sai.AddStringPermission("delete")

	roles := sai.GetRoles()
	if len(roles) != 3 || !contains(roles, "admin") || !contains(roles, "user") || !contains(roles, "auditor") {
		t.Fatalf("unexpected roles: %v", roles)
	}
	perms := sai.GetStringPermissions()
	if len(perms) != 3 || !contains(perms, "read") || !contains(perms, "write") || !contains(perms, "delete") {
		t.Fatalf("unexpected permissions: %v", perms)
	}

	if !sai.HasRole("admin") || sai.HasRole("missing") {
		t.Fatalf("HasRole incorrect")
	}
	if !sai.HasPermission("delete") || sai.HasPermission("missing") {
		t.Fatalf("HasPermission incorrect")
	}
}

// TestSimpleAuthorizationInfoDedup verifies set semantics (no duplicates).
func TestSimpleAuthorizationInfoDedup(t *testing.T) {
	sai := NewSimpleAuthorizationInfo()
	sai.AddRole("admin")
	sai.AddRole("admin")
	sai.AddStringPermission("read")
	sai.AddStringPermission("read")
	if len(sai.GetRoles()) != 1 {
		t.Fatalf("expected 1 role after dedup, got %v", sai.GetRoles())
	}
	if len(sai.GetStringPermissions()) != 1 {
		t.Fatalf("expected 1 permission after dedup, got %v", sai.GetStringPermissions())
	}
}
