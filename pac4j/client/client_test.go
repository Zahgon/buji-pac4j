package client

import "testing"

// TestNewClients covers NewClients returning an empty, usable holder.
func TestNewClients(t *testing.T) {
	c := NewClients()
	if c.Clients == nil {
		t.Fatalf("NewClients must initialise the map")
	}
	if len(c.Clients) != 0 {
		t.Fatalf("new Clients must be empty")
	}
	c.Clients["cas"] = struct{}{}
	if len(c.Clients) != 1 {
		t.Fatalf("Clients map must be writable")
	}
}
