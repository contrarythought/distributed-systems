package server

import (
	"fmt"
	"sync"
)

type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

type Log struct {
	records []Record
	mu      sync.Mutex
}

func NewLog() *Log {
	return &Log{}
}

// append a record to the log
func (l *Log) Append(record Record) uint64 {
	l.mu.Lock()
	defer l.mu.Unlock()

	record.Offset = uint64(len(l.records))
	l.records = append(l.records, record)
	return record.Offset
}

var ErrOutOfBounds = fmt.Errorf("offset out of bounds")

// retrieve a record from the log
func (l *Log) Read(offset uint64) (Record, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if offset >= uint64(len(l.records)) {
		return Record{}, ErrOutOfBounds
	}

	return l.records[offset], nil
}
