package sub

//go:generate mockgen -source=interfaces.go -destination=interfaces_mock.go -package=sub

type SensorService interface {
}
