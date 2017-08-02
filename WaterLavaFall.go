package main

import (
	"bytes"
	"log"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	mcshapes "github.com/GreenSeaTurtle/mcFunctionDev/mcShapes"
	"github.com/olekukonko/tablewriter"
	"github.com/benmcclelland/mcrender"
)


// Structure for using TOML to extract input from the user.
type mcfdFallsInputStruct struct {
	FallWidth            []int    `toml:"FallWidth"`
	FallHeight           []int    `toml:"FallHeight"`
	FallFlowBlock        []string `toml:"FallFlowBlock"`
}

//BuildWaterFalls builds n, s, e, w waterfalls
func BuildFalls(inputFile string, basepath string) error {
	// Extract pertinent input, using TOML, from the user input file
	var mcfdInput mcfdFallsInputStruct
	if _, err := toml.DecodeFile(inputFile, &mcfdInput); err != nil {
		fmt.Println(err)
		return nil
	}

	// Consistency check on the user input
	dim := [3]int{0, 0, 0}
	dim[0] = len(mcfdInput.FallWidth)
	dim[1] = len(mcfdInput.FallHeight)
	dim[2] = len(mcfdInput.FallFlowBlock)
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
		fmt.Println("BuildFalls user input FATAL ERROR")
		fmt.Println("You must specify the same number of array values for all the")
		fmt.Println("BuildFalls arrays - FallWidth, FallHeight, ...")
		return nil
	}

	// Create the falls requested by the user in the user input file.
	if maxdim > 0 {
		// First echo user input to stdout so the user knows what was done.
		// This also sets the filename to write the Minecraft function data.
		fmt.Println("\nCreating Falls Functions for Minecraft")
		fmt.Println("The following table summarizes user input for the falls:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Filename", "Width", "Height", "Flow"})
		ndirvals := 8
		filename := make([]string, maxdim*ndirvals)
		filename_rm := make([]string, maxdim*ndirvals)
		filename_cfw := make([]string, maxdim*ndirvals)
		directionValues := []string{"north", "north_ew", "east", "east_sn", "south",
			"south_ew", "west", "west_sn"}
		directionNames := []string{"NWE", "NEW", "ENS", "ESN", "SWE", "SEW", "WNS", "WSN"}
		for i := 0; i < maxdim; i++ {
			for j := 0; j < ndirvals; j++ {
				dname := directionNames[j]
				k := j + i*ndirvals
				// Minecraft functions must have a suffix of ".mcfunction"
				swidth := fmt.Sprintf("%d", mcfdInput.FallWidth[i])
				sheight := fmt.Sprintf("%d", mcfdInput.FallHeight[i])
				blkname := mcfdInput.FallFlowBlock[i] + "fall"
				filename[k] = blkname + "_" + dname + "_" + swidth + "_" +
					sheight + ".mcfunction"
				filename_rm[k] = blkname + "_" + dname + "_" + swidth + "_" +
					sheight + "_rm.mcfunction"
				filename_cfw[k] = blkname + "_" + dname + "_" + swidth + "_" +
					sheight + "_cfw.mcfunction"
				table.Append([]string{filename[k], swidth, sheight,
					mcfdInput.FallFlowBlock[i]})
				table.Append([]string{filename_rm[k], swidth, sheight,
					mcfdInput.FallFlowBlock[i]})
				table.Append([]string{filename_cfw[k], swidth, sheight,
					mcfdInput.FallFlowBlock[i]})
			}
		}
		table.Render()

		origin := mcshapes.XYZ{X: 0, Y: 0, Z: -2}
		for i := 0; i < maxdim; i++ {
			for j := 0; j < ndirvals; j++ {
				direction := directionValues[j]
				k := j + i*ndirvals
				// Create the falls functions
				fname := basepath + "/Falls/" + filename[k]
				f, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
				if err != nil {
					return fmt.Errorf("open FallsBuild %v: %v", fname, err)
				}

				falltype := mcfdInput.FallFlowBlock[i] + "fall"   // lavafall or waterfall
				obj := mcshapes.NewMCObject(mcshapes.WithOrientation(direction),
					mcshapes.WithType(falltype), mcshapes.WithWidth(mcfdInput.FallWidth[i]),
					mcshapes.WithHeight(mcfdInput.FallHeight[i]))
				wf := CreateWaterfall(origin, obj)
				err = mcshapes.WriteShapes(f, wf)
				if err != nil {
					return fmt.Errorf("BuildFalls: %v", err)
				}
				f.Close()

				var buf bytes.Buffer
				if err = mcshapes.WriteShapes(&buf, wf); err != nil {
					return fmt.Errorf("BuildFalls write to buffer: %v", err)		
				}

				// Create an stl file so we can look at the falls before going into the game.
				stlname := "stlFiles/" + strings.Replace(filename[k], "mcfunction", "stl", 1)
				err = mcrender.CreateSTLFromInput(&buf, stlname)
				if err != nil {
					return fmt.Errorf("CreateSphere render to stl file: %v", err)
				}

				// Clear out a buffer area for the falls
				fname = basepath + "/Falls/" + filename_cfw[k]
				f, err = os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
				if err != nil {
					return fmt.Errorf("open falls ClearForWall %v: %v", fname, err)
				}
				ClearForWall(mcfdInput.FallWidth[i], direction, f)
				f.Close()

				// Remove falls
				fname = basepath + "/Falls/" + filename_rm[k]
				f, err = os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
				if err != nil {
					return fmt.Errorf("open rmFalls %v: %v", fname, err)
				}
				rmFalls(mcfdInput.FallWidth[i], mcfdInput.FallHeight[i], direction, f)
				f.Close()
			}
		}
	}

	// Still under development
	err := BuildRollerCoasterFalls(basepath)
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}


