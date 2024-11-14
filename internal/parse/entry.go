package parse

import (
	"encoding/binary"
	"fmt"
	"slices"
	"strings"

	"github.com/semihbkgr/fmtdump/internal/format"
)

type Data struct {
	Field format.Field
	Value []byte
}

type Entry []Data

func (e Entry) data(n string) *Data {
	for _, d := range e {
		if d.Field.Name == n {
			return &d
		}
	}
	return nil
}

func (e Entry) String() (string, error) {
	b := strings.Builder{}
	p := len(slices.MaxFunc(e, func(a, b Data) int {
		return len(a.Field.Name) - len(b.Field.Name)
	}).Field.Name)
	for _, d := range e {
		s, err := d.String(p)
		if err != nil {
			return "", err
		}
		b.WriteString(s + "\n")
	}
	return b.String(), nil
}

func (d Data) String(p int) (string, error) {
	s, err := d.ValueString()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s| %*s: %s", hexRawBytes(d.Value), p, d.Field.Name, s), nil
}

func (d Data) ValueString() (string, error) {
	switch d.Field.Type {
	case format.IntType:
		i, err := uintBySize(uint64(len(d.Value)))
		if err != nil {
			return "", err
		}
		_, err = binary.Decode(d.Value, binaryEndian(d.Field.Encoding), i)
		if err != nil {
			return "", err
		}
		value, err := anyToInt64(i)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", value), nil

	case format.StringType:
		b := make([]byte, len(d.Value))
		_, err := binary.Decode(d.Value, binaryEndian(d.Field.Encoding), b)
		if err != nil {
			return "", err
		}
		return string(b), nil

	case format.BytesType:
		b := make([]byte, len(d.Value))
		_, err := binary.Decode(d.Value, binaryEndian(d.Field.Encoding), b)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%v", b), nil
	}
	return "", fmt.Errorf("unsupported")
}

func binaryEndian(e format.Encoding) binary.ByteOrder {
	if e == format.LittleEndianEncoding {
		return binary.LittleEndian
	}
	return binary.BigEndian
}

func hexRawBytes(bytes []byte) string {
	s := strings.Builder{}
	for i, b := range bytes {
		s.WriteString(fmt.Sprintf("%02x ", b))
		if i+1%8 == 0 && i != len(bytes)-1 {
			s.WriteString("\n")
		}
	}
	return fmt.Sprintf("%-24s", s.String())
}
