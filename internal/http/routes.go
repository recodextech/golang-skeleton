package http

import (
	"golang-skeleton/internal/http/handlers"
	"golang-skeleton/internal/http/request"
	"golang-skeleton/internal/http/responses"
	"golang-skeleton/internal/http/responses/writers"
	"net/http"

	"github.com/recodextech/container"

	"github.com/tryfix/krouter"
)

func initRoutes(r *Router, c container.Container) {
	r.router.Handle("/create-trip",
		r.krouter.NewHandler(
			handlers.ModuleCreateTrip,
			request.TripCreate{},
			c.Resolve(handlers.ModuleCreateTrip).(*handlers.TripCreateHandler).Handle,
			c.Resolve(handlers.ModuleCreateTrip).(*handlers.TripCreateHandler).PostHandler,
			krouter.HandlerWithHeader(request.HeaderAccountID.String(), request.ParamTypeAppUUID, nil),
			krouter.HandlerWithHeader(request.HeaderUserID.String(), request.ParamTypeAppUUID, nil),
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithSuccessHandlerFunc(c.
				Resolve(writers.ModuleTripResponseWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodPost)
}
