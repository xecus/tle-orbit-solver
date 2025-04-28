package main

import (
	"fmt"
	"os"
	"time"

	"starlink/pkg/model"
	"starlink/pkg/orbital"
	"starlink/pkg/tle"
)

func main() {
	// Get satellite name from command line args or use default
	satelliteName := "STARLINK-1008"
	if len(os.Args) > 1 {
		satelliteName = os.Args[1]
	}

	var satelliteElements *model.TleOrbitalElement

	// Try to fetch TLE data for the requested satellite
	if satelliteName != "" {
		fmt.Printf("Fetching TLE data for %s...\n", satelliteName)
		tleData, err := tle.FetchStarlinkTLEData()
		if err != nil {
			fmt.Printf("Error fetching TLE data: %v\nUsing default TLE data instead.\n", err)
			satelliteElements = tle.ParseTle()
		} else {
			line1, line2, err := tle.FindSatelliteByName(tleData, satelliteName)
			if err != nil {
				fmt.Printf("Error finding satellite %s: %v\nUsing default TLE data instead.\n",
					satelliteName, err)
				satelliteElements = tle.ParseTle()
			} else {
				fmt.Printf("Found TLE data for %s\n", satelliteName)
				satelliteElements = tle.ParseTleFromStrings(line1, line2)
			}
		}
	} else {
		satelliteElements = tle.ParseTle()
	}

	// Calculate satellite position at current time
	currentTime := time.Now()
	satLocation1 := orbital.CalculateSatelliteLocation(satelliteElements, currentTime)

	// Calculate satellite position 1 second later to determine velocity
	//satLocation2 := orbital.CalculateSatelliteLocation(satelliteElements, currentTime.Add(time.Second))
	// Calculate and display velocity
	//orbital.CalculateVelocity(satLocation1, satLocation2)
}
