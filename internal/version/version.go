package version

import "fmt"

// Values injected at build-time
var (
	version string = "dev"
	commit  string = "unknown"
	date    string = "unknown"
)

func BuildInfo() string {
	return fmt.Sprintf("%s, commit %s, built at %s", version, commit, date)
}
