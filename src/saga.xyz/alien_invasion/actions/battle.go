package actions

import (
	"fmt"
	"saga.xyz/alien_invasion/config"
)

func CheckForDestruction(w config.T_World) {
	for x, _ := range w.Atlas {
		for y, _ := range w.Atlas[x] {
			cityname := w.Atlas[x][y]
			fighters := []int{}
			for id, alien := range w.Invaders {
				if alien.X == x && alien.Y == y {
					fighters = append(fighters, id)
				}
				if len(fighters) > 1 {
					// Destroy City
					fmt.Printf("%s has been destroyed by: ", cityname)
					for _, warrior := range fighters {
						fmt.Printf(" Alien %s ", warrior)
					}
					fmt.Printf("\n")
				}
			}
		}
	}
}

func Move(w config.T_World) {

}
