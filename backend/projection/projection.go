package projection

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/private-running-journal/backend/filebased"
	"github.com/fafeitsch/private-running-journal/backend/shared"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Projection struct {
	directory   string
	projectors  []Projector
	fileService *filebased.Service
}

type Projector interface {
	Init(message json.RawMessage, write func())
	AddTrack(track shared.Track)
	AddJournalEntry(entry shared.JournalEntry)
	GetData() any
	ProjectionName() string
}

type Retriever func() (json.RawMessage, error)

type Rebuilder func(message json.RawMessage) error

func New(
	configDirectory string, fileService *filebased.Service, projectors ...Projector,
) *Projection {
	result := &Projection{
		directory:   configDirectory,
		projectors:  projectors,
		fileService: fileService,
	}
	return result
}

func (p *Projection) Initialized() bool {
	_, err := os.Stat(p.directory)
	return err == nil
}

func (p *Projection) Build() error {
	err := os.RemoveAll(p.directory)
	if err != nil {
		return err
	}
	err = os.MkdirAll(p.directory, 0755)
	if err != nil {
		return err
	}

	projectionsToRebuild := make([]Projector, 0)
	for i := range p.projectors {
		message, err := p.readFile(p.projectors[i].ProjectionName())
		if err != nil {
			projectionsToRebuild = append(projectionsToRebuild, p.projectors[i])
			message = nil
		}
		p.projectors[i].Init(
			message, func() {
				p.writePayload(p.projectors[i].GetData(), p.projectors[i].ProjectionName())
			},
		)
	}
	group := sync.WaitGroup{}
	group.Add(2)
	var tracksError error
	go func() {
		tracksError = p.fileService.ReadAllTracks(
			func(track shared.Track) {
				for i := range projectionsToRebuild {
					p.projectors[i].AddTrack(track)
				}
			},
		)
		group.Done()
	}()
	var journalError error
	go func() {
		var entries []shared.JournalEntry
		entries, journalError = p.fileService.ReadAllJournalEntries()
		for i := range projectionsToRebuild {
			for j := range entries {
				p.projectors[i].AddJournalEntry(entries[j])
			}
		}
		group.Done()
	}()
	group.Wait()
	if journalError != nil || tracksError != nil {
		if journalError != nil {
			return fmt.Errorf("could not read journal entries: %v", journalError)
		}
		if tracksError != nil {
			return fmt.Errorf("could not read track entries: %v", tracksError)
		}
	}
	for i := range p.projectors {
		payload := p.projectors[i].GetData()
		p.writePayload(payload, p.projectors[i].ProjectionName())
	}
	return nil
}

func (p *Projection) writePayload(payload any, name string) {
	message, _ := json.Marshal(payload)
	err := p.writeProjection(message, name)
	if err != nil {
		log.Printf("%v", fmt.Errorf("could not write projection %s: %v", name, err))
		//TODO delete file in order to reattempt writing the next time
	}
}

func (p *Projection) writeProjection(message json.RawMessage, name string) error {
	return os.WriteFile(
		filepath.Join(p.directory, name+".json"), message, 0644,
	)
}

func (p *Projection) readFile(name string) (json.RawMessage, error) {
	return os.ReadFile(filepath.Join(p.directory, name+".json"))
}
