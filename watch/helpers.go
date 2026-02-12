package watch

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//go:embed scripts/setcookies.py
var setcookiesPy string

//go:embed scripts/pahe.py
var pahePy string

//go:embed scripts/stream-impersonate.sh
var streamImpersonateSh string

type Result struct {
	Urls         []string `json:"urls"`
	Subtitles    []string `json:"subtitles"`
	SourceUsed   string   `json:"source_used"`
	TotalSources int      `json:"total_sources"`
}

type PlayResult struct {
	Urls         []string
	Subtitles    []string
	SourceUsed   string
	SourcesTried []string
	TotalSources int
	Success      bool
	Error        error
	UrlsFound    bool
}

func executePythonTask(content string, id int64, season_number int64, episode_number int64, title string, anilist_id int, anime_episode int, sports_url string, anime_title string, skipSources []string) (Result, error) {

	tmpDir, err := os.MkdirTemp("", "gonwatch-scripts-*")
	if err != nil {
		return Result{}, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	pahePath := filepath.Join(tmpDir, "pahe.py")
	if err := os.WriteFile(pahePath, []byte(pahePy), 0644); err != nil {
		return Result{}, fmt.Errorf("failed to write pahe.py: %w", err)
	}

	setcookiesPath := filepath.Join(tmpDir, "setcookies.py")
	if err := os.WriteFile(setcookiesPath, []byte(setcookiesPy), 0644); err != nil {
		return Result{}, fmt.Errorf("failed to write setcookies.py: %w", err)
	}

	skipArg := strings.Join(skipSources, ",")

	cmdArgs := []string{
		setcookiesPath,
		content,
		strconv.FormatInt(id, 10),
		strconv.Itoa(int(season_number)),
		strconv.Itoa(int(episode_number)),
		title,
		strconv.Itoa(anilist_id),
		strconv.Itoa(anime_episode),
		sports_url,
		anime_title,
		"--skip-sources", skipArg,
	}
	log.Println(cmdArgs)

	pythonPath, err := resolvePythonPath()
	if err != nil {
		return Result{}, err
	}

	cmd := exec.Command(pythonPath, cmdArgs...)

	out, err := cmd.Output()
	if err != nil {
		return Result{}, fmt.Errorf("python script failed: %w", err)
	}
	// out, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return Result{}, fmt.Errorf("python script failed: %w\n%s", err, strings.TrimSpace(string(out)))
	// }

	var result Result
	if err := json.Unmarshal(out, &result); err != nil {
		return Result{}, fmt.Errorf("failed to parse python output: %w", err)
	}

	if len(result.Subtitles) == 0 && (content != "anime" && content != "stream") {
		return Result{
			Urls:         result.Urls,
			Subtitles:    GetSubtitles(int(id), content, int(season_number), int(episode_number)),
			SourceUsed:   result.SourceUsed,
			TotalSources: result.TotalSources,
		}, nil
	}

	return result, nil
}

func resolvePythonPath() (string, error) {
	if custom := os.Getenv("GONWATCH_PYTHON"); custom != "" {
		if _, err := os.Stat(custom); err == nil {
			return custom, nil
		}
		return "", fmt.Errorf("GONWATCH_PYTHON not found: %s", custom)
	}

	home, err := os.UserHomeDir()
	if err == nil {
		venvPython := filepath.Join(home, ".local", "share", "gonwatch", "venv", "bin", "python")
		if _, err := os.Stat(venvPython); err == nil {
			return venvPython, nil
		}
	}

	if path, err := exec.LookPath("python3"); err == nil {
		return path, nil
	}
	if path, err := exec.LookPath("python"); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("python not found. Install python3 or set GONWATCH_PYTHON")
}

func openMpv(urls []string, subtitles []string) error {
	if len(urls) == 0 {
		return fmt.Errorf("no URLs to play")
	}

	var addSubtitleArgs []string
	if len(subtitles) > 0 {
		for _, subtitle := range subtitles {
			addSubtitleArgs = append(addSubtitleArgs, "--sub-file="+subtitle)
		}
	}

	var cmdArgs []string
	mpv, err := exec.LookPath("mpv")
	if err != nil {
		return fmt.Errorf("mpv not found: %w", err)
	}

	streamlink, err := exec.LookPath("streamlink")
	if err != nil {
		log.Println("streamlink not found, some streams may not work")
		streamlink = ""
	}

	var lastErr error
	for _, host := range urls {

		if strings.Contains(host, "_v7") {
			cmdArgs = []string{"--cache", "--cache-secs=5", "--demuxer-readahead-secs=5", "--profile=high-quality", "--vo=gpu", "--gpu-api=vulkan", "--scale=ewa_lanczos", "--cscale=ewa_lanczos", "--correct-downscaling=yes", "--dither-depth=auto", "--deband=no", "--hls-bitrate=max", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", "--http-header-fields=Referer: https://rapid-cloud.co/", host}

		} else if strings.Contains(host, "tools.fast4speed.rsvp") {
			cmdArgs = []string{"--cache", "--cache-secs=5", "--demuxer-readahead-secs=5", "--profile=high-quality", "--vo=gpu", "--gpu-api=vulkan", "--scale=ewa_lanczossharp", "--cscale=ewa_lanczossharp", "--correct-downscaling=yes", "--sigmoid-upscaling=yes", "--dither-depth=auto", "--deband=yes", "--deband-iterations=2", "--deband-threshold=35", "--deband-range=16", "--deband-grain=5", "--hls-bitrate=max", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", "--http-header-fields=Referer: https://allmanga.to/", host}

		} else if strings.Contains(host, "one.techparadise") {
			cmdArgs = []string{"--cache", "--cache-secs=5", "--demuxer-readahead-secs=5", "--profile=high-quality", "--vo=gpu", "--gpu-api=vulkan", "--scale=ewa_lanczos", "--cscale=ewa_lanczos", "--correct-downscaling=yes", "--dither-depth=auto", "--deband=no", "--hls-bitrate=max", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", "--http-header-fields=Referer: https://player.videasy.net/", host}

		} else if strings.Contains(host, "owocdn.top") || strings.Contains(host, "uwucdn.top") {
			cmdArgs = []string{"--cache", "--cache-secs=5", "--demuxer-readahead-secs=5", "--profile=high-quality", "--vo=gpu", "--gpu-api=vulkan", "--scale=ewa_lanczossharp", "--cscale=ewa_lanczossharp", "--correct-downscaling=yes", "--sigmoid-upscaling=yes", "--dither-depth=auto", "--deband=yes", "--deband-iterations=2", "--deband-threshold=35", "--deband-range=16", "--deband-grain=5", "--hls-bitrate=max", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", "--http-header-fields=Referer: https://kwik.cx/", host}

		} else if strings.Contains(host, "cf-master") || strings.Contains(host, "lethe399key.com") || strings.Contains(host, "slime403heq.com") {
			cmdArgs = []string{"--cache", "--cache-secs=5", "--demuxer-readahead-secs=5", "--profile=high-quality", "--vo=gpu", "--gpu-api=vulkan", "--scale=ewa_lanczos", "--cscale=ewa_lanczos", "--correct-downscaling=yes", "--dither-depth=auto", "--deband=no", "--hls-bitrate=max", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", "--http-header-fields=Referer: https://vidnest.fun/", host}

		} else if strings.Contains(host, "lightningbolt") {
			cmdArgs = []string{"--cache", "--cache-secs=5", "--demuxer-readahead-secs=5", "--profile=high-quality", "--vo=gpu", "--gpu-api=vulkan", "--scale=ewa_lanczos", "--cscale=ewa_lanczos", "--correct-downscaling=yes", "--dither-depth=auto", "--deband=no", "--hls-bitrate=max", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", "--http-header-fields=Referer: https://vidsrc.cc/", host}

		} else if strings.Contains(host, "strmd.top") {
			err := openMpvWithCurlImpersonate(host)
			if err != nil {
				return fmt.Errorf("unable to play stream: %w", err)
			}
			return nil

		} else if strings.Contains(host, "embedsports.top") || strings.Contains(host, "poocloud.in") || strings.Contains(host, "vdcast.live") {

			cmdArgs = []string{"--retry-open", "5", "--retry-streams", "5", "--stream-segment-attempts", "5", "--stream-segment-timeout", "10", "--player-continuous-http",
				"--http-no-ssl-verify", "--http-header", "Referer=https://embedsports.top/", host, "best", "-p", mpv, "-a", "--network-timeout=60 --stream-lavf-o=reconnect=1,reconnect_streamed=1,reconnect_delay_max=5 --fullscreen"}

		} else if strings.Contains(host, "storm") {
			cmdArgs = []string{"--cache", "--cache-secs=5", "--demuxer-readahead-secs=5", "--profile=high-quality", "--vo=gpu", "--gpu-api=vulkan", "--scale=ewa_lanczos", "--cscale=ewa_lanczos", "--correct-downscaling=yes", "--dither-depth=auto", "--deband=no", "--hls-bitrate=max", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", "--http-header-fields=Referer: https://vidlink.pro/", "--ytdl-raw-options=impersonate=Chrome-131:Android-14,add-header=Referer:https://vidlink.pro/", host}

		} else {
			cmdArgs = []string{"--cache", "--cache-secs=5", "--demuxer-readahead-secs=5", "--profile=high-quality", "--vo=gpu", "--gpu-api=vulkan", "--scale=ewa_lanczos", "--cscale=ewa_lanczos", "--correct-downscaling=yes", "--dither-depth=auto", "--deband=no", "--hls-bitrate=max", "--demuxer-lavf-o=fflags=+genpts", "--no-audio-pitch-correction", "--video-sync=audio", "--stream-lavf-o=reconnect=1", "--stream-lavf-o=reconnect_streamed=1", "--stream-lavf-o=reconnect_delay_max=5", "--stream-lavf-o=reconnect_on_http_error=1", "--stream-lavf-o=reconnect_on_network_error=1", "--fullscreen", "--save-position-on-quit", "--slang=en,eng", host}
		}

		if len(addSubtitleArgs) > 0 {
			cmdArgs = append(cmdArgs, addSubtitleArgs...)
		}

		var newCmd *exec.Cmd
		if strings.Contains(host, "embedsports.top") || strings.Contains(host, "strmd.top") || strings.Contains(host, "poocloud.in") {
			newCmd = exec.Command(streamlink, cmdArgs...)
			// log.Println(newCmd)
		} else {
			newCmd = exec.Command(mpv, cmdArgs...)
			// log.Println(newCmd)
		}

		err = newCmd.Start()
		if err != nil {
			lastErr = fmt.Errorf("failed to start mpv: %w", err)
			continue
		}

		done := make(chan error, 1)
		go func() {
			done <- newCmd.Wait()
		}()

		select {
		case err := <-done:
			if err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					log.Println("Exit code is: ", exitError.ExitCode())
					log.Println("Got exit code, onto the next one")
					lastErr = fmt.Errorf("mpv exited with code %d", exitError.ExitCode())
					continue
				}
				lastErr = fmt.Errorf("mpv failed: %w", err)
				continue
			}
			return nil
		case <-time.After(3 * time.Second):
			return nil
		}
	}

	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("all URLs failed to play")
}

func openMpvWithCurlImpersonate(m3u8URL string) error {
	tmpDir, err := os.MkdirTemp("", "gonwatch-stream-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	scriptPath := filepath.Join(tmpDir, "stream-impersonate.sh")
	if err := os.WriteFile(scriptPath, []byte(streamImpersonateSh), 0755); err != nil {
		return fmt.Errorf("failed to write stream-impersonate.sh: %w", err)
	}

	mpv, err := exec.LookPath("mpv")
	if err != nil {
		return fmt.Errorf("mpv not found: %w", err)
	}

	bashCmd := exec.Command("bash", scriptPath, m3u8URL)
	mpvCmd := exec.Command(mpv,
		"--no-cache",
		"--demuxer-max-bytes=50M",
		"--demuxer-readahead-secs=10",
		"--fullscreen",
		"--force-seekable=no",
		"-",
	)

	// log.Println(m3u8URL)
	// log.Println(bashCmd)

	pipe, err := bashCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create pipe: %w", err)
	}
	mpvCmd.Stdin = pipe

	if err := bashCmd.Start(); err != nil {
		return fmt.Errorf("failed to start stream script: %w", err)
	}
	if err := mpvCmd.Start(); err != nil {
		bashCmd.Process.Kill()
		return fmt.Errorf("failed to start mpv: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- mpvCmd.Wait()
	}()

	select {
	case err := <-done:
		bashCmd.Process.Kill()
		if err != nil {
			return fmt.Errorf("mpv exited with error: %w", err)
		}
		return nil
	case <-time.After(5 * time.Second):
		return nil
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
		log.Println("issue with unmarshalling subtitle. Probably not found")
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
