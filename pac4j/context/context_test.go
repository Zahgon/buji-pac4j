package context

import "testing"

// TestNoopWebContext covers NewNoopWebContext, GetRequestParameter,
// GetRequestAttribute and SetRequestAttribute.
func TestNoopWebContext(t *testing.T) {
	c := NewNoopWebContext()
	if v, ok := c.GetRequestParameter("any"); ok || v != "" {
		t.Fatalf("noop context has no parameters, got (%q,%v)", v, ok)
	}
	if _, ok := c.GetRequestAttribute("missing"); ok {
		t.Fatalf("missing attribute must report not-present")
	}
	c.SetRequestAttribute("k", 42)
	v, ok := c.GetRequestAttribute("k")
	if !ok || v != 42 {
		t.Fatalf("stored attribute not returned, got (%v,%v)", v, ok)
	}
}
