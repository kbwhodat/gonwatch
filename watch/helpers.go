package watch

import (
	// "bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	// "fmt"
	_ "embed"
	"log"
	"os/exec"
	"strconv"
	// "gonwatch/scripts"
)

func checkForErrors(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Result struct {
	Urls      []string `json:"urls"`
	Subtitles []string `json:"subtitles"`
}
func executePythonTask(content string, id int64, season_number int64, episode_number int64, title string, anilist_id int, anime_episode int, sports_url string) ([]string, []string) {

	// log.Println("going to run python script...")
	cmdArgs := []string{}

	cmdArgs = []string{"scripts/setcookies.py", content, strconv.FormatInt(id, 10), strconv.Itoa(int(season_number)), strconv.Itoa(int(episode_number)), title, strconv.Itoa(anilist_id), strconv.Itoa(anime_episode), sports_url}
	log.Println(cmdArgs)

	cmd := exec.Command("python", cmdArgs...)

	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(string(out))

	var result Result
	if err := json.Unmarshal(out, &result); err != nil {
		return []string{}, []string{}
	}

	if len(result.Subtitles) == 0 && (content != "anime" && content != "stream") {
		return result.Urls, GetSubtitles(int(id), content, int(season_number), int(episode_number))
	} else {
		return result.Urls, result.Subtitles
	}
}

func openMpv(urls []string, subtitles []string) {

	var dir string
	var err error

	var addSubtitleArgs []string
	if len(subtitles) > 0 {
		for _, subtitle := range subtitles {
			addSubtitleArgs = append(addSubtitleArgs, "--sub-file=" + subtitle)
		}
	}

	var cmdArgs []string
	dir, err = exec.LookPath("mpv")
	streamlink, err := exec.LookPath("streamlink")
	checkForErrors(err)

	for _, url := range urls {
		if strings.Contains(url, "shadowlandschronicles.com") {
			cmdArgs = []string{"--cache",  "--cache-secs=10", "--demuxer-readahead-secs=5", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", url}

		} else if strings.Contains(url, "_v7") {
			cmdArgs = []string{"--cache",  "--cache-secs=10", "--demuxer-readahead-secs=5", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", "--http-header-fields=Referer: https://rapid-cloud.co/", url}

		} else if strings.Contains(url, "lightningbolts.ru") {
			cmdArgs = []string{"--cache",  "--cache-secs=10", "--demuxer-readahead-secs=5", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", "--http-header-fields=Referer: https://vidsrc.cc/", url}

		} else if strings.Contains(url, "strmd.top") {
			log.Println(url)
			// Multiple flags needs to prevent the player (mpv) from disconnecting
			cmdArgs = []string{"--retry-open", "999", "--retry-streams", "999", "--stream-segment-attempts", "10", "--stream-segment-timeout", "10", "--player-continuous-http",
		 						"--http-no-ssl-verify", "--http-header", "Referer=https://embedsports.top/", url, "best", "-p", dir, "-a",
								"--network-timeout=60 --stream-lavf-o=reconnect=1 --stream-lavf-o=reconnect_streamed=1 --stream-lavf-o=reconnect_delay_max=5"}

		} else if strings.Contains(url, "gg.poocloud.in") {
			log.Println(url)
			// Multiple flags needs to prevent the player (mpv) from disconnecting
			cmdArgs = []string{"--retry-open", "999", "--retry-streams", "999", "--stream-segment-attempts", "10", "--stream-segment-timeout", "10", "--player-continuous-http",
		 						"--http-no-ssl-verify", "--http-header", "Referer=https://embedsports.top/", url, "best", "-p", dir, "-a",
								"--network-timeout=60 --stream-lavf-o=reconnect=1 --stream-lavf-o=reconnect_streamed=1 --stream-lavf-o=reconnect_delay_max=5"}

		} else if strings.Contains(url, "storm") {
			cmdArgs = []string{"--cache",  "--cache-secs=10", "--demuxer-readahead-secs=5", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", "--http-header-fields=Referer: https://vidlink.pro/", url}

		} else {
			cmdArgs = []string{"--cache",  "--cache-secs=10", "--demuxer-readahead-secs=5", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", url}
		}

		if len(addSubtitleArgs) > 0 {
			cmdArgs = append(cmdArgs, addSubtitleArgs...)
		}

		var newCmd *exec.Cmd
		if strings.Contains(url, "strmd.top") || strings.Contains(url, "gg.poocloud.in") {
			newCmd = exec.Command(streamlink, cmdArgs...)
			log.Println(newCmd)
		} else {
			newCmd = exec.Command(dir, cmdArgs...)
			log.Println(newCmd)
		}
		err := newCmd.Run()

		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				if exitError.ExitCode() != 0 {
					log.Println("Exit code is: ", exitError.ExitCode())
					log.Println("Got exit code, onto the next one")
					continue
				}
			} else {
				log.Fatalf("Failed to start mpv: %v", err)
			}
		} else {
			break
		}
	}
}

type SubtitlesResponse []struct {
	Url      string `json:"url"`
	Language string `json:"language"`
}
func GetSubtitles(tmdbid int, content string, season int, episode int) []string {

	subtitleList := []string{}
	var subtitle_url string
	if content == "tv" {
		subtitle_url = "https://sub.wyzie.ru/search?id="+strconv.Itoa(tmdbid)+"&season="+strconv.Itoa(season)+"&episode="+strconv.Itoa(episode)+"&format=srt"
	} else {
		subtitle_url = "https://sub.wyzie.ru/search?id="+strconv.Itoa(tmdbid)+"&format=srt"
	}

	req, _ := http.NewRequest("GET", subtitle_url, nil)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("unable to get the request")
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var subtitles SubtitlesResponse
	if err := json.Unmarshal(respBody, &subtitles); err != nil {
		log.Println("issue with unmarshalling. Probably not found")
	}

	for _, item := range subtitles {
		if item.Language == "en" {
			subtitleList = append(subtitleList, item.Url)
		}
	}

	return subtitleList
}
