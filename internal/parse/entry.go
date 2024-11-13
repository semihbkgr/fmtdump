package parse

import "github.com/semihbkgr/fmtdump/internal/format"

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
