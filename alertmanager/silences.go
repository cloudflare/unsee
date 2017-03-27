package alertmanager

import (
	"errors"
	"fmt"
	"math"

	"github.com/cloudflare/unsee/config"
	"github.com/cloudflare/unsee/models"
	"github.com/gavv/monotime"

	log "github.com/Sirupsen/logrus"
)

type silencesData struct {
	Silences      []models.AlertmanagerSilence `json:"silences"`
	TotalSilences int                          `json:"totalSilences"`
}

// SilenceAPIResponse is what Alertmanager API returns
type SilenceAPIResponse struct {
	Status    string       `json:"status"`
	Data      silencesData `json:"data"`
	ErrorType string       `json:"errorType"`
	Error     string       `json:"error"`
}

// Get will return fresh data from Alertmanager API
func (response *SilenceAPIResponse) Get() error {
	start := monotime.Now()

	url, err := joinURL(config.Config.AlertmanagerURI, "api/v1/silences")
	if err != nil {
		return err
	}
	url = fmt.Sprintf("%s?limit=%d", url, math.MaxUint32)

	err = getJSONFromURL(url, config.Config.AlertmanagerTimeout, response)
	if err != nil {
		return err
	}

	if response.Status != "success" {
		return errors.New(response.Error)
	}

	log.Infof("Got %d silences(s) in %s", len(response.Data.Silences), monotime.Since(start))
	return nil
}
