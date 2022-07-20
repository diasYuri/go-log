package application

import (
	"fmt"
	"lfa.com/logs-master/application/domain"
	"regexp"
	"strings"
	"time"
)

type LogTokenizerApplication struct {
	regexCompiled    *regexp.Regexp
	inputLogChannel  chan string
	outputLogChannel chan domain.Log
}

func NewLogTokenizerApplication(logRegexPattern string, inputLogChan chan string, outputLogChan chan domain.Log) *LogTokenizerApplication {
	return &LogTokenizerApplication{
		regexCompiled:    regexp.MustCompile(logRegexPattern),
		inputLogChannel:  inputLogChan,
		outputLogChannel: outputLogChan,
	}
}

func (l *LogTokenizerApplication) Execute(requestEndChan chan byte, signalEndChan chan byte) {
	var count int = 0

	for {
		select {
		case <-requestEndChan:
			signalEndChan <- 1
			fmt.Printf("Contador de tokenizador: %d\n", count)
			return
		case logLine := <-l.inputLogChannel:
			l.Tokenizer(LogMapper, logLine, &count)
		}
	}
}

func (l *LogTokenizerApplication) Tokenizer(mapper func(map[string]string) domain.Log, logLine string, count *int) {
	logTokens := l.regexCompiled.FindStringSubmatch(logLine)

	if len(logTokens) <= 0 {
		return
	}

	mapLog := make(map[string]string)

	for i, nameToken := range l.regexCompiled.SubexpNames() {
		if nameToken != "" {
			mapLog[nameToken] = logTokens[i]
		}
	}

	log := mapper(mapLog)

	*count++
	l.outputLogChannel <- log
}

func LogMapper(mapLog map[string]string) domain.Log {
	time, _ := time.Parse("2006-01-02 15:04:05", string(mapLog["date"]))

	return domain.Log{
		Date:    time,
		Level:   strings.TrimSpace(mapLog["level"]),
		Header:  strings.TrimSpace(mapLog["header"]),
		Message: strings.TrimSpace(mapLog["message"]),
	}
}
