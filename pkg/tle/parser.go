package tle

import (
	"strconv"
	"strings"

	"starlink/pkg/model"
	"starlink/pkg/util"
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
	// TLEフォーマットは固定長なので、位置ベースで抽出する
	// TLE format reference: https://celestrak.org/NORAD/documentation/tle-fmt.php

	// Debug output of raw TLE lines
	util.LogDebug("TLE Line 1: %s\n", str1)
	util.LogDebug("TLE Line 2: %s\n", str2)

	// Parse line 1 data - Using fixed positions as per TLE format definition
	// Example: 1 44714U 19074B   25117.42924319 -.00001157  00000+0 -58773-4 0  9990
	//          1         2         3         4         5         6         7
	//          123456789012345678901234567890123456789012345678901234567890123456789

	satelliteNumber := strings.TrimSpace(str1[2:7])

	// Get international designator
	internationalDesignator := strings.TrimSpace(str1[9:17])

	// Epoch year and day
	etYearStr := strings.TrimSpace(str1[18:20])
	etYear := 0
	// Convert 2-digit year to 4-digit (assuming 20xx for now)
	if val, err := strconv.Atoi(etYearStr); err == nil {
		etYear = 2000 + val
	}

	etDayStr := strings.TrimSpace(str1[20:32])
	etDay, _ := strconv.ParseFloat(etDayStr, 64)

	// First Time Derivative of the Mean Motion
	firstDerStr := strings.TrimSpace(str1[33:43])
	firstTimeDerivativeOfTheMeanMotion, _ := strconv.ParseFloat(firstDerStr, 64)

	// Second Time Derivative of Mean Motion (decimal point assumed)
	secondTimeDerivativeOfTheMeanMotion := strings.TrimSpace(str1[44:52])

	// B* drag term
	bstarDragTerm := strings.TrimSpace(str1[53:61])

	// Element number and checksum
	elementnum := strings.TrimSpace(str1[64:68])
	checksum1 := string(str1[68])

	// Parse line 2 data - Using fixed positions as per TLE format definition
	// Example: 2 44714  53.0517 166.3609 0001116  99.1558 260.9557 15.06400606301084
	//          1         2         3         4         5         6         7
	//          123456789012345678901234567890123456789012345678901234567890123456789

	// Orbital Inclination (degrees)
	orbitalInclination, _ := strconv.ParseFloat(strings.TrimSpace(str2[8:16]), 64)

	// Right Ascension of the Ascending Node (degrees)
	rightAscensionOfAscendingNode, _ := strconv.ParseFloat(strings.TrimSpace(str2[17:25]), 64)

	// Eccentricity (decimal point assumed at beginning)
	eccentricity, _ := strconv.ParseFloat("0."+strings.TrimSpace(str2[26:33]), 64)

	// Argument of Perigee (degrees)
	argumentOfPerigee, _ := strconv.ParseFloat(strings.TrimSpace(str2[34:42]), 64)

	// Mean Anomaly (degrees)
	meanAnomaly, _ := strconv.ParseFloat(strings.TrimSpace(str2[43:51]), 64)

	// Mean Motion (revolutions per day)
	meanMotionStr := strings.TrimSpace(str2[52:63])
	meanMotion, _ := strconv.ParseFloat(meanMotionStr, 64)

	// Revolution number at epoch and checksum
	numberOfLaps := strings.TrimSpace(str2[63:68])
	checksum2 := string(str2[68])

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

	util.LogDebug("--------------------------------------------------\n")
	util.LogDebug("satelliteNumber=%s\n", satelliteNumber)
	util.LogDebug("internationalDesignator=%s\n", internationalDesignator)
	util.LogDebug("etYear=%d\n", etYear)
	util.LogDebug("etDay=%f\n", etDay)
	util.LogDebug("firstTimeDerivativeOfTheMeanMotion=%f\n",
		firstTimeDerivativeOfTheMeanMotion)
	util.LogDebug("secondTimeDerivativeOfTheMeanMotion=%s\n",
		secondTimeDerivativeOfTheMeanMotion)
	util.LogDebug("bstarDragTerm=%s\n", bstarDragTerm)
	util.LogDebug("elementnum=%s\n", elementnum)
	util.LogDebug("checksum1=%s\n", checksum1)
	util.LogDebug("orbitalInclination=%f [Degree]\n", orbitalInclination)
	util.LogDebug("RAAN=%f [Degree]\n", rightAscensionOfAscendingNode)
	util.LogDebug("eccentricity=%f [-]\n", eccentricity)
	util.LogDebug("argumentOfPerigee=%f [Degree]\n", argumentOfPerigee)
	util.LogDebug("meanAnomaly=%f [Degree]\n", meanAnomaly)
	util.LogDebug("meanMotion=%f [Rev/Day]\n", meanMotion)
	util.LogDebug("numberOfLaps=%s [-]\n", numberOfLaps)
	util.LogDebug("checksum2=%s [-]\n", checksum2)
	util.LogDebug("--------------------------------------------------\n")
}
