// Package goparse contains logic for parsing Go files. Specifically it parses
// source and test files into domain models for generating tests.
package goparser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"strings"

	"github.com/cweill/gotests/internal/models"
	"golang.org/x/tools/go/packages"
)

// ErrEmptyFile represents an empty file error.
var ErrEmptyFile = errors.New("file is empty")

// Result representats a parsed Go file.
type Result struct {
	// The package name and imports of a Go file.
	Header *models.Header
	// All the functions and methods in a Go file.
	Funcs      []*models.Function
	InterFaces map[string][]*models.Interfaces
}

// Parser can parse Go files.
type Parser struct {
	// The importer to resolve packages from import paths.
	Importer types.Importer
}

// Parse parses a given Go file at srcPath, along any files that share the same
// package, into a domain model for generating tests.
func (p *Parser) Parse(srcPath string, files []models.Path) (*Result, error) {
	b, err := p.readFile(srcPath)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	f, err := p.parseFile(fset, srcPath)
	if err != nil {
		return nil, err
	}
	fs, err := p.parseFiles(fset, f, files)
	if err != nil {
		return nil, err
	}
	interFaces := p.parseInterface(f.Imports)
	funcs := p.parseFunctions(fset, f, fs, interFaces)
	return &Result{
		Header: &models.Header{
			Comments: parsePkgComment(f, f.Package),
			Package:  f.Name.String(),
			Imports:  parseImports(f.Imports, funcs),
			Code:     goCode(b, f),
		},
		InterFaces: interFaces,
		Funcs:      funcs,
	}, nil
}

func (p *Parser) parseInterface(imps []*ast.ImportSpec) map[string][]*models.Interfaces {
	r := make(map[string][]*models.Interfaces)
	var pathImported []string
	for _, i := range imps {
		if i.Name != nil {
			//fmt.Println("name", i.Name.Name)
		}
		if i.Path != nil {
			if i.Path.Value == "\"fmt\"" {
				continue
			}
			//fmt.Println("path", i.Path.Value)
			pathImported = append(pathImported, strings.Trim(i.Path.Value, "\""))
		}
	}
	if len(pathImported) < 1 {
		return r
	}
	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax | packages.NeedName}
	//fmt.Println(pathImported)
	pkgs, err := packages.Load(cfg, pathImported...)
	if err != nil {
		fmt.Println(err)
	}
	for _, pkg := range pkgs {
		//	fmt.Println(pkg.ID, pkg.Name, pkg.GoFiles)
		for _, f := range pkg.Syntax {
			for _, decl := range f.Decls {
				//fmt.Println(decl)
				interfaces := getInterface(decl, pkg.ID, pkg.Name, pkg.GoFiles)
				r[pkg.ID] = append(r[pkg.ID], interfaces...)
			}
		}
	}
	return r
}

func getInterface(x ast.Decl, importPath, pkgName string, filePaths []string) []*models.Interfaces {
	var r []*models.Interfaces
	if x, ok := x.(*ast.GenDecl); ok {
		if x.Tok != token.TYPE {
			return r
		}
		for _, x := range x.Specs {
			if x, ok := x.(*ast.TypeSpec); ok {
				iname := x.Name
				if x, ok := x.Type.(*ast.InterfaceType); ok {
					for _, x := range x.Methods.List {
						if len(x.Names) == 0 {
							return r
						}
						//					mname := x.Names[0].Name
						//					fmt.Println("interface:", iname, "method:", mname)
					}
					r = append(r, &models.Interfaces{
						PkgName:    pkgName,
						FilePath:   filePaths,
						ImportPath: importPath,
						Name:       iname.Name,
						MethodList: x.Methods.List,
					})
				}
			}
		}
	}
	return r
}

func (p *Parser) readFile(srcPath string) ([]byte, error) {
	b, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadFile: %v", err)
	}
	if len(b) == 0 {
		return nil, ErrEmptyFile
	}
	return b, nil
}

