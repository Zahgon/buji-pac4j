package config

import "testing"

// TestNewCasConfigurationWithLoginURL covers the constructor variant that seeds
// the login URL.
func TestNewCasConfigurationWithLoginURL(t *testing.T) {
	c := NewCasConfigurationWithLoginURL("https://cas/login")
	if c.GetLoginURL() != "https://cas/login" {
		t.Fatalf("seeded login URL not returned, got %q", c.GetLoginURL())
	}
}
