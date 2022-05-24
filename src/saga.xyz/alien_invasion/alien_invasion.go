package main

import (
	"fmt"
	"version"
)

func main() {
	fmt.Printf("Welcome to Alien Invasion! (%s - %s)\n", version.VERSION, version.BUILDID)
}
