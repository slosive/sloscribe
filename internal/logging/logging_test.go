package logging

import (
	"github.com/go-logr/stdr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidLevel(t *testing.T) {
	t.Parallel()

	t.Run("Successfully return true for debug as input", func(t *testing.T) {
		assert.True(t, IsValidLevel("debug"))
	})

	t.Run("Successfully return true for debug(with whitespace) as input", func(t *testing.T) {
		assert.True(t, IsValidLevel("debug    "))
	})

	t.Run("Successfully return true for DEBUG as input", func(t *testing.T) {
		assert.True(t, IsValidLevel("DEBUG"))
	})

	t.Run("Successfully return true for warn as input", func(t *testing.T) {
		assert.True(t, IsValidLevel("warn"))
	})

	t.Run("Successfully return true for info as input", func(t *testing.T) {
		assert.True(t, IsValidLevel("info"))
	})

	t.Run("Successfully return true for none as input", func(t *testing.T) {
		assert.True(t, IsValidLevel("none"))
	})

	t.Run("return false for invalid log-level as input", func(t *testing.T) {
		assert.False(t, IsValidLevel("none,"))
	})

	t.Run("return false for empty log-level as input", func(t *testing.T) {
		assert.False(t, IsValidLevel(""))
	})
}

func TestFindLevel(t *testing.T) {
	t.Parallel()

	t.Run("Successfully return stdr.All for debug as input", func(t *testing.T) {
		assert.Equal(t, stdr.All, findLogLevel("debug"))
	})

	t.Run("Successfully return stdr.All for debug(with whitespace) as input", func(t *testing.T) {
		assert.Equal(t, stdr.All, findLogLevel("    debug   "))
	})

	t.Run("Successfully return stdr.All for DEBUG as input", func(t *testing.T) {
		assert.Equal(t, stdr.All, findLogLevel("DEBUG"))
	})

	t.Run("Successfully return stdr.Error for warn as input", func(t *testing.T) {
		assert.Equal(t, stdr.Error, findLogLevel("warn"))
	})

	t.Run("Successfully return stdr.Info for info as input", func(t *testing.T) {
		assert.Equal(t, stdr.Info, findLogLevel("info"))
	})

	t.Run("Successfully return stdr.None for none as input", func(t *testing.T) {
		assert.Equal(t, stdr.None, findLogLevel("none"))
	})

	t.Run("return stdr.Info for invalid log-level as input", func(t *testing.T) {
		assert.Equal(t, stdr.Info, findLogLevel(""))
	})
}
