package common

type StreamTypeList struct {
    StreamTitle string
    StreamID    int64
    StreamPlot  string
    StreamYear  string
    StreamType  string
}

type SeasonsTypeList struct {
    SeasonTitle       string
    SeasonID          string
    SeasonNumber      string
    SeasonPlot        string
    EpisodeCount      int64
    SeasonReleaseDate string
    SeriesID  	      int64
    Episodes  	      []string
}

type EpisodeTypeList struct {
    EpisodeTitle       string
    SeasonID           string
    EpisodePlot        string
    EpisodeId          int64
    EpisodeTmdbID      int64
    EpisodeReleaseDate string
    SeasonNumber       int
    Country            string
}

type AnimeEpisodeTypeList struct {
    EpisodeTitle       string
    SeasonID           string
    EpisodePlot        string
    EpisodeId          string
    EpisodeTmdbID      int64
    EpisodeReleaseDate string
    SeasonNumber       int
    Country            string
    EpisodeList  	   []string
}

type AnimeTypeList struct {
	AnimeTitle       string
    AnimeID          int64
    AnimePlot        string
    AnimeTmdbID      int64
    AnimeReleaseDate string
    AnimeCountry     string
}

type SportsGenreTypeList struct {
	SportsGenreName       string
    SportsGenreID         string
    SportsType            string
    SportSources []struct {
	    SportsSourceName  string
	    SportsSourceId    string
    }
}

type VodTypeList struct {
	VodTitle       string
    VodID          int64
    VodType        string
    VodTmdbID      int64
    VodReleaseDate string
}

type LiveTypeList struct {
    LiveTitle 			string
    LiveID    			int
    LiveChannelName string
    LiveType  			string
}
