package main

import (
	"log"

	"go.einride.tech/iam/cmd/iamctl/internal/rootcmd"
)

func main() {
	log.SetFlags(0)
	rootcmd.Execute()
}
