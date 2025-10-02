package watch

func PlayMovie(content string, id int64) []string {
	urls, subtitles := executePythonTask(content, id, 0, 0)
	go func() {
		openMpv(urls, subtitles)
	}()
	return urls
}
