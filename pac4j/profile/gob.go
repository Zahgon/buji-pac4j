package profile

// GobProfile is the serializable image of a CommonProfile used when a principal
// is round-tripped through the Shiro serializer. It captures exactly the fields
// that participate in profile equality (id, client name, remembered flag,
// attributes, roles) so a decoded profile equals its original.
type GobProfile struct {
	ID         string
	ClientName string
	Attributes map[string]any
	Roles      []string
	Remembered bool
}

// ToGobProfile converts a UserProfile into its serializable image. Only the
// CommonProfile concrete type is supported, which is the sole implementation
// used by the bridge and its tests.
func ToGobProfile(p UserProfile) GobProfile {
	gp := GobProfile{
		ID:         p.GetID(),
		ClientName: p.GetClientName(),
		Remembered: p.IsRemembered(),
		Roles:      p.GetRoles(),
	}
	if cp, ok := p.(*CommonProfile); ok {
		gp.Attributes = make(map[string]any, len(cp.attributes))
		for k, v := range cp.attributes {
			gp.Attributes[k] = v
		}
	}
	return gp
}

// FromGobProfile reconstructs a CommonProfile from its serializable image,
// restoring id, client name, remembered flag, attributes, and roles.
func FromGobProfile(gp GobProfile) UserProfile {
	cp := NewCommonProfile()
	cp.SetID(gp.ID)
	cp.SetClientName(gp.ClientName)
	cp.SetRemembered(gp.Remembered)
	for k, v := range gp.Attributes {
		cp.AddAttribute(k, v)
	}
	for _, r := range gp.Roles {
		cp.AddRole(r)
	}
	return cp
}
