package main

import (
	"log"
	"os"
)

func main() {
	// set log to an output file
	f, err := os.OpenFile("redstone.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Redstone log file opened.")

	game := NewGame()
	game.Run()
}
