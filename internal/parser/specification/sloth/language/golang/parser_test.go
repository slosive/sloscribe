package golang

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
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
