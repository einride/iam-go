package main

import (
	"log"

	"go.einride.tech/iam/cmd/iamctl/rootcmd"
)

func main() {
	log.SetFlags(0)
	rootcmd.Execute()
}
