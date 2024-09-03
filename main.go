package main

import (
	"log"

	"github.com/mnsdojo/goplay/internal/player"
	"github.com/mnsdojo/goplay/internal/ui"
)

func main() {
	musicDir := "../Music"

	musicPlayer, err := player.New(musicDir)
	if err != nil {
		log.Fatalf("Error initializing player :%v\n", err)
	}
	playerUi := ui.New(musicPlayer)
	if err := playerUi.Run(); err != nil {
		log.Fatalf("Error running UI: %v\n", err)
	}

}
