package storage

import (
	"encoding/json"
	"os"
)

type Storage struct {
	filename   string
	fetchedMap map[string]bool
}

// NewStorage は新しいストレージを初期化し、ファイルから既存のデータを読み込みます
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

// Load はJSONファイルを読み込み、マップにデータを格納します
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

// IsFetched はURLが既に取得済みかどうかを確認します
func (s *Storage) IsFetched(url string) bool {
	return s.fetchedMap[url]
}

// MarkFetched はURLをメモリ上で取得済みとしてマークします
func (s *Storage) MarkFetched(url string) {
	s.fetchedMap[url] = true
}

// Save はマップの現在の状態をJSONファイルに書き戻します
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
