package generator

import (
	authorizationv1 "go.einride.tech/protoc-gen-go-authorization-policy/proto/gen/einride/authorization/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

// Run the authorization policy generator plugin.
func Run(gen *protogen.Plugin) error {
	for _, file := range gen.Files {
		if !file.Generate {
			continue
		}
		f := gen.NewGeneratedFile(file.GeneratedFilenamePrefix+"_authorization.pb.go", file.GoImportPath)
		f.Skip()
		FileHeader{File: file}.Generate(f)
		for _, service := range file.Services {
			for _, method := range service.Methods {
				policy := proto.GetExtension(method.Desc.Options(), authorizationv1.E_Policy).(*authorizationv1.Policy)
				if policy == nil {
					continue
				}
				f.Unskip()
			}
			// TODO: Generate middleware.
		}
	}
	return nil
}
