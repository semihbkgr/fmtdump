package parse

import "fmt"

func uintBySize(b uint64) (any, error) {
	switch b {
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
			return nil, fmt.Errorf("the len of the referenced field is not supported. size: %d", b)
		}
	}
}

func anyToInt64(i any) (uint64, error) {
	switch v := i.(type) {
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
