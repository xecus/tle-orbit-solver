package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"starlink/pkg/kml"
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

	// Check for flags
	outputKML := false
	kmlFilePath := "starlink_satellites.kml"
	processAllSatellites := false

	// Simple arg parsing
	i := 0
	for i < len(satellites) {
		// Check for KML output flag
		if satellites[i] == "--kml" {
			outputKML = true
			satellites = append(satellites[:i], satellites[i+1:]...)

			// Check if next arg is a filepath
			if i < len(satellites) && !strings.HasPrefix(satellites[i], "-") {
				kmlFilePath = satellites[i]
				satellites = append(satellites[:i], satellites[i+1:]...)
			}
			continue // Don't increment i since we removed an element
		}

		// Check for all satellites flag
		if satellites[i] == "--all" {
			processAllSatellites = true
			satellites = append(satellites[:i], satellites[i+1:]...)
			continue // Don't increment i since we removed an element
		}

		i++ // Move to next argument
	}

	// Default to at least one satellite if all were removed by flag parsing
	// and we're not processing all satellites
	if len(satellites) == 0 && !processAllSatellites {
		satellites = []string{"STARLINK-1008"}
	}

	// Fetch TLE data for all Starlink satellites
	fmt.Println("Fetching TLE data for Starlink satellites...!")
	tleData, err := tle.FetchStarlinkTLEData()
	if err != nil {
		fmt.Printf("Error fetching TLE data: %v\nUsing sample TLE data instead.\n", err)
		panic("Error fetching TLE data")
	}

	// Store satellite locations if needed for KML
	locations := make(map[string]*model.SatLocation)
	currentTime := time.Now()

	// If --all flag is set, get all satellites from TLE data
	if processAllSatellites {
		fmt.Println("Processing all satellites from TLE data...")
		// Get all satellite names
		allSatellites := tle.GetAllSatellites(tleData)
		if len(allSatellites) == 0 {
			fmt.Println("No satellites found in TLE data.")
			return
		}

		fmt.Printf("Found %d satellites in TLE data.\n", len(allSatellites))

		// If there are also specific satellites requested, merge them
		if len(satellites) > 0 {
			// Create a set for faster lookup
			satelliteSet := make(map[string]bool)
			for _, sat := range allSatellites {
				satelliteSet[sat] = true
			}

			// Add any requested satellites that aren't already in the list
			for _, sat := range satellites {
				if _, exists := satelliteSet[sat]; !exists {
					allSatellites = append(allSatellites, sat)
				}
			}
		}

		// Replace the satellites slice with all satellites
		satellites = allSatellites
	}

	// Process each requested satellite
	processedCount := 0
	for _, satelliteName := range satellites {
		location := processSatellite(satelliteName, tleData, nil)
		if location != nil {
			locations[satelliteName] = location
			processedCount++
		}
	}

	fmt.Printf("\nProcessed %d satellites successfully.\n", processedCount)

	// Generate KML if requested
	if outputKML && len(locations) > 0 {
		fmt.Printf("\nGenerating KML file at: %s\n", kmlFilePath)

		// For KML, we only include satellites that were successfully processed
		var validSatellites []string
		for satName := range locations {
			validSatellites = append(validSatellites, satName)
		}

		kmlContent := kml.GenerateKML(validSatellites, locations, currentTime)

		err := os.WriteFile(kmlFilePath, []byte(kmlContent), 0644)
		if err != nil {
			fmt.Printf("Error writing KML file: %v\n", err)
		} else {
			fmt.Printf("KML file created successfully with %d satellites.\n", len(validSatellites))
			fmt.Println("You can open this file in Google Earth or import it into Google Maps.")
		}
	}
}

// processSatellite processes a single satellite, calculating and displaying its position
// Returns the location for KML generation if successful
func processSatellite(satelliteName, tleData string, defaultElements *model.TleOrbitalElement) *model.SatLocation {
	fmt.Printf("\n--- Processing satellite: %s ---\n", satelliteName)

	var satelliteElements *model.TleOrbitalElement

	// If we have TLE data, try to find this satellite
	if tleData != "" {
		line1, line2, err := tle.FindSatelliteByName(tleData, satelliteName)
		if err != nil {
			fmt.Printf("Error finding satellite %s: %v\n\n", satelliteName, err)
			return nil
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

	// Calculate velocity
	satLocation2 := orbital.CalculateSatelliteLocation(satelliteElements, currentTime.Add(time.Second))
	velocity := orbital.CalculateVelocity(satLocation1, satLocation2)

	// Store velocity in the location object for KML generation
	satLocation1.Velocity = velocity

	// Display results in a more structured format
	fmt.Printf("\nResults for %s at %s:\n", satelliteName, currentTime.Format(time.RFC3339))
	fmt.Printf("  Latitude:  %.6f°\n", satLocation1.Lat)
	fmt.Printf("  Longitude: %.6f°\n", satLocation1.Lng)
	fmt.Printf("  Altitude:  %.3f km\n", satLocation1.Alt)
	fmt.Printf("  Velocity:  %.3f km/s\n", velocity)

	return satLocation1
}
