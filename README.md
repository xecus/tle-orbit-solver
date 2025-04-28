# Starlink Satellite Tracker

A Go application that tracks Starlink satellites by fetching TLE (Two-Line Element) data, calculating orbital positions, and generating KML files for visualization.

## Features

- Calculate satellite positions based on orbital mechanics
- Generate KML files for visualization in Google Earth and other GIS applications
- Support for customizable time intervals and observation periods

## Requirements

- Go 1.18 or higher
- Space-Track.org credentials (for fetching latest TLE data)

## Installation

```bash
git clone https://github.com/xecus/tle-orbit-solver.git
cd starlink
go build
```

## Usage

### Basic Usage

```bash
# Run with default settings
hiroyuki@MacBook-Pro starlink % ./starlink 
Fetching TLE data for Starlink satellites...!

--- Processing satellite: STARLINK-1008 ---
Found TLE data for STARLINK-1008

Results for STARLINK-1008 at 2025-04-28T17:42:32+09:00:
  Latitude:  20.520559°
  Longitude: -167.657067°
  Altitude:  568.407 km
  Velocity:  7.293 km/s

Processed 1 satellites successfully.
```

## How It Works

1. The application fetches the latest TLE data for Starlink satellites from Local or Space-Track.org
2. It parses the TLE data to extract orbital elements
3. Using orbital calculations, it determines satellite positions at specific times
4. The positions are converted to a KML file format for visualization
5. The resulting KML file can be opened in Google Earth or other GIS applications

## Project Structure

- `main.go`: Application entry point and command-line interface
- `pkg/`:
  - `kepler/`: Kepler's laws implementation for orbital mechanics
  - `kml/`: KML file generation utilities
  - `model/`: Data models and types
  - `orbital/`: Orbital calculations and conversions
  - `tle/`: TLE data fetching and parsing
  - `util/`: Utility functions for conversions and logging

