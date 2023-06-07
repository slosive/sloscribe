package lang

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsSupportedLanguage(t *testing.T) {
	t.Parallel()

	t.Run("Successfully return true if Go is the target language", func(t *testing.T) {
		assert.True(t, IsSupportedLanguage(Go))
	})

	t.Run("Successfully return true if Rust is the target language", func(t *testing.T) {
		assert.True(t, IsSupportedLanguage(Rust))
	})

	t.Run("Fail to return true if the language is different from the supported ones", func(t *testing.T) {
		assert.False(t, IsSupportedLanguage("Python"))
	})
}
