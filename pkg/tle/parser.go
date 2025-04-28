package tle

import (
	"fmt"
	"strconv"
	"strings"

	"starlink/pkg/model"
)

// ParseTle parses the default TLE for STARLINK-1008
func ParseTle() *model.TleOrbitalElement {
	// デフォルトは STARLINK-1008 のTLEデータを利用
	str1 := "1 44714U 19074B   25117.42924319 -.00001157  00000+0 -58773-4 0  9990"
	str2 := "2 44714  53.0517 166.3609 0001116  99.1558 260.9557 15.06400606301084"

	return ParseTleFromStrings(str1, str2)
}

// ParseTleFromStrings parses TLE data from two input strings
func ParseTleFromStrings(str1, str2 string) *model.TleOrbitalElement {
	a := strings.Split(str1, " ")
	b := strings.Split(str2, " ")

	// TLEから各パラメータを抽出
	// Debug output
	fmt.Printf("TLE line 1 split count: %d\n", len(a))
	for i, v := range a {
		fmt.Printf("a[%d] = %s\n", i, v)
	}

	// Parse line 1 data
	satelliteNumber := a[1]
	internationalDesignator := a[2]
	etYear, _ := strconv.Atoi(fmt.Sprintf("20%s", a[5][0:2]))
	etDay, _ := strconv.ParseFloat(a[5][2:], 64)
	firstTimeDerivativeOfTheMeanMotion, _ := strconv.ParseFloat(fmt.Sprintf("0%s", a[7]), 64)
	secondTimeDerivativeOfTheMeanMotion := a[9]
	bstarDragTerm := a[11]
	elementnum := a[12][0:3]
	checksum1 := a[12][3:]
	
	// Debug output
	fmt.Printf("TLE line 2 split count: %d\n", len(b))
	for i, v := range b {
		fmt.Printf("b[%d] = %s\n", i, v)
	}

	// Parse line 2 data using fixed positions for accuracy
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

	// Display all TLE parameters
	PrintTleParameters(satelliteNumber, internationalDesignator, etYear, etDay,
		firstTimeDerivativeOfTheMeanMotion, secondTimeDerivativeOfTheMeanMotion,
		bstarDragTerm, elementnum, checksum1, orbitalInclination,
		rightAscensionOfAscendingNode, eccentricity, argumentOfPerigee,
		meanAnomaly, meanMotion, numberOfLaps, checksum2)

	return &model.TleOrbitalElement{
		MeanAnomaly:        meanAnomaly,
		MeanMotion:         meanMotion,
		MeanMotionDot:      firstTimeDerivativeOfTheMeanMotion,
		Eccentricity:       eccentricity,
		EtYear:             etYear,
		EtDay:              etDay,
		OrbitalInclination: orbitalInclination,
		Raan:               rightAscensionOfAscendingNode,
		ArgumentOfPerigee:  argumentOfPerigee,
	}
}

// PrintTleParameters outputs all TLE parameters in a readable format
func PrintTleParameters(satelliteNumber, internationalDesignator string, etYear int, etDay float64,
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