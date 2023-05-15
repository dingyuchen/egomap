package egomap

import (
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
}
