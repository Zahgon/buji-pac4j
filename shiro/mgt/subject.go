// This file ports org.apache.shiro.subject.support.DelegatingSubject: the concrete
// Subject implementation that delegates authentication/authorization to the
// SecurityManager. It lives in package mgt (rather than package subject) so it can
// reference SecurityManager without creating an import cycle (package subject is
// imported by mgt).
package mgt

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/buji/pac4j/shiro/session"
	"github.com/buji/pac4j/shiro/subject"
)

// DelegatingSubject ports org.apache.shiro.subject.support.DelegatingSubject.
type DelegatingSubject struct {
	// principals holds the authenticated identity, or nil when none.
	principals subject.PrincipalCollection
	// authenticated is the full-authentication flag (Subject.isAuthenticated()).
	authenticated bool
	// remembered is the remembered-but-not-authenticated flag
	// (Subject.isRemembered()).
	remembered bool
	// session is the subject's session, created lazily.
	session session.Session
	// sessionCreationEnabled mirrors whether this subject may create a session.
	sessionCreationEnabled bool
	// securityManager performs login/logout/authorization on this subject's
	// behalf.
	securityManager *DefaultSecurityManager
}

// NewDelegatingSubject constructs a subject bound to a security manager.
func NewDelegatingSubject(sm *DefaultSecurityManager) *DelegatingSubject {
	return &DelegatingSubject{securityManager: sm, sessionCreationEnabled: true}
}

// GetPrincipal implements subject.Subject: the primary principal, or nil.
func (s *DelegatingSubject) GetPrincipal() any {
	if s.principals == nil {
		return nil
	}
	return s.principals.GetPrimaryPrincipal()
}

// GetPrincipals implements subject.Subject: the full principal collection, or nil.
func (s *DelegatingSubject) GetPrincipals() subject.PrincipalCollection {
	if s.principals == nil || s.principals.IsEmpty() {
		return nil
	}
	return s.principals
}

// IsAuthenticated implements subject.Subject.
func (s *DelegatingSubject) IsAuthenticated() bool { return s.authenticated }

// IsRemembered implements subject.Subject: remembered principals present and not
// fully authenticated.
func (s *DelegatingSubject) IsRemembered() bool {
	return s.principals != nil && !s.principals.IsEmpty() && !s.authenticated && s.remembered
}

// Login implements subject.Subject by delegating to the security manager and
// adopting the returned authenticated state.
func (s *DelegatingSubject) Login(token subject.AuthenticationToken) error {
	_, err := s.securityManager.Login(s, token)
	return err
}

// Logout implements subject.Subject by delegating to the security manager.
func (s *DelegatingSubject) Logout() { s.securityManager.Logout(s) }

// GetSession implements subject.Subject: returns the session, creating it when
// create is true and creation is enabled.
func (s *DelegatingSubject) GetSession(create bool) session.Session {
	if s.session == nil && create && s.sessionCreationEnabled {
		s.session = session.NewSimpleSession()
	}
	return s.session
}

// SetSession sets the subject's session (used by the subject factory / manager).
func (s *DelegatingSubject) SetSession(sess session.Session) { s.session = sess }

// SetSessionCreationEnabled toggles lazy session creation.
func (s *DelegatingSubject) SetSessionCreationEnabled(enabled bool) {
	s.sessionCreationEnabled = enabled
}

// HasRole implements subject.Subject via the security manager / realm.
func (s *DelegatingSubject) HasRole(role string) bool {
	if s.principals == nil {
		return false
	}
	roles, _ := s.securityManager.GetAuthorizationInfo(s.principals)
	_, ok := roles[role]
	return ok
}

// IsPermitted implements subject.Subject via the security manager / realm.
func (s *DelegatingSubject) IsPermitted(permission string) bool {
	if s.principals == nil {
		return false
	}
	_, perms := s.securityManager.GetAuthorizationInfo(s.principals)
	_, ok := perms[permission]
	return ok
}

// setPrincipals / setAuthenticated / setRemembered are used by the subject factory
// to seed a subject built from a SubjectContext.
func (s *DelegatingSubject) setPrincipals(pc subject.PrincipalCollection) { s.principals = pc }
func (s *DelegatingSubject) setAuthenticated(v bool)                      { s.authenticated = v }
func (s *DelegatingSubject) setRemembered(v bool)                         { s.remembered = v }

// SetPrincipalsCollection sets the subject's principals; used by subject factories
// in other packages to seed a subject built from a SubjectContext.
func (s *DelegatingSubject) SetPrincipalsCollection(pc subject.PrincipalCollection) {
	s.principals = pc
}

// SetAuthenticated sets the full-authentication flag.
func (s *DelegatingSubject) SetAuthenticated(v bool) { s.authenticated = v }

// SetRemembered sets the remembered flag.
func (s *DelegatingSubject) SetRemembered(v bool) { s.remembered = v }

// defaultSubjectFactory ports org.apache.shiro.mgt.DefaultSubjectFactory: builds a
// DelegatingSubject from the context's carried state.
type defaultSubjectFactory struct{}

// CreateSubject implements SubjectFactory.
func (defaultSubjectFactory) CreateSubject(context subject.SubjectContext, sm SecurityManager) subject.Subject {
	dsm, _ := sm.(*DefaultSecurityManager)
	s := NewDelegatingSubject(dsm)
	s.setAuthenticated(context.IsAuthenticated())
	if pc := context.GetPrincipals(); pc != nil {
		s.setPrincipals(pc)
	}
	s.SetSessionCreationEnabled(context.IsSessionCreationEnabled())
	if sess := context.GetSession(); sess != nil {
		s.SetSession(sess)
	}
	return s
}

// newRenewedSessionID generates a fresh session identifier for session-fixation
// protection on login.
func newRenewedSessionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "renewed-00000000000000000000000000000000"
	}
	return hex.EncodeToString(b)
}
