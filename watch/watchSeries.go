package watch

import (
	// "gonwatch/scripts"
	// "strings"
)

func PlayTv(content string, id int64, season_number int64, episode_number int64, title string) []string {
	var anilist_id, anime_episode int
	// if strings.Contains(title, "JP"){
	// 	anilist_id, anime_episode = scripts.GetEpisodesFromID(scripts.GetMappings(id), episode_number)
	// }
	// log.Println("anilist id is: ", anilist_id)
	// log.Println("episode id is: ", anime_episode)

	// return []string{}

	urls, subtitles := executePythonTask(content, id, season_number, episode_number, title, anilist_id, anime_episode)
	// urls, subtitles := executePythonTask(content, id, season_number, episode_number, title, anilist_id, anime_episode)

	go func() {
		openMpv(urls, subtitles)
	}()
	return urls
}
