package update

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"gonwatch/search"
	"gonwatch/common"
)

func InputUpdate(m textinput.Model, msg tea.Msg) (textinput.Model, tea.Cmd) {

	m, cmd := m.Update(msg)

	return m, cmd
}

func InputUpdateMsgSeries(m textinput.Model) []common.StreamTypeList {
	resultList := search.GetSeries(m.Value())
	return resultList
}

func InputUpdateMsgAnime(m textinput.Model) []common.AnimeTypeList {
	resultList := search.GetAnime(m.Value())
	return resultList
}

func InputUpdateMsgVods(m textinput.Model) []common.VodTypeList {
	resultList := search.GetMovies(m.Value())
	return resultList
}
