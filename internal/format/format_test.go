package format

import (
	"encoding/json"
	"os"
	"testing"
)

func TestParseFormat(t *testing.T) {
	file, err := os.ReadFile("./testdata/format.json")
	if err != nil {
		t.Error(err)
	}
	var format Format
	err = json.Unmarshal(file, &format)
	if err != nil {
		t.Error(err)
	}
	if err := format.Validate(); err != nil {
		t.Error(err)
	}
	if len(format) != 3 {
		t.Error("expected 3 blocks, got", len(format))
	}
}

func TestParseFormatFile(t *testing.T) {
	f, err := ParseFormatFile("./testdata/format.json")
	if err != nil {
		t.Error(err)
	}
	if err := f.Validate(); err != nil {
		t.Error(err)
	}
	if len(f) != 3 {
		t.Error("expected 3 blocks, got", len(f))
	}
}
