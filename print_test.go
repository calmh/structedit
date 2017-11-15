package structedit

import (
	"testing"
	"time"
)

type TestStruct struct {
	Hello   bool `name:"hello"`
	Goodbye bool
	Level   int
	Options []string `name:"options"`
	Color   string
	KVs     []KV
	MainKV  KV
	Map     map[string]string `name:"rymes"`
	Created time.Time
}

type KV struct {
	Key   string `structedit:"key"`
	Value int
}

func TestPrintStruct(t *testing.T) {
	Print(TestStruct{
		Hello:   true,
		Options: []string{"bgp", "ospf", "odd banana"},
		Color:   "green",
		KVs: []KV{
			{"foo", 42},
			{"bar", 43},
		},
		MainKV: KV{"main", -1},
		Map: map[string]string{
			"see you later": "alligator",
			"in":            "a while crocodile",
		},
		Created: time.Now(),
	})
}
