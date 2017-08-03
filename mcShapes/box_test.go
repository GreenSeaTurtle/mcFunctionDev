package mcshapes

import (
	"bytes"
	"testing"
)

// Test a north box
func TestNorthBox(t *testing.T) {
	expected := "fill ~1 ~2 ~3 ~4 ~5 ~6 testsurface\n"
	b := NewBox(
		WithSurface("testsurface"),
		WithCorner1(XYZ{X: 1, Y: 2, Z: 3}),
		WithCorner2(XYZ{X: 4, Y: 5, Z: 6}))

	b.Orient("north")

	var buf bytes.Buffer
	if err := b.WriteShape(&buf); err != nil {
		t.Errorf("WriteShape: %v", err)
	}

	if buf.String() != expected {
		t.Errorf("expected '%v', got '%v'", expected, buf.String())
	}
}

func TestSouthBox(t *testing.T) {
	expected := "fill ~4 ~2 ~-6 ~1 ~5 ~-3 testsurface\n"
	b := NewBox(
		WithSurface("testsurface"),
		WithCorner1(XYZ{X: 1, Y: 2, Z: 3}),
		WithCorner2(XYZ{X: 4, Y: 5, Z: 6}))

	b.Orient("south_refl")

	var buf bytes.Buffer
	if err := b.WriteShape(&buf); err != nil {
		t.Errorf("WriteShape: %v", err)
	}

	if buf.String() != expected {
		t.Errorf("expected '%v', got '%v'", expected, buf.String())
	}
}

func TestEastBox(t *testing.T) {
	expected := "fill ~-6 ~2 ~4 ~-3 ~5 ~1 testsurface\n"
	b := NewBox(
		WithSurface("testsurface"),
		WithCorner1(XYZ{X: 1, Y: 2, Z: 3}),
		WithCorner2(XYZ{X: 4, Y: 5, Z: 6}))

	b.Orient("east")

	var buf bytes.Buffer
	if err := b.WriteShape(&buf); err != nil {
		t.Errorf("WriteShape: %v", err)
	}

	if buf.String() != expected {
		t.Errorf("expected '%v', got '%v'", expected, buf.String())
	}
}

func TestWestBox(t *testing.T) {
	expected := "fill ~3 ~2 ~1 ~6 ~5 ~4 testsurface\n"
	b := NewBox(
		WithSurface("testsurface"),
		WithCorner1(XYZ{X: 1, Y: 2, Z: 3}),
		WithCorner2(XYZ{X: 4, Y: 5, Z: 6}))

	b.Orient("west_refl")

	var buf bytes.Buffer
	if err := b.WriteShape(&buf); err != nil {
		t.Errorf("WriteShape: %v", err)
	}

	if buf.String() != expected {
		t.Errorf("expected '%v', got '%v'", expected, buf.String())
	}
}
