package profile

import "testing"

// TestCommonProfileRolesAndRemembered exercises AddRole, GetRoles (sorted
// snapshot) and SetRemembered/IsRemembered.
func TestCommonProfileRolesAndRemembered(t *testing.T) {
	p := NewCommonProfile()
	p.AddRole("b")
	p.AddRole("a")
	p.AddRole("a") // duplicate is a no-op (set semantics)
	roles := p.GetRoles()
	if len(roles) != 2 || roles[0] != "a" || roles[1] != "b" {
		t.Fatalf("expected sorted [a b], got %v", roles)
	}
	if p.IsRemembered() {
		t.Fatalf("new profile must not be remembered")
	}
	p.SetRemembered(true)
	if !p.IsRemembered() {
		t.Fatalf("SetRemembered(true) did not take effect")
	}
}

// TestEqualsAndProfilesEqual exercises Equals, stringSetsEqual, attributesEqual
// and ProfilesEqual across matching and mismatching cases.
func TestEqualsAndProfilesEqual(t *testing.T) {
	mk := func(id, client string, roles ...string) *CommonProfile {
		p := NewCommonProfile()
		p.SetID(id)
		p.SetClientName(client)
		p.AddAttribute("k", "v")
		for _, r := range roles {
			p.AddRole(r)
		}
		return p
	}
	a := mk("1", "c", "r1")
	b := mk("1", "c", "r1")
	if !Equals(a, b) {
		t.Fatalf("identical profiles must be equal")
	}
	if Equals(a, mk("2", "c", "r1")) {
		t.Fatalf("different id must not be equal")
	}
	if Equals(a, mk("1", "d", "r1")) {
		t.Fatalf("different client must not be equal")
	}
	if Equals(a, mk("1", "c", "r2")) {
		t.Fatalf("different roles must not be equal")
	}
	diffAttr := mk("1", "c", "r1")
	diffAttr.AddAttribute("k", "other")
	if Equals(a, diffAttr) {
		t.Fatalf("different attributes must not be equal")
	}
	remembered := mk("1", "c", "r1")
	remembered.SetRemembered(true)
	if Equals(a, remembered) {
		t.Fatalf("different remembered flag must not be equal")
	}
	// nil handling
	if !Equals(nil, nil) {
		t.Fatalf("nil == nil")
	}
	if Equals(a, nil) || Equals(nil, a) {
		t.Fatalf("nil vs non-nil must differ")
	}
	// ProfilesEqual
	if !ProfilesEqual([]UserProfile{a}, []UserProfile{b}) {
		t.Fatalf("equal single-element lists")
	}
	if ProfilesEqual([]UserProfile{a}, []UserProfile{a, b}) {
		t.Fatalf("different length lists must not be equal")
	}
	if ProfilesEqual([]UserProfile{a}, []UserProfile{mk("9", "c")}) {
		t.Fatalf("different content must not be equal")
	}
}

// TestFlatIntoAProfileList verifies a copy preserving order is returned.
func TestFlatIntoAProfileList(t *testing.T) {
	p1 := NewCommonProfile()
	p1.SetID("1")
	p2 := NewCommonProfile()
	p2.SetID("2")
	in := []UserProfile{p1, p2}
	out := FlatIntoAProfileList(in)
	if len(out) != 2 || out[0].GetID() != "1" || out[1].GetID() != "2" {
		t.Fatalf("order not preserved: %v", out)
	}
	// mutating the copy must not affect the input slice header
	out[0] = p2
	if in[0].GetID() != "1" {
		t.Fatalf("FlatIntoAProfileList must return an independent slice")
	}
}

// TestGobRoundTrip exercises ToGobProfile/FromGobProfile preserving equality.
func TestGobRoundTrip(t *testing.T) {
	p := NewCommonProfile()
	p.SetID("id")
	p.SetClientName("client")
	p.SetRemembered(true)
	p.AddAttribute("email", "a@b.c")
	p.AddRole("admin")
	gp := ToGobProfile(p)
	if gp.ID != "id" || gp.ClientName != "client" || !gp.Remembered {
		t.Fatalf("gob image lost scalar fields: %+v", gp)
	}
	back := FromGobProfile(gp)
	if !Equals(p, back) {
		t.Fatalf("gob round-trip must preserve equality")
	}
}

// TestAssertNotNil covers both branches of AssertNotNil.
func TestAssertNotNil(t *testing.T) {
	AssertNotNil("ok", NewCommonProfile()) // must not panic
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("AssertNotNil(nil) must panic")
		}
	}()
	AssertNotNil("bad", nil)
}
