package physics

import (
	"math"
)

// calculateSpeed calculates the speed based on the given power
// stolen from https://www.endurance-data.com/en/power-to-speed-calculator-pro/
func CalculateSpeed(power float64) float64 {
	// Default values for other parameters
	weight := 82.0                // Weight in kg
	cda := 0.270                  // Drag coefficient times frontal area in m^2
	crr := 0.00366                // Coefficient of rolling resistance
	drivetrainEfficiency := 0.965 // Drivetrain efficiency as a fraction
	temperature := 20.0           // Temperature in Celsius
	gravity := 9.81               // Gravitational acceleration in m/s^2

	// Constants
	R := 8.31432                 // Universal gas constant
	M := 0.0289644               // Molar mass of air
	pressureSeaLevel := 101325.0 // Pressure at sea level in Pascals

	// Calculate air density
	tKelvin := temperature + 273.15
	airDensity := (M * pressureSeaLevel) / (R * tKelvin)

	// Effective power
	effectivePower := power * drivetrainEfficiency

	// Coefficients for the cubic equation
	a := 0.5 * cda * airDensity
	b := 0.0 // No quadratic term
	c := crr * weight * gravity
	d := -effectivePower

	// Solve the cubic equation using Cardano's method
	roots := solveCubic(a, b, c, d)

	// Find the positive real root
	var velocity float64
	for _, root := range roots {
		if root > 0 {
			velocity = root
			break
		}
	}

	// Convert velocity to km/h
	speedKmh := velocity * 3.6
	return speedKmh
}

// solveCubic solves the cubic equation ax^3 + bx^2 + cx + d = 0 and returns the real roots
func solveCubic(a, b, c, d float64) []float64 {
	// Convert to depressed cubic form t^3 + pt + q = 0
	p := (3*a*c - b*b) / (3 * a * a)
	q := (2*b*b*b - 9*a*b*c + 27*a*a*d) / (27 * a * a * a)

	var roots []float64
	if math.Abs(p) < 1e-8 { // Handle edge case where p ≈ 0
		root := math.Cbrt(-q)
		roots = append(roots, root)
	} else if math.Abs(q) < 1e-8 { // Handle edge case where q ≈ 0
		if p > 0 {
			roots = append(roots, 0, math.Sqrt(-p), -math.Sqrt(-p))
		} else {
			roots = append(roots, 0)
		}
	} else {
		D := q*q/4 + p*p*p/27
		if math.Abs(D) < 1e-8 { // One real root
			root := -1.5 * q / p
			roots = append(roots, root)
		} else if D > 0 { // One real root
			u := math.Cbrt(-q/2 - math.Sqrt(D))
			root := u - p/(3*u)
			roots = append(roots, root)
		} else { // Three real roots
			u := 2 * math.Sqrt(-p/3)
			theta := math.Acos(3*q/p/u) / 3
			roots = append(roots, u*math.Cos(theta))
			roots = append(roots, u*math.Cos(theta-2*math.Pi/3))
			roots = append(roots, u*math.Cos(theta+2*math.Pi/3))
		}
	}

	// Normalize roots by subtracting offset
	for i := range roots {
		roots[i] -= b / (3 * a)
	}
	return roots
}
