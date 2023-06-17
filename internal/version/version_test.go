package version

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildInfo(t *testing.T) {
	info := BuildInfo()
	assert.Equal(t, fmt.Sprintf("%s, commit %s, built at %s", Version, Commit, Date), info)
}
