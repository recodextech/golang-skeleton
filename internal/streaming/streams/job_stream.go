package streams

import (
	"context"
	"golang-skeleton/internal/domain"
	"golang-skeleton/internal/domain/adaptors"
	"golang-skeleton/internal/domain/adaptors/encoders"
	"golang-skeleton/internal/domain/application"
	"golang-skeleton/pkg/container"

	kstream "github.com/gmbyapa/kstream/v2/streams"
)

type JobStream struct {
	kstream.Stream
	log adaptors.Logger
}

// Init to trip stream
func (tvs *JobStream) Init(c container.Container) error {
	tvs.log = c.Resolve(application.ModuleLogger).(adaptors.Logger).
		NewLog(adaptors.LoggerPrefixed("streams.job-stream"))

	tvs.Stream = c.Resolve(application.ModuleStreamBuilder).(*kstream.StreamBuilder).KStream(
		domain.TopicTrip,
		c.Resolve(application.ModuleStringEncoder).(encoders.Encoder),
		c.Resolve(application.ModuleTripJSONEncoder).(encoders.Encoder),
	)

	tvs.Stream.Filter(func(ctx context.Context, key, value interface{}) (bool, error) {
		return value != nil, nil
	}).To(domain.TopicJob,
		kstream.ProduceWithKeyEncoder(c.Resolve(application.ModuleStringEncoder).(encoders.Encoder)),
		kstream.ProduceWithValEncoder(c.Resolve(application.ModuleTripJSONEncoder).(encoders.Encoder)))

	return nil
}
