package main

import (
	"log"
	"os"

	"github.com/0xm0-v1/sik6/internal/cli/cismoke"
)

func main() {
	if err := cismoke.Run(os.Args[1:]); err != nil {
		log.Fatalf("ci smoke error: %v", err)
	}
}
