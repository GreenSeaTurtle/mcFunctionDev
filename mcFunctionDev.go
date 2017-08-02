package main

import (
	//"flag"
	"fmt"
	//"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	//mcshapes "github.com/GreenSeaTurtle/mcFunctionDev/mcShapes"
)

// mcFunctionPath struct for reading various things from the init file
// Note that fields must start with a capital letter!!!!!!
// Example:
//    MCSavesDir - this is what TOML uses below to reference the user input.
//    mc_saves_dir - this is what appears in the init file
type mcFunctionPath struct {
	Title          string
	MCSavesDir     string `toml:"mc_saves_dir"`
	MCFunctionsDir string `toml:"mc_world_functions_dir"`
}

func main() {
	// mcFunctionDev uses two control files, init and input.
	//    init file - sets things that do not change often
	//    input file - controls what mcFunctionDev does when executed

	//
	// The TOML package is used to read and parse both the init file
	// and the input file.
	//    github.com/BurntSushi/toml
	//

	// Read and extract information from the init file.
	//
	// Right now, the only information in the init file is the path to
	// the Minecraft functions directory on this system.  The
	// output function files are written directly to the game directory
	// which saves time and hassle of copying files.  The path is
	// split into two strings just because it is typically a long path.
	gopath := os.Getenv("GOPATH")
	initfile := gopath + "/mcFunctionDev.init"
	//fmt.Println("initfile = " + initfile)
	var mcwpath mcFunctionPath
	if _, err := toml.DecodeFile(initfile, &mcwpath); err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Printf("Title: %s\n", mcwpath.Title)
	//fmt.Printf("mc_saves_dir: %s\n", mcwpath.MCSavesDir)
	//fmt.Println("mc_world_functions_dir = " + mcwpath.MCFunctionsDir)

	// Keep this for now as an example of how to get and process
	// execution line arguments.
	//flag.StringVar(&mcSavesDir, "s", "~", "Minecraft saves directory")
	//flag.StringVar(&mcWorldFuncDir, "w", "mc", "Minecraft functions directory")
	//flag.Parse()

	inputFile := "all.input"
	basepath := path.Join(mcwpath.MCSavesDir, mcwpath.MCFunctionsDir)


	//fmt.Println("basepath = " + basepath)
	//err := BuildFalls(inputFile, basepath)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//CreateClearVolDriver(inputFile, basepath)
	//CreateMWallDriver(inputFile, basepath)
	CreateSign7Driver(inputFile, basepath)
	//CreateSphereDriver(inputFile, basepath)
	//CreateWalkwayDriver(inputFile, basepath)
}

