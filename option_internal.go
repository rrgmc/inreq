package inreq

import (
	"net/http"

	"github.com/RangelReale/instruct"
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
		opt.applyDefaultOption(d)
	}
}

type typeDefaultOptions struct {
	options instruct.TypeDefaultOptions[*http.Request, DecodeContext]
	sharedDefaultOptions
}

func (d *typeDefaultOptions) apply(options ...TypeDefaultOption) {
	for _, opt := range options {
		opt.applyTypeDefaultOption(d)
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

// defaultOptions

type defaultOptionFunc func(*defaultOptions)

func (f defaultOptionFunc) isOption() {}

func (f defaultOptionFunc) applyDefaultOption(o *defaultOptions) {
	f(o)
}

// defaultAndTypeOptions

type defaultAndTypeOptionImpl struct {
	f func(o *instruct.DefaultOptions[*http.Request, DecodeContext])
}

func (f defaultAndTypeOptionImpl) isOption() {}

func (f defaultAndTypeOptionImpl) applyDefaultOption(o *defaultOptions) {
	f.f(&o.options)
}

func (f defaultAndTypeOptionImpl) applyTypeDefaultOption(o *typeDefaultOptions) {
	f.f(&o.options.DefaultOptions)
}

func defaultAndTypeOptionsFunc(f func(o *instruct.DefaultOptions[*http.Request, DecodeContext])) *defaultAndTypeOptionImpl {
	return &defaultAndTypeOptionImpl{f}
}

// defaultAndTypeSharedOptions

type defaultAndTypeSharedOptionImpl struct {
	f func(o *sharedDefaultOptions)
}

func (f defaultAndTypeSharedOptionImpl) isOption() {}

func (f defaultAndTypeSharedOptionImpl) applyDefaultOption(o *defaultOptions) {
	f.f(&o.sharedDefaultOptions)
}

func (f defaultAndTypeSharedOptionImpl) applyTypeDefaultOption(o *typeDefaultOptions) {
	f.f(&o.sharedDefaultOptions)
}

func defaultAndTypeSharedOptionFunc(f func(o *sharedDefaultOptions)) *defaultAndTypeSharedOptionImpl {
	return &defaultAndTypeSharedOptionImpl{f}
}

// typeDefaultOptions

type typeDefaultOptionFunc func(*typeDefaultOptions)

func (f typeDefaultOptionFunc) isOption() {}

func (f typeDefaultOptionFunc) applyTypeDefaultOption(o *typeDefaultOptions) {
	f(o)
}

// decodeOptions

type decodeOptionFunc func(*decodeOptions)

func (f decodeOptionFunc) isOption() {}

func (f decodeOptionFunc) applyDecodeOption(o *decodeOptions) {
	f(o)
}

// typeAndDecodeOptions

type typeAndDecodeOptionImpl struct {
	t func(o *typeDefaultOptions)
	d func(o *decodeOptions)
}

func (f typeAndDecodeOptionImpl) isOption() {}

func (f typeAndDecodeOptionImpl) applyTypeDefaultOption(o *typeDefaultOptions) {
	f.t(o)
}

func (f typeAndDecodeOptionImpl) applyDecodeOption(o *decodeOptions) {
	f.d(o)
}

func typeAndDecodeOptionFunc(t func(o *typeDefaultOptions), d func(o *decodeOptions)) *typeAndDecodeOptionImpl {
	return &typeAndDecodeOptionImpl{t, d}
}

// defaultAndTypeSharedFullOptions

type defaultAndTypeSharedFullOptionImpl struct {
	def func(o *sharedDefaultOptions)
	dec func(o *decodeOptions)
}

func (f defaultAndTypeSharedFullOptionImpl) isOption() {}

func (f defaultAndTypeSharedFullOptionImpl) applyDefaultOption(o *defaultOptions) {
	f.def(&o.sharedDefaultOptions)
}

func (f defaultAndTypeSharedFullOptionImpl) applyTypeDefaultOption(o *typeDefaultOptions) {
	f.def(&o.sharedDefaultOptions)
}

func (f defaultAndTypeSharedFullOptionImpl) applyDecodeOption(o *decodeOptions) {
	f.dec(o)
}

func defaultAndTypeSharedFullOptionFunc(def func(o *sharedDefaultOptions), dec func(o *decodeOptions)) *defaultAndTypeSharedFullOptionImpl {
	return &defaultAndTypeSharedFullOptionImpl{def, dec}
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
