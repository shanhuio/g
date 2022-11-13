package gotags

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// tagParser contains the data needed while parsing.
type tagParser struct {
	fset *token.FileSet
	tags []Tag // list of created tags

	// all types we encounter, used to determine the constructors
	types    []string
	relative bool   // should filenames be relative to basepath
	basepath string // output file directory
}

// Parse parses the source in filename and returns a list of tags. If relative
// is true, the filenames in the list of tags are relative to basepath.
func Parse(filename string, relative bool, basepath string) ([]Tag, error) {
	p := &tagParser{
		fset:     token.NewFileSet(),
		tags:     []Tag{},
		types:    make([]string, 0),
		relative: relative,
		basepath: basepath,
	}

	f, err := parser.ParseFile(p.fset, filename, nil, 0)
	if err != nil {
		return nil, err
	}

	// package
	pkgName := p.parsePackage(f)

	// imports
	p.parseImports(f)

	// declarations
	p.parseDeclarations(f, pkgName)

	return p.tags, nil
}

// parsePackage creates a package tag.
func (p *tagParser) parsePackage(f *ast.File) string {
	p.tags = append(p.tags, p.createTag(f.Name.Name, f.Name.Pos(), Package))
	return f.Name.Name
}

// parseImports creates an import tag for each import.
func (p *tagParser) parseImports(f *ast.File) {
	for _, im := range f.Imports {
		name := strings.Trim(im.Path.Value, "\"")
		p.tags = append(p.tags, p.createTag(name, im.Path.Pos(), Import))
	}
}

// parseDeclarations creates a tag for each function, type or value
// declaration.  On function symbol we will add 2 entries in the tag file, one
// with the function name only and one with the belonging module name and the
// function name.  For method symbol we will add 3 entries: method,
// receiver.method, module.receiver.method
func (p *tagParser) parseDeclarations(f *ast.File, pkgName string) {
	// first parse the type and value declarations, so that we have a list of
	// all known types before parsing the functions.
	for _, d := range f.Decls {
		if decl, ok := d.(*ast.GenDecl); ok {
			for _, s := range decl.Specs {
				switch ts := s.(type) {
				case *ast.TypeSpec:
					p.parseTypeDeclaration(ts, pkgName)
				case *ast.ValueSpec:
					p.parseValueDeclaration(ts, pkgName)
				}
			}
		}
	}

	// now parse all the functions
	for _, d := range f.Decls {
		if decl, ok := d.(*ast.FuncDecl); ok {
			p.parseFunction(decl, pkgName)
		}
	}
}

// parseFunction creates a tag for function declaration f.
func (p *tagParser) parseFunction(f *ast.FuncDecl, pkgName string) {
	tag := p.createTag(f.Name.Name, f.Pos(), Function)

	tag.Fields[Access] = getAccess(tag.Name)
	tag.Fields[Signature] = fmt.Sprintf("(%s)", getTypes(f.Type.Params, true))
	tag.Fields[TypeField] = getTypes(f.Type.Results, false)

	if f.Recv != nil && len(f.Recv.List) > 0 {
		// this function has a receiver, set the type to Method
		tag.Fields[ReceiverType] = getType(f.Recv.List[0].Type, false)
		tag.Type = Method
	} else if name, ok := p.belongsToReceiver(f.Type.Results); ok {
		// this function does not have a receiver, but it belongs to one based
		// on its return values; its type will be Function instead of Method.
		tag.Fields[ReceiverType] = name
		tag.Type = Function
	}

	p.tags = append(p.tags, tag)
}

// parseTypeDeclaration creates a tag for type declaration ts and for each
// field in case of a struct, or each method in case of an interface.
// The pkgName argument holds the name of the package the file currently parsed
// belongs to.
func (p *tagParser) parseTypeDeclaration(ts *ast.TypeSpec, pkgName string) {
	tag := p.createTag(ts.Name.Name, ts.Pos(), Type)

	tag.Fields[Access] = getAccess(tag.Name)

	switch s := ts.Type.(type) {
	case *ast.StructType:
		tag.Fields[TypeField] = "struct"
		p.parseStructFields(tag.Name, s)
		p.types = append(p.types, tag.Name)
	case *ast.InterfaceType:
		tag.Fields[TypeField] = "interface"
		tag.Type = Interface
		p.parseInterfaceMethods(tag.Name, s)
	default:
		tag.Fields[TypeField] = getType(ts.Type, true)
	}

	p.tags = append(p.tags, tag)
}

