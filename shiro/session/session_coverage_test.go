package session

import "testing"

// TestSimpleSessionLifecycle exercises the SimpleSession accessors: id renewal,
// attribute set/get/remove, stop/is-stopped, and the attribute snapshot.
func TestSimpleSessionLifecycle(t *testing.T) {
	s := NewSimpleSession()

	id1 := s.GetID().(string)
	if id1 == "" {
		t.Fatalf("expected a non-empty session id")
	}
	s.SetID("renewed-id")
	if got := s.GetID().(string); got != "renewed-id" {
		t.Fatalf("SetID: got %q want %q", got, "renewed-id")
	}

	if v := s.GetAttribute("missing"); v != nil {
		t.Fatalf("missing attribute should be nil, got %v", v)
	}
	s.SetAttribute("k", "v")
	if got := s.GetAttribute("k"); got != "v" {
		t.Fatalf("GetAttribute: got %v want v", got)
	}

	snap := s.Attributes()
	if snap["k"] != "v" {
		t.Fatalf("Attributes snapshot missing k=v: %v", snap)
	}
	// Snapshot must be an independent copy.
	snap["k"] = "mutated"
	if s.GetAttribute("k") != "v" {
		t.Fatalf("Attributes snapshot must not alias the session state")
	}

	if removed := s.RemoveAttribute("k"); removed != "v" {
		t.Fatalf("RemoveAttribute: got %v want v", removed)
	}
	if s.GetAttribute("k") != nil {
		t.Fatalf("attribute should be gone after RemoveAttribute")
	}

	if s.IsStopped() {
		t.Fatalf("session should not start stopped")
	}
	s.Stop()
	if !s.IsStopped() {
		t.Fatalf("session should be stopped after Stop()")
	}
}
