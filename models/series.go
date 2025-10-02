package models

import (
	"gonwatch/common"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/key"
	"os"
	"golang.org/x/term"
)

func SeriesModel(m *Model, items []common.StreamTypeList) {

	listItems := make([]list.Item, len(items))

	for i, item := range items {
		listItems[i] = BubbleTeaSeriesList{StreamTypeList: item}
	}

	m.List = list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	m.List.Title = "Results"

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
