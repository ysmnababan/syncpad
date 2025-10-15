package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOneReplica(t *testing.T) {
	net := &Network{
		Replicas: make([]*Replica, 0, 1),
		Queue:    make([]Envelope, 0, 1024),
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

	assert.Equal(t, 0, len(net.Queue))
}

func TestSendToNetworkQueue(t *testing.T) {
	net := &Network{
		Replicas: make([]*Replica, 0, 1),
		Queue:    make([]Envelope, 0, 1024),
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

	assert.Equal(t, 5, len(net.Queue))
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
