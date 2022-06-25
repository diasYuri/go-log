package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"lfa.com/logs-master/application"
	"lfa.com/logs-master/application/domain"
	"lfa.com/logs-master/infra"
	"log"
	"os"
	"time"
)

const pattern string = `(?P<date>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2},\d{0,3})[ ]{1,2}-[ ]{1,2}(?P<level>[A-Z]* )[ ]{1,2}(?P<header>\[.*?] )[ ]{0,1}-[ ]{0,2}(?P<message>.*)`

func main() {
	//ctx := context.Background()
	defer timeTrack(time.Now(), "Logs")

	stringLogChan := make(chan string)
	logChan := make(chan domain.Log)
	endSignalReader := make(chan byte)
	endSignalTokenizer := make(chan byte)
	endSignalPersistence := make(chan byte)

	file, err := os.Open("./logs.log")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	logPersistenceApplication, client := FactoryLogPersistenceApplication(logChan)

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	logTokenizerApplication := FactoryLogTokenizerApplication(stringLogChan, logChan)
	logReaderApplication := FactoryLogReaderApplication(file, stringLogChan)

	go logReaderApplication.Execute(endSignalReader)
	go logTokenizerApplication.Execute(endSignalReader, endSignalTokenizer)
	go logPersistenceApplication.Execute(endSignalTokenizer, endSignalPersistence)

	<-endSignalPersistence
}

func FactoryLogReaderApplication(file *os.File, stringLogChan chan string) *application.LogReaderApplication {
	readerAdapter := infra.NewReaderFile(file)

	return application.NewLogReaderApplication(readerAdapter, stringLogChan)
}
func FactoryLogTokenizerApplication(stringLogChan chan string, logChan chan domain.Log) *application.LogTokenizerApplication {
	return application.NewLogTokenizerApplication(pattern, stringLogChan, logChan)
}
func FactoryLogPersistenceApplication(logChan chan domain.Log) (*application.LogPersistenceApplication, *mongo.Client) {
	const uri = "mongodb+srv://app-user:Nbwz4CuC01grU56s@cluster0.fdeyt.mongodb.net/?retryWrites=true&w=majority"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	persistence := infra.NewPersistenceMongo(client, "logs_master", "logs_06")
	return application.NewLogPersistenceApplication(persistence, logChan), client
}
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
