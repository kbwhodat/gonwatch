package scripts

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
)

type IDMappings []struct {
	AnilistID     int `json:"anilist_id,omitempty"`
	ThemoviedbID  any `json:"themoviedb_id,omitempty"`
	AnimePlanetID any `json:"anime_planet_id,omitempty"`
}

func getAniListId(mappings IDMappings, tmdbid int64) []int {

	var anilistids []int

	for _, item := range mappings {
		switch v := item.ThemoviedbID.(type) {
		case float64:
			if int(v) == int(tmdbid) {
				anilistids = append(anilistids, int(item.AnilistID))
			}
		case string:
			if v == strconv.FormatInt(tmdbid, 10) {
				anilistids = append(anilistids, int(item.AnilistID))
			}
		}
	}
	return anilistids
}

func GetMappings(tmdbid int64) []int {
	gob.Register(IDMappings{})

	c := cache.New(168*time.Hour, 168*time.Hour)

	file, _ := os.Open("/tmp/cache")
	defer file.Close()

	if err := c.Load(file); err != nil {
		log.Println(err)
	}

	if mappings, found := c.Get("id_mappings"); found {
		val := mappings.(IDMappings)

		return getAniListId(val, tmdbid)

	} else {

		url := "https://raw.githubusercontent.com/Fribb/anime-lists/refs/heads/master/anime-list-full.json"

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
			log.Fatal("Nothing was returned from this...")
		}

		var result IDMappings
		if err := json.Unmarshal(respBody, &result); err != nil {
			log.Println("An issue occured with the unmarshalling of the response: ", err)
		}

		c.Set("id_mappings", result, cache.DefaultExpiration)

		file, err := os.Create("/tmp/cache")
		if err != nil {
			log.Println("unable to write to file")
		}
		defer file.Close()
		if err := c.Save(file); err != nil {
			log.Fatal(err)
		}

		return getAniListId(result, tmdbid)
	}
}
