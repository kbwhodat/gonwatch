package models

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"gonwatch/history"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WatchedAwareDelegate struct {
	list.DefaultDelegate
	ParentTmdbID   int64
	ParentSeason   int
	ParentSeasonID string
}

func NewWatchedAwareDelegate() WatchedAwareDelegate {
	d := WatchedAwareDelegate{
		DefaultDelegate: list.NewDefaultDelegate(),
	}
	return d
}

func NewWatchedAwareDelegateWithContext(tmdbID int64, seasonID string) WatchedAwareDelegate {
	d := NewWatchedAwareDelegate()
	d.ParentTmdbID = tmdbID
	d.ParentSeasonID = seasonID
	return d
}

func (d WatchedAwareDelegate) isItemWatched(item list.Item) bool {
	listItem, ok := item.(ListItem)
	if !ok {
		return false
	}

	switch listItem.Type() {
	case "episode":
		return history.IsEpisodeWatched(listItem.TmdbID(), listItem.SznNumber(), listItem.ID())

	case "anime episodes":
		if d.ParentSeasonID == "" {
			return false
		}
		episodeNum, _ := strconv.ParseInt(listItem.EpString(), 10, 64)
		return history.IsAnimeEpisodeWatchedBySeasonID(d.ParentTmdbID, d.ParentSeasonID, episodeNum)

	case "vods":
		return history.IsMovieWatched(listItem.TmdbID())
	}

	return false
}

func (d WatchedAwareDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		title, desc string
		isSelected  = index == m.Index()
		isWatched   = d.isItemWatched(item)
	)

	if i, ok := item.(list.DefaultItem); ok {
		title = i.Title()
		desc = i.Description()
	} else {
		return
	}

	normalTitle := d.Styles.NormalTitle.Copy()
	normalDesc := d.Styles.NormalDesc.Copy()
	selectedTitle := d.Styles.SelectedTitle.Copy()
	selectedDesc := d.Styles.SelectedDesc.Copy()
	dimmedTitle := d.Styles.DimmedTitle.Copy()
	dimmedDesc := d.Styles.DimmedDesc.Copy()

	if isWatched {
		watchedColor := lipgloss.Color("241")
		normalTitle = normalTitle.Foreground(watchedColor)
		normalDesc = normalDesc.Foreground(watchedColor)
		selectedTitle = selectedTitle.Foreground(watchedColor)
		selectedDesc = selectedDesc.Foreground(watchedColor)
		dimmedTitle = dimmedTitle.Foreground(watchedColor)
		dimmedDesc = dimmedDesc.Foreground(watchedColor)
	}

	if m.Width() > 0 {
		textwidth := m.Width() - normalTitle.GetPaddingLeft() - normalTitle.GetPaddingRight()
		title = truncate(title, textwidth)
		desc = truncate(desc, textwidth)
	}

	if isSelected {
		title = selectedTitle.Render(title)
		desc = selectedDesc.Render(desc)
	} else if index == m.Index() {
		title = dimmedTitle.Render(title)
		desc = dimmedDesc.Render(desc)
	} else {
		title = normalTitle.Render(title)
		desc = normalDesc.Render(desc)
	}

	if d.ShowDescription {
		fmt.Fprintf(w, "%s\n%s", title, desc)
	} else {
		fmt.Fprintf(w, "%s", title)
	}
}

func (d WatchedAwareDelegate) Height() int {
	return d.DefaultDelegate.Height()
}

func (d WatchedAwareDelegate) Spacing() int {
	return d.DefaultDelegate.Spacing()
}

func (d WatchedAwareDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return d.DefaultDelegate.Update(msg, m)
}

func truncate(s string, width int) string {
	if width <= 0 {
		return s
	}
	if len(s) > width {
		if width > 3 {
			return strings.TrimSpace(s[:width-3]) + "..."
		}
		return s[:width]
	}
	return s
}
