package main

import "fmt"

type Process struct {
	N     int // nodes count
	nodes []*VectorNode
}

type VectorClock struct {
	N         int // nodes count
	processes []*Process
}

func NewProcess(N int) *Process {
	return &Process{
		N:     N,
		nodes: make([]*VectorNode, 0),
	}
}

func (p *Process) getLastNode() *VectorNode {
	if len(p.nodes) == 0 {
		return NewVectorNode(p.N)
	}
	return p.nodes[len(p.nodes)-1]
}

func (p *Process) AddLocalEvent(idx int) {
	last := p.getLastNode()
	p.nodes = append(p.nodes, last.IncrementNode(idx))
}

func (p *Process) PrintProcess() {
	for _, node := range p.nodes {
		fmt.Printf("%v => ", *node)
	}
	fmt.Println()
	fmt.Println()
}

func (p *Process) AppendNode(v *VectorNode) {
	p.nodes = append(p.nodes, v)
}

func NewVectorClock(N int) *VectorClock {
	if N <= 0 {
		panic("N must be positive integers")
	}
	vc := &VectorClock{
		N: N,
	}
	for range N {
		p := NewProcess(N)
		vc.processes = append(vc.processes, p)
	}
	return vc
}

func (v *VectorClock) AddLocalEvent(i int) {
	v.processes[i].AddLocalEvent(i)
}

func (v *VectorClock) PrintAll() {
	fmt.Println("RESULT ")
	for i, p := range v.processes {
		ch := rune('a' + i)
		fmt.Printf("Process %s: \n", string(ch))
		p.PrintProcess()
	}
}

func (v *VectorClock) TransferMessage(a, b int) {
	if a == b {
		return
	}
	v.processes[a].AddLocalEvent(a)
	lNode := v.processes[a].getLastNode()
	maxNode := lNode.MaxVector(v.processes[b].getLastNode())
	newNode := maxNode.IncrementNode(b)
	v.processes[b].AppendNode(newNode)
}
