package main

import "main/internal/app"

func main() {
	err := app.Run()
	if err != nil {
		panic(err)
	}
}
