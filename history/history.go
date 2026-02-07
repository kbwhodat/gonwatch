package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type WatchedItem struct {
	Type       string    `json:"type"`                // "episode", "anime_episode", "movie"
	TmdbID     int64     `json:"tmdb_id"`
	SeasonNum  int       `json:"season_num"`
	SeasonID   string    `json:"season_id,omitempty"`
	EpisodeNum int64     `json:"episode_num"`
	Title      string    `json:"title"`
	WatchedAt  time.Time `json:"watched_at"`
}

type History struct {
	Items []WatchedItem `json:"items"`
}

var (
	cache     *History
	cacheMu   sync.RWMutex
	cacheFile string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/tmp"
	}
	cacheFile = filepath.Join(homeDir, ".config", "gonwatch", "history.json")
}

func getConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/gonwatch"
	}
	return filepath.Join(homeDir, ".config", "gonwatch")
}

func Load() (*History, error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	if cache != nil {
		return cache, nil
	}

	configDir := getConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config dir: %w", err)
	}

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			cache = &History{Items: []WatchedItem{}}
			return cache, nil
		}
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	var h History
	if err := json.Unmarshal(data, &h); err != nil {
		// If file is corrupted, start fresh
		cache = &History{Items: []WatchedItem{}}
		return cache, nil
	}

	cache = &h
	return cache, nil
}

func Save() error {
	cacheMu.RLock()
	defer cacheMu.RUnlock()

	if cache == nil {
		return nil
	}

	configDir := getConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}

	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write history file: %w", err)
	}

	return nil
}

func MarkWatched(itemType string, tmdbID int64, seasonNum int, episodeNum int64, title string) {
	h, err := Load()
	if err != nil {
		return
	}

	cacheMu.Lock()
	defer cacheMu.Unlock()

	for i, item := range h.Items {
		if item.TmdbID == tmdbID && item.SeasonNum == seasonNum && item.EpisodeNum == episodeNum && item.Type == itemType {
			h.Items[i].WatchedAt = time.Now()
			go Save()
			return
		}
	}

	h.Items = append(h.Items, WatchedItem{
		Type:       itemType,
		TmdbID:     tmdbID,
		SeasonNum:  seasonNum,
		EpisodeNum: episodeNum,
		Title:      title,
		WatchedAt:  time.Now(),
	})

	go Save()
}

func IsWatched(itemType string, tmdbID int64, seasonNum int, episodeNum int64) bool {
	h, err := Load()
	if err != nil {
		return false
	}

	cacheMu.RLock()
	defer cacheMu.RUnlock()

	for _, item := range h.Items {
		if item.TmdbID == tmdbID && item.SeasonNum == seasonNum && item.EpisodeNum == episodeNum && item.Type == itemType {
			return true
		}
	}
	return false
}

func IsMovieWatched(tmdbID int64) bool {
	return IsWatched("movie", tmdbID, 0, 0)
}

func IsEpisodeWatched(tmdbID int64, seasonNum int, episodeNum int64) bool {
	return IsWatched("episode", tmdbID, seasonNum, episodeNum)
}

func IsAnimeEpisodeWatched(tmdbID int64, seasonNum int, episodeNum int64) bool {
	return IsWatched("anime_episode", tmdbID, seasonNum, episodeNum)
}

func MarkAnimeEpisodeWatchedBySeasonID(tmdbID int64, seasonID string, episodeNum int64, title string) {
	h, err := Load()
	if err != nil {
		return
	}

	cacheMu.Lock()
	defer cacheMu.Unlock()

	for i, item := range h.Items {
		if item.TmdbID == tmdbID && item.SeasonID == seasonID && item.EpisodeNum == episodeNum && item.Type == "anime_episode" {
			h.Items[i].WatchedAt = time.Now()
			go Save()
			return
		}
	}

	h.Items = append(h.Items, WatchedItem{
		Type:       "anime_episode",
		TmdbID:     tmdbID,
		SeasonID:   seasonID,
		SeasonNum:  0,
		EpisodeNum: episodeNum,
		Title:      title,
		WatchedAt:  time.Now(),
	})

	go Save()
}

func IsAnimeEpisodeWatchedBySeasonID(tmdbID int64, seasonID string, episodeNum int64) bool {
	h, err := Load()
	if err != nil {
		return false
	}

	cacheMu.RLock()
	defer cacheMu.RUnlock()

	for _, item := range h.Items {
		if item.TmdbID == tmdbID && item.SeasonID == seasonID && item.EpisodeNum == episodeNum && item.Type == "anime_episode" {
			return true
		}
	}
	return false
}
