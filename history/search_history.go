package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type SearchHistory struct {
	Movies []string `json:"movies"`
	Series []string `json:"series"`
	Anime  []string `json:"anime"`
}

const maxSearchHistory = 50

var (
	searchCache   *SearchHistory
	searchCacheMu sync.RWMutex
	searchFile    string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/tmp"
	}
	searchFile = filepath.Join(homeDir, ".config", "gonwatch", "search_history.json")
}

func LoadSearchHistory() (*SearchHistory, error) {
	searchCacheMu.Lock()
	defer searchCacheMu.Unlock()

	if searchCache != nil {
		return searchCache, nil
	}

	configDir := getConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config dir: %w", err)
	}

	data, err := os.ReadFile(searchFile)
	if err != nil {
		if os.IsNotExist(err) {
			searchCache = &SearchHistory{
				Movies: []string{},
				Series: []string{},
				Anime:  []string{},
			}
			return searchCache, nil
		}
		return nil, fmt.Errorf("failed to read search history file: %w", err)
	}

	var h SearchHistory
	if err := json.Unmarshal(data, &h); err != nil {
		searchCache = &SearchHistory{
			Movies: []string{},
			Series: []string{},
			Anime:  []string{},
		}
		return searchCache, nil
	}

	searchCache = &h
	return searchCache, nil
}

func SaveSearchHistory() error {
	searchCacheMu.RLock()
	defer searchCacheMu.RUnlock()

	if searchCache == nil {
		return nil
	}

	configDir := getConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	data, err := json.MarshalIndent(searchCache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal search history: %w", err)
	}

	if err := os.WriteFile(searchFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write search history file: %w", err)
	}

	return nil
}

func AddSearch(category, term string) {
	term = strings.TrimSpace(term)
	if term == "" {
		return
	}

	h, err := LoadSearchHistory()
	if err != nil {
		return
	}

	searchCacheMu.Lock()
	defer searchCacheMu.Unlock()

	var list *[]string
	switch category {
	case "movies":
		list = &h.Movies
	case "series":
		list = &h.Series
	case "anime":
		list = &h.Anime
	default:
		return
	}

	newList := []string{term}
	for _, existing := range *list {
		if !strings.EqualFold(existing, term) {
			newList = append(newList, existing)
		}
		if len(newList) >= maxSearchHistory {
			break
		}
	}
	*list = newList

	go SaveSearchHistory()
}

func GetSearches(category string) []string {
	h, err := LoadSearchHistory()
	if err != nil {
		return []string{}
	}

	searchCacheMu.RLock()
	defer searchCacheMu.RUnlock()

	switch category {
	case "movies":
		return h.Movies
	case "series":
		return h.Series
	case "anime":
		return h.Anime
	default:
		return []string{}
	}
}
