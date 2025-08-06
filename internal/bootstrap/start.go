package bootstrap

import (
	"golang-skeleton/internal/domain/application"
	"golang-skeleton/pkg/container"
	"os"
	"os/signal"
)

func start(con container.AppContainer) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		con.Shutdown(
			// application.ModuleReadyIndicator,
			application.ModuleHTTPServer,
			application.ModuleHTTPRouter,
			application.ModuleStreams,
			application.ModuleMetricsReporter,
		)
	}()

	con.Start(
		application.ModuleMetricsReporter,
		application.ModuleStreams,
		application.ModuleHTTPRouter,
		application.ModuleHTTPServer,
		// application.ModuleReadyIndicator,
	)
}
