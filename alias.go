package inreq

import (
	"github.com/RangelReale/instruct"
	"github.com/RangelReale/instruct/types"
)

// map_tags.go

// MapTags is an alternative to struct tags, and can be used to override them.
type MapTags = instruct.MapTags

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
