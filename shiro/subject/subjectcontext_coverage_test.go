package subject

import (
	"testing"

	"github.com/buji/pac4j/shiro/session"
)

// stubToken is a minimal AuthenticationToken for exercising the subject context.
type stubToken struct{}

func (stubToken) GetPrincipal() any   { return "p" }
func (stubToken) GetCredentials() any { return "c" }

// TestDefaultSubjectContext exercises every SubjectContext accessor/mutator.
func TestDefaultSubjectContext(t *testing.T) {
	c := NewDefaultSubjectContext()

	if !c.IsSessionCreationEnabled() {
		t.Fatalf("session creation should default to enabled")
	}
	if c.IsAuthenticated() {
		t.Fatalf("context should not start authenticated")
	}

	c.SetAuthenticated(true)
	if !c.IsAuthenticated() {
		t.Fatalf("SetAuthenticated(true) not reflected")
	}

	var tok AuthenticationToken = stubToken{}
	c.SetAuthenticationToken(tok)
	if c.GetAuthenticationToken() != tok {
		t.Fatalf("authentication token not stored")
	}

	c.SetSessionCreationEnabled(false)
	if c.IsSessionCreationEnabled() {
		t.Fatalf("SetSessionCreationEnabled(false) not reflected")
	}

	pc := NewSimplePrincipalCollectionFor([]any{"alice"}, "realmA")
	c.SetPrincipals(pc)
	if c.GetPrincipals() != pc {
		t.Fatalf("principals not stored")
	}

	sess := session.NewSimpleSession()
	c.SetSession(sess)
	if c.GetSession() != sess {
		t.Fatalf("session not stored")
	}
}
