package domain

import "time"

type Log struct {
	Date    time.Time `bson:"data,omitempty"`
	Level   string    `bson:"level,omitempty"`
	Header  string    `bson:"header,omitempty"`
	Message string    `bson:"message,omitempty"`
}

func NewLog(date time.Time, level string, header string, message string) *Log {
	return &Log{}
}
