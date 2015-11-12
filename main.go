package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"math"
	"os"

	"github.com/dustin/go-heatmap"
	"github.com/dustin/go-heatmap/schemes"
	"gopkg.in/yaml.v2"
)

//InputFile holds information about a source for input
type InputFile struct {
	Path     string `yaml:"path"`
	FileType string `yaml:"type"`
}

//InputFiles list of files to read input data from
type InputFiles struct {
	InputFiles InputFile `yaml:"InputFile"`
}

//InputFilesHelper helper to read list of files to read input data from
type InputFilesHelper struct {
	InputFiles InputFiles `yaml:"InputFiles"`
}

//Activity keep an activity type and yhe confidence of that activity
type Activity struct {
	Type       string `json:"type"`
	Confidence int    `json:"confidence"`
}

//Activitys holds a set of different possible Activity as well as a Timestamp
type Activitys struct {
	TimestampMs string     `json:"timestampMs"`
	Activities  []Activity `json:"activities"`
}

//Location holds longitude/latitude/timestamp as well as accuracy and a list of possible activities
type Location struct {
	TimestampMs string      `json:"timestampMs"`
	LatitudeE7  int         `json:"latitudeE7"`
	LongitudeE7 int         `json:"longitudeE7"`
	Accuracy    int         `json:"accuracy"`
	Activitys   []Activitys `json:"activitys"`
}

