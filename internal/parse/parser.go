package parse

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/semihbkgr/fmtdump/internal/format"
)

type Parser struct {
	reader io.Reader
	format *format.Format
}

func NewParser(r io.Reader, f *format.Format) *Parser {
	return &Parser{
		reader: r,
		format: f,
	}
}

func (p *Parser) Next() ([]*Data, error) {
	data := make([]*Data, 0, len(p.format.Blocks))
	for _, block := range p.format.Blocks {
		size, err := getSize(block, data)
		if err != nil {
			return nil, err
		}
		buff := make([]byte, size)
		_, err = p.reader.Read(buff)
		if err != nil {
			return nil, err
		}
		data = append(data, &Data{
			Block: block,
			Value: buff,
		},
		)
	}
	return data, nil
}

func getSize(b *format.Block, data []*Data) (uint, error) {
	if b.Size.IsLiteral() {
		return b.Size.Size, nil
	}
	for _, d := range data {
		if d.Block.Name == b.Size.Reference {
			var size int16
			_, err := binary.Decode(d.Value, binary.LittleEndian, &size)
			if err != nil {
				return 0, err
			}
			return uint(size), nil
		}
	}
	return 0, errors.New("size reference cannot be found")
}
