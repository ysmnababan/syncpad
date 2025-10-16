package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOneReplica(t *testing.T) {
	net := NewNetwork(1, 1024)
	r1 := NewReplica("A")
	net.AddNewReplica(r1)

	assert.Equal(t, 1, len(net.Replicas))

	assertText := func(expected string) {
		buff := &bytes.Buffer{}
		r1.PrintTextOnly(buff)
		assert.Equal(t, expected, buff.String())
	}

	r1.Add("h", NewIDwithA(0), net)
	r1.Add("e", NewIDwithA(1), net)
	r1.Add("l", NewIDwithA(2), net)
	r1.Add("l", NewIDwithA(3), net)
	r1.Add("o", NewIDwithA(4), net)
	r1.Add("X", NewIDwithA(0), net)

	assertText("helloX\n")

	r1.Remove(NewIDwithA(6), net)
	assertText("hello\n")

	r1.Remove(NewIDwithA(2), net)
	assertText("hllo\n")

	assert.Equal(t, 0, net.Queue.ElementCount())
}

func TestSendToNetworkQueue(t *testing.T) {
	net := NewNetwork(1, 1024)

	r1 := NewReplica("A")
	net.AddNewReplica(r1)
	r2 := NewReplica("B")
	net.AddNewReplica(r2)
	assert.Equal(t, 2, len(net.Replicas))

	assertText := func(expected string) {
		buff := &bytes.Buffer{}
		r1.PrintTextOnly(buff)
		assert.Equal(t, expected, buff.String())
	}

	r1.Add("h", NewIDwithA(0), net)
	r1.Add("e", NewIDwithA(1), net)
	r1.Add("l", NewIDwithA(2), net)
	r1.Add("l", NewIDwithA(3), net)
	r1.Add("o", NewIDwithA(4), net)
	assertText("hello\n")

	assert.Equal(t, 5, net.Queue.ElementCount())
	buff := &bytes.Buffer{}
	net.ShowQueue(buff)
	expectedMessage := `Queue with length:  5
INSERT From: A To: B: Message: h
INSERT From: A To: B: Message: e
INSERT From: A To: B: Message: l
INSERT From: A To: B: Message: l
INSERT From: A To: B: Message: o
-----------
`
	assert.Equal(t, expectedMessage, buff.String())
}

func TestPush(t *testing.T) {
	queue := NewQueue[Envelope](10)
	env := Envelope{To: "X"}
	assert.Equal(t, 0, queue.readCounter)
	assert.Equal(t, 0, queue.insertCounter)
	assert.Equal(t, 10, len(queue.data))

	err := queue.Push(env)
	require.NoError(t, err)
	assert.Equal(t, 10, len(queue.data))
	assert.Equal(t, "X", queue.data[0].To)
	for range 9 {
		err := queue.Push(env)
		require.NoError(t, err)
	}
	assert.Equal(t, 10, queue.insertCounter)

	err = queue.Push(env)
	assert.Equal(t, "buffer already full", err.Error())
	queue.readCounter = 5
	env.To = "Y"
	err = queue.Push(env)
	require.NoError(t, err)
	assert.Equal(t, "Y", queue.data[0].To)
}

func TestPop(t *testing.T) {
	queue := NewQueue[Envelope](10)
	env := Envelope{To: "X"}
	assert.Equal(t, 0, queue.readCounter)
	assert.Equal(t, 0, queue.insertCounter)
	assert.Equal(t, 10, len(queue.data))
	_, err := queue.Pop()
	assert.Equal(t, "buffer is empty", err.Error())

	for range 5 {
		err := queue.Push(env)
		require.NoError(t, err)
	}

	for i := range 5 {
		data, err := queue.Pop()
		require.NoError(t, err)
		assert.NotNil(t, data)
		assert.Equal(t, i+1, queue.readCounter)
	}
	_, err = queue.Pop()
	assert.Equal(t, "buffer is empty", err.Error())
}

func TestBroadcast(t *testing.T) {
	net := NewNetwork(1, 1024)

	r1 := NewReplica("A")
	net.AddNewReplica(r1)
	r2 := NewReplica("B")
	net.AddNewReplica(r2)
	assert.Equal(t, 2, len(net.Replicas))

	assertText := func(expected string) {
		buff := &bytes.Buffer{}
		r1.PrintTextOnly(buff)
		assert.Equal(t, expected, buff.String())
	}

	r1.Add("h", NewIDwithA(0), net)
	r1.Add("e", NewIDwithA(1), net)
	r1.Add("l", NewIDwithA(2), net)
	r1.Add("l", NewIDwithA(3), net)
	r1.Add("o", NewIDwithA(4), net)
	assertText("hello\n")

	assert.Equal(t, 5, net.Queue.ElementCount())
	net.Broadcast()
	assert.Equal(t, 5, r2.Inbox.ElementCount())
	r2Inbox := []*Message{}
	for range 5 {
		msg, err := r2.Inbox.Pop()
		require.NoError(t, err)
		r2Inbox = append(r2Inbox, msg)
	}
	for i, ibx := range r2Inbox {
		assert.Equal(t, "A", ibx.Op.From)
		assert.Equal(t, i+1, ibx.Op.ID.Counter)
		assert.Equal(t, "A", ibx.Op.ID.ReplicaID)
		assert.Equal(t, "insert", ibx.Op.Type)
	}

	assert.Equal(t, "h", r2Inbox[0].Op.Value)
	assert.Equal(t, "e", r2Inbox[1].Op.Value)
	assert.Equal(t, "l", r2Inbox[2].Op.Value)
	assert.Equal(t, "l", r2Inbox[3].Op.Value)
	assert.Equal(t, "o", r2Inbox[4].Op.Value)
}

