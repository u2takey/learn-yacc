//go:generate goyacc -o parser.go parser.y
package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	for {
		if _, err := os.Stdout.WriteString("> "); err != nil {
			log.Println("WriteString: ", err)
		}
		line, err := in.ReadString('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("ReadString: %s", err)
		}
		sql, err := Parse(line)
		if err != nil {
			log.Fatalf("Parse error: %s", err)
		}
		log.Println("get sql:", sql.String())
		newEnv(sql).Run()
	}
}
