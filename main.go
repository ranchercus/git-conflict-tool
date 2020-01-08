package main

import (
	"io/ioutil"
	"log"
)

const fileName  = "changelist"

func main() {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

}