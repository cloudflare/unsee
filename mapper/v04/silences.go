package amv04

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/blang/semver"
	"github.com/cloudflare/unsee/config"
	"github.com/cloudflare/unsee/mapper"
	"github.com/cloudflare/unsee/models"
	"github.com/cloudflare/unsee/remote"
)

// Alertmanager 0.4 silence format
type silence struct {
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
type silenceAPISchema struct {
	Status string `json:"status"`
	Data   struct {
		Silences      []silence `json:"silences"`
		TotalSilences int       `json:"totalSilences"`
	} `json:"data"`
	Error string `json:"error"`
}

// V04SilenceMapper implements Alertmanager 0.4 API schema
type V04SilenceMapper struct {
	mapper.SilenceMapper
}

// IsSupported returns true if given version string is supported
func (m V04SilenceMapper) IsSupported(version string) bool {
	versionRange := semver.MustParseRange(">=0.4.0 <0.5.0")
	return versionRange(semver.MustParse(version))
}

// GetSilences will make a request to Alertmanager API and parse the response
// It will only return silences or error (if any)
func (m V04SilenceMapper) GetSilences() ([]models.Silence, error) {
	silences := []models.Silence{}
	resp := silenceAPISchema{}

	url, err := remote.JoinURL(config.Config.AlertmanagerURI, "api/v1/silences")
	if err != nil {
		return silences, err
	}

	// Alertmanager 0.4 uses pagination for silences
	url = fmt.Sprintf("%s?limit=%d", url, math.MaxUint32)
	err = remote.GetJSONFromURL(url, config.Config.AlertmanagerTimeout, &resp)
	if err != nil {
		return silences, err
	}

	if resp.Status != "success" {
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
