package num

import (
	"math"
	"testing"
)

// TestNewMatrix3 tests the creation of a new Matrix3
func TestNewMatrix3(t *testing.T) {
	m := NewMatrix3(1, 2, 3, 4, 5, 6, 7, 8, 9)

	// Check all matrix elements
	expected := [9]float64{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	}

	for i := range 9 {
		if m[i] != expected[i] {
			t.Errorf("Matrix element [%d]: expected %f, got %f", i, expected[i], m[i])
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
	result := rotation.Mul(v)
	expected := NewVector3(0, 1, 0)

	// After 90-degree rotation around Z, (1,0,0) should become approximately
	// (0,1,0)
	for i := range result {
		if !almostEqual(result[i], expected[i]) {
			t.Errorf("Rotation matrix multiplication failed: expected %s, got %s", expected, result)
		}
	}
}

// TestVectorValues tests extracting values from a vector
func TestVectorValues(t *testing.T) {
	v := NewVector3(5.5, 6.6, 7.7)
	x, y, z := v.Values()

	if x != 5.5 || y != 6.6 || z != 7.7 {
		t.Errorf("Vector.Values() failed: expected [5.5,6.6,7.7], got %s", v)
	}
}

// Helper function to compare float values with small epsilon
func almostEqual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }
