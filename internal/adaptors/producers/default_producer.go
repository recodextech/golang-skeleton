package producers

import (
	"golang-skeleton/internal/domain/application"

	"github.com/gmbyapa/kstream/v2/kafka"
	"github.com/gmbyapa/kstream/v2/kafka/adaptors/librd"
	"github.com/recodextech/container"
	tlog "github.com/tryfix/log"
	"github.com/tryfix/metrics"
)

type DefaultProducer struct {
	kafka.Producer
}

func (d *DefaultProducer) Init(c container.Container) (err error) {
	conf := c.GetGlobalConfig(application.ModuleDefaultProducer).(*Config)
	pcon := librd.NewProducerConfig()
	pcon.Logger = c.Resolve(application.ModuleBaseLogger).(tlog.Logger).NewLog(tlog.Prefixed("producers"))
	pcon.BootstrapServers = conf.BootStrapServers
	pcon.Acks = kafka.WaitForAll
	pcon.Transactional.Enabled = false
	pcon.MetricsReporter = c.Resolve(application.ModuleStreamReporter).(metrics.Reporter)
	d.Producer, err = librd.NewProducer(pcon)

	return err
}

func (d *DefaultProducer) KafkaProducer() kafka.Producer {
	return d.Producer
}
