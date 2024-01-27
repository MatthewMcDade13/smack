package interp

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	TOK_LPAREN = iota
	TOK_RPAREN
	TOK_LITERAL
	TOK_ATOM
)

type parser struct {
	toks    []string
	current uint32
}

func new_parser(source string) parser {
	toks := tokenize(source)
	current := uint32(0)
	return parser{
		toks,
		current,
	}
}

func (p *parser) peek() string {
	return p.toks[p.current]
}

func (p *parser) peek_next() string {
	return p.toks[p.current+1]
}

func (p *parser) skip(n int) {
	p.current += uint32(n)
}

func (p *parser) read_form() (Value, error) {
	tok := p.peek()

	if tok[0] == '(' {
		p.skip(1)
		return p.read_list()
	} else {
		return p.read_atom()
	}
}

func (p *parser) read_atom() (Value, error) {
	tok := p.peek()
	switch tok[0] {
	case '"':
		return NewString(tok), nil
	case ':':
		return NewSymbol(Symbol(tok)), nil
	default:
		n, err := strconv.ParseFloat(tok, 32)
		if err != nil {
			// TODO :: Probably want to put this in a global keyword map
			switch tok {
			case "true":
				return NewBool(true), nil
			case "false":
				return NewBool(false), nil
			default:
				return NewSymbol(Symbol(tok)), nil
			}
		}
		return NewNumber(n), nil
	}
}

func (p *parser) read_list() (Value, error) {
	list := make([]Value, 0)

	for {

		// TODO :: This is kind of a quick and dirty way to prevent
		// peeking beyond buffer of tokens. I hate it. Need to debug
		// where in my logic im incorrect so we dont have to do this check. Or at
		// least find a better way of doing check? idk i need to go to bed...
		if int(p.current) >= len(p.toks) {
			// End of input
			return NoValue(), fmt.Errorf("Missing matching end parenthesis ')'")
		}

		tok := p.peek()
		if tok[0] == ')' {
			break
		}

		if v, err := p.read_form(); err == nil {
			list = append(list, v)
		} else {
			return NoValue(), err
		}

		p.skip(1)
	}

	return NewList(list), nil
}

func tokenize(source string) []string {
	re := regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" + `~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" + `,;)]*)`)
	matchesRaw := re.FindAll([]byte(source), -1)

	matches := make([]string, 0, len(matchesRaw))
	for _, s := range matchesRaw {
		matches = append(matches, strings.TrimSpace(string(s)))
	}
	return matches
}
