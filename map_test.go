package egomap

import (
	"fmt"
	"math/rand"
	"testing"
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

func BenchmarkEgomap_100(b *testing.B) {
	sizes := []int{
		100, 1000, 10000, 100000,
	}
	concurrency := []int{2, 4, 8, 16, 32}

	for _, size := range sizes {
		// setup
		zipf := rand.NewZipf(rand.New(rand.NewSource(0)), 1.2, 10, uint64(size))
		keys := make([]uint64, 0, 1000000)
		for i := 0; i < cap(keys); i++ {
			keys = append(keys, zipf.Uint64())
		}
		handle := NewHandle[uint64, int](1)

		for i := 0; i < size; i++ {
			handle.Set(keys[i], rand.Int())
		}
		for _, conc := range concurrency {
			b.Run(fmt.Sprintf("size:%d|conc:%d", size, conc), func(b *testing.B) {
				b.SetParallelism(conc)
				b.RunParallel(func(p *testing.PB) {
					count := 0
					starting := rand.Int() % size
					reader := handle.Reader()
					for p.Next() {
						if count%100 == 0 {
							handle.Set(keys[starting], rand.Int())
						} else {
							reader.Get(keys[starting])
						}
						starting++
						count++
						starting %= 1000000
						count %= 100
					}
					reader.Close()
				})
			})
		}
	}
}
