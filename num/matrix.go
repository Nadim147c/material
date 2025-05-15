package num

// Matrix3 defines a 3x3 matrix of float64 values
type Matrix3 [3]Vector3

func NewMatrix3(x1, y1, z1, x2, y2, z2, x3, y3, z3 float64) Matrix3 {
	return Matrix3{
		{x1, y1, z1},
		{x2, y2, z2},
		{x3, y3, z3},
	}
}

// MultiplyXYZ takes x, y, z and creates 3D Vector3. Then Multiply m with the
// newly created Vector3. Returns the resulting vector
func (m Matrix3) MultiplyXYZ(x, y, z float64) Vector3 {
	return m.Multiply(NewVector3(x, y, z))
}

// Multiply applies the matrix to a 3D vector and returns the resulting vector
func (m Matrix3) Multiply(v Vector3) Vector3 {
	var result Vector3
	for i := range 3 {
		for j := range 3 {
			result[i] += m[i][j] * v[j]
		}
	}
	return result
}

// Vector3 defines a 3D vector
type Vector3 [3]float64

// NewVector3 create new 3D vector: Vector3
func NewVector3(x, y, z float64) Vector3 {
	return Vector3{x, y, z}
}

func (v Vector3) MultiplyMatrix(m Matrix3) Vector3 {
	var result Vector3
	for j := range 3 {
		for i := range 3 {
			result[j] += v[i] * m[i][j]
		}
	}
	return result
}

func (v Vector3) MultiplyScalar(s float64) Vector3 {
	var result Vector3
	for i := range 3 {
		result[i] = v[i] * s
	}
	return v
}

func (v Vector3) Add(vec Vector3) Vector3 {
	var result Vector3
	for i := range 3 {
		result[i] = v[i] + vec[i]
	}
	return v
}

func (v Vector3) Values() (float64, float64, float64) {
	return v[0], v[1], v[2]
}
