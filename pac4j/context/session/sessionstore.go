// Package session ports org.pac4j.core.context.session.SessionStore: the pac4j
// contract for reading and writing session-scoped state. The Shiro bridge
// implements this against the thread-bound Shiro Subject.
package session

import "github.com/buji/pac4j/pac4j/context"

// SessionStore ports org.pac4j.core.context.session.SessionStore. Method
// signatures mirror the Java Optional-returning API using (value, ok) pairs.
type SessionStore interface {
	// GetSessionID returns the session id, creating the session when
	// createSession is true (getSessionId).
	GetSessionID(webContext context.WebContext, createSession bool) (string, bool)
	// Get returns a session attribute value and whether present (get).
	Get(webContext context.WebContext, key string) (any, bool)
	// Set stores a session attribute (set).
	Set(webContext context.WebContext, key string, value any)
	// DestroySession stops the session (destroySession).
	DestroySession(webContext context.WebContext) bool
	// GetTrackableSession returns a trackable session token (getTrackableSession).
	GetTrackableSession(webContext context.WebContext) (any, bool)
	// BuildFromTrackableSession rebuilds a store from a trackable session
	// (buildFromTrackableSession).
	BuildFromTrackableSession(webContext context.WebContext, trackableSession any) (SessionStore, bool)
	// RenewSession renews the session (renewSession).
	RenewSession(webContext context.WebContext) bool
}
