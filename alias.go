package inreq

import (
	"github.com/rrgmc/instruct"
	"github.com/rrgmc/instruct/types"
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
	CoerceError                = types.CoerceError
)

// map_tags.go

// MapTags is an alternative to struct tags, and can be used to override them.
type MapTags = instruct.MapTags

// option_field.go

// StructOption can be used as a struct field to give options to the struct itself.
type StructOption = instruct.StructOption

// StructOptionMapTag corresponds to StructOption in a MapTags.
const StructOptionMapTag = instruct.StructOptionMapTag

// resolver.go

// Resolver converts strings to the type of the struct field.
type Resolver = instruct.Resolver

// tag.go

// Tag contains the options parsed from the struct tags or MapTags.
type Tag = instruct.Tag
