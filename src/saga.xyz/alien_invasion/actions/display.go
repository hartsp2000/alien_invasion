package actions

import (
	"fmt"
	"saga.xyz/alien_invasion/config"
)

func (w *T_World) DisplayWorld(showAliens bool) {
	// Clear Screen
	fmt.Printf(string(config.Cls))
	fmt.Println(config.TITLE)
	// Top Border
	fmt.Printf(string(config.Cblue) + "╔")
	for i := 0; i < w.GridSize; i++ {
		fmt.Printf("═════")
		if i < w.GridSize-1 {
			fmt.Printf("═╦═")
		}
	}
	fmt.Printf("╗\n")

	// Cities
	for y := 0; y < w.GridSize; y++ {
		fmt.Printf(string(config.Cblue) + "║")
		for x := 0; x < w.GridSize; x++ {
			road := ""
			for _, r := range w.Traversible {
				if r.X == x && r.Y == y {
					road = "  +  "
				}
			}
			city := w.Atlas[x][y]
			if len(w.Atlas[x][y]) > 5 {
				city = w.Atlas[x][y][:5]
			}
			if showAliens {
				for _, val := range w.Invaders {
					if val.X == x && val.Y == y {
						city = fmt.Sprintf(string(config.Cred)+"%5s"+string(config.Cblue), "!!")
					}
				}
			}
			if len(city) == 0 && len(road) > 0 {
				fmt.Printf(string(config.Ccyan)+"%5s"+string(config.Cblue), road)
			} else {
				fmt.Printf(string(config.Cyellow)+"%5s"+string(config.Cblue), city)
			}
			if x < w.GridSize-1 {
				fmt.Printf(" ║ ")
			}
		}
		fmt.Printf("║\n")

		// Separator Line
		fmt.Printf(string(config.Cblue) + "╟")
		for i := 0; i < w.GridSize; i++ {
			fmt.Printf("─────")
			if i < w.GridSize-1 {
				fmt.Printf("─╫─")
			}
		}
		fmt.Printf("╢\n")
	}

	// Empty Line
	fmt.Printf("║")
	for i := 0; i < w.GridSize; i++ {
		fmt.Printf("     ")
		if i < w.GridSize-1 {
			fmt.Printf(" ║ ")
		}
	}
	fmt.Printf("║\n")

	// Bottom Border
	fmt.Printf("╚")
	for i := 0; i < w.GridSize; i++ {
		fmt.Printf("═════")
		if i < w.GridSize-1 {
			fmt.Printf("═╩═")
		}
	}
	fmt.Printf("╝" + string(config.Creset) + "\n")
	fmt.Printf(config.Log)
}
