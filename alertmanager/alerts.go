package alertmanager

import (
	"errors"

	"github.com/cloudflare/unsee/config"
	"github.com/cloudflare/unsee/models"
	"github.com/gavv/monotime"

	log "github.com/Sirupsen/logrus"
)

// AlertGroupsAPIResponse is the schema of API response for /api/v1/alerts/groups
type AlertGroupsAPIResponse struct {
	Status    string                          `json:"status"`
	Groups    []models.AlertmanagerAlertGroup `json:"data"`
	ErrorType string                          `json:"errorType"`
	Error     string                          `json:"error"`
}

// Get response from Alertmanager /api/v1/alerts/groups
func (response *AlertGroupsAPIResponse) Get() error {
	start := monotime.Now()

	url, err := joinURL(config.Config.AlertmanagerURI, "api/v1/alerts/groups")
	if err != nil {
		return err
	}

	err = getJSONFromURL(url, config.Config.AlertmanagerTimeout, response)
	if err != nil {
		return err
	}

	if response.Status != "success" {
		return errors.New(response.Error)
	}

	log.Infof("Got %d alert group(s) in %s", len(response.Groups), monotime.Since(start))
	return nil
}
