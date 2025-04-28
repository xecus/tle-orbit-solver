package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/gonum/mat"
)

const EarthRadius = 6356.752 //地球の半径

type SatLocation struct {
	X   float64
	Y   float64
	Z   float64
	Lat float64
	Lng float64
	Alt float64
}

type TleOrbitalElement struct {
	meanAnomaly        float64 // M0 平均近点角 [Degree]
	meanMotion         float64 // M1 平均運動: [Rev/Day]
	MeanMotionDot      float64 // M2 平均運動変化係数: [Rev/Day2]
	eccentricity       float64 // 離心率 [-]
	etYear             int     // 元期 Epoctime [Year]
	etDay              float64 // 元期 EpocTime [Day]
	orbitalInclination float64 // 軌道傾斜角 [Degree]
	raan               float64 // 昇交点赤経: RAAN [Degree]
	argumentOfPerigee  float64 // 近地点引数 [Degree]
}

// ケプラー方程式
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

// ケプラー方程式をニュートン-ラフソン法で解く
func NewtonRaphson(e, before, a float64) (float64, float64) {
	// e: 離心率
	// M: 平均近点離角
	// a: allowable error 許容誤差
	equation := KeplerEquation(e, before)
	var after float64
	err := 100.0 // 許容誤差の初期化（0だとforが回らないため）
	count := 0
	//for err > a {
	for i := 0; i < 10; i++ {
		after = equation(before)
		err = math.Abs(after - before)
		//fmt.Println("誤差: ", err)
		before = after
		count = count + 1
		//fmt.Println("繰り返し：", count, "回目")
		//fmt.Println("----")
	}
	return after, err
}

func main() {
	testSat := parseTle()

	targetTime := time.Now()
	satLoc1 := calculateSatelliteLocation(testSat, targetTime)
	satLoc2 := calculateSatelliteLocation(testSat, targetTime.Add(time.Second))

	calculateVelocity(satLoc1, satLoc2)
}

func calculateVelocity(satLoc1, satLoc2 *SatLocation) {
	diffX := satLoc2.X - satLoc1.X
	diffY := satLoc2.Y - satLoc1.Y
	diffZ := satLoc2.Z - satLoc1.Z
	fmt.Printf("DiffX[km] = %f\n", diffX)
	fmt.Printf("DiffY[km]= %f\n", diffY)
	fmt.Printf("DiffZ[km] = %f\n", diffZ)
	fmt.Printf("V[km/s] = %f\n", math.Sqrt(diffX*diffX+diffY*diffY+diffZ*diffZ))
}

