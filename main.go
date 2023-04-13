package main

import (
	"log"
)

func main() {
	n := NewNode()
	if err := n.Node.Run(); err != nil {
		log.Fatal(err)
	}
}
