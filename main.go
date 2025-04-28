package main

import (
	"fmt"
	"os"
	"time"

	"starlink/pkg/model"
	"starlink/pkg/orbital"
	"starlink/pkg/tle"
	"starlink/pkg/util"
)

func main() {
	// Configure logging from environment variable
	util.GetLogLevelFromEnv()
	
	// Get satellite names from command line args or use default
	satellites := []string{"STARLINK-1008"}
	if len(os.Args) > 1 {
		satellites = os.Args[1:]
	}

	// Fetch TLE data for all Starlink satellites
	fmt.Println("Fetching TLE data for Starlink satellites...")
	tleData, err := tle.FetchStarlinkTLEData()
	if err != nil {
		fmt.Printf("Error fetching TLE data: %v\nUsing default TLE data instead.\n", err)
		return
	}

	// Process each requested satellite
	for _, satelliteName := range satellites {
		processSatellite(satelliteName, tleData, nil)
	}
}

// processSatellite processes a single satellite, calculating and displaying its position
func processSatellite(satelliteName, tleData string, defaultElements *model.TleOrbitalElement) {
	fmt.Printf("\n--- Processing satellite: %s ---\n", satelliteName)

	var satelliteElements *model.TleOrbitalElement

	// If we have TLE data, try to find this satellite
	if tleData != "" {
		line1, line2, err := tle.FindSatelliteByName(tleData, satelliteName)
		if err != nil {
			fmt.Printf("Error finding satellite %s: %v\n\n", satelliteName, err)
			return
		}
		fmt.Printf("Found TLE data for %s\n", satelliteName)
		satelliteElements = tle.ParseTleFromStrings(line1, line2)
	}
	if defaultElements != nil {
		// Use provided default elements if available
		satelliteElements = defaultElements
	}

	// Calculate satellite position at current time
	currentTime := time.Now()
	satLocation1 := orbital.CalculateSatelliteLocation(satelliteElements, currentTime)

	// Display results in a more structured format
	fmt.Printf("\nResults for %s at %s:\n", satelliteName, currentTime.Format(time.RFC3339))
	fmt.Printf("  Latitude:  %.6f°\n", satLocation1.Lat)
	fmt.Printf("  Longitude: %.6f°\n", satLocation1.Lng)
	fmt.Printf("  Altitude:  %.3f km\n", satLocation1.Alt)

	// If wanted, calculate and display velocity
	satLocation2 := orbital.CalculateSatelliteLocation(satelliteElements, currentTime.Add(time.Second))
	velocity := orbital.CalculateVelocity(satLocation1, satLocation2)
	fmt.Printf("  Velocity:  %.3f km/s\n", velocity)
}
