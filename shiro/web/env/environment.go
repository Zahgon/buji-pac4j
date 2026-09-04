// Package env ports the subset of org.apache.shiro.web.env used by the bridge:
// the INI model and the IniWebEnvironment base whose framework INI can be
// supplied by a subclass hook.
package env

import (
	"bufio"
	"strings"
)

// Ini ports org.apache.shiro.config.Ini: an ordered collection of named
// sections, each mapping keys to values.
type Ini struct {
	// sections maps section name to its key/value pairs.
	sections map[string]map[string]string
	// sectionOrder preserves section insertion order.
	sectionOrder []string
}

// NewIni constructs an empty INI.
func NewIni() *Ini {
	return &Ini{sections: make(map[string]map[string]string)}
}

// IniFromString parses INI text into an Ini, mirroring Ini.fromResourcePath's
// result for the bundled default resource. Lines of the form "[name]" open a
// section; "key = value" lines add entries to the current section; blank lines
// and lines beginning with '#' or ';' are ignored.
func IniFromString(text string) *Ini {
	ini := NewIni()
	current := ""
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			current = strings.TrimSpace(line[1 : len(line)-1])
			ini.ensureSection(current)
			continue
		}
		if idx := strings.Index(line, "="); idx >= 0 {
			key := strings.TrimSpace(line[:idx])
			value := strings.TrimSpace(line[idx+1:])
			ini.ensureSection(current)
			ini.sections[current][key] = value
		}
	}
	return ini
}

// ensureSection creates the named section if it does not yet exist.
func (i *Ini) ensureSection(name string) {
	if _, ok := i.sections[name]; !ok {
		i.sections[name] = make(map[string]string)
		i.sectionOrder = append(i.sectionOrder, name)
	}
}

// GetSection returns the key/value pairs of a section, or nil when absent
// (Ini.getSection(String)).
func (i *Ini) GetSection(name string) map[string]string {
	return i.sections[name]
}

// SectionNames returns section names in insertion order.
func (i *Ini) SectionNames() []string {
	out := make([]string, len(i.sectionOrder))
	copy(out, i.sectionOrder)
	return out
}

// Get returns a value from a section, or "" when absent.
func (i *Ini) Get(section, key string) string {
	if s := i.sections[section]; s != nil {
		return s[key]
	}
	return ""
}

// IsEmpty reports whether the INI has no sections with entries.
func (i *Ini) IsEmpty() bool {
	for _, s := range i.sections {
		if len(s) > 0 {
			return false
		}
	}
	return true
}

// FrameworkIniProvider is the hook a subclass implements to supply the built-in
// framework INI, mirroring the protected IniWebEnvironment.getFrameworkIni()
// template method.
type FrameworkIniProvider interface {
	// GetFrameworkIni returns the framework-supplied INI merged before any
	// user configuration.
	GetFrameworkIni() *Ini
}

// IniWebEnvironment ports org.apache.shiro.web.env.IniWebEnvironment with the
// subset the bridge relies on: it holds an optional framework INI provider and
// the resulting configuration INI.
type IniWebEnvironment struct {
	// provider supplies the framework INI (set by the embedding subclass).
	provider FrameworkIniProvider
	// ini is the resolved configuration.
	ini *Ini
}

// NewIniWebEnvironment constructs a base environment.
func NewIniWebEnvironment() *IniWebEnvironment {
	return &IniWebEnvironment{}
}

// SetFrameworkIniProvider installs the subclass hook that yields the framework
// INI. The buji Pac4jIniEnvironment registers itself here.
func (e *IniWebEnvironment) SetFrameworkIniProvider(p FrameworkIniProvider) {
	e.provider = p
}

// GetFrameworkIni returns the framework INI via the installed provider, or nil
// when none is set. This mirrors the protected getFrameworkIni() hook that
// subclasses override.
func (e *IniWebEnvironment) GetFrameworkIni() *Ini {
	if e.provider != nil {
		return e.provider.GetFrameworkIni()
	}
	return nil
}

// SetIni sets the resolved configuration INI.
func (e *IniWebEnvironment) SetIni(ini *Ini) {
	e.ini = ini
}

// GetIni returns the resolved configuration INI: the explicitly set INI when
// present, otherwise the framework INI.
func (e *IniWebEnvironment) GetIni() *Ini {
	if e.ini != nil {
		return e.ini
	}
	return e.GetFrameworkIni()
}
