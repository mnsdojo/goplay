package main

import (
	"fmt"
	"log"

	"github.com/mnsdojo/goplay/internal/player"
)

func main() {
	musicDir := "../Music"

	musicPlayer, err := player.New(musicDir)
	if err != nil {
		log.Fatalf("Error initializing player :%v\n", err)
	}
	playlist := musicPlayer.GetPlaylist()

	fmt.Printf("total songs : %d\n", len(playlist))
	for _, song := range playlist {
		fmt.Printf("Title :%s\n ", song.Title)
		fmt.Printf("Artist :%s\n ", song.Artist)

	}
}
