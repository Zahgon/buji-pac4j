package profile

import (
	"testing"

	bujirealm "github.com/buji/pac4j/buji/realm"
	bujisubject "github.com/buji/pac4j/buji/subject"
	"github.com/buji/pac4j/pac4j/config"
	pac4jcontext "github.com/buji/pac4j/pac4j/context"
	pac4jprofile "github.com/buji/pac4j/pac4j/profile"
	shiromgt "github.com/buji/pac4j/shiro/mgt"
	shiroutil "github.com/buji/pac4j/shiro/util"
)

func TestShiroProfileManagerSaveAllAndRemove(t *testing.T) {
	sm := shiromgt.NewDefaultSecurityManager(bujirealm.NewPac4jRealm().GetAuthorizingRealm())
	sm.SetSubjectFactory(bujisubject.NewPac4jSubjectFactory())
	shiroutil.Bind(sm)
	defer shiroutil.Remove()

	pm := NewShiroProfileManager(pac4jcontext.NewNoopWebContext(), nil)

	p := pac4jprofile.NewCommonProfile()
	p.SetID("id")
	p.SetClientName("clientName")
	ordered := []config.OrderedProfile{{Key: "clientName", Profile: p}}

	if err := pm.SaveAll(ordered, true); err != nil {
		t.Fatalf("SaveAll returned error: %v", err)
	}
	if !shiroutil.GetSubject().IsAuthenticated() {
		t.Fatalf("expected subject to be authenticated after SaveAll")
	}

	pm.RemoveProfiles()
	if shiroutil.GetSubject().IsAuthenticated() {
		t.Fatalf("expected subject to be logged out after RemoveProfiles")
	}
}
