package generator

import (
	"go/ast"
	"go/format"
	"go/token"
	"os"

	"github.com/momchil-atanasov/gostub/util"
)

const receiverName string = "stub"

func NewGeneratorModel(pkgName, stubName string) *GeneratorModel {
	fileBuilder := NewFileBuilder()
	fileBuilder.SetPackage(pkgName)

	structBuilder := NewStructBuilder()
	structBuilder.SetName(stubName)

	return &GeneratorModel{
		fileBuilder:   fileBuilder,
		structBuilder: structBuilder,
		structName:    stubName,
	}
}

type GeneratorModel struct {
	fileBuilder   *FileBuilder
	structBuilder *StructBuilder
	structName    string
}

// AddImport assures that the specified package name in the specified
// location will be added as an import.
// This function returns the alias to be used in selector expressions.
// If the specified location is already added, then just the alias for
// that package is returned.
func (t *GeneratorModel) AddImport(pkgName, location string) string {
	return t.fileBuilder.AddImport(pkgName, location)
}

func (t *GeneratorModel) AddMethod(config *MethodConfig) error {
	t.createMethodStubField(config)
	t.createMutexField(config)
	t.createArgsForCallField(config)
	if config.HasResults() {
		t.createReturnsField(config)
	}
	t.createStubMethod(config)
	t.createCallCountMethod(config)
	if config.HasParams() {
		t.createArgsForCallMethod(config)
	}
	if config.HasResults() {
		t.createReturnsMethod(config)
	}
	return nil
}

func (t *GeneratorModel) createMethodStubField(config *MethodConfig) {
	builder := NewMethodStubFieldBuilder()
	builder.SetFieldName(config.StubFieldName())
	builder.SetParams(config.MethodParams)
	builder.SetResults(config.MethodResults)
	t.structBuilder.AddField(builder.Build())
}

func (t *GeneratorModel) createMutexField(config *MethodConfig) {
	builder := NewMethodMutexFieldBuilder()
	builder.SetFieldName(config.MutexFieldName())
	builder.SetMutexType(t.resolveMutexType())
	t.structBuilder.AddField(builder.Build())
}

func (t *GeneratorModel) createArgsForCallField(config *MethodConfig) {
	builder := NewMethodArgsFieldBuilder()
	builder.SetFieldName(config.ArgsFieldName())
	builder.SetParams(config.MethodParams)
	t.structBuilder.AddField(builder.Build())
}

func (t *GeneratorModel) createReturnsField(config *MethodConfig) {
	builder := NewReturnsFieldBuilder()
	builder.SetFieldName(config.ReturnsFieldName())
	builder.SetResults(config.MethodResults)
	t.structBuilder.AddField(builder.Build())
}

func (t *GeneratorModel) createStubMethod(config *MethodConfig) {
	methodBuilder := t.createMethodBuilder(config, config.MethodName)
	builder := NewStubMethodBuilder(methodBuilder)
	builder.SetMutexFieldSelector(config.MutexFieldSelector())
	builder.SetArgsFieldSelector(config.ArgsFieldSelector())
	builder.SetReturnsFieldSelector(config.ReturnsFieldSelector())
	builder.SetStubFieldSelector(config.StubFieldSelector())
	builder.SetParams(config.MethodParams)
	builder.SetResults(config.MethodResults)
	t.fileBuilder.AddFunctionDeclaration(builder.Build())
}

func (t *GeneratorModel) createCallCountMethod(config *MethodConfig) {
	methodBuilder := t.createMethodBuilder(config, config.CallCountMethodName())
	builder := NewCountMethodBuilder(methodBuilder)
	builder.SetMutexFieldSelector(config.MutexFieldSelector())
	builder.SetArgsFieldSelector(config.ArgsFieldSelector())
	t.fileBuilder.AddFunctionDeclaration(builder.Build())
}

