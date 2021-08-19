package sensors

//go:generate mockgen -source=interfaces.go -destination=interfaces_mock.go -package=sensors

type Publisher interface {
	PublishSensorBatch(value []byte) error
}