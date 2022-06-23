package actions

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/hartsp2000/alien_invasion/config"
	"github.com/hartsp2000/alien_invasion/maths"
	"os"
	"strings"
)

// Optimize per comment: code design: data structure heavy
type T_World struct {
	// Geography and Alien Locations
	Atlas                map[int]map[int]map[int]string // [Key][X][Y]["City Name" or "ROAD"]
	Aliens               map[int]map[int][]T_Alien      // [X][Y][Alien Info]
	GridSizeX, GridSizeY int                            // Grid Size (Added Y For Optics)
}

type T_Alien struct {
	AlienID      int
	AlienSpecies string
}

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
			world.checkAtlasInit(index, x, y) // Checks and initalizes map if necessary
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
				world.checkAtlasInit(index, x, y+c.Distance)
				world.Atlas[index][x][y+c.Distance] = north
				index++
				for tmpy := y; tmpy < y+c.Distance; tmpy++ {
					world.checkAtlasInit(index, x, tmpy)
					world.Atlas[index][x][tmpy] = "ROAD" // Add road to Atlas
					index++
				}
			case "SOUTH":
				south = direction[1] // City name to the South
				world.checkAtlasInit(index, x, y-c.Distance)
				world.Atlas[index][x][y-c.Distance] = south
				index++
				for tmpy := y; tmpy > y-c.Distance; tmpy-- {
					world.checkAtlasInit(index, x, tmpy)
					world.Atlas[index][x][tmpy] = "ROAD"
					index++
				}
			case "EAST":
				east = direction[1] // City name to the East
				world.checkAtlasInit(index, x+c.Distance, y)
				world.Atlas[index][x+c.Distance][y] = east
				index++
				for tmpx := x; tmpx < x+c.Distance; tmpx++ {
					world.checkAtlasInit(index, tmpx, y)
					world.Atlas[index][tmpx][y] = "ROAD"
					index++
				}
			case "WEST":
				west = direction[1] // City name to the West
				world.checkAtlasInit(index, x-c.Distance, y)
				world.Atlas[index][x-c.Distance][y] = west
				index++
				for tmpx := x; tmpx > x-c.Distance; tmpx-- {
					world.checkAtlasInit(index, tmpx, y)
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
				newworld.checkAtlasInit(key, X+Xoffset, Y+Yoffset)
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
		world.checkAliensInit(aX, aY)
		randomChar1 := 'A' + rune(maths.RandBetween(0, 26))
		randomChar2 := 'A' + rune(maths.RandBetween(0, 26))
		species := fmt.Sprintf("%s%s-%d", string(randomChar1), string(randomChar2), maths.RandBetween(1, 32767))
		world.Aliens[aX][aY] = append(world.Aliens[aX][aY], T_Alien{a, species})
		index++
	}

	return nil
}

func (world *T_World) OnRoad(x, y int) bool {
	_, _, found := world.LookupAtlas(x, y)
	return found
}

func (world *T_World) BreakRoad(x, y int) {
	index, _, found := world.LookupAtlas(x, y)
	if found {
		world.Atlas[index][x][y] = ""
	}
}

func (world *T_World) PostApocalypse(c config.T_Config) error {
	for key := 0; key < len(world.Atlas); key++ {

	}

	/*
		var output string

		for x := 0; x < world.GridSize; x++ {
			for y := 0; y < world.GridSize; y++ {
				if _, ok := world.Atlas[x][y]; !ok {
					continue
				}
				city := world.Atlas[x][y]
				output += fmt.Sprintf("%s", city)
				// North
				for posY := y - 1; posY > 0; posY-- {
					if len(world.Atlas[x][posY]) > 0 {
						output += fmt.Sprintf(" north=%s", world.Atlas[x][posY])
						break
					}
				}
				// West
				for posX := x + 1; posX < world.GridSize; posX++ {
					if len(world.Atlas[posX][y]) > 0 {
						output += fmt.Sprintf(" west=%s", world.Atlas[posX][y])
						break
					}
				}
				// South
				for posY := y + 1; posY < world.GridSize; posY++ {
					if len(world.Atlas[x][posY]) > 0 {
						output += fmt.Sprintf(" south=%s", world.Atlas[x][posY])
						break
					}
				}
				// East
				for posX := x - 1; posX > 0; posX-- {
					if len(world.Atlas[posX][y]) > 0 {
						output += fmt.Sprintf(" east=%s", world.Atlas[posX][y])
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
	*/
	return nil
}

