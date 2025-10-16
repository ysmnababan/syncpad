package main

import (
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
)

type ID struct {
	Counter   int
	ReplicaID string
}

func NewIDwithA(c int) ID {
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
	Counter int
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
	Inbox     *Queue[Message]
	RgaState  *RGA
}

func NewReplica(ID string) *Replica {
	return &Replica{
		ReplicaID: ID,
		Inbox:     NewQueue[Message](1024),
		RgaState:  NewRGA(ID),
	}
}

func (r *Replica) AddString(words string, prevID ID, n *Network) {
	for i, val := range words {
		id := ID{
			Counter:   prevID.Counter + i,
			ReplicaID: prevID.ReplicaID,
		}

		r.Add(string(val), id, n)
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

func (r *RGA) ProcessIncomingOp(op Op) {
	switch op.Type {
	case "insert":
		newNode := Node{
			ID:     op.ID,
			PrevID: op.PrevID,
			Value:  op.Value,
		}
		r.Cache[op.ID] = &newNode
		if op.PrevID.Counter == 0 { // is head
			r.Head.InsertChild(&newNode)
		} else {
			n := r.Cache[op.PrevID]
			n.InsertChild(&newNode)
		}
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

func (r *Replica) ProcessInbox() {
	len := r.Inbox.ElementCount()
	for range len {
		msg, _ := r.Inbox.Pop()
		r.RgaState.ProcessIncomingOp(msg.Op)
	}
}

type Queue[T any] struct {
	size          int
	data          []T
	readCounter   int
	insertCounter int
}

func NewQueue[T any](N int) *Queue[T] {
	return &Queue[T]{
		data:          make([]T, N),
		readCounter:   0,
		insertCounter: 0,
		size:          N,
	}
}

func (q *Queue[T]) Push(env T) error {
	if q.insertCounter-q.readCounter == q.size {
		return fmt.Errorf("buffer already full")
	}

	idx := q.insertCounter % q.size
	q.data[idx] = env
	q.insertCounter++
	return nil
}

func (q *Queue[T]) Pop() (*T, error) {
	if q.insertCounter-q.readCounter == 0 {
		return nil, fmt.Errorf("buffer is empty")
	}
	idx := q.readCounter % q.size
	env := q.data[idx]
	q.readCounter++
	return &env, nil
}

func (q *Queue[T]) ElementCount() int {
	return q.insertCounter - q.readCounter
}

type Network struct {
	Replicas []*Replica
	Queue    *Queue[Envelope]
	cache    map[string]*Replica
}

func NewNetwork(numofReplica int, queueSize int) *Network {
	return &Network{
		Replicas: make([]*Replica, 0, numofReplica),
		Queue:    NewQueue[Envelope](queueSize),
		cache:    make(map[string]*Replica),
	}
}
func (n *Network) AddToQueue(op Op) {
	for _, r := range n.Replicas {
		if r.ReplicaID == op.From {
			// don't send to yourself
			continue
		}
		_ = n.Queue.Push(Envelope{
			To:  r.ReplicaID,
			Msg: Message{Op: op},
		})
	}
}

func (n *Network) AddNewReplica(r *Replica) {
	n.Replicas = append(n.Replicas, r)
	n.cache[r.ReplicaID] = r
}

func (n *Network) ShowQueue(w io.Writer) {
	initialReadCounter := n.Queue.readCounter
	initialInsertCounter := n.Queue.insertCounter
	length := n.Queue.ElementCount()
	fmt.Fprintln(w, "Queue with length: ", length)
	for range length {
		env, _ := n.Queue.Pop()
		fmt.Fprintf(w, "%s From: %s To: %s: Message: %s\n",
			strings.ToUpper(env.Msg.Op.Type),
			env.Msg.Op.From,
			env.To,
			env.Msg.Op.Value)
	}
	fmt.Fprintln(w, "-----------")
	n.Queue.insertCounter = initialInsertCounter
	n.Queue.readCounter = initialReadCounter
}

func (n *Network) Broadcast() {
	ec := n.Queue.ElementCount()
	for range ec {
		env, _ := n.Queue.Pop()
		replica, ok := n.cache[env.To]
		if !ok {
			log.Printf("Replica ID: %s not found", env.To)
			continue
		}
		_ = replica.Inbox.Push(env.Msg)
	}
}

func main() {
}
