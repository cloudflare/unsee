package alertmanager

import (
	"time"

	"github.com/blang/semver"
	"github.com/cloudflare/unsee/config"
	"github.com/cloudflare/unsee/models"

	log "github.com/Sirupsen/logrus"
)

// GetSilences will send request to Alertmanager and return list of silences
// from the API
func GetSilences(v *semver.Version) ([]models.Silence, error) {
	silences := []models.Silence{}
	start := time.Now()
	url, err := joinURL(config.Config.AlertmanagerURI, "api/v1/silences")
	if err != nil {
		return silences, err
	}

	v05, _ := semver.Make("0.5.0")
	if v.GE(v05) {
		response := silenceAPIResponseV05{}
		silences, err = response.Get(url)
	} else {
		response := silenceAPIResponseV04{}
		silences, err = response.Get(url)
	}
	if err != nil {
		return silences, err
	}
	log.Infof("Got %d silences(s) in %s", len(silences), time.Since(start))
	return silences, nil
}
