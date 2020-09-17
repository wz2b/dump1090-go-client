package store1090

import (
	"github.com/influxdata/influxdb1-client/v2"
	"log1090/pkg/log1090"
	time2 "time"
)

type Store1090 struct {
	client client.Client
}

func Create() (*Store1090, error) {
	client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
		Username: "dump1090",
		Password: "xxxxxxxx",
	})
	if err != nil {
		return nil, err
	} else {
		return &Store1090{
			client: client,
		}, err
	}
}

func (this *Store1090) Write(report log1090.Report) (error) {
	time := time2.Unix(int64(report.Now), 0)

	points, err := client.NewBatchPoints(client.BatchPointsConfig{
		Precision: "s",
		Database:  "dump1090",
	})
	if err != nil {
		return err
	}

	for _, plane := range(report.Aircraft) {
		point, err := makePoint(time, plane)
		if err == nil {
			points.AddPoint(point)
		} else {
			return err
		}
	}


	return this.client.Write(points)
}

func makePoint(t time2.Time, plane log1090.Aircraft) (*client.Point, error) {
	name := "spot"

	/*
	 * Tags
	 */
	tags := map[string]string {
		"hex": plane.Hex,
	}

	if plane.Category != nil {
		tags["category"] = *plane.Category
	}

	/*
	 * Fields
	 */
	fields := map[string]interface{} {}
	if plane.Flight != nil {
		fields["flight"] = *plane.Flight
	}

	if plane.Lat != nil {
		fields["lat"] = *plane.Lat
	}

	if plane.Lon != nil {
		fields["lon"] = *plane.Lon
	}

	if plane.Track != nil {
		track := *plane.Track
		fields["track"] = track

		if track >= 0 && track <= 360 {
			if track < 45 {
				fields["direction"] = "N"
			} else if track < 135 {
				fields["direction"] = "E"
			} else if track < 225 {
				fields["direction"] = "S"
			} else if track < 315 {
				fields["direction"] = "W"
			} else {
				fields["direction"] = "N"
			}
		}

	}

	if plane.GroundSpeed != nil {
		fields["speed"] = *plane.GroundSpeed
	}

	if plane.Rssi != nil {
		fields["rssi"] = *plane.Rssi
	}

	if plane.Squawk != nil {
		fields["squawk"] = *plane.Squawk
	}

	if plane.GeometricAltitude != nil {
		fields["geom"] = *plane.GeometricAltitude
	}
	if plane.Emergency != nil {
		fields["emergency"] = *plane.Emergency
	}


	if plane.BarometerAltitude != nil {

		/*
		 * the alt_baro field might be nil, a float, or a
		 * string "ground"
		 */
		numValue, ok := plane.BarometerAltitude.(float64)
		if ok {
			fields["altitude"] = numValue
		} else {
			stringValue, ok := plane.BarometerAltitude.(string)
			if ok {
				if stringValue == "on_ground" {
					fields["on_ground"] = true
				} else {
					fields["on_ground"] = false
				}
			}
		}
	}

	return client.NewPoint(name, tags, fields, t)
}


func (this *Store1090) Close() (error) {
	return this.client.Close()
}