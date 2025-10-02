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

func (i BubbleTeaEpisodeList) Title() string       { return i.EpisodeTitle + " (s" + strconv.Itoa(int(i.SeasonNumber)) + "e" + strconv.Itoa(int(i.EpisodeId)) + ")"}
func (i BubbleTeaEpisodeList) Description() string { return i.EpisodePlot }
func (i BubbleTeaEpisodeList) FilterValue() string { return i.EpisodeTitle + " (s" + strconv.Itoa(int(i.SeasonNumber)) + "e" + strconv.Itoa(int(i.EpisodeId)) + ")"}


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
