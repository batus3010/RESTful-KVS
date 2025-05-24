package kvs

import (
	"encoding/json"
	"fmt"
	"io"
)

type Table []KVPair

func NewTable(rdr io.Reader) ([]KVPair, error) {
	var table Table
	err := json.NewDecoder(rdr).Decode(&table)
	if err != nil {
		err = fmt.Errorf("problem parsing table, %v", err)
	}
	return table, err
}

func (table Table) Find(key string) *KVPair {
	for i, k := range table {
		if k.Key == key {
			return &table[i]
		}
	}
	return nil
}
