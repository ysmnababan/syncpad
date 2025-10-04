package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	a := NewLamportClock("A")
	b := NewLamportClock("B")
	c := NewLamportClock("C")
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Lamport's Clock Demo")
	fmt.Println("----------------------------")

	for {
		fmt.Print("==> ")
		text, _ := reader.ReadString('\n')
		text = strings.ReplaceAll(text, "\n", "")
		text = strings.ToLower(text)

		switch text {
		case "a":
			a.AddEvent(0)
		case "b":
			b.AddEvent(0)
		case "c":
			c.AddEvent(0)
		case "ab":
			a.AddEvent(0)
			aNode := a.GetLastEvent()
			b.AddEvent(aNode)
		case "ac":
			a.AddEvent(0)
			aNode := a.GetLastEvent()
			c.AddEvent(aNode)
		case "bc":
			b.AddEvent(0)
			bNode := b.GetLastEvent()
			c.AddEvent(bNode)
		case "ca":
			c.AddEvent(0)
			cNode := c.GetLastEvent()
			a.AddEvent(cNode)
		case "cb":
			c.AddEvent(0)
			cNode := c.GetLastEvent()
			b.AddEvent(cNode)
		case "ba":
			b.AddEvent(0)
			bNode := b.GetLastEvent()
			a.AddEvent(bNode)

		case "q":
			fmt.Println("exiting program ...")
			os.Exit(0)
		default:
			fmt.Println("invalid input")
		}
		a.Print()
		b.Print()
		c.Print()
	}
}
