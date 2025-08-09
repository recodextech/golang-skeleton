package bootstrap

import (
	"golang-skeleton/internal/adaptors/encoders/avro"
	"golang-skeleton/internal/adaptors/producers"
	"golang-skeleton/internal/domain/application"
	"golang-skeleton/internal/http"
	"golang-skeleton/internal/streaming"
	"golang-skeleton/pkg/log"
	"golang-skeleton/pkg/metrics"

	"github.com/recodextech/container"
)

func Boot() {
	con := container.NewContainer()

	// Application config bindings
	err := con.SetModuleGlobalConfig(
		container.ModuleConfig{Key: application.ModuleBaseLogger, Value: new(log.LoggerConf)},
		// container.ModuleConfig{Key: application.ModuleApp, Value: new(app.Config)},
		container.ModuleConfig{Key: application.ModuleStreams, Value: new(streaming.Config)},
		container.ModuleConfig{Key: application.ModuleDefaultProducer, Value: new(producers.Config)},
		container.ModuleConfig{Key: application.ModuleMetricsReporter, Value: new(metrics.Config)},
		container.ModuleConfig{Key: application.ModuleSchemaRegistry, Value: new(avro.SchemaRegistryConfig)},
		container.ModuleConfig{Key: application.ModuleHTTPRouter, Value: new(http.KRouterConf)},
		container.ModuleConfig{Key: application.ModuleHTTPServer, Value: new(http.Conf)},
	)
	if err != nil {
		panic(err)
	}

	bind(con)
	initModules(con)
	start(con)
}
