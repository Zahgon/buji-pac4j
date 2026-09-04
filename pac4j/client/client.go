// Package client ports the minimal subset of org.pac4j.core.client used by the
// default configuration: the Clients holder.
package client

// Clients ports org.pac4j.core.client.Clients: a container of authentication
// clients. The default INI wires an empty Clients into the Config; the bridge
// itself does not drive client-based authentication in the migrated tests, so
// the ported type carries the client list without the full pac4j redirection
// machinery.
type Clients struct {
	// Clients holds the configured client instances (by name).
	Clients map[string]any
}

// NewClients constructs an empty Clients holder.
func NewClients() *Clients {
	return &Clients{Clients: make(map[string]any)}
}
