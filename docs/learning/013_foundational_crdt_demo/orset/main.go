package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Node struct {
	Pos      int
	Adds     map[string]map[string]int
	Removes  map[string]map[string]int
	Ch       chan string
	isOnline bool // offline or online
}

func NewNode(pos int) *Node {
	return &Node{
		Pos:      pos,
		Adds:     make(map[string]map[string]int),
		Removes:  make(map[string]map[string]int),
		Ch:       make(chan string),
		isOnline: true,
	}
}

func (n *Node) Add(obj string) {
	defer n.Print()
	s := n.Adds[obj]
	keyVal := make(map[string]int)
	id := uuid.NewString()
	if s == nil {
		keyVal[id] = n.Pos
		n.Adds[obj] = keyVal
		return
	}
	s[id] = n.Pos
}

func (n *Node) Remove(obj string) {
	defer n.Print()
	s := n.Adds[obj]
	if s == nil {
		return
	}
	keyVal := make(map[string]int)
	for k, val := range s {
		keyVal[k] = val
	}
	n.Removes[obj] = keyVal
}
func (n *Node) IsOnline() bool {
	return n.isOnline
}

func (n *Node) Print() {
	fmt.Println("NODE with pos", n.Pos)
	fmt.Println("ADD")
	for k, v := range n.Adds {
		fmt.Println("	Obj: ", k)
		fmt.Println("	Tag: ", v)
	}
	fmt.Println("REMOVE")
	for k, v := range n.Removes {
		fmt.Println("	Obj: ", k)
		fmt.Println("	Tag: ", v)
	}
}

func unionTag(a, b map[string]int) map[string]int {
	final := make(map[string]int)
	for id, pos := range a {
		final[id] = pos
	}
	for id, pos := range b {
		final[id] = pos
	}
	return final
}

func unionObj(a, b map[string]map[string]int) map[string]map[string]int {
	newAdd := make(map[string]map[string]int)

	for obj, tag := range a {
		if _, ok := b[obj]; !ok {
			// only in A
			newAdd[obj] = tag
			continue
		}
		bTag := b[obj]
		newAdd[obj] = unionTag(tag, bTag)
	}

	for obj, tag := range b {
		if _, ok := a[obj]; !ok {
			// only in B
			newAdd[obj] = tag
			continue
		}
		aTag := a[obj]
		newAdd[obj] = unionTag(tag, aTag)
	}
	return newAdd
}

func synchronizeNode(a, b *Node) {
	if !a.IsOnline() || !b.IsOnline() {
		return
	}

	newAdd := unionObj(a.Adds, b.Adds)
	newRemove := unionObj(a.Removes, b.Removes)
	a.Adds = newAdd
	b.Adds = deepCopyUnion(newAdd)
	a.Removes = newRemove
	b.Removes = deepCopyUnion(newRemove)
	b.Print()
}

func deepCopyUnion(src map[string]map[string]int) map[string]map[string]int {
	dst := make(map[string]map[string]int, len(src))
	for k, innerMap := range src {
		newInner := make(map[string]int, len(innerMap))
		for innerK, innerV := range innerMap {
			newInner[innerK] = innerV
		}
		dst[k] = newInner
	}
	return dst
}

func getcommand(r *bufio.Reader) (string, string, error) {
	in, err := r.ReadString('\n')
	if err != nil {
		return "", "", fmt.Errorf("error reading :%w", err)
	}
	in = strings.ToLower(strings.TrimSpace(in))
	if in == "q" || in == "c" {
		return in, "", nil
	}
	pair := strings.Split(in, " ")
	if len(pair) != 2 {
		return "", "", fmt.Errorf("error format, must be <command>-<object>")
	}
	if pair[0] != "a" && pair[0] != "d" {
		return "", "", fmt.Errorf("command must be 'a' for add, or 'd' for delete")
	}

	if len(pair[1]) == 0 {
		return "", "", fmt.Errorf("object can't be empty")
	}

	return pair[0], pair[1], nil
}

type ORSet struct {
	N     int
	nodes []*Node
}

func NewORSet(N int) *ORSet {
	var nodes []*Node
	for i := range N {
		nodes = append(nodes, NewNode(i))
	}
	return &ORSet{
		nodes: nodes,
		N:     N,
	}
}

func (o *ORSet) RunNode(pos int, done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	n := o.nodes[pos]
	for {
		select {
		case <-done:
			fmt.Println("Exit from node", n.Pos)
			return
		case in := <-n.Ch:
			fmt.Println(in, n.Pos, "----")
			pair := strings.Split(in, " ")
			switch pair[0] {
			case "a":
				// fmt.Println(pair)
				n.Add(pair[1])
			case "d":
				n.Remove(pair[1])
			default:
				fmt.Println("hey")
			}

			for i := range o.nodes {
				delay := 100 * rand.Intn(15)
				if i == pos {
					continue
				}
				time.Sleep(time.Millisecond * time.Duration(delay))
				go synchronizeNode(n, o.nodes[i])
			}
			fmt.Println("========")
			fmt.Println()
		}
	}
}

func main() {
	fmt.Printf("OR-SET CRDT DEMO\n")
	fmt.Println("----------------------")
	fmt.Println("")
	N := 5
	var wg sync.WaitGroup
	wg.Add(5)
	done := make(chan struct{})
	orset := NewORSet(N)
	for i := range orset.nodes {
		go orset.RunNode(i, done, &wg)
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		cmd, obj, err := getcommand(reader)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if cmd == "q" || cmd == "c" {
			close(done)
			break
		}

		randPos := rand.Intn(N)
		orset.nodes[randPos].Ch <- fmt.Sprintf("%s %s", cmd, obj)
	}
	wg.Wait()
	fmt.Println("end program ...")
}
