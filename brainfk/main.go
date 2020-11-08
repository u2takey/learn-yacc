//go:generate goyacc -o parser.go  parser.y
package main

import (
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("please input codes")
		return
	}

	var input []byte
	var err error
	if os.Args[1] == "-f" {
		log.Println("read code from file")
		input, err = ioutil.ReadFile(os.Args[2])
		if err != nil {
			log.Panic(err)
		}
	} else {
		input = []byte(os.Args[1])
	}

	codes, err := Parse(input)
	if err != nil {
		log.Println("parse error", err)
		return
	}

	newEnv(codes).Run()

}
