package sensors

import (
	"time"
)

type Sensor struct {
	DeviceID        int                    `json:"device_id"`
	MeasurementTime time.Time              `json:"measurement_time"`
	Data            map[string]interface{} `json:",inline"`
}

type ByMeasurementTime []*Sensor

func (s ByMeasurementTime) Len() int      { return len(s) }
func (s ByMeasurementTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByMeasurementTime) Less(i, j int) bool {
	return s[i].MeasurementTime.Before(s[j].MeasurementTime)
}
