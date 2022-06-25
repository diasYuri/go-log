package infra

import (
	"bufio"
	"fmt"
	"io"
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

	for {

		line, err := r.reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return "", err
		}

		return line, nil
	}
}