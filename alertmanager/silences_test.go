package alertmanager_test

import (
	"fmt"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/cloudflare/unsee/alertmanager"
	"github.com/cloudflare/unsee/mock"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestSilenceAPIResponseGet(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, version := range alertmanager.SupportedVersions {
		httpmock.Reset()
		mock.RegisterURL("api/v1/status", version, "status")
		mock.RegisterURL("api/v1/silences", version, "silences")

		v := alertmanager.GetVersion()
		vs := fmt.Sprintf("%d.%d", v.Major, v.Minor)
		if version != vs {
			t.Errorf("GetVersion() returned '%s', expected '%s'", vs, version)
		}

		silences, err := alertmanager.GetSilences(&v)
		if err != nil {
			t.Errorf("GetSilences() failed: %s", err.Error())
		}
		if len(silences) != 2 {
			t.Errorf("Got %d silences, expected %d", len(silences), 2)
		}
	}
}