func (p *Parser) parseFile(fset *token.FileSet, srcPath string) (*ast.File, error) {
	f, err := parser.ParseFile(fset, srcPath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("target parser.ParseFile(): %v", err)
	}
	return f, nil
}

func (p *Parser) parseFiles(fset *token.FileSet, f *ast.File, files []models.Path) ([]*ast.File, error) {
	pkg := f.Name.String()
	var fs []*ast.File
	for _, file := range files {
		ff, err := parser.ParseFile(fset, string(file), nil, 0)
		if err != nil {
			return nil, fmt.Errorf("other file parser.ParseFile: %v", err)
		}
		if name := ff.Name.String(); name != pkg {
			continue
		}
		fs = append(fs, ff)
	}
	return fs, nil
}

func (p *Parser) parseFunctions(fset *token.FileSet, f *ast.File, fs []*ast.File, interfaces map[string][]*models.Interfaces) []*models.Function {
	ul, el := p.parseTypes(fset, fs)
	var funcs []*models.Function
	for _, d := range f.Decls {
		fDecl, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}
		funcs = append(funcs, parseFunc(fDecl, ul, el, interfaces))
	}
	return funcs
}

func (p *Parser) parseTypes(fset *token.FileSet, fs []*ast.File) (map[string]types.Type, map[*types.Struct]ast.Expr) {
	conf := &types.Config{
		Importer: p.Importer,
		// Adding a NO-OP error function ignores errors and performs best-effort
		// type checking. https://godoc.org/golang.org/x/tools/go/types#Config
		Error: func(error) {},
	}
	ti := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}
	// Note: conf.Check can fail, but since Info is not required data, it's ok.
	conf.Check("", fset, fs, ti)
	ul := make(map[string]types.Type)
	el := make(map[*types.Struct]ast.Expr)
	for e, t := range ti.Types {
		// Collect the underlying types.
		ul[t.Type.String()] = t.Type.Underlying()
		// Collect structs to determine the fields of a receiver.
		if v, ok := t.Type.(*types.Struct); ok {
			el[v] = e
		}
	}
	return ul, el
}

func parsePkgComment(f *ast.File, pkgPos token.Pos) []string {
	var comments []string
	var count int

	for _, comment := range f.Comments {

		if comment.End() >= pkgPos {
			break
		}
		for _, c := range comment.List {
			count += len(c.Text) + 1 // +1 for '\n'
			if count < int(c.End()) {
				n := int(c.End()) - count - 1
				comments = append(comments, strings.Repeat("\n", n))
				count++ // for last of '\n'
			}
			comments = append(comments, c.Text)
		}
	}

	if int(pkgPos)-count > 1 {
		comments = append(comments, strings.Repeat("\n", int(pkgPos)-count-2))
	}
	return comments
}

// Returns the Go code below the imports block.
func goCode(b []byte, f *ast.File) []byte {
	furthestPos := f.Name.End()
	for _, node := range f.Imports {
		if pos := node.End(); pos > furthestPos {
			furthestPos = pos
		}
	}
	if furthestPos < token.Pos(len(b)) {
		furthestPos++

		// Avoid wrong output on windows-encoded files
		if b[furthestPos-2] == '\r' && b[furthestPos-1] == '\n' && furthestPos < token.Pos(len(b)) {
			furthestPos++
		}
	}
	return b[furthestPos:]
}

func parseFunc(fDecl *ast.FuncDecl, ul map[string]types.Type, el map[*types.Struct]ast.Expr, interfaces map[string][]*models.Interfaces) *models.Function {
	f := &models.Function{
		Name:       fDecl.Name.String(),
		IsExported: fDecl.Name.IsExported(),
		Receiver:   parseReceiver(fDecl.Recv, ul, el, interfaces),
		Parameters: parseFieldList(fDecl.Type.Params, ul, interfaces),
	}
	fs := parseFieldList(fDecl.Type.Results, ul, interfaces)
	i := 0
	for _, fi := range fs {
		if fi.Type.String() == "error" {
			f.ReturnsError = true
			continue
		}
		fi.Index = i
		f.Results = append(f.Results, fi)
		i++
	}
	return f
}

