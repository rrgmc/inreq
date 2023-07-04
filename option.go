package inreq

import (
	"net/http"

	"github.com/RangelReale/instruct"
	"github.com/RangelReale/instruct/options"
)

type (
	Option        = options.Option
	AnyOption     = options.AnyOption[*http.Request, DecodeContext]
	DefaultOption = options.DefaultOption[*http.Request, DecodeContext, defaultOptions]
	DecodeOption  = options.DecodeOption[*http.Request, DecodeContext, decodeOptions]

	AnyTypeOption     = options.AnyTypeOption[*http.Request, DecodeContext]
	TypeDefaultOption = options.TypeDefaultOption[*http.Request, DecodeContext, typeDefaultOptions]
	TypeDecodeOption  = options.TypeDecodeOption[*http.Request, DecodeContext, decodeOptions]

	DefaultAndTypeDefaultOption = options.DefaultAndTypeDefaultOption[*http.Request, DecodeContext, defaultOptions, typeDefaultOptions]
	DefaultAndDecodeOption      = options.DefaultAndDecodeOption[*http.Request, DecodeContext, defaultOptions, decodeOptions]
	TypeDefaultAndDecodeOption  = options.TypeDefaultAndDecodeOption[*http.Request, DecodeContext, typeDefaultOptions, decodeOptions]
	FullOption                  = options.FullOption[*http.Request, DecodeContext, defaultOptions, typeDefaultOptions, decodeOptions, decodeOptions]
)

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
