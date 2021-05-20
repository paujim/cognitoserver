package entities

type ParameterStorer interface {
	Get(key string) (string, error)
}