func parseTle() *TleOrbitalElement {
	// STARLINK-1008のTLEデータ (2025/04/28 15:08時点)
	// https://celestrak.org/NORAD/elements/gp.php?GROUP=starlink&FORMAT=tle
	str1 := "1 44714U 19074B   25117.42924319 -.00001157  00000+0 -58773-4 0  9990"
	str2 := "2 44714  53.0517 166.3609 0001116  99.1558 260.9557 15.06400606301084"
	a := strings.Split(str1, " ")
	b := strings.Split(str2, " ")

	// TLEから各パラメータを抽出
	// aの長さを確認
	fmt.Printf("TLE line 1 split count: %d\n", len(a))
	for i, v := range a {
		fmt.Printf("a[%d] = %s\n", i, v)
	}

	satelliteNumber := a[1]
	internationalDesignator := a[2]
	etYear, _ := strconv.Atoi(fmt.Sprintf("20%s", a[5][0:2]))
	etDay, _ := strconv.ParseFloat(a[5][2:], 64)
	firstTimeDerivativeOfTheMeanMotion, _ := strconv.ParseFloat(fmt.Sprintf("0%s", a[7]), 64)
	secondTimeDerivativeOfTheMeanMotion := a[9]
	bstarDragTerm := a[11]
	elementnum := a[12][0:3]
	checksum1 := a[12][3:]
	// bの長さを確認
	fmt.Printf("TLE line 2 split count: %d\n", len(b))
	for i, v := range b {
		fmt.Printf("b[%d] = %s\n", i, v)
	}

	// TLEのフォーマットを考慮して固定位置で分割
	// TLE行2のフォーマットを考慮して正確に読み取る
	line2 := str2
	orbitalInclination, _ := strconv.ParseFloat(strings.TrimSpace(line2[8:16]), 64)
	rightAscensionOfAscendingNode, _ := strconv.ParseFloat(strings.TrimSpace(line2[17:25]), 64)
	eccentricity, _ := strconv.ParseFloat("0."+strings.TrimSpace(line2[26:33]), 64)
	argumentOfPerigee, _ := strconv.ParseFloat(strings.TrimSpace(line2[34:42]), 64)
	meanAnomaly, _ := strconv.ParseFloat(strings.TrimSpace(line2[43:51]), 64)

	// 平均運動
	meanMotionStr := strings.TrimSpace(line2[52:63])
	meanMotion, _ := strconv.ParseFloat(meanMotionStr, 64)

	// 周回数と検査数字
	numberOfLaps := strings.TrimSpace(line2[63:68])
	checksum2 := line2[68:69]

	// TLEパラメータの表示
	printTleParameters(satelliteNumber, internationalDesignator, etYear, etDay,
		firstTimeDerivativeOfTheMeanMotion, secondTimeDerivativeOfTheMeanMotion,
		bstarDragTerm, elementnum, checksum1, orbitalInclination,
		rightAscensionOfAscendingNode, eccentricity, argumentOfPerigee,
		meanAnomaly, meanMotion, numberOfLaps, checksum2)

	return &TleOrbitalElement{
		meanAnomaly:        meanAnomaly,
		meanMotion:         meanMotion,
		MeanMotionDot:      firstTimeDerivativeOfTheMeanMotion,
		eccentricity:       eccentricity,
		etYear:             etYear,
		etDay:              etDay,
		orbitalInclination: orbitalInclination,
		raan:               rightAscensionOfAscendingNode,
		argumentOfPerigee:  argumentOfPerigee,
	}
}

func printTleParameters(satelliteNumber, internationalDesignator string, etYear int, etDay float64,
	firstTimeDerivativeOfTheMeanMotion float64, secondTimeDerivativeOfTheMeanMotion string,
	bstarDragTerm, elementnum, checksum1 string, orbitalInclination float64,
	rightAscensionOfAscendingNode, eccentricity, argumentOfPerigee,
	meanAnomaly, meanMotion float64, numberOfLaps, checksum2 string) {
	fmt.Println("--------------------------------------------------")
	fmt.Printf("satelliteNumber=%s\n", satelliteNumber)
	fmt.Printf("internationalDesignator=%s\n", internationalDesignator)
	fmt.Printf("etYear=%d\n", etYear)
	fmt.Printf("etDay=%f\n", etDay)
	fmt.Printf("firstTimeDerivativeOfTheMeanMotion=%f\n",
		firstTimeDerivativeOfTheMeanMotion)
	fmt.Printf("secondTimeDerivativeOfTheMeanMotion=%s\n",
		secondTimeDerivativeOfTheMeanMotion)
	fmt.Printf("bstarDragTerm=%s\n", bstarDragTerm)
	fmt.Printf("elementnum=%s\n", elementnum)
	fmt.Printf("checksum1=%s\n", checksum1)
	fmt.Printf("orbitalInclination=%f [Degree]\n", orbitalInclination)
	fmt.Printf("RAAN=%f [Degree]\n", rightAscensionOfAscendingNode)
	fmt.Printf("eccentricity=%f [-]\n", eccentricity)
	fmt.Printf("argumentOfPerigee=%f [Degree]\n", argumentOfPerigee)
	fmt.Printf("meanAnomaly=%f [Degree]\n", meanAnomaly)
	fmt.Printf("meanMotion=%f [Rev/Day]\n", meanMotion)
	fmt.Printf("numberOfLaps=%s [-]\n", numberOfLaps)
	fmt.Printf("checksum2=%s [-]\n", checksum2)
	fmt.Println("--------------------------------------------------")
}

