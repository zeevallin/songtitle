package songtitle

type Song struct {
	Title  string
	Artist string
	Tags   []string
}

func Parse(s string) Song {
	l := lex(s)
	var tokens []token
	for {
		t, more := l.nextItem()
		if !more {
			break
		}
		tokens = append(tokens, t)
	}

	p := parse(tokens...)

	var title string
	if p.title != nil {
		title = *p.title
	}

	var artist string
	if p.artist != nil {
		artist = *p.artist
	}

	return Song{
		Title:  title,
		Artist: artist,
		Tags:   p.tags,
	}
}
