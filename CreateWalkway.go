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
type mcfdWalkwayInputStruct struct {
	WalkwayLength []int `toml:"WalkwayLength"`
}

// CreateWalkwayDriver
// Driver for creating the walkway Minecraft function files.
func CreateWalkwayDriver(inputFile string, basepath string) {
	// Extract pertinent input, using TOML, from the user input file
	var mcfdInput mcfdWalkwayInputStruct
	if _, err := toml.DecodeFile(inputFile, &mcfdInput); err != nil {
		fmt.Println(err)
		return
	}

	// Create the walkways requested by the user in the user input file.
	dim := len(mcfdInput.WalkwayLength)
	if dim > 0 {
		// First echo user input to stdout so the user knows what was done.
		// This also sets the filename to write the Minecraft function data.
		fmt.Println("\nCreating Walkway Functions for Minecraft")
		fmt.Println("The following table summarizes user input for the walkways:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Filename", "Length"})
		filename := make([]string, dim*4)
		filename_cap := make([]string, dim*4)
		directionValues := []string{"north", "east", "south", "west"}
		for i := 0; i < dim; i++ {
			for j := 0; j < 4; j++ {
				direction := directionValues[j]
				k := j + i*4
				// Minecraft functions must have a suffix of ".mcfunction"
				slen := fmt.Sprintf("%d", mcfdInput.WalkwayLength[i])
				filename[k] = "Walkway_" + direction + "_" + slen + ".mcfunction"
				filename_cap[k] = "Walkway_" + direction + "_cap.mcfunction"
				table.Append([]string{filename[k], slen})
			}
		}
		table.Render()

		// Now actually create the walkway functions
		for i := 0; i < dim; i++ {
			for j := 0; j < 4; j++ {
				direction := directionValues[j]
				k := j + i*4
				err := CreateWalkway(basepath, filename[k], direction,
					mcfdInput.WalkwayLength[i])
				if err != nil {
					log.Fatalln(err)
				}
			}
		}


		// Create the walkway cap functions
		for i := 0; i < dim; i++ {
			for j := 0; j < 4; j++ {
				direction := directionValues[j]
				k := j + i*4
				err := CreateWalkwayCap(basepath, filename_cap[k], direction)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}
}


// CreateWalkway
func CreateWalkway(basepath string, filename string, direction string,
	wlength int) error {

	fname := basepath + "/" + filename
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("CreateWalkway open %v: %v", fname, err)
	}
	defer f.Close()

	// First replace everything in the walkway with air, i.e.
	// clear it out.
	WriteWalkwayBox( -8, 0, -1,  8, 0, -wlength, "air", direction, f)
	WriteWalkwayBox( -8, 1, -1,  8, 1, -wlength, "air", direction, f)
	WriteWalkwayBox( -7, 2, -1,  7, 2, -wlength, "air", direction, f)
	WriteWalkwayBox( -7, 3, -1,  7, 3, -wlength, "air", direction, f)
	WriteWalkwayBox( -6, 4, -1,  6, 4, -wlength, "air", direction, f)
	WriteWalkwayBox( -5, 5, -1,  5, 5, -wlength, "air", direction, f)
	WriteWalkwayBox( -4, 6, -1,  4, 6, -wlength, "air", direction, f)
	WriteWalkwayBox( -2, 7, -1,  2, 7, -wlength, "air", direction, f)
	WriteWalkwayBox( -2, 8, -1,  2, 8, -wlength, "air", direction, f)

	// Layer beneath the character.
	WriteWalkwayBox( 0, -1, -1,  0, -1, -wlength, "gold_block", direction, f)
	WriteWalkwayBox(-1, -1, -1, -1, -1, -wlength, "glowstone", direction, f)
	WriteWalkwayBox( 1, -1, -1,  1, -1, -wlength, "glowstone", direction, f)
	WriteWalkwayBox(-2, -1, -1, -2, -1, -wlength, "lapis_block", direction, f)
	WriteWalkwayBox( 2, -1, -1,  2, -1, -wlength, "lapis_block", direction, f)
	WriteWalkwayBox(-3, -1, -1, -3, -1, -wlength, "redstone_block", direction, f)
	WriteWalkwayBox( 3, -1, -1,  3, -1, -wlength, "redstone_block", direction, f)
	WriteWalkwayBox(-3,  0, -1, -3,  0, -wlength, "golden_rail", direction, f)
	WriteWalkwayBox( 3,  0, -1,  3,  0, -wlength, "golden_rail", direction, f)
	WriteWalkwayBox(-4, -1, -1, -4, -1, -wlength, "lapis_block", direction, f)
	WriteWalkwayBox( 4, -1, -1,  4, -1, -wlength, "lapis_block", direction, f)
	WriteWalkwayBox(-5, -1, -1, -5, -1, -wlength, "sea_lantern", direction, f)
	WriteWalkwayBox( 5, -1, -1,  5, -1, -wlength, "sea_lantern", direction, f)
	WriteWalkwayBox(-6, -1, -1, -8, -1, -wlength, "stone 4", direction, f)
	WriteWalkwayBox( 6, -1, -1,  8, -1, -wlength, "stone 4", direction, f)

	// Work up the left and right sides.
	WriteWalkwayBox(-8,  0, -1, -8,  1, -wlength, "glass",   direction, f)
	WriteWalkwayBox( 8,  0, -1,  8,  1, -wlength, "glass",   direction, f)
	WriteWalkwayBox(-7,  2, -1, -7,  3, -wlength, "glass",   direction, f)
	WriteWalkwayBox( 7,  2, -1,  7,  3, -wlength, "glass",   direction, f)
	WriteWalkwayBox(-6,  4, -1, -6,  4, -wlength, "stone 4", direction, f)
	WriteWalkwayBox( 6,  4, -1,  6,  4, -wlength, "stone 4", direction, f)
	WriteWalkwayBox(-5,  5, -1, -5,  5, -wlength, "stone 4", direction, f)
	WriteWalkwayBox( 5,  5, -1,  5,  5, -wlength, "stone 4", direction, f)
	WriteWalkwayBox(-4,  6, -1, -3,  6, -wlength, "stone 4", direction, f)
	WriteWalkwayBox( 4,  6, -1,  3,  6, -wlength, "stone 4", direction, f)
	WriteWalkwayBox(-2,  7, -1,  2,  7, -wlength, "stone 4", direction, f)

	// Place the chandeliers
	for z := -3; z > -wlength; z -= 5 {
		WriteWalkwayBox( 0, 6, z,   0,  6, z,   "fence", direction, f)
		WriteWalkwayBox( 0, 5, z,   0,  5, z,   "fence", direction, f)
		WriteWalkwayBox(-1, 5, z,   -1, 5, z,   "fence", direction, f)
		WriteWalkwayBox( 1, 5, z,   1,  5, z,   "fence", direction, f)
		WriteWalkwayBox( 0, 5, z-1, 0,  5, z-1, "fence", direction, f)
		WriteWalkwayBox( 0, 5, z+1, 0,  5, z+1, "fence", direction, f)

		WriteWalkwayBox(-1, 4, z,   -1, 4, z,   "glowstone", direction, f)
		WriteWalkwayBox( 1, 4, z,    1, 4, z,   "glowstone", direction, f)
		WriteWalkwayBox( 0, 4, z-1,  0, 4, z-1, "glowstone", direction, f)
		WriteWalkwayBox( 0, 4, z+1,  0, 4, z+1, "glowstone", direction, f)
	}

	// Add a layer of sea lanterns above the walkway so it can be seen
	// from the air.
	WriteWalkwayBox( -2, 8, -1,  2, 8, -wlength, "sea_lantern", direction, f)

	return nil
}


