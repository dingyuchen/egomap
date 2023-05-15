package queue

import "testing"

func genArray(size int) []int {
	arr := make([]int, 0, size)
	for i := 0; i < size; i++ {
		arr = append(arr, i)
	}
	return arr
}

func TestQueue(t *testing.T) {
	type args struct {
		inputs     []int
		peek       int
		wantLength int
	}
	tests := []args{
		{inputs: []int{1, 2, 3, 4, 5}, wantLength: 5, peek: 0},
		{inputs: []int{1, 2, 3, 4, 5}, wantLength: 5, peek: 1},
		{inputs: genArray(maxInternalSlice), wantLength: maxInternalSlice, peek: maxInternalSlice - 1},
		{inputs: genArray(maxInternalSlice * 2), wantLength: maxInternalSlice * 2, peek: maxInternalSlice + initialInternalSize},
	}

	queue := New[int]()
	for _, test := range tests {
		for _, input := range test.inputs {
			queue.Enqueue(input)
		}
		if queue.Len() != test.wantLength {
			t.Errorf("Len() = %v, want %v", queue.Len(), test.wantLength)
		}
		for i, input := range test.inputs {
			if i == test.peek {
				if out := queue.Peek(); *out != input {
					t.Errorf("Peek() = %v, want %v", *out, input)
				}
			}
			if out, _ := queue.Dequeue(); out != input {
				t.Errorf("Dequeue() = %v, want %v", out, input)
			}
		}
		_, err := queue.Dequeue()
		if err != ErrEmpty {
			t.Errorf("Dequeue() on empty queue, got %v, want %v", err, ErrEmpty)
		}
		peek := queue.Peek()
		if peek != nil {
			t.Errorf("Peek() on empty queue, got %v, want nil", *peek)
		}
		queue.Enqueue(6969)
		if out, _ := queue.Dequeue(); out != 6969 {
			t.Errorf("Dequeue() = %v, want %v", out, 6969)
		}
	}
}
