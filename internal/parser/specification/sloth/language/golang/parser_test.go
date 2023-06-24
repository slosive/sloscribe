package golang

import (
	"go/ast"
	"io"
	"strings"
	"testing"

	k8sloth "github.com/slok/sloth/pkg/kubernetes/api/sloth/v1"
	sloth "github.com/slok/sloth/pkg/prometheus/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetPackages(t *testing.T) {
	t.Parallel()
	t.Run("Successfully return all the go packages in the current and subdirectories", func(t *testing.T) {
		exp := []string{"golang", "testdata", "fixtures"}
		packages, err := getAllGoPackages("./.")
		require.NoError(t, err)
		require.Len(t, packages, len(exp))

		_, ok := packages[exp[0]]
		assert.True(t, ok)
		_, ok = packages[exp[1]]
		assert.True(t, ok)
		_, ok = packages[exp[2]]
		assert.True(t, ok)
	})

	t.Run("Successfully return the go packages in testdata/fixtures", func(t *testing.T) {
		exp := []string{"fixtures"}
		packages, err := getAllGoPackages("./testdata/fixtures")
		require.NoError(t, err)
		require.Len(t, packages, len(exp))

		_, ok := packages[exp[0]]
		assert.True(t, ok)
	})

	t.Run("Fails to return the go package in a non go package", func(t *testing.T) {
		_, err := getAllGoPackages("./testdata/gofake")
		require.Error(t, err)
	})

	t.Run("Fails to return the go package in a non-existing directory", func(t *testing.T) {
		_, err := getAllGoPackages("./testdata/non-existing")
		require.Error(t, err)
	})
}

func TestGetFile(t *testing.T) {
	t.Parallel()
	t.Run("Successfully get comments from testdata/fixtures/fixture.go", func(t *testing.T) {
		f, err := getFile("./testdata/fixtures/fixture.go", nil)
		require.NoError(t, err)
		assert.Equal(t, "Package fixtures contains testdata\n", f.Comments[0].Text())
	})
	t.Run("Fail to get comments from non existing file testdata/fixtures/fake.go", func(t *testing.T) {
		_, err := getFile("./testdata/fixtures/fake.go", nil)
		require.Error(t, err)
	})
	t.Run("Successfully get comments from string reader", func(t *testing.T) {
		f, err := getFile("", io.NopCloser(strings.NewReader(`
// Package fixtures contains testdata
package fixtures
`)))
		require.NoError(t, err)
		assert.Equal(t, "Package fixtures contains testdata\n", f.Comments[0].Text())
	})
	t.Run("Successfully get comments from string reader if filename is passed", func(t *testing.T) {
		f, err := getFile("./testdata/fixtures/fake.go", io.NopCloser(strings.NewReader(`
// Package fixtures contains testdata
package fixtures
`)))
		require.NoError(t, err)
		assert.Equal(t, "Package fixtures contains testdata\n", f.Comments[0].Text())
	})
	t.Run("Fails to get comments from empty string reader", func(t *testing.T) {
		_, err := getFile("", io.NopCloser(strings.NewReader(``)))
		require.Error(t, err)
	})
}

