package version

import (
	"fmt"
)

// Values injected at build-time
var (
	Platform string = "unknown"
	Version  string = "dev"
	Commit   string = "unknown"
	Date     string = "unknown"
)

// BuildInfo returns the binary build information
func BuildInfo() string {
	return fmt.Sprintf(`
Host machine: %s

  Version:    %s
  Commit:     %s
  Built at:   %s`, Platform, Version, Commit, Date)
}

// Info returns simple version information of the binary
func Info() string {
	return Version
}