// CreateWalkwayCap
func CreateWalkwayCap(basepath string, filename_cap string,
	direction string) error {

	fname_cap := basepath + "/" + filename_cap
	f, err := os.OpenFile(fname_cap, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("CreateWalkwayCap open %v: %v", fname_cap, err)
	}
	defer f.Close()

	WriteWalkwayBox( -8, -1, -1,  8, -1, -1, "stained_glass 5", direction, f)
	WriteWalkwayBox( -8,  0, -1,  8,  0, -1, "stained_glass 5", direction, f)
	WriteWalkwayBox( -8,  1, -1,  8,  1, -1, "stained_glass 5", direction, f)
	WriteWalkwayBox( -7,  2, -1,  7,  2, -1, "stained_glass 5", direction, f)
	WriteWalkwayBox( -7,  3, -1,  7,  3, -1, "stained_glass 5", direction, f)
	WriteWalkwayBox( -6,  4, -1,  6,  4, -1, "stained_glass 5", direction, f)
	WriteWalkwayBox( -5,  5, -1,  5,  5, -1, "stained_glass 5", direction, f)
	WriteWalkwayBox( -4,  6, -1,  4,  6, -1, "stained_glass 5", direction, f)
	WriteWalkwayBox( -2,  7, -1,  2,  7, -1, "stained_glass 5", direction, f)
	WriteWalkwayBox( -2,  8, -1,  2,  8, -1, "stained_glass 5", direction, f)

	return nil
}


