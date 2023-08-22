package main

import (
	"log"

	"github.com/atgane/ggogio"
)

func main() {
	addr := ":10000"
	s := ggogio.NewServer(addr, SampleFactory{})
	err := s.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
