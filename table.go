package kvs

import (
	"encoding/json"
	"fmt"
	"io"
)

func NewTable(rdr io.Reader) ([]KVPair, error) {
	var KVPairs []KVPair
	err := json.NewDecoder(rdr).Decode(&KVPairs)
	if err != nil {
		err = fmt.Errorf("problem parsing table, %v", err)
	}
	return KVPairs, err
}
