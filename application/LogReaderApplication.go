package application

import (
	"fmt"
	"lfa.com/logs-master/application/adapters"
)

type LogReaderApplication struct {
	reader           adapters.ReaderAdapter
	outputLogChannel chan string
	signalEndChan    chan byte
}

func NewLogReaderApplication(reader adapters.ReaderAdapter, outputLogChan chan string) *LogReaderApplication {
	return &LogReaderApplication{
		reader:           reader,
		outputLogChannel: outputLogChan,
	}
}

func (l *LogReaderApplication) Execute(signalEndChan chan byte) {
	var count int = 0

	for {
		log, err := l.reader.ReadLog()
		if err != nil {
			signalEndChan <- 1
			break
		}

		count++
		l.outputLogChannel <- log
	}

	fmt.Printf("Contador de leitura: %d\n", count)
}
