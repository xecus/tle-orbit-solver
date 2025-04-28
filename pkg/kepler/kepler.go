package kepler

import "math"

// KeplerEquation returns a function that evaluates the Kepler equation for given values of e and M
func KeplerEquation(e, M float64) func(Ebefore float64) float64 {
	// Ebefore: n回目の離心近点離角
	// M: 平均近点離角
	// e: 離心率
	// Eafter: n+1回目の離心近点離角
	return func(Ebefore float64) float64 {
		FE := Ebefore - e*math.Sin(Ebefore) - M
		Eafter := Ebefore - FE/(1-e*math.Cos(Ebefore))
		return Eafter
	}
}

// NewtonRaphson solves the Kepler equation using Newton-Raphson method
func NewtonRaphson(e, before, a float64) (float64, float64) {
	// e: 離心率
	// before: 平均近点離角（初期値）
	// a: allowable error 許容誤差
	equation := KeplerEquation(e, before)
	var after float64
	err := 100.0 // 許容誤差の初期化（0だとforが回らないため）

	// Fixed number of iterations for consistency
	for i := 0; i < 10; i++ {
		after = equation(before)
		err = math.Abs(after - before)
		before = after
	}
	
	return after, err
}