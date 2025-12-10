// Package music plays music.
package music

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
)

// MusicType defines music type.
type MusicType int

const (
	// TypeOgg is type ogg.
	TypeOgg MusicType = iota

	// TypeMP3 is type MP3.
	TypeMP3
)

// SampleRate defines sample rate.
const SampleRate = 48000

func (t MusicType) String() string {
	switch t {
	case TypeOgg:
		return "Ogg"
	case TypeMP3:
		return "MP3"
	default:
		panic("not reached")
	}
}

// Player represents the current audio state.
type Player struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	seBytes      []byte
	seCh         chan []byte
	volume128    int
	musicType    MusicType
}

// NewPlayer creates a music player.
func NewPlayer(audioContext *audio.Context, musicType MusicType, src io.Reader) (*Player, error) {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}

	// bytesPerSample is the byte size for one sample (8 [bytes] = 2 [channels] * 4 [bytes] (32bit float)).
	// TODO: This should be defined in audio package.
	const bytesPerSample = 8
	var s audioStream

	switch musicType {
	case TypeOgg:
		var err error
		s, err = vorbis.DecodeF32(src)
		if err != nil {
			return nil, err
		}
	case TypeMP3:
		var err error
		s, err = mp3.DecodeF32(src)
		if err != nil {
			return nil, err
		}
	default:
		panic("not reached")
	}
	p, err := audioContext.NewPlayerF32(s)
	if err != nil {
		return nil, err
	}
	player := &Player{
		audioContext: audioContext,
		audioPlayer:  p,
		total:        time.Second * time.Duration(s.Length()) / bytesPerSample / SampleRate,
		volume128:    128,
		seCh:         make(chan []byte),
		musicType:    musicType,
	}
	if player.total == 0 {
		player.total = 1
	}

	player.audioPlayer.Play()
	go func() {
		s, err := wav.DecodeF32(bytes.NewReader(raudio.Jab_wav))
		if err != nil {
			log.Fatal(err)
			return
		}
		b, err := io.ReadAll(s)
		if err != nil {
			log.Fatal(err)
			return
		}
		player.seCh <- b
	}()
	return player, nil
}

// Close closes the player.
func (p *Player) Close() error {
	return p.audioPlayer.Close()
}

// Update updates the player.
func (p *Player) Update() error {
	select {
	case p.seBytes = <-p.seCh:
		close(p.seCh)
		p.seCh = nil
	default:
	}

	if p.audioPlayer.IsPlaying() {
		p.current = p.audioPlayer.Position()
	}

	return nil
}
