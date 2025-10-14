package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"sync"
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
	n.local[n.pos]++
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

func (gc *GlobalCounter) RunNode(pos int, done chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-gc.nodes[pos].ch:
			fmt.Println("Hey, this is from pos ", pos)
			gc.nodes[pos].Increment()
			for i := range gc.N {
				if i == pos {
					continue
				}
				gc.nodes[pos].Synchronize(gc.nodes[i])
			}
		case <-done:
			close(gc.nodes[pos].ch)
			fmt.Println("done from ", pos)
			return
		}
	}
}

func main() {
	fmt.Printf("G-COUNTER CRDT DEMO\n")
	fmt.Println("----------------------")
	fmt.Println("")
	N := 5
	var wg sync.WaitGroup
	wg.Add(N)
	gc := NewGlobalCounter(N)
	done := make(chan struct{})
	for i := range N {
		go gc.RunNode(i, done, &wg)
	}

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
			break
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
	close(done)

	wg.Wait()
	fmt.Println("done")
}
