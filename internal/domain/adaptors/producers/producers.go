package producers

import (
	"context"
	"github.com/gmbyapa/kstream/v2/kafka"
	"golang-skeleton/internal/domain/events"
)

type DefaultProducer interface {
	KafkaProducer() kafka.Producer
}

type TripProducer interface {
	Produce(ctx context.Context, trip events.TripCreate) error
}
