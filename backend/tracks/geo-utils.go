package tracks

import (
	"github.com/twpayne/go-gpx"
	"math"
)

func degreesToRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func distanceBetweenTwoPoints(latDeg1, lonDeg1, latDeg2, lonDeg2 float64) float64 {
	earthRadius := 6371.8 // Earth radius in kilometers

	lat1 := degreesToRadians(latDeg1)
	lon1 := degreesToRadians(lonDeg1)
	lat2 := degreesToRadians(latDeg2)
	lon2 := degreesToRadians(lonDeg2)

	// Haversine formula
	dLat := lat2 - lat1
	dLon := lon2 - lon1
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	return distance
}

func distanceGpx(coords []*gpx.WptType) float64 {
	result := 0.0
	for index := 0; index < len(coords)-1; index++ {
		result = result + distanceBetweenTwoPoints(
			coords[index].Lat, coords[index].Lon, coords[index+1].Lat, coords[index+1].Lon,
		)
	}
	return result
}

func distance(coords []Coordinates) float64 {
	result := 0.0
	for index := 0; index < len(coords)-1; index++ {
		result = result + distanceBetweenTwoPoints(
			coords[index].Latitude, coords[index].Longitude, coords[index+1].Latitude, coords[index+1].Longitude,
		)
	}
	return result
}
