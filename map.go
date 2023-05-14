package egomap

import "sync/atomic"

type leftRightMap[K comparable, V any] struct {
	backingMaps [2]map[K]V
	ptr         *atomic.Int32
}

func (m *leftRightMap[K, V]) readable() map[K]V {
	return m.backingMaps[m.ptr.Load()]
}

func (m *leftRightMap[K, V]) writeable() map[K]V {
	return m.backingMaps[1-m.ptr.Load()]
}

func (m *leftRightMap[K, V]) swap() {
	m.ptr.Store(1 - m.ptr.Load())
}

func NewMap[K comparable, V any]() (ReadHandler[K, V], Writer[K, V]) {
	innerMap := &leftRightMap[K, V]{
		backingMaps: [2]map[K]V{
			{},
			{},
		},
		ptr: &atomic.Int32{},
	}
	writer := NewWriter(innerMap)
	return NewReadHandler(innerMap, writer), writer
}
