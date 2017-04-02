package models

import "time"

// AlertmanagerAlert is vanilla alert object from Alertmanager
type AlertmanagerAlert struct {
	Annotations  map[string]string `json:"annotations"`
	Labels       map[string]string `json:"labels"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Inhibited    bool              `json:"inhibited"`
	Silenced     string            `json:"silenced"`
}

// AlertmanagerAlertGroup is vanilla group object from Alertmanager, exposed under api/v1/alerts/groups
type AlertmanagerAlertGroup struct {
	Labels map[string]string `json:"labels"`
	Blocks []struct {
		Alerts []AlertmanagerAlert `json:"alerts"`
	} `json:"blocks"`
}

// AlertmanagerSilence is vanilla silence object from Alertmanager, exposed under api/v1/silences
type AlertmanagerSilence struct {
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
