package generator

import (
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
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
			f.Unskip()
			Middleware{Service: service}.Generate(f)
			for _, method := range service.Methods {
				policy := proto.GetExtension(method.Desc.Options(), iamv1.E_Policy).(*iamv1.Policy)
				MiddlewareMethod{Method: method, Policy: policy}.Generate(f)
			}
		}
	}
	return nil
}
