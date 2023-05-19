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
		{inputs: []int{}, wantLength: 0, peek: 0},
		{inputs: []int{3}, wantLength: 1, peek: 0},
		{inputs: []int{4, 3}, wantLength: 2, peek: 1},
		{inputs: []int{1, 2, 3, 4, 5}, wantLength: 5, peek: 0},
		{inputs: []int{1, 2, 3, 4, 5}, wantLength: 5, peek: 1},
		{inputs: genArray(1093), wantLength: 1093, peek: 253},
		{inputs: genArray(203), wantLength: 203, peek: 100},
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
				if out := queue.Peek(); out != input {
					t.Errorf("Peek() = %v, want %v", out, input)
				}
			}
			if out := queue.Dequeue(); out != input {
				t.Errorf("Dequeue() = %v, want %v", out, input)
			}
		}
		e := queue.Dequeue()
		if e != 0 {
			t.Errorf("Dequeue() on empty queue, got %v, want %v", e, 0)
		}
		peek := queue.Peek()
		if peek != 0 {
			t.Errorf("Peek() on empty queue, got %v, want zero", peek)
		}
		queue.Enqueue(6969)
		if out := queue.Dequeue(); out != 6969 {
			t.Errorf("Dequeue() = %v, want %v", out, 6969)
		}
	}
}
