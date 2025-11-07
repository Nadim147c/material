package num

import "fmt"

// Matrix3 defines a 3x3 matrix of float64 values.
type Matrix3 [3]Vector3

// NewMatrix3 creates a 3x3 Matrix3.
func NewMatrix3(x1, y1, z1, x2, y2, z2, x3, y3, z3 float64) Matrix3 {
	return Matrix3{
		{x1, y1, z1},
		{x2, y2, z2},
		{x3, y3, z3},
	}
}

// Multiply applies the matrix to a 3D vector and returns the resulting vector.
func (m Matrix3) Multiply(v Vector3) Vector3 {
	var result Vector3
	for row := range 3 {
		for col := range 3 {
			result[row] += m[row][col] * v[col]
		}
	}
	return result
}

// Transpose transposes the Matrix3.
func (m Matrix3) Transpose() Matrix3 {
	var result Matrix3
	for row := range 3 {
		for col := range 3 {
			result[col][row] = m[row][col]
		}
	}
	return result
}

// Inverse inverses the M Matrix3. Returns false if m Matrix3 is not inversable.
func (m Matrix3) Inverse() (Matrix3, bool) {
	a, b, c := m[0][0], m[0][1], m[0][2]
	d, e, f := m[1][0], m[1][1], m[1][2]
	g, h, i := m[2][0], m[2][1], m[2][2]

	// Compute the determinant
	det := a*(e*i-f*h) - b*(d*i-f*g) + c*(d*h-e*g)
	if det == 0 {
		return Matrix3{}, false // Matrix is not invertible
	}
	invDet := 1.0 / det

	var inv Matrix3
	inv[0][0] = (e*i - f*h) * invDet
	inv[0][1] = -(b*i - c*h) * invDet
	inv[0][2] = (b*f - c*e) * invDet
	inv[1][0] = -(d*i - f*g) * invDet
	inv[1][1] = (a*i - c*g) * invDet
	inv[1][2] = -(a*f - c*d) * invDet
	inv[2][0] = (d*h - e*g) * invDet
	inv[2][1] = -(a*h - b*g) * invDet
	inv[2][2] = (a*e - b*d) * invDet

	return inv, true
}

// String implements Stringer method.
func (m Matrix3) String() string {
	return fmt.Sprintf(
		"[\n\t%.10f,%.10f,%.10f,\n\t%.10f,%.10f,%.10f,\n\t%.10f,%.10f,%.10f,\n]",
		m[0][0],
		m[0][1],
		m[0][2],
		m[1][0],
		m[1][1],
		m[1][2],
		m[2][0],
		m[2][1],
		m[2][2],
	)
}

// Vector3 defines a 3D vector.
type Vector3 [3]float64

// NewVector3 create new 3D vector: Vector3.
func NewVector3(x, y, z float64) Vector3 {
	return Vector3{x, y, z}
}

// VectorLike is any data that has 3 float64 values. Usually a color model.
type VectorLike interface {
	Values() (float64, float64, float64)
}

// NewVector create new 3D vector: Vector3.
func NewVector(v VectorLike) Vector3 {
	return NewVector3(v.Values())
}

// Transform returns a new Vector3 where each component of v has been
// transformed by the given func f. The func f is applied to each of v's
// components independently, producing a new Vector3 without modifying the
// original.
func (v Vector3) Transform(f func(float64) float64) Vector3 {
	var result Vector3
	for i := range 3 {
		result[i] = f(v[i])
	}
	return result
}

// MultiplyScalar multiply all values using a scalar float64 and retuns a new
// Vector3.
func (v Vector3) MultiplyScalar(s float64) Vector3 {
	var result Vector3
	for i := range 3 {
		result[i] = v[i] * s
	}
	return result
}

// Add adds values of v Vector3 with given Vector3 and returns new Vector3.
func (v Vector3) Add(vec Vector3) Vector3 {
	var result Vector3
	for i := range 3 {
		result[i] = v[i] + vec[i]
	}
	return result
}

// Values retuns all values of underlying Vector3.
func (v Vector3) Values() (float64, float64, float64) {
	return v[0], v[1], v[2]
}
