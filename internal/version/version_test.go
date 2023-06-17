package version

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildInfo(t *testing.T) {
	info := BuildInfo()
	assert.Equal(t, fmt.Sprintf(`
Host machine: %s

  Version:    %s
  Commit:     %s
  Built at:   %s`, Platform, Version, Commit, Date), info)
}

func TestInfo(t *testing.T) {
	info := Info()
	assert.Equal(t, Version, info)
}
