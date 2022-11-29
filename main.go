package main

import (
	"os"

	"github.com/Son0-0/micro-container/handlers"
)

func main() {
	switch os.Args[1] {
	case "run":
		handlers.Run(os.Args)
	case "child":
		handlers.Child(os.Args)
	default:
		panic("invalid command")
	}
}
