// Package profile ports io.buji.pac4j.profile.ShiroProfileManager: the pac4j
// ProfileManager that populates the Shiro subject when profiles are saved and
// logs it out when they are removed.
package profile

import (
	bujiutil "github.com/buji/pac4j/buji/util"
	pac4jcontext "github.com/buji/pac4j/pac4j/context"
	pac4jsession "github.com/buji/pac4j/pac4j/context/session"
	"github.com/buji/pac4j/pac4j/config"
	"github.com/buji/pac4j/shiro/authc"
	shiroutil "github.com/buji/pac4j/shiro/util"
)

// ShiroProfileManager ports io.buji.pac4j.profile.ShiroProfileManager.
type ShiroProfileManager struct {
	context pac4jcontext.WebContext
	store   pac4jsession.SessionStore
}

// NewShiroProfileManager ports the ShiroProfileManager(WebContext, SessionStore)
// constructor.
func NewShiroProfileManager(ctx pac4jcontext.WebContext, store pac4jsession.SessionStore) *ShiroProfileManager {
	return &ShiroProfileManager{context: ctx, store: store}
}

// SaveAll ports ShiroProfileManager.saveAll: after the base save it populates the
// Shiro subject; on an authentication failure it removes the profiles and
// propagates the error.
func (m *ShiroProfileManager) SaveAll(profiles []config.OrderedProfile, saveInSession bool) error {
	err := bujiutil.PopulateSubject(profiles)
	if err != nil {
		if _, ok := err.(*authc.AuthenticationException); ok {
			m.RemoveProfiles()
		}
		return err
	}
	return nil
}

// RemoveProfiles ports ShiroProfileManager.removeProfiles: after the base removal
// it logs the Shiro subject out.
func (m *ShiroProfileManager) RemoveProfiles() {
	shiroutil.GetSubject().Logout()
}