func (world *T_World) ListAtlas() (xlist []int, ylist []int) {
	for key := 0; key < len(world.Atlas); key++ {
		for X, _ := range world.Atlas[key] {
			for Y, _ := range world.Atlas[key][X] {
				if len(world.Atlas[key][X][Y]) > 0 {
					xlist = append(xlist, X)
					ylist = append(ylist, Y)
				}
			}
		}
	}
	return xlist, ylist
}

func (world *T_World) ListAliens() (xlist []int, ylist []int) {
	for X, _ := range world.Aliens {
		for Y, _ := range world.Aliens[X] {
			if len(world.Aliens[X][Y]) > 0 {
				xlist = append(xlist, X)
				ylist = append(ylist, Y)
			}
		}
	}
	return xlist, ylist
}

func (world *T_World) LookupAtlas(x, y int) (index int, data string, found bool) {
	found = false
	for key := 0; key < len(world.Atlas); key++ {
		if _, ok := world.Atlas[key][x][y]; ok {
			if len(world.Atlas[key][x][y]) > 0 {
				found = true
			}
			return key, world.Atlas[key][x][y], found
		}
	}
	return -1, "", found
}

func (world *T_World) LookupCity(name string) (x int, y int, err error) {
	for key, _ := range world.Atlas {
		for x, _ := range world.Atlas[key] {
			for y, _ := range world.Atlas[key][x] {
				if _, ok := world.Atlas[key][x][y]; ok {
					if world.Atlas[key][x][y] == name {
						return x, y, nil
					}
				}
			}
		}
	}
	return -1, -1, errors.New("Can't find city in world map")
}

func (world *T_World) LookupAtlasIndex(index int) (x, y int, name string) {
	for X, _ := range world.Atlas[index] {
		for Y, _ := range world.Atlas[index][x] {
			if X == x && Y == y {
				return X, Y, world.Atlas[index][x][y]
			}
		}
	}
	return -1, -1, ""
}

func (world *T_World) LookupAliens(x, y int) (invaders []T_Alien, found bool) {
	found = false
	if _, ok := world.Aliens[x][y]; ok {
		if len(world.Aliens[x][y]) > 0 {
			found = true
		}
		return world.Aliens[x][y], found
	}
	return []T_Alien{}, found
}

func (world *T_World) MoveAlien(aid, fromX, fromY, toX, toY int) {
	var inhabitants []T_Alien = world.Aliens[fromX][fromY]
	var newinhabitants []T_Alien
	world.checkAliensInit(toX, toY)
	for _, alien := range inhabitants {
		if alien.AlienID == aid {
			world.Aliens[toX][toY] = append(world.Aliens[toX][toY], alien)
			// If there is more than one alien here, we want to preserve their location
			for stillhere, _ := range inhabitants {
				if inhabitants[stillhere].AlienID != aid {
					newinhabitants = append(newinhabitants, inhabitants[stillhere])
				}
			}
			delete(world.Aliens[fromX], fromY)
			world.checkAliensInit(fromX, fromY)
			world.Aliens[fromX][fromY] = newinhabitants
		}
	}
}

func (world *T_World) checkAtlasInit(index, x, y int) {
	if _, ok := world.Atlas[index]; !ok {
		world.Atlas[index] = make(map[int]map[int]string)
	}

	if _, ok := world.Atlas[index][x]; !ok {
		world.Atlas[index][x] = make(map[int]string)
	}
}

func (world *T_World) checkAliensInit(x, y int) {
	if _, ok := world.Aliens[x]; !ok {
		world.Aliens[x] = make(map[int][]T_Alien)
	}
}
