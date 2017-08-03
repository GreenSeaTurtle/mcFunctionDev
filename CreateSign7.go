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

//**************************************************************************************************
//**************************************************************************************************
// Creating the Minecraft function files for 7 block tall letter signs
//
// This is called Sign7 because it uses 7 block tall letters. The "7" distinuishes from other
// types of signs such as those using 3 block tall letters.
//**************************************************************************************************
//**************************************************************************************************


// Structure for using TOML to extract input from the user.
//    Sign7Index          The index for naming the sign function. A possible name for the function
//                        files is to use the text or the block types but that leads to function
//                        names that are too long. Therefore simply use an index and rely on the user
//                        refering to the input file to know what text the index is associated with.
//    Sign7Text1          Text for the first line
//    Sign7Text2          Text for the second line, can be blank
//    Sign7Text3          Text for the thrid line,  can be blank
//    Sign7BackBlockType  Block for the backing of the sign. Can be "none" for letters
//                            that hang in the air.
//    Sign7EdgeBlockType  Block for edge of the sign. Can be the same as the backing.
//                           Can also be "none"
//    Sign7TextBlockType  Block for the text.
//
type mcfdSign7InputStruct struct {
	Sign7Index         []int    `toml:"Sign7Index"`
	Sign7Text1         []string `toml:"Sign7Text1"`
	Sign7Text2         []string `toml:"Sign7Text2"`
	Sign7Text3         []string `toml:"Sign7Text3"`
	Sign7BackBlockType []string `toml:"Sign7BackBlockType"`
	Sign7EdgeBlockType []string `toml:"Sign7EdgeBlockType"`
	Sign7TextBlockType []string `toml:"Sign7TextBlockType"`
}

// Number of text lines the user can input.
// Currently hardwired to 3, i.e. text1, text2, text3
const nlines_inp = 3


