package models

import (
	"database/sql"
	"gonwatch/common"
	"strconv"
)

// This file allows you to format the text on bubbletea

// SERIES
type BubbleTeaSeriesList struct {
    common.StreamTypeList
}

func (e BubbleTeaSeriesList) Type() string {
    return "series"
}
func (e BubbleTeaSeriesList) ID() int64 {
    return e.StreamID
}
func (e BubbleTeaSeriesList) SznNumber() int {
    return 0
}
func (e BubbleTeaSeriesList) SznID() int {
    return 0
}
func (e BubbleTeaSeriesList) TmdbID() int64 {
    return e.StreamID
}
func (e BubbleTeaSeriesList) EpList() []string {
    return []string{}
}
func (e BubbleTeaSeriesList) EpString() string {
    return ""
}

func (i BubbleTeaSeriesList) Title() string {
    title := i.StreamTitle
    if i.StreamYear != "" {
    	title += " (" + i.StreamYear[0:4] + ")"
    }
    return title
}
func (i BubbleTeaSeriesList) Description() string { return i.StreamPlot }
func (i BubbleTeaSeriesList) FilterValue() string { return i.StreamTitle }



// SEASONS
type BubbleTeaSeasonList struct {
    common.SeasonsTypeList
}

func (e BubbleTeaSeasonList) Type() string {
    return "season"
}
func (e BubbleTeaSeasonList) ID() int64 {
    return e.SeriesID
}
func (e BubbleTeaSeasonList) SznNumber() int {
	season_number, _ := strconv.Atoi(e.SeasonNumber)
	return season_number
}
func (e BubbleTeaSeasonList) SznID() int {
	return 0
}
func (e BubbleTeaSeasonList) TmdbID() int64 {
    return e.SeriesID
}
func (e BubbleTeaSeasonList) EpList() []string {
    return []string{}
}
func (e BubbleTeaSeasonList) EpString() string {
    return ""
}

func (i BubbleTeaSeasonList) Title() string       { return i.SeasonTitle }
func (i BubbleTeaSeasonList) Description() string { return strconv.Itoa(int(i.EpisodeCount)) + " episodes"}
func (i BubbleTeaSeasonList) FilterValue() string { return i.SeasonTitle }



// EPISODES
type BubbleTeaEpisodeList struct {
    common.EpisodeTypeList
}

func (e BubbleTeaEpisodeList) Type() string {
    return "episode"
}
func (e BubbleTeaEpisodeList) ID() int64 {
    return e.EpisodeId
}
func (e BubbleTeaEpisodeList) SznNumber() int {
    return e.SeasonNumber
}
func (e BubbleTeaEpisodeList) SznID() int {
    return 0
}
func (e BubbleTeaEpisodeList) TmdbID() int64 {
    return e.EpisodeTmdbID
}
func (e BubbleTeaEpisodeList) EpList() []string {
    return []string{}
}
func (e BubbleTeaEpisodeList) EpString() string {
    return ""
}

func (i BubbleTeaEpisodeList) Title() string       { return i.EpisodeTitle + " (s" + strconv.Itoa(int(i.SeasonNumber)) + "e" + strconv.Itoa(int(i.EpisodeId)) + ")"}
func (i BubbleTeaEpisodeList) Description() string { return i.EpisodePlot }
func (i BubbleTeaEpisodeList) FilterValue() string { return i.EpisodeTitle + " (s" + strconv.Itoa(int(i.SeasonNumber)) + "e" + strconv.Itoa(int(i.EpisodeId)) + ")" + " " + i.Country}

// ANIME
type BubbleTeaAnimeList struct {
    common.AnimeTypeList
}

func (e BubbleTeaAnimeList) Type() string {
    return "anime"
}
func (e BubbleTeaAnimeList) ID() int64 {
    return e.AnimeID
}
func (e BubbleTeaAnimeList) SznNumber() int {
    return 0
}
func (e BubbleTeaAnimeList) SznID() int {
    return 0
}
func (e BubbleTeaAnimeList) TmdbID() int64 {
    return e.AnimeID
}
func (e BubbleTeaAnimeList) EpList() []string {
    return []string{}
}
func (e BubbleTeaAnimeList) EpString() string {
    return ""
}

func (i BubbleTeaAnimeList) Title() string {
    title := i.AnimeTitle
    if i.AnimeReleaseDate != "" {
    	title += " (" + i.AnimeReleaseDate[0:4] + ")"
    }
    return title
}
func (i BubbleTeaAnimeList) Description() string { return i.AnimePlot }
func (i BubbleTeaAnimeList) FilterValue() string { return i.AnimeTitle }

