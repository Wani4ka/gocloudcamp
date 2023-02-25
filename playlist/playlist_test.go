package playlist

import (
	"fmt"
	songmodule "gocloudcamp/song"
	"math/rand"
	"testing"
	"time"
)

func validateSong(got, want *songmodule.Song) bool {
	return got != nil && want.Equal(*got)
}

func validateCurrentSong(pl *playlist, want *songmodule.Song) bool {
	return pl != nil && pl.currentSong != nil && validateSong(pl.currentSong.Data, want)
}

func TestNewPlaylist(t *testing.T) {
	got, ok := NewPlaylist().(*playlist)
	if !ok {
		t.Fatal("Couldn't convert NewPlaylist result to a *playlist")
	}
	if got.currentSong != nil && got.currentSong.Data != nil {
		t.Fatal("Freshly created playlist is initialized with a song")
	}
	if got.timer == nil {
		t.Fatal("Playlist created without a timer")
	}
}

func TestPlaylist_AddSong(t *testing.T) {
	pl, ok := NewPlaylist().(*playlist)
	if !ok {
		t.Fatal("Couldn't convert NewPlaylist result to a *playlist")
	}
	song := songmodule.NewSong("Test song", time.Second*30)
	pl.AddSong(song)
	if pl.currentSong.Data != song {
		t.Fatal("Data of the song added is not equal to the original")
	}
}

func TestPlaylist_AddSong_many(t *testing.T) {
	pl, ok := NewPlaylist().(*playlist)
	if !ok {
		t.Fatal("Couldn't convert NewPlaylist result to a *playlist")
	}
	const amount = 1000
	var songs []*songmodule.Song
	for i := 0; i < amount; i++ {
		songs = append(songs, songmodule.NewSong(fmt.Sprintf("Test song %v", i), time.Duration(rand.Int63n(int64(time.Hour)))))
	}
	start := time.Now()
	for _, song := range songs {
		pl.AddSong(song)
	}
	t.Logf("%v elapsed on adding %v songs to playlist", time.Now().Sub(start), amount)
	lastSong := pl.currentSong
	for i, song := range songs {
		if lastSong == nil || lastSong.Data == nil || !song.Equal(*lastSong.Data) {
			t.Fatalf("Song %v doesn't match the one added to playlist", i)
		}
		lastSong = lastSong.Next
	}
	if lastSong != nil && lastSong.Data != nil {
		t.Fatalf("Playlist has more songs that were added")
	}
}

func TestPlaylist_Play(t *testing.T) {
	const grace = 50 * time.Millisecond
	pl, ok := NewPlaylist().(*playlist)
	if !ok {
		t.Fatal("Couldn't convert NewPlaylist result to a *playlist")
	}
	const amount = 4
	var songs []*songmodule.Song
	for i := 0; i < amount; i++ {
		songs = append(songs, songmodule.NewSong(fmt.Sprintf("Test song %v", i), time.Duration(rand.Int63n(int64(3*time.Second)))+500*time.Millisecond))
		pl.AddSong(songs[i])
	}
	if !validateCurrentSong(pl, songs[0]) {
		t.Fatal("Playlist didn't initialize the first song")
	}
	pl.Play()
	time.Sleep(songs[0].Length + grace)
	if !validateCurrentSong(pl, songs[1]) {
		t.Fatal("Playlist didn't start the next song after the first one should've done playing")
	}
}

func TestPlaylist_Pause(t *testing.T) {
	const grace = 50 * time.Millisecond
	pl, ok := NewPlaylist().(*playlist)
	if !ok {
		t.Fatal("Couldn't convert NewPlaylist result to a *playlist")
	}
	const amount = 10
	var songs []*songmodule.Song
	for i := 0; i < amount; i++ {
		songs = append(songs, songmodule.NewSong(fmt.Sprintf("Test song %v", i), time.Duration(rand.Int63n(int64(1*time.Second)))+200*time.Millisecond))
		pl.AddSong(songs[i])
	}
	if !validateCurrentSong(pl, songs[0]) {
		t.Fatal("Playlist didn't initialize the first song")
	}
	pl.Play()
	pl.Pause()
	time.Sleep(songs[0].Length + grace)
	if !validateCurrentSong(pl, songs[0]) {
		t.Fatal("Song was lost or not paused")
	}
	pl.Play()
	time.Sleep(songs[0].Length + grace)
	if !validateCurrentSong(pl, songs[1]) {
		t.Fatal("Playlist didn't start the next song after the first one should've done playing")
	}
}

func TestPlaylist_Seek(t *testing.T) {
	pl, ok := NewPlaylist().(*playlist)
	if !ok {
		t.Fatal("Couldn't convert NewPlaylist result to a *playlist")
	}
	const amount = 10
	var songs []*songmodule.Song
	for i := 0; i < amount; i++ {
		songs = append(songs, songmodule.NewSong(fmt.Sprintf("Test song %v", i), time.Duration(rand.Int63n(int64(1*time.Second)))+200*time.Millisecond))
		pl.AddSong(songs[i])
	}
	if !validateCurrentSong(pl, songs[0]) {
		t.Fatal("Playlist didn't initialize the first song")
	}
	for i := 1; i < amount; i++ {
		pl.Next()
		if !validateCurrentSong(pl, songs[i]) {
			t.Fatalf("Playlist should play %v-th song after %v Next() method calls", i+1, i)
		}
	}
	if pl.currentSong.Next != nil && pl.currentSong.Next.Data != nil {
		t.Fatal("Current song in playlist is not last, but should be")
	}
	for i := amount - 2; i > -1; i-- {
		pl.Prev()
		if !validateCurrentSong(pl, songs[i]) {
			t.Fatalf("Playlist should play %v-th song after %v Prev() method calls", i+1, amount-i+1)
		}
	}
	if pl.currentSong.Previous != nil && pl.currentSong.Previous.Data != nil {
		t.Fatal("Current song in playlist is not first, but should be")
	}
}
