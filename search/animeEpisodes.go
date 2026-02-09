package search

//http://ky-iptv.com/player_api.php?username=walkerrodney216%40gmail.com&password=Stream1958&action=get_vod_info&vod_id=819792

import (
	"encoding/json"
	"net/http"
	"strings"

	"gonwatch/common"

	// "fmt"
	"io"
	"log"
	"net/url"
)

type TheAnimeEpisodesDBResponse struct {
	Result []struct {
		AnimeId     string `json:"id"`
		EnglishName string `json:"englishName"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
}

func GetAnimeEpisodes(tmdbid int64, query string) []common.SeasonsTypeList {

	urlBytes := []byte{104, 116, 116, 112, 115, 58, 47, 47, 104, 101, 97, 118, 101, 110, 115, 99, 97, 112, 101, 46, 118, 101, 114, 99, 101, 108, 46, 97, 112, 112, 47, 97, 112, 105, 47, 97, 110, 105, 109, 101, 47, 115, 101, 97, 114, 99, 104, 47}
	url := string(urlBytes) + url.QueryEscape(query)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

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

	var result TheAnimeSeasonsDBResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Println("cannot marshall the json")
	}

	var s common.SeasonsTypeList

	myList := []common.SeasonsTypeList{}
	for _, item := range result.Result {
		if item.EnglishName != "" {
			lengthOfString := len(query)
			if lengthOfString <= len(item.EnglishName) {
				if strings.EqualFold(item.EnglishName[0:lengthOfString], strings.ToLower(query)) {
					s.SeriesID = tmdbid
					// s.EpisodeCount      = item.EpisodeCount
					s.SeasonTitle = item.EnglishName
					s.SeasonPlot = item.Description
					// s.SeasonNumber      = strconv.Itoa(int(item.SeasonNumber))
					s.SeasonID = item.AnimeId

					myList = append(myList, s)
				}
			}
		}
	}

	return myList
}
