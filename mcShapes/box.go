package mcshapes

import (
	"fmt"
	"io"
)

// Box is a 3D rectangle of blocks (also in 1D and 2D)
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

// Orient box to new direction
func (b *Box) Orient(direction string) {
	switch direction {
	case "north":
		return

	case "east":
		b.corner1.X, b.corner2.X, b.corner1.Z, b.corner2.Z =
			-b.corner2.Z, -b.corner1.Z, b.corner2.X, b.corner1.X

	case "south":
		b.corner1.X, b.corner2.X, b.corner1.Z, b.corner2.Z =
			b.corner2.X, b.corner1.X, -b.corner2.Z, -b.corner1.Z

	case "west":
		b.corner1.X, b.corner2.X, b.corner1.Z, b.corner2.Z =
			b.corner1.Z, b.corner2.Z, b.corner1.X, b.corner2.X
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
