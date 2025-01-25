package main

import "testing"

func TestQueue(t *testing.T) {
	q := MakeQueue()

	have := [3]Node{
		{"a", "a", "a"},
		{"b", "b", "b"},
		{"c", "c", "c"},
	}
	want := [3]Node{
		{"a", "a", "a"},
		{"b", "b", "b"},
		{"c", "c", "c"},
	}
	got := [3]Node{}

	for i := range 3 {
		q.Enqueue(&have[i])
	}

	for i := range 3 {
		want_val, err := q.Peek()
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		got_val, err := q.Dequeue()
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if *got_val != *want_val {
			t.Errorf("got %v, want %v", got_val, want_val)
		}
		got[i] = *got_val
	}

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestDequeEmptyQueue(t *testing.T) {
	q := MakeQueue()

	_, err := q.Dequeue()

	if err == nil {
		t.Errorf("did not handle empty deque correctly")
	}
}
