package test

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strings"
	"testing"
)

func TestTrue(t *testing.T) {
	assert.True(t, true)
}

func TestFileAPI(t *testing.T) {
	os.Create("testfile")
	defer os.Remove("testfile")

	file, err := os.Open("testfile") // need to create first?
	if err != nil {
		log.Printf("file failed to open")
		t.Fail()
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Fail()
	}

}

type Consumer interface {
	Consume()
}

type TestConsumer struct {
	Consume func()
}

func TestStringSplitting(t *testing.T) {
	fields := strings.Split("blargh::eggs::bacon", "::")
	assert.True(t, fields[0] == "blargh")
}

func TestPipe(t *testing.T) {
	//reader, writer, _ := os.Pipe()
	reader := os.Stdin
	go func() {
		writer.Write([]byte{0, 1, 2, 3})
		err := writer.Close()
		if err != nil {
			log.Fatal("couldnt close writer")
		}
	}()
	var contents []byte
	go func() {
		reader.Read(contents)
		err := reader.Close()
		if err != nil {
			log.Fatal("couldnt close reader")
		}
	}()
	assert.True(t, contents[0] == 0)

}
