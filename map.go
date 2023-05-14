package egomap

type Map[K comparable, V any] interface {
	Get(K) (V, bool)
	Set(K, V)
}

type egomap[K comparable, V any] struct {
	backingMaps []map[K]V
}

func NewMap[K comparable, V any]() (Reader[K, V], Writer[K, V]) {
	innerMap := &egomap[K, V]{
		backingMaps: []map[K]V{
			{},
			{},
		},
	}
	return NewReader[K, V](innerMap), NewWriter[K, V](innerMap)
}

func (m *egomap[K, V]) Get(key K) (V, bool) {
	v, ok := m.backingMaps[0][key]
	return v, ok
}

func (m *egomap[K, V]) Set(key K, value V) {
	m.backingMaps[0][key] = value
}