// CreateSign7Driver
// Driver for creating the Minecraft function files for 7 block tall letter signs
func CreateSign7Driver(inputFile string, basepath string) {
	// Extract pertinent input, using TOML, from the user input file
	var mcfdInput mcfdSign7InputStruct
	if _, err := toml.DecodeFile(inputFile, &mcfdInput); err != nil {
		fmt.Println(err)
		return
	}

	// Consistency check on the user input
	dim := [7]int{0, 0, 0, 0, 0, 0, 0}
	dim[0] = len(mcfdInput.Sign7Index)
	dim[1] = len(mcfdInput.Sign7Text1)
	dim[2] = len(mcfdInput.Sign7Text2)
	dim[3] = len(mcfdInput.Sign7Text3)
	dim[4] = len(mcfdInput.Sign7BackBlockType)
	dim[5] = len(mcfdInput.Sign7EdgeBlockType)
	dim[6] = len(mcfdInput.Sign7TextBlockType)
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
		fmt.Println("CreateSign7 user input FATAL ERROR")
		fmt.Println("You must specify the same number of array values for all the")
		fmt.Println("CreateSign7 arrays - Sign7Text1, Sign7EdgeBlockType, ...")
		fmt.Println("maxdim = ", maxdim)
		fmt.Println("mindim = ", mindim)
		return
	}

	for i := 0; i < maxdim; i++ {
		if mcfdInput.Sign7Text1[i] == "none" {
			fmt.Println("CreateSign7 user input FATAL ERROR")
			fmt.Println("You must specify something other than none for Sign7Text1")
			fmt.Println("Sign7Text2 and Sign7Text3 can be none")
			return
		}
	}

	// If the user has not specified anything then there is nothing left
	// to do.
	if maxdim <= 0 {
		return
	}

	// Create the signs requested by the user in the user input file.

	// First echo user input to stdout so the user knows what was done.
	// This also sets the filename to write the Minecraft function data.
	fmt.Println("\nCreating 7 Block Tall Sign Functions for Minecraft")
	fmt.Println("The following table summarizes user input for the signs:")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Filename", "Back Blk", "Edge Blk", "Text Blk", "Index"})
	ndirvals := 4
	filename := make([]string, maxdim*ndirvals)
	filename_rm := make([]string, maxdim*ndirvals)
	directionValues := []string{"north", "east", "south_ew", "west_sn"}
	directionNames := []string{"N", "E", "S", "W"}
	for i := 0; i < maxdim; i++ {
		index_str := fmt.Sprintf("%d", i)
		for j := 0; j < ndirvals; j++ {
			dname := directionNames[j]
			k := j + i*ndirvals
			//filename[k] = "sign_" + dname + "_" + mcfdInput.Sign7BackBlockType[i] +
			//	"_" + mcfdInput.Sign7EdgeBlockType[i] +
			//	"_" + mcfdInput.Sign7TextBlockType[i] + "_" + index_str + ".mcfunction"
			filename[k] = "s_" + dname + "_" + index_str + ".mcfunction"
			filename_rm[k] = "s_" + dname + "_" + index_str + "_rm.mcfunction"
			table.Append([]string{filename[k], mcfdInput.Sign7BackBlockType[i],
				mcfdInput.Sign7EdgeBlockType[i], mcfdInput.Sign7TextBlockType[i],
				index_str})
			table.Append([]string{filename_rm[k], mcfdInput.Sign7BackBlockType[i],
				mcfdInput.Sign7EdgeBlockType[i], mcfdInput.Sign7TextBlockType[i],
				index_str})
		}
	}
	table.Render()

	// Now actually write the Sign7 functions, both create and remove.
	for i := 0; i < maxdim; i++ {
		text_inp_arr := [nlines_inp]string{}
		text_inp_arr[0] = mcfdInput.Sign7Text1[i]
		text_inp_arr[1] = mcfdInput.Sign7Text2[i]
		text_inp_arr[2] = mcfdInput.Sign7Text3[i]
		for j := 0; j < ndirvals; j++ {
			direction := directionValues[j]
			k := j + i*ndirvals
			err := CreateSign7(basepath, filename[k], filename_rm[k], direction,
				text_inp_arr[:], mcfdInput.Sign7BackBlockType[i],
				mcfdInput.Sign7EdgeBlockType[i], mcfdInput.Sign7TextBlockType[i])
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

// CreateSign7
// Create a sign that is 7 blocks high
func CreateSign7(basepath string, filename string, filename_rm string, direction string,
	text_inp_arr []string, blk_back string, blk_edge string, blk_text string) error {

	fname := basepath + "/Sign7/" + filename
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("CreateSign7 open %v: %v", fname, err)
	}
	defer f.Close()

	// Every character is in a box 5 blocks wide and 7 blocks tall (except for "I" and "1" which are 3 blocks
	// wide). This is taken to be a 5x7, X,Y grid with the 1,1 point being in the lower left corner. X goes
	// to the right, from 1 to 5, while Y goes up, from 1 to 7. All the blocks in the sign start out as
	// type blk_back. The letters/numbers are formed by replacing specified blocks with type blk_text.
	// For example, blk_back might be the nice blue of lapis while the letters/numbers are in gold.
	// For each letter/character the coords 2D array below gives the X,y coordinates in the 5x7 grid that
	// need to be type blk_text, for example gold.
	//
	// These letters/numbers were found on the internet as a JPG file showing each 5x7 letter with a white
	// block background and black letters/numbers. It was simple but tedious to extract coordinates for
	// each letter/number. The lower left coordinate was taken to be 1,1 instead of 0,0 simply because it
	// was easier to extract the letter/number coordinates when going from 1 to 5 and 1 to 7.
	//
	// The map mapdex is a map from the letter/number as a string and its index into the coords array.
	mapdex := make(map[string]int)
	coords := [][]int{}
	Aarr := []int{1,1, 1,2, 1,3, 1,4, 1,5, 1,6, 5,6, 5,1, 5,2, 5,3, 5,4, 5,5, 5,6, 2,7, 3,7, 4,7, 2,4, 3,4, 4,4}
	coords = append(coords, Aarr); mapdex["A"] = 0
	Barr := []int{1,1, 1,2, 1,3, 1,4, 1,5, 1,6, 1,7, 2,7, 3,7, 4,7, 2,4, 3,4, 4,4, 2,1, 3,1, 4,1, 5,2, 5,3, 5,5, 5,6}
	coords = append(coords, Barr); mapdex["B"] = 1
	Carr := []int{1,2, 1,3, 1,4, 1,5, 1,6, 2,7, 3,7, 4,7, 2,1, 3,1, 4,1, 5,2, 5,6}
	coords = append(coords, Carr); mapdex["C"] = 2
	Darr := []int{1,1, 1,2, 1,3, 1,4, 1,5, 1,6, 1,7, 2,7, 3,7, 4,7, 2,1, 3,1, 4,1, 5,2, 5,3, 5,4, 5,5, 5,6}
	coords = append(coords, Darr); mapdex["D"] = 3
	Earr := []int{1,1, 1,2, 1,3, 1,4, 1,5, 1,6, 1,7, 2,7, 3,7, 4,7, 5,7, 2,1, 3,1, 4,1, 5,1, 2,4, 3,4, 4,4}
	coords = append(coords, Earr); mapdex["E"] = 4
	Farr := []int{1,1, 1,2, 1,3, 1,4, 1,5, 1,6, 1,7, 2,7, 3,7, 4,7, 5,7, 2,4, 3,4, 4,4}
	coords = append(coords, Farr); mapdex["F"] = 5
	Garr := []int{1,2, 1,3, 1,4, 1,5, 1,6, 2,7, 3,7, 4,7, 5,6, 2,1, 3,1, 4,1, 5,2, 5,3, 4,3}
	coords = append(coords, Garr); mapdex["G"] = 6
	Harr := []int{1,1, 1,2, 1,3, 1,4, 1,5, 1,6, 1,7, 5,1, 5,2, 5,3, 5,4, 5,5, 5,6, 5,7, 2,4, 3,4, 4,4}
	coords = append(coords, Harr); mapdex["H"] = 7
	Iarr := []int{1,1, 2,1, 3,1, 1,7, 2,7, 3,7, 2,2, 2,3, 2,4, 2,5, 2,6}
	coords = append(coords, Iarr); mapdex["I"] = 8
	Jarr := []int{1,2, 2,1, 3,1, 4,2, 4,3, 4,4, 4,5, 4,6, 4,7, 5,7, 3,7}
	coords = append(coords, Jarr); mapdex["J"] = 9
	Karr := []int{1,1, 1,2, 1,3, 1,4, 1,5, 1,6, 1,7, 2,4, 3,5, 4,6, 5,7, 3,3, 4,2, 4,1}
	coords = append(coords, Karr); mapdex["K"] = 10
	Larr := []int{1,1, 1,2, 1,3, 1,4, 1,5, 1,6, 1,7, 2,1, 3,1, 4,1, 5,1}
	coords = append(coords, Larr); mapdex["L"] = 11
	Marr := []int{1,1, 1,2, 1,3, 1,4, 1,5, 1,6, 1,7, 5,1, 5,2, 5,3, 5,4, 5,5, 5,6, 5,7, 2,6, 3,5, 3,4, 4,6}
	coords = append(coords, Marr); mapdex["M"] = 12
	Narr := []int{1,1, 1,2, 1,3, 1,4, 1,5, 1,6, 1,7, 5,1, 5,2, 5,3, 5,4, 5,5, 5,6, 5,7, 2,5, 3,4, 4,3}
	coords = append(coords, Narr); mapdex["N"] = 13
	Oarr := []int{1,2, 1,3, 1,4, 1,5, 1,6, 5,2, 5,3, 5,4, 5,5, 5,6, 2,7, 3,7, 4,7, 2,1, 3,1, 4,1}
	coords = append(coords, Oarr); mapdex["O"] = 14
	Parr := []int{1,1, 1,2, 1,3, 1,4, 1,5, 1,6, 1,7, 2,7, 3,7, 4,7, 5,6, 5,5, 4,4, 3,4, 2,4}
	coords = append(coords, Parr); mapdex["P"] = 15
	Qarr := []int{1,2, 1,3, 1,4, 1,5, 1,6, 2,7, 3,7, 4,7, 5,6, 5,5, 5,4, 5,3, 3,3, 4,2, 5,1, 2,1, 3,1}
	coords = append(coords, Qarr); mapdex["Q"] = 16
	Rarr := []int{1,1, 1,2, 1,3, 1,4, 1,5, 1,6, 1,7, 2,7, 3,7, 4,7, 5,6, 5,5, 4,4, 3,4, 2,4, 3,3, 4,2, 5,1}
	coords = append(coords, Rarr); mapdex["R"] = 17
	Sarr := []int{1,1, 2,1, 3,1, 4,1, 5,2, 5,3, 4,4, 3,4, 2,4, 1,5, 1,6, 2,7, 3,7, 4,7, 5,7}
	coords = append(coords, Sarr); mapdex["S"] = 18
	Tarr := []int{3,1, 3,2, 3,3, 3,4, 3,5, 3,6, 3,7, 2,7, 1,7, 4,7, 5,7}
	coords = append(coords, Tarr); mapdex["T"] = 19
	Uarr := []int{2,1, 3,1, 4,1, 1,2, 1,3, 1,4, 1,5, 1,6, 1,7, 5,2, 5,3, 5,4, 5,5, 5,6, 5,7}
	coords = append(coords, Uarr); mapdex["U"] = 20
	Varr := []int{1,7, 1,6, 1,5, 1,4, 2,3, 2,2, 3,1, 4,2, 4,3, 5,4, 5,5, 5,6, 5,7}
	coords = append(coords, Varr); mapdex["V"] = 21
	Warr := []int{1,7, 1,6, 1,5, 1,4, 1,3, 1,2, 2,1, 3,2, 3,3, 3,4, 4,1, 5,2, 5,3, 5,4, 5,5, 5,6, 5,7}
	coords = append(coords, Warr); mapdex["W"] = 22
	Xarr := []int{1,7, 1,6, 2,5, 3,4, 4,3, 5,2, 5,1, 1,1, 1,2, 2,3, 4,5, 5,6, 5,7}
	coords = append(coords, Xarr); mapdex["X"] = 23
	Yarr := []int{1,7, 1,6, 2,5, 3,4, 4,5, 5,6, 5,7, 3,1, 3,2, 3,3}
	coords = append(coords, Yarr); mapdex["Y"] = 24
	Zarr := []int{1,7, 2,7, 3,7, 4,7, 5,7, 5,6, 4,5, 3,4, 2,3, 1,2, 1,1, 2,1, 3,1, 4,1, 5,1}
	coords = append(coords, Zarr); mapdex["Z"] = 25
	arr0 := []int{1,2, 1,3, 1,4, 1,5, 1,6, 5,2, 5,3, 5,4, 5,5, 5,6, 2,7, 3,7, 4,7, 2,1, 3,1, 4,1, 2,3, 3,4, 4,5}
	coords = append(coords, arr0); mapdex["0"] = 26
	arr1 := []int{1,6, 2,7, 2,6, 2,5, 2,4, 2,3, 2,2, 2,1, 3,1, 1,1}
	coords = append(coords, arr1); mapdex["1"] = 27
	arr2 := []int{1,6, 2,7, 3,7, 4,7, 5,6, 5,5, 4,4, 3,3, 2,2, 1,1, 2,1, 3,1, 4,1, 5,1}
	coords = append(coords, arr2); mapdex["2"] = 28
	arr3 := []int{1,7, 2,7, 3,7, 4,7, 5,7, 4,6, 3,5, 4,4, 5,3, 5,2, 4,1, 3,1, 2,1, 1,2}
	coords = append(coords, arr3); mapdex["3"] = 29
	arr4 := []int{1,3, 1,4, 2,5, 3,6, 4,7, 4,6, 4,5, 4,4, 4,3, 4,2, 4,1, 2,3, 3,3, 5,3}
	coords = append(coords, arr4); mapdex["4"] = 30
	arr5 := []int{5,7, 4,7, 3,7, 2,7, 1,7, 1,6, 1,5, 2,5, 3,5, 4,5, 5,4, 5,3, 5,2, 4,1, 3,1, 2,1, 1,2}
	coords = append(coords, arr5); mapdex["5"] = 31
	arr6 := []int{4,7, 3,7, 2,6, 1,5, 1,4, 1,3, 1,2, 2,4, 3,4, 4,4, 5,3, 5,2, 4,1, 3,1, 2,1}
	coords = append(coords, arr6); mapdex["6"] = 32
	arr7 := []int{1,7, 2,7, 3,7, 4,7, 5,7, 5,6, 4,5, 3,4, 2,3, 2,2, 2,1}
	coords = append(coords, arr7); mapdex["7"] = 33
	arr8 := []int{2,7, 3,7, 4,7, 5,6, 5,5, 4,4, 3,4, 2,4, 1,5, 1,6, 1,3, 1,2, 2,1, 3,1, 4,1, 5,2, 5,3}
	coords = append(coords, arr8); mapdex["8"] = 34
	arr9 := []int{2,7, 3,7, 4,7, 5,6, 5,5, 5,4, 4,4, 3,4, 2,4, 1,5, 1,6, 5,3, 4,2, 3,1, 2,1}
	coords = append(coords, arr9); mapdex["9"] = 35

	// Get the number of blocks needed to display each line of text.
	// Also get the number of lines of text the user wants in the sign.
	//    nlines_int   Max number of lines the user could specify (a constant)
	//    nlines_text  The number of lines the user actually specified
	//                 nlines_text <= nlines_int
	nblocks_text := [nlines_inp]int{}
	nlines_text := 0
	for i := 0; i < nlines_inp; i++ {
		if text_inp_arr[i] != "none" {
			nblocks_text[i] = GetNblocksInString(text_inp_arr[i])
			fmt.Println("i=",i, "  nblocks_text[i]=",nblocks_text[i])
			nlines_text++
		}
	}
	fmt.Println("nlines_text=",nlines_text)	

	// Calculate the number of blocks for the total width of the sign.
	//     nblocks_text1       Number of blocks needed for the characters
	//     len(text1) - 1      One space between characters
	//     2 + 2               Left and right edge borders
	width_tot_text := [nlines_inp]int{}
	width_tot := 0
	for i := 0; i < nlines_inp; i++ {
		if text_inp_arr[i] != "none" {
			width_tot_text[i] = nblocks_text[i] + len(text_inp_arr[i]) - 1 + 2 + 2
			if width_tot_text[i] > width_tot {
				width_tot = width_tot_text[i]
			}
		}
	}
	fmt.Println("width_tot = ", width_tot)

	// Calculate the starting block X coordinate for each line of text.
	// Currently every line of text is centered.
	xs_text := [nlines_inp]int{}
	for i := 0; i < nlines_inp; i++ {
		xs_text[i] = 2 + (width_tot - width_tot_text[i]) / 2
	}

    // Block index for total width, i.e. 0 to bw
	bw := width_tot - 1

    // Block index for total height
	//     7*nline_text     Each line of text is 7 blocks tall
	//     nlines_text - 1  Spacing between the lines, 1 block
	//     2 + 2            Top and bottom edges
	//     -1               Make bh an index, i.e. 0 to bh
	bh := 7*nlines_text + (nlines_text - 1) + (2 + 2) - 1

	// Render the back of the sign and the edges.
	WriteSign7Box(0,  0,  -2,   bw, bh, -2,  blk_back, direction, f)    // Back of the sign
	WriteSign7Box(0,  0,  -2,   bw,  0, -2,  blk_edge, direction, f)    // Lower edge
	WriteSign7Box(0,  bh, -2,   bw, bh, -2,  blk_edge, direction, f)    // Upper edge
	WriteSign7Box(0,  0,  -2,   0,  bh, -2,  blk_edge, direction, f)    // Left edge
	WriteSign7Box(bw, 0,  -2,   bw, bh, -2,  blk_edge, direction, f)    // Right edge

	// Render the text.
	xs := 2
	ys := bh - 2 - 7 + 1
	for line := 0; line < nlines_text; line++ {
		xs = xs_text[line]
		//text := text_inp_arr[line]
		if text_inp_arr[line] == "none" {continue}

		// Loop through all the characters on the line drawing each one.
		for n:= 0; n<len(text_inp_arr[line]); n++ {
			// Character being rendered.
			schar := string(text_inp_arr[line][n])

			// All characters are 5 blocks wide except for "I" and "1"
			// which are 3 blocks wide.
			nb := 5
			if schar == "I" || schar == "1" {
				nb = 3
			}

			// ic    Index for the character in the coords array
			// np    Number of points used to render the character
			ic := mapdex[schar]
			np := len(coords[ic]) / 2

			// Loop through the points and render the character.
			for i := 0; i < np; i++ {
				x := xs + coords[ic][i*2] - 1
				y := ys + coords[ic][i*2+1] - 1
				WriteSign7Box(x, y, -2,  x, y, -2, blk_text, direction, f)
			}

			// Go on to the next character
			xs += nb + 1
		}
		// Go on to the next line.
		ys -= 7 + 1
	}

	// Write the file to remove a sign.
	fname_rm := basepath + "/Sign7/" + filename_rm
	f_rm, err_rm := os.OpenFile(fname_rm, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err_rm != nil {
		return fmt.Errorf("CreateSign7 open rm file %v: %v", fname_rm, err_rm)
	}
	defer f_rm.Close()

	// Remove the sign
	WriteSign7Box(0,  0,  -2,   bw, bh, -2,  "air",  direction, f_rm)    // Back of the sign

	return nil
}


// GetNblocksInString - Get the number of blocks needed to display a string.
// This does not include spacing between the characters.
// This would be trivial if all the characters were 5 blocks wide but "I" and
// "1" are 3 blocks wide which complicates things.
func GetNblocksInString(text string) int {
	nblocks_text := 0
	for n:= 0; n<len(text); n++ {
		// All characters are 5 blocks wide except for "I" and "1"
		// which are 3 blocks wide.
		nb := 5
		schar := string(text[n])
		if schar == "I" || schar == "1" {
			nb = 3
		}
		nblocks_text += nb
	}
	return nblocks_text
}



// WriteSign7Box writes out a low level box for the sign.
func WriteSign7Box(x1 int, y1 int, z1 int, x2 int, y2 int, z2 int,
	block_type string, direction string, f *os.File) error {

	corner1 := mcshapes.XYZ{X: x1, Y: y1, Z: z1}
	corner2 := mcshapes.XYZ{X: x2, Y: y2, Z: z2}
	b := mcshapes.NewBox(mcshapes.WithCorner1(corner1), mcshapes.WithCorner2(corner2),
		mcshapes.WithSurface("minecraft:"+block_type))
	b.Orient(direction)
	err := b.WriteShape(f)
	if err != nil {
		return fmt.Errorf("CreateSign7: %v", err)
	}
	return nil
}
