package cmd

type parsedOptions struct {
	inputFile      string
	outputFilePath string
	timeSpan       int
	thumbWidth     float32
	sprit          int
}

type commandGenOptions struct {
	inputFilePath  string
	outputFilePath string
	timeSpan       int
	thumbWidth     float32
}
