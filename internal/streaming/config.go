package streaming

import (
	"golang-skeleton/pkg/env"
	"time"
)

type TopicConfig struct {
	NumOfPartitions   int32 `json:"num_of_partitions"`
	ReplicationFactor int16 `json:"replication_factor"`
}

type Config struct {
	BootstrapServers         []string      `env:"KAFKA_BOOTSTRAP_SERVERS"`
	ProducerID               string        `env:"STREAMS_KAFKA_PRODUCER_ID" envDefault:"services.hailing.golang-skeleton.producer"`
	ApplicationID            string        `env:"STREAMS_K_STREAM_APPLICATION_ID" envDefault:"services.hailing.golang-skeleton"`
	StoreHost                string        `env:"STREAMS_K_STREAM_STORES_PORT" envDefault:":9001"`
	WorkerPoolNumOfWorkers   int           `env:"STREAMS_K_STREAM_WORKER_POOL_WORKERS" envDefault:"10"`
	WorkerPoolBuffer         int           `env:"STREAMS_K_STREAM_WORKER_POOL_BUFFER" envDefault:"50"`
	ChangelogEnabled         bool          `env:"STREAMS_K_STREAM_CHANGELOG_ENABLED" envDefault:"false"`
	DefaultReplicationFactor int16         `env:"KAFKA_REPLICATION_FACTOR_DEFAULT" envDefault:"2"`
	ConsumerCount            int           `env:"KAFKA_CONSUMER_COUNT" envDefault:"2"`
	BufferSize               int           `env:"KAFKA_BUFFER_SIZE" envDefault:"100"`
	FlushInterval            time.Duration `env:"KAFKA_FLUSH_INTERVAL" envDefault:"1ms"`
}

func (r *Config) Register() error {
	return env.Parse(r)
}
