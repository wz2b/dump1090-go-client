package log1090

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Report struct {
	Now      float64    `json:"now"`
	Messages uint64     `json:"messages"`
	Aircraft []Aircraft `json:"aircraft"`
}
type Aircraft struct {
	Hex               string  `json:"hex"`
	Squawk            string  `json:"squawk"`
	Lat               float64 `json:"lat"`
	Lon               float64 `json:"lon"`
	Flight            string  `json:"flight"`
	GroundSpeed       float64 `json:"gs"`
	Track             float64 `json:"track"`
	Emergency         string  `json:"emergency"`
	Category          string  `json:"category"`
	Rssi              float32 `json:"rssi"`
	GeometricAltitude *int64  `json:"alt_geom"`
	BarometerAltitude *int64  `json:"alt_baro,omitempty"`
}

func main() {
	resp, err := http.Get("http://adsb.gis.rit.edu/dump1090-fa/data/aircraft.json")
	if err != nil {
		log.Fatal("Unable to get data ", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Unable to read body ", err)
	}

	var report Report
	err = json.Unmarshal(body, &report)
	if err != nil {
		log.Fatal("Unable to parse body ", err)
	}

	for i, a := range report.Aircraft {
		report.Aircraft[i].Flight = strings.TrimSpace(a.Flight)
	}

	fmt.Printf("%vv\n", report)
}
