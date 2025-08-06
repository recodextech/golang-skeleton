package services

import (
	"context"
	"golang-skeleton/internal/domain/adaptors/producers"
	"golang-skeleton/internal/domain/adaptors/repositories"
	"golang-skeleton/internal/domain/application"
	"golang-skeleton/internal/domain/events"
	"golang-skeleton/internal/http/request"
	"golang-skeleton/pkg/container"
	"golang-skeleton/pkg/errors"
)

type TripCreatedService struct {
	tripProducer producers.TripProducer
	tripRepo     repositories.TripRepository
}

func (t *TripCreatedService) Init(c container.Container) error {
	t.tripProducer = c.Resolve(application.ModuleTripProducer).(producers.TripProducer)
	t.tripRepo = c.Resolve(application.ModuleTripRepo).(repositories.TripRepository)
	return nil
}

func (t *TripCreatedService) CreateTrip(ctx context.Context, trip request.TripCreate) (events.TripCreate, error) {
	tripEvent := events.TripCreate{}
	exist, err := t.tripRepo.Exists(ctx, trip.PassengerID)
	if err != nil {
		return tripEvent, errors.Wrap(err, `error service layer`)
	}
	if exist {
		return tripEvent, errors.New(`passenger already assigned : ` + trip.PassengerID)
	}

	tripEvent.Meta = events.NewMetaContext(ctx)
	tripEvent.Payload.PassengerID = trip.PassengerID
	tripEvent.Payload.Jobs = trip.Jobs
	err = t.tripProducer.Produce(ctx, tripEvent)

	return tripEvent, err
}
