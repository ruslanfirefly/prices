package main

import "log"

func error_log(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