// parseValueDeclaration creates a tag for each variable or constant
// declaration, unless the declaration uses a blank identifier.
func (p *tagParser) parseValueDeclaration(v *ast.ValueSpec, pkgName string) {
	for _, d := range v.Names {
		if d.Name == "_" {
			continue
		}

		tag := p.createTag(d.Name, d.Pos(), Variable)
		tag.Fields[Access] = getAccess(tag.Name)

		if v.Type != nil {
			tag.Fields[TypeField] = getType(v.Type, true)
		}

		switch d.Obj.Kind {
		case ast.Var:
			tag.Type = Variable
		case ast.Con:
			tag.Type = Constant
		}
		p.tags = append(p.tags, tag)
	}
}

// parseStructFields creates a tag for each field in struct s, using name as the
// tags ctype.
func (p *tagParser) parseStructFields(name string, s *ast.StructType) {
	for _, f := range s.Fields.List {
		var tag Tag
		if len(f.Names) > 0 {
			for _, n := range f.Names {
				tag = p.createTag(n.Name, n.Pos(), Field)
				tag.Fields[Access] = getAccess(tag.Name)
				tag.Fields[ReceiverType] = name
				tag.Fields[TypeField] = getType(f.Type, true)
				p.tags = append(p.tags, tag)
			}
		} else {
			// embedded field
			tag = p.createTag(getType(f.Type, true), f.Pos(), Embedded)
			tag.Fields[Access] = getAccess(tag.Name)
			tag.Fields[ReceiverType] = name
			tag.Fields[TypeField] = getType(f.Type, true)
			p.tags = append(p.tags, tag)
		}
	}
}

// parseInterfaceMethods creates a tag for each method in interface s, using
// name as the tags ctype.
func (p *tagParser) parseInterfaceMethods(name string, s *ast.InterfaceType) {
	for _, f := range s.Methods.List {
		var tag Tag
		if len(f.Names) > 0 {
			tag = p.createTag(f.Names[0].Name, f.Names[0].Pos(), Method)
		} else {
			// embedded interface
			tag = p.createTag(getType(f.Type, true), f.Pos(), Embedded)
		}

		tag.Fields[Access] = getAccess(tag.Name)

		if t, ok := f.Type.(*ast.FuncType); ok {
			tag.Fields[Signature] = fmt.Sprintf(
				"(%s)", getTypes(t.Params, true),
			)
			tag.Fields[TypeField] = getTypes(t.Results, false)
		}

		tag.Fields[InterfaceType] = name

		p.tags = append(p.tags, tag)
	}
}

// createTag creates a new tag, using pos to find the filename and set the line
// number.
func (p *tagParser) createTag(name string, pos token.Pos, tagType TagType) Tag {
	f := p.fset.File(pos).Name()
	if p.relative {
		if abs, err := filepath.Abs(f); err != nil {
			fmt.Fprintf(os.Stderr, "fail to determine absolute path: %s\n", err)
		} else if rel, err := filepath.Rel(p.basepath, abs); err != nil {
			fmt.Fprintf(os.Stderr, "fail to determine relative path: %s\n", err)
		} else {
			f = rel
		}
	}
	return NewTag(name, f, p.fset.Position(pos).Line, tagType)
}

// belongsToReceiver checks if a function with these return types belongs to
// a receiver. If it belongs to a receiver, the name of that receiver will be
// returned with ok set to true. Otherwise ok will be false.
// Behavior should be similar to how go doc decides when a function belongs to
// a receiver (gosrc/pkg/go/doc/reader.go).
func (p *tagParser) belongsToReceiver(types *ast.FieldList) (string, bool) {
	if types == nil || types.NumFields() == 0 {
		return "", false
	}

	// If the first return type has more than 1 result associated with
	// it, it should not belong to that receiver.
	// Similar behavior as go doc (go source/.
	if len(types.List[0].Names) > 1 {
		return "", false
	}

	// get name of the first return type
	t := getType(types.List[0].Type, false)

	// check if it exists in the current list of known types
	for _, knownType := range p.types {
		if t == knownType {
			return knownType, true
		}
	}

	return "", false
}
