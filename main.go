package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"

	"github.com/christer79/location-visualizer/config"

	"github.com/dustin/go-heatmap"
	"github.com/dustin/go-heatmap/schemes"
)

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

func filterValues(locations Locations, filter config.Filter) Locations {
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
		os.Exit(1)
		fmt.Println(err)
	}

	dat, err := json.Marshal(locations)
	if err != nil {
		fmt.Println("error:", err)
	}
	out.Write(dat)
}

func main() {

	configPath := flag.String("config", "", "Path to configuration file")

	flag.Parse()
	var filter config.Filter
	var outputformat config.OutputFormats
	var inputs config.Inputs

	if *configPath != "" {
		fmt.Println("config: " + *configPath)
		filter = config.ParseConfigFilter(*configPath).Filter
		fmt.Println(filter)
		outputformat = config.ParseConfigOutputFormat(*configPath)
		fmt.Println(outputformat)
		inputs = config.ParseConfigInput(*configPath)
		fmt.Println(inputs)

	}
	var locations Locations
	for _, input := range inputs.Inputs {

		fmt.Printf("Input \n - type: %s \n - path: %s \n", input.Type, input.Path)

		locations = readData(input.Path)
	}
	getmaxValues(locations)

	filteredLocations := filterValues(locations, filter)

	points := []heatmap.DataPoint{}
	for _, location := range filteredLocations.Locations {
		points = append(points,
			heatmap.P(float64(location.LongitudeE7), float64(location.LatitudeE7)))
	}

	for _, format := range outputformat.Outputs {
		fmt.Printf("Output \n - type: %s \n - path: %s \n", format.Filetype, format.Filename)
		if format.Filetype == "json" {
			writeData(format.Filename, filteredLocations)
		}
		if format.Filetype == "jpeg" {
			// scheme, _ := schemes.FromImage("../schemes/fire.png")
			scheme := schemes.OMG

			img := heatmap.Heatmap(image.Rect(0, 0, format.Width, format.Height),
				points, format.CircleSize, uint8(format.Opacity), scheme)

			fmt.Println("Writing to: " + format.Filename)
			writeImgToFile(img, format.Filename)
		}
	}
}

//png.Encode(os.Stdout, img)
