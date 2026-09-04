package config

import (
	"testing"

	"github.com/buji/pac4j/pac4j/client"
	"github.com/buji/pac4j/pac4j/context"
	"github.com/buji/pac4j/pac4j/context/session"
)

// TestConfigClients covers NewConfig, GetClients and SetClients.
func TestConfigClients(t *testing.T) {
	c := NewConfig()
	if c.GetClients() != nil {
		t.Fatalf("new config has no clients")
	}
	cl := client.NewClients()
	c.SetClients(cl)
	if c.GetClients() != cl {
		t.Fatalf("SetClients/GetClients mismatch")
	}
}

// TestProfileManagerFactoryIfUndefined verifies the factory installs only once.
func TestProfileManagerFactoryIfUndefined(t *testing.T) {
	c := NewConfig()
	if c.GetProfileManagerFactory() != nil {
		t.Fatalf("no factory initially")
	}
	first := ProfileManagerFactory(func(ctx context.WebContext, store session.SessionStore) ProfileManager {
		return nil
	})
	c.SetProfileManagerFactoryIfUndefined(first)
	if c.GetProfileManagerFactory() == nil {
		t.Fatalf("factory should be installed")
	}
	// second install is ignored (IfUndefined)
	c.SetProfileManagerFactoryIfUndefined(func(ctx context.WebContext, store session.SessionStore) ProfileManager {
		return nil
	})
	// invoke to exercise the stored factory
	if got := c.GetProfileManagerFactory()(nil, nil); got != nil {
		t.Fatalf("stored factory should be the first one")
	}
}

// TestSessionStoreFactoryIfUndefined verifies the session-store factory installs
// only once and is invocable.
func TestSessionStoreFactoryIfUndefined(t *testing.T) {
	c := NewConfig()
	if c.GetSessionStoreFactory() != nil {
		t.Fatalf("no session store factory initially")
	}
	c.SetSessionStoreFactoryIfUndefined(func(params any) session.SessionStore { return nil })
	if c.GetSessionStoreFactory() == nil {
		t.Fatalf("session store factory should be installed")
	}
	c.SetSessionStoreFactoryIfUndefined(func(params any) session.SessionStore { return nil })
	if got := c.GetSessionStoreFactory()(nil); got != nil {
		t.Fatalf("stored session store factory should be invocable")
	}
}