//Locations is a list of Location
type Locations struct {
	Locations []Location `json:"locations"`
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
type OutputFormat struct {
	Filetype   string `yaml:"type"`
	Filename   string `yaml:"filename"`
	Width      int    `yaml:"width"`
	Height     int    `yaml:"height"`
	CircleSize int    `yaml:"circleSize"`
	Opacity    uint8  `yaml:"opacity"`
}

//OutputFormats helper to help parse OutputFormat from yaml-file
type OutputFormats struct {
	OutputFormat OutputFormat `yaml:"OutputFormat"`
}

func getmaxValues(locations Locations) {
	var maxLong = 0
	var maxLat = 0
	var minLong = 99999999999999
	var minLat = 999999999999999
	var maxAccuracy = 0
	var minAccuracy = 99999999999
	for _, location := range locations.Locations {
		if location.LatitudeE7 < minLat {
			minLat = location.LatitudeE7
		}
		if location.LongitudeE7 < minLong {
			minLong = location.LongitudeE7
		}
		if location.LatitudeE7 > maxLat {
			maxLat = location.LatitudeE7
		}
		if location.LongitudeE7 > maxLong {
			maxLong = location.LongitudeE7
		}
		if location.Accuracy > maxAccuracy {
			maxAccuracy = location.Accuracy
		}
		if location.Accuracy < minAccuracy {
			minAccuracy = location.Accuracy
		}
	}
	fmt.Printf("Latitude (%v, %v) - diff: %v\n", minLat, maxLat, maxLong-minLong)
	fmt.Printf("Longitude (%v, %v) - diff: %v \n", minLong, maxLong, maxLat-minLat)
	fmt.Printf("Accuracy (%v, %v)\n", minAccuracy, maxAccuracy)

}

func filterValues(locations Locations, filter Filter) Locations {
	var filtered Locations

	for _, location := range locations.Locations {
		if location.LatitudeE7 < filter.Latitude.Min {
			continue
		}
		if location.LatitudeE7 > filter.Latitude.Max {
			continue
		}
		if location.LongitudeE7 < filter.Longitude.Min {
			continue
		}
		if location.LongitudeE7 > filter.Longitude.Max {
			continue
		}

		filtered.Locations = append(filtered.Locations, location)
	}
	return filtered
}

func writeImgToFile(img image.Image, filename string) {
	out, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var opt jpeg.Options

	opt.Quality = 100
	// ok, write out the data into the new JPEG file

	err = jpeg.Encode(out, img, &opt) // put quality to 80%
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Generated image to " + filename + " \n")
}

func readData(filename string) Locations {

	dat, frErr := ioutil.ReadFile(filename)
	check(frErr)
	//fmt.Print(string(dat))

	var locations Locations
	err := json.Unmarshal(dat, &locations)
	if err != nil {
		fmt.Println("error:", err)
	}
	return locations
}

func writeData(filename string, locations Locations) {
	out, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dat, err := json.Marshal(locations)
	if err != nil {
		fmt.Println("error:", err)
	}
	out.Write(dat)
}

func parseConfigFilter(filename string) Filters {
	var filter Filters
	dat, frErr := ioutil.ReadFile(filename)
	check(frErr)
	err := yaml.Unmarshal(dat, &filter)
	if err != nil {
		fmt.Println("error:", err)
	}
	return filter
}
func parseConfigOutputFormat(filename string) OutputFormat {
	var outputformat OutputFormats
	dat, frErr := ioutil.ReadFile(filename)
	check(frErr)
	err := yaml.Unmarshal(dat, &outputformat)
	if err != nil {
		fmt.Println("error:", err)
	}
	return outputformat.OutputFormat
}

func parseConfigInput(filename string) InputFiles {
	var inputfileshelper InputFilesHelper
	dat, frErr := ioutil.ReadFile(filename)
	check(frErr)
	err := yaml.Unmarshal(dat, &inputfileshelper)
	if err != nil {
		fmt.Println("error:", err)
	}
	return inputfileshelper.InputFiles
}

func main() {

	configPath := flag.String("config", "", "Path to configuration file")

	//inputFile := flag.String("input", "Platshistorik.json", "Input file")
	outputFile := flag.String("output", "output.jpg", "Output file")

	minLong := flag.Int("minlong", math.MinInt64, "Minimum value of longitude")
	maxLong := flag.Int("maxlong", math.MaxInt64, "Maximum value of longitude")
	minLat := flag.Int("minlat", math.MinInt64, "Minimum value of latitude")
	maxLat := flag.Int("maxlat", math.MaxInt64, "Maximum value of latitude")

	jpegWidth := flag.Int("width", 1024, "Width of output file")
	jpegHeight := flag.Int("height", 1024, "Height of output file.")
	jpegOpacity := flag.Uint("opacity", 128, "Opacity of heat circles")
	jpegCircleSize := flag.Int("size", 2, "Imact circle size of data point")

	flag.Parse()
	var filter Filter
	var outputformat OutputFormat
	var inputfiles InputFiles

	if *configPath != "" {
		fmt.Println("config: " + *configPath)
		filter = parseConfigFilter(*configPath).Filter
		fmt.Println(filter)
		outputformat = parseConfigOutputFormat(*configPath)
		fmt.Println(outputformat)
		inputfiles = parseConfigInput(*configPath)
		fmt.Println(inputfiles)

	} else {
		filter = Filter{Latitude: IntervalInt{Min: *minLat, Max: *maxLat}, Longitude: IntervalInt{Min: *minLong, Max: *maxLong}}
		outputformat = OutputFormat{Filename: *outputFile, Width: *jpegWidth, Height: *jpegHeight, CircleSize: *jpegCircleSize, Opacity: uint8(*jpegOpacity)}

	}

	fmt.Println("Input:" + inputfiles.InputFiles.Path)

	locations := readData(inputfiles.InputFiles.Path)
	getmaxValues(locations)

	filteredLocations := filterValues(locations, filter)

	writeData("filtered_data.json", filteredLocations)

	points := []heatmap.DataPoint{}
	for _, location := range filteredLocations.Locations {
		points = append(points,
			heatmap.P(float64(location.LongitudeE7), float64(location.LatitudeE7)))
	}

	// scheme, _ := schemes.FromImage("../schemes/fire.png")
	scheme := schemes.OMG

	img := heatmap.Heatmap(image.Rect(0, 0, outputformat.Width, outputformat.Height),
		points, outputformat.CircleSize, uint8(outputformat.Opacity), scheme)

	fmt.Println("Writing to: " + outputformat.Filename)
	writeImgToFile(img, outputformat.Filename)
}

//png.Encode(os.Stdout, img)
