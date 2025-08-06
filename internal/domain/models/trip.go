package models

type Trip struct {
	Pickup string   `json:"pickup"`
	Drops  []string `json:"drops"`
}
