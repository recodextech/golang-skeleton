package handlers

import (
	"context"
	"golang-skeleton/internal/domain/application"
	"golang-skeleton/internal/domain/service"
	"golang-skeleton/internal/http/request"

	"github.com/recodextech/container"

	"github.com/tryfix/krouter"
)

const TripCreate = "trip.create"

type TripCreateHandler struct {
	tripService service.TripCreated
}

func (t *TripCreateHandler) Init(c container.Container) error {
	t.tripService = c.Resolve(application.ModuleTripCreateService).(service.TripCreated)
	return nil
}

func (t *TripCreateHandler) Handle(ctx context.Context, payload krouter.HttpPayload) (i interface{}, err error) {
	tr := payload.Body.(request.TripCreate)
	return t.tripService.CreateTrip(ctx, tr)
}

func (t *TripCreateHandler) PostHandler(ctx context.Context, payload krouter.Payload) error {
	return nil
}
