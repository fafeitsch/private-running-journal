package dashboard

import (
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/projection"
	"github.com/fafeitsch/private-running-journal/backend/shared"
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
	Month          int `json:"month"`
	Year           int `json:"year"`
	TotalDistance  int `json:"totalDistance"`
	MedianDistance int `json:"medianDistance"`
	TotalRuns      int `json:"totalRuns"`
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
}

func (a *Assembler) LoadDashboard(options Options) (*DashboardDto, error) {
	entries, err := a.sortedEntries.FindJournalEntryIdsBetween(options.From, options.To)
	if err != nil {
		return nil, err
	}
	sum := 0
	trackCache := make(map[string]shared.Track)
	trackCounter := make(map[string]int)
	entryPerMonth := make(map[string][]entry)
	currentMonth := options.From.AddDate(0, 0, -options.From.Day()+1)
	lastOfToMonth := options.To.AddDate(0, 1, -options.To.Day())
	for currentMonth.Before(lastOfToMonth) || currentMonth.Equal(lastOfToMonth) {
		entryPerMonth[currentMonth.Format("2006-01")] = make([]entry, 0)
		currentMonth = currentMonth.AddDate(0, 1, 0)
	}
	for _, entryId := range entries {
		loaded, err := a.fileService.ReadJournalEntry(entryId)
		if err != nil {
			return nil, err
		}
		track, ok := trackCache[loaded.TrackId]
		trackCounter[loaded.TrackId] = trackCounter[loaded.TrackId] + 1
		if !ok {
			track, err = a.fileService.ReadTrack(loaded.TrackId)
			if err != nil {
				return nil, err
			}
			trackCache[loaded.TrackId] = track
		}
		length := track.Waypoints.Length() * loaded.Laps
		if loaded.CustomLength != nil {
			length = *loaded.CustomLength
		}
		sum = sum + (length)
		month := loaded.Date.Format("2006-01")
		entryPerMonth[month] = append(entryPerMonth[month], entry{
			id:     loaded.TrackId,
			length: length,
		})
	}
	topTracks := make([]Track, 0, options.TopTracks)
	for id, track := range trackCache {
		topTracks = append(topTracks, Track{Id: id, Name: track.Name, Count: trackCounter[id], Parents: track.Parents, Length: track.Waypoints.Length()})
	}
	monthlyAnalytics := createAnalytics(entryPerMonth)
	slices.SortFunc(topTracks, compareTracks)
	return &DashboardDto{
		TotalDistance:    sum,
		TopTracks:        topTracks[:int(math.Min(float64(options.TopTracks), float64(len(topTracks))))],
		TotalRuns:        len(entries),
		MonthlyAnalytics: monthlyAnalytics,
	}, nil
}

func createAnalytics(entries map[string][]entry) []MonthlyAnalytics {
	result := make([]MonthlyAnalytics, 0, len(entries))
	for key, list := range entries {
		splitted := strings.Split(key, "-")
		year, _ := strconv.Atoi(splitted[0])
		month, _ := strconv.Atoi(splitted[1])
		slices.SortFunc(list, func(a, b entry) int {
			return a.length - b.length
		})
		median := 0
		if len(list) > 0 {
			median = list[len(list)/2].length
		}
		sum := 0
		for _, entry := range list {
			sum = sum + entry.length
		}
		result = append(result, MonthlyAnalytics{Year: year, Month: month, MedianDistance: median, TotalDistance: sum, TotalRuns: len(list)})
	}
	slices.SortFunc(result, func(a, b MonthlyAnalytics) int {
		if a.Year == b.Year {
			return a.Month - b.Month
		}
		return a.Year - b.Year
	})
	return result
}

func compareTracks(track1 Track, track2 Track) int {
	compare := track2.Count - track1.Count
	if compare == 0 {
		compare = strings.Compare(track2.Name, track1.Name)
	}
	if compare == 0 {
		compare = strings.Compare(strings.Join(track2.Parents, ""), strings.Join(track1.Parents, ""))
	}
	return compare
}
