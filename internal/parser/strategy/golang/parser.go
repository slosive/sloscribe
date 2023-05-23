package golang

import (
	"context"
	"github.com/juju/errors"
	sloth "github.com/slok/sloth/pkg/prometheus/api/v1"
	"github.com/tfadeyi/sloth-simple-comments/internal/logging"
	"github.com/tfadeyi/sloth-simple-comments/internal/parser/grammar"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type parser struct {
	Spec                *sloth.Spec
	GeneralInfoSource   string
	IncludedDirs        []string
	ApplicationPackages map[string]*ast.Package
	Logger              *logging.Logger
}

const (
	defaultSourceFile = "main.go"
)

// newParser client parser performs all checks at initialization time
func newParser(logger *logging.Logger, dirs ...string) *parser {
	pkgs := map[string]*ast.Package{}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			//skip if dir doesn't exists
			continue
		}

		foundPkgs, err := getPackages(dir)
		if err != nil {
			logger.Info(err.Error())
			continue
		}

		for pkgName, pkg := range foundPkgs {
			if _, ok := pkgs[pkgName]; !ok {
				pkgs[pkgName] = pkg
			}
		}

	}

	return &parser{
		Spec: &sloth.Spec{
			Version: sloth.Version,
			Service: "",
			Labels:  nil,
			SLOs:    nil,
		},
		GeneralInfoSource:   defaultSourceFile,
		IncludedDirs:        dirs,
		ApplicationPackages: pkgs,
		Logger:              logger,
	}
}

func getPackages(dir string) (map[string]*ast.Package, error) {
	fset := token.NewFileSet()
	pkgs, err := goparser.ParseDir(fset, dir, nil, goparser.ParseComments)
	if err != nil {
		return map[string]*ast.Package{}, err
	}

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			foundPkgs, err := goparser.ParseDir(fset, path, nil, goparser.ParseComments)
			if err != nil {
				return err
			}
			for pkgName, pkg := range foundPkgs {
				if _, ok := pkgs[pkgName]; !ok {
					pkgs[pkgName] = pkg
				}
			}
		}
		return err
	})
	return pkgs, err
}

func (p parser) Parse(ctx context.Context) (*sloth.Spec, error) {
	// collect all aloe error comments from packages and add them to the spec struct
	for _, pkg := range p.ApplicationPackages {
		for _, file := range pkg.Files {
			if err := p.parseComments(file.Comments...); err != nil {
				p.Logger.Info(err.Error())
				continue
			}
		}
	}

	return p.Spec, nil
}

func (p parser) parseComments(comments ...*ast.CommentGroup) error {
	for _, comment := range comments {
		newSpec, err := grammar.Eval(strings.TrimSpace(comment.Text()))
		switch {
		case errors.Is(err, grammar.ErrParseSource):
			continue
		case err != nil:
			p.Logger.Error(err, "")
			continue
		}

		if p.Spec.Service == "" {
			p.Spec.Service = newSpec.Service
		}
		if p.Spec.Version == "" {
			p.Spec.Version = newSpec.Version
		}
		if p.Spec.Labels == nil {
			p.Spec.Labels = newSpec.Labels
		}

		for _, slo := range newSpec.SLOs {
			if slo.Name != "" {
				p.Spec.SLOs = append(p.Spec.SLOs, slo)
			}
		}
	}
	return nil
}
