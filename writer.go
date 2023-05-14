package egomap

import "github.com/dingyuchen/egomap/internal/oplog"

type writer[K comparable, V any] struct {
	oplog oplog.Log[K, V]
}

type Writer[K comparable, V any] interface {
	Set(K, V)
}

func NewWriter[K comparable, V any](innerMap Map[K, V]) Writer[K, V] {
	return &writer[K, V]{}
}

func (w *writer[K, V]) Set(key K, value V) {
}
