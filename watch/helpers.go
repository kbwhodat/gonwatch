package watch

import (
	// "bytes"
	"strings"
	"encoding/json"
	// "fmt"
	"log"
	"os/exec"
	"strconv"
	_ "embed"
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
func executePythonTask(content string, id int64, season_number int64, episode_number int64) ([]string, []string) {

	log.Println("going to run python script...")
	cmdArgs := []string{}

	if content == "tv" {
		cmdArgs = []string{"scripts/setcookies.py", content, strconv.FormatInt(id, 10), strconv.Itoa(int(season_number)), strconv.Itoa(int(episode_number))}
	} else {
		cmdArgs = []string{"scripts/setcookies.py", content, strconv.FormatInt(id, 10)}
	}
	cmd := exec.Command("python", cmdArgs...)

	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(string(out))

	var result Result
	if err := json.Unmarshal(out, &result); err != nil {
		log.Fatal(err)
	}

	return result.Urls, result.Subtitles
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
	checkForErrors(err)

	for _, url := range urls {
		if strings.Contains(url, "shadowlandschronicles.com") {
			cmdArgs = []string{"--save-position-on-quit", "slang=en", url}
		} else {
			cmdArgs = []string{"--save-position-on-quit", "slang=en", "--http-header-fields=Referer: https://vidsrc.cc/", url}
		}

		if len(addSubtitleArgs) > 0 {
			cmdArgs = append(cmdArgs, addSubtitleArgs...)
		}

		newCmd := exec.Command(dir, cmdArgs...)
		err := newCmd.Run()

		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				if exitError.ExitCode() != 0 {
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
