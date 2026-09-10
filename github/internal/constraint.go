package internal

import "strconv"

type NameOrId interface {
	~string | ~int64
}

func TryGetId[T NameOrId](x T) (int64, bool) {
	var value interface{} = x
	if id, ok := value.(int64); ok {
		return id, true
	}

	id, err := strconv.ParseInt(value.(string), 10, 64)
	return id, err == nil
}
