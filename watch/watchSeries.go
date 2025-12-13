package watch

import (
	// "gonwatch/scripts"
	// "strings"
	// "log"
	// "reflect"
)

func PlayTv(content string, id int64, season_number int64, episode_number int64, title string, animeTitle string) []string {
	var anilist_id, anime_episode int
	// if strings.Contains(title, "JP"){
	// 	anilist_id, anime_episode = scripts.GetEpisodesFromID(scripts.GetMappings(id), episode_number)
	// }
	// log.Println("anilist id is: ", anilist_id)
	// log.Println("episode id is: ", anime_episode)

	// return []string{}

	urls, subtitles := executePythonTask(content, id, season_number, episode_number, title, anilist_id, anime_episode, "", animeTitle)
	// urls, subtitles := executePythonTask(content, id, season_number, episode_number, title, anilist_id, anime_episode)

	// url_type := reflect.TypeOf(urls)
	// log.Println(url_type)
	// switch url_type.Kind() {
	// 	case reflect.String:
	// 		log.Println("Hi i'm a string")
	// 	case reflect.Array:
	// 		log.Println("Hi I'm an arrary")
	// 	case reflect.Slice:
	// 		log.Println(len(urls))
	// }

	if len(urls) > 0 {
		go func() {
			openMpv(urls, subtitles)
		}()
	}

	return urls
}
