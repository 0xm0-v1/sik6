package main

import (
	"log"
	"os"

	"github.com/0xm0-v1/sik6/internal/cli/dbseed"
)

func main() {
	if err := dbseed.Run(os.Args[1:]); err != nil {
		log.Fatalf("db seed error: %v", err)
	}
}
