package singleflight

type Key interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Group[K Key] interface {
	Do(key K, fn func(key K) (interface{}, error)) (interface{}, error)

	Forget(key K)
}
