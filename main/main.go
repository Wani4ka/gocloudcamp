package main

import (
	"gocloudcamp/playlist"
	"gocloudcamp/song"
	"time"
)

func main() {
	list := playlist.NewPlaylist()
	list.AddSong(song.NewSong("Song 1", 17*time.Second))
	list.AddSong(song.NewSong("Song 2", 5*time.Second))
	list.Play()
	time.Sleep(3 * time.Minute)
	list.Play()
}
