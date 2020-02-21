package main

import (
	"log"
	"testing"
)

func init() {
	initializePhyKeys()
}

func TestFingerCheck(t *testing.T) {

	str := "asdf"

	v := fingerCheck(str)

	log.Printf("Value of '%s' is %v", str, v)

}

func TestHandCheck(t *testing.T) {

	str := "asdf"

	v := handCheck(str)

	log.Printf("Value of '%s' is %v", str, v)

}

func TestThing(t *testing.T) {

	words = []string{"orlando"}

	for _, word := range words {
		log.Print(dothething(word))
	}

}
