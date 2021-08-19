package sub

import (
	cluster "github.com/bsm/sarama-cluster"
	"github.com/dangkaka/go-kafka-avro"
	"log"
)

func (s *Subscriber) positionConsumer(eventChan chan string) kafka.ConsumerCallbacks {
	return kafka.ConsumerCallbacks{
		OnDataReceived: func(msg kafka.Message) {
			eventChan <- msg.Value
		},
		OnError: func(err error) {
			log.Println("Consumer error", err)
		},
		OnNotification: func(notification *cluster.Notification) {
			log.Println(notification)
		},
	}
}
