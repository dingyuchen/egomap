package egomap

type Map[K comparable, V any] interface {
	Get(K) (V, bool)
	Set(K, V)
}

type egomap[K comparable, V any] struct {
	innerMap []map[K]V
}

func NewMap[K comparable, V any]() Map[K, V] {
	return &egomap[K, V]{
		innerMap: []map[K]V{
			{},
			{},
		},
	}
}

func (m *egomap[K, V]) Get(key K) (V, bool) {
	v, ok := m.innerMap[0][key]
	return v, ok
}

func (m *egomap[K, V]) Set(key K, value V) {
	m.innerMap[0][key] = value
}
