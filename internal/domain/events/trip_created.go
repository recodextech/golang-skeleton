package events

import "golang-skeleton/internal/domain/models"

type TripCreate struct {
	Meta    Meta `json:"meta"`
	Payload struct {
		PassengerID string        `json:"passengerID"`
		Jobs        []models.Trip `json:"jobs"`
	}
}

func (t *TripCreate) GetMeta() Meta {
	return t.Meta
}
