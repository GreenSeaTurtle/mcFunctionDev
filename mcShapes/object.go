package mcshapes

import "io"

// XYZ holds the x, y, and z locations
type XYZ struct {
	X int
	Y int
	Z int
}

// ObjectWriter interface to write out shape objects
type ObjectWriter interface {
	WriteShape(w io.Writer) error
}

// MCObject is the basic attributes of a minecraft object
type MCObject struct {
	width       int
	height      int
	orientation string
	oType       string
	location    XYZ
}

// NewMCObject creates a new minecreaft object
func NewMCObject(opts ...MCOption) *MCObject {
	m := &MCObject{
		//default width is 102
		width: 102,
		//default height is 30
		height: 30,
		//default orientation is north
		orientation: "north",
		//default Wtype is water
		oType: "waterfall",
		//default Location is 0,0,0
		// no need to specify 0 value is default
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// MCOption sets various options for NewMCObject
type MCOption func(*MCObject)

// WithWidth set the width of the object
func WithWidth(width int) MCOption {
	return func(m *MCObject) { m.width = width }
}

// WithHeight set the height of the object
func WithHeight(height int) MCOption {
	return func(m *MCObject) { m.height = height }
}

// WithOrientation set the orientation of the object
// The orientation is expected to be one of:
//   north, south, east, or west
func WithOrientation(orientation string) MCOption {
	return func(m *MCObject) { m.orientation = orientation }
}

// WithType sets the type of object
// The type is expected to be one of:
//   waterfall, lavafall
func WithType(oType string) MCOption {
	return func(m *MCObject) { m.oType = oType }
}

// WithLocation sets the base location of the waterfall
func WithLocation(location XYZ) MCOption {
	return func(m *MCObject) { m.location = location }
}

// Width of the MCObject
func (m MCObject) Width() int {
	return m.width
}

// Height of the MCObject
func (m MCObject) Height() int {
	return m.height
}

// Orientation of the MCObject
func (m MCObject) Orientation() string {
	return m.orientation
}

// OType is the type of object
func (m MCObject) OType() string {
	return m.oType
}

// WriteShapes writes out shapes in minecraft format
func WriteShapes(w io.Writer, shapes []ObjectWriter) error {
	for _, s := range shapes {
		err := s.WriteShape(w)
		if err != nil {
			return err
		}
	}
	return nil
}
