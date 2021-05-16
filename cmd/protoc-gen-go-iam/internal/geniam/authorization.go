package geniam

import (
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/protobuf/compiler/protogen"
)

type authorizationCodeGenerator struct {
	gen     *protogen.Plugin
	file    *protogen.File
	service *protogen.Service
}

func (c authorizationCodeGenerator) ConstructorGoName() string {
	return "New" + c.StructGoName()
}

func (c authorizationCodeGenerator) StructGoName() string {
	return c.service.GoName + "Authorization"
}

func (c authorizationCodeGenerator) GenerateCode(g *protogen.GeneratedFile) {
	c.generateConstructor(g)
	c.generateStruct(g)
}

func (c authorizationCodeGenerator) GeneratesCode() bool {
	for _, method := range c.service.Methods {
		if getMethodAuthorizationOptions(method) != nil {
			return true
		}
	}
	return false
}

func (c authorizationCodeGenerator) serverGoName() string {
	return c.service.GoName + "Server"
}

func (c authorizationCodeGenerator) generateStruct(g *protogen.GeneratedFile) {
	g.P()
	g.P("type ", c.StructGoName(), " struct {")
	g.P("next ", c.serverGoName())
	for _, method := range c.service.Methods {
		switch getMethodAuthorizationOptions(method).GetStrategy().(type) {
		case *iamv1.MethodAuthorizationOptions_Before:
			beforeMethodAuthorization := g.QualifiedGoIdent(protogen.GoIdent{
				GoImportPath: "go.einride.tech/iam/iamauthz",
				GoName:       "BeforeMethodAuthorization",
			})
			g.P("before", method.GoName, " *", beforeMethodAuthorization)
		case *iamv1.MethodAuthorizationOptions_After:
			afterMethodAuthorization := g.QualifiedGoIdent(protogen.GoIdent{
				GoImportPath: "go.einride.tech/iam/iamauthz",
				GoName:       "AfterMethodAuthorization",
			})
			g.P("after", method.GoName, " *", afterMethodAuthorization)
		}
	}
	g.P("}")
	for _, method := range c.service.Methods {
		contextContext := g.QualifiedGoIdent(protogen.GoIdent{
			GoImportPath: "context",
			GoName:       "Context",
		})
		g.P()
		g.P("func (a *", c.StructGoName(), ") ", method.GoName, "(")
		g.P("ctx ", contextContext, ",")
		g.P("request *", method.Input.GoIdent, ",")
		g.P(") (*", method.Output.GoIdent, ", error) {")
		switch getMethodAuthorizationOptions(method).GetStrategy().(type) {
		case *iamv1.MethodAuthorizationOptions_Before:
			g.P("ctx, err := a.before", method.GoName, ".AuthorizeRequest(ctx, request)")
			g.P("if err != nil {")
			g.P("return nil, err")
			g.P("}")
			g.P("return a.next.", method.GoName, "(ctx, request)")
		case *iamv1.MethodAuthorizationOptions_After:
			g.P("response, err := a.next.", method.GoName, "(ctx, request)")
			g.P("_, errAuth := a.after", method.GoName, ".AuthorizeRequestAndResponse(ctx, request, response)")
			g.P("if errAuth != nil {")
			g.P("return nil, errAuth")
			g.P("}")
			g.P("return response, err")
		case *iamv1.MethodAuthorizationOptions_Open:
			authorize := g.QualifiedGoIdent(protogen.GoIdent{
				GoImportPath: "go.einride.tech/iam/iamauthz",
				GoName:       "Authorize",
			})
			g.P(authorize, "(ctx)")
			g.P("return a.next.", method.GoName, "(ctx, request)")
		default:
			statusError := g.QualifiedGoIdent(protogen.GoIdent{
				GoImportPath: "google.golang.org/grpc/status",
				GoName:       "Error",
			})
			codesInternal := g.QualifiedGoIdent(protogen.GoIdent{
				GoImportPath: "google.golang.org/grpc/codes",
				GoName:       "Internal",
			})
			g.P("return nil, ", statusError, "(", codesInternal, `, "authorization not configured")`)
		}
		g.P("}")
	}
}