func (t *GeneratorModel) createArgsForCallMethod(config *MethodConfig) {
	methodBuilder := t.createMethodBuilder(config, config.ArgsForCallMethodName())
	builder := NewArgsMethodBuilder(methodBuilder)
	builder.SetMutexFieldSelector(config.MutexFieldSelector())
	builder.SetArgsFieldSelector(config.ArgsFieldSelector())
	builder.SetParams(config.MethodParams)
	t.fileBuilder.AddFunctionDeclaration(builder.Build())
}

func (t *GeneratorModel) createReturnsMethod(config *MethodConfig) {
	methodBuilder := t.createMethodBuilder(config, config.ReturnsMethodName())
	builder := NewReturnsMethodBuilder(methodBuilder)
	builder.SetMutexFieldSelector(config.MutexFieldSelector())
	builder.SetReturnsFieldSelector(config.ReturnsFieldSelector())
	builder.SetResults(config.MethodResults)
	t.fileBuilder.AddFunctionDeclaration(builder.Build())
}

func (t *GeneratorModel) createMethodBuilder(config *MethodConfig, name string) *MethodBuilder {
	builder := NewMethodBuilder()
	builder.SetName(name)
	builder.SetReceiver(receiverName, t.structName)
	return builder
}

func (t *GeneratorModel) resolveMutexType() ast.Expr {
	alias := t.AddImport("sync", "sync")
	return &ast.SelectorExpr{
		X:   ast.NewIdent(alias),
		Sel: ast.NewIdent("RWMutex"),
	}
}

func (t *GeneratorModel) Save(filePath string) error {
	osFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer osFile.Close()

	t.fileBuilder.AddGeneralDeclaration(t.structBuilder.Build())
	astFile := t.fileBuilder.Build()
	err = format.Node(osFile, token.NewFileSet(), astFile)
	if err != nil {
		return err
	}

	return nil
}

// MethodConfig provides the needed information for the generation
// of a stub implementation of a given method from an interface.
type MethodConfig struct {

	// MethodName specifies the name of the method as seen in the
	// interface it came from.
	MethodName string

	// MethodParams specifies all the parameters of the method.
	// They should have been normalized (i.e. no type reuse and no
	// anonymous parameters) and resolved (i.e. all selector expressions
	// resolved against the generated stub's new namespace)
	MethodParams []*ast.Field

	// MethodResults specifies all the results of the method.
	// They should have been normalized (i.e. no type reuse and no
	// anonymous results) and resolved (i.e. all selector expressions
	// resolved against the generated stub's new namespace)
	MethodResults []*ast.Field
}

func (s *MethodConfig) HasParams() bool {
	return len(s.MethodParams) > 0
}

func (s *MethodConfig) HasResults() bool {
	return len(s.MethodResults) > 0
}

func (s *MethodConfig) MutexFieldName() string {
	return util.ToPrivate(s.MethodName + "Mutex")
}

func (s *MethodConfig) MutexFieldSelector() *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X:   ast.NewIdent(receiverName),
		Sel: ast.NewIdent(s.MutexFieldName()),
	}
}

func (s *MethodConfig) StubFieldName() string {
	return s.MethodName + "Stub"
}

func (s *MethodConfig) StubFieldSelector() *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X:   ast.NewIdent(receiverName),
		Sel: ast.NewIdent(s.StubFieldName()),
	}
}

func (s *MethodConfig) ArgsFieldName() string {
	return util.ToPrivate(s.MethodName + "ArgsForCall")
}

func (s *MethodConfig) ArgsFieldSelector() *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X:   ast.NewIdent(receiverName),
		Sel: ast.NewIdent(s.ArgsFieldName()),
	}
}

func (s *MethodConfig) ReturnsFieldName() string {
	return util.ToPrivate(s.MethodName + "Returns")
}

func (s *MethodConfig) ReturnsFieldSelector() *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X:   ast.NewIdent(receiverName),
		Sel: ast.NewIdent(s.ReturnsFieldName()),
	}
}

func (s *MethodConfig) CallCountMethodName() string {
	return s.MethodName + "CallCount"
}

func (s *MethodConfig) ArgsForCallMethodName() string {
	return s.MethodName + "ArgsForCall"
}

func (s *MethodConfig) ReturnsMethodName() string {
	return s.MethodName + "Returns"
}
