package singleflight

type Key interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Group[K Key, V any] interface {
	Do(key K, fn func(key K) (V, error)) (V, error)

	Forget(key K)
}
