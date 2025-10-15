package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOneReplica(t *testing.T) {
	net := &Network{
		Replicas: make([]*Replica, 0, 1),
		Queue:    NewQueue(1024),
	}
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
	net := &Network{
		Replicas: make([]*Replica, 0, 1),
		Queue:    NewQueue(1024),
	}
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
	queue := NewQueue(10)
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
	queue := NewQueue(10)
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
