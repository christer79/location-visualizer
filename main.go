package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/christer79/location-visualizer/config"
	"github.com/christer79/location-visualizer/googlelocationdata"

	"github.com/gorilla/mux"
)

var locations googlelocationdata.Locations

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type MapData struct {
	Title string
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	p := MapData{Title: "Test"}
	t, err := template.ParseFiles("html/view.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, p)
}
func rangesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	datebegin := vars["datebegin"]
	dateend := vars["dateend"]
	filter := config.Filter{Longitude: config.IntervalInt{Min: -10000000000, Max: 1000000000}, Latitude: config.IntervalInt{Min: -1000000000, Max: 1000000000}, Accuracy: config.IntervalInt{Min: 0, Max: 9000}, Time: config.IntervalString{Min: datebegin, Max: dateend}}
	filtered := googlelocationdata.FilterValues(locations, filter)
	retval := googlelocationdata.GetmaxValues(filtered)
	json.NewEncoder(w).Encode(retval)

}

func filterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("filterHandler: CALLED")
	vars := mux.Vars(r)
	longitude, _ := strconv.Atoi(vars["longitude"])
	latitude, _ := strconv.Atoi(vars["latitude"])
	width, _ := strconv.Atoi(vars["width"])
	height, _ := strconv.Atoi(vars["height"])
	datebegin := vars["datebegin"]
	dateend := vars["dateend"]

	filter := config.Filter{Longitude: config.IntervalInt{Min: longitude, Max: longitude + width}, Latitude: config.IntervalInt{Min: latitude, Max: latitude + height}, Accuracy: config.IntervalInt{Min: 0, Max: 9000}, Time: config.IntervalString{Min: datebegin, Max: dateend}}
	filtered := googlelocationdata.FilterValues(locations, filter)
	json.NewEncoder(w).Encode(filtered)
}

func timeregionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Time in region")
	vars := mux.Vars(r)
	longitude, _ := strconv.Atoi(vars["longitude"])
	latitude, _ := strconv.Atoi(vars["latitude"])
	width, _ := strconv.Atoi(vars["width"])
	height, _ := strconv.Atoi(vars["height"])
	datebegin := vars["datebegin"]
	dateend := vars["dateend"]

	filter := config.Filter{Longitude: config.IntervalInt{Min: longitude, Max: longitude + width}, Latitude: config.IntervalInt{Min: latitude, Max: latitude + height}, Accuracy: config.IntervalInt{Min: 0, Max: 9000}, Time: config.IntervalString{Min: datebegin, Max: dateend}}
	//filtered := googlelocationdata.FilterValues(locations, filter)
	timeinregion := googlelocationdata.TimeInRegion(locations, filter)
	json.NewEncoder(w).Encode(timeinregion)
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

	for _, input := range inputs.Inputs {

		fmt.Printf("Input \n - type: %s \n - path: %s \n", input.Type, input.Path)

		locations = googlelocationdata.ReadData(input.Path)
	}

	fmt.Println("Max values of all data")
	googlelocationdata.GetmaxValues(locations)

	/*
		filteredLocations := googlelocationdata.FilterValues(locations, filter)

		fmt.Println("Max values of filtered data")
		googlelocationdata.GetmaxValues(filteredLocations)

		points := []heatmap.DataPoint{}
		for _, location := range filteredLocations.Locations {
			points = append(points,
				heatmap.P(float64(location.LongitudeE7), float64(location.LatitudeE7)))
		}

		for _, format := range outputformat.Outputs {
			fmt.Printf("Output \n - type: %s \n - path: %s \n", format.Filetype, format.Filename)

			if format.Filetype == "json" {
				googlelocationdata.WriteData(format.Filename, filteredLocations, format.Filetype)
			}
			if format.Filetype == "csv" {
				googlelocationdata.WriteData(format.Filename, filteredLocations, format.Filetype)
			}

			if format.Filetype == "jpeg" {
				// scheme, _ := schemes.FromImage("../schemes/fire.png")
				scheme := schemes.OMG

				img := heatmap.Heatmap(image.Rect(0, 0, format.Width, format.Height),
					points, format.CircleSize, uint8(format.Opacity), scheme)

				fmt.Println("Writing to: " + format.Filename)
				writeImgToFile(img, format.Filename)
			}

			if format.Filetype == "timeinregion" {
				out, err := os.Create(format.Filename)
				if err != nil {
					os.Exit(1)
					fmt.Println(err)
				}

				defer out.Close()

				inregion := false
				var count int
				for _, location := range locations.Locations {
					if googlelocationdata.InRegion(location, format.Filter) {
						if inregion == false {
							line := fmt.Sprintf("Inside: %v - ", comparedates.ParseTimeNs(location.TimestampMs))
							if _, err = out.WriteString(line); err != nil {
								panic(err)
							}
							//		fmt.Printf(line)
						}
						count++
						inregion = true
					} else {
						if inregion == true {
							line := fmt.Sprintf(" %v   (%d)\n", comparedates.ParseTimeNs(location.TimestampMs), count)
							if _, err = out.WriteString(line); err != nil {
								panic(err)
							}
							//	fmt.Printf(line)
							count = 0
						}
						inregion = false
					}
				}
			}
		}
	*/

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/view/", viewHandler)
	router.HandleFunc("/filter/{longitude}/{latitude}/{width}/{height}/{datebegin}/{dateend}/", filterHandler)
	router.HandleFunc("/ranges/{datebegin}/{dateend}/", rangesHandler)
	router.HandleFunc("/timeregion/{longitude}/{latitude}/{width}/{height}/{datebegin}/{dateend}/", timeregionHandler)
	log.Fatal(http.ListenAndServe(":8080", router))

}