// CreateWaterfall creates a water or lava fall at the origin with attributes
func CreateWaterfall(origin mcshapes.XYZ, o *mcshapes.MCObject) []mcshapes.ObjectWriter {
	b := CreateBasin(origin, o)
	b = append(b, CreateSideWall(origin, o, "left")...)
	b = append(b, CreateSideWall(origin, o, "right")...)
	b = append(b, CreateBackWall(origin, o)...)
	b = append(b, CreateBottom(origin, o)...)
	b = append(b, CreateFrontWall(origin, o)...)
	b = append(b, CreateHeater(origin, o)...)
	b = append(b, CreateHeatExchanger(origin, o)...)
	b = append(b, CreateFalls(origin, o)...)

	return b
}

//CreateBasin creates the basin
func CreateBasin(origin mcshapes.XYZ, o *mcshapes.MCObject) []mcshapes.ObjectWriter {
	xyz := mcshapes.XYZ{X: origin.X + o.Width() - 1, Y: origin.Y, Z: origin.Z}
	b1 := mcshapes.NewBox(
		mcshapes.WithCorner1(origin),
		mcshapes.WithCorner2(xyz))
	b1.Orient(o.Orientation())

	xyz = mcshapes.XYZ{X: origin.X, Y: origin.Y, Z: origin.Z - 1}
	b2 := mcshapes.NewBox(mcshapes.WithCorner1(xyz), mcshapes.WithCorner2(xyz))
	b2.Orient(o.Orientation())

	xyz = mcshapes.XYZ{X: origin.X + o.Width() - 1, Y: origin.Y, Z: origin.Z - 1}
	b3 := mcshapes.NewBox(mcshapes.WithCorner1(xyz), mcshapes.WithCorner2(xyz))
	b3.Orient(o.Orientation())

	return append([]mcshapes.ObjectWriter{}, b1, b2, b3)
}