func calculateSatelliteLocation(sat *TleOrbitalElement, targetTime time.Time) *SatLocation {
	// TLEパラメータを変数に展開
	m0 := sat.meanAnomaly
	m1 := sat.meanMotion
	m2 := sat.MeanMotionDot
	ecc := sat.eccentricity
	angleOmegaA0 := sat.argumentOfPerigee
	angleI0 := sat.orbitalInclination
	angleOmegaB0 := sat.raan
	epocTimeYear := sat.etYear
	epocTimeDay := sat.etDay

	// UTC時間に変換
	targetTime = targetTime.UTC()
	fmt.Println("targetTime=", targetTime)

	// 元期からの経過時間計算
	t_diff := calculateTimeDifference(targetTime, epocTimeYear, epocTimeDay)
	fmt.Println("t_diff=", t_diff)

	// 軌道長半径と短半径の計算
	a, b := calculateOrbitalSemiAxes(m1)
	fmt.Println("a [km] =", a)
	fmt.Println("b [km] =", b)
	fmt.Println("ecc =", ecc)

	// 平均近点角の計算
	fracM_Radian := calculateMeanAnomaly(m0, m1, m2, t_diff)
	fmt.Println("fracM (Radian) =", fracM_Radian)

	// 離心近点角の計算（ニュートンラフソン法）
	eccentricAnomaly, valerr := NewtonRaphson(ecc, fracM_Radian, 0.00001)
	fmt.Println("eccentricAnomaly=", eccentricAnomaly)
	fmt.Println("valerr=", valerr)

	// 軌道平面上の位置計算
	u, v := calculatePositionInOrbitalPlane(a, ecc, eccentricAnomaly)
	fmt.Println("u (km)=", u)
	fmt.Println("v (km)=", v)

	// 摂動補正
	angleOmegaA_Degree, angleOmegaB_Degree := calculatePerturbationCorrection(
		angleOmegaA0, angleOmegaB0, angleI0, a, t_diff)
	fmt.Println("angleOmegaA (Degree)=", angleOmegaA_Degree)
	fmt.Println("angleOmegaB (Degree)=", angleOmegaB_Degree)

	// 角度をラジアンに変換
	angleOmegaA_Rad := deg2rad(angleOmegaA_Degree)
	angleOmegaB_Rad := deg2rad(angleOmegaB_Degree)
	angleI0_Rad := deg2rad(angleI0)

	// 軌道平面から春分点座標系への変換
	x, y, z := transformToEquatorial(u, v, angleOmegaA_Rad, angleOmegaB_Rad, angleI0_Rad)
	fmt.Println("x (km) =", x)
	fmt.Println("y (km) =", y)
	fmt.Println("z (km) =", z)

	// 春分点座標系から地球固定座標系への変換
	largeX, largeY, largeZ := transformToEarthFixed(x, y, z, targetTime)
	fmt.Println("LargeX (km) =", largeX)
	fmt.Println("LargeY (km) =", largeY)
	fmt.Println("LargeZ (km) =", largeZ)

	// 緯度経度への変換
	lat, lng := calculateLatLong(largeX, largeY, largeZ)
	fmt.Println("Fai (Degree) =", lat)
	fmt.Println("Lambda (Degree) =", lng)

	// 高度の計算
	alt := calculateAltitude(largeX, largeY, largeZ)
	fmt.Println("Alt (km) =", alt)

	return &SatLocation{
		X:   largeX,
		Y:   largeY,
		Z:   largeZ,
		Lat: lat,
		Lng: lng,
		Alt: alt,
	}
}

