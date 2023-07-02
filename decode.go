package inreq

import (
	"net/http"

	"github.com/RangelReale/instruct"
)

// Decoder decodes http requests to structs.
type Decoder struct {
	dec            *instruct.Decoder[*http.Request, DecodeContext]
	defaultOptions defaultOptions
}

// NewDecoder creates a Decoder instance with the default decode operations (query, path, header, form, body).
func NewDecoder(options ...DefaultOption) *Decoder {
	return NewCustomDecoder(concatOptionsBefore[DefaultOption](options, WithDefaultDecodeOperations())...)
}

// NewCustomDecoder creates a Decoder instance without any decode operations. At least one must be added for
// decoding to work.
func NewCustomDecoder(options ...DefaultOption) *Decoder {
	optns := defaultDefaultOptions()
	optns.apply(options...)

	return &Decoder{
		dec:            instruct.NewDecoder[*http.Request, DecodeContext](optns.options),
		defaultOptions: optns,
	}
}

// Decode decodes the http request to the struct passed in "data".
func (d *Decoder) Decode(r *http.Request, data any, options ...DecodeOption) error {
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

	return d.dec.Decode(r, data, optns.options)
}

// Decode decodes the http request to the struct passed in "data" using NewDecoder.
// Any map tags set using WithMapTags will be considered as "default" map tags. (see WithDefaultMapTags for details).
func Decode(r *http.Request, data any, options ...Option) error {
	options = concatOptionsBefore[Option](options,
		withUseDecodeMapTagsAsDefault(true),
		WithDefaultDecodeOperations(),
	)
	return NewDecoder(extractOptions[DefaultOption](options)...).Decode(r, data,
		extractOptions[DecodeOption](options)...)
}

// CustomDecode decodes the http request to the struct passed in "data" using NewCustomDecoder.
// Any map tags set using WithMapTags will be considered as "default" map tags. (see WithDefaultMapTags for details).
func CustomDecode(r *http.Request, data any, options ...Option) error {
	options = concatOptionsBefore[Option](options,
		withUseDecodeMapTagsAsDefault(true),
	)
	return NewDecoder(extractOptions[DefaultOption](options)...).Decode(r, data,
		extractOptions[DecodeOption](options)...)
}
