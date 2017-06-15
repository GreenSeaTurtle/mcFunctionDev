package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	mcshapes "github.com/GreenSeaTurtle/mcFunctionDev/mcShapes"
	"github.com/olekukonko/tablewriter"
	"github.com/benmcclelland/mcrender"
)

// Structure for using TOML to extract input from the user.
type mcfdControlStruct struct {
	SphereRadius            []int    `toml:"SphereRadius"`
	SphereExteriorBlockType []string `toml:"SphereExteriorBlockType"`
	SphereInteriorBlockType []string `toml:"SphereInteriorBlockType"`
}

// CreateSphereDriver
// Driver for creating the sphere Minecraft function files.
func CreateSphereDriver(inputFile string, basepath string) {
	// Extract pertinent input, using TOML, from the user input file
	var mcfdInput mcfdControlStruct
	if _, err := toml.DecodeFile(inputFile, &mcfdInput); err != nil {
		fmt.Println(err)
		return
	}

	// Consistency check on the user input
	dim := [3]int{0, 0, 0}
	dim[0] = len(mcfdInput.SphereRadius)
	dim[1] = len(mcfdInput.SphereExteriorBlockType)
	dim[2] = len(mcfdInput.SphereInteriorBlockType)
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
		fmt.Println("CreateSphere user input FATAL ERROR")
		fmt.Println("You must specify the same number of array values for all the")
		fmt.Println("CreateSphere arrays - SphereRadius, SphereExteriorBlockType, ...")
		return
	}

	// Create the spheres requested by the user in the user input file.
	if maxdim > 0 {
		// First echo user input to stdout so the user knows what was done.
		// This also sets the filename to write the Minecraft function data.
		fmt.Println("\nCreating Sphere Functions for Minecraft")
		fmt.Println("The following table summarizes user input for the spheres:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Filename", "Radius", "Exterior", "Interior"})
		filename := make([]string, maxdim)
		for i := 0; i < maxdim; i++ {
			// Minecraft functions must have a suffix of ".mcfunction"
			srad := fmt.Sprintf("%d", mcfdInput.SphereRadius[i])
			blkname := mcfdInput.SphereExteriorBlockType[i]
			if mcfdInput.SphereInteriorBlockType[i] != "none" {
				blkname = mcfdInput.SphereExteriorBlockType[i] + "_" +
					mcfdInput.SphereInteriorBlockType[i]
			}
			filename[i] = "Sphere_" + blkname + "_" + srad + ".mcfunction"

			table.Append([]string{filename[i], srad, mcfdInput.SphereExteriorBlockType[i],
				mcfdInput.SphereInteriorBlockType[i]})
		}
		table.Render()

		// Now actually create the sphere functions
		for i := 0; i < maxdim; i++ {
			err := CreateSphere(basepath, filename[i], mcfdInput.SphereRadius[i],
				mcfdInput.SphereExteriorBlockType[i],
				mcfdInput.SphereInteriorBlockType[i])
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

// CreateSphere
func CreateSphere(basepath string, filename string, radius int, exteriorBlockType string,
	interiorBlockType string) error {
	center := mcshapes.XYZ{X: radius, Y: 0, Z: radius + 2}

	fname := basepath + "/" + filename
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("CreateSphere open %v: %v", fname, err)
	}
	defer f.Close()

	b := mcshapes.NewSphere(mcshapes.WithRadius(radius), mcshapes.WithCenter(center),
		mcshapes.WithSphereSurface("minecraft:"+exteriorBlockType),
		mcshapes.WithSphereInteriorSurface("minecraft:"+interiorBlockType))
	err = b.WriteShape(f)
	if err != nil {
		return fmt.Errorf("CreateSphere write mcfunctions: %v", err)
	}

	var buf bytes.Buffer
	if err = b.WriteShape(&buf); err != nil {
		return fmt.Errorf("CreateSphere write to buffer: %v", err)		
	}

	stlname := "stlFiles/" + strings.Replace(filename, "mcfunction", "stl", 1)
	err = mcrender.CreateSTLFromInput(&buf, stlname)
	if err != nil {
		return fmt.Errorf("CreateSphere render to stl file: %v", err)
	}


	return nil
}
