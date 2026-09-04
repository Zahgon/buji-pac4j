// Package util ports io.buji.pac4j.util.ShiroHelper: the core bridge that pushes
// pac4j profiles into the Shiro subject, logging in or refreshing profiles as
// appropriate.
package util

import (
	bujisubject "github.com/buji/pac4j/buji/subject"
	"github.com/buji/pac4j/buji/token"
	"github.com/buji/pac4j/pac4j/authorization/authorizer"
	"github.com/buji/pac4j/pac4j/config"
	"github.com/buji/pac4j/pac4j/profile"
	shirosubject "github.com/buji/pac4j/shiro/subject"
	shiroutil "github.com/buji/pac4j/shiro/util"
)

// fullyAuthenticatedAuthorizer ports ShiroHelper.IS_FULLY_AUTHENTICATED_AUTHORIZER.
var fullyAuthenticatedAuthorizer = authorizer.NewIsFullyAuthenticatedAuthorizer()

// rememberedAuthorizer ports ShiroHelper.IS_REMEMBERED_AUTHORIZER.
var rememberedAuthorizer = authorizer.NewIsRememberedAuthorizer()

// PopulateSubject ports ShiroHelper.populateSubject: given the ordered profiles,
// it logs the subject in fully or as remembered depending on the profiles'
// authentication state.
func PopulateSubject(profiles []config.OrderedProfile) error {
	if len(profiles) == 0 {
		return nil
	}
	list := make([]profile.UserProfile, 0, len(profiles))
	for _, op := range profiles {
		list = append(list, op.Profile)
	}
	sub := shiroutil.GetSubject()
	if fullyAuthenticatedAuthorizer.IsAuthorized(list) {
		return login(sub, list, false, sub.IsAuthenticated())
	}
	if rememberedAuthorizer.IsAuthorized(list) {
		return login(sub, list, true, sub.IsRemembered())
	}
	return nil
}

// login ports ShiroHelper.login: if the subject is already logged in and its
// profiles can be refreshed in place (same user), it avoids a new login;
// otherwise it performs a fresh login with a Pac4jToken.
func login(sub shirosubject.Subject, profiles []profile.UserProfile, rememberMe bool, alreadyLoggedIn bool) error {
	if alreadyLoggedIn && refreshProfiles(sub, profiles) {
		return nil
	}
	return sub.Login(token.NewPac4jToken(profiles, rememberMe))
}

// refreshProfiles ports ShiroHelper.refreshProfiles: it updates the existing
// Pac4jPrincipal's profiles in place when the user is unchanged and the computed
// principal name is unchanged, keeping the current session; returns false when a
// full re-login is required.
func refreshProfiles(sub shirosubject.Subject, profiles []profile.UserProfile) bool {
	pcs := sub.GetPrincipals()
	if pcs == nil {
		return false
	}
	pr := pcs.OneByType((*bujisubject.Pac4jPrincipal)(nil))
	if pr == nil {
		return false
	}
	principal := pr.(*bujisubject.Pac4jPrincipal)
	if !isSameUser(principal.GetProfiles(), profiles) {
		return false
	}
	old := principal.GetProfiles()
	oldName := principal.GetName()
	principal.SetProfiles(profiles)
	if !profile.AreEquals(oldName, principal.GetName()) {
		principal.SetProfiles(old)
		return false
	}
	sess := sub.GetSession(false)
	if sess != nil && sess.GetAttribute(shirosubject.PRINCIPALS_SESSION_KEY) != nil {
		sess.SetAttribute(shirosubject.PRINCIPALS_SESSION_KEY, pcs)
	}
	return true
}

// isSameUser ports ShiroHelper.isSameUser: profiles refer to the same user when
// they have the same size and each pairwise profile shares client name and id.
func isSameUser(oldProfiles, newProfiles []profile.UserProfile) bool {
	if oldProfiles == nil || len(oldProfiles) != len(newProfiles) {
		return false
	}
	for i := range newProfiles {
		oldProfile := oldProfiles[i]
		newProfile := newProfiles[i]
		if oldProfile == nil || newProfile == nil {
			return false
		}
		if !profile.AreEqualsString(oldProfile.GetClientName(), newProfile.GetClientName()) {
			return false
		}
		if !profile.AreEqualsString(oldProfile.GetID(), newProfile.GetID()) {
			return false
		}
	}
	return true
}
