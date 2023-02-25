package song

import "time"

type Song struct {
	Name   string
	Length time.Duration
}

func NewSong(name string, length time.Duration) *Song {
	return &Song{
		Name:   name,
		Length: length,
	}
}
