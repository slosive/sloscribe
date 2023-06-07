package golang

import (
	"go/ast"
	"io"
	"strings"
	"testing"

	sloth "github.com/slok/sloth/pkg/prometheus/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	t.Run("Successfully parse the sloth annotations per single commentGroup", func(t *testing.T) {
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
		assert.Equal(t, "foobar", parser.getSpec().Service)
		assert.Equal(t, sloth.SLO{
			Name:        "availability",
			Description: "availability SLO",
			Objective:   95.0,
			Labels:      make(map[string]string),
			SLI:         sloth.SLI{},
			Alerting:    sloth.Alerting{},
		}, parser.getSpec().SLOs[0])
	})

	t.Run("Fail to parse the sloth annotations for a given SLO if no SLO name was given", func(t *testing.T) {
		parser := NewParser(nil)
		require.NoError(t, parser.parseSlothAnnotations(&ast.CommentGroup{List: []*ast.Comment{
			{
				Text: `@sloth.slo description availability SLO`,
			},
			{
				Text: `@sloth.slo objective 95.0`,
			},
		}}))
		assert.Len(t, parser.getSpec().SLOs, 0)
	})

	t.Run("Successfully parse the sloth annotations per multiple commentGroups", func(t *testing.T) {
		parser := NewParser(nil)
		require.NoError(t, parser.parseSlothAnnotations(
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
			},
			},
			&ast.CommentGroup{List: []*ast.Comment{
				{
					Text: `@sloth.slo name freshness`,
				},
				{
					Text: `@sloth.slo description freshness SLO`,
				},
				{
					Text: `@sloth.slo objective 95.0`,
				},
			}},
		))
		assert.Equal(t, sloth.SLO{
			Name:        "availability",
			Description: "availability SLO",
			Objective:   95.0,
			Labels:      make(map[string]string),
			SLI:         sloth.SLI{},
			Alerting:    sloth.Alerting{},
		}, parser.getSpec().SLOs[0])
		assert.Equal(t, sloth.SLO{
			Name:        "freshness",
			Description: "freshness SLO",
			Objective:   95.0,
			Labels:      make(map[string]string),
			SLI:         sloth.SLI{},
			Alerting:    sloth.Alerting{},
		}, parser.getSpec().SLOs[1])
	})
}
