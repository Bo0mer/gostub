package util_test

import (
	"go/ast"
	"go/token"

	. "github.com/momchil-atanasov/gostub/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AST", func() {

	Describe("CreateField", func() {
		var field *ast.Field
		var fieldType ast.Expr

		BeforeEach(func() {
			fieldType = ast.NewIdent("string")
			field = CreateField("Name", fieldType)
		})

		It("has correct name", func() {
			Ω(field.Names).ShouldNot(BeNil())
			Ω(field.Names).Should(HaveLen(1))
			Ω(field.Names[0].String()).Should(Equal("Name"))
		})

		It("has correct type", func() {
			Ω(field.Type).Should(Equal(fieldType))
		})
	})

	Describe("FieldReuseCount", func() {
		var anonymousField *ast.Field
		var field *ast.Field

		BeforeEach(func() {
			anonymousField = &ast.Field{}
			field = &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent("first"),
					ast.NewIdent("second"),
				},
			}
		})

		It("returns 1 for anonymous fields", func() {
			Ω(FieldReuseCount(anonymousField)).Should(Equal(1))
		})

		It("returns the correct count for a reused field", func() {
			Ω(FieldReuseCount(field)).Should(Equal(2))
		})
	})

	Describe("CreateFuncType", func() {
		var funcType *ast.FuncType

		BeforeEach(func() {
			funcType = CreateFuncType()
		})

		It("is not nil", func() {
			Ω(funcType).ShouldNot(BeNil())
		})

		It("has zero params", func() {
			Ω(funcType.Params).ShouldNot(BeNil())
			Ω(funcType.Params.List).ShouldNot(BeNil())
			Ω(funcType.Params.List).Should(HaveLen(0))
		})

		It("has zero results", func() {
			Ω(funcType.Results).ShouldNot(BeNil())
			Ω(funcType.Results.List).ShouldNot(BeNil())
			Ω(funcType.Results.List).Should(HaveLen(0))
		})
	})

	Describe("FuncTypeParamCount", func() {
		var funcType *ast.FuncType
		var emptyFuncType *ast.FuncType

		BeforeEach(func() {
			funcType = &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{},
						&ast.Field{},
					},
				},
			}
			emptyFuncType = &ast.FuncType{}
		})

		It("returns zero for empty func types", func() {
			Ω(FuncTypeParamCount(emptyFuncType)).Should(Equal(0))
		})

		It("return correct param count for non-empty func types", func() {
			Ω(FuncTypeParamCount(funcType)).Should(Equal(2))
		})
	})

	Describe("EachParamInFunc", func() {
		var funcType *ast.FuncType
		var firstParam *ast.Field
		var secondParam *ast.Field
		var emptyFuncType *ast.FuncType

		BeforeEach(func() {
			firstParam = &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent("first"),
				},
			}
			secondParam = &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent("second"),
				},
			}
			funcType = &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						firstParam,
						secondParam,
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{},
				},
			}
			emptyFuncType = &ast.FuncType{}
		})

		It("returns all fields", func() {
			fieldChan := EachParamInFunc(funcType)
			Ω(<-fieldChan).Should(Equal(firstParam))
			Ω(<-fieldChan).Should(Equal(secondParam))
			Eventually(fieldChan).Should(BeClosed())
		})

		It("returns no fields for empty func", func() {
			fieldChan := EachParamInFunc(emptyFuncType)
			Eventually(fieldChan).Should(BeClosed())
		})
	})

	Describe("FuncTypeResultCount", func() {
		var funcType *ast.FuncType
		var emptyFuncType *ast.FuncType

		BeforeEach(func() {
			funcType = &ast.FuncType{
				Results: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{},
						&ast.Field{},
					},
				},
			}
			emptyFuncType = &ast.FuncType{}
		})

		It("returns zero for empty func types", func() {
			Ω(FuncTypeResultCount(emptyFuncType)).Should(Equal(0))
		})

		It("return correct result count for non-empty func types", func() {
			Ω(FuncTypeResultCount(funcType)).Should(Equal(2))
		})
	})

	Describe("EachResultInFunc", func() {
		var funcType *ast.FuncType
		var firstResult *ast.Field
		var secondResult *ast.Field
		var emptyFuncType *ast.FuncType

		BeforeEach(func() {
			firstResult = &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent("first"),
				},
			}
			secondResult = &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent("second"),
				},
			}
			funcType = &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						firstResult,
						secondResult,
					},
				},
			}
			emptyFuncType = &ast.FuncType{}
		})

		It("returns all fields", func() {
			fieldChan := EachResultInFunc(funcType)
			Ω(<-fieldChan).Should(Equal(firstResult))
			Ω(<-fieldChan).Should(Equal(secondResult))
			Eventually(fieldChan).Should(BeClosed())
		})

		It("returns no fields for empty func", func() {
			fieldChan := EachResultInFunc(emptyFuncType)
			Eventually(fieldChan).Should(BeClosed())
		})
	})

	Describe("EachDeclarationInFile", func() {
		var file *ast.File
		var firstDeclaration ast.Decl
		var secondDeclaration ast.Decl
		var thirdDeclaration ast.Decl

		BeforeEach(func() {
			firstDeclaration = &ast.BadDecl{}
			secondDeclaration = &ast.FuncDecl{}
			thirdDeclaration = &ast.GenDecl{}
			file = &ast.File{
				Decls: []ast.Decl{
					firstDeclaration,
					secondDeclaration,
					thirdDeclaration,
				},
			}
		})

		It("returns all declarations", func() {
			decChan := EachDeclarationInFile(file)
			Ω(<-decChan).Should(Equal(firstDeclaration))
			Ω(<-decChan).Should(Equal(secondDeclaration))
			Ω(<-decChan).Should(Equal(thirdDeclaration))
			Eventually(decChan).Should(BeClosed())
		})
	})

	Describe("EachGenericDeclarationInFile", func() {
		var file *ast.File
		var firstDeclaration ast.Decl
		var secondDeclaration ast.Decl
		var thirdDeclaration ast.Decl
		BeforeEach(func() {
			firstDeclaration = &ast.GenDecl{
				Tok: token.IMPORT,
			}
			secondDeclaration = &ast.FuncDecl{}
			thirdDeclaration = &ast.GenDecl{
				Tok: token.CONST,
			}
			file = &ast.File{
				Decls: []ast.Decl{
					firstDeclaration,
					secondDeclaration,
					thirdDeclaration,
				},
			}
		})

		It("returns only generic declarations", func() {
			decChan := EachGenericDeclarationInFile(file)
			Ω(<-decChan).Should(Equal(firstDeclaration))
			Ω(<-decChan).Should(Equal(thirdDeclaration))
			Eventually(decChan).Should(BeClosed())
		})
	})

	Describe("EachSpecificationInGenericDeclaration", func() {
		var decl *ast.GenDecl
		var firstSpec ast.Spec
		var secondSpec ast.Spec

		BeforeEach(func() {
			firstSpec = &ast.ValueSpec{
				Type: ast.NewIdent("first"),
			}
			secondSpec = &ast.ValueSpec{
				Type: ast.NewIdent("second"),
			}
			decl = &ast.GenDecl{
				Specs: []ast.Spec{
					firstSpec,
					secondSpec,
				},
			}
		})

		It("returns all specifications", func() {
			specChan := EachSpecificationInGenericDeclaration(decl)
			Ω(<-specChan).Should(Equal(firstSpec))
			Ω(<-specChan).Should(Equal(secondSpec))
			Eventually(specChan).Should(BeClosed())
		})
	})

	Describe("EachTypeSpecificationInGenericDeclaration", func() {
		var decl *ast.GenDecl
		var firstSpec ast.Spec
		var secondSpec ast.Spec
		var thirdSpec ast.Spec

		BeforeEach(func() {
			firstSpec = &ast.TypeSpec{
				Name: ast.NewIdent("first"),
			}
			secondSpec = &ast.ValueSpec{
				Type: ast.NewIdent("second"),
			}
			thirdSpec = &ast.TypeSpec{
				Name: ast.NewIdent("third"),
			}
			decl = &ast.GenDecl{
				Specs: []ast.Spec{
					firstSpec,
					secondSpec,
					thirdSpec,
				},
			}
		})

		It("returns all specifications", func() {
			specChan := EachTypeSpecificationInGenericDeclaration(decl)
			Ω(<-specChan).Should(Equal(firstSpec))
			Ω(<-specChan).Should(Equal(thirdSpec))
			Eventually(specChan).Should(BeClosed())
		})
	})

	Describe("EachInterfaceDeclarationInFile", func() {
		var file *ast.File
		var firstSpec ast.Spec
		var thirdSpec ast.Spec

		BeforeEach(func() {
			firstSpec = &ast.TypeSpec{
				Name: ast.NewIdent("first"),
				Type: &ast.InterfaceType{},
			}
			thirdSpec = &ast.TypeSpec{
				Name: ast.NewIdent("third"),
				Type: &ast.InterfaceType{},
			}

			firstDeclaration := &ast.GenDecl{
				Specs: []ast.Spec{
					firstSpec,
				},
			}
			secondDeclaration := &ast.GenDecl{
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Type: &ast.FuncType{},
					},
				},
			}
			thirdDeclaration := &ast.GenDecl{
				Specs: []ast.Spec{
					thirdSpec,
				},
			}
			file = &ast.File{
				Decls: []ast.Decl{
					firstDeclaration,
					secondDeclaration,
					thirdDeclaration,
				},
			}
		})

		It("returns all specifications", func() {
			specChan := EachInterfaceDeclarationInFile(file)
			Ω(<-specChan).Should(Equal(firstSpec))
			Ω(<-specChan).Should(Equal(thirdSpec))
			Eventually(specChan).Should(BeClosed())
		})
	})

	Describe("EachMethodInInterfaceType", func() {
		var iFaceType *ast.InterfaceType
		var firstMethod *ast.Field
		var outlier *ast.Field
		var secondMethod *ast.Field

		BeforeEach(func() {
			firstMethod = &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent("First"),
				},
				Type: &ast.FuncType{},
			}
			secondMethod = &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent("Second"),
				},
				Type: &ast.FuncType{},
			}
			outlier = &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent("NotMethod"),
				},
				Type: &ast.SelectorExpr{},
			}
			iFaceType = &ast.InterfaceType{
				Methods: &ast.FieldList{
					List: []*ast.Field{
						firstMethod,
						outlier,
						secondMethod,
					},
				},
			}
		})

		It("returns all specifications", func() {
			funcChan := EachMethodInInterfaceType(iFaceType)
			Ω(<-funcChan).Should(Equal(firstMethod))
			Ω(<-funcChan).Should(Equal(secondMethod))
			Eventually(funcChan).Should(BeClosed())
		})
	})
})
