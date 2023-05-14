package egomap

import "sync/atomic"

type reader[K comparable, V any] struct {
	mapPtr   atomic.Int32
	innerMap Map[K, V]
}

type Reader[K comparable, V any] interface {
}

func NewReader[K comparable, V any](innerMap Map[K, V]) Reader[K, V] {
	return &reader[K, V]{
		mapPtr:   atomic.Int32{},
		innerMap: innerMap,
	}
}
