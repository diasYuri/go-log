package adapters

import (
	"lfa.com/logs-master/application/domain"
)

type PersistenceAdapter interface {
	Save(log domain.Log) error
	SaveBatch(logs []interface{}) error
}
