package player

import "time"


type Song struct {
	Title    string
	Artist   string
	Album    string
	Year     int
	Track    int
	Genre    string
	Duration time.Duration
	FilePath string
}
