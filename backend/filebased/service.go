package filebased

import (
	"fmt"
	"os"
	"path/filepath"
)

type Service struct {
	path string
}

var journalDirectory = "journal"

func NewService(path string) *Service {
	return &Service{path: path}
}

func (s *Service) Init() error {
	err := os.MkdirAll(filepath.Join(s.path, journalDirectory), os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create journal directory: %v", err)
	}
	err = os.MkdirAll(filepath.Join(s.path, tracksDirectory), os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create tracks directory: %v", err)
	}
	return nil
}
