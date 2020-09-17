# Go Parser for dump1090

This project isn't so much a library as it is a demonstration of how to poll
[PiAware Dump1090](https://flightaware.com/adsb/piaware/install)
 parse the data.  PiAware is a FlightAware client program that
 runs on a Raspberry Pi to securely transmit dump1090 ADS-B
 and Mode S data to FlightAware.  In simple terms: it receives
 position, altitude, and speed from aircraft.
 
 The main place where this project might save people time is
 the log1090 object model.  It is a go struct with json
 decorations for simple parsing.
 
 To use this, see main.go.  Briefly:
 
 ```go
	resp, err := http.Get("http://some_host/dump1090-fa/data/aircraft.json")
	if err != nil {
		log.Fatal("Unable to get data ", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Unable to read body ", err)
	}

	var report Report
	err = json.Unmarshal(body, &report)
```
 
Report is an object that gives you an array of Aircraft objects
with all the extracted info.

I'd like to reformulate into a useful library that can then
be used by other things including
[telegraf](https://www.influxdata.com/time-series-platform/telegraf/),
time permitting.  Contributions (in the form of pull requests) are welcome.




