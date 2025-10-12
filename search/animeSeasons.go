package search

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"gonwatch/common"

	// "fmt"
	"io"
	"log"
	"net/url"

	_ "github.com/marcboeker/go-duckdb"
)

type TheAnimeSeasonsDBResponse struct {
	Result []struct {
		AnimeId	    string  `json:"id"`
		EnglishName string `json:"englishName"`
		Description string `json:"description"`
		Status      string  `json:"status"`
		Episodes    struct {
			Sub []string `json:"sub"`
		} `json:"episodes"`
	}
}
func GetAnimeSeasons(tmdbid int64, query string) []common.SeasonsTypeList {

	url := "https://heavenscape.vercel.app/api/anime/search/" + url.QueryEscape(query)

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
			re := regexp.MustCompile(`[^a-z0-9 ]+`)
			englishName := re.ReplaceAllString(strings.ToLower(query), " ")
			if lengthOfString <= len(englishName) {
				if strings.EqualFold(item.EnglishName[0:lengthOfString], englishName) || strings.EqualFold(item.EnglishName[0:lengthOfString], strings.ToLower(query)) {
					s.SeriesID          = tmdbid
					// s.EpisodeCount      = item.EpisodeCount
					s.SeasonTitle       = item.EnglishName
					s.SeasonPlot        = item.Description
					// s.SeasonNumber      = strconv.Itoa(int(item.SeasonNumber))
					s.SeasonID          = item.AnimeId
					s.Episodes          = item.Episodes.Sub

					myList = append(myList, s)
				}
			}
		}
	}

	return myList
}
