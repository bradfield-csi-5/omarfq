package leveldb

import (
	"testing"
	//"github.com/omarfq/leveldb/iterator"
)

func seedDb() DB {
	return &LevelDb{
		entries: []Entry{
			{
				Key:   []byte("testkey1"),
				Value: []byte("testvalue1"),
			},
		},
	}
}

func TestLevelDb_Get_Ok(t *testing.T) {
	leveldb := seedDb()

	testValue := []byte("testvalue1")

	val, _ := leveldb.Get([]byte("testkey1"))
	if string(val) != string(testValue) {
		t.Errorf("Expected %q, got %q", testValue, val)
	}

}
