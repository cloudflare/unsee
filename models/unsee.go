package models

// UnseeSilence is vanilla silence + some additional attributes
// Unsee adds JIRA support, it can extract JIRA IDs from comments
// extracted ID is used to generate link to JIRA issue
// this means Unsee needs to store additional fields for each silence
type UnseeSilence struct {
	AlertmanagerSilence
	JiraID  string `json:"jiraID"`
	JiraURL string `json:"jiraURL"`
}

// UnseeAlert is vanilla alert + some additional attributes
// unsee extends an alert object with:
// * Links map, it's generated from annotations if annotation value is an url
//   it's pulled out of annotation map and returned under links field,
//   unsee UI used this to show links differently than other annotations
// * Fingerprint, which is a sha1 of the entire alert
type UnseeAlert struct {
	AlertmanagerAlert
	Links       map[string]string `json:"links"`
	Fingerprint string            `json:"-"`
}

// UnseeAlertList is flat list of UnseeAlert objects
type UnseeAlertList []UnseeAlert

func (a UnseeAlertList) Len() int {
	return len(a)

}
func (a UnseeAlertList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a UnseeAlertList) Less(i, j int) bool {
	// compare timestamps, if equal compare fingerprints to stable sort order
	if a[i].StartsAt.After(a[j].StartsAt) {
		return true
	}
	if a[i].StartsAt.Before(a[j].StartsAt) {
		return false
	}
	return a[i].Fingerprint < a[j].Fingerprint
}

// UnseeAlertGroup is vanilla Alertmanager group, but alerts are flattened
// There is a hash computed from all alerts, it's used by UI to quickly tell
// if there was any change in a group and it needs to refresh it
type UnseeAlertGroup struct {
	Labels          map[string]string `json:"labels"`
	Alerts          UnseeAlertList    `json:"alerts"`
	ID              string            `json:"id"`
	Hash            string            `json:"hash"`
	SilencedCount   int               `json:"silencedCount"`
	UnsilencedCount int               `json:"unsilencedCount"`
}

// UnseeFilter holds returned data on any filter passed by the user as part of the query
type UnseeFilter struct {
	Text    string `json:"text"`
	Hits    int    `json:"hits"`
	IsValid bool   `json:"isValid"`
}

// UnseeColor is used by UnseeLabelColor to reprenset colors as RGBA
type UnseeColor struct {
	Red   uint8 `json:"red"`
	Green uint8 `json:"green"`
	Blue  uint8 `json:"blue"`
	Alpha uint8 `json:"alpha"`
}

// UnseeLabelColor holds color information for labels that should be colored in the UI
// every configured label will have a distinct coloring for each value
type UnseeLabelColor struct {
	Font       UnseeColor `json:"font"`
	Background UnseeColor `json:"background"`
}

// UnseeColorMap is a map of "Label Key" -> "Label Value" -> UnseeLabelColor
type UnseeColorMap map[string]map[string]UnseeLabelColor

// UnseeCountMap is a map of "Label Key" -> "Label Value" -> number of occurence
type UnseeCountMap map[string]map[string]int

// UnseeAlertsResponse is the structure of JSON response UI will use to get alert data
type UnseeAlertsResponse struct {
	Status      string                  `json:"status"`
	Error       string                  `json:"error,omitempty"`
	Timestamp   string                  `json:"timestamp"`
	Version     string                  `json:"version"`
	AlertGroups []UnseeAlertGroup       `json:"groups"`
	Silences    map[string]UnseeSilence `json:"silences"`
	Colors      UnseeColorMap           `json:"colors"`
	Filters     []UnseeFilter           `json:"filters"`
	Counters    UnseeCountMap           `json:"counters"`
}

// UnseeAutocomplete is the structure of autocomplete object for filter hints
// this is internal represenation, not what's returned to the user
type UnseeAutocomplete struct {
	Value  string   `json:"value"`
	Tokens []string `json:"tokens"`
}
