package main

import (
	"fmt"
	"os"

	"go.einride.tech/iam/cmd/protoc-gen-go-iam/internal/geniam"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const docURL = "https://pkg.go.dev/go.einride.tech/iam"

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--help" {
		_, _ = fmt.Fprintf(os.Stdout, "See %s for usage information.\n", docURL)
		os.Exit(0)
	}
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		var files protoregistry.Files
		for _, file := range gen.Files {
			if err := files.RegisterFile(file.Desc); err != nil {
				return err
			}
		}
		for _, f := range gen.Files {
			if f.Generate {
				if err := geniam.GenerateFile(gen, &files, f); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
