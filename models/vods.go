package models

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"golang.org/x/term"
	"gonwatch/common"
	"os"
)

func VodModel(m *Model, items []common.VodTypeList) {

	listItems := make([]list.Item, len(items))

	for i, item := range items {
		listItems[i] = BubbleTeaVodsList{VodTypeList: item}
	}

	delegate := NewWatchedAwareDelegate()
	m.List = list.New(listItems, delegate, 0, 0)
	m.List.Title = "Movies"

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
