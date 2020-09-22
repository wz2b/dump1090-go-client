package main

import (
	"bufio"
	csv "encoding/csv"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/wz2b/dump1090-go-client/pkg/dump1090"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"reflect"
)

func main() {
	db, err := bolt.Open("aircraft.db", 0644, nil)
	if err != nil {
		panic(err)
	}

	if err != mergeSkynet(db) {
		panic(err)
	}
	//
	//if err !=	mergeMicro(db) {
	//	log.Fatal(err)
	//}

	db.Close()
}

/*
 * Merge in data from opensky-network.org/datasets/
 */
func mergeSkynet(db *bolt.DB) error {
	file, err := os.Open("../aircraftDatabase.csv")

	if err != nil {
		return err
	}
	defer file.Close()
	scanner := csv.NewReader(bufio.NewReader(file))
	header, err := scanner.Read()
	if err != nil {
		panic(err)
	}
	i := 0

	tx, err := db.Begin(true)
	if err != nil {
		log.Fatal(err)
	}

	record, err := scanner.Read()
	if err != nil {
		panic(err)
	}
	for ; err == nil; record, err = scanner.Read() {

		rowMap := map[string]string{}

		for i, name := range header {
			rowMap[name] = record[i]
		}

		if rowMap["icao24"] == "" {
			continue
		}

		a := new(dump1090.AircraftInfo)
		err := populate(a, rowMap)
		if err != nil {
			panic(err)
		}

		//if err := scanner.Populate(&a); err != nil {
		//	return err
		//}

		/*
		 * Commit transactions in blocks, for performance reasons.  If it's time to
		 * commit, generate a new block.
		 */
		i++
		if i%10000 == 0 {
			fmt.Println(i)
			if err := tx.Commit(); err != nil {
				return err
			}
			tx, err = db.Begin(true)
			if err != nil {
				return err
			}
		}

		bkt, err := tx.CreateBucketIfNotExists([]byte("aircraft"))
		if err != nil {
			return err
		}

		key := []byte(*a.Icao24)
		if key == nil || len(key) < 1 {
			panic("Key is invalid")
		}

		msg, err := proto.Marshal(a)
		bkt.Put(key, msg)

	}

	// Commit the last transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

/*
 * Merge data from micrtronics.de/aircraft-dtabase/export.php
 */
func mergeMicro(b *bolt.DB) error {
	file, err := os.Open("../icao24plus.txt")
	if err != nil {
		return err
	}

	bufferedReader := bufio.NewReader(file)
	reader := csv.NewReader(bufferedReader)
	reader.Comma = '\t'

	// Skip the first row - it's a date and time
	record, err := reader.Read()
	if err != nil {
		return err
	}

	for record, err = reader.Read(); record != nil; record, err = reader.Read() {
		if err != nil {
			return err
		}
	}

	return nil
}

func populate(obj interface{}, values map[string]string) (err error) {
	v := reflect.ValueOf(obj)
	t := v.Elem().Type()

	for x := 0; x < t.NumField(); x++ {
		decoration := t.Field(x).Tag.Get("csv")

		if decoration != "" {
			value := values[decoration]
			field := v.Elem().Field(x)
			fieldKind := field.Kind()

			if fieldKind == reflect.String {
				field.SetString(value)
			} else if fieldKind == reflect.Ptr {
				/*
				 * The incoming column is a string, but it is being assigned
				 * to a pointer variable.  For this purpose, "" being the empty
				 * string value should be interpreted as assigning the field to null.
				 */
				newValue := value
				if newValue != "" {
					field.Set(reflect.ValueOf(&newValue))
				}

			} else {
				return errors.New("Cannot convert this type")
			}
		}
	}
	return nil
}
