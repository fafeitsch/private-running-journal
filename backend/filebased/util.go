package filebased

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (s *Service) deleteEmptyDirectories(path string) {
	parts := strings.Split(path, string(filepath.Separator))
	if len(parts) == 0 {
		return
	}
	counter := 1
	directory := filepath.Join(parts[:len(parts)-counter]...)
	file, err := os.Open(filepath.Join(s.path, directory))
	if err != nil {
		log.Printf("could not check whether directory is empty: %v", err)
		return
	}
	_, err = file.Readdirnames(1)
	_ = file.Close()
	for err == io.EOF && counter < len(parts) {
		_ = file.Close()
		log.Printf("deleting %v", directory)
		_ = os.RemoveAll(filepath.Join(s.path, directory))
		counter = counter + 1
		directory = filepath.Join(parts[:len(parts)-counter]...)
		file, err = os.Open(filepath.Join(s.path, directory))
		if err != nil {
			log.Printf("could not check whether directory is empty: %v", err)
			return
		}
		_, err = file.Readdirnames(1)
		log.Printf("%s %v", directory, err)
	}
	_ = file.Close()
}