// SEASONS
type BubbleTeaAnimeSeasonList struct {
    common.SeasonsTypeList
}

func (e BubbleTeaAnimeSeasonList) Type() string {
    return "anime seasons"
}
func (e BubbleTeaAnimeSeasonList) ID() int64 {
    return 0
}
func (e BubbleTeaAnimeSeasonList) SznNumber() int {
	season_number, _ := strconv.Atoi(e.SeasonNumber)
	return season_number
}
func (e BubbleTeaAnimeSeasonList) SznID() int {
	return 0
}
func (e BubbleTeaAnimeSeasonList) TmdbID() int64 {
    return e.SeriesID
}
func (e BubbleTeaAnimeSeasonList) EpList() []string {
    return e.Episodes
}
func (e BubbleTeaAnimeSeasonList) EpString() string {
    return ""
}

func (i BubbleTeaAnimeSeasonList) Title()       string { return i.SeasonTitle }
func (i BubbleTeaAnimeSeasonList) Description() string { return i.SeasonPlot}
func (i BubbleTeaAnimeSeasonList) FilterValue() string { return i.SeasonID }

// EPISODES
type BubbleTeaAnimeEpisodesList struct {
    common.AnimeEpisodeTypeList
}

func (e BubbleTeaAnimeEpisodesList) Type() string {
    return "anime episodes"
}
func (e BubbleTeaAnimeEpisodesList) ID() int64 {
    return 0
}
func (e BubbleTeaAnimeEpisodesList) SznNumber() int {
	return 0
}
func (e BubbleTeaAnimeEpisodesList) SznID() int {
	return 0
}
func (e BubbleTeaAnimeEpisodesList) TmdbID() int64 {
    return 0
}
func (e BubbleTeaAnimeEpisodesList) EpList() []string {
    return e.EpisodeList
}
func (e BubbleTeaAnimeEpisodesList) EpString() string {
    return e.EpisodeId
}

func (i BubbleTeaAnimeEpisodesList) Title()       string { return "Episode: " + i.EpisodeId }
func (i BubbleTeaAnimeEpisodesList) Description() string { return ""}
func (i BubbleTeaAnimeEpisodesList) FilterValue() string { return i.SeasonID }

// VODS
type BubbleTeaVodsList struct {
    common.VodTypeList
}

func (e BubbleTeaVodsList) Type() string {
    return "vods"
}
func (e BubbleTeaVodsList) ID() int64 {
    return e.VodID
}
func (e BubbleTeaVodsList) SznNumber() int {
    return 0
}
func (e BubbleTeaVodsList) SznID() int {
    return 0
}
func (e BubbleTeaVodsList) TmdbID() int64 {
    return e.VodTmdbID
}
func (e BubbleTeaVodsList) EpList() []string {
    return []string{}
}
func (e BubbleTeaVodsList) EpString() string {
    return ""
}

func (i BubbleTeaVodsList) Title() string {
    title := i.VodTitle
    if i.VodReleaseDate != "" {
    	title += " (" + i.VodReleaseDate[0:4] + ")"
    }
    return title
}
func (i BubbleTeaVodsList) Description() string { return i.VodType }
func (i BubbleTeaVodsList) FilterValue() string { return i.VodTitle }


// LIVE
type BubbleTeaLiveList struct {
    common.LiveTypeList
}

func (e BubbleTeaLiveList) Type() string {
    return "live"
}
func (e BubbleTeaLiveList) ID() int {
    return e.LiveID
}
func (e BubbleTeaLiveList) SznNumber() int {
    return 0
}
func (e BubbleTeaLiveList) SznID() int {
    return 0
}
func (e BubbleTeaLiveList) TmdbID() sql.NullInt64 {
    return sql.NullInt64{
        Int64: 0,
        Valid: false, // The 0 here is irrelevant since Valid is false
    }
}

func (i BubbleTeaLiveList) Title() string       { return i.LiveTitle }
func (i BubbleTeaLiveList) Description() string { return i.LiveChannelName }
func (i BubbleTeaLiveList) FilterValue() string { return i.LiveTitle }


type Choices struct {
	choice string
}

func (c Choices) Title() string { return c.choice }
func (c Choices) Description() string { return "" }
func (c Choices) FilterValue() string { return c.choice }