func (c authorizationCodeGenerator) generateConstructor(g *protogen.GeneratedFile) {
	memberResolver := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "go.einride.tech/iam/iammember",
		GoName:       "Resolver",
	})
	permissionTester := g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "go.einride.tech/iam/iamauthz",
		GoName:       "PermissionTester",
	})
	g.P()
	g.P("// ", c.ConstructorGoName(), " creates a new authorization middleware for ", c.service.GoName, ".")
	g.P("func ", c.ConstructorGoName(), "(")
	g.P("next ", c.serverGoName(), ",")
	g.P("permissionTester ", permissionTester, ",")
	g.P("memberResolver ", memberResolver, ",")
	g.P(") (*", c.StructGoName(), ", error) {")
	g.P("var result ", c.StructGoName())
	g.P("result.next = next")
	for _, method := range c.service.Methods {
		switch getMethodAuthorizationOptions(method).GetStrategy().(type) {
		case *iamv1.MethodAuthorizationOptions_Before, *iamv1.MethodAuthorizationOptions_After:
			globalFiles := g.QualifiedGoIdent(protogen.GoIdent{
				GoImportPath: "google.golang.org/protobuf/reflect/protoregistry",
				GoName:       "GlobalFiles",
			})
			methodDescriptor := g.QualifiedGoIdent(protogen.GoIdent{
				GoImportPath: "google.golang.org/protobuf/reflect/protoreflect",
				GoName:       "MethodDescriptor",
			})
			fmtErrorf := g.QualifiedGoIdent(protogen.GoIdent{
				GoImportPath: "fmt",
				GoName:       "Errorf",
			})
			methodDescriptorVar := "descriptor" + method.GoName
			methodVar := "method" + method.GoName
			methodName := method.Desc.FullName()
			g.P(methodDescriptorVar, ", err := ", globalFiles, ".FindDescriptorByName(\"", methodName, "\")")
			g.P("if err != nil {")
			g.P(
				"return nil, ", fmtErrorf,
				`("new `, c.service.GoName, ` authorization: failed to find descriptor for `, method.Desc.Name(), `")`,
			)
			g.P("}")
			g.P(methodVar, ", ok := ", methodDescriptorVar, ".(", methodDescriptor, ")")
			g.P("if !ok {")
			g.P(
				"return nil, ", fmtErrorf,
				`("new `, c.service.GoName, ` authorization: got non-method descriptor for `, method.Desc.Name(), `")`,
			)
			g.P("}")
			switch getMethodAuthorizationOptions(method).GetStrategy().(type) {
			case *iamv1.MethodAuthorizationOptions_Before:
				constructor := g.QualifiedGoIdent(protogen.GoIdent{
					GoImportPath: "go.einride.tech/iam/iamauthz",
					GoName:       "NewBeforeMethodAuthorization",
				})
				g.P(
					"before", method.GoName, ", err := ", constructor,
					"(", methodVar, ", permissionTester, memberResolver)",
				)
				g.P("if err != nil {")
				g.P("return nil, ", fmtErrorf, `("new `, c.service.GoName, ` authorization: %w", err)`)
				g.P("}")
				g.P("result.before", method.GoName, " = before", method.GoName)
			case *iamv1.MethodAuthorizationOptions_After:
				constructor := g.QualifiedGoIdent(protogen.GoIdent{
					GoImportPath: "go.einride.tech/iam/iamauthz",
					GoName:       "NewAfterMethodAuthorization",
				})
				g.P(
					"after", method.GoName, ", err := ", constructor,
					"(", methodVar, ", permissionTester, memberResolver)",
				)
				g.P("if err != nil {")
				g.P("return nil, ", fmtErrorf, `("new `, c.service.GoName, ` authorization: %w", err)`)
				g.P("}")
				g.P("result.after", method.GoName, " = after", method.GoName)
			}
		}
	}
	g.P("return &result, nil")
	g.P("}")
}
