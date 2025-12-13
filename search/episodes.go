package search

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gonwatch/common"

	// "fmt"
	"strings"
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

	urlBytes := []byte{104, 116, 116, 112, 115, 58, 47, 47, 97, 112, 105, 46, 116, 104, 101, 109, 111, 118, 105, 101, 100, 98, 46, 111, 114, 103, 47, 51, 47, 116, 118, 47}
	url := string(urlBytes) + url.QueryEscape(strconv.Itoa(int(tmdbid))) + "/season/" + strconv.Itoa(int(szn_numb)) + "?language=en-US"
	bearerBytes := []byte{66, 101, 97, 114, 101, 114, 32, 101, 121, 74, 104, 98, 71, 99, 105, 79, 105, 74, 73, 85, 122, 73, 49, 78, 105, 74, 57, 46, 101, 121, 74, 104, 100, 87, 81, 105, 79, 105, 73, 119, 77, 122, 77, 48, 77, 71, 73, 48, 90, 68, 104, 107, 79, 68, 103, 53, 78, 68, 77, 120, 77, 122, 73, 52, 89, 50, 69, 119, 79, 68, 81, 48, 89, 122, 73, 51, 90, 106, 65, 51, 90, 105, 73, 115, 73, 109, 53, 105, 90, 105, 73, 54, 77, 84, 99, 120, 77, 122, 73, 122, 78, 106, 73, 120, 77, 67, 52, 53, 78, 122, 99, 115, 73, 110, 78, 49, 89, 105, 73, 54, 73, 106, 89, 50, 77, 87, 82, 108, 79, 71, 89, 121, 78, 84, 73, 52, 89, 106, 74, 108, 77, 68, 69, 48, 89, 84, 78, 108, 78, 84, 100, 109, 89, 121, 73, 115, 73, 110, 78, 106, 98, 51, 66, 108, 99, 121, 73, 54, 87, 121, 74, 104, 99, 71, 108, 102, 99, 109, 86, 104, 90, 67, 74, 100, 76, 67, 74, 50, 90, 88, 74, 122, 97, 87, 57, 117, 73, 106, 111, 120, 102, 81, 46, 100, 49, 122, 95, 101, 55, 122, 54, 105, 118, 76, 84, 50, 65, 49, 115, 75, 45, 101, 95, 98, 75, 98, 76, 119, 108, 71, 82, 112, 83, 71, 55, 102, 80, 57, 74, 81, 73, 55, 115, 69, 97, 111}
	var bearer = string(bearerBytes)

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
	urlBytes := []byte{104, 116, 116, 112, 115, 58, 47, 47, 97, 112, 105, 46, 116, 104, 101, 109, 111, 118, 105, 101, 100, 98, 46, 111, 114, 103, 47, 51, 47, 116, 118, 47}
	url := string(urlBytes) + url.QueryEscape(strconv.Itoa(int(tmdbid)))
	bearerBytes := []byte{66, 101, 97, 114, 101, 114, 32, 101, 121, 74, 104, 98, 71, 99, 105, 79, 105, 74, 73, 85, 122, 73, 49, 78, 105, 74, 57, 46, 101, 121, 74, 104, 100, 87, 81, 105, 79, 105, 73, 119, 77, 122, 77, 48, 77, 71, 73, 48, 90, 68, 104, 107, 79, 68, 103, 53, 78, 68, 77, 120, 77, 122, 73, 52, 89, 50, 69, 119, 79, 68, 81, 48, 89, 122, 73, 51, 90, 106, 65, 51, 90, 105, 73, 115, 73, 109, 53, 105, 90, 105, 73, 54, 77, 84, 99, 120, 77, 122, 73, 122, 78, 106, 73, 120, 77, 67, 52, 53, 78, 122, 99, 115, 73, 110, 78, 49, 89, 105, 73, 54, 73, 106, 89, 50, 77, 87, 82, 108, 79, 71, 89, 121, 78, 84, 73, 52, 89, 106, 74, 108, 77, 68, 69, 48, 89, 84, 78, 108, 78, 84, 100, 109, 89, 121, 73, 115, 73, 110, 78, 106, 98, 51, 66, 108, 99, 121, 73, 54, 87, 121, 74, 104, 99, 71, 108, 102, 99, 109, 86, 104, 90, 67, 74, 100, 76, 67, 74, 50, 90, 88, 74, 122, 97, 87, 57, 117, 73, 106, 111, 120, 102, 81, 46, 100, 49, 122, 95, 101, 55, 122, 54, 105, 118, 76, 84, 50, 65, 49, 115, 75, 45, 101, 95, 98, 75, 98, 76, 119, 108, 71, 82, 112, 83, 71, 55, 102, 80, 57, 74, 81, 73, 55, 115, 69, 97, 111}
	var bearer = string(bearerBytes)

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

func GetAnimeEpisodeList(episodeList []string, filterValue string) []common.AnimeEpisodeTypeList {
	var s common.AnimeEpisodeTypeList
	myList := []common.AnimeEpisodeTypeList{}
	filteredValue := strings.Split(filterValue, "|")

	for _, item := range episodeList {
		s.EpisodeId    = item
		s.SeasonID     = filteredValue[0]
		s.AnimeName    = filteredValue[1]

		myList = append(myList, s)
	}

	return myList
}
