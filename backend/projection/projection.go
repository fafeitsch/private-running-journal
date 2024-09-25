package projection

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Projection struct {
	directory   string
	trackUsages map[string]int
	projectors  []Projector
}

type Projector interface {
	BuildProjection() (json.RawMessage, error)
	ProjectionName() string
	Bootstrap(retriever Retriever, rebuilder Rebuilder)
}

type Retriever func() (json.RawMessage, error)

type Rebuilder func(message json.RawMessage) error

func New(configDirectory string, projectors ...Projector) *Projection {
	result := &Projection{directory: filepath.Join(configDirectory, ".projection"), projectors: projectors}
	for _, projector := range projectors {
		name := projector.ProjectionName()
		projector.Bootstrap(
			func() (json.RawMessage, error) {
				return result.readFile(name)
			}, func(message json.RawMessage) error {
				return result.writeProjection(message, name)
			},
		)
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

	for _, projector := range p.projectors {
		message, err := projector.BuildProjection()
		name := projector.ProjectionName()
		if err != nil {
			return fmt.Errorf("could not build projection %s: %v", name, err)
		}
		err = p.writeProjection(message, name)
		if err != nil {
			return fmt.Errorf("could not build projection %s: %v", name, err)
		}
	}
	return nil
}

func (p *Projection) writeProjection(message json.RawMessage, name string) error {
	return os.WriteFile(
		filepath.Join(p.directory, name+".json"), message, 0644,
	)
}

func (p *Projection) TrackUsages(trackId string) int {
	return p.trackUsages[trackId]
}

func (p *Projection) readFile(name string) (json.RawMessage, error) {
	return os.ReadFile(filepath.Join(p.directory, name+".json"))
}
