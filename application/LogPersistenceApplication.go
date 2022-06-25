package application

import (
	"fmt"
	"lfa.com/logs-master/application/adapters"
	"lfa.com/logs-master/application/domain"
)

type LogPersistenceApplication struct {
	inputLogChannel chan domain.Log
	logPersistence  adapters.PersistenceAdapter
}

func NewLogPersistenceApplication(persistence adapters.PersistenceAdapter, inputLogChan chan domain.Log) *LogPersistenceApplication {
	return &LogPersistenceApplication{
		inputLogChannel: inputLogChan,
		logPersistence:  persistence,
	}
}

func (l *LogPersistenceApplication) Execute(requestEndChan chan byte, signalEndChan chan byte) {
	lenBatch := 1000
	indexLogBatch := 0
	batchCount := 0
	channelReceiveCount := 0
	var logBatch = make([]interface{}, lenBatch)

	for {
		select {
		case <-requestEndChan:
			l.logEnd(indexLogBatch, batchCount, channelReceiveCount)
			l.logPersistence.SaveBatch(logBatch[0:indexLogBatch])
			signalEndChan <- 1
			return
		case log := <-l.inputLogChannel:
			if indexLogBatch < lenBatch {
				logBatch[indexLogBatch] = log
				indexLogBatch++
			}

			if indexLogBatch == lenBatch {
				l.logPersistence.SaveBatch(logBatch)
				batchCount++
				indexLogBatch = 0
			}
		}
	}
}

func (l *LogPersistenceApplication) logEnd(indexLogBatch int, batchCount int, channelReceiveCount int) {
	fmt.Printf("Contador de resto na persistencia: %d\n", indexLogBatch)
	fmt.Printf("Contador de qtd de batch salvos: %d\n", batchCount)
	fmt.Printf("Contador do channel: %d\n", channelReceiveCount)
}
