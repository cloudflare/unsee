package alertmanager

import (
	"github.com/blang/semver"
	"github.com/cloudflare/unsee/config"

	log "github.com/Sirupsen/logrus"
)

// StatusOK is the string used in successful responses
var StatusOK = "success"

// SupportedVersions is the list of versions we support
var SupportedVersions = []string{"0.4", "0.5"}

// AlertmanagerVersion is what api/v1/status returns, we only use it to check
// version, so we skip all other keys (except for status)
type alertmanagerVersion struct {
	Status string `json:"status"`
	Data   struct {
		VersionInfo struct {
			Version string `json:"version"`
		} `json:"versionInfo"`
	} `json:"data"`
}

// GetVersion returns version information of the remote Alertmanager endpoint
func GetVersion() semver.Version {
	// if everything fails assume Alertmanager is at latest possible version
	defaultVersion, _ := semver.Make("999.0.0")

	url, err := joinURL(config.Config.AlertmanagerURI, "api/v1/status")
	if err != nil {
		log.Errorf("Failed to join url '%s' and path 'api/v1/status': %s", config.Config.AlertmanagerURI, err.Error())
		return defaultVersion
	}
	ver := alertmanagerVersion{}
	err = getJSONFromURL(url, config.Config.AlertmanagerTimeout, &ver)
	if err != nil {
		log.Errorf("%s request failed: %s", url, err.Error())
		return defaultVersion
	}

	if ver.Status != StatusOK {
		log.Errorf("Request to %s returned status %s", url, ver.Status)
		return defaultVersion
	}

	if ver.Data.VersionInfo.Version == "" {
		log.Error("No version information in Alertmanager API")
		return defaultVersion
	}

	v, err := semver.Make(ver.Data.VersionInfo.Version)
	if err != nil {
		log.Warningf("Alertmanager version string ('%s') parsing failed: %s", ver.Data.VersionInfo.Version, err.Error())
		return defaultVersion
	}
	log.Infof("Remote Alertmanager version: %v", v)
	return v
}
