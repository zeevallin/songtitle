package songtitle_test

import (
	"fmt"
	"testing"

	"github.com/zeevallin/songtitle"
)

func ExampleParse() {
	// Get song title from a file name
	song := songtitle.Parse("Foo Fighters • Everlong (CC) 🎤 [Karaoke] [Instrumental Lyrics]")
	fmt.Printf("%s - %s\n", song.Artist, song.Title)
	for _, tag := range song.Tags {
		fmt.Printf("\t- %s\n", tag)
	}
	// Output:
	// Foo Fighters - Everlong
	// 	- CC
	//	- Karaoke
	//	- Instrumental Lyrics
	//
}

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		want  songtitle.Song
	}{
		{
			input: "Foo Fighters • Everlong (CC) 🎤 [Karaoke] [Instrumental Lyrics]",
			want:  songtitle.Song{Title: "Everlong", Artist: "Foo Fighters", Tags: []string{"CC", "Karaoke", "Instrumental Lyrics"}},
		},
		{
			input: "Michael Bublé - Feeling Good (Karaoke Version)",
			want:  songtitle.Song{Title: "Feeling Good", Artist: "Michael Bublé", Tags: []string{"Karaoke Version"}},
		},
		{
			input: "The Beatles - Hey Jude (Karaoke)",
			want:  songtitle.Song{Title: "Hey Jude", Artist: "The Beatles", Tags: []string{"Karaoke"}},
		},
		{
			input: "Ebba Grön- Staten och kapitalet. Karaoke",
			want:  songtitle.Song{Title: "Staten och kapitalet. Karaoke", Artist: "Ebba Grön", Tags: []string{}},
		},
		{
			input: "Eminmen ft. Rihanna - Love The Way You Lie (Karaoke Version)",
			want:  songtitle.Song{Title: "Love The Way You Lie", Artist: "Eminmen ft. Rihanna", Tags: []string{"Karaoke Version"}},
		},
		{
			input: "The Beatles - Hey Jude (Karaoke) [Official",
			want:  songtitle.Song{Title: "Hey Jude", Artist: "The Beatles", Tags: []string{"Karaoke"}},
		},
		{
			input: "The Beatles - Hey Jude (Karaoke) [Official • 4K]",
			want:  songtitle.Song{Title: "Hey Jude", Artist: "The Beatles", Tags: []string{"Karaoke"}},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := songtitle.Parse(tt.input)
			if got.Title != tt.want.Title {
				t.Errorf("got %q; want %q", got.Title, tt.want.Title)
			}
			if got.Artist != tt.want.Artist {
				t.Errorf("got %q; want %q", got.Artist, tt.want.Artist)
			}
			if len(got.Tags) != len(tt.want.Tags) {
				t.Errorf("got %v; want %v", got.Tags, tt.want.Tags)
			}
			for i := range got.Tags {
				if got.Tags[i] != tt.want.Tags[i] {
					t.Errorf("got %q; want %q", got.Tags[i], tt.want.Tags[i])
				}
			}
		})
	}
}
