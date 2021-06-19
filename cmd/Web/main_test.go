package main

import (
	"log"
	"testing"
)

func TestRun(t *testing.T) {
	err := run()
	if err != nil {
		log.Println(err)
		t.Error("Run() function failed !")
	}
}
