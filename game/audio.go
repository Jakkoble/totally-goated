package game

import (
	"io"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const sampleRate = 44100

var audioCtx *audio.Context

type SoundPlayer struct {
	data []byte
}

var (
	sfxBounce        *SoundPlayer
	sfxChargeRelease *SoundPlayer
	sfxCrumble       *SoundPlayer
	sfxDashPower5    *SoundPlayer
	sfxDashPower10   *SoundPlayer
	sfxDashPower15   *SoundPlayer
	sfxDeath         *SoundPlayer
	sfxIce           *SoundPlayer
	sfxMilestone     *SoundPlayer
	sfxWallHit       *SoundPlayer
)

func init() {
	audioCtx = audio.NewContext(sampleRate)

	sfxBounce = loadWav("assets/bounce.wav")
	sfxChargeRelease = loadWav("assets/charge_release.wav")
	sfxCrumble = loadWav("assets/crumble.wav")
	sfxDashPower5 = loadWav("assets/dash_power5.wav")
	sfxDashPower10 = loadWav("assets/dash_power10.wav")
	sfxDashPower15 = loadWav("assets/dash_power15.wav")
	sfxDeath = loadWav("assets/death.wav")
	sfxIce = loadWav("assets/ice.wav")
	sfxMilestone = loadWav("assets/milestone.wav")
	sfxWallHit = loadWav("assets/wall_hit.wav")
}

func loadWav(path string) *SoundPlayer {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("audio: open %s: %v", path, err)
	}
	defer f.Close()

	decoded, err := wav.DecodeF32(f)
	if err != nil {
		log.Fatalf("audio: decode %s: %v", path, err)
	}

	data, err := io.ReadAll(decoded)
	if err != nil {
		log.Fatalf("audio: read %s: %v", path, err)
	}

	return &SoundPlayer{data: data}
}

func (s *SoundPlayer) Play() {
	if s == nil {
		return
	}
	p := audioCtx.NewPlayerF32FromBytes(s.data)
	p.Play()
}
