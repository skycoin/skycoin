package main

import (
	"log"
	"time"
)

func parseTimeValue(strTime string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, strTime)
	if err != nil {
		log.Printf("Time conversion error. Format=%s Value=\"%s\" Error: %s", time.RFC3339, strTime, err)
	}
	return t, err
}