func CreateAngledWalkway(basepath string, nchunk int) error {
	ntot := nchunk * 10
	sntot := fmt.Sprintf("%d", ntot)

	filename := "Walkway_NW_" + sntot + ".mcfunction"
	fname := basepath + "/" + filename
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("CreateWalkway open %v: %v", fname, err)
	}
	defer f.Close()

	direction := "north"
	xt := 0
	zt := 0
	for nc := 0; nc < nchunk; nc++ {
		xt = -nc*10
		zt = xt
		WriteAngledWalkwayPath( xt-1, -1, zt-1,  10, 8, "gold_block",     direction, "n", f)
		WriteAngledWalkwayPath( xt-1, -1, zt-2,  10, 8, "gold_block",     direction, "y", f)
		WriteAngledWalkwayPath( xt+0, -1, zt-2,  10, 8, "glowstone",      direction, "y", f)
		WriteAngledWalkwayPath( xt+0, -1, zt-3,  10, 8, "glowstone",      direction, "y", f)
		WriteAngledWalkwayPath( xt+1, -1, zt-3,  10, 8, "lapis_block",    direction, "y", f)
		WriteAngledWalkwayPath( xt+1, -1, zt-4,  10, 8, "lapis_block",    direction, "y", f)
		WriteAngledWalkwayPath( xt+2, -1, zt-4,  10, 6, "redstone_block", direction, "y", f)
		WriteAngledWalkwayPath( xt+2, -1, zt-5,  10, 6, "redstone_block", direction, "y", f)
		WriteAngledWalkwayPath( xt+3, -1, zt-5,  10, 6, "lapis_block",    direction, "y", f)
		WriteAngledWalkwayPath( xt+3, -1, zt-6,  10, 6, "lapis_block",    direction, "y", f)
		WriteAngledWalkwayPath( xt+4, -1, zt-6,  10, 5, "sea_lantern",    direction, "y", f)
		WriteAngledWalkwayPath( xt+4, -1, zt-7,  10, 5, "sea_lantern",    direction, "y", f)
		WriteAngledWalkwayPath( xt+5, -1, zt-7,  10, 4, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt+5, -1, zt-8,  10, 4, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt+6, -1, zt-8,  10, 3, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt+6, -1, zt-9,  10, 3, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt+7, -1, zt-9,  10, 1, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt+7,  0, zt-9,  10, 0, "glass",          direction, "y", f)
		WriteAngledWalkwayPath( xt+7,  1, zt-9,  10, 0, "glass",          direction, "y", f)
		WriteAngledWalkwayPath( xt+6,  2, zt-8,  10, 0, "glass",          direction, "y", f)
		WriteAngledWalkwayPath( xt+6,  2, zt-9,  10, 0, "glass",          direction, "y", f)
		WriteAngledWalkwayPath( xt+6,  3, zt-8,  10, 0, "glass",          direction, "y", f)
		WriteAngledWalkwayPath( xt+6,  3, zt-9,  10, 0, "glass",          direction, "y", f)
		WriteAngledWalkwayPath( xt+5,  4, zt-7,  10, 0, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt+5,  4, zt-8,  10, 0, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt+4,  5, zt-6,  10, 0, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt+4,  5, zt-7,  10, 0, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt+3,  6, zt-5,  10, 0, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt+3,  6, zt-6,  10, 0, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt+2,  6, zt-4,  10, 0, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt+2,  6, zt-5,  10, 0, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt+1,  7, zt-3,  10, 0, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt+1,  7, zt-4,  10, 0, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt+0,  7, zt-2,  10, 0, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt+0,  7, zt-3,  10, 0, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt-1,  7, zt-2,  10, 0, "stone 4",        direction, "y", f)
		WriteAngledWalkwayPath( xt-1,  7, zt-1,  10, 0, "stone 4",        direction, "n", f)
		WriteAngledWalkwayPath( xt+1,  8, zt-3,  10, 0, "sea_lantern",    direction, "y", f)
		WriteAngledWalkwayPath( xt+1,  8, zt-4,  10, 0, "sea_lantern",    direction, "y", f)
		WriteAngledWalkwayPath( xt+0,  8, zt-2,  10, 0, "sea_lantern",    direction, "y", f)
		WriteAngledWalkwayPath( xt+0,  8, zt-3,  10, 0, "sea_lantern",    direction, "y", f)
		WriteAngledWalkwayPath( xt-1,  8, zt-2,  10, 0, "sea_lantern",    direction, "y", f)
		WriteAngledWalkwayPath( xt-1,  8, zt-1,  10, 0, "sea_lantern",    direction, "n", f)

		// Add the rails to the redstone blocks.
		WriteAngledWalkwayPath( xt+2,  0, zt-4,  10, 6, "rail",           direction, "y", f)
		WriteAngledWalkwayPath( xt+2,  0, zt-5,  10, 6, "rail",           direction, "y", f)

		WriteWalkwayBox(xt-7,  -1, zt-12,  xt-7,  -1, zt-12,  "redstone_block",  direction, f)
		WriteWalkwayBox(xt-6,   0, zt-13,  xt-6,   0, zt-13,  "air",             direction, f)
		WriteWalkwayBox(xt-6,   0, zt-12,  xt-6,   0, zt-12,  "golden_rail",     direction, f)
		WriteWalkwayBox(xt-7,   0, zt-12,  xt-7,   0, zt-12,  "rail",            direction, f)
		WriteWalkwayBox(xt-7,   0, zt-13,  xt-7,   0, zt-13,  "golden_rail",     direction, f)

		WriteWalkwayBox(zt-12,  -1, xt-7,  zt-12, -1, xt-7,   "redstone_block",  direction, f)
		WriteWalkwayBox(zt-13,   0, xt-6,  zt-13,  0, xt-6,   "air",             direction, f)
		WriteWalkwayBox(zt-12,   0, xt-6,  zt-12,  0, xt-6,   "golden_rail",     direction, f)
		WriteWalkwayBox(zt-12,   0, xt-7,  zt-12,  0, xt-7,   "rail",            direction, f)
		WriteWalkwayBox(zt-13,   0, xt-7,  zt-13,  0, xt-7,   "golden_rail",     direction, f)

		// Place the chandeliers
		for zc := -3; zc >= -8; zc -= 5 {
			xc := zc
			WriteWalkwayBox(xt+xc,   6, zt+zc,   xt+xc,   6, zt+zc,   "fence", direction, f)
			WriteWalkwayBox(xt+xc,   5, zt+zc,   xt+xc,   5, zt+zc,   "fence", direction, f)
			WriteWalkwayBox(xt+xc-1, 5, zt+zc,   xt+xc-1, 5, zt+zc,   "fence", direction, f)
			WriteWalkwayBox(xt+xc+1, 5, zt+zc,   xt+xc+1, 5, zt+zc,   "fence", direction, f)
			WriteWalkwayBox(xt+xc,   5, zt+zc-1, xt+xc,   5, zt+zc-1, "fence", direction, f)
			WriteWalkwayBox(xt+xc,   5, zt+zc+1, xt+xc,   5, zt+zc+1, "fence", direction, f)

			WriteWalkwayBox(xt+xc-1, 4, zt+zc,   xt+xc-1, 4, zt+zc,   "glowstone", direction, f)
			WriteWalkwayBox(xt+xc+1, 4, zt+zc,   xt+xc+1, 4, zt+zc,   "glowstone", direction, f)
			WriteWalkwayBox(xt+xc,   4, zt+zc-1, xt+xc,   4, zt+zc-1, "glowstone", direction, f)
			WriteWalkwayBox(xt+xc,   4, zt+zc+1, xt+xc,   4, zt+zc+1, "glowstone", direction, f)
		}
	}


	return nil
}


