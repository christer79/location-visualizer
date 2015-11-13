package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

//InputFile holds information about a source for input
type Input struct {
	Path string `yaml:"path"`
	Type string `yaml:"type"`
}

//InputFiles list of files to read input data from
type Inputs struct {
	Inputs []Input `yaml:"inputs"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// IntervalString min adn max values of type string
type IntervalString struct {
	Min string `yaml:"min"`
	Max string `yaml:"max"`
}

//IntervalInt min adn max vaue of type int
type IntervalInt struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

// Filters helper to pars Filter type from yaml configuration
type Filters struct {
	Filter Filter `yaml:"Filter"`
}

//Filter a struct to keep intervals for which to filter location data
type Filter struct {
	Latitude  IntervalInt    `yaml:"Latitude"`
	Longitude IntervalInt    `yaml:"Longitude"`
	Accuracy  IntervalInt    `yaml:"Accuracy"`
	Time      IntervalString `yaml:"Time"`
}

//OutputFormat struct with parameters for how to format the output data
type Format struct {
	Filetype   string `yaml:"type"`
	Filename   string `yaml:"filename"`
	Width      int    `yaml:"width"`
	Height     int    `yaml:"height"`
	CircleSize int    `yaml:"circleSize"`
	Opacity    uint8  `yaml:"opacity"`
}

//OutputFormats helper to help parse OutputFormat from yaml-file
type OutputFormats struct {
	Outputs []Format `ymal:"outputs"`
}

func ParseConfigFilter(filename string) Filters {
	var filter Filters
	dat, frErr := ioutil.ReadFile(filename)
	check(frErr)
	err := yaml.Unmarshal(dat, &filter)
	if err != nil {
		fmt.Println("error:", err)
	}
	return filter
}
func ParseConfigOutputFormat(filename string) OutputFormats {
	var outputformat OutputFormats
	dat, frErr := ioutil.ReadFile(filename)
	check(frErr)
	err := yaml.Unmarshal(dat, &outputformat)
	if err != nil {
		fmt.Println("error:", err)
	}
	return outputformat
}

func ParseConfigInput(filename string) Inputs {
	var inputs Inputs
	dat, frErr := ioutil.ReadFile(filename)
	check(frErr)
	err := yaml.Unmarshal(dat, &inputs)
	if err != nil {
		fmt.Println("error:", err)
	}
	return inputs

}
