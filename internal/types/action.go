package types

type ActionInterface interface {
	Apply(value interface{}) interface{}
}
