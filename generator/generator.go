package generator

import (
	"fmt"
	"image"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"sort"

	"image/draw"
	_ "image/jpeg"

	"github.com/yochidros/generate-thumbnails/helper"
)

func generateFFmpegCommand(inputFilePah string, outputDirPah string, timespan int, width float32, info *VideoInfo) *exec.Cmd {
	commandSlices := []string{
		"-ss",
		fmt.Sprintf("%0.04f", info.Start+0.0001),
		"-i",
		fmt.Sprintf("%s", inputFilePah),
		"-y",
		"-an",
		"-sn",
		"-vsync",
		"0",
		"-q:v",
		"5",
		"-threads",
		"1",
		"-vf",
		fmt.Sprintf("scale=%f:-1,select=not(mod(n\\,%d))", width, info.Tbr*timespan),
		fmt.Sprintf("%s/tthumbnails-%%04d.jpg", outputDirPah),
	}
	cmd := exec.Command("ffmpeg", commandSlices...)
	fmt.Println(cmd.String())
	return cmd
}

func GenerateThumbnails(input string, outputDirPah string, span int, width float32, sprit int) {
	fmt.Println("# Start generate thumbnails")
	_, err := os.Stat(outputDirPah)
	if err == nil {
		err = os.RemoveAll(outputDirPah)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	err = os.MkdirAll(outputDirPah+"/thumbnails", 0777)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	videoInfo := GetVideoInfo(input)
	fmt.Printf("\n# Completed get video information. \nduration: %d\nstart: %f\nframe-rate: %d", videoInfo.Seconds, videoInfo.Start, videoInfo.Tbr)

	fmt.Println("\n\n# Start Generate sprit images using ffmpeg")
	command := generateFFmpegCommand(input, outputDirPah+"/thumbnails", span, width, &videoInfo)
	_, err = command.CombinedOutput()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(outputDirPah + "/thumbnails")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, outputDirPah+"/thumbnails/"+file.Name())
	}
	totalFileCount := len(fileNames)
	if totalFileCount == 0 {
		fmt.Println("ERROR: thumbnail is empty")
		os.Exit(1)
	}
	sort.Strings(fileNames)

	thumbsAcross := math.Min(float64(totalFileCount), float64(sprit))
	rows := math.Ceil(float64(totalFileCount) / float64(thumbsAcross))

	fmt.Println("# Completed Generate sprit images")
	fmt.Printf("total files: %d\nacross: %f\nrows: %f\n", totalFileCount, thumbsAcross, rows)

	fmt.Println("\n\n# Starting Create JPEG Image and VTT file")
	w, h := helper.GetImageDimension(fileNames[0])

	var srcImages []image.Image
	for _, pah := range fileNames {
		img, err := helper.GetImage(pah)
		if err != nil {
			fmt.Println(err)
			continue
		}
		srcImages = append(srcImages, img)
	}

	dstImage := image.NewRGBA(
		image.Rect(0, 0, int(float64(w)*thumbsAcross), int(float64(h)*rows)),
	)

	vtt := "WEBVTT\n\n"
	for rx, ry, s, f := 0, -1, 0, 0; f < totalFileCount; f++ {
		t1 := fmt.Sprintf("%02d:%02d:%02d.000", s/3600, (s / 60 % 60), s%60)
		s += span
		t2 := fmt.Sprintf("%02d:%02d:%02d.000", s/3600, (s / 60 % 60), s%60)

		if f%int(thumbsAcross) == 0 {
			rx = 0
			ry++
		}

		draw.Draw(
			dstImage,
			image.Rect(rx*w, ry*h, (rx*w)+w, (ry+1)*h),
			srcImages[f],
			image.Point{0, 0},
			draw.Src,
		)
		rx++

		vtt += fmt.Sprintf("%s --> %s\nthumbnails.jpg#xywh=%d,%d,%d,%d", t1, t2, rx*w, ry*h, w, h)
		vtt += "\n\n"
	}

	helper.CreateJPEGImage(outputDirPah+"/thumbnails.jpg", dstImage, 100)
	helper.WriteString(outputDirPah+"/thumbnails.vtt", vtt)
	fmt.Printf("#Process Completed!!\nOutput: %s", outputDirPah)
	os.Exit(0)
}
