package playlist

import (
	"gocloudcamp/song"
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
	lastSong    *Song
	timer       Timer
}

func NewPlaylist() Playlist {
	singleSong := &Song{}
	return &playlist{
		timer:       NewTimer(),
		currentSong: singleSong,
		lastSong:    singleSong,
	}
}

func (playlist *playlist) Play() {
	if playlist.timer.IsPaused() {
		playlist.timer.Resume()
	} else if playlist.currentSong != nil && playlist.currentSong.Data != nil {
		playlist.timer.Schedule(playlist.currentSong.Data.Length, playlist.Next)
	}
}

func (playlist *playlist) Pause() {
	if !playlist.timer.IsPaused() {
		playlist.timer.Pause()
	}
}

func (playlist *playlist) AddSong(song *song.Song) {
	playlist.lastSong.define(song)
	playlist.lastSong = playlist.lastSong.Next
}

func (playlist *playlist) Next() {
	playlist.timer.Stop()
	if playlist.currentSong.Next != nil {
		playlist.currentSong = playlist.currentSong.Next
	}
	if playlist.currentSong.Data != nil {
		playlist.Play()
	}
}

func (playlist *playlist) Prev() {
	playlist.timer.Stop()
	if playlist.currentSong.Previous != nil {
		playlist.currentSong = playlist.currentSong.Previous
		playlist.Play()
	}
}
