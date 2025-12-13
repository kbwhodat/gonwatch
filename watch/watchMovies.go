package watch

func PlayMovie(content string, id int64) []string {
	urls, subtitles := executePythonTask(content, id, 0, 0, "placeholder", 0, 0, "", "")
	if len(urls) == 0 {
		return urls

	} else {
		go func() {
			openMpv(urls, subtitles)
		}()
	}

	return urls
}
