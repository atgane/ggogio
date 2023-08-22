package main

import "github.com/atgane/ggogio"

type SampleFactory struct {
}

func (s SampleFactory) Create() ggogio.Client {
	return new(SampleClient)
}
