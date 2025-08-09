package repositories

import (
	"context"
	"golang-skeleton/internal/domain/adaptors"
	"golang-skeleton/internal/domain/application"

	"github.com/recodextech/container"

	store "github.com/gmbyapa/kstream/v2/streams/stores"
)

type TripRepository struct {
	log       adaptors.Logger
	tripStore store.Store
}

func (repo *TripRepository) Init(c container.Container) error {
	repo.log = c.Resolve(application.ModuleLogger).(adaptors.Logger).
		NewLog(adaptors.LoggerPrefixed("repositories.trip-store"))
	repo.tripStore = c.Resolve(application.ModuleTripStore).(store.Store)

	return nil
}

func (repo *TripRepository) Exists(ctx context.Context, trip string) (exists bool, err error) {
	v, err := repo.tripStore.Get(ctx, trip)
	if err != nil {
		return exists, err
	}

	return v != nil, nil
}