//CreateSideWall creates either left or right side wall
func CreateSideWall(origin mcshapes.XYZ, o *mcshapes.MCObject, side string) []mcshapes.ObjectWriter {
	var x int
	switch side {
	case "left":
		x = origin.X
	case "right":
		x = origin.X + o.Width() - 1
	}

	xyz1 := mcshapes.XYZ{X: x, Y: origin.Y, Z: origin.Z - 2}
	xyz2 := mcshapes.XYZ{X: x, Y: origin.Y + o.Height(), Z: origin.Z - 2}
	b1 := mcshapes.NewBox(mcshapes.WithCorner1(xyz1), mcshapes.WithCorner2(xyz2),
		mcshapes.WithSurface("minecraft:stone 4"))
	b1.Orient(o.Orientation())

	xyz1 = mcshapes.XYZ{X: x, Y: origin.Y + o.Height() - 3, Z: origin.Z - 4}
	xyz2 = mcshapes.XYZ{X: x, Y: origin.Y + o.Height(), Z: origin.Z - 3}
	b2 := mcshapes.NewBox(
		mcshapes.WithCorner1(xyz1),
		mcshapes.WithCorner2(xyz2),
		mcshapes.WithSurface("minecraft:stone 4"))
	b2.Orient(o.Orientation())

	return append([]mcshapes.ObjectWriter{}, b1, b2)
}

//CreateBackWall creates the back wall
func CreateBackWall(origin mcshapes.XYZ, o *mcshapes.MCObject) []mcshapes.ObjectWriter {
	xyz1 := mcshapes.XYZ{X: origin.X, Y: origin.Y + o.Height(), Z: origin.Z - 4}
	xyz2 := mcshapes.XYZ{X: origin.X + o.Width() - 2, Y: origin.Y + o.Height() - 3, Z: origin.Z - 4}
	b := mcshapes.NewBox(
		mcshapes.WithCorner1(xyz1),
		mcshapes.WithCorner2(xyz2),
		mcshapes.WithSurface("minecraft:stone 4"))
	b.Orient(o.Orientation())

	return append([]mcshapes.ObjectWriter{}, b)
}

//CreateBottom creates the bottom of the falls
func CreateBottom(origin mcshapes.XYZ, o *mcshapes.MCObject) []mcshapes.ObjectWriter {
	xyz1 := mcshapes.XYZ{X: origin.X + 1, Y: origin.Y + o.Height() - 3, Z: origin.Z - 3}
	xyz2 := mcshapes.XYZ{X: origin.X + o.Width() - 2, Y: origin.Y + o.Height() - 3, Z: origin.Z - 3}
	b := mcshapes.NewBox(
		mcshapes.WithCorner1(xyz1),
		mcshapes.WithCorner2(xyz2),
		mcshapes.WithSurface("minecraft:stone 4"))
	b.Orient(o.Orientation())

	return append([]mcshapes.ObjectWriter{}, b)
}

//CreateFrontWall creates the front wall for the water to cascade down
func CreateFrontWall(origin mcshapes.XYZ, o *mcshapes.MCObject) []mcshapes.ObjectWriter {
	xyz1 := mcshapes.XYZ{X: origin.X + 1, Y: origin.Y, Z: origin.Z - 2}
	xyz2 := mcshapes.XYZ{X: origin.X + o.Width() - 2, Y: origin.Y + o.Height() - 1, Z: origin.Z - 2}
	b := mcshapes.NewBox(mcshapes.WithCorner1(xyz1), mcshapes.WithCorner2(xyz2))
	b.Orient(o.Orientation())

	return append([]mcshapes.ObjectWriter{}, b)
}

//CreateHeater lava is needed to prevent freezing
func CreateHeater(origin mcshapes.XYZ, o *mcshapes.MCObject) []mcshapes.ObjectWriter {
	xyz1 := mcshapes.XYZ{X: origin.X + 1, Y: origin.Y + o.Height() - 2, Z: origin.Z - 3}
	xyz2 := mcshapes.XYZ{X: origin.X + o.Width() - 2, Y: origin.Y + o.Height() - 2, Z: origin.Z - 3}
	b := mcshapes.NewBox(
		mcshapes.WithCorner1(xyz1),
		mcshapes.WithCorner2(xyz2),
		mcshapes.WithSurface("minecraft:flowing_lava"))
	b.Orient(o.Orientation())

	return append([]mcshapes.ObjectWriter{}, b)
}

