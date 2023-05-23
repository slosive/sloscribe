
export function hello_world(): string {
    return `
{
    "glossary": {
        "title": "example glossary",
\t\t"GlossDiv": {
            "title": "S",
\t\t\t"GlossList": {
                "GlossEntry": {
                    "ID": "SGML",
\t\t\t\t\t"SortAs": "SGML",
\t\t\t\t\t"GlossTerm": "Standard Generalized Markup Language",
\t\t\t\t\t"Acronym": "SGML",
\t\t\t\t\t"Abbrev": "ISO 8879:1986",
\t\t\t\t\t"GlossDef": {
                        "para": "A meta-markup language, used to create markup languages such as DocBook.",
\t\t\t\t\t\t"GlossSeeAlso": ["GML", "XML"]
                    },
\t\t\t\t\t"GlossSee": "markup"
                }
            }
        }
    }
}
`;
}
