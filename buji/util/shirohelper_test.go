// Package util tests port io.buji.pac4j.util.ShiroHelperTests: they exercise the
// core bridge behaviour of pushing pac4j profiles into the Shiro subject,
// covering session-id renewal on login, in-place profile refresh that keeps the
// session, and re-authentication on identity/client/principal-name changes.
package util

import (
	bujirealm "github.com/buji/pac4j/buji/realm"
	bujisubject "github.com/buji/pac4j/buji/subject"
	bujitoken "github.com/buji/pac4j/buji/token"
	"github.com/buji/pac4j/pac4j/config"
	"github.com/buji/pac4j/pac4j/profile"
	"github.com/buji/pac4j/shiro/authc"
	shgiromgt "github.com/buji/pac4j/shiro/mgt"
	shirorealm "github.com/buji/pac4j/shiro/realm"
	shirosubject "github.com/buji/pac4j/shiro/subject"
	shiroutil "github.com/buji/pac4j/shiro/util"

	"testing"
)

const (
	clientName  = "clientName"
	clientName2 = "clientName2"
	id          = "id"
	accessToken = "access_token"
	email       = "email"
)

// countingSource ports the ShiroHelperTests.CountingPac4jRealm subclass: it wraps
// the real Pac4jRealm authentication hook and counts how many times a genuine
// authentication (Subject.login) occurs.
type countingSource struct {
	nb       int
	delegate *bujirealm.Pac4jRealm
}

// DoGetAuthenticationInfo increments the authentication counter then delegates to
// the wrapped realm (mirroring "nbAuthentications++; return super...").
func (c *countingSource) DoGetAuthenticationInfo(token authc.AuthenticationToken) (authc.AuthenticationInfo, error) {
	c.nb++
	return c.delegate.DoGetAuthenticationInfo(token)
}

// fixture bundles the per-test Shiro wiring that ShiroHelperTests builds in
// setUp: the counting realm, the pac4j realm (for principal-name config), and the
// security manager bound to the current flow.
type fixture struct {
	pr       *bujirealm.Pac4jRealm
	counting *countingSource
}

// setUp ports ShiroHelperTests.setUp: build a CountingPac4jRealm, wrap it in a
// DefaultSecurityManager with the Pac4jSubjectFactory, and bind it to the flow.
func setUp(t *testing.T) *fixture {
	t.Helper()
	pr := bujirealm.NewPac4jRealm()
	counting := &countingSource{delegate: pr}
	ar := shirorealm.NewAuthorizingRealm(counting, pr)
	ar.SetTokenSupports(func(tok authc.AuthenticationToken) bool {
		_, ok := tok.(*bujitoken.Pac4jToken)
		return ok
	})
	sm := shgiromgt.NewDefaultSecurityManager(ar)
	sm.SetSubjectFactory(bujisubject.NewPac4jSubjectFactory())
	shiroutil.Bind(sm)
	return &fixture{pr: pr, counting: counting}
}

// tearDown ports ShiroHelperTests.tearDown: unbind the security manager.
func tearDown() { shiroutil.Remove() }

// makeProfile ports ShiroHelperTests.profile(clientName, id, accessToken).
func makeProfile(cn, profileID, at string) profile.UserProfile {
	p := profile.NewCommonProfile()
	p.SetID(profileID)
	p.SetClientName(cn)
	p.AddAttribute(accessToken, at)
	return p
}

// orderedProfiles ports ShiroHelperTests.profiles(...): a LinkedHashMap keyed by
// client name, preserving insertion order.
func orderedProfiles(profiles ...profile.UserProfile) []config.OrderedProfile {
	out := make([]config.OrderedProfile, 0, len(profiles))
	for _, p := range profiles {
		out = append(out, config.OrderedProfile{Key: p.GetClientName(), Profile: p})
	}
	return out
}

// sessionID ports ShiroHelperTests.sessionId(): the current subject's session id.
func sessionID(t *testing.T) string {
	t.Helper()
	id := shiroutil.GetSubject().GetSession(true).GetID()
	s, ok := id.(string)
	if !ok {
		t.Fatalf("session id is not a string: %v", id)
	}
	return s
}

