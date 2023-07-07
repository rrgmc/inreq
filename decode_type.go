package inreq

import (
	"net/http"

	"github.com/RangelReale/instruct"
	inoptions "github.com/RangelReale/instruct/options"
)

// TypeDecoder decodes http requests to structs.
type TypeDecoder[T any] struct {
	dec            *instruct.TypeDecoder[*http.Request, DecodeContext, T]
	defaultOptions typeDefaultOptions
}

// NewTypeDecoder creates a Decoder instance with the default decode operations (query, path, header, form, body).
func NewTypeDecoder[T any](options ...TypeDefaultOption) *TypeDecoder[T] {
	return NewCustomTypeDecoder[T](inoptions.ConcatOptionsBefore[TypeDefaultOption](options, WithDefaultDecodeOperations())...)
}

// NewCustomTypeDecoder creates a Decoder instance without any decode operations. At least one must be added for
// decoding to work.
func NewCustomTypeDecoder[T any](options ...TypeDefaultOption) *TypeDecoder[T] {
	optns := defaultTypeDefaultOptions()
	optns.apply(options...)

	return &TypeDecoder[T]{
		dec:            instruct.NewTypeDecoder[*http.Request, DecodeContext, T](optns.options),
		defaultOptions: optns,
	}
}

// Decode decodes the http request to the struct passed in "data".
func (d *TypeDecoder[T]) Decode(r *http.Request, options ...TypeDecodeOption) (T, error) {
	optns := defaultDecodeOptions()
	optns.applyType(options...)

	optns.options.Ctx = &decodeContext{
		DefaultDecodeContext: instruct.NewDefaultDecodeContext(d.defaultOptions.options.FieldNameMapper),
		pathValue:            d.defaultOptions.pathValue,
		bodyDecoder:          d.defaultOptions.bodyDecoder,
		sliceSplitSeparator:  d.defaultOptions.sliceSplitSeparator,
		allowReadBody:        optns.allowReadBody,
		ensureAllQueryUsed:   optns.ensureAllQueryUsed,
		ensureAllFormUsed:    optns.ensureAllFormUsed,
	}

	return d.dec.Decode(r, optns.options)
}

// DecodeType decodes the http request to the struct passed in "data" using NewDecoder.
// Any map tags set using WithMapTags will be considered as "default" map tags. (see WithDefaultMapTags for details).
func DecodeType[T any](r *http.Request, options ...AnyTypeOption) (T, error) {
	options = inoptions.ConcatOptionsBefore[AnyTypeOption](options,
		WithDefaultDecodeOperations(),
	)
	return NewTypeDecoder[T](inoptions.ExtractOptions[TypeDefaultOption](options)...).Decode(r,
		inoptions.ExtractOptions[TypeDecodeOption](options)...)
}

// CustomDecodeType decodes the http request to the struct passed in "data" using NewCustomDecoder.
// Any map tags set using WithMapTags will be considered as "default" map tags. (see WithDefaultMapTags for details).
func CustomDecodeType[T any](r *http.Request, options ...AnyTypeOption) (T, error) {
	return NewTypeDecoder[T](inoptions.ExtractOptions[TypeDefaultOption](options)...).Decode(r,
		inoptions.ExtractOptions[TypeDecodeOption](options)...)
}
