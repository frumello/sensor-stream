package pub

import (
	"github.com/dangkaka/go-kafka-avro"
	"time"
)

var (
	kafkaServers          = []string{"localhost:9092"}
	schemaRegistryServers = []string{"http://localhost:8081"}
	topic                 = "sensor-output"
	schema                = `{
								  "name": "SensorBatch",
								  "type": "array",
								  "items": {
										"name": "Sensor",
										"type": "record",
										"fields": [
											  {
												"name": "measurement_time",
												"type": "int"
											  },
											  {
												"name": "device_id",
												"type": "int"
											  }
										]
								  }
							}`
)

//var _ teams.Publisher = (*Publisher)(nil)

type Publisher struct {
	producer *kafka.AvroProducer
}

func New() (*Publisher, error) {
	producer, err := kafka.NewAvroProducer(kafkaServers, schemaRegistryServers)
	if err != nil {
		return nil, err
	}
	return &Publisher{producer: producer}, nil
}

func (p *Publisher) PublishSensorBatch(value []byte) error {
	key := time.Now().String()
	if err := p.producer.Add(topic, schema, []byte(key), value); err != nil {
		return err
	}
	return nil
}
