package parser

import (
	"fmt"
	"strings"
	"unicode"
)

func dummy() {
	fmt.Println("sasidaren")
	unicode.IsDigit('9')
}

func lexStart(l *lexer) *Token {
	l.acceptRun(whitespaceSansNewline)
	l.ignore()
	if strings.HasPrefix(l.input[l.pos:], linecomment) {
		return lexLineComment(l)
	}
	if strings.HasPrefix(l.input[l.pos:], spancomment) {
		return lexSpanComment(l, 0)
	}
	if strings.HasPrefix(l.input[l.pos:], multilinequote) {
		l.accept(quote)
		l.accept(quote)
		l.accept(quote)
		l.ignore()
		return lexMultiLineStringIn(l)
	}
	// it could be a symbol literal or a character literal at this point
	if l.peek() == '\'' {
		if l.peekNth(2) == '\'' || (l.peekNth(1) == '\\' && l.peekNth(3) == '\'') {
			l.skip()
			return lexCharacterLiteral(l)
		}
		// TODO(sundaram): handle error case
		// consume single quote
		l.next()
		lexPlainId(l)
		return l.emit(SYMBOL, lexStart)
	}
	if l.accept(letter) {
		l.backup()
		return lexLetter(l)
	}
	if l.accept(num) {
		l.backup()
		return lexNumber(l)
	}
	if l.accept(quote) {
		l.ignore()
		return lexStringIn(l)
	}
	if l.accept(paren) {
		tokenType := parenToTokenType[l.val()]
		switch tokenType {
		case L_PAREN, L_BRACKET, L_CURLY:
			l.regionStack = append(l.regionStack, tokenType)
		case R_PAREN, R_BRACKET, R_CURLY:
			if len(l.regionStack) == 0 {
				return l.emitError("closing paren found without a matching opening bracket", l.val())
			}
			if false == isMatchingParen(l.regionStack[len(l.regionStack)-1], tokenType) {
				return l.emitError("Unmatched paren type", l.val())
			}
			l.regionStack = l.regionStack[:len(l.regionStack)-1]
		}
		return l.emit(tokenType, lexStart)
	}
	if l.accept(semicolon) {
		return l.emit(SEMICOLON, lexStart)
	}
	if l.accept(".") {
		if l.accept(num) {
			return lexNumber(l)
		}
		return l.emit(DOT, lexStart)
	}
	if l.accept(backtick) {
		return lexStringIdIn(l)
	}
	if l.accept(newline) {
		return lexNewline(l)
	}
	if isOpOrDelim(l.input[l.pos:]) {
		return lexOpDelim(l)
	}
	if isOperatorCharacter(l.peek()) {
		lexOp(l)
		return l.emit(IDENTIFIER, lexStart)
	}
	if l.peek() == eof {
		return lexEof(l)
	}

	// leave the curren token and start from the next
	res := l.emitError("unknown lexemes left: ", l.input[l.pos:])
	return res
}

func lexOpDelim(l *lexer) *Token {
	for _, od := range opDelims {
		if strings.HasPrefix(l.input[l.pos:], od) {
			l.pos += len(od)
			return l.emit(OPORDELIM, lexStart)
		}
	}
	return l.emitError("Unexpected extry into lexOpDelim")
}

func lexLineComment(l *lexer) *Token {
	l.acceptRunAllBut(newline)
	l.ignore()
	return l.emit(COMMENT, lexStart)
}

func lexOp(l *lexer) {
	l.acceptRunAllBut(whitespace + alphaLower + alphaUpper + num + paren + delim + newline)
}

func lexSpanComment(l *lexer, level int) *Token {
	l.acceptRunAllBut("*/")
	if l.accept("/") {
		if l.accept("*") {
			return lexSpanComment(l, level+1)
		}
		return lexSpanComment(l, level)
	}
	if l.accept("*") {
		if l.accept("/") {
			if level == 1 {
				return lexStart(l)
			}
			return lexSpanComment(l, level-1)
		}
		return lexSpanComment(l, level)
	}
	return lexStart(l)
}

func lexNewline(l *lexer) *Token {
	l.accept(whitespaceSansNewline)
	// check if the next line is a blank line
	blankLine := l.accept(newline)
	l.acceptRun(whitespace)
	l.ignore()
	if shouldIntroduceNewLine(l) {
		if blankLine {
			return l.emit(NEWLINES, lexStart)
		}
		return l.emit(NEWLINE, lexStart)
	}
	return lexStart(l)
}

