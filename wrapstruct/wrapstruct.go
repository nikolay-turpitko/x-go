package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/types"
	"html/template"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/ernesto-jimenez/gogen/imports"

	"golang.org/x/tools/go/loader"
)

func main() {
	prevUsage := flag.Usage
	flag.Usage = func() {
		fmt.Println(`
wrapstruct generates wrapper around provided struct
to hide fields but to provide getters instead.
It can be used in rare cases when you have a struct at hand,
but need to provide an interface from it, having methods with same
names and returning types as fileds of original structure.
WARNING: automates creation of non-idiomatic code.
`)
		prevUsage()
	}
	out := flag.String("o", "", "what file to write")
	src := flag.String("src", "", "source type (struct)")
	dst := flag.String("dst", "", "destination type (wrapper)")
	srcpkgpath := flag.String("srcpkg", ".", "source type's (struct) package")
	dstpkgpath := flag.String("dstpkg", ".", "destination type (wrapper) package")
	tag := flag.String("tag", "wrapstruct", "destination type (wrapper) package")
	flag.Parse()
	if *src == "" {
		log.Fatal("need to specify a struct name (src)")
	}
	if *dst == "" {
		log.Fatal("need to specify a type name (dst)")
	}
	*srcpkgpath = cleanPkgpath(*srcpkgpath)
	*dstpkgpath = cleanPkgpath(*dstpkgpath)
	log.Printf(`Generating wrapper "%s".%s for  "%s".%s`, *dstpkgpath, *dst, *srcpkgpath, *src)
	srcpkg := loadPkg(*srcpkgpath)
	dstpkg := loadPkg(*dstpkgpath)
	srctype := structType(srcpkg, *src)

	c := context{
		SrcName: *src,
		DstPkg:  dstpkg,
		DstName: *dst,
		SrcType: srctype,
	}
	c.init(*tag)

	var ow io.Writer
	if *out == "" {
		ow = os.Stdout
	} else {
		fw, err := os.Create(*out)
		if err != nil {
			log.Fatal(err)
		}
		defer fw.Close()
		ow = fw
	}
	w := bufio.NewWriter(ow)
	defer w.Flush()
	fileTmpl.Execute(w, c)
}

func cleanPkgpath(path string) string {
	if path != "" && path[0] != '.' {
		return path
	}
	var err error
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}
	path, err = filepath.Abs(filepath.Clean(path))
	if err != nil {
		log.Fatal(err)
	}
	return removeGopath(path)
}

func removeGopath(p string) string {
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
		p = strings.Replace(p, path.Join(gopath, "src")+"/", "", 1)
	}
	return p
}

func loadPkg(pkgpath string) *types.Package {
	// Ignore all type checker errors, because tool can be used on unfinished code
	conf := loader.Config{AllowErrors: true}
	conf.TypeChecker.Error = func(error) {}
	conf.Import(pkgpath)
	lprog, err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}
	return lprog.Package(pkgpath).Pkg
}

func structType(pkg *types.Package, name string) *types.Struct {
	obj := pkg.Scope().Lookup(name)
	if obj == nil {
		log.Fatalf("%s.%s not found", pkg.Path(), name)
	}

	typ := obj.Type().Underlying()
	strtyp, ok := typ.(*types.Struct)
	if !ok {
		log.Fatalf("%s should be an struct, was %s", name, typ)
	}
	return strtyp
}

type context struct {
	SrcName string
	DstPkg  *types.Package
	DstName string
	SrcType *types.Struct
	Fields  []field
	Imports map[string]string
}

func (c *context) init(tag string) {
	c.initFields(tag)
	c.initImports()
}

func (c *context) initFields(tag string) {
	s := c.SrcType
	numFields := s.NumFields()
	var fields []field
	for i := 0; i < numFields; i++ {
		n := reflect.StructTag(s.Tag(i)).Get(tag)
		if n != "-" {
			f := field{n, s.Field(i), c.DstPkg.Name()}
			fields = append(fields, f)
		}
	}
	c.Fields = fields
}

func (c *context) initImports() {
	imports := imports.New(c.DstPkg.String())
	fields := c.Fields
	for _, f := range fields {
		imports.AddImportsFrom(f.Field.Type())
		imports.AddImportsFrom(f.UnderlyingType())
	}
	c.Imports = imports.Imports()
	delete(c.Imports, c.DstPkg.Path())
}

type field struct {
	nameFromTag string
	Field       *types.Var
	DstPkgName  string
}

func (f field) Name() string {
	if f.nameFromTag != "" {
		return f.nameFromTag
	}
	return strings.Title(f.Field.Name())
}

func (f field) Type() string {
	return types.TypeString(f.Field.Type(), func(*types.Package) string {
		s := f.Field.Type().String()
		i := strings.LastIndex(s, ".")
		if i > 0 {
			s = s[:i]
		}
		i = strings.LastIndex(s, "/")
		if i > 0 {
			s = s[i+1:]
		}
		s = strings.TrimPrefix(s, f.DstPkgName)
		return s
	})
}

func (f field) UnderlyingType() types.Type {
	switch t := f.Field.Type().(type) {
	case *types.Array:
		switch t := t.Elem().(type) {
		case *types.Pointer:
			return t.Elem()
		}
		return t.Elem()
	case *types.Slice:
		switch t := t.Elem().(type) {
		case *types.Pointer:
			return t.Elem()
		}
		return t.Elem()
	}
	return nil
}

var (
	fileTmpl = template.Must(template.New("file").Parse(`
/*
* CODE GENERATED AUTOMATICALLY WITH github.com/nikolay-turpitko/x-go/wrapstruct
* THIS FILE SHOULD NOT BE EDITED BY HAND
*/

package {{.DstPkg.Name}}
{{ if .Imports }}
import (
{{- range $path, $name := .Imports}}
	"{{$path}}"{{end}}
)
{{ end }}
type {{ .DstName }} struct {
	w *{{ .SrcName }}
}
{{ range .Fields }}
func (w *{{ $.DstName }}) {{ .Name }}() {{ .Type }} {
	return w.w.{{ .Field.Name }}
}
{{ end }}
`))
)
