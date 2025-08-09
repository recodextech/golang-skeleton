package tables

import (
	"golang-skeleton/internal/domain"
	"golang-skeleton/internal/domain/adaptors/encoders"
	"golang-skeleton/internal/domain/application"
	"golang-skeleton/internal/streaming/stores"

	"github.com/recodextech/container"

	kstream "github.com/gmbyapa/kstream/v2/streams"
)

type TripTable struct {
	kstream.GlobalTable
}

func (t *TripTable) Init(c container.Container) error {
	t.GlobalTable = c.Resolve(application.ModuleStreamBuilder).(*kstream.StreamBuilder).GlobalTable(
		domain.TopicTrip,
		c.Resolve(application.ModuleStringEncoder).(encoders.Encoder),
		c.Resolve(application.ModuleTripJSONEncoder).(encoders.Encoder),
		stores.StoreNameTrip)
	return nil
}
