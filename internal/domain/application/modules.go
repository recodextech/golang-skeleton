package application

// Module bindings
const (
	ModuleApp             string = "modules.application"
	ModuleLogger          string = "modules.logger"
	ModuleBaseLogger      string = "modules.base-logger"
	ModuleSchemaRegistry  string = "modules.schema-registry"
	ModuleStreams         string = "modules.streams"
	ModuleStreamBuilder   string = "modules.abstractions-builder"
	ModuleStoreRegistry   string = "modules.store-registry"
	ModuleMetricsReporter string = "modules.metrics-reporter"
	ModuleStreamReporter  string = "modules.base-metrics-reporter"
	ModuleSQL             string = "module-sql"

	ModuleHTTP           string = "modules.http"
	ModuleHTTPRouter     string = "modules.router"
	ModuleHTTPServer     string = "modules.http.server"
	ModuleReadyIndicator string = "module-ready-indicator"
	ModulePprofIndicator string = "module-pprof-indicator"
	ModuleErrorHandler   string = "http-error-handler"
)

// Encoder
const (
	ModuleEncoders        string = "encoders"
	ModuleJSONEncoder     string = "encoders.json"
	ModuleTripJSONEncoder string = "encoders.trip.json"
	ModuleStringEncoder   string = "encoders.string"
	ModuleUUIDEncoder     string = "encoders.uuid"
)

// Repositories
const (
	ModuleTripRepo string = "repositories.trip"
)

// Stores
const (
	ModuleTripStore string = "stores.trip"
)

// Global Tables
const (
	ModuletripGlobalTable string = `streams.gtables.trip`
)

// ModuleTimelineVersionStream Streams
const (
	ModuleJobStream string = `streams.trip`
)

// Use cases
const (
	ModuleTripCreateUsecase string = "useCases.trip-create"
)

// Services
const (
	ModuleTripCreateService string = "services.trip-create"
)

// Producer
const (
	ModuleDefaultProducer string = `producers.default-producer`
	ModuleTripProducer    string = `producers.trip-producer`
)
