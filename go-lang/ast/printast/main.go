package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"github.com/ironzhang/tlog"
)

type Options struct {
	Input  string
	Output string
}

func (p *Options) Setup() {
	flag.StringVar(&p.Input, "input", "", "input file")
	flag.StringVar(&p.Output, "output", "", "output file")
	flag.Parse()

	if p.Input == "" {
		flag.Usage()
		os.Exit(-1)
	}
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stdout, "printast FILENAME [FILENAME...]")
		return
	}

	for _, filename := range os.Args[1:] {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			tlog.Errorw("parse file", "filename", filename, "error", err)
			return
		}

		fmt.Fprintf(os.Stdout, "%s:\n", filename)
		ast.Fprint(os.Stdout, fset, f, ast.NotNilFilter)
	}
}
