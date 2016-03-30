package teamsnap

import (
	"os"
	"strings"
	"testing"
)

var authToken = os.Getenv("AuthToken")

func TestAuthToken(t *testing.T) {
	if authToken == "" {
		t.Error("Please set AuthToken environment variable.")
	}
}

func TestInitialize(t *testing.T) {
	teamSnap := &TeamSnap{AuthToken: authToken}
	teamSnap.Initialize()
	if !strings.HasPrefix(teamSnap.Version, "3.") {
		t.Error("Version Changed: Expected 3.206.5 but got", teamSnap.Version)
	}
}

func TestFindHref(t *testing.T) {
	teamSnap := &TeamSnap{AuthToken: authToken}
	teamSnap.Initialize()
	href, err := teamSnap.findHref("me", teamSnap.RootLinks)
	if err != nil {
		t.Error("Error Occured.", err.Error())
		return
	}
	if href != "https://api.teamsnap.com/v3/me" {
		t.Error("findHref for me expected: https://api.teamsnap.com/v3/me - returned: ", href)
	}
}

func TestMe(t *testing.T) {
	teamSnap := &TeamSnap{AuthToken: authToken}
	teamSnap.Initialize()
	resp, err := teamSnap.FetchRoot("me")
	if err != nil {
		t.Error("Error Occured.", err.Error())
		return
	}
	items := resp.Collection.Items
	if len(items) == 0 {
		t.Error("No Items returned")
		return
	}
	emailAddress, _ := items[0].DataValueString("email")
	if emailAddress == "" {
		t.Error("No email address found in 'me' response")
	}
}
