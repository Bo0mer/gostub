package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/mokiat/gostub/generator"
	"github.com/mokiat/gostub/resolution"
	"github.com/mokiat/gostub/util"
)

func main() {
	locator := resolution.NewLocator()

	interfacePath := "github.com/mokiat/gostub/cmd/hoho/example"
	interfaceName := "Example"
	context := resolution.NewSingleLocationContext(interfacePath)
	d, err := locator.FindIdentType(context, ast.NewIdent(interfaceName))
	if err != nil {
		log.Fatal(err)
	}

	typeName := fmt.Sprintf("monitoring%s", interfaceName)

	fileBuilder := generator.NewFileBuilder()
	fileBuilder.SetPackage("examplemw")

	model := NewModel(interfacePath, interfaceName, typeName, fileBuilder)
	generator := Generator{
		Model:    model,
		Locator:  locator,
		Resolver: generator.NewResolver(model, locator),
	}

	err = generator.ProcessInterface(d)
	if err != nil {
		log.Fatal(err)
	}

	err = model.WriteSource(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

// FieldBuilder builds a struct field.
type FieldBuilder struct {
	fieldName string
	fieldType ast.Expr
}

func NewFieldBuilder(fieldName, typePackage, fieldType string) *FieldBuilder {
	return &FieldBuilder{
		fieldName: fieldName,
		fieldType: &ast.SelectorExpr{
			X:   ast.NewIdent(typePackage),
			Sel: ast.NewIdent(fieldType),
		},
	}
}

func (b *FieldBuilder) Build() *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{
			ast.NewIdent(b.fieldName),
		},
		Type: b.fieldType,
	}
}

func splitPackageType(s string) (pkg, typ string) {
	i := strings.LastIndex(s, ".")
	if i < 0 {
		return "", s
	}
	return s[:i], s[i+1:] // skip the dot
}

type Model struct {
	fileBuilder   *generator.FileBuilder
	structBuilder *generator.StructBuilder
	structName    string

	timePackageAlias string
	totalOps         *ast.SelectorExpr
	failedOps        *ast.SelectorExpr
	opsDuration      *ast.SelectorExpr
}

// NewModel creates new model.
func NewModel(interfacePath, interfaceName, structName string, fileBuilder *generator.FileBuilder) *Model {
	structBuilder := generator.NewStructBuilder()
	structBuilder.SetName(structName)

	fileBuilder.AddDeclarationBuilder(structBuilder)

	m := &Model{
		fileBuilder:   fileBuilder,
		structBuilder: structBuilder,
		structName:    structName,

		totalOps: &ast.SelectorExpr{
			X:   ast.NewIdent("m"),        // receiver name
			Sel: ast.NewIdent("totalOps"), // member name
		},
		failedOps: &ast.SelectorExpr{
			X:   ast.NewIdent("m"),         // receiver name
			Sel: ast.NewIdent("failedOps"), // member name
		},
		opsDuration: &ast.SelectorExpr{
			X:   ast.NewIdent("m"),            // receiver name
			Sel: ast.NewIdent("opsDuraation"), // member name
		},
	}

	m.timePackageAlias = m.AddImport("", "time")

	structBuilder.AddFieldBuilder(NewFieldBuilder("next", m.AddImport("", interfacePath), interfaceName))
	metricsAlias := m.AddImport("", "githbu.com/go-kit/kit/metrics")
	structBuilder.AddFieldBuilder(NewFieldBuilder("totalOps", metricsAlias, "Counter"))
	structBuilder.AddFieldBuilder(NewFieldBuilder("failedOps", metricsAlias, "Counter"))
	structBuilder.AddFieldBuilder(NewFieldBuilder("opsDuration", metricsAlias, "Histogram"))

	return m
}

func (m *Model) WriteSource(w io.Writer) error {
	astFile := m.fileBuilder.Build()

	//printer.Fprint(os.Stdout, token.NewFileSet(), astFile)

	if err := format.Node(w, token.NewFileSet(), astFile); err != nil {
		return err
	}
	return nil
}

func (m *Model) AddImport(pkgName, location string) string {
	return m.fileBuilder.AddImport(pkgName, location)
}

func (m *Model) AddMethod(method *MethodConfig) error {
	mmb := NewMonitoringMethodBuilder(m.structName, method)

	mmb.SetTotalOps(m.totalOps)
	mmb.SetFailedOps(m.failedOps)
	mmb.SetOpsDuration(m.opsDuration)
	mmb.SetTimePackageAlias(m.timePackageAlias)

	m.fileBuilder.AddDeclarationBuilder(mmb)
	return nil
}

func (m *Model) createMethodBuilder(config *MethodConfig) *generator.MethodBuilder {
	mb := generator.NewMethodBuilder()
	mb.SetName(config.MethodName)
	mb.SetReceiver("m", m.structName)

	return mb
}

