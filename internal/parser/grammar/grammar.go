package grammar

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/juju/errors"
	sloth "github.com/slok/sloth/pkg/prometheus/api/v1"
	"reflect"
	"strconv"
	"strings"
)

type (
	ServiceInfoGrammar struct {
		Stmts []*Statement `@@*`
	}

	Grammar struct {
		SloStmts []*Statement `@@*`
	}

	Statement struct {
		Key   Key    `@@`
		Value string `Whitespace* @(String (Whitespace|EOL)*)+`
	}

	Key struct {
		Type  string `(Sloth @((".alerting"|".sli")(".page"|".ticket")?)?)`
		Value string `Whitespace* @Keyword`
	}
)

const (
	serviceAttr            = "service"
	versionAttr            = "version"
	sliErrorQueryAttr      = "error_query"
	sliTotalQueryAttr      = "total_query"
	sliErrorRatioQueryAttr = "error_ratio_query"
	nameAttr               = "name"
	descriptionAttr        = "description"
	objectiveAttr          = "objective"
	labelAttr              = "labels"
	annotationAttr         = "annotations"
	disableAttr            = "disable"
)

var (
	ErrMissingRequiredField = errors.New("missing required application field(s)")
	ErrParseSource          = errors.New("error parsing source material")
	keywords                = []string{
		sliErrorQueryAttr,
		sliTotalQueryAttr,
		sliErrorRatioQueryAttr,
		nameAttr,
		descriptionAttr,
		objectiveAttr,
		labelAttr,
		annotationAttr,
		disableAttr,
		serviceAttr,
		versionAttr,
	}
)

func (k Key) GetStmtType() string {
	return k.Type
}

func IsMapValue(attr string) bool {
	return attr == labelAttr || attr == annotationAttr
}

func (g ServiceInfoGrammar) getAttribute(attribute string) (map[string]string, bool) {
	attributes := make(map[string]string)
	for _, attr := range g.Stmts {
		if strings.ToLower(attr.Key.Value) == attribute {
			if !IsMapValue(attribute) {
				return map[string]string{
					attribute: strings.TrimSpace(attr.Value),
				}, true
			}
			// get label name
			maps := strings.SplitN(attr.Value, " ", 1)
			name := strings.TrimSpace(maps[0])
			value := strings.TrimSpace(maps[1])
			attributes[name] = value
		}
	}
	return attributes, len(attributes) > 0
}

func parseFields(attr string, value string, fields []reflect.StructField, pValue reflect.Value) error {
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
							panic(err)
						}
						v.SetBool(b)
					case reflect.Float64:
						f, err := strconv.ParseFloat(value, 64)
						if err != nil {
							panic(err)
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

func (g Grammar) parseSLO() (*sloth.SLO, error) {
	var slo = &sloth.SLO{
		Name:        "",
		Description: "",
		Objective:   0,
		Labels:      map[string]string{},
		SLI: sloth.SLI{
			Raw:    &sloth.SLIRaw{},
			Events: &sloth.SLIEvents{},
			Plugin: nil,
		},
	}
	for _, attr := range g.SloStmts {
		switch attr.Key.GetStmtType() {
		case ".alerting.ticket":
			alert := &sloth.Alert{
				Disable:     false,
				Labels:      map[string]string{},
				Annotations: map[string]string{},
			}
			fields := reflect.VisibleFields(reflect.TypeOf(*alert))
			pValue := reflect.ValueOf(alert).Elem()
			parseFields(strings.ToLower(attr.Key.Value), strings.TrimSpace(attr.Value), fields, pValue)
			slo.Alerting.TicketAlert = *alert
		case ".alerting.page":
			alert := &sloth.Alert{
				Disable:     false,
				Labels:      map[string]string{},
				Annotations: map[string]string{},
			}
			fields := reflect.VisibleFields(reflect.TypeOf(*alert))
			pValue := reflect.ValueOf(alert).Elem()
			parseFields(strings.ToLower(attr.Key.Value), strings.TrimSpace(attr.Value), fields, pValue)
			slo.Alerting.PageAlert = *alert
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
			parseFields(strings.ToLower(attr.Key.Value), strings.TrimSpace(attr.Value), fields, pValue)
			slo.Alerting = *alerting
		case ".sli":
			// SLI
			switch attr.Key.Value {
			case sliTotalQueryAttr:
				slo.SLI.Events.TotalQuery = strings.TrimSpace(attr.Value)
			case sliErrorQueryAttr:
				slo.SLI.Events.ErrorQuery = strings.TrimSpace(attr.Value)
			case sliErrorRatioQueryAttr:
				slo.SLI.Raw.ErrorRatioQuery = strings.TrimSpace(attr.Value)
			}
		default:
			fields := reflect.VisibleFields(reflect.TypeOf(*slo))
			pValue := reflect.ValueOf(slo).Elem()
			parseFields(strings.ToLower(attr.Key.Value), strings.TrimSpace(attr.Value), fields, pValue)
		}
	}
	return slo, nil
}

var lexerDefinition = lexer.MustSimple([]lexer.SimpleRule{
	{"EOL", `[\n\r]+`},
	{"Keyword", strings.Join(keywords, "|")},
	{"Sloth", `@sloth`},
	{"String", `([a-zA-Z_0-9\.\/:,\-\'\(\)~\[\]\{\}=\"\|%])\w*`},
	{"Whitespace", `[ \t]+`},
})

func eval(filename, source string, options ...participle.ParseOption) (*Grammar, error) {
	ast, err := participle.Build[Grammar](
		participle.Lexer(lexerDefinition),
	)
	if err != nil {
		return nil, err
	}

	return ast.ParseString(filename, source, options...)
}

func Eval(source string, options ...participle.ParseOption) (map[string]sloth.SLO, error) {
	grammar, err := eval("", source, options...)
	if err != nil {
		return nil, err
	}

	foundSLOs := make(map[string]sloth.SLO)
	newSLO := sloth.SLO{}

	// SLO
	s, _ := grammar.parseSLO()
	newSLO = *s

	// TODO checks on the required fields
	if newSLO.Objective == 0 {
		return nil, errors.New("SLO's objective is missing")
	}
	if newSLO.Name == "" {
		return nil, errors.New("SLO's name is missing")
	}

	foundSLOs[newSLO.Name] = newSLO
	return foundSLOs, nil
}

func evalService(filename, source string, options ...participle.ParseOption) (*ServiceInfoGrammar, error) {
	ast, err := participle.Build[ServiceInfoGrammar](
		participle.Lexer(lexerDefinition),
	)
	if err != nil {
		return nil, err
	}

	return ast.ParseString(filename, source, options...)
}

func EvalService(source string, options ...participle.ParseOption) (*sloth.Spec, error) {
	grammar, err := evalService("", source, options...)
	if err != nil {
		return nil, err
	}

	spec := &sloth.Spec{
		Version: sloth.Version,
		Service: "",
		Labels:  nil,
	}

	// Spec
	if attrs, ok := grammar.getAttribute(serviceAttr); ok {
		spec.Service = attrs[serviceAttr]
	}
	if attrs, ok := grammar.getAttribute(versionAttr); ok {
		spec.Version = attrs[versionAttr]
	}
	if attrs, ok := grammar.getAttribute(labelAttr); ok {
		spec.Labels = attrs
	}

	return spec, nil
}
