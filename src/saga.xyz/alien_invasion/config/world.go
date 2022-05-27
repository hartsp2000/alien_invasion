package config

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

type T_World struct {
	Cities           map[string]T_CityData
	CitiesRelativeTo []T_CityRelativeLocation
	Rows             int                    // Y axis
	Cols             int                    // X axis
	Atlas            map[int]map[int]string // World Map
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
	Destroyed   bool
	Coordinates T_CityCoordinates
}

type T_CityCoordinates struct {
	X int
	Y int
}

type T_Alien struct {
	X int
	Y int
}

func (c *T_Config) Genesis() (world T_World, err error) {
	// Load the data line-by-line
	file, err := os.Open(c.MapInfile)
	if err != nil {
		return world, err
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
			distance := RandBetween(1, c.Distance) // Choose a random distance between cities
			direction := strings.Split(directionline, "=")
			switch direct := strings.ToUpper(direction[0]); direct {
			case "NORTH":
				north = direction[1]
				world.Cities[north] = T_CityData{Coordinates: T_CityCoordinates{x, y + distance}}
			case "SOUTH":
				south = direction[1]
				world.Cities[south] = T_CityData{Coordinates: T_CityCoordinates{x, y - distance}}
			case "EAST":
				east = direction[1]
				world.Cities[east] = T_CityData{Coordinates: T_CityCoordinates{x + distance, y}}
			case "WEST":
				west = direction[1]
				world.Cities[west] = T_CityData{Coordinates: T_CityCoordinates{x - 1, distance}}
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

	cols := Abs(xList[0]) + Abs(xList[len(xList)-1]) + 1
	rows := Abs(yList[0]) + Abs(yList[len(yList)-1]) + 1

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

	// Reticulating splines (not really -- but it sounds cool)

	// Initialize the map
	world.Atlas = make(map[int]map[int]string, maxsize)
	for x := 0; x < maxsize; x++ {
		world.Atlas[x] = make(map[int]string, maxsize)
	}

	// Populate map with cities
	for city, data := range world.Cities {
		x := data.Coordinates.X + offset
		y := data.Coordinates.Y + offset
		world.Atlas[x][y] = city
	}

	// Create Alien Invaders
	for a := 0; a < c.Aliens; a++ {
		aX := RandBetween(0, maxsize)
		aY := RandBetween(0, maxsize)
		world.Invaders = append(world.Invaders, T_Alien{aX, aY})
	}

	return world, nil
}
