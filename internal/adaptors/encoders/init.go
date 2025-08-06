package encoders

import (
	"golang-skeleton/internal/adaptors/encoders/avro"
	"golang-skeleton/internal/adaptors/encoders/json"
	"golang-skeleton/internal/adaptors/encoders/misc"
	"golang-skeleton/internal/domain/application"
	"golang-skeleton/pkg/container"
)

type Encoders struct{}

func (e *Encoders) Init(c container.Container) error {
	// bind encoder modules
	c.Bind(application.ModuleSchemaRegistry, new(avro.SchemaRegistry))
	c.Bind(application.ModuleStringEncoder, new(misc.StringEncoder))
	c.Bind(application.ModuleUUIDEncoder, new(misc.UUIDEncoder))
	c.Bind(application.ModuleJSONEncoder, new(json.Encoder))
	c.Bind(application.ModuleTripJSONEncoder, new(json.TripCreateJSONEncoder))

	// init schema registry
	// c.Init(application.ModuleSchemaRegistry)

	return nil
}
