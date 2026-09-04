// This file ports org.apache.shiro.subject.support.DefaultSubjectContext.
package subject

import "github.com/buji/pac4j/shiro/session"

// DefaultSubjectContext ports
// org.apache.shiro.subject.support.DefaultSubjectContext.
type DefaultSubjectContext struct {
	authenticated         bool
	authenticationToken   AuthenticationToken
	sessionCreationEnabled bool
	principals            PrincipalCollection
	sess                  session.Session
}

// NewDefaultSubjectContext constructs an empty context with session creation
// enabled by default (matching Shiro's default).
func NewDefaultSubjectContext() *DefaultSubjectContext {
	return &DefaultSubjectContext{sessionCreationEnabled: true}
}

// IsAuthenticated implements SubjectContext.
func (c *DefaultSubjectContext) IsAuthenticated() bool { return c.authenticated }

// SetAuthenticated implements SubjectContext.
func (c *DefaultSubjectContext) SetAuthenticated(authenticated bool) { c.authenticated = authenticated }

// GetAuthenticationToken implements SubjectContext.
func (c *DefaultSubjectContext) GetAuthenticationToken() AuthenticationToken {
	return c.authenticationToken
}

// SetAuthenticationToken implements SubjectContext.
func (c *DefaultSubjectContext) SetAuthenticationToken(token AuthenticationToken) {
	c.authenticationToken = token
}

// SetSessionCreationEnabled implements SubjectContext.
func (c *DefaultSubjectContext) SetSessionCreationEnabled(enabled bool) {
	c.sessionCreationEnabled = enabled
}

// IsSessionCreationEnabled implements SubjectContext.
func (c *DefaultSubjectContext) IsSessionCreationEnabled() bool { return c.sessionCreationEnabled }

// GetPrincipals implements SubjectContext.
func (c *DefaultSubjectContext) GetPrincipals() PrincipalCollection { return c.principals }

// SetPrincipals implements SubjectContext.
func (c *DefaultSubjectContext) SetPrincipals(principals PrincipalCollection) {
	c.principals = principals
}

// GetSession implements SubjectContext.
func (c *DefaultSubjectContext) GetSession() session.Session { return c.sess }

// SetSession implements SubjectContext.
func (c *DefaultSubjectContext) SetSession(s session.Session) { c.sess = s }
