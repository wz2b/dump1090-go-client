package main

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/wz2b/dump1090-go-client/pkg/dump1090"
	bolt "go.etcd.io/bbolt"
)

func main() {

	db, err := bolt.Open("aircraft.db", 0644, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("aircraft"))
		bucket.ForEach(func(key []byte, data []byte) error {

			a := dump1090.AircraftInfo{}
			err = proto.Unmarshal(data, &a)
			if a.Registration != nil {
				fmt.Printf("%s %s\n", *a.Icao24, *a.Registration)
			}

			return nil
		})
		return nil
	})

}