func (m *Model) resolveInterfaceType(location, name string) *ast.SelectorExpr {
	alias := m.AddImport("", location)
	return &ast.SelectorExpr{
		X:   ast.NewIdent(alias),
		Sel: ast.NewIdent(name),
	}
}

// MethodConfig provides the needed information for the generation
// of a stub implementation of a given method from an interface.
type MethodConfig struct {
	// MethodName specifies the name of the method as seen in the interface it
	// came from.
	MethodName string

	// MethodParams specifies all the parameters of the method.
	// They should have been normalized (i.e. no type reuse and no anonymous
	// parameters) and resolved (i.e. all selector expressions resolved against
	// the generated stub's new namespace)
	MethodParams []*ast.Field

	// MethodResults specifies all the results of the method.
	// They should have been normalized (i.e. no type reuse and no anonymous
	// results) and resolved (i.e. all selector expressions resolved against
	// the generated stub's new namespace)
	MethodResults []*ast.Field
}

func (s *MethodConfig) HasParams() bool {
	return len(s.MethodParams) > 0
}

func (s *MethodConfig) HasResults() bool {
	return len(s.MethodResults) > 0
}

// MonitoringMethodBuilder is responsible for creating a method that implements
// the original method from the interface and does all the measurement and
// recording logic.
type MonitoringMethodBuilder struct {
	method        *MethodConfig
	methodBuilder *generator.MethodBuilder

	typeName    string            // name of the type
	totalOps    *ast.SelectorExpr // selector for the struct member
	failedOps   *ast.SelectorExpr // selector for the struct member
	opsDuration *ast.SelectorExpr // selector for the struct member

	timePackageAlias string

	params  []*ast.Field // method params
	results []*ast.Field // method results
}

func NewMonitoringMethodBuilder(structName string, method *MethodConfig) *MonitoringMethodBuilder {
	mb := generator.NewMethodBuilder()
	mb.SetName(method.MethodName)
	// TODO(borshukov): Propagate type name via constructor argument.
	mb.SetReceiver("m", structName)

	return &MonitoringMethodBuilder{
		method:        method,
		methodBuilder: mb,
		params:        method.MethodParams,
		results:       method.MethodResults,
	}
}

func (b *MonitoringMethodBuilder) SetTotalOps(totalOps *ast.SelectorExpr) {
	b.totalOps = totalOps
}

func (b *MonitoringMethodBuilder) SetFailedOps(failedOps *ast.SelectorExpr) {
	b.failedOps = failedOps
}

func (b *MonitoringMethodBuilder) SetOpsDuration(opsDuration *ast.SelectorExpr) {
	b.opsDuration = opsDuration
}

func (b *MonitoringMethodBuilder) SetTimePackageAlias(alias string) {
	b.timePackageAlias = alias
}

func (b *MonitoringMethodBuilder) Build() ast.Decl {
	b.methodBuilder.SetType(&ast.FuncType{
		Params: &ast.FieldList{
			List: b.params,
		},
		Results: &ast.FieldList{
			List: util.FieldsAsAnonymous(b.results),
		},
	})

	// Add increase total operations statement
	//   m.totalOps.Add(1)
	increaseTotalOps := &CounterAddAction{counterField: b.totalOps, operationName: b.method.MethodName}
	b.methodBuilder.AddStatementBuilder(increaseTotalOps)

	// Add statement to capture current time
	//   start := time.Now()
	b.methodBuilder.AddStatementBuilder(RecordStartTime(b.timePackageAlias))

	// Add method invocation:
	//   result1, result2 := m.next.Method(arg1, arg2)
	methodInvocation := NewMethodInvocation(b.method)
	methodInvocation.SetReceiver(&ast.SelectorExpr{
		X:   ast.NewIdent("m"), // receiver name
		Sel: ast.NewIdent("next"),
	})
	b.methodBuilder.AddStatementBuilder(methodInvocation)

	// Record operation duration
	//   m.opsDuration.Observe(time.Since(start))
	b.methodBuilder.AddStatementBuilder(NewRecordOpDuraton(b.timePackageAlias, b.opsDuration))

	// Add increase failed operations statement
	//   if err != nil { m.failedOps.Add(1) }
	increaseFailedOps := NewIncreaseFailedOps(b.method, b.failedOps)
	b.methodBuilder.AddStatementBuilder(increaseFailedOps)

	// Add return statement
	//   return result1, result2
	returnResults := NewReturnResults(b.method)
	b.methodBuilder.AddStatementBuilder(returnResults)

	return b.methodBuilder.Build()
}

type CounterAddAction struct {
	counterField  *ast.SelectorExpr
	operationName string
}

