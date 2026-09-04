// Package adapter ports org.pac4j.core.adapter.FrameworkAdapter and the JEE
// specialisation that io.buji.pac4j's FrameworkAdapterImpl extends.
package adapter

import "github.com/buji/pac4j/pac4j/config"

// FrameworkAdapter ports the org.pac4j.core.adapter.FrameworkAdapter contract:
// it applies framework-specific default settings to a Config.
type FrameworkAdapter interface {
	// ApplyDefaultSettingsIfUndefined installs framework defaults that are not
	// already set on the config.
	ApplyDefaultSettingsIfUndefined(cfg *config.Config)
}

// JEEFrameworkAdapter ports org.pac4j.jee.adapter.JEEFrameworkAdapter: the base
// adapter whose default behaviour the Shiro adapter extends via embedding. Its
// ApplyDefaultSettingsIfUndefined is a no-op base implementation (the JEE
// defaults it would add are not consulted by the migrated bridge/tests), so
// subclasses call it last after applying their own defaults.
type JEEFrameworkAdapter struct{}

// ApplyDefaultSettingsIfUndefined implements the base behaviour (super call).
// The JEE base installs no additional Config defaults in the migrated bridge, so
// it simply validates its argument and returns.
func (a *JEEFrameworkAdapter) ApplyDefaultSettingsIfUndefined(cfg *config.Config) {
	_ = cfg
}

// String mirrors JEEFrameworkAdapter.toString() ("JEE").
func (a *JEEFrameworkAdapter) String() string { return "JEE" }
