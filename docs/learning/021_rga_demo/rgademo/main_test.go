package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOneReplica(t *testing.T) {
	net := &Network{
		Replicas: make([]*Replica, 0, 3),
		Queue:    make([]Envelope, 0, 1024),
	}
	r1 := NewReplica("A")
	net.AddNewReplica(r1)

	assert.Equal(t, 1, len(net.Replicas))

	r1.Add("h", NewIDwithA(0), net)
	r1.Add("e", NewIDwithA(1), net)
	r1.Add("l", NewIDwithA(2), net)
	r1.Add("l", NewIDwithA(3), net)
	r1.Add("o", NewIDwithA(4), net)
	r1.Add("X", NewIDwithA(0), net)

	assertText := func(expected string) {
		buff := &bytes.Buffer{}
		r1.PrintTextOnly(buff)
		assert.Equal(t, expected, buff.String())
	}
	assertText("helloX\n")

	r1.Remove(NewIDwithA(6), net)
	assertText("hello\n")

	r1.Remove(NewIDwithA(2), net)
	assertText("hllo\n")

	assert.Equal(t, 0, len(net.Queue))
}
