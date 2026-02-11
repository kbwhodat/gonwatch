package models

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var choiceList = []Choices{
	{choice: "recently watched"},
	{choice: "trending"},
	{choice: "movies"},
	{choice: "series"},
	{choice: "anime"},
	{choice: "sports"},
}

func ChoiceModel() *Model {
	var (
		titleStyle      = lipgloss.NewStyle().MarginLeft(2)
		paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
		helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	)

	listItems := make([]list.Item, len(choiceList))
	for i, item := range choiceList {
		listItems[i] = item
	}

	const defaultWidth = 20
	const listHeight = 14

	l := list.New(listItems, list.NewDefaultDelegate(), defaultWidth, listHeight)
	l.Title = "What do you want to do today?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := &Model{ // Use & to return a pointer
		List: l,
		Mode: "select",
	}

	return m
}
