package streaming

import (
	"context"
	"fmt"
	"golang-skeleton/internal/domain"
	"golang-skeleton/internal/domain/adaptors"
	"golang-skeleton/internal/domain/application"
	"golang-skeleton/internal/domain/events"
	"golang-skeleton/internal/streaming/stores"
	"golang-skeleton/internal/streaming/streams"
	"golang-skeleton/internal/streaming/tables"
	"golang-skeleton/pkg/container"
	"golang-skeleton/pkg/errors"
	metrics2 "golang-skeleton/pkg/metrics"
	"time"

	"github.com/gmbyapa/kstream/v2/kafka"
	"github.com/gmbyapa/kstream/v2/kafka/adaptors/librd"
	"github.com/gmbyapa/kstream/v2/kafka/adaptors/sarama"
	kstream "github.com/gmbyapa/kstream/v2/streams"
	"github.com/gmbyapa/kstream/v2/streams/topology"
	tlog "github.com/tryfix/log"
	"github.com/tryfix/metrics"
)

var visualize = func(topology string) {}

type StreamConfig struct {
	BootstrapServers []string `json:"bootstrap_servers"`
	Group            string   `json:"group"`
	Logger           adaptors.Logger
	Reporter         metrics.Reporter
}

type Streams struct {
	logger   adaptors.Logger
	builder  *kstream.StreamBuilder
	runner   kstream.Runner
	topology topology.Topology
	config   *Config
	admin    kafka.Admin
}

// todo mermur partition
// partitioner
var AccIDPartitioner = kstream.Partitioner(func(ctx context.Context, key, val interface{}, numPartitions int32) (int32, error) {
	return 0, nil
})

func (s *Streams) Init(c container.Container) (err error) {
	s.logger = c.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(
		adaptors.LoggerPrefixed("streaming"))
	s.config = c.GetGlobalConfig(application.ModuleStreams).(*Config)

	baseLogger := c.Resolve(application.ModuleBaseLogger).(tlog.Logger)
	c.Bind(`partitioner`, AccIDPartitioner)

	conf := librd.NewProducerConfig()
	conf.BootstrapServers = s.config.BootstrapServers
	conf.Id = s.config.ProducerID
	conf.Idempotent = true
	conf.Logger = baseLogger.NewLog(tlog.Prefixed(`DLQHandler`))
	producer, err := librd.NewProducer(conf)
	if err != nil {
		return err
	}

	// configure kStream
	kconf := kstream.NewStreamBuilderConfig()
	kconf.BootstrapServers = s.config.BootstrapServers
	kconf.ApplicationId = s.config.ApplicationID
	kconf.DefaultPartitioner = AccIDPartitioner
	kconf.Logger = baseLogger

	kconf.ChangelogTopicNameFormatter = func(storeName string) func(ctx topology.BuilderContext) string {
		return func(ctx topology.BuilderContext) string {
			return fmt.Sprintf(`%s.%s.internal.changelog`, ctx.ApplicationId(), storeName)
		}
	}
	kconf.RepartitionTopicNameFormatter = func(topic string) func(ctx topology.BuilderContext, nodeId topology.NodeId) string {
		return func(ctx topology.BuilderContext, nodeId topology.NodeId) string {
			return fmt.Sprintf(`%s.%s.repartitioned`, ctx.ApplicationId(), topic)
		}
	}

	kconf.Processing.Buffer.Size = 4
	kconf.Processing.ConsumerCount = 2
	kconf.Processing.FailedMessageHandler = (&DLQHandler{
		producer: producer,
		topic:    fmt.Sprintf(`%s.dlq`, s.config.ApplicationID),
		logger:   s.logger,
	}).Handle

	kconf.Processing.Buffer.FlushInterval = 100 * time.Microsecond
	kconf.Consumer.ContextExtractor = contextExtractor

	cc := c.GetGlobalConfig(application.ModuleMetricsReporter).(*metrics2.Config)

	rep := metrics.PrometheusReporter(metrics.ReporterConf{
		System:      cc.System,
		Subsystem:   cc.SubSystem,
		ConstLabels: map[string]string{},
	})
	kconf.MetricsReporter = rep
	// offset to latest
	kconf.Consumer.Offsets.Initial = kafka.OffsetLatest

	kconf.Store.Http.Host = s.config.StoreHost
	kconf.Store.Http.Enabled = true

	s.builder = kstream.NewStreamBuilder(kconf)

	// setup kafka admin
	s.admin, err = sarama.NewAdmin(s.config.BootstrapServers,
		sarama.WithLogger(c.Resolve(application.ModuleBaseLogger).(tlog.Logger)))
	if err != nil {
		return errors.Wrap(err, "streaming: cannot create kafka admin")
	}

	// create topics
	// if err = s.createTopics(c); err != nil {
	// 	return errors.Wrap(err, "streaming: cannot create topics")
	// }

	s.bindAndInit(c)

	tp, err := s.builder.Build()
	if err != nil {
		s.logger.Error(`Error building streams`, s.builder.Topology().Describe())
		return errors.Wrap(err, "error building streams")
	}

	s.logger.Info(tp.Describe())
	if application.DebugMode() {
		visualize(tp.Describe())
	}

	s.topology = tp

	return nil
}

func (s *Streams) Run() error {
	s.runner = s.builder.NewRunner()
	if err := s.runner.Run(s.topology); err != nil {
		return errors.Wrap(err, `cannot start streams`)
	}

	return nil
}

func (s *Streams) Ready() chan bool {
	return nil
}

func (s *Streams) Stop() error {
	s.logger.Info(`streaming stopping`)
	return s.runner.Stop()
}

func (s *Streams) bindAndInit(c container.Container) {
	// bind package internal dependencies
	c.Bind(application.ModuleStoreRegistry, s.builder.StoreRegistry())
	c.Bind(application.ModuleStreamBuilder, s.builder)

	// bind stores
	c.Bind(stores.ModuleTripStoreResolver, new(stores.TripStore))

	// bind global tables
	c.Bind(application.ModuletripGlobalTable, new(tables.TripTable))

	// bind streams
	c.Bind(application.ModuleJobStream, new(streams.JobStream))

	c.Init(
		// stores
		stores.ModuleTripStoreResolver,

		// global tables
		application.ModuletripGlobalTable,

		// streams
		application.ModuleJobStream,
	)
}

func contextExtractor(record kafka.Record) context.Context {
	ctx := context.WithValue(context.Background(),
		domain.ContextKeyAccountID, string(record.Headers().Read([]byte(events.EventHeaderAccountID))))
	ctx = context.WithValue(ctx,
		domain.ContextKeyTraceID, string(record.Headers().Read([]byte(events.EventHeaderTraceID))))

	return ctx
}
