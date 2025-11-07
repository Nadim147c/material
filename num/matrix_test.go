package num

import (
	"math"
	"testing"
)

// TestNewMatrix3 tests the creation of a new Matrix3
func TestNewMatrix3(t *testing.T) {
	m := NewMatrix3(1, 2, 3, 4, 5, 6, 7, 8, 9)

	// Check all matrix elements
	expected := [3][3]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	for i := range 3 {
		for j := range 3 {
			if m[i][j] != expected[i][j] {
				t.Errorf(
					"Matrix element [%d][%d]: expected %f, got %f",
					i,
					j,
					expected[i][j],
					m[i][j],
				)
			}
		}
	}
}

// TestNewVector3 tests the creation of a new Vector3
func TestNewVector3(t *testing.T) {
	v := NewVector3(1, 2, 3)

	if v[0] != 1 || v[1] != 2 || v[2] != 3 {
		t.Errorf("Vector creation failed: expected [1, 2, 3], got %v", v)
	}
}

// TestMatrixMultiply tests multiplying a matrix with a vector
func TestMatrixMultiply(t *testing.T) {
	// Rotation matrix (90 degrees around Z axis)
	rotation := NewMatrix3(
		0, -1, 0,
		1, 0, 0,
		0, 0, 1,
	)

	v := NewVector3(1, 0, 0)
	result := rotation.Multiply(v)
	x, y, z := result.Values()

	// After 90-degree rotation around Z, (1,0,0) should become approximately
	// (0,1,0)
	if !almostEqual(x, 0) || !almostEqual(y, 1) || !almostEqual(z, 0) {
		t.Errorf(
			"Rotation matrix multiplication failed: expected (0,1,0), got (%f,%f,%f)",
			x,
			y,
			z,
		)
	}
}

// TestVectorValues tests extracting values from a vector
func TestVectorValues(t *testing.T) {
	v := NewVector3(5.5, 6.6, 7.7)
	x, y, z := v.Values()

	if x != 5.5 || y != 6.6 || z != 7.7 {
		t.Errorf(
			"Vector.Values() failed: expected (5.5,6.6,7.7), got (%f,%f,%f)",
			x,
			y,
			z,
		)
	}
}

// Helper function to compare float values with small epsilon
func almostEqual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }
