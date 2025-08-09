package http

import (
	"golang-skeleton/internal/domain/application"
	"golang-skeleton/internal/http/handlers"
	"golang-skeleton/internal/http/responses/writers"

	"github.com/recodextech/container"
)

type HTTP struct{}

// Init initializes the http module.
func (h *HTTP) Init(c container.Container) error {
	// Http validators

	// response writers
	c.Bind(writers.ModuleTripResponseWriter, new(writers.TripResponseWriter))

	// Http handlers
	c.Bind(handlers.ModuleCreateTrip, new(handlers.TripCreateHandler))

	// Http error handler
	c.Bind(handlers.ModuleErrorHandler, new(handlers.ErrorHandler))

	c.Bind(application.ModuleHTTPRouter, new(Router))
	c.Bind(application.ModuleHTTPServer, new(Server))

	c.Init(
		// Http request validators

		// response writers
		writers.ModuleTripResponseWriter,

		// Http handlers
		handlers.ModuleCreateTrip,

		// Http error handler
		handlers.ModuleErrorHandler,

		// Http Server Init
		application.ModuleHTTPRouter,
		application.ModuleHTTPServer,
	)

	return nil
}
