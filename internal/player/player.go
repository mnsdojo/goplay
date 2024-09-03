package player

import (
	"fmt"
	"os/exec"

	"github.com/mnsdojo/goplay/internal/metadata"
)

type Player struct {
	playlist     []Song
	currentIndex int
	isPlaying    bool
	cmd          *exec.Cmd
}

func New(musicDir string) (*Player, error) {
	extractor := metadata.NewExtractor()
	metaList, err := extractor.ExtractAll(musicDir)
	if err != nil {
		return nil, fmt.Errorf("Error extracting meta :%v", err)
	}
	fmt.Printf("Debug: Number of metadata entries extracted: %d\n", len(metaList))

	var playlist []Song
	for _, meta := range metaList {
		song := Song{
			Title:    meta.Title,
			Artist:   meta.Artist,
			Year:     meta.Year,
			Album:     meta.Album,
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

func (p *Player) Play() {
	if len(p.playlist) == 0 {
		return
	}
	// kill the process
	if p.cmd != nil {
		p.cmd.Process.Kill()
	}
	p.cmd = exec.Command("mpv", "--no-video", "--no-terminal", p.playlist[p.currentIndex].FilePath)
	err := p.cmd.Start()
	if err != nil {
		fmt.Printf("Error starting mpv %v\n", err)
		return
	}
	p.isPlaying = true

}

func (p *Player) Stop() {
	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
		p.cmd.Wait()
		p.cmd = nil
	}
	p.isPlaying = false
}

func (p *Player) Pause() {
	if p.cmd != nil && p.cmd.Process != nil {
		exec.Command("pkill", "-STOP", "mpv").Run()
		p.isPlaying = false
	}
}
func (p *Player) Next() {
	if len(p.playlist) == 0 {
		return
	}
	p.currentIndex = (p.currentIndex + 1) % len(p.playlist)
	p.Play()
}
func (p *Player)SetCurrentIndex(index int){
	p.currentIndex= index
}
func (p *Player) Prev() {
	if len(p.playlist) == 0 {
		return
	}
	p.currentIndex = (p.currentIndex - 1 + len(p.playlist)) % len(p.playlist)
	p.Play()
}

func (p *Player) IsPlaying() bool {
	return p.isPlaying
}
