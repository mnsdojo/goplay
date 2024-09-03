package player

import "time"

type Song struct {
	Title    string
	Artist   string
	Year     int
	Track    int
	Genre    string
	Album    string
	Duration time.Duration
	FilePath string
}
