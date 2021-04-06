package main

import (
	"go.einride.tech/protoc-gen-go-authorization-policy/internal/generator"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	protogen.Options{}.Run(generator.Run)
}
