package music

import (
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

//Music - the music object
type Music struct {
	streamer beep.StreamSeekCloser
	format   beep.Format
}

//Load - load a track
func (m *Music) Load(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	m.streamer, m.format, err = mp3.Decode(f)
	if err != nil {
		panic(err)
	}

}

//LoadRandom - load a random track
func (m *Music) LoadRandom() {
	matches, err := filepath.Glob("music/*.mp3")
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	choice := rand.Intn(len(matches))
	f, err := os.Open(matches[choice])
	if err != nil {
		panic(err)
	}
	m.streamer, m.format, err = mp3.Decode(f)
	if err != nil {
		panic(err)
	}

}

//Play - plays the loaded song
func (m *Music) Play() {
	speaker.Init(m.format.SampleRate, m.format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	ctrl := &beep.Ctrl{Streamer: beep.Seq(m.streamer, beep.Callback(func() {
		done <- true
	})), Paused: false}
	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     10,
		Volume:   -10 / 10,
		Silent:   false,
	}
	speaker.Play(volume)
	<-done
}

//Close - closes the file
func (m *Music) Close() {
	m.streamer.Close()
}
