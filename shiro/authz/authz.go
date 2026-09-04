// Package authz ports the subset of org.apache.shiro.authz used by the bridge:
// the AuthorizationInfo contract and its simple implementation carrying roles and
// string permissions.
package authz

// AuthorizationInfo ports org.apache.shiro.authz.AuthorizationInfo: the roles and
// permissions a realm associates with a set of principals.
type AuthorizationInfo interface {
	// GetRoles returns the assigned role names (AuthorizationInfo.getRoles()).
	GetRoles() []string
	// GetStringPermissions returns the assigned permission strings
	// (AuthorizationInfo.getStringPermissions()).
	GetStringPermissions() []string
}

// SimpleAuthorizationInfo ports org.apache.shiro.authz.SimpleAuthorizationInfo.
// Roles and permissions are stored as sets (Java uses Set), so duplicate adds are
// coalesced; iteration order is not part of the observable contract.
type SimpleAuthorizationInfo struct {
	roles       map[string]struct{}
	permissions map[string]struct{}
}

// NewSimpleAuthorizationInfo constructs an empty authorization info, mirroring the
// no-arg SimpleAuthorizationInfo() constructor.
func NewSimpleAuthorizationInfo() *SimpleAuthorizationInfo {
	return &SimpleAuthorizationInfo{
		roles:       make(map[string]struct{}),
		permissions: make(map[string]struct{}),
	}
}

// AddRoles mirrors SimpleAuthorizationInfo.addRoles(Collection).
func (s *SimpleAuthorizationInfo) AddRoles(roles []string) {
	for _, r := range roles {
		s.roles[r] = struct{}{}
	}
}

// AddRole mirrors SimpleAuthorizationInfo.addRole(String).
func (s *SimpleAuthorizationInfo) AddRole(role string) {
	s.roles[role] = struct{}{}
}

// AddStringPermissions mirrors
// SimpleAuthorizationInfo.addStringPermissions(Collection).
func (s *SimpleAuthorizationInfo) AddStringPermissions(permissions []string) {
	for _, p := range permissions {
		s.permissions[p] = struct{}{}
	}
}

// AddStringPermission mirrors
// SimpleAuthorizationInfo.addStringPermission(String).
func (s *SimpleAuthorizationInfo) AddStringPermission(permission string) {
	s.permissions[permission] = struct{}{}
}

// GetRoles implements AuthorizationInfo.
func (s *SimpleAuthorizationInfo) GetRoles() []string {
	return keys(s.roles)
}

// GetStringPermissions implements AuthorizationInfo.
func (s *SimpleAuthorizationInfo) GetStringPermissions() []string {
	return keys(s.permissions)
}

// Has reports whether the given role is present (used by Subject.hasRole).
func (s *SimpleAuthorizationInfo) HasRole(role string) bool {
	_, ok := s.roles[role]
	return ok
}

// HasPermission reports whether the given string permission is present (used by
// Subject.isPermitted).
func (s *SimpleAuthorizationInfo) HasPermission(permission string) bool {
	_, ok := s.permissions[permission]
	return ok
}

// keys returns the keys of a string set.
func keys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
