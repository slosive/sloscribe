package grammar

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/juju/errors"
	sloth "github.com/slok/sloth/pkg/prometheus/api/v1"
)

type (
	// Grammar is the participle grammar use to parse the Sloth comment groups in source files
	Grammar struct {
		// Stmts is a list of Sloth grammar Statements
		Stmts []*Statement `@@*`
	}
	// Statement is any comment starting with @sloth keyword
	Statement struct {
		Scope Scope  `@@`
		Value string `Whitespace* @(String (Whitespace|EOL)*)+`
	}
	// Scope defines the statement scope, similar to a code function
	Scope struct {
		// Type is the specification struct a statement refers to
		Type  string `(Sloth @((".alerting"(".page"|".ticket")?|".sli"|".slo"))?)`
		Value string `Whitespace* @("service"|"version"|"error_query"|"total_query"|"error_ratio_query"|"name"|"description"|"objective"|"labels"|"annotations"|"disable")`
	}
)

const (
	sliErrorQueryAttr      = "error_query"
	sliTotalQueryAttr      = "total_query"
	sliErrorRatioQueryAttr = "error_ratio_query"
)

var (
	ErrMissingRequiredField = errors.New("missing required application field(s)")
	ErrParseSource          = errors.New("error parsing source material")
)

// GetType returns the type of the statement scope
func (k Scope) GetType() string {
	return k.Type
}

// parseAndAssignStructFields
func parseAndAssignStructFields(attr string, value string, fields []reflect.StructField, pValue reflect.Value) error {
	for _, field := range fields {
		tag, ok := field.Tag.Lookup("yaml")
		if !ok {
			return nil
		}
		key := strings.Split(tag, ",")[0]
		if attr == key {
			// set field value
			v := pValue.FieldByName(field.Name)
			if v.IsValid() {
				if v.CanSet() {
					switch v.Kind() {
					case reflect.Bool:
						b, err := strconv.ParseBool(value)
						if err != nil {
							return err
						}
						v.SetBool(b)
					case reflect.Float64:
						f, err := strconv.ParseFloat(value, 64)
						if err != nil {
							return err
						}
						v.SetFloat(f)
					case reflect.Map:
						// label or annotation
						m := strings.Split(value, " ")
						v.SetMapIndex(reflect.ValueOf(m[0]), reflect.ValueOf(m[1]))
					default:
						v.Set(reflect.ValueOf(value))
					}
				}
			}
		}
	}
	return nil
}

func (g Grammar) parse() (*sloth.Spec, error) {
	var spec = &sloth.Spec{
		Version: sloth.Version,
		Service: "",
	}
	var slo = &sloth.SLO{
		Name:        "",
		Description: "",
		Objective:   0,
		Labels:      map[string]string{},
		SLI: sloth.SLI{
			Raw:    nil,
			Events: nil,
			Plugin: nil,
		},
	}

	for _, attr := range g.Stmts {
		switch attr.Scope.GetType() {
		case ".alerting.ticket":
			alert := &sloth.Alert{
				Disable:     false,
				Labels:      map[string]string{},
				Annotations: map[string]string{},
			}
			fields := reflect.VisibleFields(reflect.TypeOf(*alert))
			pValue := reflect.ValueOf(alert).Elem()
			if err := parseAndAssignStructFields(strings.ToLower(attr.Scope.Value), strings.TrimSpace(attr.Value), fields, pValue); err == nil {
				slo.Alerting.TicketAlert = *alert
			}
		case ".alerting.page":
			alert := &sloth.Alert{
				Disable:     false,
				Labels:      map[string]string{},
				Annotations: map[string]string{},
			}
			fields := reflect.VisibleFields(reflect.TypeOf(*alert))
			pValue := reflect.ValueOf(alert).Elem()
			if err := parseAndAssignStructFields(strings.ToLower(attr.Scope.Value), strings.TrimSpace(attr.Value), fields, pValue); err == nil {
				slo.Alerting.PageAlert = *alert
			}
		case ".alerting":
			alerting := &sloth.Alerting{
				Name:        "",
				Labels:      map[string]string{},
				Annotations: map[string]string{},
				PageAlert:   sloth.Alert{},
				TicketAlert: sloth.Alert{},
			}
			fields := reflect.VisibleFields(reflect.TypeOf(*alerting))
			pValue := reflect.ValueOf(alerting).Elem()
			if err := parseAndAssignStructFields(strings.ToLower(attr.Scope.Value), strings.TrimSpace(attr.Value), fields, pValue); err == nil && alerting.Name != "" {
				slo.Alerting = *alerting
			}
		case ".sli":
			// SLI
			switch attr.Scope.Value {
			case sliTotalQueryAttr:
				if slo.SLI.Events == nil {
					slo.SLI.Events = &sloth.SLIEvents{}
				}
				slo.SLI.Events.TotalQuery = strings.TrimSpace(attr.Value)
			case sliErrorQueryAttr:
				if slo.SLI.Events == nil {
					slo.SLI.Events = &sloth.SLIEvents{}
				}
				slo.SLI.Events.ErrorQuery = strings.TrimSpace(attr.Value)
			case sliErrorRatioQueryAttr:
				if slo.SLI.Raw == nil {
					slo.SLI.Raw = &sloth.SLIRaw{}
				}
				slo.SLI.Raw.ErrorRatioQuery = strings.TrimSpace(attr.Value)
			}
		case ".slo":
			fields := reflect.VisibleFields(reflect.TypeOf(*slo))
			pValue := reflect.ValueOf(slo).Elem()
			if err := parseAndAssignStructFields(strings.ToLower(attr.Scope.Value), strings.TrimSpace(attr.Value), fields, pValue); err != nil {
				continue
			}
		default:
			fields := reflect.VisibleFields(reflect.TypeOf(*spec))
			pValue := reflect.ValueOf(spec).Elem()
			if err := parseAndAssignStructFields(strings.ToLower(attr.Scope.Value), strings.TrimSpace(attr.Value), fields, pValue); err != nil {
				continue
			}
		}
	}

	if slo.Name != "" {
		spec.SLOs = []sloth.SLO{*slo}
	}

	return spec, nil
}

func createGrammar(filename, source string, options ...participle.ParseOption) (*Grammar, error) {
	ast, err := participle.Build[Grammar](
		participle.Lexer(lexerDefinition),
	)
	if err != nil {
		return nil, err
	}

	return ast.ParseString(filename, source, options...)
}

// Eval evaluates the source input against the grammar and returns an instance of *sloth.spec
func Eval(source string, options ...participle.ParseOption) (*sloth.Spec, error) {
	grammar, err := createGrammar("", source, options...)
	if err != nil {
		return nil, err
	}

	// Spec
	s, _ := grammar.parse()
	newSpec := *s

	return &newSpec, nil
}
