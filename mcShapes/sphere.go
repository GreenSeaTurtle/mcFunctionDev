package mcshapes

import (
	"io"
	"math"
)

// Sphere is a hollow sphere defined by a center
// point and a radius with a given surface
type Sphere struct {
	surface string
	interior_surface string
	radius  int
	center  XYZ
}

// NewSphere creates a new sphere
func NewSphere(opts ...SphereOption) *Sphere {
	s := &Sphere{
		surface: "minecraft:glass",
		interior_surface:  "nothing",   // "nothing" means no interior
		radius: 30,
		center: XYZ{Y: 30}, //default center to bring whole sphere on surface
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// SphereOption sets various options for NewSphere
type SphereOption func(*Sphere)

// WithRadius set the radius of the sphere
func WithRadius(r int) SphereOption {
	return func(s *Sphere) { s.radius = r }
}

// WithSphereSurface set the surface of the sphere
func WithSphereSurface(surface string) SphereOption {
	return func(s *Sphere) { s.surface = surface }
}

// WithSphereInteriorSurface set the surface of the interior of the sphere
// If this has the special value of "nothing" then the interior will be
// left empty.
func WithSphereInteriorSurface(surface string) SphereOption {
	return func(s *Sphere) { s.interior_surface = surface }
}

// WithCenter set the center point of the sphere
// note that the center should be at least radius
// for Y, otherwise sphere is below ground level
func WithCenter(c XYZ) SphereOption {
	return func(s *Sphere) { s.center = c }
}

// WriteShape satisfies ObjectWriter interface
func (s *Sphere) WriteShape(w io.Writer) error {
	var voxels []ObjectWriter
	for x := -s.radius; x <= s.radius; x++ {
		for y := -s.radius; y <= s.radius; y++ {
			for z := -s.radius; z <= s.radius; z++ {
				sqs := math.Pow(float64(x), 2) +
					math.Pow(float64(y), 2) +
					math.Pow(float64(z), 2)
				outline := math.Sqrt(sqs)
				if outline >= float64(s.radius-2) && outline <= float64(s.radius) {
					b := NewBox(
						At(XYZ{X: x + s.center.X, Y: y + s.center.Y, Z: z + s.center.Z}),
						WithSurface(s.surface))
					voxels = append(voxels, b)
				}
				if outline < float64(s.radius-2) {
					if s.interior_surface != "nothing" {
						b := NewBox(
							At(XYZ{X: x + s.center.X, Y: y + s.center.Y, Z: z + s.center.Z}),
							WithSurface(s.interior_surface))
						voxels = append(voxels, b)
					}
				}
			}
		}
	}
	return WriteShapes(w, voxels)
}
