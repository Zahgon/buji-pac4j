package adapter

import (
	"testing"

	"github.com/buji/pac4j/pac4j/config"
)

// TestJEEFrameworkAdapter covers ApplyDefaultSettingsIfUndefined (no-op base)
// and String.
func TestJEEFrameworkAdapter(t *testing.T) {
	a := &JEEFrameworkAdapter{}
	// no-op base must not panic and must leave config untouched
	cfg := config.NewConfig()
	a.ApplyDefaultSettingsIfUndefined(cfg)
	if cfg.GetProfileManagerFactory() != nil || cfg.GetSessionStoreFactory() != nil {
		t.Fatalf("base adapter must not install factories")
	}
	if a.String() != "JEE" {
		t.Fatalf("JEEFrameworkAdapter.String() must be JEE, got %q", a.String())
	}
	// confirm it satisfies the FrameworkAdapter contract
	var _ FrameworkAdapter = a
}
