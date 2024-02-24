package backend

import (
	"github.com/twpayne/go-gpx"
	"log"
	"math"
)

func degreesToRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func distanceBetweenTwoPoints(coord1, coord2 *gpx.WptType) float64 {
	earthRadius := 6371.8 // Earth radius in kilometers

	lat1 := degreesToRadians(coord1.Lat)
	lon1 := degreesToRadians(coord1.Lon)
	lat2 := degreesToRadians(coord2.Lat)
	lon2 := degreesToRadians(coord2.Lon)

	// Haversine formula
	dLat := lat2 - lat1
	dLon := lon2 - lon1
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c
	log.Printf("%f", dLat)

	return distance
}

func distance(coords []*gpx.WptType) float64 {
	result := 0.0
	for index := 0; index < len(coords)-1; index++ {
		result = result + distanceBetweenTwoPoints(coords[index], coords[index+1])
	}
	return result
}
