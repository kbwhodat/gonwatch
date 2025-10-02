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
    // SeasonID          int64
    SeasonNumber    string
    EpisodeCount      int64
    SeasonReleaseDate string
    SeriesID  	      int64
}

type EpisodeTypeList struct {
    EpisodeTitle       string
    SeasonID           string
    EpisodePlot        string
    EpisodeId          int64
    EpisodeTmdbID      int64
    EpisodeReleaseDate string
    SeasonNumber       int
    Runtime            int
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