//CreateHeatExchanger protects the lava from the water
func CreateHeatExchanger(origin mcshapes.XYZ, o *mcshapes.MCObject) []mcshapes.ObjectWriter {
	xyz1 := mcshapes.XYZ{X: origin.X + 1, Y: origin.Y + o.Height() - 1, Z: origin.Z - 3}
	xyz2 := mcshapes.XYZ{X: origin.X + o.Width() - 2, Y: origin.Y + o.Height() - 1, Z: origin.Z - 3}
	b := mcshapes.NewBox(
		mcshapes.WithCorner1(xyz1),
		mcshapes.WithCorner2(xyz2),
		mcshapes.WithSurface("minecraft:glass"))
	b.Orient(o.Orientation())

	return append([]mcshapes.ObjectWriter{}, b)
}

//CreateFalls creates either the lava or water falls
func CreateFalls(origin mcshapes.XYZ, o *mcshapes.MCObject) []mcshapes.ObjectWriter {
	var surface string
	switch o.OType() {
	case "waterfall":
		surface = "minecraft:flowing_water"
	case "lavafall":
		surface = "minecraft:flowing_lava"
	}

	xyz1 := mcshapes.XYZ{X: origin.X + 1, Y: origin.Y + o.Height(), Z: origin.Z - 3}
	xyz2 := mcshapes.XYZ{X: origin.X + o.Width() - 2, Y: origin.Y + o.Height(), Z: origin.Z - 3}
	b := mcshapes.NewBox(
		mcshapes.WithCorner1(xyz1),
		mcshapes.WithCorner2(xyz2),
		mcshapes.WithSurface(surface))
	b.Orient(o.Orientation())

	return append([]mcshapes.ObjectWriter{}, b)
}




// rmFalls removes north, south, east, west waterfalls
// After placing a water or lava fall somewhere in the Minecraft game, it is sometimes
// necessary to remove it. Perhaps, for example, it was placed in the wrong location and
// needs to be moved. This function writes out a Minecraft function to do this using the
// Minecraft fill command.
// An example of this command is:
//    fill ~0 ~0 ~-2 ~101 ~30 ~-6 minecraft:air
// The ~ refers to the players current position in the game.
// Yes, a fall could be removed by hand inside the game, but this is very tedious, thus
// the need for this function.
func rmFalls(width int, height int, direction string, f *os.File) error {
	origin := mcshapes.XYZ{X: 0, Y: 0, Z: -2}
	// Use a loop because Minecraft has a limit on total number of blocks per fill command.
	for h := 0; h <= height; h++ {
		corner1 := mcshapes.XYZ{X: origin.X,             Y: origin.Y + h, Z: origin.Z}
		corner2 := mcshapes.XYZ{X: origin.X + width - 1, Y: origin.Y + h, Z: origin.Z - 4}
		b := mcshapes.NewBox(mcshapes.WithCorner1(corner1), mcshapes.WithCorner2(corner2),
			mcshapes.WithSurface("minecraft:air"))
		b.Orient(direction)
		err := b.WriteShape(f)
		if err != nil {
			return fmt.Errorf("rm fall: %v", err)
		}
	}

	return nil
}


