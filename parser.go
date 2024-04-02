package songtitle

import (
	"strings"
)

type parserFn func(*parser) parserFn

type parser struct {
	input  []token
	pos    int
	start  int
	artist *string
	title  *string
	tags   []string
}

func parse(tokens ...token) *parser {
	p := &parser{
		input: tokens,
	}
	p.run()
	return p
}

func (p *parser) run() {
	for state := parseExpr; state != nil; {
		state = state(p)
	}
}

func (p *parser) next() token {
	if p.pos >= len(p.input) {
		return token{}
	}
	t := p.input[p.pos]
	p.pos++
	return t
}

func (p *parser) backup() {
	p.pos--
}

func (p *parser) ignore() {
	p.start = p.pos
}

func parseExpr(p *parser) parserFn {
	for {
		if p.pos >= len(p.input) {
			break
		}
		switch t := p.next(); t.kind {
		case kindWord:
			return parseBareWords
		case kindSep, kindSpace, kindClose:
			p.ignore()
		case kindOpen:
			return parseTag
		default:
			return nil
		}
	}
	return nil
}

func parseBareWords(p *parser) parserFn {
	for {
		switch t := p.next(); t.kind {
		case kindWord, kindSpace:
			continue
		default:
			p.backup()

			v := summiseBareWords(p.input[p.start:p.pos]...)
			if p.title == nil && p.artist == nil {
				p.title = &v
			} else if p.artist == nil && p.title != nil {
				p.artist = p.title
				p.title = &v
			}

			p.start = p.pos
			return parseExpr
		}
	}
}

func parseTag(p *parser) parserFn {
	for {
		switch t := p.next(); t.kind {
		case kindWord, kindSpace:
			continue
		case kindClose:
			p.tags = append(p.tags, summiseBareWords(p.input[p.start+1:p.pos-1]...))
			p.start = p.pos
			return parseExpr
		default:
			return nil
		}
	}
}

func summiseBareWords(tokens ...token) string {
	var s string
	for _, t := range tokens {
		s += t.literal
	}
	return strings.TrimSpace(s)
}
