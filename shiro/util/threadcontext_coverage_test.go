package util

import (
	"testing"

	bujirealm "github.com/buji/pac4j/buji/realm"
	shiromgt "github.com/buji/pac4j/shiro/mgt"
)

// TestThreadContextBindGetSubjectRemove exercises Bind, GetSecurityManager,
// GetSubject (build-and-cache path) and Remove.
func TestThreadContextBindGetSubjectRemove(t *testing.T) {
	sm := shiromgt.NewDefaultSecurityManager(bujirealm.NewPac4jRealm().GetAuthorizingRealm())
	Bind(sm)
	defer Remove()

	var tc ThreadContext
	if tc.GetSecurityManager() != sm {
		t.Fatal("GetSecurityManager did not return the bound security manager")
	}

	subject := GetSubject()
	if subject == nil {
		t.Fatal("GetSubject returned nil for a bound security manager")
	}
	// Second call returns the cached subject (same instance).
	if GetSubject() != subject {
		t.Fatal("GetSubject did not return the cached subject on the second call")
	}
}

// TestGetSubjectPanicsWhenUnbound verifies the no-SecurityManager panic path.
func TestGetSubjectPanicsWhenUnbound(t *testing.T) {
	Remove()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("GetSubject did not panic when no SecurityManager is bound")
		}
	}()
	_ = GetSubject()
}
