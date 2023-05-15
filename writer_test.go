package egomap

// import (
// 	"reflect"
// 	"sync"
// 	"sync/atomic"
// 	"testing"
//
// 	"github.com/dingyuchen/egomap/internal/oplog"
// 	"github.com/dingyuchen/egomap/internal/queue"
// )
//
// func TestNewWriter(t *testing.T) {
// 	type args struct {
// 		innerMap *leftRightMap[string, int]
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want *writer[string, int]
// 	}{
// 		// {}
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := NewWriter(tt.args.innerMap); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewWriter() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func Test_writer_Delete(t *testing.T) {
// 	type fields struct {
// 		mu       *sync.RWMutex
// 		innerMap *leftRightMap[string, int]
// 		oplog    oplog.Log[string, int]
// 		readers  map[int]*reader[string, int]
// 		scan     queue.Queue[*reader[string, int]]
// 	}
// 	type args struct {
// 		key string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		{
// 			"deletes item",
// 			fields{
// 				mu: new(sync.RWMutex),
// 				innerMap: &leftRightMap[string, int]{
// 					backingMaps: [2]map[string]int{
// 						{},
// 						{},
// 					},
// 					ptr: &atomic.Int32{},
// 				},
// 			},
// 			args{
// 				"alice",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := &writer[string, int]{
// 				mu:       tt.fields.mu,
// 				innerMap: tt.fields.innerMap,
// 				oplog:    tt.fields.oplog,
// 				readers:  tt.fields.readers,
// 				scan:     tt.fields.scan,
// 			}
// 			w.Delete(tt.args.key)
// 		})
// 	}
// }
//
// func Test_writer_Refresh(t *testing.T) {
// 	type fields struct {
// 		mu       *sync.RWMutex
// 		innerMap *leftRightMap[K, V]
// 		oplog    oplog.Log
// 		readers  map[int]*reader[K, V]
// 		scan     queue.Queue
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := &writer{
// 				mu:       tt.fields.mu,
// 				innerMap: tt.fields.innerMap,
// 				oplog:    tt.fields.oplog,
// 				readers:  tt.fields.readers,
// 				scan:     tt.fields.scan,
// 			}
// 			w.Refresh()
// 		})
// 	}
// }
//
// func Test_writer_Set(t *testing.T) {
// 	type fields struct {
// 		mu       *sync.RWMutex
// 		innerMap *leftRightMap[K, V]
// 		oplog    oplog.Log
// 		readers  map[int]*reader[K, V]
// 		scan     queue.Queue
// 	}
// 	type args struct {
// 		key   K
// 		value V
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := &writer{
// 				mu:       tt.fields.mu,
// 				innerMap: tt.fields.innerMap,
// 				oplog:    tt.fields.oplog,
// 				readers:  tt.fields.readers,
// 				scan:     tt.fields.scan,
// 			}
// 			w.Set(tt.args.key, tt.args.value)
// 		})
// 	}
// }
//
// func Test_writer_applyWrites(t *testing.T) {
// 	type fields struct {
// 		mu       *sync.RWMutex
// 		innerMap *leftRightMap[K, V]
// 		oplog    oplog.Log
// 		readers  map[int]*reader[K, V]
// 		scan     queue.Queue
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := &writer{
// 				mu:       tt.fields.mu,
// 				innerMap: tt.fields.innerMap,
// 				oplog:    tt.fields.oplog,
// 				readers:  tt.fields.readers,
// 				scan:     tt.fields.scan,
// 			}
// 			w.applyWrites()
// 		})
// 	}
// }
//
// func Test_writer_register(t *testing.T) {
// 	type fields struct {
// 		mu       *sync.RWMutex
// 		innerMap *leftRightMap[K, V]
// 		oplog    oplog.Log
// 		readers  map[int]*reader[K, V]
// 		scan     queue.Queue
// 	}
// 	type args struct {
// 		r *reader[K, V]
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := &writer{
// 				mu:       tt.fields.mu,
// 				innerMap: tt.fields.innerMap,
// 				oplog:    tt.fields.oplog,
// 				readers:  tt.fields.readers,
// 				scan:     tt.fields.scan,
// 			}
// 			w.register(tt.args.r)
// 		})
// 	}
// }
//
// func Test_writer_unregister(t *testing.T) {
// 	type fields struct {
// 		mu       *sync.RWMutex
// 		innerMap *leftRightMap[K, V]
// 		oplog    oplog.Log
// 		readers  map[int]*reader[K, V]
// 		scan     queue.Queue
// 	}
// 	type args struct {
// 		id int
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := &writer{
// 				mu:       tt.fields.mu,
// 				innerMap: tt.fields.innerMap,
// 				oplog:    tt.fields.oplog,
// 				readers:  tt.fields.readers,
// 				scan:     tt.fields.scan,
// 			}
// 			w.unregister(tt.args.id)
// 		})
// 	}
// }
