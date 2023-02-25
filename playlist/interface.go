package playlist

import (
	"gocloudcamp/song"
)

type Playlist interface {
	Play()
	Pause()
	AddSong(song *song.Song)
	Next()
	Prev()
}
