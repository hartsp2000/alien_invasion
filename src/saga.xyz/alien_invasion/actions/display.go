package actions

import (
	"fmt"
	"saga.xyz/alien_invasion/config"
)

func DisplayWorld(w config.T_World, showAliens bool) {
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
	for y := w.GridSize; y > 0; y-- {
		fmt.Printf(string(config.Cblue) + "║")
		for x := 0; x < w.GridSize; x++ {
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
			fmt.Printf(string(config.Cyellow)+"%5s"+string(config.Cblue), city)
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

}
