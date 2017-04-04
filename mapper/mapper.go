package mapper

import (
	"fmt"

	"github.com/cloudflare/unsee/models"
)

var (
	alertMappers   = []AlertMapper{}
	silenceMappers = []SilenceMapper{}
)

func init() {

}

// AlertMapper impement Alertmanager -> unsee alert data mapping that works
// for a specific range of Alertmanager versions
type AlertMapper interface {
	IsSupported(version string) bool
	GetAlerts() ([]models.AlertGroup, error)
}

// RegisterAlertMapper allows to register mapper implementing alert data
// handling for specific Alertmanager versions
func RegisterAlertMapper(m AlertMapper) {
	alertMappers = append(alertMappers, m)
}

// GetAlertMapper returns mapper for given version
func GetAlertMapper(version string) (AlertMapper, error) {
	for _, m := range alertMappers {
		if m.IsSupported(version) {
			return m, nil
		}
	}
	return nil, fmt.Errorf("Can't find alert mapper for Alertmanager %s", version)
}

// SilenceMapper impement Alertmanager -> unsee silence data mapping that works
// for a specific range of Alertmanager versions
type SilenceMapper interface {
	Release() string
	IsSupported(version string) bool
	GetSilences() ([]models.Silence, error)
}

// RegisterSilenceMapper allows to register mapper implementing silence data
// handling for specific Alertmanager versions
func RegisterSilenceMapper(m SilenceMapper) {
	silenceMappers = append(silenceMappers, m)
}

// GetSilenceMapper returns mapper for given version
func GetSilenceMapper(version string) (SilenceMapper, error) {
	for _, m := range silenceMappers {
		if m.IsSupported(version) {
			return m, nil
		}
	}
	return nil, fmt.Errorf("Can't find silence mapper for Alertmanager %s", version)
}
