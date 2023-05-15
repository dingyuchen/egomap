package egomap

import (
	"reflect"
	"sync/atomic"
	"testing"
)

func TestNewReadHandler(t *testing.T) {
	type args struct {
		innerMap *leftRightMap[K, V]
		writer   *writer[K, V]
	}
	tests := []struct {
		name string
		args args
		want ReadHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReadHandler(tt.args.innerMap, tt.args.writer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReadHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reader_Close(t *testing.T) {
	type fields struct {
		id         int
		innerMap   *leftRightMap[K, V]
		epoch      *atomic.Uint32
		removeSelf func()
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &reader{
				id:         tt.fields.id,
				innerMap:   tt.fields.innerMap,
				epoch:      tt.fields.epoch,
				removeSelf: tt.fields.removeSelf,
			}
			r.Close()
		})
	}
}

func Test_reader_Get(t *testing.T) {
	type fields struct {
		id         int
		innerMap   *leftRightMap[K, V]
		epoch      *atomic.Uint32
		removeSelf func()
	}
	type args struct {
		key K
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   V
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &reader{
				id:         tt.fields.id,
				innerMap:   tt.fields.innerMap,
				epoch:      tt.fields.epoch,
				removeSelf: tt.fields.removeSelf,
			}
			got, got1 := r.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_readhandler_Reader(t *testing.T) {
	type fields struct {
		counter  int
		innerMap *leftRightMap[K, V]
		writer   writeHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   Reader
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rh := &readhandler{
				counter:  tt.fields.counter,
				innerMap: tt.fields.innerMap,
				writer:   tt.fields.writer,
			}
			if got := rh.Reader(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader() = %v, want %v", got, tt.want)
			}
		})
	}
}
