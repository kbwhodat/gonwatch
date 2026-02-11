package models

import (
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"golang.org/x/term"
	"gonwatch/history"
)

func RecentlyWatchedModel(m *Model) {
	items := history.GetRecentlyWatched(50)

	listItems := make([]list.Item, len(items))

	for i, item := range items {
		listItems[i] = BubbleTeaRecentlyWatchedList{
			ItemType:   item.Type,
			ItemTmdbID: item.TmdbID,
			SeasonNum:  item.SeasonNum,
			SeasonID:   item.SeasonID,
			EpisodeNum: item.EpisodeNum,
			ItemTitle:  item.Title,
			WatchedAt:  formatTimeAgo(item.WatchedAt),
		}
	}

	delegate := list.NewDefaultDelegate()
	m.List = list.New(listItems, delegate, 0, 0)
	m.List.Title = "Recently Watched"

	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	m.List.SetSize(w-2, h-2)

	m.Mode = "list"

	m.List.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{GoBackKeys, GoForwardKeys}
	}

	m.List.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{GoBackKeys, GoForwardKeys}
	}
}

func formatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return formatInt(mins) + " minutes ago"
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return formatInt(hours) + " hours ago"
	case diff < 48*time.Hour:
		return "yesterday"
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		return formatInt(days) + " days ago"
	case diff < 30*24*time.Hour:
		weeks := int(diff.Hours() / 24 / 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return formatInt(weeks) + " weeks ago"
	default:
		return t.Format("Jan 2, 2006")
	}
}

func formatInt(n int) string {
	return strconv.Itoa(n)
}
