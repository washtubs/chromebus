package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	stdinReader := bufio.NewReader(os.Stdin)
	readBuf, _, err := stdinReader.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(string(readBuf))

}
