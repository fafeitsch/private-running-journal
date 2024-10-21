package shared

import (
	"math"
)

type Waypoints []Coordinates

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

func (w Waypoints) Length() int {
	result := 0.0
	for index := 0; index < len(w)-1; index++ {
		result = result + distanceBetweenTwoPoints(
			w[index].Latitude, w[index].Longitude, w[index+1].Latitude, w[index+1].Longitude,
		)
	}
	return int(result * 1000)
}

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

func (w Waypoints) DistanceMarkers() []DistanceMarker {
	result := make([]DistanceMarker, 0)
	steps := float64(1000)
	total := 0.0
	for index := 0; index < len(w)-1; index++ {
		distance := distanceBetweenTwoPoints(
			w[index].Latitude, w[index].Longitude, w[index+1].Latitude, w[index+1].Longitude,
		) * 1000 // convert to meters

		if total+distance < steps {
			total = total + distance
			continue
		}
		remainingDistance := steps - total
		ratio := remainingDistance / distance

		lat := w[index].Latitude + ratio*(w[index+1].Latitude-w[index].Latitude)
		lon := w[index].Longitude + ratio*(w[index+1].Longitude-w[index].Longitude)

		result = append(
			result, DistanceMarker{
				Coordinates: Coordinates{Latitude: lat, Longitude: lon},
				Distance:    (len(result) + 1) * int(steps),
			},
		)
		total = distance - remainingDistance
	}

	return result
}

type DistanceMarker struct {
	Coordinates
	Distance int `json:"distance"`
}

type Track struct {
	Waypoints Waypoints
	Id        string   `json:"id"`
	Length    int      `json:"length"`
	Name      string   `json:"name"`
	Parents   []string `json:"parents"`
}

type SaveTrack struct {
	Waypoints Waypoints
	Id        string
	Length    int
	Name      string
	Parents   []string
}

type JournalEntry struct {
	TrackId string `json:"trackId"`
	Id      string `json:"id"`
}
