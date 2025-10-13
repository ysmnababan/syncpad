package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type Node struct {
	Pos      int
	Adds     map[string]*[]string
	Removes  map[string]*[]string
	Ch       chan string
	isOnline bool // offline or online
}

func NewNode(pos int) *Node {
	return &Node{
		Pos:      pos,
		Adds:     make(map[string]*[]string),
		Removes:  make(map[string]*[]string),
		Ch:       make(chan string),
		isOnline: true,
	}
}

func (n *Node) Add(obj string) {
	s := n.Adds[obj]
	in := fmt.Sprintf("%d:%s", n.Pos, uuid.NewString())
	if s == nil {
		n.Adds[obj] = &[]string{in}
		return
	}

	*s = append(*s, in)
	// n.Print()
}

func (n *Node) Remove(obj string) {
	s := n.Adds[obj]
	// in := fmt.Sprintf("%d:%s", n.Pos, uuid.NewString())
	if s == nil {
		return
	}
	dst := make([]string, len(*s))
	copy(dst, *s)
	n.Removes["obj"] = &dst
}
func (n *Node) IsOnline() bool {
	return n.isOnline
}

func (n *Node) Print() {
	fmt.Println("NODE with pos", n.Pos)
	fmt.Println("ADD")
	for k, v := range n.Adds {
		fmt.Println("	Obj: ", k)
		fmt.Println("	Tag: ", *v)
	}
	fmt.Println("REMOVE")
	for k, v := range n.Removes {
		fmt.Println("	Obj: ", k)
		fmt.Println("	Tag: ", *v)
	}
}


func synchronizeNode(a, b *Node) {
	if !a.IsOnline() || !b.IsOnline() {
		return
	}

	newAdd := make(map[string]*[]string)
	for key, val := range a.Adds {
		if _, ok := b.Adds[key]; !ok {
			newAdd[key] = val
			continue
		}

	}

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
			// fmt.Println(in)
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
				if i == pos {
					continue
				}
				synchronizeNode(n, o.nodes[i])
			}
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
		// n := orset.nodes[randPos]
		// fmt.Println("randpos", randPos)
		orset.nodes[randPos].Ch <- fmt.Sprintf("%s %s", cmd, obj)
		// orset.nodes[randPos].Print()
	}
	wg.Wait()
	fmt.Println("end program ...")
}