// shiroProfile ports ShiroHelperTests.shiroProfile(): the profile carried by the
// Pac4jPrincipal in the subject's principal collection.
func shiroProfile(t *testing.T) profile.UserProfile {
	t.Helper()
	pcs := shiroutil.GetSubject().GetPrincipals()
	if pcs == nil {
		t.Fatalf("no principals on subject")
	}
	pr := pcs.OneByType((*bujisubject.Pac4jPrincipal)(nil))
	if pr == nil {
		t.Fatalf("no Pac4jPrincipal in collection")
	}
	return pr.(*bujisubject.Pac4jPrincipal).GetProfile()
}

// shiroPrincipal returns the Pac4jPrincipal object from the subject's principal
// collection.
func shiroPrincipal(t *testing.T) *bujisubject.Pac4jPrincipal {
	t.Helper()
	pcs := shiroutil.GetSubject().GetPrincipals()
	if pcs == nil {
		t.Fatalf("no principals on subject")
	}
	pr := pcs.OneByType((*bujisubject.Pac4jPrincipal)(nil))
	if pr == nil {
		t.Fatalf("no Pac4jPrincipal in collection")
	}
	return pr.(*bujisubject.Pac4jPrincipal)
}

// sessionProfile ports ShiroHelperTests.sessionProfile(): the profile carried by
// the Pac4jPrincipal stored under the session's principals key.
func sessionProfile(t *testing.T) profile.UserProfile {
	t.Helper()
	sess := shiroutil.GetSubject().GetSession(true)
	stored := sess.GetAttribute(shirosubject.PRINCIPALS_SESSION_KEY)
	if stored == nil {
		t.Fatalf("no principals stored in session")
	}
	pcs, ok := stored.(shirosubject.PrincipalCollection)
	if !ok {
		t.Fatalf("session principals are not a PrincipalCollection: %T", stored)
	}
	pr := pcs.OneByType((*bujisubject.Pac4jPrincipal)(nil))
	if pr == nil {
		t.Fatalf("no Pac4jPrincipal in session collection")
	}
	return pr.(*bujisubject.Pac4jPrincipal).GetProfile()
}

// mustPopulate runs the core bridge call, failing the test on error.
func mustPopulate(t *testing.T, profiles []config.OrderedProfile) {
	t.Helper()
	if err := PopulateSubject(profiles); err != nil {
		t.Fatalf("PopulateSubject failed: %v", err)
	}
}

// TestLoginRenewsTheSessionId ports testLoginRenewsTheSessionId: a genuine login
// authenticates the subject, counts one authentication, and renews the session id.
func TestLoginRenewsTheSessionId(t *testing.T) {
	f := setUp(t)
	defer tearDown()

	sessionIDBefore := sessionID(t)
	mustPopulate(t, orderedProfiles(makeProfile(clientName, id, "at1")))

	if !shiroutil.GetSubject().IsAuthenticated() {
		t.Fatalf("subject should be authenticated")
	}
	if f.counting.nb != 1 {
		t.Fatalf("nbAuthentications = %d, want 1", f.counting.nb)
	}
	if sessionID(t) == sessionIDBefore {
		t.Fatalf("session id should have been renewed on login")
	}
}

// TestProfileRenewalKeepsTheSessionId ports testProfileRenewalKeepsTheSessionId:
// re-populating with the same identity but a new access token refreshes the
// profile in place, keeping the session id and avoiding a second authentication.
func TestProfileRenewalKeepsTheSessionId(t *testing.T) {
	f := setUp(t)
	defer tearDown()

	mustPopulate(t, orderedProfiles(makeProfile(clientName, id, "at1")))
	sessionIDAfterLogin := sessionID(t)

	mustPopulate(t, orderedProfiles(makeProfile(clientName, id, "at2")))

	if sessionID(t) != sessionIDAfterLogin {
		t.Fatalf("session id should be unchanged on profile refresh")
	}
	if !shiroutil.GetSubject().IsAuthenticated() {
		t.Fatalf("subject should be authenticated")
	}
	if f.counting.nb != 1 {
		t.Fatalf("nbAuthentications = %d, want 1", f.counting.nb)
	}
	if got := shiroProfile(t).GetAttribute(accessToken); got != "at2" {
		t.Fatalf("shiro profile access_token = %v, want at2", got)
	}
	if got := sessionProfile(t).GetAttribute(accessToken); got != "at2" {
		t.Fatalf("session profile access_token = %v, want at2", got)
	}
}

