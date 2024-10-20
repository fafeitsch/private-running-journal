package trackEditor

import (
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/projection"
)

type CoordinateDto struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type DistanceMarker struct {
	CoordinateDto
	Distance int `json:"distance"`
}

type TrackDto struct {
	Id              string           `json:"id"`
	Name            string           `json:"name"`
	Length          int              `json:"length"`
	Waypoints       []CoordinateDto  `json:"waypoints"`
	DistanceMarkers []DistanceMarker `json:"distanceMarkers"`
	Parents         []string         `json:"parents"`
	Usages          []string         `json:"usages"`
}

type TrackEditor struct {
	service     *filebased.Service
	trackUsages *projection.TrackUsagesProjector
}

func New(service *filebased.Service, trackUsages *projection.TrackUsagesProjector) *TrackEditor {
	return &TrackEditor{service: service, trackUsages: trackUsages}
}

func (t *TrackEditor) GetTrack(id string) (TrackDto, error) {
	file, err := t.service.ReadTrack(id)
	if err != nil {
		return TrackDto{}, err
	}
	waypoints := make([]CoordinateDto, 0)
	for _, waypoint := range file.Waypoints {
		waypoints = append(waypoints, CoordinateDto{Latitude: waypoint.Latitude, Longitude: waypoint.Longitude})
	}
	distanceMarkers := make([]DistanceMarker, 0)
	for _, dm := range file.DistanceMarkers() {
		distanceMarkers = append(
			distanceMarkers, DistanceMarker{
				CoordinateDto: CoordinateDto{
					Latitude:  dm.Latitude,
					Longitude: dm.Longitude,
				}, Distance: dm.Distance,
			},
		)
	}
	usages, err := t.trackUsages.GetUsages(id)
	if err != nil {
		return TrackDto{}, err
	}
	return TrackDto{
		Id:              file.Name,
		Name:            file.Name,
		Length:          file.Waypoints.Length(),
		Waypoints:       waypoints,
		DistanceMarkers: distanceMarkers,
		Parents:         file.Parents,
		Usages:          usages,
	}, nil
}
