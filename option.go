package inreq

import (
	"net/http"

	"github.com/RangelReale/instruct"
)

type Option interface {
	isOption()
}

type DefaultOption interface {
	TypeDefaultOption
	applyDefaultOption(*defaultOptions)
}

type TypeDefaultOption interface {
	Option
	applyTypeDefaultOption(*typeDefaultOptions)
}

type DecodeOption interface {
	Option
	applyDecodeOption(*decodeOptions)
}

type TypeDefaultAndDecodeOption interface {
	TypeDefaultOption
	DecodeOption
}

type FullOption interface {
	DefaultOption
	DecodeOption
}

// PathValue is used by the "path" operation to extract the path from the request. Usually this is stored
// in the context by libraries like "gorilla/mux".
type PathValue interface {
	GetRequestPath(r *http.Request, name string) (found bool, value any, err error)
}

type PathValueFunc func(r *http.Request, name string) (found bool, value any, err error)

func (p PathValueFunc) GetRequestPath(r *http.Request, name string) (found bool, value any, err error) {
	return p(r, name)
}

// FieldNameMapper maps a struct field name to the header/query/form field name.
// The default one uses [strings.ToLower].
type FieldNameMapper = instruct.FieldNameMapper
