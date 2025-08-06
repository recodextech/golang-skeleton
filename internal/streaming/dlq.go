package streaming

import (
	"context"
	"fmt"
	"golang-skeleton/internal/domain/adaptors"
	"time"

	"github.com/gmbyapa/kstream/v2/kafka"
)

type DLQHandler struct {
	producer kafka.Producer
	topic    string
	logger   adaptors.Logger
}

func (d *DLQHandler) Handle(err error, record kafka.Record) {
	headers := append(record.Headers(),
		kafka.RecordHeader{
			Key:   []byte("__dlq_trace"),
			Value: []byte(err.Error()),
		},
		kafka.RecordHeader{
			Key:   []byte("__dlq_offset"),
			Value: []byte(fmt.Sprint(record.Offset())),
		},
		kafka.RecordHeader{
			Key:   []byte("__dlq_topic"),
			Value: []byte(record.Topic()),
		},
		kafka.RecordHeader{
			Key:   []byte("__dlq_partition"),
			Value: []byte(fmt.Sprint(record.Partition())),
		},
		kafka.RecordHeader{
			Key:   []byte("__dlq_timestamp"),
			Value: []byte(record.Timestamp().Format(time.RFC3339Nano)),
		},
	)

	rec := d.producer.NewRecord(
		record.Ctx(),
		record.Key(),
		record.Value(),
		d.topic,
		kafka.PartitionAny,
		time.Now(),
		headers,
		``,
	)

	_, _, err = d.producer.ProduceSync(context.Background(), rec)
	if err != nil {
		d.logger.ErrorContext(record.Ctx(), fmt.Sprintf(`DLQ sent failed due to %s`, err))
	}

}
