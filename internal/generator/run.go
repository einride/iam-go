package generator

import (
	"google.golang.org/protobuf/compiler/protogen"
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
				MiddlewareMethod{Method: method, Authorization: getAuthorization(method)}.Generate(f)
			}
		}
	}
	return nil
}
