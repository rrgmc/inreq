package inreq

import (
	"reflect"
)

// WithTagName sets the tag name to check on structs. The default is "inreq".
func WithTagName(tagName string) DefaultOption {
	return defaultOptionFunc(func(o *defaultOptions) {
		o.options.TagName = tagName
	})
}

// WithDefaultRequired sets whether the default for fields should be "required" or "not required"
func WithDefaultRequired(defaultRequired bool) DefaultOption {
	return defaultOptionFunc(func(o *defaultOptions) {
		o.options.DefaultRequired = defaultRequired
	})
}

// WithSliceSplitSeparator sets the string to be used as separator on string-to-array conversion. Default is ",".
func WithSliceSplitSeparator(sep string) DefaultOption {
	return defaultOptionFunc(func(o *defaultOptions) {
		o.sliceSplitSeparator = sep
	})
}

// WithFieldNameMapper sets the field name mapper. Default one uses [strings.ToLower].
func WithFieldNameMapper(fieldNameMapper FieldNameMapper) DefaultOption {
	return defaultOptionFunc(func(o *defaultOptions) {
		o.options.FieldNameMapper = fieldNameMapper
	})
}

// WithPathValue sets the function used to extract the path from the request.
func WithPathValue(pathValue PathValue) DefaultOption {
	return defaultOptionFunc(func(o *defaultOptions) {
		o.pathValue = pathValue
	})
}

// WithDefaultDecodeOperations adds the default operations (query, path, header, form and body).
// If the non-"Custom" calls are used, this option is added by default.
func WithDefaultDecodeOperations() DefaultOption {
	return defaultOptionFunc(func(o *defaultOptions) {
		o.options.DecodeOperations[OperationQuery] = &DecodeOperationQuery{}
		o.options.DecodeOperations[OperationPath] = &DecodeOperationPath{}
		o.options.DecodeOperations[OperationHeader] = &DecodeOperationHeader{}
		o.options.DecodeOperations[OperationForm] = &DecodeOperationForm{}
		o.options.DecodeOperations[OperationBody] = &DecodeOperationBody{}
	})
}

// WithDecodeOperation adds a decode operation.
func WithDecodeOperation(name string, operation DecodeOperation) DefaultOption {
	return defaultOptionFunc(func(o *defaultOptions) {
		if operation == nil {
			delete(o.options.DecodeOperations, name)
		} else {
			o.options.DecodeOperations[name] = operation
		}
	})
}

// WithResolver sets the decode Resolver.
func WithResolver(resolver Resolver) DefaultOption {
	return defaultOptionFunc(func(o *defaultOptions) {
		o.options.Resolver = resolver
	})
}

// WithDefaultMapTags adds a "default" MapTags. The default one is checked alongside the tags, so the
// check for unused fields takes both in account. Passing a struct without any struct tags and using
// WithMapTags will result in "field configuration not found" errors (except in free-standing functions like
// Decode, CustomDecode, DecodeType and CustomDecodeType.
func WithDefaultMapTags(dataForType any, tags MapTags) DefaultOption {
	return defaultOptionFunc(func(o *defaultOptions) {
		o.options.DefaultMapTagsSet(reflect.TypeOf(dataForType), tags)
	})
}

// WithDefaultMapTagsType is the same as WithDefaultMapTags using a reflect.Type.
func WithDefaultMapTagsType(typ reflect.Type, tags MapTags) DefaultOption {
	return defaultOptionFunc(func(o *defaultOptions) {
		o.options.DefaultMapTagsSet(typ, tags)
	})
}

// WithStructInfoCache sets whether to cache info for structs on parse. Default is false.
func WithStructInfoCache(cache bool) DefaultOption {
	return defaultOptionFunc(func(o *defaultOptions) {
		o.options.StructInfoCache(cache)
	})
}

// WithAllowReadBody sets whether operations are allowed to read the request body. Default is false.
func WithAllowReadBody(allowReadBody bool) FullOption {
	return &withAllowReadBody{allowReadBody: allowReadBody}
}

// WithEnsureAllQueryUsed sets whether to check if all query parameters were used.
func WithEnsureAllQueryUsed(ensureAllQueryUsed bool) FullOption {
	return &withEnsureAllQueryUsed{ensureAllQueryUsed: ensureAllQueryUsed}
}

// WithEnsureAllFormUsed sets whether to check if all form parameters were used.
func WithEnsureAllFormUsed(ensureAllFormUsed bool) FullOption {
	return &withEnsureAllFormUsed{ensureAllFormUsed: ensureAllFormUsed}
}

// WithMapTags sets decode-operation-specific MapTags. These override the default cached struct information
// but don't change the original one. This should be used to override configurations on each call.
func WithMapTags(tags MapTags) DecodeOption {
	return decodeOptionFunc(func(o *decodeOptions) {
		o.options.MapTags = tags
	})
}

// withUseDecodeMapTagsAsDefault is an internal option to allow WithMapTags to set default map tags for
// free-standing Decode functions.
func withUseDecodeMapTagsAsDefault(useDecodeMapTagsAsDefault bool) DecodeOption {
	return decodeOptionFunc(func(o *decodeOptions) {
		o.options.UseDecodeMapTagsAsDefault = useDecodeMapTagsAsDefault
	})
}

type withAllowReadBody struct {
	allowReadBody bool
}

func (w withAllowReadBody) isOption() {}

func (w withAllowReadBody) applyDefaultOption(o *defaultOptions) {
	o.defaultDecodeOptions.allowReadBody = w.allowReadBody
}

func (w withAllowReadBody) applyDecodeOption(o *decodeOptions) {
	o.allowReadBody = w.allowReadBody
}

type withEnsureAllQueryUsed struct {
	ensureAllQueryUsed bool
}

func (w withEnsureAllQueryUsed) isOption() {}

func (w withEnsureAllQueryUsed) applyDefaultOption(o *defaultOptions) {
	o.defaultDecodeOptions.ensureAllQueryUsed = w.ensureAllQueryUsed
}

func (w withEnsureAllQueryUsed) applyDecodeOption(o *decodeOptions) {
	o.ensureAllQueryUsed = w.ensureAllQueryUsed
}

type withEnsureAllFormUsed struct {
	ensureAllFormUsed bool
}

func (w withEnsureAllFormUsed) isOption() {}

func (w withEnsureAllFormUsed) applyDefaultOption(o *defaultOptions) {
	o.defaultDecodeOptions.ensureAllFormUsed = w.ensureAllFormUsed
}

func (w withEnsureAllFormUsed) applyDecodeOption(o *decodeOptions) {
	o.ensureAllFormUsed = w.ensureAllFormUsed
}