func TestParseAnnotations(t *testing.T) {
	t.Parallel()

	t.Run("Successfully parse the sloth annotations per single commentGroup, should return 1 specification", func(t *testing.T) {
		parser := NewParser(nil)
		require.NoError(t, parser.parseSlothAnnotations(&ast.CommentGroup{List: []*ast.Comment{{Text: `@sloth service foobar`}}},
			&ast.CommentGroup{List: []*ast.Comment{
				{
					Text: `@sloth.slo name availability`,
				},
				{
					Text: `@sloth.slo description availability SLO`,
				},
				{
					Text: `@sloth.slo objective 95.0`,
				},
			}}))
		assert.Equal(t, "foobar", (parser.specs["foobar"].(*sloth.Spec)).Service)
		assert.Equal(t, sloth.SLO{
			Name:        "availability",
			Description: "availability SLO",
			Objective:   95.0,
			Labels:      make(map[string]string),
			SLI:         sloth.SLI{},
			Alerting:    sloth.Alerting{},
		}, (parser.specs["foobar"].(*sloth.Spec)).SLOs[0])
	})

	t.Run("Successfully parse the sloth annotations per single commentGroup, should return 1 specification", func(t *testing.T) {
		parser := NewParser(nil)
		require.NoError(t, parser.parseSlothAnnotations(&ast.CommentGroup{List: []*ast.Comment{
			{
				Text: `@sloth service foobar`,
			},
			{
				Text: `@sloth.slo name availability`,
			},
			{
				Text: `@sloth.slo description availability SLO`,
			},
			{
				Text: `@sloth.slo objective 95.0`,
			},
		}}))
		assert.Equal(t, "foobar", (parser.specs["foobar"].(*sloth.Spec)).Service)
		assert.Equal(t, sloth.SLO{
			Name:        "availability",
			Description: "availability SLO",
			Objective:   95.0,
			Labels:      make(map[string]string),
			SLI:         sloth.SLI{},
			Alerting:    sloth.Alerting{},
		}, (parser.specs["foobar"].(*sloth.Spec)).SLOs[0])
	})

	t.Run("Successfully parse sloth service if service name is defined after SLO definition", func(t *testing.T) {
		parser := NewParser(nil)
		comments := []*ast.CommentGroup{
			{List: []*ast.Comment{
				{
					Text: `@sloth.slo name availability`,
				},
				{
					Text: `@sloth.slo description availability SLO for foobar service`,
				},
				{
					Text: `@sloth.slo objective 95.0`,
				},
			}},
			{List: []*ast.Comment{
				{
					Text: `@sloth service foobar`,
				},
			}},
		}
		require.NoError(t, parser.parseSlothAnnotations(comments...))
		require.Len(t, parser.specs, 1)
		resultSpec := parser.specs

		expected := []*sloth.Spec{
			{
				Version: sloth.Version,
				Service: "foobar",
				Labels:  make(map[string]string),
				SLOs: []sloth.SLO{
					{
						Name:        "availability",
						Description: "availability SLO for foobar service",
						Objective:   95.0,
						Labels:      make(map[string]string),
						SLI: sloth.SLI{
							Raw:    nil,
							Events: nil,
							Plugin: nil,
						},
						Alerting: sloth.Alerting{
							Name:        "",
							Labels:      nil,
							Annotations: nil,
							PageAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
							TicketAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
						},
					},
				},
			},
		}

		for _, exp := range expected {
			actual, ok := resultSpec[exp.Service]
			require.True(t, ok)
			assert.Equal(t, exp, actual)
		}
	})

	t.Run("Successfully parse multiple Sloth services with single SLO defined, should return 3 specifications", func(t *testing.T) {
		parser := NewParser(nil)
		comments := []*ast.CommentGroup{
			{List: []*ast.Comment{
				{
					Text: `@sloth service foobar`,
				},
				{
					Text: `@sloth.slo name availability`,
				},
				{
					Text: `@sloth.slo description availability SLO for foobar service`,
				},
				{
					Text: `@sloth.slo objective 95.0`,
				},
			}},
			{List: []*ast.Comment{
				{
					Text: `@sloth service foo`,
				},
				{
					Text: `@sloth.slo name availability`,
				},
				{
					Text: `@sloth.slo description availability SLO for foo service`,
				},
				{
					Text: `@sloth.slo objective 95.0`,
				},
			}},
			{List: []*ast.Comment{
				{
					Text: `@sloth service bar`,
				},
				{
					Text: `@sloth.slo name availability`,
				},
				{
					Text: `@sloth.slo description availability SLO for bar service`,
				},
				{
					Text: `@sloth.slo objective 95.0`,
				},
			}},
		}
		require.NoError(t, parser.parseSlothAnnotations(comments...))
		require.Len(t, parser.specs, 3)
		resultSpec := parser.specs

		expected := []*sloth.Spec{
			{
				Version: sloth.Version,
				Service: "foo",
				Labels:  make(map[string]string),
				SLOs: []sloth.SLO{
					{
						Name:        "availability",
						Description: "availability SLO for foo service",
						Objective:   95.0,
						Labels:      make(map[string]string),
						SLI: sloth.SLI{
							Raw:    nil,
							Events: nil,
							Plugin: nil,
						},
						Alerting: sloth.Alerting{
							Name:        "",
							Labels:      nil,
							Annotations: nil,
							PageAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
							TicketAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
						},
					},
				},
			},
			{
				Version: sloth.Version,
				Service: "bar",
				Labels:  make(map[string]string),
				SLOs: []sloth.SLO{
					{
						Name:        "availability",
						Description: "availability SLO for bar service",
						Objective:   95.0,
						Labels:      make(map[string]string),
						SLI: sloth.SLI{
							Raw:    nil,
							Events: nil,
							Plugin: nil,
						},
						Alerting: sloth.Alerting{
							Name:        "",
							Labels:      nil,
							Annotations: nil,
							PageAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
							TicketAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
						},
					},
				},
			},
			{
				Version: sloth.Version,
				Service: "foobar",
				Labels:  make(map[string]string),
				SLOs: []sloth.SLO{
					{
						Name:        "availability",
						Description: "availability SLO for foobar service",
						Objective:   95.0,
						Labels:      make(map[string]string),
						SLI: sloth.SLI{
							Raw:    nil,
							Events: nil,
							Plugin: nil,
						},
						Alerting: sloth.Alerting{
							Name:        "",
							Labels:      nil,
							Annotations: nil,
							PageAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
							TicketAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
						},
					},
				},
			},
		}

		for _, exp := range expected {
			actual, ok := resultSpec[exp.Service]
			require.True(t, ok)
			assert.Equal(t, exp, actual)
		}
	})

	t.Run("Successfully parse multiple Sloth services with multiple SLOs defined, should return 3 specifications", func(t *testing.T) {
		parser := NewParser(nil)
		comments := []*ast.CommentGroup{
			{List: []*ast.Comment{
				{
					Text: `@sloth service foobar`,
				},
				{
					Text: `@sloth.slo name availability`,
				},
				{
					Text: `@sloth.slo description availability SLO for foobar service`,
				},
				{
					Text: `@sloth.slo objective 95.0`,
				},
			}},
			{List: []*ast.Comment{
				{
					Text: `@sloth.slo name correctness`,
				},
				{
					Text: `@sloth.slo description correctness SLO for foobar service`,
				},
				{
					Text: `@sloth.slo objective 55.0`,
				},
			}},
			{List: []*ast.Comment{
				{
					Text: `@sloth service foo`,
				},
				{
					Text: `@sloth.slo name availability`,
				},
				{
					Text: `@sloth.slo description availability SLO for foo service`,
				},
				{
					Text: `@sloth.slo objective 95.0`,
				},
			}},
			{List: []*ast.Comment{
				{
					Text: `@sloth.slo name correctness`,
				},
				{
					Text: `@sloth.slo description correctness SLO for foo service`,
				},
				{
					Text: `@sloth.slo objective 85.0`,
				},
			}},
			{List: []*ast.Comment{
				{
					Text: `@sloth.slo name freshness`,
				},
				{
					Text: `@sloth.slo description freshness SLO for foo service`,
				},
				{
					Text: `@sloth.slo objective 99.999`,
				},
			}},
			{List: []*ast.Comment{
				{
					Text: `@sloth service bar`,
				},
				{
					Text: `@sloth.slo name availability`,
				},
				{
					Text: `@sloth.slo description availability SLO for bar service`,
				},
				{
					Text: `@sloth.slo objective 95.0`,
				},
			}},
		}
		require.NoError(t, parser.parseSlothAnnotations(comments...))
		require.Len(t, parser.specs, 3)
		resultSpec := parser.specs

		expected := []*sloth.Spec{
			{
				Version: sloth.Version,
				Service: "foo",
				Labels:  make(map[string]string),
				SLOs: []sloth.SLO{
					{
						Name:        "availability",
						Description: "availability SLO for foo service",
						Objective:   95.0,
						Labels:      make(map[string]string),
						SLI: sloth.SLI{
							Raw:    nil,
							Events: nil,
							Plugin: nil,
						},
						Alerting: sloth.Alerting{
							Name:        "",
							Labels:      nil,
							Annotations: nil,
							PageAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
							TicketAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
						},
					},
					{
						Name:        "correctness",
						Description: "correctness SLO for foo service",
						Objective:   85.0,
						Labels:      make(map[string]string),
						SLI: sloth.SLI{
							Raw:    nil,
							Events: nil,
							Plugin: nil,
						},
						Alerting: sloth.Alerting{
							Name:        "",
							Labels:      nil,
							Annotations: nil,
							PageAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
							TicketAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
						},
					},
					{
						Name:        "freshness",
						Description: "freshness SLO for foo service",
						Objective:   99.999,
						Labels:      make(map[string]string),
						SLI: sloth.SLI{
							Raw:    nil,
							Events: nil,
							Plugin: nil,
						},
						Alerting: sloth.Alerting{
							Name:        "",
							Labels:      nil,
							Annotations: nil,
							PageAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
							TicketAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
						},
					},
				},
			},
			{
				Version: sloth.Version,
				Service: "bar",
				Labels:  make(map[string]string),
				SLOs: []sloth.SLO{
					{
						Name:        "availability",
						Description: "availability SLO for bar service",
						Objective:   95.0,
						Labels:      make(map[string]string),
						SLI: sloth.SLI{
							Raw:    nil,
							Events: nil,
							Plugin: nil,
						},
						Alerting: sloth.Alerting{
							Name:        "",
							Labels:      nil,
							Annotations: nil,
							PageAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
							TicketAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
						},
					},
				},
			},
			{
				Version: sloth.Version,
				Service: "foobar",
				Labels:  make(map[string]string),
				SLOs: []sloth.SLO{
					{
						Name:        "availability",
						Description: "availability SLO for foobar service",
						Objective:   95.0,
						Labels:      make(map[string]string),
						SLI: sloth.SLI{
							Raw:    nil,
							Events: nil,
							Plugin: nil,
						},
						Alerting: sloth.Alerting{
							Name:        "",
							Labels:      nil,
							Annotations: nil,
							PageAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
							TicketAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
						},
					},
					{
						Name:        "correctness",
						Description: "correctness SLO for foobar service",
						Objective:   55.0,
						Labels:      make(map[string]string),
						SLI: sloth.SLI{
							Raw:    nil,
							Events: nil,
							Plugin: nil,
						},
						Alerting: sloth.Alerting{
							Name:        "",
							Labels:      nil,
							Annotations: nil,
							PageAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
							TicketAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
						},
					},
				},
			},
		}

		for _, exp := range expected {
			actual, ok := resultSpec[exp.Service]
			require.True(t, ok)
			assert.Equal(t, exp, actual)
		}
	})

	t.Run("Fail to parse the sloth spec SLO item if sloth annotation name for a given SLO is missing", func(t *testing.T) {
		parser := NewParser(nil)
		require.NoError(t, parser.parseSlothAnnotations(&ast.CommentGroup{List: []*ast.Comment{
			{
				Text: `@sloth service bar`,
			},
			{
				Text: `@sloth.slo description availability SLO`,
			},
			{
				Text: `@sloth.slo objective 95.0`,
			},
		}}))
		assert.Len(t, (parser.specs["bar"].(*sloth.Spec)).SLOs, 0)
	})

	t.Run("Successfully parse and merge duplicate Sloth service", func(t *testing.T) {
		parser := NewParser(nil)
		comments := []*ast.CommentGroup{
			{List: []*ast.Comment{
				{
					Text: `@sloth service foobar`,
				},
				{
					Text: `@sloth.slo name availability`,
				},
				{
					Text: `@sloth.slo description availability SLO for foobar service`,
				},
				{
					Text: `@sloth.slo objective 95.0`,
				},
			}},
			{List: []*ast.Comment{
				{
					Text: `@sloth service foobar`,
				},
				{
					Text: `@sloth.slo name foobar_availability`,
				},
				{
					Text: `@sloth.slo description availability SLO for foobar service`,
				},
				{
					Text: `@sloth.slo objective 99.0`,
				},
			}},
		}
		require.NoError(t, parser.parseSlothAnnotations(comments...))
		require.Len(t, parser.specs, 1)
		resultSpec := parser.specs

		expected := []*sloth.Spec{
			{
				Version: sloth.Version,
				Service: "foobar",
				Labels:  make(map[string]string),
				SLOs: []sloth.SLO{
					{
						Name:        "availability",
						Description: "availability SLO for foobar service",
						Objective:   95.0,
						Labels:      make(map[string]string),
						SLI: sloth.SLI{
							Raw:    nil,
							Events: nil,
							Plugin: nil,
						},
						Alerting: sloth.Alerting{
							Name:        "",
							Labels:      nil,
							Annotations: nil,
							PageAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
							TicketAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
						},
					},
					{
						Name:        "foobar_availability",
						Description: "availability SLO for foobar service",
						Objective:   99.0,
						Labels:      make(map[string]string),
						SLI: sloth.SLI{
							Raw:    nil,
							Events: nil,
							Plugin: nil,
						},
						Alerting: sloth.Alerting{
							Name:        "",
							Labels:      nil,
							Annotations: nil,
							PageAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
							TicketAlert: sloth.Alert{
								Disable:     false,
								Labels:      nil,
								Annotations: nil,
							},
						},
					},
				},
			},
		}

		for _, exp := range expected {
			actual, ok := resultSpec[exp.Service]
			require.True(t, ok)
			assert.Equal(t, exp, actual)
		}
	})
}

