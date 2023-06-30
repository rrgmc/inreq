package inreq

import (
	"net/http"

	"github.com/RangelReale/instruct"
)

type defaultOptions struct {
	options              instruct.DefaultOptions[*http.Request, DecodeContext]
	sliceSplitSeparator  string        // string to be used as separator on string-to-array conversion. Default is ",".
	pathValue            PathValue     // function used to extract the path from the request.
	defaultDecodeOptions decodeOptions // default decode options.
}

func (d *defaultOptions) apply(options ...DefaultOption) {
	for _, opt := range options {
		opt.applyDefaultOption(d)
	}
}

type decodeOptions struct {
	options            instruct.DecodeOptions[*http.Request, DecodeContext]
	allowReadBody      bool // whether operations are allowed to read the request body.
	ensureAllQueryUsed bool // whether to check if all query parameters were used.
	ensureAllFormUsed  bool // whether to check if all form parameters were used.
}

func (d *decodeOptions) apply(options ...DecodeOption) {
	for _, opt := range options {
		opt.applyDecodeOption(d)
	}
}

func defaultDefaultOptions() defaultOptions {
	ret := defaultOptions{
		options:              instruct.NewDefaultOptions[*http.Request, DecodeContext](),
		sliceSplitSeparator:  ",",
		defaultDecodeOptions: defaultDecodeOptions(),
	}
	ret.options.TagName = "inreq"
	return ret
}

func defaultDecodeOptions() decodeOptions {
	return decodeOptions{
		options:       instruct.NewDecodeOptions[*http.Request, DecodeContext](),
		allowReadBody: true,
	}
}

// helpers

type defaultOptionFunc func(*defaultOptions)

func (f defaultOptionFunc) isOption() {}

func (f defaultOptionFunc) applyDefaultOption(o *defaultOptions) {
	f(o)
}

type decodeOptionFunc func(*decodeOptions)

func (f decodeOptionFunc) isOption() {}

func (f decodeOptionFunc) applyDecodeOption(o *decodeOptions) {
	f(o)
}

// extractOptions extracts only options of a specific type.
func extractOptions[T Option](options []Option) []T {
	var ret []T
	for _, opt := range options {
		if o, ok := opt.(T); ok {
			ret = append(ret, o)
		}
	}
	return ret
}

// concatOptionsBefore returns an array with "options" before "source".
func concatOptionsBefore[T Option](source []T, options ...T) []T {
	return append(append([]T{}, options...), source...)
}
