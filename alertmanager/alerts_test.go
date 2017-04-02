package alertmanager_test

import (
	"fmt"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/cloudflare/unsee/alertmanager"
	"github.com/cloudflare/unsee/mock"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestAlertGroupsAPIResponseGet(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, version := range alertmanager.SupportedVersions {
		httpmock.Reset()
		mock.RegisterURL("api/v1/status", version, "status")
		mock.RegisterURL("api/v1/alerts/groups", version, "alerts/groups")

		v := alertmanager.GetVersion()
		vs := fmt.Sprintf("%d.%d", v.Major, v.Minor)
		if version != vs {
			t.Errorf("GetVersion() returned '%s', expected '%s'", vs, version)
		}

		groups, err := alertmanager.GetAlerts(&v)
		if err != nil {
			t.Errorf("GetAlerts() failed: %s", err.Error())
		}
		if len(groups) != 4 {
			t.Errorf("Got %d groups, expected %d", len(groups), 4)
		}
	}
}
