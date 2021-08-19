package sub

import (
	"github.com/dangkaka/go-kafka-avro"
)

var (
	kafkaServers                                = []string{"localhost:9092"}
	schemaRegistryServers                       = []string{"http://localhost:8081"}
	positionTopic, temperatureTopic, powerTopic = "position", "temperature", "power"
	group                                       = "sensor-group"
)

type Subscriber struct {
	sensorService SensorService
	eventsChan    []<-chan string
}

func New(sensorService SensorService) *Subscriber {
	eventsChan := make([]<-chan string, 0)
	return &Subscriber{sensorService: sensorService, eventsChan: eventsChan}
}

func (s *Subscriber) Subscribe() error {
	eventChan := make(chan string)
	positionConsumer, err := kafka.NewAvroConsumer(
		kafkaServers, schemaRegistryServers, positionTopic, group, s.positionConsumer(eventChan))
	if err != nil {
		return err
	}
	s.eventsChan = append(s.eventsChan, eventChan)
	positionConsumer.Consume()
	return nil
}

func (s *Subscriber) EventsChan() []<-chan string {
	return s.eventsChan
}
