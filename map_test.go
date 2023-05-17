package egomap

import (
	"fmt"
	"math/rand"
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

func BenchmarkEgomap(b *testing.B) {
	writeFreq := []int{
		100, 1000, 10000, 100000,
	}
	// concurrency := []int{2, 4, 8, 16} // i only have 10 cores, so 16 is not very useful
	size := 1000000

	for _, freq := range writeFreq {
		// setup
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
		handle := NewHandle[uint64, int](1)

		for i := 0; i < size; i++ {
			handle.Set(uint64(i), rand.Int())
		}
		b.ResetTimer()
		// for _, conc := range concurrency {
		b.Run(fmt.Sprintf("write_per:%d", freq), func(b *testing.B) {
			// b.SetParallelism(conc)
			b.RunParallel(func(p *testing.PB) {
				idx := rand.Int() % size
				reader := handle.Reader()
				for p.Next() {
					action := keys[idx]
					if action.do == write {
						handle.Set(action.key, idx)
					} else {
						reader.Get(action.key)
					}
					idx++
					idx %= size
				}
				reader.Close()
			})
		})
		// }
	}
}
