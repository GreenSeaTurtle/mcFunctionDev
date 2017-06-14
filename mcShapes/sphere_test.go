package mcshapes

import (
	"bytes"
	"testing"
	//"errors"

	"github.com/benmcclelland/mcrender"
)

// Test a sphere
func TestSphere(t *testing.T) {
	b := NewSphere(WithSphereSurface("testsurface"),
		WithRadius(50),
		//hovering sphere
		WithCenter(XYZ{Y: 70}))

	var buf bytes.Buffer
	if err := b.WriteShape(&buf); err != nil {
		t.Errorf("Create sphere test WriteShape: %v", err)
	}

	err := mcrender.CreateSTLFromInput(&buf, "spheretest.stl")

	// To test if the testing code actually works uncomment the next line
	// and uncomment the errors import above
	//err = errors.New("fail set intentionally")

	if err != nil {
		t.Errorf("Create sphere STL test: %v", err)
	}
}
