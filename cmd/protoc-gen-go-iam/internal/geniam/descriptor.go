package geniam

import (
	"google.golang.org/protobuf/compiler/protogen"
)

type descriptorCodeGenerator struct {
	gen     *protogen.Plugin
	file    *protogen.File
	service *protogen.Service
}

func (c descriptorCodeGenerator) ConstructorGoName() string {
	return "New" + c.service.GoName + "IAMDescriptor"
}

func (c descriptorCodeGenerator) GenerateCode(g *protogen.GeneratedFile) {
	c.generateConstructor(g)
}

func (c descriptorCodeGenerator) GeneratesCode() bool {
	if getPredefinedRoles(c.service) != nil {
		return true
	}
	for _, method := range c.service.Methods {
		if getMethodAuthorizationOptions(method) != nil {
			return true
		}
	}
	return false
}

func (c descriptorCodeGenerator) generateConstructor(g *protogen.GeneratedFile) {
	iamDescriptor := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "go.einride.tech/iam/iamreflect",
		GoName:       "IAMDescriptor",
	})
	newIAMDescriptor := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "go.einride.tech/iam/iamreflect",
		GoName:       "NewIAMDescriptor",
	})
	globalFiles := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "google.golang.org/protobuf/reflect/protoregistry",
		GoName:       "GlobalFiles",
	})
	serviceDescriptor := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "google.golang.org/protobuf/reflect/protoreflect",
		GoName:       "ServiceDescriptor",
	})
	fmtErrorf := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "fmt",
		GoName:       "Errorf",
	})
	g.P()
	g.P("// ", c.ConstructorGoName(), " returns a new ", c.service.GoName, " IAM descriptor.")
	g.P("func ", c.ConstructorGoName(), "() (*", iamDescriptor, ", error) {")
	g.P("descriptor, err := ", globalFiles, ".FindDescriptorByName(\"", c.service.Desc.FullName(), "\")")
	g.P("if err != nil {")
	g.P("return nil, ", fmtErrorf, `("new `, c.service.GoName, ` IAM descriptor: %w", err)`)
	g.P("}")
	g.P("service, ok := descriptor.(", serviceDescriptor, ")")
	g.P("if !ok {")
	g.P("return nil, ", fmtErrorf, `("new `, c.service.GoName, `IAM descriptor: got non-service descriptor")`)
	g.P("}")
	g.P("return ", newIAMDescriptor, "(service, ", globalFiles, ")")
	g.P("}")
}
