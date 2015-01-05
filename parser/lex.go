package parser

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

type TokenType int

type stateFn func(*lexer) *Token

type lexer struct {
	input       string // string being scanned
	start       int    // string being scanned
	pos         int    // current position of input
	width       int    // width of last rune read
	lastToken   *Token
	lastStateFn stateFn
	regionStack []TokenType
}

type Token struct {
	Typ TokenType
	Val string
}

// for debugging purposes
func (t Token) String() string {
	return fmt.Sprintf("(%s %s)", t.Typ.String(), t.Val)
}

func Lexer(src string) *lexer {
	return &lexer{input: src, lastStateFn: lexStart}
}

func (l *lexer) LexTillDone() []*Token {
	var res []*Token
	for token := l.Lex(); token != nil && token.Typ != EOF; token = l.Lex() {
		res = append(res, token)
	}
	return res
}

func (l *lexer) Lex() *Token {
	return l.lastStateFn(l)
}

func isBoolean(val string) bool {
	return val == "true" || val == "false"
}

func isKeyword(val string) bool {
	for _, k := range keywords {
		if val == k {
			return true
		}
	}
	return false
}

func isDigit(c rune) bool {
	return strings.IndexRune(num, c) != -1
}

func isOpOrDelim(val string) bool {
	for _, od := range opDelims {
		if strings.HasPrefix(val, od) {
			return true
		}
	}
	return false
}

func isOperatorCharacter(c rune) bool {
	if c == eof {
		return false
	}
	if !unicode.IsPrint(c) {
		return false
	}
	if strings.IndexRune(whitespace+letter+num+paren+delim+newline, c) != -1 {
		return false
	}
	return true
}

// reads & returns the next rune, steps width forward
func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, s := utf8.DecodeRuneInString(l.input[l.pos:])
	if r == utf8.RuneError && s == 1 {
		log.Fatal("input error")
	}
	l.width = s
	l.pos += l.width
	return r
}

func (l *lexer) skip() {
	l.next()
	l.ignore()
}

// can only be called once after each next
func (l *lexer) backup() {
	l.pos -= l.width
}

// accepts single rune in accepted
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) != -1 {
		return true
	}
	l.backup()
	return false
}

// accepts all runes in valid
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) != -1 {
	}
	l.backup()
}

func (l *lexer) acceptAllBut(invalid string) bool {
	for c := l.next(); c != eof && strings.IndexRune(invalid, c) == -1; {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRunAllBut(invalid string) {
	for c := l.next(); c != eof && strings.IndexRune(invalid, c) == -1; c = l.next() {
	}
	l.backup()
}

func isCharacterLiteral(val string) bool {
	re := regexp.MustCompile(`'([[:print:]]|\\.)'.*`)
	return re.MatchString(val)
}

func (l *lexer) lastAccepted() rune {
	// check if we did consume some character
	if l.width > 0 {
		l.backup()
		return l.next()
	}
	return eof
}

func (l *lexer) peek() rune {
	return l.peekNth(0)
}

// 0th peek gives you
func (l *lexer) peekNth(n int) rune {
	var res rune
	currentPos := l.pos
	for j := 0; j <= n; j++ {
		if (currentPos) >= len(l.input) {
			return eof
		}
		r, s := utf8.DecodeRuneInString(l.input[currentPos:])
		if r == utf8.RuneError && s == 1 {
			log.Fatal("input error")
		}
		currentPos += s
		res = r
	}
	return res
}

func (l *lexer) ignore() {
	l.start = l.pos
}

// Contract: Current token has to be consumed
// lookAhead(0) gives the current token from lexer.Pos
func (l *lexer) lookAhead(i int) (*Token, error) {
	var res *Token
	var err error
	if l.pos > l.start {
		return nil, errors.New("previous token was not consumed completely for the next token to be looked up")
	}
	// save the lexer state
	prev_start := l.start
	prev_pos := l.pos
	prev_width := l.width
	prev_token := l.lastToken
	prev_stateFn := l.lastStateFn
	// execute the state machine
	curr_token := lexStart(l)
	if i == 0 {
		res = curr_token
	} else {
		r, e := l.lookAhead(i - 1)
		if e == nil {
			res = r
		} else {
			err = e
		}
	}
	// invert the lexer changes
	switch curr_token.Typ {
	case L_PAREN, L_CURLY, L_BRACKET:
		// this would have pushed the paren into the stack
		// time to pop
		l.regionStack = l.regionStack[:len(l.regionStack)-1]
	case R_PAREN, R_CURLY, R_BRACKET:
		// this would have popped something from the stack
		// let's push it back
		invertedParen, e := invertParen(curr_token.Typ)
		if e == nil {
			l.regionStack = append(l.regionStack, invertedParen)
		} else {
			err = e
		}
	}
	// restore the lexer state
	l.start = prev_start
	l.pos = prev_pos
	l.width = prev_width
	l.lastToken = prev_token
	l.lastStateFn = prev_stateFn
	// return the token
	return res, err
}

func (l *lexer) emit(t TokenType, fn stateFn) *Token {
	v := l.input[l.start:l.pos]
	// reset the state
	l.start = l.pos
	l.width = 0
	// create the resultant token
	l.lastToken = &Token{Typ: t, Val: v}
	l.lastStateFn = fn
	return l.lastToken
}

// func (l *lexer) emitErrorf(format string, a ...interface{}) {
// 	l.Tokens <- Token{Typ: Error, Val: fmt.Sprintf(format, a)}
// }

func (l *lexer) emitError(a ...interface{}) *Token {
	// ignore up to whatever has been parsed
	if l.pos > l.start {
		l.ignore()
	} else {
		l.skip()
	}
	l.lastStateFn = lexStart
	return &Token{Typ: ERROR, Val: fmt.Sprint(a)}
}

// func (l *lexer) emitEof() {
// 	return Token{Typ: EOF}
// }

// func (l *lexer) emitSemicolon() {
// 	l.start = l.pos
// 	l.lastToken = &Token{Typ: OpOrDelim, Val: ";"}
// 	l.Tokens <- Token{Typ: OpOrDelim, Val: ";"}
// }

// peeks at the lexer's current value, without emitting it or changing
// the position.
func (l *lexer) val() string {
	return l.input[l.start:l.pos]
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func invertParen(l TokenType) (TokenType, error) {
	switch l {
	case L_PAREN:
		return R_PAREN, nil
	case L_BRACKET:
		return R_BRACKET, nil
	case L_CURLY:
		return R_CURLY, nil
	case R_PAREN:
		return L_PAREN, nil
	case R_BRACKET:
		return L_BRACKET, nil
	case R_CURLY:
		return L_CURLY, nil
	default:
		log.Fatal("invertParen: input error")
		return ERROR, errors.New("Unknown token type " + l.String())
	}
}

func isMatchingParen(l TokenType, r TokenType) bool {
	switch {
	case l == L_PAREN && r == R_PAREN:
		return true
	case l == L_BRACKET && r == R_BRACKET:
		return true
	case l == L_CURLY && r == R_CURLY:
		return true
	}
	return false
}
