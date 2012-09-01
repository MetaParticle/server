package structures

import (
	"math"
)

const (
	MAXTHETA = math.Pi
	MAXPHI = 2*math.Pi
)

type SphVec3 struct {
	r, phi, theta float64
}

func (c *SphVec3) Set(r, phi, theta float64) *SphVec3 {
	c.SetR(r)
	c.SetPhi(phi)
	c.SetTheta(theta)
	return c
}


/**
 * Sets the radius part of the coordinates
 */
func (c *SphVec3) SetR(n float64) *SphVec3 {
	c.r = n
	return c
}

/**
 * Sets the phi part of the coordinates
 */
func (c *SphVec3) SetPhi(n float64) *SphVec3 {
	// TODO Make more fluid transitions
	
	if n < 0 {
		return c.SetPhi(n+MAXPHI)
	}
	
	if n > MAXPHI {
		return c.SetPhi(n-MAXPHI)
	} else {
		c.phi = n
	}
	return c
}

/**
 * Sets the theta part of the coordinates
 */
func (c *SphVec3) SetTheta(n float64) *SphVec3 {
	// TODO Make more fluid transitions
	
	if n < 0 {
		return c.SetTheta(n+MAXTHETA)
	}
	
	if n > MAXTHETA {
		return c.SetTheta(n-MAXTHETA)
	} else {
		c.theta = n
	}
	return c
}

func (c *SphVec3) LengthSq() float64 {
	return c.r*c.r
}

func (c *SphVec3) Length() float64 {
	return c.r
}

// TODO ADD

// TODO SUBTRACT

func (c *SphVec3) Multiply(n float64) *SphVec3 {
	c.r *= n
	return c
}

func (c *SphVec3) Divide(n float64) *SphVec3 {
	if n == 0 {
		c.r = 0
	} else {
		c.r /= n
	}
	return c
}

/*
 * A SphVec3.SetLength(n) method mirroring the Vec3 one.
 * For a similar interface.
 */
func (c *SphVec3) SetLength(n float64) *SphVec3 {
	return c.SetR(n)
}

// MOVED Vec3ToSpherical(Vec3) -> Vec3.ToSpherical()

func (c *SphVec3) ToVec3() (v *Vec3) {
	v.x = c.r*math.Cos(c.phi)*math.Sin(c.theta)
	v.y = c.r*math.Sin(c.phi)*math.Sin(c.theta)
	v.z = c.r*math.Sin(c.theta)
	return
}