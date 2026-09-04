// Package resources embeds the default framework configuration shipped with the
// bridge, mirroring src/main/resources/buji-pac4j-default.ini in the original
// project.
package resources

import _ "embed"

// DefaultIni holds the contents of buji-pac4j-default.ini: the default Shiro INI
// configuration loaded by the bridge's web environment.
//
//go:embed buji-pac4j-default.ini
var DefaultIni string
