package context

import (
	"testing"

	bujirealm "github.com/buji/pac4j/buji/realm"
	bujisubject "github.com/buji/pac4j/buji/subject"
	pac4jcontext "github.com/buji/pac4j/pac4j/context"
	shiromgt "github.com/buji/pac4j/shiro/mgt"
	shiroutil "github.com/buji/pac4j/shiro/util"
)

func boundManager() *shiromgt.DefaultSecurityManager {
	sm := shiromgt.NewDefaultSecurityManager(bujirealm.NewPac4jRealm().GetAuthorizingRealm())
	sm.SetSubjectFactory(bujisubject.NewPac4jSubjectFactory())
	return sm
}

func TestShiroSessionStoreWithBoundManager(t *testing.T) {
	sm := boundManager()
	shiroutil.Bind(sm)
	defer shiroutil.Remove()

	ctx := pac4jcontext.NewNoopWebContext()

	INSTANCE.Set(ctx, "k", "v")
	if v, ok := INSTANCE.Get(ctx, "k"); !ok || v != "v" {
		t.Fatalf("expected Get(k)==(v,true), got (%v,%v)", v, ok)
	}
	if v, ok := INSTANCE.Get(ctx, "missing"); ok || v != nil {
		t.Fatalf("expected Get(missing)==(nil,false), got (%v,%v)", v, ok)
	}

	id, ok := INSTANCE.GetSessionID(ctx, true)
	if !ok || id == "" {
		t.Fatalf("expected a session id, got (%q,%v)", id, ok)
	}

	if _, ok := INSTANCE.GetTrackableSession(ctx); ok {
		t.Fatalf("GetTrackableSession must be empty")
	}
	if _, ok := INSTANCE.BuildFromTrackableSession(ctx, nil); ok {
		t.Fatalf("BuildFromTrackableSession must be empty")
	}
	if INSTANCE.RenewSession(ctx) {
		t.Fatalf("RenewSession must be false")
	}
	if !INSTANCE.DestroySession(ctx) {
		t.Fatalf("DestroySession must return true")
	}
}

func TestShiroSessionStoreNoManager(t *testing.T) {
	// Ensure nothing is bound: getSession must recover the panic and return nil.
	shiroutil.Remove()
	ctx := pac4jcontext.NewNoopWebContext()
	if id, ok := INSTANCE.GetSessionID(ctx, false); ok || id != "" {
		t.Fatalf("expected ('',false) with no bound manager, got (%q,%v)", id, ok)
	}
}
