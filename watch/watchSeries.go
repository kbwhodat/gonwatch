package watch

import (
	"log"
)

func PlayTv(content string, id int64, season_number int64, episode_number int64, title string, animeTitle string, skipSources []string) PlayResult {
	var anilist_id, anime_episode int

	result, err := executePythonTask(content, id, season_number, episode_number, title, anilist_id, anime_episode, "", animeTitle, skipSources)
	if err != nil {
		log.Println("Python task failed:", err)
		return PlayResult{
			SourcesTried: skipSources,
			TotalSources: result.TotalSources,
			UrlsFound:    false,
			Error:        err,
		}
	}

	newSourcesTried := skipSources
	if result.SourceUsed != "" {
		newSourcesTried = append(skipSources, result.SourceUsed)
	}

	if len(result.Urls) == 0 {
		return PlayResult{
			SourceUsed:   result.SourceUsed,
			SourcesTried: newSourcesTried,
			TotalSources: result.TotalSources,
			UrlsFound:    false,
			Success:      false,
		}
	}

	err = openMpv(result.Urls, result.Subtitles)
	if err != nil {
		log.Println("MPV playback failed:", err)
		return PlayResult{
			Urls:         result.Urls,
			Subtitles:    result.Subtitles,
			SourceUsed:   result.SourceUsed,
			SourcesTried: newSourcesTried,
			TotalSources: result.TotalSources,
			UrlsFound:    true,
			Success:      false,
			Error:        err,
		}
	}

	return PlayResult{
		Urls:         result.Urls,
		Subtitles:    result.Subtitles,
		SourceUsed:   result.SourceUsed,
		SourcesTried: newSourcesTried,
		TotalSources: result.TotalSources,
		UrlsFound:    true,
		Success:      true,
	}
}
