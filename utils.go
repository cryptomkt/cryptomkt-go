package client

import (
	"errors"
	"fmt"
)

func Mapss(mapsi map[string]interface{}) (map[string]string, error) {
	mapss := make(map[string]string)
	for k, v := range mapsi {
		switch vv := v.(type) {
		case string:
			mapss[k] = vv
		case int:
			mapss[k] = string(vv)
		case float64:
			mapss[k] = fmt.Sprintf("%f", vv)
		default:
			return nil, errors.New("unexpected format, argument map only accepts strings, ints and floats")
		}
	}
	return mapss, nil
}