func TestParseK8SAnnotations(t *testing.T) {
	t.Parallel()

	t.Run("Successfully parse sloth service if service name is defined after slo definition", func(t *testing.T) {
		parser := NewParser(nil)
		comments := []*ast.CommentGroup{
			{List: []*ast.Comment{
				{
					Text: `@sloth.slo name availability`,
				},
				{
					Text: `@sloth.slo description availability SLO for foobar service`,
				},
				{
					Text: `@sloth.slo objective 95.0`,
				},
			}},
			{List: []*ast.Comment{
				{
					Text: `@sloth service foobar`,
				},
			}},
		}
		require.NoError(t, parser.parseK8SlothAnnotations(comments...))
		require.Len(t, parser.specs, 1)
		resultSpec := parser.specs

		expected := []*k8sloth.PrometheusServiceLevel{
			{
				TypeMeta: v1.TypeMeta{
					Kind:       "PrometheusServiceLevel",
					APIVersion: "sloth.slok.dev/v1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:   "foobar",
					Labels: make(map[string]string),
				},
				Spec: k8sloth.PrometheusServiceLevelSpec{
					Service: "foobar",
					Labels:  make(map[string]string),
					SLOs: []k8sloth.SLO{
						{
							Name:        "availability",
							Description: "availability SLO for foobar service",
							Objective:   95.0,
							Labels:      make(map[string]string),
							SLI: k8sloth.SLI{
								Raw:    nil,
								Events: nil,
								Plugin: nil,
							},
							Alerting: k8sloth.Alerting{
								Name:        "",
								Labels:      nil,
								Annotations: nil,
								PageAlert: k8sloth.Alert{
									Disable:     false,
									Labels:      nil,
									Annotations: nil,
								},
								TicketAlert: k8sloth.Alert{
									Disable:     false,
									Labels:      nil,
									Annotations: nil,
								},
							},
						},
					},
				},
			},
		}

		for _, exp := range expected {
			actual, ok := resultSpec[exp.Name]
			require.True(t, ok)
			assert.Equal(t, exp, actual)
		}
	})

	t.Run("Successfully parse multiple Sloth services, should return 3 specifications", func(t *testing.T) {
		parser := NewParser(nil)
		comments := []*ast.CommentGroup{
			{List: []*ast.Comment{
				{
					Text: `@sloth service foobar`,
				},
				{
					Text: `@sloth.slo name availability`,
				},
				{
					Text: `@sloth.slo description availability SLO for foobar service`,
				},
				{
					Text: `@sloth.slo objective 95.0`,
				},
			}},
			{List: []*ast.Comment{
				{
					Text: `@sloth service foo`,
				},
				{
					Text: `@sloth.slo name availability`,
				},
				{
					Text: `@sloth.slo description availability SLO for foo service`,
				},
				{
					Text: `@sloth.slo objective 95.0`,
				},
			}},
			{List: []*ast.Comment{
				{
					Text: `@sloth service bar`,
				},
				{
					Text: `@sloth.slo name availability`,
				},
				{
					Text: `@sloth.slo description availability SLO for bar service`,
				},
				{
					Text: `@sloth.slo objective 95.0`,
				},
			}},
		}
		require.NoError(t, parser.parseK8SlothAnnotations(comments...))
		require.Len(t, parser.specs, 3)
		resultSpec := parser.specs

		expected := []*k8sloth.PrometheusServiceLevel{
			{
				TypeMeta: v1.TypeMeta{
					Kind:       "PrometheusServiceLevel",
					APIVersion: "sloth.slok.dev/v1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:   "bar",
					Labels: make(map[string]string),
				},
				Spec: k8sloth.PrometheusServiceLevelSpec{
					Service: "bar",
					Labels:  make(map[string]string),
					SLOs: []k8sloth.SLO{
						{
							Name:        "availability",
							Description: "availability SLO for bar service",
							Objective:   95.0,
							Labels:      make(map[string]string),
							SLI: k8sloth.SLI{
								Raw:    nil,
								Events: nil,
								Plugin: nil,
							},
							Alerting: k8sloth.Alerting{
								Name:        "",
								Labels:      nil,
								Annotations: nil,
								PageAlert: k8sloth.Alert{
									Disable:     false,
									Labels:      nil,
									Annotations: nil,
								},
								TicketAlert: k8sloth.Alert{
									Disable:     false,
									Labels:      nil,
									Annotations: nil,
								},
							},
						},
					},
				},
			},
			{
				TypeMeta: v1.TypeMeta{
					Kind:       "PrometheusServiceLevel",
					APIVersion: "sloth.slok.dev/v1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:   "foobar",
					Labels: make(map[string]string),
				},
				Spec: k8sloth.PrometheusServiceLevelSpec{
					Service: "foobar",
					Labels:  make(map[string]string),
					SLOs: []k8sloth.SLO{
						{
							Name:        "availability",
							Description: "availability SLO for foobar service",
							Objective:   95.0,
							Labels:      make(map[string]string),
							SLI: k8sloth.SLI{
								Raw:    nil,
								Events: nil,
								Plugin: nil,
							},
							Alerting: k8sloth.Alerting{
								Name:        "",
								Labels:      nil,
								Annotations: nil,
								PageAlert: k8sloth.Alert{
									Disable:     false,
									Labels:      nil,
									Annotations: nil,
								},
								TicketAlert: k8sloth.Alert{
									Disable:     false,
									Labels:      nil,
									Annotations: nil,
								},
							},
						},
					},
				},
			},
			{
				TypeMeta: v1.TypeMeta{
					Kind:       "PrometheusServiceLevel",
					APIVersion: "sloth.slok.dev/v1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:   "foo",
					Labels: make(map[string]string),
				},
				Spec: k8sloth.PrometheusServiceLevelSpec{
					Service: "foo",
					Labels:  make(map[string]string),
					SLOs: []k8sloth.SLO{
						{
							Name:        "availability",
							Description: "availability SLO for foo service",
							Objective:   95.0,
							Labels:      make(map[string]string),
							SLI: k8sloth.SLI{
								Raw:    nil,
								Events: nil,
								Plugin: nil,
							},
							Alerting: k8sloth.Alerting{
								Name:        "",
								Labels:      nil,
								Annotations: nil,
								PageAlert: k8sloth.Alert{
									Disable:     false,
									Labels:      nil,
									Annotations: nil,
								},
								TicketAlert: k8sloth.Alert{
									Disable:     false,
									Labels:      nil,
									Annotations: nil,
								},
							},
						},
					},
				},
			},
		}

		for _, exp := range expected {
			actual, ok := resultSpec[exp.Name]
			require.True(t, ok)
			assert.Equal(t, exp, actual)
		}
	})
}
