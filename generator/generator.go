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

func generateFFmpegCommand(inputFilepath string, outputDirPath string, timespan int, width float32, info *VideoInfo) *exec.Cmd {
	commandSlices := []string{
		"-ss",
		fmt.Sprintf("%0.04f", info.Start+0.0001),
		"-i",
		fmt.Sprintf("%s", inputFilepath),
		"-y",
		"-an",
		"-sn",
		"-vsync",
		"0",
		"-q:v",
		"5",
		"-threads",
		"2",
		"-vf",
		fmt.Sprintf("scale=%f:-1,select=not(mod(n\\,%d))", width, info.Tbr*timespan),
		fmt.Sprintf("%s/thumbnails-%%04d.jpg", outputDirPath),
	}
	cmd := exec.Command("ffmpeg", commandSlices...)
	fmt.Println(cmd.String())
	return cmd
}

func splitImages(inputPath string, outputDirPath string, span int, width float32) error {
	fmt.Println("# Start generate thumbnails")
	_, err := os.Stat(outputDirPath)
	if err == nil {
		err = os.RemoveAll(outputDirPath)
		if err != nil {
			return err
		}
	}
	err = os.MkdirAll(outputDirPath+"/thumbnails", 0777)
	if err != nil {
		return err
	}

	videoInfo := GetVideoInfo(inputPath)
	fmt.Printf("\n# Completed get video information. \nduration: %d\nstart: %f\nframe-rate: %d", videoInfo.Seconds, videoInfo.Start, videoInfo.Tbr)

	fmt.Println("\n\n# Start Generate sprit images using ffmpeg")
	command := generateFFmpegCommand(inputPath, outputDirPath+"/thumbnails", span, width, &videoInfo)
	_, err = command.CombinedOutput()

	if err != nil {
		return err
	}
	return nil
}

func GenerateThumbnails(input string, outputDirPath string, span int, width float32, sprit int, outputDir string, debug int) {
	fmt.Println("# Start generate thumbnails")

	_, err := os.Stat(outputDirPath)
	if err == nil && debug == 0 {
		err = os.RemoveAll(outputDirPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	err = os.MkdirAll(outputDirPath+"/thumbnails", 0777)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	videoInfo := GetVideoInfo(input)
	fmt.Printf("\n# Completed get video information. \nduration: %d\nstart: %f\nframe-rate: %d", videoInfo.Seconds, videoInfo.Start, videoInfo.Tbr)

	fmt.Println("\n\n# Start Generate sprit images using ffmpeg")
	if debug == 0 {
		command := generateFFmpegCommand(input, outputDirPath+"/thumbnails", span, width, &videoInfo)

		_, err = command.CombinedOutput()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	files, err := ioutil.ReadDir(outputDirPath + "/thumbnails")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, outputDirPath+"/thumbnails/"+file.Name())
	}
	if debug == 0 {
		// remote all generate sprit images
		defer os.RemoveAll(outputDirPath + "/thumbnails/")
	}

	totalFileCount := len(fileNames)
	if totalFileCount == 0 {
		fmt.Println("ERROR: thumbnail is empty")
		os.Exit(1)
	}
	sort.Strings(fileNames)

	thumbsAcross := math.Min(float64(totalFileCount), float64(sprit))
	rows := math.Ceil(float64(totalFileCount) / float64(thumbsAcross))
	w, h := helper.GetImageDimension(fileNames[0])

	fmt.Println("# Completed Generate sprit images")
	fmt.Printf("total files: %d\nacross: %f\nrows: %f\nwidth: %d\nheight: %d", totalFileCount, thumbsAcross, rows, w, h)
	fmt.Println("\n\n# Starting Create JPEG Image and VTT file")

	tmpTotal, index := 0, 0

	var srcImages [][]image.Image
	var srcImage []image.Image
	srcImages = append(srcImages, srcImage)

	for _, path := range fileNames {
		img, err := helper.GetImage(path)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if (tmpTotal%100) == 0 && tmpTotal != 0 {
			srcImage = nil
			srcImages = append(srcImages, srcImage)
			tmpTotal = 0
			index++
		}

		srcImages[index] = append(srcImages[index], img)
		tmpTotal++
	}

	fmt.Println("Image Splited Length: ", len(srcImages))

	vtt := "WEBVTT\n\n"
	totalSeconds := 0
	for i, src := range srcImages {

		row := math.Ceil(float64(len(src)) / thumbsAcross)

		fmt.Printf("Creating image...\n===size===\nwidth: %d\nheight: %d\n\n", int(float64(w)*thumbsAcross), int(float64(h)*(row)))

		dstImage := image.NewRGBA(
			image.Rect(0, 0, int(float64(w)*thumbsAcross), int(float64(h)*(row))),
		)

		tmpX, tmpY := 0, 0
		for rx, ry, s, f := 0, -1, totalSeconds, 0; f < len(src); f++ {
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
				src[f],
				image.Point{0, 0},
				draw.Src,
			)
			vtt += fmt.Sprintf("%s --> %s\nthumbnails%d.jpg#xywh=%d,%d,%d,%d", t1, t2, i, rx*w, ry*h, w, h)

			rx++
			vtt += "\n\n"
			tmpX = rx * w
			tmpY = ry * h
		}
		if debug != 0 {
			fmt.Printf("generated image last file image position: \nx: %d, y: %d\n\n", tmpX, tmpY)
		}
		totalSeconds += len(src) * span

		dstPath := fmt.Sprintf("%s/thumbnails%d.jpg", outputDirPath, i)
		helper.CreateJPEGImage(dstPath, dstImage, 100)
	}
	helper.WriteString(outputDirPath+"/thumbnails.vtt", vtt)
	fmt.Printf("#Process Completed!!\nOutput: %s", outputDirPath)
	os.Exit(0)
}
