package parse

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/semihbkgr/fmtdump/internal/format"
)

type Parser struct {
	reader io.Reader
	format format.Format
}

func NewParser(r io.Reader, f format.Format) *Parser {
	return &Parser{
		reader: r,
		format: f,
	}
}

func (p *Parser) Next() (Entry, error) {
	entry := make(Entry, 0, len(p.format))
	for _, field := range p.format {
		size, err := size(field, entry)
		if err != nil {
			return nil, err
		}
		buff := make([]byte, size)
		_, err = p.reader.Read(buff)
		if err != nil {
			return nil, err
		}
		entry = append(entry, Data{
			Field: field,
			Value: buff,
		},
		)
	}
	return entry, nil
}

func size(b format.Field, entry Entry) (uint64, error) {
	if !b.IsVarSized() {
		return *b.Size, nil
	}
	for _, d := range entry {
		if d.Field.Name == *b.SizeRef {
			size, err := sizeInt(d.Field, entry)
			if err != nil {
				return 0, err
			}
			_, err = binary.Decode(d.Value, binaryEndian(d.Field.Encoding), size)
			if err != nil {
				return 0, err
			}
			return anyToInt64(size)
		}
	}
	return 0, errors.New("size reference cannot be found")
}

func sizeInt(b format.Field, entry Entry) (any, error) {
	if b.IsVarSized() {
		panic("todo")
	}
	refData := entry.data(b.Name)
	return uintBySize(*refData.Field.Size)
}
