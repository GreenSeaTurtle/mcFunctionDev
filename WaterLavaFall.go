package main

import (
	"bytes"
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
		filename := make([]string, maxdim*4)
		directionValues := []string{"north", "east", "south", "west"}
		for i := 0; i < maxdim; i++ {
			for j := 0; j < 4; j++ {
				direction := directionValues[j]
				k := j + i*4
				// Minecraft functions must have a suffix of ".mcfunction"
				swidth := fmt.Sprintf("%d", mcfdInput.FallWidth[i])
				sheight := fmt.Sprintf("%d", mcfdInput.FallHeight[i])
				blkname := mcfdInput.FallFlowBlock[i] + "fall"
				filename[k] = blkname + "_" + direction + "_" + swidth + "_" +
					sheight + ".mcfunction"
				table.Append([]string{filename[k], swidth, sheight,
					mcfdInput.FallFlowBlock[i]})
					
			}
		}
		table.Render()

		origin := mcshapes.XYZ{X: 0, Y: 0, Z: -2}
		for i := 0; i < maxdim; i++ {
			for j := 0; j < 4; j++ {
				direction := directionValues[j]
				k := j + i*4
				fname := path.Join(basepath, filename[k])
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

				stlname := "stlFiles/" + strings.Replace(filename[k], "mcfunction", "stl", 1)
				err = mcrender.CreateSTLFromInput(&buf, stlname)
				if err != nil {
					return fmt.Errorf("CreateSphere render to stl file: %v", err)
				}
			}
		}
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
