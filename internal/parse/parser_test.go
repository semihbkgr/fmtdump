package parse

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/semihbkgr/fmtdump/internal/format"
)

func TestParser(t *testing.T) {
	f, err := format.ParseFormatFile("./testdata/format.json")
	if err != nil {
		t.Error(err)
	}
	if err := f.Validate(); err != nil {
		t.Error(err)
	}
	file, err := os.Open("./testdata/data.bin")
	if err != nil {
		t.Error(err)
	}
	p := NewParser(file, f)

	expectedEntries := []map[string][]byte{
		{
			"crc":     {0x75, 0x59, 0x96, 0x0d},
			"len":     {0x07, 0x00},
			"type":    {0x00},
			"payload": {0x65, 0x6e, 0x74, 0x72, 0x79, 0x2d, 0x31},
		},
		{
			"crc":     {0xcf, 0x08, 0x9f, 0x94},
			"len":     {0x07, 0x00},
			"type":    {0x00},
			"payload": {0x65, 0x6e, 0x74, 0x72, 0x79, 0x2d, 0x32},
		},
		{
			"crc":     {0x59, 0x38, 0x98, 0xe3},
			"len":     {0x07, 0x00},
			"type":    {0x00},
			"payload": {0x65, 0x6e, 0x74, 0x72, 0x79, 0x2d, 0x33},
		},
	}

	i := 0
	for entry, err := p.Next(); err != io.EOF; entry, err = p.Next() {
		if err != nil {
			t.Error(err)
		}
		expectedEntry := expectedEntries[i]
		for field, value := range expectedEntry {
			data := entry.data(field)
			if data == nil {
				t.Fatalf("expected field %s, got nil", field)
			}
			if !bytes.Equal(data.Value, value) {
				t.Errorf("expected field %s to be %x, got %x", field, value, data)
			}
		}
		i++
	}

}
