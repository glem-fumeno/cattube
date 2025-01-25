package main

import (
	"fmt"
	"log"
)

type Node struct {
	Url      string `json:"url"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
}

type Queue struct {
	nodes []*Node
}

func MakeQueue() Queue {
	return Queue{[]*Node{}}
}

func (q *Queue) Enqueue(n *Node) {
	log.Printf("Enqueueing (%X) %v\n", &n, n)
	q.nodes = append(q.nodes, n)
}
func (q *Queue) Dequeue() (*Node, error) {
	if len(q.nodes) == 0 {
		return nil, fmt.Errorf("queue empty")
	}
	ret := q.nodes[0]
	log.Printf("Dequeueing (%X) %v\n", &ret, ret)
	q.nodes = q.nodes[1:]
	return ret, nil
}
func (q Queue) Peek() (*Node, error) {
	if len(q.nodes) == 0 {
		return nil, fmt.Errorf("queue empty")
	}
	return q.nodes[0], nil
}
func (q Queue) IsEmpty() bool {
	return len(q.nodes) == 0
}

func (q Queue) GetAll() []Node {
	ret := make([]Node, len(q.nodes))
	for i, n := range q.nodes {
		ret[i] = *n
	}
	return ret
}
