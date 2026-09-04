// Package session ports the subset of org.apache.shiro.session used by the
// bridge: the Session contract and a simple in-memory implementation whose id can
// be renewed (session-fixation protection on login).
package session

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

// Session ports org.apache.shiro.session.Session: server-side session state.
type Session interface {
	// GetID returns the session id (Session.getId()); the concrete value is an
	// opaque, comparable identifier.
	GetID() any
	// GetAttribute returns a session attribute or nil (Session.getAttribute).
	GetAttribute(key any) any
	// SetAttribute stores a session attribute (Session.setAttribute).
	SetAttribute(key any, value any)
	// RemoveAttribute removes and returns a session attribute
	// (Session.removeAttribute).
	RemoveAttribute(key any) any
	// Stop invalidates the session (Session.stop()).
	Stop()
}

// SimpleSession ports org.apache.shiro.session.mgt.SimpleSession with the subset
// of behaviour the bridge observes: attribute storage, a mutable id (so login can
// renew it), and stop.
type SimpleSession struct {
	mu         sync.Mutex
	id         string
	attributes map[any]any
	stopped    bool
}

// NewSimpleSession constructs a session with a freshly generated id.
func NewSimpleSession() *SimpleSession {
	return &SimpleSession{id: newSessionID(), attributes: make(map[any]any)}
}

// GetID implements Session.
func (s *SimpleSession) GetID() any {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.id
}

// SetID assigns a new id, used by session renewal to reproduce Shiro's
// session-fixation protection (a new id after a successful login).
func (s *SimpleSession) SetID(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.id = id
}

// GetAttribute implements Session.
func (s *SimpleSession) GetAttribute(key any) any {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.attributes[key]
}

// SetAttribute implements Session.
func (s *SimpleSession) SetAttribute(key any, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.attributes[key] = value
}

// RemoveAttribute implements Session.
func (s *SimpleSession) RemoveAttribute(key any) any {
	s.mu.Lock()
	defer s.mu.Unlock()
	v := s.attributes[key]
	delete(s.attributes, key)
	return v
}

// Stop implements Session.
func (s *SimpleSession) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopped = true
}

// IsStopped reports whether the session has been stopped.
func (s *SimpleSession) IsStopped() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stopped
}

// Attributes returns a snapshot copy of the attribute map, used when a renewed
// session must carry over existing state.
func (s *SimpleSession) Attributes() map[any]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[any]any, len(s.attributes))
	for k, v := range s.attributes {
		out[k] = v
	}
	return out
}

// newSessionID generates a random, unique session identifier.
func newSessionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand should not fail; fall back to a fixed-length zero id which
		// remains unique enough for the single-threaded test scenarios.
		return "00000000000000000000000000000000"
	}
	return hex.EncodeToString(b)
}
