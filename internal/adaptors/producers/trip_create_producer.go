package producers

import (
	"context"
	"fmt"
	"golang-skeleton/internal/domain"
	"golang-skeleton/internal/domain/adaptors"
	"golang-skeleton/internal/domain/adaptors/encoders"
	domprod "golang-skeleton/internal/domain/adaptors/producers"
	"golang-skeleton/internal/domain/application"
	"golang-skeleton/internal/domain/events"
	"golang-skeleton/pkg/errors"
	"time"

	"github.com/recodextech/container"

	"github.com/gmbyapa/kstream/v2/kafka"
)

type DemandTripProducer struct {
	producer kafka.Producer
	logger   adaptors.Logger
	encoders struct {
		key, value encoders.Encoder
	}
}

// Init initializes the DemandTripProducer.
func (d *DemandTripProducer) Init(container container.Container) error {
	d.producer = container.Resolve(application.ModuleDefaultProducer).(domprod.DefaultProducer).KafkaProducer()
	d.logger = container.Resolve(application.ModuleLogger).(adaptors.Logger).
		NewLog(adaptors.LoggerPrefixed(`producers.trip`))
	d.encoders.key = container.Resolve(application.ModuleStringEncoder).(encoders.Encoder)
	d.encoders.value = container.Resolve(application.ModuleTripJSONEncoder).(encoders.Encoder)

	return nil
}

// Produce publish streams to the demand trips
func (d *DemandTripProducer) Produce(ctx context.Context, trip events.TripCreate) error {

	key, err := d.encoders.key.Encode(trip.Payload.PassengerID)
	if err != nil {
		return ProducerError{errors.Wrap(err, ErrMsgKeyEncode), errors.CodeKeyEncode,
			ErrMsgKeyEncode}
	}

	value, err := d.encoders.value.Encode(trip)
	if err != nil {
		return ProducerError{errors.Wrap(err, ErrMsgValEncode), errors.CodeValEncode, ErrMsgValEncode}
	}

	headers := kafka.RecordHeaders{
		kafka.RecordHeader{
			Key:   []byte(events.EventHeaderAccountID),
			Value: []byte(trip.GetMeta().AccountID.String()),
		},
		kafka.RecordHeader{
			Key:   []byte(events.EventHeaderTraceID),
			Value: []byte(trip.GetMeta().TraceID),
		},
	}
	record := d.producer.NewRecord(ctx, key, value, domain.TopicTrip, 0, time.Now(), headers, ``)

	partition, offset, err := d.producer.ProduceSync(ctx, record)
	if err != nil {
		return ProducerError{errors.Wrap(err, ErrMsgResourceWrite), errors.CodeResWFailed, ErrMsgResourceWrite}
	}

	d.logger.TraceContext(ctx, fmt.Sprintf(`Trip create failed for passenger id %s and budget type %s`,
		trip.Payload.PassengerID, trip.Payload.PassengerID),
		d.logger.Params(`topic`, record.Topic()),
		d.logger.Params(`partition`, fmt.Sprintf("%v", partition)),
		d.logger.Params(`offset`, fmt.Sprintf("%v", offset)))

	return nil
}
