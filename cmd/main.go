package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"lfa.com/logs-master/application"
	"lfa.com/logs-master/application/Job"
	"lfa.com/logs-master/application/domain"
	"lfa.com/logs-master/infra"
	"log"
	"os"
	"time"
)

const pattern string = `(?P<date>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2},\d{0,3})[ ]{1,2}-[ ]{1,2}(?P<level>[A-Z]* )[ ]{1,2}(?P<header>\[.*?] )[ ]{0,1}-[ ]{0,2}(?P<message>.*)`

func main() {
	ctx := context.Background()
	defer infra.TimeTrack(time.Now(), "Application")

	LoadEnv()

	stringLogChan := make(chan string)
	logChan := make(chan domain.Log)
	logBatchChan := make(chan []interface{})
	endSignalReader := make(chan byte)
	endSignalTokenizer := make(chan byte)
	endSignalPersistence := make(chan byte)
	endSignalPersistenceJob := make(chan byte)

	defer close(stringLogChan)
	defer close(logChan)
	defer close(logBatchChan)
	defer close(endSignalReader)
	defer close(endSignalTokenizer)
	defer close(endSignalPersistence)
	defer close(endSignalPersistenceJob)

	file := OpenFile(os.Getenv("FILE_LOG_PATH"))
	defer file.Close()

	client := FactoryMongoClient(os.Getenv("MONGO_CONNECTION_STRING"))
	defer client.Disconnect(ctx)

	logReaderApplication := FactoryLogReaderApplication(file, stringLogChan)
	logTokenizerApplication := FactoryLogTokenizerApplication(stringLogChan, logChan)
	logPersistenceApplication := FactoryLogPersistenceApplication(logChan, logBatchChan)
	logPersistenceJob := FactoryLogPersistenceJob(logBatchChan, client, os.Getenv("MONGO_DATABASE"), os.Getenv("MONGO_COLLECTION"))

	go logReaderApplication.Execute(endSignalReader)
	go logTokenizerApplication.Execute(endSignalReader, endSignalTokenizer)
	go logPersistenceApplication.Execute(endSignalTokenizer, endSignalPersistence)
	go StartJobPersistence(logPersistenceJob, endSignalPersistence, endSignalPersistenceJob)

	<-endSignalPersistenceJob
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func OpenFile(logPath string) *os.File {
	file, err := os.Open(logPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return file
}

func StartJobPersistence(job *Job.LogPersistenceJob, requestEndChan chan byte, signalEndChan chan byte) {
	var listRequestEndSignalPersistenceJob []chan byte
	var listEndSignalPersistenceJob []chan byte

	concurrencyPersistenceJob := 10
	for i := 0; i < concurrencyPersistenceJob; i++ {
		requestEndSignalPersistenceJob := make(chan byte)
		endSignalPersistenceJob := make(chan byte)

		listRequestEndSignalPersistenceJob = append(listRequestEndSignalPersistenceJob, requestEndSignalPersistenceJob)
		listEndSignalPersistenceJob = append(listEndSignalPersistenceJob, endSignalPersistenceJob)

		go job.Execute(requestEndSignalPersistenceJob, endSignalPersistenceJob)
	}

	<-requestEndChan

	for _, requestEnd := range listRequestEndSignalPersistenceJob {
		requestEnd <- 1
	}

	for _, endSignal := range listEndSignalPersistenceJob {
		<-endSignal
	}

	for _, requestEnd := range listRequestEndSignalPersistenceJob {
		close(requestEnd)
	}

	for _, endSignal := range listEndSignalPersistenceJob {
		close(endSignal)
	}

	signalEndChan <- 1
}

func FactoryLogReaderApplication(file *os.File, stringLogChan chan string) *application.LogReaderApplication {
	readerAdapter := infra.NewReaderFile(file)

	return application.NewLogReaderApplication(readerAdapter, stringLogChan)
}
func FactoryLogTokenizerApplication(stringLogChan chan string, logChan chan domain.Log) *application.LogTokenizerApplication {
	return application.NewLogTokenizerApplication(domain.GetPattern(os.Getenv("REGEX_PATTERN")), stringLogChan, logChan)
}

func FactoryLogPersistenceApplication(logChan chan domain.Log, logBatchChan chan []interface{}) *application.LogBatchForPersistenceApplication {
	return application.NewLogBatchForPersistenceApplication(logChan, logBatchChan, 10000)
}
func FactoryLogPersistenceJob(logBatchChan chan []interface{}, client *mongo.Client, database string, collection string) *Job.LogPersistenceJob {
	persistence := infra.NewPersistenceMongo(client, database, collection)

	return Job.NewLogPersistenceJob(persistence, logBatchChan)
}

func FactoryMongoClient(uri string) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client
}
