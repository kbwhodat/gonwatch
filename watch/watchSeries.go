package watch

func PlayTv(content string, id int64, season_number int64, episode_number int64) []string {
	urls, subtitles := executePythonTask(content, id, season_number, episode_number)
	go func() {
		openMpv(urls, subtitles)
	}()
	return urls
}
