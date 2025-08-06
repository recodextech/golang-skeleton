package json

import (
	"encoding/json"
	"golang-skeleton/internal/domain/events"
	"golang-skeleton/pkg/errors"
)

type TripCreateJSONEncoder struct {
}

func (t TripCreateJSONEncoder) Encode(data interface{}) ([]byte, error) {
	byt, err := json.Marshal(data)
	if err != nil {
		return nil, ValueEncodingError{errors.Wrap(err, failedToEncodeValue)}
	}

	return byt, nil
}

func (t TripCreateJSONEncoder) Decode(data []byte) (interface{}, error) {
	py := events.TripCreate{}
	err := json.Unmarshal(data, &py)
	if err != nil {
		return nil, ValueEncodingError{errors.Wrap(err, failedToDecodeValue)}
	}

	return py, nil
}
