package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

type GlobalCounter struct {
	N     int
	nodes []*Node
}

type Node struct {
	N     int
	pos   int
	local []int
	ch    chan struct{}
}

func NewNode(N int, pos int) *Node {
	return &Node{
		local: make([]int, N),
		N:     N,
		pos:   pos,
		ch:    make(chan struct{}),
	}
}
func (n *Node) Increment() {
	inc := n.local[n.pos]
	inc++
	n.local[n.pos] = inc
	n.ShowTotalCounter()
}

func (n *Node) ShowTotalCounter() {
	total := 0
	for _, val := range n.local {
		total += val
	}

	fmt.Printf("Total: %d from NODE %d   ->  %v\n", total, n.pos, n.local)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (n *Node) Synchronize(t *Node) {
	delay := rand.IntN(15) * 100
	time.Sleep(time.Duration(delay) * time.Millisecond)
	for i := range n.local {
		newMax := max(n.local[i], t.local[i])
		n.local[i] = newMax
		t.local[i] = newMax
	}
	t.ShowTotalCounter()
}

func NewGlobalCounter(N int) *GlobalCounter {
	var nodes []*Node
	for i := range N {
		nodes = append(nodes, NewNode(N, i))
	}

	return &GlobalCounter{
		nodes: nodes,
		N:     N,
	}
}

func (gc *GlobalCounter) RunNode(pos int) {
	for range gc.nodes[pos].ch {
		fmt.Println("Hey, this is from pos ", pos)
		gc.nodes[pos].Increment()
		for i := range gc.N {
			if i == pos {
				continue
			}
			gc.nodes[pos].Synchronize(gc.nodes[i])
		}
	}
}

func main() {
	fmt.Printf("G COUNTER CRDT DEMO\n")
	fmt.Println("----------------------")
	fmt.Println("")
	N := 5
	gc := NewGlobalCounter(N)
	done := make(chan struct{}, 1)
	pos := make(chan int)
	for i := range N {
		go gc.RunNode(i)
	}
	go func() {
		totalCounter := 0
		for {
			reader := bufio.NewReader(os.Stdin)
			in, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("can't read:", err)
			}
			fmt.Println()
			fmt.Println("-------")
			in = strings.Trim(in, "\n")
			if in == "c" {
				close(done)
				return
			}
			posInt, err := strconv.Atoi(in)
			if err != nil {
				fmt.Println("invalid input")
				continue
			}
			if posInt >= N {
				fmt.Println("out of range")
				continue
			}
			totalCounter++
			fmt.Println("TOTAL COUNTER: ", totalCounter)
			gc.nodes[posInt].ch <- struct{}{}
		}
	}()

	for {
		select {
		case <-done:
			fmt.Println("end")
			return
		case p := <-pos:
			fmt.Println(p)
		}

	}
}
