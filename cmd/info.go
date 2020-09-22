package main

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/wz2b/dump1090-go-client/pkg/dump1090"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"time"
)

func main() {
	icao := os.Args[1]

	db, err := bolt.Open("aircraft.db", 0644, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("aircraft"))

		start := time.Now()
		data := bucket.Get([]byte(icao))

		a := dump1090.AircraftInfo{}
		err = proto.Unmarshal(data, &a)
		stop := time.Now()
		elapsed := stop.Sub(start)
		fmt.Printf("Lookup time: %d\n", elapsed.Nanoseconds())

		fmt.Printf("transponder: %s\n", *a.Icao24)

		if a.Registration != nil {
			fmt.Printf("Registration: %s\n", *a.Registration)
		}

		if a.Model != nil {
			fmt.Printf("Model: %s\n", *a.Model)
		}

		if a.TypeCode != nil {
			fmt.Printf("Type: %s\n", *a.TypeCode)
		}

		if a.IcaoAircraftType != nil {
			fmt.Printf("ICAO Type: %s\n", *a.IcaoAircraftType)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

}
