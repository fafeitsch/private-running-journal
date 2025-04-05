package dashboard

import (
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"github.com/labstack/gommon/log"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"
)

type DashboardDto struct {
	TotalDistance    int                `json:"totalDistance"`
	TopTracks        []Track            `json:"topTracks"`
	TotalRuns        int                `json:"totalRuns"`
	MedianDistance   int                `json:"medianDistance"`
	AverageDistance  int                `json:"averageDistance"`
	MonthlyAnalytics []MonthlyAnalytics `json:"analytics"`
}

type Track struct {
	Name    string   `json:"name"`
	Id      string   `json:"id"`
	Parents []string `json:"parents"`
	Count   int      `json:"count"`
	Length  int      `json:"length"`
}

type MonthlyAnalytics struct {
	Month           int `json:"month"`
	Year            int `json:"year"`
	TotalDistance   int `json:"totalDistance"`
	MedianDistance  int `json:"medianDistance"`
	AverageDistance int `json:"averageDistance"`
	TotalRuns       int `json:"totalRuns"`
}

type Assembler struct {
	sortedEntries *projection.SortedJournalEntries
	fileService   *filebased.Service
}

func NewAssembler(sortedEntries *projection.SortedJournalEntries, fileService *filebased.Service) *Assembler {
	return &Assembler{sortedEntries: sortedEntries, fileService: fileService}
}

type Options struct {
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
	TopTracks int       `json:"topTracks"`
}

type entry struct {
	id     string
	length int
	date   time.Time
}

func (a *Assembler) LoadDashboard(options Options) (*DashboardDto, error) {
	runsPerDay, tracks, err := a.readRunsPerDay(options)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	trackCounter := make(map[string]int)
	entryPerMonth := make(map[string][]int)
	lengths := make([]int, 0, 0)
	for _, entries := range runsPerDay {
		if len(entries) == 0 {
			continue
		}
		length := 0
		for _, entry := range entries {
			length = length + entry.length
			trackCounter[entry.id] = trackCounter[entry.id] + 1
		}
		lengths = append(lengths, length)
		month := entries[0].date.Format("2006-01")
		entryPerMonth[month] = append(entryPerMonth[month], length)
	}
	topTracks := make([]Track, 0, options.TopTracks)
	for id, track := range tracks {
		topTracks = append(
			topTracks, Track{
				Id:      id,
				Name:    track.Name,
				Count:   trackCounter[id],
				Parents: track.Parents,
				Length:  track.Waypoints.Length(),
			},
		)
	}
	monthlyAnalytics := createAnalytics(entryPerMonth)
	slices.SortFunc(topTracks, compareTracks)
	sum := 0
	slices.Sort(lengths)
	for _, length := range lengths {
		sum = sum + length
	}
	median := 0
	average := 0
	if len(lengths) > 0 {
		median = lengths[len(lengths)/2]
		average = sum / len(lengths)
	}

	return &DashboardDto{
		TotalDistance:    sum,
		MedianDistance:   median,
		AverageDistance:  average,
		TopTracks:        topTracks[:int(math.Min(float64(options.TopTracks), float64(len(topTracks))))],
		TotalRuns:        len(lengths),
		MonthlyAnalytics: monthlyAnalytics,
	}, nil
}

func (a *Assembler) readRunsPerDay(options Options) (map[string][]entry, map[string]shared.Track, error) {
	entries, err := a.sortedEntries.FindJournalEntryIdsBetween(options.From, options.To)
	if err != nil {
		return nil, nil, err
	}
	trackCache := make(map[string]shared.Track)
	entryPerDay := make(map[string][]entry)
	currentMonth := options.From.AddDate(0, 0, -options.From.Day()+1)
	lastOfToMonth := options.To.AddDate(0, 1, -options.To.Day())
	for currentMonth.Before(lastOfToMonth) || currentMonth.Equal(lastOfToMonth) {
		entryPerDay[currentMonth.Format(time.DateOnly)] = make([]entry, 0)
		currentMonth = currentMonth.AddDate(0, 1, 0)
	}
	lengths := make([]int, 0, 0)
	for _, entryId := range entries {
		loaded, err := a.fileService.ReadJournalEntry(entryId)
		if err != nil {
			return nil, nil, err
		}
		track, ok := trackCache[loaded.TrackId]
		if !ok {
			track, err = a.fileService.ReadTrack(loaded.TrackId)
			if err != nil {
				return nil, nil, err
			}
			trackCache[loaded.TrackId] = track
		}
		length := track.Waypoints.Length() * loaded.Laps
		if loaded.CustomLength != nil {
			length = *loaded.CustomLength
		}
		lengths = append(lengths, length)
		month := loaded.Date.Format(time.DateOnly)
		entryPerDay[month] = append(
			entryPerDay[month], entry{
				id:     loaded.TrackId,
				length: length,
				date:   loaded.Date,
			},
		)
	}
	return entryPerDay, trackCache, nil
}

func createAnalytics(entries map[string][]int) []MonthlyAnalytics {
	result := make([]MonthlyAnalytics, 0, len(entries))
	for key, list := range entries {
		splitted := strings.Split(key, "-")
		year, _ := strconv.Atoi(splitted[0])
		month, _ := strconv.Atoi(splitted[1])
		slices.Sort(list)
		sum := 0
		for _, length := range list {
			sum = sum + length
		}
		median := 0
		average := 0
		if len(list) > 0 {
			median = list[len(list)/2]
			average = sum / len(list)
		}
		result = append(
			result, MonthlyAnalytics{
				Year:            year,
				Month:           month,
				MedianDistance:  median,
				TotalDistance:   sum,
				TotalRuns:       len(list),
				AverageDistance: average,
			},
		)
	}
	slices.SortFunc(
		result, func(a, b MonthlyAnalytics) int {
			if a.Year == b.Year {
				return a.Month - b.Month
			}
			return a.Year - b.Year
		},
	)
	return result
}

func compareTracks(track1 Track, track2 Track) int {
	compare := track2.Count - track1.Count
	if compare == 0 {
		compare = strings.Compare(track1.Name, track2.Name)
	}
	if compare == 0 {
		compare = strings.Compare(strings.Join(track1.Parents, ""), strings.Join(track2.Parents, ""))
	}
	return compare
}
