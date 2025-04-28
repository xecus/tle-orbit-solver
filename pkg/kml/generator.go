package kml

import (
	"fmt"
	"time"

	"starlink/pkg/model"
)

// GenerateKML creates a KML document for satellite location visualization
func GenerateKML(satellites []string, locations map[string]*model.SatLocation, timestamp time.Time) string {
	// KML header
	kml := `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
<Document>
	<name>Starlink Satellite Locations</name>
	<description>Satellite positions at %s</description>
	<Style id="satellite">
		<IconStyle>
			<Icon>
				<href>http://maps.google.com/mapfiles/kml/shapes/satellite.png</href>
			</Icon>
			<scale>1.0</scale>
		</IconStyle>
		<LabelStyle>
			<scale>0.8</scale>
		</LabelStyle>
	</Style>
`
	// Format the timestamp for description
	kml = fmt.Sprintf(kml, timestamp.Format(time.RFC3339))

	// Add each satellite as a placemark
	for _, satName := range satellites {
		location, exists := locations[satName]
		if !exists {
			continue // Skip satellites that weren't found
		}

		// Add placemark for this satellite
		placemark := `	<Placemark>
		<name>%s</name>
		<description>
			Altitude: %.3f km
			Velocity: %.3f km/s
		</description>
		<styleUrl>#satellite</styleUrl>
		<Point>
			<coordinates>%.6f,%.6f,%.0f</coordinates>
		</Point>
	</Placemark>
`
		// Append formatted placemark to KML
		placemark = fmt.Sprintf(placemark, 
			satName, 
			location.Alt, 
			location.Velocity, 
			location.Lng,  // Longitude goes first in KML
			location.Lat,  // Then latitude
			location.Alt*1000) // Altitude in meters for KML
		
		kml += placemark
	}

	// Add footer
	kml += `</Document>
</kml>`

	return kml
}