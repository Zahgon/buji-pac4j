// Package config ports the subset of org.pac4j.cas.config.CasConfiguration
// exercised by IniTests: a configuration object exposing a settable loginUrl.
//
// The Java test uses commons-beanutils' FluentPropertyBeanIntrospector to prove
// that Lombok generated a fluent, writable setter for the loginUrl property. In
// Go the equivalent guarantee is that the CasConfiguration type exposes a
// settable LoginUrl property; the ported test verifies SetLoginUrl exists and
// takes effect.
package config

// CasConfiguration ports org.pac4j.cas.config.CasConfiguration (the subset used
// by the migrated tests): it holds the CAS server login URL.
type CasConfiguration struct {
	loginURL string
}

// NewCasConfiguration constructs an empty CasConfiguration, mirroring the Java
// no-arg constructor.
func NewCasConfiguration() *CasConfiguration { return &CasConfiguration{} }

// NewCasConfigurationWithLoginURL mirrors new CasConfiguration(String loginUrl).
func NewCasConfigurationWithLoginURL(loginURL string) *CasConfiguration {
	return &CasConfiguration{loginURL: loginURL}
}

// GetLoginURL returns the configured login URL (CasConfiguration.getLoginUrl).
func (c *CasConfiguration) GetLoginURL() string { return c.loginURL }

// SetLoginURL sets the login URL, mirroring the fluent writable property the
// Java test asserts exists (CasConfiguration.setLoginUrl).
func (c *CasConfiguration) SetLoginURL(loginURL string) *CasConfiguration {
	c.loginURL = loginURL
	return c
}
