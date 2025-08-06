package request

import (
	"encoding/json"
	"golang-skeleton/internal/domain/models"
)

type TripCreate struct {
	PassengerID string        `json:"passengerId"`
	Jobs        []models.Trip `json:"jobs"`
}

func (t TripCreate) Encode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (t TripCreate) Decode(data []byte) (interface{}, error) {
	req := TripCreate{}
	err := json.Unmarshal(data, &req)

	return req, err
}
