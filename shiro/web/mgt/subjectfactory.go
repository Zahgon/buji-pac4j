// Package mgt (shiro/web/mgt) ports the subset of
// org.apache.shiro.web.mgt used by the bridge: DefaultWebSubjectFactory, the base
// subject factory that Pac4jSubjectFactory extends.
package mgt

import (
	shiromgt "github.com/buji/pac4j/shiro/mgt"
	"github.com/buji/pac4j/shiro/subject"
)

// DefaultWebSubjectFactory ports
// org.apache.shiro.web.mgt.DefaultWebSubjectFactory. For the behaviour observed by
// the bridge it builds a subject from the context exactly like the default subject
// factory. Pac4jSubjectFactory embeds this type and calls CreateSubject as "super".
type DefaultWebSubjectFactory struct{}

// NewDefaultWebSubjectFactory constructs the base web subject factory.
func NewDefaultWebSubjectFactory() *DefaultWebSubjectFactory {
	return &DefaultWebSubjectFactory{}
}

// CreateSubject ports DefaultWebSubjectFactory.createSubject(SubjectContext): it
// builds a DelegatingSubject seeded from the context's carried state. This is the
// "super.createSubject(context)" invoked by Pac4jSubjectFactory.
func (f *DefaultWebSubjectFactory) CreateSubject(context subject.SubjectContext, sm shiromgt.SecurityManager) subject.Subject {
	dsm, _ := sm.(*shiromgt.DefaultSecurityManager)
	s := shiromgt.NewDelegatingSubject(dsm)
	s.SetAuthenticated(context.IsAuthenticated())
	if pc := context.GetPrincipals(); pc != nil {
		s.SetPrincipalsCollection(pc)
	}
	s.SetSessionCreationEnabled(context.IsSessionCreationEnabled())
	if sess := context.GetSession(); sess != nil {
		s.SetSession(sess)
	}
	return s
}
