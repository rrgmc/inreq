package inreq

import (
	"net/http"
	"reflect"
)

// DecodeOperationHeader is a DecodeOperation that gets values from HTTP headers.
type DecodeOperationHeader struct {
}

func (d *DecodeOperationHeader) Decode(ctx DecodeContext, r *http.Request, field reflect.Value,
	typ reflect.Type, tag *Tag) (bool, any, error) {
	values := r.Header.Values(tag.Name)

	if len(values) == 0 {
		return false, nil, nil
	}

	if field.Kind() == reflect.Slice {
		return true, values, nil
	}
	return true, values[0], nil
}
