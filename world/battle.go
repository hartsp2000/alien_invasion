package world

import (
	"fmt"
	"github.com/hartsp2000/alien_invasion/config"
	"github.com/hartsp2000/alien_invasion/maths"
)

func (w *T_World) Fight() {
	for x, _ := range w.Aliens {
		for y, _ := range w.Aliens[x] {
			fighters, _ := w.LookupAliens(x, y)
			if len(fighters) > 1 {
				_, cityname, found := w.LookupAtlas(x, y)
				if found {
					if cityname != "ROAD" {
						msg := fmt.Sprintf("%s has been destroyed by: ", cityname)
						for war, _ := range w.Aliens[x][y] {
							msg += fmt.Sprintf("Alien %d (Species: %s)    ", fighters[war].AlienID, fighters[war].AlienSpecies)
						}
						msg += fmt.Sprintf("\n")
						fmt.Printf(msg)
						config.Log += msg
						w.BreakRoad(x, y)
					}
				}
			}
		}
	}

}

func (w *T_World) Iterate() {
	var newpos int
	for x, _ := range w.Aliens {
		for y, _ := range w.Aliens[x] {
			fighters, _ := w.LookupAliens(x, y)
			for traveler, _ := range fighters {
				switch direction := maths.RandBetween(1, 4); direction {
				case 1: // Move North
					newpos = y + 1
					if w.OnRoad(x, newpos) {
						w.MoveAlien(fighters[traveler].AlienID, x, y, x, newpos)
					}

				case 2: // Move South
					newpos = y - 1
					if w.OnRoad(x, newpos) {
						w.MoveAlien(fighters[traveler].AlienID, x, y, x, newpos)
					}

				case 3: // Move East
					newpos = x + 1
					if w.OnRoad(newpos, y) {
						w.MoveAlien(fighters[traveler].AlienID, x, y, newpos, y)
					}

				case 4: // Move West
					newpos = x - 1
					if w.OnRoad(newpos, y) {
						w.MoveAlien(fighters[traveler].AlienID, x, y, newpos, y)
					}

				default:
				}
			}
		}
	}
}
