package golang

import (
	"context"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/juju/errors"
	sloth "github.com/slok/sloth/pkg/prometheus/api/v1"
	"github.com/tfadeyi/slotalk/internal/logging"
	"github.com/tfadeyi/slotalk/internal/parser/specification/sloth/grammar"
)

type parser struct {
	specs         map[string]*sloth.Spec
	sourceFile    string
	sourceContent io.ReadCloser
	includedDirs  []string
	logger        *logging.Logger
}

// Options contains the configuration options available to the Parser
type Options struct {
	Logger           *logging.Logger
	SourceFile       string
	SourceContent    io.ReadCloser
	InputDirectories []string
}

func NewOptions() *Options {
	l := logging.NewStandardLogger()
	return &Options{
		Logger:           &l,
		SourceFile:       "",
		SourceContent:    nil,
		InputDirectories: nil,
	}
}

// NewParser client parser performs all checks at initialization time
func NewParser(opts *Options) *parser {
	// create default options, these will be overridden
	if opts == nil {
		opts = NewOptions()
	}

	logger := opts.Logger
	dirs := opts.InputDirectories
	sourceFile := opts.SourceFile
	sourceContent := opts.SourceContent

	return &parser{
		specs:         map[string]*sloth.Spec{},
		sourceFile:    sourceFile,
		sourceContent: sourceContent,
		includedDirs:  dirs,
		logger:        logger,
	}
}

// getAllGoPackages fetches all the available golang packages in the target directory and subdirectories
func getAllGoPackages(dir string) (map[string]*ast.Package, error) {
	fset := token.NewFileSet()
	pkgs, err := goparser.ParseDir(fset, dir, nil, goparser.ParseComments)
	if err != nil {
		return map[string]*ast.Package{}, err
	}

	// walk through the directories and parse the not already found go packages
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
	if err != nil {
		return nil, err
	}

	if len(pkgs) == 0 {
		return nil, errors.Errorf("no go packages were found in the target directory and subdirectories: %s", dir)
	}

	return pkgs, nil
}

// getFile returns the ast go file struct given filename or an io.Reader. If an io.Reader is passed it will take precedence
// over the filename
func getFile(name string, file io.ReadCloser) (*ast.File, error) {
	fset := token.NewFileSet()
	if file != nil {
		defer file.Close()
	}
	return goparser.ParseFile(fset, name, file, goparser.ParseComments)
}

// parseSlothAnnotations parses the source code comments for sloth annotations using the sloth grammar.
// It expects only SLO definition per comment group
func (p parser) parseSlothAnnotations(comments ...*ast.CommentGroup) error {
	var currentServiceSpec *sloth.Spec

	for _, comment := range comments {
		// partialServiceSpec contains the partially parsed sloth Specification for a given comment group
		// this means the parsed spec will only contain data for the fields that are present in the comments, making the spec only partially accurate
		partialServiceSpec, err := grammar.Eval(strings.TrimSpace(comment.Text()))
		if err != nil {
			p.warn(err)
			continue
		}

		// if the comment group contains a reference to the service name
		// check if service was parsed before else add it the collection of specs.
		// Set the found service spec as the current service spec.
		if partialServiceSpec.Service != "" {
			if currentServiceSpec != nil && (currentServiceSpec.Service == partialServiceSpec.Service || currentServiceSpec.Service == "") {
				p.specs[partialServiceSpec.Service] = currentServiceSpec
			}
			spec, ok := p.specs[partialServiceSpec.Service]
			if !ok {
				p.specs[partialServiceSpec.Service] = partialServiceSpec
				currentServiceSpec = partialServiceSpec
			} else {
				currentServiceSpec = spec
			}
		}

		if currentServiceSpec == nil {
			currentServiceSpec = &sloth.Spec{
				Version: "",
				Service: "",
				Labels:  nil,
				SLOs:    nil,
			}
		}

		if currentServiceSpec.Service == "" {
			currentServiceSpec.Service = partialServiceSpec.Service
		}
		if currentServiceSpec.Version == "" {
			currentServiceSpec.Version = partialServiceSpec.Version
		}
		if currentServiceSpec.Labels == nil {
			for key, label := range partialServiceSpec.Labels {
				currentServiceSpec.Labels[key] = label
			}
		}
		if currentServiceSpec.SLOs == nil {
			currentServiceSpec.SLOs = append(currentServiceSpec.SLOs, partialServiceSpec.SLOs...)
		}
	}
	return nil
}

// Parse will parse the source code for sloth annotations.
// In case of error during parsing, Parse returns an empty sloth.Spec
func (p parser) Parse(ctx context.Context) (map[string]*sloth.Spec, error) {
	// collect all sloth annotations from the file and add them to the spec struct
	if p.sourceFile != "" || p.sourceContent != nil {
		file, err := getFile(p.sourceFile, p.sourceContent)
		if err != nil {
			// error hard as we can't extract more data for the spec
			return nil, err
		}
		if err := p.parseSlothAnnotations(file.Comments...); err != nil {
			return nil, err
		}
		return p.specs, nil
	}

	applicationPackages := map[string]*ast.Package{}
	for _, dir := range p.includedDirs {
		// handle signals with context
		select {
		case <-ctx.Done():
			return nil, errors.New("termination signal was received, terminating process...")
		default:
		}
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			// skip if dir doesn't exists
			p.warn(err)
			continue
		}

		foundPkgs, err := getAllGoPackages(dir)
		if err != nil {
			p.warn(err)
			continue
		}

		for pkgName, pkg := range foundPkgs {
			if _, ok := applicationPackages[pkgName]; !ok {
				applicationPackages[pkgName] = pkg
			}
		}
	}

	// collect all sloth annotations from packages and add them to the spec struct
	if len(applicationPackages) > 0 {
		for _, pkg := range applicationPackages {
			for _, file := range pkg.Files {
				// handle signals with context
				select {
				case <-ctx.Done():
					return nil, errors.New("termination signal was received, terminating process...")
				default:
				}
				if err := p.parseSlothAnnotations(file.Comments...); err != nil {
					p.warn(err)
					continue
				}
			}
		}
	}

	return p.specs, nil
}

func (p parser) warn(err error, keyValues ...interface{}) {
	if p.logger != nil {
		p.logger.Warn(err, keyValues...)
	}
}
