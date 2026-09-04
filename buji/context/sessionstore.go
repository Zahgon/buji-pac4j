// Package context ports io.buji.pac4j.context.ShiroSessionStore: the pac4j
// SessionStore backed by the current Shiro subject's session.
package context

import (
	"fmt"

	pac4jcontext "github.com/buji/pac4j/pac4j/context"
	pac4jsession "github.com/buji/pac4j/pac4j/context/session"
	shirosession "github.com/buji/pac4j/shiro/session"
	shiroutil "github.com/buji/pac4j/shiro/util"
)

// ShiroSessionStore ports io.buji.pac4j.context.ShiroSessionStore: it adapts the
// current Shiro subject's session to the pac4j SessionStore contract.
type ShiroSessionStore struct{}

// INSTANCE ports ShiroSessionStore.INSTANCE: the shared singleton used by the
// framework adapter.
var INSTANCE = &ShiroSessionStore{}

// getSession ports ShiroSessionStore.getSession(boolean): it returns the current
// subject's session, or nil when no session is available (mirroring the
// DisabledSessionException path, here surfaced as a recovered panic).
func (s *ShiroSessionStore) getSession(create bool) (sess shirosession.Session) {
	defer func() {
		if r := recover(); r != nil {
			sess = nil
		}
	}()
	sub := shiroutil.GetSubject()
	return sub.GetSession(create)
}

// GetSessionID ports ShiroSessionStore.getSessionId: the session id as a string,
// present only when a session exists.
func (s *ShiroSessionStore) GetSessionID(ctx pac4jcontext.WebContext, createSession bool) (string, bool) {
	sess := s.getSession(createSession)
	if sess != nil {
		return fmt.Sprint(sess.GetID()), true
	}
	return "", false
}

// Get ports ShiroSessionStore.get: a session attribute, present only when a
// session exists and the attribute is non-nil.
func (s *ShiroSessionStore) Get(ctx pac4jcontext.WebContext, key string) (any, bool) {
	sess := s.getSession(false)
	if sess != nil {
		v := sess.GetAttribute(key)
		return v, v != nil
	}
	return nil, false
}

// Set ports ShiroSessionStore.set: stores a session attribute, creating the
// session if necessary.
func (s *ShiroSessionStore) Set(ctx pac4jcontext.WebContext, key string, value any) {
	sess := s.getSession(true)
	if sess != nil {
		sess.SetAttribute(key, value)
	}
}

// DestroySession ports ShiroSessionStore.destroySession: stops the session.
func (s *ShiroSessionStore) DestroySession(ctx pac4jcontext.WebContext) bool {
	sess := s.getSession(true)
	if sess != nil {
		sess.Stop()
	}
	return true
}

// GetTrackableSession ports ShiroSessionStore.getTrackableSession: always empty.
func (s *ShiroSessionStore) GetTrackableSession(ctx pac4jcontext.WebContext) (any, bool) {
	return nil, false
}

// BuildFromTrackableSession ports
// ShiroSessionStore.buildFromTrackableSession: always empty.
func (s *ShiroSessionStore) BuildFromTrackableSession(ctx pac4jcontext.WebContext, trackableSession any) (pac4jsession.SessionStore, bool) {
	return nil, false
}

// RenewSession ports ShiroSessionStore.renewSession: always false.
func (s *ShiroSessionStore) RenewSession(ctx pac4jcontext.WebContext) bool {
	return false
}
