package alertmanager

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/cloudflare/unsee/config"
	"github.com/cloudflare/unsee/models"
)

// Alertmanager 0.4 silence format
type silenceV04 struct {
	ID       int `json:"id"`
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

// silenceAPIResponseV04 is what Alertmanager 0.4 API returns
type silenceAPIResponseV04 struct {
	Status string `json:"status"`
	Data   struct {
		Silences      []silenceV04 `json:"silences"`
		TotalSilences int          `json:"totalSilences"`
	} `json:"data"`
	Error string `json:"error"`
}

// Get will make a request to Alertmanager API and parse the response
// It will only return silences or error (if any)
func (resp *silenceAPIResponseV04) Get(url string) ([]models.Silence, error) {
	silences := []models.Silence{}
	// Alertmanager 0.4 uses pagination, request max number of silences
	url = fmt.Sprintf("%s?limit=%d", url, math.MaxUint32)
	err := getJSONFromURL(url, config.Config.AlertmanagerTimeout, &resp)
	if err != nil {
		return silences, err
	}
	if resp.Status != StatusOK {
		return silences, errors.New(resp.Error)
	}
	for _, s := range resp.Data.Silences {
		us := models.Silence{
			ID:        string(s.ID),
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
