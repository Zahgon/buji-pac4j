// Package mgt ports the subset of org.apache.shiro.mgt used by the bridge: the
// SecurityManager contract and DefaultSecurityManager, including the
// session-fixation protection (session id renewal) applied on a successful login,
// which Shiro 2.2+ performs via beforeSuccessfulLogin.
package mgt

import (
	"github.com/buji/pac4j/shiro/authc"
	"github.com/buji/pac4j/shiro/realm"
	"github.com/buji/pac4j/shiro/session"
	"github.com/buji/pac4j/shiro/subject"
)

// SecurityManager ports org.apache.shiro.mgt.SecurityManager as consumed by the
// bridge: it authenticates subjects, creates them from a context, and logs them
// out.
type SecurityManager interface {
	// Login authenticates the subject against the realm and returns the resulting
	// authenticated subject, or an error mirroring a thrown
	// AuthenticationException (SecurityManager.login(Subject, AuthenticationToken)).
	Login(current subject.Subject, token subject.AuthenticationToken) (subject.Subject, error)
	// Logout clears the subject's authenticated state
	// (SecurityManager.logout(Subject)).
	Logout(s subject.Subject)
	// CreateSubject builds a Subject from a context
	// (SecurityManager.createSubject(SubjectContext)).
	CreateSubject(context subject.SubjectContext) subject.Subject
	// GetAuthorizationInfo exposes the realm's authorization lookup for a
	// principal collection, backing Subject.hasRole / isPermitted.
	GetAuthorizationInfo(principals subject.PrincipalCollection) (roles map[string]struct{}, permissions map[string]struct{})
}

// SubjectFactory ports org.apache.shiro.mgt.SubjectFactory: builds a Subject from
// a SubjectContext.
type SubjectFactory interface {
	// CreateSubject builds a Subject (SubjectFactory.createSubject(SubjectContext)).
	CreateSubject(context subject.SubjectContext, sm SecurityManager) subject.Subject
}

// DefaultSecurityManager ports org.apache.shiro.mgt.DefaultSecurityManager
// (the web variant reuses the same login path for these tests).
type DefaultSecurityManager struct {
	// realm is the single realm consulted for authentication/authorization.
	realm *realm.AuthorizingRealm
	// subjectFactory builds subjects; Pac4jSubjectFactory is installed by the
	// tests and by the default ini.
	subjectFactory SubjectFactory
}

// NewDefaultSecurityManager mirrors new DefaultSecurityManager(Realm): wires the
// realm and installs the default subject factory.
func NewDefaultSecurityManager(r *realm.AuthorizingRealm) *DefaultSecurityManager {
	return &DefaultSecurityManager{
		realm:          r,
		subjectFactory: defaultSubjectFactory{},
	}
}

// SetSubjectFactory mirrors DefaultSecurityManager.setSubjectFactory.
func (m *DefaultSecurityManager) SetSubjectFactory(f SubjectFactory) { m.subjectFactory = f }

// GetRealm returns the configured realm.
func (m *DefaultSecurityManager) GetRealm() *realm.AuthorizingRealm { return m.realm }

// CreateSubject implements SecurityManager by delegating to the subject factory.
func (m *DefaultSecurityManager) CreateSubject(context subject.SubjectContext) subject.Subject {
	return m.subjectFactory.CreateSubject(context, m)
}

// Login implements SecurityManager. It authenticates the token via the realm,
// and on success:
//   - renews the session id (Shiro 2.2+ beforeSuccessfulLogin, session-fixation
//     protection): a new session id is generated while attributes are preserved;
//   - binds the resulting principals into the subject and marks it authenticated;
//   - stores the principals into the session under PRINCIPALS_SESSION_KEY.
//
// On failure it returns the (wrapped) AuthenticationException and leaves the
// subject unauthenticated.
func (m *DefaultSecurityManager) Login(current subject.Subject, token subject.AuthenticationToken) (subject.Subject, error) {
	authcToken, ok := token.(authc.AuthenticationToken)
	if !ok {
		return current, &authc.AuthenticationException{Message: "unsupported token"}
	}
	if !m.realm.Supports(authcToken) {
		return current, &authc.AuthenticationException{Message: "realm does not support token"}
	}
	info, err := m.realm.GetAuthenticationInfo(authcToken)
	if err != nil {
		return current, err
	}
	if info == nil {
		return current, &authc.AuthenticationException{Message: "no account data found"}
	}

	ds, ok := current.(*DelegatingSubject)
	if !ok {
		// The bridge always operates on DelegatingSubject instances.
		return current, &authc.AuthenticationException{Message: "unexpected subject type"}
	}

	// Session-fixation protection: renew the session id before associating the
	// authenticated principals with it (beforeSuccessfulLogin in Shiro 2.2+).
	if sess := ds.session; sess != nil {
		if ss, ok := sess.(*session.SimpleSession); ok {
			ss.SetID(newRenewedSessionID())
		}
	}

	ds.principals = info.GetPrincipals()
	ds.authenticated = true

	// Persist the principals in the session, mirroring Shiro storing the
	// principal collection under DefaultSubjectContext.PRINCIPALS_SESSION_KEY.
	sess := ds.GetSession(true)
	if sess != nil {
		sess.SetAttribute(subject.PRINCIPALS_SESSION_KEY, info.GetPrincipals())
		sess.SetAttribute(subject.AUTHENTICATED_SESSION_KEY, true)
	}
	return ds, nil
}

// Logout implements SecurityManager: clears principals and authenticated state and
// stops the session.
func (m *DefaultSecurityManager) Logout(s subject.Subject) {
	ds, ok := s.(*DelegatingSubject)
	if !ok {
		return
	}
	if ds.session != nil {
		ds.session.RemoveAttribute(subject.PRINCIPALS_SESSION_KEY)
		ds.session.RemoveAttribute(subject.AUTHENTICATED_SESSION_KEY)
		ds.session.Stop()
		ds.session = nil
	}
	ds.principals = nil
	ds.authenticated = false
	ds.remembered = false
}

// GetAuthorizationInfo implements SecurityManager: it consults the realm for the
// roles and permissions of the given principals.
func (m *DefaultSecurityManager) GetAuthorizationInfo(principals subject.PrincipalCollection) (map[string]struct{}, map[string]struct{}) {
	roles := make(map[string]struct{})
	perms := make(map[string]struct{})
	if principals == nil {
		return roles, perms
	}
	info := m.realm.GetAuthorizationInfo(principals)
	if info == nil {
		return roles, perms
	}
	for _, r := range info.GetRoles() {
		roles[r] = struct{}{}
	}
	for _, p := range info.GetStringPermissions() {
		perms[p] = struct{}{}
	}
	return roles, perms
}
