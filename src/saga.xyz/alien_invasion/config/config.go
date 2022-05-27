package config

// ANSI Basic Color Map
const (
	Cls    = "\033[H\033[2J"
	Creset = "\033[0m"

	Cred    = "\033[31m"
	Cgreen  = "\033[32m"
	Cyellow = "\033[33m"
	Cblue   = "\033[34m"
	Cpurple = "\033[35m"
	Ccyan   = "\033[36m"
	Cwhite  = "\033[37m"
)

var TITLE string = string(Cgreen) + "Alien Invasion - A Wargames Simulated Invasion" + string(Creset)
var Log string

// Configurable Options
type T_Config struct {
	MapInfile    string
	MapOutfile   string
	Aliens       int
	AllowedMoves int
	Distance     int
	ShowMap      bool
	Refresh      string
}
