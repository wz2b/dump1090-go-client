package main

import (
	"fmt"
	"github.com/wz2b/dump1090-go-client/pkg/dump1090"
	"log"
	"time"
)

func main() {

	for {
		fmt.Println("get")
		report, err := dump1090.GetReport()

		if err != nil {
			log.Printf("Error retrieving report: %s\n", err)
		} else {
			fmt.Printf("%+v\n", report)
		}

		time.Sleep(15 * time.Second)

	}

}
