// This file ports io.buji.pac4j.subject.Pac4jSubjectFactory: the subject factory
// that adjusts authentication state for remember-me tokens before delegating to
// the web subject factory.
package subject

import (
	shiromgt "github.com/buji/pac4j/shiro/mgt"
	"github.com/buji/pac4j/shiro/subject"
	webmgt "github.com/buji/pac4j/shiro/web/mgt"
	"github.com/buji/pac4j/buji/token"
)

// Pac4jSubjectFactory ports io.buji.pac4j.subject.Pac4jSubjectFactory: a
// DefaultWebSubjectFactory that demotes remember-me logins to non-authenticated
// so Shiro treats them as remembered rather than fully authenticated.
type Pac4jSubjectFactory struct {
	*webmgt.DefaultWebSubjectFactory
}

// NewPac4jSubjectFactory constructs the factory.
func NewPac4jSubjectFactory() *Pac4jSubjectFactory {
	return &Pac4jSubjectFactory{DefaultWebSubjectFactory: webmgt.NewDefaultWebSubjectFactory()}
}

// CreateSubject ports Pac4jSubjectFactory.createSubject: when the context is
// authenticated by a remember-me Pac4jToken it clears the authenticated flag,
// always enables session creation, then delegates to the web subject factory.
func (f *Pac4jSubjectFactory) CreateSubject(context subject.SubjectContext, sm shiromgt.SecurityManager) subject.Subject {
	authenticated := context.IsAuthenticated()
	if authenticated {
		tok := context.GetAuthenticationToken()
		if pt, ok := tok.(*token.Pac4jToken); ok && pt.IsRememberMe() {
			context.SetAuthenticated(false)
		}
	}
	context.SetSessionCreationEnabled(true)
	return f.DefaultWebSubjectFactory.CreateSubject(context, sm)
}
