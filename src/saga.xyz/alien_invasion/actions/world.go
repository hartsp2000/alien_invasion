package actions

import (
	"bufio"
	"fmt"
	"os"
	"saga.xyz/alien_invasion/config"
	"saga.xyz/alien_invasion/maths"
	"sort"
	"strings"
)

type T_World struct {
	Cities           map[string]T_CityData
	CitiesRelativeTo []T_CityRelativeLocation
	Rows             int                        // Y axis
	Cols             int                        // X axis
	Atlas            map[int]map[int]string     // World Map
	Traversible      []T_TraversibleCoordinates // List of coordinates that can be traversed
	Invaders         []T_Alien
	GridSize         int // GridSize x GridSize
}

type T_CityRelativeLocation struct {
	City  string
	North string
	South string
	East  string
	West  string
}

type T_CityData struct {
	Coordinates T_CityCoordinates
}

type T_CityCoordinates struct {
	X int
	Y int
}

type T_TraversibleCoordinates struct {
	X int
	Y int
}

type T_Alien struct {
	X int
	Y int
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
	world.Cities = make(map[string]T_CityData)

	var x int = 0
	var y int = 0

	var firstcity bool = true

	for _, linedata := range data {
		columndata := strings.Split(linedata, " ")

		var cityname, north, south, east, west string

		cityname = columndata[0] // Column 0 is assumed to *ALWAYS* be the City Name
		// Directions are case-insensitive but must be spelled out fully and correctly
		// City names *ARE* case sensitive
		if firstcity {
			world.Cities[cityname] = T_CityData{Coordinates: T_CityCoordinates{x, y}}
			firstcity = false
		}

		x = world.Cities[cityname].Coordinates.X
		y = world.Cities[cityname].Coordinates.Y

		for _, directionline := range columndata[1:] {
			//distance := maths.RandBetween(1, c.Distance) // Choose a random distance between cities
			direction := strings.Split(directionline, "=")
			switch direct := strings.ToUpper(direction[0]); direct {
			case "NORTH":
				north = direction[1]
				world.Cities[north] = T_CityData{Coordinates: T_CityCoordinates{x, y - c.Distance}}
				for tmpy := y; tmpy > y-c.Distance; tmpy-- {
					world.Traversible = append(world.Traversible, T_TraversibleCoordinates{x, tmpy})
				}
			case "SOUTH":
				south = direction[1]
				world.Cities[south] = T_CityData{Coordinates: T_CityCoordinates{x, y + c.Distance}}
				for tmpy := y; tmpy < y+c.Distance; tmpy++ {
					world.Traversible = append(world.Traversible, T_TraversibleCoordinates{x, tmpy})
				}
			case "EAST":
				east = direction[1]
				world.Cities[east] = T_CityData{Coordinates: T_CityCoordinates{x + c.Distance, y}}
				for tmpx := x; tmpx < x+c.Distance; tmpx++ {
					world.Traversible = append(world.Traversible, T_TraversibleCoordinates{tmpx, y})
				}
			case "WEST":
				west = direction[1]
				world.Cities[west] = T_CityData{Coordinates: T_CityCoordinates{x - c.Distance, y}}
				for tmpx := x; tmpx > x-c.Distance; tmpx-- {
					world.Traversible = append(world.Traversible, T_TraversibleCoordinates{tmpx, y})
				}
			default:
			}
		}

		world.CitiesRelativeTo = append(world.CitiesRelativeTo, T_CityRelativeLocation{City: cityname, North: north, South: south, East: east, West: west})
	}

	// Calculate Grid Size
	var xList []int
	var yList []int

	for _, val := range world.Cities {
		xList = append(xList, val.Coordinates.X)
		yList = append(yList, val.Coordinates.Y)
	}
	sort.Sort(sort.IntSlice(xList))
	sort.Sort(sort.IntSlice(yList))

	cols := maths.Abs(xList[0]) + maths.Abs(xList[len(xList)-1]) + 1
	rows := maths.Abs(yList[0]) + maths.Abs(yList[len(yList)-1]) + 1

	world.Cols = cols
	world.Rows = rows

	// Calculate an offset to make the grid positive
	// Assumption: I want a square grid so I will calculate
	// an arbitrary equal xy maxsize (It gives aliens more
	// space to roam)

	var offset int
	var maxsize int

	offset = rows
	maxsize = offset + rows

	if cols > rows {
		offset = cols
		maxsize = offset + cols
	}

	world.GridSize = maxsize

	// Transform the traversible terrian with the offset
	for n, _ := range world.Traversible {
		world.Traversible[n].X = world.Traversible[n].X + offset
		world.Traversible[n].Y = world.Traversible[n].Y + offset
	}

	// Reticulating splines (not really -- but it sounds cool)

	// Initialize the map
	world.Atlas = make(map[int]map[int]string, maxsize)
	for x := 0; x < maxsize; x++ {
		world.Atlas[x] = make(map[int]string, maxsize)
	}

	// Transform with offset and Populate map with cities
	for city, data := range world.Cities {
		x := data.Coordinates.X + offset
		y := data.Coordinates.Y + offset
		world.Atlas[x][y] = city
		world.Cities[city] = T_CityData{Coordinates: T_CityCoordinates{x, y}}
	}

	// Create Alien Invaders
	for a := 0; a < c.Aliens; a++ {
		aX := maths.RandBetween(0, maxsize)
		aY := maths.RandBetween(0, maxsize)

		// Keep looping until we get some coordinates on traversible terrain
		for !world.OnRoad(aX, aY) {
			aX = maths.RandBetween(0, maxsize)
			aY = maths.RandBetween(0, maxsize)
		}
		world.Invaders = append(world.Invaders, T_Alien{aX, aY})
	}

	return nil
}

func (world *T_World) OnRoad(x, y int) bool {
	for n, _ := range world.Traversible {
		if world.Traversible[n].X == x && world.Traversible[n].Y == y {
			return true
		}
	}
	return false
}

func (world *T_World) BreakRoad(x, y int) {
	var newworld T_World

	for n := 0; n < len(world.Traversible); n++ {
		if world.Traversible[n].X == x && world.Traversible[n].Y == y {
			continue
		}
		newworld.Traversible = append(newworld.Traversible, T_TraversibleCoordinates{world.Traversible[n].X, world.Traversible[n].Y})
	}
	world.Traversible = newworld.Traversible
}

func (world *T_World) PostApocalypse(c config.T_Config) error {
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

	return nil
}
