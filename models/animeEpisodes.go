package models

import (
	"os"
	"gonwatch/common"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"golang.org/x/term"
)

func AnimeEpisodesModel(m *Model, items []common.AnimeEpisodeTypeList) {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = BubbleTeaAnimeEpisodesList{AnimeEpisodeTypeList: item}
	}

	m.List = list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	m.List.Title = "Anime Episodes"

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