func WriteAngledWalkwayPath(xs int, ys int, zs int, nblocks int, ymax int,
	block_type string, direction string, reflect string, f *os.File) error {

	yv := ys
	for n:=0; n < nblocks; n++ {
		x1 := xs - n
		z1 := zs - n
		if yv == -1 {
			for y := 0; y <= ymax; y++ {
				WriteWalkwayBox( x1, y, z1, x1, y, z1, "air", direction, f)
			}
		}
		WriteWalkwayBox( x1, yv, z1, x1, yv, z1, block_type, direction, f)
		if reflect == "y" {
			z2 := x1
			x2 := z1
			if yv == -1 {
				for y := 0; y <= ymax; y++ {
					WriteWalkwayBox( x2, y, z2, x2, y, z2, "air", direction, f)
				}
			}
			WriteWalkwayBox( x2, yv, z2, x2, yv, z2, block_type, direction, f)
		}
	}

	return nil
}


// WriteWalkwayBox writes out a low level box for the walkway.
func WriteWalkwayBox(x1 int, y1 int, z1 int, x2 int, y2 int, z2 int,
	block_type string, direction string, f *os.File) error {

	corner1 := mcshapes.XYZ{X: x1, Y: y1, Z: z1}
	corner2 := mcshapes.XYZ{X: x2, Y: y2, Z: z2}
	b := mcshapes.NewBox(mcshapes.WithCorner1(corner1), mcshapes.WithCorner2(corner2),
		mcshapes.WithSurface("minecraft:"+block_type))
	b.Orient(direction)
	err := b.WriteShape(f)
	if err != nil {
		return fmt.Errorf("CreateWalkway: %v", err)
	}
	return nil
}
