package lang

type (
	// SourceLanguage represents the language for the source file parsed by the generator
	SourceLanguage string
)

const (
	Go   = SourceLanguage("go")
	Wasm = SourceLanguage("wasm")
)

// IsSupportedLanguage returns true is the input language is a supported language
func IsSupportedLanguage(l SourceLanguage) bool {
	switch l {
	case Go, Wasm:
		return true
	}
	return false
}
