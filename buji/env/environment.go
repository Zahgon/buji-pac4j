// Package env ports io.buji.pac4j.env.Pac4jIniEnvironment: the Shiro web
// environment that supplies the bundled default framework INI.
package env

import (
	"github.com/buji/pac4j/buji/resources"
	shiroenv "github.com/buji/pac4j/shiro/web/env"
)

// Pac4jIniEnvironment ports io.buji.pac4j.env.Pac4jIniEnvironment: an
// IniWebEnvironment whose framework INI is the embedded buji-pac4j-default.ini.
type Pac4jIniEnvironment struct {
	*shiroenv.IniWebEnvironment
}

// NewPac4jIniEnvironment constructs the environment and registers itself as the
// framework-INI provider.
func NewPac4jIniEnvironment() *Pac4jIniEnvironment {
	e := &Pac4jIniEnvironment{IniWebEnvironment: shiroenv.NewIniWebEnvironment()}
	e.SetFrameworkIniProvider(e)
	return e
}

// GetFrameworkIni ports Pac4jIniEnvironment.getFrameworkIni: it parses the
// embedded default INI (Ini.fromResourcePath("classpath:buji-pac4j-default.ini")).
func (e *Pac4jIniEnvironment) GetFrameworkIni() *shiroenv.Ini {
	return shiroenv.IniFromString(resources.DefaultIni)
}
