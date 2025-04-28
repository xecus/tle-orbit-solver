package orbital

import (
	"fmt"
	"math"
	"time"

	"gonum.org/v1/gonum/mat"

	"starlink/pkg/kepler"
	"starlink/pkg/model"
	"starlink/pkg/util"
)

// CalculateSatelliteLocation calculates the position of a satellite based on its orbital elements
func CalculateSatelliteLocation(sat *model.TleOrbitalElement, targetTime time.Time) *model.SatLocation {
	// Extract orbital parameters
	m0 := sat.MeanAnomaly
	m1 := sat.MeanMotion
	m2 := sat.MeanMotionDot
	ecc := sat.Eccentricity
	angleOmegaA0 := sat.ArgumentOfPerigee
	angleI0 := sat.OrbitalInclination
	angleOmegaB0 := sat.Raan
	epocTimeYear := sat.EtYear
	epocTimeDay := sat.EtDay

	// Convert to UTC
	targetTime = targetTime.UTC()
	fmt.Println("targetTime=", targetTime)

	// Calculate time difference from epoch
	t_diff := calculateTimeDifference(targetTime, epocTimeYear, epocTimeDay)
	fmt.Println("t_diff=", t_diff)

	// Calculate orbital semi-axes
	a, b := calculateOrbitalSemiAxes(m1)
	fmt.Println("a [km] =", a)
	fmt.Println("b [km] =", b)
	fmt.Println("ecc =", ecc)

	// Calculate mean anomaly
	fracM_Radian := calculateMeanAnomaly(m0, m1, m2, t_diff)
	fmt.Println("fracM (Radian) =", fracM_Radian)

	// Solve Kepler's equation for eccentric anomaly
	eccentricAnomaly, valerr := kepler.NewtonRaphson(ecc, fracM_Radian, 0.00001)
	fmt.Println("eccentricAnomaly=", eccentricAnomaly)
	fmt.Println("valerr=", valerr)

	// Calculate position in orbital plane
	u, v := calculatePositionInOrbitalPlane(a, ecc, eccentricAnomaly)
	fmt.Println("u (km)=", u)
	fmt.Println("v (km)=", v)

	// Apply perturbation corrections
	angleOmegaA_Degree, angleOmegaB_Degree := calculatePerturbationCorrection(
		angleOmegaA0, angleOmegaB0, angleI0, a, t_diff)
	fmt.Println("angleOmegaA (Degree)=", angleOmegaA_Degree)
	fmt.Println("angleOmegaB (Degree)=", angleOmegaB_Degree)

	// Convert angles to radians
	angleOmegaA_Rad := util.Deg2Rad(angleOmegaA_Degree)
	angleOmegaB_Rad := util.Deg2Rad(angleOmegaB_Degree)
	angleI0_Rad := util.Deg2Rad(angleI0)

	// Transform to equatorial coordinate system
	x, y, z := transformToEquatorial(u, v, angleOmegaA_Rad, angleOmegaB_Rad, angleI0_Rad)
	fmt.Println("x (km) =", x)
	fmt.Println("y (km) =", y)
	fmt.Println("z (km) =", z)

	// Transform to Earth-fixed coordinate system
	largeX, largeY, largeZ := transformToEarthFixed(x, y, z, targetTime)
	fmt.Println("LargeX (km) =", largeX)
	fmt.Println("LargeY (km) =", largeY)
	fmt.Println("LargeZ (km) =", largeZ)

	// Calculate latitude and longitude
	lat, lng := calculateLatLong(largeX, largeY, largeZ)
	fmt.Println("Fai (Degree) =", lat)
	fmt.Println("Lambda (Degree) =", lng)

	// Calculate altitude
	alt := calculateAltitude(largeX, largeY, largeZ)
	fmt.Println("Alt (km) =", alt)

	return &model.SatLocation{
		X:   largeX,
		Y:   largeY,
		Z:   largeZ,
		Lat: lat,
		Lng: lng,
		Alt: alt,
	}
}

