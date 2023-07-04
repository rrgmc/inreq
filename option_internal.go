package inreq

import (
	"net/http"

	"github.com/RangelReale/instruct"
	"github.com/RangelReale/instruct/options"
)

const (
	DefaultTagName = "inreq"
)

type sharedDefaultOptions struct {
	sliceSplitSeparator  string        // string to be used as separator on string-to-array conversion. Default is ",".
	pathValue            PathValue     // function used to extract the path from the request.
	defaultDecodeOptions decodeOptions // default decode options.
}

type defaultOptions struct {
	options instruct.DefaultOptions[*http.Request, DecodeContext]
	sharedDefaultOptions
}

func (d *defaultOptions) apply(options ...DefaultOption) {
	for _, opt := range options {
		opt.ApplyDefaultOption(d)
	}
}

type typeDefaultOptions struct {
	options instruct.TypeDefaultOptions[*http.Request, DecodeContext]
	sharedDefaultOptions
}

func (d *typeDefaultOptions) apply(options ...TypeDefaultOption) {
	for _, opt := range options {
		opt.ApplyTypeDefaultOption(d)
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
		opt.ApplyDecodeOption(d)
	}
}

func defaultSharedDefaultOptions() sharedDefaultOptions {
	ret := sharedDefaultOptions{
		sliceSplitSeparator:  ",",
		defaultDecodeOptions: defaultDecodeOptions(),
	}
	return ret
}

func defaultDefaultOptions() defaultOptions {
	ret := defaultOptions{
		options:              instruct.NewDefaultOptions[*http.Request, DecodeContext](),
		sharedDefaultOptions: defaultSharedDefaultOptions(),
	}
	ret.options.TagName = DefaultTagName
	return ret
}

func defaultTypeDefaultOptions() typeDefaultOptions {
	ret := typeDefaultOptions{
		options:              instruct.NewTypeDefaultOptions[*http.Request, DecodeContext](),
		sharedDefaultOptions: defaultSharedDefaultOptions(),
	}
	ret.options.TagName = DefaultTagName
	return ret
}

func defaultDecodeOptions() decodeOptions {
	return decodeOptions{
		options:       instruct.NewDecodeOptions[*http.Request, DecodeContext](),
		allowReadBody: true,
	}
}

// helpers

// DefaultOption -> instruct.DefaultOptions

func defaultOptionFunc(f func(o *instruct.DefaultOptions[*http.Request, DecodeContext])) DefaultOption {
	return options.DefaultOptionFunc[*http.Request, DecodeContext, defaultOptions, typeDefaultOptions](func(o *defaultOptions) {
		f(&o.options)
	}, func(o *typeDefaultOptions) {
		f(&o.options.DefaultOptions)
	})
}

func decodeOptionFunc(f func(o *decodeOptions)) DecodeOption {
	return options.DecodeOptionFunc[*http.Request, DecodeContext, decodeOptions](func(o *decodeOptions) {
		f(o)
	})
}

func defaultSharedOptionFunc(f func(o *sharedDefaultOptions)) DefaultOption {
	return options.DefaultOptionFunc[*http.Request, DecodeContext, defaultOptions, typeDefaultOptions](func(o *defaultOptions) {
		f(&o.sharedDefaultOptions)
	}, func(o *typeDefaultOptions) {
		f(&o.sharedDefaultOptions)
	})
}

func typeAndDecodeOptionFunc(tf func(o *typeDefaultOptions), cf func(o *decodeOptions)) TypeDefaultAndDecodeOption {
	return options.TypeDefaultAndDecodeOptionFunc[*http.Request, DecodeContext, typeDefaultOptions, decodeOptions](func(o *typeDefaultOptions) {
		tf(o)
	}, func(o *decodeOptions) {
		cf(o)
	})
}

func fullSharedOptionFunc(def func(o *sharedDefaultOptions), dec func(o *decodeOptions)) FullOption {
	return options.FullOptionFunc[*http.Request, DecodeContext, defaultOptions, typeDefaultOptions, decodeOptions](func(o *defaultOptions) {
		def(&o.sharedDefaultOptions)
	}, func(o *typeDefaultOptions) {
		def(&o.sharedDefaultOptions)
	}, func(o *decodeOptions) {
		dec(o)
	})
}

// concatOptionsBefore returns an array with "options" before "source".
func concatOptionsBefore[T Option](source []T, optns ...T) []T {
	return options.ConcatOptionsBefore[*http.Request, DecodeContext](source, optns...)
}

// extractOptions extracts only options of a specific type.
func extractOptions[T Option](optns []Option) []T {
	return options.ExtractOptions[*http.Request, DecodeContext, T](optns)
}
