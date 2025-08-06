package misc

import (
	"golang-skeleton/pkg/errors"
	"golang-skeleton/pkg/uuid"
)

type UUIDEncoder struct{}

func (u UUIDEncoder) Encode(data interface{}) ([]byte, error) {
	id, ok := data.(uuid.UUID)
	if !ok {
		return nil, errors.New("key encoder")
	}

	return []byte(id.String()), nil
}

func (u UUIDEncoder) Decode(data []byte) (interface{}, error) {
	if data == nil {
		return uuid.Nil, nil
	}
	idString := string(data)

	return uuid.Parse(idString)
}
