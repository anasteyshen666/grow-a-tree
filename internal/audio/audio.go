// Package audio synthesizes simple 8-bit sound effects and a calm looping
// chiptune, all procedurally — no asset files needed.
package audio

import (
	"bytes"
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const sampleRate = 44100

var (
	ctx     *audio.Context
	players map[string]*audio.Player
	music   *audio.Player
)

// Init sets up the audio context and prepares every sound. Safe to call once.
func Init() {
	if ctx != nil {
		return
	}
	ctx = audio.NewContext(sampleRate)
	players = map[string]*audio.Player{
		"click": ctx.NewPlayerFromBytes(clickPCM()),
		"grow":  ctx.NewPlayerFromBytes(growPCM()),
		"rot":   ctx.NewPlayerFromBytes(rotPCM()),
		"eat":   ctx.NewPlayerFromBytes(eatPCM()),
	}
	for _, p := range players {
		p.SetVolume(0.45)
	}
	m := musicPCM()
	loop := audio.NewInfiniteLoop(bytes.NewReader(m), int64(len(m)))
	if p, err := ctx.NewPlayer(loop); err == nil {
		music = p
		music.SetVolume(0.35)
	}
}

func play(name string) {
	if p := players[name]; p != nil {
		p.Rewind()
		p.Play()
	}
}

func PlayClick() { play("click") }
func PlayGrow()  { play("grow") }
func PlayRot()   { play("rot") }
func PlayEat()   { play("eat") }

// StartMusic begins (or resumes) the looping background music.
func StartMusic() {
	if music != nil && !music.IsPlaying() {
		music.Play()
	}
}

// --- synthesis ---

type wave func(cycles float64) float64
type env func(p float64) float64

func square(c float64) float64 {
	if c-math.Floor(c) < 0.5 {
		return 1
	}
	return -1
}

func triangle(c float64) float64 {
	f := c - math.Floor(c)
	return 4*math.Abs(f-0.5) - 1
}

// pluck is a fast-attack, linear-decay envelope for short blips.
func pluck(p float64) float64 {
	if p < 0.05 {
		return p / 0.05
	}
	return (1 - p) / 0.95
}

// pad is a gentle attack/sustain/release envelope for the music.
func pad(p float64) float64 {
	const a, r = 0.12, 0.25
	switch {
	case p < a:
		return p / a
	case p > 1-r:
		return (1 - p) / r
	default:
		return 1
	}
}

type synth struct{ buf []byte }

// tone appends one note (16-bit stereo PCM) to the buffer.
func (s *synth) tone(freq, dur, vol float64, w wave, e env) {
	n := int(sampleRate * dur)
	for i := 0; i < n; i++ {
		v := w(freq*float64(i)/sampleRate) * vol * e(float64(i)/float64(n))
		if v > 1 {
			v = 1
		} else if v < -1 {
			v = -1
		}
		q := int16(v * 32767)
		lo, hi := byte(q), byte(q>>8)
		s.buf = append(s.buf, lo, hi, lo, hi)
	}
}

func clickPCM() []byte {
	s := &synth{}
	s.tone(880, 0.05, 0.35, square, pluck)
	return s.buf
}

func growPCM() []byte {
	s := &synth{}
	s.tone(523.25, 0.04, 0.30, triangle, pluck)
	s.tone(783.99, 0.07, 0.30, triangle, pluck)
	return s.buf
}

func rotPCM() []byte {
	s := &synth{}
	s.tone(174.61, 0.14, 0.35, square, pluck)
	return s.buf
}

func eatPCM() []byte {
	s := &synth{}
	s.tone(440, 0.05, 0.30, square, pluck)
	s.tone(294, 0.05, 0.30, square, pluck)
	s.tone(196, 0.07, 0.30, square, pluck)
	return s.buf
}

// musicPCM builds a calm arpeggio over Am - F - C - G that loops seamlessly.
func musicPCM() []byte {
	s := &synth{}
	chords := [][]float64{
		{220.00, 261.63, 329.63, 440.00}, // Am
		{174.61, 220.00, 261.63, 349.23}, // F
		{261.63, 329.63, 392.00, 523.25}, // C
		{196.00, 246.94, 293.66, 392.00}, // G
	}
	for _, ch := range chords {
		for _, f := range ch {
			s.tone(f, 0.42, 0.16, triangle, pad)
		}
	}
	return s.buf
}
