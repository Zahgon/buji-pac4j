package env

import "testing"

func TestPac4jIniEnvironmentFrameworkIni(t *testing.T) {
	e := NewPac4jIniEnvironment()
	if e == nil {
		t.Fatalf("NewPac4jIniEnvironment returned nil")
	}

	ini := e.GetFrameworkIni()
	if ini == nil {
		t.Fatalf("GetFrameworkIni returned nil")
	}
	if ini.IsEmpty() {
		t.Fatalf("framework ini must not be empty")
	}
	if ini.GetSection("main") == nil {
		t.Fatalf("framework ini must contain a [main] section")
	}
}
