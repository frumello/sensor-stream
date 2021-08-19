package sensors

import (
	"encoding/json"
	"log"
	"sort"
	"sync"
	"time"
)

type Service struct {
	publisher Publisher
}

func New(publisher Publisher) *Service {
	return &Service{publisher: publisher}
}

// MergeSensors receives messages from N channels as string and convert to *Sensor,
// then sends the *Sensor to an output channel
func (s *Service) MergeSensors(inputsChan ...<-chan string) <-chan *Sensor {
	output := make(chan *Sensor)

	var wg sync.WaitGroup
	wg.Add(len(inputsChan))

	for _, input := range inputsChan {
		go func(input <-chan string) {
			for v := range input {
				sensor := new(Sensor)
				if err := json.Unmarshal([]byte(v), sensor); err != nil {
					log.Println("error unmarshal", err)
				}
				output <- sensor
			}
			wg.Done()
		}(input)
	}
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// CollectAndOrder receives *Sensor from a channel ordering them by MeasurementTime,
// then sends the ordered *Sensor to another channel
func (s *Service) CollectAndOrder(sensorChan <-chan *Sensor) <-chan *Sensor{
	sensorList := make([]*Sensor, 0)
	for sensor := range sensorChan {
		sensorList = append(sensorList, sensor)
	}
	sort.Sort(ByMeasurementTime(sensorList))

	output := make(chan *Sensor, len(sensorList))
	for _, n := range sensorList {
		output <- n
	}
	close(output)
	return output
}

// BatchSensorsByDuration Accumulates all *Sensor for a time.Duration period, after the period expires
// if any *Sensor were received it sends them to a channel
func (s *Service) BatchSensorsByDuration(sensorChan <-chan *Sensor, timeout time.Duration) chan []*Sensor {
	batchesChan := make(chan []*Sensor)

	go func() {
		defer close(batchesChan)

		for keepGoing := true; keepGoing; {
			var batch []*Sensor
			expire := time.After(timeout)
			for {
				select {
				case sensor, ok := <-sensorChan:
					if !ok {
						keepGoing = false
						goto done
					}
					batch = append(batch, sensor)
				case <-expire:
					goto done
				}
			}
		done:
			if len(batch) > 0 {
				batchesChan <- batch
			}
		}
	}()

	return batchesChan
}

// SendSensorBatch marshal the []*Sensor and async sends to the publisher,
// that will be sent to the "sensor-output" topic
func (s *Service) SendSensorBatch(batch []*Sensor) error {
	data, err := json.Marshal(batch)
	if err != nil {
		return err
	}
	go s.publishSensorBatch(data)
	return nil
}

func (s *Service) publishSensorBatch(data []byte) {
	if err := s.publisher.PublishSensorBatch(data); err != nil {
		log.Println("can't publish the event")
	}
}