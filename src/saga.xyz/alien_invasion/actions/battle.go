package actions

import (
	"fmt"
	"saga.xyz/alien_invasion/config"
	"saga.xyz/alien_invasion/maths"
)

func (w *T_World) CheckForCasualties() (completed bool) {
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
					config.Log += fmt.Sprintf("%s has been destroyed by: ", cityname)
					for _, warrior := range fighters {
						fmt.Printf(" Alien %d ", warrior)
						config.Log += fmt.Sprintf(" Alien %d ", warrior)
					}
					fmt.Printf("\n")
					config.Log += fmt.Sprintf("\n")
					w.Destroy(cityname)
					return false
				}
			}
		}
	}
	return true
}

// Assumption: A move off road is not allowed and is a lost/wasted turn
func (w *T_World) Advance() {
	var newpos int

	for id, _ := range w.Invaders {
		direction := maths.RandBetween(1, 4) // 1 = North, 2 = South, 3 = East, 4 = West
		if direction == 1 {
			newpos = w.Invaders[id].Y + 1
			if !w.OnRoad(w.Invaders[id].X, newpos) {
				return
			}
			w.Invaders[id].Y = newpos
		}

		if direction == 2 {
			newpos = w.Invaders[id].Y - 1
			if !w.OnRoad(w.Invaders[id].X, newpos) {
				return
			}
			w.Invaders[id].Y = newpos
		}

		if direction == 3 {
			newpos = w.Invaders[id].X + 1
			if !w.OnRoad(newpos, w.Invaders[id].Y) {
				return
			}
			w.Invaders[id].X = newpos
		}

		if direction == 4 {
			newpos = w.Invaders[id].X - 1
			if !w.OnRoad(newpos, w.Invaders[id].Y) {
				return
			}
			w.Invaders[id].X = newpos
		}

	}
}

func (w *T_World) Destroy(city string) {
	var nw T_World

	nw.Atlas = make(map[int]map[int]string, w.GridSize)
	for i := 0; i < w.GridSize; i++ {
		nw.Atlas[i] = make(map[int]string, w.GridSize)
	}

	for x, _ := range w.Atlas {
		for y, _ := range w.Atlas[x] {
			if city != w.Atlas[x][y] {
				nw.Atlas[x][y] = w.Atlas[x][y]
			} else {
				w.BreakRoad(x, y)
			}
		}
	}

	w.Atlas = nw.Atlas
	delete(w.Cities, city)

	for _, val := range w.CitiesRelativeTo {
		tmp_city := val.City
		tmp_north := val.North
		tmp_south := val.South
		tmp_east := val.East
		tmp_west := val.West

		if tmp_city == city {
			continue
		}

		if tmp_north == city {
			tmp_north = ""
		}

		if tmp_south == city {
			tmp_south = ""
		}

		if tmp_east == city {
			tmp_east = ""
		}

		if tmp_west == city {
			tmp_west = ""
		}

		nw.CitiesRelativeTo = append(nw.CitiesRelativeTo, T_CityRelativeLocation{City: tmp_city,
			North: tmp_north,
			South: tmp_south,
			East:  tmp_east,
			West:  tmp_west})

	}

	w.CitiesRelativeTo = nw.CitiesRelativeTo

	return
}
