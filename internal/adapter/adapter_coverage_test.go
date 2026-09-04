package adapter

import (
	"testing"

	"github.com/buji/pac4j/pac4j/config"
)

func TestFrameworkAdapterImplDefaults(t *testing.T) {
	a := NewFrameworkAdapterImpl()
	if a == nil {
		t.Fatalf("NewFrameworkAdapterImpl returned nil")
	}
	if a.String() != "Shiro" {
		t.Fatalf("expected String()==Shiro, got %q", a.String())
	}

	cfg := config.NewConfig()
	a.ApplyDefaultSettingsIfUndefined(cfg)

	if cfg.GetProfileManagerFactory() == nil {
		t.Fatalf("expected profile manager factory to be installed")
	}
	if cfg.GetSessionStoreFactory() == nil {
		t.Fatalf("expected session store factory to be installed")
	}

	// Exercise the embedded base adapter's no-op ApplyDefaultSettingsIfUndefined
	// (the JEE "super" call path) so it is covered as executed.
	base := config.NewConfig()
	a.JEEFrameworkAdapter.ApplyDefaultSettingsIfUndefined(base)
	if base.GetProfileManagerFactory() != nil || base.GetSessionStoreFactory() != nil {
		t.Fatalf("base adapter must not install factories")
	}
}
