package uuid

import (
	"golang-skeleton/pkg/errors"

	"github.com/google/uuid"
)

type UUID struct {
	uuid uuid.UUID
}

var Nil = UUID{uuid: uuid.Nil}

func New() UUID {
	return UUID{uuid: uuid.New()}
}

func Parse(s string) (UUID, error) {
	var id uuid.UUID
	if s == "" {
		id = uuid.Nil

		return UUID{uuid: id}, nil
	}

	i, err := uuid.Parse(s)
	if err != nil {
		return UUID{}, ParseError{errors.New(ErrMsgInvalidUUID), ErrMsgInvalidUUID}
	}

	return UUID{uuid: i}, nil
}

func (u UUID) String() string {
	return u.uuid.String()
}

func (u UUID) UUID() uuid.UUID {
	return u.uuid
}

func (u *UUID) UnmarshalJSON(b []byte) error {
	if string(b) == `""` || string(b) == "null" {
		u.uuid = uuid.Nil

		return nil
	}

	u.uuid = uuid.UUID{}

	return u.uuid.UnmarshalText(b)
}

func (u UUID) MarshalJSON() ([]byte, error) {
	if u.uuid == uuid.Nil {
		return []byte(`""`), nil
	}

	return []byte(`"` + u.uuid.String() + `"`), nil
}
