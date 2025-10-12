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
	Networks []struct {
		OriginCountry string `json:"origin_country"`
	} `json:"networks"`
}

type episodeResult struct {
	Seasons []struct {
		EpisodeCount int64 `json:"episode_count"`
		SeasonNumber int64  `json:"season_number"`
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
		log.Println("cannot unmarshall the json")
	}

	var s common.EpisodeTypeList


	var country string
	for _, item := range result.Networks {
		country = item.OriginCountry
	}

	myList := []common.EpisodeTypeList{}
	for _, item := range result.Episodes {
		s.EpisodeTitle       = item.EpisodeName
		s.EpisodePlot        = item.EpisodeOverview
		s.EpisodeId          = item.EpisodeNumber
		s.EpisodeReleaseDate = item.ReleaseDate
		s.EpisodeTmdbID      = item.EpisodeId
		s.SeasonNumber       = item.SeasonNumber
		s.Country            = country

		myList = append(myList, s)
	}

	return myList
}

func GetAbsoluteEpisode(tmdbid int64, season int64, episode int64) int64 {
	url := "https://api.themoviedb.org/3/tv/" + url.QueryEscape(strconv.Itoa(int(tmdbid)))
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

	var result episodeResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Println("unable to unmarshalll this")
	}

	var flat_episode int64
	for _, item := range result.Seasons {
		if item.SeasonNumber != 0 {
			if item.SeasonNumber < season {
				flat_episode += item.EpisodeCount
			}
		}
	}

	if episode > flat_episode {
		return episode
	}
	flat_episode += episode
	return flat_episode

}

func GetAnimeEpisodeList(episodeList []string, seasonId string) []common.AnimeEpisodeTypeList {
	var s common.AnimeEpisodeTypeList
	myList := []common.AnimeEpisodeTypeList{}

	for _, item := range episodeList {
		s.EpisodeId = item
		s.SeasonID = seasonId

		myList = append(myList, s)
	}


	return myList

}
