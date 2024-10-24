package trackEditor

import (
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"github.com/fafeitsch/private-running-journal/backend/tracks"
	"path/filepath"
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
	PolylineMeta
	Id        string          `json:"id"`
	Name      string          `json:"name"`
	Waypoints []CoordinateDto `json:"waypoints"`
	Parents   []string        `json:"parents"`
	Usages    []string        `json:"usages"`
}

type PolylineMeta struct {
	Length          int              `json:"length"`
	DistanceMarkers []DistanceMarker `json:"distanceMarkers"`
}

type TrackEditor struct {
	service     *filebased.Service
	trackUsages *projection.TrackUsagesProjector
	trackIdMap  *tracks.TrackIdMapProjection
}

func New(
	service *filebased.Service, trackUsages *projection.TrackUsagesProjector, trackIdMap *tracks.TrackIdMapProjection,
) *TrackEditor {
	return &TrackEditor{service: service, trackUsages: trackUsages, trackIdMap: trackIdMap}
}

func (t *TrackEditor) GetTrack(id string) (TrackDto, error) {
	parents, err := t.trackIdMap.GetTrackLocation(id)
	if err != nil {
		return TrackDto{}, err
	}
	path := filepath.Join(parents...)
	file, err := t.service.ReadTrack(path)
	if err != nil {
		return TrackDto{}, err
	}
	waypoints := make([]CoordinateDto, 0)
	for _, waypoint := range file.Waypoints {
		waypoints = append(waypoints, CoordinateDto{Latitude: waypoint.Latitude, Longitude: waypoint.Longitude})
	}
	distanceMarkers := mapDistanceMarkerToDto(file.Waypoints)
	usages, err := t.trackUsages.GetUsages(id)
	if err != nil {
		return TrackDto{}, err
	}
	return TrackDto{
		Id:        file.Id,
		Name:      file.Name,
		Waypoints: waypoints,
		PolylineMeta: PolylineMeta{
			Length:          file.Waypoints.Length(),
			DistanceMarkers: distanceMarkers,
		},
		Parents: file.Parents,
		Usages:  usages,
	}, nil
}

func mapDistanceMarkerToDto(coordinates shared.Waypoints) []DistanceMarker {
	distanceMarkers := make([]DistanceMarker, 0)
	for _, dm := range coordinates.DistanceMarkers() {
		distanceMarkers = append(
			distanceMarkers, DistanceMarker{
				CoordinateDto: CoordinateDto{
					Latitude:  dm.Latitude,
					Longitude: dm.Longitude,
				}, Distance: dm.Distance,
			},
		)
	}
	return distanceMarkers
}

func (t *TrackEditor) GetPolylineMeta(dtos []CoordinateDto) PolylineMeta {
	coordinates := make(shared.Waypoints, 0)
	for _, dto := range dtos {
		coordinates = append(coordinates, shared.Coordinates{Longitude: dto.Longitude, Latitude: dto.Latitude})
	}
	return PolylineMeta{
		DistanceMarkers: mapDistanceMarkerToDto(coordinates),
		Length:          coordinates.Length(),
	}
}

type SaveTrackDto struct {
	Id        string          `json:"id"`
	Name      string          `json:"name"`
	Waypoints []CoordinateDto `json:"waypoints"`
	Parents   []string        `json:"parents"`
}

func (t *TrackEditor) SaveTrack(track SaveTrackDto) error {
	oldPath, err := t.trackIdMap.GetTrackLocation(track.Id)
	if err != nil {
		return fmt.Errorf("could not read old track location: %v", err)
	}
	if oldPath != nil {
		err = t.service.DeleteTrackDirectory(oldPath)
		if err != nil {
			return fmt.Errorf("could not delete old track file: %v", err)
		}
	}
	wp := make(shared.Waypoints, 0)
	for _, waypoint := range track.Waypoints {
		wp = append(wp, shared.Coordinates{Longitude: waypoint.Longitude, Latitude: waypoint.Latitude})
	}
	saveTrack := shared.SaveTrack{
		Id:        track.Id,
		Name:      track.Name,
		Waypoints: wp,
		Parents:   track.Parents,
	}
	err = t.service.SaveTrack(saveTrack)
	shared.SendEvent(shared.TrackUpsertedEvent{SaveTrack: &saveTrack})
	return err
}

func (t *TrackEditor) DeleteTrack(id string) error {
	path, err := t.trackIdMap.GetTrackLocation(id)
	if err != nil {
		return fmt.Errorf("could not read track location: %v", err)
	}
	if path == nil {
		return fmt.Errorf("track with id %s does not exist", id)
	}
	err = t.service.DeleteTrackDirectory(path)
	if err != nil {
		return fmt.Errorf("could not delete track directory: %v", err)
	}
	shared.SendEvent(shared.TrackDeletedEvent{Id: id})
	return err
}
