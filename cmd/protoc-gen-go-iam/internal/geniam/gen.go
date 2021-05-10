package geniam

import (
	"google.golang.org/protobuf/compiler/protogen"
)

const Version = "development"

const generatedFilenameSuffix = "_iam.go"

func GenerateFile(gen *protogen.Plugin, f *protogen.File) {
	filename := f.GeneratedFilenamePrefix + generatedFilenameSuffix
	g := gen.NewGeneratedFile(filename, f.GoImportPath)
	g.Skip()
	g.P("package ", f.GoPackageName)
	g.P()
	for _, service := range f.Services {
		descriptor := iamServiceDescriptorCodeGenerator{gen: gen, file: f, service: service}
		descriptor.GenerateCode(g)
	}
	g.P()
	g.P("func init() {")
	g.P("// This init function runs after the init function of ", f.GeneratedFilenamePrefix, ".pb.go.")
	g.P("// We depend on the Go initialization order to ensure this.")
	g.P("// See: https://golang.org/ref/spec#Package_initialization")
	for _, service := range f.Services {
		descriptor := iamServiceDescriptorCodeGenerator{gen: gen, file: f, service: service}
		descriptor.GenerateInitFunctionCalls(g)
	}
	g.P("}")
}
