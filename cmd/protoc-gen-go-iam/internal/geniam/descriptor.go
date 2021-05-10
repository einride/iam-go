package geniam

import (
	"strconv"

	"google.golang.org/protobuf/compiler/protogen"
)

type iamServiceDescriptorCodeGenerator struct {
	gen     *protogen.Plugin
	file    *protogen.File
	service *protogen.Service
}

func (c iamServiceDescriptorCodeGenerator) GlobalFunctionGoName() string {
	return c.service.GoName + "IAM"
}

func (c iamServiceDescriptorCodeGenerator) ServiceDescriptorInterfaceGoName() string {
	return c.service.GoName + "IAMDescriptor"
}

func (c iamServiceDescriptorCodeGenerator) ServiceDescriptorStructGoName() string {
	return private(c.service.GoName + "IAMDescriptor")
}

func (c iamServiceDescriptorCodeGenerator) ServiceDescriptorVariableGoName() string {
	return fileVarName(c.file, "iamDescriptor_"+c.service.GoName)
}

func (c iamServiceDescriptorCodeGenerator) ServiceDescriptorVariableInitFunctionGoName() string {
	return fileVarName(c.file, "init_iamDescriptor_"+c.service.GoName)
}

func (c iamServiceDescriptorCodeGenerator) GenerateCode(g *protogen.GeneratedFile) {
	c.generateGlobalFunction(g)
	c.generateServiceDescriptorInterface(g)
	c.generateServiceDescriptorStruct(g)
	c.generateServiceDescriptorVariable(g)
}

func (c iamServiceDescriptorCodeGenerator) GenerateInitFunctionCalls(g *protogen.GeneratedFile) {
	g.P(c.ServiceDescriptorVariableInitFunctionGoName(), "()")
}

func (c iamServiceDescriptorCodeGenerator) generateGlobalFunction(g *protogen.GeneratedFile) {
	g.P()
	g.P(
		"// ", c.GlobalFunctionGoName(),
		" returns a descriptor for the ", c.service.Desc.Name(), " IAM configuration.",
	)
	g.P("func ", c.GlobalFunctionGoName(), "() ", c.ServiceDescriptorInterfaceGoName(), " {")
	g.P("return &", c.ServiceDescriptorVariableGoName())
	g.P("}")
}

func (c iamServiceDescriptorCodeGenerator) generateServiceDescriptorInterface(g *protogen.GeneratedFile) {
	g.P()
	g.P(
		"// ", c.ServiceDescriptorInterfaceGoName(),
		" describes the ", c.service.Desc.Name(), " IAM configuration.",
	)
	g.P("type ", c.ServiceDescriptorInterfaceGoName(), " interface {")
	if getPredefinedRoles(c.service) != nil {
		g.Unskip()
		iamv1Roles := g.QualifiedGoIdent(protogen.GoIdent{
			GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
			GoName:       "Roles",
		})
		g.P("PredefinedRoles() *", iamv1Roles)
	}
	for _, method := range c.service.Methods {
		if getMethodAuthorizationOptions(method) == nil {
			continue
		}
		g.Unskip()
		iamv1MethodAuthorizationOptions := g.QualifiedGoIdent(protogen.GoIdent{
			GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
			GoName:       "MethodAuthorizationOptions",
		})
		g.P(method.GoName, "() *", iamv1MethodAuthorizationOptions)
	}
	g.P("}")
}

