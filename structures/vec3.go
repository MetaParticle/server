package structures

import (
	"math"
)

type Vector interface {
	Set(a, b,c float64) *Vector
	SetLength(n float64) *Vector
	LengthSq() float64
	Length() float64
	Add(v Vector) *Vector
	Subtract(v Vector) *Vector
	Multiply(v Vector) *Vector
	Divide(v Vector) *Vector
	Copy() *Vector
}

type Vec3 struct {
	x, y, z float64
}

func (v *Vec3) Set(x, y, z float64) *Vec3 {
	v.x, v.y, v.z = x, y, z
	return v
}


/**
 * Sets the x part of the vector
 */
func (v *Vec3) SetX(n float64) *Vec3 {
	v.x = n
	return v
}

/**
 * Sets the y part of the vector
 */
func (v *Vec3) SetY(n float64) *Vec3 {
	v.y = n
	return v
}

/**
 * Sets the z part of the vector
 */
func (v *Vec3) SetZ(n float64) *Vec3 {
	v.z = n
	return v
}

func (v *Vec3) LengthSq() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

func (v *Vec3) Length() float64 {
	return math.Sqrt(v.LengthSq())
}

/**
 * Adds the vector u to vector v and returns vector v.
 */
func (v *Vec3) Add(u *Vec3) *Vec3 {
	v.x += u.x
	v.y += u.y
	v.z += u.z
	return v
}

/**
 * Subtracts vector u from vector v and returns vector v.
 */
func (v *Vec3) Subtract(u *Vec3) *Vec3 {
	clone := u.Copy()
	return v.Add(clone.Multiply(-1))
}

func (v *Vec3) Multiply(n float64) *Vec3 {
	v.x *= n
	v.y *= n
	v.z *= n
	return v
}

func (v *Vec3) Divide(n float64) *Vec3 {
	if n != 0 {
		v.x /= n
		v.y /= n
		v.z /= n
	} else {
		v.x, v.y, v.z = 0, 0, 0
	}
	return v
}

func (v *Vec3) Normalize() *Vec3 {
	return v.Divide(v.Length())
}

func (v *Vec3) SetLength(n float64) *Vec3 {
	return v.Normalize().Multiply(n)
}

func (v *Vec3) Copy() (w *Vec3) {
	w.x = v.x
	w.y = v.y
	w.z = v.z
	return
}

func (v *Vec3) ToSpherical() (c *SphVec3) {
	// from Wikipedia article
	c.r = v.Length()
	c.phi = math.Atan2(v.y, v.x)
	c.theta = math.Atan(v.z / c.r)
	return
}
