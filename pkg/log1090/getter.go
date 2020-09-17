package log1090

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func GetReport() (Report, error) {
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

	/*
	 * Clean up the flight ID by removing leading and trailing spaces
	 */
	for i, a := range report.Aircraft {
		if a.Flight != nil {
			trimmed := strings.TrimSpace(*a.Flight)
			report.Aircraft[i].Flight = &trimmed
		}

		//fmt.Println(reflect.TypeOf(a.alt_baro))
	}

	return report, nil
}