package tle

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// FetchStarlinkTLEData fetches TLE data for all Starlink satellites from CelesTrak
func FetchStarlinkTLEData() (string, error) {
	url := "https://celestrak.org/NORAD/elements/gp.php?GROUP=starlink&FORMAT=tle"

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request returned status: %s", resp.Status)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
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
