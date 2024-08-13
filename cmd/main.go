package main

import (
	"flag"
	"main/internal/app"
)

var configPath = flag.String("config", "", "a path to config")

func main() {
	flag.Parse()
	err := app.Run(*configPath)
	if err != nil {
		panic(err)
	}
}
