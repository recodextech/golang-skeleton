package bootstrap

import (
	"golang-skeleton/internal/adaptors/encoders"
	"golang-skeleton/internal/adaptors/producers"
	"golang-skeleton/internal/adaptors/repositories"
	"golang-skeleton/internal/domain/application"
	"golang-skeleton/internal/http"
	"golang-skeleton/internal/services"
	"golang-skeleton/internal/streaming"
	"golang-skeleton/internal/usecases"
	"golang-skeleton/pkg/container"
	log2 "golang-skeleton/pkg/log"
	"golang-skeleton/pkg/metrics"
)

func bind(c container.Container) {
	c.Bind(application.ModuleBaseLogger, new(log2.BaseLogger))
	c.Bind(application.ModuleLogger, new(log2.Logger))
	c.Bind(application.ModuleMetricsReporter, new(metrics.Reporter))

	// streams
	c.Bind(application.ModuleEncoders, new(encoders.Encoders))
	c.Bind(application.ModuleStreams, new(streaming.Streams))

	//Producers
	c.Bind(application.ModuleDefaultProducer, new(producers.DefaultProducer))
	c.Bind(application.ModuleTripProducer, new(producers.DemandTripProducer))

	// Repositories
	c.Bind(application.ModuleTripRepo, new(repositories.TripRepository))

	// Adapters

	// UseCases
	c.Bind(application.ModuleTripCreateUsecase, new(usecases.TripCreateUsecase))

	// Services
	c.Bind(application.ModuleTripCreateService, new(services.TripCreatedService))

	// http
	c.Bind(application.ModuleHTTP, new(http.HTTP))
}
