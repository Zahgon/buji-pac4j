// Package context ports the subset of org.pac4j.core.context used by the bridge:
// the WebContext contract and the CallContext carrier.
package context

// WebContext ports org.pac4j.core.context.WebContext: the abstraction over an
// HTTP request/response that pac4j operates against. The bridge only passes it
// through to the SessionStore, which for Shiro ignores it and uses the
// thread-bound Subject instead, so the ported contract is intentionally minimal.
type WebContext interface {
	// GetRequestParameter returns a request parameter value and whether present.
	GetRequestParameter(name string) (string, bool)
	// GetRequestAttribute returns a request-scoped attribute and whether present.
	GetRequestAttribute(name string) (any, bool)
	// SetRequestAttribute stores a request-scoped attribute.
	SetRequestAttribute(name string, value any)
}

// noopWebContext is a minimal WebContext used where the Shiro SessionStore
// ignores the context (it resolves state from the thread-bound Subject).
type noopWebContext struct {
	attributes map[string]any
}

// NewNoopWebContext returns a WebContext with no backing request, suitable for
// Shiro's session store which does not consult the context.
func NewNoopWebContext() WebContext {
	return &noopWebContext{attributes: make(map[string]any)}
}

func (c *noopWebContext) GetRequestParameter(string) (string, bool) { return "", false }

func (c *noopWebContext) GetRequestAttribute(name string) (any, bool) {
	v, ok := c.attributes[name]
	return v, ok
}

func (c *noopWebContext) SetRequestAttribute(name string, value any) {
	c.attributes[name] = value
}
