package parser

import "flag"

type OptionModel struct {
	GraphFilename          string
	WarmupRequestCount     int
	EvaluationRequestCount int
	SpectrumBitSize        int
	UseShortestPath        bool
	OutputFormat           string
	Quiet                  bool
	UseReferenceRanks      bool
	GA                     bool // add
	InsertNewContents      bool // add
}

var GA bool //add
var Options *OptionModel

func PrintDefaults() {
	flag.PrintDefaults()
}

func ParseArgs() *OptionModel {
	Options = new(OptionModel)
	flag.StringVar(&Options.GraphFilename, "G", "", "[required] specify graph configuration file for simulation.")
	flag.IntVar(&Options.WarmupRequestCount, "W", 5000, "set the number of request for warmup.")
	flag.IntVar(&Options.EvaluationRequestCount, "E", 1000, "set the number of request for evaluation.")
	flag.IntVar(&Options.SpectrumBitSize, "B", 4, "set the number of spectrums for iris cache algorithm.")
	flag.BoolVar(&Options.UseShortestPath, "S", false, "force using shortest path routing with iris cache algorithm.")
	flag.StringVar(&Options.OutputFormat, "F", "plain", "set output format {json|plain}.")
	flag.BoolVar(&Options.Quiet, "q", false, "suppress debug output.")
	flag.BoolVar(&Options.UseReferenceRanks, "R", false, "use reference ranks instead of calculating sub-optimal separator ranks.")

	flag.BoolVar(&Options.GA, "GA", false, "Using GA.")                          // add
	flag.BoolVar(&Options.InsertNewContents, "I", false, "Insert new contents.") // add

	flag.Parse()
	return Options
}
