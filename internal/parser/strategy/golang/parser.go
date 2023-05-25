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
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type parser struct {
	spec                *sloth.Spec
	sourceFile          string
	sourceContent       io.ReadCloser
	includedDirs        []string
	applicationPackages map[string]*ast.Package
	logger              *logging.Logger
}

// newParser client parser performs all checks at initialization time
func newParser(logger *logging.Logger, sourceFile string,
	sourceContent io.ReadCloser, dirs ...string) *parser {
	pkgs := map[string]*ast.Package{}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			// skip if dir doesn't exists
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
		spec: &sloth.Spec{
			Version: sloth.Version,
			Service: "",
			Labels:  nil,
			SLOs:    nil,
		},
		sourceFile:          sourceFile,
		sourceContent:       sourceContent,
		includedDirs:        dirs,
		applicationPackages: pkgs,
		logger:              logger,
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

func getFile(name string, file io.ReadCloser) (*ast.File, error) {
	fset := token.NewFileSet()
	if file != nil {
		defer file.Close()
	}
	return goparser.ParseFile(fset, name, file, goparser.ParseComments)
}

func (p parser) Parse(ctx context.Context) (*sloth.Spec, error) {
	// collect all aloe error comments from packages and add them to the spec struct
	if p.sourceFile != "" || p.sourceContent != nil {
		file, err := getFile(p.sourceFile, p.sourceContent)
		if err != nil {
			return nil, err
		}
		if err := p.parseComments(file.Comments...); err != nil {
			return nil, err
		}
		return p.spec, nil
	}

	if len(p.applicationPackages) > 0 {
		for _, pkg := range p.applicationPackages {
			for _, file := range pkg.Files {
				if err := p.parseComments(file.Comments...); err != nil {
					p.logger.Info(err.Error())
					continue
				}
			}
		}
	}

	return p.spec, nil
}

func (p parser) parseComments(comments ...*ast.CommentGroup) error {
	for _, comment := range comments {
		newSpec, err := grammar.Eval(strings.TrimSpace(comment.Text()))
		switch {
		case errors.Is(err, grammar.ErrParseSource):
			continue
		case err != nil:
			p.logger.Error(err, "")
			continue
		}

		if p.spec.Service == "" {
			p.spec.Service = newSpec.Service
		}
		if p.spec.Version == "" {
			p.spec.Version = newSpec.Version
		}
		if p.spec.Labels == nil {
			p.spec.Labels = newSpec.Labels
		}

		for _, slo := range newSpec.SLOs {
			if slo.Name != "" {
				p.spec.SLOs = append(p.spec.SLOs, slo)
			}
		}
	}
	return nil
}
