package main

import (
	"log"
	"os"

	"github.com/koxanybak/quit-torrent/torrent"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Not enough arguments. Torrent file path should be the second arg.")
	}
	fpath := args[1]

	process, err := torrent.NewProcess(fpath)
	if err != nil {
		log.Fatal("Error creating a new process: ", err)
	}
	process.Start()
}