package watch

import (
	// "bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
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

func executePythonTask(content string, id int64, season_number int64, episode_number int64, title string, anilist_id int, anime_episode int, sports_url string, anime_title string) ([]string, []string) {

	cmdArgs := []string{}

	cmdArgs = []string{"scripts/setcookies.py", content, strconv.FormatInt(id, 10), strconv.Itoa(int(season_number)), strconv.Itoa(int(episode_number)), title, strconv.Itoa(anilist_id), strconv.Itoa(anime_episode), sports_url, anime_title}
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

	// log.Println(urls)

	var mpv string
	var err error

	var addSubtitleArgs []string
	if len(subtitles) > 0 {
		for _, subtitle := range subtitles {
			addSubtitleArgs = append(addSubtitleArgs, "--sub-file="+subtitle)
		}
	}

	var cmdArgs []string
	mpv, err = exec.LookPath("mpv")
	streamlink, err := exec.LookPath("streamlink")
	// xdg_open, err := exec.LookPath("xdg-open")
	checkForErrors(err)

	for _, host := range urls {
		if strings.Contains(host, "_v7") {
			cmdArgs = []string{"--cache",  "--cache-secs=5", "--demuxer-readahead-secs=5", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", "--http-header-fields=Referer: https://rapid-cloud.co/", host}

		} else if strings.Contains(host, "tools.fast4speed.rsvp") {
			cmdArgs = []string{"--cache",  "--cache-secs=5", "--demuxer-readahead-secs=5", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", "--http-header-fields=Referer: https://allmanga.to/", host}

		} else if strings.Contains(host, "owocdn.top") || strings.Contains(host, "uwucdn.top"){
			cmdArgs = []string{"--cache",  "--cache-secs=5", "--demuxer-readahead-secs=5", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", "--http-header-fields=Referer: https://kwik.cx/", host}

		} else if strings.Contains(host, "cf-master") || strings.Contains(host, "lethe399key.com") {
			cmdArgs = []string{"--cache",  "--cache-secs=5", "--demuxer-readahead-secs=5", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", "--http-header-fields=Referer: https://vidnest.fun/", host}

		} else if strings.Contains(host, "lightningbolts.ru") {
			cmdArgs = []string{"--cache",  "--cache-secs=5", "--demuxer-readahead-secs=5", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", "--http-header-fields=Referer: https://vidsrc.cc/", host}

		} else if strings.Contains(host, "embedsports.top") || strings.Contains(host, "strmd.top") || strings.Contains(host, "poocloud.in") || strings.Contains(host, "vdcast.live"){
			// Multiple flags needs to prevent the player (mpv) from disconnecting
			cmdArgs = []string{"--retry-open", "5", "--retry-streams", "5", "--stream-segment-attempts", "5", "--stream-segment-timeout", "10", "--player-continuous-http",
		 						"--http-no-ssl-verify", "--http-header", "Referer=https://embedsports.top/", host, "best", "-p", mpv}
			// cmdArgs = []string{url}

		} else if strings.Contains(host, "storm") {
			type Headers struct {
				Referer string `json:"referer"`
				Origin  string `json:"origin"`
			}

			u, err := url.Parse(host)
			if err != nil {
				log.Fatal("unable to parse")
			}

			q := u.Query()

			headersRaw := q.Get("headers")

			var h Headers
			if err := json.Unmarshal([]byte(headersRaw), &h); err != nil {
				log.Fatal(err)
			}

			new_url := strings.Replace(host, "https://storm.vodvidl.site/proxy", q.Get("host"), 1)

			cmdArgs = []string{"--cache",  "--cache-secs=5", "--demuxer-readahead-secs=5", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", "--http-header-fields=Referer: " + h.Referer, new_url}


		} else {
			cmdArgs = []string{"--cache",  "--cache-secs=5", "--demuxer-readahead-secs=5", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", host}
		}

		if len(addSubtitleArgs) > 0 {
			cmdArgs = append(cmdArgs, addSubtitleArgs...)
		}

		var newCmd *exec.Cmd
		if strings.Contains(host, "embedsports.top") || strings.Contains(host, "strmd.top") || strings.Contains(host, "poocloud.in") {
			// newCmd = exec.Command(streamlink, cmdArgs...)
			newCmd = exec.Command(streamlink, cmdArgs...)
			log.Println(newCmd)
		} else {
			newCmd = exec.Command(mpv, cmdArgs...)
			// log.Println(newCmd)
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
		subtitle_url = string([]byte{104, 116, 116, 112, 115, 58, 47, 47, 115, 117, 98, 46, 119, 121, 122, 105, 101, 46, 114, 117, 47, 115, 101, 97, 114, 99, 104, 63, 105, 100, 61}) + strconv.Itoa(tmdbid) + string([]byte{38, 115, 101, 97, 115, 111, 110, 61}) + strconv.Itoa(season) + string([]byte{38, 101, 112, 105, 115, 111, 100, 101, 61}) + strconv.Itoa(episode) + string([]byte{38, 102, 111, 114, 109, 97, 116, 61, 115, 114, 116})
	} else {
		subtitle_url = string([]byte{104, 116, 116, 112, 115, 58, 47, 47, 115, 117, 98, 46, 119, 121, 122, 105, 101, 46, 114, 117, 47, 115, 101, 97, 114, 99, 104, 63, 105, 100, 61}) + strconv.Itoa(tmdbid) + string([]byte{38, 102, 111, 114, 109, 97, 116, 61, 115, 114, 116})
	}

	log.Println(subtitle_url)

	req, _ := http.NewRequest("GET", subtitle_url, nil)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("unable to get the request")
		return []string{}
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var subtitles SubtitlesResponse
	if err := json.Unmarshal(respBody, &subtitles); err != nil {
		log.Println("issue with unmarshalling. Probably not found")
	}

	for _, item := range subtitles {
		if item.Language == "en" {
			if len(subtitleList) <= 2 {
				subtitleList = append(subtitleList, item.Url)
			} else {
				break
			}
		}
	}

	return subtitleList
}
