package parser

import (
	"fmt"
	"testing"
)

var lexTests = []struct {
	input    string
	expected []TokenType
}{
	{"// what is you name\nidentifier", []TokenType{COMMENT, IDENTIFIER}},
	{"__system", []TokenType{IDENTIFIER}},
	{"_MAX_LEN_", []TokenType{IDENTIFIER}},
	{"`yield`", []TokenType{IDENTIFIER}},
	{"αρετη", []TokenType{IDENTIFIER}},
	{"αρετη()", []TokenType{IDENTIFIER, L_PAREN, R_PAREN}},
	{"'αρετη()", []TokenType{SYMBOL, L_PAREN, R_PAREN}},
	{"// what is you name\nidentifier identifier", []TokenType{COMMENT, IDENTIFIER, IDENTIFIER}},
	{"whiles while", []TokenType{IDENTIFIER, WHILE}},
	{"/*sundaram is an idiot\n whatever works for you\n india is my country\n", []TokenType{}},
	{"/*sundaram is an idiot\n whatever works for you\n india is my country\n*/", []TokenType{}},
	{"/*sundaram is an idiot\n whatever works for you\n india is my country\n*/identifier", []TokenType{IDENTIFIER}},
	{"dot_product_* __system", []TokenType{IDENTIFIER, IDENTIFIER}},
	{"0.1234", []TokenType{NUMBER}},
	{"/*sundaram is an idiot\n whatever works for you\n india is my country\n*/123l", []TokenType{NUMBER}},
	{"0l", []TokenType{NUMBER}},
	{"1e+4", []TokenType{NUMBER}},
	{"1e+e", []TokenType{ERROR}},
	{"\"Hello,\nWorld!\"", []TokenType{STRING}},
	{"\"This string contains a \\\" character.\"", []TokenType{STRING}},
	{`  """the present string
	    spans three
	    lines."""`, []TokenType{STRING}},
	{`  """the present string
	    spans three
	    lines."""""""""""""sundaram`, []TokenType{STRING, IDENTIFIER}},
	{"truee true", []TokenType{IDENTIFIER, BOOLEAN}},
	{"'a''\\n'", []TokenType{CHARACTER, CHARACTER}},
	{"1e30f", []TokenType{NUMBER}},
	{"3.14159f", []TokenType{NUMBER}},
	{"1.0e-100", []TokenType{NUMBER}},
	{".1", []TokenType{NUMBER}},
	{"_;", []TokenType{IDENTIFIER, SEMICOLON}},
	{";", []TokenType{SEMICOLON}},
	{"22.`yield`", []TokenType{NUMBER, DOT, IDENTIFIER}},
	{`
		  val vcard3StrMoreAttributes =
    """
      |BEGIN:VCARD
      |VERSION:3.0
      |N:Gump;Forrest;Mr.
      |FN:Forrest Gump
      |ORG:Bubba Gump Shrimp Co.
      |TITLE:Shrimp Man
      |PHOTO;VALUE=URL;TYPE=GIF:http://www.example.com/dir_photos/my_photo.gif
      |TEL;TYPE=WORK,VOICE:+1-111-555-1212
      |TEL;TYPE=HOME,VOICE:+1-404-555-1212
      |EMAIL;TYPE=PREF,INTERNET:forrestgump@example.com
      |ADR;TYPE=work;LABEL="100 Waters Edge\nBaytown, LA 30314\nUnited States of America"
      |  :;;100 Waters Edge;Baytown;LA;30314;United States of America
      |ADR;TYPE=home;LABEL="42 Plantation St.\nBaytown, LA 30314\nUnited States of America"
      | :;;42 Plantation St.;Baytown;LA;30314;United States of America
      |
      |REV:2008-04-24T19:52:43Z
      |END:VCARD
    """.stripMargin`, []TokenType{VAL, IDENTIFIER, OPORDELIM, STRING, DOT, IDENTIFIER}},
	{`
  /* The following joins are generated with this code:
  scala -e '
  val meths = for (end <- ''b'' to ''v''; ps = ''a'' to end) yield
      """/**
 * Join %d futures. The returned future is complete when all
 * underlying futures complete. It fails immediately if any of them
 * do.
 */
def join[%s](%s): Future[(%s)] = join(Seq(%s)) map { _ => (%s) }""".format(
        ps.size,
        ps map (_.toUpper) mkString ",",
        ps map(p => "%c: Future[%c]".format(p, p.toUpper)) mkString ",",
        ps map (_.toUpper) mkString ",",
        ps mkString ",",
        ps map(p => "Await.result("+p+")") mkString ","
      )

  meths foreach println
  '
  */Identifier`, []TokenType{IDENTIFIER}},
	{"'defined", []TokenType{SYMBOL}},
}

func getTokenTypes(tokens []*Token) []TokenType {
	var res []TokenType
	for _, token := range tokens {
		res = append(res, token.Typ)
	}
	return res
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		lexer := Lexer(test.input)
		tokens := lexer.LexTillDone()

		strSlice1 := fmt.Sprintf("%v", getTokenTypes(tokens))
		strSlice2 := fmt.Sprintf("%v", test.expected)
		// fmt.Printf("returned = %s, expected = %s\n", strSlice1, strSlice2)
		if strSlice1 != strSlice2 {
			t.Errorf("Lex(%s) = %s, Expected = %s", test.input, tokens, test.expected)
		}
	}
}
