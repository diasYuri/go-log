package application

import (
	"lfa.com/logs-master/application/domain"
	"log"
)

type LogBatchForPersistenceApplication struct {
	receiveLogChannel      chan domain.Log
	persistenceLogsChannel chan []interface{}
	batchSize              int
}

func NewLogBatchForPersistenceApplication(receiveLogChan chan domain.Log, persistenceLogsChan chan []interface{}, batchSize int) *LogBatchForPersistenceApplication {
	return &LogBatchForPersistenceApplication{
		receiveLogChannel:      receiveLogChan,
		persistenceLogsChannel: persistenceLogsChan,
		batchSize:              batchSize,
	}
}

func (l *LogBatchForPersistenceApplication) Execute(requestEndChan chan byte, signalEndChan chan byte) {
	indexLogBatch := 0
	batchCount := 0
	channelReceiveCount := 0
	logBatch := make([]interface{}, l.batchSize)

	for {
		select {
		case <-requestEndChan:
			logEndOperation(indexLogBatch, batchCount, channelReceiveCount)
			l.persistenceLogsChannel <- CopyBatch(logBatch[0:indexLogBatch])
			signalEndChan <- 1
			return
		case logReceive := <-l.receiveLogChannel:
			l.PutInBatchLog(&indexLogBatch, logBatch, &logReceive, &batchCount)
		}
	}
}

func (l *LogBatchForPersistenceApplication) PutInBatchLog(indexLogBatch *int, logBatch []interface{}, logReceive *domain.Log, batchCount *int) {
	if l.BatchHasSpace(indexLogBatch) {
		logBatch[*indexLogBatch] = *logReceive
		*indexLogBatch++
	}
	if l.BatchIsFull(indexLogBatch) {
		copyOfLogBatch := CopyBatch(logBatch)
		l.persistenceLogsChannel <- copyOfLogBatch
		*batchCount++
		*indexLogBatch = 0
	}
}

func (l *LogBatchForPersistenceApplication) BatchHasSpace(indexLogBatch *int) bool {
	return *indexLogBatch < l.batchSize
}

func (l *LogBatchForPersistenceApplication) BatchIsFull(indexLogBatch *int) bool {
	return *indexLogBatch == l.batchSize
}

func CopyBatch(logBatch []interface{}) []interface{} {
	copyOfLogBatch := make([]interface{}, len(logBatch))
	copy(copyOfLogBatch, logBatch)
	return copyOfLogBatch
}

func logEndOperation(indexLogBatch int, batchCount int, channelReceiveCount int) {
	log.Printf("Contador de resto na persistencia: %d\n", indexLogBatch)
	log.Printf("Contador de qtd de batch salvos: %d\n", batchCount)
	log.Printf("Contador do channel: %d\n", channelReceiveCount)
}
