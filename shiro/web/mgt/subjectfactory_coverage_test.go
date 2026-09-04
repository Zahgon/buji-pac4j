package mgt

import (
	"testing"

	"github.com/buji/pac4j/shiro/authc"
	"github.com/buji/pac4j/shiro/authz"
	shiromgt "github.com/buji/pac4j/shiro/mgt"
	"github.com/buji/pac4j/shiro/realm"
	"github.com/buji/pac4j/shiro/session"
	"github.com/buji/pac4j/shiro/subject"
)

// stubAuthnSource / stubAuthzSource provide a minimal realm so a
// DefaultSecurityManager can be constructed without importing the buji layer.
type stubAuthnSource struct{}

func (stubAuthnSource) DoGetAuthenticationInfo(token authc.AuthenticationToken) (authc.AuthenticationInfo, error) {
	pc := subject.NewSimplePrincipalCollectionFor([]any{"alice"}, "stubRealm")
	return authc.NewSimpleAuthenticationInfo(pc, 0), nil
}

type stubAuthzSource struct{}

func (stubAuthzSource) DoGetAuthorizationInfo(principals subject.PrincipalCollection) authz.AuthorizationInfo {
	return authz.NewSimpleAuthorizationInfo()
}

func buildSM() *shiromgt.DefaultSecurityManager {
	ar := realm.NewAuthorizingRealm(stubAuthnSource{}, stubAuthzSource{})
	return shiromgt.NewDefaultSecurityManager(ar)
}

func TestDefaultWebSubjectFactoryCreateSubject(t *testing.T) {
	f := NewDefaultWebSubjectFactory()
	if f == nil {
		t.Fatal("NewDefaultWebSubjectFactory returned nil")
	}
	sm := buildSM()

	ctx := subject.NewDefaultSubjectContext()
	ctx.SetAuthenticated(true)
	pc := subject.NewSimplePrincipalCollectionFor([]any{"alice"}, "stubRealm")
	ctx.SetPrincipals(pc)
	sess := session.NewSimpleSession()
	ctx.SetSession(sess)

	s := f.CreateSubject(ctx, sm)
	if s == nil {
		t.Fatal("CreateSubject returned nil")
	}
	if !s.IsAuthenticated() {
		t.Fatal("expected authenticated flag carried from context")
	}
	if s.GetPrincipals() == nil {
		t.Fatal("expected principals carried from context")
	}
	if s.GetSession(false) != sess {
		t.Fatal("expected session carried from context")
	}
}

func TestDefaultWebSubjectFactoryCreateSubjectMinimalContext(t *testing.T) {
	f := NewDefaultWebSubjectFactory()
	sm := buildSM()

	ctx := subject.NewDefaultSubjectContext()
	s := f.CreateSubject(ctx, sm)
	if s == nil {
		t.Fatal("CreateSubject returned nil")
	}
	if s.IsAuthenticated() {
		t.Fatal("expected unauthenticated for empty context")
	}
}
