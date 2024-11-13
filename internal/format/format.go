package format

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Format struct {
	Blocks []*Block
}

func ParseFormat(blocks []string) (*Format, error) {
	format := Format{
		Blocks: make([]*Block, 0, len(blocks)),
	}
	for _, block := range blocks {
		b, err := ParseBlock(block)
		if err != nil {
			return nil, err
		}
		if !b.Size.IsLiteral() && !format.HasBlock(b.Size.Reference) {
			return nil, fmt.Errorf("referenced block not exists, %s", b.Size.Reference)
		}
		format.Blocks = append(format.Blocks, b)
	}
	return &format, nil
}

func (f *Format) HasBlock(n string) bool {
	for _, b := range f.Blocks {
		if b.Name == n {
			return true
		}
	}
	return false
}

func (f *Format) String() string {
	var s strings.Builder
	for _, b := range f.Blocks {
		if b.Size.IsLiteral() {
			s.WriteString(fmt.Sprintf("%s=%d\n", b.Name, b.Size.Size))
		} else {
			s.WriteString(fmt.Sprintf("%s=%s\n", b.Name, b.Size.Reference))
		}
	}
	return s.String()
}

type Block struct {
	Name string
	Size *BlockSize
}

func ParseBlock(block string) (*Block, error) {
	name, sizeStr, found := strings.Cut(block, "=")
	if !found {
		return nil, errors.New("error on parsing the format")
	}
	size, err := ParseBlockSize(sizeStr)
	if err != nil {
		return nil, err
	}
	return &Block{
		Name: name,
		Size: size,
	}, nil
}

type BlockSize struct {
	Size      uint
	Reference string
}

func ParseBlockSize(s string) (*BlockSize, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return &BlockSize{
			Reference: s,
		}, nil
	}
	if i < 1 {
		return nil, errors.New("block size must be a positive value")
	}
	return &BlockSize{
		Size: uint(i),
	}, nil
}

func (s *BlockSize) IsLiteral() bool {
	return s.Size > 0
}