func TestProcessIncomingOp(t *testing.T) {
	net := NewNetwork(1, 1024)

	r1 := NewReplica("A")
	net.AddNewReplica(r1)
	r2 := NewReplica("B")
	net.AddNewReplica(r2)

	assertText := func(expected string) {
		buff := &bytes.Buffer{}
		r2.PrintTextOnly(buff)
		assert.Equal(t, expected, buff.String())
	}
	assertText("\n")
	r1.Add("h", NewIDwithA(0), net)
	net.Broadcast()

	assert.Equal(t, 1, r2.Inbox.ElementCount())
	msg, err := r2.Inbox.Pop()
	require.NoError(t, err)
	op := msg.Op
	assert.Equal(t, "h", op.Value)
	assert.Equal(t, 1, op.ID.Counter)
	assert.Equal(t, "A", op.ID.ReplicaID)
	r2.RgaState.ProcessIncomingOp(op)
	assertText("h\n")

	r1.Add("i", NewIDwithA(1), net)
	r1.Add("k", NewIDwithA(2), net)
	r1.Add("s", NewIDwithA(3), net)
	net.Broadcast()
	for range 3 {
		msg, err := r2.Inbox.Pop()
		require.NoError(t, err)
		r2.RgaState.ProcessIncomingOp(msg.Op)
	}

	assertText("hiks\n")
	r1.Remove(NewIDwithA(2), net)
	net.Broadcast()
	msg, err = r2.Inbox.Pop()
	require.NoError(t, err)
	fmt.Println(msg.Op.Type)
	r2.RgaState.ProcessIncomingOp(msg.Op)
	assertText("hks\n")
}

func TestProcessInbox(t *testing.T) {
	net := NewNetwork(1, 1024)

	r1 := NewReplica("A")
	net.AddNewReplica(r1)
	r2 := NewReplica("B")
	net.AddNewReplica(r2)
	r1.Add("h", NewIDwithA(0), net)
	r1.Add("e", NewIDwithA(1), net)
	r1.Add("l", NewIDwithA(2), net)
	r1.Add("l", NewIDwithA(3), net)
	r1.Add("o", NewIDwithA(4), net)

	net.Broadcast()
	r2.ProcessInbox()
	assertText := func(expected string) {
		buff := &bytes.Buffer{}
		r2.PrintTextOnly(buff)
		assert.Equal(t, expected, buff.String())
	}
	assertText("hello\n")

	r1.Remove(NewIDwithA(1), net)
	r1.Remove(NewIDwithA(2), net)
	r1.Remove(NewIDwithA(3), net)
	r1.Remove(NewIDwithA(4), net)
	net.Broadcast()
	r2.ProcessInbox()
	assertText("o\n")
}

func TestAddString(t *testing.T) {
	// setup
	net := NewNetwork(1, 1024)
	r1 := NewReplica("A")
	net.AddNewReplica(r1)
	r2 := NewReplica("B")
	net.AddNewReplica(r2)

	assertText := func(r *Replica, expected string) {
		buff := &bytes.Buffer{}
		r.PrintTextOnly(buff)
		assert.Equal(t, expected, buff.String())
	}

	// execute
	r1.AddString("hello", NewIDwithA(0), net)

	//assert
	assertText(r1, "hello\n")

	net.Broadcast()
	r2.ProcessInbox()
	assertText(r2, "hello\n")
}

func TestEditSamePos(t *testing.T) {
	// setup
	net := NewNetwork(1, 1024)

	r1 := NewReplica("A")
	net.AddNewReplica(r1)
	r2 := NewReplica("B")
	net.AddNewReplica(r2)
	newID := func(id string, counter int) ID {
		return ID{
			ReplicaID: id,
			Counter:   counter,
		}
	}
	r1.AddString("hello", newID("A", 0), net)

	net.Broadcast()
	r2.ProcessInbox()
	assertText := func(r *Replica, expected string) {
		buff := &bytes.Buffer{}
		r.PrintTextOnly(buff)
		assert.Equal(t, expected, buff.String())
	}

	// assert
	assertText(r1, "hello\n")
	assertText(r2, "hello\n")

	r1.AddString(" world", newID("A", 5), net)
	assertText(r1, "hello world\n")

	r2.Add(" ", newID("A", 5), net)
	r2.AddString("dunia", newID("B", 1), net)
	assertText(r2, "hello dunia\n")

	// try to sync
	net.Broadcast()
	r1.ProcessInbox()
	assertText(r1, "hello dunia world\n")
	r2.ProcessInbox()
	assertText(r2, "hello dunia world\n")

}
