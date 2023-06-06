package grammar

import "github.com/alecthomas/participle/v2/lexer"

var lexerDefinition = lexer.MustSimple([]lexer.SimpleRule{
	{"EOL", `[\n\r]+`},
	{"Sloth", `@sloth`},
	{"String", `([a-zA-Z_0-9\.\/:,\-\'\(\)~\[\]\{\}=\"\|%])\w*`},
	{"Whitespace", `[ \t]+`},
})
