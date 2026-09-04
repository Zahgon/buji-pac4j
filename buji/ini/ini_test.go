// Package ini ports io.buji.pac4j.ini.IniTests: it verifies that a pac4j
// configuration object exposes a settable property, mirroring the original
// test's check that CasConfiguration.loginUrl is writable.
package ini

import (
	"testing"

	casconfig "github.com/buji/pac4j/pac4j/cas/config"
)

// TestIniWithLombok ports IniTests.testIniWithLombok. The Java test used
// commons-beanutils to assert that CasConfiguration exposes a writable
// 'loginUrl' property (a fluent setter generated for the field). The Go
// equivalent asserts the same intent: the configuration property is settable
// and the fluent setter returns the configuration for chaining.
func TestIniWithLombok(t *testing.T) {
	c := casconfig.NewCasConfiguration()
	if c == nil {
		t.Fatal("NewCasConfiguration returned nil")
	}

	// The property must be writable and read back the value that was set.
	c.SetLoginURL("http://login")
	if got := c.GetLoginURL(); got != "http://login" {
		t.Fatalf("GetLoginURL() = %q, want %q", got, "http://login")
	}

	// The setter must be fluent: it returns the same configuration so callers
	// can chain, and the returned value reflects the update.
	c2 := c.SetLoginURL("http://example.org/cas/login")
	if c2 == nil {
		t.Fatal("SetLoginURL returned nil; expected fluent *CasConfiguration")
	}
	if got := c2.GetLoginURL(); got != "http://example.org/cas/login" {
		t.Fatalf("fluent GetLoginURL() = %q, want %q", got, "http://example.org/cas/login")
	}
}
