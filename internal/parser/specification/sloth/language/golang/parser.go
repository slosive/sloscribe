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
	k8sloth "github.com/slok/sloth/pkg/kubernetes/api/sloth/v1"
	sloth "github.com/slok/sloth/pkg/prometheus/api/v1"
	"github.com/slosive/sloscribe/internal/logging"
	"github.com/slosive/sloscribe/internal/parser/specification/sloth/grammar"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type parser struct {
	specs         map[string]any
	sourceFile    string
	sourceContent io.ReadCloser
	includedDirs  []string
	logger        *logging.Logger
	kubernetes    bool
}

// Options contains the configuration options available to the Parser
type Options struct {
	Logger           *logging.Logger
	SourceFile       string
	SourceContent    io.ReadCloser
	InputDirectories []string
	Kubernetes       bool
}

func NewOptions() *Options {
	l := logging.NewStandardLogger()
	return &Options{
		Logger:           &l,
		SourceFile:       "",
		SourceContent:    nil,
		InputDirectories: nil,
		Kubernetes:       false,
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
		specs:         map[string]any{},
		sourceFile:    sourceFile,
		sourceContent: sourceContent,
		includedDirs:  dirs,
		logger:        logger,
		kubernetes:    opts.Kubernetes,
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

func (p parser) parseK8SlothAnnotations(comments ...*ast.CommentGroup) error {
	var currentServiceSpec *k8sloth.PrometheusServiceLevel

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
			if currentServiceSpec != nil && (currentServiceSpec.Name == partialServiceSpec.Service || currentServiceSpec.Name == "") {
				p.specs[partialServiceSpec.Service] = currentServiceSpec
			}
			spec, ok := p.specs[partialServiceSpec.Service]
			if !ok {
				tmpSpec := &k8sloth.PrometheusServiceLevel{
					TypeMeta: v1.TypeMeta{
						Kind:       "PrometheusServiceLevel",
						APIVersion: "sloth.slok.dev/v1",
					},
					ObjectMeta: v1.ObjectMeta{
						Name:   partialServiceSpec.Service,
						Labels: map[string]string{},
					},
					Spec: k8sloth.PrometheusServiceLevelSpec{
						Service: partialServiceSpec.Service,
						Labels:  partialServiceSpec.Labels,
						SLOs:    toKubernetes(partialServiceSpec.SLOs...),
					},
				}
				p.specs[partialServiceSpec.Service] = tmpSpec
				currentServiceSpec = tmpSpec
			} else {
				currentServiceSpec = spec.(*k8sloth.PrometheusServiceLevel)
			}
		}

		if currentServiceSpec == nil {
			currentServiceSpec = &k8sloth.PrometheusServiceLevel{
				TypeMeta: v1.TypeMeta{
					Kind:       "PrometheusServiceLevel",
					APIVersion: "sloth.slok.dev/v1",
				},
				Spec: k8sloth.PrometheusServiceLevelSpec{
					Service: "",
					Labels:  nil,
					SLOs:    nil,
				},
			}
		}

		if currentServiceSpec.Name == "" {
			currentServiceSpec.Name = partialServiceSpec.Service
		}
		if currentServiceSpec.Labels == nil {
			for key, label := range partialServiceSpec.Labels {
				currentServiceSpec.Labels[key] = label
			}
		}
		if currentServiceSpec.Spec.Service == "" {
			currentServiceSpec.Spec.Service = partialServiceSpec.Service
		}
		if currentServiceSpec.Spec.Labels == nil {
			for key, label := range partialServiceSpec.Labels {
				currentServiceSpec.Spec.Labels[key] = label
			}
		}
		if currentServiceSpec.Spec.SLOs == nil {
			currentServiceSpec.Spec.SLOs = append(currentServiceSpec.Spec.SLOs, toKubernetes(partialServiceSpec.SLOs...)...)
		}
	}
	return nil
}

func toKubernetes(slos ...sloth.SLO) []k8sloth.SLO {
	var k8SLOs []k8sloth.SLO
	for _, slo := range slos {
		result := k8sloth.SLO{
			Name:        slo.Name,
			Description: slo.Description,
			Objective:   slo.Objective,
			Labels:      slo.Labels,
			SLI: k8sloth.SLI{
				Raw:    (*k8sloth.SLIRaw)(slo.SLI.Raw),
				Events: (*k8sloth.SLIEvents)(slo.SLI.Events),
			},
			Alerting: k8sloth.Alerting{
				Name:        slo.Alerting.Name,
				Labels:      slo.Alerting.Labels,
				Annotations: slo.Alerting.Annotations,
				PageAlert: k8sloth.Alert{
					Disable:     slo.Alerting.PageAlert.Disable,
					Labels:      slo.Alerting.PageAlert.Labels,
					Annotations: slo.Alerting.PageAlert.Annotations,
				},
				TicketAlert: k8sloth.Alert{
					Disable:     slo.Alerting.TicketAlert.Disable,
					Labels:      slo.Alerting.TicketAlert.Labels,
					Annotations: slo.Alerting.TicketAlert.Annotations,
				},
			},
		}

		k8SLOs = append(k8SLOs, result)
	}

	return k8SLOs
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
				currentServiceSpec = spec.(*sloth.Spec)
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
func (p parser) Parse(ctx context.Context) (map[string]any, error) {
	// collect all sloth annotations from the file and add them to the spec struct
	if p.sourceFile != "" || p.sourceContent != nil {
		file, err := getFile(p.sourceFile, p.sourceContent)
		if err != nil {
			// error hard as we can't extract more data for the spec
			return nil, err
		}
		if p.kubernetes {
			if err := p.parseK8SlothAnnotations(file.Comments...); err != nil {
				return nil, err
			}
		} else {
			if err := p.parseSlothAnnotations(file.Comments...); err != nil {
				return nil, err
			}
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