func calculateTimeDifference(targetTime time.Time, epocTimeYear int, epocTimeDay float64) float64 {
	var t_diff1, t_diff2 float64

	// 予測したい時間の日数経過
	tmp := time.Date(targetTime.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	t_diff1 = float64(targetTime.Unix()-tmp.Unix())/86400.0 + 1.0
	fmt.Println("t_diff1=", t_diff1)

	// 元期(EpocTime)
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

func calculateOrbitalSemiAxes(meanMotion float64) (float64, float64) {
	// 軌道長半径aの算出 (ケプラーの第三法則)
	a := math.Cbrt((2.975537 * math.Pow(10, 15)) / (4 * math.Pi * math.Pi * meanMotion * meanMotion))
	b := math.Sqrt(a*a - a*a*0.0001679*0.0001679)

	return a, b
}

func calculateMeanAnomaly(m0, m1, m2, t_diff float64) float64 {
	// M
	m := (m0 / 360.0) + m1*t_diff + 0.5*m2*t_diff*t_diff
	fmt.Println("m (Rev) =", m)

	// RevからRadへの変換
	fracM := m - float64(int64(m))
	fracM_Degree := fracM * 360.0
	fracM_Radian := fracM * (2 * math.Pi)

	fmt.Println("fracM (Degree) =", fracM_Degree)

	return fracM_Radian
}

func calculatePositionInOrbitalPlane(a, ecc, eccentricAnomaly float64) (float64, float64) {
	// 軌道平面上の位置を算出する
	u := a*math.Cos(eccentricAnomaly) - a*ecc
	v := a * math.Sqrt(1-ecc*ecc) * math.Sin(eccentricAnomaly)

	return u, v
}

func calculatePerturbationCorrection(angleOmegaA0, angleOmegaB0, angleI0, a, t_diff float64) (float64, float64) {
	// 摂動を補正する
	angleOmegaA_Degree := angleOmegaA0 +
		(180*0.174*(2-2.5*math.Pow(math.Sin(angleI0*math.Pi/180.0), 2)))/(math.Pi*math.Pow(a/EarthRadius, 3.5))*t_diff
	angleOmegaB_Degree := angleOmegaB0 -
		(180*0.174*math.Cos(angleI0*math.Pi/180.0))/(math.Pi*math.Pow(a/EarthRadius, 3.5))*t_diff

	return angleOmegaA_Degree, angleOmegaB_Degree
}

func deg2rad(deg float64) float64 {
	return deg / 180.0 * math.Pi
}

func rad2deg(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

func transformToEquatorial(u, v float64, angleOmegaA_Rad, angleOmegaB_Rad, angleI0_Rad float64) (float64, float64, float64) {
	// 軌道平面から春分点をX軸とする三次元座標系へ変換する
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

func transformToEarthFixed(x, y, z float64, targetTime time.Time) (float64, float64, float64) {
	// 春分点基軸の座標系から本初子午線を基軸とした座標系へ
	tmp := time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC)
	sheta0 := 0.27644444444
	t_diff1 := float64(targetTime.Unix()-tmp.Unix()) / 86400.0
	sheta := sheta0 + 1.002737909*t_diff1

	// 少数部を使って角度へ変換
	shetaG_Deg := (sheta - float64(int64(sheta))) * 360.0
	shetaG_Rad := (sheta - float64(int64(sheta))) * 2.0 * math.Pi

	fmt.Println("ShetaG (Degree) = ", shetaG_Deg)

	// z軸周りの回転行列を逆方向の角度に適用
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

func calculateLatLong(x, y, z float64) (float64, float64) {
	// 緯度・経緯に変換する (注意: 地球を球と仮定)
	fai := math.Asin(z / math.Sqrt(x*x+y*y+z*z))
	lambda := math.Atan2(y, x)

	// ラジアンから度に変換
	fai_Degree := rad2deg(fai)
	lambda_Degree := rad2deg(lambda)

	return fai_Degree, lambda_Degree
}

func calculateAltitude(x, y, z float64) float64 {
	// 高度の算出 (注意: 地球を球と仮定)
	return math.Sqrt(x*x+y*y+z*z) - EarthRadius
}
