package main

import (
	//"bytes"
	"fmt"
	"log"
	"os"
	//"strings"

	"github.com/BurntSushi/toml"
	mcshapes "github.com/GreenSeaTurtle/mcFunctionDev/mcShapes"
	"github.com/olekukonko/tablewriter"
	//"github.com/benmcclelland/mcrender"
)

// Structure for using TOML to extract input from the user.
// CVHeight      Height to clear, should be large enough to clear everything    default=100
// CVWidth       Width, X axis, left to right (west to east) when facing north  default=50
// CVDepth       Depth, Z axis, near to far (south to north) when facing north  default=50
// CVBlockType   Replace all blocks in the clear volume with this block, default="air"
//
// Clears width*depth*height in front of the player
// This capability can be used to create a rectangular volume of whatever is desired.
type mcfdClearVolInputStruct struct {
	ClearVolHeight     []int    `toml:"ClearVolHeight"`
	ClearVolWidth      []int    `toml:"ClearVolWidth"`
	ClearVolDepth      []int    `toml:"ClearVolDepth"`
	ClearVolBlockType  []string `toml:"ClearVolBlockType"`
}

// CreateClearVolDriver
// Driver for creating the Minecraft function files for clearing volumes
func CreateClearVolDriver(inputFile string, basepath string) {
	// Extract pertinent input, using TOML, from the user input file
	var mcfdInput mcfdClearVolInputStruct
	if _, err := toml.DecodeFile(inputFile, &mcfdInput); err != nil {
		fmt.Println(err)
		return
	}

	// Consistency check on the user input
	dim := [4]int{0, 0, 0, 0}
	dim[0] = len(mcfdInput.ClearVolHeight)
	dim[1] = len(mcfdInput.ClearVolWidth)
	dim[2] = len(mcfdInput.ClearVolDepth)
	dim[3] = len(mcfdInput.ClearVolBlockType)
	maxdim := dim[0]
	mindim := dim[0]
	for _, v := range dim {
		if v < mindim {
			mindim = v
		}
		if v > maxdim {
			maxdim = v
		}
	}
	if maxdim != mindim {
		fmt.Println("CreateClearVol user input FATAL ERROR")
		fmt.Println("You must specify the same number of array values for all the")
		fmt.Println("CreateClearVol arrays - ClearVolHeight, ClearVolBlockType, ...")
		return
	}

	// If the user has not specified anything then there is nothing left
	// to do.
	if maxdim <= 0 {
		return
	}

	// Create the clear volumes requested by the user in the user input file.

	// First echo user input to stdout so the user knows what was done.
	// This also sets the filename to write the Minecraft function data.
	fmt.Println("\nCreating ClearVol Functions for Minecraft")
	fmt.Println("The following table summarizes user input for the volumes:")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Filename", "Width", "Depth", "Height", "Block"})
	ndirvals := 4
	filename := make([]string, maxdim*ndirvals)
	directionValues := []string{"north", "east", "south", "west"}
	directionNames  := []string{"N",     "E",    "S",     "W"}
	for i := 0; i < maxdim; i++ {
		for j := 0; j < ndirvals; j++ {
			dname := directionNames[j]
			k := j + i*ndirvals
			// Minecraft functions must have a suffix of ".mcfunction"
			height := mcfdInput.ClearVolHeight[i]
			height_str := fmt.Sprintf("%d", mcfdInput.ClearVolHeight[i])
			sheight := ""
			if height != 100 {
				sheight = "_" + height_str
			}
			width_str := fmt.Sprintf("%d", mcfdInput.ClearVolWidth[i])
			swidth := "_" + width_str
			depth_str  := fmt.Sprintf("%d", mcfdInput.ClearVolDepth[i])
			sdepth := "_" + depth_str
			bname := mcfdInput.ClearVolBlockType[i]
			sbname := ""
			if bname != "air" {
				sbname = "_" + bname
			}
			filename[k] = "cv_" + dname + swidth + sdepth + sheight + sbname + ".mcfunction"
			table.Append([]string{filename[k], width_str, depth_str, height_str, bname})
		}
	}
	table.Render()

	// Now actually create the ClearVol functions
	for i := 0; i < maxdim; i++ {
		for j := 0; j < ndirvals; j++ {
			direction := directionValues[j]
			k := j + i*ndirvals
			err := CreateClearVol(basepath, filename[k], direction,
				mcfdInput.ClearVolHeight[i],
				mcfdInput.ClearVolWidth[i],
				mcfdInput.ClearVolDepth[i],
				mcfdInput.ClearVolBlockType[i])
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

// CreateClearVol
// Clear a volume given a direction and user input.
func CreateClearVol(basepath string, filename string, direction string,
	height int, width int, depth int, btype string) error {

	fname := basepath + "/ClearVol/" + filename
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("CreateClearVol open %v: %v", fname, err)
	}
	defer f.Close()

	// When facing north, the depth is the negative Z coordinate and the width is positive X
	// height is the Y coordinate
    x1 := -width/2
	x2 := -x1 - 1
	r := width - 2 * (width/2)
	if r != 0 {
		x2 += 1
	}
	z1 := -2
	z2 := z1 - depth + 1

	for y := 0; y < height; y++ {
		WriteClearVolBox(x1, y, z1, x2, y, z2, btype, direction, f)
	}

	return nil
}

// WriteClearVolBox writes out a low level box for the wall.
func WriteClearVolBox(x1 int, y1 int, z1 int, x2 int, y2 int, z2 int,
	block_type string, direction string, f *os.File) error {

	corner1 := mcshapes.XYZ{X: x1, Y: y1, Z: z1}
	corner2 := mcshapes.XYZ{X: x2, Y: y2, Z: z2}
	b := mcshapes.NewBox(mcshapes.WithCorner1(corner1), mcshapes.WithCorner2(corner2),
		mcshapes.WithSurface("minecraft:"+block_type))
	b.Orient(direction)
	err := b.WriteShape(f)
	if err != nil {
		return fmt.Errorf("CreateClearVol: %v", err)
	}
	return nil
}
