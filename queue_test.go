package main

import "testing"

func TestQueueEnqueue(t *testing.T) {
	q := MakeQueue()

	have := [3]Node{
		{"a", "a", "a", 1},
		{"b", "b", "b", 2},
		{"c", "c", "c", 3},
	}
	want := [3]Node{
		{"a", "a", "a", 1},
		{"b", "b", "b", 2},
		{"c", "c", "c", 3},
	}
	got := [3]Node{}

	for i := range 3 {
		q.Enqueue(have[i])
	}

	for i := range 3 {
		val, err := q.Dequeue()
		if err != nil {
			t.Errorf("should not get here")
		}
		got[i] = *val
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
