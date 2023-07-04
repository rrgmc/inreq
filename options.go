package inreq

import (
	"net/http"
	"reflect"

	"github.com/RangelReale/instruct"
)

// WithTagName sets the tag name to check on structs. The default is "inreq".
func WithTagName(tagName string) DefaultOption {
	return defaultAndTypeDefaultOptionFunc(func(o *instruct.DefaultOptions[*http.Request, DecodeContext]) {
		o.TagName = tagName
	})
}

// WithDefaultRequired sets whether the default for fields should be "required" or "not required"
func WithDefaultRequired(defaultRequired bool) DefaultOption {
	return defaultAndTypeDefaultOptionFunc(func(o *instruct.DefaultOptions[*http.Request, DecodeContext]) {
		o.DefaultRequired = defaultRequired
	})
}

// WithSliceSplitSeparator sets the string to be used as separator on string-to-array conversion. Default is ",".
func WithSliceSplitSeparator(sep string) DefaultOption {
	return defaultAndTypeDefaultSharedOptionFunc(func(o *sharedDefaultOptions) {
		o.sliceSplitSeparator = sep
	})
}

// WithFieldNameMapper sets the field name mapper. Default one uses [strings.ToLower].
func WithFieldNameMapper(fieldNameMapper FieldNameMapper) DefaultOption {
	return defaultAndTypeDefaultOptionFunc(func(o *instruct.DefaultOptions[*http.Request, DecodeContext]) {
		o.FieldNameMapper = fieldNameMapper
	})
}

// WithPathValue sets the function used to extract the path from the request.
func WithPathValue(pathValue PathValue) DefaultOption {
	return defaultAndTypeDefaultSharedOptionFunc(func(o *sharedDefaultOptions) {
		o.pathValue = pathValue
	})
}

// WithDefaultDecodeOperations adds the default operations (query, path, header, form and body).
// If the non-"Custom" calls are used, this option is added by default.
func WithDefaultDecodeOperations() DefaultAndTypeDefaultOption {
	return defaultAndTypeDefaultOptionFunc(func(o *instruct.DefaultOptions[*http.Request, DecodeContext]) {
		o.DecodeOperations[OperationQuery] = &DecodeOperationQuery{}
		o.DecodeOperations[OperationPath] = &DecodeOperationPath{}
		o.DecodeOperations[OperationHeader] = &DecodeOperationHeader{}
		o.DecodeOperations[OperationForm] = &DecodeOperationForm{}
		o.DecodeOperations[OperationBody] = &DecodeOperationBody{}
	})
}

// WithDecodeOperation adds a decode operation.
func WithDecodeOperation(name string, operation DecodeOperation) DefaultOption {
	return defaultAndTypeDefaultOptionFunc(func(o *instruct.DefaultOptions[*http.Request, DecodeContext]) {
		if operation == nil {
			delete(o.DecodeOperations, name)
		} else {
			o.DecodeOperations[name] = operation
		}
	})
}

// WithResolver sets the decode Resolver.
func WithResolver(resolver Resolver) DefaultOption {
	return defaultAndTypeDefaultOptionFunc(func(o *instruct.DefaultOptions[*http.Request, DecodeContext]) {
		o.Resolver = resolver
	})
}

// WithDefaultMapTags adds a "default" MapTags. The default one is checked alongside the tags, so the
// check for unused fields takes both in account. Passing a struct without any struct tags and using
// WithMapTags will result in "field configuration not found" errors (except in free-standing functions like
// Decode, CustomDecode, DecodeType and CustomDecodeType.
func WithDefaultMapTags(dataForType any, tags MapTags) DefaultOption {
	return defaultAndTypeDefaultOptionFunc(func(o *instruct.DefaultOptions[*http.Request, DecodeContext]) {
		o.DefaultMapTagsSet(reflect.TypeOf(dataForType), tags)
	})
}

// WithDefaultMapTagsType is the same as WithDefaultMapTags using a reflect.Type.
func WithDefaultMapTagsType(typ reflect.Type, tags MapTags) DefaultOption {
	return defaultAndTypeDefaultOptionFunc(func(o *instruct.DefaultOptions[*http.Request, DecodeContext]) {
		o.DefaultMapTagsSet(typ, tags)
	})
}

// WithStructInfoCache sets whether to cache info for structs on parse. Default is false.
func WithStructInfoCache(cache bool) DefaultOption {
	return defaultAndTypeDefaultOptionFunc(func(o *instruct.DefaultOptions[*http.Request, DecodeContext]) {
		o.StructInfoCache(cache)
	})
}

// WithAllowReadBody sets whether operations are allowed to read the request body. Default is false.
func WithAllowReadBody(allowReadBody bool) FullOption {
	return fullSharedOptionFunc(func(o *sharedDefaultOptions) {
		o.defaultDecodeOptions.allowReadBody = allowReadBody
	}, func(o *decodeOptions) {
		o.allowReadBody = allowReadBody
	})
}

// WithEnsureAllQueryUsed sets whether to check if all query parameters were used.
func WithEnsureAllQueryUsed(ensureAllQueryUsed bool) FullOption {
	return fullSharedOptionFunc(func(o *sharedDefaultOptions) {
		o.defaultDecodeOptions.ensureAllQueryUsed = ensureAllQueryUsed
	}, func(o *decodeOptions) {
		o.ensureAllQueryUsed = ensureAllQueryUsed
	})
}

// WithEnsureAllFormUsed sets whether to check if all form parameters were used.
func WithEnsureAllFormUsed(ensureAllFormUsed bool) FullOption {
	return fullSharedOptionFunc(func(o *sharedDefaultOptions) {
		o.defaultDecodeOptions.ensureAllFormUsed = ensureAllFormUsed
	}, func(o *decodeOptions) {
		o.ensureAllFormUsed = ensureAllFormUsed
	})
}

// WithMapTags sets decode-operation-specific MapTags. These override the default cached struct information
// but don't change the original one. This should be used to override configurations on each call.
func WithMapTags(tags MapTags) TypeDefaultAndDecodeOption {
	return typeAndDecodeOptionFunc(func(o *typeDefaultOptions) {
		o.options.MapTags = tags
	}, func(o *decodeOptions) {
		o.options.MapTags = tags
	})
}

func WithX(x int) TypeDefaultOption {
	return typeDefaultOptionFunc(func(o *typeDefaultOptions) {
		o.x = x
	})
}

// withUseDecodeMapTagsAsDefault is an internal option to allow WithMapTags to set default map tags for
// free-standing Decode functions.
func withUseDecodeMapTagsAsDefault(useDecodeMapTagsAsDefault bool) DecodeOption {
	return decodeOptionFunc(func(o *decodeOptions) {
		o.options.UseDecodeMapTagsAsDefault = useDecodeMapTagsAsDefault
	})
}