func (c iamServiceDescriptorCodeGenerator) generateServiceDescriptorStruct(g *protogen.GeneratedFile) {
	g.P()
	g.P("type ", c.ServiceDescriptorStructGoName(), " struct {")
	if getPredefinedRoles(c.service) != nil {
		iamv1Roles := g.QualifiedGoIdent(protogen.GoIdent{
			GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
			GoName:       "Roles",
		})
		g.P("predefinedRoles *", iamv1Roles)
	}
	for _, method := range c.service.Methods {
		if getMethodAuthorizationOptions(method) == nil {
			continue
		}
		iamv1MethodAuthorizationOptions := g.QualifiedGoIdent(protogen.GoIdent{
			GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
			GoName:       "MethodAuthorizationOptions",
		})
		g.P(private(method.GoName), " *", iamv1MethodAuthorizationOptions)
	}
	g.P("}")
	if getPredefinedRoles(c.service) != nil {
		iamv1Roles := g.QualifiedGoIdent(protogen.GoIdent{
			GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
			GoName:       "Roles",
		})
		g.P()
		g.P("func (d *", c.ServiceDescriptorStructGoName(), ") PredefinedRoles() *", iamv1Roles, " {")
		g.P("return d.predefinedRoles")
		g.P("}")
	}
	for _, method := range c.service.Methods {
		if getMethodAuthorizationOptions(method) == nil {
			continue
		}
		methodOptions := g.QualifiedGoIdent(protogen.GoIdent{
			GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
			GoName:       "MethodAuthorizationOptions",
		})
		g.P()
		g.P("func (d *", c.ServiceDescriptorStructGoName(), ") ", method.GoName, "() *", methodOptions, " {")
		g.P("return d.", private(method.GoName))
		g.P("}")
	}
}

func (c iamServiceDescriptorCodeGenerator) generateServiceDescriptorVariable(g *protogen.GeneratedFile) {
	globalFiles := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "google.golang.org/protobuf/reflect/protoregistry",
		GoName:       "GlobalFiles",
	})
	getExtension := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "google.golang.org/protobuf/proto",
		GoName:       "GetExtension",
	})
	g.P()
	g.P("var ", c.ServiceDescriptorVariableGoName(), " ", c.ServiceDescriptorStructGoName())
	g.P()
	g.P("func ", c.ServiceDescriptorVariableInitFunctionGoName(), "() {")
	if getPredefinedRoles(c.service) != nil {
		roles := g.QualifiedGoIdent(protogen.GoIdent{
			GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
			GoName:       "Roles",
		})
		ePredefinedRoles := g.QualifiedGoIdent(protogen.GoIdent{
			GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
			GoName:       "E_PredefinedRoles",
		})
		g.P("// Init predefined roles.")
		g.P("serviceDesc, err := ", globalFiles, ".FindDescriptorByName(")
		g.P(strconv.Quote(string(c.service.Desc.FullName())), ",")
		g.P(")")
		g.P("if err != nil {")
		g.P("panic(\"unable to find service descriptor ", c.service.Desc.FullName(), "\")")
		g.P("}")
		g.P(c.ServiceDescriptorVariableGoName(), ".predefinedRoles = ", getExtension, "(")
		g.P("serviceDesc.Options(),")
		g.P(ePredefinedRoles, ",")
		g.P(").(*", roles, ")")
	}
	for _, method := range c.service.Methods {
		if getMethodAuthorizationOptions(method) == nil {
			continue
		}
		methodOptions := g.QualifiedGoIdent(protogen.GoIdent{
			GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
			GoName:       "MethodAuthorizationOptions",
		})
		eMethodAuthorization := g.QualifiedGoIdent(protogen.GoIdent{
			GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
			GoName:       "E_MethodAuthorization",
		})
		g.P("// Init ", method.Desc.Name(), " authorization options.")
		g.P("methodDesc", method.Desc.Name(), ", err := ", globalFiles, ".FindDescriptorByName(")
		g.P(strconv.Quote(string(method.Desc.FullName())), ",")
		g.P(")")
		g.P("if err != nil {")
		g.P("panic(\"unable to find method descriptor ", method.Desc.FullName(), "\")")
		g.P("}")
		g.P(c.ServiceDescriptorVariableGoName(), ".", private(method.GoName), " = ", getExtension, "(")
		g.P("methodDesc", method.Desc.Name(), ".Options(),")
		g.P(eMethodAuthorization, ",")
		g.P(").(*", methodOptions, ")")
		// TODO: For methods with per-resource permissions, resolve the resource descriptors here.
	}
	g.P("}")
}
