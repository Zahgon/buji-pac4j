// Package adapter ports org.pac4j.framework.adapter.FrameworkAdapterImpl: the
// framework adapter that installs the Shiro-backed default profile manager and
// session store into a pac4j Config.
package adapter

import (
	bujicontext "github.com/buji/pac4j/buji/context"
	bujiprofile "github.com/buji/pac4j/buji/profile"
	"github.com/buji/pac4j/pac4j/config"
	pac4jcontext "github.com/buji/pac4j/pac4j/context"
	pac4jsession "github.com/buji/pac4j/pac4j/context/session"
	pac4jadapter "github.com/buji/pac4j/pac4j/framework/adapter"
	"github.com/buji/pac4j/pac4j/profile"
)

const frameworkName = "Shiro"

// FrameworkAdapterImpl ports org.pac4j.framework.adapter.FrameworkAdapterImpl.
type FrameworkAdapterImpl struct {
	*pac4jadapter.JEEFrameworkAdapter
}

// NewFrameworkAdapterImpl constructs the adapter.
func NewFrameworkAdapterImpl() *FrameworkAdapterImpl {
	return &FrameworkAdapterImpl{JEEFrameworkAdapter: &pac4jadapter.JEEFrameworkAdapter{}}
}

// ApplyDefaultSettingsIfUndefined ports
// FrameworkAdapterImpl.applyDefaultSettingsIfUndefined: it installs the
// ShiroProfileManager factory and the ShiroSessionStore singleton as defaults
// before delegating to the base adapter.
func (a *FrameworkAdapterImpl) ApplyDefaultSettingsIfUndefined(cfg *config.Config) {
	profile.AssertNotNil("config", cfg)
	cfg.SetProfileManagerFactoryIfUndefined(func(ctx pac4jcontext.WebContext, store pac4jsession.SessionStore) config.ProfileManager {
		return bujiprofile.NewShiroProfileManager(ctx, store)
	})
	cfg.SetSessionStoreFactoryIfUndefined(func(params any) pac4jsession.SessionStore {
		return bujicontext.INSTANCE
	})
}

// String ports FrameworkAdapterImpl.toString: "Shiro".
func (a *FrameworkAdapterImpl) String() string { return frameworkName }
