package main

import mcshapes "github.com/GreenSeaTurtle/mcFunctionDev/mcShapes"

//CreateWaterfall creates a water at origin with attributes
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
	b1 := mcshapes.NewBox(
		mcshapes.WithCorner1(xyz1),
		mcshapes.WithCorner2(xyz2),
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