func parseImports(imps []*ast.ImportSpec, funcs []*models.Function) []*models.Import {
	var is []*models.Import
	for _, imp := range imps {
		var n string
		if imp.Name != nil {
			n = imp.Name.String()
		}
		isInterface := false
		for _, f := range funcs {
			if f.Receiver == nil {
				continue
			}
			for _, field := range f.Receiver.Fields {
				if field.InterfacesInfo == nil {
					continue
				}
				//fmt.Println(field.InterfacesInfo.PkgName == n, n, field.InterfacesInfo.PkgName, field.InterfacesInfo.PkgPath, field.InterfacesInfo.PkgPath == strings.Trim(imp.Path.Value, "\""), imp.Path.Value)
				if field.InterfacesInfo.PkgPath == strings.Trim(imp.Path.Value, "\"") {
					isInterface = true
					goto OUT
				}
			}
		}
	OUT:
		path := strings.Trim(imp.Path.Value, "\"")
		paths := strings.Split(path, "/")
		tname := ""
		if len(paths) > 0 {
			tname = paths[len(paths)-1]
		}
		//fmt.Println(n, imp.Path.Value, isInterface)
		is = append(is, &models.Import{
			Name:        n,
			TailName:    tname,
			Path:        path,
			IsInterface: isInterface,
			IsEmpty:     path == "",
		})
	}
	return is
}

func parseReceiver(fl *ast.FieldList, ul map[string]types.Type, el map[*types.Struct]ast.Expr, interfaces map[string][]*models.Interfaces) *models.Receiver {
	if fl == nil {
		return nil
	}
	r := &models.Receiver{
		Field: parseFieldList(fl, ul, interfaces)[0],
	}
	t, ok := ul[r.Type.Value]
	if !ok {
		return r
	}
	s, ok := t.(*types.Struct)
	if !ok {
		return r
	}
	st, found := el[s]
	if !found {
		return r
	}
	r.Fields = append(r.Fields, parseFieldList(st.(*ast.StructType).Fields, ul, interfaces)...)
	for i, f := range r.Fields {
		// https://github.com/cweill/gotests/issues/69
		if i >= s.NumFields() {
			break
		}
		f.Name = s.Field(i).Name()
	}
	return r

}

func parseFieldList(fl *ast.FieldList, ul map[string]types.Type, interfaces map[string][]*models.Interfaces) []*models.Field {
	if fl == nil {
		return nil
	}
	i := 0
	var fs []*models.Field
	for _, f := range fl.List {
		for _, pf := range parseFields(f, ul, interfaces) {
			pf.Index = i
			fs = append(fs, pf)
			i++
		}
	}
	return fs
}

func parseFields(f *ast.Field, ul map[string]types.Type, interfaces map[string][]*models.Interfaces) []*models.Field {
	t := parseExpr(f.Type, ul)
	if len(f.Names) == 0 {
		return []*models.Field{{
			Type: t,
		}}
	}
	var fs []*models.Field
	for _, n := range f.Names {
		//fmt.Println(n.Name, t)
		interfacesInfo := getInterfaceInfo(t, interfaces)

		//fmt.Println("interface :", interfacesInfo != nil)
		fs = append(fs, &models.Field{
			Name:           n.Name,
			Type:           t,
			IsInterface:    interfacesInfo != nil,
			InterfacesInfo: interfacesInfo,
		})
	}
	return fs
}

func isBasicOrSrcType(t string) bool {
	switch t {
	case "bool", "string", "int", "int8", "int16", "int32", "int64", "uint",
		"uint8", "uint16", "uint32", "uint64", "uintptr", "byte", "rune",
		"float32", "float64", "complex64", "complex128", "error":
		return true
	default:
		return false
	}
}

