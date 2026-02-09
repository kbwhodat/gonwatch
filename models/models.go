package models

import (
	"fmt"
	"gonwatch/history"
	"gonwatch/search"
	"gonwatch/update"
	"gonwatch/view"
	"gonwatch/watch"
	"log"
	"slices"
	"strconv"
	"strings"

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
	TextInput      textinput.Model
	List           list.Model
	Mode           string
	Err            error
	PreviousStates []ModelStateSnapshot
	Cursor         int
	Choice         Choices
	Altscreen      bool
	Id             int
	TmdbID         int
	SelectedItem   string

	spinner      spinner.Model
	loading      bool
	loadingLabel string

	playingItem     ListItem
	playingTmdbID   int64
	playingSeason   int
	playingSeasonID string

	sourcesTried  []string
	currentSource int
	totalSources  int
	streamUrl     string
}

type ListItem interface {
	ID() int64
	TmdbID() int64
	Type() string
	SznNumber() int
	SznID() int
	EpList() []string
	EpString() string
	SportName() string
	SportId() string
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
	found        bool
	sourcesTried []string
	totalSources int
}

type playbackRetryMsg struct {
	sourcesTried []string
	totalSources int
}

func fetchEpisodeCmd(item ListItem, m *Model, skipSources []string) tea.Cmd {
	return func() tea.Msg {
		selectedItem, _ := m.List.SelectedItem().(ListItem)

		var result watch.PlayResult
		if selectedItem.Type() == "anime episodes" {
			episode_number, _ := strconv.Atoi(selectedItem.EpString())
			filteredValue := strings.Split(m.List.SelectedItem().FilterValue(), "|")
			result = watch.PlayTv("anime", item.TmdbID(), int64(item.SznNumber()), int64(episode_number), filteredValue[0], filteredValue[1], skipSources)
		} else {
			result = watch.PlayTv("tv", item.TmdbID(), int64(item.SznNumber()), item.ID(), m.List.SelectedItem().FilterValue(), "", skipSources)
		}

		if result.Error != nil || !result.Success {
			if result.UrlsFound && len(result.SourcesTried) < result.TotalSources {
				return playbackRetryMsg{
					sourcesTried: result.SourcesTried,
					totalSources: result.TotalSources,
				}
			}
			return linkFetchedMsg{
				found:        false,
				sourcesTried: result.SourcesTried,
				totalSources: result.TotalSources,
			}
		}

		return linkFetchedMsg{
			found:        true,
			sourcesTried: result.SourcesTried,
			totalSources: result.TotalSources,
		}
	}
}

func fetchStreamCmd(url string, skipSources []string) tea.Cmd {
	return func() tea.Msg {
		result := watch.PlayStream("stream", url, skipSources)

		if result.Error != nil || !result.Success {
			if result.UrlsFound && len(result.SourcesTried) < result.TotalSources {
				return playbackRetryMsg{
					sourcesTried: result.SourcesTried,
					totalSources: result.TotalSources,
				}
			}
			return linkFetchedMsg{
				found:        false,
				sourcesTried: result.SourcesTried,
				totalSources: result.TotalSources,
			}
		}

		return linkFetchedMsg{
			found:        true,
			sourcesTried: result.SourcesTried,
			totalSources: result.TotalSources,
		}
	}
}

