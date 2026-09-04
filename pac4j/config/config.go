// Package config ports the subset of org.pac4j.core.config.Config used by the
// default INI and the framework adapter: the holder of clients plus the
// pluggable profile-manager and session-store factories.
package config

import (
	"github.com/buji/pac4j/pac4j/client"
	"github.com/buji/pac4j/pac4j/context"
	"github.com/buji/pac4j/pac4j/context/session"
	"github.com/buji/pac4j/pac4j/profile"
)

// ProfileManager ports the profile-manager contract that a factory produces.
// Its concrete Shiro implementation lives in the buji profile package; here it is
// an opaque interface so the Config can hold a factory without importing Shiro.
type ProfileManager interface {
	// SaveAll persists the ordered profiles, optionally into the session
	// (ProfileManager.save / saveAll).
	SaveAll(profiles []OrderedProfile, saveInSession bool) error
	// RemoveProfiles clears the stored profiles (ProfileManager.removeProfiles).
	RemoveProfiles()
}

// OrderedProfile pairs a client-name key with its profile, preserving the
// insertion order of the Java LinkedHashMap<String, UserProfile>.
type OrderedProfile struct {
	// Key is the client-name key.
	Key string
	// Profile is the associated user profile.
	Profile profile.UserProfile
}

// ProfileManagerFactory mirrors the functional interface used by
// setProfileManagerFactoryIfUndefined: builds a ProfileManager from a context and
// session store.
type ProfileManagerFactory func(ctx context.WebContext, store session.SessionStore) ProfileManager

// SessionStoreFactory mirrors the functional interface used by
// setSessionStoreFactoryIfUndefined: builds a SessionStore from parameters.
type SessionStoreFactory func(params any) session.SessionStore

// Config ports org.pac4j.core.config.Config.
type Config struct {
	clients               *client.Clients
	profileManagerFactory ProfileManagerFactory
	sessionStoreFactory   SessionStoreFactory
}

// NewConfig constructs an empty Config.
func NewConfig() *Config { return &Config{} }

// GetClients returns the configured clients (Config.getClients).
func (c *Config) GetClients() *client.Clients { return c.clients }

// SetClients sets the configured clients (Config.setClients).
func (c *Config) SetClients(clients *client.Clients) { c.clients = clients }

// GetProfileManagerFactory returns the profile-manager factory.
func (c *Config) GetProfileManagerFactory() ProfileManagerFactory { return c.profileManagerFactory }

// SetProfileManagerFactoryIfUndefined mirrors
// Config.setProfileManagerFactoryIfUndefined: installs the factory only when none
// is already set.
func (c *Config) SetProfileManagerFactoryIfUndefined(f ProfileManagerFactory) {
	if c.profileManagerFactory == nil {
		c.profileManagerFactory = f
	}
}

// GetSessionStoreFactory returns the session-store factory.
func (c *Config) GetSessionStoreFactory() SessionStoreFactory { return c.sessionStoreFactory }

// SetSessionStoreFactoryIfUndefined mirrors
// Config.setSessionStoreFactoryIfUndefined: installs the factory only when none
// is already set.
func (c *Config) SetSessionStoreFactoryIfUndefined(f SessionStoreFactory) {
	if c.sessionStoreFactory == nil {
		c.sessionStoreFactory = f
	}
}
