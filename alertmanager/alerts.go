package alertmanager

import (
	"time"

	"github.com/blang/semver"
	"github.com/cloudflare/unsee/config"
	"github.com/cloudflare/unsee/models"

	log "github.com/Sirupsen/logrus"
)

// GetAlerts will send request to Alertmanager and return list of alert groups
// from the API
func GetAlerts(v *semver.Version) ([]models.AlertGroup, error) {
	groups := []models.AlertGroup{}
	start := time.Now()
	url, err := joinURL(config.Config.AlertmanagerURI, "api/v1/alerts/groups")
	if err != nil {
		return groups, err
	}

	v05, _ := semver.Make("0.5.0")
	if v.GE(v05) {
		response := alertsGroupsAPIResponseV05{}
		groups, err = response.Get(url)
	} else {
		response := alertsGroupsAPIResponseV04{}
		groups, err = response.Get(url)
	}
	if err != nil {
		return groups, err
	}
	log.Infof("Got %d alert group(s) in %s", len(groups), time.Since(start))
	return groups, nil
}
