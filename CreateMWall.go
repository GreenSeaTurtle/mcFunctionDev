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
type mcfdMWallInputStruct struct {
	MWallHeight         []int    `toml:"MWallHeight"`
	MWallWidth          []int    `toml:"MWallWidth"`
	MWallDepth          []int    `toml:"MWallDepth"`
	MWallWoodBlockType  []string `toml:"MWallWoodBlockType"`
	MWallBrickBlockType []string `toml:"MWallBrickBlockType"`
}

// The construction unit for MWall is 2 blocks wide. This unit is duplicated
// as needed to achieve the total desired width.
var conun_width int = 2

// CreateMWallDriver
// Driver for creating the Minecraft function files for this type of
// castle wall.
func CreateMWallDriver(inputFile string, basepath string) {
	// Extract pertinent input, using TOML, from the user input file
	var mcfdInput mcfdMWallInputStruct
	if _, err := toml.DecodeFile(inputFile, &mcfdInput); err != nil {
		fmt.Println(err)
		return
	}

	// Consistency check on the user input
	dim := [5]int{0, 0, 0, 0, 0}
	dim[0] = len(mcfdInput.MWallHeight)
	dim[1] = len(mcfdInput.MWallWidth)
	dim[2] = len(mcfdInput.MWallDepth)
	dim[3] = len(mcfdInput.MWallWoodBlockType)
	dim[4] = len(mcfdInput.MWallBrickBlockType)
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
		fmt.Println("CreateMWall user input FATAL ERROR")
		fmt.Println("You must specify the same number of array values for all the")
		fmt.Println("CreateMWall arrays - MWallHeight, MWallWoodBlockType, ...")
		return
	}


	// If the user has not specified anything then there is nothing left
	// to do.
	if maxdim <= 0 {
		return
	}

	// Create the walls requested by the user in the user input file.

	// First echo user input to stdout so the user knows what was done.
	// This also sets the filename to write the Minecraft function data.
	fmt.Println("\nCreating MWall Functions for Minecraft")
	fmt.Println("The following table summarizes user input for the m walls:")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Filename", "Height", "Width", "Depth", "Wood", "Brick"})
	ndirvals := 8
	filename := make([]string, maxdim*ndirvals)
	filename_rm := make([]string, maxdim*ndirvals)
	directionValues := []string{"north", "north_refl", "east", "east_refl", "south_refl",
		 "south", "west_refl", "west"}
	directionNames := []string{"NWE", "NEW", "ENS", "ESN", "SWE", "SEW", "WNS", "WSN"}
	for i := 0; i < maxdim; i++ {
		for j := 0; j < ndirvals; j++ {
			dname := directionNames[j]
			k := j + i*ndirvals
			// Minecraft functions must have a suffix of ".mcfunction"
			sheight := fmt.Sprintf("%d", mcfdInput.MWallHeight[i])
			swidth := fmt.Sprintf("%d", mcfdInput.MWallWidth[i])
			sdepth := fmt.Sprintf("%d", mcfdInput.MWallDepth[i])
			wood_blkname := mcfdInput.MWallWoodBlockType[i]
			brick_blkname := mcfdInput.MWallBrickBlockType[i]
			//filename[k] = "MWall_" + direction + "_" + sheight + "_" + swidth + "_" +
			//	sdepth + "_" + wood_blkname + "_" + brick_blkname + ".mcfunction"
			filename[k] = "mw_" + dname + "_" + sheight + "_" + swidth + ".mcfunction"
			filename_rm[k] = "mw_" + dname + "_" + sheight + "_" + swidth + "_rm.mcfunction"

			table.Append([]string{filename[k], sheight, swidth, sdepth,
				wood_blkname, brick_blkname})
			table.Append([]string{filename_rm[k], sheight, swidth, sdepth,
				wood_blkname, brick_blkname})
		}
	}
	table.Render()

	// Now actually create the MWall functions
	for i := 0; i < maxdim; i++ {
		for j := 0; j < ndirvals; j++ {
			direction := directionValues[j]
			k := j + i*ndirvals
			err := CreateMWall(basepath, filename[k], direction,
				mcfdInput.MWallHeight[i],
				mcfdInput.MWallWidth[i],
				mcfdInput.MWallDepth[i],
				mcfdInput.MWallWoodBlockType[i],
				mcfdInput.MWallBrickBlockType[i])
			if err != nil {
				log.Fatalln(err)
			}

			err = RmMWall(basepath, filename_rm[k], direction,
				mcfdInput.MWallHeight[i],
				mcfdInput.MWallWidth[i],
				mcfdInput.MWallDepth[i])
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

// CreateMWall
// Create one M wall given a direction and user input for the wall.
func CreateMWall(basepath string, filename string, direction string,
	total_height int, width int, depth int, wood_btype string,
	brick_btype string) error {

	fname := basepath + "/MWall/" + filename
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("CreateMWall open %v: %v", fname, err)
	}
	defer f.Close()

	// When facing north, the depth is the negative Z coordinate and the width is positive X
	// height is the Y coordinate
	height := total_height - 2      // Height of the wood
	total_depth := depth + 3 + 3    // Depth of the wall

	// A single construction unit is 2 columns wide (conun_width). This is repeated
	// as many times as needed to get the desired width.
	nc := width/conun_width

	// near and far are along the negative Z axis when facing north.
	// near is > far, for example near = -2, far = -8
	near_bf := -2                         // Near brick face
	far_bf  := near_bf - total_depth + 1  // Far brick face
	near_wf := -3                         // Near wood face
	far_wf  := near_wf - 3 - depth        // Far wood face

	// The following is done for one construction unit. The low level function
	// duplicates over all construction units.
	// The following is done facing north, i.e. a north wall. The low level function
	// takes care of rotating it for a east, south, and west wall.

	// Clear out the space first
	WriteMWallBox(0, 0, near_bf, 1, total_height-1, far_bf, nc, "air", direction, f)

	// The lower two wood pieces
	WriteMWallBox(0, 0, near_wf, 1, 0, near_wf, nc, wood_btype, direction, f)
	WriteMWallBox(0, 0, far_wf,  1, 0, far_wf,  nc, wood_btype, direction, f)

	// The lower bricks going from near to far
	WriteMWallBox(0, 1, near_bf,   0, 1, far_bf,    nc, brick_btype, direction, f)
	WriteMWallBox(1, 1, near_bf-2, 1, 1, far_bf+2,  nc, brick_btype, direction, f)

	// The wood and brick vertical columns
	WriteMWallBox(1, 1, near_wf, 1, height-2, near_wf, nc, wood_btype, direction, f)
	WriteMWallBox(1, 1, far_wf,  1, height-2, far_wf,  nc, wood_btype, direction, f)

	WriteMWallBox(0, 2, near_bf, 0, height-3, near_bf, nc, brick_btype, direction, f)
	WriteMWallBox(0, 2, far_bf,  0, height-3, far_bf,  nc, brick_btype, direction, f)

	WriteMWallBox(0, 2, near_bf-2, 1, height-3, near_bf-2, nc, brick_btype, direction, f)
	WriteMWallBox(0, 2, far_bf+2, 1, height-3, far_bf+2, nc, brick_btype, direction, f)

	// The upper bricks going from near to far
	WriteMWallBox(0, height-2, near_bf,   0, height-2, far_bf,    nc, brick_btype, direction, f)
	WriteMWallBox(1, height-2, near_bf-2, 1, height-2, far_bf+2,  nc, brick_btype, direction, f)

	// The top two wood peices
	WriteMWallBox(0, height-1, near_wf, 1, height-1, near_wf, nc, wood_btype, direction, f)
	WriteMWallBox(0, height-1, far_wf,  1, height-1, far_wf,  nc, wood_btype, direction, f)

	WriteMWallBox(0, height, near_wf, 0, height, near_wf, nc, "gold_block", direction, f)
	WriteMWallBox(0, height, far_wf,  0, height, far_wf,  nc, "gold_block", direction, f)
	WriteMWallBox(0, height+1, near_wf, 0, height+1, near_wf, nc, "torch", direction, f)
	WriteMWallBox(0, height+1, far_wf,  0, height+1, far_wf,  nc, "torch", direction, f)

	return nil
}

// RmMWall
// Remove one M wall given a direction and user input for the wall.
func RmMWall(basepath string, filename string, direction string,
	total_height int, width int, depth int) error {

	fname := basepath + "/MWall/" + filename
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("CreateMWall open %v: %v", fname, err)
	}
	defer f.Close()

	nc := width/conun_width
	total_depth := depth + 3 + 3

	near_bf := -2                         // Near brick face
	far_bf  := near_bf - total_depth + 1  // Far brick face

	// Clear out the space first
	WriteMWallBox(0, 0, near_bf, 1, total_height-1, far_bf, nc, "air", direction, f)

	return nil
}


// WriteMWallBox writes out a low level box for the wall.
// Duplicate for all the contruction units (nconun)
func WriteMWallBox(x1 int, y1 int, z1 int, x2 int, y2 int, z2 int,
	nconun int, block_type string, direction string, f *os.File) error {

	xt := 0
	for n:=0; n<nconun; n++ {
		xt = n*conun_width
		corner1 := mcshapes.XYZ{X: xt+x1, Y: y1, Z: z1}
		corner2 := mcshapes.XYZ{X: xt+x2, Y: y2, Z: z2}
		b := mcshapes.NewBox(mcshapes.WithCorner1(corner1), mcshapes.WithCorner2(corner2),
			mcshapes.WithSurface("minecraft:"+block_type))
		b.Orient(direction)
		err := b.WriteShape(f)
		if err != nil {
			return fmt.Errorf("CreateMWall: %v", err)
		}
	}
	return nil
}
