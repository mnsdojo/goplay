package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mnsdojo/goplay/internal/player"
	"github.com/rivo/tview"
)

type MusicPlayerUI struct {
	app        *tview.Application
	player     *player.Player
	playlist   *tview.List
	nowPlaying *tview.TextView
	controls   *tview.TextView
	statusBar  *tview.TextView
	mainView   *tview.TextView
}

func New(player *player.Player) *MusicPlayerUI {
	return &MusicPlayerUI{
		app:        tview.NewApplication(),
		player:     player,
		playlist:   tview.NewList().ShowSecondaryText(false),
		nowPlaying: tview.NewTextView().SetTextAlign(tview.AlignCenter),
		controls:   tview.NewTextView().SetTextAlign(tview.AlignCenter),
		statusBar:  tview.NewTextView().SetTextAlign(tview.AlignLeft),
		mainView:   tview.NewTextView().SetTextAlign(tview.AlignCenter),
	}
}

func (ui *MusicPlayerUI) setupUI() {
	ui.playlist.SetTitle("Playlist").SetBorder(true).SetTitleColor(tcell.ColorGreen)
	ui.nowPlaying.SetTitle("Now Playing").SetBorder(true).SetTitleColor(tcell.ColorYellow)
	ui.controls.SetBorder(true).SetTitleColor(tcell.ColorDarkCyan)
	ui.statusBar.SetBorder(true)
	ui.mainView.SetTitle("Music Info").SetBorder(true).SetTitleColor(tcell.ColorFuchsia)

	for i, song := range ui.player.GetPlaylist() {
		ui.playlist.AddItem(fmt.Sprintf("%3d. %s", i+1, truncateString(song.Title, 25)), "", 0, nil)
	}

	ui.nowPlaying.SetText("No song playing")
	ui.controls.SetText("[ p ] Play/Pause  [ n ] Next  [ b ] Previous  [ q ] Quit")
	ui.statusBar.SetText("Ready")
	ui.mainView.SetText("Welcome to GoPlay Music Player!\n\nSelect a song from the playlist and press 'p' to start playing.")

	grid := tview.NewGrid().
		SetRows(3, 0, 3, 1).
		SetColumns(30, 0).
		SetBorders(false)

	grid.AddItem(ui.nowPlaying, 0, 0, 1, 2, 0, 0, false)
	grid.AddItem(ui.playlist, 1, 0, 1, 1, 0, 0, true)
	grid.AddItem(ui.mainView, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(ui.controls, 2, 0, 1, 2, 0, 0, false)
	grid.AddItem(ui.statusBar, 3, 0, 1, 2, 0, 0, false)

	ui.app.SetRoot(grid, true).EnableMouse(true)
}

func (ui *MusicPlayerUI) Run() error {
	ui.setupUI()
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			ui.player.Stop()
			ui.app.Stop()
			
		case 'p':
			if ui.player.IsPlaying() {
				ui.player.Pause()
				ui.statusBar.SetText("Paused")
			} else {
				ui.player.Play()
				ui.statusBar.SetText("Playing")
			}
			ui.updateNowPlaying()
		case 'n':
			ui.player.Next()
			ui.updateNowPlaying()
		case 'b':
			ui.player.Prev()
			ui.updateNowPlaying()
		}
	
		return event
	})
	return ui.app.Run()
}

func (ui *MusicPlayerUI) updateNowPlaying() {
	currentSong := ui.player.CurrentSong()
	ui.nowPlaying.SetText(fmt.Sprintf("%s - %s", currentSong.Title, currentSong.Artist))
	ui.updateMainView(currentSong)
}

func (ui *MusicPlayerUI) updateMainView(song player.Song) {
	info := fmt.Sprintf("Title: %s\nArtist: %s\nAlbum: %s\nYear: %d\nDuration: %s",
		song.Title, song.Artist, song.Album, song.Year, formatDuration(int(song.Duration)))
	ui.mainView.SetText(info)
}

func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

func formatDuration(seconds int) string {
	minutes := seconds / 60
	remainingSeconds := seconds % 60
	return fmt.Sprintf("%d:%02d", minutes, remainingSeconds)
}
