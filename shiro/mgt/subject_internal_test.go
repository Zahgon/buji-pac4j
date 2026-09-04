package mgt

import (
	"testing"

	"github.com/buji/pac4j/shiro/subject"
)

// TestDelegatingSubjectUnexportedSetters exercises the unexported white-box
// setters (setPrincipals, setAuthenticated, setRemembered) directly, since these
// are internal helpers not reachable from external-package tests.
func TestDelegatingSubjectUnexportedSetters(t *testing.T) {
	s := NewDelegatingSubject(nil)

	pc := subject.NewSimplePrincipalCollectionFor([]any{"alice"}, "realmA")
	s.setPrincipals(pc)
	if s.GetPrincipals() == nil {
		t.Fatal("setPrincipals did not set principals")
	}

	s.setAuthenticated(true)
	if !s.IsAuthenticated() {
		t.Fatal("setAuthenticated did not set authenticated flag")
	}

	s.setRemembered(true)
	s.setAuthenticated(false)
	if !s.IsRemembered() {
		t.Fatal("setRemembered did not set remembered flag")
	}
}

// TestDefaultSubjectFactoryWithPrincipals ensures the factory's principals path
// (setPrincipals branch) is exercised.
func TestDefaultSubjectFactoryWithPrincipals(t *testing.T) {
	ctx := subject.NewDefaultSubjectContext()
	ctx.SetAuthenticated(true)
	ctx.SetPrincipals(subject.NewSimplePrincipalCollectionFor([]any{"bob"}, "realmB"))

	f := defaultSubjectFactory{}
	s := f.CreateSubject(ctx, (*DefaultSecurityManager)(nil))
	if s == nil {
		t.Fatal("CreateSubject returned nil")
	}
	if s.GetPrincipals() == nil {
		t.Fatal("expected principals seeded from context")
	}
}
