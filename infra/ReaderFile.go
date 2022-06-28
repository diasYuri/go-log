package infra

import (
	"bufio"
	"io"
	"log"
)

type ReaderFile struct {
	reader           *bufio.Reader
	outputLogChannel chan string
}

func NewReaderFile(reader io.Reader) *ReaderFile {
	return &ReaderFile{
		reader: bufio.NewReader(reader),
	}
}

func (r *ReaderFile) ReadLog() (string, error) {
	line, err := r.reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return "", err
	}
	return line, nil
}
