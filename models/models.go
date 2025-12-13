package models

import (
	"gonwatch/search"
	"gonwatch/update"
	"gonwatch/view"
	"gonwatch/watch"
	"log"
	"slices"
	"strconv"
	"strings"

	// "log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("204")).Background(lipgloss.Color("235"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

var sportGenres = []string{"basketball", "football", "american-football", "hockey", "baseball", "motor-sports", "fight", "tennis", "rugby", "golf", "billiards", "afl", "darts", "cricket", "other"}

type Model struct {
	TextInput 	   textinput.Model
	List 		   list.Model
	Mode		   string
	Err       	   error
	PreviousStates []ModelStateSnapshot
	Cursor 		   int
	Choice 		   Choices
	Altscreen 	   bool
	Id			   int
	TmdbID		   int
	SelectedItem   string

	spinner        spinner.Model
	loading        bool
	loadingLabel   string
}

type ListItem interface {
	ID()            int64
	TmdbID()        int64
	Type() 	        string
	SznNumber()     int
	SznID() 	    int
	EpList() 	    []string
	EpString() 	    string
	SportName()     string
	SportId() 	    string
	OriginCountry() string
}

type SportsSources interface {
	Sources() []string
}

func (m *Model) Init() tea.Cmd {
	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Dot
	return m.spinner.Tick
}

type linkFetchedMsg struct {
	found bool
}

func fetchEpisodeCmd(item ListItem, m *Model) tea.Cmd {
	return func() tea.Msg {

		selectedItem, _ := m.List.SelectedItem().(ListItem)
		// var episode int64
		var ok bool
		if selectedItem.Type() == "anime episodes" {
			episode_number, _ := strconv.Atoi(selectedItem.EpString())
			filteredValue := strings.Split(m.List.SelectedItem().FilterValue(), "|")
			ok = len(watch.PlayTv("anime", item.TmdbID(), int64(item.SznNumber()), int64(episode_number), filteredValue[0], filteredValue[1])[0]) > 0
		} else {
			ok = len(watch.PlayTv("tv", item.TmdbID(), int64(item.SznNumber()), item.ID(), m.List.SelectedItem().FilterValue(), "")) > 0
		}

		return linkFetchedMsg{found: ok}
	}
}
func fetchStreamCmd(url string) tea.Cmd {

	return func() tea.Msg {
		ok := len(watch.PlayStream("stream", url)) > 0
		// ok := false
		return linkFetchedMsg{found: ok}
	}
}
func fetchMovieCmd(item ListItem) tea.Cmd {

	return func() tea.Msg {
		ok := len(watch.PlayMovie("movie", item.ID())) > 0
		return linkFetchedMsg{found: ok}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case linkFetchedMsg:
		m.loading = false
		if msg.found {
			m.Altscreen = true
			m.Mode = "fullscreen"
			return m, tea.EnterAltScreen
		} else {
			m.Altscreen = true
			m.Mode = "linknotfoundscreen"
			return m, tea.EnterAltScreen
		}
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+u":
			m.Altscreen = true
			m.Mode = "fullscreen"
			m.SelectedItem = "Please wait, updating database with new media..."

		case "ctrl+c", "esc":
			return m, tea.Quit

		case "left":
			m := m.restorePreviousState()
			return m, nil

		case "down", "j":
			m.Cursor++
			if m.Mode == "trending" {
				if m.Cursor >= len(TrendingChoiceList) {
					m.Cursor = 0
				}
			} else {
				if m.Cursor >= len(choiceList) {
					m.Cursor = 0
				}
			}

		case "up", "k":
			m.Cursor--
			if m.Mode == "trending" {
				if m.Cursor < 0 {
					m.Cursor = len(TrendingChoiceList) - 1
				}
			} else {
				if m.Cursor < 0 {
					m.Cursor = len(choiceList) - 1
				}
			}

		case "enter", "right":

			m.saveCurrentState()
			if m.Mode == "list" {
				selectedItem, ok := m.List.SelectedItem().(ListItem)
				if ok {
					switch selectedItem.Type() {

					case "movie":
						log.Println("hello there kind sir")

					case "streams":
						m.loading = true
						m.loadingLabel = "Fetching content…"
						m.Mode = "loading"
						// log.Println(m.List.SelectedItem().FilterValue())
						// log.Println(selectedItem.SportName()) //admin
						url := search.GetStreamLink(selectedItem.SportName(), m.List.SelectedItem().FilterValue())
						return m, tea.Batch(
							m.spinner.Tick,
							fetchStreamCmd(url),
						)

					case "sports":
						if slices.Contains(sportGenres, m.List.SelectedItem().FilterValue()) {
							matches := search.ListSportMatches(m.List.SelectedItem().FilterValue())
							// log.Println(matches)
							MatchesModel(m, matches)
						} else {
							sel := m.List.SelectedItem()
							match, ok := sel.(SportsSources)
					        if !ok {
					            log.Printf("unexpected item type: %T\n", sel)
					            break
					        }

							streamlist := search.ListStreams(match.Sources())
							MatchesModel(m, streamlist)

						}

					case "series":
						seasonList := search.GetSeasons(selectedItem.ID())
						SeasonModel(m, seasonList)

					case "trending":
						log.Println("in the trending ting!")

					case "season":
						episodeList := search.GetEpisodes(selectedItem.TmdbID(), selectedItem.SznNumber())
						EpisodeModel(m, episodeList)

					case "anime":
						seasonList := search.GetAnimeSeasons(selectedItem.ID(), m.List.SelectedItem().FilterValue())
						AnimeSeasonModel(m, seasonList)

					case "anime seasons":
						episodeList := search.GetAnimeEpisodeList(selectedItem.EpList(), m.List.SelectedItem().FilterValue())
						AnimeEpisodesModel(m, episodeList)

					case "anime episodes":
						m.loading = true
						m.loadingLabel = "Fetching content…"
						m.Mode = "loading"
						return m, tea.Batch(
							m.spinner.Tick,
							fetchEpisodeCmd(selectedItem, m),
						)

					case "episode":
						m.loading = true
						m.loadingLabel = "Fetching content…"
						m.Mode = "loading"
						return m, tea.Batch(
							m.spinner.Tick,
							fetchEpisodeCmd(selectedItem, m),
						)

					case "vods":
						m.loading = true
						m.loadingLabel = "Fetching content…"
						m.Mode = "loading"
						return m, tea.Batch(
							m.spinner.Tick,
							fetchMovieCmd(selectedItem),
						)
					}
				}
				return m, nil
			}

			if m.Mode == "input" {
				selectedItem := m.Choice.choice
				switch selectedItem {
					case "movies":
						resultList := update.InputUpdateMsgVods(m.TextInput)
						VodModel(m, resultList)
					case "series":
						resultList := update.InputUpdateMsgSeries(m.TextInput)
						SeriesModel(m, resultList)
					case "anime":
						resultList := update.InputUpdateMsgAnime(m.TextInput)
						AnimeModel(m, resultList)
				}
			}

			if m.Mode == "fullscreen" {
				m.Mode = "fullscreen"
			}

			if m.Mode == "trending" {
				selected := TrendingChoiceList[m.Cursor]
				if selected.FilterValue() == "movie" {
					moviesList := search.GetTrendingMovies()
					VodModel(m, moviesList)
				} else {
					tvList := search.GetTrendingTv()
					SeriesModel(m, tvList)
				}
			}

			if m.Mode == "select" {
				m.Choice = choiceList[m.Cursor]
				if m.Choice.FilterValue() == "sports" {
					genreList := search.ListSports()
					SportsModel(m, genreList)

				} else if m.Choice.FilterValue() == "trending" {
					TrendingModel(m)

				} else {
					InputModel(m)
				}
			}
		}
	}

	switch m.Mode {
	case "input":
		model, cmd := update.InputUpdate(m.TextInput, msg)
		m.TextInput = model
		return m, cmd

	case "list":
		model, cmd := update.ListUpdate(m.List, msg)
		m.List = model
		return m, cmd

	case "select":
		model, cmd := update.SelectUpdate(m.List, msg)
		m.List = model
		return m, cmd

	case "loading":
		var c tea.Cmd
		m.spinner, c = m.spinner.Update(msg)
		return m, c

	}

	return m, cmd
}

func (m Model) View() string {
	switch m.Mode {
		case "fullscreen":
		filteredValue := strings.Split(m.List.SelectedItem().FilterValue(), "|")
		if len(filteredValue) > 1 {
			return view.FullscreenView(
				keywordStyle.Render(filteredValue[1]),
				helpStyle.Render("\n\n\nleft/h: go back • q: exit/quit\n"),
			)
		} else {
			return view.FullscreenView(
				keywordStyle.Render(m.List.SelectedItem().FilterValue()),
				helpStyle.Render("\n\n\nleft/h: go back • q: exit/quit\n"),
			)
		}

	case "linknotfoundscreen":
		return view.LinkNotFoundView(
			keywordStyle.Render(m.List.SelectedItem().FilterValue()),
			helpStyle.Render("\n\n\nleft/h: go back • q: exit/quit\n"),
		)

	case "input":
		return view.InputView(m.TextInput.View())

	case "list":
		return view.ListView(m.List.View())

	case "select":
		return view.SelectView(m.Cursor)

	case "trending":
		return view.TrendingSelectView(m.Cursor)

	case "loading":
		return "\n\n  " + m.spinner.View() + " " + m.loadingLabel
	}

	return "unsupported"
}
