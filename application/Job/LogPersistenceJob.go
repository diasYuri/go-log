package Job

import "lfa.com/logs-master/application/adapters"

type LogPersistenceJob struct {
	persistenceLog       adapters.PersistenceAdapter
	inputLogBatchChannel chan []interface{}
}

func NewLogPersistenceJob(persistenceLog adapters.PersistenceAdapter, inputLogBatchChannel chan []interface{}) *LogPersistenceJob {
	return &LogPersistenceJob{
		persistenceLog:       persistenceLog,
		inputLogBatchChannel: inputLogBatchChannel,
	}
}

func (l *LogPersistenceJob) Execute(requestEndChan chan byte, signalEndChan chan byte) {

	for {
		select {
		case <-requestEndChan:
			signalEndChan <- 1
			return

		case logBatch := <-l.inputLogBatchChannel:
			l.persistenceLog.SaveBatch(logBatch)
		}
	}

}
