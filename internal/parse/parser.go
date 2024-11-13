package parse

import (
	"encoding/binary"
	"errors"
	"fmt"
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
			_, err = binary.Decode(d.Value, binary.LittleEndian, size)
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
	switch *refData.Field.Size {
	case 1:
		{
			var i uint8
			return &i, nil
		}
	case 2:
		{
			var i uint16
			return &i, nil
		}
	case 4:
		{
			var i uint32
			return &i, nil
		}
	case 8:
		{
			var i uint64
			return &i, nil
		}
	default:
		{
			//todo
			return nil, fmt.Errorf("the len of the referenced field is not supported. len: %d", *refData.Field.Size)
		}
	}
}

func anyToInt64(a any) (uint64, error) {
	switch v := a.(type) {
	case *uint8:
		return uint64(*v), nil
	case *uint16:
		return uint64(*v), nil
	case *uint32:
		return uint64(*v), nil
	case *uint64:
		return *v, nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}
}
