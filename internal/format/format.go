package format

import (
	"fmt"
)

type Format []*Block

func (f *Format) Validate() error {
	names := make(map[string]struct{})
	for i, b := range *f {
		if b.Name == "" {
			return fmt.Errorf("block name cannot be empty. index: %d", i)
		}
		if _, found := names[b.Name]; found {
			return fmt.Errorf("duplicated block name. index: %d", i)
		}
		if b.Size == nil && b.SizeRef == nil {
			return fmt.Errorf("either size or sizeRef must be provided. index: %d", i)
		}
		if b.Size != nil && b.SizeRef != nil {
			return fmt.Errorf("both size and sizeRef cannot be provided. index: %d", i)
		}
		if b.SizeRef != nil {
			if _, found := names[*b.SizeRef]; !found {
				return fmt.Errorf("sizeRef '%s' does not exist. index: %d", *b.SizeRef, i)
			}
		}
		if b.Size != nil && *b.Size < 1 {
			return fmt.Errorf("size must be greater than 0. index: %d", i)
		}
		if !b.Type.isValid() {
			return fmt.Errorf("invalid type '%s'. index: %d", b.Type, i)
		}
		if b.Encoding.isValid() {
			return fmt.Errorf("invalid encoding '%s'. index: %d", b.Encoding, i)
		}
		names[b.Name] = struct{}{}
	}
	return nil
}

type Block struct {
	Name     string   `json:"name"`
	Size     *int     `json:"size"`
	SizeRef  *string  `json:"sizeRef"`
	Encoding Encoding `json:"encoding"`
	Type     Type     `json:"type"`
}

func (b *Block) IsVarSized() bool {
	return b.SizeRef != nil
}

type Type string

const (
	IntType    Type = "int"
	StringType Type = "string"
	BytesType  Type = "bytes"
)

func (t Type) isValid() bool {
	return t == IntType || t == StringType || t == BytesType
}

type Encoding string

const (
	LittleEndianEncoding Encoding = "LittleEndian"
	BigEndianEncoding    Encoding = "BigEndian"
)

func (e Encoding) isValid() bool {
	return e == LittleEndianEncoding || e == BigEndianEncoding
}
