package player

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/mnsdojo/goplay/internal/metadata"
)

type Player struct {
	playlist     []Song
	currentIndex int
	isPlaying    bool
	isPaused     bool
	cmd          *exec.Cmd
}

func New(musicDir string) (*Player, error) {
	extractor := metadata.NewExtractor()
	metaList, err := extractor.ExtractAll(musicDir)
	if err != nil {
		return nil, fmt.Errorf("error extracting metadata: %v", err)
	}
	fmt.Printf("Debug: Number of metadata entries extracted: %d\n", len(metaList))

	playlist := make([]Song, 0, len(metaList))
	for _, meta := range metaList {
		song := Song{
			Title:    meta.Title,
			Artist:   meta.Artist,
			Year:     meta.Year,
			Album:    meta.Album,
			Track:    meta.Track,
			Genre:    meta.Genre,
			Duration: meta.Duration,
			FilePath: meta.FilePath,
		}
		playlist = append(playlist, song)
	}
	return &Player{
		playlist:     playlist,
		currentIndex: 0,
		isPlaying:    false,
		isPaused:     false,
	}, nil
}

func (p *Player) GetPlaylist() []Song {
	return p.playlist
}

func (p *Player) CurrentSong() Song {
	if len(p.playlist) == 0 {
		return Song{}
	}
	return p.playlist[p.currentIndex]
}

func (p *Player) Play() error {
	if len(p.playlist) == 0 {
		return fmt.Errorf("playlist is empty")
	}
	if p.isPaused {
		return p.Resume()
	}
	p.Stop()
	p.cmd = exec.Command("mpv", "--no-video", "--no-terminal", p.playlist[p.currentIndex].FilePath)
	err := p.cmd.Start()
	if err != nil {
		return fmt.Errorf("error starting mpv: %v", err)
	}
	p.isPlaying = true
	p.isPaused = false
	return nil
}

func (p *Player) Stop() {
	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Signal(os.Interrupt)
		timeout := time.AfterFunc(5*time.Second, func() {
			p.cmd.Process.Kill()
		})
		p.cmd.Wait()
		timeout.Stop()
		p.cmd = nil
	}
	p.isPlaying = false
	p.isPaused = false
}

func (p *Player) Pause() error {
	if !p.isPlaying || p.isPaused {
		return nil
	}
	if p.cmd != nil && p.cmd.Process != nil {
		err := exec.Command("pkill", "-STOP", "mpv").Run()
		if err != nil {
			return fmt.Errorf("error pausing playback: %v", err)
		}
		p.isPaused = true
		p.isPlaying = false
	}
	return nil
}

func (p *Player) Resume() error {
	if !p.isPaused {
		return nil
	}
	if p.cmd != nil && p.cmd.Process != nil {
		err := exec.Command("pkill", "-CONT", "mpv").Run()
		if err != nil {
			return fmt.Errorf("error resuming playback: %v", err)
		}
		p.isPaused = false
		p.isPlaying = true
	}
	return nil
}

func (p *Player) Next() error {
	if len(p.playlist) == 0 {
		return fmt.Errorf("playlist is empty")
	}
	p.currentIndex = (p.currentIndex + 1) % len(p.playlist)
	return p.Play()
}

func (p *Player) SetCurrentIndex(index int) error {
	if index < 0 || index >= len(p.playlist) {
		return fmt.Errorf("invalid index: %d", index)
	}
	p.currentIndex = index
	return nil
}

func (p *Player) Prev() error {
	if len(p.playlist) == 0 {
		return fmt.Errorf("playlist is empty")
	}
	p.currentIndex = (p.currentIndex - 1 + len(p.playlist)) % len(p.playlist)
	return p.Play()
}

func (p *Player) IsPlaying() bool {
	return p.isPlaying
}

func (p *Player) IsPaused() bool {
	return p.isPaused
}

func (p *Player) Cleanup() {
	p.Stop()
	exec.Command("pkill", "mpv").Run()
}
