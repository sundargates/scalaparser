package parser

import "fmt"

const eof = -1

const (
	semicolon             = ";"
	newline               = "\n"
	space                 = " "
	tab                   = "\t"
	whitespaceSansNewline = space + tab
	whitespace            = whitespaceSansNewline + newline
	quote                 = "\""
	singlequote           = "'"
	multilinequote        = "\"\"\""
	backslash             = "\\"
	linecomment           = "//"
	spancomment           = "/*"
	alphaLower            = "abcdefghijklmnopqrstuvwxyz"
	alphaUpper            = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alpha                 = alphaLower + alphaUpper
	letter                = alpha + "_" + "$"
	numSansZero           = "123456789"
	num                   = numSansZero + "0"
	alphaNum              = alpha + num
	letterNum             = letter + num
	paren                 = "()[]{}"
	delim                 = "`'\".;,"
	backtick              = "`"
)

var opDelims = [...]string{"=>", "<-", ">:", "<:", "<%", "-", ",", "#", "@", ":", "*", "=", "|", "~", "!", "+", "-"}
var keywords = [...]string{"abstract", "case", "catch", "class", "def", "do", "else", "extends", "final", "finally", "for", "forSome", "if", "implicit", "import", "lazy", "match", "new", "null", "object", "override", "package", "private", "protected", "return", "sealed", "super", "this", "throw", "trait", "try", "type", "val", "var", "while", "with", "yield"}

const (
	NIL TokenType = iota
	EOF
	ERROR
	SYMBOL
	IDENTIFIER
	NUMBER
	BOOLEAN
	CHARACTER
	STRING
	WHITESPACE
	COMMENT
	NEWLINE
	NEWLINES
	// operators and punctuation
	OPORDELIM
	SEMICOLON
	DOT
	L_PAREN
	R_PAREN
	L_BRACKET
	R_BRACKET
	L_CURLY
	R_CURLY
	// langauge keywords
	ABSTRACT
	CASE
	CATCH
	CLASS
	DEF
	DO
	ELSE
	EXTENDS
	FALSE
	FINAL
	FINALLY
	FOR
	FORSOME
	IF
	IMPLICIT
	IMPORT
	LAZY
	MATCH
	NEW
	NULL
	OBJECT
	OVERRIDE
	PACKAGE
	PRIVATE
	PROTECTED
	RETURN
	SEALED
	SUPER
	THIS
	THROW
	TRAIT
	TRUE
	TRY
	TYPE
	VAL
	VAR
	WHILE
	WITH
	YIELD
)

var parenToTokenType = map[string]TokenType{
	"(": L_PAREN,
	")": R_PAREN,
	"[": L_BRACKET,
	"]": R_BRACKET,
	"{": L_CURLY,
	"}": R_CURLY,
}

var keywordsToTokenType = map[string]TokenType{
	"abstract":  ABSTRACT,
	"case":      CASE,
	"catch":     CATCH,
	"class":     CLASS,
	"def":       DEF,
	"do":        DO,
	"else":      ELSE,
	"extends":   EXTENDS,
	"false":     FALSE,
	"final":     FINAL,
	"finally":   FINALLY,
	"for":       FOR,
	"forsome":   FORSOME,
	"if":        IF,
	"implicit":  IMPLICIT,
	"import":    IMPORT,
	"lazy":      LAZY,
	"match":     MATCH,
	"new":       NEW,
	"null":      NULL,
	"object":    OBJECT,
	"override":  OVERRIDE,
	"package":   PACKAGE,
	"private":   PRIVATE,
	"protected": PROTECTED,
	"return":    RETURN,
	"sealed":    SEALED,
	"super":     SUPER,
	"this":      THIS,
	"throw":     THROW,
	"trait":     TRAIT,
	"true":      TRUE,
	"try":       TRY,
	"type":      TYPE,
	"val":       VAL,
	"var":       VAR,
	"while":     WHILE,
	"with":      WITH,
	"yield":     YIELD,
}

func (i TokenType) String() string {
	switch i {
	case NIL:
		return "NIL"
	case EOF:
		return "EOF"
	case ERROR:
		return "ERROR"
	case CHARACTER:
		return "CHARACTER"
	case IDENTIFIER:
		return "IDENTIFIER"
	case SYMBOL:
		return "SYMBOL"
	case NUMBER:
		return "NUMBER"
	case BOOLEAN:
		return "BOOLEAN"
	case STRING:
		return "STRING"
	case WHITESPACE:
		return "WHITESPACE"
	case COMMENT:
		return "COMMENT"
	case NEWLINE:
		return "NEWLINE"
	case NEWLINES:
		return "NEWLINES"
	// operators and punctuation: return "// operators and punctuation"
	case OPORDELIM:
		return "OPERATOR"
	case SEMICOLON:
		return ";"
	case DOT:
		return "."
	case L_PAREN:
		return "("
	case R_PAREN:
		return ")"
	case L_CURLY:
		return "{"
	case R_CURLY:
		return "}"
	case L_BRACKET:
		return "["
	case R_BRACKET:
		return "]"
	// langauge keywords: return "// langauge keywords"
	case ABSTRACT:
		return "ABSTRACT"
	case CASE:
		return "CASE"
	case CATCH:
		return "CATCH"
	case CLASS:
		return "CLASS"
	case DEF:
		return "DEF"
	case DO:
		return "DO"
	case ELSE:
		return "ELSE"
	case EXTENDS:
		return "EXTENDS"
	case FINAL:
		return "FINAL"
	case FINALLY:
		return "FINALLY"
	case FOR:
		return "FOR"
	case FORSOME:
		return "FORSOME"
	case IF:
		return "IF"
	case IMPLICIT:
		return "IMPLICIT"
	case IMPORT:
		return "IMPORT"
	case LAZY:
		return "LAZY"
	case MATCH:
		return "MATCH"
	case NEW:
		return "NEW"
	case NULL:
		return "NULL"
	case OBJECT:
		return "OBJECT"
	case OVERRIDE:
		return "OVERRIDE"
	case PACKAGE:
		return "PACKAGE"
	case PRIVATE:
		return "PRIVATE"
	case PROTECTED:
		return "PROTECTED"
	case RETURN:
		return "RETURN"
	case SEALED:
		return "SEALED"
	case SUPER:
		return "SUPER"
	case THIS:
		return "THIS"
	case THROW:
		return "THROW"
	case TRAIT:
		return "TRAIT"
	case TRY:
		return "TRY"
	case TYPE:
		return "TYPE"
	case VAL:
		return "VAL"
	case VAR:
		return "VAR"
	case WHILE:
		return "WHILE"
	case WITH:
		return "WITH"
	case YIELD:
		return "YIELD"
	}
	fmt.Printf("TokenType = %d\n", i)
	return "Whooops"
}
