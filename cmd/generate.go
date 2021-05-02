package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yochidros/generate-thumbnails/generator"
	"github.com/yochidros/generate-thumbnails/helper"
)

type createOptions struct {
	input  string
	output string
	span   int
	width  float32
	sprit  int
	debug  int
}

var ffmpegCommand = "ffmpeg -ss %0.04f -i %s -y -an -sn -vsync 0 -q:v 5 -threads 1 -vf scale=%d:-1,select=%d %s/%s/%s-%%04d.jpg"

func NewCmdGenerate() *cobra.Command {
	var (
		o = &createOptions{}
	)
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate thumbnails",
		Run: func(cmd *cobra.Command, args []string) {
			runGenerateThumbnails(o)
		},
	}
	cmd.Flags().StringVarP(&o.input, "input", "i", "", "input file path")
	cmd.Flags().StringVarP(&o.output, "output", "o", "", "output file path")
	cmd.Flags().IntVarP(&o.span, "time-span", "t", 1, "time span")
	cmd.Flags().Float32VarP(&o.width, "width", "w", 120.0, "thumbnails width")
	cmd.Flags().IntVarP(&o.sprit, "sprit", "s", 10, "thumbnails sprit col length")
	cmd.Flags().IntVarP(&o.debug, "debug", "d", 0, "debug mode")
	return cmd
}

func init() {
}

func runGenerateThumbnails(option *createOptions) {
	if !helper.RequireDepends() {
		fmt.Fprintf(os.Stderr, "Error: Dependencies not installed, please install ffmpeg, ffprobe")
		os.Exit(1)
	}

	po := parsedOptions{}

	if len(option.input) == 0 {
		fmt.Println("Error: Input file is Empty Be Required")
		os.Exit(1)
	}

	cwd, _ := os.Getwd()

	inputFilePath := cwd + "/" + option.input

	po.inputFile = inputFilePath

	if len(option.output) == 0 {
		fmt.Println("Error: Output dirctory is Empty Be Required")
		os.Exit(1)
	}

	outputFileDirPath := cwd + "/out"
	_, err := os.Stat(outputFileDirPath)

	// create output dir if not exist.
	if err != nil {
		err = os.Mkdir(outputFileDirPath, 0777)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err)
		}
		fmt.Println("Create output Directory: ", outputFileDirPath)
	}

	po.outputFilePath = outputFileDirPath + "/" + option.output
	po.timeSpan = option.span
	po.thumbWidth = option.width
	po.sprit = option.sprit

	generator.GenerateThumbnails(po.inputFile, po.outputFilePath, po.timeSpan, po.thumbWidth, po.sprit, "out/"+option.output, option.debug)
}
