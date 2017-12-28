package main

import (
	"data"
	"rest"
)

func main() {
	// init mongo connection
	err := data.InitMongo("localhost")
	if err != nil {
		panic(err)
	}
	// init http server
	router := rest.InitHackAPI()
	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
}
