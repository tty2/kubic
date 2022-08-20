/*
Package css is responsible for parsing css file for theme customization.

This implementation is taken from https://github.com/napsy/go-css
parser.go file

The implementation has been changed a little bit case the original one has bugs and is not maintained.
*/
package css

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/scanner"
)

type tokenType int

const (
	tokenFirstToken tokenType = iota - 1
	tokenBlockStart
	tokenBlockEnd
	tokenValue
	tokenSelector
	tokenStyleSeparator
	tokenStatementEnd
)

// Selector is a string type that represents a CSS rule.
type Selector string

type tokenEntry struct {
	value string
	pos   scanner.Position
}

type tokenizer struct {
	s *scanner.Scanner
}

// Type returns the rule type, which can be a class, id or a tag.
func (rule Selector) Type() string {
	if strings.HasPrefix(string(rule), ".") {
		return "class"
	}
	if strings.HasPrefix(string(rule), "#") {
		return "id"
	}

	return "tag"
}

func (e tokenEntry) typ() tokenType {
	return newTokenType(e.value)
}

// nolint gocognit: not my implementation. not so important part on start TODO refactor
func (t *tokenizer) next() (tokenEntry, error) {
	token := t.s.Scan()
	if token == scanner.EOF {
		return tokenEntry{}, errors.New("EOF")
	}
	value := t.s.TokenText()
	pos := t.s.Pos()
	if newTokenType(value).String() == "STYLE_SEPARATOR" {
		t.s.IsIdentRune = func(ch rune, i int) bool { // property value can contain spaces
			if ch == -1 || ch == '\n' || ch == '\t' || ch == ':' || ch == ';' {
				return false
			}

			return true
		}
	} else {
		t.s.IsIdentRune = func(ch rune, i int) bool { // other tokens can't contain spaces
			if ch == -1 || ch == '.' || ch == '#' || ch == '\n' || ch == ' ' || ch == '\t' || ch == ':' || ch == ';' {
				return false
			}

			return true
		}
	}

	return tokenEntry{
		value: value,
		pos:   pos,
	}, nil
}

func (t tokenType) String() string {
	switch t {
	case tokenBlockStart:
		return "BLOCK_START"
	case tokenBlockEnd:
		return "BLOCK_END"
	case tokenStyleSeparator:
		return "STYLE_SEPARATOR"
	case tokenStatementEnd:
		return "STATEMENT_END"
	case tokenSelector:
		return "SELECTOR"
	}

	return "VALUE"
}

func newTokenType(typ string) tokenType {
	switch typ {
	case "{":
		return tokenBlockStart
	case "}":
		return tokenBlockEnd
	case ":":
		return tokenStyleSeparator
	case ";":
		return tokenStatementEnd
	case ".", "#":
		return tokenSelector
	}

	return tokenValue
}

func newTokenizer(r io.Reader) *tokenizer {
	s := &scanner.Scanner{}
	s.Init(r)

	return &tokenizer{
		s: s,
	}
}

func buildList(r io.Reader) *list.List {
	l := list.New()
	t := newTokenizer(r)
	for {
		token, err := t.next()
		if err != nil {
			break
		}
		l.PushBack(token)
	}

	return l
}

// nolint funlen: not my implementation. not so important part on start. TODO refactor
func parse(l *list.List) (map[Selector]map[string]string, error) {

	var (
		// Information about the current block that is parsed.
		selectors []string
		style     string
		value     string
		selector  string

		isBlock bool

		// Parsed styles.
		css    = make(map[Selector]map[string]string)
		styles = make(map[string]string)

		// Previous token for the state machine.
		prevToken = tokenType(tokenFirstToken)
	)

	for e := l.Front(); e != nil; e = l.Front() {
		token := e.Value.(tokenEntry)

		l.Remove(e)
		switch token.typ() {
		case tokenValue:
			switch prevToken {
			case tokenFirstToken, tokenBlockEnd:
				selectors = append(selectors, token.value)
			case tokenSelector:
				selectors = append(selectors, selector+token.value)
			case tokenBlockStart, tokenStatementEnd:
				style = token.value
			case tokenStyleSeparator:
				value = token.value
			case tokenValue:
				selectors = append(selectors, token.value)
			default:
				return css, fmt.Errorf("line %d: invalid syntax", token.pos.Line)
			}
		case tokenSelector:
			selector = token.value
		case tokenBlockStart:
			if prevToken != tokenValue {
				return css, fmt.Errorf("line %d: block is missing rule identifier", token.pos.Line)
			}
			isBlock = true
		case tokenStatementEnd:
			if prevToken != tokenValue || style == "" || value == "" {
				return css, fmt.Errorf("line %d: expected style before semicolon", token.pos.Line)
			}
			styles[style] = value
		case tokenBlockEnd:
			if !isBlock {
				return css, fmt.Errorf("line %d: rule block ends without a beginning", token.pos.Line)
			}
			for i := range selectors {
				oldSelector, ok := css[Selector(selectors[i])]
				if ok {
					// merge selectors
					for style, value := range oldSelector {
						if _, ok := styles[style]; !ok {
							styles[style] = value
						}
					}
				}
				css[Selector(selectors[i])] = styles

			}
			styles = map[string]string{}
			style, value = "", ""
			selectors = []string{}
			isBlock = false
		}
		prevToken = token.typ()
	}

	return css, nil
}

// Unmarshal will take a byte slice, containing sylesheet rules and return
// a map of a rules map.
func Unmarshal(b []byte) (map[Selector]map[string]string, error) {
	return parse(buildList(bytes.NewReader(b)))
}

// CSSStyle returns an error-checked parsed style, or an error if the
// style is unknown. Most of the styles are not supported yet.
func CSSStyle(name string, styles map[string]string) string {
	value, ok := styles[name]
	if !ok {
		return ""
	}

	return value
}
