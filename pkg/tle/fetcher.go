package tle

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func FetchStarlinkTLEData() (string, error) {
	data, err := FetchStarlinkTLEDataFromLocalFile()
	if err != nil {
		return "", err
	}
	return data, err
}

func FetchStarlinkTLEDataFromLocalFile() (string, error) {
	data, err := os.ReadFile("tle.txt")
	if err != nil {
		return "", fmt.Errorf("failed to read local TLE file: %w", err)
	}
	return string(data), nil
}

// FetchStarlinkTLEData fetches TLE data for all Starlink satellites from CelesTrak
func FetchStarlinkTLEDataFromCelesTrak() (string, error) {
	// Try to fetch data from CelesTrak
	url := "https://celestrak.org/NORAD/elements/gp.php?GROUP=starlink&FORMAT=tle"

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Get Error")
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("HTTP Error")
		return "", errors.New("status Error")
	}
	defer resp.Body.Close()

	// Read response body
	fmt.Println("ReadAll...")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// FindSatelliteByName finds a specific satellite in TLE data by name
func FindSatelliteByName(tleData, satelliteName string) (string, string, error) {
	lines := strings.Split(tleData, "\n")

	// TLE format consists of three lines per satellite: name, line1, line2
	for i := 0; i < len(lines)-2; i++ {
		name := strings.TrimSpace(lines[i])
		if strings.Contains(name, satelliteName) {
			return lines[i+1], lines[i+2], nil
		}
	}

	return "", "", errors.New("satellite not found")
}

// GetAllSatellites extracts all satellites from TLE data
func GetAllSatellites(tleData string) []string {
	lines := strings.Split(tleData, "\n")
	var satellites []string

	// TLE format consists of three lines per satellite: name, line1, line2
	for i := 0; i < len(lines)-2; i++ {
		// Check for TLE lines pattern
		if i+2 < len(lines) &&
			strings.HasPrefix(lines[i+1], "1 ") &&
			strings.HasPrefix(lines[i+2], "2 ") {
			// This line should be a satellite name
			satName := strings.TrimSpace(lines[i])
			if satName != "" {
				satellites = append(satellites, satName)
			}
			// Skip the next two lines since we already processed them
			i += 2
		}
	}

	return satellites
}

// GetSatelliteAndLines returns all satellite names and their TLE lines
func GetSatelliteAndLines(tleData string) map[string][]string {
	lines := strings.Split(tleData, "\n")
	result := make(map[string][]string)

	// TLE format consists of three lines per satellite: name, line1, line2
	for i := 0; i < len(lines)-2; i++ {
		// Check for TLE lines pattern
		if i+2 < len(lines) &&
			strings.HasPrefix(lines[i+1], "1 ") &&
			strings.HasPrefix(lines[i+2], "2 ") {
			// This line should be a satellite name
			satName := strings.TrimSpace(lines[i])
			if satName != "" {
				result[satName] = []string{lines[i+1], lines[i+2]}
			}
			// Skip the next two lines since we already processed them
			i += 2
		}
	}

	return result
}