func (c *CounterAddAction) Build() ast.Stmt {
	callWithExpr := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   c.counterField,
			Sel: ast.NewIdent("With"),
		},
		Args: []ast.Expr{
			&ast.BasicLit{Kind: token.STRING, Value: `"operation"`},
			&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf(`"%s"`, toSnakeCase(c.operationName))},
		},
	}

	callAddExpr := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   callWithExpr,
			Sel: ast.NewIdent("Add"),
		},
		Args: []ast.Expr{
			&ast.BasicLit{Kind: token.FLOAT, Value: "1"},
		},
	}

	return &ast.ExprStmt{
		X: callAddExpr,
	}
}

type MethodInvocation struct {
	receiver *ast.SelectorExpr
	method   *MethodConfig
}

func (m *MethodInvocation) SetReceiver(s *ast.SelectorExpr) {
	m.receiver = s
}

func NewMethodInvocation(method *MethodConfig) *MethodInvocation {
	return &MethodInvocation{method: method}
}

func (m *MethodInvocation) Build() ast.Stmt {
	resultSelectors := []ast.Expr{}
	for _, result := range m.method.MethodResults {
		resultSelectors = append(resultSelectors, ast.NewIdent(result.Names[0].String()))
	}

	paramSelectors := []ast.Expr{}
	for _, param := range m.method.MethodParams {
		paramSelectors = append(paramSelectors, ast.NewIdent(param.Names[0].String()))
	}

	callExpr := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   m.receiver,
			Sel: ast.NewIdent(m.method.MethodName),
		},
		Args: paramSelectors,
	}

	return &ast.AssignStmt{
		Lhs: resultSelectors,
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			callExpr,
		},
	}
}

type IncreaseFailedOps struct {
	method       *MethodConfig
	counterField *ast.SelectorExpr
}

func NewIncreaseFailedOps(m *MethodConfig, counterField *ast.SelectorExpr) *IncreaseFailedOps {
	return &IncreaseFailedOps{m, counterField}
}

func (i *IncreaseFailedOps) Build() ast.Stmt {
	var errorResult ast.Expr
	for _, result := range i.method.MethodResults {
		if id, ok := result.Type.(*ast.Ident); ok {
			if id.Name == "error" {
				errorResult = ast.NewIdent(result.Names[0].String())
				break
			}
		}
	}

	if errorResult == nil {
		return &ast.EmptyStmt{}
	}

	callExpr := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   i.counterField,
			Sel: ast.NewIdent("Add"),
		},
		Args: []ast.Expr{&ast.BasicLit{
			Kind:  token.FLOAT,
			Value: "1",
		}},
	}

	callStmt := &ast.ExprStmt{
		X: callExpr,
	}

	return &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  errorResult,
			Op: token.NEQ,
			Y:  ast.NewIdent("nil"),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{callStmt},
		},
	}
}

type ReturnResults struct {
	method *MethodConfig
}

func NewReturnResults(m *MethodConfig) *ReturnResults {
	return &ReturnResults{m}
}

func (r *ReturnResults) Build() ast.Stmt {
	resultSelectors := []ast.Expr{}
	for _, result := range r.method.MethodResults {
		resultSelectors = append(resultSelectors, ast.NewIdent(result.Names[0].String()))
	}

	return &ast.ReturnStmt{
		Results: resultSelectors,
	}
}

type startTimeRecorder struct {
	timePackageAlias string
}

func RecordStartTime(timePackageAlias string) *startTimeRecorder {
	return &startTimeRecorder{timePackageAlias}
}

func (r *startTimeRecorder) Build() ast.Stmt {
	callExpr := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent(r.timePackageAlias),
			Sel: ast.NewIdent("Now"),
		},
	}

	return &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent("_start")},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			callExpr,
		},
	}
}

type RecordOpDuration struct {
	timePackageAlias string
	opsDuration      *ast.SelectorExpr
}

func NewRecordOpDuraton(timePackageAlias string, opsDuration *ast.SelectorExpr) *RecordOpDuration {
	return &RecordOpDuration{
		timePackageAlias: timePackageAlias,
		opsDuration:      opsDuration,
	}
}

func (r *RecordOpDuration) Build() ast.Stmt {
	timeSinceCallExpr := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent(r.timePackageAlias),
			Sel: ast.NewIdent("Since"),
		},
		Args: []ast.Expr{ast.NewIdent("_start")},
	}

	observeCallExpr := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   r.opsDuration,
			Sel: ast.NewIdent("Observe"),
		},
		Args: []ast.Expr{timeSinceCallExpr},
	}

	return &ast.ExprStmt{X: observeCallExpr}
}

func toSnakeCase(in string) string {
	runes := []rune(in)

	var out []rune
	for i := 0; i < len(runes); i++ {
		if i > 0 && (unicode.IsUpper(runes[i]) || unicode.IsNumber(runes[i])) && ((i+1 < len(runes) && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}
