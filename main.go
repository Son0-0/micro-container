package main

import (
	"os"

	"github.com/Son0-0/micro-container/handlers"
)

func main() {
	switch os.Args[1] {
	case "build":
		handlers.Build(os.Args[2:])
	case "run":
		handlers.Run(os.Args[2:])
	case "child":
		handlers.Child(os.Args)
	default:
		panic("invalid command")
	}
}
