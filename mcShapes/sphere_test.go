package mcshapes

import (
	"bytes"
	"testing"

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
		t.Errorf("WriteShape: %v", err)
	}

	err := mcrender.CreateSTLFromInput(&buf, "spheretest.stl")
	if err != nil {
		t.Errorf("create STL: %v", err)
	}
}
