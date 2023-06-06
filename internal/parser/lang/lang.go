package lang

type (
	// Target represents the language for the source file parsed by the generator
	Target string
)

const (
	Go   = Target("go")
	Rust = Target("rust")
)

// IsSupportedLanguage returns true is the input language is a supported language
func IsSupportedLanguage(l Target) bool {
	switch l {
	case Go, Rust:
		return true
	}
	return false
}
