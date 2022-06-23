package world

import (
	"fmt"
	"github.com/hartsp2000/alien_invasion/config"
)

func (w *T_World) DisplayWorld(showAliens bool) {
	// Clear Screen
	fmt.Printf(string(config.Cls))
	fmt.Println(config.TITLE)
	// Top Border
	fmt.Printf(string(config.Cblue) + "╔")
	for i := 0; i < w.GridSizeX; i++ {
		fmt.Printf("═════")
		if i < w.GridSizeX-1 {
			fmt.Printf("═╦═")
		}
	}
	fmt.Printf("╗\n")

	// Cities
	for y := w.GridSizeY - 1; y >= 0; y-- {
		fmt.Printf(string(config.Cblue) + "║")
		for x := 0; x < w.GridSizeX; x++ {
			block := ""
			// Display Cities and Roads
			_, data, found := w.LookupAtlas(x, y)
			if found {
				if data == "ROAD" {
					block = fmt.Sprintf(string(config.Ccyan)+"%5s"+string(config.Cblue), "  +  ")
				} else {
					if len(data) < 5 {
						block = fmt.Sprintf(string(config.Cyellow)+"%5s"+string(config.Cblue), data[:len(data)])
					} else {
						block = fmt.Sprintf(string(config.Cyellow)+"%5s"+string(config.Cblue), data[:5])
					}
				}
			}

			// Display Aliens
			if showAliens {
				aliens, found := w.LookupAliens(x, y)
				if found {
					block = fmt.Sprintf(string(config.Cred)+" (%d) "+string(config.Cblue), len(aliens))
				}
			}

			fmt.Printf("%5s", block)

			if x < w.GridSizeX-1 {
				fmt.Printf(" ║ ")
			}
		}

		fmt.Printf("║\n")

		// Separator Line
		if y > 0 {
			fmt.Printf(string(config.Cblue) + "╟")
			for i := 0; i < w.GridSizeX; i++ {
				fmt.Printf("─────")
				if i < w.GridSizeX-1 {
					fmt.Printf("─╫─")
				}
			}
			fmt.Printf("╢\n")
		}
	}

	// Bottom Border
	fmt.Printf("╚")
	for i := 0; i < w.GridSizeX; i++ {
		fmt.Printf("═════")
		if i < w.GridSizeX-1 {
			fmt.Printf("═╩═")
		}
	}
	fmt.Printf("╝" + string(config.Creset) + "\n")
	fmt.Printf(config.Log)

}