// ClearForWall - clear space for a wall.
// A wall in this context is meant to surround some area, such as a Mincraft village, and
// provides protection from Minecraft Hostile Mobs (zombies, creeper, spiders, ...). The wall
// should be at least 3 blocks high and it needs an overhang to keep the spiders out (spiders
// can crawl up a wall but cannot get past a ledge). The area inside the wall needs to be lit
// up so Hostile Mobs will not spawn. There needs to be clear space on the outside of the wall
// so the Hostile Mobs will not be able to jump to the top of the wall and thus into the secure
// area. Space is also left on the inside of the wall so villagers will not accidentally find
// their way outside the wall.
// The lava and water falls provide an excellent wall. Such falls are tall enough and come
// with a ledge on the outside. They are also visually stunning. The ledge keeps spiders from
// crawling up and over the wall.
// This function clears space for the wall. The width, height, and depth parameters specify the
// extent of the cleared area. The wall is put in the middle of the cleared area.
func ClearForWall(width int, direction string, f *os.File) error {
	origin := mcshapes.XYZ{X: 0, Y: 0, Z: -2}

	// Minecraft will not accept a width that is too large. 150 is too large, 100 works.
	// Probably has something to do with the size of Minecraft chunks and how many chunks
	// are active and/or being visualized.
	height := 50
	depth := 17

	// Use a loop because Minecraft has a limit on total number of blocks per fill command.
	for h := -1; h <= height; h++ {
		corner1 := mcshapes.XYZ{X: origin.X, Y: origin.Y + h, Z: origin.Z}
		corner2 := mcshapes.XYZ{X: origin.X + width - 1, Y: origin.Y + h, Z: origin.Z - depth + 1}
		block_type := "air"
		if h == -1 {
			block_type = "sea_lantern"
		}
		b := mcshapes.NewBox(mcshapes.WithCorner1(corner1), mcshapes.WithCorner2(corner2),
			mcshapes.WithSurface("minecraft:"+block_type))
		b.Orient(direction)
		err := b.WriteShape(f)
		if err != nil {
			return fmt.Errorf("ClearForWall: %v", err)
		}
	}

	return nil
}



//**************************************************************************************************
//**************************************************************************************************
// Two falls with a roller coaster between them - still under development.
//**************************************************************************************************
//**************************************************************************************************

// BuildRollerCoasterFalls builds two falls next to each other separated by only one
// block. It adds redstone and track to make it a roller coaster ride.
func BuildRollerCoasterFalls(basepath string) error {
	// Create the file that will contain both the north and south waterfalls.
	fname := path.Join(basepath, "Falls/waterfall_rc_north_south.mcfunction")
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open %v: %v", fname, err)
	}
	defer f.Close()

	// Build the north fall - faces south, runs west to east.
	origin := mcshapes.XYZ{X: 2, Y: 0, Z: -2}
	direction := "north"
	obj := mcshapes.NewMCObject(mcshapes.WithOrientation(direction))
	wf := CreateWaterfall(origin, obj)
	err = mcshapes.WriteShapes(f, wf)
	if err != nil {
		return fmt.Errorf("build waterfall rc north fall: %v", err)
	}

	// Build the south fall - faces north, runs west to east.
	// The north and south falls are -1 blocks apart and so they share the same blocks for
	// the front of the basin. In fact, the south falls overwrites what the north falls
	// for the front of the basin.
	// This ends up producing two sheets of water, one block apart. The roller coaster
	// goes between those sheets.
	origin = mcshapes.XYZ{X: 2, Y: 0, Z: 2}
	direction = "south"
	obj = mcshapes.NewMCObject(mcshapes.WithOrientation(direction))
	wf = CreateWaterfall(origin, obj)
	err = mcshapes.WriteShapes(f, wf)
	if err != nil {
		return fmt.Errorf("build waterfall rc south fall: %v", err)
	}

	// At this point we have two waterfalls facing each other, separated by one row of blocks,
	// i.e. the front of the basin which defaults to sandstone.
	// Change those blocks to be redstone in preparation for putting tracks on them.
	// Replace the sandstone with redstone to power the rails.
	width := 102
	corner1 := mcshapes.XYZ{X: origin.X, Y: origin.Y, Z: origin.Z - 4}
	corner2 := mcshapes.XYZ{X: origin.X + width - 1, Y: origin.Y, Z: origin.Z - 4}
	b := mcshapes.NewBox(mcshapes.WithCorner1(corner1), mcshapes.WithCorner2(corner2),
		mcshapes.WithSurface("minecraft:redstone_block"))
	err = b.WriteShape(f)
	if err != nil {
		return fmt.Errorf("build waterfall rc redstone: %v", err)
	}

	// Lay down powered track on top of the redstone.
	corner1.Y += 1
	corner2.Y += 1
	b = mcshapes.NewBox(mcshapes.WithCorner1(corner1), mcshapes.WithCorner2(corner2),
		mcshapes.WithSurface("minecraft:golden_rail"))
	err = b.WriteShape(f)
	if err != nil {
		return fmt.Errorf("build waterfall rc track: %v", err)
	}

	return nil
}

