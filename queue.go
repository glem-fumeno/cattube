package main

import (
	"fmt"
)

type Node struct {
	Url      string `json:"url"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
	Size     int    `json:"size"`
}

type Queue struct {
	nodes []Node
}

func MakeQueue() Queue {
	return Queue{
		[]Node{},
	}
}

func (q *Queue) Enqueue(n Node) {
	q.nodes = append(q.nodes, n)
}
func (q *Queue) Dequeue() (*Node, error) {
	if len(q.nodes) == 0 {
		return nil, fmt.Errorf("queue empty")
	}
	ret := q.nodes[0]
	q.nodes = q.nodes[1:]
	return &ret, nil
}