func lexPlainId(l *lexer) error {
	if isOperatorCharacter(l.peek()) {
		lexOp(l)
		return nil
	}
	if l.accept(letter) {
		l.acceptRun(letterNum)
		if l.lastAccepted() == '_' {
			// after '_' we could have optional op characters
			lexOp(l)
		}
		return nil
	}
	return fmt.Errorf("Unexpected entry into lexPlainId %s", l.peek())
}

func lexLetter(l *lexer) *Token {
	if l.accept(letter) {
		l.backup()
		err := lexPlainId(l)
		if err == nil {
			if isBoolean(l.val()) {
				return l.emit(BOOLEAN, lexStart)
			}
			if isKeyword(l.val()) {
				return l.emit(keywordsToTokenType[l.val()], lexStart)
			}
			return l.emit(IDENTIFIER, lexStart)
		}
		return l.emitError(err.Error())
	}
	return l.emitError("unknown entry into lexLetter")
	// TODO(sundarama): Change this to emit error
	return nil
}

func lexNumber(l *lexer) *Token {
	if !l.scanNumber() {
		return l.emitError("bad number syntax: %q", l.input[l.start:l.pos])
	}
	return l.emit(NUMBER, lexStart)
}

func (l *lexer) scanNumber() bool {
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept("lL") {
		// this means we are at the end of decimal number
		// Next thing mustn't be alphanumeric.
		if isAlphaNumeric(l.peek()) {
			l.next()
			return false
		}
		return true
	}
	if l.peek() == '.' {
		if isDigit(l.peekNth(1)) {
			l.accept(".")
			l.acceptRun(num)
		} else {
			// TODO(sundarama): handle case of 5.f here
			// return early because the dot is
			// an operator in this usage
			return true
		}
	}
	// exponent part
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	// floatType
	l.accept("fFdD")
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

func lexStringIn(l *lexer) *Token {
	l.acceptRunAllBut(quote + backslash)
	if l.peek() == '\\' {
		lexStringBackslash(l)
		return lexStringIn(l)
	}
	if l.peek() == '"' {
		return l.emit(STRING, lexIgnoreNextCharacter)
	}
	return l.emitError("Unknown state within lexStringIn")
}

func lexStringIdIn(l *lexer) *Token {
	l.acceptRunAllBut(backtick)
	// TODO(sundarama): Investigate if this is needed
	// if l.peek() == '\\' {
	// 	lexStringBackslash(l)
	// 	return lexStringIdIn(l)
	// }
	if l.peek() == '`' {
		return l.emit(IDENTIFIER, lexIgnoreNextCharacter)
	}
	return l.emitError("Unknown state within lexStringIdIn" + l.input[l.pos:])
}

func lexMultiLineStringIn(l *lexer) *Token {
	l.acceptRunAllBut(quote)
	if l.peekNth(0) == '"' && l.peekNth(1) == '"' && l.peekNth(2) == '"' {
		return l.emit(STRING, lexMultiLineStringOut)
	}
	if l.next() != eof {
		return lexMultiLineStringIn(l)
	} else {
		return l.emitError("unexpected end of multi line string")
	}
}

func lexCharacterLiteral(l *lexer) *Token {
	l.accept("\\")
	l.next()
	if l.accept(singlequote) {
		// remove single quote
		l.backup()
		res := l.emit(CHARACTER, lexStart)
		// accept single quote
		l.next()
		return res
	}
	return l.emitError("consumed open single quote; expected end single quote; but found", l.val())
}

// TODO turn \[A-Z] into char code
func lexStringBackslash(l *lexer) {
	l.next() // eat backslash
	l.next() // eat next rune
}

func lexIgnoreNextCharacter(l *lexer) *Token {
	l.next() // eat quote
	l.ignore()
	return lexStart(l)
}

func lexMultiLineStringOut(l *lexer) *Token {
	if strings.HasPrefix(l.input[l.pos:], multilinequote) {
		l.acceptRun(quote)
		l.ignore()
		return lexStart(l)
	}
	return l.emitError("Unknown state within lexMultiStringOut")
}

func lexEof(l *lexer) *Token {
	return l.emit(EOF, nil)
}
