package storage

import (
	"encoding/json"
	"os"
)

type Storage struct {
	filename   string
	fetchedMap map[string]bool
}

// NewStorage initializes a new storage and loads existing data from the file
func NewStorage(filename string) (*Storage, error) {
	s := &Storage{
		filename:   filename,
		fetchedMap: make(map[string]bool),
	}

	err := s.Load()
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return s, nil
}

// Load reads the JSON file and populates the map
func (s *Storage) Load() error {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		return err
	}

	var urls []string
	if err := json.Unmarshal(data, &urls); err != nil {
		return err
	}

	for _, url := range urls {
		s.fetchedMap[url] = true
	}

	return nil
}

// IsFetched checks if a URL has already been fetched
func (s *Storage) IsFetched(url string) bool {
	return s.fetchedMap[url]
}

// MarkFetched marks a URL as fetched in memory
func (s *Storage) MarkFetched(url string) {
	s.fetchedMap[url] = true
}

// Save writes the current state of the map back to the JSON file
func (s *Storage) Save() error {
	urls := make([]string, 0, len(s.fetchedMap))
	for url := range s.fetchedMap {
		urls = append(urls, url)
	}

	data, err := json.MarshalIndent(urls, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}
