package env

import "testing"

// TestIniParsingAndAccessors covers NewIni, IniFromString, ensureSection,
// GetSection, SectionNames, Get and IsEmpty.
func TestIniParsingAndAccessors(t *testing.T) {
	if !NewIni().IsEmpty() {
		t.Fatalf("empty INI must report empty")
	}
	text := "# comment\n; also comment\n[main]\nkey = value\nother=x\n[second]\na = b\n"
	ini := IniFromString(text)
	if ini.IsEmpty() {
		t.Fatalf("parsed INI must not be empty")
	}
	sec := ini.GetSection("main")
	if sec == nil || sec["key"] != "value" || sec["other"] != "x" {
		t.Fatalf("main section parsed wrong: %v", sec)
	}
	if ini.GetSection("missing") != nil {
		t.Fatalf("absent section must be nil")
	}
	if ini.Get("main", "key") != "value" {
		t.Fatalf("Get(main,key) wrong")
	}
	if ini.Get("main", "absent") != "" || ini.Get("nope", "k") != "" {
		t.Fatalf("Get for absent keys/sections must be empty")
	}
	names := ini.SectionNames()
	if len(names) != 2 || names[0] != "main" || names[1] != "second" {
		t.Fatalf("section order wrong: %v", names)
	}
}

// stubProvider implements FrameworkIniProvider for the environment test.
type stubProvider struct{ ini *Ini }

func (s *stubProvider) GetFrameworkIni() *Ini { return s.ini }

// TestIniWebEnvironment covers NewIniWebEnvironment, SetFrameworkIniProvider,
// GetFrameworkIni, SetIni and GetIni (both branches).
func TestIniWebEnvironment(t *testing.T) {
	e := NewIniWebEnvironment()
	if e.GetFrameworkIni() != nil {
		t.Fatalf("no provider => nil framework INI")
	}
	if e.GetIni() != nil {
		t.Fatalf("no ini and no provider => nil")
	}
	fw := IniFromString("[main]\nk = v\n")
	e.SetFrameworkIniProvider(&stubProvider{ini: fw})
	if e.GetFrameworkIni() != fw {
		t.Fatalf("framework INI should come from provider")
	}
	// GetIni falls back to framework INI when no explicit INI set
	if e.GetIni() != fw {
		t.Fatalf("GetIni should return framework INI as fallback")
	}
	explicit := NewIni()
	e.SetIni(explicit)
	if e.GetIni() != explicit {
		t.Fatalf("GetIni should prefer the explicitly set INI")
	}
}
