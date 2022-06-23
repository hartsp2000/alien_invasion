package world

import (
	"bufio"
	"fmt"
	"github.com/hartsp2000/alien_invasion/config"
	"os"
)

func (world *T_World) Apocalypse(c config.T_Config) error {
	var output string

	for x := 0; x < world.GridSizeX; x++ {
		for y := 0; y < world.GridSizeY; y++ {
			_, city, found := world.LookupAtlas(x, y)

			if !found || len(city) < 1 || city == "ROAD" {
				continue
			}

			output += fmt.Sprintf("%s", city)

			// North
			for posY := y + 1; posY < world.GridSizeY; posY++ {
				_, nextcity, _ := world.LookupAtlas(x, posY)
				if len(nextcity) > 0 && nextcity != "ROAD" {
					output += fmt.Sprintf(" north=%s", nextcity)
					break
				}
			}
			// West
			for posX := x - 1; posX > 0; posX-- {
				_, nextcity, _ := world.LookupAtlas(posX, y)
				if len(nextcity) > 0 && nextcity != "ROAD" {
					output += fmt.Sprintf(" west=%s", nextcity)
					break
				}
			}
			// South
			for posY := y - 1; posY > 0; posY-- {
				_, nextcity, _ := world.LookupAtlas(x, posY)
				if len(nextcity) > 0 && nextcity != "ROAD" {
					output += fmt.Sprintf(" south=%s", nextcity)
					break
				}
			}
			// East
			for posX := x + 1; posX < world.GridSizeX; posX++ {
				_, nextcity, _ := world.LookupAtlas(posX, y)
				if len(nextcity) > 0 && nextcity != "ROAD" {
					output += fmt.Sprintf(" east=%s", nextcity)
					break
				}
			}
			output += fmt.Sprintf("\n")
		}
	}

	file, err := os.Create((c.MapOutfile))

	if err != nil {
		return err
	}

	dw := bufio.NewWriter(file)
	dw.WriteString(output)
	dw.Flush()
	file.Close()

	return nil
}