// TestProfileRenewalKeepsTheSessionIdForMultiProfiles ports
// testProfileRenewalKeepsTheSessionIdForMultiProfiles: refreshing one of two
// same-identity profiles keeps the session and does not re-authenticate.
func TestProfileRenewalKeepsTheSessionIdForMultiProfiles(t *testing.T) {
	f := setUp(t)
	defer tearDown()

	mustPopulate(t, orderedProfiles(
		makeProfile(clientName, id, "at1"),
		makeProfile(clientName2, id, "at1"),
	))
	sessionIDAfterLogin := sessionID(t)

	mustPopulate(t, orderedProfiles(
		makeProfile(clientName, id, "at2"),
		makeProfile(clientName2, id, "at1"),
	))

	if sessionID(t) != sessionIDAfterLogin {
		t.Fatalf("session id should be unchanged on profile refresh")
	}
	if f.counting.nb != 1 {
		t.Fatalf("nbAuthentications = %d, want 1", f.counting.nb)
	}
	if got := len(shiroPrincipal(t).GetProfiles()); got != 2 {
		t.Fatalf("principal profiles size = %d, want 2", got)
	}
	if got := shiroProfile(t).GetAttribute(accessToken); got != "at2" {
		t.Fatalf("shiro profile access_token = %v, want at2", got)
	}
}

// TestNewIdentityRenewsTheSessionId ports testNewIdentityRenewsTheSessionId: a
// changed profile id is a new user, forcing re-authentication and session renewal.
func TestNewIdentityRenewsTheSessionId(t *testing.T) {
	f := setUp(t)
	defer tearDown()

	mustPopulate(t, orderedProfiles(makeProfile(clientName, id, "at1")))
	sessionIDAfterLogin := sessionID(t)

	mustPopulate(t, orderedProfiles(makeProfile(clientName, "id2", "at1")))

	if f.counting.nb != 2 {
		t.Fatalf("nbAuthentications = %d, want 2", f.counting.nb)
	}
	if got := shiroProfile(t).GetID(); got != "id2" {
		t.Fatalf("shiro profile id = %q, want id2", got)
	}
	if sessionID(t) == sessionIDAfterLogin {
		t.Fatalf("session id should have been renewed on new identity")
	}
}

// TestNewPrincipalNameRenewsTheSessionId ports
// testNewPrincipalNameRenewsTheSessionId: when the computed principal name changes
// (via the email attribute), the refresh is rejected and a new login occurs.
func TestNewPrincipalNameRenewsTheSessionId(t *testing.T) {
	f := setUp(t)
	defer tearDown()

	f.pr.SetPrincipalNameAttribute(email)

	p1 := makeProfile(clientName, id, "at1")
	p1.AddAttribute(email, "john@example.com")
	mustPopulate(t, orderedProfiles(p1))

	if got := shiroutil.GetSubject().GetPrincipal(); got != "john@example.com" {
		t.Fatalf("principal = %v, want john@example.com", got)
	}
	sessionIDAfterLogin := sessionID(t)

	p2 := makeProfile(clientName, id, "at2")
	p2.AddAttribute(email, "jane@example.com")
	mustPopulate(t, orderedProfiles(p2))

	if f.counting.nb != 2 {
		t.Fatalf("nbAuthentications = %d, want 2", f.counting.nb)
	}
	if got := shiroutil.GetSubject().GetPrincipal(); got != "jane@example.com" {
		t.Fatalf("principal = %v, want jane@example.com", got)
	}
	if sessionID(t) == sessionIDAfterLogin {
		t.Fatalf("session id should have been renewed on new principal name")
	}
}

// TestNewClientRenewsTheSessionId ports testNewClientRenewsTheSessionId: a
// changed client name is a new user, forcing re-authentication and session renewal.
func TestNewClientRenewsTheSessionId(t *testing.T) {
	f := setUp(t)
	defer tearDown()

	mustPopulate(t, orderedProfiles(makeProfile(clientName, id, "at1")))
	sessionIDAfterLogin := sessionID(t)

	mustPopulate(t, orderedProfiles(makeProfile(clientName2, id, "at1")))

	if f.counting.nb != 2 {
		t.Fatalf("nbAuthentications = %d, want 2", f.counting.nb)
	}
	if sessionID(t) == sessionIDAfterLogin {
		t.Fatalf("session id should have been renewed on new client")
	}
}
