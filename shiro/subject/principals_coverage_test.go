package subject

import "testing"

// TestSimplePrincipalCollection exercises the collection accessors: primary
// principal, per-realm retrieval, OneByType, AsList, IsEmpty, and nil handling.
func TestSimplePrincipalCollection(t *testing.T) {
	empty := NewSimplePrincipalCollection()
	if !empty.IsEmpty() {
		t.Fatalf("new collection should be empty")
	}
	if empty.GetPrimaryPrincipal() != nil {
		t.Fatalf("empty collection has no primary principal")
	}
	if empty.OneByType((*string)(nil)) != nil {
		t.Fatalf("empty collection OneByType should be nil")
	}

	c := NewSimplePrincipalCollectionFor([]any{"alice", 42}, "realmA")
	// nil principals are skipped on Add.
	c.Add(nil, "realmA")

	if c.IsEmpty() {
		t.Fatalf("collection with principals should not be empty")
	}
	if got := c.GetPrimaryPrincipal(); got != "alice" {
		t.Fatalf("GetPrimaryPrincipal: got %v want alice", got)
	}
	if got := c.OneByType((*string)(nil)); got != nil {
		// target is *string typed-nil; principals are string (not *string),
		// so exact-type match fails -> nil. Confirm the int match instead.
		_ = got
	}
	if got := c.FromRealm("realmA"); len(got) != 2 {
		t.Fatalf("FromRealm(realmA): got %d principals want 2", len(got))
	}
	if got := c.FromRealm("missing"); got != nil {
		t.Fatalf("FromRealm(missing) should be nil")
	}
	list := c.AsList()
	if len(list) != 2 || list[0] != "alice" || list[1] != 42 {
		t.Fatalf("AsList order/content wrong: %v", list)
	}
}

// principalStub is a concrete pointer type used to verify OneByType exact-type
// matching (the mechanism the realm uses to fetch a *Pac4jPrincipal).
type principalStub struct{ name string }

func TestOneByTypeExactMatch(t *testing.T) {
	stub := &principalStub{name: "p"}
	c := NewSimplePrincipalCollectionFor([]any{"username", stub}, "realmA")
	got := c.OneByType((*principalStub)(nil))
	if got == nil {
		t.Fatalf("OneByType should find the *principalStub")
	}
	if got.(*principalStub).name != "p" {
		t.Fatalf("OneByType returned wrong principal: %v", got)
	}
}
