package search

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gonwatch/common"

	// "fmt"
	"io"
	"log"
	"net/url"
)

type TheEpisodesDBResponse struct {
	Episodes []struct {
		EpisodeId	    int64  `json:"show_id"`
		EpisodeName     string `json:"name"`
		EpisodeOverview string `json:"overview"`
		EpisodeNumber   int64  `json:"episode_number"`
		ReleaseDate	    string `json:"air_date"`
		SeasonNumber    int    `json:"season_number"`
		Runtime         int    `json:"runtime"`
	}
}
func GetEpisodes(tmdbid int64, szn_numb int) []common.EpisodeTypeList {

	url := "https://api.themoviedb.org/3/tv/" + url.QueryEscape(strconv.Itoa(int(tmdbid))) + "/season/" + strconv.Itoa(int(szn_numb)) + "?language=en-US"
	var bearer = "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIwMzM0MGI0ZDhkODg5NDMxMzI4Y2EwODQ0YzI3ZjA3ZiIsIm5iZiI6MTcxMzIzNjIxMC45NzcsInN1YiI6IjY2MWRlOGYyNTI4YjJlMDE0YTNlNTdmYyIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.d1z_e7z6ivLT2A1sK-e_bKbLwlGRpSG7fP9JQI7sEao"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result TheEpisodesDBResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Println("cannot marshall the json")
	}

	var s common.EpisodeTypeList

	myList := []common.EpisodeTypeList{}
	for _, item := range result.Episodes {
		s.EpisodeTitle       = item.EpisodeName
		s.EpisodePlot        = item.EpisodeOverview
		s.EpisodeId          = item.EpisodeNumber
		s.EpisodeReleaseDate = item.ReleaseDate
		s.EpisodeTmdbID      = item.EpisodeId
		s.SeasonNumber       = item.SeasonNumber
		s.Runtime            = item.Runtime

		myList = append(myList, s)
	}

	return myList
}
