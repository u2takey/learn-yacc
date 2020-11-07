//go:generate goyacc -o parser.go  parser.y
package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("please input codes")
		return
	}

	codes, err := Parse([]byte(os.Args[1]))
	if err != nil {
		log.Println("parse error", err)
		return
	}
	//log.Println(codes)

	newEnv(codes).Run()

}
