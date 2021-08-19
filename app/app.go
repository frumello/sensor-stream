package app

import (
	"log"
	"sensor-stream/eventbus/pub"
	"sensor-stream/eventbus/sub"
	"sensor-stream/sensors"
	"time"
)

type App struct {
	service    *sensors.Service
	subscriber *sub.Subscriber
}

// New creates configuration for instances
func New() (*App, error) {
	publisher, err := pub.New()
	if err != nil {
		return nil, err
	}
	service := sensors.New(publisher)
	subscriber := sub.New(service)
	if err = subscriber.Subscribe(); err != nil {
		return nil, err
	}
	return &App{service: service, subscriber: subscriber}, nil
}

// Process run all the service methods
func (a *App) Process() error {

	sensorChan := a.service.MergeSensors(a.subscriber.EventsChan()...)
	orderedSensorChan := a.service.CollectAndOrder(sensorChan)
	batchChan := a.service.BatchSensorsByDuration(orderedSensorChan, time.Minute)

	for batch := range batchChan {
		if err := a.service.SendSensorBatch(batch); err != nil {
			log.Println("error during batch", err)
		}
	}
	return nil
}
