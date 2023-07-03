package inreq

import (
	"net/http"

	"github.com/RangelReale/instruct"
)

// TypeDecoder decodes http requests to structs.
type TypeDecoder[T any] struct {
	dec            *instruct.TypeDecoder[*http.Request, DecodeContext, T]
	defaultOptions defaultOptions
}

// NewTypeDecoder creates a Decoder instance with the default decode operations (query, path, header, form, body).
func NewTypeDecoder[T any](options ...DefaultOption) *TypeDecoder[T] {
	return NewCustomTypeDecoder[T](concatOptionsBefore[DefaultOption](options, WithDefaultDecodeOperations())...)
}

// NewCustomTypeDecoder creates a Decoder instance without any decode operations. At least one must be added for
// decoding to work.
func NewCustomTypeDecoder[T any](options ...DefaultOption) *TypeDecoder[T] {
	optns := defaultDefaultOptions()
	optns.apply(options...)

	return &TypeDecoder[T]{
		dec:            instruct.NewTypeDecoder[*http.Request, DecodeContext, T](optns.options),
		defaultOptions: optns,
	}
}

// Decode decodes the http request to the struct passed in "data".
func (d *TypeDecoder[T]) Decode(r *http.Request, options ...DecodeOption) (T, error) {
	optns := defaultDecodeOptions()
	optns.apply(options...)

	optns.options.Ctx = &decodeContext{
		DefaultDecodeContext: instruct.NewDefaultDecodeContext(d.defaultOptions.options.FieldNameMapper),
		pathValue:            d.defaultOptions.pathValue,
		sliceSplitSeparator:  d.defaultOptions.sliceSplitSeparator,
		allowReadBody:        optns.allowReadBody,
		ensureAllQueryUsed:   optns.ensureAllQueryUsed,
		ensureAllFormUsed:    optns.ensureAllFormUsed,
	}

	return d.dec.Decode(r, optns.options)
}

// DecodeType decodes the http request to the struct passed in "data" using NewDecoder.
// Any map tags set using WithMapTags will be considered as "default" map tags. (see WithDefaultMapTags for details).
func DecodeType[T any](r *http.Request, options ...Option) (T, error) {
	options = concatOptionsBefore[Option](options,
		withUseDecodeMapTagsAsDefault(true),
		WithDefaultDecodeOperations(),
	)
	return NewTypeDecoder[T](extractOptions[DefaultOption](options)...).Decode(r,
		extractOptions[DecodeOption](options)...)
}

// CustomDecodeType decodes the http request to the struct passed in "data" using NewCustomDecoder.
// Any map tags set using WithMapTags will be considered as "default" map tags. (see WithDefaultMapTags for details).
func CustomDecodeType[T any](r *http.Request, options ...Option) (T, error) {
	options = concatOptionsBefore[Option](options,
		withUseDecodeMapTagsAsDefault(true),
	)
	return NewTypeDecoder[T](extractOptions[DefaultOption](options)...).Decode(r,
		extractOptions[DecodeOption](options)...)
}
