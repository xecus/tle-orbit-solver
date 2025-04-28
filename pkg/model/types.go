package model

// SatLocation represents satellite location and position in space
type SatLocation struct {
	X         float64 // X coordinate in Earth-fixed frame [km]
	Y         float64 // Y coordinate in Earth-fixed frame [km]
	Z         float64 // Z coordinate in Earth-fixed frame [km]
	Lat       float64 // Latitude [degree]
	Lng       float64 // Longitude [degree]
	Alt       float64 // Altitude from Earth surface [km]
	Velocity  float64 // Velocity [km/s], optional
}

// TleOrbitalElement contains the orbital elements parsed from a TLE
type TleOrbitalElement struct {
	MeanAnomaly        float64 // M0 平均近点角 [Degree]
	MeanMotion         float64 // M1 平均運動: [Rev/Day]
	MeanMotionDot      float64 // M2 平均運動変化係数: [Rev/Day2]
	Eccentricity       float64 // 離心率 [-]
	EtYear             int     // 元期 Epoctime [Year]
	EtDay              float64 // 元期 EpocTime [Day]
	OrbitalInclination float64 // 軌道傾斜角 [Degree]
	Raan               float64 // 昇交点赤経: RAAN [Degree]
	ArgumentOfPerigee  float64 // 近地点引数 [Degree]
}