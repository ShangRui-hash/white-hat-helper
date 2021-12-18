package main

import (
	"log"
	"white-hat-helper/pkg/hackflow"
)

func main() {
	if err := hackflow.GetPython().Run("main.py"); err != nil {
		log.Fatal(err)
	}
}
