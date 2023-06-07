package generate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidOutputFormat(t *testing.T) {
	t.Run("Successfully return true if the input format is json", func(t *testing.T) {
		assert.True(t, IsValidOutputFormat("json"))
	})
	t.Run("Successfully return true if the input format is yaml", func(t *testing.T) {
		assert.True(t, IsValidOutputFormat("yaml"))
	})
	t.Run("Successfully return true if the input format is YAML", func(t *testing.T) {
		assert.True(t, IsValidOutputFormat("YAML"))
	})
	t.Run("Successfully return true if the input format is YAML with whitespace", func(t *testing.T) {
		assert.True(t, IsValidOutputFormat("   YAML  "))
	})
	t.Run("Fail to return true if the input format is not supported", func(t *testing.T) {
		assert.False(t, IsValidOutputFormat("toml"))
	})
}