func getInterfaceInfo(t *models.Expression, interfaces map[string][]*models.Interfaces) *models.InterfaceTypeInfo {
	if t.Value == "" {
		return nil
	}
	v := strings.Split(strings.Trim(t.Value, " "), ".")
	if len(v) < 2 || v[0] == "" || v[1] == "" {
		//@TODO current package
		return nil
	}
	for _, ifs := range interfaces {
		for _, intf := range ifs {
			if intf.Name == v[1] && intf.PkgName == v[0] {
				//fmt.Println("type raw info:", t.Value, v, ifs, interfaces)

				var ml []*models.MethodInfo
				for _, x := range intf.MethodList {
					if len(x.Names) == 0 {
						continue
					}
					mname := x.Names[0].Name
					//fmt.Printf("%#v", x.Type)

					var params, results []string
					if fn, ok := x.Type.(*ast.FuncType); ok {
						//fmt.Printf("%#v  %#v", fn.Params, fn.Results)
						for _, p := range fn.Params.List {
							ll := 1
							if len(p.Names) > 1 {
								ll = len(p.Names)
							}
							for i := 0; i < ll; i++ {
								switch tp := p.Type.(type) {
								case *ast.ArrayType:
									pkg := intf.PkgName + "."
									if isBasicOrSrcType(tp.Elt.(*ast.Ident).Name) {
										pkg = ""
									}
									params = append(params, "[]"+pkg+tp.Elt.(*ast.Ident).Name)
								case *ast.Ident:
									pkg := intf.PkgName + "."
									if isBasicOrSrcType(tp.Name) {
										pkg = ""
									}
									params = append(params, pkg+tp.Name)
								}
							}
						}
						for _, r := range fn.Results.List {
							ll := 1
							if len(r.Names) > 1 {
								ll = len(r.Names)
							}
							for i := 0; i < ll; i++ {
								switch tp := r.Type.(type) {
								case *ast.ArrayType:
									pkg := intf.PkgName + "."
									if isBasicOrSrcType(tp.Elt.(*ast.Ident).Name) {
										pkg = ""
									}
									results = append(results, "[]"+pkg+tp.Elt.(*ast.Ident).Name)
								case *ast.Ident:
									pkg := intf.PkgName + "."
									if isBasicOrSrcType(tp.Name) {
										pkg = ""
									}
									results = append(results, pkg+tp.Name)
								}
							}
						}
					}
					//fmt.Println("func decsribe:", params, results)
					ml = append(ml, &models.MethodInfo{
						Name:   mname,
						Params: params,
						Return: results,
					})
				}
				return &models.InterfaceTypeInfo{
					PkgName:        v[0],
					Name:           v[1],
					PkgPath:        intf.ImportPath,
					MethodList:     intf.MethodList,
					MethodInfoList: ml,
				}
			}
		}
	}
	return nil
}

func parseExpr(e ast.Expr, ul map[string]types.Type) *models.Expression {
	switch v := e.(type) {
	case *ast.StarExpr:
		val := types.ExprString(v.X)
		return &models.Expression{
			Value:      val,
			IsStar:     true,
			Underlying: underlying(val, ul),
		}
	case *ast.Ellipsis:
		exp := parseExpr(v.Elt, ul)
		return &models.Expression{
			Value:      exp.Value,
			IsStar:     exp.IsStar,
			IsVariadic: true,
			Underlying: underlying(exp.Value, ul),
		}
	default:
		val := types.ExprString(e)
		return &models.Expression{
			Value:      val,
			Underlying: underlying(val, ul),
			IsWriter:   val == "io.Writer",
		}
	}
}

func underlying(val string, ul map[string]types.Type) string {
	if ul[val] != nil {
		return ul[val].String()
	}
	return ""
}