// CalculateVelocity calculates the satellite's velocity based on position delta
func CalculateVelocity(satLoc1, satLoc2 *model.SatLocation) float64 {
	diffX := satLoc2.X - satLoc1.X
	diffY := satLoc2.Y - satLoc1.Y
	diffZ := satLoc2.Z - satLoc1.Z
	
	fmt.Printf("DiffX[km] = %f\n", diffX)
	fmt.Printf("DiffY[km]= %f\n", diffY)
	fmt.Printf("DiffZ[km] = %f\n", diffZ)
	
	velocity := math.Sqrt(diffX*diffX + diffY*diffY + diffZ*diffZ)
	fmt.Printf("V[km/s] = %f\n", velocity)
	
	return velocity
}

// calculateTimeDifference calculates the time difference between target time and epoch
func calculateTimeDifference(targetTime time.Time, epocTimeYear int, epocTimeDay float64) float64 {
	var t_diff1, t_diff2 float64

	// Target time days from year start
	tmp := time.Date(targetTime.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	t_diff1 = float64(targetTime.Unix()-tmp.Unix())/86400.0 + 1.0
	fmt.Println("t_diff1=", t_diff1)

	// Epoch time (EpocTime)
	tmp = time.Date(epocTimeYear, time.January, 1, 0, 0, 0, 0, time.UTC)
	tmp = tmp.Add(time.Duration(float64(time.Second) * 86400.0 * epocTimeDay))
	tmp2 := time.Date(epocTimeYear, time.January, 1, 0, 0, 0, 0, time.UTC)
	t_diff2 = float64(tmp.Unix()-tmp2.Unix()) / 86400.0
	fmt.Println("t_diff2=", t_diff2)

	t_diff := t_diff1 - t_diff2
	if t_diff < 0.0 {
		panic("TargetTime < ET")
	}

	return t_diff
}

// calculateOrbitalSemiAxes calculates semi-major and semi-minor axes of orbit
func calculateOrbitalSemiAxes(meanMotion float64) (float64, float64) {
	// Calculate semi-major axis (Kepler's third law)
	a := math.Cbrt((2.975537 * math.Pow(10, 15)) / (4 * math.Pi * math.Pi * meanMotion * meanMotion))
	b := math.Sqrt(a*a - a*a*0.0001679*0.0001679)

	return a, b
}

// calculateMeanAnomaly calculates the mean anomaly at target time
func calculateMeanAnomaly(m0, m1, m2, t_diff float64) float64 {
	// M
	m := (m0 / 360.0) + m1*t_diff + 0.5*m2*t_diff*t_diff
	fmt.Println("m (Rev) =", m)

	// Convert from revs to radians
	fracM := m - float64(int64(m))
	fracM_Degree := fracM * 360.0
	fracM_Radian := fracM * (2 * math.Pi)

	fmt.Println("fracM (Degree) =", fracM_Degree)

	return fracM_Radian
}

// calculatePositionInOrbitalPlane calculates position in the orbital plane
func calculatePositionInOrbitalPlane(a, ecc, eccentricAnomaly float64) (float64, float64) {
	// Calculate position in orbital plane
	u := a*math.Cos(eccentricAnomaly) - a*ecc
	v := a * math.Sqrt(1-ecc*ecc) * math.Sin(eccentricAnomaly)

	return u, v
}

// calculatePerturbationCorrection applies perturbation corrections to orbital elements
func calculatePerturbationCorrection(angleOmegaA0, angleOmegaB0, angleI0, a, t_diff float64) (float64, float64) {
	// Apply perturbations
	angleOmegaA_Degree := angleOmegaA0 +
		(180*0.174*(2-2.5*math.Pow(math.Sin(angleI0*math.Pi/180.0), 2)))/(math.Pi*math.Pow(a/util.EarthRadius, 3.5))*t_diff
	angleOmegaB_Degree := angleOmegaB0 -
		(180*0.174*math.Cos(angleI0*math.Pi/180.0))/(math.Pi*math.Pow(a/util.EarthRadius, 3.5))*t_diff

	return angleOmegaA_Degree, angleOmegaB_Degree
}

// transformToEquatorial transforms coordinates from orbital plane to equatorial reference frame
func transformToEquatorial(u, v float64, angleOmegaA_Rad, angleOmegaB_Rad, angleI0_Rad float64) (float64, float64, float64) {
	// Transform from orbital plane to equatorial reference frame
	elemA := []float64{
		math.Cos(angleOmegaB_Rad), -math.Sin(angleOmegaB_Rad), 0,
		math.Sin(angleOmegaB_Rad), math.Cos(angleOmegaB_Rad), 0,
		0, 0, 1}
	elemB := []float64{
		1, 0, 0,
		0, math.Cos(angleI0_Rad), -math.Sin(angleI0_Rad),
		0, math.Sin(angleI0_Rad), math.Cos(angleI0_Rad)}
	elemC := []float64{
		math.Cos(angleOmegaA_Rad), -math.Sin(angleOmegaA_Rad), 0,
		math.Sin(angleOmegaA_Rad), math.Cos(angleOmegaA_Rad), 0,
		0, 0, 1}
	elemD := []float64{u, v, 0}

	elemTempA := make([]float64, 9)
	elemTempC := make([]float64, 3)
	matTempA := mat.NewDense(3, 3, elemTempA)
	matTempC := mat.NewDense(3, 1, elemTempC)
	matA := mat.NewDense(3, 3, elemA)
	matB := mat.NewDense(3, 3, elemB)
	matC := mat.NewDense(3, 3, elemC)
	matD := mat.NewDense(3, 1, elemD)

	matTempA.Mul(matA, matB)
	matTempA.Mul(matTempA, matC)
	matTempC.Mul(matTempA, matD)

	x := matTempC.At(0, 0)
	y := matTempC.At(1, 0)
	z := matTempC.At(2, 0)

	return x, y, z
}

// transformToEarthFixed transforms from equatorial to Earth-fixed reference frame
func transformToEarthFixed(x, y, z float64, targetTime time.Time) (float64, float64, float64) {
	// Transform from equatorial to Earth-fixed reference frame
	tmp := time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC)
	sheta0 := 0.27644444444
	t_diff1 := float64(targetTime.Unix()-tmp.Unix()) / 86400.0
	sheta := sheta0 + 1.002737909*t_diff1

	// Convert fractional part to angle
	shetaG_Deg := (sheta - float64(int64(sheta))) * 360.0
	shetaG_Rad := (sheta - float64(int64(sheta))) * 2.0 * math.Pi

	fmt.Println("ShetaG (Degree) = ", shetaG_Deg)

	// Apply rotation matrix around z-axis in reverse direction
	elemA := []float64{
		math.Cos(-shetaG_Rad), -math.Sin(-shetaG_Rad), 0,
		math.Sin(-shetaG_Rad), math.Cos(-shetaG_Rad), 0,
		0, 0, 1}
	elemB := []float64{x, y, z}
	elemC := make([]float64, 3)

	matA := mat.NewDense(3, 3, elemA)
	matB := mat.NewDense(3, 1, elemB)
	matC := mat.NewDense(3, 1, elemC)

	matC.Mul(matA, matB)

	largeX := matC.At(0, 0)
	largeY := matC.At(1, 0)
	largeZ := matC.At(2, 0)

	return largeX, largeY, largeZ
}

// calculateLatLong converts Cartesian coordinates to latitude and longitude
func calculateLatLong(x, y, z float64) (float64, float64) {
	// Calculate latitude and longitude (assuming spherical Earth)
	fai := math.Asin(z / math.Sqrt(x*x+y*y+z*z))
	lambda := math.Atan2(y, x)

	// Convert from radians to degrees
	fai_Degree := util.Rad2Deg(fai)
	lambda_Degree := util.Rad2Deg(lambda)

	return fai_Degree, lambda_Degree
}

// calculateAltitude calculates the altitude above Earth's surface
func calculateAltitude(x, y, z float64) float64 {
	// Calculate altitude (assuming spherical Earth)
	return math.Sqrt(x*x+y*y+z*z) - util.EarthRadius
}