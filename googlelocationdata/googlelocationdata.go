package googlelocationdata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/christer79/location-visualizer/comparedates"
	"github.com/christer79/location-visualizer/config"
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

func GetmaxValues(locations Locations) {
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
func InRegion(location Location, filter config.Filter) bool {
	var beginTime = comparedates.ParseTimeStr(filter.Time.Min)
	var endTime = comparedates.ParseTimeStr(filter.Time.Max)
	var compareTime = false
	if filter.Time.Min != "" && filter.Time.Max != "" {
		compareTime = true
	}
	var timeStamp time.Time

	if location.LatitudeE7 < filter.Latitude.Min {
		return false
	}
	if location.LatitudeE7 > filter.Latitude.Max {
		return false
	}
	if location.LongitudeE7 < filter.Longitude.Min {
		return false
	}
	if location.LongitudeE7 > filter.Longitude.Max {
		return false
	}
	if compareTime {
		timeStamp = comparedates.ParseTimeNs(location.TimestampMs)
		if !comparedates.InTimespan(beginTime, endTime, timeStamp) {
			return false
		}
	}
	return true

}
func FilterValues(locations Locations, filter config.Filter) Locations {
	var filtered Locations

	for _, location := range locations.Locations {
		if InRegion(location, filter) {
			filtered.Locations = append(filtered.Locations, location)
		}
	}
	return filtered
}

func ReadData(filename string) Locations {

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

func WriteData(filename string, locations Locations) {
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
