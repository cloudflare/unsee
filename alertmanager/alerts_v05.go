package alertmanager

import (
	"errors"
	"time"

	"github.com/cloudflare/unsee/config"
	"github.com/cloudflare/unsee/models"
)

// AlertmanagerAlert is vanilla alert object from Alertmanager 0.5
type alertV05 struct {
	Annotations  map[string]string `json:"annotations"`
	Labels       map[string]string `json:"labels"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Inhibited    bool              `json:"inhibited"`
	Silenced     string            `json:"silenced"`
}

// alertsGroupsV05 is vanilla group object from Alertmanager, exposed under api/v1/alerts/groups
type alertsGroupsV05 struct {
	Labels map[string]string `json:"labels"`
	Blocks []struct {
		Alerts []alertV05 `json:"alerts"`
	} `json:"blocks"`
}

// alertsGroupsAPIResponseV05 is the schema of API response for /api/v1/alerts/groups
type alertsGroupsAPIResponseV05 struct {
	Status string            `json:"status"`
	Groups []alertsGroupsV05 `json:"data"`
	Error  string            `json:"error"`
}

// Get will make a request to Alertmanager API and parse the response
// It will only return alerts or error (if any)
func (resp *alertsGroupsAPIResponseV05) Get(url string) ([]models.AlertGroup, error) {
	groups := []models.AlertGroup{}
	err := getJSONFromURL(url, config.Config.AlertmanagerTimeout, &resp)
	if err != nil {
		return groups, err
	}
	if resp.Status != StatusOK {
		return groups, errors.New(resp.Error)
	}
	for _, g := range resp.Groups {
		alertList := models.AlertList{}
		for _, b := range g.Blocks {
			for _, a := range b.Alerts {
				us := models.Alert{
					Annotations:  a.Annotations,
					Labels:       a.Labels,
					StartsAt:     a.StartsAt,
					EndsAt:       a.EndsAt,
					GeneratorURL: a.GeneratorURL,
					Inhibited:    a.Inhibited,
					Silenced:     a.Silenced,
				}
				alertList = append(alertList, us)
			}
		}
		ug := models.AlertGroup{
			Labels: g.Labels,
			Alerts: alertList,
		}
		groups = append(groups, ug)
	}
	return groups, nil
}
