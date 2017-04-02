package alertmanager

import (
	"errors"
	"time"

	"github.com/cloudflare/unsee/config"
	"github.com/cloudflare/unsee/models"
)

// Alertmanager 0.5 silence format
type silence05 struct {
	ID       string `json:"id"`
	Matchers []struct {
		Name    string `json:"name"`
		Value   string `json:"value"`
		IsRegex bool   `json:"isRegex"`
	} `json:"matchers"`
	StartsAt  time.Time `json:"startsAt"`
	EndsAt    time.Time `json:"endsAt"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
	Comment   string    `json:"comment"`
}

// SilenceAPIResponseV05 is what Alertmanager 0.5 API returns
type silenceAPIResponseV05 struct {
	Status string      `json:"status"`
	Data   []silence05 `json:"data"`
	Error  string      `json:"error"`
}

// Get will make a request to Alertmanager API and parse the response
// It will only return silences or error (if any)
func (resp *silenceAPIResponseV05) Get(url string) ([]models.Silence, error) {
	silences := []models.Silence{}
	err := getJSONFromURL(url, config.Config.AlertmanagerTimeout, &resp)
	if err != nil {
		return silences, err
	}
	if resp.Status != StatusOK {
		return silences, errors.New(resp.Error)
	}
	for _, s := range resp.Data {
		us := models.Silence{
			ID:        s.ID,
			Matchers:  s.Matchers,
			StartsAt:  s.StartsAt,
			EndsAt:    s.EndsAt,
			CreatedAt: s.CreatedAt,
			CreatedBy: s.CreatedBy,
			Comment:   s.Comment,
		}
		silences = append(silences, us)
	}
	return silences, nil
}
