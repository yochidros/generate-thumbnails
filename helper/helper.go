package helper

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
)

func RequireDepends() bool {
	return CommandExist("ffmpeg") && CommandExist("ffprobe")
}

func CommandExist(name string) bool {
	cmd := exec.Command(name, "-h")
	if err := cmd.Run(); err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func GetImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open File ERROR: %v\n", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func GetImageDimension(path string) (int, int) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open File ERROR: %v\n", err)
	}
	defer file.Close()

	conf, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File decode Config: %s: %v\n", path, err)
		return 0, 0
	}
	return conf.Width, conf.Height
}

func IntMax(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func CreateJPEGImage(path string, img image.Image, quality int) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	opts := &jpeg.Options{Quality: quality}
	err = jpeg.Encode(f, img, opts)
	return err
}

func WriteString(path string, content string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}
