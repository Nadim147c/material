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
				t.Errorf("Matrix element [%d][%d]: expected %f, got %f", i, j, expected[i][j], m[i][j])
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

// TestMatrixMultiplyXYZ tests multiplying a matrix with xyz coordinates
func TestMatrixMultiplyXYZ(t *testing.T) {
	// Identity matrix
	identity := NewMatrix3(
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	)

	// Test with identity matrix - should return same coordinates
	result := identity.MultiplyXYZ(1, 2, 3)
	x, y, z := result.Values()

	if x != 1 || y != 2 || z != 3 {
		t.Errorf("Identity matrix multiplication failed: expected (1,2,3), got (%f,%f,%f)", x, y, z)
	}

	// Test with a scaling matrix
	scaling := NewMatrix3(
		2, 0, 0,
		0, 3, 0,
		0, 0, 4,
	)

	result = scaling.MultiplyXYZ(1, 2, 3)
	x, y, z = result.Values()

	if x != 2 || y != 6 || z != 12 {
		t.Errorf("Scaling matrix multiplication failed: expected (2,6,12), got (%f,%f,%f)", x, y, z)
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

	// After 90-degree rotation around Z, (1,0,0) should become approximately (0,1,0)
	if !almostEqual(x, 0) || !almostEqual(y, 1) || !almostEqual(z, 0) {
		t.Errorf("Rotation matrix multiplication failed: expected (0,1,0), got (%f,%f,%f)", x, y, z)
	}
}

// TestVectorMultiplyMatrix tests multiplying a vector with a matrix
func TestVectorMultiplyMatrix(t *testing.T) {
	// Test matrix
	m := NewMatrix3(
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	)

	v := NewVector3(2, 3, 4)
	result := v.MultiplyMatrix(m)
	x, y, z := result.Values()

	// Expected: (2*1 + 3*4 + 4*7, 2*2 + 3*5 + 4*8, 2*3 + 3*6 + 4*9) = (42, 51, 60)
	if !almostEqual(x, 42) || !almostEqual(y, 51) || !almostEqual(z, 60) {
		t.Errorf("Vector-matrix multiplication failed: expected (42,51,60), got (%f,%f,%f)", x, y, z)
	}
}

// TestVectorValues tests extracting values from a vector
func TestVectorValues(t *testing.T) {
	v := NewVector3(5.5, 6.6, 7.7)
	x, y, z := v.Values()

	if x != 5.5 || y != 6.6 || z != 7.7 {
		t.Errorf("Vector.Values() failed: expected (5.5,6.6,7.7), got (%f,%f,%f)", x, y, z)
	}
}

// TestMatrixVectorConsistency tests that m.Multiply(v) and v.MultiplyMatrix(m)
// are not the same - they represent different mathematical operations
func TestMatrixVectorConsistency(t *testing.T) {
	m := NewMatrix3(
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	)

	v := NewVector3(2, 3, 4)

	// M*v - standard matrix-vector multiplication
	result1 := m.Multiply(v)
	// v*M - vector as row vector multiplying matrix
	result2 := v.MultiplyMatrix(m)

	// These should be different operations
	x1, y1, z1 := result1.Values()
	x2, y2, z2 := result2.Values()

	// Matrix * Vector = (2*1 + 3*2 + 4*3, 2*4 + 3*5 + 4*6, 2*7 + 3*8 + 4*9) = (20, 47, 74)
	// Vector * Matrix = see previous test = (42, 51, 60)

	if almostEqual(x1, x2) && almostEqual(y1, y2) && almostEqual(z1, z2) {
		t.Errorf("Matrix*Vector and Vector*Matrix should be different but got the same result")
	}

	// Check specific expected values
	if !almostEqual(x1, 20) || !almostEqual(y1, 47) || !almostEqual(z1, 74) {
		t.Errorf("Matrix*Vector calculation incorrect: expected (20,47,74), got (%f,%f,%f)", x1, y1, z1)
	}

	if !almostEqual(x2, 42) || !almostEqual(y2, 51) || !almostEqual(z2, 60) {
		t.Errorf("Vector*Matrix calculation incorrect: expected (42,51,60), got (%f,%f,%f)", x2, y2, z2)
	}
}

// Helper function to compare float values with small epsilon
func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}
