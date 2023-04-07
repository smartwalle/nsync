package singleflight

type Group interface {
	Do(key string, fn func(key string) (interface{}, error)) (interface{}, error)

	Forget(key string)
}
