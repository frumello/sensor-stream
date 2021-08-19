package sensors

import (
	"encoding/json"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestService_MergeSensors(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	//setUp

	service := New(nil)

	t.Run("merge sensors success", func(t *testing.T) {
		// given
		now := time.Now()
		sensorList := newSensorList(now)
		sensorsData := make([]string, 4)

		for i, sensor := range sensorList {
			data, err := json.Marshal(sensor)
			require.NoError(t, err)
			sensorsData[i] = string(data)
		}

		stringChan1 := make(chan string)
		stringChan2 := make(chan string)
		stringChan3 := make(chan string)
		chanList := []<-chan string{stringChan1, stringChan2, stringChan3}

		go func() {
			stringChan1 <- sensorsData[0]
			stringChan3 <- sensorsData[1]
			stringChan2 <- sensorsData[2]
			stringChan3 <- sensorsData[3]
			close(stringChan1)
			close(stringChan2)
			close(stringChan3)
		}()

		// when
		sensorsChan := service.MergeSensors(chanList...)
		mergedSensors := make([]*Sensor, 0)
		for sensor := range sensorsChan {
			mergedSensors = append(mergedSensors, sensor)
		}
		// then
		require.Equal(t, len(sensorsData), len(mergedSensors))
	})
}

func TestService_CollectAndOrder(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	//setUp

	service := New(nil)

	t.Run("collect and order success", func(t *testing.T) {
		// given
		now := time.Now()
		sensorList := newSensorList(now)
		sensorChan := make(chan *Sensor)
		go func() {
			sensorChan <- sensorList[2]
			sensorChan <- sensorList[3]
			sensorChan <- sensorList[1]
			sensorChan <- sensorList[0]
			close(sensorChan)
		}()

		// when
		orderedSensorsChan := service.CollectAndOrder(sensorChan)

		// then
		orderedSensors := make([]*Sensor, 0)
		sort.Sort(ByMeasurementTime(sensorList))
		for sensors := range orderedSensorsChan {
			orderedSensors = append(orderedSensors, sensors)
		}
		require.Equal(t, sensorList, orderedSensors)
	})
}

func TestService_BatchSensorsByDuration(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	//setUp

	service := New(nil)

	t.Run("batch sensors by duration success", func(t *testing.T) {
		// given
		now := time.Now()
		sensorList := newSensorList(now)
		sensorChan := make(chan *Sensor)

		go func() {
			sensorChan <- sensorList[0]
			sensorChan <- sensorList[1]
			sensorChan <- sensorList[2]

			time.Sleep(time.Second) // hit timeout

			sensorChan <- sensorList[3]
			close(sensorChan)
		}()

		// when
		sensorsChan := service.BatchSensorsByDuration(sensorChan, time.Second)

		// then
		i := 0
		for sensors := range sensorsChan {
			if i == 0 {
				require.Equal(t, sensorList[:3], sensors)
			}
			if i == 1 {
				require.Equal(t, sensorList[3:4], sensors)
			}
			i++
		}
	})
}

func TestService_SendSensorBatch(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	//setUp

	publisher := NewMockPublisher(controller)
	service := New(publisher)

	t.Run("send sensor batch success", func(t *testing.T) {
		//given
		now := time.Now()
		sensorList := newSensorList(now)
		expectedTeam, err := json.Marshal(sensorList)
		require.NoError(t, err)

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(1)

		publisher.EXPECT().PublishSensorBatch(gomock.Any()).Do(func(data []byte) {
			require.Equal(t, expectedTeam, data)
			waitGroup.Done()
		})
		err = service.SendSensorBatch(sensorList)
		waitGroup.Wait()

		// then
		require.NoError(t, err)
	})
}

func newSensorList(now time.Time) []*Sensor {
	sensorList := []*Sensor{
		{
			DeviceID:        12,
			MeasurementTime: now,
			Data:            nil,
		},
		{
			DeviceID:        34,
			MeasurementTime: now.Add(time.Second),
			Data:            nil,
		},
		{
			DeviceID:        56,
			MeasurementTime: now.Add(time.Second * 2),
			Data:            nil,
		},
		{
			DeviceID:        78,
			MeasurementTime: now.Add(time.Second * 3),
			Data:            nil,
		},
	}
	return sensorList
}
