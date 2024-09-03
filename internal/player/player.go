package player

import (
	"fmt"

	"github.com/mnsdojo/goplay/internal/metadata"
)

type Player struct {
	playlist     []Song
	currentIndex int
	isPlaying    bool
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
			Album:    meta.Album,
			Year:     meta.Year,
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
	p.isPlaying = true
}

func (p *Player) Pause() {
	p.isPlaying = false
}
func (p *Player) Next() {
	if len(p.playlist) == 0 {
		return
	}
	p.currentIndex = (p.currentIndex + 1) % len(p.playlist)
}

func (p *Player) Prev() {
	if len(p.playlist) == 0 {
		return
	}
	p.currentIndex = (p.currentIndex - 1 + len(p.playlist)) % len(p.playlist)
}

func (p *Player) IsPlaying() bool {
	return p.isPlaying
}
