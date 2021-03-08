package generator

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/yochidros/generate-thumbnails/helper"
)

func GetVideoInfo(filename string) VideoInfo {
	commandSlices := []string{
		"-hide_banner",
		"-i",
		fmt.Sprintf("%s", filename),
	}
	if !helper.CommandExist("ffprobe") {
		os.Exit(1)
	}

	// get video infomation duration and tbr.
	cmd := exec.Command("ffprobe", commandSlices...)
	fmt.Println(cmd.String())

	out, _ := cmd.CombinedOutput()

	info := VideoInfo{}

	regex, _ := regexp.Compile("Duration: ((\\d+):(\\d+):(\\d+))\\.\\d+, start: ([^,]*)")
	result := regex.FindSubmatch(out)
	if len(result) >= 4 {
		hour, _ := strconv.Atoi(string(result[2]))
		minute, _ := strconv.Atoi(string(result[3]))
		seconds, _ := strconv.Atoi(string(result[4]))
		start, _ := strconv.ParseFloat(string(result[5]), 64)
		s := hour*3600 + minute*60 + seconds
		info.Start = start
		info.Seconds = s
	}

	// get video framerate
	regex, _ = regexp.Compile("(\\d+(?:\\.\\d+)?) tbr")
	result = regex.FindSubmatch(out)

	if len(result) >= 1 {
		tbr, _ := strconv.ParseFloat(string(result[1]), 64)
		info.Tbr = int(math.Ceil(tbr))
	}

	return info
}
