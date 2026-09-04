// Package util ports the subset of org.apache.shiro.util (ThreadContext) and
// org.apache.shiro.SecurityUtils used by the bridge: binding a SecurityManager to
// the current execution context and resolving the current Subject.
//
// Java uses a thread-local (ThreadContext) to hold the bound SecurityManager and
// the current Subject. Go has no thread-local storage; the tests bind a manager in
// setUp and remove it in tearDown, and each test runs on a single goroutine. This
// port therefore keeps the bound SecurityManager and the lazily-built Subject in
// package-level state guarded by a mutex, which faithfully reproduces the observed
// single-flow behaviour (bind -> use -> remove) of the Shiro tests.
package util

import (
	"sync"

	"github.com/buji/pac4j/shiro/mgt"
	"github.com/buji/pac4j/shiro/subject"
)

// threadState holds the currently bound security manager and the subject resolved
// for it, mirroring the ThreadContext resources map.
type threadState struct {
	mu              sync.Mutex
	securityManager *mgt.DefaultSecurityManager
	currentSubject  subject.Subject
}

// state is the package-level ThreadContext analogue.
var state threadState

// ThreadContext ports org.apache.shiro.util.ThreadContext: the holder of the
// bound SecurityManager (and, implicitly, the current Subject) for the active flow.
type ThreadContext struct{}

// Bind ports ThreadContext.bind(SecurityManager): associates the given manager
// with the current flow and clears any previously resolved subject so a fresh one
// is built on next access.
func (ThreadContext) Bind(sm *mgt.DefaultSecurityManager) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.securityManager = sm
	state.currentSubject = nil
}

// Remove ports ThreadContext.remove(): clears all bound resources.
func (ThreadContext) Remove() {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.securityManager = nil
	state.currentSubject = nil
}

// GetSecurityManager returns the bound SecurityManager, or nil.
func (ThreadContext) GetSecurityManager() *mgt.DefaultSecurityManager {
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.securityManager
}

// Bind is a package-level convenience mirroring ThreadContext.bind, matching the
// static call style used by the Shiro tests (ThreadContext.bind(sm)).
func Bind(sm *mgt.DefaultSecurityManager) { ThreadContext{}.Bind(sm) }

// Remove is a package-level convenience mirroring ThreadContext.remove().
func Remove() { ThreadContext{}.Remove() }

// SecurityUtils ports org.apache.shiro.SecurityUtils: the entry point for
// obtaining the current Subject.
type SecurityUtils struct{}

// GetSubject ports SecurityUtils.getSubject(): returns the Subject bound to the
// current flow, building one from the bound SecurityManager on first access and
// caching it (as Shiro caches the subject on the thread). It panics if no
// SecurityManager is bound, mirroring Shiro's UnavailableSecurityManagerException.
func (SecurityUtils) GetSubject() subject.Subject {
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.currentSubject != nil {
		return state.currentSubject
	}
	if state.securityManager == nil {
		panic("no SecurityManager accessible to the calling code")
	}
	ctx := subject.NewDefaultSubjectContext()
	s := state.securityManager.CreateSubject(ctx)
	state.currentSubject = s
	return s
}

// GetSubject is a package-level convenience mirroring the static
// SecurityUtils.getSubject() call style used across the bridge.
func GetSubject() subject.Subject { return SecurityUtils{}.GetSubject() }
