package watch

func PlayStream(content string, stream_url string) []string {
	// stream_url = "https://embedsports.top/embed/alpha/los-angeles-rams-vs-tampa-bay-buccaneers/1"
	urls, subtitles := executePythonTask(content, 0, 0, 0, "placeholder", 0, 0, stream_url, "")
	if len(urls) == 0 {
		return urls

	} else {
		go func() {
			openMpv(urls, subtitles)
		}()
	}

	return urls
}
