package songtitle

import (
	"fmt"
	"slices"
	"unicode/utf8"
)

type kind string

const (
	kindWord  kind = "word"
	kindSep   kind = "sep"
	kindOpen  kind = "open"
	kindClose kind = "close"
	kindSpace kind = "space"
	kindEOF   kind = "eof"
)

type token struct {
	kind    kind
	literal string
}

func (t token) String() string {
	return fmt.Sprintf("{%s %q}", t.kind, t.literal)
}

func lex(input string) *lexer {
	l := &lexer{
		input:  input,
		tokens: make(chan token),
	}
	go l.run()
	return l
}

type stateFn func(*lexer) stateFn

type lexer struct {
	input  string
	start  int
	pos    int
	width  int
	tokens chan token
}

func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *lexer) nextItem() (token, bool) {
	item, closed := <-l.tokens
	return item, closed
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return 0
	}
	r, w := rune(l.input[l.pos]), 1
	if r >= 0x80 {
		r, w = utf8.DecodeRuneInString(l.input[l.pos:])
	}
	l.width = w
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) emit(kind kind) {
	l.tokens <- token{kind, l.input[l.start:l.pos]}
	l.start = l.pos
}

func lexText(l *lexer) stateFn {
	for {
		if l.pos >= len(l.input) {
			l.emit(kindEOF)
			break
		}
		switch r := l.next(); {
		case isSep(r):
			l.emit(kindSep)
		case isSpace(r):
			l.emit(kindSpace)
		case isOpenTag(r):
			l.emit(kindOpen)
		case isCloseTag(r):
			l.emit(kindClose)
		case isNonReserved(r):
			return lexWord
		default:
			return nil
		}
	}
	return nil
}

func lexWord(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isNonReserved(r) && !isEOF(r):
			continue
		default:
			l.backup()
			l.emit(kindWord)
			return lexText
		}
	}
}

func isNonReserved(r rune) bool {
	return !isSep(r) && !isOpenTag(r) && !isCloseTag(r) && !isSpace(r)
}

func isSep(r rune) bool {
	runes := []rune{'•', '-', '–', '—', '―', '‒', '−', '˗', '⁃', '∼', '∾', '∿', '〜', '〰', '﹏', '︱', '︲', '﹍', '﹎', '﹉', '﹊', '﹋', '﹌', '＋', '－', '／', '＼', '＝', '＿', '～', '｜', '￨', '￩', '￪', '￫', '￬', '￭', '￮', '﹢', '﹣', '﹤', '﹥', '﹦', '﹨', '﹩', '﹪', '﹫'}
	return slices.Contains[[]rune](runes, r)
}

func isOpenTag(r rune) bool {
	runes := []rune{'(', '[', '{', '<', '«', '“', '„', '‹', '「', '『', '〈', '《', '【', '〔', '〖', '〘', '〚', '〝', '〟', '﹁', '﹃', '［', '｛', '｟', '｢'}
	return slices.Contains[[]rune](runes, r)
}

func isCloseTag(r rune) bool {
	runes := []rune{')', ']', '}', '>', '»', '”', '„', '›', '」', '』', '〉', '》', '】', '〕', '〗', '〙', '〛', '〞', '〟', '﹂', '﹄', '］', '｝', '｠', '｣'}
	return slices.Contains[[]rune](runes, r)
}

func isSpace(r rune) bool {
	runes := []rune{' ', '\t', '\n', '\v', '\f', '\r', 0x85, 0xA0, 0x1680, 0x2000, 0x2001, 0x2002, 0x2003, 0x2004, 0x2005, 0x2006, 0x2007, 0x2008, 0x2009, 0x200A, 0x2028, 0x2029, 0x202F, 0x205F, 0x3000, 0xFEFF}
	return slices.Contains[[]rune](runes, r)
}

func isEOF(r rune) bool {
	return r == 0
}
