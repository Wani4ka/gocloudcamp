package playlist

import (
	"gocloudcamp/song"
	"log"
	"os"
)

type Song struct {
	Previous *Song
	Data     *song.Song
	Next     *Song
}

func (song *Song) define(data *song.Song) {
	song.Data = data
	if song.Next == nil {
		song.Next = &Song{Previous: song}
	}
}

type playlist struct {
	currentSong *Song
	timer       Timer
	logger      *log.Logger
}

func NewPlaylist() Playlist {
	return &playlist{
		timer:       NewTimer(),
		currentSong: &Song{},
		logger:      log.New(os.Stderr, "PLAYLIST\t", log.Ltime),
	}
}

func (playlist *playlist) Play() {
	if playlist.timer.IsPaused() {
		playlist.timer.Resume()
	} else if playlist.currentSong != nil && playlist.currentSong.Data != nil {
		playlist.timer.Schedule(playlist.currentSong.Data.Length, playlist.Next)
		playlist.logger.Printf("Song %v (%v) started playing", playlist.currentSong.Data.Name, playlist.currentSong.Data.Length)
	}
}

func (playlist *playlist) Pause() {
	if !playlist.timer.IsPaused() {
		playlist.timer.Pause()
		playlist.logger.Printf("Song %v (%v) is now paused", playlist.currentSong.Data.Name, playlist.currentSong.Data.Length)
	}
}

func (playlist *playlist) AddSong(song *song.Song) {
	lastSong := playlist.currentSong
	for lastSong.Data != nil {
		lastSong = lastSong.Next
	}
	lastSong.define(song)
	playlist.logger.Printf("Song %v (%v) added to the end of playlist", lastSong.Data.Name, lastSong.Data.Length)
}

func (playlist *playlist) Next() {
	playlist.timer.Stop()
	if playlist.currentSong.Next != nil {
		playlist.currentSong = playlist.currentSong.Next
	}
	if playlist.currentSong.Data != nil {
		playlist.logger.Printf("Playing next song: %v (%v)", playlist.currentSong.Data.Name, playlist.currentSong.Data.Length)
		playlist.Play()
	} else {
		playlist.logger.Println("No next song to play")
	}
}

func (playlist *playlist) Prev() {
	playlist.timer.Stop()
	if playlist.currentSong.Previous != nil {
		playlist.currentSong = playlist.currentSong.Previous
		playlist.Play()
		playlist.logger.Printf("Playing previous song: %v (%v)", playlist.currentSong.Data.Name, playlist.currentSong.Data.Length)
	} else {
		playlist.logger.Println("No previous song to play")
	}
}
