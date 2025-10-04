package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	clock := NewVectorClock(3)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Vector Clock Demo")
	fmt.Println("----------------------------")

	for {
		fmt.Print("==> ")
		text, _ := reader.ReadString('\n')
		text = strings.ReplaceAll(text, "\n", "")
		text = strings.ToLower(text)
		switch text {
		case "a":
			clock.AddLocalEvent(0)
		case "b":
			clock.AddLocalEvent(1)
		case "c":
			clock.AddLocalEvent(2)
		case "ab":
			clock.TransferMessage(0, 1)
		case "ac":
			clock.TransferMessage(0, 2)
		case "bc":
			clock.TransferMessage(1, 2)
		case "ca":
			clock.TransferMessage(2, 0)
		case "cb":
			clock.TransferMessage(2, 1)
		case "ba":
			clock.TransferMessage(1, 0)
		case "q":
			fmt.Println("exiting program ...")
			os.Exit(0)
		default:
			fmt.Println("invalid input")
		}
		clock.PrintAll()
	}

}
