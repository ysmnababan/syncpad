package main

import (
	"fmt"
	"sort"
)

type ID struct {
	Counter   uint64
	ReplicaID string
}

func NewIDwithA(c uint64) ID {
	return ID{
		Counter:   c,
		ReplicaID: "A",
	}
}

type Op struct {
	Type   string // insert or delete
	ID     ID
	Value  string
	PrevID ID
	From   string
}

type Message struct {
	Op Op
}

type Envelope struct {
	To  string // replicaID
	Msg Message
}

type Node struct {
	ID        ID
	PrevID    ID // or parent
	Value     string
	Children  []*Node // must be sorted
	Tombstone bool
}

func (p *Node) InsertChild(c *Node) {
	p.Children = append(p.Children, c)
	sort.Slice(p.Children, func(i, j int) bool {
		if p.Children[i].ID.Counter != p.Children[j].ID.Counter {
			return p.Children[i].ID.Counter < p.Children[j].ID.Counter // sort by counter first
		}
		return p.Children[i].ID.ReplicaID < p.Children[j].ID.ReplicaID // then by replicaID
	})
}
func (n *Node) PrintText() {
	if !n.Tombstone {
		fmt.Printf("%s", n.Value)
	}
	for _, ch := range n.Children {
		ch.PrintText()
	}

	if n.IsHead() {
		fmt.Println()
		fmt.Println("-end of sequential text-")
	}
}

func (n *Node) IsHead() bool {
	return n.ID.Counter == 0
}

type RGA struct {
	Counter uint64
	Head    *Node
	Cache   map[ID]*Node
}

func NewRGA(replicaID string) *RGA {
	head := Node{
		ID: ID{
			Counter:   0,
			ReplicaID: replicaID,
		},
		Value:     "",
		Children:  []*Node{},
		Tombstone: false,
	}
	cache := make(map[ID]*Node)
	cache[head.ID] = &head
	return &RGA{
		Counter: 0,
		Head:    &head,
		Cache:   cache,
	}
}

func (r *RGA) GetNewID(replicaID string) ID {
	r.Counter++
	return ID{
		Counter:   r.Counter,
		ReplicaID: replicaID,
	}
}

func (r *RGA) UpdateNode(op Op) {
	switch op.Type {
	case "insert":
		prev, ok := r.Cache[op.PrevID]
		if !ok {
			fmt.Println("missed cache")
			return
		}
		node := &Node{
			ID:        op.ID,
			PrevID:    op.PrevID,
			Value:     op.Value,
			Children:  []*Node{},
			Tombstone: false,
		}
		prev.InsertChild(node)
		r.Cache[node.ID] = node
	case "delete":
		currentNode, ok := r.Cache[op.ID]
		if !ok {
			fmt.Println("missed cache")
			return
		}
		currentNode.Tombstone = true
	default:
		fmt.Println("undefined")
	}
}

type Replica struct {
	ReplicaID string
	Inbox     []Message
	RgaState  *RGA
}

func NewReplica(ID string) *Replica {
	return &Replica{
		ReplicaID: ID,
		Inbox:     make([]Message, 0),
		RgaState:  NewRGA(ID),
	}
}

func (r *Replica) Add(val string, prev ID, n *Network) {
	if val == "" {
		return
	}
	Op := Op{
		Type:   "insert",
		ID:     r.RgaState.GetNewID(r.ReplicaID),
		Value:  val,
		PrevID: prev,
		From:   r.ReplicaID,
	}
	r.RgaState.UpdateNode(Op)

	n.AddToQueue(Op)
}

func (r *Replica) Remove(id ID, n *Network) {
	Op := Op{
		Type: "delete",
		ID:   id,
		From: r.ReplicaID,
	}
	r.RgaState.UpdateNode(Op)
	n.AddToQueue(Op)
}

func (r *Replica) PrintTextOnly() {
	r.RgaState.Head.PrintText()
}

type Network struct {
	Replicas []*Replica
	Queue    []Envelope
}

func (n *Network) AddToQueue(op Op) {
	for _, r := range n.Replicas {
		if r.ReplicaID == op.From {
			// don't send to yourself
			continue
		}
		n.Queue = append(n.Queue, Envelope{
			To:  r.ReplicaID,
			Msg: Message{Op: op},
		})
	}
}

func (n *Network) AddNewReplica(r *Replica) {
	n.Replicas = append(n.Replicas, r)
}

func (n *Network) ShowQueue() {
	fmt.Println("Queue with lenght: ", len(n.Queue))
	for _, env := range n.Queue {
		fmt.Printf("From:%s To: %s: Message: %s\n", env.Msg.Op.From, env.To, env.Msg.Op.Value)
	}
	fmt.Println("-----------")
}

func main() {
	net := &Network{
		Replicas: make([]*Replica, 0, 3),
		Queue:    make([]Envelope, 0, 1024),
	}
	r1 := NewReplica("A")
	net.AddNewReplica(r1)

	r1.Add("h", NewIDwithA(0), net)
	r1.Add("e", NewIDwithA(1), net)
	r1.Add("l", NewIDwithA(2), net)
	r1.Add("l", NewIDwithA(3), net)
	r1.Add("o", NewIDwithA(4), net)
	r1.Add("X", NewIDwithA(0), net)
	r1.PrintTextOnly()

	r1.Remove(NewIDwithA(6), net)
	r1.PrintTextOnly()
	r1.Remove(NewIDwithA(2), net)
	r1.PrintTextOnly()

	net.ShowQueue()
}