func fetchMovieCmd(item ListItem, skipSources []string) tea.Cmd {
	return func() tea.Msg {
		result := watch.PlayMovie("movie", item.ID(), skipSources)

		if result.Error != nil || !result.Success {
			if result.UrlsFound && len(result.SourcesTried) < result.TotalSources {
				return playbackRetryMsg{
					sourcesTried: result.SourcesTried,
					totalSources: result.TotalSources,
				}
			}
			return linkFetchedMsg{
				found:        false,
				sourcesTried: result.SourcesTried,
				totalSources: result.TotalSources,
			}
		}

		return linkFetchedMsg{
			found:        true,
			sourcesTried: result.SourcesTried,
			totalSources: result.TotalSources,
		}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case playbackRetryMsg:
		m.sourcesTried = msg.sourcesTried
		m.totalSources = msg.totalSources
		m.currentSource = len(msg.sourcesTried) + 1
		m.loadingLabel = fmt.Sprintf("Trying source %d of %d...", m.currentSource, m.totalSources)

		if m.playingItem != nil {
			switch m.playingItem.Type() {
			case "episode", "anime episodes":
				return m, tea.Batch(
					m.spinner.Tick,
					fetchEpisodeCmd(m.playingItem, m, msg.sourcesTried),
				)
			case "vods":
				return m, tea.Batch(
					m.spinner.Tick,
					fetchMovieCmd(m.playingItem, msg.sourcesTried),
				)
			case "streams":
				return m, tea.Batch(
					m.spinner.Tick,
					fetchStreamCmd(m.streamUrl, msg.sourcesTried),
				)
			}
		}
		return m, nil

	case linkFetchedMsg:
		m.loading = false
		m.sourcesTried = nil
		m.currentSource = 0
		m.totalSources = 0

		if msg.found {
			if m.playingItem != nil {
				switch m.playingItem.Type() {
				case "episode":
					history.MarkWatched(
						"episode",
						m.playingItem.TmdbID(),
						m.playingItem.SznNumber(),
						m.playingItem.ID(),
						m.List.SelectedItem().FilterValue(),
					)
				case "anime episodes":
					episodeNum, _ := strconv.ParseInt(m.playingItem.EpString(), 10, 64)
					history.MarkAnimeEpisodeWatchedBySeasonID(
						m.playingTmdbID,
						m.playingSeasonID,
						episodeNum,
						m.List.SelectedItem().FilterValue(),
					)
				case "vods":
					history.MarkWatched(
						"movie",
						m.playingItem.TmdbID(),
						0,
						0,
						m.List.SelectedItem().FilterValue(),
					)
				}
			}
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
						log.Println("")

					case "streams":
						m.loading = true
						m.loadingLabel = "Fetching content..."
						m.Mode = "loading"
						m.playingItem = selectedItem
						m.sourcesTried = nil
						url := search.GetStreamLink(selectedItem.SportName(), m.List.SelectedItem().FilterValue())
						m.streamUrl = url
						return m, tea.Batch(
							m.spinner.Tick,
							fetchStreamCmd(url, nil),
						)

					case "sports":
						if slices.Contains(sportGenres, m.List.SelectedItem().FilterValue()) {
							matches := search.ListSportMatches(m.List.SelectedItem().FilterValue())
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
						seasonFilterValue := strings.Split(m.List.SelectedItem().FilterValue(), "|")
						seasonID := ""
						if len(seasonFilterValue) > 0 {
							seasonID = seasonFilterValue[0]
						}
						AnimeEpisodesModel(m, episodeList, selectedItem.TmdbID(), seasonID)

					case "anime episodes":
						m.loading = true
						m.loadingLabel = "Fetching content..."
						m.Mode = "loading"
						m.playingItem = selectedItem
						m.sourcesTried = nil
						return m, tea.Batch(
							m.spinner.Tick,
							fetchEpisodeCmd(selectedItem, m, nil),
						)

					case "episode":
						m.loading = true
						m.loadingLabel = "Fetching content..."
						m.Mode = "loading"
						m.playingItem = selectedItem
						m.sourcesTried = nil
						return m, tea.Batch(
							m.spinner.Tick,
							fetchEpisodeCmd(selectedItem, m, nil),
						)

					case "vods":
						m.loading = true
						m.loadingLabel = "Fetching content..."
						m.Mode = "loading"
						m.playingItem = selectedItem
						m.sourcesTried = nil
						return m, tea.Batch(
							m.spinner.Tick,
							fetchMovieCmd(selectedItem, nil),
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
