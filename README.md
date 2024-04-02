# Songtitle
Library to get a song title, artist and other tags from a file name or a youtube video.

## Usage

```go

import (
    "fmt"
    "github.com/zeevallin/songtitle"
)

func main() {
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
```