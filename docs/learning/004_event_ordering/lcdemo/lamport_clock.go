package main

import "fmt"

type LamportClock struct {
	name  string
	event []int
}

func NewLamportClock(name string) *LamportClock {
	if len(name) == 0 {
		panic("name can't be empty")
	}
	return &LamportClock{
		name:  name,
		event: make([]int, 0),
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (lc *LamportClock) AddEvent(x int) {
	in := max(x, lc.GetLastEvent())
	lc.event = append(lc.event, in+1)
}

func (lc *LamportClock) GetLastEvent() int {
	if len(lc.event) == 0 {
		return 0
	}
	return lc.event[len(lc.event)-1]
}

func (lc *LamportClock) Print() {
	fmt.Printf("Node %v: \n", lc.name)
	fmt.Println(lc.event)
	fmt.Println("")
}
