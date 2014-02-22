package main

import (
	"log"
	"strings"
)

func error_log(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func prepareStrings(str string) string{
	str = strings.Trim(str, " ")
	return str
}
