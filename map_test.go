package egomap

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	readHandle, writer := New[string, int]()
	reader := readHandle.Reader()
	writer.Set("alice", 25)

	if v, ok := reader.Get("alice"); ok || v != 0 {
		t.Errorf("refresh not applied, expected %v, %v; got %v, %v", 0, false, v, ok)
	}

	writer.Refresh()
	if v, ok := reader.Get("alice"); !ok || v != 25 {
		t.Errorf("refresh applied, expected %v, %v; got %v, %v", 25, true, v, ok)
	}
	writer.Refresh()
	if v, ok := reader.Get("alice"); !ok || v != 25 {
		t.Errorf("refresh reapplied, expected %v, %v; got %v, %v", 25, true, v, ok)
	}

	reader2 := readHandle.Reader()
	writer.Set("bob", 30)

	writer.Refresh()
	writer.Refresh()

	if v, ok := reader.Get("bob"); !ok || v != 30 {
		t.Errorf("refresh applied, expected %v, %v; got %v, %v", 30, true, v, ok)
	}
	if v, ok := reader2.Get("bob"); !ok || v != 30 {
		t.Errorf("refresh reapplied, expected %v, %v; got %v, %v", 30, true, v, ok)
	}
}

const (
	read = iota
	write
)

type op int

type operation struct {
	do  op
	key uint64
}

const (
	size = 1000000
)

var (
	writeFreq = []int{
		2, 100, 1000, 10000, 100000,
	}
)

func gen(freq int) []operation {
	zipf := rand.NewZipf(rand.New(rand.NewSource(time.Now().UnixNano())), 1.2, 100, uint64(size))
	keys := make([]operation, 0, size)
	for i := 0; i < size; i++ {
		var do op = read
		if rand.Int31n(int32(freq)) == 0 {
			do = write
		}
		op := operation{
			do:  do,
			key: zipf.Uint64(),
		}
		keys = append(keys, op)
	}
	return keys
}

func mapTester(b *testing.B, m MapHandle[uint64, int]) {
	for _, freq := range writeFreq {
		// setup
		keys := gen(freq)
		for i := 0; i < size; i++ {
			m.Set(uint64(i), rand.Int())
		}
		b.ResetTimer()

		b.Run(fmt.Sprintf("write_per:%d", freq), func(b *testing.B) {
			b.RunParallel(func(p *testing.PB) {
				idx := rand.Int() % size
				reader := m.Reader()
				for p.Next() {
					action := keys[idx]
					if action.do == write {
						m.Set(action.key, idx)
					} else {
						reader.Get(action.key)
					}
					idx++
					idx %= size
				}
				reader.Close()
			})
		})
	}
}

func BenchmarkEgomap(b *testing.B) {
	m := NewHandle[uint64, int](1)
	mapTester(b, m)
}

type syncMapWrapper[K comparable, V any] struct {
	m *sync.Map
}

func (w *syncMapWrapper[K, V]) Set(key K, val V) {
	w.m.Store(key, val)
}

func (w *syncMapWrapper[K, V]) Reader() Reader[K, V] {
	return w
}

func (w *syncMapWrapper[K, V]) Close() {}

func (w *syncMapWrapper[K, V]) Get(key K) (V, bool) {
	v, ok := w.m.Load(key)
	return v.(V), ok
}

func (w *syncMapWrapper[K, V]) Delete(key K) {
	w.m.Delete(key)
}

func (w *syncMapWrapper[K, V]) Refresh() {}

func BenchmarkSyncMap(b *testing.B) {
	mapTester(b, &syncMapWrapper[uint64, int]{
		m: new(sync.Map),
	})
}

func mapTestWrite(b *testing.B, m MapHandle[uint64, int]) {
	keys := gen(size)
	for i := 0; i < size; i++ {
		m.Set(uint64(i), rand.Int())
	}
	b.ResetTimer()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go func() {
			r := m.Reader()
			idx := rand.Int() % size
			for i := 0; ; idx = (idx + 1) % size {
				select {
				case <-ctx.Done():
					r.Close()
					return
				default:
					r.Get(keys[i].key)
				}
			}
		}()
	}

	id := rand.Int() % size
	for i := 0; i < b.N; i++ {
		m.Set(keys[id].key, i)
		id = (id + 1) % size
	}
}

func BenchmarkEgomap_Write(b *testing.B) {
	m := NewHandle[uint64, int](1)
	mapTestWrite(b, m)
}

func BenchmarkSyncMap_Write(b *testing.B) {
	m := &syncMapWrapper[uint64, int]{
		m: new(sync.Map),
	}
	mapTestWrite(b, m)
}
