package mcshapes

import (
	"fmt"
	"io"
)

// Box is a 3D rectangle of blocks (also in 1D, 2D)
// The edge lengths of a Box do not have to be equal.
// All blocks within a Box have to be of the same type.
// The Box geometry is fully specified with the XYZ coordinates of two
// opposite corners.
type Box struct {
	surface string
	corner1 XYZ
	corner2 XYZ
}

// NewBox creates a new box
func NewBox(opts ...BoxOption) *Box {
	b := &Box{
		//default surface is "minecraft:sandstone"
		surface: "minecraft:sandstone",
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

// BoxOption sets various options for NewBox
type BoxOption func(*Box)

// WithSurface set the surface of the block
func WithSurface(surface string) BoxOption {
	return func(b *Box) { b.surface = surface }
}

// WithCorner1 sets the location of the first corner
func WithCorner1(xyz XYZ) BoxOption {
	return func(b *Box) { b.corner1 = xyz }
}

// WithCorner2 sets the location of the opposite corner
func WithCorner2(xyz XYZ) BoxOption {
	return func(b *Box) { b.corner2 = xyz }
}

// At sets the location for a single voxel box
func At(xyz XYZ) BoxOption {
	return func(b *Box) {
		b.corner1 = xyz
		b.corner2 = xyz
	}
}

// Orient box to new direction

// By convention, the user will construct the object while facing
// north. Normally the starting position is 2 blocks in front of the
// player, so when facing north this is z = -2. Examples of objects
// are water/lava falls, spheres, castle walls, ...

// After constructing the object facing north it is often desirable
// to orient the object in some other direction. This is particularly
// true for castle walls, for example, where you typically want the
// walls to be on all 4 sides of an enclosure, north, east, south, and
// west. This function achieves the different orientations using rotations.

// In addition, the user may want to reflect about an axis after the
// rotation. A north castle wall, for example, might run west to east
// starting at the player's position. It is convenient to also have a
// capability to run from east to west. This function also implements
// reflections.

// For this app we are only concerned with 2D rotations in the ZX plane.
// Normally we think of a XY plane but this app is ZX
// With "a" being the angle of rotation, the general rotation matrix is
//    cos(a) -sin(a)
//    sin(a) cos(a)
// For a point, ZX, the rotated point, ZrXr, is
//    Zr = Z*cos(a) - X*sin(a)
//    Xr = Z*sin(a) + X*cos(a)
// We envision being up in the air looking down on the ZX plane.
// Positive angle rotations are counterclockwise, negative are clockwise.
// For now, we only care about the special cases of +/- 90, 180, 270 degrees.
//     90    Zr = -X    Xr = Z      |   -90    Zr = X     Xr = -Z
//    180    Zr = -Z    Xr = -X     |  -180    Zr = -Z    Xr = -X
//    270    Zr = X     Xr = -Z     |  -270    Zr = -X    Xr = Z

// At some point in the future we may want to implement general rotations
// and see if they are useful.

func (b *Box) Orient(direction string) {
	switch direction {
	// No rotation required
	// runs west to east
	case "north":
		return

	// I am amazed that the syntax below of setting 4 values in two lines
	// actually works. It must copy the 4 original values to tmps and then
	// set the new values from the tmps. Doing the 4 operations sequentially
	// with no tmp values would, in some cases, lead to incorrect values.

	// No rotation required, but reflect about Z
	// runs east to west
	case "north_ew":
		b.corner1.X, b.corner2.X = -b.corner1.X, -b.corner2.X


	// 270 ( or -90) degree rotation
	// runs north to south
	case "east":
		b.corner1.X, b.corner2.X, b.corner1.Z, b.corner2.Z =
			-b.corner1.Z, -b.corner2.Z, b.corner1.X, b.corner2.X

	// 270 ( or -90) degree rotation followed by a reflection about X
	// runs south to north
	case "east_sn":
		b.corner1.X, b.corner2.X, b.corner1.Z, b.corner2.Z =
			-b.corner1.Z, -b.corner2.Z, b.corner1.X, b.corner2.X
		b.corner1.Z, b.corner2.Z = -b.corner1.Z, -b.corner2.Z

	// 180 ( or -180) degree rotation
	// runs east to west
	case "south_ew":
		b.corner1.X, b.corner2.X, b.corner1.Z, b.corner2.Z =
			-b.corner1.X, -b.corner2.X, -b.corner1.Z, -b.corner2.Z

	// 180 ( or -180) degree rotation followed by a reflection about Z
	// runs west to east
	case "south":
		b.corner1.X, b.corner2.X, b.corner1.Z, b.corner2.Z =
			-b.corner1.X, -b.corner2.X, -b.corner1.Z, -b.corner2.Z
		b.corner1.X, b.corner2.X = -b.corner1.X, -b.corner2.X

	// 90 ( or -270) degree rotation
	// runs south to north
	case "west_sn":
		b.corner1.X, b.corner2.X, b.corner1.Z, b.corner2.Z =
			b.corner1.Z, b.corner2.Z, -b.corner1.X, -b.corner2.X

	// 90 ( or -270) degree rotation followed by a reflection about X
	// runs north to south
	case "west":
		b.corner1.X, b.corner2.X, b.corner1.Z, b.corner2.Z =
			b.corner1.Z, b.corner2.Z, -b.corner1.X, -b.corner2.X
		b.corner1.Z, b.corner2.Z = -b.corner1.Z, -b.corner2.Z


	// Original transformations.
	// North fall runs west to east
	// East fall runs north to south
	// South fall runs west to east
	// West fall runs north to south

	//case "east":
	//	b.corner1.X, b.corner2.X, b.corner1.Z, b.corner2.Z =
	//		-b.corner2.Z, -b.corner1.Z, b.corner2.X, b.corner1.X

	//case "south":
	//	b.corner1.X, b.corner2.X, b.corner1.Z, b.corner2.Z =
	//		b.corner2.X, b.corner1.X, -b.corner2.Z, -b.corner1.Z

	//case "west":
	//	b.corner1.X, b.corner2.X, b.corner1.Z, b.corner2.Z =
	//		b.corner1.Z, b.corner2.Z, b.corner1.X, b.corner2.X

	}
}

// WriteShape satisfies ObjectWriter interface
func (b *Box) WriteShape(w io.Writer) error {
	s := fmt.Sprintf("fill ~%d ~%d ~%d ~%d ~%d ~%d %s\n",
		b.corner1.X, b.corner1.Y, b.corner1.Z,
		b.corner2.X, b.corner2.Y, b.corner2.Z,
		b.surface)
	_, err := w.Write([]byte(s))
	if err != nil {
		return err
	}
	return nil
}
