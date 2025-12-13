package models

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var TrendingChoiceList = []Choices{
	{choice: "movie"},
	{choice: "tv"},
}

func TrendingModel(m *Model) {
	var (
		titleStyle        = lipgloss.NewStyle().MarginLeft(2)
		paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
		helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	)

	listItems := make([]list.Item, len(TrendingChoiceList))
	for i, item := range TrendingChoiceList {
		listItems[i] = item
	}

	const defaultWidth = 20
	const listHeight = 14

	l := list.New(listItems, list.NewDefaultDelegate(), defaultWidth, listHeight)
	// l.Title = "Which trending content would you like to view?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m.List = l
	m.Mode = "trending"
	m.Cursor = 0

}
