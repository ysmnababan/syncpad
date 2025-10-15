package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
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

func (n *Node) PrintText(w io.Writer) {
	if !n.Tombstone {
		fmt.Fprintf(w, "%s", n.Value)
	}
	for _, ch := range n.Children {
		ch.PrintText(w)
	}

	if n.IsHead() {
		fmt.Fprintln(w)
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

func (r *Replica) PrintTextOnly(w io.Writer) {
	r.RgaState.Head.PrintText(w)
}

type Queue struct {
	size          int
	data          []Envelope
	readCounter   int
	insertCounter int
}

func NewQueue(N int) *Queue {
	return &Queue{
		data:          make([]Envelope, N),
		readCounter:   0,
		insertCounter: 0,
		size:          N,
	}
}

func (q *Queue) Push(env Envelope) error {
	if q.insertCounter-q.readCounter == q.size {
		return fmt.Errorf("buffer already full")
	}

	idx := q.insertCounter % q.size
	q.data[idx] = env
	q.insertCounter++
	return nil
}

func (q *Queue) Pop() (*Envelope, error) {
	if q.insertCounter-q.readCounter == 0 {
		return nil, fmt.Errorf("buffer is empty")
	}
	idx := q.readCounter % q.size
	env := q.data[idx]
	q.readCounter++
	return &env, nil
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

func (n *Network) ShowQueue(w io.Writer) {
	fmt.Fprintln(w, "Queue with length: ", len(n.Queue))
	for _, env := range n.Queue {
		fmt.Fprintf(w, "%s From: %s To: %s: Message: %s\n",
			strings.ToUpper(env.Msg.Op.Type),
			env.Msg.Op.From,
			env.To,
			env.Msg.Op.Value)
	}
	fmt.Fprintln(w, "-----------")
}

func (n *Network) Broadcast() {
	// for _, env := range n.Queue {
	// env.
	// }
}

func main() {
}
