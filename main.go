package main

import (
	"fmt"
	"log"
	"log1090/pkg/log1090"
	"log1090/pkg/store1090"
	"time"
)

func main() {
	logger, err := store1090.Create()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()



	for {
		fmt.Println("get")
		report, err := log1090.GetReport()

		if err != nil {
			log.Printf("Error retrieving report: %s\n", err)
		} else {
			fmt.Println("store")
			err = logger.Write(report)
			if err != nil {
				log.Printf("Error storing report: %s\n", err)
			}
		}

		time.Sleep(15 * time.Second)

	}



}