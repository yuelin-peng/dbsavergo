package util

import "strconv"

func Interface2Int64(v interface{}) int64 {
	if v == nil {
		return int64(0)
	}
	switch t := v.(type) {
	case int32:
		return int64(t)
	case int64:
		return int64(t)
	case int:
		return int64(t)
	case int8:
		return int64(t)
	case uint8:
		return int64(t)
	case uint:
		return int64(t)
	case uint32:
		return int64(t)
	case uint64:
		return int64(t)
	case string:
		i, _ := strconv.ParseInt(t, 10, 64)
		return i
	default:
		return int64(0)
	}
}
