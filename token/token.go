package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILEGAL"
	EOF     = "EOF"
	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1234566
	FLOAT = "FLOAT" // 120.224253456
	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	BANG     = "!"
	MINUS    = "-"
	SLASH    = "/"
	ASTERISK = "*"
	MOD      = "%"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="
	AND      = "&&"
	OR       = "||"
	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
	DOT       = "."
	// Keywords
	FUNCTION = "FUNCTION"
	STRING   = "STRING"
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	FOR      = "FOR"
	SLEEP    = "SLEEP"
	WHILE    = "WHILE"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	TYPEOF   = "TYPEOF"
	EXEC     = "EXEC"
	NEW      = "NEW"
	CLASS    = "CLASS"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"for":    FOR,
	"sleep":  SLEEP,
	"exec":   EXEC,
	"while":  WHILE,
	"typeof": TYPEOF,
	"new":    NEW,
	"class":  CLASS,
}

func LookupIdent(ident string) TokenType {
	if tokenType, ok := keywords[ident]; ok {
		return tokenType
	}
	return IDENT
}
