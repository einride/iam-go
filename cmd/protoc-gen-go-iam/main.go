package main

import (
	"go.einride.tech/iam/internal/generator"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	protogen.Options{}.Run(generator.Run)
}
