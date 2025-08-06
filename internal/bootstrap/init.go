package bootstrap

import (
	"golang-skeleton/internal/domain/application"
	"golang-skeleton/pkg/container"
)

func initModules(c container.Container) {
	c.Init(
		// main application dependencies
		application.ModuleBaseLogger,
		application.ModuleLogger,
		application.ModuleMetricsReporter,
		application.ModuleReadyIndicator,
		application.ModulePprofIndicator,

		//encoders
		application.ModuleEncoders,

		// Producers
		application.ModuleDefaultProducer,
		application.ModuleTripProducer,

		// Streams
		application.ModuleStreams,

		// Repositories
		application.ModuleTripRepo,

		// Adapters

		// UseCases
		application.ModuleTripCreateUsecase,

		// Services
		application.ModuleTripCreateService,

		// Http
		application.ModuleHTTP,
	)
}
