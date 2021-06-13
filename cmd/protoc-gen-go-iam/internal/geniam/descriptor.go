package geniam

import (
	"fmt"

	"go.einride.tech/iam/iamannotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type descriptorCodeGenerator struct {
	gen     *protogen.Plugin
	files   *protoregistry.Files
	service *protogen.Service
}

func (c descriptorCodeGenerator) StructGoName() string {
	return c.service.GoName + "IAMDescriptor"
}

func (c descriptorCodeGenerator) ConstructorGoName() string {
	return "New" + c.service.GoName + "IAMDescriptor"
}

func (c descriptorCodeGenerator) GenerateCode(g *protogen.GeneratedFile) error {
	c.generateStruct(g)
	if err := c.generateConstructor(g); err != nil {
		return err
	}
	return nil
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

func (c descriptorCodeGenerator) generateConstructor(g *protogen.GeneratedFile) error {
	predefinedRoles := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
		GoName:       "PredefinedRoles",
	})
	longRunningOperationsAuthorizationOptions := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
		GoName:       "LongRunningOperationsAuthorizationOptions",
	})
	methodAuthorizationOptions := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
		GoName:       "MethodAuthorizationOptions",
	})
	protoUnmarshal := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "google.golang.org/protobuf/proto",
		GoName:       "Unmarshal",
	})
	fmtErrorf := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "fmt",
		GoName:       "Errorf",
	})
	pReturnErr := func(format string, args ...interface{}) {
		g.P(
			"return nil, ", fmtErrorf, `("new `, c.service.Desc.Name(),
			` IAM descriptor: `, fmt.Sprintf(format, args...), `: %w", err)`,
		)
	}
	g.P()
	g.P("// ", c.ConstructorGoName(), " returns a new ", c.service.GoName, " IAM descriptor.")
	g.P("func ", c.ConstructorGoName(), "() (*", c.StructGoName(), ", error) {")
	g.P("result := ", c.StructGoName(), "{")
	if getPredefinedRoles(c.service) != nil {
		g.P("PredefinedRoles: &", predefinedRoles, "{},")
	}
	if getLongRunningOperationsAuthorizationOptions(c.service) != nil {
		g.P("LongRunningOperationsAuthorization: &", longRunningOperationsAuthorizationOptions, "{},")
	}
	for _, method := range c.service.Methods {
		if getMethodAuthorizationOptions(method) != nil {
			g.P(method.GoName, "Authorization: &", methodAuthorizationOptions, "{},")
		}
	}
	g.P("}")
	if options := getPredefinedRoles(c.service); options != nil {
		data, err := proto.Marshal(options)
		if err != nil {
			return err
		}
		g.P("if err := ", protoUnmarshal, "(")
		g.P(fmt.Sprintf("%#v", data), ",")
		g.P("result.PredefinedRoles,")
		g.P("); err != nil {")
		pReturnErr("unmarshal pre-defined roles")
		g.P("}")
	}
	if options := getLongRunningOperationsAuthorizationOptions(c.service); options != nil {
		resolvedOptions, err := iamannotations.ResolveLongRunningOperationsAuthorizationOptions(
			options, c.files, c.service.Desc.ParentFile().Package(),
		)
		if err != nil {
			return err
		}
		data, err := proto.Marshal(resolvedOptions)
		if err != nil {
			return err
		}
		g.P("if err := ", protoUnmarshal, "(")
		g.P(fmt.Sprintf("%#v", data), ",")
		g.P("result.LongRunningOperationsAuthorization,")
		g.P("); err != nil {")
		pReturnErr("unmarshal long-running operations authorization")
		g.P("}")
	}
	for _, method := range c.service.Methods {
		options := getMethodAuthorizationOptions(method)
		if options == nil {
			continue
		}
		resolvedOptions, err := iamannotations.ResolveMethodAuthorizationOptions(
			options, c.files, method.Desc.ParentFile().Package(),
		)
		if err != nil {
			return err
		}
		data, err := proto.Marshal(resolvedOptions)
		if err != nil {
			return err
		}
		g.P("if err := ", protoUnmarshal, "(")
		g.P(fmt.Sprintf("%#v", data), ",")
		g.P("result.", method.GoName, "Authorization,")
		g.P("); err != nil {")
		pReturnErr("unmarshal %s method authorization", method.Desc.Name())
		g.P("}")
	}
	g.P("return &result, nil")
	g.P("}")
	return nil
}

func (c descriptorCodeGenerator) generateStruct(g *protogen.GeneratedFile) {
	predefinedRoles := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
		GoName:       "PredefinedRoles",
	})
	longRunningOperationsAuthorizationOptions := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
		GoName:       "LongRunningOperationsAuthorizationOptions",
	})
	methodAuthorizationOptions := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
		GoName:       "MethodAuthorizationOptions",
	})
	g.P("type ", c.StructGoName(), " struct {")
	if getPredefinedRoles(c.service) != nil {
		g.P("PredefinedRoles *", predefinedRoles)
	}
	if getLongRunningOperationsAuthorizationOptions(c.service) != nil {
		g.P("LongRunningOperationsAuthorization *", longRunningOperationsAuthorizationOptions)
	}
	for _, method := range c.service.Methods {
		if getMethodAuthorizationOptions(method) != nil {
			g.P(method.GoName, "Authorization *", methodAuthorizationOptions)
		}
	}
	g.P("}")
}
