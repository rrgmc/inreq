package inreq

import (
	"github.com/RangelReale/instruct"
	"github.com/RangelReale/instruct/types"
)

// error.go

var (
	ErrCoerceInvalid     = types.ErrCoerceInvalid
	ErrCoerceOverflow    = types.ErrCoerceOverflow
	ErrCoerceUnsupported = types.ErrCoerceUnsupported
	ErrCoerceUnknown     = types.ErrCoerceUnknown
)

type (
	ValuesNotUsedError         = types.ValuesNotUsedError
	InvalidDecodeError         = types.InvalidDecodeError
	RequiredError              = types.RequiredError
	OperationNotSupportedError = types.OperationNotSupportedError
)

// map_tags.go

// MapTags is an alternative to struct tags, and can be used to override them.
type MapTags = instruct.MapTags

// resolver.go

// Resolver converts strings to the type of the struct field.
type Resolver = instruct.Resolver

// tag.go

// Tag contains the options parsed from the struct tags or MapTags.
type Tag = instruct.Tag
