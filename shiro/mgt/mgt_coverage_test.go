package mgt_test

import (
	"testing"

	bujirealm "github.com/buji/pac4j/buji/realm"
	bujisubject "github.com/buji/pac4j/buji/subject"
	bujitoken "github.com/buji/pac4j/buji/token"
	"github.com/buji/pac4j/pac4j/profile"
	"github.com/buji/pac4j/shiro/mgt"
	"github.com/buji/pac4j/shiro/session"
	"github.com/buji/pac4j/shiro/subject"
)

// buildManager wires a real Pac4jRealm-backed AuthorizingRealm into a
// DefaultSecurityManager, exercising GetRealm and the default subject factory.
func buildManager() *mgt.DefaultSecurityManager {
	pr := bujirealm.NewPac4jRealm()
	return mgt.NewDefaultSecurityManager(pr.GetAuthorizingRealm())
}

func TestDefaultSecurityManagerCreateSubjectAndGetRealm(t *testing.T) {
	sm := buildManager()
	if sm.GetRealm() == nil {
		t.Fatal("GetRealm returned nil")
	}
	ctx := subject.NewDefaultSubjectContext()
	ctx.SetAuthenticated(true)
	s := sm.CreateSubject(ctx)
	if s == nil {
		t.Fatal("CreateSubject returned nil")
	}
	if !s.IsAuthenticated() {
		t.Fatal("expected subject authenticated from context flag")
	}
}

func TestDefaultSecurityManagerLoginRenewsSessionAndLogout(t *testing.T) {
	sm := buildManager()
	ctx := subject.NewDefaultSubjectContext()
	s := sm.CreateSubject(ctx)

	// Establish a session before login so renewal can be observed.
	sess := s.GetSession(true)
	if sess == nil {
		t.Fatal("expected a session")
	}
	idBefore := sess.GetID()

	p := profile.NewCommonProfile()
	p.SetID("id")
	p.SetClientName("clientName")
	token := bujitoken.NewPac4jToken([]profile.UserProfile{p}, false)

	if err := s.Login(token); err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if !s.IsAuthenticated() {
		t.Fatal("expected authenticated after login")
	}
	if s.GetSession(true).GetID() == idBefore {
		t.Fatal("expected session id renewed on login")
	}
	// principals present -> GetAuthorizationInfo path
	roles, perms := sm.GetAuthorizationInfo(s.GetPrincipals())
	if roles == nil || perms == nil {
		t.Fatal("expected non-nil role/perm sets")
	}

	sm.Logout(s)
	if s.IsAuthenticated() {
		t.Fatal("expected logout to clear authenticated state")
	}
	if s.GetPrincipals() != nil {
		t.Fatal("expected principals cleared after logout")
	}
}

func TestGetAuthorizationInfoNilPrincipals(t *testing.T) {
	sm := buildManager()
	roles, perms := sm.GetAuthorizationInfo(nil)
	if len(roles) != 0 || len(perms) != 0 {
		t.Fatal("expected empty sets for nil principals")
	}
}

func TestGetAuthorizationInfoWithRolesAndPerms(t *testing.T) {
	sm := buildManager()
	p := profile.NewCommonProfile()
	p.SetID("id")
	p.SetClientName("clientName")
	p.AddRole("admin")
	p.AddAttribute(bujirealm.SHIRO_PERMISSIONS, []string{"read", "write"})
	principal := bujisubject.NewPac4jPrincipal([]profile.UserProfile{p})
	pc := subject.NewSimplePrincipalCollectionFor([]any{"id", principal}, "pac4jRealm")

	roles, perms := sm.GetAuthorizationInfo(pc)
	if _, ok := roles["admin"]; !ok {
		t.Fatal("expected admin role")
	}
	if _, ok := perms["read"]; !ok {
		t.Fatal("expected read permission")
	}
	if _, ok := perms["write"]; !ok {
		t.Fatal("expected write permission")
	}
}

// stubToken is an unsupported token type used to reach Login's error branches.
type stubToken struct{}

func (stubToken) GetPrincipal() any   { return nil }
func (stubToken) GetCredentials() any { return nil }

func TestLoginUnsupportedToken(t *testing.T) {
	sm := buildManager()
	ctx := subject.NewDefaultSubjectContext()
	s := sm.CreateSubject(ctx)
	if _, err := sm.Login(s, stubToken{}); err == nil {
		t.Fatal("expected error for unsupported token")
	}
}

func TestDelegatingSubjectHasRoleIsPermittedNoPrincipals(t *testing.T) {
	sm := buildManager()
	s := mgt.NewDelegatingSubject(sm)
	if s.HasRole("admin") {
		t.Fatal("expected false with no principals")
	}
	if s.IsPermitted("read") {
		t.Fatal("expected false with no principals")
	}
	if s.IsRemembered() {
		t.Fatal("expected not remembered initially")
	}
}

func TestDelegatingSubjectRememberedAndSetters(t *testing.T) {
	sm := buildManager()
	s := mgt.NewDelegatingSubject(sm)
	pc := subject.NewSimplePrincipalCollectionFor([]any{"alice"}, "pac4jRealm")
	s.SetPrincipalsCollection(pc)
	s.SetRemembered(true)
	s.SetAuthenticated(false)
	if !s.IsRemembered() {
		t.Fatal("expected remembered when principals present, not authenticated, remembered set")
	}
	sess := session.NewSimpleSession()
	s.SetSession(sess)
	if s.GetSession(false) != sess {
		t.Fatal("expected the session we set")
	}
}
