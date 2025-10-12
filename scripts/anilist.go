package scripts

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const endpoint = "https://graphql.anilist.co"

type AniListResponse struct {
	Data struct {
		Page struct {
			Media []struct {
				Episodes int `json:"episodes"`
				AnilistId int `json:"id"`
			}
		} `json:"Page"`
	} `json:"data"`
}

func GetEpisodesFromID(anilist_ids []int, flatepisode_number int64) (int, int) {

	query := `
		query ($ids: [Int]) {
		  Page {
		    media(id_in: $ids, type: ANIME, format_in: TV) {
		      episodes
		      id
		    }
		  }
		}
	`

	variables := map[string]any{
		"ids": anilist_ids,
	}

	payload := map[string]any{
		"query": query,
		"variables": variables,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("There was an issue with the graphql request")
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result AniListResponse

	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Fatal("there was an issue unmarshalling the data")
	}

	var current_episode int
	// var season_number int
	var anilist_id int
	remaining := int(flatepisode_number)
	for _, item := range result.Data.Page.Media {
		// season_number = (i + 1)
		if item.Episodes != 0 {
			if remaining > item.Episodes {
				log.Printf("%d > %d", remaining, item.Episodes)
				remaining -= item.Episodes
			} else {
				current_episode = remaining
				anilist_id = item.AnilistId
				break
			}
		} else {
			current_episode = int(flatepisode_number)
			anilist_id = item.AnilistId
			break
		}
	}

	return anilist_id, current_episode

}
