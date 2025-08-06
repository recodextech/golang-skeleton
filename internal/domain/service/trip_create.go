package service

import (
	"context"
	"golang-skeleton/internal/domain/events"
	"golang-skeleton/internal/http/request"
)

type TripCreated interface {
	CreateTrip(ctx context.Context, trip request.TripCreate) (events.TripCreate, error)
}
