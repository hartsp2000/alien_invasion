package world

import (
	"bufio"
	"fmt"
	"github.com/hartsp2000/alien_invasion/config"
	"github.com/hartsp2000/alien_invasion/maths"
	"os"
	"strings"
)

func (world *T_World) Genesis(c config.T_Config) (err error) {
	// Load the data line-by-line
	file, err := os.Open(c.MapInfile)
	if err != nil {
		return err
	}
	defer file.Close()

	var data []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	// Process the data and create the world
	world.Atlas = make(map[int]map[int]map[int]string)
	world.Aliens = make(map[int]map[int][]T_Alien)

	var x int = 0
	var y int = 0
	var index int = 0
	var cityerr error
	var firstcity bool = true // The first city will be at inital coordinates (0, 0)

	for _, linedata := range data {
		columndata := strings.Split(linedata, " ")

		var cityname, north, south, east, west string

		cityname = columndata[0] // Column 0 is assumed to *ALWAYS* be the City Name
		// Directions are case-insensitive but must be spelled out fully and correctly
		// City names *ARE* case sensitive
		if firstcity {
			world.CheckAtlasInit(index, x, y) // Checks and initalizes map if necessary
			world.Atlas[index][x][y] = cityname
			index++
			firstcity = false
		} else {
			x, y, cityerr = world.LookupCity(cityname)
			if err != nil {
				return cityerr
			}
		}

		for _, directionline := range columndata[1:] {
			direction := strings.Split(directionline, "=")
			switch direct := strings.ToUpper(direction[0]); direct {
			case "NORTH":
				north = direction[1] // City name to the North
				world.CheckAtlasInit(index, x, y+c.Distance)
				world.Atlas[index][x][y+c.Distance] = north
				index++
				for tmpy := y; tmpy < y+c.Distance; tmpy++ {
					world.CheckAtlasInit(index, x, tmpy)
					world.Atlas[index][x][tmpy] = "ROAD" // Add road to Atlas
					index++
				}
			case "SOUTH":
				south = direction[1] // City name to the South
				world.CheckAtlasInit(index, x, y-c.Distance)
				world.Atlas[index][x][y-c.Distance] = south
				index++
				for tmpy := y; tmpy > y-c.Distance; tmpy-- {
					world.CheckAtlasInit(index, x, tmpy)
					world.Atlas[index][x][tmpy] = "ROAD"
					index++
				}
			case "EAST":
				east = direction[1] // City name to the East
				world.CheckAtlasInit(index, x+c.Distance, y)
				world.Atlas[index][x+c.Distance][y] = east
				index++
				for tmpx := x; tmpx < x+c.Distance; tmpx++ {
					world.CheckAtlasInit(index, tmpx, y)
					world.Atlas[index][tmpx][y] = "ROAD"
					index++
				}
			case "WEST":
				west = direction[1] // City name to the West
				world.CheckAtlasInit(index, x-c.Distance, y)
				world.Atlas[index][x-c.Distance][y] = west
				index++
				for tmpx := x; tmpx > x-c.Distance; tmpx-- {
					world.CheckAtlasInit(index, tmpx, y)
					world.Atlas[index][tmpx][y] = "ROAD"
					index++
				}
			default:
			}
		}
	}

	// Calculate Grid Size
	var xList []int
	var yList []int

	for key := 0; key < len(world.Atlas); key++ {
		for X, _ := range world.Atlas[key] {
			for Y, _ := range world.Atlas[key][X] {
				xList = append(xList, X)
				yList = append(yList, Y)
			}
		}
	}

	// Optimize per comment: finds map width by sorting and grabbing indices, could have found max and min in loop instead
	xmin, xmax := maths.MinMax(xList)
	ymin, ymax := maths.MinMax(yList)

	cols := maths.Count(xmin, xmax)
	rows := maths.Count(ymin, ymax)

	// Calculate gridsize and an offset that will make the grid positive
	world.GridSizeX = cols + 1
	world.GridSizeY = rows + 1

	var Xoffset int = maths.Offset(xmin)
	var Yoffset int = maths.Offset(ymin)

	// Transform the Atlas with the offset (move the map to only the positive side of the grid)
	var newworld T_World
	newworld.Atlas = make(map[int]map[int]map[int]string)

	for key := 0; key < len(world.Atlas); key++ {
		for X, _ := range world.Atlas[key] {
			for Y, name := range world.Atlas[key][X] {
				newworld.CheckAtlasInit(key, X+Xoffset, Y+Yoffset)
				newworld.Atlas[key][X+Xoffset][Y+Yoffset] = name
			}
		}
	}
	world.Atlas = newworld.Atlas

	// Create Alien Invaders
	index = 0

	for a := 0; a < c.Aliens; a++ {
		aX := maths.RandBetween(0, cols)
		aY := maths.RandBetween(0, rows)

		// Keep looping until they land on a city or road
		for !world.OnRoad(aX, aY) {
			aX = maths.RandBetween(0, cols)
			aY = maths.RandBetween(0, rows)
		}
		world.CheckAliensInit(aX, aY)
		randomChar1 := 'A' + rune(maths.RandBetween(0, 26))
		randomChar2 := 'A' + rune(maths.RandBetween(0, 26))
		species := fmt.Sprintf("%s%s-%d", string(randomChar1), string(randomChar2), maths.RandBetween(1, 32767))
		world.Aliens[aX][aY] = append(world.Aliens[aX][aY], T_Alien{a, species})
		index++
	}

	return nil
}
